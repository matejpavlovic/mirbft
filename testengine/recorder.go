/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package testengine

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"hash"

	"github.com/IBM/mirbft"
	pb "github.com/IBM/mirbft/mirbftpb"
	tpb "github.com/IBM/mirbft/testengine/testenginepb"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Hasher func() hash.Hash

func uint64ToBytes(value uint64) []byte {
	byteValue := make([]byte, 8)
	binary.LittleEndian.PutUint64(byteValue, value)
	return byteValue
}

func bytesToUint64(value []byte) uint64 {
	return binary.LittleEndian.Uint64(value)
}

type RecorderNode struct {
	PlaybackNode         *PlaybackNode
	State                *NodeState
	Config               *tpb.NodeConfig
	AwaitingProcessEvent bool
}

type RecorderClient struct {
	Config            *ClientConfig
	LastNodeReqNoSend []uint64
}

func (rc *RecorderClient) RequestByReqNo(reqNo uint64) *pb.RequestData {
	if reqNo > rc.Config.Total {
		// We've sent all we should
		return nil
	}

	var buffer bytes.Buffer
	buffer.Write(rc.Config.ID)
	buffer.Write([]byte("-"))
	buffer.Write(uint64ToBytes(reqNo))

	return &pb.RequestData{
		ClientId: rc.Config.ID,
		ReqNo:    reqNo,
		Data:     buffer.Bytes(),
	}
}

type NodeState struct {
	LastCommittedSeqNo uint64
	OutstandingCommits []*mirbft.Commit
	Hasher             hash.Hash
	Value              []byte
	Length             uint64
}

func (ns *NodeState) Commit(commits []*mirbft.Commit, node uint64) []*tpb.Checkpoint {
	for _, commit := range commits {
		if commit.QEntry.SeqNo <= ns.LastCommittedSeqNo {
			panic(fmt.Sprintf("trying to commit seqno=%d, but we've already committed seqno=%d", commit.QEntry.SeqNo, ns.LastCommittedSeqNo))
		}
		index := commit.QEntry.SeqNo - ns.LastCommittedSeqNo
		for index >= uint64(len(ns.OutstandingCommits)) {
			ns.OutstandingCommits = append(ns.OutstandingCommits, nil)
		}
		ns.OutstandingCommits[index-1] = commit
	}

	var results []*tpb.Checkpoint

	i := 0
	for _, commit := range ns.OutstandingCommits {
		if commit == nil {
			break
		}
		i++

		for _, request := range commit.QEntry.Requests {
			ns.Hasher.Write(request.Digest)
			ns.Length++
		}

		if commit.Checkpoint {
			results = append(results, &tpb.Checkpoint{
				SeqNo: commit.QEntry.SeqNo,
				Value: ns.Hasher.Sum(nil),
			})
		}

		ns.LastCommittedSeqNo = commit.QEntry.SeqNo
	}

	k := 0
	for j := i; j < len(ns.OutstandingCommits); j++ {
		ns.OutstandingCommits[k] = ns.OutstandingCommits[j]
		k++
	}
	ns.OutstandingCommits = ns.OutstandingCommits[:k]

	ns.Value = ns.Hasher.Sum(nil)

	return results

}

type ClientConfig struct {
	ID          []byte
	TxLatency   uint64
	MaxInFlight int
	Total       uint64
}

type Recorder struct {
	NetworkConfig *pb.NetworkConfig
	NodeConfigs   []*tpb.NodeConfig
	ClientConfigs []*ClientConfig
	Logger        *zap.Logger
	Hasher        Hasher
}

func (r *Recorder) Recording() (*Recording, error) {
	eventLog := &EventLog{
		InitialConfig: r.NetworkConfig,
		NodeConfigs:   r.NodeConfigs,
	}

	player, err := NewPlayer(eventLog, r.Logger)
	if err != nil {
		return nil, errors.WithMessage(err, "could not construct player")
	}

	nodes := make([]*RecorderNode, len(r.NodeConfigs))
	for i, nodeConfig := range r.NodeConfigs {
		nodeID := uint64(i)
		eventLog.InsertTick(nodeID, uint64(nodeConfig.TickInterval))
		nodes[i] = &RecorderNode{
			State: &NodeState{
				Hasher: r.Hasher(),
			},
			PlaybackNode: player.Nodes[i],
			Config:       nodeConfig,
		}
	}

	clients := make([]*RecorderClient, len(r.ClientConfigs))
	for i, clientConfig := range r.ClientConfigs {
		client := &RecorderClient{
			Config:            clientConfig,
			LastNodeReqNoSend: make([]uint64, len(nodes)),
		}

		clients[i] = client

		for i := 1; i <= clientConfig.MaxInFlight; i++ {
			req := client.RequestByReqNo(uint64(i))
			if req == nil {
				continue
			}
			for j := range nodes {
				client.LastNodeReqNoSend[uint64(j)] = uint64(i)
				eventLog.InsertPropose(uint64(j), req, client.Config.TxLatency)
			}
		}
	}

	return &Recording{
		Hasher:   r.Hasher,
		EventLog: player.EventLog,
		Player:   player,
		Nodes:    nodes,
		Clients:  clients,
	}, nil
}

type Recording struct {
	Hasher   Hasher
	EventLog *EventLog
	Player   *Player
	Nodes    []*RecorderNode
	Clients  []*RecorderClient
}

func (r *Recording) Step() error {
	err := r.Player.Step()
	if err != nil {
		return errors.WithMessagef(err, "could not step recorder's underlying player")
	}

	lastEvent := r.Player.LastEvent
	node := r.Nodes[int(lastEvent.Target)]
	nodeConfig := node.Config
	playbackNode := node.PlaybackNode
	nodeState := node.State

	switch lastEvent.Type.(type) {
	case *tpb.Event_Apply_:
		nodeStatus := node.PlaybackNode.Status
		for _, rw := range nodeStatus.RequestWindows {
			for _, client := range r.Clients {
				if !bytes.Equal(client.Config.ID, rw.ClientID) {
					continue
				}

				for i := client.LastNodeReqNoSend[lastEvent.Target] + 1; i <= rw.HighWatermark; i++ {
					req := client.RequestByReqNo(i)
					if req == nil {
						continue
					}
					client.LastNodeReqNoSend[lastEvent.Target] = i
					r.EventLog.InsertPropose(lastEvent.Target, req, client.Config.TxLatency)
				}
			}
		}
	case *tpb.Event_Receive_:
	case *tpb.Event_Process_:
		if !node.AwaitingProcessEvent {
			return errors.Errorf("node %d was not awaiting a processing message, but got one", lastEvent.Target)
		}
		node.AwaitingProcessEvent = false
		processing := playbackNode.Processing
		for _, msg := range processing.Broadcast {
			for i := range r.Player.Nodes {
				if uint64(i) != lastEvent.Target {
					r.EventLog.InsertRecv(uint64(i), lastEvent.Target, msg, uint64(nodeConfig.LinkLatency))
				} else {
					// Send it to ourselves with no latency
					err := playbackNode.Node.Step(context.Background(), lastEvent.Target, msg)
					if err != nil {
						return errors.WithMessagef(err, "node %d could not step message to self", lastEvent.Target)
					}
				}
			}
		}

		for _, unicast := range processing.Unicast {
			if unicast.Target != lastEvent.Target {
				r.EventLog.InsertRecv(unicast.Target, lastEvent.Target, unicast.Msg, uint64(nodeConfig.LinkLatency))
			} else {
				// Send it to ourselves with no latency
				err := playbackNode.Node.Step(context.Background(), lastEvent.Target, unicast.Msg)
				if err != nil {
					return errors.WithMessagef(err, "node %d could not step message to self", lastEvent.Target)
				}
			}
		}

		apply := &tpb.Event_Apply{
			Preprocessed: make([]*tpb.Request, len(processing.Preprocess)),
			Processed:    make([]*tpb.Batch, len(processing.Process)),
		}

		for i, preprocess := range processing.Preprocess {
			hasher := r.Hasher()
			hasher.Write(preprocess.ClientRequest.ClientId)
			hasher.Write(uint64ToBytes(preprocess.ClientRequest.ReqNo))
			hasher.Write(preprocess.ClientRequest.Data)
			hasher.Write(preprocess.ClientRequest.Signature)

			apply.Preprocessed[i] = &tpb.Request{
				ClientId:  preprocess.ClientRequest.ClientId,
				ReqNo:     preprocess.ClientRequest.ReqNo,
				Data:      preprocess.ClientRequest.Data,
				Signature: preprocess.ClientRequest.Signature,
				Digest:    hasher.Sum(nil),
			}
		}

		for i, process := range processing.Process {
			hasher := r.Hasher()
			requests := make([]*pb.Request, len(process.Requests))
			for i, request := range process.Requests {
				hasher.Write(request.Digest)
				requests[i] = &pb.Request{
					ClientId: request.RequestData.ClientId,
					ReqNo:    request.RequestData.ReqNo,
					Digest:   request.Digest,
				}
			}

			apply.Processed[i] = &tpb.Batch{
				Source:   process.Source,
				Epoch:    process.Epoch,
				SeqNo:    process.SeqNo,
				Digest:   hasher.Sum(nil),
				Requests: requests,
			}
		}

		apply.Checkpoints = nodeState.Commit(processing.Commits, lastEvent.Target)

		r.EventLog.InsertApply(lastEvent.Target, apply, uint64(nodeConfig.ReadyLatency))
	case *tpb.Event_Propose_:
	case *tpb.Event_Tick_:
		r.EventLog.InsertTick(lastEvent.Target, uint64(nodeConfig.TickInterval))
	}

	if playbackNode.Processing == nil &&
		!playbackNode.Actions.IsEmpty() &&
		!node.AwaitingProcessEvent {
		r.EventLog.InsertProcess(lastEvent.Target, uint64(nodeConfig.ProcessLatency))
		node.AwaitingProcessEvent = true
	}

	return nil
}

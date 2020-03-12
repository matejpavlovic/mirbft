/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mirbft

import (
	"bytes"
	"fmt"

	pb "github.com/IBM/mirbft/mirbftpb"
)

type request struct {
	requestData *pb.RequestData
	digest      []byte
	state       SequenceState
	seqNo       uint64
}

type clientWindow struct {
	lowWatermark  uint64
	highWatermark uint64
	requests      []*request
	clientWaiter  *clientWaiter // Used to throttle clients
}

type clientWaiter struct {
	lowWatermark  uint64
	highWatermark uint64
	expired       chan struct{}
}

func newRequestWindow(lowWatermark, highWatermark uint64) *clientWindow {
	return &clientWindow{
		lowWatermark:  lowWatermark,
		highWatermark: highWatermark,
		requests:      make([]*request, int(highWatermark-lowWatermark)+1),
		clientWaiter: &clientWaiter{
			lowWatermark:  lowWatermark,
			highWatermark: highWatermark,
			expired:       make(chan struct{}),
		},
	}
}

func (rw *clientWindow) garbageCollect(maxSeqNo uint64) {
	newRequests := make([]*request, int(rw.highWatermark-rw.lowWatermark)+1)
	i := 0
	j := uint64(0)
	copying := false
	for _, request := range rw.requests {
		if request == nil || request.state != Committed || request.seqNo > maxSeqNo {
			copying = true
		}

		if copying {
			newRequests[i] = request
			i++
		} else {
			if request.seqNo == 0 {
				panic("this should be initialized if here")
			}
			j++
		}

	}

	rw.lowWatermark += j
	rw.highWatermark += j
	rw.requests = newRequests
	close(rw.clientWaiter.expired)
	rw.clientWaiter = &clientWaiter{
		lowWatermark:  rw.lowWatermark,
		highWatermark: rw.highWatermark,
		expired:       make(chan struct{}),
	}
}

func (rw *clientWindow) allocate(requestData *pb.RequestData, digest []byte) {
	reqNo := requestData.ReqNo
	if reqNo > rw.highWatermark {
		panic(fmt.Sprintf("unexpected: %d > %d", reqNo, rw.highWatermark))
	}

	if reqNo < rw.lowWatermark {
		panic(fmt.Sprintf("unexpected: %d < %d", reqNo, rw.lowWatermark))
	}

	offset := int(reqNo - rw.lowWatermark)
	if rw.requests[offset] != nil && !bytes.Equal(rw.requests[offset].digest, digest) {
		panic("we don't handle byzantine clients yet, but two different requests for the same reqno")
	}

	rw.requests[offset] = &request{
		requestData: requestData,
		digest:      digest,
	}
}

func (rw *clientWindow) request(reqNo uint64) *request {
	if reqNo > rw.highWatermark {
		panic(fmt.Sprintf("unexpected: %d > %d", reqNo, rw.highWatermark))
	}

	if reqNo < rw.lowWatermark {
		panic(fmt.Sprintf("unexpected: %d < %d", reqNo, rw.lowWatermark))
	}

	offset := int(reqNo - rw.lowWatermark)

	return rw.requests[offset]
}

func (rw *clientWindow) status() *RequestWindowStatus {
	allocated := make([]uint64, len(rw.requests))
	for i, request := range rw.requests {
		if request == nil {
			continue
		}
		if request.state == Committed {
			allocated[i] = 2
		} else {
			allocated[i] = 1
		}
		// allocated[i] = bytesToUint64(request.preprocessResult.Proposal.Data)
	}

	return &RequestWindowStatus{
		LowWatermark:  rw.lowWatermark,
		HighWatermark: rw.highWatermark,
		Allocated:     allocated,
	}
}

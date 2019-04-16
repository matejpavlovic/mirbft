// Code generated by protoc-gen-go.
// source: mirbft.proto
// DO NOT EDIT!

/*
Package mirbftpb is a generated protocol buffer package.

It is generated from these files:
	mirbft.proto

It has these top-level messages:
	Msg
	Preprepare
	Prepare
	Commit
*/
package mirbftpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Msg struct {
	// Types that are valid to be assigned to Type:
	//	*Msg_Preprepare
	//	*Msg_Prepare
	//	*Msg_Commit
	Type isMsg_Type `protobuf_oneof:"type"`
}

func (m *Msg) Reset()                    { *m = Msg{} }
func (m *Msg) String() string            { return proto.CompactTextString(m) }
func (*Msg) ProtoMessage()               {}
func (*Msg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type isMsg_Type interface {
	isMsg_Type()
}

type Msg_Preprepare struct {
	Preprepare *Preprepare `protobuf:"bytes,1,opt,name=preprepare,oneof"`
}
type Msg_Prepare struct {
	Prepare *Prepare `protobuf:"bytes,2,opt,name=prepare,oneof"`
}
type Msg_Commit struct {
	Commit *Commit `protobuf:"bytes,3,opt,name=commit,oneof"`
}

func (*Msg_Preprepare) isMsg_Type() {}
func (*Msg_Prepare) isMsg_Type()    {}
func (*Msg_Commit) isMsg_Type()     {}

func (m *Msg) GetType() isMsg_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *Msg) GetPreprepare() *Preprepare {
	if x, ok := m.GetType().(*Msg_Preprepare); ok {
		return x.Preprepare
	}
	return nil
}

func (m *Msg) GetPrepare() *Prepare {
	if x, ok := m.GetType().(*Msg_Prepare); ok {
		return x.Prepare
	}
	return nil
}

func (m *Msg) GetCommit() *Commit {
	if x, ok := m.GetType().(*Msg_Commit); ok {
		return x.Commit
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Msg) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Msg_OneofMarshaler, _Msg_OneofUnmarshaler, _Msg_OneofSizer, []interface{}{
		(*Msg_Preprepare)(nil),
		(*Msg_Prepare)(nil),
		(*Msg_Commit)(nil),
	}
}

func _Msg_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Msg)
	// type
	switch x := m.Type.(type) {
	case *Msg_Preprepare:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Preprepare); err != nil {
			return err
		}
	case *Msg_Prepare:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Prepare); err != nil {
			return err
		}
	case *Msg_Commit:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Commit); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Msg.Type has unexpected type %T", x)
	}
	return nil
}

func _Msg_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Msg)
	switch tag {
	case 1: // type.preprepare
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Preprepare)
		err := b.DecodeMessage(msg)
		m.Type = &Msg_Preprepare{msg}
		return true, err
	case 2: // type.prepare
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Prepare)
		err := b.DecodeMessage(msg)
		m.Type = &Msg_Prepare{msg}
		return true, err
	case 3: // type.commit
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Commit)
		err := b.DecodeMessage(msg)
		m.Type = &Msg_Commit{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Msg_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Msg)
	// type
	switch x := m.Type.(type) {
	case *Msg_Preprepare:
		s := proto.Size(x.Preprepare)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Msg_Prepare:
		s := proto.Size(x.Prepare)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Msg_Commit:
		s := proto.Size(x.Commit)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type Preprepare struct {
	SeqNo   uint64   `protobuf:"varint,1,opt,name=seq_no,json=seqNo" json:"seq_no,omitempty"`
	View    uint64   `protobuf:"varint,2,opt,name=view" json:"view,omitempty"`
	Payload [][]byte `protobuf:"bytes,3,rep,name=payload,proto3" json:"payload,omitempty"`
}

func (m *Preprepare) Reset()                    { *m = Preprepare{} }
func (m *Preprepare) String() string            { return proto.CompactTextString(m) }
func (*Preprepare) ProtoMessage()               {}
func (*Preprepare) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type Prepare struct {
	SeqNo  uint64 `protobuf:"varint,1,opt,name=seq_no,json=seqNo" json:"seq_no,omitempty"`
	View   uint64 `protobuf:"varint,2,opt,name=view" json:"view,omitempty"`
	Digest []byte `protobuf:"bytes,3,opt,name=digest,proto3" json:"digest,omitempty"`
}

func (m *Prepare) Reset()                    { *m = Prepare{} }
func (m *Prepare) String() string            { return proto.CompactTextString(m) }
func (*Prepare) ProtoMessage()               {}
func (*Prepare) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type Commit struct {
	SeqNo  uint64 `protobuf:"varint,1,opt,name=seq_no,json=seqNo" json:"seq_no,omitempty"`
	View   uint64 `protobuf:"varint,2,opt,name=view" json:"view,omitempty"`
	Digest []byte `protobuf:"bytes,3,opt,name=digest,proto3" json:"digest,omitempty"`
}

func (m *Commit) Reset()                    { *m = Commit{} }
func (m *Commit) String() string            { return proto.CompactTextString(m) }
func (*Commit) ProtoMessage()               {}
func (*Commit) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func init() {
	proto.RegisterType((*Msg)(nil), "mirbftpb.Msg")
	proto.RegisterType((*Preprepare)(nil), "mirbftpb.Preprepare")
	proto.RegisterType((*Prepare)(nil), "mirbftpb.Prepare")
	proto.RegisterType((*Commit)(nil), "mirbftpb.Commit")
}

func init() { proto.RegisterFile("mirbft.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 228 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xc9, 0xcd, 0x2c, 0x4a,
	0x4a, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x80, 0xf0, 0x0a, 0x92, 0x94, 0x16,
	0x30, 0x72, 0x31, 0xfb, 0x16, 0xa7, 0x0b, 0x99, 0x71, 0x71, 0x15, 0x14, 0xa5, 0x82, 0x50, 0x62,
	0x51, 0xaa, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0xb7, 0x91, 0x88, 0x1e, 0x4c, 0x99, 0x5e, 0x00, 0x5c,
	0xce, 0x83, 0x21, 0x08, 0x49, 0xa5, 0x90, 0x2e, 0x17, 0x3b, 0x4c, 0x13, 0x13, 0x58, 0x93, 0x20,
	0xaa, 0x26, 0x88, 0x0e, 0x98, 0x1a, 0x21, 0x2d, 0x2e, 0xb6, 0xe4, 0xfc, 0xdc, 0xdc, 0xcc, 0x12,
	0x09, 0x66, 0xb0, 0x6a, 0x01, 0x84, 0x6a, 0x67, 0xb0, 0xb8, 0x07, 0x43, 0x10, 0x54, 0x85, 0x13,
	0x1b, 0x17, 0x4b, 0x49, 0x65, 0x41, 0xaa, 0x52, 0x20, 0x17, 0x17, 0xc2, 0x7a, 0x21, 0x51, 0x2e,
	0xb6, 0xe2, 0xd4, 0xc2, 0xf8, 0xbc, 0x7c, 0xb0, 0x23, 0x59, 0x82, 0x58, 0x8b, 0x53, 0x0b, 0xfd,
	0xf2, 0x85, 0x84, 0xb8, 0x58, 0xca, 0x32, 0x53, 0xcb, 0xc1, 0x8e, 0x60, 0x09, 0x02, 0xb3, 0x85,
	0x24, 0xb8, 0xd8, 0x0b, 0x12, 0x2b, 0x73, 0xf2, 0x13, 0x53, 0x24, 0x98, 0x15, 0x98, 0x35, 0x78,
	0x82, 0x60, 0x5c, 0x25, 0x1f, 0x2e, 0xf6, 0x00, 0xd2, 0xcd, 0x13, 0xe3, 0x62, 0x4b, 0xc9, 0x4c,
	0x4f, 0x2d, 0x86, 0x38, 0x9e, 0x27, 0x08, 0xca, 0x53, 0xf2, 0xe6, 0x62, 0x83, 0x38, 0x9e, 0x0a,
	0x86, 0x25, 0xb1, 0x81, 0x63, 0xc8, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x8a, 0xa3, 0x7a, 0xb6,
	0xb1, 0x01, 0x00, 0x00,
}

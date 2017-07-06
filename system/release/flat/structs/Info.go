// automatically generated by the FlatBuffers compiler, do not modify

package structs

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Info struct {
	_tab flatbuffers.Table
}

func GetRootAsInfo(buf []byte, offset flatbuffers.UOffsetT) *Info {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Info{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *Info) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Info) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Info) ID() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Info) IDLike() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Info) PrettyName() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Info) Version() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Info) VersionID() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Info) HomeURL() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Info) BugReportURL() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func InfoStart(builder *flatbuffers.Builder) { builder.StartObject(8) }
func InfoAddName(builder *flatbuffers.Builder, Name flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(Name), 0) }
func InfoAddID(builder *flatbuffers.Builder, ID flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(ID), 0) }
func InfoAddIDLike(builder *flatbuffers.Builder, IDLike flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(IDLike), 0) }
func InfoAddPrettyName(builder *flatbuffers.Builder, PrettyName flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(PrettyName), 0) }
func InfoAddVersion(builder *flatbuffers.Builder, Version flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(Version), 0) }
func InfoAddVersionID(builder *flatbuffers.Builder, VersionID flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(VersionID), 0) }
func InfoAddHomeURL(builder *flatbuffers.Builder, HomeURL flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(6, flatbuffers.UOffsetT(HomeURL), 0) }
func InfoAddBugReportURL(builder *flatbuffers.Builder, BugReportURL flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(7, flatbuffers.UOffsetT(BugReportURL), 0) }
func InfoEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
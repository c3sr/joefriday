// automatically generated, do not modify

package flat

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Util struct {
	_tab flatbuffers.Table
}

func (rcv *Util) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Util) ID() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Util) Usage() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Util) User() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Util) Nice() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Util) System() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Util) Idle() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Util) IOWait() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0
}

func UtilStart(builder *flatbuffers.Builder) { builder.StartObject(7) }
func UtilAddID(builder *flatbuffers.Builder, ID flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(ID), 0) }
func UtilAddUsage(builder *flatbuffers.Builder, Usage float32) { builder.PrependFloat32Slot(1, Usage, 0) }
func UtilAddUser(builder *flatbuffers.Builder, User float32) { builder.PrependFloat32Slot(2, User, 0) }
func UtilAddNice(builder *flatbuffers.Builder, Nice float32) { builder.PrependFloat32Slot(3, Nice, 0) }
func UtilAddSystem(builder *flatbuffers.Builder, System float32) { builder.PrependFloat32Slot(4, System, 0) }
func UtilAddIdle(builder *flatbuffers.Builder, Idle float32) { builder.PrependFloat32Slot(5, Idle, 0) }
func UtilAddIOWait(builder *flatbuffers.Builder, IOWait float32) { builder.PrependFloat32Slot(6, IOWait, 0) }
func UtilEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }

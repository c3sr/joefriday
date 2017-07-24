// automatically generated by the FlatBuffers compiler, do not modify

package structs

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type CPU struct {
	_tab flatbuffers.Table
}

func (rcv *CPU) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *CPU) Processor() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *CPU) CPUMHz() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *CPU) PhysicalID() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *CPU) CoreID() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *CPU) APICID() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func CPUStart(builder *flatbuffers.Builder) { builder.StartObject(5) }
func CPUAddProcessor(builder *flatbuffers.Builder, Processor int16) { builder.PrependInt16Slot(0, Processor, 0) }
func CPUAddCPUMHz(builder *flatbuffers.Builder, CPUMHz float32) { builder.PrependFloat32Slot(1, CPUMHz, 0.0) }
func CPUAddPhysicalID(builder *flatbuffers.Builder, PhysicalID int16) { builder.PrependInt16Slot(2, PhysicalID, 0) }
func CPUAddCoreID(builder *flatbuffers.Builder, CoreID int16) { builder.PrependInt16Slot(3, CoreID, 0) }
func CPUAddAPICID(builder *flatbuffers.Builder, APICID int16) { builder.PrependInt16Slot(4, APICID, 0) }
func CPUEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
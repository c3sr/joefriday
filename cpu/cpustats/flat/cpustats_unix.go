// Copyright 2016 Joel Scoble and The JoeFriday authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package cpustats handles Flatbuffer based processing of kernel activity,
// /proc/stat. The first Stats.CPU element aggregates the values for all other
// CPU elements. The values are aggregated since system boot. Instead of
// returning a Go struct, it returns Flatbuffer serialized bytes. A function to
// deserialize the Flatbuffer serialized bytes into a cpustats.CPUStats struct
// is provided.
//
// Note: the package name is cpustats and not the final element of the import
// path (flat). 
package cpustats

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/c3sr/joefriday"
	stats "github.com/c3sr/joefriday/cpu/cpustats"
	"github.com/c3sr/joefriday/cpu/cpustats/flat/structs"
)

// Profiler is used to process the /proc/stats file as Flatbuffer serialized
// bytes.
type Profiler struct {
	*stats.Profiler
	*fb.Builder
}

// Returns an initialized profiler that uses Flatbuffers.
func NewProfiler() (prof *Profiler, err error) {
	p, err := stats.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns information about current kernel activity as Flatbuffer
// serialized bytes.
func (prof *Profiler) Get() ([]byte, error) {
	stts, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(stts), nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns information about current kernel activity as Flatbuffer
// serialized bytes using the package's global Profiler.
func Get() (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	} else {
		std.Builder.Reset()
	}

	return std.Get()
}

// Serialize cpustats.CPUStats using Flatbuffers.
func (prof *Profiler) Serialize(stts *stats.CPUStats) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	cpusF := make([]fb.UOffsetT, len(stts.CPU))
	ids := make([]fb.UOffsetT, len(stts.CPU))
	for i := 0; i < len(ids); i++ {
		ids[i] = prof.Builder.CreateString(stts.CPU[i].ID)
	}
	for i := 0; i < len(cpusF); i++ {
		structs.CPUStart(prof.Builder)
		structs.CPUAddID(prof.Builder, ids[i])
		structs.CPUAddUser(prof.Builder, stts.CPU[i].User)
		structs.CPUAddNice(prof.Builder, stts.CPU[i].Nice)
		structs.CPUAddSystem(prof.Builder, stts.CPU[i].System)
		structs.CPUAddIdle(prof.Builder, stts.CPU[i].Idle)
		structs.CPUAddIOWait(prof.Builder, stts.CPU[i].IOWait)
		structs.CPUAddIRQ(prof.Builder, stts.CPU[i].IRQ)
		structs.CPUAddSoftIRQ(prof.Builder, stts.CPU[i].SoftIRQ)
		structs.CPUAddSteal(prof.Builder, stts.CPU[i].Steal)
		structs.CPUAddQuest(prof.Builder, stts.CPU[i].Quest)
		structs.CPUAddQuestNice(prof.Builder, stts.CPU[i].QuestNice)
		cpusF[i] = structs.CPUEnd(prof.Builder)
	}
	structs.CPUStatsStartCPUVector(prof.Builder, len(cpusF))
	for i := len(cpusF) - 1; i >= 0; i-- {
		prof.Builder.PrependUOffsetT(cpusF[i])
	}
	cpusV := prof.Builder.EndVector(len(cpusF))
	structs.CPUStatsStart(prof.Builder)
	structs.CPUStatsAddClkTck(prof.Builder, stts.ClkTck)
	structs.CPUStatsAddTimestamp(prof.Builder, stts.Timestamp)
	structs.CPUStatsAddCtxt(prof.Builder, stts.Ctxt)
	structs.CPUStatsAddBTime(prof.Builder, stts.BTime)
	structs.CPUStatsAddProcesses(prof.Builder, stts.Processes)
	structs.CPUStatsAddCPU(prof.Builder, cpusV)
	prof.Builder.Finish(structs.CPUStatsEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize cpustats.CPUStats with Flatbuffers using the package's global
// Profiler.
func Serialize(stts *stats.CPUStats) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(stts), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserializes them as
// cpustats.CPUStats.
func Deserialize(p []byte) *stats.CPUStats {
	statsS := &stats.CPUStats{}
	cpuF := &structs.CPU{}
	statsF := structs.GetRootAsCPUStats(p, 0)
	statsS.ClkTck = statsF.ClkTck()
	statsS.Timestamp = statsF.Timestamp()
	statsS.Ctxt = statsF.Ctxt()
	statsS.BTime = statsF.BTime()
	statsS.Processes = statsF.Processes()
	len := statsF.CPULength()
	statsS.CPU = make([]stats.CPU, len)
	for i := 0; i < len; i++ {
		var cpu stats.CPU
		if statsF.CPU(cpuF, i) {
			cpu.ID = string(cpuF.ID())
			cpu.User = cpuF.User()
			cpu.Nice = cpuF.Nice()
			cpu.System = cpuF.System()
			cpu.Idle = cpuF.Idle()
			cpu.IOWait = cpuF.IOWait()
			cpu.IRQ = cpuF.IRQ()
			cpu.SoftIRQ = cpuF.SoftIRQ()
			cpu.Steal = cpuF.Steal()
			cpu.Quest = cpuF.Quest()
			cpu.QuestNice = cpuF.QuestNice()
		}
		statsS.CPU[i] = cpu
	}
	return statsS
}

// Ticker delivers the system's kernel activity at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan []byte
	*Profiler
}

// NewTicker returns a new Ticker containing a Data channel that delivers the
// data at intervals and an error channel that delivers any errors encountered.
// Stop the ticker to signal the ticker to stop running. Stopping the ticker
// does not close the Data channel; call Close to close both the ticker and the
// data channel.
func NewTicker(d time.Duration) (joe.Tocker, error) {
	p, err := NewProfiler()
	if err != nil {
		return nil, err
	}
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan []byte), Profiler: p}
	go t.Run()
	return &t, nil
}

// Run runs the ticker.
func (t *Ticker) Run() {
	for {
		select {
		case <-t.Done:
			return
		case <-t.C:
			p, err := t.Get()
			if err != nil {
				t.Errs <- err
				continue
			}
			t.Data <- p
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}

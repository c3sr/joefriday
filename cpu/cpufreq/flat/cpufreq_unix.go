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

// Package cpufreq provides the current CPU frequency, in MHz, as reported by
// /proc/cpuinfo. Instead of returning a Go struct, it returns Flatbuffer
// serialized bytes. A function to deserialize the Flatbuffer serialized bytes
// into a cpufreq.Frequency struct is provided.
//
// Note: the package name is cpufreq and not the final element of the import
// path (flat). 
package cpufreq

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/rai-project/joefriday"
	freq "github.com/rai-project/joefriday/cpu/cpufreq"
	"github.com/rai-project/joefriday/cpu/cpufreq/flat/structs"
)

// Profiler is used to process the frequency information as Flatbuffers
// serialized bytes.
type Profiler struct {
	*freq.Profiler
	*fb.Builder
}

// Initializes and returns a cpuinfo profiler.
func NewProfiler() (p *Profiler, err error) {
	prof, err := freq.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: prof, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the Frequency as Flatbuffer serialized bytes.
func (p *Profiler) Get() ([]byte, error) {
	inf, err := p.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return p.Serialize(inf), nil
}

var std *Profiler    // global for convenience; lazily instantiated.
var stdMu sync.Mutex // protects access

// Get returns the Frequency as Flatbuffer serialized bytes using the package's
// global profiler.
func Get() (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Get()
}

// Serialize serializes Frequency using Flatbuffers.
func (p *Profiler) Serialize(f *freq.Frequency) []byte {
	// ensure the Builder is in a usable state.
	p.Builder.Reset()
	uoffs := make([]fb.UOffsetT, len(f.CPU))
	for i, cpu := range f.CPU {
		uoffs[i] = p.SerializeCPU(&cpu)
	}
	structs.FrequencyStartCPUVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	cpusV := p.Builder.EndVector(len(uoffs))
	structs.FrequencyStart(p.Builder)
	structs.FrequencyAddTimestamp(p.Builder, f.Timestamp)
	structs.FrequencyAddSockets(p.Builder, f.Sockets)
	structs.FrequencyAddCPU(p.Builder, cpusV)
	p.Builder.Finish(structs.FrequencyEnd(p.Builder))
	b := p.Builder.Bytes[p.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}

// Serialize serializes a CPU using flatbuffers and returns the resulting
// UOffsetT.
func (p *Profiler) SerializeCPU(cpu *freq.CPU) fb.UOffsetT {
	structs.CPUStart(p.Builder)
	structs.CPUAddProcessor(p.Builder, cpu.Processor)
	structs.CPUAddCPUMHz(p.Builder, cpu.CPUMHz)
	structs.CPUAddPhysicalID(p.Builder, cpu.PhysicalID)
	structs.CPUAddCoreID(p.Builder, cpu.CoreID)
	structs.CPUAddAPICID(p.Builder, cpu.APICID)
	return structs.CPUEnd(p.Builder)
}

// Serialize cpufreq.Frequency using the package global profiler.
func Serialize(f *freq.Frequency) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(f), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserializes them
// as cpufreq.Frequency.
func Deserialize(p []byte) *freq.Frequency {
	ff := structs.GetRootAsFrequency(p, 0)
	l := ff.CPULength()
	f := &freq.Frequency{}
	fCPU := &structs.CPU{}
	cpu := freq.CPU{}
	f.Timestamp = ff.Timestamp()
	f.Sockets = ff.Sockets()
	for i := 0; i < l; i++ {
		if !ff.CPU(fCPU, i) {
			continue
		}
		cpu.Processor = fCPU.Processor()
		cpu.CPUMHz = fCPU.CPUMHz()
		cpu.PhysicalID = fCPU.PhysicalID()
		cpu.CoreID = fCPU.CoreID()
		cpu.APICID = fCPU.APICID()
		f.CPU = append(f.CPU, cpu)
	}
	return f
}

// Ticker delivers the CPU Frequencies at intervals.
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
	// ticker
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

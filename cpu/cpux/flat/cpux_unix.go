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

// Package cpux provides information about a system's cpus, where X is the
// integer of each CPU on the system, e.g. cpu0, cpu1, etc. On linux systems
// this comes from the sysfs filesystem. Not all paths are available on all
// systems, e.g. /sys/devices/system/cpu/cpuX/cpufreq and its children may not
// exist on some systems. If the system doesn't have a particular path, the
// field's value will be the type's zero value. Instead of returning a Go
// struct, Flatbuffer serialized bytes are returned. A function to deserialize
// the Flatbuffer serialized bytes into a cpux.CPU struct is provided.
//
// Note: the package name is cpux and not the final element of the import path
// (flat).
package cpux

import (
	"sync"

	fb "github.com/google/flatbuffers/go"
	"github.com/rai-project/joefriday/cpu/cpux"
	"github.com/rai-project/joefriday/cpu/cpux/flat/structs"
)

// Profiler is used to process the cpux information as Flatbuffers serialized
// bytes.
type Profiler struct {
	*cpux.Profiler
	*fb.Builder
}

// Initializes and returns a cpux profiler.
func NewProfiler() *Profiler {
	prof := cpux.NewProfiler()
	return &Profiler{Profiler: prof, Builder: fb.NewBuilder(0)}
}

// Get returns the cpux information as Flatbuffer serialized bytes.
func (p *Profiler) Get() ([]byte, error) {
	cpus, err := p.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return p.Serialize(cpus), nil
}

var std *Profiler    // global for convenience; lazily instantiated.
var stdMu sync.Mutex // protects access

// Get returns the cpux information as Flatbuffer serialized bytes using the
// package's global profiler.
func Get() (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std = NewProfiler()
	}
	return std.Get()
}

// Serialize serializes cpux.CPUs using Flatbuffers.
func (p *Profiler) Serialize(cpus *cpux.CPUs) []byte {
	// ensure the Builder is in a usable state.
	p.Builder.Reset()
	possible := p.Builder.CreateString(cpus.Possible)
	online := p.Builder.CreateString(cpus.Online)
	offline := p.Builder.CreateString(cpus.Offline)
	present := p.Builder.CreateString(cpus.Present)
	uoffs := make([]fb.UOffsetT, len(cpus.CPU))
	for i, cpu := range cpus.CPU {
		uoffs[i] = p.SerializeCPU(&cpu)
	}
	structs.CPUsStartCPUVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	cpusV := p.Builder.EndVector(len(uoffs))
	structs.CPUsStart(p.Builder)
	structs.CPUsAddSockets(p.Builder, cpus.Sockets)
	structs.CPUsAddPossible(p.Builder, possible)
	structs.CPUsAddOnline(p.Builder, online)
	structs.CPUsAddOffline(p.Builder, offline)
	structs.CPUsAddPresent(p.Builder, present)
	structs.CPUsAddCPU(p.Builder, cpusV)
	p.Builder.Finish(structs.CPUsEnd(p.Builder))
	b := p.Builder.Bytes[p.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}

// SerializeCPU serializes a CPU using flatbuffers and returns the resulting
// UOffsetT.
func (p *Profiler) SerializeCPU(cpu *cpux.CPU) fb.UOffsetT {
	uoffs := make([]fb.UOffsetT, len(cpu.CacheIDs))
	// serialize cache info in order; the flatbuffer table will have the info
	// in order so a separate cache ID list isn't necessary for flatbuffers.
	for i, id := range cpu.CacheIDs {
		// If the ID doesn't exist, the 0 value will be used
		inf := cpu.Cache[id]
		uoffs[i] = p.SerializeCache(id, inf)
	}
	structs.CPUStartCacheVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	cache := p.Builder.EndVector(len(uoffs))

	structs.CPUStart(p.Builder)
	structs.CPUAddPhysicalPackageID(p.Builder, cpu.PhysicalPackageID)
	structs.CPUAddCoreID(p.Builder, cpu.CoreID)
	structs.CPUAddMHzMin(p.Builder, cpu.MHzMin)
	structs.CPUAddMHzMax(p.Builder, cpu.MHzMax)
	structs.CPUAddCache(p.Builder, cache)
	return structs.CPUEnd(p.Builder)
}

// SerializeCache serializes a cache entry using flatbuffers and returns the
// resulting UOffsetT.
func (p *Profiler) SerializeCache(id, inf string) fb.UOffsetT {
	cID := p.Builder.CreateString(id)
	cInf := p.Builder.CreateString(inf)
	structs.CacheInfStart(p.Builder)
	structs.CacheInfAddID(p.Builder, cID)
	structs.CacheInfAddSize(p.Builder, cInf)
	return structs.CacheInfEnd(p.Builder)
}

// Serialize cpux.CPUs using the package global profiler.
func Serialize(cpus *cpux.CPUs) (p []byte) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std = NewProfiler()
	}
	return std.Serialize(cpus)
}

// Deserialize takes some Flatbuffer serialized bytes and deserializes them
// as cpux.CPUs.
func Deserialize(p []byte) *cpux.CPUs {
	fcpus := structs.GetRootAsCPUs(p, 0)
	l := fcpus.CPULength()
	cpus := &cpux.CPUs{}
	fCPU := &structs.CPU{}
	fCache := &structs.CacheInf{}
	cpu := cpux.CPU{}
	cpus.Sockets = fcpus.Sockets()
	cpus.Possible = string(fcpus.Possible())
	cpus.Online = string(fcpus.Online())
	cpus.Present = string(fcpus.Present())
	for i := 0; i < l; i++ {
		if !fcpus.CPU(fCPU, i) {
			continue
		}
		cpu.PhysicalPackageID = fCPU.PhysicalPackageID()
		cpu.CoreID = fCPU.CoreID()
		cpu.MHzMin = fCPU.MHzMin()
		cpu.MHzMax = fCPU.MHzMax()
		caches := fCPU.CacheLength()
		cpu.CacheIDs = make([]string, 0, caches)
		cpu.Cache = make(map[string]string, caches)
		for j := 0; j < caches; j++ {
			if !fCPU.Cache(fCache, j) {
				continue
			}
			cpu.CacheIDs = append(cpu.CacheIDs, string(fCache.ID()))
			cpu.Cache[cpu.CacheIDs[j]] = string(fCache.Size())
		}
		cpus.CPU = append(cpus.CPU, cpu)
	}
	return cpus
}

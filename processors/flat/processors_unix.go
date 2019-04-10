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

// Package processors gathers information about the physical processors on a
// system by parsing the information from /procs/cpuinfo and the sysfs. This
// package gathers basic information about sockets, physical processors, etc.
// on the system. For multi-socket systems, it is assumed that all of the
// processors are the same. Instead of returning a Go struct, Flatbuffer
// serialized bytes are returned. A function to deserialize the Flatbuffer
// serialized bytes into a processors.Processors struct is provided.
//
// CPUMHz currently provides the current speed of the first core encountered
// for each physical processor. Modern x86/x86-64 cores have the ability to
// shift their speed so this is just a point in time data point for that core;
// there may be other cores on the processor that are at higher and lower
// speeds at the time the data is read. This field is more useful for other
// architectures. For x86/x86-64 cores, the MHzMin and MHzMax fields provide
// information about the range of speeds that are possible for the cores.
//
// Note: the package name is processors and not the final element of the import
// path (flat).
package processors

import (
	"sync"

	fb "github.com/google/flatbuffers/go"
	"github.com/rai-project/joefriday/node"
	"github.com/rai-project/joefriday/processors"
	"github.com/rai-project/joefriday/processors/flat/structs"
)

// Profiler is used to get the processor information, as Flatbuffers serialized
// bytes, by processing the /proc/cpuinfo file.
type Profiler struct {
	*processors.Profiler
	*fb.Builder
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (p *Profiler, err error) {
	prof, err := processors.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: prof, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the processor information as Flatbuffer serialized bytes.
func (p *Profiler) Get() ([]byte, error) {
	proc, err := p.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return p.Serialize(proc), nil
}

var std *Profiler    // global for convenience; lazily instantiated.
var stdMu sync.Mutex // protects access

// Get returns the current processor info as Flatbuffer serialized bytes using
// the package's global Profiler.
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

// Serialize serializes Processors using Flatbuffers.
func (p *Profiler) Serialize(procs *processors.Processors) []byte {
	// ensure the Builder is in a usable state.
	p.Builder.Reset()
	architecture := p.Builder.CreateString(procs.Architecture)
	vendorID := p.Builder.CreateString(procs.VendorID)
	cpuFamily := p.Builder.CreateString(procs.CPUFamily)
	model := p.Builder.CreateString(procs.Model)
	modelName := p.Builder.CreateString(procs.ModelName)
	stepping := p.Builder.CreateString(procs.Stepping)
	microcode := p.Builder.CreateString(procs.Microcode)
	cacheSize := p.Builder.CreateString(procs.CacheSize)
	possible := p.Builder.CreateString(procs.Possible)
	present := p.Builder.CreateString(procs.Present)
	offline := p.Builder.CreateString(procs.Offline)
	online := p.Builder.CreateString(procs.Online)
	virtualization := p.Builder.CreateString(procs.Virtualization)

	uoffs := make([]fb.UOffsetT, len(procs.CacheIDs))
	// serialize cache info in order; the flatbuffer table will have the info
	// in order so a separate cache ID list isn't necessary for flatbuffers.
	for i, id := range procs.CacheIDs {
		// If the ID doesn't exist, the 0 value will be used
		inf := procs.Cache[id]
		uoffs[i] = p.SerializeCache(id, inf)
	}
	structs.ProcessorsStartCacheVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	cache := p.Builder.EndVector(len(uoffs))

	uoffs = make([]fb.UOffsetT, len(procs.Flags))
	for i, flag := range procs.Flags {
		uoffs[i] = p.Builder.CreateString(flag)
	}
	structs.ProcessorsStartFlagsVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	flags := p.Builder.EndVector(len(uoffs))

	uoffs = make([]fb.UOffsetT, len(procs.Bugs))
	for i, bug := range procs.Bugs {
		uoffs[i] = p.Builder.CreateString(bug)
	}
	structs.ProcessorsStartBugsVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	bugs := p.Builder.EndVector(len(uoffs))

	uoffs = make([]fb.UOffsetT, len(procs.OpModes))
	for i := range procs.OpModes {
		uoffs[i] = p.Builder.CreateString(procs.OpModes[i])
	}
	structs.ProcessorsStartOpModesVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	modes := p.Builder.EndVector(len(uoffs))
	uoffs = make([]fb.UOffsetT, len(procs.NumaNodeCPUs))
	for i := range procs.NumaNodeCPUs {
		uoffs[i] = p.SerializeNumaNodeCPUs(&procs.NumaNodeCPUs[i])
	}
	structs.ProcessorsStartNumaNodeCPUsVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	nodeCPUs := p.Builder.EndVector(len(uoffs))
	structs.ProcessorsStart(p.Builder)
	structs.ProcessorsAddTimestamp(p.Builder, procs.Timestamp)
	structs.ProcessorsAddArchitecture(p.Builder, architecture)
	structs.ProcessorsAddCPUs(p.Builder, int32(procs.CPUs))
	structs.ProcessorsAddPossible(p.Builder, possible)
	structs.ProcessorsAddPresent(p.Builder, present)
	structs.ProcessorsAddOffline(p.Builder, offline)
	structs.ProcessorsAddOnline(p.Builder, online)
	structs.ProcessorsAddSockets(p.Builder, procs.Sockets)
	structs.ProcessorsAddCoresPerSocket(p.Builder, procs.CoresPerSocket)
	structs.ProcessorsAddThreadsPerCore(p.Builder, procs.ThreadsPerCore)
	structs.ProcessorsAddVendorID(p.Builder, vendorID)
	structs.ProcessorsAddCPUFamily(p.Builder, cpuFamily)
	structs.ProcessorsAddModel(p.Builder, model)
	structs.ProcessorsAddModelName(p.Builder, modelName)
	structs.ProcessorsAddStepping(p.Builder, stepping)
	structs.ProcessorsAddMicrocode(p.Builder, microcode)
	structs.ProcessorsAddCPUMHz(p.Builder, procs.CPUMHz)
	structs.ProcessorsAddMHzMin(p.Builder, procs.MHzMin)
	structs.ProcessorsAddMHzMax(p.Builder, procs.MHzMax)
	structs.ProcessorsAddBogoMIPS(p.Builder, procs.BogoMIPS)
	structs.ProcessorsAddCacheSize(p.Builder, cacheSize)
	structs.ProcessorsAddCache(p.Builder, cache)
	structs.ProcessorsAddFlags(p.Builder, flags)
	structs.ProcessorsAddBugs(p.Builder, bugs)
	structs.ProcessorsAddOpModes(p.Builder, modes)
	structs.ProcessorsAddVirtualization(p.Builder, virtualization)
	structs.ProcessorsAddNumaNodes(p.Builder, procs.NumaNodes)
	structs.ProcessorsAddNumaNodeCPUs(p.Builder, nodeCPUs)
	p.Builder.Finish(structs.ProcessorsEnd(p.Builder))
	b := p.Builder.Bytes[p.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
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

func (p *Profiler) SerializeNumaNodeCPUs(n *node.Node) fb.UOffsetT {
	list := p.Builder.CreateString(n.CPUList)
	structs.NodeStart(p.Builder)
	structs.NodeAddID(p.Builder, n.ID)
	structs.NodeAddCPUList(p.Builder, list)
	return structs.NodeEnd(p.Builder)
}

// Serialize processors information.
func Serialize(proc *processors.Processors) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(proc), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserializes them
// as processors.Processors.
func Deserialize(p []byte) *processors.Processors {
	flatP := structs.GetRootAsProcessors(p, 0)
	procs := &processors.Processors{}
	flatCache := &structs.CacheInf{}
	procs.Timestamp = flatP.Timestamp()
	procs.Architecture = string(flatP.Architecture())
	procs.CPUs = flatP.CPUs()
	procs.Possible = string(flatP.Possible())
	procs.Present = string(flatP.Present())
	procs.Offline = string(flatP.Offline())
	procs.Online = string(flatP.Online())
	procs.Sockets = flatP.Sockets()
	procs.CoresPerSocket = flatP.CoresPerSocket()
	procs.ThreadsPerCore = flatP.ThreadsPerCore()
	procs.VendorID = string(flatP.VendorID())
	procs.CPUFamily = string(flatP.CPUFamily())
	procs.Model = string(flatP.Model())
	procs.ModelName = string(flatP.ModelName())
	procs.Stepping = string(flatP.Stepping())
	procs.Microcode = string(flatP.Microcode())
	procs.CPUMHz = flatP.CPUMHz()
	procs.MHzMin = flatP.MHzMin()
	procs.MHzMax = flatP.MHzMax()
	procs.BogoMIPS = flatP.BogoMIPS()
	procs.CacheSize = string(flatP.CacheSize())
	procs.CacheIDs = make([]string, 0, flatP.CacheLength())
	procs.Cache = make(map[string]string, flatP.CacheLength())
	for j := 0; j < flatP.CacheLength(); j++ {
		if !flatP.Cache(flatCache, j) {
			continue
		}
		procs.CacheIDs = append(procs.CacheIDs, string(flatCache.ID()))
		procs.Cache[procs.CacheIDs[j]] = string(flatCache.Size())
	}
	procs.Flags = make([]string, flatP.FlagsLength())
	for i := 0; i < len(procs.Flags); i++ {
		procs.Flags[i] = string(flatP.Flags(i))
	}
	procs.Bugs = make([]string, flatP.BugsLength())
	for i := 0; i < len(procs.Bugs); i++ {
		procs.Bugs[i] = string(flatP.Bugs(i))
	}
	procs.OpModes = make([]string, flatP.OpModesLength())
	for i := 0; i < len(procs.OpModes); i++ {
		procs.OpModes[i] = string(flatP.OpModes(i))
	}
	procs.Virtualization = string(flatP.Virtualization())
	procs.NumaNodes = flatP.NumaNodes()
	var n structs.Node
	procs.NumaNodeCPUs = make([]node.Node, flatP.NumaNodeCPUsLength())
	for i := 0; i < len(procs.NumaNodeCPUs); i++ {
		if !flatP.NumaNodeCPUs(&n, i) {
			continue
		}
		procs.NumaNodeCPUs[i].ID = n.ID()
		procs.NumaNodeCPUs[i].CPUList = string(n.CPUList())
	}
	return procs
}

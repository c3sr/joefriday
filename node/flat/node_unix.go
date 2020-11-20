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

// Package node gets information about the system's NUMA nodes. This looks
// at the sysfs's node tree and extracts information about each node. If the
// node tree doesn't exist on the system, instead of node information, an
// os.ErrNotExist will be returned. Instead of returning a Go struct, the data
// will be returned as Flatbuffer serialized bytes. A function to deserialize
// the Flatbuffer serialized bytes into a node.Nodes struct is provided.
//
// Note: the package name is node and not the final element of the import path
// (flat).
package node

import (
	"sync"

	fb "github.com/google/flatbuffers/go"
	numa "github.com/c3sr/joefriday/node"
	"github.com/c3sr/joefriday/node/flat/structs"
)

// Profiler is used to process the node information as Flatbuffers serialized
// bytes.
type Profiler struct {
	*numa.Profiler
	*fb.Builder
}

// Initializes and returns a node profiler.
func NewProfiler() (p *Profiler) {
	prof := numa.NewProfiler()
	return &Profiler{Profiler: prof, Builder: fb.NewBuilder(0)}
}

// Get returns the node information as Flatbuffer serialized bytes.
func (p *Profiler) Get() ([]byte, error) {
	nodes, err := p.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return p.Serialize(nodes), nil
}

var std *Profiler    // global for convenience; lazily instantiated.
var stdMu sync.Mutex // protects access

// Get returns the node information as Flatbuffer serialized bytes using the
// package's global profiler.
func Get() (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std = NewProfiler()
	}
	return std.Get()
}

// Serialize serializes node.Nodes using Flatbuffers.
func (p *Profiler) Serialize(nodes *numa.Nodes) []byte {
	// ensure the Builder is in a usable state.
	p.Builder.Reset()
	uoffs := make([]fb.UOffsetT, len(nodes.Node))
	for i, node := range nodes.Node {
		uoffs[i] = p.SerializeNode(&node)
	}
	structs.NodesStartNodeVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	nodeV := p.Builder.EndVector(len(uoffs))
	structs.NodesStart(p.Builder)
	structs.NodesAddNode(p.Builder, nodeV)
	p.Builder.Finish(structs.NodesEnd(p.Builder))
	b := p.Builder.Bytes[p.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}

// SerializeNode serializes a Node using flatbuffers and returns the resulting
// UOffsetT.
func (p *Profiler) SerializeNode(node *numa.Node) fb.UOffsetT {
	cpuList := p.Builder.CreateString(node.CPUList)
	structs.NodeStart(p.Builder)
	structs.NodeAddID(p.Builder, node.ID)
	structs.NodeAddCPUList(p.Builder, cpuList)
	return structs.NodeEnd(p.Builder)
}

// Serialize node.Nodes using the package global profiler.
func Serialize(nodes *numa.Nodes) (p []byte) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std = NewProfiler()
	}
	return std.Serialize(nodes)
}

// Deserialize takes some Flatbuffer serialized bytes and deserializes them
// as node.Nodes.
func Deserialize(p []byte) *numa.Nodes {
	fnodes := structs.GetRootAsNodes(p, 0)
	l := fnodes.NodeLength()
	nodes := &numa.Nodes{}
	fNode := &structs.Node{}
	node := numa.Node{}
	nodes.Node = make([]numa.Node, 0, l)
	for i := 0; i < l; i++ {
		if !fnodes.Node(fNode, i) {
			continue
		}
		node.ID = fNode.ID()
		node.CPUList = string(fNode.CPUList())
		nodes.Node = append(nodes.Node, node)
	}
	return nodes
}

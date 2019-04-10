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
// os.ErrNotExist will be returned. Instead of returning a Go struct, JSON
// serialized bytes are returned. A function to deserialize the JSON serialized
// bytes into a node.Nodes struct is provided.
//
// Note: the package name is node and not the final element of the import path
// (json).
package node

import (
	"encoding/json"
	"sync"

	numa "github.com/rai-project/joefriday/node"
)

// Profiler is used to process the node information.
type Profiler struct {
	*numa.Profiler
}

// Initializes and returns a cpuinfo profiler.
func NewProfiler() (prof *Profiler) {
	p := numa.NewProfiler()
	return &Profiler{Profiler: p}
}

// Get returns the current cpuinfo as JSON serialized bytes.
func (prof *Profiler) Get() (p []byte, err error) {
	inf, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf)
}

var std *Profiler
var stdMu sync.Mutex //protects standard to prevent data race on checking/instantiation

// Get returns the current node as JSON serialized bytes using the package's
// global Profiler.
func Get() (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std = NewProfiler()
	}
	return std.Get()
}

// Serialize node.Nodes as JSON.
func (prof *Profiler) Serialize(nodes *numa.Nodes) ([]byte, error) {
	return json.Marshal(nodes)
}

// Serialize node.Nodes as JSON using package globals.
func Serialize(nodes *numa.Nodes) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std = NewProfiler()
	}
	return std.Serialize(nodes)
}

// Marshal is an alias for serialize.
func (prof *Profiler) Marshal(nodes *numa.Nodes) ([]byte, error) {
	return prof.Serialize(nodes)
}

// Marshal is an alias for Serialize using package globals.
func Marshal(nodes *numa.Nodes) ([]byte, error) {
	return std.Serialize(nodes)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// node.Nodes.
func Deserialize(p []byte) (*numa.Nodes, error) {
	nodes := &numa.Nodes{}
	err := json.Unmarshal(p, nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

// Unmarshal is an alias for Deserialize using package globals.
func Unmarshal(p []byte) (*numa.Nodes, error) {
	return Deserialize(p)
}

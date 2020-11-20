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
// os.ErrNotExist will be returned.
package node

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/c3sr/joefriday"
)

const CPUList = "cpulist"

type Nodes struct {
	Node []Node `json:"node"`
}

// NumaNodes returns the number of numa nodes for the system.
func (n *Nodes) NumaNodes() int {
	return len(n.Node)
}

// Information about a specific node.
type Node struct {
	ID      int32  `json:"id"` // max_numa_node returns an int (in C)
	CPUList string `json:"cpu_list"`
}

// Profiler is used to process the system's sysfs node information.
type Profiler struct {
	sysFSSystemPath string
	// the path of th4e sysfs node tree; cached so it doesn't need to be
	// generated for every use.
	nodePath string
}

// Returns an initialized Profiler.
func NewProfiler() (prof *Profiler) {
	prof = &Profiler{}
	prof.SysFSSystemPath(joefriday.SysFSSystem)
	return prof
}

// Rest resources: this does nothing for this implemtation.
func (prof *Profiler) Reset() error {
	return nil
}

// Get the node information. If the node tree doesn't exist an os.ErrNotExist
// will be returned. During processing, any error will be returned along with a
// nil for nodes.
func (prof *Profiler) Get() (nodes *Nodes, err error) {
	nodes = &Nodes{}
	var x int32 // index of nodeX currently being processed.

	// First see if the node dir exists, return any error.
	_, err = os.Stat(prof.nodePath)
	if err != nil {
		return nil, err
	}
	// Loop and increment x after each dir read; stop when the nodeX dir to be
	// processed doesn't exist.
	for {
		var n Node
		p := prof.nodeXPath(x)
		// get the cpulist
		n.ID = x
		n.CPUList, err = prof.CPUList(p)
		if err != nil {
			// if the dir didn't exist; there are no more nodes to process
			if os.IsNotExist(err) {
				return nodes, nil
			}
			// any other error will be passed back
			return nil, err
		}
		nodes.Node = append(nodes.Node, n)
		x++
	}
}

func (prof *Profiler) nodeXPath(x int32) string {
	return filepath.Join(prof.nodePath, fmt.Sprintf("node%d", x))
}

// CPUList returns the string found in the CPUList file or any error that
// occurs.
func (prof *Profiler) CPUList(path string) (string, error) {
	p, err := ioutil.ReadFile(filepath.Join(path, CPUList))
	if err != nil {
		return "", err
	}
	// the list is everything except for the trailing new line
	return string(p[:len(p)-1]), nil
}

// SysFSSystemPath enables overriding the default value. This is for testing
// and should be used outside of tests.
func (prof *Profiler) SysFSSystemPath(s string) {
	prof.sysFSSystemPath = s
	prof.setNodePath()
}

// the path to the sysfs node tree is cached so that it doesn't need to be
// regenerated for each use.
func (prof *Profiler) setNodePath() {
	prof.nodePath = filepath.Join(prof.sysFSSystemPath, "node")
}

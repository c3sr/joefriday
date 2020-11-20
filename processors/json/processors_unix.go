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
// processors are the same. Instead of returning a Go struct, JSON serialized
// bytes are returned. A function to deserialize the JSON serialized bytes into
// a processors.Processors struct is provided.
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
// path (json).
package processors

import (
	"encoding/json"
	"sync"

	procs "github.com/c3sr/joefriday/processors"
)

// Profiler is used to get the processor information, as JSON serialized bytes,
// by processing the /proc/cpuinfo file.
type Profiler struct {
	*procs.Profiler
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (p *Profiler, err error) {
	prof, err := procs.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: prof}, nil
}

// Get returns the processor information as JSON serialized bytes.
func (p *Profiler) Get() (b []byte, err error) {
	proc, err := p.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return p.Serialize(proc)
}

var std *Profiler
var stdMu sync.Mutex //protects standard to prevent a data race on checking/instantiation

// Get returns the processor information as JSON serialized bytes using the
// package's global Profiler.
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

// Serialize processor information.
func (p *Profiler) Serialize(proc *procs.Processors) ([]byte, error) {
	return json.Marshal(proc)
}

// Serialize processor information.
func Serialize(proc *procs.Processors) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(proc)
}

// Marshal is an alias for serialize.
func (p *Profiler) Marshal(proc *procs.Processors) ([]byte, error) {
	return p.Serialize(proc)
}

// Marshal is an alias for Serialize.
func Marshal(proc *procs.Processors) ([]byte, error) {
	return std.Serialize(proc)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// processors.Processors
func Deserialize(p []byte) (*procs.Processors, error) {
	proc := &procs.Processors{}
	err := json.Unmarshal(p, proc)
	if err != nil {
		return nil, err
	}
	return proc, nil
}

// Unmarshal is an alias for Deserialize
func Unmarshal(p []byte) (*procs.Processors, error) {
	return Deserialize(p)
}

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

// Package cpuinfo (json) handles JSON based processing of /proc/cpuinfo.
// Instead of returning a Go struct, it returns JSON serialized bytes. A
// function to deserialize the JSON serialized bytes into a cpuinfo.CPUInfo
// struct is provided.
//
// Note: the package name is cpuinfo and not the final element of the import
// path (json). 
package cpuinfo

import (
	"encoding/json"
	"sync"

	info "github.com/c3sr/joefriday/cpu/cpuinfo"
)

// Profiler is used to process the /proc/cpuinfo file as JSON serialized bytes.
type Profiler struct {
	*info.Profiler
}

// Initializes and returns a cpuinfo profiler.
func NewProfiler() (prof *Profiler, err error) {
	p, err := info.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p}, nil
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

// Get returns the current cpuinfo as JSON serialized bytes using the package's
// global Profiler.
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

// Serialize cpuinfo as JSON.
func (prof *Profiler) Serialize(inf *info.CPUInfo) ([]byte, error) {
	return json.Marshal(inf)
}

// Serialize cpuinfo as JSON using package globals.
func Serialize(inf *info.CPUInfo) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(inf)
}

// Marshal is an alias for serialize.
func (prof *Profiler) Marshal(inf *info.CPUInfo) ([]byte, error) {
	return prof.Serialize(inf)
}

// Marshal is an alias for Serialize using package globals.
func Marshal(inf *info.CPUInfo) ([]byte, error) {
	return std.Serialize(inf)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// cpuinfo.CPUInfo.
func Deserialize(p []byte) (*info.CPUInfo, error) {
	inf := &info.CPUInfo{}
	err := json.Unmarshal(p, inf)
	if err != nil {
		return nil, err
	}
	return inf, nil
}

// Unmarshal is an alias for Deserialize using package globals.
func Unmarshal(p []byte) (*info.CPUInfo, error) {
	return Deserialize(p)
}

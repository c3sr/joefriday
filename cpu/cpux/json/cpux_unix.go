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
// struct, JSON serialized bytes are returned. A function to deserialize the
// JSON serialized bytes into a cpux.CPUs struct is provided.
//
// Note: the package name is cpux and not the final element of the import path
// (json).
package cpux

import (
	"encoding/json"
	"sync"

	x "github.com/c3sr/joefriday/cpu/cpux"
)

// Profiler is used to process the cpuX information.
type Profiler struct {
	*x.Profiler
}

// Initializes and returns a cpuinfo profiler.
func NewProfiler() *Profiler {
	p := x.NewProfiler()
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

// Get returns the current cpux as JSON serialized bytes using the package's
// global Profiler.
func Get() (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std = NewProfiler()
	}
	return std.Get()
}

// Serialize cpux.CPUs as JSON.
func (prof *Profiler) Serialize(cpus *x.CPUs) ([]byte, error) {
	return json.Marshal(cpus)
}

// Serialize cpux.CPUs as JSON using package globals.
func Serialize(cpus *x.CPUs) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std = NewProfiler()
	}
	return std.Serialize(cpus)
}

// Marshal is an alias for serialize.
func (prof *Profiler) Marshal(cpus *x.CPUs) ([]byte, error) {
	return prof.Serialize(cpus)
}

// Marshal is an alias for Serialize using package globals.
func Marshal(cpus *x.CPUs) ([]byte, error) {
	return std.Serialize(cpus)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// cpux.CPUs.
func Deserialize(p []byte) (*x.CPUs, error) {
	cpus := &x.CPUs{}
	err := json.Unmarshal(p, cpus)
	if err != nil {
		return nil, err
	}
	return cpus, nil
}

// Unmarshal is an alias for Deserialize using package globals.
func Unmarshal(p []byte) (*x.CPUs, error) {
	return Deserialize(p)
}

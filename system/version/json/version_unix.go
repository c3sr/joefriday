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

// Package version gets the kernel and version information from the
// /proc/version file. Instead of returning a Go struct, it returns JSON
// serialized bytes. A function to deserialize the JSON serialized bytes into a
// version.Kernel struct is provided.
//
// Note: the package name is version and not the final element of the import
// path (json). 
package version

import (
	"encoding/json"
	"sync"

	v "github.com/c3sr/joefriday/system/version"
)

// Profiler processes the version information, /proc/version, using
// JSON.
type Profiler struct {
	*v.Profiler
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	p, err := v.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p}, nil
}

// Get gets the kernel information from the /proc/version file as JSON
// serialized bytes.
func (prof *Profiler) Get() (p []byte, err error) {
	inf, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf)
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get gets the kernel information from the /proc/version file as JSON
// serialized bytes using the package's global Profiler.
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

// Serialize version.Kernel as JSON.
func (prof *Profiler) Serialize(k *v.Kernel) ([]byte, error) {
	return json.Marshal(k)
}

// Serialize version.Kernel as JSON using the package's global Profiler.
func Serialize(k *v.Kernel) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(k)
}

// Marshal is an alias for Serialize.
func (prof *Profiler) Marshal(k *v.Kernel) ([]byte, error) {
	return prof.Serialize(k)
}

// Marshal is an alias for Serialize that uses the package's global profiler.
func Marshal(k *v.Kernel) ([]byte, error) {
	return Serialize(k)
}

// Deserialize takes some JSON serialized bytes and deserializes them as
// version.Kernel.
func Deserialize(p []byte) (*v.Kernel, error) {
	k := &v.Kernel{}
	err := json.Unmarshal(p, k)
	if err != nil {
		return nil, err
	}
	return k, nil
}

// Unmarshal is an alias for Deserialize.
func Unmarshal(p []byte) (*v.Kernel, error) {
	return Deserialize(p)
}

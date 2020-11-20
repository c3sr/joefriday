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

// Package os provides OS Release information, /etc/os-release. Instead of
// returning a Go struct, it returns JSON serialized bytes. A function to
// deserialize the JSON serialized bytes into am os.OS struct is provided.
//
// Note: the package name is os and not the final element of the import
// path (json).
package os

import (
	"encoding/json"
	"sync"

	o "github.com/c3sr/joefriday/system/os"
)

// Profiler processes the OS release information, /etc/os-release,
// using JSON.
type Profiler struct {
	*o.Profiler
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	p, err := o.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p}, nil
}

// Get gets the OS release information, /etc/os-release, as JSON serialized
// bytes.
func (prof *Profiler) Get() (p []byte, err error) {
	k, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(k)
}

var std *Profiler
var stdMu sync.Mutex //protects standard to prevent a data race on checking/instantiation

// Get gets the OS release information, /etc/os-release, as JSON serialized
// bytes using the package's global Profiler.
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

// Serialize os.OS as JSON
func (prof *Profiler) Serialize(os *o.OS) ([]byte, error) {
	return json.Marshal(os)
}

// Serialize os.OS as JSON using the package's global Profiler.
func Serialize(os *o.OS) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(os)
}

// Marshal is an alias for Serialize.
func (prof *Profiler) Marshal(os *o.OS) ([]byte, error) {
	return prof.Serialize(os)
}

// Marshal is an alias for Serialize using the package's global profiler.
func Marshal(os *o.OS) ([]byte, error) {
	return Serialize(os)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as os.OS.
func Deserialize(p []byte) (*o.OS, error) {
	os := &o.OS{}
	err := json.Unmarshal(p, os)
	if err != nil {
		return nil, err
	}
	return os, nil
}

// Unmarshal is an alias for Deserialize.
func Unmarshal(p []byte) (*o.OS, error) {
	return Deserialize(p)
}

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

// Package json handles JSON based processing of CPU stats.  Instead of
// returning a Go struct, it returns JSON serialized bytes.  A function to
// deserialize the JSON serialized bytes into a stats.Stats struct is
// provided.
package json

import (
	"encoding/json"
	"sync"

	"github.com/mohae/joefriday/cpu/stats"
)

// Profiler is used to process the /proc/stats file, as Stats, using JSON.
type Profiler struct {
	Prof *stats.Profiler
}

// Initializes and returns a cpu Stats profiler.
func New() (prof *Profiler, err error) {
	p, err := stats.New()
	if err != nil {
		return nil, err
	}
	return &Profiler{Prof: p}, nil
}

// Get returns the current Stats as JSON serialized bytes.
func (prof *Profiler) Get() (p []byte, err error) {
	prof.Prof.Reset()
	st, err := prof.Prof.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(st)
}

// TODO: is it even worth it to have this as a global?  Should GetInfo()
// just instantiate a local version and use that?  InfoTicker does...
var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current Stats as JSON serialized bytes using the package's
// global Profiler.
func Get() (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = New()
		if err != nil {
			return nil, err
		}
	}
	return std.Get()
}

// Serialize cpu Stats as JSON
func (prof *Profiler) Serialize(st *stats.Stats) ([]byte, error) {
	return json.Marshal(st)
}

// Marshal is an alias for Serialize
func (prof *Profiler) Marshal(st *stats.Stats) ([]byte, error) {
	return prof.Serialize(st)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// stats.Stats
func Deserialize(p []byte) (*stats.Stats, error) {
	st := &stats.Stats{}
	err := json.Unmarshal(p, st)
	if err != nil {
		return nil, err
	}
	return st, nil
}

// Unmarshal is an alias for Deserialize
func Unmarshal(p []byte) (*stats.Stats, error) {
	return Deserialize(p)
}

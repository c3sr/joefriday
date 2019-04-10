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

// Package membasic processes a subset of the /proc/meminfo file. Instead of
// returning a Go struct, it returns JSON serialized bytes. A function to
// deserialize the JSON serialized bytes into a membasic.Info struct is
// provided. For more detailed information about a system's memory, use the
// meminfo package.
//
// Note: the package name is membasic and not the final element of the import
// path (json). 
package membasic

import (
	"encoding/json"
	"sync"
	"time"

	joe "github.com/rai-project/joefriday"
	basic "github.com/rai-project/joefriday/mem/membasic"
)

// Profiler is used to get the basic memory information, as JSON, by processing
// the /proc/meminfo file.
type Profiler struct {
	*basic.Profiler
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	p, err := basic.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p}, nil
}

// Get returns the current basic memory information as JSON serialized bytes.
func (prof *Profiler) Get() (p []byte, err error) {
	inf, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf)
}

var std *Profiler
var stdMu sync.Mutex //protects standard to prevent a data race on checking/instantiation

// Get returns the current basic memory information as JSON serialized bytes
// using the package's global Profiler.
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

// Serialize the basic memory information using JSON.
func (prof *Profiler) Serialize(inf *basic.Info) ([]byte, error) {
	return json.Marshal(inf)
}

// Serialize the basic memory information using JSON with the package's global
// Profiler.
func Serialize(inf *basic.Info) (p []byte, err error) {
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

// Marshal is an alias for Serialize.
func (prof *Profiler) Marshal(inf *basic.Info) ([]byte, error) {
	return prof.Serialize(inf)
}

// Marshal is an alias for Serialize that uses the package's global profiler.
func Marshal(inf *basic.Info) ([]byte, error) {
	return Serialize(inf)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// membasic.Info.
func Deserialize(p []byte) (*basic.Info, error) {
	info := &basic.Info{}
	err := json.Unmarshal(p, info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// Unmarshal is an alias for Deserialize.
func Unmarshal(p []byte) (*basic.Info, error) {
	return Deserialize(p)
}

// Ticker delivers the system's basic memory information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan []byte
	*Profiler
}

// NewTicker returns a new Ticker containing a Data channel that delivers the
// data at intervals and an error channel that delivers any errors encountered.
// Stop the ticker to signal the ticker to stop running. Stopping the ticker
// does not close the Data channel; call Close to close both the ticker and the
// data channel.
func NewTicker(d time.Duration) (joe.Tocker, error) {
	p, err := NewProfiler()
	if err != nil {
		return nil, err
	}
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan []byte), Profiler: p}
	go t.Run()
	return &t, nil
}

// Run runs the ticker.
func (t *Ticker) Run() {
	for {
		select {
		case <-t.Done:
			return
		case <-t.C:
			p, err := t.Get()
			if err != nil {
				t.Errs <- err
				continue
			}
			t.Data <- p
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}

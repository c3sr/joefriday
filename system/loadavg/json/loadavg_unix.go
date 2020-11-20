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

// Package loadAvg gets loadavg information from the /proc/loadavg file.
// Instead of returning a Go struct, it returns JSON serialized bytes. A
// function to deserialize the JSON serialized bytes into an loadavg.LoadAvg
// struct is provided.
//
// Note: the package name is loadavg and not the final element of the import
// path (json). 
package loadavg

import (
	"encoding/json"
	"sync"
	"time"

	joe "github.com/c3sr/joefriday"
	l "github.com/c3sr/joefriday/system/loadavg"
)

// Profiler is used to process the loadavg information, /proc/loadavg, using
// JSON.
type Profiler struct {
	*l.Profiler
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	p, err := l.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p}, nil
}

// Get returns the current loadavg information as JSON serialized bytes.
func (prof *Profiler) Get() (p []byte, err error) {
	k, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(k)
}

var std *Profiler
var stdMu sync.Mutex //protects standard to prevent a data race on checking/instantiation

// Get returns the current loadavg information as JSON serialized bytes using
// the package's global Profiler.
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

// Serialize loadavg.LoadAvg using JSON.
func (prof *Profiler) Serialize(la l.LoadAvg) ([]byte, error) {
	return json.Marshal(la)
}

// Serialize loadavg.LoadAvg using JSON with the package's global Profiler.
func Serialize(la l.LoadAvg) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(la)
}

// Marshal is an alias for Serialize.
func (prof *Profiler) Marshal(la l.LoadAvg) ([]byte, error) {
	return prof.Serialize(la)
}

// Marshal is an alias for Serialize using the package's global profiler.
func Marshal(la l.LoadAvg) ([]byte, error) {
	return Serialize(la)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// loadavg.LoadAvg.
func Deserialize(p []byte) (la l.LoadAvg, err error) {
	err = json.Unmarshal(p, &la)
	if err != nil {
		return la, err
	}
	return la, nil
}

// Unmarshal is an alias for Deserialize.
func Unmarshal(p []byte) (l.LoadAvg, error) {
	return Deserialize(p)
}

// Ticker delivers the system's loadavg information at intervals.
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

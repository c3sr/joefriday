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

// Package uptime gets the current uptime from the /proc/uptime file. Instead
// of returning a Go struct, it returns JSON serialized bytes. A function to
// deserialize the JSON serialized bytes into an uptime.Uptime struct is
// provided.
//
// Note: the package name is uptime and not the final element of the import
// path (json). 
package uptime

import (
	"encoding/json"
	"sync"
	"time"

	joe "github.com/c3sr/joefriday"
	u "github.com/c3sr/joefriday/system/uptime"
)

// Profiler processes uptime information, /proc/uptime, using JSON.
type Profiler struct {
	*u.Profiler
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	p, err := u.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p}, nil
}

// Get gets the current uptime, /proc/uptime, as JSON serialized bytes.
func (prof *Profiler) Get() (p []byte, err error) {
	k, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(k)
}

var std *Profiler
var stdMu sync.Mutex //protects standard to prevent a data race on checking/instantiation

// Get gets the current uptime, /proc/uptime, as JSON serialized bytes using
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

// Serialize uptime.Uptime as JSON.
func (prof *Profiler) Serialize(up u.Uptime) ([]byte, error) {
	return json.Marshal(up)
}

// Serialize uptime.Uptime as JSON using the package's global Profiler.
func Serialize(up u.Uptime) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(up)
}

// Marshal is an alias for Serialize.
func (prof *Profiler) Marshal(up u.Uptime) ([]byte, error) {
	return prof.Serialize(up)
}

// Marshal is an alias for Serialize that uses the package's global profiler.
func Marshal(up u.Uptime) ([]byte, error) {
	return Serialize(up)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// uptime.Uptime.
func Deserialize(p []byte) (up u.Uptime, err error) {
	err = json.Unmarshal(p, &up)
	if err != nil {
		return up, err
	}
	return up, nil
}

// Unmarshal is an alias for Deserialize.
func Unmarshal(p []byte) (up u.Uptime, err error) {
	return Deserialize(p)
}

// Ticker delivers the system's uptime at intervals.
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

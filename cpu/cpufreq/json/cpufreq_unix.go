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

// Package cpufreq provides the current CPU frequency, in MHz, as reported by
// /proc/cpuinfo. Instead of returning a Go struct, it returns JSON serialized
// bytes. A function to deserialize the JSON serialized bytes into a
// cpufreq.Frequency struct is provided.
//
// Note: the package name is cpufreq and not the final element of the import
// path (json). 
package cpufreq

import (
	"encoding/json"
	"sync"
	"time"

	joe "github.com/rai-project/joefriday"
	freq "github.com/rai-project/joefriday/cpu/cpufreq"
)

// Profiler is used to process the frequency information as JSON serialized
// bytes.
type Profiler struct {
	*freq.Profiler
}

// Initializes and returns a cpufreq profiler.
func NewProfiler() (prof *Profiler, err error) {
	p, err := freq.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p}, nil
}

// Get returns the frequency as JSON serialized bytes.
func (prof *Profiler) Get() (p []byte, err error) {
	f, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(f)
}

var std *Profiler
var stdMu sync.Mutex //protects standard to prevent data race on checking/instantiation

// Get returns the frequency as JSON serialized bytes using the package's
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

// Serialize Frequency as JSON.
func (prof *Profiler) Serialize(f *freq.Frequency) ([]byte, error) {
	return json.Marshal(f)
}

// Serialize Frequency as JSON using package globals.
func Serialize(f *freq.Frequency) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(f)
}

// Marshal is an alias for serialize.
func (prof *Profiler) Marshal(f *freq.Frequency) ([]byte, error) {
	return prof.Serialize(f)
}

// Marshal is an alias for Serialize using package globals.
func Marshal(f *freq.Frequency) ([]byte, error) {
	return std.Serialize(f)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// cpufreq.Frequency.
func Deserialize(p []byte) (*freq.Frequency, error) {
	f := &freq.Frequency{}
	err := json.Unmarshal(p, f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Unmarshal is an alias for Deserialize using package globals.
func Unmarshal(p []byte) (*freq.Frequency, error) {
	return Deserialize(p)
}

// Ticker delivers the CPU Frequencies at intervals.
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
	// ticker
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

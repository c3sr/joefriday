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
// Instead of returning a Go struct, it returns Flatbuffer serialized bytes. A
// function to deserialize the Flatbuffer serialized bytes into a
// loadavg.LoadAvg struct is provided.
//
// Note: the package name is loadavg and not the final element of the import
// path (flat). 
package loadavg

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/rai-project/joefriday"
	l "github.com/rai-project/joefriday/system/loadavg"
	"github.com/rai-project/joefriday/system/loadavg/flat/structs"
)

// Profiler is used to process the loadavg information, /proc/loadavg, using
// Flatbuffers.
type Profiler struct {
	*l.Profiler
	*fb.Builder
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	p, err := l.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current loadavg information as Flatbuffer serialized bytes.
func (prof *Profiler) Get() ([]byte, error) {
	inf, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf), nil
}

var std *Profiler
var stdMu sync.Mutex //protects standard to prevent a data race on checking/instantiation

// Get returns the current loadavg information as Flatbuffer serialized bytes
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

// Serialize serializes loadavg information using Flatbuffers.
func (prof *Profiler) Serialize(la l.LoadAvg) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	structs.LoadAvgStart(prof.Builder)
	structs.LoadAvgAddTimestamp(prof.Builder, la.Timestamp)
	structs.LoadAvgAddMinute(prof.Builder, la.Minute)
	structs.LoadAvgAddFive(prof.Builder, la.Five)
	structs.LoadAvgAddFifteen(prof.Builder, la.Fifteen)
	structs.LoadAvgAddRunning(prof.Builder, la.Running)
	structs.LoadAvgAddTotal(prof.Builder, la.Total)
	structs.LoadAvgAddPID(prof.Builder, la.PID)
	prof.Builder.Finish(structs.LoadAvgEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize serializes LoadAvg information using Flatbuffers using the
// package's global Profiler.
func Serialize(la l.LoadAvg) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(la), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserializes them as
// loadavg.LoadAvg.
func Deserialize(p []byte) l.LoadAvg {
	flatLA := structs.GetRootAsLoadAvg(p, 0)
	var la l.LoadAvg
	la.Timestamp = flatLA.Timestamp()
	la.Minute = flatLA.Minute()
	la.Five = flatLA.Five()
	la.Fifteen = flatLA.Fifteen()
	la.Running = flatLA.Running()
	la.Total = flatLA.Total()
	la.PID = flatLA.PID()
	return la
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

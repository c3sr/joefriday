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

// Package cpuutil handles JSON based processing of CPU (kernel) utilization
// information. This information is calculated using the difference between
// two CPU (kernel) stats snapshots, /proc/stat, and represented as a
// percentage. The time elapsed between the two snapshots is stored in the
// TimeDelta field. Instead of returning a Go struct, it returns JSON
// serialized bytes. For convenience, a function to deserialize the JSON
// serialized bytes into a cpuutil.CPUUtil struct is provided.
//
// Note: the package name is cpuutil and not the final element of the import
// path (json). 
package cpuutil

import (
	"encoding/json"
	"sync"
	"time"

	joe "github.com/rai-project/joefriday"
	util "github.com/rai-project/joefriday/cpu/cpuutil"
)

// Profiler is used to process the /proc/stats file and calculate utilization
// information, returning the data as JSON serialized bytes.
type Profiler struct {
	*util.Profiler
}

// Initializes and returns a cpu utlization profiler.
func NewProfiler() (prof *Profiler, err error) {
	p, err := util.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p}, nil
}

// Get returns the cpu utilization as JSON serialized bytes. Utilization
// calculations requires two snapshots. This func gets the current snapshot of
// /proc/stat and calculates the utilization using the difference between the
// current snapshot and the prior one. The current snapshot is stored and for
// use as the prior snapshot on the next Get call. If ongoing utilitzation
// information is desired, the Ticker should be used; it's better suited for
// ongoing utilization information.
func (prof *Profiler) Get() (p []byte, err error) {
	st, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(st)
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current cpu utilization as JSON serialized bytes using the
// package's global Profiler. The Profiler is instantiated lazily. If the
// profiler doesn't already exist, the first usage information will not be
// useful due to minimal time elapsing between the initial and second snapshots
// used for usage calculations; the results of the first call should be
// discarded.
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

// Serialize cpu Utilization using JSON.
func (prof *Profiler) Serialize(ut *util.CPUUtil) ([]byte, error) {
	return json.Marshal(ut)
}

// Serialize the CPU Utilization as JSON using the package global Profiler.
func Serialize(ut *util.CPUUtil) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(ut)
}

// Marshal is an alias for Serialize
func (prof *Profiler) Marshal(ut *util.CPUUtil) ([]byte, error) {
	return prof.Serialize(ut)
}

// Marsha is an alias for Serialize using the package global Profiler.
func Marshal(ut *util.CPUUtil) ([]byte, error) {
	return Serialize(ut)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// cpuutil.CPUUtil.
func Deserialize(p []byte) (*util.CPUUtil, error) {
	ut := &util.CPUUtil{}
	err := json.Unmarshal(p, ut)
	if err != nil {
		return nil, err
	}
	return ut, nil
}

// Unmarshal is an alias for Deserialize.
func Unmarshal(p []byte) (*util.CPUUtil, error) {
	return Deserialize(p)
}

// Ticker delivers the system's CPU utilization information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan []byte
	*Profiler
}

// NewTicker returns a new Ticker containing a Data channel that delivers
// the data at intervals and an error channel that delivers any errors
// encountered. Stop the ticker to signal the ticker to stop running. Stopping
// the ticker does not close the Data channel; call Close to close both the
// ticker and the data channel.
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

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

// Package netdev handles JSON based processing of network device information.
// Instead of returning a Go struct, it returns JSON serialized bytes. A
// function to deserialize the JSON serialized bytes into a structs.DevInfo
// struct is provided.
//
// Note: the package name is netdev and not the final element of the import
// path (flat). 
package netdev

import (
	"encoding/json"
	"sync"
	"time"

	joe "github.com/mohae/joefriday"
	dev "github.com/mohae/joefriday/net/netdev"
	"github.com/mohae/joefriday/net/structs"
)

// Profiler is used to process the /proc/net/dev file using JSON.
type Profiler struct {
	*dev.Profiler
}

// Initializes and returns a network devices profiler.
func NewProfiler() (prof *Profiler, err error) {
	p, err := dev.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p}, nil
}

// Get returns the current network devices as JSON serialized bytes.
func (prof *Profiler) Get() (p []byte, err error) {
	inf, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf)
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current network devices information as JSON serialized bytes
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

// Serialize network devicese information using JSON.
func (prof *Profiler) Serialize(inf *structs.DevInfo) ([]byte, error) {
	return json.Marshal(inf)
}

// Serialize network devices information using JSON with the package global
// Profiler.
func Serialize(inf *structs.DevInfo) (p []byte, err error) {
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
func (prof *Profiler) Marshal(inf *structs.DevInfo) ([]byte, error) {
	return prof.Serialize(inf)
}

// Marshal is an alias for Serialize; uses the package global Profiler.
func Marshal(inf *structs.DevInfo) ([]byte, error) {
	return Serialize(inf)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// structs.Info
func Deserialize(p []byte) (*structs.DevInfo, error) {
	info := &structs.DevInfo{}
	err := json.Unmarshal(p, info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// Unmarshal is an alias for Deserialize.
func Unmarshal(p []byte) (*structs.DevInfo, error) {
	return Deserialize(p)
}

// Ticker delivers the system's net devices information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan []byte
	*Profiler
}

// NewTicker returns a new Ticker contianing a Data channel that delivers
// the data at intervals and an error channel that delivers any errors
// encountered.  Stop the ticker to signal the ticker to stop running; it
// does not close the Data channel.  Close the ticker to close all ticker
// channels.
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

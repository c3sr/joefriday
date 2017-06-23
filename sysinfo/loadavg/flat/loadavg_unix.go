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

// Package loadavg handles Flatbuffer based processing of load using syscall.
// Instead of returning a Go struct, it returns Flatbuffer serialized bytes.
// A function to deserialize the Flatbuffer serialized bytes into a
// loadavg.Info struct is provided. After the first use, the flatbuffer
// builder is reused.
//
// Note: the loadavg name is processors and not the final element of the import
// path (flat). 
package loadavg

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	load "github.com/mohae/joefriday/sysinfo/loadavg"
	"github.com/mohae/joefriday/sysinfo/loadavg/flat/structs"
)

var builder = fb.NewBuilder(0)
var mu sync.Mutex

// Get returns the current load as Flatbuffer serialized bytes.
func Get() (p []byte, err error) {
	var l load.Info
	err = l.Get()
	if err != nil {
		return nil, err
	}
	return Serialize(&l), nil
}

// Serialize loadAvg.Info using Flatbuffers.
func Serialize(l *load.Info) []byte {
	mu.Lock()
	defer mu.Unlock()
	// ensure the Builder is in a usable state.
	builder.Reset()
	structs.InfoStart(builder)
	structs.InfoAddTimestamp(builder, l.Timestamp)
	structs.InfoAddOne(builder, l.One)
	structs.InfoAddFive(builder, l.Five)
	structs.InfoAddFifteen(builder, l.Fifteen)
	builder.Finish(structs.InfoEnd(builder))
	p := builder.Bytes[builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them
// as loadavg.Info.
func Deserialize(p []byte) *load.Info {
	lF := structs.GetRootAsInfo(p, 0)
	l := &load.Info{}
	l.Timestamp = lF.Timestamp()
	l.One = lF.One()
	l.Five = lF.Five()
	l.Fifteen = lF.Fifteen()
	return l
}

// Ticker delivers loadavg.Info as Flatbuffers serialized bytes at
// intervals.
type Ticker struct {
	*joe.Ticker
	Data chan []byte
}

// NewTicker returns a new Ticker continaing a Data channel that delivers
// the data at intervals and an error channel that delivers any errors
// encountered.  Stop the ticker to signal the ticker to stop running; it
// does not close the Data channel.  Close the ticker to close all ticker
// channels.
func NewTicker(d time.Duration) (joe.Tocker, error) {
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan []byte)}
	go t.Run()
	return &t, nil
}

// Run runs the ticker.
func (t *Ticker) Run() {
	// read until done signal is received
	for {
		select {
		case <-t.Done:
			return
		case <-t.Ticker.C:
			s, err := Get()
			if err != nil {
				t.Errs <- err
				continue
			}
			t.Data <- s
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}

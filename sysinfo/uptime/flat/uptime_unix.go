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

// Package uptime gets the system's uptime using syscall. Instead of returning
// a Go struct, it returns Flatbuffer serialized bytes. A function to
// deserialize the Flatbuffer serialized bytes into an uptime.Info struct is
// provided.
//
// Note: the package name is uptime and not the final element of the import
// path (flat). 
package uptime

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/c3sr/joefriday"
	up "github.com/c3sr/joefriday/sysinfo/uptime"
	"github.com/c3sr/joefriday/sysinfo/uptime/flat/structs"
)

var builder = fb.NewBuilder(0)
var mu sync.Mutex

// Get returns the current uptime as Flatbuffer serialized bytes.
func Get() (p []byte, err error) {
	var u up.Uptime
	err = u.Get()
	if err != nil {
		return nil, err
	}
	return Serialize(&u), nil
}

// Serialize uptime.Uptime using Flatbuffers.
func Serialize(u *up.Uptime) []byte {
	mu.Lock()
	defer mu.Unlock()
	// ensure the Builder is in a usable state.
	builder.Reset()
	structs.UptimeStart(builder)
	structs.UptimeAddTimestamp(builder, u.Timestamp)
	structs.UptimeAddUptime(builder, u.Uptime)
	builder.Finish(structs.UptimeEnd(builder))
	p := builder.Bytes[builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Deserialize takes some Flatbuffer serialized bytes and deserializes them
// as uptime.Uptime.
func Deserialize(p []byte) *up.Uptime {
	uF := structs.GetRootAsUptime(p, 0)
	var u up.Uptime
	u.Timestamp = uF.Timestamp()
	u.Uptime = uF.Uptime()
	return &u
}

// Ticker delivers the uptime as Flatbuffers serialized bytes at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan []byte
}

// NewTicker returns a new Ticker containing a Data channel that delivers the
// data at intervals and an error channel that delivers any errors encountered.
// Stop the ticker to signal the ticker to stop running. Stopping the ticker
// does not close the Data channel; call Close to close both the ticker and the
// data channel.
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
			p, err := Get()
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

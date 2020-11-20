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
// a Go struct, it returns JSON serialized bytes. A function to deserialize the
// JSON serialized bytes into a uptime.Uptime struct is provided.
//
// Note: the package name is uptime and not the final element of the import
// path (json). 
package uptime

import (
	"encoding/json"
	"time"

	joe "github.com/c3sr/joefriday"
	up "github.com/c3sr/joefriday/sysinfo/uptime"
)

// Get returns the current uptime as JSON serialized bytes.
func Get() (p []byte, err error) {
	var u up.Uptime
	err = u.Get()
	if err != nil {
		return nil, err
	}
	return json.Marshal(&u)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// uptime.Uptime.
func Deserialize(p []byte) (*up.Uptime, error) {
	var u up.Uptime
	err := json.Unmarshal(p, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Unmarshal is an alias for Deserialize.
func Unmarshal(p []byte) (*up.Uptime, error) {
	return Deserialize(p)
}

// Ticker delivers uptime.Uptime as JSON serialized bytes at intervals.
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

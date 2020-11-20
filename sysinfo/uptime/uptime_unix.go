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

// Package uptime gets the system's uptime using syscall.
package uptime

import (
	"syscall"
	"time"

	joe "github.com/c3sr/joefriday"
)

// Uptime holds the current uptime and timestamp.
type Uptime struct {
	Timestamp int64
	Uptime    int64 // sorry for the stutter
}

// Get gets the current uptime information.
func (u *Uptime) Get() error {
	var sysinfo syscall.Sysinfo_t
	err := syscall.Sysinfo(&sysinfo)
	if err != nil {
		return err
	}
	u.Timestamp = time.Now().UTC().UnixNano()
	u.Uptime = sysinfo.Uptime
	return nil
}

// Get gets the current uptime information.
func Get() (u Uptime, err error) {
	err = u.Get()
	return u, err
}

// Ticker deliivers the uptime at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan Uptime
}

// NewTicker returns a new Ticker containing a Data channel that delivers the
// data at intervals and an error channel that delivers any errors encountered.
// Stop the ticker to signal the ticker to stop running. Stopping the ticker
// does not close the Data channel; call Close to close both the ticker and the
// data channel.
func NewTicker(d time.Duration) (joe.Tocker, error) {
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan Uptime)}
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
			u, err := Get()
			if err != nil {
				t.Errs <- err
				continue
			}
			t.Data <- u
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}

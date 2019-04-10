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

// Package loadavg provides the system's loadavg information using a syscall.
package loadavg

import (
	"syscall"
	"time"

	joe "github.com/rai-project/joefriday"
)

const LoadsScale = 65536

// LoadAvg holds loadavg information and the timestamp of the data.
type LoadAvg struct {
	Timestamp int64
	One       float64
	Five      float64
	Fifteen   float64
}

// Get the load average for the last 1, 5, and 15 minutes.
func (l *LoadAvg) Get() error {
	var sysinfo syscall.Sysinfo_t
	err := syscall.Sysinfo(&sysinfo)
	if err != nil {
		return err
	}
	l.Timestamp = time.Now().UTC().UnixNano()
	l.One = float64(sysinfo.Loads[0]) / LoadsScale
	l.Five = float64(sysinfo.Loads[1]) / LoadsScale
	l.Fifteen = float64(sysinfo.Loads[2]) / LoadsScale
	return nil
}

// Get returns LoadAvg populated with the 1, 5, and 15 minute values.
func Get() (LoadAvg, error) {
	var l LoadAvg
	err := l.Get()
	return l, err
}

// Ticker delivers the LoadAvg at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan LoadAvg
}

// NewTicker returns a new Ticker containing a Data channel that delivers the
// data at intervals and an error channel that delivers any errors encountered.
// Stop the ticker to signal the ticker to stop running. Stopping the ticker
// does not close the Data channel; call Close to close both the ticker and the
// data channel.
func NewTicker(d time.Duration) (joe.Tocker, error) {
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan LoadAvg)}
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

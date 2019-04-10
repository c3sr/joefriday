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

// Package mem returns memory information using syscalls. Only basic memory
// information is provided by this package.
package mem

import (
	"syscall"
	"time"

	joe "github.com/rai-project/joefriday"
)

// MemInfo holds information about system memory.
type MemInfo struct {
	Timestamp int64
	TotalRAM  uint64
	FreeRAM   uint64
	SharedRAM uint64
	BufferRAM uint64
	TotalSwap uint64
	FreeSwap  uint64
}

// Get gets the system's memory information.
func (m *MemInfo) Get() error {
	var sysinfo syscall.Sysinfo_t
	err := syscall.Sysinfo(&sysinfo)
	if err != nil {
		return err
	}
	m.Timestamp = time.Now().UTC().UnixNano()
	m.TotalRAM = sysinfo.Totalram
	m.FreeRAM = sysinfo.Freeram
	m.SharedRAM = sysinfo.Sharedram
	m.BufferRAM = sysinfo.Bufferram
	m.TotalSwap = sysinfo.Totalswap
	m.FreeSwap = sysinfo.Freeswap
	return nil
}

// Get gets the system's memory information.
func Get() (m MemInfo, err error) {
	err = m.Get()
	return m, err
}

// Ticker delivers the system's memory information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan MemInfo
}

// NewTicker returns a new Ticker containing a Data channel that delivers the
// data at intervals and an error channel that delivers any errors encountered.
// Stop the ticker to signal the ticker to stop running. Stopping the ticker
// does not close the Data channel; call Close to close both the ticker and the
// data channel.
func NewTicker(d time.Duration) (joe.Tocker, error) {
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan MemInfo)}
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

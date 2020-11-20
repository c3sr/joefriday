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

// Package netdev gets the system's network device information: /proc/net/dev.
package netdev

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/c3sr/joefriday/helpers"
	joe "github.com/c3sr/joefriday"
	"github.com/c3sr/joefriday/net/structs"
)

// ProcFile is the file used by the netdev Profiler.
const ProcFile = "/proc/net/dev"

// Profiler is used to process the network device information using the
// /proc/net/dev file.
type Profiler struct {
	joe.Procer
	*joe.Buffer
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	proc, err := joe.NewProc(ProcFile)
	if err != nil {
		return nil, err
	}
	return &Profiler{Procer: proc, Buffer: joe.NewBuffer()}, nil
}

// Reset resources: after reset, the profiler is ready to be used again.
func (prof *Profiler) Reset() error {
	prof.Buffer.Reset()
	return prof.Procer.Reset()
}

// Get returns the current network device information.
func (prof *Profiler) Get() (*structs.DevInfo, error) {
	var (
		i, pos, line, fieldNum int
		n                      uint64
		v                      byte
		dev                  structs.Device
	)
	err := prof.Reset()
	if err != nil {
		return nil, err
	}
	// there's, usually, at least 2 devices
	nDev := &structs.DevInfo{Timestamp: time.Now().UTC().UnixNano(), Device: make([]structs.Device, 0, 2)}
	for {
		prof.Line, err = prof.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, &joe.ReadError{Err: err}
		}
		line++
		if line < 3 {
			continue
		}
		prof.Val = prof.Val[:0]
		// first grab the interface name (everything up to the ':')
		for i, v = range prof.Line {
			if v == 0x3A {
				pos = i + 1
				break
			}
			// skip spaces
			if v == 0x20 {
				continue
			}
			prof.Val = append(prof.Val, v)
		}
		dev.Name = string(prof.Val[:])
		fieldNum = 0
		// process the rest of the line
		for {
			fieldNum++
			// skip all spaces
			for i, v = range prof.Line[pos:] {
				if v != 0x20 {
					pos += i
					break
				}
			}

			// grab the numbers
			for i, v = range prof.Line[pos:] {
				if v == 0x20 || v == '\n' {
					break
				}
			}
			// any conversion error results in 0
			n, err = helpers.ParseUint(prof.Line[pos : pos+i])
			pos += i
			if err != nil {
				return nil, &joe.ParseError{Info: fmt.Sprintf("line %d: field %d", line, fieldNum), Err: err}
			}
			if fieldNum < 9 {
				if fieldNum < 5 {
					if fieldNum < 3 {
						if fieldNum == 1 {
							dev.RBytes = int64(n)
							continue
						}
						dev.RPackets = int64(n) // must be 2
						continue
					}
					if fieldNum == 3 {
						dev.RErrs = int64(n)
						continue
					}
					dev.RDrop = int64(n) // must be 4
					continue
				}
				if fieldNum < 7 {
					if fieldNum == 5 {
						dev.RFIFO = int64(n)
						continue
					}
					dev.RFrame = int64(n) // must be 6
					continue
				}
				if fieldNum == 7 {
					dev.RCompressed = int64(n)
					continue
				}
				dev.RMulticast = int64(n) // must be 8
				continue
			}
			if fieldNum < 13 {
				if fieldNum < 11 {
					if fieldNum == 9 {
						dev.TBytes = int64(n)
						continue
					}
					dev.TPackets = int64(n) // must be 10
					continue
				}
				if fieldNum == 11 {
					dev.TErrs = int64(n)
					continue
				}
				dev.TDrop = int64(n) // must be 12
				continue
			}
			if fieldNum < 15 {
				if fieldNum == 13 {
					dev.TFIFO = int64(n)
					continue
				}
				dev.TColls = int64(n) // must be 14
				continue
			}
			if fieldNum == 15 {
				dev.TCarrier = int64(n)
				continue
			}
			dev.TCompressed = int64(n)
			break
		}
		nDev.Device = append(nDev.Device, dev)
	}
	return nDev, nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current network device information using the package's
// global Profiler.
func Get() (inf *structs.DevInfo, err error) {
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

// Ticker delivers the system's network device information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan *structs.DevInfo
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
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan *structs.DevInfo), Profiler: p}
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
			s, err := t.Get()
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

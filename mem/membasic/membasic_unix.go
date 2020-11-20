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

// Package membasic processes a subset of the /proc/meminfo file. For more
// detailed information about a system's memory, use the meminfo package.
package membasic

import (
	"io"
	"sync"
	"time"

	"github.com/c3sr/joefriday/helpers"
	joe "github.com/c3sr/joefriday"
)

const procFile = "/proc/meminfo"

// Info holds the basic meminfo information.
type Info struct {
	Timestamp    int64  `json:"timestamp"`
	Active       uint64 `json:"active"`
	Inactive     uint64 `json:"inactive"`
	Mapped       uint64 `json:"mapped"`
	MemAvailable uint64 `json:"mem_available"`
	MemFree      uint64 `json:"mem_free"`
	MemTotal     uint64 `json:"mem_total"`
	SwapCached   uint64 `json:"swap_cached"`
	SwapFree     uint64 `json:"swap_free"`
	SwapTotal    uint64 `json:"swap_total"`
}

// Profiler is used to get the basic memory information by processing the
// /proc/meminfo file.
type Profiler struct {
	joe.Procer
	*joe.Buffer
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	proc, err := joe.NewProc(procFile)
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

// Get returns the current basic memory information.
func (prof *Profiler) Get() (inf *Info, err error) {
	var (
		i, pos, nameLen int
		v               byte
		n               uint64
	)
	err = prof.Reset()
	if err != nil {
		return nil, err
	}
	inf = &Info{}
	inf.Timestamp = time.Now().UTC().UnixNano()
	for {
		prof.Val = prof.Val[:0]
		prof.Line, err = prof.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return inf, &joe.ReadError{Err: err}
		}
		// first grab the key name (everything up to the ':')
		for i, v = range prof.Line {
			if v == ':' {
				pos = i + 1
				break
			}
			prof.Val = append(prof.Val, v)
		}
		nameLen = len(prof.Val)

		// skip all spaces
		for i, v = range prof.Line[pos:] {
			if v != ' ' {
				pos += i
				break
			}
		}

		// grab the numbers
		for _, v = range prof.Line[pos:] {
			if v == ' ' || v == '\n' {
				break
			}
			prof.Val = append(prof.Val, v)
		}
		// any conversion error results in 0
		n, err = helpers.ParseUint(prof.Val[nameLen:])
		if err != nil {
			return inf, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
		}

		v = prof.Val[0]
		// evaluate the key
		if v == 'A' {
			if prof.Val[5] == 'e' && nameLen == 6 {
				inf.Active = n
			}
			continue
		}
		if v == 'I' {
			if nameLen == 8 {
				inf.Inactive = n
			}
			continue
		}
		if v == 'M' {
			v = prof.Val[3]
			if nameLen < 8 {
				if v == 'p' {
					inf.Mapped = n
					continue
				}
				if v == 'F' {
					inf.MemFree = n
				}
				continue
			}
			if v == 'A' {
				inf.MemAvailable = n
				continue
			}
			inf.MemTotal = n
			continue
		}
		if v == 'S' {
			v = prof.Val[1]
			if v == 'w' {
				if prof.Val[4] == 'C' {
					inf.SwapCached = n
					continue
				}
				if prof.Val[4] == 'F' {
					inf.SwapFree = n
					continue
				}
				inf.SwapTotal = n
				continue
			}
		}
	}
	return inf, nil
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current basic memory information using the package's global
// Profiler.
func Get() (inf *Info, err error) {
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

// Ticker delivers the system's basic memory information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan Info
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
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan Info), Profiler: p}
	go t.Run()
	return &t, nil
}

// Run runs the ticker.
func (t *Ticker) Run() {
	// predeclare some vars
	var (
		i, pos, nameLen int
		v               byte
		n               uint64
		err             error
		inf             Info
	)
	// ticker
	for {
		select {
		case <-t.Done:
			return
		case <-t.C:
			err = t.Profiler.Reset()
			if err != nil {
				t.Errs <- err
				continue
			}
			inf.Timestamp = time.Now().UTC().UnixNano()
			for {
				t.Val = t.Val[:0]
				t.Line, err = t.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					t.Errs <- &joe.ReadError{Err: err}
				}
				// first grab the key name (everything up to the ':')
				for i, v = range t.Line {
					if v == ':' {
						pos = i + 1
						break
					}
					t.Val = append(t.Val, v)
				}
				nameLen = len(t.Val)

				// skip all spaces
				for i, v = range t.Line[pos:] {
					if v != ' ' {
						pos += i
						break
					}
				}

				// grab the numbers
				for _, v = range t.Line[pos:] {
					if v == ' ' || v == '\n' {
						break
					}
					t.Val = append(t.Val, v)
				}
				// any conversion error results in 0
				n, err = helpers.ParseUint(t.Val[nameLen:])
				if err != nil {
					t.Errs <- &joe.ParseError{Info: string(t.Val[:nameLen]), Err: err}
				}

				v = t.Val[0]
				// evaluate the key
				if v == 'A' {
					if t.Val[5] == 'e' && nameLen == 6 {
						inf.Active = n
					}
					continue
				}
				if v == 'I' {
					if nameLen == 8 {
						inf.Inactive = n
					}
					continue
				}
				if v == 'M' {
					v = t.Val[3]
					if nameLen < 8 {
						if v == 'p' {
							inf.Mapped = n
							continue
						}
						if v == 'F' {
							inf.MemFree = n
						}
						continue
					}
					if v == 'A' {
						inf.MemAvailable = n
						continue
					}
					inf.MemTotal = n
					continue
				}
				if v == 'S' {
					v = t.Val[1]
					if v == 'w' {
						if t.Val[4] == 'C' {
							inf.SwapCached = n
							continue
						}
						if t.Val[4] == 'F' {
							inf.SwapFree = n
							continue
						}
						inf.SwapTotal = n
						continue
					}
				}
			}
			t.Data <- inf
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}

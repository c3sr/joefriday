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

// Package loadAvg gets loadavg information from the /proc/loadavg file.
package loadavg

import (
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/c3sr/joefriday/helpers"
	joe "github.com/c3sr/joefriday"
)

const procFile = "/proc/loadavg"

// LoadAvg holds loadavg information
type LoadAvg struct {
	Timestamp int64
	Minute    float32
	Five      float32
	Fifteen   float32
	Running   int32
	Total     int32
	PID       int32
}

// Profiler processes the loadavg information.
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

// Get populates LoadAvg with /proc/loadavg information.
func (prof *Profiler) Get() (la LoadAvg, err error) {
	err = prof.Reset()
	if err != nil {
		return la, err
	}
	var (
		i, priorPos, pos, line, fieldNum int
		n                                uint64
		f                                float64
		v                                byte
	)
	la.Timestamp = time.Now().UTC().UnixNano()
	for {
		prof.Line, err = prof.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return la, &joe.ReadError{Err: err}
		}
		line++
		for {
			// space delimits the values
			for i, v = range prof.Line[pos:] {
				if v == 0x20 {
					priorPos, pos = pos, pos+i+1
					break
				}
			}
			fieldNum++
			if fieldNum <= 3 {
				f, err = strconv.ParseFloat(string(prof.Line[priorPos:pos-1]), 64)
				if err != nil {
					return la, &joe.ParseError{Info: fmt.Sprintf("line %d: field %d", line, fieldNum), Err: err}
				}
				if fieldNum == 1 {
					la.Minute = float32(f)
					continue
				}
				if fieldNum == 2 {
					la.Five = float32(f)
					continue
				}
				la.Fifteen = float32(f)
				continue
			}
			if fieldNum == 4 {
				// get the process information: separated by /
				for i, v = range prof.Line[priorPos:pos] {
					if v == '/' {
						break
					}
				}
				n, err = helpers.ParseUint(prof.Line[priorPos : priorPos+i])
				if err != nil {
					return la, &joe.ParseError{Info: fmt.Sprintf("line %d: field %d", line, fieldNum), Err: err}
				}
				la.Running = int32(n)
				n, err = helpers.ParseUint(prof.Line[priorPos+i+1 : pos-1])
				if err != nil {
					return la, &joe.ParseError{Info: fmt.Sprintf("line %d: field %d", line, fieldNum), Err: err}
				}
				la.Total = int32(n)
				continue
			}
			n, err = helpers.ParseUint(prof.Line[pos : len(prof.Line)-1])
			if err != nil {
				return la, err
			}
			la.PID = int32(n)
			break
		}
	}
	return la, nil
}

var std *Profiler
var stdMu sync.Mutex

// Get gets the loadavg information using the package's global Profiler.
func Get() (la LoadAvg, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return la, err
		}
	}
	return std.Get()
}

// Ticker delivers the system's loadavg information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan LoadAvg
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
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan LoadAvg), Profiler: p}
	go t.Run()
	return &t, nil
}

// Run runs the ticker.
func (t *Ticker) Run() {
	var (
		i, priorPos, pos, line, fieldNum int
		n                                uint64
		f                                float64
		v                                byte
		la                              LoadAvg
		err                              error
	)
	// ticker
	for {
		select {
		case <-t.Done:
			return
		case <-t.C:
			la.Timestamp = time.Now().UTC().UnixNano()
			err = t.Reset()
			line = 0
		tick:
			for {
				t.Line, err = t.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					t.Errs <- &joe.ReadError{Err: err}
					break
				}
				line++
				for {
					// space delimits the values
					for i, v = range t.Line[pos:] {
						if v == 0x20 {
							priorPos, pos = pos, pos+i+1
							break
						}
					}
					fieldNum++
					if fieldNum <= 3 {
						f, err = strconv.ParseFloat(string(t.Line[priorPos:pos-1]), 64)
						if err != nil {
							t.Errs <- &joe.ParseError{Info: fmt.Sprintf("line %d: field %d", line, fieldNum), Err: err}
							break tick
						}
						if fieldNum == 1 {
							la.Minute = float32(f)
							continue
						}
						if fieldNum == 2 {
							la.Five = float32(f)
							continue
						}
						la.Fifteen = float32(f)
						continue
					}
					if fieldNum == 4 {
						// get the process information: separated by /
						for i, v = range t.Line[priorPos:pos] {
							if v == '/' {
								break
							}
						}
						n, err = helpers.ParseUint(t.Line[priorPos : priorPos+i])
						if err != nil {
							t.Errs <- &joe.ParseError{Info: fmt.Sprintf("line %d: field %d", line, fieldNum), Err: err}
							break tick
						}
						la.Running = int32(n)
						n, err = helpers.ParseUint(t.Line[priorPos+i+1 : pos-1])
						if err != nil {
							t.Errs <- &joe.ParseError{Info: fmt.Sprintf("line %d: field %d", line, fieldNum), Err: err}
							break tick
						}
						la.Total = int32(n)
						continue
					}
					n, err = helpers.ParseUint(t.Line[pos : len(t.Line)-1])
					if err != nil {
						t.Errs <- &joe.ParseError{Info: fmt.Sprintf("line %d: field %d", line, fieldNum), Err: err}
						break tick
					}
					la.PID = int32(n)
					break
				}
			}
			t.Data <- la
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}

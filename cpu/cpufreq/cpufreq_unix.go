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

// Package cpufreq provides the current CPU frequency, in MHz, as reported by
// /proc/cpuinfo.
package cpufreq

import (
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/c3sr/joefriday/helpers"
	joe "github.com/c3sr/joefriday"
)

const procFile = "/proc/cpuinfo"

// Frequency holds information about the current frequency of a system's cpus,
// in MHz. The reported values are the current speeds as reported by
// /proc/cpuinfo.
type Frequency struct {
	Timestamp int64 `json:"timestamp"`
	Sockets   int32 `json:"sockets"`
	CPU       []CPU `json:"cpu"`
}

// CPU holds the clock info for a single processor.
type CPU struct {
	Processor  int32   `json:"processor"`
	CPUMHz     float32 `json:"cpu_mhz"`
	PhysicalID int32   `json:"physical_id"`
	CoreID     int32   `json:"core_id"`
	APICID     int32   `json:"apicid"`
}

// Profiler is used to process the frequency information.
type Profiler struct {
	joe.Procer
	*joe.Buffer
	Frequency // this is used too hold the socket/cpu info so that everything doesn't have to be reprocessed.
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	proc, err := joe.NewProc(procFile)
	if err != nil {
		return nil, err
	}
	prof = &Profiler{Procer: proc, Buffer: joe.NewBuffer()}
	err = prof.InitFrequency()
	if err != nil {
		return nil, err
	}
	return prof, nil
}

// Reset resources; after reset the profiler is ready to be used again.
func (prof *Profiler) Reset() error {
	prof.Buffer.Reset()
	return prof.Procer.Reset()
}

// InitFrequency sets the profiler's frequency with the static information so
// that everything doesn't need to be reprocessed every time the frequency is
// requested. This assumes that cpuinfo returns processor information in the
// same order every time.
//
// This shouldn't be used; it's exported for testing reasons.
func (prof *Profiler) InitFrequency() error {
	var (
		err         error
		n           uint64
		pos, cpuCnt int
		pidFound    bool
		physIDs     []int32 // tracks unique physical IDs encountered
		cpu         CPU
	)

	prof.Frequency = Frequency{}
	for {
		prof.Line, err = prof.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return &joe.ReadError{Err: err}
		}
		prof.Val = prof.Val[:0]
		// First grab the attribute name; everything up to the ':'.  The key may have
		// spaces and has trailing spaces; that gets trimmed.
		for i, v := range prof.Line {
			if v == 0x3A {
				prof.Val = prof.Line[:i]
				pos = i + 1
				break
			}
			//prof.Val = append(prof.Val, v)
		}
		prof.Val = joe.TrimTrailingSpaces(prof.Val[:])
		nameLen := len(prof.Val)
		// if there's no name; skip.
		if nameLen == 0 {
			continue
		}
		// if there's anything left, the value is everything else; trim spaces
		if pos+1 < len(prof.Line) {
			prof.Val = append(prof.Val, joe.TrimTrailingSpaces(prof.Line[pos+1:])...)
		}
		if prof.Val[0] == 'a' {
			if prof.Val[1] == 'p' { // apicid
				n, err = helpers.ParseUint(prof.Val[nameLen:])
				if err != nil {
					return &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}
				cpu.APICID = int32(n)
			}
			continue
		}
		if prof.Val[0] == 'c' {
			if prof.Val[1] == 'o' { // core id
				n, err = helpers.ParseUint(prof.Val[nameLen:])
				if err != nil {
					return &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}
				cpu.CoreID = int32(n)
			}
			continue
		}
		if prof.Val[0] == 'p' {
			if prof.Val[1] == 'h' { // physical id
				n, err = helpers.ParseUint(prof.Val[nameLen:])
				if err != nil {
					return &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}
				cpu.PhysicalID = int32(n)
				for i := range physIDs {
					if physIDs[i] == cpu.PhysicalID {
						pidFound = true
						break
					}
				}
				if pidFound {
					pidFound = false // reset for next use
				} else {
					// physical id hasn't been encountered yet; add it
					physIDs = append(physIDs, cpu.PhysicalID)
				}
				continue
			}
			// processor starts information about a processor.
			if prof.Val[1] == 'r' { // processor
				if cpuCnt > 0 {
					prof.Frequency.CPU = append(prof.Frequency.CPU, cpu)
				}
				cpuCnt++
				n, err = helpers.ParseUint(prof.Val[nameLen:])
				if err != nil {
					return &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}
				cpu = CPU{Processor: int32(n)}
			}
		}
		continue
	}
	// append the current processor informatin
	prof.Frequency.CPU = append(prof.Frequency.CPU, cpu)
	prof.Frequency.Sockets = int32(len(physIDs))
	return nil
}

// returns a copy of the profiler's frequency.
func (prof *Profiler) newFrequency() *Frequency {
	f := &Frequency{Timestamp: time.Now().UTC().UnixNano(), Sockets: prof.Frequency.Sockets, CPU: make([]CPU, len(prof.Frequency.CPU))}
	copy(f.CPU, prof.Frequency.CPU)
	return f
}

// Get returns Frequency information.
func (prof *Profiler) Get() (f *Frequency, err error) {
	f = prof.newFrequency()
	err = prof.Reset()
	if err != nil {
		return nil, err
	}

	var (
		i, pos, nameLen int
		v               byte
		x               float64
	)
	processor := -1 // start at -1 because it'll be incremented before use as it's the first line encountered
	for {
		prof.Line, err = prof.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, &joe.ReadError{Err: err}
		}
		prof.Val = prof.Val[:0]
		// First grab the attribute name; everything up to the ':'.  The key may have
		// spaces and has trailing spaces; that gets trimmed.
		for i, v = range prof.Line {
			if v == 0x3A {
				prof.Val = prof.Line[:i]
				pos = i + 1
				break
			}
			//prof.Val = append(prof.Val, v)
		}
		prof.Val = joe.TrimTrailingSpaces(prof.Val[:])
		nameLen = len(prof.Val)
		// if there's no name; skip.
		if nameLen == 0 {
			continue
		}
		// if there's anything left, the value is everything else; trim spaces
		if pos+1 < len(prof.Line) {
			prof.Val = append(prof.Val, joe.TrimTrailingSpaces(prof.Line[pos+1:])...)
		}
		if prof.Val[0] == 'c' {
			if prof.Val[4] == 'M' { // cpu MHz
				x, err = strconv.ParseFloat(string(prof.Val[nameLen:]), 32)
				if err != nil {
					return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}
				f.CPU[processor].CPUMHz = float32(x)
			}
			continue
		}
		if prof.Val[0] == 'p' {
			// processor starts information about a processor.
			if prof.Val[1] == 'r' { // processor
				processor++
			}
		}
	}
	return f, nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns Frequency using the package's global Profiler.
func Get() (f *Frequency, err error) {
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

// Ticker delivers the CPU Frequencies at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan *Frequency
	*Profiler
	Sockets uint8
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
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan *Frequency), Profiler: p}
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

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

// Package cpustats handles the processing of information about kernel activity,
// /proc/stat. The first CPUStats.CPU element aggregates the values for all
// other CPU elements. The values are aggregated since system boot.
package cpustats

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/c3sr/joefriday/helpers"
	joe "github.com/c3sr/joefriday"
)

const procFile = "/proc/stat"

var CLK_TCK int32    // the ticks per clock cycle
var tckMu sync.Mutex //protects CLK_TCK

// Set CLK_TCK.
func ClkTck() error {
	tckMu.Lock()
	defer tckMu.Unlock()
	var out bytes.Buffer
	cmd := exec.Command("getconf", "CLK_TCK")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("get CLK_TCK: %s", err)
	}
	b, err := out.ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("err CLK_TCK: %s", err)
	}
	v, err := strconv.Atoi(string(b[:len(b)-1]))
	if err != nil {
		return fmt.Errorf("process CLK_TCK output: %s", err)
	}
	atomic.StoreInt32(&CLK_TCK, int32(v))
	return nil
}

// CPUStats holds the kernel activity information; /proc/stat. The first CPU
// element's values are the aggregates of all other CPU elements. The stats are
// aggregated from sytem boot.
type CPUStats struct {
	ClkTck    int16  `json:"clk_tck"`
	Timestamp int64  `json:"timestamp"`
	Ctxt      int64  `json:"ctxt"`
	BTime     int64  `json:"btime"`
	Processes int64  `json:"processes"`
	CPU       []CPU `json:"cpu"`
}

// CPU holds the stats for a single CPU entry in the /proc/stat file.
type CPU struct {
	ID        string `json:"ID"`
	User      int64  `json:"user"`
	Nice      int64  `json:"nice"`
	System    int64  `json:"system"`
	Idle      int64  `json:"idle"`
	IOWait    int64  `json:"io_wait"`
	IRQ       int64  `json:"irq"`
	SoftIRQ   int64  `json:"soft_irq"`
	Steal     int64  `json:"steal"`
	Quest     int64  `json:"quest"`
	QuestNice int64  `json:"quest_nice"`
}

// Profiler is used to process the /proc/stats file.
type Profiler struct {
	joe.Procer
	*joe.Buffer
	ClkTck int16
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	// if it hasn't been set, set it.
	if atomic.LoadInt32(&CLK_TCK) == 0 {
		err = ClkTck()
		if err != nil {
			return nil, err
		}
	}
	proc, err := joe.NewProc(procFile)
	if err != nil {
		return nil, err
	}
	return &Profiler{Procer: proc, Buffer: joe.NewBuffer(), ClkTck: int16(atomic.LoadInt32(&CLK_TCK))}, nil
}

// Reset resources: after reset, the profiler is ready to be used again.
func (prof *Profiler) Reset() error {
	prof.Buffer.Reset()
	return prof.Procer.Reset()
}

// Get returns information about current kernel activity.
func (prof *Profiler) Get() (stats *CPUStats, err error) {
	err = prof.Reset()
	if err != nil {
		return nil, err
	}
	var (
		i, j, pos, fieldNum int
		n                   uint64
		v                   byte
		stop                bool
	)

	stats = &CPUStats{Timestamp: time.Now().UTC().UnixNano(), ClkTck: prof.ClkTck, CPU: make([]CPU, 0, 2)}

	// read each line until eof
	for {
		prof.Line, err = prof.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, &joe.ReadError{Err: err}
		}
		prof.Val = prof.Val[:0]
		// Get everything up to the first space, this is the key.  Not all keys are processed.
		for i, v = range prof.Line {
			if v == 0x20 {
				prof.Val = prof.Line[:i]
				pos = i + 1
				break
			}
		}
		// skip the intr line
		if prof.Val[0] == 'i' {
			continue
		}
		if prof.Val[0] == 'c' {
			if prof.Val[1] == 'p' { // process CPU
				cpu := CPU{ID: string(prof.Val[:])}
				j = 0
				// skip over any remaining spaces
				for i, v = range prof.Line[pos:] {
					if v != 0x20 {
						break
					}
					j++
				}
				fieldNum = 0
				pos, j = j+pos, j+pos
				// space is the field separator
				for i, v = range prof.Line[pos:] {
					if v == '\n' {
						stop = true
					}
					if v == 0x20 || stop {
						fieldNum++
						n, err = helpers.ParseUint(prof.Line[j : pos+i])
						if err != nil {
							return stats, &joe.ParseError{Info: string(prof.Val[:]), Err: err}
						}
						j = pos + i + 1
						if fieldNum < 6 {
							if fieldNum < 4 {
								if fieldNum == 1 {
									cpu.User = int64(n)
									continue
								}
								if fieldNum == 2 {
									cpu.Nice = int64(n)
									continue
								}
								cpu.System = int64(n) // 3
								continue
							}
							if fieldNum == 4 {
								cpu.Idle = int64(n)
								continue
							}
							cpu.IOWait = int64(n) // 5
							continue
						}
						if fieldNum < 8 {
							if fieldNum == 6 {
								cpu.IRQ = int64(n)
								continue
							}
							cpu.SoftIRQ = int64(n) // 7
							continue
						}
						if fieldNum == 8 {
							cpu.Steal = int64(n)
							continue
						}
						if fieldNum == 9 {
							cpu.Quest = int64(n)
							continue
						}
						cpu.QuestNice = int64(n) // 10
					}
				}
				stats.CPU = append(stats.CPU, cpu)
				stop = false
				continue
			}
			// Otherwise it's ctxt info; rest of the line is the data.
			n, err = helpers.ParseUint(prof.Line[pos : len(prof.Line)-1])
			if err != nil {
				return stats, &joe.ParseError{Info: string(prof.Val[:]), Err: err}
			}
			stats.Ctxt = int64(n)
			continue
		}
		if prof.Val[0] == 'b' {
			// rest of the line is the data
			n, err = helpers.ParseUint(prof.Line[pos : len(prof.Line)-1])
			if err != nil {
				return stats, &joe.ParseError{Info: string(prof.Val[:]), Err: err}
			}
			stats.BTime = int64(n)
			continue
		}
		if prof.Val[0] == 'p' && prof.Val[4] == 'e' { // processes info
			// rest of the line is the data
			n, err = helpers.ParseUint(prof.Line[pos : len(prof.Line)-1])
			if err != nil {
				return stats, &joe.ParseError{Info: string(prof.Val[:]), Err: err}
			}
			stats.Processes = int64(n)
			continue
		}
	}
	return stats, nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current kernel activity information using the package's
// global Profiler.
func Get() (stat *CPUStats, err error) {
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

// Ticker delivers the system's kernel activity information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan *CPUStats
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
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan *CPUStats), Profiler: p}
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

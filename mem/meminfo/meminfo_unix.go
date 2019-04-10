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

// Package meminfo gets the system's memory information, /proc/meminfo.
package meminfo

import (
	"io"
	"sync"
	"time"

	"github.com/rai-project/joefriday/helpers"
	joe "github.com/rai-project/joefriday"
)

const procFile = "/proc/meminfo"

// Info holds the memory information.
type Info struct {
	Timestamp         int64  `json:"timestamp"`
	Active            uint64 `json:"active"`
	ActiveAnon        uint64 `json:"active_anon"`
	ActiveFile        uint64 `json:"active_file"`
	AnonHugePages     uint64 `json:"anon_huge_pages"`
	AnonPages         uint64 `json:"anon_pages"`
	Bounce            uint64 `json:"bounce"`
	Buffers           uint64 `json:"buffers"`
	Cached            uint64 `json:"cached"`
	CommitLimit       uint64 `json:"commit_limit"`
	CommittedAS       uint64 `json:"commited_as"`
	DirectMap4K       uint64 `json:"direct_map_4k"`
	DirectMap2M       uint64 `json:"direct_map_2m"`
	Dirty             uint64 `json:"dirty"`
	HardwareCorrupted uint64 `json:"hardware_corrupted"`
	HugePagesFree     uint64 `json:"huge_pages_free"`
	HugePagesRsvd     uint64 `json:"huge_pages_rsvd"`
	HugePagesSize     uint64 `json:"huge_pages_size"`
	HugePagesSurp     uint64 `json:"huge_pages_surp"`
	HugePagesTotal    uint64 `json:"huge_pages_total"`
	Inactive          uint64 `json:"inactive"`
	InactiveAnon      uint64 `json:"inactive_anon"`
	InactiveFile      uint64 `json:"inactive_file"`
	KernelStack       uint64 `json:"kernel_stack"`
	Mapped            uint64 `json:"mapped"`
	MemAvailable      uint64 `json:"mem_available"`
	MemFree           uint64 `json:"mem_free"`
	MemTotal          uint64 `json:"mem_total"`
	Mlocked           uint64 `json:"mlocked"`
	NFSUnstable       uint64 `json:"nfs_unstable"`
	PageTables        uint64 `json:"page_tables"`
	Shmem             uint64 `json:"shmem"`
	Slab              uint64 `json:"slab"`
	SReclaimable      uint64 `json:"s_reclaimable"`
	SUnreclaim        uint64 `json:"s_unreclaim"`
	SwapCached        uint64 `json:"swap_cached"`
	SwapFree          uint64 `json:"swap_free"`
	SwapTotal         uint64 `json:"swap_total"`
	Unevictable       uint64 `json:"unevictable"`
	VmallocChunk      uint64 `json:"vmalloc_chunk"`
	VmallocTotal      uint64 `json:"vmalloc_total"`
	VmallocUsed       uint64 `json:"vmalloc_used"`
	Writeback         uint64 `json:"writeback"`
	WritebackTmp      uint64 `json:"writeback_tmp"`
}

// Profiler is used to get the memory information by processing the
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

// Get returns the current memory information.
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
			if prof.Val[5] == 'e' {
				if nameLen == 6 {
					inf.Active = n
					continue
				}
				if prof.Val[7] == 'a' {
					inf.ActiveAnon = n
					continue
				}
				inf.ActiveFile = n
				continue
			}
			if nameLen == 9 {
				inf.AnonPages = n
				continue
			}
			inf.AnonHugePages = n
			continue
		}
		if v == 'C' {
			if nameLen == 6 {
				inf.Cached = n
				continue
			}
			if nameLen == 11 {
				inf.CommitLimit = n
				continue
			}
			inf.CommittedAS = n
			continue
		}
		if v == 'D' {
			if nameLen == 5 {
				inf.Dirty = n
				continue
			}
			if prof.Val[10] == 'k' {
				inf.DirectMap4K = n
				continue
			}
			inf.DirectMap2M = n
			continue
		}
		if v == 'H' {
			if nameLen == 14 {
				if prof.Val[10] == 'F' {
					inf.HugePagesFree = n
					continue
				}
				if prof.Val[10] == 'R' {
					inf.HugePagesRsvd = n
					continue
				}
				inf.HugePagesSurp = n
			}
			if prof.Val[1] == 'a' {
				inf.HardwareCorrupted = n
				continue
			}
			if prof.Val[9] == 'i' {
				inf.HugePagesSize = n
				continue
			}
			inf.HugePagesTotal = n
			continue
		}
		if v == 'I' {
			if nameLen == 8 {
				inf.Inactive = n
				continue
			}
			if prof.Val[9] == 'a' {
				inf.InactiveAnon = n
				continue
			}
			inf.InactiveFile = n
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
					continue
				}
				inf.Mlocked = n
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
			if v == 'h' {
				inf.Shmem = n
				continue
			}
			if v == 'l' {
				inf.Slab = n
				continue
			}
			if v == 'R' {
				inf.SReclaimable = n
				continue
			}
			inf.SUnreclaim = n
			continue
		}
		if v == 'V' {
			if prof.Val[8] == 'C' {
				inf.VmallocChunk = n
				continue
			}
			if prof.Val[8] == 'T' {
				inf.VmallocTotal = n
				continue
			}
			inf.VmallocUsed = n
			continue
		}
		if v == 'W' {
			if nameLen == 9 {
				inf.Writeback = n
				continue
			}
			inf.WritebackTmp = n
			continue
		}
		if v == 'B' {
			if nameLen == 6 {
				inf.Bounce = n
				continue
			}
			inf.Buffers = n
			continue
		}
		if v == 'K' {
			inf.KernelStack = n
			continue
		}
		if v == 'N' {
			inf.NFSUnstable = n
			continue
		}
		if v == 'P' {
			inf.PageTables = n
		}
		inf.Unevictable = n
	}
	return inf, nil
}

var std *Profiler
var stdMu sync.Mutex //protects standard to prevent a data race on checking/instantiation

// Get returns the current memory information using the package's global
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

// Ticker delivers the system's memory information at intervals.
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
					if t.Val[5] == 'e' {
						if nameLen == 6 {
							inf.Active = n
							continue
						}
						if t.Val[7] == 'a' {
							inf.ActiveAnon = n
							continue
						}
						inf.ActiveFile = n
						continue
					}
					if nameLen == 9 {
						inf.AnonPages = n
						continue
					}
					inf.AnonHugePages = n
					continue
				}
				if v == 'C' {
					if nameLen == 6 {
						inf.Cached = n
						continue
					}
					if nameLen == 11 {
						inf.CommitLimit = n
						continue
					}
					inf.CommittedAS = n
					continue
				}
				if v == 'D' {
					if nameLen == 5 {
						inf.Dirty = n
						continue
					}
					if t.Val[10] == 'k' {
						inf.DirectMap4K = n
						continue
					}
					inf.DirectMap2M = n
					continue
				}
				if v == 'H' {
					if nameLen == 14 {
						if t.Val[10] == 'F' {
							inf.HugePagesFree = n
							continue
						}
						if t.Val[10] == 'R' {
							inf.HugePagesRsvd = n
							continue
						}
						inf.HugePagesSurp = n
					}
					if t.Val[1] == 'a' {
						inf.HardwareCorrupted = n
						continue
					}
					if t.Val[9] == 'i' {
						inf.HugePagesSize = n
						continue
					}
					inf.HugePagesTotal = n
					continue
				}
				if v == 'I' {
					if nameLen == 8 {
						inf.Inactive = n
						continue
					}
					if t.Val[9] == 'a' {
						inf.InactiveAnon = n
						continue
					}
					inf.InactiveFile = n
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
							continue
						}
						inf.Mlocked = n
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
					if v == 'h' {
						inf.Shmem = n
						continue
					}
					if v == 'l' {
						inf.Slab = n
						continue
					}
					if v == 'R' {
						inf.SReclaimable = n
						continue
					}
					inf.SUnreclaim = n
					continue
				}
				if v == 'V' {
					if t.Val[8] == 'C' {
						inf.VmallocChunk = n
						continue
					}
					if t.Val[8] == 'T' {
						inf.VmallocTotal = n
						continue
					}
					inf.VmallocUsed = n
					continue
				}
				if v == 'W' {
					if nameLen == 9 {
						inf.Writeback = n
						continue
					}
					inf.WritebackTmp = n
					continue
				}
				if v == 'B' {
					if nameLen == 6 {
						inf.Bounce = n
						continue
					}
					inf.Buffers = n
					continue
				}
				if v == 'K' {
					inf.KernelStack = n
					continue
				}
				if v == 'N' {
					inf.NFSUnstable = n
					continue
				}
				if v == 'P' {
					inf.PageTables = n
				}
				inf.Unevictable = n
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

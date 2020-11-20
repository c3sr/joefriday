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

// Package meminfo processes a subset of the /proc/meminfo file. Instead of
// returning a Go struct, it returns Flatbuffer serialized bytes. A function to
// deserialize the Flatbuffer serialized bytes into a meminfo.Info struct is
// provided.
//
// Note: the package name is meminfo and not the final element of the import
// path (flat). 
package meminfo

import (
	"io"
	"sync"
	"time"

	"github.com/c3sr/joefriday/helpers"
	fb "github.com/google/flatbuffers/go"
	joe "github.com/c3sr/joefriday"
	mem "github.com/c3sr/joefriday/mem/meminfo"
	"github.com/c3sr/joefriday/mem/meminfo/flat/structs"
)

// Profiler is used to get the memory information as Flatbuffer serialized
// bytes by processing the /proc/meminfo file.
type Profiler struct {
	*mem.Profiler
	*fb.Builder
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	p, err := mem.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current memory information as Flatbuffer serialized bytes.
func (prof *Profiler) Get() ([]byte, error) {
	inf, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf), nil
}

var std *Profiler
var stdMu sync.Mutex //protects standard to prevent a data race on checking/instantiation

// Get returns the current memory information as Flatbuffer serialized bytes
// using the package's global Profiler.
func Get() (p []byte, err error) {
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

// Serialize the memory information using Flatbuffers.
func (prof *Profiler) Serialize(inf *mem.Info) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	structs.InfoStart(prof.Builder)
	structs.InfoAddTimestamp(prof.Builder, inf.Timestamp)
	structs.InfoAddActive(prof.Builder, inf.Active)
	structs.InfoAddActiveAnon(prof.Builder, inf.ActiveAnon)
	structs.InfoAddActiveFile(prof.Builder, inf.ActiveFile)
	structs.InfoAddAnonHugePages(prof.Builder, inf.AnonHugePages)
	structs.InfoAddAnonPages(prof.Builder, inf.AnonPages)
	structs.InfoAddBounce(prof.Builder, inf.Bounce)
	structs.InfoAddBuffers(prof.Builder, inf.Buffers)
	structs.InfoAddCached(prof.Builder, inf.Cached)
	structs.InfoAddCommitLimit(prof.Builder, inf.CommitLimit)
	structs.InfoAddCommittedAS(prof.Builder, inf.CommittedAS)
	structs.InfoAddDirectMap4K(prof.Builder, inf.DirectMap4K)
	structs.InfoAddDirectMap2M(prof.Builder, inf.DirectMap2M)
	structs.InfoAddDirty(prof.Builder, inf.Dirty)
	structs.InfoAddHardwareCorrupted(prof.Builder, inf.HardwareCorrupted)
	structs.InfoAddHugePagesFree(prof.Builder, inf.HugePagesFree)
	structs.InfoAddHugePagesRsvd(prof.Builder, inf.HugePagesRsvd)
	structs.InfoAddHugePagesSize(prof.Builder, inf.HugePagesSize)
	structs.InfoAddHugePagesSurp(prof.Builder, inf.HugePagesSurp)
	structs.InfoAddHugePagesTotal(prof.Builder, inf.HugePagesTotal)
	structs.InfoAddInactive(prof.Builder, inf.Inactive)
	structs.InfoAddInactiveAnon(prof.Builder, inf.InactiveAnon)
	structs.InfoAddInactiveFile(prof.Builder, inf.InactiveFile)
	structs.InfoAddKernelStack(prof.Builder, inf.KernelStack)
	structs.InfoAddMapped(prof.Builder, inf.Mapped)
	structs.InfoAddMemAvailable(prof.Builder, inf.MemAvailable)
	structs.InfoAddMemFree(prof.Builder, inf.MemFree)
	structs.InfoAddMemTotal(prof.Builder, inf.MemTotal)
	structs.InfoAddMlocked(prof.Builder, inf.Mlocked)
	structs.InfoAddNFSUnstable(prof.Builder, inf.NFSUnstable)
	structs.InfoAddPageTables(prof.Builder, inf.PageTables)
	structs.InfoAddShmem(prof.Builder, inf.Shmem)
	structs.InfoAddSlab(prof.Builder, inf.Slab)
	structs.InfoAddSReclaimable(prof.Builder, inf.SReclaimable)
	structs.InfoAddSUnreclaim(prof.Builder, inf.SUnreclaim)
	structs.InfoAddSwapCached(prof.Builder, inf.SwapCached)
	structs.InfoAddSwapFree(prof.Builder, inf.SwapFree)
	structs.InfoAddSwapTotal(prof.Builder, inf.SwapTotal)
	structs.InfoAddUnevictable(prof.Builder, inf.Unevictable)
	structs.InfoAddVmallocChunk(prof.Builder, inf.VmallocChunk)
	structs.InfoAddVmallocTotal(prof.Builder, inf.VmallocTotal)
	structs.InfoAddVmallocUsed(prof.Builder, inf.VmallocUsed)
	structs.InfoAddWriteback(prof.Builder, inf.Writeback)
	structs.InfoAddWritebackTmp(prof.Builder, inf.WritebackTmp)
	prof.Builder.Finish(structs.InfoEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize the memory information using Flatbuffers with the package's global
// Profiler.
func Serialize(inf *mem.Info) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(inf), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them
// as meminfo.Info.
func Deserialize(p []byte) *mem.Info {
	infoFlat := structs.GetRootAsInfo(p, 0)
	info := &mem.Info{}
	info.Timestamp = infoFlat.Timestamp()
	info.Active = infoFlat.Active()
	info.ActiveAnon = infoFlat.ActiveAnon()
	info.ActiveFile = infoFlat.ActiveFile()
	info.AnonHugePages = infoFlat.AnonHugePages()
	info.AnonPages = infoFlat.AnonPages()
	info.Bounce = infoFlat.Bounce()
	info.Buffers = infoFlat.Buffers()
	info.Cached = infoFlat.Cached()
	info.CommitLimit = infoFlat.CommitLimit()
	info.CommittedAS = infoFlat.CommittedAS()
	info.DirectMap4K = infoFlat.DirectMap4K()
	info.DirectMap2M = infoFlat.DirectMap2M()
	info.Dirty = infoFlat.Dirty()
	info.HardwareCorrupted = infoFlat.HardwareCorrupted()
	info.HugePagesFree = infoFlat.HugePagesFree()
	info.HugePagesRsvd = infoFlat.HugePagesRsvd()
	info.HugePagesSize = infoFlat.HugePagesSize()
	info.HugePagesSurp = infoFlat.HugePagesSurp()
	info.HugePagesTotal = infoFlat.HugePagesTotal()
	info.Inactive = infoFlat.Inactive()
	info.InactiveAnon = infoFlat.InactiveAnon()
	info.InactiveFile = infoFlat.InactiveFile()
	info.KernelStack = infoFlat.KernelStack()
	info.Mapped = infoFlat.Mapped()
	info.MemAvailable = infoFlat.MemAvailable()
	info.MemFree = infoFlat.MemFree()
	info.MemTotal = infoFlat.MemTotal()
	info.Mlocked = infoFlat.Mlocked()
	info.NFSUnstable = infoFlat.NFSUnstable()
	info.PageTables = infoFlat.PageTables()
	info.Shmem = infoFlat.Shmem()
	info.Slab = infoFlat.Slab()
	info.SReclaimable = infoFlat.SReclaimable()
	info.SUnreclaim = infoFlat.SUnreclaim()
	info.SwapCached = infoFlat.SwapCached()
	info.SwapFree = infoFlat.SwapFree()
	info.SwapTotal = infoFlat.SwapTotal()
	info.Unevictable = infoFlat.Unevictable()
	info.VmallocChunk = infoFlat.VmallocChunk()
	info.VmallocTotal = infoFlat.VmallocTotal()
	info.VmallocUsed = infoFlat.VmallocUsed()
	info.Writeback = infoFlat.Writeback()
	info.WritebackTmp = infoFlat.WritebackTmp()
	return info
}

// Ticker delivers the system's memory information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan []byte
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
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan []byte), Profiler: p}
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
	)
	// ticker
Tick:
	for {
		select {
		case <-t.Done:
			return
		case <-t.C:
			t.Builder.Reset()
			err = t.Profiler.Profiler.Reset()
			if err != nil {
				t.Errs <- err
				continue
			}
			structs.InfoStart(t.Builder)
			structs.InfoAddTimestamp(t.Builder, time.Now().UTC().UnixNano())
			for {
				t.Val = t.Val[:0]
				t.Line, err = t.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					// An error results in sending error message and stop processing of this tick.
					t.Errs <- &joe.ReadError{Err: err}
					continue Tick
				}
				// first grab the key name (everything up to the ':')
				for i, v = range t.Line {
					if v == 0x3A {
						t.Val = t.Line[:i]
						pos = i + 1 // skip the :
						break
					}
				}
				nameLen = len(t.Val)
				// skip all spaces
				for i, v = range t.Line[pos:] {
					if v != 0x20 {
						pos += i
						break
					}
				}

				// grab the numbers
				for _, v = range t.Line[pos:] {
					if v == 0x20 || v == '\n' {
						break
					}
					t.Val = append(t.Val, v)
				}
				// any conversion error results in 0
				n, err = helpers.ParseUint(t.Val[nameLen:])
				if err != nil {
					t.Errs <- &joe.ParseError{Info: string(t.Val[:nameLen]), Err: err}
					continue
				}
				v = t.Val[0]
				if v == 'A' {
					if t.Val[5] == 'e' {
						if nameLen == 6 {
							structs.InfoAddActive(t.Builder, n)
							continue
						}
						if t.Val[7] == 'a' {
							structs.InfoAddActiveAnon(t.Builder, n)
							continue
						}
						structs.InfoAddActiveFile(t.Builder, n)
						continue
					}
					if nameLen == 9 {
						structs.InfoAddAnonPages(t.Builder, n)
						continue
					}
					structs.InfoAddAnonHugePages(t.Builder, n)
					continue
				}
				if v == 'C' {
					if nameLen == 6 {
						structs.InfoAddCached(t.Builder, n)
						continue
					}
					if nameLen == 11 {
						structs.InfoAddCommitLimit(t.Builder, n)
						continue
					}
					structs.InfoAddCommittedAS(t.Builder, n)
					continue
				}
				if v == 'D' {
					if nameLen == 5 {
						structs.InfoAddDirty(t.Builder, n)
						continue
					}
					if t.Val[10] == 'k' {
						structs.InfoAddDirectMap4K(t.Builder, n)
						continue
					}
					structs.InfoAddDirectMap2M(t.Builder, n)
					continue
				}
				if v == 'H' {
					if nameLen == 14 {
						if t.Val[10] == 'F' {
							structs.InfoAddHugePagesFree(t.Builder, n)
							continue
						}
						if t.Val[10] == 'R' {
							structs.InfoAddHugePagesRsvd(t.Builder, n)
							continue
						}
						structs.InfoAddHugePagesSurp(t.Builder, n)
					}
					if t.Val[1] == 'a' {
						structs.InfoAddHardwareCorrupted(t.Builder, n)
						continue
					}
					if t.Val[9] == 'i' {
						structs.InfoAddHugePagesSize(t.Builder, n)
						continue
					}
					structs.InfoAddHugePagesTotal(t.Builder, n)
					continue
				}
				if v == 'I' {
					if nameLen == 8 {
						structs.InfoAddInactive(t.Builder, n)
						continue
					}
					if t.Val[9] == 'a' {
						structs.InfoAddInactiveAnon(t.Builder, n)
						continue
					}
					structs.InfoAddInactiveFile(t.Builder, n)
				}
				if v == 'M' {
					v = t.Val[3]
					if nameLen < 8 {
						if v == 'p' {
							structs.InfoAddMapped(t.Builder, n)
							continue
						}
						if v == 'F' {
							structs.InfoAddMemFree(t.Builder, n)
							continue
						}
						structs.InfoAddMlocked(t.Builder, n)
						continue
					}
					if v == 'A' {
						structs.InfoAddMemAvailable(t.Builder, n)
						continue
					}
					structs.InfoAddMemTotal(t.Builder, n)
					continue
				}
				if v == 'S' {
					v = t.Val[1]
					if v == 'w' {
						if t.Val[4] == 'C' {
							structs.InfoAddSwapCached(t.Builder, n)
							continue
						}
						if t.Val[4] == 'F' {
							structs.InfoAddSwapFree(t.Builder, n)
							continue
						}
						structs.InfoAddSwapTotal(t.Builder, n)
						continue
					}
					if v == 'h' {
						structs.InfoAddShmem(t.Builder, n)
						continue
					}
					if v == 'l' {
						structs.InfoAddSlab(t.Builder, n)
						continue
					}
					if v == 'R' {
						structs.InfoAddSReclaimable(t.Builder, n)
						continue
					}
					structs.InfoAddSUnreclaim(t.Builder, n)
					continue
				}
				if v == 'V' {
					if t.Val[8] == 'C' {
						structs.InfoAddVmallocChunk(t.Builder, n)
						continue
					}
					if t.Val[8] == 'T' {
						structs.InfoAddVmallocTotal(t.Builder, n)
						continue
					}
					structs.InfoAddVmallocUsed(t.Builder, n)
					continue
				}
				if v == 'W' {
					if nameLen == 9 {
						structs.InfoAddWriteback(t.Builder, n)
						continue
					}
					structs.InfoAddWritebackTmp(t.Builder, n)
					continue
				}
				if v == 'B' {
					if nameLen == 6 {
						structs.InfoAddBounce(t.Builder, n)
						continue
					}
					structs.InfoAddBuffers(t.Builder, n)
					continue
				}
				if v == 'K' {
					structs.InfoAddKernelStack(t.Builder, n)
					continue
				}
				if v == 'N' {
					structs.InfoAddNFSUnstable(t.Builder, n)
					continue
				}
				if v == 'P' {
					structs.InfoAddPageTables(t.Builder, n)
				}
				structs.InfoAddUnevictable(t.Builder, n)
			}
			t.Builder.Finish(structs.InfoEnd(t.Builder))
			t.Data <- t.Profiler.Builder.Bytes[t.Builder.Head():]
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}

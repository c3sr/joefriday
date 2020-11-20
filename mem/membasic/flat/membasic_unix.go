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

// Package membasic processes a subset of the /proc/meminfo file. Instead of
// returning a Go struct, it returns Flatbuffer serialized bytes. A function to
// deserialize the Flatbuffer serialized bytes into a membasic.Info struct is
// provided. For more detailed information about a system's memory, use the
// meminfo package.
//
// Note: the package name is membasic and not the final element of the import
// path (flat).
package membasic

import (
	"io"
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/c3sr/joefriday"
	"github.com/c3sr/joefriday/helpers"
	basic "github.com/c3sr/joefriday/mem/membasic"
	"github.com/c3sr/joefriday/mem/membasic/flat/structs"
)

// Profiler is used to get the basic memory information as Flatbuffer
// serialized by processing the /proc/meminfo file.
type Profiler struct {
	*basic.Profiler
	*fb.Builder
}

// Returns an initialized Profiler that utilizes FlatBuffers; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	p, err := basic.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current basic memory information as Flatbuffer serialized
// bytes.
func (prof *Profiler) Get() ([]byte, error) {
	inf, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf), nil
}

var std *Profiler
var stdMu sync.Mutex //protects standard to prevent a data race on checking/instantiation

// Get returns the current basic memory information as Flatbuffer serialized
// bytes using the package's global Profiler.
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

// Serialize the basic memory information using Flatbuffers.
func (prof *Profiler) Serialize(inf *basic.Info) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	structs.InfoStart(prof.Builder)
	structs.InfoAddTimestamp(prof.Builder, inf.Timestamp)
	structs.InfoAddActive(prof.Builder, inf.Active)
	structs.InfoAddInactive(prof.Builder, inf.Inactive)
	structs.InfoAddMapped(prof.Builder, inf.Mapped)
	structs.InfoAddMemAvailable(prof.Builder, inf.MemAvailable)
	structs.InfoAddMemFree(prof.Builder, inf.MemFree)
	structs.InfoAddMemTotal(prof.Builder, inf.MemTotal)
	structs.InfoAddSwapCached(prof.Builder, inf.SwapCached)
	structs.InfoAddSwapFree(prof.Builder, inf.SwapFree)
	structs.InfoAddSwapTotal(prof.Builder, inf.SwapTotal)
	prof.Builder.Finish(structs.InfoEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize the basic memory information using Flatbuffers with the package's
// global Profiler.
func Serialize(inf *basic.Info) (p []byte, err error) {
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

// Deserialize takes some Flatbuffer serialized bytes and deserializes them
// as membasic.Info.
func Deserialize(p []byte) *basic.Info {
	infoFlat := structs.GetRootAsInfo(p, 0)
	info := &basic.Info{}
	info.Timestamp = infoFlat.Timestamp()
	info.Active = infoFlat.Active()
	info.Inactive = infoFlat.Inactive()
	info.Mapped = infoFlat.Mapped()
	info.MemAvailable = infoFlat.MemAvailable()
	info.MemFree = infoFlat.MemFree()
	info.MemTotal = infoFlat.MemTotal()
	info.SwapCached = infoFlat.SwapCached()
	info.SwapFree = infoFlat.SwapFree()
	info.SwapTotal = infoFlat.SwapTotal()
	return info
}

// Ticker delivers the system's basic memory information at intervals.
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
					if t.Val[5] == 'e' && nameLen == 6 {
						structs.InfoAddActive(t.Builder, n)
					}
					continue
				}
				if v == 'I' {
					if nameLen == 8 {
						structs.InfoAddInactive(t.Builder, n)
					}
					continue
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
						}
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
					}
				}
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

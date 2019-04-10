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

// Package netusage gets the usage of the network devices. Usage is calculated
// by taking the difference between two network device snapshots,
// /proc/net/dev. The time elapsed between the two snapshots is stored in the
// TimeDelta field.
package netusage

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/rai-project/joefriday/helpers"
	joe "github.com/rai-project/joefriday"
	"github.com/rai-project/joefriday/net/netdev"
	"github.com/rai-project/joefriday/net/structs"
)

// Profiler is used to process the network device usage.
type Profiler struct {
	*netdev.Profiler
	prior structs.DevInfo
}

// Returns an initialized Profiler; ready to use. Upon creation, a
// /proc/net/dev snapshot is taken so that any Get() will return valid
// information.
func NewProfiler() (prof *Profiler, err error) {
	p, err := netdev.NewProfiler()
	if err != nil {
		return nil, err
	}
	prior, err := p.Get()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, prior: *prior}, nil
}

// Get returns the current network device usage. Calculating usage requires two
// snapshots. This func gets the current snapshot of /proc/net/dev and
// calculates the difference between that and the prior snapshot. The current
// snapshot is stored for use as the prior snapshot on the next Get call. If
// ongoing usage information is desired, the Ticker should be used; it's better
// suited for ongoing usage information.
func (prof *Profiler) Get() (u *structs.DevUsage, err error) {
	infCur, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	u = prof.CalculateUsage(infCur)
	prof.prior = *infCur
	return u, nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current network device usage using the package's global
// Profiler. The profiler is lazily instantiated. If the profiler doesn't
// already exist, the first usage information will not be useful due to minimal
// time elapsing between the initial and second snapshots used for usage
// calculations; the results of the first call should be discarded.
func Get() (u *structs.DevUsage, err error) {
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

// CalculateUsage returns the difference between the current /proc/net/dev
// snapshot and the prior one.
func (prof *Profiler) CalculateUsage(cur *structs.DevInfo) *structs.DevUsage {
	u := &structs.DevUsage{
		Timestamp:  cur.Timestamp,
		TimeDelta:  cur.Timestamp - prof.prior.Timestamp,
		Device: make([]structs.Device, len(cur.Device)),
	}
	for i := 0; i < len(cur.Device); i++ {
		u.Device[i].Name = cur.Device[i].Name
		u.Device[i].RBytes = cur.Device[i].RBytes - prof.prior.Device[i].RBytes
		u.Device[i].RPackets = cur.Device[i].RPackets - prof.prior.Device[i].RPackets
		u.Device[i].RErrs = cur.Device[i].RErrs - prof.prior.Device[i].RErrs
		u.Device[i].RDrop = cur.Device[i].RDrop - prof.prior.Device[i].RDrop
		u.Device[i].RFIFO = cur.Device[i].RFIFO - prof.prior.Device[i].RFIFO
		u.Device[i].RFrame = cur.Device[i].RFrame - prof.prior.Device[i].RFrame
		u.Device[i].RCompressed = cur.Device[i].RCompressed - prof.prior.Device[i].RCompressed
		u.Device[i].RMulticast = cur.Device[i].RMulticast - prof.prior.Device[i].RMulticast
		u.Device[i].TBytes = cur.Device[i].TBytes - prof.prior.Device[i].TBytes
		u.Device[i].TPackets = cur.Device[i].TPackets - prof.prior.Device[i].TPackets
		u.Device[i].TErrs = cur.Device[i].TErrs - prof.prior.Device[i].TErrs
		u.Device[i].TDrop = cur.Device[i].TDrop - prof.prior.Device[i].TDrop
		u.Device[i].TFIFO = cur.Device[i].TFIFO - prof.prior.Device[i].TFIFO
		u.Device[i].TColls = cur.Device[i].TColls - prof.prior.Device[i].TColls
		u.Device[i].TCarrier = cur.Device[i].TCarrier - prof.prior.Device[i].TCarrier
		u.Device[i].TCompressed = cur.Device[i].TCompressed - prof.prior.Device[i].TCompressed
	}
	return u
}

// Ticker delivers the system's network device usage at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan *structs.DevUsage
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
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan *structs.DevUsage), Profiler: p}
	go t.Run()
	return &t, nil
}

// Run runs the ticker.
func (t *Ticker) Run() {
	var (
		i, pos, line, fieldNum int
		n                      uint64
		v                      byte
		err                    error
		cur                    structs.DevInfo
		dev                 structs.Device
	)
	// ticker
	for {
		select {
		case <-t.Done:
			return
		case <-t.C:
			cur.Timestamp = time.Now().UTC().UnixNano()
			err = t.Reset()
			if err != nil {
				t.Errs <- err
				break
			}
			line = 0
			cur.Device = cur.Device[:0]
			// read each line until eof
			for {
				t.Val = t.Val[:0]
				t.Line, err = t.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					t.Errs <- &joe.ReadError{Err: err}
					break
				}
				line++
				if line < 3 {
					continue
				}
				// first grab the interface name (everything up to the ':')
				for i, v = range t.Line {
					if v == 0x3A {
						pos = i + 1
						break
					}
					// skip spaces
					if v != 0x20 {
						continue
					}
					t.Val = append(t.Val, v)
				}
				dev.Name = string(t.Val[:])
				fieldNum = 0
				// process the rest of the line
				for {
					fieldNum++
					// skip all spaces
					for i, v = range t.Line[pos:] {
						if v != 0x20 {
							pos += i
							break
						}
					}
					// grab the numbers
					for i, v = range t.Line[pos:] {
						if v == 0x20 || v == '\n' {
							break
						}
					}
					// any conversion error results in 0
					n, err = helpers.ParseUint(t.Line[pos : pos+1])
					pos += i
					if err != nil {
						t.Errs <- &joe.ParseError{Info: fmt.Sprintf("line %d: field %d", line, fieldNum), Err: err}
						continue
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
					if fieldNum == 16 {
						dev.TCompressed = int64(n)
						break
					}
				}
				cur.Device = append(cur.Device, dev)
			}
			t.Data <- t.CalculateUsage(&cur)
			t.prior.Timestamp = cur.Timestamp
			if len(t.prior.Device) != len(cur.Device) {
				t.prior.Device = make([]structs.Device, len(cur.Device))
			}
			copy(t.prior.Device, cur.Device)
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}

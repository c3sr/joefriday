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

// Package mem gets and processes /proc/meminfo, returning the data in the
// appropriate format.
package mem

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/SermoDigital/helpers"
	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/mem/flat"
)

const procMemInfo = "/proc/meminfo"

type InfoProfiler struct {
	joe.Proc
	val []byte
}

var std *InfoProfiler

func NewInfoProfiler() (proc *InfoProfiler, err error) {
	f, err := os.Open(procMemInfo)
	if err != nil {
		return nil, err
	}
	return &InfoProfiler{Proc: joe.Proc{File: f, Buf: bufio.NewReader(f)}, val: make([]byte, 32)}, nil
}

// It is expected that the caller has the lock.
func (p *InfoProfiler) reset() error {
	_, err := p.File.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}
	p.Buf.Reset(p.File)
	return nil
}

// Get returns some of the results of /proc/meminfo.
func (p *InfoProfiler) Get() (inf *Info, err error) {
	p.Lock()
	defer p.Unlock()

	var (
		i, pos, nameLen int
		v               byte
	)
	inf = &Info{}
	for l := 0; l < 16; l++ {
		p.Line, err = p.Buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return inf, fmt.Errorf("error reading output bytes: %s", err)
		}
		if l > 8 && l < 14 {
			continue
		}
		// first grab the key name (everything up to the ':')
		for i, v = range p.Line {
			if v == ':' {
				pos = i + 1
				break
			}
			p.val = append(p.val, v)
		}
		nameLen = len(p.val)

		// skip all spaces
		for i, v = range p.Line[pos:] {
			if v != ' ' {
				pos += i
				break
			}
		}

		// grab the numbers
		for _, v = range p.Line[pos:] {
			if v == ' ' || v == '\n' {
				break
			}
			p.val = append(p.val, v)
		}
		// any conversion error results in 0
		n, err := helpers.ParseUint(p.val[nameLen:])
		if err != nil {
			return inf, fmt.Errorf("%s: %s", p.val[:nameLen], err)
		}

		v = p.val[0]

		// Reduce evaluations.
		if v == 'M' {
			v = p.val[3]
			if v == 'T' {
				inf.MemTotal = int64(n)
			} else if v == 'F' {
				inf.MemFree = int64(n)
			} else {
				inf.MemAvailable = int64(n)
			}
		} else if v == 'S' {
			v = p.val[4]
			if v == 'C' {
				inf.SwapCached = int64(n)
			} else if v == 'T' {
				inf.SwapTotal = int64(n)
			} else if v == 'F' {
				inf.SwapFree = int64(n)
			}
		} else if v == 'B' {
			inf.Buffers = int64(n)
		} else if v == 'I' {
			inf.Inactive = int64(n)
		} else if v == 'C' {
			inf.Cached = int64(n)
		} else if v == 'A' {
			inf.Active = int64(n)
		}
		p.val = p.val[:0]
	}
	inf.Timestamp = time.Now().UTC().UnixNano()
	return inf, nil
}

// GetInfo get's the current meminfo.
func GetInfo() (inf *Info, err error) {
	if std == nil {
		std, err = NewInfoProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Get()
}

// GetFlat returns the current meminfo as flatbuffer serialized bytes.
func (p *InfoProfiler) GetFlat() ([]byte, error) {
	inf, err := p.Get()
	if err != nil {
		return nil, err
	}
	return inf.SerializeFlat(), nil
}

// GetInfoFlat returns the current meminfo as flatbuffer serialized bytes.
func GetInfoFlat() ([]byte, error) {
	var err error
	if std == nil {
		std, err = NewInfoProfiler()
		if err != nil {
			return nil, err
		}
	}
	inf, err := std.Get()
	if err != nil {
		return nil, err
	}
	return inf.SerializeFlat(), nil
}

// GetJSON returns the current meminfo as JSON serialized bytes.
func (p *InfoProfiler) GetJSON() ([]byte, error) {
	inf, err := p.Get()
	if err != nil {
		return nil, err
	}
	return inf.SerializeJSON()
}

// GetInfoJSON returns the current meminfo as JSON serialized bytes.
func GetInfoJSON() ([]byte, error) {
	var err error
	if std == nil {
		std, err = NewInfoProfiler()
		if err != nil {
			return nil, err
		}
	}
	inf, err := std.Get()
	if err != nil {
		return nil, err
	}
	return inf.SerializeJSON()
}

// Ticker gathers the meminfo on a ticker, whose interval is defined by the
// received duration, and sends the results to the channel.  The output is
// Flatbuffer serialized bytes of Info.  Any error encountered during
// processing is sent to the error channel; processing will continue.
//
// If an error occurs while opening /proc/meminfo, the error will be sent
// to the errs channel and this func will exit.
//
// To stop processing and exit; send a signal on the done channel.  This
// will cause the function to stop the ticker, close the out channel and
// return.
func (p *InfoProfiler) Ticker(interval time.Duration, out chan Info, done chan struct{}, errs chan error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(out)
	// predeclare some vars
	var (
		l, i, pos, nameLen int
		v                  byte
		n                  uint64
		err                error
		inf                Info
	)
	// Lock now because the for loop unlocks to simplify unlock logic when
	// a continue occurs (instead of the tick completing.)
	p.Lock()
	// ticker
	for {
		p.Unlock()
		select {
		case <-done:
			return
		case <-ticker.C:
			p.Lock()
			err = p.reset()
			if err != nil {
				errs <- joe.Error{Type: "mem", Op: "seek byte 0: /proc/meminfo", Err: err}
				continue
			}
			p.Line, err = p.Buf.ReadSlice('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				errs <- fmt.Errorf("error reading output bytes: %s", err)
				continue
			}
			if l > 8 && l < 14 {
				continue
			}
			// first grab the key name (everything up to the ':')
			for i, v = range p.Line {
				if v == ':' {
					pos = i + 1
					break
				}
				p.val = append(p.val, v)
			}
			nameLen = len(p.val)

			// skip all spaces
			for i, v = range p.Line[pos:] {
				if v != ' ' {
					pos += i
					break
				}
			}

			// grab the numbers
			for _, v = range p.Line[pos:] {
				if v == ' ' || v == '\n' {
					break
				}
				p.val = append(p.val, v)
			}
			// any conversion error results in 0
			n, err = helpers.ParseUint(p.val[nameLen:])
			if err != nil {
				errs <- fmt.Errorf("%s: %s", p.val[:nameLen], err)
			}
			v = p.val[0]

			// Reduce evaluations.
			if v == 'M' {
				v = p.val[3]
				if v == 'T' {
					inf.MemTotal = int64(n)
				} else if v == 'F' {
					inf.MemFree = int64(n)
				} else {
					inf.MemAvailable = int64(n)
				}
			} else if v == 'S' {
				v = p.val[4]
				if v == 'C' {
					inf.SwapCached = int64(n)
				} else if v == 'T' {
					inf.SwapTotal = int64(n)
				} else if v == 'F' {
					inf.SwapFree = int64(n)
				}
			} else if v == 'B' {
				inf.Buffers = int64(n)
			} else if v == 'I' {
				inf.Inactive = int64(n)
			} else if v == 'C' {
				inf.Cached = int64(n)
			} else if v == 'A' {
				inf.Active = int64(n)
			}
			p.val = p.val[:0]
		}
		inf.Timestamp = time.Now().UTC().UnixNano()
		out <- inf
	}
}

// InfoTicker gathers the meminfo on a ticker, whose interval is defined by
// the received duration, and sends the results to the channel.  The output
// is Flatbuffer serialized bytes of Info.  Any error encountered during
// processing is sent to the error channel; processing will continue.
//
// If an error occurs while opening /proc/meminfo, the error will be sent
// to the errs channel and this func will exit.
//
// To stop processing and exit; send a signal on the done channel.  This
// will cause the function to stop the ticker, close the out channel and
// return.
//
// This func uses a local InfoProfiler.  If an error occurs during the
// creation of the InfoProfiler, it will be sent to errs and exit.
func InfoTicker(interval time.Duration, out chan Info, done chan struct{}, errs chan error) {
	p, err := NewInfoProfiler()
	if err != nil {
		errs <- err
		return
	}
	p.Ticker(interval, out, done, errs)
}

// TickerFlat gathers the meminfo on a ticker, whose interval is defined by
// the received duration, and sends the results to the channel.  The output
// is Flatbuffer serialized bytes of Info.  Any error encountered during
// processing is sent to the error channel; processing will continue.
//
// If an error occurs while opening /proc/meminfo, the error will be sent
// to the errs channel and this func will exit.
//
// To stop processing and exit; send a signal on the done channel.  This
// will cause the function to stop the ticker, close the out channel and
// return.
func (p *InfoProfiler) TickerFlat(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(out)
	// predeclare some vars
	var (
		l, i, pos, nameLen int
		v                  byte
		n                  uint64
		err                error
		line               []byte
	)
	// just reset the bldr at the end of every ticker
	bldr := fb.NewBuilder(0)
	// Lock now because the for loop unlocks to simplify unlock logic when
	// a continue occurs (instead of the tick completing.)
	p.Lock()
	// ticker
Tick:
	for {
		p.Unlock()
		select {
		case <-done:
			return
		case <-ticker.C:
			p.Lock()
			err = p.reset()
			if err != nil {
				errs <- joe.Error{Type: "mem", Op: "seek byte 0: /proc/meminfo", Err: err}
				continue
			}
			flat.InfoStart(bldr)
			flat.InfoAddTimestamp(bldr, time.Now().UTC().UnixNano())
			for l = 0; l < 16; l++ {
				line, err = p.Buf.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					// An error results in sending error message and stop processing of this tick.
					errs <- joe.Error{Type: "mem", Op: "read output bytes", Err: err}
					continue Tick
				}
				if l > 7 && l < 14 {
					continue
				}
				// first grab the key name (everything up to the ':')
				for i, v = range line {
					if v == 0x3A {
						p.val = line[:i]
						break
					}
				}
				nameLen = len(p.val)
				// skip all spaces
				for i, v = range line[pos:] {
					if v != 0x20 {
						pos += i
						break
					}
				}

				// grab the numbers
				for _, v = range line[pos:] {
					if v == 0x20 || v == '\n' {
						break
					}
					p.val = append(p.val, v)
				}
				// any conversion error results in 0
				n, err = helpers.ParseUint(p.val[nameLen:])
				if err != nil {
					errs <- joe.Error{Type: "mem", Op: fmt.Sprintf("convert %s", p.val[:nameLen]), Err: err}
					continue
				}
				v = p.val[0]
				if v == 'M' {
					v = p.val[3]
					if v == 'T' {
						flat.InfoAddMemTotal(bldr, int64(n))
					} else if v == 'F' {
						flat.InfoAddMemFree(bldr, int64(n))
					} else {
						flat.InfoAddMemAvailable(bldr, int64(n))
					}
				} else if v == 'S' {
					v = p.val[4]
					if v == 'C' {
						flat.InfoAddSwapCached(bldr, int64(n))
					} else if v == 'T' {
						flat.InfoAddSwapTotal(bldr, int64(n))
					} else if v == 'F' {
						flat.InfoAddSwapFree(bldr, int64(n))
					}
				} else if v == 'B' {
					flat.InfoAddBuffers(bldr, int64(n))
				} else if v == 'I' {
					flat.InfoAddInactive(bldr, int64(n))
				} else if v == 'C' {
					flat.InfoAddMemAvailable(bldr, int64(n))
				} else if v == 'A' {
					flat.InfoAddInactive(bldr, int64(n))
				}
			}
			bldr.Finish(flat.InfoEnd(bldr))
			inf := bldr.Bytes[bldr.Head():]
			out <- inf
		}
	}
}

// TODO: should InfoTickerFlat use std or have a local proc?
// InfoTickerFlat gathers the meminfo on a ticker, whose interval is defined
// by the received duration, and sends the results to the channel.  The
// output is Flatbuffer serialized bytes of Info.  Any error encountered
// during processing is sent to the error channel; processing will continue.
//
// If an error occurs while opening /proc/meminfo, the error will be sent
// to the errs channel and this func will exit.
//
// To stop processing and exit; send a signal on the done channel.  This
// will cause the function to stop the ticker, close the out channel and
// return.
//
// This func uses a local InfoProfiler.  If an error occurs during the
// creation of the InfoProfiler, it will be sent to errs and exit.
func InfoTickerFlat(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	p, err := NewInfoProfiler()
	if err != nil {
		errs <- err
		return
	}
	p.TickerFlat(interval, out, done, errs)
}

// TickerFlat gathers the meminfo on a ticker, whose interval is defined by
// the received duration, and sends the results to the channel.  The output
// is Flatbuffer serialized bytes of Info.  Any error encountered during
// processing is sent to the error channel; processing will continue.
//
// If an error occurs while opening /proc/meminfo, the error will be sent
// to the errs channel and this func will exit.
//
// To stop processing and exit; send a signal on the done channel.  This
// will cause the function to stop the ticker, close the out channel and
// return.
func (p *InfoProfiler) TickerJSON(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	outCh := make(chan Info)
	defer close(outCh)
	go p.Ticker(interval, outCh, done, errs)
	for {
		select {
		case inf, ok := <-outCh:
			if !ok {
				return
			}
			b, err := json.Marshal(inf)
			if err != nil {
				errs <- err
				continue
			}
			out <- b
		}
	}
}

// InfoTickerJSON gathers the meminfo on a ticker, whose interval is defined
// by the received duration, and sends the results to the channel.  The
// output is the JSON serialized bytes of Info.  Any error encountered
// during processing is sent to the error channel; processing will continue.
//
// If an error occurs while opening /proc/meminfo, the error will be sent
// to the errs channel and this func will exit.
//
// To stop processing and exit; send a signal on the done channel.  This
// will cause the function to stop the ticker, close the out channel and
// return.
//
// This func uses a local InfoProfiler.  If an error occurs during the
// creation of the InfoProfiler, it will be sent to errs and exit.
func InfoTickerJSON(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	p, err := NewInfoProfiler()
	if err != nil {
		errs <- err
		return
	}
	p.TickerJSON(interval, out, done, errs)
}

// InfoProfilerFlat wraps InfoProfiler and provides a builder; enabling reuse.
type InfoProfilerFlat struct {
	InfoProfiler
	bldr *fb.Builder
}

func NewInfoProfilerFlat() (proc *InfoProfilerFlat, err error) {
	f, err := os.Open(procMemInfo)
	if err != nil {
		return nil, err
	}
	return &InfoProfilerFlat{InfoProfiler: InfoProfiler{Proc: joe.Proc{File: f, Buf: bufio.NewReader(f)}, val: make([]byte, 32)}, bldr: fb.NewBuilder(0)}, nil
}

// GetFlat returns the current meminfo as flatbuffer serialized bytes.
func (p *InfoProfilerFlat) GetFlat() ([]byte, error) {
	inf, err := p.InfoProfiler.Get()
	if err != nil {
		return nil, err
	}
	p.bldr.Reset()
	return inf.SerializeFlatBuilder(p.bldr), nil
}

// func (i *InfoFlat) String() string {
// 	return fmt.Sprintf("Timestamp: %v\nMemTotal:\t%d\tMemFree:\t%d\tMemAvailable:\t%d\tActive:\t%d\tInactive:\t%d\nCached:\t\t%d\tBuffers\t:%d\nSwapTotal:\t%d\tSwapCached:\t%d\tSwapFree:\t%d\n", time.Unix(0, i.Timestamp()).UTC(), i.MemTotal(), i.MemFree(), i.MemAvailable(), i.Active(), i.Inactive(), i.Cached(), i.Buffers(), i.SwapTotal(), i.SwapCached(), i.SwapFree())
// }

type Info struct {
	Timestamp    int64 `json:"timestamp"`
	MemTotal     int64 `json:"mem_total"`
	MemFree      int64 `json:"mem_free"`
	MemAvailable int64 `json:"mem_available"`
	Buffers      int64 `json:"buffers"`
	Cached       int64 `json:"cached"`
	SwapCached   int64 `json:"swap_cached"`
	Active       int64 `json:"active"`
	Inactive     int64 `json:"inactive"`
	SwapTotal    int64 `json:"swap_total"`
	SwapFree     int64 `json:"swap_free"`
}

// Serialize serializes the Info using flatbuffers.
func (i *Info) SerializeFlat() []byte {
	bldr := fb.NewBuilder(0)
	return i.SerializeFlatBuilder(bldr)
}

func (i *Info) SerializeFlatBuilder(bldr *fb.Builder) []byte {
	flat.InfoStart(bldr)
	flat.InfoAddTimestamp(bldr, int64(i.Timestamp))
	flat.InfoAddMemTotal(bldr, int64(i.MemTotal))
	flat.InfoAddMemFree(bldr, int64(i.MemFree))
	flat.InfoAddMemAvailable(bldr, int64(i.MemAvailable))
	flat.InfoAddBuffers(bldr, int64(i.Buffers))
	flat.InfoAddCached(bldr, int64(i.Cached))
	flat.InfoAddSwapCached(bldr, int64(i.SwapCached))
	flat.InfoAddActive(bldr, int64(i.Active))
	flat.InfoAddInactive(bldr, int64(i.Inactive))
	flat.InfoAddSwapTotal(bldr, int64(i.SwapTotal))
	flat.InfoAddSwapFree(bldr, int64(i.SwapFree))
	bldr.Finish(flat.InfoEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

// Marshal Info as JSON
func (i *Info) SerializeJSON() ([]byte, error) {
	return json.Marshal(i)
}

// DeserializeInfoFlat deserializes bytes serialized with Flatbuffers from
// InfoFlat into *Info.
func DeserializeInfoFlat(p []byte) *Info {
	infoFlat := flat.GetRootAsInfo(p, 0)
	info := &Info{}
	info.Timestamp = infoFlat.Timestamp()
	info.MemTotal = infoFlat.MemTotal()
	info.MemFree = infoFlat.MemFree()
	info.MemAvailable = infoFlat.MemAvailable()
	info.Buffers = infoFlat.Buffers()
	info.Cached = infoFlat.Cached()
	info.SwapCached = infoFlat.SwapCached()
	info.Active = infoFlat.Active()
	info.Inactive = infoFlat.Inactive()
	info.SwapTotal = infoFlat.SwapTotal()
	info.SwapFree = infoFlat.SwapFree()
	return info
}

// UnmarshalInfoJSON unmarshals JSON into *Info.
func UnmarshalInfoJSON(p []byte) (*Info, error) {
	info := &Info{}
	err := json.Unmarshal(p, info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (i *Info) String() string {
	return fmt.Sprintf("Timestamp: %v\nMemTotal:\t%d\tMemFree:\t%d\tMemAvailable:\t%d\tActive:\t%d\tInactive:\t%d\nCached:\t\t%d\tBuffers\t:%d\nSwapTotal:\t%d\tSwapCached:\t%d\tSwapFree:\t%d\n", time.Unix(0, i.Timestamp).UTC(), i.MemTotal, i.MemFree, i.MemAvailable, i.Active, i.Inactive, i.Cached, i.Buffers, i.SwapTotal, i.SwapCached, i.SwapFree)
}

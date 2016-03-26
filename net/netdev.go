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

// Package net gets and processes /proc/net/dev, returning the infoFlat in the
// appropriate format.
package net

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/net/flat"
)

type Info struct {
	Timestamp  int64       `json:"timestamp"`
	Interfaces []Interface `json:"interfaces"`
}

// Interface: contains information for a given network interface.
type Interface struct {
	Name        string `json:"name"`
	RBytes      int64  `json:"receive_bytes"`
	RPackets    int64  `json:"receive_packets"`
	RErrs       int64  `json:"receive_errs"`
	RDrop       int64  `json:"receive_drop"`
	RFIFO       int64  `json:"recieve_fifo"`
	RFrame      int64  `json:"receive_frame"`
	RCompressed int64  `json:"receive_compressed"`
	RMulticast  int64  `json:"receive_multicast"`
	TBytes      int64  `json:"transmit_bytes"`
	TPackets    int64  `json:"transmit_packets"`
	TErrs       int64  `json:"transmit_errs"`
	TDrop       int64  `json:"transmit_drop"`
	TFIFO       int64  `json:"transmit_fifo"`
	TColls      int64  `json:"transmit_colls"`
	TCarrier    int64  `json:"transmit_carrier"`
	TCompressed int64  `json:"transmit_compressed"`
}

// SerializeFlat serializes the Info using flatbuffers.
func (inf *Info) SerializeFlat() []byte {
	bldr := fb.NewBuilder(0)
	return inf.SerializeFlatBuilder(bldr)
}

// SerializeFlatBuilder serializes the Info using flatbuffers.
func (inf *Info) SerializeFlatBuilder(bldr *fb.Builder) []byte {
	ifaces := make([]fb.UOffsetT, len(inf.Interfaces))
	names := make([]fb.UOffsetT, len(inf.Interfaces))
	for i := 0; i < len(inf.Interfaces); i++ {
		names[i] = bldr.CreateString(inf.Interfaces[i].Name)
	}
	for i := 0; i < len(inf.Interfaces); i++ {
		flat.InterfaceStart(bldr)
		flat.InterfaceAddName(bldr, names[i])
		flat.InterfaceAddRBytes(bldr, inf.Interfaces[i].RBytes)
		flat.InterfaceAddRPackets(bldr, inf.Interfaces[i].RPackets)
		flat.InterfaceAddRErrs(bldr, inf.Interfaces[i].RErrs)
		flat.InterfaceAddRDrop(bldr, inf.Interfaces[i].RDrop)
		flat.InterfaceAddRFIFO(bldr, inf.Interfaces[i].RFIFO)
		flat.InterfaceAddRFrame(bldr, inf.Interfaces[i].RFrame)
		flat.InterfaceAddRCompressed(bldr, inf.Interfaces[i].RCompressed)
		flat.InterfaceAddRMulticast(bldr, inf.Interfaces[i].RMulticast)
		flat.InterfaceAddTBytes(bldr, inf.Interfaces[i].TBytes)
		flat.InterfaceAddTPackets(bldr, inf.Interfaces[i].TPackets)
		flat.InterfaceAddTErrs(bldr, inf.Interfaces[i].TErrs)
		flat.InterfaceAddTDrop(bldr, inf.Interfaces[i].TDrop)
		flat.InterfaceAddTFIFO(bldr, inf.Interfaces[i].TFIFO)
		flat.InterfaceAddTColls(bldr, inf.Interfaces[i].TColls)
		flat.InterfaceAddTCarrier(bldr, inf.Interfaces[i].TCarrier)
		flat.InterfaceAddTCompressed(bldr, inf.Interfaces[i].TCompressed)
		ifaces[i] = flat.InterfaceEnd(bldr)
	}
	flat.InfoStartInterfacesVector(bldr, len(ifaces))
	for i := len(inf.Interfaces) - 1; i >= 0; i-- {
		bldr.PrependUOffsetT(ifaces[i])
	}
	ifacesV := bldr.EndVector(len(ifaces))
	flat.InfoStart(bldr)
	flat.InfoAddTimestamp(bldr, inf.Timestamp)
	flat.InfoAddInterfaces(bldr, ifacesV)
	bldr.Finish(flat.InfoEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

// DeserializeInfoFlat deserializes Flatbuffer serialized bytes as *Info.
func DeserializeInfoFlat(p []byte) *Info {
	infoFlat := flat.GetRootAsInfo(p, 0)
	// get the # of interfaces
	iLen := infoFlat.InterfacesLength()
	info := &Info{Timestamp: infoFlat.Timestamp(), Interfaces: make([]Interface, iLen)}
	iFace := &flat.Interface{}
	iface := Interface{}
	for i := 0; i < iLen; i++ {
		if infoFlat.Interfaces(iFace, i) {
			iface.Name = string(iFace.Name())
			iface.RBytes = iFace.RBytes()
			iface.RPackets = iFace.RPackets()
			iface.RErrs = iFace.RErrs()
			iface.RDrop = iFace.RDrop()
			iface.RFIFO = iFace.RFIFO()
			iface.RFrame = iFace.RFrame()
			iface.RCompressed = iFace.RCompressed()
			iface.RMulticast = iFace.RMulticast()
			iface.TBytes = iFace.TBytes()
			iface.TPackets = iFace.TPackets()
			iface.TErrs = iFace.TErrs()
			iface.TDrop = iFace.TDrop()
			iface.TFIFO = iFace.TFIFO()
			iface.TColls = iFace.TColls()
			iface.TCarrier = iFace.TCarrier()
			iface.TCompressed = iFace.TCompressed()
		}
		info.Interfaces[i] = iface
	}
	return info
}

func (inf Info) String() string {
	var buf bytes.Buffer
	buf.WriteString(time.Unix(0, inf.Timestamp).UTC().String())
	buf.WriteByte('\n')
	for i := 0; i < len(inf.Interfaces); i++ {
		buf.WriteString(joe.Column(8, inf.Interfaces[i].Name))
		buf.WriteString(joe.Int64Column(22, inf.Interfaces[i].RBytes))
		buf.WriteString(joe.Int64Column(22, inf.Interfaces[i].RPackets))
		buf.WriteString(joe.Int64Column(22, inf.Interfaces[i].RErrs))
		buf.WriteString(joe.Int64Column(22, inf.Interfaces[i].RDrop))
		buf.WriteString(joe.Int64Column(22, inf.Interfaces[i].RFIFO))
		buf.WriteString(joe.Int64Column(22, inf.Interfaces[i].RFrame))
		buf.WriteString(joe.Int64Column(22, inf.Interfaces[i].RCompressed))
		buf.WriteString(joe.Int64Column(22, inf.Interfaces[i].RMulticast))
		buf.WriteString(joe.Int64Column(22, inf.Interfaces[i].TBytes))
		buf.WriteString(joe.Int64Column(22, inf.Interfaces[i].TPackets))
		buf.WriteString(joe.Int64Column(22, inf.Interfaces[i].TErrs))
		buf.WriteString(joe.Int64Column(22, inf.Interfaces[i].TDrop))
		buf.WriteString(joe.Int64Column(22, inf.Interfaces[i].TFIFO))
		buf.WriteString(joe.Int64Column(22, inf.Interfaces[i].TColls))
		buf.WriteString(joe.Int64Column(22, inf.Interfaces[i].TCarrier))
		buf.WriteString(joe.Int64Column(22, inf.Interfaces[i].TCompressed))
		buf.WriteByte('\n')
	}
	return buf.String()
}

// GetInfo returns some of the results of /proc/net/dev
func GetInfo() (*Info, error) {
	var l, i, pos, fieldNum, fieldVal int
	var v byte
	t := time.Now().UTC().UnixNano()
	f, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	// there's always at least 2 interfaces (I think)
	inf := &Info{Timestamp: t, Interfaces: make([]Interface, 0, 2)}
	val := make([]byte, 0, 32)
	for {
		line, err := buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading output bytes: %s", err)
		}
		l++
		if l < 3 {
			continue
		}
		var iInfo Interface

		// first grab the interface name (everything up to the ':')
		for i, v = range line {
			if v == 0x3A {
				pos = i + 1
				break
			}
			val = append(val, v)
		}
		iInfo.Name = string(val[:])
		val = val[:0]
		fieldNum = 0
		// process the rest of the line
		for {
			fieldNum++
			// skip all spaces
			for i, v = range line[pos:] {
				if v != 0x20 {
					pos += i
					break
				}
			}

			// grab the numbers
			for i, v = range line[pos:] {
				if v == 0x20 || v == '\n' {
					pos += i
					break
				}
				val = append(val, v)
			}
			// any conversion error results in 0
			fieldVal, err = strconv.Atoi(string(val[:]))
			if err != nil {
				return nil, fmt.Errorf("%s: %s", iInfo.Name, err)
			}
			val = val[:0]
			if fieldNum == 1 {
				iInfo.RBytes = int64(fieldVal)
				continue
			}
			if fieldNum == 2 {
				iInfo.RPackets = int64(fieldVal)
				continue
			}
			if fieldNum == 3 {
				iInfo.RErrs = int64(fieldVal)
				continue
			}
			if fieldNum == 4 {
				iInfo.RDrop = int64(fieldVal)
				continue
			}
			if fieldNum == 5 {
				iInfo.RFIFO = int64(fieldVal)
				continue
			}
			if fieldNum == 6 {
				iInfo.RFrame = int64(fieldVal)
				continue
			}
			if fieldNum == 7 {
				iInfo.RCompressed = int64(fieldVal)
				continue
			}
			if fieldNum == 8 {
				iInfo.RMulticast = int64(fieldVal)
				continue
			}
			if fieldNum == 9 {
				iInfo.TBytes = int64(fieldVal)
				continue
			}
			if fieldNum == 10 {
				iInfo.TPackets = int64(fieldVal)
				continue
			}
			if fieldNum == 11 {
				iInfo.TErrs = int64(fieldVal)
				continue
			}
			if fieldNum == 12 {
				iInfo.TDrop = int64(fieldVal)
				continue
			}
			if fieldNum == 13 {
				iInfo.TFIFO = int64(fieldVal)
				continue
			}
			if fieldNum == 14 {
				iInfo.TColls = int64(fieldVal)
				continue
			}
			if fieldNum == 15 {
				iInfo.TCarrier = int64(fieldVal)
				continue
			}
			if fieldNum == 16 {
				iInfo.TCompressed = int64(fieldVal)
				break
			}
		}
		inf.Interfaces = append(inf.Interfaces, iInfo)
	}
	return inf, nil
}

// GetInfoFlat returns the current meminfo as flatbuffer serialized bytes.
func GetInfoFlat() ([]byte, error) {
	inf, err := GetInfo()
	if err != nil {
		return nil, err
	}
	return inf.SerializeFlat(), nil
}

// InfoTickerFlat gathers the net/dev infoFlat on a ticker, whose interval is defined
// by the received duration, and sends the results to the channel.  The output
// is Flatbuffers serialized Info.  Any error encountered during processing
// is sent to the error channel.  Processing will continue
//
// Sending a struct{}{} to the done channel will result in function exit.
//
// This pre-allocates the builder and everything other than the []byte that
// gets sent to the out channel to reduce allocations, as this is expected
// to be both a frequent and a long-running process.
func InfoTickerFlat(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(out)
	// predeclare some vars
	var l, i, pos, fieldNum, fieldVal int
	var v byte
	var t int64
	var iInfo Interface
	// premake some temp slices
	val := make([]byte, 0, 32)
	// just reset the bldr at the end of every ticker
	bldr := fb.NewBuilder(0)
	// Some hopes to jump through to ensure we don't get a ErrBufferFull.
	var bs []byte
	tmp := bytes.NewBuffer(bs)
	// If there are a lot of interfaces, and a high volume of bytes being
	// processed this could be inadequate
	// TODO: revisit this sizing/maybe make it configurable?
	buf := bufio.NewReaderSize(tmp, 4096)
	tmp = nil
	inf := &Info{Interfaces: make([]Interface, 0, 4)}
	// ticker
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			t = time.Now().UTC().UnixNano()
			f, err := os.Open("/proc/net/dev")
			if err != nil {
				return
			}
			defer f.Close()
			buf.Reset(f)
			inf.Interfaces = inf.Interfaces[:0]
			// there's always at least 2 interfaces (I think)
			inf.Timestamp = t
			for {
				line, err := buf.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					errs <- fmt.Errorf("/proc/mem/dev: read output bytes: %s", err)
					break
				}
				l++
				if l < 3 {
					continue
				}

				// first grab the interface name (everything up to the ':')
				for i, v = range line {
					if v == 0x3A {
						pos = i + 1
						break
					}
					val = append(val, v)
				}
				iInfo.Name = string(val[:])
				val = val[:0]
				fieldNum = 0
				// process the rest of the line
				for {
					fieldNum++
					// skip all spaces
					for i, v = range line[pos:] {
						if v != 0x20 {
							pos += i
							break
						}
					}

					// grab the numbers
					for i, v = range line[pos:] {
						if v == 0x20 || v == '\n' {
							pos += i
							break
						}
						val = append(val, v)
					}
					// any conversion error results in 0
					fieldVal, err = strconv.Atoi(string(val[:]))
					if err != nil {
						errs <- fmt.Errorf("/proc/net/dev ticker: %s: %s", iInfo.Name, err)
						continue
					}
					val = val[:0]
					if fieldNum == 1 {
						iInfo.RBytes = int64(fieldVal)
						continue
					}
					if fieldNum == 2 {
						iInfo.RPackets = int64(fieldVal)
						continue
					}
					if fieldNum == 3 {
						iInfo.RErrs = int64(fieldVal)
						continue
					}
					if fieldNum == 4 {
						iInfo.RDrop = int64(fieldVal)
						continue
					}
					if fieldNum == 5 {
						iInfo.RFIFO = int64(fieldVal)
						continue
					}
					if fieldNum == 6 {
						iInfo.RFrame = int64(fieldVal)
						continue
					}
					if fieldNum == 7 {
						iInfo.RCompressed = int64(fieldVal)
						continue
					}
					if fieldNum == 8 {
						iInfo.RMulticast = int64(fieldVal)
						continue
					}
					if fieldNum == 9 {
						iInfo.TBytes = int64(fieldVal)
						continue
					}
					if fieldNum == 10 {
						iInfo.TPackets = int64(fieldVal)
						continue
					}
					if fieldNum == 11 {
						iInfo.TErrs = int64(fieldVal)
						continue
					}
					if fieldNum == 12 {
						iInfo.TDrop = int64(fieldVal)
						continue
					}
					if fieldNum == 13 {
						iInfo.TFIFO = int64(fieldVal)
						continue
					}
					if fieldNum == 14 {
						iInfo.TColls = int64(fieldVal)
						continue
					}
					if fieldNum == 15 {
						iInfo.TCarrier = int64(fieldVal)
						continue
					}
					if fieldNum == 16 {
						iInfo.TCompressed = int64(fieldVal)
						break
					}
				}
				inf.Interfaces = append(inf.Interfaces, iInfo)
			}
			f.Close()
			bldr.Reset()
			infoFlat := inf.SerializeFlatBuilder(bldr)
			out <- infoFlat
			l = 0
		}
	}
}

// Usage holds the difference between network IO snapshots.
type Usage Info

func (u *Usage) SerializeFlat() []byte {
	bldr := fb.NewBuilder(0)
	return u.SerializeFlatBuilder(bldr)
}

func (u *Usage) SerializeFlatBuilder(bldr *fb.Builder) []byte {
	inf := Info(*u)
	return inf.SerializeFlatBuilder(bldr)
}

func DeserializeUsageFlat(p []byte) *Usage {
	u := Usage(*(DeserializeInfoFlat(p)))
	return &u
}

// Usage gets the number of recieve/transmit information for the given
func GetUsage(t time.Duration) (Usage, error) {
	snap1, err := GetInfo()
	if err != nil {
		return Usage{}, err
	}
	time.Sleep(t)
	snap2, err := GetInfo()
	if err != nil {
		return Usage{}, err
	}

	return calculateUsage(snap1, snap2), nil
}

func UsageTicker(period time.Duration, out chan Usage, done chan struct{}, errs chan error) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	defer close(out)
	// predeclare some vars
	var i, l, pos, fieldNum, fieldVal int
	var v byte
	var iInfo Interface
	val := make([]byte, 0, 32)
	prior := &Info{}
	// first get Info as the baseline
	cur, err := GetInfo()
	if err != nil {
		errs <- err
	}
	// Some hoops to jump through to ensure we don't get a ErrBufferFull.
	// TODO: should bufio.Scanner be used instead?
	var bs []byte
	tmp := bytes.NewBuffer(bs)
	buf := bufio.NewReaderSize(tmp, 4096)
	tmp = nil
	// ticker
tick:
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			prior.Timestamp = cur.Timestamp
			if len(prior.Interfaces) != len(cur.Interfaces) {
				prior.Interfaces = make([]Interface, len(cur.Interfaces))
			}
			copy(prior.Interfaces, cur.Interfaces)
			cur.Timestamp = time.Now().UTC().UnixNano()
			f, err := os.Open("/proc/net/dev")
			if err != nil {
				errs <- joe.Error{Type: "net", Op: "usage ticker", Err: err}
				continue tick
			}
			defer f.Close()
			cur.Interfaces = cur.Interfaces[:0]
			buf.Reset(f)
			// read each line until eof
			for {
				line, err := buf.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					errs <- fmt.Errorf("/proc/mem/dev: read output bytes: %s", err)
					break
				}
				l++
				if l < 3 {
					continue
				}

				// skip leading spaces
				for i, v = range line {
					if v != 0x20 {
						pos = i
						break
					}
				}
				// first grab the interface name (everything up to the ':')
				for i, v = range line[pos:] {
					if v == 0x3A {
						iInfo.Name = string(line[pos : pos+i])
						pos += i + 1
						break
					}
				}
				fieldNum = 0
				// process the rest of the line
				for {
					fieldNum++
					// skip all spaces
					for i, v = range line[pos:] {
						if v != 0x20 {
							pos += i
							break
						}
					}

					// grab the numbers
					for i, v = range line[pos:] {
						if v == 0x20 || v == '\n' {
							val = line[pos : pos+i]
							pos += i + 1
							break
						}
					}
					// any conversion error results in 0
					fieldVal, err = strconv.Atoi(string(val[:]))
					if err != nil {
						errs <- fmt.Errorf("/proc/net/dev ticker: %s: %s", iInfo.Name, err)
						continue
					}
					val = val[:0]
					if fieldNum == 1 {
						iInfo.RBytes = int64(fieldVal)
						continue
					}
					if fieldNum == 2 {
						iInfo.RPackets = int64(fieldVal)
						continue
					}
					if fieldNum == 3 {
						iInfo.RErrs = int64(fieldVal)
						continue
					}
					if fieldNum == 4 {
						iInfo.RDrop = int64(fieldVal)
						continue
					}
					if fieldNum == 5 {
						iInfo.RFIFO = int64(fieldVal)
						continue
					}
					if fieldNum == 6 {
						iInfo.RFrame = int64(fieldVal)
						continue
					}
					if fieldNum == 7 {
						iInfo.RCompressed = int64(fieldVal)
						continue
					}
					if fieldNum == 8 {
						iInfo.RMulticast = int64(fieldVal)
						continue
					}
					if fieldNum == 9 {
						iInfo.TBytes = int64(fieldVal)
						continue
					}
					if fieldNum == 10 {
						iInfo.TPackets = int64(fieldVal)
						continue
					}
					if fieldNum == 11 {
						iInfo.TErrs = int64(fieldVal)
						continue
					}
					if fieldNum == 12 {
						iInfo.TDrop = int64(fieldVal)
						continue
					}
					if fieldNum == 13 {
						iInfo.TFIFO = int64(fieldVal)
						continue
					}
					if fieldNum == 14 {
						iInfo.TColls = int64(fieldVal)
						continue
					}
					if fieldNum == 15 {
						iInfo.TCarrier = int64(fieldVal)
						continue
					}
					if fieldNum == 16 {
						iInfo.TCompressed = int64(fieldVal)
						break
					}
				}
				cur.Interfaces = append(cur.Interfaces, iInfo)
			}
			f.Close()
			out <- calculateUsage(prior, cur)
			l = 0
		}
	}
}

func calculateUsage(prior, cur *Info) Usage {
	u := Usage{Timestamp: cur.Timestamp, Interfaces: make([]Interface, len(cur.Interfaces))}
	for i := 0; i < len(cur.Interfaces); i++ {
		u.Interfaces[i].Name = cur.Interfaces[i].Name
		u.Interfaces[i].RBytes = cur.Interfaces[i].RBytes - prior.Interfaces[i].RBytes
		u.Interfaces[i].RPackets = cur.Interfaces[i].RPackets - prior.Interfaces[i].RPackets
		u.Interfaces[i].RErrs = cur.Interfaces[i].RErrs - prior.Interfaces[i].RErrs
		u.Interfaces[i].RDrop = cur.Interfaces[i].RDrop - prior.Interfaces[i].RDrop
		u.Interfaces[i].RFIFO = cur.Interfaces[i].RFIFO - prior.Interfaces[i].RFIFO
		u.Interfaces[i].RFrame = cur.Interfaces[i].RFrame - prior.Interfaces[i].RFrame
		u.Interfaces[i].RCompressed = cur.Interfaces[i].RCompressed - prior.Interfaces[i].RCompressed
		u.Interfaces[i].RMulticast = cur.Interfaces[i].RMulticast - prior.Interfaces[i].RMulticast
		u.Interfaces[i].TBytes = cur.Interfaces[i].TBytes - prior.Interfaces[i].TBytes
		u.Interfaces[i].TPackets = cur.Interfaces[i].TPackets - prior.Interfaces[i].TPackets
		u.Interfaces[i].TErrs = cur.Interfaces[i].TErrs - prior.Interfaces[i].TErrs
		u.Interfaces[i].TDrop = cur.Interfaces[i].TDrop - prior.Interfaces[i].TDrop
		u.Interfaces[i].TFIFO = cur.Interfaces[i].TFIFO - prior.Interfaces[i].TFIFO
		u.Interfaces[i].TColls = cur.Interfaces[i].TColls - prior.Interfaces[i].TColls
		u.Interfaces[i].TCarrier = cur.Interfaces[i].TCarrier - prior.Interfaces[i].TCarrier
		u.Interfaces[i].TCompressed = cur.Interfaces[i].TCompressed - prior.Interfaces[i].TCompressed
	}
	return u
}

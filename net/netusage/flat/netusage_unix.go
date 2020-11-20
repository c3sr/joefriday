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
// TimeDelta field. Instead of returning a Go struct, it returns Flatbuffer
// serialized bytes. A function to deserialize the Flatbuffer serialized bytes
// into a structs.DevUsage struct is provided.
//
// Note: the package name is netusage and not the final element of the import
// path (flat)
package netusage

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/c3sr/joefriday"
	"github.com/c3sr/joefriday/net/structs"
	"github.com/c3sr/joefriday/net/structs/flat"
	usage "github.com/c3sr/joefriday/net/netusage"
)

// Profiler is used to process the network device usage.
type Profiler struct {
	*usage.Profiler
	*fb.Builder
}

// Returns an initialized Profiler; ready to use. Upon creation, a
// /proc/net/dev snapshot is taken so that any Get() will return valid
// information.
func NewProfiler() (prof *Profiler, err error) {
	p, err := usage.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current network device usage as Flatbuffer serialized bytes.
// Calculating usage requires two snapshots. This func gets the current
// snapshot of /proc/net/dev and calculates the difference between that and the
// prior snapshot. The current snapshot is stored for use as the prior snapshot
// on the next Get call. If ongoing usage information is desired, the Ticker
// should be used; it's better suited for ongoing usage information..
func (prof *Profiler) Get() (p []byte, err error) {
	u, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(u), nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current network device usage as Flatbuffer serialized bytes
// using the package's global Profiler. The profiler is lazily instantiated.
// If the profiler doesn't already exist, the first usage information will not
// be useful due to the minimal time elapsing between the initial and second
// snapshots used for usage calculations; the results of the first call should
// be discarded.
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

// Serialize network device usage using Flatbuffers.
func (prof *Profiler) Serialize(u *structs.DevUsage) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	devs := make([]fb.UOffsetT, len(u.Device))
	names := make([]fb.UOffsetT, len(u.Device))
	for i := 0; i < len(u.Device); i++ {
		names[i] = prof.Builder.CreateString(u.Device[i].Name)
	}
	for i := 0; i < len(u.Device); i++ {
		flat.DeviceStart(prof.Builder)
		flat.DeviceAddName(prof.Builder, names[i])
		flat.DeviceAddRBytes(prof.Builder, u.Device[i].RBytes)
		flat.DeviceAddRPackets(prof.Builder, u.Device[i].RPackets)
		flat.DeviceAddRErrs(prof.Builder, u.Device[i].RErrs)
		flat.DeviceAddRDrop(prof.Builder, u.Device[i].RDrop)
		flat.DeviceAddRFIFO(prof.Builder, u.Device[i].RFIFO)
		flat.DeviceAddRFrame(prof.Builder, u.Device[i].RFrame)
		flat.DeviceAddRCompressed(prof.Builder, u.Device[i].RCompressed)
		flat.DeviceAddRMulticast(prof.Builder, u.Device[i].RMulticast)
		flat.DeviceAddTBytes(prof.Builder, u.Device[i].TBytes)
		flat.DeviceAddTPackets(prof.Builder, u.Device[i].TPackets)
		flat.DeviceAddTErrs(prof.Builder, u.Device[i].TErrs)
		flat.DeviceAddTDrop(prof.Builder, u.Device[i].TDrop)
		flat.DeviceAddTFIFO(prof.Builder, u.Device[i].TFIFO)
		flat.DeviceAddTColls(prof.Builder, u.Device[i].TColls)
		flat.DeviceAddTCarrier(prof.Builder, u.Device[i].TCarrier)
		flat.DeviceAddTCompressed(prof.Builder, u.Device[i].TCompressed)
		devs[i] = flat.DeviceEnd(prof.Builder)
	}
	flat.DevUsageStartDeviceVector(prof.Builder, len(devs))
	for i := len(u.Device) - 1; i >= 0; i-- {
		prof.Builder.PrependUOffsetT(devs[i])
	}
	devsV := prof.Builder.EndVector(len(devs))
	flat.DevUsageStart(prof.Builder)
	flat.DevUsageAddTimestamp(prof.Builder, u.Timestamp)
	flat.DevUsageAddTimeDelta(prof.Builder, u.TimeDelta)
	flat.DevUsageAddDevice(prof.Builder, devsV)
	prof.Builder.Finish(flat.DevUsageEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize network device usage using Flatbuffers with the package's global
// Profiler.
func Serialize(u *structs.DevUsage) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(u), nil
}

// Deserialize deserializes Flatbuffer serialized bytes as structs.DevUsage.
func Deserialize(p []byte) *structs.DevUsage {
	uFlat := flat.GetRootAsDevUsage(p, 0)
	// get the # of interfaces
	iLen := uFlat.DeviceLength()
	u := &structs.DevUsage{
		Timestamp:  uFlat.Timestamp(),
		TimeDelta:  uFlat.TimeDelta(),
		Device: make([]structs.Device, iLen),
	}
	fDev := &flat.Device{}
	sDev := structs.Device{}
	for i := 0; i < iLen; i++ {
		if uFlat.Device(fDev, i) {
			sDev.Name = string(fDev.Name())
			sDev.RBytes = fDev.RBytes()
			sDev.RPackets = fDev.RPackets()
			sDev.RErrs = fDev.RErrs()
			sDev.RDrop = fDev.RDrop()
			sDev.RFIFO = fDev.RFIFO()
			sDev.RFrame = fDev.RFrame()
			sDev.RCompressed = fDev.RCompressed()
			sDev.RMulticast = fDev.RMulticast()
			sDev.TBytes = fDev.TBytes()
			sDev.TPackets = fDev.TPackets()
			sDev.TErrs = fDev.TErrs()
			sDev.TDrop = fDev.TDrop()
			sDev.TFIFO = fDev.TFIFO()
			sDev.TColls = fDev.TColls()
			sDev.TCarrier = fDev.TCarrier()
			sDev.TCompressed = fDev.TCompressed()
		}
		u.Device[i] = sDev
	}
	return u
}

// Ticker delivers the system's net devices usage at intervals.
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
	for {
		select {
		case <-t.Done:
			return
		case <-t.C:
			p, err := t.Get()
			if err != nil {
				t.Errs <- err
				continue
			}
			t.Data <- p
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}

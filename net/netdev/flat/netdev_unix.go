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

// Package netdev gets the system's network device information: /proc/net/dev.
// Instead of returning a Go struct, it returns Flatbuffer serialized bytes. A
// function to deserialize the Flatbuffer serialized bytes into a
// structs.DevInfo struct is provided.
//
// Note: the package name is netdev and not the final element of the import
// path (flat). 
package netdev

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/rai-project/joefriday"
	dev "github.com/rai-project/joefriday/net/netdev"
	"github.com/rai-project/joefriday/net/structs"
	"github.com/rai-project/joefriday/net/structs/flat"
)

// Profiler is used to process the network device information as Flatbuffer
// serialized bytes using the /proc/net/dev file.
type Profiler struct {
	*dev.Profiler
	*fb.Builder
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	p, err := dev.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current network device information as Flatbuffer serialized
// bytes.
func (prof *Profiler) Get() ([]byte, error) {
	inf, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf), nil
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current network device information as Flatbuffer serialized
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

// Serialize network device information as Flatbuffer serialized bytes.
func (prof *Profiler) Serialize(inf *structs.DevInfo) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	devs := make([]fb.UOffsetT, len(inf.Device))
	names := make([]fb.UOffsetT, len(inf.Device))
	for i := 0; i < len(inf.Device); i++ {
		names[i] = prof.Builder.CreateString(inf.Device[i].Name)
	}
	for i := 0; i < len(inf.Device); i++ {
		flat.DeviceStart(prof.Builder)
		flat.DeviceAddName(prof.Builder, names[i])
		flat.DeviceAddRBytes(prof.Builder, inf.Device[i].RBytes)
		flat.DeviceAddRPackets(prof.Builder, inf.Device[i].RPackets)
		flat.DeviceAddRErrs(prof.Builder, inf.Device[i].RErrs)
		flat.DeviceAddRDrop(prof.Builder, inf.Device[i].RDrop)
		flat.DeviceAddRFIFO(prof.Builder, inf.Device[i].RFIFO)
		flat.DeviceAddRFrame(prof.Builder, inf.Device[i].RFrame)
		flat.DeviceAddRCompressed(prof.Builder, inf.Device[i].RCompressed)
		flat.DeviceAddRMulticast(prof.Builder, inf.Device[i].RMulticast)
		flat.DeviceAddTBytes(prof.Builder, inf.Device[i].TBytes)
		flat.DeviceAddTPackets(prof.Builder, inf.Device[i].TPackets)
		flat.DeviceAddTErrs(prof.Builder, inf.Device[i].TErrs)
		flat.DeviceAddTDrop(prof.Builder, inf.Device[i].TDrop)
		flat.DeviceAddTFIFO(prof.Builder, inf.Device[i].TFIFO)
		flat.DeviceAddTColls(prof.Builder, inf.Device[i].TColls)
		flat.DeviceAddTCarrier(prof.Builder, inf.Device[i].TCarrier)
		flat.DeviceAddTCompressed(prof.Builder, inf.Device[i].TCompressed)
		devs[i] = flat.DeviceEnd(prof.Builder)
	}
	flat.DevInfoStartDeviceVector(prof.Builder, len(devs))
	for i := len(inf.Device) - 1; i >= 0; i-- {
		prof.Builder.PrependUOffsetT(devs[i])
	}
	devsV := prof.Builder.EndVector(len(devs))
	flat.DevInfoStart(prof.Builder)
	flat.DevInfoAddTimestamp(prof.Builder, inf.Timestamp)
	flat.DevInfoAddDevice(prof.Builder, devsV)
	prof.Builder.Finish(flat.DevInfoEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize network device information as Flatbuffer serialized bytes using the
// package's global Profiler.
func Serialize(inf *structs.DevInfo) (p []byte, err error) {
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

// Deserialize takes some Flatbuffer serialized bytes and deserializes them as
// structs.DevInfo.
func Deserialize(p []byte) *structs.DevInfo {
	devInfo := flat.GetRootAsDevInfo(p, 0)
	// get the # of interfaces
	dLen := devInfo.DeviceLength()
	info := &structs.DevInfo{Timestamp: devInfo.Timestamp(), Device: make([]structs.Device, dLen)}
	fDev := &flat.Device{}
	sDev := structs.Device{}
	for i := 0; i < dLen; i++ {
		if devInfo.Device(fDev, i) {
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
		info.Device[i] = sDev
	}
	return info
}

// Ticker delivers the system's network device information at intervals.
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

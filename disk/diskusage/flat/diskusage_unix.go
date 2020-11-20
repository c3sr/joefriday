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

// Package diskusage calculates IO usage of the block devices. Usage is
// calculated by taking the difference between two snapshots of IO statistics
// of the block devices, /proc/diskstats. The time elapsed between the two
// snapshots is stored in the TimeDelta field. Instead of returning a Go
// struct, it returns Flatbuffer serialized bytes. A function to deserialize
// the Flatbuffer serialized bytes into a struct.DiskUsage struct is provided.
//
// Note: the package name is diskusage and not the final element of the import
// path (flat). 
package diskusage

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/c3sr/joefriday"
	"github.com/c3sr/joefriday/disk/structs"
	"github.com/c3sr/joefriday/disk/structs/flat"
	usage "github.com/c3sr/joefriday/disk/diskusage"
)

// Profiler is used to process IO usage of the block devices using Flatbuffers.
type Profiler struct {
	*usage.Profiler
	*fb.Builder
}

// Returns an initialized Profiler; ready to use. Upon creation, a
// /proc/diskstats snapshot is taken so that any Get() will return valid
// information
func NewProfiler() (prof *Profiler, err error) {
	p, err := usage.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current IO usage of the block devices as Flatbuffer
// serialized bytes. Calculating usage requires two snapshots. This func gets
// the current snapshot of /proc/diskstats and calculates the difference
// between that and the prior snapshot. The current snapshot is stored for use
// as the prior snapshot on the next Get call. If ongoing usage information is
// desired, the Ticker should be used; it's better suited for ongoing usage
// information.
func (prof *Profiler) Get() ([]byte, error) {
	u, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(u), nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current IO usage of the block devices as Flatbuffer
// serialized bytes using the package's global Profiler. The Profiler is
// instantiated lazily. If the profiler doesn't already exist, the first
// utilization information will not be useful due to minimal time elapsing
// between the initial and second snapshots used for utilization calculations;
// the results of the first call should be discarded.
func Get() (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	} else {
		std.Builder.Reset()
	}

	return std.Get()
}

// Serialize IO usage of the block devices using Flatbuffers.
func (prof *Profiler) Serialize(u *structs.DiskUsage) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	devF := make([]fb.UOffsetT, len(u.Device))
	names := make([]fb.UOffsetT, len(u.Device))
	for i := 0; i < len(names); i++ {
		names[i] = prof.Builder.CreateString(u.Device[i].Name)
	}
	for i := 0; i < len(devF); i++ {
		flat.DeviceStart(prof.Builder)
		flat.DeviceAddMajor(prof.Builder, u.Device[i].Major)
		flat.DeviceAddMinor(prof.Builder, u.Device[i].Minor)
		flat.DeviceAddName(prof.Builder, names[i])
		flat.DeviceAddReadsCompleted(prof.Builder, u.Device[i].ReadsCompleted)
		flat.DeviceAddReadsMerged(prof.Builder, u.Device[i].ReadsMerged)
		flat.DeviceAddReadSectors(prof.Builder, u.Device[i].ReadSectors)
		flat.DeviceAddReadingTime(prof.Builder, u.Device[i].ReadingTime)
		flat.DeviceAddWritesCompleted(prof.Builder, u.Device[i].WritesCompleted)
		flat.DeviceAddWritesMerged(prof.Builder, u.Device[i].WritesMerged)
		flat.DeviceAddWrittenSectors(prof.Builder, u.Device[i].WrittenSectors)
		flat.DeviceAddWritingTime(prof.Builder, u.Device[i].WritingTime)
		flat.DeviceAddIOInProgress(prof.Builder, u.Device[i].IOInProgress)
		flat.DeviceAddIOTime(prof.Builder, u.Device[i].IOTime)
		flat.DeviceAddWeightedIOTime(prof.Builder, u.Device[i].WeightedIOTime)
		devF[i] = flat.DeviceEnd(prof.Builder)
	}
	flat.DiskUsageStartDeviceVector(prof.Builder, len(devF))
	for i := len(devF) - 1; i >= 0; i-- {
		prof.Builder.PrependUOffsetT(devF[i])
	}
	devV := prof.Builder.EndVector(len(devF))
	flat.DiskUsageStart(prof.Builder)
	flat.DiskUsageAddTimestamp(prof.Builder, u.Timestamp)
	flat.DiskUsageAddTimeDelta(prof.Builder, u.TimeDelta)
	flat.DiskUsageAddDevice(prof.Builder, devV)
	prof.Builder.Finish(flat.DiskUsageEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize IO usage of the block devices as Flatbuffer serialized bytes using
// the package's global Profiler.
func Serialize(u *structs.DiskUsage) (p []byte, err error) {
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

// Deserialize takes some Flatbuffer serialized bytes and deserializes them
// as a structs.DiskUsage.
func Deserialize(p []byte) *structs.DiskUsage {
	u := &structs.DiskUsage{}
	devF := &flat.Device{}
	uF := flat.GetRootAsDiskUsage(p, 0)
	u.Timestamp = uF.Timestamp()
	u.TimeDelta = uF.TimeDelta()
	len := uF.DeviceLength()
	u.Device = make([]structs.Device, len)
	for i := 0; i < len; i++ {
		var dev structs.Device
		if uF.Device(devF, i) {
			dev.Major = devF.Major()
			dev.Minor = devF.Minor()
			dev.Name = string(devF.Name())
			dev.ReadsCompleted = devF.ReadsCompleted()
			dev.ReadsMerged = devF.ReadsMerged()
			dev.ReadSectors = devF.ReadSectors()
			dev.ReadingTime = devF.ReadingTime()
			dev.WritesCompleted = devF.WritesCompleted()
			dev.WritesMerged = devF.WritesMerged()
			dev.WrittenSectors = devF.WrittenSectors()
			dev.WritingTime = devF.WritingTime()
			dev.IOInProgress = devF.IOInProgress()
			dev.IOTime = devF.IOTime()
			dev.WeightedIOTime = devF.WeightedIOTime()
		}
		u.Device[i] = dev
	}
	return u
}

// Ticker delivers the system's IO usage of the block devices at intervals.
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

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

// Package diskstats handles processing of IO statistics of each block device,
// /proc/diskstats. Instead of returning a Go struct, it returns Flatbuffer
// serialized bytes. A function to deserialize the Flatbuffer serialized bytes
// into a structs.DiskStats struct is provided.
//
// Note: the package name is diskstats and not the final element of the import
// path (flat). 
package diskstats

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/rai-project/joefriday"
	stats "github.com/rai-project/joefriday/disk/diskstats"
	"github.com/rai-project/joefriday/disk/structs"
	"github.com/rai-project/joefriday/disk/structs/flat"
)

// Profiler is used to process the /proc/diskstast file.
type Profiler struct {
	*stats.Profiler
	*fb.Builder
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	p, err := stats.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns information about current IO statistics of the block devices as
// Flatbuffer serialized bytes.
func (prof *Profiler) Get() ([]byte, error) {
	stts, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(stts), nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns information about current IO statistics of the block devices as
// Flatbuffer serialized bytes using the package's global Profiler.
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

// Serialize serializes structs.DiskStats as Flatbuffer serialized bytes.
func (prof *Profiler) Serialize(stts *structs.DiskStats) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	devF := make([]fb.UOffsetT, len(stts.Device))
	names := make([]fb.UOffsetT, len(stts.Device))
	for i := 0; i < len(names); i++ {
		names[i] = prof.Builder.CreateString(stts.Device[i].Name)
	}
	for i := 0; i < len(devF); i++ {
		flat.DeviceStart(prof.Builder)
		flat.DeviceAddMajor(prof.Builder, stts.Device[i].Major)
		flat.DeviceAddMinor(prof.Builder, stts.Device[i].Minor)
		flat.DeviceAddName(prof.Builder, names[i])
		flat.DeviceAddReadsCompleted(prof.Builder, stts.Device[i].ReadsCompleted)
		flat.DeviceAddReadsMerged(prof.Builder, stts.Device[i].ReadsMerged)
		flat.DeviceAddReadSectors(prof.Builder, stts.Device[i].ReadSectors)
		flat.DeviceAddReadingTime(prof.Builder, stts.Device[i].ReadingTime)
		flat.DeviceAddWritesCompleted(prof.Builder, stts.Device[i].WritesCompleted)
		flat.DeviceAddWritesMerged(prof.Builder, stts.Device[i].WritesMerged)
		flat.DeviceAddWrittenSectors(prof.Builder, stts.Device[i].WrittenSectors)
		flat.DeviceAddWritingTime(prof.Builder, stts.Device[i].WritingTime)
		flat.DeviceAddIOInProgress(prof.Builder, stts.Device[i].IOInProgress)
		flat.DeviceAddIOTime(prof.Builder, stts.Device[i].IOTime)
		flat.DeviceAddWeightedIOTime(prof.Builder, stts.Device[i].WeightedIOTime)
		devF[i] = flat.DeviceEnd(prof.Builder)
	}
	flat.DiskStatsStartDeviceVector(prof.Builder, len(devF))
	for i := len(devF) - 1; i >= 0; i-- {
		prof.Builder.PrependUOffsetT(devF[i])
	}
	devV := prof.Builder.EndVector(len(devF))
	flat.DiskStatsStart(prof.Builder)
	flat.DiskStatsAddTimestamp(prof.Builder, stts.Timestamp)
	flat.DiskStatsAddDevice(prof.Builder, devV)
	prof.Builder.Finish(flat.DiskStatsEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize serializes structs.DiskStats as Flatbuffer serialized bytes using
// the package's global Profiler.
func Serialize(stts *structs.DiskStats) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(stts), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them
// as a structs.DiskStats.
func Deserialize(p []byte) *structs.DiskStats {
	stts := &structs.DiskStats{}
	devF := &flat.Device{}
	statsFlat := flat.GetRootAsDiskStats(p, 0)
	stts.Timestamp = statsFlat.Timestamp()
	len := statsFlat.DeviceLength()
	stts.Device = make([]structs.Device, len)
	for i := 0; i < len; i++ {
		var dev structs.Device
		if statsFlat.Device(devF, i) {
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
		stts.Device[i] = dev
	}
	return stts
}

// Ticker delivers the system's IO statistics of the block devices at
// intervals.
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

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

// Package cpuutil handles Flatbuffer based processing of CPU (kernel)
// utilization information. This information is calculated using the
// difference between two CPU (kernel) stats snapshots, /proc/stat, and
// represented as a percentage. The time elapsed between the two snapshots is
// stored in the TimeDelta field. Instead of returning a Go struct, it returns
// the data as Flatbuffer serialized bytes. For convenience, a function to
// deserialize the Flatbuffer serialized bytes into a cpuutil.Utilization
// struct is provided. After the first use, the flatbuffer builder is reused.
// 
// Note: the package name is cpuutil and not the final element of the import
// path (flat). 
package cpuutil

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	util "github.com/mohae/joefriday/cpu/cpuutil"
	"github.com/mohae/joefriday/cpu/cpuutil/flat/structs"
)

// Profiler is used to process the /proc/stats file and calculate utilization
// information, returning the data as Flatbuffer serialized bytes.
type Profiler struct {
	*util.Profiler
	*fb.Builder
}

// Initializes and returns a cpu utilization profiler that uses FlatBuffers.
func NewProfiler() (prof *Profiler, err error) {
	p, err := util.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the cpu utilization as Flatbuffer serialized bytes. Utilization
// calculations requires two snapshots. This func gets the current snapshot of
// /proc/stat and calculates the utilization using the difference between the
// current snapshot and the prior one. The current snapshot is stored and for
// use as the prior snapshot on the next Get call. If ongoing utilitzation
// information is desired, the Ticker should be used; it's better suited for
// ongoing utilization information.
func (prof *Profiler) Get() (p []byte, err error) {
	u, err := prof.Profiler.Get()
	return prof.Serialize(u), nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current cpu utilization as Flatbuffer serialized bytes using
// the package's global Profiler. The Profiler is instantiated lazily. If the
// profiler doesn't already exist, the first usage information will not be
// useful due to minimal time elapsing between the initial and second snapshots
// used for usage calculations; the results of the first call should be
// discarded.
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

// Serialize cpu utilization using Flatbuffers.
func (prof *Profiler) Serialize(u *util.Utilization) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	utils := make([]fb.UOffsetT, len(u.CPU))
	ids := make([]fb.UOffsetT, len(u.CPU))
	for i := 0; i < len(ids); i++ {
		ids[i] = prof.Builder.CreateString(u.CPU[i].ID)
	}
	for i := 0; i < len(utils); i++ {
		structs.UtilStart(prof.Builder)
		structs.UtilAddID(prof.Builder, ids[i])
		structs.UtilAddUsage(prof.Builder, u.CPU[i].Usage)
		structs.UtilAddUser(prof.Builder, u.CPU[i].User)
		structs.UtilAddNice(prof.Builder, u.CPU[i].Nice)
		structs.UtilAddSystem(prof.Builder, u.CPU[i].System)
		structs.UtilAddIdle(prof.Builder, u.CPU[i].Idle)
		structs.UtilAddIOWait(prof.Builder, u.CPU[i].IOWait)
		utils[i] = structs.UtilEnd(prof.Builder)
	}
	structs.UtilizationStartCPUVector(prof.Builder, len(utils))
	for i := len(utils) - 1; i >= 0; i-- {
		prof.Builder.PrependUOffsetT(utils[i])
	}
	utilsV := prof.Builder.EndVector(len(utils))
	structs.UtilizationStart(prof.Builder)
	structs.UtilizationAddTimestamp(prof.Builder, u.Timestamp)
	structs.UtilizationAddTimeDelta(prof.Builder, u.TimeDelta)
	structs.UtilizationAddBTimeDelta(prof.Builder, u.BTimeDelta)
	structs.UtilizationAddCtxtDelta(prof.Builder, u.CtxtDelta)
	structs.UtilizationAddProcesses(prof.Builder, u.Processes)
	structs.UtilizationAddCPU(prof.Builder, utilsV)
	prof.Builder.Finish(structs.UtilizationEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize the CPU Utilization using the package global Profiler.
func Serialize(u *util.Utilization) (p []byte, err error) {
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

// Deserialize takes some Flatbuffer serialized bytes and deserializes them as
// cpuutil.Utilization.
func Deserialize(p []byte) *util.Utilization {
	u := &util.Utilization{}
	uF := &structs.Util{}
	flatUtil := structs.GetRootAsUtilization(p, 0)
	u.Timestamp = flatUtil.Timestamp()
	u.TimeDelta = flatUtil.TimeDelta()
	u.CtxtDelta = flatUtil.CtxtDelta()
	u.BTimeDelta = flatUtil.BTimeDelta()
	u.Processes = flatUtil.Processes()
	len := flatUtil.CPULength()
	u.CPU = make([]util.Util, len)
	for i := 0; i < len; i++ {
		var util util.Util
		if flatUtil.CPU(uF, i) {
			util.ID = string(uF.ID())
			util.Usage = uF.Usage()
			util.User = uF.User()
			util.Nice = uF.Nice()
			util.System = uF.System()
			util.Idle = uF.Idle()
			util.IOWait = uF.IOWait()
		}
		u.CPU[i] = util
	}
	return u
}

// Ticker delivers the system's CPU utilization information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan []byte
	*Profiler
}

// NewTicker returns a new Ticker containing a Data channel that delivers
// the data at intervals and an error channel that delivers any errors
// encountered. Stop the ticker to signal the ticker to stop running. Stopping
// the ticker does not close the Data channel; call Close to close both the
// ticker and the data channel.
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

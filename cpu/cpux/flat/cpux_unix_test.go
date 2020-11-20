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

package cpux

import (
	"encoding/json"
	"testing"

	fb "github.com/google/flatbuffers/go"
	"github.com/c3sr/joefriday/cpu/cpux"
	"github.com/c3sr/joefriday/testinfo"
)

func TestCPUX(t *testing.T) {
	// set up test stuff w/o freq
	tSysFS := testinfo.NewTempSysFS()
	tSysFS.PhysicalPackageCount = 1
	tSysFS.CoresPerPhysicalPackage = 2
	tSysFS.ThreadsPerCore = 2
	err := tSysFS.SetSysFS("")
	if err != nil {
		t.Fatalf("setting up sysfs tree: %s", err)
	}
	defer tSysFS.Clean()

	err = tSysFS.CreateCPU()
	if err != nil {
		t.Errorf("setting up cpux testing info: %s", err)
		return
	}
	prof := &Profiler{Builder: fb.NewBuilder(0), Profiler: &cpux.Profiler{NumCPU: int(tSysFS.CPUs())}}
	prof.SysFSSystemPath(tSysFS.Path())
	p, err := prof.Get()
	if err != nil {
		t.Error(err)
	}
	cpus := Deserialize(p)

	//compare results cpufreq
	err = tSysFS.ValidateCPUX(cpus)
	if err != nil {
		t.Error(err)
	}
	js, _ := json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))

	// cleanup for next
	err = tSysFS.CleanCPU()
	if err != nil {
		t.Error(err)
	}

	// set up test stuff w/o freq
	tSysFS.Freq = false
	err = tSysFS.CreateCPU()
	if err != nil {
		t.Error("setting up cpux testing info: %s", err)
		goto multiSocket
	}

	p, err = prof.Get()
	if err != nil {
		t.Error(err)
	}
	cpus = Deserialize(p)

	js, _ = json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))
	// compare results with frequency
	err = tSysFS.ValidateCPUX(cpus)
	if err != nil {
		t.Error(err)
	}

	// cleanup for next
	err = tSysFS.CleanCPU()
	if err != nil {
		t.Error(err)
	}

multiSocket:
	// 2 sockets
	tSysFS.PhysicalPackageCount = 2
	prof.NumCPU = int(tSysFS.CPUs())
	tSysFS.Freq = true
	err = tSysFS.CreateCPU()
	if err != nil {
		t.Error("setting up cpux testing info: %s", err)
		goto noFreq
	}

	p, err = prof.Get()
	if err != nil {
		t.Error(err)
	}
	cpus = Deserialize(p)

	//compare results cpufreq
	err = tSysFS.ValidateCPUX(cpus)
	if err != nil {
		t.Error(err)
	}
	js, _ = json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))

	// cleanup for next
	err = tSysFS.CleanCPU()
	if err != nil {
		t.Error(err)
	}

noFreq:
	// set up test stuff w/o freq
	tSysFS.Freq = false
	err = tSysFS.CreateCPU()
	if err != nil {
		t.Error("setting up cpux testing info: %s", err)
		return
	}

	p, err = prof.Get()
	if err != nil {
		t.Error(err)
	}
	cpus = Deserialize(p)

	js, _ = json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))
	// compare results with frequency
	err = tSysFS.ValidateCPUX(cpus)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkGet(b *testing.B) {
	var jsn []byte
	p := NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Get()
	}
	_ = jsn
}

func BenchmarkSerialize(b *testing.B) {
	var f []byte
	p := NewProfiler()
	v, _ := p.Profiler.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f = p.Serialize(v)
	}
	_ = f
}

func BenchmarkDeserialize(b *testing.B) {
	var cpus *cpux.CPUs
	p := NewProfiler()
	cpusb, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cpus = Deserialize(cpusb)
	}
	_ = cpus
}

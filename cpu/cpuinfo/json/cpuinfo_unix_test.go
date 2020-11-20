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

package cpuinfo

import (
	"testing"

	"github.com/c3sr/joefriday"
	"github.com/c3sr/joefriday/cpu/cpuinfo"
	"github.com/c3sr/joefriday/testinfo"
)

func TestGeti75600u(t *testing.T) {
	tProc, err := joefriday.NewTempFileProc("intel", "i9700u", testinfo.I75600uCPUInfo)
	if err != nil {
		t.Fatal(err)
	}
	defer tProc.Remove()
	prof, err := NewProfiler()
	if err != nil {
		t.Fatal(err)
	}
	prof.Profiler.Procer = tProc
	inf, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	info, err := Unmarshal(inf)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	err = testinfo.ValidateI75600uCPUInfo(info)
	if err != nil {
		t.Error(err)
	}
}

func TestGetR71800xJSON(t *testing.T) {
	tProc, err := joefriday.NewTempFileProc("amd", "r71800x", testinfo.R71800xCPUInfo)
	if err != nil {
		t.Fatal(err)
	}
	defer tProc.Remove()
	prof, err := NewProfiler()
	if err != nil {
		t.Fatal(err)
	}
	prof.Procer = tProc
	inf, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	info, err := Unmarshal(inf)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	err = testinfo.ValidateR71800xCPUInfo(info)
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}

func BenchmarkGet(b *testing.B) {
	var jsn []byte
	p, _ := NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Get()
	}
	_ = jsn
}

func BenchmarkSerialize(b *testing.B) {
	var jsn []byte
	p, _ := NewProfiler()
	v, _ := p.Profiler.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Serialize(v)
	}
	_ = jsn
}

func BenchmarkMarshal(b *testing.B) {
	var jsn []byte
	p, _ := NewProfiler()
	v, _ := p.Profiler.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Marshal(v)
	}
	_ = jsn
}

func BenchmarkDeserialize(b *testing.B) {
	var inf *cpuinfo.CPUInfo
	p, _ := NewProfiler()
	infB, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = Deserialize(infB)
	}
	_ = inf
}

func BenchmarkUnmarshal(b *testing.B) {
	var inf *cpuinfo.CPUInfo
	p, _ := NewProfiler()
	infB, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = Unmarshal(infB)
	}
	_ = inf
}

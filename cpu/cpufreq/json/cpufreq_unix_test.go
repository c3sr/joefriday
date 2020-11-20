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

package cpufreq

import (
	"testing"
	"time"

	"github.com/c3sr/joefriday"
	freq "github.com/c3sr/joefriday/cpu/cpufreq"
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
	err = prof.InitFrequency()
	f, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	freq, err := Unmarshal(f)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	err = testinfo.ValidateI75600uCPUFreq(freq)
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
	err = prof.InitFrequency()
	f, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	ff, err := Unmarshal(f)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	err = testinfo.ValidateR71800xCPUFreq(ff)
	if err != nil {
		t.Error(err)
	}
	t.Log(ff)
}

func TestGetXeonE52690(t *testing.T) {
	tProc, err := joefriday.NewTempFileProc("intel", "xeon_e52690", testinfo.XeonE52690CPUInfo)
	if err != nil {
		t.Fatal(err)
	}
	defer tProc.Remove()
	prof, err := NewProfiler()
	if err != nil {
		t.Fatal(err)
	}
	prof.Profiler.Procer = tProc
	err = prof.InitFrequency()
	f, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	freq, err := Unmarshal(f)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	err = testinfo.ValidateXeonE52690CPUFreq(freq)
	if err != nil {
		t.Error(err)
	}
}

func TestTicker(t *testing.T) {
	tkr, err := NewTicker(time.Millisecond)
	if err != nil {
		t.Error(err)
		return
	}
	tProc, err := joefriday.NewTempFileProc("intel", "i9700u", testinfo.I75600uCPUInfo)
	if err != nil {
		t.Fatal(err)
	}
	defer tProc.Remove()
	prof, err := NewProfiler()
	if err != nil {
		t.Fatal(err)
	}
	prof.Procer = tProc
	err = prof.InitFrequency()
	tk := tkr.(*Ticker)
	tk.Profiler = prof

	for i := 0; i < 5; i++ {
		select {
		case <-tk.Done:
			break
		case v, ok := <-tk.Data:
			if !ok {
				break
			}
			f, err := Deserialize(v)
			if err != nil {
				t.Error(err)
				continue
			}
			err = testinfo.ValidateI75600uCPUFreq(f)
			if err != nil {
				t.Error(err)
			}
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
	tk.Close()
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
	var f *freq.Frequency
	p, _ := NewProfiler()
	fB, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, _ = Deserialize(fB)
	}
	_ = f
}

func BenchmarkUnmarshal(b *testing.B) {
	var f *freq.Frequency
	p, _ := NewProfiler()
	fB, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, _ = Unmarshal(fB)
	}
	_ = f
}

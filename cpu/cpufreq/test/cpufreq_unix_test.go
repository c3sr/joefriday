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

package cpufreqtest

import (
	"testing"
	"time"

	"github.com/c3sr/joefriday"
	"github.com/c3sr/joefriday/cpu/cpufreq"
	"github.com/c3sr/joefriday/testinfo"
)

func TestGeti75600u(t *testing.T) {
	tProc, err := joefriday.NewTempFileProc("intel", "i75600u", testinfo.I75600uCPUInfo)
	if err != nil {
		t.Fatal(err)
	}
	defer tProc.Remove()
	prof, err := cpufreq.NewProfiler()
	if err != nil {
		t.Fatal(err)
	}
	prof.Procer = tProc
	err = prof.InitFrequency()
	if err != nil {
		t.Fatal(err)
	}

	f, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	err = testinfo.ValidateI75600uCPUFreq(f)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", f)
}

func TestGetR71800x(t *testing.T) {
	tProc, err := joefriday.NewTempFileProc("amd", "r71800x", testinfo.R71800xCPUInfo)
	if err != nil {
		t.Fatal(err)
	}
	defer tProc.Remove()
	prof, err := cpufreq.NewProfiler()
	if err != nil {
		t.Fatal(err)
	}
	prof.Procer = tProc
	err = prof.InitFrequency()
	if err != nil {
		t.Fatal(err)
	}

	f, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	err = testinfo.ValidateR71800xCPUFreq(f)
	if err != nil {
		t.Error(err)
	}
	t.Log(f)
}

func TestGetXeonE52690(t *testing.T) {
	tProc, err := joefriday.NewTempFileProc("intel", "xeonE52690", testinfo.XeonE52690CPUInfo)
	if err != nil {
		t.Fatal(err)
	}
	defer tProc.Remove()
	prof, err := cpufreq.NewProfiler()
	if err != nil {
		t.Fatal(err)
	}
	prof.Procer = tProc
	err = prof.InitFrequency()
	if err != nil {
		t.Fatal(err)
	}

	f, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	err = testinfo.ValidateXeonE52690CPUFreq(f)
	if err != nil {
		t.Error(err)
	}
	t.Log(f)
}

func TestTicker(t *testing.T) {
	tkr, err := cpufreq.NewTicker(time.Millisecond)
	if err != nil {
		t.Error(err)
		return
	}
	tProc, err := joefriday.NewTempFileProc("intel", "i9700u", testinfo.I75600uCPUInfo)
	if err != nil {
		t.Fatal(err)
	}
	defer tProc.Remove()
	prof, err := cpufreq.NewProfiler()
	if err != nil {
		t.Fatal(err)
	}
	prof.Procer = tProc
	err = prof.InitFrequency()
	if err != nil {
		t.Fatal(err)
	}

	tk := tkr.(*cpufreq.Ticker)
	tk.Profiler = prof
	for i := 0; i < 5; i++ {
		select {
		case <-tk.Done:
			break
		case v, ok := <-tk.Data:
			if !ok {
				break
			}
			err = testinfo.ValidateI75600uCPUFreq(v)
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

func BenchmarkGeti75600u(b *testing.B) {
	var f *cpufreq.Frequency
	p, _ := cpufreq.NewProfiler()
	tProc, err := joefriday.NewTempFileProc("intel", "i9700u", testinfo.I75600uCPUInfo)
	if err != nil {
		b.Fatal(err)
	}
	defer tProc.Remove()
	prof, err := cpufreq.NewProfiler()
	if err != nil {
		b.Fatal(err)
	}
	prof.Procer = tProc
	err = prof.InitFrequency()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, _ = p.Get()
	}
	_ = f
}

func BenchmarkGetXeonE52690(b *testing.B) {
	var f *cpufreq.Frequency
	p, _ := cpufreq.NewProfiler()
	tProc, err := joefriday.NewTempFileProc("intel", "xeonE5290", testinfo.XeonE52690CPUInfo)
	if err != nil {
		b.Fatal(err)
	}
	defer tProc.Remove()
	prof, err := cpufreq.NewProfiler()
	if err != nil {
		b.Fatal(err)
	}
	prof.Procer = tProc
	err = prof.InitFrequency()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, _ = p.Get()
	}
	_ = f
}
func BenchmarkGetR71800x(b *testing.B) {
	var f *cpufreq.Frequency
	p, _ := cpufreq.NewProfiler()
	tProc, err := joefriday.NewTempFileProc("amd", "r71800x", testinfo.R71800xCPUInfo)
	if err != nil {
		b.Fatal(err)
	}
	defer tProc.Remove()
	prof, err := cpufreq.NewProfiler()
	if err != nil {
		b.Fatal(err)
	}
	prof.Procer = tProc
	err = prof.InitFrequency()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, _ = p.Get()
	}
	_ = f
}

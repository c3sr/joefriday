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

package cpuutil

import (
	"testing"
	"time"

	util "github.com/rai-project/joefriday/cpu/cpuutil"
)

func TestGet(t *testing.T) {
	p, err := NewProfiler()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	time.Sleep(time.Duration(300) * time.Millisecond)
	u, err := p.Get()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	ut, err := Unmarshal(u)
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	checkCPUUtil("get", ut, t)
	t.Logf("%#v\n", ut)
}

func TestTicker(t *testing.T) {
	tkr, err := NewTicker(time.Duration(200) * time.Millisecond)
	if err != nil {
		t.Error(err)
		return
	}
	tk := tkr.(*Ticker)
	for i := 0; i < 5; i++ {
		select {
		case <-tk.Done:
			break
		case v, ok := <-tk.Data:
			if !ok {
				break
			}
			st, err := Deserialize(v)
			if err != nil {
				t.Error(err)
				continue
			}
			checkCPUUtil("ticker", st, t)
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
	tk.Close()
}

func checkCPUUtil(name string, u *util.CPUUtil, t *testing.T) {
	if u.Timestamp == 0 {
		t.Errorf("%s: timestamp: expected on-zero", name)
	}
	if u.TimeDelta == 0 {
		t.Errorf("%s: TimeDelta: expected non-zero value, got 0", name)
	}
	if u.CtxtDelta == 0 {
		t.Errorf("%s: CtxtDelta: expected non-zero value, got 0", name)
	}
	if u.BTimeDelta == 0 {
		t.Errorf("%s: BTimeDelta: expected non-zero value, got 0", name)
	}
	if u.Processes == 0 {
		t.Errorf("%s: Processes: expected non-zero value, got 0", name)
	}
	if len(u.CPU) < 2 {
		t.Errorf("%s: cpu: got %d, want at least 2", name, len(u.CPU))
	}
	for i, v := range u.CPU {
		if v.ID == "" {
			t.Errorf("%d: %s: expected ID to have a value, was empty", i, name)
		}
	}
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
	var u *util.CPUUtil
	p, _ := NewProfiler()
	uB, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u, _ = Deserialize(uB)
	}
	_ = u
}

func BenchmarkUnmarshal(b *testing.B) {
	var u *util.CPUUtil
	p, _ := NewProfiler()
	uB, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u, _ = Unmarshal(uB)
	}
	_ = u
}

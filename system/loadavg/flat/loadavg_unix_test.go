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

package loadavg

import (
	"testing"
	"time"

	l "github.com/rai-project/joefriday/system/loadavg"
)

func TestGet(t *testing.T) {
	p, err := Get()
	if err != nil {
		t.Errorf("Get(): got %s, want nil", err)
		return
	}
	l := Deserialize(p)
	checkLoad("get", l, t)
}

func TestTicker(t *testing.T) {
	tkr, err := NewTicker(time.Millisecond)
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
			inf := Deserialize(v)
			if err != nil {
				t.Error(err)
				continue
			}
			checkLoad("ticker", inf, t)
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
	tk.Close()
}

func checkLoad(n string, la l.LoadAvg, t *testing.T) {
	if la.Timestamp == 0 {
		t.Errorf("%s: expected Timestamp to be a non-zero value; got 0", n)
	}
	if la.Minute == 0 {
		t.Errorf("%s: expected Minute to be a non-zero value; got 0", n)
	}
	if la.Five == 0 {
		t.Errorf("%s: expected Five to be a non-zero value; got 0", n)
	}
	if la.Fifteen == 0 {
		t.Errorf("%s: expected Fifteen to be a non-zero value; got 0", n)
	}
	if la.Running == 0 {
		t.Errorf("%s: expected Running to be a non-zero value; got 0", n)
	}
	if la.Total == 0 {
		t.Errorf("%s: expected Total to be a non-zero value; got 0", n)
	}
	if la.PID == 0 {
		t.Errorf("%s: expected PID to be a non-zero value; got 0", n)
	}
}

func BenchmarkGet(b *testing.B) {
	var tmp []byte
	p, _ := NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func BenchmarkSerialize(b *testing.B) {
	var tmp []byte
	p, _ := NewProfiler()
	l, _ := p.Profiler.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = Serialize(l)
	}
	_ = tmp
}

func BenchmarkDeserialize(b *testing.B) {
	var la l.LoadAvg
	p, _ := NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		la = Deserialize(tmp)
	}
	_ = la
}

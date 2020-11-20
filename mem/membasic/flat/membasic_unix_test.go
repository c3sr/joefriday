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

package membasic

import (
	"testing"
	"time"

	basic "github.com/c3sr/joefriday/mem/membasic"
)

func TestSerializeDeserialize(t *testing.T) {
	p, err := Get()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	inf := Deserialize(p)
	checkInfo("get", *inf, t)
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
			checkInfo("ticker", *inf, t)
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
	tk.Close()
}

func checkInfo(n string, i basic.Info, t *testing.T) {
	if i.Timestamp == 0 {
		t.Errorf("%s: expected timestamp to be a non-zero value, got 0", n)
	}
	if i.Active == 0 {
		t.Errorf("%s: expected Active to be a non-zero value, got 0", n)
	}
	if i.Inactive == 0 {
		t.Errorf("%s: expected Inactive to be a non-zero value, got 0", n)
	}
	if i.Mapped == 0 {
		t.Errorf("%s: expected Mapped to be a non-zero value, got 0", n)
	}
	if i.MemAvailable == 0 {
		t.Errorf("%s: expected MemAvailable to be a non-zero value, got 0", n)
	}
	if i.MemFree == 0 {
		t.Errorf("%s: expected MemFree to be a non-zero value, got 0", n)
	}
	if i.MemTotal == 0 {
		t.Errorf("%s: expected MemTotal to be a non-zero value, got 0", n)
	}
	if i.SwapFree == 0 {
		t.Errorf("%s: expected SwapFree to be a non-zero value, got 0", n)
	}
	if i.SwapTotal == 0 {
		t.Errorf("%s: expected SwapTotal to be a non-zero value, got 0", n)
	}
	t.Logf("%#v\n", i)
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
	inf, _ := p.Profiler.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = Serialize(inf)
	}
	_ = tmp
}

var inf *basic.Info

func BenchmarkDeserialize(b *testing.B) {
	p, _ := NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inf = Deserialize(tmp)
	}
	_ = inf
}

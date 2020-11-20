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

package meminfo

import (
	"testing"
	"time"

	mem "github.com/c3sr/joefriday/mem/meminfo"
)

func TestGet(t *testing.T) {
	nf, err := Get()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	inf, err := Unmarshal(nf)
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
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
			inf, err := Unmarshal(v)
			if err != nil {
				t.Errorf("got %s, want nil", err)
				return
			}
			checkInfo("ticker", *inf, t)
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
	tk.Close()
}

func checkInfo(n string, i mem.Info, t *testing.T) {
	if i.Timestamp == 0 {
		t.Errorf("%s: expected timestamp to be a non-zero value, got 0", n)
	}
	if i.Active == 0 {
		t.Errorf("%s: expected Active to be a non-zero value, got 0", n)
	}
	if i.ActiveAnon == 0 {
		t.Errorf("%s: expected ActiveAnon to be a non-zero value, got 0", n)
	}
	if i.ActiveFile == 0 {
		t.Errorf("%s: expected ActiveFile to be a non-zero value, got 0", n)
	}
	if i.AnonPages == 0 {
		t.Errorf("%s: expected AnonPages to be a non-zero value, got 0", n)
	}
	if i.Buffers == 0 {
		t.Errorf("%s: expected Buffers to be a non-zero value, got 0", n)
	}
	if i.Cached == 0 {
		t.Errorf("%s: expected Cached to be a non-zero value, got 0", n)
	}
	if i.CommitLimit == 0 {
		t.Errorf("%s: expected CommitLimit to be a non-zero value, got 0", n)
	}
	if i.CommittedAS == 0 {
		t.Errorf("%s: expected CommittedAS to be a non-zero value, got 0", n)
	}
	if i.DirectMap4K == 0 {
		t.Errorf("%s: expected DirectMap4K to be a non-zero value, got 0", n)
	}
	if i.DirectMap2M == 0 {
		t.Errorf("%s: expected DirectMap2M to be a non-zero value, got 0", n)
	}
	if i.HugePagesSize == 0 {
		t.Errorf("%s: expected HugePagesSize to be a non-zero value, got 0", n)
	}
	if i.Inactive == 0 {
		t.Errorf("%s: expected Inactive to be a non-zero value, got 0", n)
	}
	if i.InactiveAnon == 0 {
		t.Errorf("%s: expected InactiveAnon to be a non-zero value, got 0", n)
	}
	if i.InactiveFile == 0 {
		t.Errorf("%s: expected InactiveFile to be a non-zero value, got 0", n)
	}
	if i.KernelStack == 0 {
		t.Errorf("%s: expected KernelStack to be a non-zero value, got 0", n)
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
	if i.PageTables == 0 {
		t.Errorf("%s: expected PageTables to be a non-zero value, got 0", n)
	}
	if i.Shmem == 0 {
		t.Errorf("%s: expected Shmem to be a non-zero value, got 0", n)
	}
	if i.Slab == 0 {
		t.Errorf("%s: expected Slab to be a non-zero value, got 0", n)
	}
	if i.SReclaimable == 0 {
		t.Errorf("%s: expected SReclaimable to be a non-zero value, got 0", n)
	}
	if i.SUnreclaim == 0 {
		t.Errorf("%s: expected SReclaimable to be a non-zero value, got 0", n)
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

var inf *mem.Info

func BenchmarkDeserialize(b *testing.B) {
	p, _ := NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = Deserialize(tmp)
	}
	_ = inf
}

func BenchmarkUnmarshal(b *testing.B) {
	p, _ := NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = Unmarshal(tmp)
	}
	_ = inf
}

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

package diskusage

import (
	"testing"
	"time"

	"github.com/c3sr/joefriday/disk/structs"
)

func TestSerializeDeserialize(t *testing.T) {
	p, err := NewProfiler()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	b, err := p.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	u := Deserialize(b)
	checkUsage("get", u, t)
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
			u := Deserialize(v)
			checkUsage("ticker", u, t)
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
	tk.Close()
}

func checkUsage(n string, u *structs.DiskUsage, t *testing.T) {
	if u.Timestamp == 0 {
		t.Errorf("%s: Timestamp: wanted non-zero value; got 0", n)
	}
	if u.TimeDelta == 0 {
		t.Errorf("%s: TimeDelta: wanted non-zero value; got 0", n)
	}
	if len(u.Device) == 0 {
		t.Errorf("%s: expected there to be devices; didn't get any", n)
	}
	for i := 0; i < len(u.Device); i++ {
		if u.Device[i].Major == 0 {
			t.Errorf("%s: Device %d: Major: wanted a non-zero value, was 0", n, i)
		}
		if u.Device[i].Name == "" {
			t.Errorf("%s: Device %d: Name: wanted a non-empty value; was empty", n, i)
		}
	}
}

func BenchmarkSerialize(b *testing.B) {
	var tmp []byte
	p, _ := NewProfiler()
	u, _ := p.Profiler.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = Serialize(u)
	}
	_ = tmp
}

func BenchmarkDeserialize(b *testing.B) {
	var u *structs.DiskUsage
	p, _ := NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u = Deserialize(tmp)
	}
	_ = u
}

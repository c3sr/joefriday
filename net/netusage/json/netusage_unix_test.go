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

package netusage

import (
	"testing"
	"time"

	"github.com/rai-project/joefriday/net/structs"
)

func TestGet(t *testing.T) {
	p, err := NewProfiler()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	b, err := p.Get()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	u, err := Deserialize(b)
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	checkUsage("get", u, t)
	t.Logf("%#v\n", u)
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
			u, err := Deserialize(v)
			if err != nil {
				t.Error(err)
				continue
			}
			checkUsage("ticker", u, t)
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
	tk.Close()
}

func checkUsage(n string, u *structs.DevUsage, t *testing.T) {
	if u.Timestamp == 0 {
		t.Errorf("%s: expected timestamp to be a non-zero value; was 0", n)
	}
	if u.TimeDelta == 0 {
		t.Errorf("%s: expected TimeDelta to be a non-zero value; was 0", n)
	}
	if len(u.Device) == 0 {
		t.Errorf("%s: expected devices; got none", n)
		return
	}
	// check name
	for i, v := range u.Device {
		if v.Name == "" {
			t.Errorf("%s: %d: expected device to have a name; was empty", n, i)
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
	var u *structs.DevUsage
	p, _ := NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u, _ = Deserialize(tmp)
	}
	_ = u
}

func BenchmarkUnmarshal(b *testing.B) {
	var u *structs.DevUsage
	p, _ := NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u, _ = Unmarshal(tmp)
	}
	_ = u
}

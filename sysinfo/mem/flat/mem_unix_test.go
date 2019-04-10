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

package mem

import (
	"testing"
	"time"

	m "github.com/rai-project/joefriday/sysinfo/mem"
)

func TestSerializeDeserialize(t *testing.T) {
	p, err := Get()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	m := Deserialize(p)
	checkMemInfo("get", m, t)
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
		case p, ok := <-tk.Data:
			if !ok {
				break
			}
			m := Deserialize(p)
			checkMemInfo("ticker", m, t)
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
	tk.Close()
}

func checkMemInfo(n string, inf *m.MemInfo, t *testing.T) {
	if inf.Timestamp == 0 {
		t.Errorf("%s: expected the Timestamp to be non-zero, was 0", n)
	}
	if inf.TotalRAM == 0 {
		t.Errorf("%s: expected the TotalRAM to be non-zero, was 0", n)
	}
	if inf.FreeRAM == 0 {
		t.Errorf("%s: expected the FreeRAM to be non-zero, was 0", n)
	}
	t.Logf("%#v\n", inf)
}

func BenchmarkGet(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		tmp, _ = Get()
	}
	_ = tmp
}

func BenchmarkSerialize(b *testing.B) {
	var tmp []byte
	var inf m.MemInfo
	inf.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp = Serialize(&inf)
	}
	_ = tmp
}

var inf *m.MemInfo

func BenchmarkDeserialize(b *testing.B) {
	p, _ := Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inf = Deserialize(p)
	}
	_ = inf
}

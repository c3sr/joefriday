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

	load "github.com/c3sr/joefriday/sysinfo/loadavg"
)

func TestSerializeDeserialize(t *testing.T) {
	p, err := Get()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	l, err := Deserialize(p)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	checkMemInfo("get", l, t)
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
			l, err := Deserialize(p)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				continue
			}
			checkMemInfo("ticker", l, t)
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
	tk.Close()
}

func checkMemInfo(n string, l *load.LoadAvg, t *testing.T) {
	if l.Timestamp == 0 {
		t.Errorf("%s: expected the Timestamp to be non-zero, was 0", n)
	}
	if l.One == 0 {
		t.Errorf("%s: expected the One to be non-zero, was 0", n)
	}
	if l.Five == 0 {
		t.Errorf("%s: expected the Five to be non-zero, was 0", n)
	}
	if l.Fifteen == 0 {
		t.Errorf("%s: expected the Fifteen to be non-zero, was 0", n)
	}
	t.Logf("%#v\n", l)
}

func BenchmarkGet(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		tmp, _ = Get()
	}
	_ = tmp
}

func BenchmarkDeserialize(b *testing.B) {
	var l *load.LoadAvg
	p, _ := Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l, _ = Deserialize(p)
	}
	_ = l
}

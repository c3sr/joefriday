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

package flat

import (
	"testing"

	"github.com/mohae/joefriday/cpu/facts"
)

func TestSerialize(t *testing.T) {
	facts, err := Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	factsD := Deserialize(facts)
	if factsD.Timestamp == 0 {
		t.Error("timestamp: expected non-zero timestamp")
	}
	if len(factsD.CPU) == 0 {
		t.Error("expected at least 1 cpu entries; got none")
	}
	for i := 0; i < len(factsD.CPU); i++ {
		if factsD.CPU[i].VendorID == "" {
			t.Errorf("%d: VendorID: expected Vendor ID to not be empty, it was", i)
		}
		if factsD.CPU[i].Model == "" {
			t.Errorf("%d: Model: expected model to not be empty; it was", i)
		}
		if factsD.CPU[i].CPUCores == 0 {
			t.Errorf("%d: CPUCores: expected non-zero value; was 0", i)
		}
		if len(factsD.CPU[i].Flags) == 0 {
			t.Errorf("%d: Flags: expected some flags, none were found", i)
		}
	}
	_, err = Serialize(factsD)
	if err != nil {
		t.Errorf("unexpected serialization error: %s", err)
		return
	}
}

func BenchmarkGet(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func BenchmarkSerialize(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := NewProfiler()
	fct, _ := p.Profiler.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = Serialize(fct)
	}
	_ = tmp
}

func BenchmarkDeserialize(b *testing.B) {
	var fct *facts.Facts
	b.StopTimer()
	p, _ := NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		fct = Deserialize(tmp)
	}
	_ = fct
}
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

package json

import (
	"testing"

	"github.com/mohae/joefriday/cpu/facts"
)

func TestGet(t *testing.T) {
	fct, err := Get()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	fcts, err := Unmarshal(fct)
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	if fcts.Timestamp == 0 {
		t.Error("expected timestamp to be a non-zero value; got 0")
	}
	if len(fcts.CPU) == 0 {
		t.Error("expected CPUs to be a non-zero value; got 0")
	}
	for i, v := range fcts.CPU {
		if v.VendorID == "" {
			t.Error("expected vendor_id to have a value; it was empty")
		}
		if v.Flags == "" {
			t.Error("expected flags to have values; it was empty")
		}
	}
	t.Logf("%#v\n", fcts)
}

var fct *facts.Facts

func BenchmarkGet(b *testing.B) {
	var jsn []byte
	p, _ := New()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Get()
	}
	_ = jsn
}

func BenchmarkUnmarshal(b *testing.B) {
	p, _ := New()
	fctB, _ := p.Get()
	for i := 0; i < b.N; i++ {
		fct, _ = Unmarshal(fctB)
	}
	_ = fct
}

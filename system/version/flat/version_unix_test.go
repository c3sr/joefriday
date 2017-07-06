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

package version

import (
	"testing"

	v "github.com/mohae/joefriday/system/version"
)

func TestSerializeDeserialize(t *testing.T) {
	p, err := Get()
	if err != nil {
		t.Errorf("Get(): got %s, want nil", err)
		return
	}
	inf, err := v.Get()
	if err != nil {
		t.Errorf("version.Get(): got %s, want nil", err)
		return
	}
	infD := Deserialize(p)
	if inf.OS != infD.OS {
		t.Errorf("OS: got %s; want %s", infD.OS, inf.OS)
	}
	if inf.Version != infD.Version {
		t.Errorf("Version: got %s; want %s", infD.Version, inf.Version)
	}
	if inf.CompileUser != infD.CompileUser {
		t.Errorf("CompileUser: got %s; want %s", infD.CompileUser, inf.CompileUser)
	}
	if inf.GCC != infD.GCC {
		t.Errorf("GCC: got %s; want %s", infD.GCC, inf.GCC)
	}
	if inf.OSGCC != infD.OSGCC {
		t.Errorf("Version: got %s; want %s", infD.OSGCC, inf.OSGCC)
	}
	if inf.Type != infD.Type {
		t.Errorf("Version: got %s; want %s", infD.Type, inf.Type)
	}
	if inf.CompileDate != infD.CompileDate {
		t.Errorf("CompileDate: got %s; want %s", infD.CompileDate, inf.CompileDate)
	}
	if inf.Arch != infD.Arch {
		t.Errorf("Arch: got %s; want %s", infD.Arch, inf.Arch)
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
	inf, _ := p.Profiler.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = Serialize(inf)
	}
	_ = tmp
}

func BenchmarkDeserialize(b *testing.B) {
	var inf *v.Info
	b.StopTimer()
	p, _ := NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		inf = Deserialize(tmp)
	}
	_ = inf
}
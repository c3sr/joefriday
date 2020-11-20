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

package cpuinfo

import (
	"testing"

	"github.com/c3sr/joefriday"
	"github.com/c3sr/joefriday/cpu/cpuinfo"
	"github.com/c3sr/joefriday/testinfo"
)

func TestGeti75600u(t *testing.T) {
	tProc, err := joefriday.NewTempFileProc("intel", "i75600u", testinfo.I75600uCPUInfo)
	if err != nil {
		t.Fatal(err)
	}
	defer tProc.Remove()
	prof, err := cpuinfo.NewProfiler()
	if err != nil {
		t.Fatal(err)
	}
	prof.Procer = tProc
	inf, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	err = testinfo.ValidateI75600uCPUInfo(inf)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", inf)
}

func TestGetR71800x(t *testing.T) {
	tProc, err := joefriday.NewTempFileProc("amd", "r71800x", testinfo.R71800xCPUInfo)
	if err != nil {
		t.Fatal(err)
	}
	defer tProc.Remove()
	prof, err := cpuinfo.NewProfiler()
	if err != nil {
		t.Fatal(err)
	}
	prof.Procer = tProc
	inf, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	err = testinfo.ValidateR71800xCPUInfo(inf)
	if err != nil {
		t.Error(err)
	}
	t.Log(inf)
}

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
// limitations under the License
package node

import (
	"os"
	"testing"

	"github.com/rai-project/joefriday/node"
	"github.com/rai-project/joefriday/testinfo"
)

func TestNoNode(t *testing.T) {
	prof := node.NewProfiler()
	// local relative path can be used to make sure it doesn't exist
	prof.SysFSSystemPath("")
	n, err := prof.Get()
	if !os.IsNotExist(err) {
		t.Errorf("got %q; want ErrNotExist", err)
	}
	if n != nil {
		t.Errorf("got %#v; want nil", n)
	}

}

func TestNodeX(t *testing.T) {
	tests := []struct {
		sockets      int32
		cores        int32
		threads      int32
		expectedList []string
	}{
		{1, 2, 2, []string{"0-3"}},
		{2, 8, 2, []string{"0-15", "16-31"}},
	}
	tSysFS := testinfo.NewTempSysFS()
	// use a randomly generated temp dir
	err := tSysFS.SetSysFS("")
	if err != nil {
		t.Fatalf("settiing up sysfs tree: %s", err)
	}
	defer tSysFS.Clean()
	prof := NewProfiler()
	prof.SysFSSystemPath(tSysFS.Path())

	for _, test := range tests {
		tSysFS.PhysicalPackageCount = test.sockets
		tSysFS.CoresPerPhysicalPackage = test.cores
		tSysFS.ThreadsPerCore = test.threads
		err = tSysFS.CreateNode()
		if err != nil {
			t.Errorf("%d sockets test setup: unexpected err: %s", test.sockets, err)
			continue
		}
		p, err := prof.Get()
		if err != nil {
			t.Errorf("%d sockets get: unexpected err: %s", test.sockets, err)
			continue
		}
		n := Deserialize(p)

		if int32(len(n.Node)) != test.sockets {
			t.Errorf("%d sockets: got %d nodes; want %d", test.sockets, len(n.Node), test.sockets)
			continue
		}

		for i, v := range n.Node {
			if v.CPUList != test.expectedList[i] {
				t.Errorf("%d socket test: node %d: got %q; want %q", test.sockets, i, v.CPUList, test.expectedList[i])
			}
			if v.ID != int32(i) {
				t.Errorf("%d socket test: node %d: ID: got %d; want %d", test.sockets, i, v.ID, i)
			}
		}

		tSysFS.CleanNode()
	}
}

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
	"time"

	"github.com/mohae/joefriday/cpu/utilization"
)

func TestGet(t *testing.T) {
	b, err := Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	u := Deserialize(b)
	checkUtilization(u, t)
}

func TestGetTicker(t *testing.T) {
	results := make(chan []byte)
	errs := make(chan error)
	done := make(chan struct{})
	go Ticker(time.Duration(400)*time.Millisecond, results, done, errs)
	var x int
	for {
		if x > 0 {
			close(done)
			break
		}
		select {
		case b, ok := <-results:
			if !ok {
				break
			}
			u := Deserialize(b)
			checkUtilization(u, t)
			t.Logf("%#v\n", u)
		case err := <-errs:
			t.Errorf("unexpected error: %s", err)
		}
		x++
	}
}

func checkUtilization(u *utilization.Utilization, t *testing.T) {
	if u.Timestamp == 0 {
		t.Error("timestamp: expected on-zero")
	}
	if u.CtxtDelta == 0 {
		t.Error("CtxtDelta: expected non-zero value, got 0")
	}
	if u.BTimeDelta == 0 {
		t.Error("BTimeDelta: expected non-zero value, got 0")
	}
	if u.Processes == 0 {
		t.Error("Processes: expected non-zero value, got 0")
	}
	if len(u.CPU) < 2 {
		t.Errorf("cpu: got %d, want at least 2", len(u.CPU))
	}
	for i, v := range u.CPU {
		if v.ID == "" {
			t.Errorf("%d: expected ID to have a value, was empty", i)
		}
	}
}
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

package net

import (
	"testing"
	"time"

	fb "github.com/google/flatbuffers/go"
	json "github.com/mohae/customjson"
)

func TestGetInfo(t *testing.T) {
	inf, err := GetInfo()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	// test flatbuffers stuff
	infS := inf.SerializeFlat()
	infD := DeserializeInfoFlat(infS)
	// compare
	if inf.Timestamp != infD.Timestamp {
		t.Errorf("got %d; want %d", inf.Timestamp, infD.Timestamp)
	}
	for i := 0; i < len(inf.Interfaces); i++ {
		if inf.Interfaces[i].RBytes != infD.Interfaces[i].RBytes {
			t.Errorf("%d: Rbytes: got %d; want %d", i, infD.Interfaces[i].RBytes, inf.Interfaces[i].RBytes)
		}
		if inf.Interfaces[i].RPackets != infD.Interfaces[i].RPackets {
			t.Errorf("%d: RPackets: got %d; want %d", i, infD.Interfaces[i].RPackets, inf.Interfaces[i].RPackets)
		}
		if inf.Interfaces[i].RErrs != infD.Interfaces[i].RErrs {
			t.Errorf("%d: RErrs: got %d; want %d", i, infD.Interfaces[i].RErrs, inf.Interfaces[i].RErrs)
		}
		if inf.Interfaces[i].RDrop != infD.Interfaces[i].RDrop {
			t.Errorf("%d: RDrop: got %d; want %d", i, infD.Interfaces[i].RDrop, inf.Interfaces[i].RDrop)
		}
		if inf.Interfaces[i].RFIFO != infD.Interfaces[i].RFIFO {
			t.Errorf("%d: RFIFO: got %d; want %d", i, infD.Interfaces[i].RFIFO, inf.Interfaces[i].RFIFO)
		}
		if inf.Interfaces[i].RFrame != infD.Interfaces[i].RFrame {
			t.Errorf("%d: RFrame: got %d; want %d", i, infD.Interfaces[i].RFrame, inf.Interfaces[i].RFrame)
		}
		if inf.Interfaces[i].RCompressed != infD.Interfaces[i].RCompressed {
			t.Errorf("%d: RCompressed: got %d; want %d", i, infD.Interfaces[i].RCompressed, inf.Interfaces[i].RCompressed)
		}
		if inf.Interfaces[i].RMulticast != infD.Interfaces[i].RMulticast {
			t.Errorf("%d: RMulticast: got %d; want %d", i, infD.Interfaces[i].RMulticast, inf.Interfaces[i].RMulticast)
		}
		if inf.Interfaces[i].TBytes != infD.Interfaces[i].TBytes {
			t.Errorf("%d: TBytes: got %d; want %d", i, infD.Interfaces[i].TBytes, inf.Interfaces[i].TBytes)
		}
		if inf.Interfaces[i].TPackets != infD.Interfaces[i].TPackets {
			t.Errorf("%d: TPackets: got %d; want %d", i, infD.Interfaces[i].TPackets, inf.Interfaces[i].TPackets)
		}
		if inf.Interfaces[i].TErrs != infD.Interfaces[i].TErrs {
			t.Errorf("%d: TErrs: got %d; want %d", i, infD.Interfaces[i].TErrs, inf.Interfaces[i].TErrs)
		}
		if inf.Interfaces[i].TDrop != infD.Interfaces[i].TDrop {
			t.Errorf("%d: TDrop: got %d; want %d", i, infD.Interfaces[i].TDrop, inf.Interfaces[i].TDrop)
		}
		if inf.Interfaces[i].TFIFO != infD.Interfaces[i].TFIFO {
			t.Errorf("%d: TFIFO: got %d; want %d", i, infD.Interfaces[i].TFIFO, inf.Interfaces[i].TFIFO)
		}
		if inf.Interfaces[i].TColls != infD.Interfaces[i].TColls {
			t.Errorf("%d: TColls: got %d; want %d", i, infD.Interfaces[i].TColls, inf.Interfaces[i].TColls)
		}
		if inf.Interfaces[i].TCarrier != infD.Interfaces[i].TCarrier {
			t.Errorf("%d: TCarrier: got %d; want %d", i, infD.Interfaces[i].TCarrier, inf.Interfaces[i].TCarrier)
		}
		if inf.Interfaces[i].TCompressed != infD.Interfaces[i].TCompressed {
			t.Errorf("%d: TCompressed: got %d; want %d", i, infD.Interfaces[i].TCompressed, inf.Interfaces[i].TCompressed)
		}
	}
}

func TestInfo(t *testing.T) {
	inf, _ := GetInfo()
	bldr := fb.NewBuilder(0)
	b := inf.SerializeFlatBuilder(bldr)
	infD := DeserializeInfoFlat(b)
	if json.MarshalToString(inf) != json.MarshalToString(infD) {
		t.Errorf("serialize/deserialize flatbuffers: got %v, want %v", infD, inf)
	}
}

func TestInfoTicker(t *testing.T) {
	results := make(chan []byte)
	errs := make(chan error)
	done := make(chan struct{})
	go InfoTickerFlat(time.Second, results, done, errs)
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
			inf := DeserializeInfoFlat(b)
			if len(inf.Interfaces) < 2 {
				t.Errorf("expected at least 2 interfaces, got %d", len(inf.Interfaces))
			}
		case err := <-errs:
			t.Errorf("unexpected error: %s", err)
		}
		x++
	}
}

func TestGetUsage(t *testing.T) {
	u, err := GetUsage(time.Duration(100) * time.Millisecond)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if u.Timestamp == 0 {
		t.Error("expected timestamp to have a non-zero value, it didn't")
	}
	// just check names because we can't guarantee any butes were passed
	// during the test.
	for i, v := range u.Interfaces {
		if v.Name == "" {
			t.Errorf("%d: expected the interface to have a name, it was empty", i)
		}
	}
}

func TestGetUsageSerializeDeserialize(t *testing.T) {
	u, err := GetUsage(time.Duration(100) * time.Millisecond)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	b := u.SerializeFlat()
	uD := DeserializeUsageFlat(b)
	if json.MarshalToString(u) != json.MarshalToString(uD) {
		t.Errorf("flatbuffers serialize/deserialize usage: got %v, want %v", uD, u)
	}
}

func TestUsageTicker(t *testing.T) {
	results := make(chan Usage)
	errs := make(chan error)
	done := make(chan struct{})
	go UsageTicker(time.Second, results, done, errs)
	var x int
	for {
		if x > 0 {
			close(done)
			break
		}
		select {
		case u, ok := <-results:
			if !ok {
				break
			}
			if len(u.Interfaces) < 2 {
				t.Errorf("expected at least 2 interfaces, got %d", len(u.Interfaces))
			}
			for i, v := range u.Interfaces {
				if v.Name == "" {
					t.Errorf("%d: expected name to have a value; was empty", i)
				}
			}
		case err := <-errs:
			t.Errorf("unexpected error: %s", err)
		}
		x++
	}
}

/*
func TestUsageTicker(t *testing.T) {
	results := make(chan Usage)
	errs := make(chan error)
	done := make(chan struct{})
	go UsageTicker(time.Second, results, done, errs)
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
			inf := Deserialize(b)
			if len(inf.Interfaces) < 2 {
				t.Errorf("expected at least 2 interfaces, got %d", len(inf.Interfaces))
			}
		case err := <-errs:
			t.Errorf("unexpected error: %s", err)
		}
		x++
	}
}
*/

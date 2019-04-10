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

package main

import (
	"testing"

	"github.com/mohae/benchutil"
	load "github.com/rai-project/joefriday/sysinfo/loadavg"
	lfb "github.com/rai-project/joefriday/sysinfo/loadavg/flat"
	ljson "github.com/rai-project/joefriday/sysinfo/loadavg/json"
	"github.com/rai-project/joefriday/sysinfo/mem"
	mfb "github.com/rai-project/joefriday/sysinfo/mem/flat"
	mjson "github.com/rai-project/joefriday/sysinfo/mem/json"
	"github.com/rai-project/joefriday/sysinfo/uptime"
	ufb "github.com/rai-project/joefriday/sysinfo/uptime/flat"
	ujson "github.com/rai-project/joefriday/sysinfo/uptime/json"
)

const (
	SysinfoLoadAvg = "Sysinfo LoadAvg"
	SysinfoMemInfo = "Sysinfo MemInfo"
	SysinfoUptime  = "Sysinfo Uptime"
)

func runSysinfoBenchmarks(bench benchutil.Benchmarker) {
	b := SysinfoLoadAvgGet()
	bench.Append(b)

	b = SysinfoLoadAvgGetFB()
	bench.Append(b)

	b = SysinfoLoadAvgSerializeFB()
	bench.Append(b)

	b = SysinfoLoadAvgDeserializeFB()
	bench.Append(b)

	b = SysinfoLoadAvgGetJSON()
	bench.Append(b)

	b = SysinfoLoadAvgDeserializeJSON()
	bench.Append(b)

	// Mem Info
	b = SysinfoMemInfoGet()
	bench.Append(b)

	b = SysinfoMemInfoGetFB()
	bench.Append(b)

	b = SysinfoMemInfoSerializeFB()
	bench.Append(b)

	b = SysinfoMemInfoDeserializeFB()
	bench.Append(b)

	b = SysinfoMemInfoGetJSON()
	bench.Append(b)

	b = SysinfoMemInfoDeserializeJSON()
	bench.Append(b)

	// Uptime
	b = SysinfoUptimeGet()
	bench.Append(b)

	b = SysinfoUptimeGetFB()
	bench.Append(b)

	b = SysinfoUptimeSerializeFB()
	bench.Append(b)

	b = SysinfoUptimeDeserializeFB()
	bench.Append(b)

	b = SysinfoUptimeGetJSON()
	bench.Append(b)

	b = SysinfoUptimeDeserializeJSON()
	bench.Append(b)
}

// LoadAvg
func BenchSysinfoLoadAvgGet(b *testing.B) {
	l := &load.LoadAvg{}
	for i := 0; i < b.N; i++ {
		l.Get()
	}
	_ = l
}

func SysinfoLoadAvgGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = SysinfoLoadAvg
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoLoadAvgGet))
	return bench
}

func BenchSysinfoLoadAvgGetFB(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		tmp, _ = ljson.Get()
	}
	_ = tmp
}

func SysinfoLoadAvgGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = SysinfoLoadAvg
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoLoadAvgGetFB))
	return bench
}

func BenchSysinfoLoadAvgSerializeFB(b *testing.B) {
	var tmp []byte
	l := &load.LoadAvg{}
	l.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp = lfb.Serialize(l)
	}
	_ = tmp
}

func SysinfoLoadAvgSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = SysinfoLoadAvg
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoLoadAvgSerializeFB))
	return bench
}

func BenchSysinfoLoadAvgDeserializeFB(b *testing.B) {
	var l *load.LoadAvg
	tmp, _ := lfb.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l = lfb.Deserialize(tmp)
	}
	_ = l
}

func SysinfoLoadAvgDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = SysinfoLoadAvg
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoLoadAvgDeserializeFB))
	return bench
}

func BenchSysinfoLoadAvgGetJSON(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		tmp, _ = ljson.Get()
	}
	_ = tmp
}

func SysinfoLoadAvgGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = SysinfoLoadAvg
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoLoadAvgGetJSON))
	return bench
}

func BenchSysinfoLoadAvgDeserializeJSON(b *testing.B) {
	var l *load.LoadAvg
	tmp, _ := ljson.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l, _ = ljson.Deserialize(tmp)
	}
	_ = l
}

func SysinfoLoadAvgDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = SysinfoLoadAvg
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoLoadAvgDeserializeJSON))
	return bench
}

// Mem Info
func BenchSysinfoMemInfoGet(b *testing.B) {
	m := &mem.MemInfo{}
	for i := 0; i < b.N; i++ {
		m.Get()
	}
	_ = m
}

func SysinfoMemInfoGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = SysinfoMemInfo
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoMemInfoGet))
	return bench
}

func BenchSysinfoMemInfoGetFB(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		tmp, _ = mfb.Get()
	}
	_ = tmp
}

func SysinfoMemInfoGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = SysinfoMemInfo
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoMemInfoGetFB))
	return bench
}

func BenchSysinfoMemInfoSerializeFB(b *testing.B) {
	var tmp []byte
	m := &mem.MemInfo{}
	m.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp = mfb.Serialize(m)
	}
	_ = tmp
}

func SysinfoMemInfoSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = SysinfoMemInfo
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoMemInfoSerializeFB))
	return bench
}

func BenchSysinfoMemInfoDeserializeFB(b *testing.B) {
	var m *mem.MemInfo
	tmp, _ := mfb.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m = mfb.Deserialize(tmp)
	}
	_ = m
}

func SysinfoMemInfoDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = SysinfoMemInfo
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoMemInfoDeserializeFB))
	return bench
}

func BenchSysinfoMemInfoGetJSON(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		tmp, _ = mjson.Get()
	}
	_ = tmp
}

func SysinfoMemInfoGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = SysinfoMemInfo
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoMemInfoGetJSON))
	return bench
}

func BenchSysinfoMemInfoDeserializeJSON(b *testing.B) {
	var m *mem.MemInfo
	tmp, _ := mjson.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, _ = mjson.Deserialize(tmp)
	}
	_ = m
}

func SysinfoMemInfoDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = SysinfoMemInfo
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoMemInfoDeserializeJSON))
	return bench
}

// Uptime
func BenchSysinfoUptimeGet(b *testing.B) {
	u := &uptime.Uptime{}
	for i := 0; i < b.N; i++ {
		u.Get()
	}
	_ = u
}

func SysinfoUptimeGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = SysinfoUptime
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoUptimeGet))
	return bench
}

func BenchSysinfoUptimeGetFB(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		tmp, _ = ufb.Get()
	}
	_ = tmp
}

func SysinfoUptimeGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = SysinfoUptime
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoUptimeGetFB))
	return bench
}

func BenchSysinfoUptimeSerializeFB(b *testing.B) {
	var tmp []byte
	u := &uptime.Uptime{}
	u.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp = ufb.Serialize(u)
	}
	_ = tmp
}

func SysinfoUptimeSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = SysinfoUptime
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoUptimeSerializeFB))
	return bench
}

func BenchSysinfoUptimeDeserializeFB(b *testing.B) {
	var u *uptime.Uptime
	tmp, _ := ufb.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u = ufb.Deserialize(tmp)
	}
	_ = u
}

func SysinfoUptimeDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = SysinfoUptime
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoUptimeDeserializeFB))
	return bench
}

func BenchSysinfoUptimeGetJSON(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		tmp, _ = ujson.Get()
	}
	_ = tmp
}

func SysinfoUptimeGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = SysinfoUptime
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoUptimeGetJSON))
	return bench
}

func BenchSysinfoUptimeDeserializeJSON(b *testing.B) {
	var u *uptime.Uptime
	tmp, _ := ujson.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u, _ = ujson.Deserialize(tmp)
	}
	_ = u
}

func SysinfoUptimeDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = SysinfoUptime
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoUptimeDeserializeJSON))
	return bench
}

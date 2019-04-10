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
	"github.com/rai-project/joefriday/net/netdev"
	dfb "github.com/rai-project/joefriday/net/netdev/flat"
	djson "github.com/rai-project/joefriday/net/netdev/json"
	"github.com/rai-project/joefriday/net/structs"
	"github.com/rai-project/joefriday/net/netusage"
	ufb "github.com/rai-project/joefriday/net/netusage/flat"
	ujson "github.com/rai-project/joefriday/net/netusage/json"
)

const (
	NetDev  = "Network Devices"
	NetUsage = "Network Usage"
)

func runNetBenchmarks(bench benchutil.Benchmarker) {
	b := NetDevGet()
	bench.Append(b)

	b = NetDevGetFB()
	bench.Append(b)

	b = NetDevSerializeFB()
	bench.Append(b)

	b = NetDevDeserializeFB()
	bench.Append(b)

	b = NetDevGetSON()
	bench.Append(b)

	b = NetDevSerializeJSON()
	bench.Append(b)

	b = NetDevDeserializeJSON()
	bench.Append(b)

	b = NetUsageGet()
	bench.Append(b)

	b = NetUsageGetFB()
	bench.Append(b)

	b = NetUsageSerializeFB()
	bench.Append(b)

	b = NetUsageDeserializeFB()
	bench.Append(b)

	b = NetUsageGetJSON()
	bench.Append(b)

	b = NetUsageSerializeJSON()
	bench.Append(b)

	b = NetUsageDeserializeJSON()
	bench.Append(b)
}

func BenchNetDevGet(b *testing.B) {
	var inf *structs.DevInfo
	p, _ := netdev.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = p.Get()
	}
	_ = inf
}

func NetDevGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = NetDev
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetDevGet))
	return bench
}

func BenchNetDevGetFB(b *testing.B) {
	var tmp []byte
	p, _ := dfb.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func NetDevGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = NetDev
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetDevGetFB))
	return bench
}

func BenchNetDevSerializeFB(b *testing.B) {
	var tmp []byte
	p, _ := netdev.NewProfiler()
	inf, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = dfb.Serialize(inf)
	}
	_ = tmp
}

func NetDevSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = NetDev
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetDevSerializeFB))
	return bench
}

func BenchNetDevDeserializeFB(b *testing.B) {
	var inf *structs.DevInfo
	p, _ := dfb.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inf = dfb.Deserialize(tmp)
	}
	_ = inf
}

func NetDevDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = NetDev
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetDevDeserializeFB))
	return bench
}

func BenchNetDevGetJSON(b *testing.B) {
	var tmp []byte
	p, _ := djson.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func NetDevGetSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = NetDev
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetDevGetJSON))
	return bench
}

func BenchNetDevSerializeJSON(b *testing.B) {
	var tmp []byte
	p, _ := netdev.NewProfiler()
	sts, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = djson.Serialize(sts)
	}
	_ = tmp
}

func NetDevSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = NetDev
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetDevSerializeJSON))
	return bench
}

func BenchNetDevDeserializeJSON(b *testing.B) {
	var inf *structs.DevInfo
	p, _ := djson.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = djson.Deserialize(tmp)
	}
	_ = inf
}

func NetDevDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = NetDev
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetDevDeserializeJSON))
	return bench
}

// Usage
func BenchNetUsageGet(b *testing.B) {
	var u *structs.DevUsage
	p, _ := netusage.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u, _ = p.Get()
	}
	_ = u
}

func NetUsageGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = NetUsage
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetUsageGet))
	return bench
}

func BenchNetUsageGetFB(b *testing.B) {
	var tmp []byte
	p, _ := ufb.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func NetUsageGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = NetUsage
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetUsageGetFB))
	return bench
}

func BenchNetUsageSerializeFB(b *testing.B) {
	var tmp []byte
	p, _ := netusage.NewProfiler()
	u, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = ufb.Serialize(u)
	}
	_ = tmp
}

func NetUsageSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = NetUsage
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetUsageSerializeFB))
	return bench
}

func BenchNetUsageDeserializeFB(b *testing.B) {
	var u *structs.DevUsage
	p, _ := ufb.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u = ufb.Deserialize(tmp)
	}
	_ = u
}

func NetUsageDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = NetUsage
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetUsageDeserializeFB))
	return bench
}

func BenchNetUsageGetJSON(b *testing.B) {
	var tmp []byte
	p, _ := ujson.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func NetUsageGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = NetUsage
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetUsageGetJSON))
	return bench
}

func BenchNetUsageSerializeJSON(b *testing.B) {
	var tmp []byte
	p, _ := netusage.NewProfiler()
	u, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = ujson.Serialize(u)
	}
	_ = tmp
}

func NetUsageSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = NetUsage
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetUsageSerializeJSON))
	return bench
}

func BenchNetUsageDeserializeJSON(b *testing.B) {
	var u *structs.DevUsage
	p, _ := ujson.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u, _ = ujson.Deserialize(tmp)
	}
	_ = u
}

func NetUsageDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = NetUsage
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetUsageDeserializeJSON))
	return bench
}

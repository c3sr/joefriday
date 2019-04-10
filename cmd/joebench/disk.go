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
	stats "github.com/rai-project/joefriday/disk/diskstats"
	sfb "github.com/rai-project/joefriday/disk/diskstats/flat"
	sjson "github.com/rai-project/joefriday/disk/diskstats/json"
	"github.com/rai-project/joefriday/disk/structs"
	usage "github.com/rai-project/joefriday/disk/diskusage"
	ufb "github.com/rai-project/joefriday/disk/diskusage/flat"
	ujson "github.com/rai-project/joefriday/disk/diskusage/json"
)

const (
	DiskStats = "Disk Stats"
	DiskUsage = "Disk Usage"
)

func runDiskBenchmarks(bench benchutil.Benchmarker) {
	b := DiskGetStats()
	bench.Append(b)

	b = DiskGetStatsJSON()
	bench.Append(b)

	b = DiskStatsSerializeJSON()
	bench.Append(b)

	b = DiskStatsDeserializeJSON()
	bench.Append(b)

	b = DiskGetUsage()
	bench.Append(b)

	b = DiskGetUsageFB()
	bench.Append(b)

	b = DiskUsageSerializeFB()
	bench.Append(b)

	b = DiskUsageDeserializeFB()
	bench.Append(b)

	b = DiskGetUsageJSON()
	bench.Append(b)

	b = DiskUsageSerializeJSON()
	bench.Append(b)

	b = DiskUsageDeserializeJSON()
	bench.Append(b)
}

func BenchDiskGetStats(b *testing.B) {
	var stts *structs.DiskStats
	p, _ := stats.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stts, _ = p.Get()
	}
	_ = stts
}

func DiskGetStats() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = DiskStats
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDiskGetStats))
	return bench
}

func BenchDiskGetStatsFB(b *testing.B) {
	var tmp []byte
	p, _ := sfb.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func DiskGetStatsFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = DiskStats
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDiskGetStatsFB))
	return bench
}

func BenchDiskStatsSerializeFB(b *testing.B) {
	var tmp []byte
	p, _ := stats.NewProfiler()
	sts, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = sfb.Serialize(sts)
	}
	_ = tmp
}

func DiskStatsSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = DiskStats
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDiskStatsSerializeFB))
	return bench
}

func BenchDiskStatsDeserializeFB(b *testing.B) {
	var sts *structs.DiskStats
	p, _ := sfb.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sts = sfb.Deserialize(tmp)
	}
	_ = sts
}

func DiskStatsDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = DiskStats
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDiskStatsDeserializeFB))
	return bench
}

func BenchDiskStatsGetJSON(b *testing.B) {
	var tmp []byte
	p, _ := sjson.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func DiskGetStatsJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = DiskStats
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDiskStatsGetJSON))
	return bench
}

func BenchDiskStatsSerializeJSON(b *testing.B) {
	var tmp []byte
	p, _ := stats.NewProfiler()
	sts, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = sjson.Serialize(sts)
	}
	_ = tmp
}

func DiskStatsSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = DiskStats
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDiskStatsSerializeJSON))
	return bench
}

func BenchDiskStatsDeserializeJSON(b *testing.B) {
	var sts *structs.DiskStats
	p, _ := sjson.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sts, _ = sjson.Deserialize(tmp)
	}
	_ = sts
}

func DiskStatsDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = DiskStats
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDiskStatsDeserializeJSON))
	return bench
}

// Usage
func BenchDiskGetUsage(b *testing.B) {
	var u *structs.DiskUsage
	p, _ := usage.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u, _ = p.Get()
	}
	_ = u
}

func DiskGetUsage() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = DiskUsage
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDiskGetUsage))
	return bench
}

func BenchDiskGetUsageFB(b *testing.B) {
	var tmp []byte
	p, _ := sfb.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func DiskGetUsageFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = DiskUsage
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDiskGetUsageFB))
	return bench
}

func BenchDiskUsageSerializeFB(b *testing.B) {
	var tmp []byte
	p, _ := usage.NewProfiler()
	u, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = ufb.Serialize(u)
	}
	_ = tmp
}

func DiskUsageSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = DiskUsage
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDiskUsageSerializeFB))
	return bench
}

func BenchDiskUsageDeserializeFB(b *testing.B) {
	var u *structs.DiskUsage
	p, _ := ufb.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u = ufb.Deserialize(tmp)
	}
	_ = u
}

func DiskUsageDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = DiskUsage
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDiskUsageDeserializeFB))
	return bench
}

func BenchDiskUsageGetJSON(b *testing.B) {
	var tmp []byte
	p, _ := ujson.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func DiskGetUsageJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = DiskUsage
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDiskUsageGetJSON))
	return bench
}

func BenchDiskUsageSerializeJSON(b *testing.B) {
	var tmp []byte
	p, _ := usage.NewProfiler()
	u, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = ujson.Serialize(u)
	}
	_ = tmp
}

func DiskUsageSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = DiskUsage
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDiskUsageSerializeJSON))
	return bench
}

func BenchDiskUsageDeserializeJSON(b *testing.B) {
	var u *structs.DiskUsage
	p, _ := ujson.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u, _ = ujson.Deserialize(tmp)
	}
	_ = u
}

func DiskUsageDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = DiskUsage
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDiskUsageDeserializeJSON))
	return bench
}

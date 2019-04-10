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
	"github.com/rai-project/joefriday/mem/meminfo"
	mfb "github.com/rai-project/joefriday/mem/meminfo/flat"
	mjson "github.com/rai-project/joefriday/mem/meminfo/json"
)

const (
	MemInfo = "Mem Info"
)

func runMemBenchmarks(bench benchutil.Benchmarker) {
	b := MemInfoGet()
	bench.Append(b)

	b = MemInfoGetFB()
	bench.Append(b)

	b = MemInfoSerializeFB()
	bench.Append(b)

	b = MemInfoDeserializeFB()
	bench.Append(b)

	b = MemInfoGetSON()
	bench.Append(b)

	b = MemInfoSerializeJSON()
	bench.Append(b)

	b = MemInfoDeserializeJSON()
	bench.Append(b)
}

func BenchMemInfoGet(b *testing.B) {
	var inf *meminfo.Info
	p, _ := meminfo.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = p.Get()
	}
	_ = inf
}

func MemInfoGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = MemInfo
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchMemInfoGet))
	return bench
}

func BenchMemInfoGetFB(b *testing.B) {
	var tmp []byte
	p, _ := mfb.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func MemInfoGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = MemInfo
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchMemInfoGetFB))
	return bench
}

func BenchMemInfoSerializeFB(b *testing.B) {
	var tmp []byte
	p, _ := meminfo.NewProfiler()
	inf, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = mfb.Serialize(inf)
	}
	_ = tmp
}

func MemInfoSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = MemInfo
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchMemInfoSerializeFB))
	return bench
}

func BenchMemInfoDeserializeFB(b *testing.B) {
	var inf *meminfo.Info
	p, _ := mfb.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inf = mfb.Deserialize(tmp)
	}
	_ = inf
}

func MemInfoDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = MemInfo
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchMemInfoDeserializeFB))
	return bench
}

func BenchMemInfoGetJSON(b *testing.B) {
	var tmp []byte
	p, _ := mjson.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func MemInfoGetSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = MemInfo
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchMemInfoGetJSON))
	return bench
}

func BenchMemInfoSerializeJSON(b *testing.B) {
	var tmp []byte
	p, _ := meminfo.NewProfiler()
	sts, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = mjson.Serialize(sts)
	}
	_ = tmp
}

func MemInfoSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = MemInfo
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchMemInfoSerializeJSON))
	return bench
}

func BenchMemInfoDeserializeJSON(b *testing.B) {
	var inf *meminfo.Info
	p, _ := mjson.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = mjson.Deserialize(tmp)
	}
	_ = inf
}

func MemInfoDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = MemInfo
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchMemInfoDeserializeJSON))
	return bench
}

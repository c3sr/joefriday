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
	"github.com/c3sr/joefriday/cpu/cpufreq"
	freqfb "github.com/c3sr/joefriday/cpu/cpufreq/flat"
	freqjson "github.com/c3sr/joefriday/cpu/cpufreq/json"
	"github.com/c3sr/joefriday/cpu/cpuinfo"
	infofb "github.com/c3sr/joefriday/cpu/cpuinfo/flat"
	infojson "github.com/c3sr/joefriday/cpu/cpuinfo/json"
	"github.com/c3sr/joefriday/cpu/cpustats"
	statsfb "github.com/c3sr/joefriday/cpu/cpustats/flat"
	statsjson "github.com/c3sr/joefriday/cpu/cpustats/json"
	"github.com/c3sr/joefriday/cpu/cpuutil"
	utilfb "github.com/c3sr/joefriday/cpu/cpuutil/flat"
	utiljson "github.com/c3sr/joefriday/cpu/cpuutil/json"
)

const (
	CPUFreq  = "CPU Frequency"
	CPUInfo  = "CPU Info"
	CPUStats = "CPU Stats"
	CPUUtil  = "CPU Utilization"
)

func runCPUBenchmarks(bench benchutil.Benchmarker) {
	b := CPUFreqGet()
	bench.Append(b)

	b = CPUFreqGetFB()
	bench.Append(b)

	b = CPUFreqSerializeFB()
	bench.Append(b)

	b = CPUFreqDeserializeFB()
	bench.Append(b)

	b = CPUFreqGetJSON()
	bench.Append(b)

	b = CPUFreqSerializeJSON()
	bench.Append(b)

	b = CPUFreqDeserializeJSON()
	bench.Append(b)

	b = CPUInfoGet()
	bench.Append(b)

	b = CPUInfoGetFB()
	bench.Append(b)

	b = CPUInfoSerializeFB()
	bench.Append(b)

	b = CPUInfoDeserializeFB()
	bench.Append(b)

	b = CPUInfoGetJSON()
	bench.Append(b)

	b = CPUInfoSerializeJSON()
	bench.Append(b)

	b = CPUInfoDeserializeJSON()
	bench.Append(b)

	b = CPUStatsGet()
	bench.Append(b)

	b = CPUStatsGetFB()
	bench.Append(b)

	b = CPUStatsSerializeFB()
	bench.Append(b)

	b = CPUStatsDeserializeFB()
	bench.Append(b)

	b = CPUStatsGetJSON()
	bench.Append(b)

	b = CPUStatsSerializeJSON()
	bench.Append(b)

	b = CPUStatsDeserializeJSON()
	bench.Append(b)

	b = CPUUtilGet()
	bench.Append(b)

	b = CPUUtilGetFB()
	bench.Append(b)

	b = CPUUtilSerializeFB()
	bench.Append(b)

	b = CPUUtilDeserializeFB()
	bench.Append(b)

	b = CPUUtilGetJSON()
	bench.Append(b)

	b = CPUUtilSerializeJSON()
	bench.Append(b)

	b = CPUUtilDeserializeJSON()
	bench.Append(b)
}

// Freq
func BenchCPUFreqGet(b *testing.B) {
	var f *cpufreq.Frequency
	p, _ := cpufreq.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, _ = p.Get()
	}
	_ = f
}

func CPUFreqGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = CPUFreq
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUFreqGet))
	return bench
}

func BenchCPUFreqGetFB(b *testing.B) {
	var tmp []byte
	p, _ := freqfb.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUFreqGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = CPUFreq
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUFreqGetFB))
	return bench
}

func BenchCPUFreqSerializeFB(b *testing.B) {
	var tmp []byte
	p, _ := cpufreq.NewProfiler()
	fct, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = freqfb.Serialize(fct)
	}
	_ = tmp
}

func CPUFreqSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = CPUFreq
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUFreqSerializeFB))
	return bench
}

func BenchCPUFreqDeserializeFB(b *testing.B) {
	var f *cpufreq.Frequency
	p, _ := freqfb.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f = freqfb.Deserialize(tmp)
	}
	_ = f
}

func CPUFreqDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = CPUFreq
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUFreqDeserializeFB))
	return bench
}

func BenchCPUFreqGetJSON(b *testing.B) {
	var tmp []byte
	p, _ := freqjson.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUFreqGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = CPUFreq
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUFreqGetJSON))
	return bench
}

func BenchCPUFreqSerializeJSON(b *testing.B) {
	var tmp []byte
	p, _ := freqjson.NewProfiler()
	fct, _ := p.Profiler.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = freqjson.Serialize(fct)
	}
	_ = tmp
}

func CPUFreqSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = CPUFreq
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUFreqSerializeJSON))
	return bench
}

func BenchCPUFreqDeserializeJSON(b *testing.B) {
	var f *cpufreq.Frequency
	p, _ := freqjson.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, _ = freqjson.Deserialize(tmp)
	}
	_ = f
}

func CPUFreqDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = CPUFreq
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUFreqDeserializeJSON))
	return bench
}

// Info
func BenchCPUInfoGet(b *testing.B) {
	var inf *cpuinfo.CPUInfo
	p, _ := cpuinfo.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = p.Get()
	}
	_ = inf
}

func CPUInfoGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = CPUInfo
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUInfoGet))
	return bench
}

func BenchCPUInfoGetFB(b *testing.B) {
	var tmp []byte
	p, _ := infofb.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUInfoGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = CPUInfo
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUInfoGetFB))
	return bench
}

func BenchCPUInfoSerializeFB(b *testing.B) {
	var tmp []byte
	p, _ := cpuinfo.NewProfiler()
	fct, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = infofb.Serialize(fct)
	}
	_ = tmp
}

func CPUInfoSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = CPUInfo
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUInfoSerializeFB))
	return bench
}

func BenchCPUInfoDeserializeFB(b *testing.B) {
	var inf *cpuinfo.CPUInfo
	p, _ := infofb.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inf = infofb.Deserialize(tmp)
	}
	_ = inf
}

func CPUInfoDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = CPUInfo
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUInfoDeserializeFB))
	return bench
}

func BenchCPUInfoGetJSON(b *testing.B) {
	var tmp []byte
	p, _ := infojson.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUInfoGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = CPUInfo
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUInfoGetJSON))
	return bench
}

func BenchCPUInfoSerializeJSON(b *testing.B) {
	var tmp []byte
	p, _ := infojson.NewProfiler()
	fct, _ := p.Profiler.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = infojson.Serialize(fct)
	}
	_ = tmp
}

func CPUInfoSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = CPUInfo
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUInfoSerializeJSON))
	return bench
}

func BenchCPUInfoDeserializeJSON(b *testing.B) {
	var fct *cpuinfo.CPUInfo
	p, _ := infojson.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fct, _ = infojson.Deserialize(tmp)
	}
	_ = fct
}

func CPUInfoDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = CPUInfo
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUInfoDeserializeJSON))
	return bench
}

// Stats
func BenchCPUStatsGet(b *testing.B) {
	var sts *cpustats.CPUStats
	p, _ := cpustats.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sts, _ = p.Get()
	}
	_ = sts
}

func CPUStatsGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = CPUStats
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUStatsGet))
	return bench
}

func BenchCPUStatsGetFB(b *testing.B) {
	var tmp []byte
	p, _ := statsfb.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUStatsGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = CPUStats
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUStatsGetFB))
	return bench
}

func BenchCPUStatsSerializeFB(b *testing.B) {
	var tmp []byte
	p, _ := cpustats.NewProfiler()
	sts, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = statsfb.Serialize(sts)
	}
	_ = tmp
}

func CPUStatsSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = CPUStats
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUStatsSerializeFB))
	return bench
}

func BenchCPUStatsDeserializeFB(b *testing.B) {
	var sts *cpustats.CPUStats
	p, _ := statsfb.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sts = statsfb.Deserialize(tmp)
	}
	_ = sts
}

func CPUStatsDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = CPUStats
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUStatsDeserializeFB))
	return bench
}

func BenchCPUStatsGetJSON(b *testing.B) {
	var tmp []byte
	p, _ := statsjson.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUStatsGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = CPUStats
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUStatsGetJSON))
	return bench
}

func BenchCPUStatsSerializeJSON(b *testing.B) {
	var tmp []byte
	p, _ := cpustats.NewProfiler()
	sts, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = statsjson.Serialize(sts)
	}
	_ = tmp
}

func CPUStatsSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = CPUStats
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUStatsSerializeJSON))
	return bench
}

func BenchCPUStatsDeserializeJSON(b *testing.B) {
	var sts *cpustats.CPUStats
	p, _ := statsjson.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sts, _ = statsjson.Deserialize(tmp)
	}
	_ = sts
}

func CPUStatsDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = CPUStats
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUStatsDeserializeJSON))
	return bench
}

// Utilization
func BenchCPUUtilGet(b *testing.B) {
	var u *cpuutil.CPUUtil
	p, _ := cpuutil.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u, _ = p.Get()
	}
	_ = u
}

func CPUUtilGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = CPUUtil
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilGet))
	return bench
}

func BenchCPUUtilGetFB(b *testing.B) {
	var tmp []byte
	p, _ := utilfb.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUUtilGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = CPUUtil
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilGetFB))
	return bench
}

func BenchCPUUtilSerializeFB(b *testing.B) {
	var tmp []byte
	p, _ := cpuutil.NewProfiler()
	u, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = utilfb.Serialize(u)
	}
	_ = tmp
}

func CPUUtilSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = CPUUtil
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilSerializeFB))
	return bench
}

func BenchCPUUtilDeserializeFB(b *testing.B) {
	var u *cpuutil.CPUUtil
	p, _ := utilfb.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u = utilfb.Deserialize(tmp)
	}
	_ = u
}

func CPUUtilDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = CPUUtil
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilDeserializeFB))
	return bench
}

func BenchCPUUtilGetJSON(b *testing.B) {
	var tmp []byte
	p, _ := utiljson.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUUtilGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = CPUUtil
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilGetJSON))
	return bench
}

func BenchCPUUtilSerializeJSON(b *testing.B) {
	var tmp []byte
	p, _ := cpuutil.NewProfiler()
	u, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = utiljson.Serialize(u)
	}
	_ = tmp
}

func CPUUtilSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = CPUUtil
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilSerializeJSON))
	return bench
}

func BenchCPUUtilDeserializeJSON(b *testing.B) {
	var u *cpuutil.CPUUtil
	p, _ := utiljson.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u, _ = utiljson.Deserialize(tmp)
	}
	_ = u
}

func CPUUtilDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = CPUUtil
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilDeserializeJSON))
	return bench
}

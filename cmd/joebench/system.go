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
	"github.com/rai-project/joefriday/system/loadavg"
	lfb "github.com/rai-project/joefriday/system/loadavg/flat"
	ljson "github.com/rai-project/joefriday/system/loadavg/json"
	"github.com/rai-project/joefriday/system/os"
	osfb "github.com/rai-project/joefriday/system/os/flat"
	osjson "github.com/rai-project/joefriday/system/os/json"
	"github.com/rai-project/joefriday/system/uptime"
	ufb "github.com/rai-project/joefriday/system/uptime/flat"
	ujson "github.com/rai-project/joefriday/system/uptime/json"
	"github.com/rai-project/joefriday/system/version"
	vfb "github.com/rai-project/joefriday/system/version/flat"
	vjson "github.com/rai-project/joefriday/system/version/json"
	
)

const (
	SystemVersion  = "System Version"
	SystemLoadAvg = "System LoadAvg"
	SystemOS = "System OS"
	SystemUptime  = "System Uptime"
)

func runSystemBenchmarks(bench benchutil.Benchmarker) {
	b := SystemLoadAvgGet()
	bench.Append(b)

	b = SystemLoadAvgGetFB()
	bench.Append(b)

	b = SystemLoadAvgSerializeFB()
	bench.Append(b)

	b = SystemLoadAvgDeserializeFB()
	bench.Append(b)

	b = SystemLoadAvgGetJSON()
	bench.Append(b)

	b = SystemLoadAvgSerializeJSON()
	bench.Append(b)

	b = SystemLoadAvgDeserializeJSON()
	bench.Append(b)

	b = SystemOSGet()
	bench.Append(b)

	b = SystemOSGetFB()
	bench.Append(b)

	b = SystemOSSerializeFB()
	bench.Append(b)

	b = SystemOSDeserializeFB()
	bench.Append(b)

	b = SystemOSGetJSON()
	bench.Append(b)

	b = SystemOSSerializeJSON()
	bench.Append(b)

	b = SystemOSDeserializeJSON()
	bench.Append(b)

	b = SystemUptimeGet()
	bench.Append(b)

	b = SystemUptimeGetFB()
	bench.Append(b)

	b = SystemUptimeSerializeFB()
	bench.Append(b)

	b = SystemUptimeDeserializeFB()
	bench.Append(b)

	b = SystemUptimeGetJSON()
	bench.Append(b)

	b = SystemUptimeSerializeJSON()
	bench.Append(b)

	b = SystemUptimeDeserializeJSON()
	bench.Append(b)

	b = SystemVersionGet()
	bench.Append(b)

	b = SystemVersionGetFB()
	bench.Append(b)

	b = SystemVersionSerializeFB()
	bench.Append(b)

	b = SystemVersionDeserializeFB()
	bench.Append(b)

	b = SystemVersionGetJSON()
	bench.Append(b)

	b = SystemVersionSerializeJSON()
	bench.Append(b)

	b = SystemVersionDeserializeJSON()
	bench.Append(b)
}

// LoadAvg
func BenchSystemLoadAvgGet(b *testing.B) {
	var l loadavg.LoadAvg
	p, _ := loadavg.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l, _ = p.Get()
	}
	_ = l
}

func SystemLoadAvgGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = SystemLoadAvg
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemLoadAvgGet))
	return bench
}

func BenchSystemLoadAvgGetFB(b *testing.B) {
	var tmp []byte
	p, _ := lfb.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func SystemLoadAvgGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = SystemLoadAvg
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemLoadAvgGetFB))
	return bench
}

func BenchSystemLoadAvgSerializeFB(b *testing.B) {
	var tmp []byte
	p, _ := loadavg.NewProfiler()
	l, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = lfb.Serialize(l)
	}
	_ = tmp
}

func SystemLoadAvgSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = SystemLoadAvg
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemLoadAvgSerializeFB))
	return bench
}

func BenchSystemLoadAvgDeserializeFB(b *testing.B) {
	var l loadavg.LoadAvg
	p, _ := lfb.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l = lfb.Deserialize(tmp)
	}
	_ = l
}

func SystemLoadAvgDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = SystemLoadAvg
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemLoadAvgDeserializeFB))
	return bench
}

func BenchSystemLoadAvgGetJSON(b *testing.B) {
	var tmp []byte
	p, _ := ljson.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func SystemLoadAvgGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = SystemLoadAvg
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemLoadAvgGetJSON))
	return bench
}

func BenchSystemLoadAvgSerializeJSON(b *testing.B) {
	var tmp []byte
	p, _ := loadavg.NewProfiler()
	l, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = ljson.Serialize(l)
	}
	_ = tmp
}

func SystemLoadAvgSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = SystemLoadAvg
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemLoadAvgSerializeJSON))
	return bench
}

func BenchSystemLoadAvgDeserializeJSON(b *testing.B) {
	var l loadavg.LoadAvg
	p, _ := ljson.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l, _ = ljson.Deserialize(tmp)
	}
	_ = l
}

func SystemLoadAvgDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = SystemLoadAvg
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemLoadAvgDeserializeJSON))
	return bench
}

// OS
func BenchSystemOSGet(b *testing.B) {
	var o *os.OS
	p, _ := os.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o, _ = p.Get()
	}
	_ = o
}

func SystemOSGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = SystemOS
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemOSGet))
	return bench
}

func BenchSystemOSGetFB(b *testing.B) {
	var tmp []byte
	p, _ := osfb.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func SystemOSGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = SystemOS
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemOSGetFB))
	return bench
}

func BenchSystemOSSerializeFB(b *testing.B) {
	var tmp []byte
	p, _ := os.NewProfiler()
	l, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = osfb.Serialize(l)
	}
	_ = tmp
}

func SystemOSSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = SystemOS
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemOSSerializeFB))
	return bench
}

func BenchSystemOSDeserializeFB(b *testing.B) {
	var o *os.OS
	p, _ := osfb.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o = osfb.Deserialize(tmp)
	}
	_ = o
}

func SystemOSDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = SystemOS
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemOSDeserializeFB))
	return bench
}

func BenchSystemOSGetJSON(b *testing.B) {
	var tmp []byte
	p, _ := osjson.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func SystemOSGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = SystemOS
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemOSGetJSON))
	return bench
}

func BenchSystemOSSerializeJSON(b *testing.B) {
	var tmp []byte
	p, _ := os.NewProfiler()
	l, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = osjson.Serialize(l)
	}
	_ = tmp
}

func SystemOSSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = SystemOS
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemOSSerializeJSON))
	return bench
}

func BenchSystemOSDeserializeJSON(b *testing.B) {
	var o *os.OS
	p, _ := osjson.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o, _ = osjson.Deserialize(tmp)
	}
	_ = o
}

func SystemOSDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = SystemOS
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemOSDeserializeJSON))
	return bench
}

// uptime
func BenchSystemUptimeGet(b *testing.B) {
	var u uptime.Uptime
	p, _ := uptime.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u, _ = p.Get()
	}
	_ = u
}

func SystemUptimeGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = SystemUptime
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemUptimeGet))
	return bench
}

func BenchSystemUptimeGetFB(b *testing.B) {
	var tmp []byte
	p, _ := ufb.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func SystemUptimeGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = SystemUptime
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemUptimeGetFB))
	return bench
}

func BenchSystemUptimeSerializeFB(b *testing.B) {
	var tmp []byte
	p, _ := uptime.NewProfiler()
	l, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = ufb.Serialize(l)
	}
	_ = tmp
}

func SystemUptimeSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = SystemUptime
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemUptimeSerializeFB))
	return bench
}

func BenchSystemUptimeDeserializeFB(b *testing.B) {
	var u uptime.Uptime
	p, _ := ufb.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u = ufb.Deserialize(tmp)
	}
	_ = u
}

func SystemUptimeDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = SystemUptime
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemUptimeDeserializeFB))
	return bench
}

func BenchSystemUptimeGetJSON(b *testing.B) {
	var tmp []byte
	p, _ := ujson.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func SystemUptimeGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = SystemUptime
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemUptimeGetJSON))
	return bench
}

func BenchSystemUptimeSerializeJSON(b *testing.B) {
	var tmp []byte
	p, _ := uptime.NewProfiler()
	u, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = ujson.Serialize(u)
	}
	_ = tmp
}

func SystemUptimeSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = SystemUptime
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemUptimeSerializeJSON))
	return bench
}

func BenchSystemUptimeDeserializeJSON(b *testing.B) {
	var u uptime.Uptime
	p, _ := ujson.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u, _ = ujson.Deserialize(tmp)
	}
	_ = u
}

func SystemUptimeDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = SystemUptime
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemUptimeDeserializeJSON))
	return bench
}

// Version
func BenchSystemVersionGet(b *testing.B) {
	var k *version.Kernel
	p, _ := version.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, _ = p.Get()
	}
	_ = k
}

func SystemVersionGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = SystemVersion
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemVersionGet))
	return bench
}

func BenchSystemVersionGetFB(b *testing.B) {
	var tmp []byte
	p, _ := vfb.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func SystemVersionGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = SystemVersion
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemVersionGetFB))
	return bench
}

func BenchSystemVersionSerializeFB(b *testing.B) {
	var tmp []byte
	p, _ := version.NewProfiler()
	k, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = vfb.Serialize(k)
	}
	_ = tmp
}

func SystemVersionSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = SystemVersion
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemVersionSerializeFB))
	return bench
}

func BenchSystemVersionDeserializeFB(b *testing.B) {
	var k *version.Kernel
	p, _ := vfb.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = vfb.Deserialize(tmp)
	}
	_ = k
}

func SystemVersionDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = SystemVersion
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemVersionDeserializeFB))
	return bench
}

func BenchSystemVersionGetJSON(b *testing.B) {
	var tmp []byte
	p, _ := vjson.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func SystemVersionGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = SystemVersion
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemVersionGetJSON))
	return bench
}

func BenchSystemVersionSerializeJSON(b *testing.B) {
	var tmp []byte
	p, _ := version.NewProfiler()
	k, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = vjson.Serialize(k)
	}
	_ = tmp
}

func SystemVersionSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = SystemVersion
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemVersionSerializeJSON))
	return bench
}

func BenchSystemVersionDeserializeJSON(b *testing.B) {
	var k *version.Kernel
	p, _ := vjson.NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, _ = vjson.Deserialize(tmp)
	}
	_ = k
}

func SystemVersionDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = SystemVersion
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSystemVersionDeserializeJSON))
	return bench
}

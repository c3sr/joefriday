package cpu

import (
	"testing"

	"github.com/DataDog/gohai/cpu"
	"github.com/mohae/benchutil"
	"github.com/rai-project/joefriday/cpu/cpuinfo"
	"github.com/rai-project/joefriday/cpu/cpustats"
	gopsutilcpu "github.com/shirou/gopsutil/cpu"
)

const CPUGroup = "CPU"

func BenchJoeFridayGetCPUInfo(b *testing.B) {
	var inf *cpuinfo.CPUInfo
	p, _ := cpuinfo.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = p.Get()
	}
	_ = inf
}

func JoeFridayGetCPUInfo() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/cpu/cpuinfo.Get")
	bench.Group = CPUGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetCPUInfo))
	return bench
}

func BenchJoeFridayGetCPUStats(b *testing.B) {
	var st *cpustats.CPUStats
	p, _ := cpustats.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st, _ = p.Get()
	}
	_ = st
}

func JoeFridayGetCPUStats() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/cpu/cpustats.Get")
	bench.Group = CPUGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetCPUStats))
	return bench
}

func BenchDataDogGohaiCPU(b *testing.B) {
	type Collector interface {
		Name() string
		Collect() (interface{}, error)
	}
	var collector = &cpu.Cpu{}
	var c interface{}
	for i := 0; i < b.N; i++ {
		c, _ = collector.Collect()
	}
	_ = c
}

func DataDogGohaiCPU() benchutil.Bench {
	bench := benchutil.NewBench("DataDog/gohai/cpu.Cpu.Collect")
	bench.Group = CPUGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDataDogGohaiCPU))
	return bench
}

func BenchShirouGopsutilInfoStat(b *testing.B) {
	var st []gopsutilcpu.InfoStat
	for i := 0; i < b.N; i++ {
		st, _ = gopsutilcpu.Info()
	}
	_ = st
}

func ShirouGopsutilInfoStat() benchutil.Bench {
	bench := benchutil.NewBench("shirou/gopsutil/cpu.Info")
	bench.Group = CPUGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchShirouGopsutilInfoStat))
	return bench
}

func BenchShirouGopsutilTimeStat(b *testing.B) {
	var st []gopsutilcpu.TimesStat
	for i := 0; i < b.N; i++ {
		st, _ = gopsutilcpu.Times(true)
	}
	_ = st
}

func ShirouGopsutilTimeStat() benchutil.Bench {
	bench := benchutil.NewBench("shirou/gopsutil/cpu.Times")
	bench.Group = CPUGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchShirouGopsutilTimeStat))
	return bench
}

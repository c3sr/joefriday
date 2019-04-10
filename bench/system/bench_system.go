package system

import (
	"testing"

	"github.com/DataDog/gohai/platform"
	"github.com/cloudfoundry/gosigar"
	"github.com/mohae/benchutil"
	joeloadavg "github.com/rai-project/joefriday/system/loadavg"
	joeos "github.com/rai-project/joefriday/system/os"
	joeversion "github.com/rai-project/joefriday/system/version"
	joeuptime "github.com/rai-project/joefriday/system/uptime"
	sysload "github.com/rai-project/joefriday/sysinfo/loadavg"
	sysuptime "github.com/rai-project/joefriday/sysinfo/uptime"
	"github.com/shirou/gopsutil/load"
)

const SystemGroup = "System"

func BenchJoeFridayGetVersion(b *testing.B) {
	var fct *joeversion.Kernel
	p, _ := joeversion.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fct, _ = p.Get()
	}
	_ = fct
}

func JoeFridayGetVersion() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/system/version.Get")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetVersion))
	return bench
}

func BenchJoeFridayGetOS(b *testing.B) {
	var os *joeos.OS
	p, _ := joeos.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		os, _ = p.Get()
	}
	_ = os
}

func JoeFridayGetOS() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/system/os.Get")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetOS))
	return bench
}

func BenchDataDogGohaiplatform(b *testing.B) {
	type Collector interface {
		Name() string
		Collect() (interface{}, error)
	}
	var collector = &platform.Platform{}
	var c interface{}
	for i := 0; i < b.N; i++ {
		c, _ = collector.Collect()
	}
	_ = c
}

func DataDogGohaiplatform() benchutil.Bench {
	bench := benchutil.NewBench("DataDog/gohai/platform.Platform.Collect")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDataDogGohaiplatform))
	return bench
}

func BenchJoeFridayGetLoadAvg(b *testing.B) {
	var tmp joeloadavg.LoadAvg
	p, _ := joeloadavg.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func JoeFridayGetLoadAvg() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/system/loadavg.Get")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetLoadAvg))
	return bench
}

func BenchJoeFridayGetSysinfoLoadAvg(b *testing.B) {
	var tmp sysload.LoadAvg
	for i := 0; i < b.N; i++ {
		_ = tmp.Get()
	}
	_ = tmp
}

func JoeFridayGetSysinfoLoadAvg() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/sysinfo/loadavg.LoadAvg.Get")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetSysinfoLoadAvg))
	return bench
}

func BenchCloudFoundryGoSigarLoadAverage(b *testing.B) {
	var tmp sigar.LoadAverage
	for i := 0; i < b.N; i++ {
		tmp.Get()
	}
	_ = tmp
}

func CloudFoundryGoSigarLoadAverage() benchutil.Bench {
	bench := benchutil.NewBench("cloudfoundry/gosigar.LoadAverage.Get")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCloudFoundryGoSigarLoadAverage))
	return bench
}

func BenchShirouGopsutilLoadAvg(b *testing.B) {
	var tmp *load.AvgStat
	for i := 0; i < b.N; i++ {
		tmp, _ = load.Avg()
	}
	_ = tmp
}

func ShirouGopsutilLoadAvg() benchutil.Bench {
	bench := benchutil.NewBench("shirou/gopsutil/load.Avg")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchShirouGopsutilLoadAvg))
	return bench
}

func BenchShirouGopsutilLoadMisc(b *testing.B) {
	var tmp *load.MiscStat
	for i := 0; i < b.N; i++ {
		tmp, _ = load.Misc()
	}
	_ = tmp
}

func ShirouGopsutilLoadMisc() benchutil.Bench {
	bench := benchutil.NewBench("shirou/gopsutil/load.Misc")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchShirouGopsutilLoadMisc))
	return bench
}

func BenchJoeFridayGetUptime(b *testing.B) {
	var tmp joeuptime.Uptime
	p, _ := joeuptime.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func JoeFridayGetUptime() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/system/uptime.Get")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetUptime))
	return bench
}

func BenchJoeFridayGetSysinfoUptime(b *testing.B) {
	var tmp sysuptime.Uptime
	for i := 0; i < b.N; i++ {
		_ = tmp.Get()
	}
	_ = tmp
}

func JoeFridayGetSysinfoUptime() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/sysinfo/uptime.Uptime.Get")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetSysinfoUptime))
	return bench
}

func BenchCloudFoundryGoSigarUptime(b *testing.B) {
	var tmp sigar.Uptime
	for i := 0; i < b.N; i++ {
		tmp.Get()
	}
	_ = tmp
}

func CloudFoundryGoSigarUptime() benchutil.Bench {
	bench := benchutil.NewBench("cloudfoundry/gosigar.Uptime.Get")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCloudFoundryGoSigarUptime))
	return bench
}

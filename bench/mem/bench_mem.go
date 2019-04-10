package mem

import (
	"testing"

	"github.com/DataDog/gohai/memory"
	"github.com/cloudfoundry/gosigar"
	meminfo "github.com/guillermo/go.procmeminfo"
	"github.com/mohae/benchutil"
	basic "github.com/rai-project/joefriday/mem/membasic"
	info "github.com/rai-project/joefriday/mem/meminfo"
	sysmem "github.com/rai-project/joefriday/sysinfo/mem"
	gopsutilmem "github.com/shirou/gopsutil/mem"
)

const MemGroup = "Memory"

func BenchJoeFridayGetMemBasic(b *testing.B) {
	var mem *basic.Info
	p, _ := basic.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mem, _ = p.Get()
	}
	_ = mem
}

func JoeFridayGetMemBasic() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/mem/membasic.Get")
	bench.Group = MemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetMemBasic))
	return bench
}

func BenchJoeFridayGetMemInfo(b *testing.B) {
	var mem *info.Info
	p, _ := info.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mem, _ = p.Get()
	}
	_ = mem
}

func JoeFridayGetMemInfo() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/mem/meminfo.Get")
	bench.Group = MemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetMemInfo))
	return bench
}

func BenchJoeFridayGetSysinfoMemInfo(b *testing.B) {
	var mem sysmem.MemInfo
	for i := 0; i < b.N; i++ {
		_ = mem.Get()
	}
	_ = mem
}

func JoeFridayGetSysinfoMemInfo() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/sysinfo/mem.MemInfo.Get")
	bench.Group = MemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetSysinfoMemInfo))
	return bench
}

func BenchCloudFoundryGoSigarMem(b *testing.B) {
	var mem sigar.Mem
	for i := 0; i < b.N; i++ {
		mem.Get()
	}
	_ = mem
}

func CloudFoundryGoSigarMem() benchutil.Bench {
	bench := benchutil.NewBench("cloudfoundry/gosigar.Mem.Get")
	bench.Group = MemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCloudFoundryGoSigarMem))
	return bench
}

func BenchDataDogGohaiMem(b *testing.B) {
	type Collector interface {
		Name() string
		Collect() (interface{}, error)
	}
	var collector = &memory.Memory{}
	var c interface{}
	for i := 0; i < b.N; i++ {
		c, _ = collector.Collect()
	}
	_ = c
}

func DataDogGohaiMem() benchutil.Bench {
	bench := benchutil.NewBench("DataDog/gohai/memory.Memory.Collect")
	bench.Group = MemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDataDogGohaiMem))
	return bench
}

func BenchGuillermoMemInfo(b *testing.B) {
	mem := meminfo.MemInfo{}
	for i := 0; i < b.N; i++ {
		mem.Update()
	}
	_ = mem
}

func GuillermoMemInfo() benchutil.Bench {
	bench := benchutil.NewBench("guillermo/go.procmeminfo.MemInfo.Update")
	bench.Group = MemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchGuillermoMemInfo))
	return bench
}

func BenchShirouGopsutilMem(b *testing.B) {
	var mem *gopsutilmem.VirtualMemoryStat
	for i := 0; i < b.N; i++ {
		mem, _ = gopsutilmem.VirtualMemory()
	}
	_ = mem
}

func ShirouGopsutilMem() benchutil.Bench {
	bench := benchutil.NewBench("shirou/gopsutil/mem.VirtualMemory")
	bench.Group = MemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchShirouGopsutilMem))
	return bench
}

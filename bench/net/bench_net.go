package net

import (
	"testing"

	"github.com/DataDog/gohai/network"
	"github.com/mohae/benchutil"
	"github.com/rai-project/joefriday/net/netdev"
	"github.com/rai-project/joefriday/net/structs"
	"github.com/rai-project/joefriday/net/netusage"
	gopsutilnet "github.com/shirou/gopsutil/net"
)

const NetGroup = "Network"

func BenchJoeFridayGetNetDev(b *testing.B) {
	var inf *structs.DevInfo
	p, _ := netdev.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = p.Get()
	}
	_ = inf
}

func JoeFridayGetNetDev() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/net/netdev.Get")
	bench.Group = NetGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetNetDev))
	return bench
}

func BenchJoeFridayGetNetUsage(b *testing.B) {
	var u *structs.DevUsage
	p, _ := netusage.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u, _ = p.Get()
	}
	_ = u
}

func JoeFridayGetNetUsage() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/net/netusage.Get")
	bench.Group = NetGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetNetUsage))
	return bench
}

func BenchDataDogGohaiNetwork(b *testing.B) {
	type Collector interface {
		Name() string
		Collect() (interface{}, error)
	}
	var collector = &network.Network{}
	var c interface{}
	for i := 0; i < b.N; i++ {
		c, _ = collector.Collect()
	}
	_ = c
}

func DataDogGohaiNetwork() benchutil.Bench {
	bench := benchutil.NewBench("DataDog/gohai/network")
	bench.Group = NetGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDataDogGohaiNetwork))
	return bench
}

func BenchShirouGopsutilNetInterfaces(b *testing.B) {
	var st []gopsutilnet.InterfaceStat
	for i := 0; i < b.N; i++ {
		st, _ = gopsutilnet.Interfaces()
	}
	_ = st
}

func ShirouGopsutilNetInterfaces() benchutil.Bench {
	bench := benchutil.NewBench("shirou/gopsutil/net")
	bench.Group = NetGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchShirouGopsutilNetInterfaces))
	return bench
}

func BenchShirouGopsutilIOCounters(b *testing.B) {
	var st []gopsutilnet.IOCountersStat
	for i := 0; i < b.N; i++ {
		st, _ = gopsutilnet.IOCounters(true)
	}
	_ = st
}

func ShirouGopsutilIOCounters() benchutil.Bench {
	bench := benchutil.NewBench("shirou/gopsutil/net/IOCounters")
	bench.Group = NetGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchShirouGopsutilIOCounters))
	return bench
}

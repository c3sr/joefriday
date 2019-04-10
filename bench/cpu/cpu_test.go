package cpu

import (
	"encoding/json"
	"testing"

	"github.com/DataDog/gohai/cpu"
	"github.com/rai-project/joefriday/cpu/cpuinfo"
	"github.com/rai-project/joefriday/cpu/cpustats"
	gopsutilcpu "github.com/shirou/gopsutil/cpu"
)

func BenchmarkJoeFridayGetInfo(b *testing.B) {
	var inf *cpuinfo.CPUInfo
	p, _ := cpuinfo.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = p.Get()
	}
	_ = inf
}

func BenchmarkJoeFridayGetStats(b *testing.B) {
	var st *cpustats.CPUStats
	p, _ := cpustats.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st, _ = p.Get()
	}
	_ = st
}

func BenchmarkDataDogGohaiCPU(b *testing.B) {
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

func BenchmarkShirouGopsutilInfoStat(b *testing.B) {
	var st []gopsutilcpu.InfoStat
	for i := 0; i < b.N; i++ {
		st, _ = gopsutilcpu.Info()
	}
	_ = st
}

func BenchmarkShirouGopsutilTimeStat(b *testing.B) {
	var st []gopsutilcpu.TimesStat
	for i := 0; i < b.N; i++ {
		st, _ = gopsutilcpu.Times(true)
	}
	_ = st
}

// These tests exist to print out the data that is collected; not to test the
// methods themselves.  Run with the -v flag.
func TestJoeFridayGetCPUInfo(t *testing.T) {
	prof, err := cpuinfo.NewProfiler()
	if err != nil {
		t.Error(err)
		return
	}
	inf, err := prof.Get()
	if err != nil {
		t.Error(err)
		return
	}
	p, err := json.MarshalIndent(inf, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

func TestJoeFridayGetCPUStats(t *testing.T) {
	prof, err := cpustats.NewProfiler()
	if err != nil {
		t.Error(err)
		return
	}
	st, err := prof.Get()
	if err != nil {
		t.Error(err)
		return
	}
	p, err := json.MarshalIndent(st, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

func TestDataDogGohaiCPU(t *testing.T) {
	type Collector interface {
		Name() string
		Collect() (interface{}, error)
	}
	var collector = &cpu.Cpu{}
	c, err := collector.Collect()
	if err != nil {
		t.Error(err)
		return
	}
	p, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

func TestShirouGopsutilInfoStat(t *testing.T) {
	st, err := gopsutilcpu.Info()
	if err != nil {
		t.Error(err)
		return
	}
	p, err := json.MarshalIndent(st, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

func TestShirouGopsutilTimeStat(t *testing.T) {
	st, err := gopsutilcpu.Times(true)
	if err != nil {
		t.Error(err)
		return
	}
	p, err := json.MarshalIndent(st, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

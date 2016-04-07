package mem

import (
	"encoding/json"
	"testing"

	"github.com/DataDog/gohai/memory"
	"github.com/cloudfoundry/gosigar"
	meminfo "github.com/guillermo/go.procmeminfo"
	joe "github.com/mohae/joefriday/mem"
	gopsutilmem "github.com/shirou/gopsutil/mem"
)

func BenchmarkJoeFridayGet(b *testing.B) {
	var mem *joe.Info
	b.StopTimer()
	p, _ := joe.New()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mem, _ = p.Get()
	}
	_ = mem
}

func BenchmarkCloudFoundryGoSigarMem(b *testing.B) {
	var mem sigar.Mem
	for i := 0; i < b.N; i++ {
		mem.Get()
	}
	_ = mem
}

func BenchmarkDataDogGohaiMem(b *testing.B) {
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

func BenchmarkGuillermoMemInfo(b *testing.B) {
	mem := meminfo.MemInfo{}
	for i := 0; i < b.N; i++ {
		mem.Update()
	}
	_ = mem
}

func BenchmarkShirouGopsutilMem(b *testing.B) {
	var mem *gopsutilmem.VirtualMemoryStat
	for i := 0; i < b.N; i++ {
		mem, _ = gopsutilmem.VirtualMemory()
	}
	_ = mem
}

// These tests exist to print out the data that is collected; not to test the
// methods themselves.  Run with the -v flag.
func TestJoeFridayGet(t *testing.T) {
	prof, _ := joe.New()
	mem, err := prof.Get()
	if err != nil {
		t.Error(err)
		return
	}
	p, err := json.MarshalIndent(mem, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

func TestCloudFoundryGoSigarMem(t *testing.T) {
	var mem sigar.Mem
	mem.Get()
	p, err := json.MarshalIndent(mem, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

func TestDataDogGohaiMem(t *testing.T) {
	type Collector interface {
		Name() string
		Collect() (interface{}, error)
	}
	var collector = &memory.Memory{}
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

func TestGuillermoMemInfo(t *testing.T) {
	mem := meminfo.MemInfo{}
	mem.Update()
	p, err := json.MarshalIndent(mem, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

func TestShirouGopsutilMem(t *testing.T) {
	mem, err := gopsutilmem.VirtualMemory()
	if err != nil {
		t.Error(err)
		return
	}
	p, err := json.MarshalIndent(mem, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

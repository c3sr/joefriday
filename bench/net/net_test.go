package net

import (
	"encoding/json"
	"testing"

	"github.com/DataDog/gohai/network"
	"github.com/rai-project/joefriday/net/netdev"
	"github.com/rai-project/joefriday/net/structs"
	"github.com/rai-project/joefriday/net/netusage"
	gopsutilnet "github.com/shirou/gopsutil/net"
)

func BenchmarkJoeFridayGetDevInfo(b *testing.B) {
	var inf *structs.DevInfo
	p, _ := netdev.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = p.Get()
	}
	_ = inf
}

func BenchmarkJoeFridayGetDevUsage(b *testing.B) {
	var u *structs.DevUsage
	p, _ := netusage.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u, _ = p.Get()
	}
	_ = u
}

func BenchmarkDataDogGohaiNetwork(b *testing.B) {
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

func BenchmarkShirouGopsutilNetInterfaces(b *testing.B) {
	var st []gopsutilnet.InterfaceStat
	for i := 0; i < b.N; i++ {
		st, _ = gopsutilnet.Interfaces()
	}
	_ = st
}

func BenchmarkShirouGopsutilTimeStat(b *testing.B) {
	var st []gopsutilnet.IOCountersStat
	for i := 0; i < b.N; i++ {
		st, _ = gopsutilnet.IOCounters(true)
	}
	_ = st
}

// These tests exist to print out the data that is collected; not to test the
// methods themselves.  Run with the -v flag.
func TestJoeFridayGetDevInfo(t *testing.T) {
	prof, err := netdev.NewProfiler()
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

func TestJoeFridayGetNetUsage(t *testing.T) {
	prof, err := netusage.NewProfiler()
	if err != nil {
		t.Error(err)
		return
	}
	u, err := prof.Get()
	if err != nil {
		t.Error(err)
		return
	}
	p, err := json.MarshalIndent(u, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

func TestDataDogGohaiNetwork(t *testing.T) {
	type Collector interface {
		Name() string
		Collect() (interface{}, error)
	}
	var collector = &network.Network{}
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

func TestShirouGopsutilNetInterfaces(t *testing.T) {
	st, err := gopsutilnet.Interfaces()
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
	st, err := gopsutilnet.IOCounters(true)
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

package system

import (
	"encoding/json"
	"testing"

	"github.com/DataDog/gohai/platform"
	"github.com/cloudfoundry/gosigar"
	joeloadavg "github.com/c3sr/joefriday/system/loadavg"
	joeos "github.com/c3sr/joefriday/system/os"
	joeversion "github.com/c3sr/joefriday/system/version"
	joeuptime "github.com/c3sr/joefriday/system/uptime"
	sysload "github.com/c3sr/joefriday/sysinfo/loadavg"
	sysuptime "github.com/c3sr/joefriday/sysinfo/uptime"
	"github.com/shirou/gopsutil/load"
)

func BenchmarkJoeFridayGetVersion(b *testing.B) {
	var fct *joeversion.Kernel
	p, _ := joeversion.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fct, _ = p.Get()
	}
	_ = fct
}

func BenchmarkCloudFoundryGoSigarLoadAverage(b *testing.B) {
	var tmp sigar.LoadAverage
	for i := 0; i < b.N; i++ {
		tmp.Get()
	}
	_ = tmp
}

func BenchmarkJoeFridayGetLoadAvg(b *testing.B) {
	var tmp joeloadavg.LoadAvg
	p, _ := joeloadavg.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func BenchmarkJoeFridayGetSysinfoLoadAvg(b *testing.B) {
	var tmp sysload.LoadAvg
	for i := 0; i < b.N; i++ {
		_ = tmp.Get()
	}
	_ = tmp
}

func BenchmarkShirouGopsutilLoadAvg(b *testing.B) {
	var tmp *load.AvgStat
	for i := 0; i < b.N; i++ {
		tmp, _ = load.Avg()
	}
	_ = tmp
}

func BenchmarkShirouGopsutilLoadMisc(b *testing.B) {
	var tmp *load.MiscStat
	for i := 0; i < b.N; i++ {
		tmp, _ = load.Misc()
	}
	_ = tmp
}

func BenchmarkJoeFridayGetOS(b *testing.B) {
	var st *joeos.OS
	p, _ := joeos.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st, _ = p.Get()
	}
	_ = st
}

func BenchmarkDataDogGohaiplatform(b *testing.B) {
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

func BenchmarkJoeFridayGetUptime(b *testing.B) {
	var tmp joeuptime.Uptime
	p, _ := joeuptime.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func BenchmarkJoeFridayGetSysinfoUptime(b *testing.B) {
	var tmp sysuptime.Uptime
	for i := 0; i < b.N; i++ {
		_ = tmp.Get()
	}
	_ = tmp
}

func BenchmarkCloudFoundryGoSigarUptime(b *testing.B) {
	var tmp sigar.Uptime
	for i := 0; i < b.N; i++ {
		tmp.Get()
	}
	_ = tmp
}

// These tests exist to print out the data that is collected; not to test the
// methods themselves.  Run with the -v flag.
func TestJoeFridayGetVersion(t *testing.T) {
	prof, err := joeversion.NewProfiler()
	if err != nil {
		t.Error(err)
		return
	}
	fct, err := prof.Get()
	if err != nil {
		t.Error(err)
		return
	}
	p, err := json.MarshalIndent(fct, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

func TestCloudFoundryGoSigarLoadAverage(t *testing.T) {
	var tmp sigar.LoadAverage
	tmp.Get()
	p, err := json.MarshalIndent(tmp, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

func TestJoeFridayGetLoadAvg(t *testing.T) {
	prof, err := joeloadavg.NewProfiler()
	if err != nil {
		t.Error(err)
		return
	}
	tmp, err := prof.Get()
	if err != nil {
		t.Error(err)
		return
	}
	p, err := json.MarshalIndent(tmp, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

func TestShirouGopsutilLoadAvg(t *testing.T) {
	tmp, err := load.Avg()
	if err != nil {
		t.Error(err)
		return
	}
	p, err := json.MarshalIndent(tmp, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

func TestShirouGopsutilLoadMisc(t *testing.T) {
	tmp, err := load.Misc()
	if err != nil {
		t.Error(err)
		return
	}
	p, err := json.MarshalIndent(tmp, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

func TestJoeFridayGetOS(t *testing.T) {
	prof, err := joeos.NewProfiler()
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

func TestDataDogGohaiPlatform(t *testing.T) {
	type Collector interface {
		Name() string
		Collect() (interface{}, error)
	}
	var collector = &platform.Platform{}
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

func TestCloudFoundryGoSigarUptime(t *testing.T) {
	var tmp sigar.Uptime
	tmp.Get()
	p, err := json.MarshalIndent(tmp, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

func TestJoeFridayGetUptime(t *testing.T) {
	prof, err := joeuptime.NewProfiler()
	if err != nil {
		t.Error(err)
		return
	}
	tmp, err := prof.Get()
	if err != nil {
		t.Error(err)
		return
	}
	p, err := json.MarshalIndent(tmp, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

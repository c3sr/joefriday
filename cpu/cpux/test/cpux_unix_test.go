package cpux

import (
	"encoding/json"
	"testing"

	"github.com/rai-project/joefriday/cpu/cpux"
	"github.com/rai-project/joefriday/testinfo"
)

func TestCPUX(t *testing.T) {
	// set up test stuff w/o freq
	tSysFS := testinfo.NewTempSysFS()

	// use a randomly generated temp dir
	err := tSysFS.SetSysFS("")
	if err != nil {
		t.Fatalf("settiing up sysfs tree: %s", err)
	}
	defer tSysFS.Clean()

	tSysFS.PhysicalPackageCount = 1
	tSysFS.CoresPerPhysicalPackage = 2
	tSysFS.ThreadsPerCore = 2
	err = tSysFS.CreateCPU()
	if err != nil {
		t.Error("setting up cpux testing info: %s", err)
		return
	}
	prof := &cpux.Profiler{NumCPU: int(tSysFS.CPUs())}
	prof.SysFSSystemPath(tSysFS.Path())
	cpus, err := prof.Get()
	if err != nil {
		t.Error(err)
	}

	//compare results cpufreq
	err = tSysFS.ValidateCPUX(cpus)
	if err != nil {
		t.Error(err)
	}
	js, _ := json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))

	// cleanup for next
	err = tSysFS.CleanCPU()
	if err != nil {
		t.Error(err)
	}

	// set up test stuff w/o freq
	tSysFS.Freq = false
	err = tSysFS.CreateCPU()
	if err != nil {
		t.Error("setting up cpux testing info: %s", err)
		goto multiSocket
	}
	cpus, err = prof.Get()
	if err != nil {
		t.Error(err)
	}
	js, _ = json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))
	// compare results with frequency
	err = tSysFS.ValidateCPUX(cpus)
	if err != nil {
		t.Error(err)
	}

	// cleanup for next
	err = tSysFS.CleanCPU()
	if err != nil {
		t.Error(err)
	}

multiSocket:
	// 2 sockets
	tSysFS.PhysicalPackageCount = 2
	prof.NumCPU = int(tSysFS.CPUs())
	tSysFS.Freq = true
	err = tSysFS.CreateCPU()
	if err != nil {
		t.Error("setting up cpux testing info: %s", err)
		goto noFreq
	}
	cpus, err = prof.Get()
	if err != nil {
		t.Error(err)
	}
	//compare results cpufreq
	err = tSysFS.ValidateCPUX(cpus)
	if err != nil {
		t.Error(err)
	}
	js, _ = json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))

	// cleanup for next
	err = tSysFS.CleanCPU()
	if err != nil {
		t.Error(err)
	}

noFreq:
	// set up test stuff w/o freq
	tSysFS.Freq = false
	err = tSysFS.CreateCPU()
	if err != nil {
		t.Error("setting up cpux testing info: %s", err)
		goto noOffline
	}
	cpus, err = prof.Get()
	if err != nil {
		t.Error(err)
	}
	js, _ = json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))
	// compare results with frequency
	err = tSysFS.ValidateCPUX(cpus)
	if err != nil {
		t.Error(err)
	}

	// cleanup for next
	err = tSysFS.CleanCPU()
	if err != nil {
		t.Error(err)
	}

	// no offline file
noOffline:
	tSysFS.OfflineFile = false
	// set up test stuff w/o freq
	err = tSysFS.CreateCPU()
	if err != nil {
		t.Error("setting up cpux testing info: %s", err)
		return
	}
	cpus, err = prof.Get()
	if err != nil {
		t.Error(err)
	}
	js, _ = json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))
	// compare results with frequency
	err = tSysFS.ValidateCPUX(cpus)
	if err != nil {
		t.Error(err)
	}
}

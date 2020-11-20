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

// Package cpux provides information about a system's cpus, where X is the
// integer of each CPU on the system, e.g. cpu0, cpu1, etc. On linux systems
// this comes from the sysfs filesystem. Not all paths are available on all
// systems, e.g. /sys/devices/system/cpu/cpuX/cpufreq and its children may not
// exist on some systems. If the system doesn't have a particular path within
// this path, the field's value will be the type's zero value.
//
// This package does not currently have a ticker implementation.
package cpux

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/c3sr/joefriday"
)

const (
	CPUFreq  = "cpufreq"
	Offline  = "offline"
	Online   = "online"
	Possible = "possible"
	Present  = "present"
)

type CPUs struct {
	Sockets  int32  `json:"sockets"`
	Possible string `json:"possible"`
	Online   string `json:"online"`
	Offline  string `json:"offline"`
	Present  string `json:"present"`
	CPU      []CPU  `json:"cpu"`
}

type CPU struct {
	PhysicalPackageID int32             `json:"physical_package_id"`
	CoreID            int32             `json:"core_id"`
	MHzMin            float32           `json:"mhz_min"`
	MHzMax            float32           `json:"mhz_max"`
	Cache             map[string]string `json:"cache:`
	// a sorted list of caches so that the cache info can be pulled out in order.
	CacheIDs []string `json:"cache_id"`
}

// GetCPU returns the cpu information for the provided physical_package_id
// (pID) and core_id (coreID). A false will be returned if an entry matching
// the physical_package_id and core_id is not found.
func (c *CPUs) GetCPU(pID, coreID int32) (cpu CPU, found bool) {
	for i := 0; i < len(c.CPU); i++ {
		if c.CPU[i].PhysicalPackageID == pID && c.CPU[i].CoreID == coreID {
			return c.CPU[i], true
		}
	}
	return CPU{}, false
}

// Profiler is used to process the system's cpuX information.
type Profiler struct {
	// this is an exported fied for testing purposes. It should not be set in
	// non-test usage
	NumCPU          int
	sysFSSystemPath string
	// path of the sysfs cpu tree; cached so it doesn't need to be constantly
	// redone.
	cpuPath string
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler) {
	// NumCPU provides the number of logical cpus usable by the current process.
	// Is this sufficient, or will there ever be a delta between that and either
	// what /proc/cpuinfo reports or what is available on /sys/devices/system/cpu/
	prof = &Profiler{NumCPU: runtime.NumCPU()}
	prof.SysFSSystemPath(joefriday.SysFSSystem)
	return prof
}

// Reset resources: this does nothing for this implemenation.
func (prof *Profiler) Reset() error {
	return nil
}

// Get the cpuX info for each cpu. Currently only min and max frequency are
// implemented.
func (prof *Profiler) Get() (*CPUs, error) {
	cpus := &CPUs{CPU: make([]CPU, prof.NumCPU)}
	var err error
	var pids []int32 // the physical ids encountered

	hasFreq := prof.hasCPUFreq()
	for x := 0; x < prof.NumCPU; x++ {
		var cpu CPU
		var found bool

		cpu.PhysicalPackageID, err = prof.physicalPackageID(x)
		if err != nil {
			return nil, err
		}
		// see if this is a new physical id; if so, add it to the inventory
		for _, v := range pids {
			if v == cpu.PhysicalPackageID {
				found = true
				break
			}
		}
		if !found {
			pids = append(pids, cpu.PhysicalPackageID)
		}
		cpu.CoreID, err = prof.coreID(x)
		if err != nil {
			return nil, err
		}
		err := prof.cache(x, &cpu)
		if err != nil {
			return nil, err
		}
		if hasFreq {
			cpu.MHzMin, err = prof.cpuMHzMin(x)
			if err != nil {
				return nil, err
			}
			cpu.MHzMax, err = prof.cpuMHzMax(x)
			if err != nil {
				return nil, err
			}
		}
		cpus.CPU[x] = cpu
	}
	cpus.Sockets = int32(len(pids))
	cpus.Possible, err = prof.Possible()
	if err != nil {
		return nil, err
	}

	cpus.Online, err = prof.Online()
	if err != nil {
		return nil, err
	}

	cpus.Offline, err = prof.Offline()
	if err != nil {
		return nil, err
	}

	cpus.Present, err = prof.Present()
	if err != nil {
		return nil, err
	}

	return cpus, nil
}

// cpuXPath returns the system's cpuX path for a given cpu number.
func (prof *Profiler) cpuXPath(x int) string {
	return filepath.Join(prof.cpuPath, fmt.Sprintf("cpu%d", x))
}

// coreIDPath returns the path of the core_id file for the given cpuX.
func (prof *Profiler) coreIDPath(x int) string {
	return fmt.Sprintf("%s/topology/core_id", prof.cpuXPath(x))
}

// physicalPackageIDPath returns the path of the physical_package_id file for
// the given cpuX.
func (prof *Profiler) physicalPackageIDPath(x int) string {
	return fmt.Sprintf("%s/topology/physical_package_id", prof.cpuXPath(x))
}

// cpuInfoFreqMaxPath returns the path for the cpuinfo_max_freq file of the
// given cpuX.
func (prof *Profiler) cpuInfoFreqMaxPath(x int) string {
	return fmt.Sprintf("%s/cpufreq/cpuinfo_max_freq", prof.cpuXPath(x))
}

// cpuInfoFreqMinPath returns the path for the cpuinfo_min_freq file of the
// given cpuX.
func (prof *Profiler) cpuInfoFreqMinPath(x int) string {
	return fmt.Sprintf("%s/cpufreq/cpuinfo_min_freq", prof.cpuXPath(x))
}

// cachePath returns the path for the cache dir
func (prof *Profiler) cachePath(x int) string {
	return fmt.Sprintf("%s/cache", prof.cpuXPath(x))
}

// hasCPUFreq returns if the system has cpufreq information:
func (prof *Profiler) hasCPUFreq() bool {
	_, err := os.Stat(filepath.Join(prof.cpuPath, CPUFreq))
	if err == nil {
		return true
	}
	return false
}

// gets the core_id of cpuX
func (prof *Profiler) coreID(x int) (int32, error) {
	v, err := ioutil.ReadFile(prof.coreIDPath(x))
	if err != nil {
		return 0, err
	}
	id, err := strconv.Atoi(string(v[:len(v)-1]))
	if err != nil {
		return 0, fmt.Errorf("cpu%d core_id: conversion error: %s", x, err)
	}
	return int32(id), nil
}

// gets the physical_package_id of cpuX
func (prof *Profiler) physicalPackageID(x int) (int32, error) {
	v, err := ioutil.ReadFile(prof.physicalPackageIDPath(x))
	if err != nil {
		return 0, err
	}
	id, err := strconv.Atoi(string(v[:len(v)-1]))
	if err != nil {
		return 0, fmt.Errorf("cpu%d physical_package_id: conversion error: %s", x, err)
	}
	return int32(id), nil
}

// gets the cpu_mhz_min information
func (prof *Profiler) cpuMHzMin(x int) (float32, error) {
	v, err := ioutil.ReadFile(prof.cpuInfoFreqMinPath(x))
	if err != nil {
		return 0, err
	}
	// insert the . in the appropriate spot
	v = append(v[:len(v)-4], append([]byte{'.'}, v[len(v)-4:len(v)-1]...)...)
	m, err := strconv.ParseFloat(string(v[:len(v)-1]), 32)
	if err != nil {
		return 0, fmt.Errorf("cpu%d MHz min: conversion error: %s", x, err)
	}
	return float32(m), nil
}

// gets the cpu_mhz_max information
func (prof *Profiler) cpuMHzMax(x int) (float32, error) {
	v, err := ioutil.ReadFile(prof.cpuInfoFreqMaxPath(x))
	if err != nil {
		return 0, err
	}
	// insert the . in the appropriate spot
	v = append(v[:len(v)-4], append([]byte{'.'}, v[len(v)-4:len(v)-1]...)...)
	m, err := strconv.ParseFloat(string(v[:len(v)-1]), 32)
	if err != nil {
		return 0, fmt.Errorf("cpu%d MHz max: conversion error: %s", x, err)
	}
	return float32(m), nil
}

// Get the cache info for the given cpuX entry
func (prof *Profiler) cache(x int, cpu *CPU) error {
	cpu.Cache = map[string]string{}
	//go through all the entries in cpuX/cache
	p := prof.cachePath(x)
	dirs, err := ioutil.ReadDir(p)
	if err != nil {
		return err
	}
	var cacheID string
	// all the entries should be dirs with their contents holding the cache info
	for _, d := range dirs {
		if !d.IsDir() {
			continue // this shouldn't happen but if it does we just skip the entry
		}
		// cache level
		l, err := ioutil.ReadFile(filepath.Join(p, d.Name(), "level"))
		if err != nil {
			return err
		}

		t, err := ioutil.ReadFile(filepath.Join(p, d.Name(), "type"))
		if err != nil {
			return err
		}

		// cache type: unified entries aren't decorated, otherwise the first letter is used
		// like what lscpu does.
		if t[0] != 'U' && t[0] != 'u' {
			cacheID = fmt.Sprintf("L%s%s cache", string(l[:len(l)-1]), strings.ToLower(string(t[0])))
		} else {
			cacheID = fmt.Sprintf("L%s cache", string(l[:len(l)-1]))
		}

		// cache size
		s, err := ioutil.ReadFile(filepath.Join(p, d.Name(), "size"))
		if err != nil {
			return err
		}

		// add the info
		cpu.Cache[cacheID] = string(s[:len(s)-1])
		cpu.CacheIDs = append(cpu.CacheIDs, cacheID)
	}
	// sort the cache names
	sort.Strings(cpu.CacheIDs)

	return nil
}

// Possible: CPUs that have been allocated resources and can be brought online
// if they are present. [cpu_possible_mask]
// from: Documentation/cputopology.txt
func (prof *Profiler) Possible() (string, error) {
	p, err := ioutil.ReadFile(filepath.Join(prof.cpuPath, Possible))
	if err != nil {
		return "", err
	}
	return string(p[:len(p)-1]), nil
}

// Present: CPUs that have been identified as being present in the system.
// [cpu_present_mask]
// from: Documentation/cputopology.txt
func (prof *Profiler) Present() (string, error) {
	p, err := ioutil.ReadFile(filepath.Join(prof.cpuPath, Present))
	if err != nil {
		return "", err
	}
	return string(p[:len(p)-1]), nil
}

// Online: CPUs that are online and being scheduled [cpu_online_mask]
// from: Documentation/cputopology.txt
func (prof *Profiler) Online() (string, error) {
	p, err := ioutil.ReadFile(filepath.Join(prof.cpuPath, Online))
	if err != nil {
		return "", err
	}
	return string(p[:len(p)-1]), nil
}

// Offline: CPUs that are not online because they have been HOTPLUGGED off
// (see cpu-hotplug.txt) or exceed the limit of CPUs allowed by the kernel
// configuration (kernel_max above). [~cpu_online_mask + cpus >= NR_CPUS].
// from: Documentation/cputopology.txt
//
// This file may not exist or may only contain a new line char, '\n', neither
// of these conditions are error states and will result in an empty string.
func (prof *Profiler) Offline() (string, error) {
	p, err := ioutil.ReadFile(filepath.Join(prof.cpuPath, Offline))
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(p[:len(p)-1]), nil
}

// SysFSSystemPath enables overriding the default value. This is for testing
// and should be used outside of tests.
func (prof *Profiler) SysFSSystemPath(s string) {
	prof.sysFSSystemPath = s
	prof.setCPUPath()
}

// the path to the sysfs cpu tree is cached so that it doesn't need to be
// regenerated for each use.
func (prof *Profiler) setCPUPath() {
	prof.cpuPath = filepath.Join(prof.sysFSSystemPath, "cpu")
}

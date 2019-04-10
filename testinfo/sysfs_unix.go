package testinfo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rai-project/joefriday/cpu/cpux"
	"github.com/rai-project/joefriday/node"
)

const (
	topology = "topology"
)

var (
	cache       = []string{"32K", "32K", "256K", "6144K"}
	cacheIDs    = []string{"L1d cache", "L1i cache", "L2 cache", "L3 cache"}
	cacheTypes  = []string{"Data", "Instruction", "Unified", "Unified"}
	cacheLevels = []string{"1", "1", "2", "3"}
)

// TempSysFS handles the creation of sysfs trees related to cpus and nodes in a
// temp directory for testing purposes. When usage of the temp info is done,
// Clean() should be called to remove everything that was created by Create().
// By default, the information will be created in its own temp directory within
// the system's temp dir. If the information is to be created in a specific
// location, set the Dir explicitly prior to calling Create. When the path is
// empty when Create is called, Dir will be populated with the path to the temp
// dir that represents the temp SysFS.
//
// PhysicalPackageCount, CoresPerPhysicalPackage, and ThreadsPerCore should be
// set if anything other than their default values are desired. Each of these
// values are required to be > 0. Create will not check to see if they are > 0;
// no CPUx directories will be created as the product of multiplying by 0 is 0.
//
// The number of nodes created will be equal to the PhysicalPackageCount.
//
// If the Freq flag is true, cpufreq information will be written. This
// information is not available on all systems so tests should cover both the
// cpufreq path existing and not existing.
type TempSysFS struct {
	path                    string
	Freq                    bool
	PhysicalPackageCount    int32
	CoresPerPhysicalPackage int32
	ThreadsPerCore          int32
	OfflineFile             bool
	cpuPath                 string
	nodePath                string
	tmpDir                  bool // set if the path is a randomly generated temp dir
}

// NewTempSysFS returns a new TempSysFS set to use the system's temp dir,
// populate cpufreq information with defaults of:
//   PhysicalPackageCount: 1
//   CoresPerPhysicalPackage: 2
//   ThreadsPerCore: 2
//   OfflineFile: true
//   Freq: true
func NewTempSysFS() TempSysFS {
	return TempSysFS{Freq: true, OfflineFile: true, PhysicalPackageCount: 1, CoresPerPhysicalPackage: 2, ThreadsPerCore: 2}
}

// SetSysFS sets the directory to use for generation of sysfs tree stuff. If an
// empty string is passed, a randomly generated temp dir with the prefix of
// "sysfs" will be created in the system's temp directory. If the passed string
// is not empty, it is assumed that the dir already exists.
//
// If the sysfs path was already set; calling this again will not remove the
// prior sysfs path or its contents.
func (t *TempSysFS) SetSysFS(s string) (err error) {
	if s == "" {
		s, err = ioutil.TempDir("", "sysfs")
		if err != nil {
			return err
		}
		t.tmpDir = true
	}
	t.path = s
	t.cpuPath = filepath.Join(s, "cpu")
	t.nodePath = filepath.Join(s, "node")
	return nil
}

// Returns the currently configured path
func (t *TempSysFS) Path() string {
	return t.path
}

// returns the number of CPUs per configuration:
//   PhysicalPackageCount * CoresPerPhysicalPackage * ThreadsPerCore
func (t *TempSysFS) CPUs() int32 {
	return t.PhysicalPackageCount * t.CoresPerPhysicalPackage * t.ThreadsPerCore
}

// CreateCPU creates the cpu tree for sysfs cpu tests. If the SysFS path wasn't
// set, a randomly generated directory, prefixed with sysfs, will be created in
// the system's temp directory.
func (t *TempSysFS) CreateCPU() (err error) {
	if t.path == "" {
		err = t.SetSysFS("")
		if err != nil {
			return err
		}
	}
	err = os.MkdirAll(t.cpuPath, 0777)
	if err != nil {
		return err
	}
	if t.Freq {
		// instead of always checking each cpuX dir for cpufreq information, cpuX
		// processing looks for the existence of the cpufreq path for determining
		// if cpu frequency information is available for processing
		err = os.MkdirAll(filepath.Join(t.cpuPath, cpux.CPUFreq), 0777)
		if err != nil {
			return err
		}
	}

	// add Possible information:
	err = ioutil.WriteFile(filepath.Join(t.cpuPath, cpux.Possible), []byte(fmt.Sprintf("%s\n", t.Possible())), 0777)
	if err != nil {
		return err
	}
	// add online info; use the same value as possible.
	err = ioutil.WriteFile(filepath.Join(t.cpuPath, cpux.Online), []byte(fmt.Sprintf("%s\n", t.Possible())), 0777)
	if err != nil {
		return err
	}

	// add prsent info; use the same value as possible.
	err = ioutil.WriteFile(filepath.Join(t.cpuPath, cpux.Present), []byte(fmt.Sprintf("%s\n", t.Possible())), 0777)
	if err != nil {
		return err
	}

	// if OfflineFile; create one with only a newline char.
	if t.OfflineFile {
		err = ioutil.WriteFile(filepath.Join(t.cpuPath, cpux.Offline), []byte("\n"), 0777)
		if err != nil {
			return err
		}
	}

	var x int // tracks current cpu X value

	// Add CPU info for each physical package count
	for i := 0; i < int(t.PhysicalPackageCount); i++ {
		cpusPerSocket := int(t.CoresPerPhysicalPackage * t.ThreadsPerCore)
		for j := 0; j < cpusPerSocket; j++ {
			cpuX := fmt.Sprintf("cpu%d", x)
			x++
			// set the topology core id is in reverse order of cpuX
			tmp := filepath.Join(t.cpuPath, cpuX, topology)
			err = os.MkdirAll(tmp, 0777)
			if err != nil {
				goto cleanup
			}
			err = ioutil.WriteFile(filepath.Join(tmp, "core_id"), []byte(fmt.Sprintf("%d\n", j/2)), 0777)
			if err != nil {
				goto cleanup
			}
			err = ioutil.WriteFile(filepath.Join(tmp, "physical_package_id"), []byte(fmt.Sprintf("%d\n", i)), 0777)
			if err != nil {
				goto cleanup
			}

			// cache entries
			for k := range cacheLevels {
				cD := filepath.Join(t.cpuPath, cpuX, "cache", fmt.Sprintf("index%d", k))
				err = os.MkdirAll(cD, 0777)
				if err != nil {
					goto cleanup
				}
				err = ioutil.WriteFile(filepath.Join(cD, "level"), []byte(cacheLevels[k]+"\n"), 0777)
				if err != nil {
					goto cleanup
				}
				err = ioutil.WriteFile(filepath.Join(cD, "type"), []byte(cacheTypes[k]+"\n"), 0777)
				if err != nil {
					goto cleanup
				}
				err = ioutil.WriteFile(filepath.Join(cD, "size"), []byte(cache[k]+"\n"), 0777)
				if err != nil {
					goto cleanup
				}

			}
			if !t.Freq {
				continue
			}
			tmp = filepath.Join(t.cpuPath, cpuX, cpux.CPUFreq)
			err = os.MkdirAll(tmp, 0777)
			if err != nil {
				goto cleanup
			}
			err = ioutil.WriteFile(filepath.Join(tmp, "cpuinfo_min_freq"), []byte("1600000\n"), 0777)
			if err != nil {
				goto cleanup

			}
			err = ioutil.WriteFile(filepath.Join(tmp, "cpuinfo_max_freq"), []byte("2800000\n"), 0777)
			if err != nil {
				goto cleanup

			}
		}
	}

	return nil

cleanup:
	os.RemoveAll(t.cpuPath)
	return err
}

// ValidateCPUX verifies that the info in the struct for cpuX processing is
// consistent with the test data.
func (t *TempSysFS) ValidateCPUX(cpus *cpux.CPUs) error {
	if len(cpus.CPU) != int(t.CPUs()) {
		return fmt.Errorf("CPU: got %d; want %d", len(cpus.CPU), t.CPUs())
	}
	if cpus.Sockets != t.PhysicalPackageCount {
		return fmt.Errorf("Sockets: got %d; want %d", cpus.Sockets, t.PhysicalPackageCount)
	}
	if cpus.Possible != t.Possible() {
		return fmt.Errorf("possible: got %q; want %q", cpus.Possible, t.Possible())
	}

	if cpus.Online != t.Possible() {
		return fmt.Errorf("online: got %q; want %q", cpus.Online, t.Possible())
	}

	if cpus.Present != t.Possible() {
		return fmt.Errorf("present: got %q; want %q", cpus.Present, t.Possible())
	}
	// should always be empty; this will need to change if offline files with values in them are tested.
	if cpus.Offline != "" {
		return fmt.Errorf("offline: got %q; want an empty string", cpus.Offline)
	}
	for i, cpu := range cpus.CPU {
		// find the core_id
		if cpu.CoreID < 0 || cpu.CoreID >= t.CPUs() {
			return fmt.Errorf("%d: core_id: got %d; want [0-%d]", i, cpu.CoreID, t.CPUs())
		}
		// get the cache info
		for i, v := range cpu.CacheIDs {
			if v != cacheIDs[i] {
				return fmt.Errorf("%d: got cache %s; want %s", i, v, cacheIDs[i])
			}
			c, ok := cpu.Cache[v]
			if !ok {
				return fmt.Errorf("%d: %s: expected it to exist in the cache map; it didn't", i, v)
			}
			if c != cache[i] {
				return fmt.Errorf("%d: %s: got %s; want %s", i, v, c, cache[i])
			}
		}
		if t.Freq {
			if int(cpu.MHzMax) != 2800 {
				return fmt.Errorf("%d: MHzMax: want 2800.000; got %.3f", i, cpu.MHzMax)
			}
			if int(cpu.MHzMin) != 1600 {
				return fmt.Errorf("%d: MHzMin: want 1600.000; got %.3f", i, cpu.MHzMin)
			}
		} else {
			if int(cpu.MHzMax) != 0 {
				return fmt.Errorf("%d: MHzMax: want 0.000; got %.3f", i, cpu.MHzMax)
			}
			if int(cpu.MHzMin) != 0 {
				return fmt.Errorf("%d: MHzMin: want 0.000; got %.3f", i, cpu.MHzMin)
			}
		}
	}
	return nil
}

// Clean removes the temp sysfs tree that was created during Create or
// CreateCPU. If TempSysFS created a randomly generated tmp dir, this will
// remove everything including the temp sysfs dir. If the directory to use for
// the sysfs was passed, that directory already existed and should not be
// removed; only the child trees that TempSysFS created.
//
// If you only want to clean the test data for another test, use CleanCPU
// instead.
func (t *TempSysFS) Clean() error {
	if t.tmpDir {
		err := os.RemoveAll(t.path)
		return fmt.Errorf("TempSysFS.Clean: %s: %s", t.path, err)
	}
	// The dir was passed; TempSysFS didn't create it.
	err := os.RemoveAll(t.nodePath)
	if err != nil {
		return fmt.Errorf("TempSysFS.Clean: %s: %s", t.nodePath, err)
	}
	err = os.RemoveAll(t.cpuPath)
	if err != nil {
		return fmt.Errorf("TempSysFS.Clean: %s: %s", t.cpuPath, err)
	}
	return nil
}

// CleanCPU cleans up the CPU tree that was created during CreateCPU. To clean
// everything, call Clean instead.
func (t *TempSysFS) CleanCPU() error {
	err := t.clean(t.cpuPath)
	if err != nil {
		return fmt.Errorf("TempSysFS.CleanCPU: %s", err)
	}
	return nil
}

// CleanNode cleans up the node tree that was created during CreateNode. To
// clean everything, call Clean instead.
func (t *TempSysFS) CleanNode() error {
	err := t.clean(t.nodePath)
	if err != nil {
		return fmt.Errorf("TempSysFS.CleanNode: %s", err)
	}
	return nil
}

func (t *TempSysFS) clean(path string) error {
	// otherwise get the contents of t.Dir and delete that
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("%s: nothing deleted: %s", path, err)
	}
	for _, fi := range fis {
		p := filepath.Join(path, fi.Name())
		if fi.IsDir() {
			err = os.RemoveAll(p)
			if err != nil {
				return fmt.Errorf("%s: not all files were deleted: %s", path, err)
			}
			continue
		}
		err = os.Remove(p)
		if err != nil {
			return fmt.Errorf("%s: not all files were deleted: %s", path, err)
		}
	}
	return nil
}

// Possible generates the possible string
func (t *TempSysFS) Possible() string {
	return fmt.Sprintf("0-%d", (t.PhysicalPackageCount*t.CoresPerPhysicalPackage*t.ThreadsPerCore)-1)
}

// CreateNode creates the sysfs node tree.  If the SysFS path wasn't set, a
// randomly generated directory, prefixed with sysfs, will be created in the
// system's temp directory.
func (t *TempSysFS) CreateNode() error {
	if t.path == "" {
		err := t.SetSysFS("")
		if err != nil {
			return err
		}
	}
	err := os.MkdirAll(t.nodePath, 0777)
	if err != nil {
		return err
	}
	var low int // the low end of the cpulist range
	cpusPerSocket := int(t.CoresPerPhysicalPackage * t.ThreadsPerCore)

	for i := 0; i < int(t.PhysicalPackageCount); i++ {
		nodeX := fmt.Sprintf("node%d", i)
		// set the topology core id is in reverse order of cpuX
		tmp := filepath.Join(t.nodePath, nodeX)
		err = os.MkdirAll(tmp, 0777)
		if err != nil {
			goto cleanup
		}
		err = ioutil.WriteFile(filepath.Join(tmp, node.CPUList), []byte(fmt.Sprintf("%s\n", t.cpuListString(low, (i+1)*cpusPerSocket))), 0777)
		if err != nil {
			return err
		}
		low = (i + 1) * cpusPerSocket
	}

	return nil

cleanup:
	os.RemoveAll(t.nodePath)
	return err
}

func (t *TempSysFS) cpuListString(x, y int) string {
	return fmt.Sprintf("%d-%d", x, y-1)
}

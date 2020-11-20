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

// Package cpuinfo handles processing of /proc/cpuinfo. The Info struct will
// have one entry per processor.
package cpuinfo

import (
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/c3sr/joefriday/helpers"
	joe "github.com/c3sr/joefriday"
)

const procFile = "/proc/cpuinfo"

// CPUInfo holds information about the system's cpus; CPU will have one entry
// per processor.
type CPUInfo struct {
	Timestamp int64
	Sockets   int32
	CPU       []CPU `json:"cpus"`
}

// CPU holds the /proc/cpuinfo for a single processor.
type CPU struct {
	Processor       int32    `json:"processor"`
	VendorID        string   `json:"vendor_id"`
	CPUFamily       string   `json:"cpu_family"`
	Model           string   `json:"model"`
	ModelName       string   `json:"model_name"`
	Stepping        string   `json:"stepping"`
	Microcode       string   `json:"microcode"`
	CPUMHz          float32  `json:"cpu_mhz"`
	CacheSize       string   `json:"cache_size"`
	PhysicalID      int32    `json:"physical_id"`
	Siblings        int8     `json:"siblings"`
	CoreID          int32    `json:"core_id"`
	CPUCores        int32    `json:"cpu_cores"`
	APICID          int32    `json:"apicid"`
	InitialAPICID   int32    `json:"initial_apicid"`
	FPU             string   `json:"fpu"`
	FPUException    string   `json:"fpu_exception"`
	CPUIDLevel      string   `json:"cpuid_level"`
	WP              string   `json:"wp"`
	Flags           []string `json:"flags"`
	BogoMIPS        float32  `json:"bogomips"`
	Bugs            []string `json:"bugs"`
	CLFlushSize     uint16   `json:"clflush_size"`
	CacheAlignment  uint16   `json:"cache_alignment"`
	AddressSizes    []string `json:"address_sizes"`
	PowerManagement []string `json:"power_management"`
	TLBSize         string   `json:"tlb_size"`
}

// Profiler is used to process the /proc/cpuinfo file.
type Profiler struct {
	joe.Procer
	*joe.Buffer
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	proc, err := joe.NewProc(procFile)
	if err != nil {
		return nil, err
	}
	return &Profiler{Procer: proc, Buffer: joe.NewBuffer()}, nil
}

// Reset resources; after reset the profiler is ready to be used again.
func (prof *Profiler) Reset() error {
	prof.Buffer.Reset()
	return prof.Procer.Reset()
}

// Get returns the current cpuinfo.
func (prof *Profiler) Get() (inf *CPUInfo, err error) {
	var (
		cpuCnt, i, pos, nameLen int
		n                       uint64
		physIDs                 []int32 // tracks unique physical IDs encountered
		pidFound                bool
		v                       byte
		tmp                     string
		cpu                     CPU
	)
	err = prof.Reset()
	if err != nil {
		return nil, err
	}
	inf = &CPUInfo{Timestamp: time.Now().UTC().UnixNano()}
	for {
		prof.Line, err = prof.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, &joe.ReadError{Err: err}
		}
		prof.Val = prof.Val[:0]
		// First grab the attribute name; everything up to the ':'.  The key may have
		// spaces and has trailing spaces; that gets trimmed.
		for i, v = range prof.Line {
			if v == 0x3A {
				prof.Val = prof.Line[:i]
				pos = i + 1
				break
			}
			//prof.Val = append(prof.Val, v)
		}
		prof.Val = joe.TrimTrailingSpaces(prof.Val[:])
		nameLen = len(prof.Val)
		// if there's no name; skip.
		if nameLen == 0 {
			continue
		}
		// if there's anything left, the value is everything else; trim spaces
		if pos+1 < len(prof.Line) {
			prof.Val = append(prof.Val, joe.TrimTrailingSpaces(prof.Line[pos+1:])...)
		}
		v = prof.Val[0]
		if v == 'a' {
			v = prof.Val[1]
			if v == 'd' { // address sizes
				cpu.AddressSizes = strings.Split(string(prof.Val[nameLen:]), ", ")
				continue
			}
			if v == 'p' { // apicid
				n, err = helpers.ParseUint(prof.Val[nameLen:])
				if err != nil {
					return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}
				cpu.APICID = int32(n)
			}
			continue
		}
		if v == 'c' {
			v = prof.Val[1]
			if v == 'p' {
				v = prof.Val[4]
				if v == 'c' { // cpu cores
					n, err = helpers.ParseUint(prof.Val[nameLen:])
					if err != nil {
						return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
					}
					cpu.CPUCores = int32(n)
					continue
				}
				if v == 'f' { // cpu family
					cpu.CPUFamily = string(prof.Val[nameLen:])
					continue
				}
				if v == 'M' { // cpu MHz
					f, err := strconv.ParseFloat(string(prof.Val[nameLen:]), 32)
					if err != nil {
						return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
					}
					cpu.CPUMHz = float32(f)
					continue
				}
				if v == 'd' { // cpuid level
					cpu.CPUIDLevel = string(prof.Val[nameLen:])
				}
				continue
			}
			v = prof.Val[5]
			if v == '_' { // cache_alignment
				n, err = helpers.ParseUint(prof.Val[nameLen:])
				if err != nil {
					return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}

				cpu.CacheAlignment = uint16(n)
				continue
			}
			if v == ' ' { // cache size
				cpu.CacheSize = string(prof.Val[nameLen:])
				continue
			}
			if v == 's' { // clflush size
				n, err = helpers.ParseUint(prof.Val[nameLen:])
				if err != nil {
					return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}
				cpu.CLFlushSize = uint16(n)
				continue
			}
			if v == 'i' { // core id
				n, err = helpers.ParseUint(prof.Val[nameLen:])
				if err != nil {
					return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}
				cpu.CoreID = int32(n)
			}
			continue
		}
		if v == 'f' {
			v = prof.Val[1]
			if v == 'l' { // flags
				cpu.Flags = strings.Split(string(prof.Val[nameLen:]), " ")
				continue
			}
			if v == 'p' {
				if nameLen == 3 { // fpu
					cpu.FPU = string(prof.Val[nameLen:])
				} else { // fpu_exception
					cpu.FPUException = string(prof.Val[nameLen:])
				}
			}
			continue
		}
		if v == 'm' {
			v = prof.Val[1]
			if v == 'i' { // microcode
				cpu.Microcode = string(prof.Val[nameLen:])
				continue
			}
			if v == 'o' {
				if nameLen == 5 { // model
					cpu.Model = string(prof.Val[nameLen:])
					continue
				}
				cpu.ModelName = string(prof.Val[nameLen:])
			}
			continue
		}
		if v == 'p' {
			v = prof.Val[1]
			if v == 'h' { // physical id
				n, err = helpers.ParseUint(prof.Val[nameLen:])
				if err != nil {
					return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}
				cpu.PhysicalID = int32(n)
				for i := range physIDs {
					if physIDs[i] == cpu.PhysicalID {
						pidFound = true
						break
					}
				}
				if pidFound {
					pidFound = false // reset for next use
				} else {
					// physical id hasn't been encountered yet; add it
					physIDs = append(physIDs, cpu.PhysicalID)
				}
				continue
			}
			if v == 'o' { // power management
				tmp = string(prof.Val[nameLen:])
				if tmp == "" {
					continue
				}
				cpu.PowerManagement = strings.Split(tmp, " ")
				continue
			}
			// processor starts information about a processor.
			if v == 'r' { // processor
				if cpuCnt > 0 {
					inf.CPU = append(inf.CPU, cpu)
				}
				cpuCnt++
				n, err = helpers.ParseUint(prof.Val[nameLen:])
				if err != nil {
					return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}
				cpu = CPU{Processor: int32(n)}
			}
			continue
		}
		if v == 's' {
			v = prof.Val[1]
			if v == 'i' { // siblings
				n, err = helpers.ParseUint(prof.Val[nameLen:])
				if err != nil {
					return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}
				cpu.Siblings = int8(n)
				continue
			}
			if v == 't' { // stepping
				cpu.Stepping = string(prof.Val[nameLen:])
			}
			continue
		}
		if v == 'b' {
			if prof.Val[1] == 'o' { // bogomips
				f, err := strconv.ParseFloat(string(prof.Val[nameLen:]), 32)
				if err != nil {
					return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}
				cpu.BogoMIPS = float32(f)
				continue
			}
			if prof.Val[1] == 'u' { // bugs
				tmp = string(prof.Val[nameLen:])
				if tmp != "" {
					cpu.Bugs = strings.Split(tmp, " ")
				}
			}
			continue
		}
		if v == 'i' { // initial apicid
			n, err = helpers.ParseUint(prof.Val[nameLen:])
			if err != nil {
				return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
			}
			cpu.InitialAPICID = int32(n)
			continue
		}
		if v == 'w' { // WP
			cpu.WP = string(prof.Val[nameLen:])
			continue
		}
		if v == 'v' { // vendor_id
			cpu.VendorID = string(prof.Val[nameLen:])
		}
		if v == 'T' { //tlb size
			cpu.TLBSize = string(prof.Val[nameLen:])
		}
	}
	// append the current processor informatin
	inf.CPU = append(inf.CPU, cpu)
	inf.Sockets = int32(len(physIDs))
	return inf, nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current cpuinfo using the package's global Profiler.
func Get() (inf *CPUInfo, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Get()
}

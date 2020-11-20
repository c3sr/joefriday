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

// Package version gets the kernel and version information from the
// /proc/version file.
package version

import (
	"io"
	"sync"

	joe "github.com/c3sr/joefriday"
)

const procFile = "/proc/version"

// Kernel holds information about the kernel and version.
type Kernel struct {
	OS          string `json:"os"`
	Version     string `json:"version"`
	CompileUser string `json:"compile_user"`
	GCC         string `json:"gcc"`
	OSGCC       string `json:"os_gcc"`
	Type        string `json:"type"`
	CompileDate string `json:"compile_date"`
	Arch        string `json:"arch"`
}

// Profiler processes the version information, /proc/version.
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

// Reset resources: after reset, the profiler is ready to be used again.
func (prof *Profiler) Reset() error {
	prof.Buffer.Reset()
	return prof.Procer.Reset()
}

// Get gets the kernel information from the /proc/version file.
func (prof *Profiler) Get() (k *Kernel, err error) {
	var (
		i, pos, pos2 int
		v            byte
	)
	err = prof.Reset()
	if err != nil {
		return nil, err
	}
	// This will always be linux, I think.
	k = &Kernel{OS: "linux"}
	for {
		prof.Line, err = prof.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, &joe.ReadError{Err: err}
		}
		// The version is everything from the space, 0x20, prior to the version string, up to the first '(', 0x28, - 1 byte
		for i, v = range prof.Line {
			if v == 0x28 {
				// get the OS
				k.Version = string(prof.Line[pos2+1 : i-1])
				pos = i + 1
				break
			}
			// keep track of the last space encountered
			if v == 0x20 {
				pos2 = pos
				pos = i
			}
		}
		// Set the arch
		k.SetArch()
		// The CompileUser is everything up to the next ')', 0x29
		for i, v = range prof.Line[pos:] {
			if v == 0x29 {
				k.CompileUser = string(prof.Line[pos : pos+i])
				pos += i + 3
				break
			}
		}

		var inOSGCC bool
		// GCC info; this may include os specific gcc info
		for i, v = range prof.Line[pos:] {
			if v == 0x28 {
				inOSGCC = true
				k.GCC = string(prof.Line[pos : pos+i-1])
				pos2 = i + pos + 1
				continue
			}
			if v == 0x29 {
				if inOSGCC {
					k.OSGCC = string(prof.Line[pos2 : pos+i])
					inOSGCC = false
					continue
				}
				pos, pos2 = pos+i+2, pos
				break
			}
		}
		// Check if GCC is empty, this happens if there wasn't an OSGCC value
		if k.GCC == "" {
			k.GCC = string(prof.Line[pos2 : pos-1])
		}
		// Get the type information, everything up to '('
		for i, v = range prof.Line[pos:] {
			if v == 0x28 {
				k.Type = string(prof.Line[pos : pos+i-1])
				pos += i + 1
				break
			}
		}
		// The rest is the compile date.
		k.CompileDate = string(prof.Line[pos : len(prof.Line)-2])
	}
	return k, nil
}

var std *Profiler
var stdMu sync.Mutex

// Get gets the kernel information from the /proc/version file using the
// package's global profiler.
func Get() (k *Kernel, err error) {
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

// Set the Version's architecture information.  This is the last segment of
// the Version.
func (k *Kernel) SetArch() {
	// get everything after the last -
	for i := len(k.Version) - 1; i > 0; i-- {
		if k.Version[i] == '-' {
			k.Arch = string(k.Version[i+1:])
			return
		}
	}
}

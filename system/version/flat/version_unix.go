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
// /proc/version file. Instead of returning a Go struct, it returns Flatbuffer
// serialized bytes. A function to deserialize the Flatbuffer serialized bytes
// into a version.Kernel struct is provided.
//
// Note: the package name is version and not the final element of the import
// path (flat). 
package version

import (
	"sync"

	fb "github.com/google/flatbuffers/go"
	v "github.com/rai-project/joefriday/system/version"
	"github.com/rai-project/joefriday/system/version/flat/structs"
)

// Profiler processes the version information, /proc/version, using
// Flatbuffers.
type Profiler struct {
	*v.Profiler
	*fb.Builder
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	p, err := v.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get gets the kernel information from the /proc/version file as Flatbuffer
// serialized bytes.
func (prof *Profiler) Get() ([]byte, error) {
	inf, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf), nil
}

var std *Profiler
var stdMu sync.Mutex //protects standard to prevent a data race on checking/instantiation

// Get gets the kernel information from the /proc/version file as Flatbuffer
// serialized bytes using the package's global Profiler.
func Get() (p []byte, err error) {
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

// Serialize version.Kernel as Flatbuffers.
func (prof *Profiler) Serialize(k *v.Kernel) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	os := prof.Builder.CreateString(k.OS)
	version := prof.Builder.CreateString(k.Version)
	compileUser := prof.Builder.CreateString(k.CompileUser)
	gcc := prof.Builder.CreateString(k.GCC)
	osgcc := prof.Builder.CreateString(k.OSGCC)
	typ := prof.Builder.CreateString(k.Type)
	compileDate := prof.Builder.CreateString(k.CompileDate)
	arch := prof.Builder.CreateString(k.Arch)
	structs.KernelStart(prof.Builder)
	structs.KernelAddOS(prof.Builder, os)
	structs.KernelAddVersion(prof.Builder, version)
	structs.KernelAddCompileUser(prof.Builder, compileUser)
	structs.KernelAddGCC(prof.Builder, gcc)
	structs.KernelAddOSGCC(prof.Builder, osgcc)
	structs.KernelAddType(prof.Builder, typ)
	structs.KernelAddCompileDate(prof.Builder, compileDate)
	structs.KernelAddArch(prof.Builder, arch)
	prof.Builder.Finish(structs.KernelEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize version.Kernel as Flatbuffers using the package's global Profiler.
func Serialize(k *v.Kernel) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(k), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserializes them as
// version.Kernel.
func Deserialize(p []byte) *v.Kernel {
	flatK := structs.GetRootAsKernel(p, 0)
	var k v.Kernel
	k.OS = string(flatK.OS())
	k.Version = string(flatK.Version())
	k.CompileUser = string(flatK.CompileUser())
	k.GCC = string(flatK.GCC())
	k.OSGCC = string(flatK.OSGCC())
	k.Type = string(flatK.Type())
	k.CompileDate = string(flatK.CompileDate())
	k.Arch = string(flatK.Arch())
	return &k
}

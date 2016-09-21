// Copyright 2016 The JoeFriday authors.
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

package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mohae/benchutil"
	"github.com/mohae/joefriday/bench/cpu"
	"github.com/mohae/joefriday/bench/mem"
	"github.com/mohae/joefriday/bench/net"
	"github.com/mohae/joefriday/bench/platform"
)

// flags
var (
	output         string
	format         string
	section        bool
	sectionHeaders bool
	nameSections   bool
	systemInfo     bool
)

func init() {
	flag.StringVar(&output, "output", "stdout", "output destination")
	flag.StringVar(&output, "o", "stdout", "output destination (short)")
	flag.StringVar(&format, "format", "txt", "format of output")
	flag.StringVar(&format, "f", "txt", "format of output")
	flag.BoolVar(&nameSections, "namesections", false, "use group as section name: some restrictions apply")
	flag.BoolVar(&nameSections, "n", false, "use group as section name: some restrictions apply")
	flag.BoolVar(&section, "sections", false, "don't separate groups of tests into sections")
	flag.BoolVar(&section, "s", false, "don't separate groups of tests into sections")
	flag.BoolVar(&sectionHeaders, "sectionheader", false, "if there are sections, add a section header row")
	flag.BoolVar(&sectionHeaders, "h", false, "if there are sections, add a section header row")
	flag.BoolVar(&systemInfo, "sysinfo", false, "add the system information to the output")
	flag.BoolVar(&systemInfo, "i", false, "add the system information to the output")
}

func main() {
	flag.Parse()
	done := make(chan struct{})
	go benchutil.Dot(done)

	// set the output
	var w io.Writer
	var err error
	switch output {
	case "stdout":
		w = os.Stdout
	default:
		w, err = os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer w.(*os.File).Close()
	}
	// get the benchmark for the desired format
	// process the output
	var bench benchutil.Benchmarker
	switch format {
	case "csv":
		bench = benchutil.NewCSVBench(w)
	case "md":
		bench = benchutil.NewMDBench(w)
		bench.NameSections(nameSections)
	default:
		bench = benchutil.NewStringBench(w)
	}
	bench.SectionPerGroup(section)
	bench.SectionHeaders(sectionHeaders)
	bench.IncludeSystemInfo(systemInfo)
	// CPU
	runCPUBenches(bench)

	// Mem
	runMemBenches(bench)

	// Net
	runNetBenches(bench)

	// Platform
	runPlatformBenches(bench)

	fmt.Println("\ngenerating output...\n")
	err = bench.Out()
	if err != nil {
		fmt.Printf("error generating output: %s\n", err)
	}
}

func runCPUBenches(bench benchutil.Benchmarker) {
	b := cpu.JoeFridayGetFacts()
	bench.Append(b)

	b = cpu.JoeFridayGetStats()
	bench.Append(b)

	b = cpu.DataDogGohaiCPU()
	bench.Append(b)

	b = cpu.ShirouGopsutilInfoStat()
	bench.Append(b)

	b = cpu.ShirouGopsutilTimeStat()
	bench.Append(b)
}

func runMemBenches(bench benchutil.Benchmarker) {
	b := mem.JoeFridayGetMemInfo()
	bench.Append(b)

	b = mem.JoeFridayGetSysinfoMemInfo()
	bench.Append(b)

	b = mem.CloudFoundryGoSigarMem()
	bench.Append(b)

	b = mem.DataDogGohaiMem()
	bench.Append(b)

	b = mem.GuillermoMemInfo()
	bench.Append(b)

	b = mem.ShirouGopsutilMem()
	bench.Append(b)
}

func runNetBenches(bench benchutil.Benchmarker) {
	b := net.JoeFridayGetInfo()
	bench.Append(b)

	b = net.JoeFridayGetUsage()
	bench.Append(b)

	b = net.DataDogGohaiNetwork()
	bench.Append(b)

	b = net.ShirouGopsutilNetInterfaces()
	bench.Append(b)

	b = net.ShirouGopsutilIOCounters()
	bench.Append(b)
}

func runPlatformBenches(bench benchutil.Benchmarker) {
	b := platform.JoeFridayGetKernel()
	bench.Append(b)

	b = platform.JoeFridayGetRelease()
	bench.Append(b)

	b = platform.DataDogGohaiplatform()
	bench.Append(b)

	b = platform.JoeFridayGetLoadAvg()
	bench.Append(b)

	b = platform.JoeFridayGetSysinfoLoadAvg()
	bench.Append(b)

	b = platform.CloudFoundryGoSigarLoadAverage()
	bench.Append(b)

	b = platform.ShirouGopsutilLoadAvg()
	bench.Append(b)

	b = platform.ShirouGopsutilLoadMisc()
	bench.Append(b)

	b = platform.JoeFridayGetUptime()
	bench.Append(b)

	b = platform.JoeFridayGetSysinfoUptime()
	bench.Append(b)

	b = platform.CloudFoundryGoSigarUptime()
	bench.Append(b)
}

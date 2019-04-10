package testinfo

import (
	"errors"
	"fmt"

	"github.com/rai-project/joefriday/cpu/cpufreq"
	"github.com/rai-project/joefriday/cpu/cpuinfo"
	"github.com/rai-project/joefriday/processors"
)

const XeonE52690ModelName = "Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz"

var XeonE52690CPUInfo = []byte(`processor	: 0
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 3273.714
cache size	: 20480 KB
physical id	: 0
siblings	: 16
core id		: 0
cpu cores	: 8
apicid		: 0
initial apicid	: 0
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5786.25
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 1
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 3431.175
cache size	: 20480 KB
physical id	: 0
siblings	: 16
core id		: 1
cpu cores	: 8
apicid		: 2
initial apicid	: 2
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5786.25
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 2
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 1669.085
cache size	: 20480 KB
physical id	: 0
siblings	: 16
core id		: 2
cpu cores	: 8
apicid		: 4
initial apicid	: 4
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5786.25
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 3
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 1637.820
cache size	: 20480 KB
physical id	: 0
siblings	: 16
core id		: 3
cpu cores	: 8
apicid		: 6
initial apicid	: 6
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5786.25
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 4
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 3548.648
cache size	: 20480 KB
physical id	: 0
siblings	: 16
core id		: 4
cpu cores	: 8
apicid		: 8
initial apicid	: 8
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5786.25
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 5
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 1341.136
cache size	: 20480 KB
physical id	: 0
siblings	: 16
core id		: 5
cpu cores	: 8
apicid		: 10
initial apicid	: 10
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5786.25
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 6
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 3370.796
cache size	: 20480 KB
physical id	: 0
siblings	: 16
core id		: 6
cpu cores	: 8
apicid		: 12
initial apicid	: 12
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5786.25
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 7
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 1663.308
cache size	: 20480 KB
physical id	: 0
siblings	: 16
core id		: 7
cpu cores	: 8
apicid		: 14
initial apicid	: 14
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5786.25
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 8
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 1590.808
cache size	: 20480 KB
physical id	: 1
siblings	: 16
core id		: 0
cpu cores	: 8
apicid		: 32
initial apicid	: 32
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5787.75
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 9
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2498.871
cache size	: 20480 KB
physical id	: 1
siblings	: 16
core id		: 1
cpu cores	: 8
apicid		: 34
initial apicid	: 34
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5787.75
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 10
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 1593.414
cache size	: 20480 KB
physical id	: 1
siblings	: 16
core id		: 2
cpu cores	: 8
apicid		: 36
initial apicid	: 36
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5787.75
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 11
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2679.214
cache size	: 20480 KB
physical id	: 1
siblings	: 16
core id		: 3
cpu cores	: 8
apicid		: 38
initial apicid	: 38
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5787.75
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 12
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2923.222
cache size	: 20480 KB
physical id	: 1
siblings	: 16
core id		: 4
cpu cores	: 8
apicid		: 40
initial apicid	: 40
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5787.75
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 13
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 1620.828
cache size	: 20480 KB
physical id	: 1
siblings	: 16
core id		: 5
cpu cores	: 8
apicid		: 42
initial apicid	: 42
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5787.75
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 14
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 1872.199
cache size	: 20480 KB
physical id	: 1
siblings	: 16
core id		: 6
cpu cores	: 8
apicid		: 44
initial apicid	: 44
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5787.75
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 15
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 1989.445
cache size	: 20480 KB
physical id	: 1
siblings	: 16
core id		: 7
cpu cores	: 8
apicid		: 46
initial apicid	: 46
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5787.75
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 16
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 3293.199
cache size	: 20480 KB
physical id	: 0
siblings	: 16
core id		: 0
cpu cores	: 8
apicid		: 1
initial apicid	: 1
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5786.25
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 17
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 1551.273
cache size	: 20480 KB
physical id	: 0
siblings	: 16
core id		: 1
cpu cores	: 8
apicid		: 3
initial apicid	: 3
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5786.25
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 18
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 3479.433
cache size	: 20480 KB
physical id	: 0
siblings	: 16
core id		: 2
cpu cores	: 8
apicid		: 5
initial apicid	: 5
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5786.25
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 19
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2072.140
cache size	: 20480 KB
physical id	: 0
siblings	: 16
core id		: 3
cpu cores	: 8
apicid		: 7
initial apicid	: 7
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5786.25
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 20
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 3321.292
cache size	: 20480 KB
physical id	: 0
siblings	: 16
core id		: 4
cpu cores	: 8
apicid		: 9
initial apicid	: 9
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5786.25
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 21
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 3296.371
cache size	: 20480 KB
physical id	: 0
siblings	: 16
core id		: 5
cpu cores	: 8
apicid		: 11
initial apicid	: 11
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5786.25
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 22
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2182.703
cache size	: 20480 KB
physical id	: 0
siblings	: 16
core id		: 6
cpu cores	: 8
apicid		: 13
initial apicid	: 13
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5786.25
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 23
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2269.023
cache size	: 20480 KB
physical id	: 0
siblings	: 16
core id		: 7
cpu cores	: 8
apicid		: 15
initial apicid	: 15
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5786.25
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 24
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 1465.632
cache size	: 20480 KB
physical id	: 1
siblings	: 16
core id		: 0
cpu cores	: 8
apicid		: 33
initial apicid	: 33
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5787.75
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 25
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 1644.503
cache size	: 20480 KB
physical id	: 1
siblings	: 16
core id		: 1
cpu cores	: 8
apicid		: 35
initial apicid	: 35
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5787.75
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 26
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 1899.839
cache size	: 20480 KB
physical id	: 1
siblings	: 16
core id		: 2
cpu cores	: 8
apicid		: 37
initial apicid	: 37
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5787.75
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 27
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 1643.257
cache size	: 20480 KB
physical id	: 1
siblings	: 16
core id		: 3
cpu cores	: 8
apicid		: 39
initial apicid	: 39
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5787.75
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 28
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 1515.250
cache size	: 20480 KB
physical id	: 1
siblings	: 16
core id		: 4
cpu cores	: 8
apicid		: 41
initial apicid	: 41
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5787.75
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 29
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 1883.414
cache size	: 20480 KB
physical id	: 1
siblings	: 16
core id		: 5
cpu cores	: 8
apicid		: 43
initial apicid	: 43
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5787.75
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 30
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 1235.332
cache size	: 20480 KB
physical id	: 1
siblings	: 16
core id		: 6
cpu cores	: 8
apicid		: 45
initial apicid	: 45
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5787.75
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 31
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2690 0 @ 2.90GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 1276.906
cache size	: 20480 KB
physical id	: 1
siblings	: 16
core id		: 7
cpu cores	: 8
apicid		: 47
initial apicid	: 47
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm epb tpr_shadow vnmi flexpriority ept vpid xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5787.75
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:
`)

// ValidateXeonE52690CPUInfo verifies that the info in the struct info is
// consistent with the above data. If everything verifies a nil is returned,
// otherwise an error is returned. This is used for testing.
//
// processor, core id, apicid, initial apicid fields are not checked as they
// vary per entry.
func ValidateXeonE52690CPUInfo(inf *cpuinfo.CPUInfo) error {
	if inf.Timestamp == 0 {
		return errors.New("expected Timestamp to have a non-zero value; it didn't")
	}
	if len(inf.CPU) != 32 {
		return fmt.Errorf("CPU: got %d; want 32", len(inf.CPU))
	}
	if inf.Sockets != 2 {
		return fmt.Errorf("got %d socket; want 2", len(inf.CPU))
	}

	for i, cpu := range inf.CPU {
		if cpu.VendorID != GenuineIntel {
			return fmt.Errorf("%d: vendor_id: got %q; want %q", i, cpu.VendorID, GenuineIntel)
		}
		if cpu.CPUFamily != "6" {
			return fmt.Errorf("%d: cpu family: got %q; want \"6\"", i, cpu.CPUFamily)
		}
		if cpu.Model != "45" {
			return fmt.Errorf("%d: model: got %q; want \"45\"", i, cpu.Model)
		}
		if cpu.ModelName != XeonE52690ModelName {
			return fmt.Errorf("%d: model name: got %q; want %q", i, cpu.ModelName, XeonE52690ModelName)
		}
		if cpu.Stepping != "7" {
			return fmt.Errorf("%d: stepping: got %q; want \"7\"", i, cpu.Stepping)
		}
		if cpu.Microcode != "0x710" {
			return fmt.Errorf("%d: microcode: got %q; want \"0x710\"", i, cpu.Microcode)
		}
		if int(cpu.CPUMHz) < 1200 {
			return fmt.Errorf("%d: cpu MHz: got %.3f; want a value >= 1200", i, cpu.CPUMHz)
		}
		if cpu.CacheSize != "20480 KB" {
			return fmt.Errorf("%d: cache size: got %q; want \"20480 KB\"", i, cpu.CacheSize)
		}
		if cpu.PhysicalID != 0 || cpu.PhysicalID != 1 {
			return fmt.Errorf("%d: physical id: got %d; want 0 or 1", i, cpu.PhysicalID)
		}
		if cpu.Siblings != 16 {
			return fmt.Errorf("%d: siblings: got %d; want 16", i, cpu.Siblings)
		}
		if cpu.CPUCores != 8 {
			return fmt.Errorf("%d: cpu cores: got %d; want 8", i, cpu.CPUCores)
		}
		y := "yes"
		if cpu.FPU != y {
			return fmt.Errorf("%d: fpu: got %q; want %q", i, cpu.FPU, y)
		}
		if cpu.FPUException != y {
			return fmt.Errorf("%d: fpu exception: got %q; want %q", i, cpu.FPUException, y)
		}
		if cpu.CPUIDLevel != "13" {
			return fmt.Errorf("%d: cpuid level: got %q; want \"130\"", i, cpu.CPUIDLevel)
		}
		if cpu.WP != y {
			return fmt.Errorf("%d: wp: got %q; want %q", i, cpu.WP, y)
		}
		if len(cpu.Flags) != 99 {
			return fmt.Errorf("%d: flags: got %d; want 99", i, len(cpu.Flags))
		}
		if len(cpu.Bugs) != 0 {
			return fmt.Errorf("%d: bugs: got %d; want 0", i, len(cpu.Bugs))
		}
		if int(cpu.BogoMIPS) < 5786 {
			return fmt.Errorf("%d: bogomips: got %.3f; want a value >= 5786", i, cpu.BogoMIPS)
		}
		if int(cpu.CLFlushSize) != 64 {
			return fmt.Errorf("%d: clflush size: got %q; want 64", i, cpu.CLFlushSize)
		}
		if int(cpu.CacheAlignment) != 64 {
			return fmt.Errorf("%d: cache alignment size: got %q; want 64", i, cpu.CacheAlignment)
		}
		if len(cpu.PowerManagement) != 0 {
			return fmt.Errorf("%d: power management: got %d; wanted 0", i, len(cpu.PowerManagement))
		}
		if len(cpu.AddressSizes) != 2 {
			return fmt.Errorf("%d: address sizes: got %d; want 2", i, len(cpu.AddressSizes))
		}
	}
	return nil
}

func ValidateXeonE52690CPUFreq(f *cpufreq.Frequency) error {
	if f.Timestamp == 0 {
		return errors.New("expected Timestamp to have a non-zero value; it didn't")
	}
	if len(f.CPU) != 32 {
		return fmt.Errorf("CPU: got %d; want 32", len(f.CPU))
	}
	if f.Sockets != 2 {
		return fmt.Errorf("got %d socket; want 2", len(f.CPU))
	}
	for i, cpu := range f.CPU {
		if int(cpu.CPUMHz) < 1200 {
			return fmt.Errorf("%d: cpu MHz: got %.3f; want a value >= 1200", i, cpu.CPUMHz)
		}
		if cpu.PhysicalID != 0 && cpu.PhysicalID != 1 {
			return fmt.Errorf("%d: physical id: got %d; want 0 or 1", i, cpu.PhysicalID)
		}
	}
	return nil
}

// ValidateXeonE52690Proc verifies that the info in the struct is consistent
// with the above data and the generated test sys frequency info. If everything
// verifies a nil is returned, otherwise an error is returned. This is used for
// testing.
func ValidateXeonE52690Proc(proc *processors.Processors, freq bool) error {
	// arch relies on a syscall so just check that it's not empty
	if proc.Architecture == "" {
		return errors.New("architecture: got an empty string; want a non-empty string")
	}
	if proc.Sockets != 2 {
		return fmt.Errorf("sockets: got %d; want 2", proc.Sockets)
	}

	if int(proc.CoresPerSocket) != 8 {
		return fmt.Errorf("cores per socket: got %d; want 8", proc.CoresPerSocket)
	}

	if proc.CPUs != 32 {
		return fmt.Errorf("CPUs: got %d; want 32", proc.CPUs)
	}
	if proc.Possible != "0-31" {
		return fmt.Errorf("possible: got %q; want \"0-31\"", proc.Possible)
	}

	if proc.VendorID != GenuineIntel {
		return fmt.Errorf("vendor_id: got %q; want %q", proc.VendorID, GenuineIntel)
	}
	if proc.CPUFamily != "6" {
		return fmt.Errorf("cpu family: got %q; want \"6\"", proc.CPUFamily)
	}
	if proc.Model != "45" {
		return fmt.Errorf("model: got %q; want \"45\"", proc.Model)
	}
	if proc.ModelName != XeonE52690ModelName {
		return fmt.Errorf("model name: got %q; want %q", proc.ModelName, XeonE52690ModelName)
	}
	if proc.Stepping != "7" {
		return fmt.Errorf("stepping: got %q; want \"7\"", proc.Stepping)
	}
	if proc.Microcode != "0x710" {
		return fmt.Errorf("microcode: got %q; want \"0x710\"", proc.Microcode)
	}
	if int(proc.CPUMHz) < 1200 {
		return fmt.Errorf("cpu MHz: got %.3f; want a value >= 1200", proc.CPUMHz)
	}
	if freq {
		if int(proc.MHzMin) != 1600 {
			return fmt.Errorf("MHzMin: got %.3f; wanted 1600.000", proc.MHzMin)
		}
		if int(proc.MHzMax) != 2800 {
			return fmt.Errorf("MHzMax: got %.3f; wanted 2800.000", proc.MHzMax)
		}
	} else {
		if int(proc.MHzMin) != 0 {
			return fmt.Errorf("MHzMin: got %.3f; wanted 0", proc.MHzMin)
		}
		if int(proc.MHzMax) != 0 {
			return fmt.Errorf("MHzMax: got %.3f; wanted 0", proc.MHzMax)
		}
	}
	if proc.ThreadsPerCore != 2 {
		fmt.Errorf("threads per core: got %d; want 2", proc.ThreadsPerCore)
	}
	if len(proc.Flags) != 79 {
		return fmt.Errorf("flags: got %d; want 99", len(proc.Flags))
	}
	if len(proc.Bugs) > 0 {
		return fmt.Errorf("bugs: got %d elements; wanted 0", len(proc.Bugs))
	}
	if len(proc.OpModes) != 2 {
		return fmt.Errorf("op modes: got %d; want 2", len(proc.OpModes))
	}
	if proc.OpModes[0] != "32-bit" {
		return fmt.Errorf("OpModes: got %q; want \"32-bit\"", proc.OpModes[0])
	}
	if proc.OpModes[1] != "64-bit" {
		return fmt.Errorf("OpModes: got %q; want \"64-bit\"", proc.OpModes[1])
	}
	if int(proc.BogoMIPS) < 5786 {
		return fmt.Errorf("bogomips: got %.3f; want a value >= 5786", proc.BogoMIPS)
	}
	if len(proc.CacheIDs) != 4 {
		return fmt.Errorf("CacheIDs: got %d elements; wanted 4", len(proc.CacheIDs))
	}
	if proc.CacheSize != "20480 KB" {
		return fmt.Errorf("cache size: got %q; want \"20480 KB\"", proc.CacheSize)
	}
	if len(proc.Cache) != 4 {
		return fmt.Errorf("Cache: got %d entries; wanted 4", len(proc.Cache))
	}
	for _, v := range proc.CacheIDs {
		val, ok := proc.Cache[v]
		if !ok {
			return fmt.Errorf("Cache: expected information about %q; none found", v)
		}
		if val == "" {
			return fmt.Errorf("Cache: expected information about %q; got an empty string", v)
		}
	}
	if proc.Virtualization != processors.VTx {
		return fmt.Errorf("virtualization: got %q; want %q", proc.Virtualization, processors.VTx)
	}
	if proc.NumaNodes != 2 {
		return fmt.Errorf("numa nodes: got %d; want 2", proc.NumaNodes)
	}

	if len(proc.NumaNodeCPUs) != 2 {
		return fmt.Errorf("numa node cpus: got %d elements; wanted 2", len(proc.NumaNodeCPUs))
	}
	for i, v := range proc.NumaNodeCPUs {
		if v.ID != int32(i) {
			return fmt.Errorf("numa node id: got %d; want %d", v.ID, i)
		}
		switch i {
		case 0:
			if v.CPUList != "0-15" {
				return fmt.Errorf("numa node cpulist: got %s; want \"0-15\"", v.CPUList)
			}
		case 1:
			if v.CPUList != "16-31" {
				return fmt.Errorf("numa node cpulist: got %s; want \"16-31\"", v.CPUList)
			}
		default:
			return fmt.Errorf("numa node cpus: unexpected index: %d", i)
		}

	}

	return nil
}

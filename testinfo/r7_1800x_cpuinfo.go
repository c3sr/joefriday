package testinfo

import (
	"errors"
	"fmt"

	"github.com/rai-project/joefriday/cpu/cpufreq"
	"github.com/rai-project/joefriday/cpu/cpuinfo"
	"github.com/rai-project/joefriday/processors"
)

const (
	AuthenticAMD   = "AuthenticAMD"
	R71800xTLBSize = "2560 4K pages"
)

var R71800xCPUInfo = []byte(`processor	: 0
vendor_id	: AuthenticAMD
cpu family	: 23
model		: 1
model name	: AMD Ryzen 7 1800X Eight-Core Processor
stepping	: 1
microcode	: 0x800111c
cpu MHz		: 2200.000
cache size	: 512 KB
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
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm constant_tsc rep_good nopl nonstop_tsc extd_apicid aperfmperf pni pclmulqdq monitor ssse3 fma cx16 sse4_1 sse4_2 movbe popcnt aes xsave avx f16c rdrand lahf_lm cmp_legacy svm extapic cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw skinit wdt tce topoext perfctr_core perfctr_nb bpext perfctr_l2 mwaitx hw_pstate vmmcall fsgsbase bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero irperf arat npt lbrv svm_lock nrip_save tsc_scale vmcb_clean flushbyasid decodeassists pausefilter pfthreshold avic overflow_recov succor smca
bugs		: fxsave_leak sysret_ss_attrs null_seg
bogomips	: 7186.42
TLB size	: 2560 4K pages
clflush size	: 64
cache_alignment	: 64
address sizes	: 48 bits physical, 48 bits virtual
power management: ts ttp tm hwpstate eff_freq_ro [13] [14]

processor	: 1
vendor_id	: AuthenticAMD
cpu family	: 23
model		: 1
model name	: AMD Ryzen 7 1800X Eight-Core Processor
stepping	: 1
microcode	: 0x800111c
cpu MHz		: 2200.000
cache size	: 512 KB
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
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm constant_tsc rep_good nopl nonstop_tsc extd_apicid aperfmperf pni pclmulqdq monitor ssse3 fma cx16 sse4_1 sse4_2 movbe popcnt aes xsave avx f16c rdrand lahf_lm cmp_legacy svm extapic cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw skinit wdt tce topoext perfctr_core perfctr_nb bpext perfctr_l2 mwaitx hw_pstate vmmcall fsgsbase bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero irperf arat npt lbrv svm_lock nrip_save tsc_scale vmcb_clean flushbyasid decodeassists pausefilter pfthreshold avic overflow_recov succor smca
bugs		: fxsave_leak sysret_ss_attrs null_seg
bogomips	: 7180.28
TLB size	: 2560 4K pages
clflush size	: 64
cache_alignment	: 64
address sizes	: 48 bits physical, 48 bits virtual
power management: ts ttp tm hwpstate eff_freq_ro [13] [14]

processor	: 2
vendor_id	: AuthenticAMD
cpu family	: 23
model		: 1
model name	: AMD Ryzen 7 1800X Eight-Core Processor
stepping	: 1
microcode	: 0x800111c
cpu MHz		: 2200.000
cache size	: 512 KB
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
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm constant_tsc rep_good nopl nonstop_tsc extd_apicid aperfmperf pni pclmulqdq monitor ssse3 fma cx16 sse4_1 sse4_2 movbe popcnt aes xsave avx f16c rdrand lahf_lm cmp_legacy svm extapic cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw skinit wdt tce topoext perfctr_core perfctr_nb bpext perfctr_l2 mwaitx hw_pstate vmmcall fsgsbase bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero irperf arat npt lbrv svm_lock nrip_save tsc_scale vmcb_clean flushbyasid decodeassists pausefilter pfthreshold avic overflow_recov succor smca
bugs		: fxsave_leak sysret_ss_attrs null_seg
bogomips	: 7180.28
TLB size	: 2560 4K pages
clflush size	: 64
cache_alignment	: 64
address sizes	: 48 bits physical, 48 bits virtual
power management: ts ttp tm hwpstate eff_freq_ro [13] [14]

processor	: 3
vendor_id	: AuthenticAMD
cpu family	: 23
model		: 1
model name	: AMD Ryzen 7 1800X Eight-Core Processor
stepping	: 1
microcode	: 0x800111c
cpu MHz		: 2200.000
cache size	: 512 KB
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
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm constant_tsc rep_good nopl nonstop_tsc extd_apicid aperfmperf pni pclmulqdq monitor ssse3 fma cx16 sse4_1 sse4_2 movbe popcnt aes xsave avx f16c rdrand lahf_lm cmp_legacy svm extapic cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw skinit wdt tce topoext perfctr_core perfctr_nb bpext perfctr_l2 mwaitx hw_pstate vmmcall fsgsbase bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero irperf arat npt lbrv svm_lock nrip_save tsc_scale vmcb_clean flushbyasid decodeassists pausefilter pfthreshold avic overflow_recov succor smca
bugs		: fxsave_leak sysret_ss_attrs null_seg
bogomips	: 7180.28
TLB size	: 2560 4K pages
clflush size	: 64
cache_alignment	: 64
address sizes	: 48 bits physical, 48 bits virtual
power management: ts ttp tm hwpstate eff_freq_ro [13] [14]

processor	: 4
vendor_id	: AuthenticAMD
cpu family	: 23
model		: 1
model name	: AMD Ryzen 7 1800X Eight-Core Processor
stepping	: 1
microcode	: 0x800111c
cpu MHz		: 2200.000
cache size	: 512 KB
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
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm constant_tsc rep_good nopl nonstop_tsc extd_apicid aperfmperf pni pclmulqdq monitor ssse3 fma cx16 sse4_1 sse4_2 movbe popcnt aes xsave avx f16c rdrand lahf_lm cmp_legacy svm extapic cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw skinit wdt tce topoext perfctr_core perfctr_nb bpext perfctr_l2 mwaitx hw_pstate vmmcall fsgsbase bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero irperf arat npt lbrv svm_lock nrip_save tsc_scale vmcb_clean flushbyasid decodeassists pausefilter pfthreshold avic overflow_recov succor smca
bugs		: fxsave_leak sysret_ss_attrs null_seg
bogomips	: 7180.28
TLB size	: 2560 4K pages
clflush size	: 64
cache_alignment	: 64
address sizes	: 48 bits physical, 48 bits virtual
power management: ts ttp tm hwpstate eff_freq_ro [13] [14]

processor	: 5
vendor_id	: AuthenticAMD
cpu family	: 23
model		: 1
model name	: AMD Ryzen 7 1800X Eight-Core Processor
stepping	: 1
microcode	: 0x800111c
cpu MHz		: 2200.000
cache size	: 512 KB
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
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm constant_tsc rep_good nopl nonstop_tsc extd_apicid aperfmperf pni pclmulqdq monitor ssse3 fma cx16 sse4_1 sse4_2 movbe popcnt aes xsave avx f16c rdrand lahf_lm cmp_legacy svm extapic cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw skinit wdt tce topoext perfctr_core perfctr_nb bpext perfctr_l2 mwaitx hw_pstate vmmcall fsgsbase bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero irperf arat npt lbrv svm_lock nrip_save tsc_scale vmcb_clean flushbyasid decodeassists pausefilter pfthreshold avic overflow_recov succor smca
bugs		: fxsave_leak sysret_ss_attrs null_seg
bogomips	: 7180.28
TLB size	: 2560 4K pages
clflush size	: 64
cache_alignment	: 64
address sizes	: 48 bits physical, 48 bits virtual
power management: ts ttp tm hwpstate eff_freq_ro [13] [14]

processor	: 6
vendor_id	: AuthenticAMD
cpu family	: 23
model		: 1
model name	: AMD Ryzen 7 1800X Eight-Core Processor
stepping	: 1
microcode	: 0x800111c
cpu MHz		: 2200.000
cache size	: 512 KB
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
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm constant_tsc rep_good nopl nonstop_tsc extd_apicid aperfmperf pni pclmulqdq monitor ssse3 fma cx16 sse4_1 sse4_2 movbe popcnt aes xsave avx f16c rdrand lahf_lm cmp_legacy svm extapic cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw skinit wdt tce topoext perfctr_core perfctr_nb bpext perfctr_l2 mwaitx hw_pstate vmmcall fsgsbase bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero irperf arat npt lbrv svm_lock nrip_save tsc_scale vmcb_clean flushbyasid decodeassists pausefilter pfthreshold avic overflow_recov succor smca
bugs		: fxsave_leak sysret_ss_attrs null_seg
bogomips	: 7180.28
TLB size	: 2560 4K pages
clflush size	: 64
cache_alignment	: 64
address sizes	: 48 bits physical, 48 bits virtual
power management: ts ttp tm hwpstate eff_freq_ro [13] [14]

processor	: 7
vendor_id	: AuthenticAMD
cpu family	: 23
model		: 1
model name	: AMD Ryzen 7 1800X Eight-Core Processor
stepping	: 1
microcode	: 0x800111c
cpu MHz		: 2200.000
cache size	: 512 KB
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
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm constant_tsc rep_good nopl nonstop_tsc extd_apicid aperfmperf pni pclmulqdq monitor ssse3 fma cx16 sse4_1 sse4_2 movbe popcnt aes xsave avx f16c rdrand lahf_lm cmp_legacy svm extapic cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw skinit wdt tce topoext perfctr_core perfctr_nb bpext perfctr_l2 mwaitx hw_pstate vmmcall fsgsbase bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero irperf arat npt lbrv svm_lock nrip_save tsc_scale vmcb_clean flushbyasid decodeassists pausefilter pfthreshold avic overflow_recov succor smca
bugs		: fxsave_leak sysret_ss_attrs null_seg
bogomips	: 7180.28
TLB size	: 2560 4K pages
clflush size	: 64
cache_alignment	: 64
address sizes	: 48 bits physical, 48 bits virtual
power management: ts ttp tm hwpstate eff_freq_ro [13] [14]

processor	: 8
vendor_id	: AuthenticAMD
cpu family	: 23
model		: 1
model name	: AMD Ryzen 7 1800X Eight-Core Processor
stepping	: 1
microcode	: 0x800111c
cpu MHz		: 2200.000
cache size	: 512 KB
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
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm constant_tsc rep_good nopl nonstop_tsc extd_apicid aperfmperf pni pclmulqdq monitor ssse3 fma cx16 sse4_1 sse4_2 movbe popcnt aes xsave avx f16c rdrand lahf_lm cmp_legacy svm extapic cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw skinit wdt tce topoext perfctr_core perfctr_nb bpext perfctr_l2 mwaitx hw_pstate vmmcall fsgsbase bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero irperf arat npt lbrv svm_lock nrip_save tsc_scale vmcb_clean flushbyasid decodeassists pausefilter pfthreshold avic overflow_recov succor smca
bugs		: fxsave_leak sysret_ss_attrs null_seg
bogomips	: 7180.28
TLB size	: 2560 4K pages
clflush size	: 64
cache_alignment	: 64
address sizes	: 48 bits physical, 48 bits virtual
power management: ts ttp tm hwpstate eff_freq_ro [13] [14]

processor	: 9
vendor_id	: AuthenticAMD
cpu family	: 23
model		: 1
model name	: AMD Ryzen 7 1800X Eight-Core Processor
stepping	: 1
microcode	: 0x800111c
cpu MHz		: 2200.000
cache size	: 512 KB
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
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm constant_tsc rep_good nopl nonstop_tsc extd_apicid aperfmperf pni pclmulqdq monitor ssse3 fma cx16 sse4_1 sse4_2 movbe popcnt aes xsave avx f16c rdrand lahf_lm cmp_legacy svm extapic cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw skinit wdt tce topoext perfctr_core perfctr_nb bpext perfctr_l2 mwaitx hw_pstate vmmcall fsgsbase bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero irperf arat npt lbrv svm_lock nrip_save tsc_scale vmcb_clean flushbyasid decodeassists pausefilter pfthreshold avic overflow_recov succor smca
bugs		: fxsave_leak sysret_ss_attrs null_seg
bogomips	: 7180.28
TLB size	: 2560 4K pages
clflush size	: 64
cache_alignment	: 64
address sizes	: 48 bits physical, 48 bits virtual
power management: ts ttp tm hwpstate eff_freq_ro [13] [14]

processor	: 10
vendor_id	: AuthenticAMD
cpu family	: 23
model		: 1
model name	: AMD Ryzen 7 1800X Eight-Core Processor
stepping	: 1
microcode	: 0x800111c
cpu MHz		: 2200.000
cache size	: 512 KB
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
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm constant_tsc rep_good nopl nonstop_tsc extd_apicid aperfmperf pni pclmulqdq monitor ssse3 fma cx16 sse4_1 sse4_2 movbe popcnt aes xsave avx f16c rdrand lahf_lm cmp_legacy svm extapic cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw skinit wdt tce topoext perfctr_core perfctr_nb bpext perfctr_l2 mwaitx hw_pstate vmmcall fsgsbase bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero irperf arat npt lbrv svm_lock nrip_save tsc_scale vmcb_clean flushbyasid decodeassists pausefilter pfthreshold avic overflow_recov succor smca
bugs		: fxsave_leak sysret_ss_attrs null_seg
bogomips	: 7180.28
TLB size	: 2560 4K pages
clflush size	: 64
cache_alignment	: 64
address sizes	: 48 bits physical, 48 bits virtual
power management: ts ttp tm hwpstate eff_freq_ro [13] [14]

processor	: 11
vendor_id	: AuthenticAMD
cpu family	: 23
model		: 1
model name	: AMD Ryzen 7 1800X Eight-Core Processor
stepping	: 1
microcode	: 0x800111c
cpu MHz		: 2200.000
cache size	: 512 KB
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
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm constant_tsc rep_good nopl nonstop_tsc extd_apicid aperfmperf pni pclmulqdq monitor ssse3 fma cx16 sse4_1 sse4_2 movbe popcnt aes xsave avx f16c rdrand lahf_lm cmp_legacy svm extapic cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw skinit wdt tce topoext perfctr_core perfctr_nb bpext perfctr_l2 mwaitx hw_pstate vmmcall fsgsbase bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero irperf arat npt lbrv svm_lock nrip_save tsc_scale vmcb_clean flushbyasid decodeassists pausefilter pfthreshold avic overflow_recov succor smca
bugs		: fxsave_leak sysret_ss_attrs null_seg
bogomips	: 7180.28
TLB size	: 2560 4K pages
clflush size	: 64
cache_alignment	: 64
address sizes	: 48 bits physical, 48 bits virtual
power management: ts ttp tm hwpstate eff_freq_ro [13] [14]

processor	: 12
vendor_id	: AuthenticAMD
cpu family	: 23
model		: 1
model name	: AMD Ryzen 7 1800X Eight-Core Processor
stepping	: 1
microcode	: 0x800111c
cpu MHz		: 2200.000
cache size	: 512 KB
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
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm constant_tsc rep_good nopl nonstop_tsc extd_apicid aperfmperf pni pclmulqdq monitor ssse3 fma cx16 sse4_1 sse4_2 movbe popcnt aes xsave avx f16c rdrand lahf_lm cmp_legacy svm extapic cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw skinit wdt tce topoext perfctr_core perfctr_nb bpext perfctr_l2 mwaitx hw_pstate vmmcall fsgsbase bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero irperf arat npt lbrv svm_lock nrip_save tsc_scale vmcb_clean flushbyasid decodeassists pausefilter pfthreshold avic overflow_recov succor smca
bugs		: fxsave_leak sysret_ss_attrs null_seg
bogomips	: 7180.28
TLB size	: 2560 4K pages
clflush size	: 64
cache_alignment	: 64
address sizes	: 48 bits physical, 48 bits virtual
power management: ts ttp tm hwpstate eff_freq_ro [13] [14]

processor	: 13
vendor_id	: AuthenticAMD
cpu family	: 23
model		: 1
model name	: AMD Ryzen 7 1800X Eight-Core Processor
stepping	: 1
microcode	: 0x800111c
cpu MHz		: 2200.000
cache size	: 512 KB
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
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm constant_tsc rep_good nopl nonstop_tsc extd_apicid aperfmperf pni pclmulqdq monitor ssse3 fma cx16 sse4_1 sse4_2 movbe popcnt aes xsave avx f16c rdrand lahf_lm cmp_legacy svm extapic cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw skinit wdt tce topoext perfctr_core perfctr_nb bpext perfctr_l2 mwaitx hw_pstate vmmcall fsgsbase bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero irperf arat npt lbrv svm_lock nrip_save tsc_scale vmcb_clean flushbyasid decodeassists pausefilter pfthreshold avic overflow_recov succor smca
bugs		: fxsave_leak sysret_ss_attrs null_seg
bogomips	: 7180.28
TLB size	: 2560 4K pages
clflush size	: 64
cache_alignment	: 64
address sizes	: 48 bits physical, 48 bits virtual
power management: ts ttp tm hwpstate eff_freq_ro [13] [14]

processor	: 14
vendor_id	: AuthenticAMD
cpu family	: 23
model		: 1
model name	: AMD Ryzen 7 1800X Eight-Core Processor
stepping	: 1
microcode	: 0x800111c
cpu MHz		: 2200.000
cache size	: 512 KB
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
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm constant_tsc rep_good nopl nonstop_tsc extd_apicid aperfmperf pni pclmulqdq monitor ssse3 fma cx16 sse4_1 sse4_2 movbe popcnt aes xsave avx f16c rdrand lahf_lm cmp_legacy svm extapic cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw skinit wdt tce topoext perfctr_core perfctr_nb bpext perfctr_l2 mwaitx hw_pstate vmmcall fsgsbase bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero irperf arat npt lbrv svm_lock nrip_save tsc_scale vmcb_clean flushbyasid decodeassists pausefilter pfthreshold avic overflow_recov succor smca
bugs		: fxsave_leak sysret_ss_attrs null_seg
bogomips	: 7180.28
TLB size	: 2560 4K pages
clflush size	: 64
cache_alignment	: 64
address sizes	: 48 bits physical, 48 bits virtual
power management: ts ttp tm hwpstate eff_freq_ro [13] [14]

processor	: 15
vendor_id	: AuthenticAMD
cpu family	: 23
model		: 1
model name	: AMD Ryzen 7 1800X Eight-Core Processor
stepping	: 1
microcode	: 0x800111c
cpu MHz		: 2200.000
cache size	: 512 KB
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
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm constant_tsc rep_good nopl nonstop_tsc extd_apicid aperfmperf pni pclmulqdq monitor ssse3 fma cx16 sse4_1 sse4_2 movbe popcnt aes xsave avx f16c rdrand lahf_lm cmp_legacy svm extapic cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw skinit wdt tce topoext perfctr_core perfctr_nb bpext perfctr_l2 mwaitx hw_pstate vmmcall fsgsbase bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero irperf arat npt lbrv svm_lock nrip_save tsc_scale vmcb_clean flushbyasid decodeassists pausefilter pfthreshold avic overflow_recov succor smca
bugs		: fxsave_leak sysret_ss_attrs null_seg
bogomips	: 7151.61
TLB size	: 2560 4K pages
clflush size	: 64
cache_alignment	: 64
address sizes	: 48 bits physical, 48 bits virtual
power management: ts ttp tm hwpstate eff_freq_ro [13] [14]
`)

func ValidateR71800xCPUInfo(inf *cpuinfo.CPUInfo) error {
	if inf.Timestamp == 0 {
		return errors.New("expected timestamp to have a nonzero value, it didn't")
	}
	if len(inf.CPU) != 16 {
		return fmt.Errorf("Expected 16 CPU entries, got %d", len(inf.CPU))
	}
	if inf.Sockets != 1 {
		return fmt.Errorf("got %d socket; want 1", len(inf.CPU))
	}
	for i, cpu := range inf.CPU {
		if cpu.VendorID != AuthenticAMD {
			return fmt.Errorf("%d: VendorID: got %q; want %q", i, cpu.VendorID, AuthenticAMD)
		}
		if cpu.CPUFamily != "23" {
			return fmt.Errorf("%d: CPUFamily: got %q; want \"23\"", i, cpu.CPUFamily)
		}
		if cpu.Model != "1" {
			return fmt.Errorf("%d: Model: got %q; want \"1\"", i, cpu.Model)
		}
		if cpu.ModelName != "AMD Ryzen 7 1800X Eight-Core Processor" {
			return fmt.Errorf("%d: ModelName: got %q; want \"AMD Ryzen 7 1800X Eight-Core Processor\"", i, cpu.ModelName)
		}
		if cpu.Stepping != "1" {
			return fmt.Errorf("%d: Stepping: got %q; want \"1\"", i, cpu.Stepping)
		}
		if cpu.Microcode != "0x800111c" {
			return fmt.Errorf("%d: MicroCode: got %q; want \"0x800111c\"", i, cpu.Microcode)
		}
		if int(cpu.CPUMHz) != 2200 {
			return fmt.Errorf("%d: CPUMHz: got %d; want 2200", i, int(cpu.CPUMHz))
		}
		if cpu.CacheSize != "512 KB" {
			return fmt.Errorf("%d: CacheSize: got %q; want \"512 KB\"", i, cpu.CacheSize)
		}
		if cpu.Siblings != 16 {
			return fmt.Errorf("%d: Siblings: got %d; want 16", i, cpu.Siblings)
		}
		if cpu.CPUCores != 8 {
			return fmt.Errorf("%d: CPUCores: got %d; want 8", i, cpu.CPUCores)
		}
		if cpu.FPU != "yes" {
			return fmt.Errorf("%d: FPU: got %q; want \"yes\"", i, cpu.FPU)
		}
		if cpu.FPUException != "yes" {
			return fmt.Errorf("%d: FPUException: got %q; want \"yes\"", i, cpu.FPUException)
		}
		if cpu.CPUIDLevel != "13" {
			return fmt.Errorf("%d: CPUIDLevel: got %q; want \"13\"", i, cpu.CPUIDLevel)
		}
		if cpu.WP != "yes" {
			return fmt.Errorf("%d: WP: got %q; want \"yes\"", i, cpu.WP)
		}
		if len(cpu.Flags) != 103 {
			return fmt.Errorf("%d: Flags: got %q; want 103", i, len(cpu.Flags))
		}
		if len(cpu.Bugs) != 3 {
			return fmt.Errorf("%d: Bugs: got %d; want 3", i, len(cpu.Bugs))
		}
		if int(cpu.CLFlushSize) != 64 {
			return fmt.Errorf("%d: Siblings: got %d; want 64", i, cpu.CLFlushSize)
		}
		if int(cpu.CacheAlignment) != 64 {
			return fmt.Errorf("%d: Siblings: got %d; want 64", i, cpu.CacheAlignment)
		}
		if cpu.TLBSize != R71800xTLBSize {
			return fmt.Errorf("%d: tlb size: got %q; want %q", i, cpu.TLBSize, R71800xTLBSize)
		}
		if len(cpu.PowerManagement) != 7 {
			return fmt.Errorf("%d: power management: got %d; want 7", i, len(cpu.PowerManagement))
		}
		if len(cpu.AddressSizes) != 2 {
			return fmt.Errorf("%d: address sizes: got %d; want 2", i, len(cpu.AddressSizes))
		}
		if int(cpu.BogoMIPS) < 7100 {
			return fmt.Errorf("bogomips: got %.3f; want a value >= 7100", cpu.BogoMIPS)
		}
	}
	return nil
}

// ValidateR71800xCPUFreq verifies that the info in the struct info is
// consistent with relevant parts of the above data. If everything verifies a
// nil is returned, otherwise an error is returned. This is used for testing.
func ValidateR71800xCPUFreq(f *cpufreq.Frequency) error {
	if f.Timestamp == 0 {
		return errors.New("expected timestamp to have a nonzero value, it didn't")
	}
	if len(f.CPU) != 16 {
		return fmt.Errorf("Expected 16 CPU entries, got %d", len(f.CPU))
	}

	for i, cpu := range f.CPU {
		if int(cpu.CPUMHz) != 2200 {
			return fmt.Errorf("%d: CPUMHz: got %d; want 2200", i, int(cpu.CPUMHz))
		}
		if int(cpu.PhysicalID) != 0 {
			return fmt.Errorf("%d: physical id: got %d; want 0", i, cpu.PhysicalID)
		}
		if int(cpu.Processor) != i {
			return fmt.Errorf("%d: processor: got %d; want %d", i, cpu.Processor, i)
		}
		// APICID happens to be consistent with i
		if int(cpu.APICID) != i {
			return fmt.Errorf("%d: apicid: got %d; want %d", i, cpu.APICID, i)
		}
		if int(cpu.CoreID) != i/2 {
			return fmt.Errorf("%d: core id: got %d; want %d", i, cpu.CoreID, i/2)
		}
	}
	return nil
}

// ValidateR71800xProc verifies that the info in the struct is consistent with
// the above data and the generated test sys frequency info. If everything
// verifies a nil is returned, otherwise an error is returned. This is used for
// testing.
func ValidateR71800xProc(proc *processors.Processors, freq bool) error {
	// arch relies on a syscall so just check that it's not empty
	if proc.Architecture == "" {
		return errors.New("architecture: got an empty string; want a non-empty string")
	}
	if proc.Sockets != 1 {
		return fmt.Errorf("sockets: got %d; want 1", proc.Sockets)
	}

	if int(proc.CoresPerSocket) != 8 {
		return fmt.Errorf("cores per socket: got %d; want 8", proc.CoresPerSocket)
	}
	if proc.CPUs != 16 {
		return fmt.Errorf("CPUs: got %d; want 16", proc.CPUs)
	}
	if proc.Possible != "0-15" {
		return fmt.Errorf("possible: got %q; want \"0-15\"", proc.Possible)
	}

	if proc.VendorID != AuthenticAMD {
		return fmt.Errorf("vendor_id: got %q; want %q", proc.VendorID, AuthenticAMD)
	}
	if proc.CPUFamily != "23" {
		return fmt.Errorf("CPUFamily: got %q; want \"23\"", proc.CPUFamily)
	}
	if proc.Model != "1" {
		return fmt.Errorf("Model: got %q; want \"1\"", proc.Model)
	}
	if proc.ModelName != "AMD Ryzen 7 1800X Eight-Core Processor" {
		return fmt.Errorf("ModelName: got %q; want \"AMD Ryzen 7 1800X Eight-Core Processor\"", proc.ModelName)
	}
	if proc.Stepping != "1" {
		return fmt.Errorf("Stepping: got %q; want \"1\"", proc.Stepping)
	}
	if proc.Microcode != "0x800111c" {
		return fmt.Errorf("MicroCode: got %q; want \"0x800111c\"", proc.Microcode)
	}
	if int(proc.CPUMHz) != 2200 {
		return fmt.Errorf("CPUMHz: got %d; want 2200", int(proc.CPUMHz))
	}
	if len(proc.Flags) != 103 {
		return fmt.Errorf("Flags: got %q; want 103", len(proc.Flags))
	}
	if len(proc.Bugs) != 3 {
		return fmt.Errorf("Bugs: got %d; want 3", len(proc.Bugs))
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
	if proc.CacheSize != "512 KB" {
		return fmt.Errorf("CacheSize: got %q; want \"512 KB\"", proc.CacheSize)
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
	if proc.ThreadsPerCore != 2 {
		fmt.Errorf("threads per core: got %d; want 2", proc.ThreadsPerCore)
	}
	if len(proc.OpModes) != 2 {
		return fmt.Errorf("op modes: got %d; want 2", len(proc.OpModes))
	}
	if proc.OpModes[0] != "32-bit" {
		return fmt.Errorf("OpModes: got %q; want \"32-bit\"", proc.OpModes[2])
	}
	if proc.OpModes[1] != "64-bit" {
		return fmt.Errorf("OpModes: got %q; want \"64-bit\"", proc.OpModes[1])
	}
	if int(proc.BogoMIPS) < 7100 {
		return fmt.Errorf("bogomips: got %.3f; want a value >= 7100", proc.BogoMIPS)
	}
	if proc.Virtualization != processors.AMDV {
		return fmt.Errorf("virtualization: got %q; want %q", proc.Virtualization, processors.AMDV)
	}
	if proc.NumaNodes != 1 {
		return fmt.Errorf("numa nodes: got %d; want 1", proc.NumaNodes)
	}

	if len(proc.NumaNodeCPUs) != 1 {
		return fmt.Errorf("numa node cpus: got %d elements; wanted 1", len(proc.NumaNodeCPUs))
	}
	if proc.NumaNodeCPUs[0].ID != int32(0) {
		return fmt.Errorf("numa node id: got %d; want 0", proc.NumaNodeCPUs[0].ID)
	}
	if proc.NumaNodeCPUs[0].CPUList != "0-15" {
		return fmt.Errorf("numa node cpulist: got %s; want \"0-15\"", proc.NumaNodeCPUs[0].CPUList)
	}

	return nil
}

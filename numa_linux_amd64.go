package numa

import (
	"syscall"
	"unsafe"

	"github.com/intel-go/cpuid"
)

var fastway = cpuid.HasFeature(cpuid.RDTSCP)

//go:noescape
func getCPUAndNodeRDTSCP() (cpu int, node int)

// GetCPUAndNode returns the node id and cpu id which current caller running on.
// https://man7.org/linux/man-pages/man2/getcpu.2.html
//
// equal:
//
//	if fastway {
//		call RDTSCP
//	 The linux kernel will fill the node cpu id in the private data of each cpu.
//	 arch/x86/kernel/vsyscall_64.c@vsyscall_set_cpu
//	}
func GetCPUAndNode() (cpu int, node int) {
	if fastway {
		return getCPUAndNodeRDTSCP()
	}
	const sysGetCPU = 309
	_, _, errno := syscall.RawSyscall(sysGetCPU,
		uintptr(unsafe.Pointer(&cpu)),
		uintptr(unsafe.Pointer(&node)),
		0)
	if errno != 0 {
		cpu = 0
		node = 0
	}
	return
}

package main

import (
    "fmt"
    "syscall"
    "unsafe"
)


// for amd64
const MMAP_FLAGS = syscall.MAP_ANONYMOUS


type JIT struct {
    placeholder int
}


type MachineCode []uint8


/**
 * WARRING:
 *     The Go runtime stack always runs copystack and this will cause variables memory address changes.
 *     So that, the returned unsafe function can ONLY use sync.Pool variables (memory address fixed) as parameter. 
 */
func (jit JIT) Emit(m MachineCode) unsafe.Pointer {
    bin, err := syscall.Mmap(
        -1,
        0,
        len(m),
        syscall.PROT_READ | syscall.PROT_WRITE | syscall.PROT_EXEC, 
        syscall.MAP_PRIVATE | MMAP_FLAGS,
    )
    if err != nil {
        log.Fatalf("mmap failed: %v", err)
    }
    for i, b := range m {
        bin[i] = b
    }
    unsafeFunc := (uintptr)(unsafe.Pointer(&bin))
    return unsafe.Pointer(&unsafeFunc) 
}



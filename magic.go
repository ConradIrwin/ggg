package main

import (
	"syscall"
	"unsafe"
)

// http://ref.x86asm.net/coder64.html
// 0x48 => 64-bit instruction
// 0x89 => MOV
// 0x64 => %rsp to SIB + n'
// 0x24 => SIB for %rsp
// 0x08 => n = 8
var FUNC = []byte{
    0x48, 0x89, 0x64, 0x24, 0x08, // movq %rsp, 0x08(%rsp)
    0xc3, // ret
}

type magic uintptr

var Magic func() magic

func makePage() []byte {
    pagesize := uintptr(syscall.Getpagesize())
    page := make([]byte, pagesize)

    f := unsafe.Pointer(&page)
    data_addr := unsafe.Pointer(*(*uintptr)(f))
    delta := uintptr(data_addr) % pagesize

    if delta != 0 {
        panic("byte buffer not page aligned, needs moar codez")
    }

    copy(page, FUNC)
    err := syscall.Mprotect(page, syscall.PROT_EXEC|syscall.PROT_READ)

    if err != nil {
        panic(err)
    }

    return page[0:pagesize]
}

func init() {
    var page = makePage()

    p := unsafe.Pointer(&Magic)
    q := unsafe.Pointer(&page)

    *(*uintptr)(p) = (uintptr)(q)

    p = unsafe.Pointer(&Magic)
}

package main

/*
#cgo CFLAGS: -g -w -std=c99 -O3 -Dnetwork -D_BSD_SOURCE -I./header
#cgo LDFLAGS: -L. mrsh.a
#include <main.h>
#include <config.h>
#include <bloomfilter.h>
#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <fingerprint.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// FingerprintFile returns a string represenation of the mrsh-v2 fingerprint of the file
func FingerprintFile(path string) string {
	cpath := C.CString("./NOTICE.txt")
	defer C.free(unsafe.Pointer(cpath))
	cres := C.fingerprint_file(cpath)
	defer C.free(unsafe.Pointer(cres))
	res := C.GoString(cres)
	return res
}

func main() {
	fp := FingerprintFile("NOTICE.txt")
	fmt.Println("Result:", fp)
}

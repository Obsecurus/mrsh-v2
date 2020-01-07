package mrsh

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

extern int init_mrsh_done = false;

void init_mrsh_mode(){
	if (init_mrsh_done) {
		return;
	}
	mode = (MODES *)malloc(sizeof(MODES));
	mode->compare = false;
	mode->gen_compare = false;
	mode->compareLists = false;
	mode->file_comparison = true;
	mode->helpmessage = false;
	mode->print = false;
	mode->threshold = 1;
	mode->recursive = false;
	mode->path_list_compare = false;
	init_mrsh_done = true;
}
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

// Fingerprint is a wrapped C pointer
type Fingerprint struct {
	ptr *C.FINGERPRINT
}

// Init ensures that the Fingerprint has a Finalizer set on it, which allows the C struct to be freed.
// This should be called on every single Fingerprint as it is constructed.
func (fp *Fingerprint) Init() {
	runtime.SetFinalizer(fp, func(x interface{}) {
		fp := x.(Fingerprint)
		C.fingerprint_destroy(fp.ptr)
	})
}

func (fp *Fingerprint) String() string {
	cs := C.stringify_fingerprint(fp.ptr)
	return C.GoString(cs)
}

// FingerprintFile returns a pointer I think?
func FingerprintFile(path string) Fingerprint {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	cres := C.init_fingerprint_for_path(cpath)
	fp := Fingerprint{cres}
	fp.Init()
	return fp
}

// FingerprintCompare returns the similarity score between two fingerprints
func FingerprintCompare(a, b Fingerprint) int {
	C.init_mrsh_mode()
	fmt.Println("A:", a.String())
	fmt.Println("B:", b.String())
	res := C.fingerprint_compare(a.ptr, b.ptr)
	fmt.Println("Raw result:", res)
	return int(res)
}

// FingerprintRead converts a string representation of a fingerprint into a Fingerprint struct
func FingerprintRead(fpstring string) Fingerprint {
	ptr := C.init_empty_fingerprint()
	cs := C.CString(fpstring)
	C.read_fingerprint_string(ptr, cs)
	fp := Fingerprint{ptr}
	fp.Init()
	return fp
}
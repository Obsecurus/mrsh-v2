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

func main() {
	// one hash comes from reading the file
	original := FingerprintFile("original.js")

	// other hash comes from parsing a string
	mods := "mod.js:59212:1:139:6139212A40810046010A82C88D74855812" +
		"891C61213107935A04C5488A402582802083400607104C244004D10268B45" +
		"C220CB860684736175429404088880250C28412A13A2082C34251A4B00700" +
		"31AF046C13009909160C807825D4A4810CEA53020202840104007BE971920" +
		"E6010042BA88944062100820A6E0000182856283800E5C46118B486511C8A" +
		"360B08CC4C41B60000680A14016C147EC98D3E9C040C0084144A50A85A01C" +
		"F6580310B024E1141C113809480A9361681B100483141001400C203702002" +
		"5E17202165645E80409DA03034F011202304077A647158088E4C249901A31" +
		"59DC010473004224D200B2052B51072AC0000A9D1742B9C69E80283780129" +
		"405020820A0B23019ED81016C024558A0003042208218C02100240111231A" +
		"602A32D4E264061082204520D1A810483082393844C80028260836070A280" +
		"00201C341815A40820950010D0142104BF1260AA228D9044011920603B270" +
		"0C6010022B4222404681E904180A1004029592AACE0134B8E008B28924406" +
		"8003000044280444329807862394458000060A110119011F0010C60A56202" +
		"020927492EE64020330CB5345C5874174221522021C20518080A184402083" +
		"81204460A20C0883044C14C826148211408382B318530C2C37082A0321806" +
		"4005912A0094200310608040087100B500A0157C00F0E10A2A09C00004905" +
		"8010049B0100C1"
	mod := FingerprintRead(mods)
	fmt.Println("Comparison:", FingerprintCompare(original, mod))
}

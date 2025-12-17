package main

/*
#include <stdint.h>
*/
import "C"

import (
	"unsafe"

	"github.com/bluesentinelsec/FatStd/pkg/fatstrings"
)

var fatstdHandles = newHandleRegistry()

func fatstdStringNewFromGoString(value string) uintptr {
	return fatstdHandles.register(fatstrings.NewUTF8(value))
}

func fatstdStringFromHandle(handle uintptr) *fatstrings.String {
	if handle == 0 {
		panic("fatstdStringFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdStringFromHandle: invalid handle")
	}
	s, ok := value.(*fatstrings.String)
	if !ok {
		panic("fatstdStringFromHandle: handle is not a fat string")
	}
	return s
}

//export fatstd_go_string_new_utf8_cstr
func fatstd_go_string_new_utf8_cstr(cstr *C.char) C.uintptr_t {
	if cstr == nil {
		panic("fatstd_go_string_new_utf8_cstr: cstr is NULL")
	}
	handle := fatstdStringNewFromGoString(C.GoString(cstr))
	return C.uintptr_t(handle)
}

//export fatstd_go_string_new_utf8_n
func fatstd_go_string_new_utf8_n(bytes *C.char, len C.size_t) C.uintptr_t {
	if bytes == nil {
		if len == 0 {
			return C.uintptr_t(fatstdStringNewFromGoString(""))
		}
		panic("fatstd_go_string_new_utf8_n: bytes is NULL but len > 0")
	}
	if len > C.size_t(2147483647) {
		panic("fatstd_go_string_new_utf8_n: len too large")
	}

	buf := C.GoBytes(unsafe.Pointer(bytes), C.int(len))
	handle := fatstdStringNewFromGoString(string(buf))
	return C.uintptr_t(handle)
}

//export fatstd_go_string_clone
func fatstd_go_string_clone(handle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(handle))
	cloned := fatstrings.Clone(s.Value())
	return C.uintptr_t(fatstdStringNewFromGoString(cloned))
}

//export fatstd_go_string_contains
func fatstd_go_string_contains(aHandle C.uintptr_t, bHandle C.uintptr_t) C.int {
	a := fatstdStringFromHandle(uintptr(aHandle))
	b := fatstdStringFromHandle(uintptr(bHandle))
	if fatstrings.Contains(a.Value(), b.Value()) {
		return 1
	}
	return 0
}

//export fatstd_go_string_free
func fatstd_go_string_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_string_free: handle is 0")
	}

	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_string_free: invalid handle")
	}
	if _, ok := value.(*fatstrings.String); !ok {
		panic("fatstd_go_string_free: handle is not a fat string")
	}
}

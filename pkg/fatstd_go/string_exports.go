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

func fatstdStringArrayNew(values []string) uintptr {
	return fatstdHandles.register(fatstrings.NewStringArray(values))
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

func fatstdStringArrayFromHandle(handle uintptr) *fatstrings.StringArray {
	if handle == 0 {
		panic("fatstdStringArrayFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdStringArrayFromHandle: invalid handle")
	}
	a, ok := value.(*fatstrings.StringArray)
	if !ok {
		panic("fatstdStringArrayFromHandle: handle is not a fat string array")
	}
	return a
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

//export fatstd_go_string_has_prefix
func fatstd_go_string_has_prefix(sHandle C.uintptr_t, prefixHandle C.uintptr_t) C.int {
	s := fatstdStringFromHandle(uintptr(sHandle))
	prefix := fatstdStringFromHandle(uintptr(prefixHandle))
	if fatstrings.HasPrefix(s.Value(), prefix.Value()) {
		return 1
	}
	return 0
}

//export fatstd_go_string_has_suffix
func fatstd_go_string_has_suffix(sHandle C.uintptr_t, suffixHandle C.uintptr_t) C.int {
	s := fatstdStringFromHandle(uintptr(sHandle))
	suffix := fatstdStringFromHandle(uintptr(suffixHandle))
	if fatstrings.HasSuffix(s.Value(), suffix.Value()) {
		return 1
	}
	return 0
}

//export fatstd_go_string_trim_space
func fatstd_go_string_trim_space(handle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(handle))
	trimmed := fatstrings.TrimSpace(s.Value())
	return C.uintptr_t(fatstdStringNewFromGoString(trimmed))
}

//export fatstd_go_string_trim
func fatstd_go_string_trim(sHandle C.uintptr_t, cutsetHandle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	cutset := fatstdStringFromHandle(uintptr(cutsetHandle))
	trimmed := fatstrings.Trim(s.Value(), cutset.Value())
	return C.uintptr_t(fatstdStringNewFromGoString(trimmed))
}

//export fatstd_go_string_split
func fatstd_go_string_split(sHandle C.uintptr_t, sepHandle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	sep := fatstdStringFromHandle(uintptr(sepHandle))
	return C.uintptr_t(fatstdStringArrayNew(fatstrings.Split(s.Value(), sep.Value())))
}

//export fatstd_go_string_split_n
func fatstd_go_string_split_n(sHandle C.uintptr_t, sepHandle C.uintptr_t, n C.int) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	sep := fatstdStringFromHandle(uintptr(sepHandle))
	return C.uintptr_t(fatstdStringArrayNew(fatstrings.SplitN(s.Value(), sep.Value(), int(n))))
}

//export fatstd_go_string_array_len
func fatstd_go_string_array_len(arrayHandle C.uintptr_t) C.size_t {
	a := fatstdStringArrayFromHandle(uintptr(arrayHandle))
	return C.size_t(a.Len())
}

//export fatstd_go_string_array_get
func fatstd_go_string_array_get(arrayHandle C.uintptr_t, index C.size_t) C.uintptr_t {
	a := fatstdStringArrayFromHandle(uintptr(arrayHandle))
	if index > C.size_t(2147483647) {
		panic("fatstd_go_string_array_get: index too large")
	}
	value := a.Get(int(index))
	return C.uintptr_t(fatstdStringNewFromGoString(value))
}

//export fatstd_go_string_join
func fatstd_go_string_join(arrayHandle C.uintptr_t, sepHandle C.uintptr_t) C.uintptr_t {
	a := fatstdStringArrayFromHandle(uintptr(arrayHandle))
	sep := fatstdStringFromHandle(uintptr(sepHandle))
	return C.uintptr_t(fatstdStringNewFromGoString(fatstrings.Join(a.Values(), sep.Value())))
}

//export fatstd_go_string_replace
func fatstd_go_string_replace(sHandle C.uintptr_t, oldHandle C.uintptr_t, newHandle C.uintptr_t, n C.int) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	old := fatstdStringFromHandle(uintptr(oldHandle))
	newValue := fatstdStringFromHandle(uintptr(newHandle))
	replaced := fatstrings.Replace(s.Value(), old.Value(), newValue.Value(), int(n))
	return C.uintptr_t(fatstdStringNewFromGoString(replaced))
}

//export fatstd_go_string_replace_all
func fatstd_go_string_replace_all(sHandle C.uintptr_t, oldHandle C.uintptr_t, newHandle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	old := fatstdStringFromHandle(uintptr(oldHandle))
	newValue := fatstdStringFromHandle(uintptr(newHandle))
	replaced := fatstrings.ReplaceAll(s.Value(), old.Value(), newValue.Value())
	return C.uintptr_t(fatstdStringNewFromGoString(replaced))
}

//export fatstd_go_string_array_free
func fatstd_go_string_array_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_string_array_free: handle is 0")
	}

	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_string_array_free: invalid handle")
	}
	if _, ok := value.(*fatstrings.StringArray); !ok {
		panic("fatstd_go_string_array_free: handle is not a fat string array")
	}
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

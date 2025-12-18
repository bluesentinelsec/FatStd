package main

/*
#include <stddef.h>
#include <stdbool.h>
#include <stdint.h>
*/
import "C"

import (
	"unsafe"

	"github.com/bluesentinelsec/FatStd/pkg/fatbytes"
)

func fatstdBytesNewFromGoBytes(value []byte) uintptr {
	return fatstdHandles.register(fatbytes.New(value))
}

func fatstdBytesArrayNew(values [][]byte) uintptr {
	return fatstdHandles.register(fatbytes.NewArray(values))
}

func fatstdBytesFromHandle(handle uintptr) *fatbytes.Bytes {
	if handle == 0 {
		panic("fatstdBytesFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdBytesFromHandle: invalid handle")
	}
	b, ok := value.(*fatbytes.Bytes)
	if !ok {
		panic("fatstdBytesFromHandle: handle is not fat bytes")
	}
	return b
}

func fatstdBytesArrayFromHandle(handle uintptr) *fatbytes.BytesArray {
	if handle == 0 {
		panic("fatstdBytesArrayFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdBytesArrayFromHandle: invalid handle")
	}
	a, ok := value.(*fatbytes.BytesArray)
	if !ok {
		panic("fatstdBytesArrayFromHandle: handle is not fat bytes array")
	}
	return a
}

//export fatstd_go_bytes_new_n
func fatstd_go_bytes_new_n(bytesPtr *C.char, len C.size_t) C.uintptr_t {
	if bytesPtr == nil {
		if len == 0 {
			return C.uintptr_t(fatstdBytesNewFromGoBytes([]byte{}))
		}
		panic("fatstd_go_bytes_new_n: bytes is NULL but len > 0")
	}
	if len > C.size_t(2147483647) {
		panic("fatstd_go_bytes_new_n: len too large")
	}
	buf := C.GoBytes(unsafe.Pointer(bytesPtr), C.int(len))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(buf))
}

//export fatstd_go_bytes_len
func fatstd_go_bytes_len(handle C.uintptr_t) C.size_t {
	b := fatstdBytesFromHandle(uintptr(handle))
	return C.size_t(len(b.Value()))
}

//export fatstd_go_bytes_copy_out
func fatstd_go_bytes_copy_out(handle C.uintptr_t, dst *C.char, dstLen C.size_t) C.size_t {
	if dst == nil {
		if dstLen == 0 {
			return 0
		}
		panic("fatstd_go_bytes_copy_out: dst is NULL but dstLen > 0")
	}
	if dstLen > C.size_t(2147483647) {
		panic("fatstd_go_bytes_copy_out: dstLen too large")
	}

	b := fatstdBytesFromHandle(uintptr(handle))
	src := b.Value()
	n := len(src)
	if n > int(dstLen) {
		n = int(dstLen)
	}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(dst)), int(dstLen))[:n], src[:n])
	return C.size_t(n)
}

//export fatstd_go_bytes_clone
func fatstd_go_bytes_clone(handle C.uintptr_t) C.uintptr_t {
	b := fatstdBytesFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(fatbytes.Clone(b.Value())))
}

//export fatstd_go_bytes_contains
func fatstd_go_bytes_contains(aHandle C.uintptr_t, bHandle C.uintptr_t) C.int {
	a := fatstdBytesFromHandle(uintptr(aHandle))
	sub := fatstdBytesFromHandle(uintptr(bHandle))
	if fatbytes.Contains(a.Value(), sub.Value()) {
		return 1
	}
	return 0
}

//export fatstd_go_bytes_has_prefix
func fatstd_go_bytes_has_prefix(sHandle C.uintptr_t, prefixHandle C.uintptr_t) C.int {
	s := fatstdBytesFromHandle(uintptr(sHandle))
	prefix := fatstdBytesFromHandle(uintptr(prefixHandle))
	if fatbytes.HasPrefix(s.Value(), prefix.Value()) {
		return 1
	}
	return 0
}

//export fatstd_go_bytes_has_suffix
func fatstd_go_bytes_has_suffix(sHandle C.uintptr_t, suffixHandle C.uintptr_t) C.int {
	s := fatstdBytesFromHandle(uintptr(sHandle))
	suffix := fatstdBytesFromHandle(uintptr(suffixHandle))
	if fatbytes.HasSuffix(s.Value(), suffix.Value()) {
		return 1
	}
	return 0
}

//export fatstd_go_bytes_trim_space
func fatstd_go_bytes_trim_space(sHandle C.uintptr_t) C.uintptr_t {
	s := fatstdBytesFromHandle(uintptr(sHandle))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(fatbytes.TrimSpace(s.Value())))
}

//export fatstd_go_bytes_trim
func fatstd_go_bytes_trim(sHandle C.uintptr_t, cutsetHandle C.uintptr_t) C.uintptr_t {
	s := fatstdBytesFromHandle(uintptr(sHandle))
	cutset := fatstdStringFromHandle(uintptr(cutsetHandle))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(fatbytes.Trim(s.Value(), cutset.Value())))
}

//export fatstd_go_bytes_trim_prefix
func fatstd_go_bytes_trim_prefix(sHandle C.uintptr_t, prefixHandle C.uintptr_t) C.uintptr_t {
	s := fatstdBytesFromHandle(uintptr(sHandle))
	prefix := fatstdBytesFromHandle(uintptr(prefixHandle))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(fatbytes.TrimPrefix(s.Value(), prefix.Value())))
}

//export fatstd_go_bytes_trim_suffix
func fatstd_go_bytes_trim_suffix(sHandle C.uintptr_t, suffixHandle C.uintptr_t) C.uintptr_t {
	s := fatstdBytesFromHandle(uintptr(sHandle))
	suffix := fatstdBytesFromHandle(uintptr(suffixHandle))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(fatbytes.TrimSuffix(s.Value(), suffix.Value())))
}

//export fatstd_go_bytes_cut
func fatstd_go_bytes_cut(sHandle C.uintptr_t, sepHandle C.uintptr_t, beforeOut *C.uintptr_t, afterOut *C.uintptr_t) C.int {
	if beforeOut == nil {
		panic("fatstd_go_bytes_cut: beforeOut is NULL")
	}
	if afterOut == nil {
		panic("fatstd_go_bytes_cut: afterOut is NULL")
	}

	s := fatstdBytesFromHandle(uintptr(sHandle))
	sep := fatstdBytesFromHandle(uintptr(sepHandle))
	before, after, found := fatbytes.Cut(s.Value(), sep.Value())

	*beforeOut = C.uintptr_t(fatstdBytesNewFromGoBytes(before))
	*afterOut = C.uintptr_t(fatstdBytesNewFromGoBytes(after))

	if found {
		return 1
	}
	return 0
}

//export fatstd_go_bytes_cut_prefix
func fatstd_go_bytes_cut_prefix(sHandle C.uintptr_t, prefixHandle C.uintptr_t, afterOut *C.uintptr_t) C.int {
	if afterOut == nil {
		panic("fatstd_go_bytes_cut_prefix: afterOut is NULL")
	}

	s := fatstdBytesFromHandle(uintptr(sHandle))
	prefix := fatstdBytesFromHandle(uintptr(prefixHandle))
	after, found := fatbytes.CutPrefix(s.Value(), prefix.Value())

	*afterOut = C.uintptr_t(fatstdBytesNewFromGoBytes(after))
	if found {
		return 1
	}
	return 0
}

//export fatstd_go_bytes_cut_suffix
func fatstd_go_bytes_cut_suffix(sHandle C.uintptr_t, suffixHandle C.uintptr_t, afterOut *C.uintptr_t) C.int {
	if afterOut == nil {
		panic("fatstd_go_bytes_cut_suffix: afterOut is NULL")
	}

	s := fatstdBytesFromHandle(uintptr(sHandle))
	suffix := fatstdBytesFromHandle(uintptr(suffixHandle))
	after, found := fatbytes.CutSuffix(s.Value(), suffix.Value())

	*afterOut = C.uintptr_t(fatstdBytesNewFromGoBytes(after))
	if found {
		return 1
	}
	return 0
}

//export fatstd_go_bytes_split
func fatstd_go_bytes_split(sHandle C.uintptr_t, sepHandle C.uintptr_t) C.uintptr_t {
	s := fatstdBytesFromHandle(uintptr(sHandle))
	sep := fatstdBytesFromHandle(uintptr(sepHandle))
	return C.uintptr_t(fatstdBytesArrayNew(fatbytes.Split(s.Value(), sep.Value())))
}

//export fatstd_go_bytes_fields
func fatstd_go_bytes_fields(sHandle C.uintptr_t) C.uintptr_t {
	s := fatstdBytesFromHandle(uintptr(sHandle))
	return C.uintptr_t(fatstdBytesArrayNew(fatbytes.Fields(s.Value())))
}

//export fatstd_go_bytes_array_len
func fatstd_go_bytes_array_len(arrayHandle C.uintptr_t) C.size_t {
	a := fatstdBytesArrayFromHandle(uintptr(arrayHandle))
	return C.size_t(a.Len())
}

//export fatstd_go_bytes_array_get
func fatstd_go_bytes_array_get(arrayHandle C.uintptr_t, index C.size_t) C.uintptr_t {
	a := fatstdBytesArrayFromHandle(uintptr(arrayHandle))
	if index > C.size_t(2147483647) {
		panic("fatstd_go_bytes_array_get: index too large")
	}
	return C.uintptr_t(fatstdBytesNewFromGoBytes(a.Get(int(index))))
}

//export fatstd_go_bytes_join
func fatstd_go_bytes_join(arrayHandle C.uintptr_t, sepHandle C.uintptr_t) C.uintptr_t {
	a := fatstdBytesArrayFromHandle(uintptr(arrayHandle))
	sep := fatstdBytesFromHandle(uintptr(sepHandle))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(fatbytes.Join(a.Values(), sep.Value())))
}

//export fatstd_go_bytes_replace_all
func fatstd_go_bytes_replace_all(sHandle C.uintptr_t, oldHandle C.uintptr_t, newHandle C.uintptr_t) C.uintptr_t {
	s := fatstdBytesFromHandle(uintptr(sHandle))
	old := fatstdBytesFromHandle(uintptr(oldHandle))
	newValue := fatstdBytesFromHandle(uintptr(newHandle))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(fatbytes.ReplaceAll(s.Value(), old.Value(), newValue.Value())))
}

//export fatstd_go_bytes_replace
func fatstd_go_bytes_replace(sHandle C.uintptr_t, oldHandle C.uintptr_t, newHandle C.uintptr_t, n C.int) C.uintptr_t {
	s := fatstdBytesFromHandle(uintptr(sHandle))
	old := fatstdBytesFromHandle(uintptr(oldHandle))
	newValue := fatstdBytesFromHandle(uintptr(newHandle))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(fatbytes.Replace(s.Value(), old.Value(), newValue.Value(), int(n))))
}

//export fatstd_go_bytes_repeat
func fatstd_go_bytes_repeat(sHandle C.uintptr_t, count C.int) C.uintptr_t {
	s := fatstdBytesFromHandle(uintptr(sHandle))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(fatbytes.Repeat(s.Value(), int(count))))
}

//export fatstd_go_bytes_to_lower
func fatstd_go_bytes_to_lower(sHandle C.uintptr_t) C.uintptr_t {
	s := fatstdBytesFromHandle(uintptr(sHandle))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(fatbytes.ToLower(s.Value())))
}

//export fatstd_go_bytes_to_upper
func fatstd_go_bytes_to_upper(sHandle C.uintptr_t) C.uintptr_t {
	s := fatstdBytesFromHandle(uintptr(sHandle))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(fatbytes.ToUpper(s.Value())))
}

//export fatstd_go_bytes_index_byte
func fatstd_go_bytes_index_byte(sHandle C.uintptr_t, c C.uchar) C.int {
	s := fatstdBytesFromHandle(uintptr(sHandle))
	return C.int(fatbytes.IndexByte(s.Value(), byte(c)))
}

//export fatstd_go_bytes_index_any
func fatstd_go_bytes_index_any(sHandle C.uintptr_t, charsHandle C.uintptr_t) C.int {
	s := fatstdBytesFromHandle(uintptr(sHandle))
	chars := fatstdStringFromHandle(uintptr(charsHandle))
	return C.int(fatbytes.IndexAny(s.Value(), chars.Value()))
}

//export fatstd_go_bytes_to_valid_utf8
func fatstd_go_bytes_to_valid_utf8(sHandle C.uintptr_t, replacementHandle C.uintptr_t) C.uintptr_t {
	s := fatstdBytesFromHandle(uintptr(sHandle))
	repl := fatstdBytesFromHandle(uintptr(replacementHandle))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(fatbytes.ToValidUTF8(s.Value(), repl.Value())))
}

//export fatstd_go_bytes_index
func fatstd_go_bytes_index(sHandle C.uintptr_t, sepHandle C.uintptr_t) C.int {
	s := fatstdBytesFromHandle(uintptr(sHandle))
	sep := fatstdBytesFromHandle(uintptr(sepHandle))
	return C.int(fatbytes.Index(s.Value(), sep.Value()))
}

//export fatstd_go_bytes_count
func fatstd_go_bytes_count(sHandle C.uintptr_t, sepHandle C.uintptr_t) C.int {
	s := fatstdBytesFromHandle(uintptr(sHandle))
	sep := fatstdBytesFromHandle(uintptr(sepHandle))
	return C.int(fatbytes.Count(s.Value(), sep.Value()))
}

//export fatstd_go_bytes_compare
func fatstd_go_bytes_compare(aHandle C.uintptr_t, bHandle C.uintptr_t) C.int {
	a := fatstdBytesFromHandle(uintptr(aHandle))
	b := fatstdBytesFromHandle(uintptr(bHandle))
	return C.int(fatbytes.Compare(a.Value(), b.Value()))
}

//export fatstd_go_bytes_equal
func fatstd_go_bytes_equal(aHandle C.uintptr_t, bHandle C.uintptr_t) C.int {
	a := fatstdBytesFromHandle(uintptr(aHandle))
	b := fatstdBytesFromHandle(uintptr(bHandle))
	if fatbytes.Equal(a.Value(), b.Value()) {
		return 1
	}
	return 0
}

//export fatstd_go_bytes_free
func fatstd_go_bytes_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_bytes_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_bytes_free: invalid handle")
	}
	if _, ok := value.(*fatbytes.Bytes); !ok {
		panic("fatstd_go_bytes_free: handle is not fat bytes")
	}
}

//export fatstd_go_bytes_array_free
func fatstd_go_bytes_array_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_bytes_array_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_bytes_array_free: invalid handle")
	}
	if _, ok := value.(*fatbytes.BytesArray); !ok {
		panic("fatstd_go_bytes_array_free: handle is not fat bytes array")
	}
}

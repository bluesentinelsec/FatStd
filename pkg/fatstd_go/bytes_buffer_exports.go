package main

/*
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
*/
import "C"

import (
	"io"
	"unsafe"

	"github.com/bluesentinelsec/FatStd/pkg/fatbytes"
)

func fatstdBytesBufferNewFromGoBytes(value []byte) uintptr {
	return fatstdHandles.register(fatbytes.NewBuffer(value))
}

func fatstdBytesBufferNewFromGoString(value string) uintptr {
	return fatstdHandles.register(fatbytes.NewBufferString(value))
}

func fatstdBytesBufferFromHandle(handle uintptr) *fatbytes.Buffer {
	if handle == 0 {
		panic("fatstdBytesBufferFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdBytesBufferFromHandle: invalid handle")
	}
	b, ok := value.(*fatbytes.Buffer)
	if !ok {
		panic("fatstdBytesBufferFromHandle: handle is not fat bytes buffer")
	}
	return b
}

//export fatstd_go_bytes_buffer_new
func fatstd_go_bytes_buffer_new() C.uintptr_t {
	return C.uintptr_t(fatstdBytesBufferNewFromGoBytes([]byte{}))
}

//export fatstd_go_bytes_buffer_new_bytes
func fatstd_go_bytes_buffer_new_bytes(bytesHandle C.uintptr_t) C.uintptr_t {
	b := fatstdBytesFromHandle(uintptr(bytesHandle))
	return C.uintptr_t(fatstdBytesBufferNewFromGoBytes(fatbytes.Clone(b.Value())))
}

//export fatstd_go_bytes_buffer_new_n
func fatstd_go_bytes_buffer_new_n(bytesPtr *C.char, len C.size_t) C.uintptr_t {
	if bytesPtr == nil {
		if len == 0 {
			return C.uintptr_t(fatstdBytesBufferNewFromGoBytes([]byte{}))
		}
		panic("fatstd_go_bytes_buffer_new_n: bytes is NULL but len > 0")
	}
	if len > C.size_t(2147483647) {
		panic("fatstd_go_bytes_buffer_new_n: len too large")
	}
	buf := C.GoBytes(unsafe.Pointer(bytesPtr), C.int(len))
	return C.uintptr_t(fatstdBytesBufferNewFromGoBytes(buf))
}

//export fatstd_go_bytes_buffer_new_string
func fatstd_go_bytes_buffer_new_string(sHandle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	return C.uintptr_t(fatstdBytesBufferNewFromGoString(s.Value()))
}

//export fatstd_go_bytes_buffer_len
func fatstd_go_bytes_buffer_len(handle C.uintptr_t) C.size_t {
	b := fatstdBytesBufferFromHandle(uintptr(handle))
	return C.size_t(b.Len())
}

//export fatstd_go_bytes_buffer_cap
func fatstd_go_bytes_buffer_cap(handle C.uintptr_t) C.size_t {
	b := fatstdBytesBufferFromHandle(uintptr(handle))
	return C.size_t(b.Cap())
}

//export fatstd_go_bytes_buffer_grow
func fatstd_go_bytes_buffer_grow(handle C.uintptr_t, n C.size_t) {
	if n > C.size_t(2147483647) {
		panic("fatstd_go_bytes_buffer_grow: n too large")
	}
	b := fatstdBytesBufferFromHandle(uintptr(handle))
	b.Grow(int(n))
}

//export fatstd_go_bytes_buffer_reset
func fatstd_go_bytes_buffer_reset(handle C.uintptr_t) {
	b := fatstdBytesBufferFromHandle(uintptr(handle))
	b.Reset()
}

//export fatstd_go_bytes_buffer_truncate
func fatstd_go_bytes_buffer_truncate(handle C.uintptr_t, n C.size_t) {
	if n > C.size_t(2147483647) {
		panic("fatstd_go_bytes_buffer_truncate: n too large")
	}
	b := fatstdBytesBufferFromHandle(uintptr(handle))
	b.Truncate(int(n))
}

//export fatstd_go_bytes_buffer_write
func fatstd_go_bytes_buffer_write(handle C.uintptr_t, bytes *C.char, len C.size_t) C.size_t {
	if bytes == nil {
		if len == 0 {
			return 0
		}
		panic("fatstd_go_bytes_buffer_write: bytes is NULL but len > 0")
	}
	if len > C.size_t(2147483647) {
		panic("fatstd_go_bytes_buffer_write: len too large")
	}
	b := fatstdBytesBufferFromHandle(uintptr(handle))
	buf := C.GoBytes(unsafe.Pointer(bytes), C.int(len))
	return C.size_t(b.Write(buf))
}

//export fatstd_go_bytes_buffer_write_byte
func fatstd_go_bytes_buffer_write_byte(handle C.uintptr_t, c C.uchar) {
	b := fatstdBytesBufferFromHandle(uintptr(handle))
	b.WriteByte(byte(c))
}

//export fatstd_go_bytes_buffer_write_rune
func fatstd_go_bytes_buffer_write_rune(handle C.uintptr_t, r C.uint32_t) C.size_t {
	b := fatstdBytesBufferFromHandle(uintptr(handle))
	return C.size_t(b.WriteRune(rune(r)))
}

//export fatstd_go_bytes_buffer_write_string
func fatstd_go_bytes_buffer_write_string(bufferHandle C.uintptr_t, sHandle C.uintptr_t) C.size_t {
	b := fatstdBytesBufferFromHandle(uintptr(bufferHandle))
	s := fatstdStringFromHandle(uintptr(sHandle))
	return C.size_t(b.WriteString(s.Value()))
}

//export fatstd_go_bytes_buffer_bytes
func fatstd_go_bytes_buffer_bytes(handle C.uintptr_t) C.uintptr_t {
	b := fatstdBytesBufferFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(fatbytes.Clone(b.Bytes())))
}

//export fatstd_go_bytes_buffer_string
func fatstd_go_bytes_buffer_string(handle C.uintptr_t) C.uintptr_t {
	b := fatstdBytesBufferFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdStringNewFromGoString(b.String()))
}

//export fatstd_go_bytes_buffer_read
func fatstd_go_bytes_buffer_read(handle C.uintptr_t, dst *C.char, dstLen C.size_t, eofOut *C.bool) C.size_t {
	if eofOut == nil {
		panic("fatstd_go_bytes_buffer_read: eofOut is NULL")
	}
	if dst == nil {
		if dstLen == 0 {
			*eofOut = false
			return 0
		}
		panic("fatstd_go_bytes_buffer_read: dst is NULL but dstLen > 0")
	}
	if dstLen > C.size_t(2147483647) {
		panic("fatstd_go_bytes_buffer_read: dstLen too large")
	}

	b := fatstdBytesBufferFromHandle(uintptr(handle))
	out := unsafe.Slice((*byte)(unsafe.Pointer(dst)), int(dstLen))

	n, err := b.Read(out)
	if err == io.EOF {
		*eofOut = true
		return C.size_t(n)
	}
	if err != nil {
		panic("fatstd_go_bytes_buffer_read: unexpected error")
	}
	*eofOut = false
	return C.size_t(n)
}

//export fatstd_go_bytes_buffer_next
func fatstd_go_bytes_buffer_next(handle C.uintptr_t, n C.size_t) C.uintptr_t {
	if n > C.size_t(2147483647) {
		panic("fatstd_go_bytes_buffer_next: n too large")
	}
	b := fatstdBytesBufferFromHandle(uintptr(handle))
	chunk := b.Next(int(n))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(fatbytes.Clone(chunk)))
}

//export fatstd_go_bytes_buffer_read_byte
func fatstd_go_bytes_buffer_read_byte(handle C.uintptr_t, byteOut *C.uchar, eofOut *C.bool) C.bool {
	if byteOut == nil {
		panic("fatstd_go_bytes_buffer_read_byte: byteOut is NULL")
	}
	if eofOut == nil {
		panic("fatstd_go_bytes_buffer_read_byte: eofOut is NULL")
	}

	b := fatstdBytesBufferFromHandle(uintptr(handle))
	c, err := b.ReadByte()
	if err == io.EOF {
		*eofOut = true
		return false
	}
	if err != nil {
		panic("fatstd_go_bytes_buffer_read_byte: unexpected error")
	}
	*byteOut = C.uchar(c)
	*eofOut = false
	return true
}

//export fatstd_go_bytes_buffer_write_to_bytes_buffer
func fatstd_go_bytes_buffer_write_to_bytes_buffer(srcHandle C.uintptr_t, dstHandle C.uintptr_t) C.longlong {
	src := fatstdBytesBufferFromHandle(uintptr(srcHandle))
	dst := fatstdBytesBufferFromHandle(uintptr(dstHandle))
	n, err := src.WriteTo(dst.Underlying())
	if err != nil {
		panic("fatstd_go_bytes_buffer_write_to_bytes_buffer: unexpected error")
	}
	return C.longlong(n)
}

//export fatstd_go_bytes_buffer_read_from_string_reader
func fatstd_go_bytes_buffer_read_from_string_reader(dstHandle C.uintptr_t, readerHandle C.uintptr_t) C.longlong {
	dst := fatstdBytesBufferFromHandle(uintptr(dstHandle))
	r := fatstdStringReaderFromHandle(uintptr(readerHandle))
	n, err := dst.ReadFrom(r)
	if err != nil {
		panic("fatstd_go_bytes_buffer_read_from_string_reader: unexpected error")
	}
	return C.longlong(n)
}

//export fatstd_go_bytes_buffer_free
func fatstd_go_bytes_buffer_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_bytes_buffer_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_bytes_buffer_free: invalid handle")
	}
	if _, ok := value.(*fatbytes.Buffer); !ok {
		panic("fatstd_go_bytes_buffer_free: handle is not fat bytes buffer")
	}
}


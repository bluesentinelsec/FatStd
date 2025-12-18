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

func fatstdBytesReaderNewFromGoBytes(value []byte) uintptr {
	return fatstdHandles.register(fatbytes.NewReader(value))
}

func fatstdBytesReaderFromHandle(handle uintptr) *fatbytes.Reader {
	if handle == 0 {
		panic("fatstdBytesReaderFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdBytesReaderFromHandle: invalid handle")
	}
	r, ok := value.(*fatbytes.Reader)
	if !ok {
		panic("fatstdBytesReaderFromHandle: handle is not fat bytes reader")
	}
	return r
}

//export fatstd_go_bytes_reader_new
func fatstd_go_bytes_reader_new(bytesHandle C.uintptr_t) C.uintptr_t {
	b := fatstdBytesFromHandle(uintptr(bytesHandle))
	return C.uintptr_t(fatstdBytesReaderNewFromGoBytes(fatbytes.Clone(b.Value())))
}

//export fatstd_go_bytes_reader_len
func fatstd_go_bytes_reader_len(handle C.uintptr_t) C.size_t {
	r := fatstdBytesReaderFromHandle(uintptr(handle))
	return C.size_t(r.Len())
}

//export fatstd_go_bytes_reader_size
func fatstd_go_bytes_reader_size(handle C.uintptr_t) C.longlong {
	r := fatstdBytesReaderFromHandle(uintptr(handle))
	return C.longlong(r.Size())
}

//export fatstd_go_bytes_reader_reset
func fatstd_go_bytes_reader_reset(readerHandle C.uintptr_t, bytesHandle C.uintptr_t) {
	r := fatstdBytesReaderFromHandle(uintptr(readerHandle))
	b := fatstdBytesFromHandle(uintptr(bytesHandle))
	r.Reset(b.Value())
}

//export fatstd_go_bytes_reader_read
func fatstd_go_bytes_reader_read(readerHandle C.uintptr_t, dst *C.char, dstLen C.size_t, eofOut *C.bool) C.size_t {
	if eofOut == nil {
		panic("fatstd_go_bytes_reader_read: eofOut is NULL")
	}
	if dst == nil {
		if dstLen == 0 {
			*eofOut = false
			return 0
		}
		panic("fatstd_go_bytes_reader_read: dst is NULL but dstLen > 0")
	}
	if dstLen > C.size_t(2147483647) {
		panic("fatstd_go_bytes_reader_read: dstLen too large")
	}

	r := fatstdBytesReaderFromHandle(uintptr(readerHandle))
	out := unsafe.Slice((*byte)(unsafe.Pointer(dst)), int(dstLen))

	n, err := r.Read(out)
	if err == io.EOF {
		*eofOut = true
		return C.size_t(n)
	}
	if err != nil {
		panic("fatstd_go_bytes_reader_read: unexpected error")
	}
	*eofOut = false
	return C.size_t(n)
}

//export fatstd_go_bytes_reader_read_at
func fatstd_go_bytes_reader_read_at(readerHandle C.uintptr_t, dst *C.char, dstLen C.size_t, off C.longlong, eofOut *C.bool) C.size_t {
	if eofOut == nil {
		panic("fatstd_go_bytes_reader_read_at: eofOut is NULL")
	}
	if dst == nil {
		if dstLen == 0 {
			*eofOut = false
			return 0
		}
		panic("fatstd_go_bytes_reader_read_at: dst is NULL but dstLen > 0")
	}
	if dstLen > C.size_t(2147483647) {
		panic("fatstd_go_bytes_reader_read_at: dstLen too large")
	}

	r := fatstdBytesReaderFromHandle(uintptr(readerHandle))
	out := unsafe.Slice((*byte)(unsafe.Pointer(dst)), int(dstLen))

	n, err := r.ReadAt(out, int64(off))
	if err == io.EOF {
		*eofOut = true
		return C.size_t(n)
	}
	if err != nil {
		panic("fatstd_go_bytes_reader_read_at: unexpected error")
	}
	*eofOut = false
	return C.size_t(n)
}

//export fatstd_go_bytes_reader_read_byte
func fatstd_go_bytes_reader_read_byte(readerHandle C.uintptr_t, byteOut *C.uchar, eofOut *C.bool) C.bool {
	if byteOut == nil {
		panic("fatstd_go_bytes_reader_read_byte: byteOut is NULL")
	}
	if eofOut == nil {
		panic("fatstd_go_bytes_reader_read_byte: eofOut is NULL")
	}

	r := fatstdBytesReaderFromHandle(uintptr(readerHandle))
	b, err := r.ReadByte()
	if err == io.EOF {
		*eofOut = true
		return false
	}
	if err != nil {
		panic("fatstd_go_bytes_reader_read_byte: unexpected error")
	}
	*byteOut = C.uchar(b)
	*eofOut = false
	return true
}

//export fatstd_go_bytes_reader_unread_byte
func fatstd_go_bytes_reader_unread_byte(readerHandle C.uintptr_t) {
	r := fatstdBytesReaderFromHandle(uintptr(readerHandle))
	if err := r.UnreadByte(); err != nil {
		panic("fatstd_go_bytes_reader_unread_byte: invalid unread")
	}
}

//export fatstd_go_bytes_reader_seek
func fatstd_go_bytes_reader_seek(readerHandle C.uintptr_t, offset C.longlong, whence C.int) C.longlong {
	r := fatstdBytesReaderFromHandle(uintptr(readerHandle))
	pos, err := r.Seek(int64(offset), int(whence))
	if err != nil {
		panic("fatstd_go_bytes_reader_seek: invalid seek")
	}
	return C.longlong(pos)
}

//export fatstd_go_bytes_reader_write_to_bytes_buffer
func fatstd_go_bytes_reader_write_to_bytes_buffer(readerHandle C.uintptr_t, bufferHandle C.uintptr_t) C.longlong {
	r := fatstdBytesReaderFromHandle(uintptr(readerHandle))
	b := fatstdBytesBufferFromHandle(uintptr(bufferHandle))
	n, err := r.WriteTo(b.Underlying())
	if err != nil {
		panic("fatstd_go_bytes_reader_write_to_bytes_buffer: unexpected error")
	}
	return C.longlong(n)
}

//export fatstd_go_bytes_reader_free
func fatstd_go_bytes_reader_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_bytes_reader_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_bytes_reader_free: invalid handle")
	}
	if _, ok := value.(*fatbytes.Reader); !ok {
		panic("fatstd_go_bytes_reader_free: handle is not fat bytes reader")
	}
}


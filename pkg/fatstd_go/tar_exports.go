package main

/*
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
*/
import "C"

import (
	"archive/tar"
	"bytes"
	"errors"
	"io"
	"os"
	"unsafe"
)

const (
	fatStatusEOF = 3

	fatTarErrCodeIO  = 110
	fatTarErrCodeTar = 111
)

func fatstdTarStatusFromError(err error) C.int {
	if err == nil {
		return fatStatusOK
	}
	if errors.Is(err, io.EOF) {
		return fatStatusEOF
	}
	if errors.Is(err, tar.ErrHeader) || errors.Is(err, tar.ErrInsecurePath) {
		return fatStatusSyntax
	}
	return fatStatusOther
}

type fatTarReader struct {
	file *os.File
	tr   *tar.Reader
	data []byte
}

type fatTarHeader struct {
	hdr *tar.Header
}

type fatTarWriter struct {
	tw *tar.Writer
}

func fatstdTarReaderFromHandle(handle uintptr) *fatTarReader {
	if handle == 0 {
		panic("fatstdTarReaderFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdTarReaderFromHandle: invalid handle")
	}
	r, ok := value.(*fatTarReader)
	if !ok {
		panic("fatstdTarReaderFromHandle: handle is not tar reader")
	}
	return r
}

func fatstdTarHeaderFromHandle(handle uintptr) *fatTarHeader {
	if handle == 0 {
		panic("fatstdTarHeaderFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdTarHeaderFromHandle: invalid handle")
	}
	h, ok := value.(*fatTarHeader)
	if !ok {
		panic("fatstdTarHeaderFromHandle: handle is not tar header")
	}
	return h
}

func fatstdTarWriterFromHandle(handle uintptr) *fatTarWriter {
	if handle == 0 {
		panic("fatstdTarWriterFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdTarWriterFromHandle: invalid handle")
	}
	w, ok := value.(*fatTarWriter)
	if !ok {
		panic("fatstdTarWriterFromHandle: handle is not tar writer")
	}
	return w
}

//export fatstd_go_tar_reader_new_bytes
func fatstd_go_tar_reader_new_bytes(bytesHandle C.uintptr_t, outReader *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outReader == nil {
		panic("fatstd_go_tar_reader_new_bytes: outReader is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_tar_reader_new_bytes: outErr is NULL")
	}

	b := fatstdBytesFromHandle(uintptr(bytesHandle))
	data := append([]byte(nil), b.Value()...)
	tr := tar.NewReader(bytes.NewReader(data))

	handle := fatstdHandles.register(&fatTarReader{file: nil, tr: tr, data: data})
	*outReader = C.uintptr_t(handle)
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_tar_reader_open_path_utf8
func fatstd_go_tar_reader_open_path_utf8(path *C.char, outReader *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outReader == nil {
		panic("fatstd_go_tar_reader_open_path_utf8: outReader is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_tar_reader_open_path_utf8: outErr is NULL")
	}
	if path == nil {
		panic("fatstd_go_tar_reader_open_path_utf8: path is NULL")
	}

	f, err := os.Open(C.GoString(path))
	if err != nil {
		*outReader = 0
		*outErr = C.uintptr_t(fatstdNewError(fatTarErrCodeIO, err.Error()))
		return fatstdTarStatusFromError(err)
	}

	tr := tar.NewReader(f)
	handle := fatstdHandles.register(&fatTarReader{file: f, tr: tr, data: nil})
	*outReader = C.uintptr_t(handle)
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_tar_reader_free
func fatstd_go_tar_reader_free(handle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_tar_reader_free: outErr is NULL")
	}
	if handle == 0 {
		panic("fatstd_go_tar_reader_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_tar_reader_free: invalid handle")
	}
	r, ok := value.(*fatTarReader)
	if !ok {
		panic("fatstd_go_tar_reader_free: handle is not tar reader")
	}
	if r.file != nil {
		if err := r.file.Close(); err != nil {
			*outErr = C.uintptr_t(fatstdNewError(fatTarErrCodeIO, err.Error()))
			return fatStatusOther
		}
	}
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_tar_reader_next
func fatstd_go_tar_reader_next(handle C.uintptr_t, outHdr *C.uintptr_t, outEOF *C.bool, outErr *C.uintptr_t) C.int {
	if outHdr == nil {
		panic("fatstd_go_tar_reader_next: outHdr is NULL")
	}
	if outEOF == nil {
		panic("fatstd_go_tar_reader_next: outEOF is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_tar_reader_next: outErr is NULL")
	}

	r := fatstdTarReaderFromHandle(uintptr(handle))
	hdr, err := r.tr.Next()
	if err == io.EOF {
		*outHdr = 0
		*outEOF = true
		*outErr = 0
		return fatStatusEOF
	}
	if err != nil {
		*outHdr = 0
		*outEOF = false
		*outErr = C.uintptr_t(fatstdNewError(fatTarErrCodeTar, err.Error()))
		return fatstdTarStatusFromError(err)
	}

	hdrCopy := *hdr
	*outHdr = C.uintptr_t(fatstdHandles.register(&fatTarHeader{hdr: &hdrCopy}))
	*outEOF = false
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_tar_reader_read
func fatstd_go_tar_reader_read(
	handle C.uintptr_t,
	dst *C.char,
	dstLen C.size_t,
	outN *C.size_t,
	outEOF *C.bool,
	outErr *C.uintptr_t,
) C.int {
	if outN == nil {
		panic("fatstd_go_tar_reader_read: outN is NULL")
	}
	if outEOF == nil {
		panic("fatstd_go_tar_reader_read: outEOF is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_tar_reader_read: outErr is NULL")
	}
	if dst == nil && dstLen != 0 {
		panic("fatstd_go_tar_reader_read: dst is NULL but dstLen > 0")
	}
	if dstLen > C.size_t(2147483647) {
		panic("fatstd_go_tar_reader_read: dstLen too large")
	}

	r := fatstdTarReaderFromHandle(uintptr(handle))
	buf := unsafe.Slice((*byte)(unsafe.Pointer(dst)), int(dstLen))
	n, err := r.tr.Read(buf)
	*outN = C.size_t(n)

	if err == io.EOF {
		*outEOF = true
		*outErr = 0
		return fatStatusEOF
	}
	if err != nil {
		*outEOF = false
		*outErr = C.uintptr_t(fatstdNewError(fatTarErrCodeIO, err.Error()))
		return fatstdTarStatusFromError(err)
	}

	*outEOF = false
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_tar_header_free
func fatstd_go_tar_header_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_tar_header_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_tar_header_free: invalid handle")
	}
	if _, ok := value.(*fatTarHeader); !ok {
		panic("fatstd_go_tar_header_free: handle is not tar header")
	}
}

//export fatstd_go_tar_header_name
func fatstd_go_tar_header_name(handle C.uintptr_t) C.uintptr_t {
	h := fatstdTarHeaderFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdStringNewFromGoString(h.hdr.Name))
}

//export fatstd_go_tar_header_typeflag
func fatstd_go_tar_header_typeflag(handle C.uintptr_t) C.uchar {
	h := fatstdTarHeaderFromHandle(uintptr(handle))
	return C.uchar(h.hdr.Typeflag)
}

//export fatstd_go_tar_header_size
func fatstd_go_tar_header_size(handle C.uintptr_t) C.longlong {
	h := fatstdTarHeaderFromHandle(uintptr(handle))
	return C.longlong(h.hdr.Size)
}

//export fatstd_go_tar_writer_new_to_bytes_buffer
func fatstd_go_tar_writer_new_to_bytes_buffer(dstBufferHandle C.uintptr_t, outWriter *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outWriter == nil {
		panic("fatstd_go_tar_writer_new_to_bytes_buffer: outWriter is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_tar_writer_new_to_bytes_buffer: outErr is NULL")
	}

	dst := fatstdBytesBufferFromHandle(uintptr(dstBufferHandle))
	tw := tar.NewWriter(dst.Underlying())
	*outWriter = C.uintptr_t(fatstdHandles.register(&fatTarWriter{tw: tw}))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_tar_writer_add_bytes
func fatstd_go_tar_writer_add_bytes(writerHandle C.uintptr_t, nameHandle C.uintptr_t, dataHandle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_tar_writer_add_bytes: outErr is NULL")
	}

	w := fatstdTarWriterFromHandle(uintptr(writerHandle))
	name := fatstdStringFromHandle(uintptr(nameHandle))
	data := fatstdBytesFromHandle(uintptr(dataHandle))

	hdr := &tar.Header{
		Name:     name.Value(),
		Mode:     0644,
		Size:     int64(len(data.Value())),
		Typeflag: tar.TypeReg,
	}
	if err := w.tw.WriteHeader(hdr); err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatTarErrCodeTar, err.Error()))
		return fatstdTarStatusFromError(err)
	}
	if _, err := w.tw.Write(data.Value()); err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatTarErrCodeIO, err.Error()))
		return fatstdTarStatusFromError(err)
	}

	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_tar_writer_flush
func fatstd_go_tar_writer_flush(handle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_tar_writer_flush: outErr is NULL")
	}
	w := fatstdTarWriterFromHandle(uintptr(handle))
	if err := w.tw.Flush(); err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatTarErrCodeIO, err.Error()))
		return fatstdTarStatusFromError(err)
	}
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_tar_writer_close
func fatstd_go_tar_writer_close(handle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_tar_writer_close: outErr is NULL")
	}
	if handle == 0 {
		panic("fatstd_go_tar_writer_close: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_tar_writer_close: invalid handle")
	}
	w, ok := value.(*fatTarWriter)
	if !ok {
		panic("fatstd_go_tar_writer_close: handle is not tar writer")
	}
	if err := w.tw.Close(); err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatTarErrCodeIO, err.Error()))
		return fatstdTarStatusFromError(err)
	}
	*outErr = 0
	return fatStatusOK
}

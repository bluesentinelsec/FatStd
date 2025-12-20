package main

/*
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
*/
import "C"

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"unsafe"
)

const (
	fatStatusZipOK   = 0
	fatStatusZipEOF  = 3
	fatStatusZipErr  = 100
	fatZipErrCodeIO  = 100
	fatZipErrCodeZip = 101
)

func fatstdZipStatusFromError(err error) C.int {
	if err == nil {
		return fatStatusZipOK
	}
	if errors.Is(err, io.EOF) {
		return fatStatusZipEOF
	}
	// Zip parsing errors are treated as recoverable failures.
	if errors.Is(err, zip.ErrFormat) {
		return 1 // FAT_ERR_SYNTAX
	}
	return fatStatusZipErr
}

type fatZipReader struct {
	rc   *zip.ReadCloser // non-nil for path-based readers
	r    *zip.Reader     // always non-nil
	data []byte          // keeps bytes-backed archives alive
}

type fatZipFile struct {
	readerHandle uintptr
	index        int
}

type fatZipFileReader struct {
	rc io.ReadCloser
}

type fatZipWriter struct {
	w *zip.Writer
}

func fatstdZipReaderFromHandle(handle uintptr) *fatZipReader {
	if handle == 0 {
		panic("fatstdZipReaderFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdZipReaderFromHandle: invalid handle")
	}
	r, ok := value.(*fatZipReader)
	if !ok {
		panic("fatstdZipReaderFromHandle: handle is not zip reader")
	}
	return r
}

func fatstdZipFileFromHandle(handle uintptr) *fatZipFile {
	if handle == 0 {
		panic("fatstdZipFileFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdZipFileFromHandle: invalid handle")
	}
	f, ok := value.(*fatZipFile)
	if !ok {
		panic("fatstdZipFileFromHandle: handle is not zip file")
	}
	return f
}

func fatstdZipFileReaderFromHandle(handle uintptr) *fatZipFileReader {
	if handle == 0 {
		panic("fatstdZipFileReaderFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdZipFileReaderFromHandle: invalid handle")
	}
	r, ok := value.(*fatZipFileReader)
	if !ok {
		panic("fatstdZipFileReaderFromHandle: handle is not zip file reader")
	}
	return r
}

func fatstdZipWriterFromHandle(handle uintptr) *fatZipWriter {
	if handle == 0 {
		panic("fatstdZipWriterFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdZipWriterFromHandle: invalid handle")
	}
	w, ok := value.(*fatZipWriter)
	if !ok {
		panic("fatstdZipWriterFromHandle: handle is not zip writer")
	}
	return w
}

func fatstdZipFileResolve(f *fatZipFile) *zip.File {
	reader := fatstdZipReaderFromHandle(f.readerHandle)
	if f.index < 0 || f.index >= len(reader.r.File) {
		panic("fatstdZipFileResolve: index out of range")
	}
	return reader.r.File[f.index]
}

//export fatstd_go_zip_reader_open_path_utf8
func fatstd_go_zip_reader_open_path_utf8(path *C.char, outReader *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outReader == nil {
		panic("fatstd_go_zip_reader_open_path_utf8: outReader is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_zip_reader_open_path_utf8: outErr is NULL")
	}
	if path == nil {
		panic("fatstd_go_zip_reader_open_path_utf8: path is NULL")
	}

	rc, err := zip.OpenReader(C.GoString(path))
	if err != nil {
		*outReader = 0
		*outErr = C.uintptr_t(fatstdNewError(fatZipErrCodeZip, err.Error()))
		return fatstdZipStatusFromError(err)
	}

	handle := fatstdHandles.register(&fatZipReader{rc: rc, r: &rc.Reader})
	*outReader = C.uintptr_t(handle)
	*outErr = 0
	return fatStatusZipOK
}

//export fatstd_go_zip_reader_new_bytes
func fatstd_go_zip_reader_new_bytes(bytesHandle C.uintptr_t, outReader *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outReader == nil {
		panic("fatstd_go_zip_reader_new_bytes: outReader is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_zip_reader_new_bytes: outErr is NULL")
	}

	b := fatstdBytesFromHandle(uintptr(bytesHandle))
	data := append([]byte(nil), b.Value()...)
	br := bytes.NewReader(data)

	r, err := zip.NewReader(br, int64(len(data)))
	if err != nil {
		*outReader = 0
		*outErr = C.uintptr_t(fatstdNewError(fatZipErrCodeZip, err.Error()))
		return fatstdZipStatusFromError(err)
	}

	handle := fatstdHandles.register(&fatZipReader{rc: nil, r: r, data: data})
	*outReader = C.uintptr_t(handle)
	*outErr = 0
	return fatStatusZipOK
}

//export fatstd_go_zip_reader_free
func fatstd_go_zip_reader_free(handle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_zip_reader_free: outErr is NULL")
	}
	if handle == 0 {
		panic("fatstd_go_zip_reader_free: handle is 0")
	}

	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_zip_reader_free: invalid handle")
	}
	r, ok := value.(*fatZipReader)
	if !ok {
		panic("fatstd_go_zip_reader_free: handle is not zip reader")
	}

	if r.rc != nil {
		if err := r.rc.Close(); err != nil {
			*outErr = C.uintptr_t(fatstdNewError(fatZipErrCodeIO, err.Error()))
			return fatStatusZipErr
		}
	}
	*outErr = 0
	return fatStatusZipOK
}

//export fatstd_go_zip_reader_num_files
func fatstd_go_zip_reader_num_files(handle C.uintptr_t) C.size_t {
	r := fatstdZipReaderFromHandle(uintptr(handle))
	return C.size_t(len(r.r.File))
}

//export fatstd_go_zip_reader_file_by_index
func fatstd_go_zip_reader_file_by_index(readerHandle C.uintptr_t, idx C.size_t) C.uintptr_t {
	r := fatstdZipReaderFromHandle(uintptr(readerHandle))
	if idx > C.size_t(2147483647) {
		panic("fatstd_go_zip_reader_file_by_index: idx too large")
	}
	i := int(idx)
	if i < 0 || i >= len(r.r.File) {
		panic("fatstd_go_zip_reader_file_by_index: idx out of range")
	}
	return C.uintptr_t(fatstdHandles.register(&fatZipFile{
		readerHandle: uintptr(readerHandle),
		index:        i,
	}))
}

//export fatstd_go_zip_file_free
func fatstd_go_zip_file_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_zip_file_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_zip_file_free: invalid handle")
	}
	if _, ok := value.(*fatZipFile); !ok {
		panic("fatstd_go_zip_file_free: handle is not zip file")
	}
}

//export fatstd_go_zip_file_name
func fatstd_go_zip_file_name(handle C.uintptr_t) C.uintptr_t {
	f := fatstdZipFileFromHandle(uintptr(handle))
	zf := fatstdZipFileResolve(f)
	return C.uintptr_t(fatstdStringNewFromGoString(zf.Name))
}

//export fatstd_go_zip_file_open
func fatstd_go_zip_file_open(handle C.uintptr_t, outReader *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outReader == nil {
		panic("fatstd_go_zip_file_open: outReader is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_zip_file_open: outErr is NULL")
	}

	f := fatstdZipFileFromHandle(uintptr(handle))
	zf := fatstdZipFileResolve(f)
	rc, err := zf.Open()
	if err != nil {
		*outReader = 0
		*outErr = C.uintptr_t(fatstdNewError(fatZipErrCodeZip, err.Error()))
		return fatstdZipStatusFromError(err)
	}
	*outReader = C.uintptr_t(fatstdHandles.register(&fatZipFileReader{rc: rc}))
	*outErr = 0
	return fatStatusZipOK
}

//export fatstd_go_zip_file_reader_read
func fatstd_go_zip_file_reader_read(
	handle C.uintptr_t,
	dst *C.char,
	dstLen C.size_t,
	outN *C.size_t,
	outEOF *C.bool,
	outErr *C.uintptr_t,
) C.int {
	if outN == nil {
		panic("fatstd_go_zip_file_reader_read: outN is NULL")
	}
	if outEOF == nil {
		panic("fatstd_go_zip_file_reader_read: outEOF is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_zip_file_reader_read: outErr is NULL")
	}
	if dst == nil && dstLen != 0 {
		panic("fatstd_go_zip_file_reader_read: dst is NULL but dstLen > 0")
	}
	if dstLen > C.size_t(2147483647) {
		panic("fatstd_go_zip_file_reader_read: dstLen too large")
	}

	r := fatstdZipFileReaderFromHandle(uintptr(handle))
	buf := unsafe.Slice((*byte)(unsafe.Pointer(dst)), int(dstLen))
	n, err := r.rc.Read(buf)
	*outN = C.size_t(n)

	if err == io.EOF {
		*outEOF = true
		*outErr = 0
		return fatStatusZipEOF
	}
	if err != nil {
		*outEOF = false
		*outErr = C.uintptr_t(fatstdNewError(fatZipErrCodeIO, err.Error()))
		return fatStatusZipErr
	}

	*outEOF = false
	*outErr = 0
	return fatStatusZipOK
}

//export fatstd_go_zip_file_reader_close
func fatstd_go_zip_file_reader_close(handle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_zip_file_reader_close: outErr is NULL")
	}
	if handle == 0 {
		panic("fatstd_go_zip_file_reader_close: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_zip_file_reader_close: invalid handle")
	}
	r, ok := value.(*fatZipFileReader)
	if !ok {
		panic("fatstd_go_zip_file_reader_close: handle is not zip file reader")
	}
	if err := r.rc.Close(); err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatZipErrCodeIO, err.Error()))
		return fatStatusZipErr
	}
	*outErr = 0
	return fatStatusZipOK
}

//export fatstd_go_zip_writer_new_to_bytes_buffer
func fatstd_go_zip_writer_new_to_bytes_buffer(dstBufferHandle C.uintptr_t, outWriter *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outWriter == nil {
		panic("fatstd_go_zip_writer_new_to_bytes_buffer: outWriter is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_zip_writer_new_to_bytes_buffer: outErr is NULL")
	}

	dst := fatstdBytesBufferFromHandle(uintptr(dstBufferHandle))
	w := zip.NewWriter(dst.Underlying())
	*outWriter = C.uintptr_t(fatstdHandles.register(&fatZipWriter{w: w}))
	*outErr = 0
	return fatStatusZipOK
}

//export fatstd_go_zip_writer_add_bytes
func fatstd_go_zip_writer_add_bytes(writerHandle C.uintptr_t, nameHandle C.uintptr_t, dataHandle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_zip_writer_add_bytes: outErr is NULL")
	}

	w := fatstdZipWriterFromHandle(uintptr(writerHandle))
	name := fatstdStringFromHandle(uintptr(nameHandle))
	data := fatstdBytesFromHandle(uintptr(dataHandle))

	entryWriter, err := w.w.Create(name.Value())
	if err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatZipErrCodeZip, err.Error()))
		return fatstdZipStatusFromError(err)
	}

	if _, err := entryWriter.Write(data.Value()); err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatZipErrCodeIO, err.Error()))
		return fatstdZipStatusFromError(err)
	}

	*outErr = 0
	return fatStatusZipOK
}

//export fatstd_go_zip_writer_close
func fatstd_go_zip_writer_close(handle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_zip_writer_close: outErr is NULL")
	}
	if handle == 0 {
		panic("fatstd_go_zip_writer_close: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_zip_writer_close: invalid handle")
	}
	w, ok := value.(*fatZipWriter)
	if !ok {
		panic("fatstd_go_zip_writer_close: handle is not zip writer")
	}
	if err := w.w.Close(); err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatZipErrCodeIO, err.Error()))
		return fatStatusZipErr
	}
	*outErr = 0
	return fatStatusZipOK
}


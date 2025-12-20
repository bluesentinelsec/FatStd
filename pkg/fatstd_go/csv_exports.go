package main

/*
#include <stdbool.h>
#include <stdint.h>
*/
import "C"

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"unsafe"
)

const (
	fatCsvErrCodeParse = 150
	fatCsvErrCodeIO    = 151
)

type fatCsvReader struct {
	r    *csv.Reader
	data []byte
}

type fatCsvWriter struct {
	w *csv.Writer
}

func fatstdCsvReaderFromHandle(handle uintptr) *fatCsvReader {
	if handle == 0 {
		panic("fatstdCsvReaderFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdCsvReaderFromHandle: invalid handle")
	}
	r, ok := value.(*fatCsvReader)
	if !ok {
		panic("fatstdCsvReaderFromHandle: handle is not csv reader")
	}
	if r.r == nil {
		panic("fatstdCsvReaderFromHandle: reader is nil")
	}
	return r
}

func fatstdCsvWriterFromHandle(handle uintptr) *fatCsvWriter {
	if handle == 0 {
		panic("fatstdCsvWriterFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdCsvWriterFromHandle: invalid handle")
	}
	w, ok := value.(*fatCsvWriter)
	if !ok {
		panic("fatstdCsvWriterFromHandle: handle is not csv writer")
	}
	if w.w == nil {
		panic("fatstdCsvWriterFromHandle: writer is nil")
	}
	return w
}

func fatstdCsvStatusFromError(err error) C.int {
	if err == nil {
		return fatStatusOK
	}
	if errors.Is(err, io.EOF) {
		return fatStatusEOF
	}
	var pe *csv.ParseError
	if errors.As(err, &pe) {
		return fatStatusSyntax
	}
	if errors.Is(err, csv.ErrBareQuote) || errors.Is(err, csv.ErrQuote) || errors.Is(err, csv.ErrFieldCount) {
		return fatStatusSyntax
	}
	return fatStatusOther
}

//export fatstd_go_csv_reader_new_bytes
func fatstd_go_csv_reader_new_bytes(bytesHandle C.uintptr_t) C.uintptr_t {
	b := fatstdBytesFromHandle(uintptr(bytesHandle))
	data := append([]byte(nil), b.Value()...)
	r := csv.NewReader(bytes.NewReader(data))
	return C.uintptr_t(fatstdHandles.register(&fatCsvReader{r: r, data: data}))
}

//export fatstd_go_csv_reader_free
func fatstd_go_csv_reader_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_csv_reader_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_csv_reader_free: invalid handle")
	}
	if _, ok := value.(*fatCsvReader); !ok {
		panic("fatstd_go_csv_reader_free: handle is not csv reader")
	}
}

//export fatstd_go_csv_reader_read
func fatstd_go_csv_reader_read(handle C.uintptr_t, outRecord *C.uintptr_t, outEOF *C.bool, outErr *C.uintptr_t) C.int {
	if outRecord == nil {
		panic("fatstd_go_csv_reader_read: outRecord is NULL")
	}
	if outEOF == nil {
		panic("fatstd_go_csv_reader_read: outEOF is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_csv_reader_read: outErr is NULL")
	}

	r := fatstdCsvReaderFromHandle(uintptr(handle))
	rec, err := r.r.Read()
	if err == io.EOF {
		*outRecord = 0
		*outEOF = true
		*outErr = 0
		return fatStatusEOF
	}
	if err != nil {
		*outRecord = 0
		*outEOF = false
		*outErr = C.uintptr_t(fatstdNewError(fatCsvErrCodeParse, err.Error()))
		return fatstdCsvStatusFromError(err)
	}

	*outRecord = C.uintptr_t(fatstdStringArrayNew(rec))
	*outEOF = false
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_csv_reader_field_pos
func fatstd_go_csv_reader_field_pos(handle C.uintptr_t, field C.int, outLine *C.int, outCol *C.int) {
	if outLine == nil {
		panic("fatstd_go_csv_reader_field_pos: outLine is NULL")
	}
	if outCol == nil {
		panic("fatstd_go_csv_reader_field_pos: outCol is NULL")
	}
	r := fatstdCsvReaderFromHandle(uintptr(handle))
	line, col := r.r.FieldPos(int(field))
	*outLine = C.int(line)
	*outCol = C.int(col)
}

//export fatstd_go_csv_reader_input_offset
func fatstd_go_csv_reader_input_offset(handle C.uintptr_t) C.longlong {
	r := fatstdCsvReaderFromHandle(uintptr(handle))
	return C.longlong(r.r.InputOffset())
}

//export fatstd_go_csv_writer_new_to_bytes_buffer
func fatstd_go_csv_writer_new_to_bytes_buffer(dstBufferHandle C.uintptr_t) C.uintptr_t {
	dst := fatstdBytesBufferFromHandle(uintptr(dstBufferHandle))
	w := csv.NewWriter(dst.Underlying())
	return C.uintptr_t(fatstdHandles.register(&fatCsvWriter{w: w}))
}

//export fatstd_go_csv_writer_write_record
func fatstd_go_csv_writer_write_record(handle C.uintptr_t, fields *C.uintptr_t, n C.size_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_csv_writer_write_record: outErr is NULL")
	}
	if fields == nil && n != 0 {
		panic("fatstd_go_csv_writer_write_record: fields is NULL but n > 0")
	}
	if n > C.size_t(2147483647) {
		panic("fatstd_go_csv_writer_write_record: n too large")
	}

	w := fatstdCsvWriterFromHandle(uintptr(handle))
	fieldHandles := unsafe.Slice((*C.uintptr_t)(unsafe.Pointer(fields)), int(n))
	record := make([]string, int(n))
	for i := 0; i < int(n); i++ {
		s := fatstdStringFromHandle(uintptr(fieldHandles[i]))
		record[i] = s.Value()
	}

	if err := w.w.Write(record); err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatCsvErrCodeIO, err.Error()))
		return fatStatusOther
	}
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_csv_writer_flush
func fatstd_go_csv_writer_flush(handle C.uintptr_t) {
	w := fatstdCsvWriterFromHandle(uintptr(handle))
	w.w.Flush()
}

//export fatstd_go_csv_writer_error
func fatstd_go_csv_writer_error(handle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_csv_writer_error: outErr is NULL")
	}
	w := fatstdCsvWriterFromHandle(uintptr(handle))
	if err := w.w.Error(); err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatCsvErrCodeIO, err.Error()))
		return fatStatusOther
	}
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_csv_writer_free
func fatstd_go_csv_writer_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_csv_writer_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_csv_writer_free: invalid handle")
	}
	if _, ok := value.(*fatCsvWriter); !ok {
		panic("fatstd_go_csv_writer_free: handle is not csv writer")
	}
}

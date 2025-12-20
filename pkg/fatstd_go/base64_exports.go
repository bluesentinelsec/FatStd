package main

/*
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
*/
import "C"

import (
	"encoding/base64"
	"errors"
	"io"
	"unicode/utf8"
	"unsafe"

	"github.com/bluesentinelsec/FatStd/pkg/fatbytes"
)

const (
	fatBase64ErrCodeConfig = 130
	fatBase64ErrCodeDecode = 131
	fatBase64ErrCodeIO     = 132
)

type fatBase64Encoding struct {
	enc *base64.Encoding
}

type fatBase64Encoder struct {
	w io.WriteCloser
}

func fatstdBase64EncodingFromHandle(handle uintptr) *base64.Encoding {
	if handle == 0 {
		panic("fatstdBase64EncodingFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdBase64EncodingFromHandle: invalid handle")
	}
	e, ok := value.(*fatBase64Encoding)
	if !ok {
		panic("fatstdBase64EncodingFromHandle: handle is not base64 encoding")
	}
	if e.enc == nil {
		panic("fatstdBase64EncodingFromHandle: encoding is nil")
	}
	return e.enc
}

func fatstdBase64EncoderFromHandle(handle uintptr) *fatBase64Encoder {
	if handle == 0 {
		panic("fatstdBase64EncoderFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdBase64EncoderFromHandle: invalid handle")
	}
	e, ok := value.(*fatBase64Encoder)
	if !ok {
		panic("fatstdBase64EncoderFromHandle: handle is not base64 encoder")
	}
	if e.w == nil {
		panic("fatstdBase64EncoderFromHandle: encoder writer is nil")
	}
	return e
}

func fatstdBase64StatusFromError(err error) C.int {
	if err == nil {
		return fatStatusOK
	}
	var cie base64.CorruptInputError
	if errors.As(err, &cie) {
		return fatStatusSyntax
	}
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return fatStatusSyntax
	}
	return fatStatusOther
}

//export fatstd_go_base64_encoding_new_utf8
func fatstd_go_base64_encoding_new_utf8(alphabet *C.char, outEnc *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outEnc == nil {
		panic("fatstd_go_base64_encoding_new_utf8: outEnc is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_base64_encoding_new_utf8: outErr is NULL")
	}
	if alphabet == nil {
		panic("fatstd_go_base64_encoding_new_utf8: alphabet is NULL")
	}

	s := C.GoString(alphabet)
	if len(s) != 64 {
		*outEnc = 0
		*outErr = C.uintptr_t(fatstdNewError(fatBase64ErrCodeConfig, "base64: alphabet must be exactly 64 bytes"))
		return fatStatusRange
	}

	enc := base64.NewEncoding(s)
	*outEnc = C.uintptr_t(fatstdHandles.register(&fatBase64Encoding{enc: enc}))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_base64_encoding_strict
func fatstd_go_base64_encoding_strict(encHandle C.uintptr_t) C.uintptr_t {
	enc := fatstdBase64EncodingFromHandle(uintptr(encHandle))
	return C.uintptr_t(fatstdHandles.register(&fatBase64Encoding{enc: enc.Strict()}))
}

//export fatstd_go_base64_encoding_with_padding
func fatstd_go_base64_encoding_with_padding(encHandle C.uintptr_t, padding C.int32_t, outEnc *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outEnc == nil {
		panic("fatstd_go_base64_encoding_with_padding: outEnc is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_base64_encoding_with_padding: outErr is NULL")
	}

	enc := fatstdBase64EncodingFromHandle(uintptr(encHandle))
	p := rune(padding)
	if p != base64.NoPadding && (p < 0 || p > utf8.MaxRune) {
		*outEnc = 0
		*outErr = C.uintptr_t(fatstdNewError(fatBase64ErrCodeConfig, "base64: invalid padding rune"))
		return fatStatusRange
	}

	defer func() {
		if r := recover(); r != nil {
			*outEnc = 0
			*outErr = C.uintptr_t(fatstdNewError(fatBase64ErrCodeConfig, "base64: invalid padding"))
		}
	}()

	out := enc.WithPadding(p)
	*outEnc = C.uintptr_t(fatstdHandles.register(&fatBase64Encoding{enc: out}))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_base64_encoded_len
func fatstd_go_base64_encoded_len(encHandle C.uintptr_t, n C.int) C.int {
	enc := fatstdBase64EncodingFromHandle(uintptr(encHandle))
	return C.int(enc.EncodedLen(int(n)))
}

//export fatstd_go_base64_decoded_len
func fatstd_go_base64_decoded_len(encHandle C.uintptr_t, n C.int) C.int {
	enc := fatstdBase64EncodingFromHandle(uintptr(encHandle))
	return C.int(enc.DecodedLen(int(n)))
}

//export fatstd_go_base64_encode_to_string
func fatstd_go_base64_encode_to_string(encHandle C.uintptr_t, srcHandle C.uintptr_t) C.uintptr_t {
	enc := fatstdBase64EncodingFromHandle(uintptr(encHandle))
	src := fatstdBytesFromHandle(uintptr(srcHandle))
	return C.uintptr_t(fatstdStringNewFromGoString(enc.EncodeToString(src.Value())))
}

//export fatstd_go_base64_encode
func fatstd_go_base64_encode(encHandle C.uintptr_t, srcHandle C.uintptr_t) C.uintptr_t {
	enc := fatstdBase64EncodingFromHandle(uintptr(encHandle))
	src := fatstdBytesFromHandle(uintptr(srcHandle))
	dst := make([]byte, enc.EncodedLen(len(src.Value())))
	enc.Encode(dst, src.Value())
	return C.uintptr_t(fatstdBytesNewFromGoBytes(dst))
}

//export fatstd_go_base64_decode_string
func fatstd_go_base64_decode_string(encHandle C.uintptr_t, sHandle C.uintptr_t, outBytes *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outBytes == nil {
		panic("fatstd_go_base64_decode_string: outBytes is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_base64_decode_string: outErr is NULL")
	}

	enc := fatstdBase64EncodingFromHandle(uintptr(encHandle))
	s := fatstdStringFromHandle(uintptr(sHandle))
	b, err := enc.DecodeString(s.Value())
	if err != nil {
		*outBytes = 0
		*outErr = C.uintptr_t(fatstdNewError(fatBase64ErrCodeDecode, err.Error()))
		return fatstdBase64StatusFromError(err)
	}
	*outBytes = C.uintptr_t(fatstdBytesNewFromGoBytes(b))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_base64_decode
func fatstd_go_base64_decode(encHandle C.uintptr_t, srcHandle C.uintptr_t, outBytes *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outBytes == nil {
		panic("fatstd_go_base64_decode: outBytes is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_base64_decode: outErr is NULL")
	}

	enc := fatstdBase64EncodingFromHandle(uintptr(encHandle))
	src := fatstdBytesFromHandle(uintptr(srcHandle))
	buf := make([]byte, enc.DecodedLen(len(src.Value())))
	n, err := enc.Decode(buf, src.Value())
	if err != nil {
		*outBytes = 0
		*outErr = C.uintptr_t(fatstdNewError(fatBase64ErrCodeDecode, err.Error()))
		return fatstdBase64StatusFromError(err)
	}
	*outBytes = C.uintptr_t(fatstdBytesNewFromGoBytes(buf[:n]))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_base64_append_encode
func fatstd_go_base64_append_encode(encHandle C.uintptr_t, dstHandle C.uintptr_t, srcHandle C.uintptr_t) C.uintptr_t {
	enc := fatstdBase64EncodingFromHandle(uintptr(encHandle))
	dst := fatstdBytesFromHandle(uintptr(dstHandle))
	src := fatstdBytesFromHandle(uintptr(srcHandle))
	dstBytes := fatbytes.Clone(dst.Value())
	need := enc.EncodedLen(len(src.Value()))
	out := append(dstBytes, make([]byte, need)...)
	enc.Encode(out[len(dstBytes):], src.Value())
	return C.uintptr_t(fatstdBytesNewFromGoBytes(out))
}

//export fatstd_go_base64_append_decode
func fatstd_go_base64_append_decode(encHandle C.uintptr_t, dstHandle C.uintptr_t, srcHandle C.uintptr_t, outBytes *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outBytes == nil {
		panic("fatstd_go_base64_append_decode: outBytes is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_base64_append_decode: outErr is NULL")
	}

	enc := fatstdBase64EncodingFromHandle(uintptr(encHandle))
	dst := fatstdBytesFromHandle(uintptr(dstHandle))
	src := fatstdBytesFromHandle(uintptr(srcHandle))

	buf := make([]byte, enc.DecodedLen(len(src.Value())))
	n, err := enc.Decode(buf, src.Value())
	if err != nil {
		*outBytes = 0
		*outErr = C.uintptr_t(fatstdNewError(fatBase64ErrCodeDecode, err.Error()))
		return fatstdBase64StatusFromError(err)
	}
	out := append(fatbytes.Clone(dst.Value()), buf[:n]...)
	*outBytes = C.uintptr_t(fatstdBytesNewFromGoBytes(out))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_base64_encoding_free
func fatstd_go_base64_encoding_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_base64_encoding_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_base64_encoding_free: invalid handle")
	}
	if _, ok := value.(*fatBase64Encoding); !ok {
		panic("fatstd_go_base64_encoding_free: handle is not base64 encoding")
	}
}

//export fatstd_go_base64_encoder_new_to_bytes_buffer
func fatstd_go_base64_encoder_new_to_bytes_buffer(encHandle C.uintptr_t, dstBufferHandle C.uintptr_t, outEncoder *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outEncoder == nil {
		panic("fatstd_go_base64_encoder_new_to_bytes_buffer: outEncoder is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_base64_encoder_new_to_bytes_buffer: outErr is NULL")
	}
	enc := fatstdBase64EncodingFromHandle(uintptr(encHandle))
	dst := fatstdBytesBufferFromHandle(uintptr(dstBufferHandle))
	w := base64.NewEncoder(enc, dst.Underlying())
	*outEncoder = C.uintptr_t(fatstdHandles.register(&fatBase64Encoder{w: w}))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_base64_encoder_write
func fatstd_go_base64_encoder_write(handle C.uintptr_t, bytesPtr *C.char, len C.size_t, outN *C.size_t, outErr *C.uintptr_t) C.int {
	if outN == nil {
		panic("fatstd_go_base64_encoder_write: outN is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_base64_encoder_write: outErr is NULL")
	}
	if bytesPtr == nil && len != 0 {
		panic("fatstd_go_base64_encoder_write: bytes is NULL but len > 0")
	}
	if len > C.size_t(2147483647) {
		panic("fatstd_go_base64_encoder_write: len too large")
	}

	e := fatstdBase64EncoderFromHandle(uintptr(handle))
	buf := unsafe.Slice((*byte)(unsafe.Pointer(bytesPtr)), int(len))
	n, err := e.w.Write(buf)
	*outN = C.size_t(n)
	if err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatBase64ErrCodeIO, err.Error()))
		return fatstdBase64StatusFromError(err)
	}
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_base64_encoder_close
func fatstd_go_base64_encoder_close(handle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_base64_encoder_close: outErr is NULL")
	}
	if handle == 0 {
		panic("fatstd_go_base64_encoder_close: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_base64_encoder_close: invalid handle")
	}
	e, ok := value.(*fatBase64Encoder)
	if !ok {
		panic("fatstd_go_base64_encoder_close: handle is not base64 encoder")
	}
	if err := e.w.Close(); err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatBase64ErrCodeIO, err.Error()))
		return fatstdBase64StatusFromError(err)
	}
	*outErr = 0
	return fatStatusOK
}

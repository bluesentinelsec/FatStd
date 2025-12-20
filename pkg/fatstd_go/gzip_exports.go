package main

/*
#include <stdint.h>
*/
import "C"

import (
	"errors"
	"compress/gzip"

	fatgzip "github.com/bluesentinelsec/FatStd/pkg/compress/gzip"
)

const (
	fatGzipErrCode = 220
)

func fatstdGzipStatusFromError(err error) C.int {
	if err == nil {
		return fatStatusCompressOK
	}
	if errors.Is(err, gzip.ErrHeader) {
		return fatStatusCompressSyntax
	}
	return fatStatusCompressErr
}

//export fatstd_go_gzip_compress
func fatstd_go_gzip_compress(bytesHandle C.uintptr_t, outBytes *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outBytes == nil {
		panic("fatstd_go_gzip_compress: outBytes is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_gzip_compress: outErr is NULL")
	}

	b := fatstdBytesFromHandle(uintptr(bytesHandle))
	compressed, err := fatgzip.Compress(b.Value())
	if err != nil {
		*outBytes = 0
		*outErr = C.uintptr_t(fatstdNewError(fatGzipErrCode, err.Error()))
		return fatStatusCompressErr
	}

	*outBytes = C.uintptr_t(fatstdBytesNewFromGoBytes(compressed))
	*outErr = 0
	return fatStatusCompressOK
}

//export fatstd_go_gzip_decompress
func fatstd_go_gzip_decompress(bytesHandle C.uintptr_t, outBytes *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outBytes == nil {
		panic("fatstd_go_gzip_decompress: outBytes is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_gzip_decompress: outErr is NULL")
	}

	b := fatstdBytesFromHandle(uintptr(bytesHandle))
	decompressed, err := fatgzip.Decompress(b.Value())
	if err != nil {
		*outBytes = 0
		*outErr = C.uintptr_t(fatstdNewError(fatGzipErrCode, err.Error()))
		return fatstdGzipStatusFromError(err)
	}

	*outBytes = C.uintptr_t(fatstdBytesNewFromGoBytes(decompressed))
	*outErr = 0
	return fatStatusCompressOK
}

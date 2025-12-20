package main

/*
#include <stdint.h>
*/
import "C"

import (
	"compress/zlib"
	"errors"

	fatzlib "github.com/bluesentinelsec/FatStd/pkg/compress/zlib"
)

const (
	fatZlibErrCode = 240
)

func fatstdZlibStatusFromError(err error) C.int {
	if err == nil {
		return fatStatusCompressOK
	}
	if errors.Is(err, zlib.ErrHeader) {
		return fatStatusCompressSyntax
	}
	return fatStatusCompressErr
}

//export fatstd_go_zlib_compress
func fatstd_go_zlib_compress(bytesHandle C.uintptr_t, outBytes *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outBytes == nil {
		panic("fatstd_go_zlib_compress: outBytes is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_zlib_compress: outErr is NULL")
	}

	b := fatstdBytesFromHandle(uintptr(bytesHandle))
	compressed, err := fatzlib.Compress(b.Value())
	if err != nil {
		*outBytes = 0
		*outErr = C.uintptr_t(fatstdNewError(fatZlibErrCode, err.Error()))
		return fatStatusCompressErr
	}

	*outBytes = C.uintptr_t(fatstdBytesNewFromGoBytes(compressed))
	*outErr = 0
	return fatStatusCompressOK
}

//export fatstd_go_zlib_decompress
func fatstd_go_zlib_decompress(bytesHandle C.uintptr_t, outBytes *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outBytes == nil {
		panic("fatstd_go_zlib_decompress: outBytes is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_zlib_decompress: outErr is NULL")
	}

	b := fatstdBytesFromHandle(uintptr(bytesHandle))
	decompressed, err := fatzlib.Decompress(b.Value())
	if err != nil {
		*outBytes = 0
		*outErr = C.uintptr_t(fatstdNewError(fatZlibErrCode, err.Error()))
		return fatstdZlibStatusFromError(err)
	}

	*outBytes = C.uintptr_t(fatstdBytesNewFromGoBytes(decompressed))
	*outErr = 0
	return fatStatusCompressOK
}

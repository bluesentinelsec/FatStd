package main

/*
#include <stdint.h>
*/
import "C"

import (
	"github.com/bluesentinelsec/FatStd/pkg/compress/flate"
)

const (
	fatFlateErrCode = 210
)

//export fatstd_go_flate_compress
func fatstd_go_flate_compress(bytesHandle C.uintptr_t, outBytes *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outBytes == nil {
		panic("fatstd_go_flate_compress: outBytes is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_flate_compress: outErr is NULL")
	}

	b := fatstdBytesFromHandle(uintptr(bytesHandle))
	compressed, err := flate.Compress(b.Value())
	if err != nil {
		*outBytes = 0
		*outErr = C.uintptr_t(fatstdNewError(fatFlateErrCode, err.Error()))
		return fatStatusCompressErr
	}

	*outBytes = C.uintptr_t(fatstdBytesNewFromGoBytes(compressed))
	*outErr = 0
	return fatStatusCompressOK
}

//export fatstd_go_flate_decompress
func fatstd_go_flate_decompress(bytesHandle C.uintptr_t, outBytes *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outBytes == nil {
		panic("fatstd_go_flate_decompress: outBytes is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_flate_decompress: outErr is NULL")
	}

	b := fatstdBytesFromHandle(uintptr(bytesHandle))
	decompressed, err := flate.Decompress(b.Value())
	if err != nil {
		*outBytes = 0
		*outErr = C.uintptr_t(fatstdNewError(fatFlateErrCode, err.Error()))
		return fatStatusCompressSyntax
	}

	*outBytes = C.uintptr_t(fatstdBytesNewFromGoBytes(decompressed))
	*outErr = 0
	return fatStatusCompressOK
}

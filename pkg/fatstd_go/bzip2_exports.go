package main

/*
#include <stdint.h>
*/
import "C"

import (
	"github.com/bluesentinelsec/FatStd/pkg/compress/bzip2"
)

const (
	fatBzip2ErrCode = 200
)

//export fatstd_go_bzip2_decompress
func fatstd_go_bzip2_decompress(bytesHandle C.uintptr_t, outBytes *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outBytes == nil {
		panic("fatstd_go_bzip2_decompress: outBytes is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_bzip2_decompress: outErr is NULL")
	}

	b := fatstdBytesFromHandle(uintptr(bytesHandle))
	decompressed, err := bzip2.Decompress(b.Value())
	if err != nil {
		*outBytes = 0
		*outErr = C.uintptr_t(fatstdNewError(fatBzip2ErrCode, err.Error()))
		return fatStatusCompressSyntax
	}

	*outBytes = C.uintptr_t(fatstdBytesNewFromGoBytes(decompressed))
	*outErr = 0
	return fatStatusCompressOK
}

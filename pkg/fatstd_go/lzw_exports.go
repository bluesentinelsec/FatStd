package main

/*
#include <stdint.h>
*/
import "C"

import (
	"compress/lzw"

	fatlzw "github.com/bluesentinelsec/FatStd/pkg/compress/lzw"
)

const (
	fatLzwErrCode = 230
)

func fatstdLzwOrderFromC(order C.int, outErr *C.uintptr_t) (lzw.Order, bool) {
	switch order {
	case 0:
		return lzw.LSB, true
	case 1:
		return lzw.MSB, true
	default:
		*outErr = C.uintptr_t(fatstdNewError(fatLzwErrCode, "invalid LZW order"))
		return lzw.LSB, false
	}
}

func fatstdLzwLitWidthFromC(litWidth C.uchar, outErr *C.uintptr_t) (int, bool) {
	width := int(litWidth)
	if width < 2 || width > 8 {
		*outErr = C.uintptr_t(fatstdNewError(fatLzwErrCode, "invalid LZW literal width"))
		return 0, false
	}
	return width, true
}

//export fatstd_go_lzw_compress
func fatstd_go_lzw_compress(bytesHandle C.uintptr_t, order C.int, litWidth C.uchar, outBytes *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outBytes == nil {
		panic("fatstd_go_lzw_compress: outBytes is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_lzw_compress: outErr is NULL")
	}

	resolvedOrder, ok := fatstdLzwOrderFromC(order, outErr)
	if !ok {
		*outBytes = 0
		return fatStatusCompressRange
	}
	resolvedWidth, ok := fatstdLzwLitWidthFromC(litWidth, outErr)
	if !ok {
		*outBytes = 0
		return fatStatusCompressRange
	}

	b := fatstdBytesFromHandle(uintptr(bytesHandle))
	compressed, err := fatlzw.Compress(b.Value(), resolvedOrder, resolvedWidth)
	if err != nil {
		*outBytes = 0
		*outErr = C.uintptr_t(fatstdNewError(fatLzwErrCode, err.Error()))
		return fatStatusCompressErr
	}

	*outBytes = C.uintptr_t(fatstdBytesNewFromGoBytes(compressed))
	*outErr = 0
	return fatStatusCompressOK
}

//export fatstd_go_lzw_decompress
func fatstd_go_lzw_decompress(bytesHandle C.uintptr_t, order C.int, litWidth C.uchar, outBytes *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outBytes == nil {
		panic("fatstd_go_lzw_decompress: outBytes is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_lzw_decompress: outErr is NULL")
	}

	resolvedOrder, ok := fatstdLzwOrderFromC(order, outErr)
	if !ok {
		*outBytes = 0
		return fatStatusCompressRange
	}
	resolvedWidth, ok := fatstdLzwLitWidthFromC(litWidth, outErr)
	if !ok {
		*outBytes = 0
		return fatStatusCompressRange
	}

	b := fatstdBytesFromHandle(uintptr(bytesHandle))
	decompressed, err := fatlzw.Decompress(b.Value(), resolvedOrder, resolvedWidth)
	if err != nil {
		*outBytes = 0
		*outErr = C.uintptr_t(fatstdNewError(fatLzwErrCode, err.Error()))
		return fatStatusCompressSyntax
	}

	*outBytes = C.uintptr_t(fatstdBytesNewFromGoBytes(decompressed))
	*outErr = 0
	return fatStatusCompressOK
}

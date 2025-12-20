package main

/*
#include <stdbool.h>
#include <stdint.h>
*/
import "C"

import (
	"strconv"

	"github.com/bluesentinelsec/FatStd/pkg/fatbytes"
	"github.com/bluesentinelsec/FatStd/pkg/fatconv"
)

const (
	fatStatusOK      = 0
	fatStatusSyntax  = 1
	fatStatusRange   = 2
	fatStatusOther   = 100
)

func fatstdConvStatusFromError(err error) C.int {
	switch fatconv.ClassifyParseError(err) {
	case fatStatusOK:
		return fatStatusOK
	case fatStatusSyntax:
		return fatStatusSyntax
	case fatStatusRange:
		return fatStatusRange
	default:
		return fatStatusOther
	}
}

//export fatstd_go_conv_append_bool
func fatstd_go_conv_append_bool(dstHandle C.uintptr_t, b C.int) C.uintptr_t {
	dst := fatstdBytesFromHandle(uintptr(dstHandle))
	out := fatconv.AppendBool(fatbytes.Clone(dst.Value()), b != 0)
	return C.uintptr_t(fatstdBytesNewFromGoBytes(out))
}

//export fatstd_go_conv_append_int
func fatstd_go_conv_append_int(dstHandle C.uintptr_t, i C.longlong, base C.int) C.uintptr_t {
	dst := fatstdBytesFromHandle(uintptr(dstHandle))
	out := fatconv.AppendInt(fatbytes.Clone(dst.Value()), int64(i), int(base))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(out))
}

//export fatstd_go_conv_append_uint
func fatstd_go_conv_append_uint(dstHandle C.uintptr_t, i C.ulonglong, base C.int) C.uintptr_t {
	dst := fatstdBytesFromHandle(uintptr(dstHandle))
	out := fatconv.AppendUint(fatbytes.Clone(dst.Value()), uint64(i), int(base))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(out))
}

//export fatstd_go_conv_append_float
func fatstd_go_conv_append_float(dstHandle C.uintptr_t, f C.double, fmt C.uchar, prec C.int, bitSize C.int) C.uintptr_t {
	dst := fatstdBytesFromHandle(uintptr(dstHandle))
	out := fatconv.AppendFloat(fatbytes.Clone(dst.Value()), float64(f), byte(fmt), int(prec), int(bitSize))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(out))
}

//export fatstd_go_conv_append_quote
func fatstd_go_conv_append_quote(dstHandle C.uintptr_t, sHandle C.uintptr_t) C.uintptr_t {
	dst := fatstdBytesFromHandle(uintptr(dstHandle))
	s := fatstdStringFromHandle(uintptr(sHandle))
	out := fatconv.AppendQuote(fatbytes.Clone(dst.Value()), s.Value())
	return C.uintptr_t(fatstdBytesNewFromGoBytes(out))
}

//export fatstd_go_conv_append_quote_rune
func fatstd_go_conv_append_quote_rune(dstHandle C.uintptr_t, r C.uint32_t) C.uintptr_t {
	dst := fatstdBytesFromHandle(uintptr(dstHandle))
	out := fatconv.AppendQuoteRune(fatbytes.Clone(dst.Value()), rune(r))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(out))
}

//export fatstd_go_conv_append_quote_rune_to_ascii
func fatstd_go_conv_append_quote_rune_to_ascii(dstHandle C.uintptr_t, r C.uint32_t) C.uintptr_t {
	dst := fatstdBytesFromHandle(uintptr(dstHandle))
	out := fatconv.AppendQuoteRuneToASCII(fatbytes.Clone(dst.Value()), rune(r))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(out))
}

//export fatstd_go_conv_append_quote_rune_to_graphic
func fatstd_go_conv_append_quote_rune_to_graphic(dstHandle C.uintptr_t, r C.uint32_t) C.uintptr_t {
	dst := fatstdBytesFromHandle(uintptr(dstHandle))
	out := fatconv.AppendQuoteRuneToGraphic(fatbytes.Clone(dst.Value()), rune(r))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(out))
}

//export fatstd_go_conv_append_quote_to_ascii
func fatstd_go_conv_append_quote_to_ascii(dstHandle C.uintptr_t, sHandle C.uintptr_t) C.uintptr_t {
	dst := fatstdBytesFromHandle(uintptr(dstHandle))
	s := fatstdStringFromHandle(uintptr(sHandle))
	out := fatconv.AppendQuoteToASCII(fatbytes.Clone(dst.Value()), s.Value())
	return C.uintptr_t(fatstdBytesNewFromGoBytes(out))
}

//export fatstd_go_conv_append_quote_to_graphic
func fatstd_go_conv_append_quote_to_graphic(dstHandle C.uintptr_t, sHandle C.uintptr_t) C.uintptr_t {
	dst := fatstdBytesFromHandle(uintptr(dstHandle))
	s := fatstdStringFromHandle(uintptr(sHandle))
	out := fatconv.AppendQuoteToGraphic(fatbytes.Clone(dst.Value()), s.Value())
	return C.uintptr_t(fatstdBytesNewFromGoBytes(out))
}

//export fatstd_go_conv_can_backquote
func fatstd_go_conv_can_backquote(sHandle C.uintptr_t) C.int {
	s := fatstdStringFromHandle(uintptr(sHandle))
	if fatconv.CanBackquote(s.Value()) {
		return 1
	}
	return 0
}

//export fatstd_go_conv_format_bool
func fatstd_go_conv_format_bool(b C.int) C.uintptr_t {
	return C.uintptr_t(fatstdStringNewFromGoString(fatconv.FormatBool(b != 0)))
}

//export fatstd_go_conv_format_int
func fatstd_go_conv_format_int(i C.longlong, base C.int) C.uintptr_t {
	return C.uintptr_t(fatstdStringNewFromGoString(fatconv.FormatInt(int64(i), int(base))))
}

//export fatstd_go_conv_format_uint
func fatstd_go_conv_format_uint(i C.ulonglong, base C.int) C.uintptr_t {
	return C.uintptr_t(fatstdStringNewFromGoString(fatconv.FormatUint(uint64(i), int(base))))
}

//export fatstd_go_conv_format_float
func fatstd_go_conv_format_float(f C.double, fmt C.uchar, prec C.int, bitSize C.int) C.uintptr_t {
	return C.uintptr_t(fatstdStringNewFromGoString(fatconv.FormatFloat(float64(f), byte(fmt), int(prec), int(bitSize))))
}

//export fatstd_go_conv_format_complex
func fatstd_go_conv_format_complex(re C.double, im C.double, fmt C.uchar, prec C.int, bitSize C.int) C.uintptr_t {
	c := complex(float64(re), float64(im))
	return C.uintptr_t(fatstdStringNewFromGoString(fatconv.FormatComplex(c, byte(fmt), int(prec), int(bitSize))))
}

//export fatstd_go_conv_itoa
func fatstd_go_conv_itoa(i C.int) C.uintptr_t {
	return C.uintptr_t(fatstdStringNewFromGoString(fatconv.Itoa(int(i))))
}

//export fatstd_go_conv_quote
func fatstd_go_conv_quote(sHandle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	return C.uintptr_t(fatstdStringNewFromGoString(fatconv.Quote(s.Value())))
}

//export fatstd_go_conv_quote_rune
func fatstd_go_conv_quote_rune(r C.uint32_t) C.uintptr_t {
	return C.uintptr_t(fatstdStringNewFromGoString(fatconv.QuoteRune(rune(r))))
}

//export fatstd_go_conv_quote_rune_to_ascii
func fatstd_go_conv_quote_rune_to_ascii(r C.uint32_t) C.uintptr_t {
	return C.uintptr_t(fatstdStringNewFromGoString(fatconv.QuoteRuneToASCII(rune(r))))
}

//export fatstd_go_conv_quote_rune_to_graphic
func fatstd_go_conv_quote_rune_to_graphic(r C.uint32_t) C.uintptr_t {
	return C.uintptr_t(fatstdStringNewFromGoString(fatconv.QuoteRuneToGraphic(rune(r))))
}

//export fatstd_go_conv_quote_to_ascii
func fatstd_go_conv_quote_to_ascii(sHandle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	return C.uintptr_t(fatstdStringNewFromGoString(fatconv.QuoteToASCII(s.Value())))
}

//export fatstd_go_conv_quote_to_graphic
func fatstd_go_conv_quote_to_graphic(sHandle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	return C.uintptr_t(fatstdStringNewFromGoString(fatconv.QuoteToGraphic(s.Value())))
}

//export fatstd_go_conv_is_graphic
func fatstd_go_conv_is_graphic(r C.uint32_t) C.int {
	if fatconv.IsGraphic(rune(r)) {
		return 1
	}
	return 0
}

//export fatstd_go_conv_is_print
func fatstd_go_conv_is_print(r C.uint32_t) C.int {
	if fatconv.IsPrint(rune(r)) {
		return 1
	}
	return 0
}

//export fatstd_go_conv_parse_bool
func fatstd_go_conv_parse_bool(sHandle C.uintptr_t, outValue *C.int, outErr *C.uintptr_t) C.int {
	if outValue == nil {
		panic("fatstd_go_conv_parse_bool: outValue is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_conv_parse_bool: outErr is NULL")
	}
	s := fatstdStringFromHandle(uintptr(sHandle))
	v, err := fatconv.ParseBool(s.Value())
	if err != nil {
		*outErr = C.uintptr_t(fatstdNewError(int32(fatconv.ClassifyParseError(err)), err.Error()))
		return fatstdConvStatusFromError(err)
	}
	*outValue = 0
	if v {
		*outValue = 1
	}
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_conv_parse_int
func fatstd_go_conv_parse_int(sHandle C.uintptr_t, base C.int, bitSize C.int, outValue *C.longlong, outErr *C.uintptr_t) C.int {
	if outValue == nil {
		panic("fatstd_go_conv_parse_int: outValue is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_conv_parse_int: outErr is NULL")
	}
	s := fatstdStringFromHandle(uintptr(sHandle))
	v, err := fatconv.ParseInt(s.Value(), int(base), int(bitSize))
	if err != nil {
		*outErr = C.uintptr_t(fatstdNewError(int32(fatconv.ClassifyParseError(err)), err.Error()))
		return fatstdConvStatusFromError(err)
	}
	*outValue = C.longlong(v)
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_conv_parse_uint
func fatstd_go_conv_parse_uint(sHandle C.uintptr_t, base C.int, bitSize C.int, outValue *C.ulonglong, outErr *C.uintptr_t) C.int {
	if outValue == nil {
		panic("fatstd_go_conv_parse_uint: outValue is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_conv_parse_uint: outErr is NULL")
	}
	s := fatstdStringFromHandle(uintptr(sHandle))
	v, err := fatconv.ParseUint(s.Value(), int(base), int(bitSize))
	if err != nil {
		*outErr = C.uintptr_t(fatstdNewError(int32(fatconv.ClassifyParseError(err)), err.Error()))
		return fatstdConvStatusFromError(err)
	}
	*outValue = C.ulonglong(v)
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_conv_parse_float
func fatstd_go_conv_parse_float(sHandle C.uintptr_t, bitSize C.int, outValue *C.double, outErr *C.uintptr_t) C.int {
	if outValue == nil {
		panic("fatstd_go_conv_parse_float: outValue is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_conv_parse_float: outErr is NULL")
	}
	s := fatstdStringFromHandle(uintptr(sHandle))
	v, err := fatconv.ParseFloat(s.Value(), int(bitSize))
	if err != nil {
		*outErr = C.uintptr_t(fatstdNewError(int32(fatconv.ClassifyParseError(err)), err.Error()))
		return fatstdConvStatusFromError(err)
	}
	*outValue = C.double(v)
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_conv_parse_complex
func fatstd_go_conv_parse_complex(sHandle C.uintptr_t, bitSize C.int, outRe *C.double, outIm *C.double, outErr *C.uintptr_t) C.int {
	if outRe == nil || outIm == nil {
		panic("fatstd_go_conv_parse_complex: outRe/outIm is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_conv_parse_complex: outErr is NULL")
	}
	s := fatstdStringFromHandle(uintptr(sHandle))
	v, err := fatconv.ParseComplex(s.Value(), int(bitSize))
	if err != nil {
		*outErr = C.uintptr_t(fatstdNewError(int32(fatconv.ClassifyParseError(err)), err.Error()))
		return fatstdConvStatusFromError(err)
	}
	*outRe = C.double(real(v))
	*outIm = C.double(imag(v))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_conv_atoi
func fatstd_go_conv_atoi(sHandle C.uintptr_t, outValue *C.longlong, outErr *C.uintptr_t) C.int {
	if outValue == nil {
		panic("fatstd_go_conv_atoi: outValue is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_conv_atoi: outErr is NULL")
	}
	s := fatstdStringFromHandle(uintptr(sHandle))
	v, err := fatconv.Atoi(s.Value())
	if err != nil {
		*outErr = C.uintptr_t(fatstdNewError(int32(fatconv.ClassifyParseError(err)), err.Error()))
		return fatstdConvStatusFromError(err)
	}
	*outValue = C.longlong(v)
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_conv_unquote
func fatstd_go_conv_unquote(sHandle C.uintptr_t, outStr *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outStr == nil {
		panic("fatstd_go_conv_unquote: outStr is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_conv_unquote: outErr is NULL")
	}
	s := fatstdStringFromHandle(uintptr(sHandle))
	v, err := fatconv.Unquote(s.Value())
	if err != nil {
		*outErr = C.uintptr_t(fatstdNewError(int32(fatconv.ClassifyParseError(err)), err.Error()))
		return fatstdConvStatusFromError(err)
	}
	*outStr = C.uintptr_t(fatstdStringNewFromGoString(v))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_conv_quoted_prefix
func fatstd_go_conv_quoted_prefix(sHandle C.uintptr_t, outStr *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outStr == nil {
		panic("fatstd_go_conv_quoted_prefix: outStr is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_conv_quoted_prefix: outErr is NULL")
	}
	s := fatstdStringFromHandle(uintptr(sHandle))
	v, err := fatconv.QuotedPrefix(s.Value())
	if err != nil {
		*outErr = C.uintptr_t(fatstdNewError(int32(fatconv.ClassifyParseError(err)), err.Error()))
		return fatstdConvStatusFromError(err)
	}
	*outStr = C.uintptr_t(fatstdStringNewFromGoString(v))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_conv_unquote_char
func fatstd_go_conv_unquote_char(sHandle C.uintptr_t, quote C.uchar, outRune *C.uint32_t, outMultibyte *C.bool, outTail *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outRune == nil || outMultibyte == nil || outTail == nil || outErr == nil {
		panic("fatstd_go_conv_unquote_char: out param is NULL")
	}
	s := fatstdStringFromHandle(uintptr(sHandle))
	value, multibyte, tail, err := fatconv.UnquoteChar(s.Value(), byte(quote))
	if err != nil {
		*outErr = C.uintptr_t(fatstdNewError(int32(fatconv.ClassifyParseError(err)), err.Error()))
		return fatstdConvStatusFromError(err)
	}
	*outRune = C.uint32_t(value)
	*outMultibyte = C.bool(multibyte)
	*outTail = C.uintptr_t(fatstdStringNewFromGoString(tail))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_conv_int_size
func fatstd_go_conv_int_size() C.int {
	return C.int(strconv.IntSize)
}

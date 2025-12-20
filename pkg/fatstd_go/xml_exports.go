package main

/*
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
*/
import "C"

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"os"
)

const (
	fatXmlErrCodeSyntax = 190
	fatXmlErrCodeIO     = 191
	fatXmlErrCodeOther  = 192
)

type fatXmlDecoder struct {
	file *os.File
	dec  *xml.Decoder
	data []byte
}

type fatXmlEncoder struct {
	enc *xml.Encoder
}

type fatXmlToken struct {
	kind int
	tok  xml.Token
}

func fatstdXmlStatusFromError(err error) C.int {
	if err == nil {
		return fatStatusOK
	}
	if errors.Is(err, io.EOF) {
		return fatStatusEOF
	}
	var se *xml.SyntaxError
	if errors.As(err, &se) {
		return fatStatusSyntax
	}
	return fatStatusOther
}

func fatstdXmlErrCodeFromError(err error) int32 {
	if err == nil {
		return fatXmlErrCodeOther
	}
	if errors.Is(err, io.EOF) {
		return fatXmlErrCodeOther
	}
	var se *xml.SyntaxError
	if errors.As(err, &se) {
		return fatXmlErrCodeSyntax
	}
	return fatXmlErrCodeIO
}

func fatstdXmlDecoderFromHandle(handle uintptr) *fatXmlDecoder {
	if handle == 0 {
		panic("fatstdXmlDecoderFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdXmlDecoderFromHandle: invalid handle")
	}
	d, ok := value.(*fatXmlDecoder)
	if !ok {
		panic("fatstdXmlDecoderFromHandle: handle is not xml decoder")
	}
	if d.dec == nil {
		panic("fatstdXmlDecoderFromHandle: decoder is nil")
	}
	return d
}

func fatstdXmlEncoderFromHandle(handle uintptr) *fatXmlEncoder {
	if handle == 0 {
		panic("fatstdXmlEncoderFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdXmlEncoderFromHandle: invalid handle")
	}
	e, ok := value.(*fatXmlEncoder)
	if !ok {
		panic("fatstdXmlEncoderFromHandle: handle is not xml encoder")
	}
	if e.enc == nil {
		panic("fatstdXmlEncoderFromHandle: encoder is nil")
	}
	return e
}

func fatstdXmlTokenFromHandle(handle uintptr) *fatXmlToken {
	if handle == 0 {
		panic("fatstdXmlTokenFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdXmlTokenFromHandle: invalid handle")
	}
	t, ok := value.(*fatXmlToken)
	if !ok {
		panic("fatstdXmlTokenFromHandle: handle is not xml token")
	}
	if t.tok == nil {
		panic("fatstdXmlTokenFromHandle: token is nil")
	}
	return t
}

func fatstdXmlTokenKind(tok xml.Token) int {
	switch tok.(type) {
	case xml.StartElement:
		return 1
	case xml.EndElement:
		return 2
	case xml.CharData:
		return 3
	case xml.Comment:
		return 4
	case xml.Directive:
		return 5
	case xml.ProcInst:
		return 6
	default:
		panic("fatstdXmlTokenKind: unsupported token type")
	}
}

//export fatstd_go_xml_decoder_new_bytes
func fatstd_go_xml_decoder_new_bytes(bytesHandle C.uintptr_t) C.uintptr_t {
	b := fatstdBytesFromHandle(uintptr(bytesHandle))
	data := append([]byte(nil), b.Value()...)
	dec := xml.NewDecoder(bytes.NewReader(data))
	return C.uintptr_t(fatstdHandles.register(&fatXmlDecoder{file: nil, dec: dec, data: data}))
}

//export fatstd_go_xml_decoder_open_path_utf8
func fatstd_go_xml_decoder_open_path_utf8(path *C.char, outDec *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outDec == nil {
		panic("fatstd_go_xml_decoder_open_path_utf8: outDec is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_xml_decoder_open_path_utf8: outErr is NULL")
	}
	if path == nil {
		panic("fatstd_go_xml_decoder_open_path_utf8: path is NULL")
	}

	f, err := os.Open(C.GoString(path))
	if err != nil {
		*outDec = 0
		*outErr = C.uintptr_t(fatstdNewError(fatXmlErrCodeIO, err.Error()))
		return fatstdXmlStatusFromError(err)
	}
	dec := xml.NewDecoder(f)
	*outDec = C.uintptr_t(fatstdHandles.register(&fatXmlDecoder{file: f, dec: dec, data: nil}))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_xml_decoder_free
func fatstd_go_xml_decoder_free(handle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_xml_decoder_free: outErr is NULL")
	}
	if handle == 0 {
		panic("fatstd_go_xml_decoder_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_xml_decoder_free: invalid handle")
	}
	d, ok := value.(*fatXmlDecoder)
	if !ok {
		panic("fatstd_go_xml_decoder_free: handle is not xml decoder")
	}
	if d.file != nil {
		if err := d.file.Close(); err != nil {
			*outErr = C.uintptr_t(fatstdNewError(fatXmlErrCodeIO, err.Error()))
			return fatStatusOther
		}
	}
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_xml_decoder_input_offset
func fatstd_go_xml_decoder_input_offset(handle C.uintptr_t) C.longlong {
	d := fatstdXmlDecoderFromHandle(uintptr(handle))
	return C.longlong(d.dec.InputOffset())
}

//export fatstd_go_xml_decoder_input_pos
func fatstd_go_xml_decoder_input_pos(handle C.uintptr_t, outLine *C.int, outCol *C.int) {
	if outLine == nil {
		panic("fatstd_go_xml_decoder_input_pos: outLine is NULL")
	}
	if outCol == nil {
		panic("fatstd_go_xml_decoder_input_pos: outCol is NULL")
	}
	d := fatstdXmlDecoderFromHandle(uintptr(handle))
	line, col := d.dec.InputPos()
	*outLine = C.int(line)
	*outCol = C.int(col)
}

func fatstdXmlReadToken(tok xml.Token, err error, outTok *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if err == io.EOF {
		*outTok = 0
		*outErr = 0
		return fatStatusEOF
	}
	if err != nil {
		*outTok = 0
		*outErr = C.uintptr_t(fatstdNewError(fatstdXmlErrCodeFromError(err), err.Error()))
		return fatstdXmlStatusFromError(err)
	}
	*outTok = C.uintptr_t(fatstdHandles.register(&fatXmlToken{kind: fatstdXmlTokenKind(tok), tok: tok}))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_xml_decoder_token
func fatstd_go_xml_decoder_token(handle C.uintptr_t, outTok *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outTok == nil {
		panic("fatstd_go_xml_decoder_token: outTok is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_xml_decoder_token: outErr is NULL")
	}
	d := fatstdXmlDecoderFromHandle(uintptr(handle))
	tok, err := d.dec.Token()
	return fatstdXmlReadToken(tok, err, outTok, outErr)
}

//export fatstd_go_xml_decoder_raw_token
func fatstd_go_xml_decoder_raw_token(handle C.uintptr_t, outTok *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outTok == nil {
		panic("fatstd_go_xml_decoder_raw_token: outTok is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_xml_decoder_raw_token: outErr is NULL")
	}
	d := fatstdXmlDecoderFromHandle(uintptr(handle))
	tok, err := d.dec.RawToken()
	return fatstdXmlReadToken(tok, err, outTok, outErr)
}

//export fatstd_go_xml_decoder_skip
func fatstd_go_xml_decoder_skip(handle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_xml_decoder_skip: outErr is NULL")
	}
	d := fatstdXmlDecoderFromHandle(uintptr(handle))
	if err := d.dec.Skip(); err != nil {
		if errors.Is(err, io.EOF) {
			*outErr = 0
			return fatStatusEOF
		}
		*outErr = C.uintptr_t(fatstdNewError(fatstdXmlErrCodeFromError(err), err.Error()))
		return fatstdXmlStatusFromError(err)
	}
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_xml_token_free
func fatstd_go_xml_token_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_xml_token_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_xml_token_free: invalid handle")
	}
	if _, ok := value.(*fatXmlToken); !ok {
		panic("fatstd_go_xml_token_free: handle is not xml token")
	}
}

//export fatstd_go_xml_token_kind
func fatstd_go_xml_token_kind(handle C.uintptr_t) C.int {
	t := fatstdXmlTokenFromHandle(uintptr(handle))
	return C.int(t.kind)
}

func fatstdXmlTokenName(tok xml.Token) xml.Name {
	switch t := tok.(type) {
	case xml.StartElement:
		return t.Name
	case xml.EndElement:
		return t.Name
	default:
		panic("fatstdXmlTokenName: token is not start/end element")
	}
}

//export fatstd_go_xml_token_name_local
func fatstd_go_xml_token_name_local(handle C.uintptr_t) C.uintptr_t {
	t := fatstdXmlTokenFromHandle(uintptr(handle))
	name := fatstdXmlTokenName(t.tok)
	return C.uintptr_t(fatstdStringNewFromGoString(name.Local))
}

//export fatstd_go_xml_token_name_space
func fatstd_go_xml_token_name_space(handle C.uintptr_t) C.uintptr_t {
	t := fatstdXmlTokenFromHandle(uintptr(handle))
	name := fatstdXmlTokenName(t.tok)
	return C.uintptr_t(fatstdStringNewFromGoString(name.Space))
}

//export fatstd_go_xml_start_element_attr_count
func fatstd_go_xml_start_element_attr_count(handle C.uintptr_t) C.size_t {
	t := fatstdXmlTokenFromHandle(uintptr(handle))
	se, ok := t.tok.(xml.StartElement)
	if !ok {
		panic("fatstd_go_xml_start_element_attr_count: token is not start element")
	}
	return C.size_t(len(se.Attr))
}

//export fatstd_go_xml_start_element_attr_get
func fatstd_go_xml_start_element_attr_get(handle C.uintptr_t, idx C.size_t, outNameLocal *C.uintptr_t, outNameSpace *C.uintptr_t, outValue *C.uintptr_t) {
	if outNameLocal == nil {
		panic("fatstd_go_xml_start_element_attr_get: outNameLocal is NULL")
	}
	if outNameSpace == nil {
		panic("fatstd_go_xml_start_element_attr_get: outNameSpace is NULL")
	}
	if outValue == nil {
		panic("fatstd_go_xml_start_element_attr_get: outValue is NULL")
	}
	t := fatstdXmlTokenFromHandle(uintptr(handle))
	se, ok := t.tok.(xml.StartElement)
	if !ok {
		panic("fatstd_go_xml_start_element_attr_get: token is not start element")
	}
	if idx > C.size_t(2147483647) {
		panic("fatstd_go_xml_start_element_attr_get: idx too large")
	}
	i := int(idx)
	if i < 0 || i >= len(se.Attr) {
		panic("fatstd_go_xml_start_element_attr_get: idx out of range")
	}
	a := se.Attr[i]
	*outNameLocal = C.uintptr_t(fatstdStringNewFromGoString(a.Name.Local))
	*outNameSpace = C.uintptr_t(fatstdStringNewFromGoString(a.Name.Space))
	*outValue = C.uintptr_t(fatstdStringNewFromGoString(a.Value))
}

//export fatstd_go_xml_token_bytes
func fatstd_go_xml_token_bytes(handle C.uintptr_t) C.uintptr_t {
	t := fatstdXmlTokenFromHandle(uintptr(handle))
	switch v := t.tok.(type) {
	case xml.CharData:
		return C.uintptr_t(fatstdBytesNewFromGoBytes(bytes.Clone([]byte(v))))
	case xml.Comment:
		return C.uintptr_t(fatstdBytesNewFromGoBytes(bytes.Clone([]byte(v))))
	case xml.Directive:
		return C.uintptr_t(fatstdBytesNewFromGoBytes(bytes.Clone([]byte(v))))
	default:
		panic("fatstd_go_xml_token_bytes: token does not carry bytes")
	}
}

//export fatstd_go_xml_proc_inst_target
func fatstd_go_xml_proc_inst_target(handle C.uintptr_t) C.uintptr_t {
	t := fatstdXmlTokenFromHandle(uintptr(handle))
	pi, ok := t.tok.(xml.ProcInst)
	if !ok {
		panic("fatstd_go_xml_proc_inst_target: token is not proc inst")
	}
	return C.uintptr_t(fatstdStringNewFromGoString(pi.Target))
}

//export fatstd_go_xml_proc_inst_inst_bytes
func fatstd_go_xml_proc_inst_inst_bytes(handle C.uintptr_t) C.uintptr_t {
	t := fatstdXmlTokenFromHandle(uintptr(handle))
	pi, ok := t.tok.(xml.ProcInst)
	if !ok {
		panic("fatstd_go_xml_proc_inst_inst_bytes: token is not proc inst")
	}
	return C.uintptr_t(fatstdBytesNewFromGoBytes(bytes.Clone(pi.Inst)))
}

//export fatstd_go_xml_escape_to_bytes_buffer
func fatstd_go_xml_escape_to_bytes_buffer(dstBufferHandle C.uintptr_t, srcHandle C.uintptr_t) {
	dst := fatstdBytesBufferFromHandle(uintptr(dstBufferHandle))
	src := fatstdBytesFromHandle(uintptr(srcHandle))
	xml.Escape(dst.Underlying(), src.Value())
}

//export fatstd_go_xml_escape_text_to_bytes_buffer
func fatstd_go_xml_escape_text_to_bytes_buffer(dstBufferHandle C.uintptr_t, srcHandle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_xml_escape_text_to_bytes_buffer: outErr is NULL")
	}
	dst := fatstdBytesBufferFromHandle(uintptr(dstBufferHandle))
	src := fatstdBytesFromHandle(uintptr(srcHandle))
	if err := xml.EscapeText(dst.Underlying(), src.Value()); err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatXmlErrCodeOther, err.Error()))
		return fatstdXmlStatusFromError(err)
	}
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_xml_encoder_new_to_bytes_buffer
func fatstd_go_xml_encoder_new_to_bytes_buffer(dstBufferHandle C.uintptr_t) C.uintptr_t {
	dst := fatstdBytesBufferFromHandle(uintptr(dstBufferHandle))
	enc := xml.NewEncoder(dst.Underlying())
	return C.uintptr_t(fatstdHandles.register(&fatXmlEncoder{enc: enc}))
}

//export fatstd_go_xml_encoder_indent
func fatstd_go_xml_encoder_indent(handle C.uintptr_t, prefixHandle C.uintptr_t, indentHandle C.uintptr_t) {
	e := fatstdXmlEncoderFromHandle(uintptr(handle))
	prefix := fatstdStringFromHandle(uintptr(prefixHandle))
	indent := fatstdStringFromHandle(uintptr(indentHandle))
	e.enc.Indent(prefix.Value(), indent.Value())
}

//export fatstd_go_xml_encoder_encode_token
func fatstd_go_xml_encoder_encode_token(encHandle C.uintptr_t, tokHandle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_xml_encoder_encode_token: outErr is NULL")
	}
	e := fatstdXmlEncoderFromHandle(uintptr(encHandle))
	t := fatstdXmlTokenFromHandle(uintptr(tokHandle))
	if err := e.enc.EncodeToken(t.tok); err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatXmlErrCodeOther, err.Error()))
		return fatstdXmlStatusFromError(err)
	}
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_xml_encoder_flush
func fatstd_go_xml_encoder_flush(handle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_xml_encoder_flush: outErr is NULL")
	}
	e := fatstdXmlEncoderFromHandle(uintptr(handle))
	if err := e.enc.Flush(); err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatXmlErrCodeOther, err.Error()))
		return fatstdXmlStatusFromError(err)
	}
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_xml_encoder_close
func fatstd_go_xml_encoder_close(handle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_xml_encoder_close: outErr is NULL")
	}
	if handle == 0 {
		panic("fatstd_go_xml_encoder_close: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_xml_encoder_close: invalid handle")
	}
	e, ok := value.(*fatXmlEncoder)
	if !ok {
		panic("fatstd_go_xml_encoder_close: handle is not xml encoder")
	}
	if err := e.enc.Close(); err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatXmlErrCodeOther, err.Error()))
		return fatstdXmlStatusFromError(err)
	}
	*outErr = 0
	return fatStatusOK
}

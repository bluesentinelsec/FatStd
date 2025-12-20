package main

/*
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
*/
import "C"

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"sort"
	"strconv"
)

const (
	fatJsonErrCodeSyntax = 170
	fatJsonErrCodeOther  = 171
)

type fatJsonValue struct {
	v any
}

type fatJsonDecoder struct {
	dec *json.Decoder
}

type fatJsonEncoder struct {
	enc *json.Encoder
}

func fatstdJsonStatusFromError(err error) C.int {
	if err == nil {
		return fatStatusOK
	}
	if errors.Is(err, io.EOF) {
		return fatStatusEOF
	}
	var se *json.SyntaxError
	if errors.As(err, &se) {
		return fatStatusSyntax
	}
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return fatStatusSyntax
	}
	return fatStatusOther
}

func fatstdJsonValueFromHandle(handle uintptr) *fatJsonValue {
	if handle == 0 {
		panic("fatstdJsonValueFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdJsonValueFromHandle: invalid handle")
	}
	v, ok := value.(*fatJsonValue)
	if !ok {
		panic("fatstdJsonValueFromHandle: handle is not json value")
	}
	return v
}

func fatstdJsonDecoderFromHandle(handle uintptr) *fatJsonDecoder {
	if handle == 0 {
		panic("fatstdJsonDecoderFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdJsonDecoderFromHandle: invalid handle")
	}
	d, ok := value.(*fatJsonDecoder)
	if !ok {
		panic("fatstdJsonDecoderFromHandle: handle is not json decoder")
	}
	if d.dec == nil {
		panic("fatstdJsonDecoderFromHandle: decoder is nil")
	}
	return d
}

func fatstdJsonEncoderFromHandle(handle uintptr) *fatJsonEncoder {
	if handle == 0 {
		panic("fatstdJsonEncoderFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdJsonEncoderFromHandle: invalid handle")
	}
	e, ok := value.(*fatJsonEncoder)
	if !ok {
		panic("fatstdJsonEncoderFromHandle: handle is not json encoder")
	}
	if e.enc == nil {
		panic("fatstdJsonEncoderFromHandle: encoder is nil")
	}
	return e
}

//export fatstd_go_json_valid
func fatstd_go_json_valid(dataHandle C.uintptr_t) C.int {
	b := fatstdBytesFromHandle(uintptr(dataHandle))
	if json.Valid(b.Value()) {
		return 1
	}
	return 0
}

//export fatstd_go_json_compact
func fatstd_go_json_compact(srcHandle C.uintptr_t, outBytes *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outBytes == nil {
		panic("fatstd_go_json_compact: outBytes is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_json_compact: outErr is NULL")
	}

	src := fatstdBytesFromHandle(uintptr(srcHandle))
	var buf bytes.Buffer
	if err := json.Compact(&buf, src.Value()); err != nil {
		*outBytes = 0
		*outErr = C.uintptr_t(fatstdNewError(fatJsonErrCodeSyntax, err.Error()))
		return fatstdJsonStatusFromError(err)
	}
	*outBytes = C.uintptr_t(fatstdBytesNewFromGoBytes(bytes.Clone(buf.Bytes())))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_json_indent
func fatstd_go_json_indent(srcHandle C.uintptr_t, prefixHandle C.uintptr_t, indentHandle C.uintptr_t, outBytes *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outBytes == nil {
		panic("fatstd_go_json_indent: outBytes is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_json_indent: outErr is NULL")
	}

	src := fatstdBytesFromHandle(uintptr(srcHandle))
	prefix := fatstdStringFromHandle(uintptr(prefixHandle))
	indent := fatstdStringFromHandle(uintptr(indentHandle))

	var buf bytes.Buffer
	if err := json.Indent(&buf, src.Value(), prefix.Value(), indent.Value()); err != nil {
		*outBytes = 0
		*outErr = C.uintptr_t(fatstdNewError(fatJsonErrCodeSyntax, err.Error()))
		return fatstdJsonStatusFromError(err)
	}
	*outBytes = C.uintptr_t(fatstdBytesNewFromGoBytes(bytes.Clone(buf.Bytes())))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_json_html_escape
func fatstd_go_json_html_escape(srcHandle C.uintptr_t) C.uintptr_t {
	src := fatstdBytesFromHandle(uintptr(srcHandle))
	var buf bytes.Buffer
	json.HTMLEscape(&buf, src.Value())
	return C.uintptr_t(fatstdBytesNewFromGoBytes(bytes.Clone(buf.Bytes())))
}

func fatstdJsonDecodeAnyUseNumber(data []byte) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}

//export fatstd_go_json_unmarshal
func fatstd_go_json_unmarshal(dataHandle C.uintptr_t, outValue *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outValue == nil {
		panic("fatstd_go_json_unmarshal: outValue is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_json_unmarshal: outErr is NULL")
	}

	data := fatstdBytesFromHandle(uintptr(dataHandle))
	v, err := fatstdJsonDecodeAnyUseNumber(data.Value())
	if err != nil {
		*outValue = 0
		*outErr = C.uintptr_t(fatstdNewError(fatJsonErrCodeSyntax, err.Error()))
		return fatstdJsonStatusFromError(err)
	}
	*outValue = C.uintptr_t(fatstdHandles.register(&fatJsonValue{v: v}))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_json_marshal
func fatstd_go_json_marshal(valueHandle C.uintptr_t, outBytes *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outBytes == nil {
		panic("fatstd_go_json_marshal: outBytes is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_json_marshal: outErr is NULL")
	}
	v := fatstdJsonValueFromHandle(uintptr(valueHandle))
	b, err := json.Marshal(v.v)
	if err != nil {
		*outBytes = 0
		*outErr = C.uintptr_t(fatstdNewError(fatJsonErrCodeOther, err.Error()))
		return fatstdJsonStatusFromError(err)
	}
	*outBytes = C.uintptr_t(fatstdBytesNewFromGoBytes(b))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_json_marshal_indent
func fatstd_go_json_marshal_indent(valueHandle C.uintptr_t, prefixHandle C.uintptr_t, indentHandle C.uintptr_t, outBytes *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outBytes == nil {
		panic("fatstd_go_json_marshal_indent: outBytes is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_json_marshal_indent: outErr is NULL")
	}
	v := fatstdJsonValueFromHandle(uintptr(valueHandle))
	prefix := fatstdStringFromHandle(uintptr(prefixHandle))
	indent := fatstdStringFromHandle(uintptr(indentHandle))
	b, err := json.MarshalIndent(v.v, prefix.Value(), indent.Value())
	if err != nil {
		*outBytes = 0
		*outErr = C.uintptr_t(fatstdNewError(fatJsonErrCodeOther, err.Error()))
		return fatstdJsonStatusFromError(err)
	}
	*outBytes = C.uintptr_t(fatstdBytesNewFromGoBytes(b))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_json_value_free
func fatstd_go_json_value_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_json_value_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_json_value_free: invalid handle")
	}
	if _, ok := value.(*fatJsonValue); !ok {
		panic("fatstd_go_json_value_free: handle is not json value")
	}
}

//export fatstd_go_json_value_type
func fatstd_go_json_value_type(handle C.uintptr_t) C.int {
	v := fatstdJsonValueFromHandle(uintptr(handle))
	switch v.v.(type) {
	case nil:
		return 0
	case bool:
		return 1
	case json.Number, float64:
		return 2
	case string:
		return 3
	case []any:
		return 4
	case map[string]any:
		return 5
	default:
		panic("fatstd_go_json_value_type: unsupported value type")
	}
}

//export fatstd_go_json_value_as_bool
func fatstd_go_json_value_as_bool(handle C.uintptr_t, outValue *C.int) {
	if outValue == nil {
		panic("fatstd_go_json_value_as_bool: outValue is NULL")
	}
	v := fatstdJsonValueFromHandle(uintptr(handle))
	b, ok := v.v.(bool)
	if !ok {
		panic("fatstd_go_json_value_as_bool: value is not bool")
	}
	*outValue = 0
	if b {
		*outValue = 1
	}
}

//export fatstd_go_json_value_as_string
func fatstd_go_json_value_as_string(handle C.uintptr_t) C.uintptr_t {
	v := fatstdJsonValueFromHandle(uintptr(handle))
	s, ok := v.v.(string)
	if !ok {
		panic("fatstd_go_json_value_as_string: value is not string")
	}
	return C.uintptr_t(fatstdStringNewFromGoString(s))
}

//export fatstd_go_json_value_as_number_string
func fatstd_go_json_value_as_number_string(handle C.uintptr_t) C.uintptr_t {
	v := fatstdJsonValueFromHandle(uintptr(handle))
	switch n := v.v.(type) {
	case json.Number:
		return C.uintptr_t(fatstdStringNewFromGoString(n.String()))
	case float64:
		return C.uintptr_t(fatstdStringNewFromGoString(strconv.FormatFloat(n, 'g', -1, 64)))
	default:
		panic("fatstd_go_json_value_as_number_string: value is not number")
	}
}

func fatstdJsonArray(v any) []any {
	switch a := v.(type) {
	case []any:
		return a
	default:
		panic("fatstdJsonArray: not array")
	}
}

//export fatstd_go_json_array_len
func fatstd_go_json_array_len(handle C.uintptr_t) C.size_t {
	v := fatstdJsonValueFromHandle(uintptr(handle))
	a := fatstdJsonArray(v.v)
	return C.size_t(len(a))
}

//export fatstd_go_json_array_get
func fatstd_go_json_array_get(handle C.uintptr_t, idx C.size_t) C.uintptr_t {
	v := fatstdJsonValueFromHandle(uintptr(handle))
	a := fatstdJsonArray(v.v)
	if idx > C.size_t(2147483647) {
		panic("fatstd_go_json_array_get: idx too large")
	}
	i := int(idx)
	if i < 0 || i >= len(a) {
		panic("fatstd_go_json_array_get: idx out of range")
	}
	return C.uintptr_t(fatstdHandles.register(&fatJsonValue{v: a[i]}))
}

func fatstdJsonObject(v any) map[string]any {
	switch m := v.(type) {
	case map[string]any:
		return m
	default:
		panic("fatstdJsonObject: not object")
	}
}

//export fatstd_go_json_object_keys
func fatstd_go_json_object_keys(handle C.uintptr_t) C.uintptr_t {
	v := fatstdJsonValueFromHandle(uintptr(handle))
	m := fatstdJsonObject(v.v)
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return C.uintptr_t(fatstdStringArrayNew(keys))
}

//export fatstd_go_json_object_get
func fatstd_go_json_object_get(handle C.uintptr_t, keyHandle C.uintptr_t, outFound *C.bool, outValue *C.uintptr_t) {
	if outFound == nil {
		panic("fatstd_go_json_object_get: outFound is NULL")
	}
	if outValue == nil {
		panic("fatstd_go_json_object_get: outValue is NULL")
	}
	v := fatstdJsonValueFromHandle(uintptr(handle))
	key := fatstdStringFromHandle(uintptr(keyHandle))
	m := fatstdJsonObject(v.v)
	val, ok := m[key.Value()]
	if !ok {
		*outFound = false
		*outValue = 0
		return
	}
	*outFound = true
	*outValue = C.uintptr_t(fatstdHandles.register(&fatJsonValue{v: val}))
}

//export fatstd_go_json_decoder_new_bytes_reader
func fatstd_go_json_decoder_new_bytes_reader(readerHandle C.uintptr_t) C.uintptr_t {
	r := fatstdBytesReaderFromHandle(uintptr(readerHandle))
	dec := json.NewDecoder(r)
	return C.uintptr_t(fatstdHandles.register(&fatJsonDecoder{dec: dec}))
}

//export fatstd_go_json_decoder_free
func fatstd_go_json_decoder_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_json_decoder_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_json_decoder_free: invalid handle")
	}
	if _, ok := value.(*fatJsonDecoder); !ok {
		panic("fatstd_go_json_decoder_free: handle is not json decoder")
	}
}

//export fatstd_go_json_decoder_use_number
func fatstd_go_json_decoder_use_number(handle C.uintptr_t) {
	d := fatstdJsonDecoderFromHandle(uintptr(handle))
	d.dec.UseNumber()
}

//export fatstd_go_json_decoder_disallow_unknown_fields
func fatstd_go_json_decoder_disallow_unknown_fields(handle C.uintptr_t) {
	d := fatstdJsonDecoderFromHandle(uintptr(handle))
	d.dec.DisallowUnknownFields()
}

//export fatstd_go_json_decoder_input_offset
func fatstd_go_json_decoder_input_offset(handle C.uintptr_t) C.longlong {
	d := fatstdJsonDecoderFromHandle(uintptr(handle))
	return C.longlong(d.dec.InputOffset())
}

//export fatstd_go_json_decoder_more
func fatstd_go_json_decoder_more(handle C.uintptr_t) C.int {
	d := fatstdJsonDecoderFromHandle(uintptr(handle))
	if d.dec.More() {
		return 1
	}
	return 0
}

//export fatstd_go_json_decoder_buffered_bytes
func fatstd_go_json_decoder_buffered_bytes(handle C.uintptr_t) C.uintptr_t {
	d := fatstdJsonDecoderFromHandle(uintptr(handle))
	b, err := io.ReadAll(d.dec.Buffered())
	if err != nil {
		panic("fatstd_go_json_decoder_buffered_bytes: unexpected error")
	}
	return C.uintptr_t(fatstdBytesNewFromGoBytes(b))
}

//export fatstd_go_json_decoder_decode_value
func fatstd_go_json_decoder_decode_value(handle C.uintptr_t, outValue *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outValue == nil {
		panic("fatstd_go_json_decoder_decode_value: outValue is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_json_decoder_decode_value: outErr is NULL")
	}
	d := fatstdJsonDecoderFromHandle(uintptr(handle))
	var v any
	if err := d.dec.Decode(&v); err != nil {
		*outValue = 0
		if errors.Is(err, io.EOF) {
			*outErr = 0
			return fatStatusEOF
		}
		*outErr = C.uintptr_t(fatstdNewError(fatJsonErrCodeSyntax, err.Error()))
		return fatstdJsonStatusFromError(err)
	}
	*outValue = C.uintptr_t(fatstdHandles.register(&fatJsonValue{v: v}))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_json_encoder_new_to_bytes_buffer
func fatstd_go_json_encoder_new_to_bytes_buffer(dstBufferHandle C.uintptr_t) C.uintptr_t {
	dst := fatstdBytesBufferFromHandle(uintptr(dstBufferHandle))
	enc := json.NewEncoder(dst.Underlying())
	return C.uintptr_t(fatstdHandles.register(&fatJsonEncoder{enc: enc}))
}

//export fatstd_go_json_encoder_free
func fatstd_go_json_encoder_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_json_encoder_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_json_encoder_free: invalid handle")
	}
	if _, ok := value.(*fatJsonEncoder); !ok {
		panic("fatstd_go_json_encoder_free: handle is not json encoder")
	}
}

//export fatstd_go_json_encoder_set_escape_html
func fatstd_go_json_encoder_set_escape_html(handle C.uintptr_t, on C.int) {
	e := fatstdJsonEncoderFromHandle(uintptr(handle))
	e.enc.SetEscapeHTML(on != 0)
}

//export fatstd_go_json_encoder_set_indent
func fatstd_go_json_encoder_set_indent(handle C.uintptr_t, prefixHandle C.uintptr_t, indentHandle C.uintptr_t) {
	e := fatstdJsonEncoderFromHandle(uintptr(handle))
	prefix := fatstdStringFromHandle(uintptr(prefixHandle))
	indent := fatstdStringFromHandle(uintptr(indentHandle))
	e.enc.SetIndent(prefix.Value(), indent.Value())
}

//export fatstd_go_json_encoder_encode_value
func fatstd_go_json_encoder_encode_value(handle C.uintptr_t, valueHandle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_json_encoder_encode_value: outErr is NULL")
	}
	e := fatstdJsonEncoderFromHandle(uintptr(handle))
	v := fatstdJsonValueFromHandle(uintptr(valueHandle))
	if err := e.enc.Encode(v.v); err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatJsonErrCodeOther, err.Error()))
		return fatstdJsonStatusFromError(err)
	}
	*outErr = 0
	return fatStatusOK
}

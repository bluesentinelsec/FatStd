package json

import (
	"bytes"
	stdjson "encoding/json"
	"errors"
	"io"
)

func Compact(dst *bytes.Buffer, src []byte) error {
	return stdjson.Compact(dst, src)
}

func HTMLEscape(dst *bytes.Buffer, src []byte) {
	stdjson.HTMLEscape(dst, src)
}

func Indent(dst *bytes.Buffer, src []byte, prefix, indent string) error {
	return stdjson.Indent(dst, src, prefix, indent)
}

func Marshal(v any) ([]byte, error) {
	return stdjson.Marshal(v)
}

func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return stdjson.MarshalIndent(v, prefix, indent)
}

func Unmarshal(data []byte, v any) error {
	return stdjson.Unmarshal(data, v)
}

func Valid(data []byte) bool {
	return stdjson.Valid(data)
}

type Decoder = stdjson.Decoder

func NewDecoder(r io.Reader) *Decoder {
	return stdjson.NewDecoder(r)
}

type Delim = stdjson.Delim

type Encoder = stdjson.Encoder

func NewEncoder(w io.Writer) *Encoder {
	return stdjson.NewEncoder(w)
}

// Deprecated: kept for parity with encoding/json.
type InvalidUTF8Error = stdjson.InvalidUTF8Error

type InvalidUnmarshalError = stdjson.InvalidUnmarshalError

type Marshaler = stdjson.Marshaler

type MarshalerError = stdjson.MarshalerError

type Number = stdjson.Number

type RawMessage = stdjson.RawMessage

type SyntaxError = stdjson.SyntaxError

type Token = stdjson.Token

// Deprecated: kept for parity with encoding/json.
type UnmarshalFieldError = stdjson.UnmarshalFieldError

type UnmarshalTypeError = stdjson.UnmarshalTypeError

type Unmarshaler = stdjson.Unmarshaler

type UnsupportedTypeError = stdjson.UnsupportedTypeError

type UnsupportedValueError = stdjson.UnsupportedValueError

// Ensure these error strings stay consistent if we ever stop aliasing stdlib.
var _ = errors.New


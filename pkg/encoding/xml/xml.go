package xml

import (
	stdxml "encoding/xml"
	"errors"
	"io"
)

func Escape(w io.Writer, s []byte) {
	stdxml.Escape(w, s)
}

func EscapeText(w io.Writer, s []byte) error {
	return stdxml.EscapeText(w, s)
}

func Marshal(v any) ([]byte, error) {
	return stdxml.Marshal(v)
}

func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return stdxml.MarshalIndent(v, prefix, indent)
}

func Unmarshal(data []byte, v any) error {
	return stdxml.Unmarshal(data, v)
}

type Attr = stdxml.Attr

type CharData = stdxml.CharData

type Comment = stdxml.Comment

type Decoder = stdxml.Decoder

func NewDecoder(r io.Reader) *Decoder {
	return stdxml.NewDecoder(r)
}

func NewTokenDecoder(t TokenReader) *Decoder {
	return stdxml.NewTokenDecoder(t)
}

type Directive = stdxml.Directive

type Encoder = stdxml.Encoder

func NewEncoder(w io.Writer) *Encoder {
	return stdxml.NewEncoder(w)
}

type EndElement = stdxml.EndElement

type Marshaler = stdxml.Marshaler

type MarshalerAttr = stdxml.MarshalerAttr

type Name = stdxml.Name

type ProcInst = stdxml.ProcInst

type StartElement = stdxml.StartElement

type SyntaxError = stdxml.SyntaxError

type TagPathError = stdxml.TagPathError

type Token = stdxml.Token

func CopyToken(t Token) Token {
	return stdxml.CopyToken(t)
}

type TokenReader = stdxml.TokenReader

type UnmarshalError = stdxml.UnmarshalError

type Unmarshaler = stdxml.Unmarshaler

type UnmarshalerAttr = stdxml.UnmarshalerAttr

type UnsupportedTypeError = stdxml.UnsupportedTypeError

// Ensure these error strings stay consistent if we ever stop aliasing stdlib.
var _ = errors.New


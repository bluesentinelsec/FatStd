package csv

import (
	stdcsv "encoding/csv"
	"errors"
	"io"
)

type ParseError = stdcsv.ParseError

type Reader = stdcsv.Reader

type Writer = stdcsv.Writer

var (
	ErrBareQuote  = stdcsv.ErrBareQuote
	ErrQuote      = stdcsv.ErrQuote
	ErrFieldCount = stdcsv.ErrFieldCount

	// Deprecated: ErrTrailingComma is no longer used.
	ErrTrailingComma = stdcsv.ErrTrailingComma
)

func NewReader(r io.Reader) *Reader {
	return stdcsv.NewReader(r)
}

func NewWriter(w io.Writer) *Writer {
	return stdcsv.NewWriter(w)
}

// Ensure these error strings stay consistent if we ever stop aliasing stdlib.
var _ = errors.New


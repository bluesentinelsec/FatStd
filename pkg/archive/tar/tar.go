package tar

import (
	stdtar "archive/tar"
	"errors"
	"io"
	"io/fs"
)

// FileInfoNames is the same type as archive/tar.FileInfoNames.
type FileInfoNames = stdtar.FileInfoNames

// Format is the same type as archive/tar.Format.
type Format = stdtar.Format

// Header is the same type as archive/tar.Header.
type Header = stdtar.Header

// Reader is the same type as archive/tar.Reader.
type Reader = stdtar.Reader

// Writer is the same type as archive/tar.Writer.
type Writer = stdtar.Writer

const (
	// Type '0' indicates a regular file.
	TypeReg = stdtar.TypeReg

	// Deprecated: Use TypeReg instead.
	TypeRegA = stdtar.TypeRegA

	// Type '1' to '6' are header-only flags and may not have a data body.
	TypeLink    = stdtar.TypeLink    // Hard link
	TypeSymlink = stdtar.TypeSymlink // Symbolic link
	TypeChar    = stdtar.TypeChar    // Character device node
	TypeBlock   = stdtar.TypeBlock   // Block device node
	TypeDir     = stdtar.TypeDir     // Directory
	TypeFifo    = stdtar.TypeFifo    // FIFO node

	// Type '7' is reserved.
	TypeCont = stdtar.TypeCont

	// Type 'x' is used by the PAX format to store key-value records that
	// are only relevant to the next file.
	// This package transparently handles these types.
	TypeXHeader = stdtar.TypeXHeader

	// Type 'g' is used by the PAX format to store key-value records that
	// are relevant to all subsequent files.
	// This package only supports parsing and composing such headers,
	// but does not currently support persisting the global state across files.
	TypeXGlobalHeader = stdtar.TypeXGlobalHeader

	// Type 'S' indicates a sparse file in the GNU format.
	TypeGNUSparse = stdtar.TypeGNUSparse

	// Types 'L' and 'K' are used by the GNU format for a meta file
	// used to store the path or link name for the next file.
	// This package transparently handles these types.
	TypeGNULongName = stdtar.TypeGNULongName
	TypeGNULongLink = stdtar.TypeGNULongLink
)

var (
	ErrHeader          = stdtar.ErrHeader
	ErrWriteTooLong    = stdtar.ErrWriteTooLong
	ErrFieldTooLong    = stdtar.ErrFieldTooLong
	ErrWriteAfterClose = stdtar.ErrWriteAfterClose
	ErrInsecurePath    = stdtar.ErrInsecurePath
)

func FileInfoHeader(fi fs.FileInfo, link string) (*Header, error) {
	return stdtar.FileInfoHeader(fi, link)
}

func NewReader(r io.Reader) *Reader {
	return stdtar.NewReader(r)
}

func NewWriter(w io.Writer) *Writer {
	return stdtar.NewWriter(w)
}

// Ensure these error strings stay consistent if we ever stop aliasing stdlib.
var _ = errors.New


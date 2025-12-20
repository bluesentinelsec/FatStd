package zip

import (
	stdzip "archive/zip"
	"io"
	"io/fs"
)

// Compressor is the same type as archive/zip.Compressor.
type Compressor = stdzip.Compressor

// Decompressor is the same type as archive/zip.Decompressor.
type Decompressor = stdzip.Decompressor

// File is the same type as archive/zip.File.
type File = stdzip.File

// FileHeader is the same type as archive/zip.FileHeader.
type FileHeader = stdzip.FileHeader

// ReadCloser is the same type as archive/zip.ReadCloser.
type ReadCloser = stdzip.ReadCloser

// Reader is the same type as archive/zip.Reader.
type Reader = stdzip.Reader

// Writer is the same type as archive/zip.Writer.
type Writer = stdzip.Writer

const (
	Store   uint16 = stdzip.Store   // no compression
	Deflate uint16 = stdzip.Deflate // DEFLATE compressed
)

var (
	ErrFormat       = stdzip.ErrFormat
	ErrAlgorithm    = stdzip.ErrAlgorithm
	ErrChecksum     = stdzip.ErrChecksum
	ErrInsecurePath = stdzip.ErrInsecurePath
)

func RegisterCompressor(method uint16, comp Compressor) {
	stdzip.RegisterCompressor(method, comp)
}

func RegisterDecompressor(method uint16, dcomp Decompressor) {
	stdzip.RegisterDecompressor(method, dcomp)
}

func FileInfoHeader(fi fs.FileInfo) (*FileHeader, error) {
	return stdzip.FileInfoHeader(fi)
}

func OpenReader(name string) (*ReadCloser, error) {
	return stdzip.OpenReader(name)
}

func NewReader(r io.ReaderAt, size int64) (*Reader, error) {
	return stdzip.NewReader(r, size)
}

func NewWriter(w io.Writer) *Writer {
	return stdzip.NewWriter(w)
}

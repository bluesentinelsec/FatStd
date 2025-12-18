package fatbytes

import (
	"bytes"
	"io"
)

type Reader struct {
	reader bytes.Reader
}

func NewReader(b []byte) *Reader {
	return &Reader{reader: *bytes.NewReader(b)}
}

func (r *Reader) Len() int {
	if r == nil {
		panic("fatbytes.Reader.Len: receiver is nil")
	}
	return r.reader.Len()
}

func (r *Reader) Read(p []byte) (int, error) {
	if r == nil {
		panic("fatbytes.Reader.Read: receiver is nil")
	}
	return r.reader.Read(p)
}

func (r *Reader) ReadAt(p []byte, off int64) (int, error) {
	if r == nil {
		panic("fatbytes.Reader.ReadAt: receiver is nil")
	}
	return r.reader.ReadAt(p, off)
}

func (r *Reader) ReadByte() (byte, error) {
	if r == nil {
		panic("fatbytes.Reader.ReadByte: receiver is nil")
	}
	return r.reader.ReadByte()
}

func (r *Reader) Reset(b []byte) {
	if r == nil {
		panic("fatbytes.Reader.Reset: receiver is nil")
	}
	r.reader.Reset(b)
}

func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	if r == nil {
		panic("fatbytes.Reader.Seek: receiver is nil")
	}
	return r.reader.Seek(offset, whence)
}

func (r *Reader) Size() int64 {
	if r == nil {
		panic("fatbytes.Reader.Size: receiver is nil")
	}
	return r.reader.Size()
}

func (r *Reader) UnreadByte() error {
	if r == nil {
		panic("fatbytes.Reader.UnreadByte: receiver is nil")
	}
	return r.reader.UnreadByte()
}

func (r *Reader) WriteTo(w io.Writer) (int64, error) {
	if r == nil {
		panic("fatbytes.Reader.WriteTo: receiver is nil")
	}
	return r.reader.WriteTo(w)
}


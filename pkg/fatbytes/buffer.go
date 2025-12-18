package fatbytes

import (
	"bytes"
	"io"
)

type Buffer struct {
	buf bytes.Buffer
}

func NewBuffer(init []byte) *Buffer {
	return &Buffer{buf: *bytes.NewBuffer(init)}
}

func NewBufferString(s string) *Buffer {
	return &Buffer{buf: *bytes.NewBufferString(s)}
}

func (b *Buffer) Underlying() *bytes.Buffer {
	if b == nil {
		panic("fatbytes.Buffer.Underlying: receiver is nil")
	}
	return &b.buf
}

func (b *Buffer) Bytes() []byte {
	if b == nil {
		panic("fatbytes.Buffer.Bytes: receiver is nil")
	}
	return b.buf.Bytes()
}

func (b *Buffer) String() string {
	if b == nil {
		panic("fatbytes.Buffer.String: receiver is nil")
	}
	return b.buf.String()
}

func (b *Buffer) Len() int {
	if b == nil {
		panic("fatbytes.Buffer.Len: receiver is nil")
	}
	return b.buf.Len()
}

func (b *Buffer) Cap() int {
	if b == nil {
		panic("fatbytes.Buffer.Cap: receiver is nil")
	}
	return b.buf.Cap()
}

func (b *Buffer) Grow(n int) {
	if b == nil {
		panic("fatbytes.Buffer.Grow: receiver is nil")
	}
	b.buf.Grow(n)
}

func (b *Buffer) Reset() {
	if b == nil {
		panic("fatbytes.Buffer.Reset: receiver is nil")
	}
	b.buf.Reset()
}

func (b *Buffer) Truncate(n int) {
	if b == nil {
		panic("fatbytes.Buffer.Truncate: receiver is nil")
	}
	b.buf.Truncate(n)
}

func (b *Buffer) Write(p []byte) int {
	if b == nil {
		panic("fatbytes.Buffer.Write: receiver is nil")
	}
	n, _ := b.buf.Write(p)
	return n
}

func (b *Buffer) WriteByte(c byte) {
	if b == nil {
		panic("fatbytes.Buffer.WriteByte: receiver is nil")
	}
	_ = b.buf.WriteByte(c)
}

func (b *Buffer) WriteRune(r rune) int {
	if b == nil {
		panic("fatbytes.Buffer.WriteRune: receiver is nil")
	}
	n, _ := b.buf.WriteRune(r)
	return n
}

func (b *Buffer) WriteString(s string) int {
	if b == nil {
		panic("fatbytes.Buffer.WriteString: receiver is nil")
	}
	n, _ := b.buf.WriteString(s)
	return n
}

func (b *Buffer) Read(p []byte) (int, error) {
	if b == nil {
		panic("fatbytes.Buffer.Read: receiver is nil")
	}
	return b.buf.Read(p)
}

func (b *Buffer) Next(n int) []byte {
	if b == nil {
		panic("fatbytes.Buffer.Next: receiver is nil")
	}
	return b.buf.Next(n)
}

func (b *Buffer) ReadByte() (byte, error) {
	if b == nil {
		panic("fatbytes.Buffer.ReadByte: receiver is nil")
	}
	return b.buf.ReadByte()
}

func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	if b == nil {
		panic("fatbytes.Buffer.WriteTo: receiver is nil")
	}
	return b.buf.WriteTo(w)
}

func (b *Buffer) ReadFrom(r io.Reader) (int64, error) {
	if b == nil {
		panic("fatbytes.Buffer.ReadFrom: receiver is nil")
	}
	return b.buf.ReadFrom(r)
}


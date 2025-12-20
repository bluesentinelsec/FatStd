package bzip2

import (
	"bytes"
	"compress/bzip2"
	"io"
)

// Decompress inflates a bzip2-compressed buffer.
func Decompress(data []byte) ([]byte, error) {
	r := bzip2.NewReader(bytes.NewReader(data))
	return io.ReadAll(r)
}

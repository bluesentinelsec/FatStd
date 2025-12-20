package lzw

import (
	"bytes"
	"compress/lzw"
	"io"
)

// Compress encodes data using LZW with the specified order and literal width.
func Compress(data []byte, order lzw.Order, litWidth int) ([]byte, error) {
	var buf bytes.Buffer
	w := lzw.NewWriter(&buf, order, litWidth)
	if _, err := w.Write(data); err != nil {
		_ = w.Close()
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decompress decodes LZW data using the specified order and literal width.
func Decompress(data []byte, order lzw.Order, litWidth int) ([]byte, error) {
	r := lzw.NewReader(bytes.NewReader(data), order, litWidth)
	defer r.Close()
	return io.ReadAll(r)
}

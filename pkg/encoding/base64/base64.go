package base64

import (
	stdb64 "encoding/base64"
	"io"
)

// CorruptInputError matches encoding/base64.CorruptInputError.
type CorruptInputError = stdb64.CorruptInputError

// Encoding is a wrapper over encoding/base64.Encoding that also provides
// AppendEncode and AppendDecode for Go versions that lack those methods.
type Encoding struct {
	enc *stdb64.Encoding
}

func NewDecoder(enc *Encoding, r io.Reader) io.Reader {
	if enc == nil || enc.enc == nil {
		panic("base64.NewDecoder: enc is nil")
	}
	return stdb64.NewDecoder(enc.enc, r)
}

func NewEncoder(enc *Encoding, w io.Writer) io.WriteCloser {
	if enc == nil || enc.enc == nil {
		panic("base64.NewEncoder: enc is nil")
	}
	return stdb64.NewEncoder(enc.enc, w)
}

func NewEncoding(encoder string) *Encoding {
	return &Encoding{enc: stdb64.NewEncoding(encoder)}
}

func (enc *Encoding) AppendDecode(dst, src []byte) ([]byte, error) {
	if enc == nil || enc.enc == nil {
		panic("base64.Encoding.AppendDecode: receiver is nil")
	}
	buf := make([]byte, enc.enc.DecodedLen(len(src)))
	n, err := enc.enc.Decode(buf, src)
	if err != nil {
		return nil, err
	}
	return append(dst, buf[:n]...), nil
}

func (enc *Encoding) AppendEncode(dst, src []byte) []byte {
	if enc == nil || enc.enc == nil {
		panic("base64.Encoding.AppendEncode: receiver is nil")
	}
	need := enc.enc.EncodedLen(len(src))
	out := append(dst, make([]byte, need)...)
	enc.enc.Encode(out[len(dst):], src)
	return out
}

func (enc *Encoding) Decode(dst, src []byte) (n int, err error) {
	if enc == nil || enc.enc == nil {
		panic("base64.Encoding.Decode: receiver is nil")
	}
	return enc.enc.Decode(dst, src)
}

func (enc *Encoding) DecodeString(s string) ([]byte, error) {
	if enc == nil || enc.enc == nil {
		panic("base64.Encoding.DecodeString: receiver is nil")
	}
	return enc.enc.DecodeString(s)
}

func (enc *Encoding) DecodedLen(n int) int {
	if enc == nil || enc.enc == nil {
		panic("base64.Encoding.DecodedLen: receiver is nil")
	}
	return enc.enc.DecodedLen(n)
}

func (enc *Encoding) Encode(dst, src []byte) {
	if enc == nil || enc.enc == nil {
		panic("base64.Encoding.Encode: receiver is nil")
	}
	enc.enc.Encode(dst, src)
}

func (enc *Encoding) EncodeToString(src []byte) string {
	if enc == nil || enc.enc == nil {
		panic("base64.Encoding.EncodeToString: receiver is nil")
	}
	return enc.enc.EncodeToString(src)
}

func (enc *Encoding) EncodedLen(n int) int {
	if enc == nil || enc.enc == nil {
		panic("base64.Encoding.EncodedLen: receiver is nil")
	}
	return enc.enc.EncodedLen(n)
}

func (enc Encoding) Strict() *Encoding {
	if enc.enc == nil {
		panic("base64.Encoding.Strict: receiver is nil")
	}
	return &Encoding{enc: enc.enc.Strict()}
}

func (enc Encoding) WithPadding(padding rune) *Encoding {
	if enc.enc == nil {
		panic("base64.Encoding.WithPadding: receiver is nil")
	}
	return &Encoding{enc: enc.enc.WithPadding(padding)}
}


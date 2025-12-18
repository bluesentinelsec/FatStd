package fatbytes

import "bytes"

type Bytes struct {
	value []byte
}

type BytesArray struct {
	values [][]byte
}

func New(value []byte) *Bytes {
	return &Bytes{value: value}
}

func NewArray(values [][]byte) *BytesArray {
	return &BytesArray{values: values}
}

func (b *Bytes) Value() []byte {
	if b == nil {
		panic("fatbytes.Bytes.Value: receiver is nil")
	}
	return b.value
}

func (a *BytesArray) Len() int {
	if a == nil {
		panic("fatbytes.BytesArray.Len: receiver is nil")
	}
	return len(a.values)
}

func (a *BytesArray) Get(index int) []byte {
	if a == nil {
		panic("fatbytes.BytesArray.Get: receiver is nil")
	}
	if index < 0 || index >= len(a.values) {
		panic("fatbytes.BytesArray.Get: index out of range")
	}
	return a.values[index]
}

func (a *BytesArray) Values() [][]byte {
	if a == nil {
		panic("fatbytes.BytesArray.Values: receiver is nil")
	}
	return a.values
}

func Clone(b []byte) []byte {
	return bytes.Clone(b)
}

func Contains(b, subslice []byte) bool {
	return bytes.Contains(b, subslice)
}

func HasPrefix(s, prefix []byte) bool {
	return bytes.HasPrefix(s, prefix)
}

func HasSuffix(s, suffix []byte) bool {
	return bytes.HasSuffix(s, suffix)
}

func TrimSpace(s []byte) []byte {
	return bytes.TrimSpace(s)
}

func Trim(s []byte, cutset string) []byte {
	return bytes.Trim(s, cutset)
}

func TrimPrefix(s, prefix []byte) []byte {
	return bytes.TrimPrefix(s, prefix)
}

func TrimSuffix(s, suffix []byte) []byte {
	return bytes.TrimSuffix(s, suffix)
}

func Cut(s, sep []byte) (before, after []byte, found bool) {
	return bytes.Cut(s, sep)
}

func CutPrefix(s, prefix []byte) (after []byte, found bool) {
	return bytes.CutPrefix(s, prefix)
}

func CutSuffix(s, suffix []byte) (after []byte, found bool) {
	return bytes.CutSuffix(s, suffix)
}

func Split(s, sep []byte) [][]byte {
	return bytes.Split(s, sep)
}

func Fields(s []byte) [][]byte {
	return bytes.Fields(s)
}

func Join(s [][]byte, sep []byte) []byte {
	return bytes.Join(s, sep)
}

func ReplaceAll(s, old, newValue []byte) []byte {
	return bytes.ReplaceAll(s, old, newValue)
}

func Replace(s, old, newValue []byte, n int) []byte {
	return bytes.Replace(s, old, newValue, n)
}

func Repeat(b []byte, count int) []byte {
	return bytes.Repeat(b, count)
}

func ToLower(s []byte) []byte {
	return bytes.ToLower(s)
}

func ToUpper(s []byte) []byte {
	return bytes.ToUpper(s)
}

func IndexByte(b []byte, c byte) int {
	return bytes.IndexByte(b, c)
}

func IndexAny(s []byte, chars string) int {
	return bytes.IndexAny(s, chars)
}

func ToValidUTF8(s, replacement []byte) []byte {
	return bytes.ToValidUTF8(s, replacement)
}

func Index(s, sep []byte) int {
	return bytes.Index(s, sep)
}

func Count(s, sep []byte) int {
	return bytes.Count(s, sep)
}

func Compare(a, b []byte) int {
	return bytes.Compare(a, b)
}

func Equal(a, b []byte) bool {
	return bytes.Equal(a, b)
}

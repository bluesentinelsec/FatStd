package fatstrings

import (
	"io"
	"strings"
)

type String struct {
	value string
}

type StringArray struct {
	values []string
}

type Builder struct {
	builder strings.Builder
}

type Reader struct {
	reader strings.Reader
}

func NewUTF8(value string) *String {
	return &String{value: value}
}

func NewStringArray(values []string) *StringArray {
	return &StringArray{values: values}
}

func NewBuilder() *Builder {
	return &Builder{}
}

func NewReader(s string) *Reader {
	return &Reader{reader: *strings.NewReader(s)}
}

func (s *String) Value() string {
	if s == nil {
		panic("fatstrings.String.Value: receiver is nil")
	}
	return s.value
}

func (a *StringArray) Len() int {
	if a == nil {
		panic("fatstrings.StringArray.Len: receiver is nil")
	}
	return len(a.values)
}

func (a *StringArray) Get(index int) string {
	if a == nil {
		panic("fatstrings.StringArray.Get: receiver is nil")
	}
	if index < 0 || index >= len(a.values) {
		panic("fatstrings.StringArray.Get: index out of range")
	}
	return a.values[index]
}

func (a *StringArray) Values() []string {
	if a == nil {
		panic("fatstrings.StringArray.Values: receiver is nil")
	}
	return a.values
}

func (b *Builder) Cap() int {
	if b == nil {
		panic("fatstrings.Builder.Cap: receiver is nil")
	}
	return b.builder.Cap()
}

func (b *Builder) Grow(n int) {
	if b == nil {
		panic("fatstrings.Builder.Grow: receiver is nil")
	}
	b.builder.Grow(n)
}

func (b *Builder) Len() int {
	if b == nil {
		panic("fatstrings.Builder.Len: receiver is nil")
	}
	return b.builder.Len()
}

func (b *Builder) Reset() {
	if b == nil {
		panic("fatstrings.Builder.Reset: receiver is nil")
	}
	b.builder.Reset()
}

func (b *Builder) String() string {
	if b == nil {
		panic("fatstrings.Builder.String: receiver is nil")
	}
	return b.builder.String()
}

func (b *Builder) Write(p []byte) int {
	if b == nil {
		panic("fatstrings.Builder.Write: receiver is nil")
	}
	n, _ := b.builder.Write(p)
	return n
}

func (b *Builder) WriteByte(c byte) {
	if b == nil {
		panic("fatstrings.Builder.WriteByte: receiver is nil")
	}
	_ = b.builder.WriteByte(c)
}

func (b *Builder) WriteString(s string) int {
	if b == nil {
		panic("fatstrings.Builder.WriteString: receiver is nil")
	}
	n, _ := b.builder.WriteString(s)
	return n
}

func (b *Builder) Underlying() *strings.Builder {
	if b == nil {
		panic("fatstrings.Builder.Underlying: receiver is nil")
	}
	return &b.builder
}

func (r *Reader) Len() int {
	if r == nil {
		panic("fatstrings.Reader.Len: receiver is nil")
	}
	return r.reader.Len()
}

func (r *Reader) Read(p []byte) (int, error) {
	if r == nil {
		panic("fatstrings.Reader.Read: receiver is nil")
	}
	return r.reader.Read(p)
}

func (r *Reader) ReadAt(p []byte, off int64) (int, error) {
	if r == nil {
		panic("fatstrings.Reader.ReadAt: receiver is nil")
	}
	return r.reader.ReadAt(p, off)
}

func (r *Reader) ReadByte() (byte, error) {
	if r == nil {
		panic("fatstrings.Reader.ReadByte: receiver is nil")
	}
	return r.reader.ReadByte()
}

func (r *Reader) Reset(s string) {
	if r == nil {
		panic("fatstrings.Reader.Reset: receiver is nil")
	}
	r.reader.Reset(s)
}

func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	if r == nil {
		panic("fatstrings.Reader.Seek: receiver is nil")
	}
	return r.reader.Seek(offset, whence)
}

func (r *Reader) Size() int64 {
	if r == nil {
		panic("fatstrings.Reader.Size: receiver is nil")
	}
	return r.reader.Size()
}

func (r *Reader) UnreadByte() error {
	if r == nil {
		panic("fatstrings.Reader.UnreadByte: receiver is nil")
	}
	return r.reader.UnreadByte()
}

func (r *Reader) WriteTo(w io.Writer) (int64, error) {
	if r == nil {
		panic("fatstrings.Reader.WriteTo: receiver is nil")
	}
	return r.reader.WriteTo(w)
}

func Clone(s string) string {
	return strings.Clone(s)
}

func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

func Trim(s, cutset string) string {
	return strings.Trim(s, cutset)
}

func Split(s, sep string) []string {
	return strings.Split(s, sep)
}

func SplitN(s, sep string, n int) []string {
	return strings.SplitN(s, sep, n)
}

func Join(elems []string, sep string) string {
	return strings.Join(elems, sep)
}

func Replace(s, old, new string, n int) string {
	return strings.Replace(s, old, new, n)
}

func ReplaceAll(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func Index(s, substr string) int {
	return strings.Index(s, substr)
}

func Count(s, substr string) int {
	return strings.Count(s, substr)
}

func Compare(a, b string) int {
	return strings.Compare(a, b)
}

func EqualFold(s, t string) bool {
	return strings.EqualFold(s, t)
}

func TrimPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}

func TrimSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix)
}

func Cut(s, sep string) (before, after string, found bool) {
	return strings.Cut(s, sep)
}

func CutPrefix(s, prefix string) (after string, found bool) {
	return strings.CutPrefix(s, prefix)
}

func CutSuffix(s, suffix string) (after string, found bool) {
	return strings.CutSuffix(s, suffix)
}

func Fields(s string) []string {
	return strings.Fields(s)
}

func Repeat(s string, count int) string {
	return strings.Repeat(s, count)
}

func ContainsAny(s, chars string) bool {
	return strings.ContainsAny(s, chars)
}

func IndexAny(s, chars string) bool {
	return strings.IndexAny(s, chars) >= 0
}

func ToValidUTF8(s, replacement string) string {
	return strings.ToValidUTF8(s, replacement)
}

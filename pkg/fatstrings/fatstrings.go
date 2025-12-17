package fatstrings

import "strings"

type String struct {
	value string
}

type StringArray struct {
	values []string
}

func NewUTF8(value string) *String {
	return &String{value: value}
}

func NewStringArray(values []string) *StringArray {
	return &StringArray{values: values}
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

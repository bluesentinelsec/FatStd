package fatstrings

import "strings"

type String struct {
	value string
}

func NewUTF8(value string) *String {
	return &String{value: value}
}

func (s *String) Value() string {
	if s == nil {
		panic("fatstrings.String.Value: receiver is nil")
	}
	return s.value
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

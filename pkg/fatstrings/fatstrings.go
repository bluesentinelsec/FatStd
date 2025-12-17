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

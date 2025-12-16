package fatstrings

type String struct {
	value string
}

func NewUTF8(value string) *String {
	return &String{value: value}
}


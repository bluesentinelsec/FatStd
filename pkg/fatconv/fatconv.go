package fatconv

import (
	"errors"
	"strconv"
	"unicode"
)

var ErrRange = strconv.ErrRange
var ErrSyntax = strconv.ErrSyntax

const IntSize = strconv.IntSize

func AppendBool(dst []byte, b bool) []byte { return strconv.AppendBool(dst, b) }
func AppendFloat(dst []byte, f float64, fmt byte, prec, bitSize int) []byte {
	return strconv.AppendFloat(dst, f, fmt, prec, bitSize)
}
func AppendInt(dst []byte, i int64, base int) []byte { return strconv.AppendInt(dst, i, base) }
func AppendUint(dst []byte, i uint64, base int) []byte { return strconv.AppendUint(dst, i, base) }

func AppendQuote(dst []byte, s string) []byte { return strconv.AppendQuote(dst, s) }
func AppendQuoteRune(dst []byte, r rune) []byte { return strconv.AppendQuoteRune(dst, r) }
func AppendQuoteRuneToASCII(dst []byte, r rune) []byte { return strconv.AppendQuoteRuneToASCII(dst, r) }
func AppendQuoteRuneToGraphic(dst []byte, r rune) []byte { return strconv.AppendQuoteRuneToGraphic(dst, r) }
func AppendQuoteToASCII(dst []byte, s string) []byte { return strconv.AppendQuoteToASCII(dst, s) }
func AppendQuoteToGraphic(dst []byte, s string) []byte { return strconv.AppendQuoteToGraphic(dst, s) }

func CanBackquote(s string) bool { return strconv.CanBackquote(s) }

func FormatBool(b bool) string { return strconv.FormatBool(b) }
func FormatComplex(c complex128, fmt byte, prec, bitSize int) string {
	return strconv.FormatComplex(c, fmt, prec, bitSize)
}
func FormatFloat(f float64, fmt byte, prec, bitSize int) string {
	return strconv.FormatFloat(f, fmt, prec, bitSize)
}
func FormatInt(i int64, base int) string { return strconv.FormatInt(i, base) }
func FormatUint(i uint64, base int) string { return strconv.FormatUint(i, base) }

func IsGraphic(r rune) bool { return unicode.IsGraphic(r) }
func IsPrint(r rune) bool { return strconv.IsPrint(r) }

func Itoa(i int) string { return strconv.Itoa(i) }

func Quote(s string) string { return strconv.Quote(s) }
func QuoteRune(r rune) string { return strconv.QuoteRune(r) }
func QuoteRuneToASCII(r rune) string { return strconv.QuoteRuneToASCII(r) }
func QuoteRuneToGraphic(r rune) string { return strconv.QuoteRuneToGraphic(r) }
func QuoteToASCII(s string) string { return strconv.QuoteToASCII(s) }
func QuoteToGraphic(s string) string { return strconv.QuoteToGraphic(s) }

func Unquote(s string) (string, error) { return strconv.Unquote(s) }
func QuotedPrefix(s string) (string, error) { return strconv.QuotedPrefix(s) }

func UnquoteChar(s string, quote byte) (value rune, multibyte bool, tail string, err error) {
	return strconv.UnquoteChar(s, quote)
}

func Atoi(s string) (int, error) { return strconv.Atoi(s) }
func ParseBool(str string) (bool, error) { return strconv.ParseBool(str) }
func ParseFloat(s string, bitSize int) (float64, error) { return strconv.ParseFloat(s, bitSize) }
func ParseInt(s string, base int, bitSize int) (int64, error) { return strconv.ParseInt(s, base, bitSize) }
func ParseUint(s string, base int, bitSize int) (uint64, error) { return strconv.ParseUint(s, base, bitSize) }
func ParseComplex(s string, bitSize int) (complex128, error) { return strconv.ParseComplex(s, bitSize) }

func ClassifyParseError(err error) int32 {
	if err == nil {
		return 0
	}
	if errors.Is(err, ErrSyntax) {
		return 1
	}
	if errors.Is(err, ErrRange) {
		return 2
	}
	return 100
}

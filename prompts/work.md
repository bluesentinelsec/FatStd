Refer to the following documents for context on the project:
1. docs/design.md
2. docs/onboarding_and_testing_functions.md
3. docs/documenting_code.md
4. docs/error_strategy.md

Implement the following functions in fatstd:

encoding/csv

type ParseError
func (e *ParseError) Error() string
func (e *ParseError) Unwrap() error
type Reader
func NewReader(r io.Reader) *Reader
func (r *Reader) FieldPos(field int) (line, column int)
func (r *Reader) InputOffset() int64
func (r *Reader) Read() (record []string, err error)
func (r *Reader) ReadAll() (records [][]string, err error)
type Writer
func NewWriter(w io.Writer) *Writer
func (w *Writer) Error() error
func (w *Writer) Flush()
func (w *Writer) Write(record []string) error
func (w *Writer) WriteAll(records [][]string) error

var (
	ErrBareQuote  = errors.New("bare \" in non-quoted-field")
	ErrQuote      = errors.New("extraneous or missing \" in quoted-field")
	ErrFieldCount = errors.New("wrong number of fields")

	// Deprecated: ErrTrailingComma is no longer used.
	ErrTrailingComma = errors.New("extra delimiter at end of line")
)

type ParseError struct {
	StartLine int   // Line where the record starts
	Line      int   // Line where the error occurred
	Column    int   // Column (1-based byte index) where the error occurred
	Err       error // The actual error
}

I expect the Go bindings, C bindings, and unit tests in Python.
If any of the functions are a poor fit for C, use an alternative that honors the design.

When finished, add a brief tutorial doc showing how to use this module from the perspective of the caller under docs/
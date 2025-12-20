Refer to the following documents for context on the project:
1. docs/design.md
2. docs/onboarding_and_testing_functions.md
3. docs/documenting_code.md
4. docs/error_strategy.md

Implement the following functions in fatstd:

encoding/json


func Compact(dst *bytes.Buffer, src []byte) error
func HTMLEscape(dst *bytes.Buffer, src []byte)
func Indent(dst *bytes.Buffer, src []byte, prefix, indent string) error
func Marshal(v any) ([]byte, error)
func MarshalIndent(v any, prefix, indent string) ([]byte, error)
func Unmarshal(data []byte, v any) error
func Valid(data []byte) bool
type Decoder
func NewDecoder(r io.Reader) *Decoder
func (dec *Decoder) Buffered() io.Reader
func (dec *Decoder) Decode(v any) error
func (dec *Decoder) DisallowUnknownFields()
func (dec *Decoder) InputOffset() int64
func (dec *Decoder) More() bool
func (dec *Decoder) Token() (Token, error)
func (dec *Decoder) UseNumber()
type Delim
func (d Delim) String() string
type Encoder
func NewEncoder(w io.Writer) *Encoder
func (enc *Encoder) Encode(v any) error
func (enc *Encoder) SetEscapeHTML(on bool)
func (enc *Encoder) SetIndent(prefix, indent string)
type InvalidUTF8ErrorDEPRECATED
func (e *InvalidUTF8Error) Error() string
type InvalidUnmarshalError
func (e *InvalidUnmarshalError) Error() string
type Marshaler
type MarshalerError
func (e *MarshalerError) Error() string
func (e *MarshalerError) Unwrap() error
type Number
func (n Number) Float64() (float64, error)
func (n Number) Int64() (int64, error)
func (n Number) String() string
type RawMessage
func (m RawMessage) MarshalJSON() ([]byte, error)
func (m *RawMessage) UnmarshalJSON(data []byte) error
type SyntaxError
func (e *SyntaxError) Error() string
type Token
type UnmarshalFieldErrorDEPRECATED
func (e *UnmarshalFieldError) Error() string
type UnmarshalTypeError
func (e *UnmarshalTypeError) Error() string
type Unmarshaler
type UnsupportedTypeError
func (e *UnsupportedTypeError) Error() string
type UnsupportedValueError
func (e *UnsupportedValueError) Error() string

I expect the Go bindings, C bindings, and unit tests in Python.
If any of the functions are a poor fit for C, use an alternative that honors the design.

When finished, add a brief tutorial doc showing how to use this module from the perspective of the caller under docs/
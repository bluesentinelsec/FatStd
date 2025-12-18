Refer to the following documents for context on the project:
1. docs/design.md
2. docs/onboarding_and_testing_functions.md

Let's continue working on the fat strings subsystem.
I want to onboard these functions:

type Reader
func NewReader(s string) *Reader
func (r *Reader) Len() int
func (r *Reader) Read(b []byte) (n int, err error)
func (r *Reader) ReadAt(b []byte, off int64) (n int, err error)
func (r *Reader) ReadByte() (byte, error)
func (r *Reader) Reset(s string)
func (r *Reader) Seek(offset int64, whence int) (int64, error)
func (r *Reader) Size() int64
func (r *Reader) UnreadByte() error
func (r *Reader) WriteTo(w io.Writer) (n int64, err error)

I expect the Go bindings, C bindings, and a unit test in Python.
Warn me if any of the specified functions violate the design or are impractical for C.
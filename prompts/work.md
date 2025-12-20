Refer to the following documents for context on the project:
1. docs/design.md
2. docs/onboarding_and_testing_functions.md
3. docs/documenting_code.md
4. docs/error_strategy.md

Implement the following archive/tar functions in fatstd.

type FileInfoNames
type Format
func (f Format) String() string
type Header
func FileInfoHeader(fi fs.FileInfo, link string) (*Header, error)
func (h *Header) FileInfo() fs.FileInfo
type Reader
func NewReader(r io.Reader) *Reader
func (tr *Reader) Next() (*Header, error)
func (tr *Reader) Read(b []byte) (int, error)
type Writer
func NewWriter(w io.Writer) *Writer
func (tw *Writer) AddFS(fsys fs.FS) error
func (tw *Writer) Close() error
func (tw *Writer) Flush() error
func (tw *Writer) Write(b []byte) (int, error)
func (tw *Writer) WriteHeader(hdr *Header) error

const (
	// Type '0' indicates a regular file.
	TypeReg = '0'

	// Deprecated: Use TypeReg instead.
	TypeRegA = '\x00'

	// Type '1' to '6' are header-only flags and may not have a data body.
	TypeLink    = '1' // Hard link
	TypeSymlink = '2' // Symbolic link
	TypeChar    = '3' // Character device node
	TypeBlock   = '4' // Block device node
	TypeDir     = '5' // Directory
	TypeFifo    = '6' // FIFO node

	// Type '7' is reserved.
	TypeCont = '7'

	// Type 'x' is used by the PAX format to store key-value records that
	// are only relevant to the next file.
	// This package transparently handles these types.
	TypeXHeader = 'x'

	// Type 'g' is used by the PAX format to store key-value records that
	// are relevant to all subsequent files.
	// This package only supports parsing and composing such headers,
	// but does not currently support persisting the global state across files.
	TypeXGlobalHeader = 'g'

	// Type 'S' indicates a sparse file in the GNU format.
	TypeGNUSparse = 'S'

	// Types 'L' and 'K' are used by the GNU format for a meta file
	// used to store the path or link name for the next file.
	// This package transparently handles these types.
	TypeGNULongName = 'L'
	TypeGNULongLink = 'K'
)

var (
	ErrHeader          = errors.New("archive/tar: invalid tar header")
	ErrWriteTooLong    = errors.New("archive/tar: write too long")
	ErrFieldTooLong    = errors.New("archive/tar: header field too long")
	ErrWriteAfterClose = errors.New("archive/tar: write after close")
	ErrInsecurePath    = errors.New("archive/tar: insecure file path")
)


I expect the Go bindings, C bindings, and unit tests in Python.
If any of the functions are a poor fit for C, use an alternative that honors the design.
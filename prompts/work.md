Refer to the following documents for context on the project:
1. docs/design.md
2. docs/onboarding_and_testing_functions.md
3. docs/documenting_code.md
4. docs/error_strategy.md

Implement the following archive/zip functions in fatstd.

func RegisterCompressor(method uint16, comp Compressor)
func RegisterDecompressor(method uint16, dcomp Decompressor)
type Compressor
type Decompressor
type File
func (f *File) DataOffset() (offset int64, err error)
func (f *File) Open() (io.ReadCloser, error)
func (f *File) OpenRaw() (io.Reader, error)
type FileHeader
func FileInfoHeader(fi fs.FileInfo) (*FileHeader, error)
func (h *FileHeader) FileInfo() fs.FileInfo
func (h *FileHeader) ModTime() time.TimeDEPRECATED
func (h *FileHeader) Mode() (mode fs.FileMode)
func (h *FileHeader) SetModTime(t time.Time)DEPRECATED
func (h *FileHeader) SetMode(mode fs.FileMode)
type ReadCloser
func OpenReader(name string) (*ReadCloser, error)
func (rc *ReadCloser) Close() error
type Reader
func NewReader(r io.ReaderAt, size int64) (*Reader, error)
func (r *Reader) Open(name string) (fs.File, error)
func (r *Reader) RegisterDecompressor(method uint16, dcomp Decompressor)
type Writer
func NewWriter(w io.Writer) *Writer
func (w *Writer) AddFS(fsys fs.FS) error
func (w *Writer) Close() error
func (w *Writer) Copy(f *File) error
func (w *Writer) Create(name string) (io.Writer, error)
func (w *Writer) CreateHeader(fh *FileHeader) (io.Writer, error)
func (w *Writer) CreateRaw(fh *FileHeader) (io.Writer, error)
func (w *Writer) Flush() error
func (w *Writer) RegisterCompressor(method uint16, comp Compressor)
func (w *Writer) SetComment(comment string) error
func (w *Writer) SetOffset(n int64)

const (
	Store   uint16 = 0 // no compression
	Deflate uint16 = 8 // DEFLATE compressed
)

var (
	ErrFormat       = errors.New("zip: not a valid zip file")
	ErrAlgorithm    = errors.New("zip: unsupported compression algorithm")
	ErrChecksum     = errors.New("zip: checksum error")
	ErrInsecurePath = errors.New("zip: insecure file path")
)


I expect the Go bindings, C bindings, and unit tests in Python.
If any of the functions are a poor fit for C, use an alternative that honors the design.
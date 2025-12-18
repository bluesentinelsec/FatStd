package main

/*
#include <stdint.h>
#include <stdbool.h>
*/
import "C"

import (
	"io"
	"unsafe"

	"github.com/bluesentinelsec/FatStd/pkg/fatstrings"
)

var fatstdHandles = newHandleRegistry()

func fatstdStringNewFromGoString(value string) uintptr {
	return fatstdHandles.register(fatstrings.NewUTF8(value))
}

func fatstdStringArrayNew(values []string) uintptr {
	return fatstdHandles.register(fatstrings.NewStringArray(values))
}

func fatstdStringBuilderNew() uintptr {
	return fatstdHandles.register(fatstrings.NewBuilder())
}

func fatstdStringReaderNew(s string) uintptr {
	return fatstdHandles.register(fatstrings.NewReader(s))
}

func fatstdStringFromHandle(handle uintptr) *fatstrings.String {
	if handle == 0 {
		panic("fatstdStringFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdStringFromHandle: invalid handle")
	}
	s, ok := value.(*fatstrings.String)
	if !ok {
		panic("fatstdStringFromHandle: handle is not a fat string")
	}
	return s
}

func fatstdStringArrayFromHandle(handle uintptr) *fatstrings.StringArray {
	if handle == 0 {
		panic("fatstdStringArrayFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdStringArrayFromHandle: invalid handle")
	}
	a, ok := value.(*fatstrings.StringArray)
	if !ok {
		panic("fatstdStringArrayFromHandle: handle is not a fat string array")
	}
	return a
}

func fatstdStringBuilderFromHandle(handle uintptr) *fatstrings.Builder {
	if handle == 0 {
		panic("fatstdStringBuilderFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdStringBuilderFromHandle: invalid handle")
	}
	b, ok := value.(*fatstrings.Builder)
	if !ok {
		panic("fatstdStringBuilderFromHandle: handle is not a fat string builder")
	}
	return b
}

func fatstdStringReaderFromHandle(handle uintptr) *fatstrings.Reader {
	if handle == 0 {
		panic("fatstdStringReaderFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdStringReaderFromHandle: invalid handle")
	}
	r, ok := value.(*fatstrings.Reader)
	if !ok {
		panic("fatstdStringReaderFromHandle: handle is not a fat string reader")
	}
	return r
}

//export fatstd_go_string_new_utf8_cstr
func fatstd_go_string_new_utf8_cstr(cstr *C.char) C.uintptr_t {
	if cstr == nil {
		panic("fatstd_go_string_new_utf8_cstr: cstr is NULL")
	}
	handle := fatstdStringNewFromGoString(C.GoString(cstr))
	return C.uintptr_t(handle)
}

//export fatstd_go_string_new_utf8_n
func fatstd_go_string_new_utf8_n(bytes *C.char, len C.size_t) C.uintptr_t {
	if bytes == nil {
		if len == 0 {
			return C.uintptr_t(fatstdStringNewFromGoString(""))
		}
		panic("fatstd_go_string_new_utf8_n: bytes is NULL but len > 0")
	}
	if len > C.size_t(2147483647) {
		panic("fatstd_go_string_new_utf8_n: len too large")
	}

	buf := C.GoBytes(unsafe.Pointer(bytes), C.int(len))
	handle := fatstdStringNewFromGoString(string(buf))
	return C.uintptr_t(handle)
}

//export fatstd_go_string_clone
func fatstd_go_string_clone(handle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(handle))
	cloned := fatstrings.Clone(s.Value())
	return C.uintptr_t(fatstdStringNewFromGoString(cloned))
}

//export fatstd_go_string_contains
func fatstd_go_string_contains(aHandle C.uintptr_t, bHandle C.uintptr_t) C.int {
	a := fatstdStringFromHandle(uintptr(aHandle))
	b := fatstdStringFromHandle(uintptr(bHandle))
	if fatstrings.Contains(a.Value(), b.Value()) {
		return 1
	}
	return 0
}

//export fatstd_go_string_has_prefix
func fatstd_go_string_has_prefix(sHandle C.uintptr_t, prefixHandle C.uintptr_t) C.int {
	s := fatstdStringFromHandle(uintptr(sHandle))
	prefix := fatstdStringFromHandle(uintptr(prefixHandle))
	if fatstrings.HasPrefix(s.Value(), prefix.Value()) {
		return 1
	}
	return 0
}

//export fatstd_go_string_has_suffix
func fatstd_go_string_has_suffix(sHandle C.uintptr_t, suffixHandle C.uintptr_t) C.int {
	s := fatstdStringFromHandle(uintptr(sHandle))
	suffix := fatstdStringFromHandle(uintptr(suffixHandle))
	if fatstrings.HasSuffix(s.Value(), suffix.Value()) {
		return 1
	}
	return 0
}

//export fatstd_go_string_trim_space
func fatstd_go_string_trim_space(handle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(handle))
	trimmed := fatstrings.TrimSpace(s.Value())
	return C.uintptr_t(fatstdStringNewFromGoString(trimmed))
}

//export fatstd_go_string_trim
func fatstd_go_string_trim(sHandle C.uintptr_t, cutsetHandle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	cutset := fatstdStringFromHandle(uintptr(cutsetHandle))
	trimmed := fatstrings.Trim(s.Value(), cutset.Value())
	return C.uintptr_t(fatstdStringNewFromGoString(trimmed))
}

//export fatstd_go_string_split
func fatstd_go_string_split(sHandle C.uintptr_t, sepHandle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	sep := fatstdStringFromHandle(uintptr(sepHandle))
	return C.uintptr_t(fatstdStringArrayNew(fatstrings.Split(s.Value(), sep.Value())))
}

//export fatstd_go_string_split_n
func fatstd_go_string_split_n(sHandle C.uintptr_t, sepHandle C.uintptr_t, n C.int) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	sep := fatstdStringFromHandle(uintptr(sepHandle))
	return C.uintptr_t(fatstdStringArrayNew(fatstrings.SplitN(s.Value(), sep.Value(), int(n))))
}

//export fatstd_go_string_array_len
func fatstd_go_string_array_len(arrayHandle C.uintptr_t) C.size_t {
	a := fatstdStringArrayFromHandle(uintptr(arrayHandle))
	return C.size_t(a.Len())
}

//export fatstd_go_string_array_get
func fatstd_go_string_array_get(arrayHandle C.uintptr_t, index C.size_t) C.uintptr_t {
	a := fatstdStringArrayFromHandle(uintptr(arrayHandle))
	if index > C.size_t(2147483647) {
		panic("fatstd_go_string_array_get: index too large")
	}
	value := a.Get(int(index))
	return C.uintptr_t(fatstdStringNewFromGoString(value))
}

//export fatstd_go_string_join
func fatstd_go_string_join(arrayHandle C.uintptr_t, sepHandle C.uintptr_t) C.uintptr_t {
	a := fatstdStringArrayFromHandle(uintptr(arrayHandle))
	sep := fatstdStringFromHandle(uintptr(sepHandle))
	return C.uintptr_t(fatstdStringNewFromGoString(fatstrings.Join(a.Values(), sep.Value())))
}

//export fatstd_go_string_replace
func fatstd_go_string_replace(sHandle C.uintptr_t, oldHandle C.uintptr_t, newHandle C.uintptr_t, n C.int) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	old := fatstdStringFromHandle(uintptr(oldHandle))
	newValue := fatstdStringFromHandle(uintptr(newHandle))
	replaced := fatstrings.Replace(s.Value(), old.Value(), newValue.Value(), int(n))
	return C.uintptr_t(fatstdStringNewFromGoString(replaced))
}

//export fatstd_go_string_replace_all
func fatstd_go_string_replace_all(sHandle C.uintptr_t, oldHandle C.uintptr_t, newHandle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	old := fatstdStringFromHandle(uintptr(oldHandle))
	newValue := fatstdStringFromHandle(uintptr(newHandle))
	replaced := fatstrings.ReplaceAll(s.Value(), old.Value(), newValue.Value())
	return C.uintptr_t(fatstdStringNewFromGoString(replaced))
}

//export fatstd_go_string_to_lower
func fatstd_go_string_to_lower(handle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdStringNewFromGoString(fatstrings.ToLower(s.Value())))
}

//export fatstd_go_string_to_upper
func fatstd_go_string_to_upper(handle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdStringNewFromGoString(fatstrings.ToUpper(s.Value())))
}

//export fatstd_go_string_index
func fatstd_go_string_index(sHandle C.uintptr_t, substrHandle C.uintptr_t) C.int {
	s := fatstdStringFromHandle(uintptr(sHandle))
	substr := fatstdStringFromHandle(uintptr(substrHandle))
	return C.int(fatstrings.Index(s.Value(), substr.Value()))
}

//export fatstd_go_string_count
func fatstd_go_string_count(sHandle C.uintptr_t, substrHandle C.uintptr_t) C.int {
	s := fatstdStringFromHandle(uintptr(sHandle))
	substr := fatstdStringFromHandle(uintptr(substrHandle))
	return C.int(fatstrings.Count(s.Value(), substr.Value()))
}

//export fatstd_go_string_compare
func fatstd_go_string_compare(aHandle C.uintptr_t, bHandle C.uintptr_t) C.int {
	a := fatstdStringFromHandle(uintptr(aHandle))
	b := fatstdStringFromHandle(uintptr(bHandle))
	return C.int(fatstrings.Compare(a.Value(), b.Value()))
}

//export fatstd_go_string_equal_fold
func fatstd_go_string_equal_fold(sHandle C.uintptr_t, tHandle C.uintptr_t) C.int {
	s := fatstdStringFromHandle(uintptr(sHandle))
	t := fatstdStringFromHandle(uintptr(tHandle))
	if fatstrings.EqualFold(s.Value(), t.Value()) {
		return 1
	}
	return 0
}

//export fatstd_go_string_trim_prefix
func fatstd_go_string_trim_prefix(sHandle C.uintptr_t, prefixHandle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	prefix := fatstdStringFromHandle(uintptr(prefixHandle))
	return C.uintptr_t(fatstdStringNewFromGoString(fatstrings.TrimPrefix(s.Value(), prefix.Value())))
}

//export fatstd_go_string_trim_suffix
func fatstd_go_string_trim_suffix(sHandle C.uintptr_t, suffixHandle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	suffix := fatstdStringFromHandle(uintptr(suffixHandle))
	return C.uintptr_t(fatstdStringNewFromGoString(fatstrings.TrimSuffix(s.Value(), suffix.Value())))
}

//export fatstd_go_string_cut
func fatstd_go_string_cut(sHandle C.uintptr_t, sepHandle C.uintptr_t, beforeOut *C.uintptr_t, afterOut *C.uintptr_t) C.int {
	if beforeOut == nil {
		panic("fatstd_go_string_cut: beforeOut is NULL")
	}
	if afterOut == nil {
		panic("fatstd_go_string_cut: afterOut is NULL")
	}

	s := fatstdStringFromHandle(uintptr(sHandle))
	sep := fatstdStringFromHandle(uintptr(sepHandle))
	before, after, found := fatstrings.Cut(s.Value(), sep.Value())

	*beforeOut = C.uintptr_t(fatstdStringNewFromGoString(before))
	*afterOut = C.uintptr_t(fatstdStringNewFromGoString(after))

	if found {
		return 1
	}
	return 0
}

//export fatstd_go_string_cut_prefix
func fatstd_go_string_cut_prefix(sHandle C.uintptr_t, prefixHandle C.uintptr_t, afterOut *C.uintptr_t) C.int {
	if afterOut == nil {
		panic("fatstd_go_string_cut_prefix: afterOut is NULL")
	}

	s := fatstdStringFromHandle(uintptr(sHandle))
	prefix := fatstdStringFromHandle(uintptr(prefixHandle))
	after, found := fatstrings.CutPrefix(s.Value(), prefix.Value())

	*afterOut = C.uintptr_t(fatstdStringNewFromGoString(after))
	if found {
		return 1
	}
	return 0
}

//export fatstd_go_string_cut_suffix
func fatstd_go_string_cut_suffix(sHandle C.uintptr_t, suffixHandle C.uintptr_t, afterOut *C.uintptr_t) C.int {
	if afterOut == nil {
		panic("fatstd_go_string_cut_suffix: afterOut is NULL")
	}

	s := fatstdStringFromHandle(uintptr(sHandle))
	suffix := fatstdStringFromHandle(uintptr(suffixHandle))
	after, found := fatstrings.CutSuffix(s.Value(), suffix.Value())

	*afterOut = C.uintptr_t(fatstdStringNewFromGoString(after))
	if found {
		return 1
	}
	return 0
}

//export fatstd_go_string_fields
func fatstd_go_string_fields(sHandle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	return C.uintptr_t(fatstdStringArrayNew(fatstrings.Fields(s.Value())))
}

//export fatstd_go_string_repeat
func fatstd_go_string_repeat(sHandle C.uintptr_t, count C.int) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	return C.uintptr_t(fatstdStringNewFromGoString(fatstrings.Repeat(s.Value(), int(count))))
}

//export fatstd_go_string_contains_any
func fatstd_go_string_contains_any(sHandle C.uintptr_t, charsHandle C.uintptr_t) C.int {
	s := fatstdStringFromHandle(uintptr(sHandle))
	chars := fatstdStringFromHandle(uintptr(charsHandle))
	if fatstrings.ContainsAny(s.Value(), chars.Value()) {
		return 1
	}
	return 0
}

//export fatstd_go_string_index_any
func fatstd_go_string_index_any(sHandle C.uintptr_t, charsHandle C.uintptr_t) C.int {
	s := fatstdStringFromHandle(uintptr(sHandle))
	chars := fatstdStringFromHandle(uintptr(charsHandle))
	if fatstrings.IndexAny(s.Value(), chars.Value()) {
		return 1
	}
	return 0
}

//export fatstd_go_string_to_valid_utf8
func fatstd_go_string_to_valid_utf8(sHandle C.uintptr_t, replacementHandle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	replacement := fatstdStringFromHandle(uintptr(replacementHandle))
	return C.uintptr_t(fatstdStringNewFromGoString(fatstrings.ToValidUTF8(s.Value(), replacement.Value())))
}

//export fatstd_go_string_builder_new
func fatstd_go_string_builder_new() C.uintptr_t {
	return C.uintptr_t(fatstdStringBuilderNew())
}

//export fatstd_go_string_builder_cap
func fatstd_go_string_builder_cap(handle C.uintptr_t) C.size_t {
	b := fatstdStringBuilderFromHandle(uintptr(handle))
	return C.size_t(b.Cap())
}

//export fatstd_go_string_builder_grow
func fatstd_go_string_builder_grow(handle C.uintptr_t, n C.size_t) {
	if n > C.size_t(2147483647) {
		panic("fatstd_go_string_builder_grow: n too large")
	}
	b := fatstdStringBuilderFromHandle(uintptr(handle))
	b.Grow(int(n))
}

//export fatstd_go_string_builder_len
func fatstd_go_string_builder_len(handle C.uintptr_t) C.size_t {
	b := fatstdStringBuilderFromHandle(uintptr(handle))
	return C.size_t(b.Len())
}

//export fatstd_go_string_builder_reset
func fatstd_go_string_builder_reset(handle C.uintptr_t) {
	b := fatstdStringBuilderFromHandle(uintptr(handle))
	b.Reset()
}

//export fatstd_go_string_builder_string
func fatstd_go_string_builder_string(handle C.uintptr_t) C.uintptr_t {
	b := fatstdStringBuilderFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdStringNewFromGoString(b.String()))
}

//export fatstd_go_string_builder_write
func fatstd_go_string_builder_write(handle C.uintptr_t, bytes *C.char, len C.size_t) C.size_t {
	if bytes == nil {
		if len == 0 {
			return 0
		}
		panic("fatstd_go_string_builder_write: bytes is NULL but len > 0")
	}
	if len > C.size_t(2147483647) {
		panic("fatstd_go_string_builder_write: len too large")
	}
	b := fatstdStringBuilderFromHandle(uintptr(handle))
	buf := C.GoBytes(unsafe.Pointer(bytes), C.int(len))
	return C.size_t(b.Write(buf))
}

//export fatstd_go_string_builder_write_byte
func fatstd_go_string_builder_write_byte(handle C.uintptr_t, c C.uchar) {
	b := fatstdStringBuilderFromHandle(uintptr(handle))
	b.WriteByte(byte(c))
}

//export fatstd_go_string_builder_write_string
func fatstd_go_string_builder_write_string(builderHandle C.uintptr_t, sHandle C.uintptr_t) C.size_t {
	b := fatstdStringBuilderFromHandle(uintptr(builderHandle))
	s := fatstdStringFromHandle(uintptr(sHandle))
	return C.size_t(b.WriteString(s.Value()))
}

//export fatstd_go_string_builder_free
func fatstd_go_string_builder_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_string_builder_free: handle is 0")
	}

	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_string_builder_free: invalid handle")
	}
	if _, ok := value.(*fatstrings.Builder); !ok {
		panic("fatstd_go_string_builder_free: handle is not a fat string builder")
	}
}

//export fatstd_go_string_reader_new
func fatstd_go_string_reader_new(sHandle C.uintptr_t) C.uintptr_t {
	s := fatstdStringFromHandle(uintptr(sHandle))
	return C.uintptr_t(fatstdStringReaderNew(s.Value()))
}

//export fatstd_go_string_reader_len
func fatstd_go_string_reader_len(handle C.uintptr_t) C.size_t {
	r := fatstdStringReaderFromHandle(uintptr(handle))
	return C.size_t(r.Len())
}

//export fatstd_go_string_reader_size
func fatstd_go_string_reader_size(handle C.uintptr_t) C.longlong {
	r := fatstdStringReaderFromHandle(uintptr(handle))
	return C.longlong(r.Size())
}

//export fatstd_go_string_reader_reset
func fatstd_go_string_reader_reset(readerHandle C.uintptr_t, sHandle C.uintptr_t) {
	r := fatstdStringReaderFromHandle(uintptr(readerHandle))
	s := fatstdStringFromHandle(uintptr(sHandle))
	r.Reset(s.Value())
}

//export fatstd_go_string_reader_read
func fatstd_go_string_reader_read(readerHandle C.uintptr_t, bytes *C.char, len C.size_t, eofOut *C.bool) C.size_t {
	if eofOut == nil {
		panic("fatstd_go_string_reader_read: eofOut is NULL")
	}
	if bytes == nil {
		if len == 0 {
			*eofOut = false
			return 0
		}
		panic("fatstd_go_string_reader_read: bytes is NULL but len > 0")
	}
	if len > C.size_t(2147483647) {
		panic("fatstd_go_string_reader_read: len too large")
	}

	r := fatstdStringReaderFromHandle(uintptr(readerHandle))
	dst := unsafe.Slice((*byte)(unsafe.Pointer(bytes)), int(len))

	n, err := r.Read(dst)
	if err == io.EOF {
		*eofOut = true
		return C.size_t(n)
	}
	if err != nil {
		panic("fatstd_go_string_reader_read: unexpected error")
	}
	*eofOut = false
	return C.size_t(n)
}

//export fatstd_go_string_reader_read_at
func fatstd_go_string_reader_read_at(readerHandle C.uintptr_t, bytes *C.char, len C.size_t, off C.longlong, eofOut *C.bool) C.size_t {
	if eofOut == nil {
		panic("fatstd_go_string_reader_read_at: eofOut is NULL")
	}
	if bytes == nil {
		if len == 0 {
			*eofOut = false
			return 0
		}
		panic("fatstd_go_string_reader_read_at: bytes is NULL but len > 0")
	}
	if len > C.size_t(2147483647) {
		panic("fatstd_go_string_reader_read_at: len too large")
	}

	r := fatstdStringReaderFromHandle(uintptr(readerHandle))
	dst := unsafe.Slice((*byte)(unsafe.Pointer(bytes)), int(len))

	n, err := r.ReadAt(dst, int64(off))
	if err == io.EOF {
		*eofOut = true
		return C.size_t(n)
	}
	if err != nil {
		panic("fatstd_go_string_reader_read_at: unexpected error")
	}
	*eofOut = false
	return C.size_t(n)
}

//export fatstd_go_string_reader_read_byte
func fatstd_go_string_reader_read_byte(readerHandle C.uintptr_t, byteOut *C.uchar, eofOut *C.bool) C.bool {
	if byteOut == nil {
		panic("fatstd_go_string_reader_read_byte: byteOut is NULL")
	}
	if eofOut == nil {
		panic("fatstd_go_string_reader_read_byte: eofOut is NULL")
	}

	r := fatstdStringReaderFromHandle(uintptr(readerHandle))
	b, err := r.ReadByte()
	if err == io.EOF {
		*eofOut = true
		return false
	}
	if err != nil {
		panic("fatstd_go_string_reader_read_byte: unexpected error")
	}
	*byteOut = C.uchar(b)
	*eofOut = false
	return true
}

//export fatstd_go_string_reader_unread_byte
func fatstd_go_string_reader_unread_byte(readerHandle C.uintptr_t) {
	r := fatstdStringReaderFromHandle(uintptr(readerHandle))
	if err := r.UnreadByte(); err != nil {
		panic("fatstd_go_string_reader_unread_byte: invalid unread")
	}
}

//export fatstd_go_string_reader_seek
func fatstd_go_string_reader_seek(readerHandle C.uintptr_t, offset C.longlong, whence C.int) C.longlong {
	r := fatstdStringReaderFromHandle(uintptr(readerHandle))
	pos, err := r.Seek(int64(offset), int(whence))
	if err != nil {
		panic("fatstd_go_string_reader_seek: invalid seek")
	}
	return C.longlong(pos)
}

//export fatstd_go_string_reader_write_to_builder
func fatstd_go_string_reader_write_to_builder(readerHandle C.uintptr_t, builderHandle C.uintptr_t) C.longlong {
	r := fatstdStringReaderFromHandle(uintptr(readerHandle))
	b := fatstdStringBuilderFromHandle(uintptr(builderHandle))
	n, err := r.WriteTo(b.Underlying())
	if err != nil {
		panic("fatstd_go_string_reader_write_to_builder: unexpected error")
	}
	return C.longlong(n)
}

//export fatstd_go_string_reader_free
func fatstd_go_string_reader_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_string_reader_free: handle is 0")
	}

	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_string_reader_free: invalid handle")
	}
	if _, ok := value.(*fatstrings.Reader); !ok {
		panic("fatstd_go_string_reader_free: handle is not a fat string reader")
	}
}

//export fatstd_go_string_array_free
func fatstd_go_string_array_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_string_array_free: handle is 0")
	}

	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_string_array_free: invalid handle")
	}
	if _, ok := value.(*fatstrings.StringArray); !ok {
		panic("fatstd_go_string_array_free: handle is not a fat string array")
	}
}

//export fatstd_go_string_free
func fatstd_go_string_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_string_free: handle is 0")
	}

	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_string_free: invalid handle")
	}
	if _, ok := value.(*fatstrings.String); !ok {
		panic("fatstd_go_string_free: handle is not a fat string")
	}
}

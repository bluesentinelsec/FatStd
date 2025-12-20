#include "fat/zip.h"

#include "fatstd_go.h"

fat_Status fat_ZipReaderOpenPathUTF8(const char *path, fat_ZipReader *out_reader, fat_Error *out_err) {
  return (fat_Status)fatstd_go_zip_reader_open_path_utf8((char *)path, (uintptr_t *)out_reader, (uintptr_t *)out_err);
}

fat_Status fat_ZipReaderNewBytes(fat_Bytes zip_bytes, fat_ZipReader *out_reader, fat_Error *out_err) {
  return (fat_Status)fatstd_go_zip_reader_new_bytes((uintptr_t)zip_bytes, (uintptr_t *)out_reader, (uintptr_t *)out_err);
}

fat_Status fat_ZipReaderFree(fat_ZipReader r, fat_Error *out_err) {
  return (fat_Status)fatstd_go_zip_reader_free((uintptr_t)r, (uintptr_t *)out_err);
}

size_t fat_ZipReaderNumFiles(fat_ZipReader r) {
  return (size_t)fatstd_go_zip_reader_num_files((uintptr_t)r);
}

fat_ZipFile fat_ZipReaderFileByIndex(fat_ZipReader r, size_t idx) {
  return (fat_ZipFile)fatstd_go_zip_reader_file_by_index((uintptr_t)r, idx);
}

void fat_ZipFileFree(fat_ZipFile f) {
  fatstd_go_zip_file_free((uintptr_t)f);
}

fat_String fat_ZipFileName(fat_ZipFile f) {
  return (fat_String)fatstd_go_zip_file_name((uintptr_t)f);
}

fat_Status fat_ZipFileOpen(fat_ZipFile f, fat_ZipFileReader *out_reader, fat_Error *out_err) {
  return (fat_Status)fatstd_go_zip_file_open((uintptr_t)f, (uintptr_t *)out_reader, (uintptr_t *)out_err);
}

fat_Status fat_ZipFileReaderRead(
  fat_ZipFileReader r,
  void *dst,
  size_t dst_len,
  size_t *out_n,
  bool *out_eof,
  fat_Error *out_err
) {
  return (fat_Status)fatstd_go_zip_file_reader_read((uintptr_t)r, (char *)dst, dst_len, out_n, (_Bool *)out_eof,
                                                    (uintptr_t *)out_err);
}

fat_Status fat_ZipFileReaderClose(fat_ZipFileReader r, fat_Error *out_err) {
  return (fat_Status)fatstd_go_zip_file_reader_close((uintptr_t)r, (uintptr_t *)out_err);
}

fat_Status fat_ZipWriterNewToBytesBuffer(fat_BytesBuffer dst, fat_ZipWriter *out_writer, fat_Error *out_err) {
  return (fat_Status)fatstd_go_zip_writer_new_to_bytes_buffer((uintptr_t)dst, (uintptr_t *)out_writer,
                                                              (uintptr_t *)out_err);
}

fat_Status fat_ZipWriterAddBytes(fat_ZipWriter w, fat_String name, fat_Bytes data, fat_Error *out_err) {
  return (fat_Status)fatstd_go_zip_writer_add_bytes((uintptr_t)w, (uintptr_t)name, (uintptr_t)data, (uintptr_t *)out_err);
}

fat_Status fat_ZipWriterClose(fat_ZipWriter w, fat_Error *out_err) {
  return (fat_Status)fatstd_go_zip_writer_close((uintptr_t)w, (uintptr_t *)out_err);
}


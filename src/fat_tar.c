#include "fat/tar.h"

#include "fatstd_go.h"

fat_Status fat_TarReaderNewBytes(fat_Bytes tar_bytes, fat_TarReader *out_reader, fat_Error *out_err) {
  return (fat_Status)fatstd_go_tar_reader_new_bytes((uintptr_t)tar_bytes, (uintptr_t *)out_reader, (uintptr_t *)out_err);
}

fat_Status fat_TarReaderOpenPathUTF8(const char *path, fat_TarReader *out_reader, fat_Error *out_err) {
  return (fat_Status)fatstd_go_tar_reader_open_path_utf8((char *)path, (uintptr_t *)out_reader, (uintptr_t *)out_err);
}

fat_Status fat_TarReaderFree(fat_TarReader r, fat_Error *out_err) {
  return (fat_Status)fatstd_go_tar_reader_free((uintptr_t)r, (uintptr_t *)out_err);
}

fat_Status fat_TarReaderNext(fat_TarReader r, fat_TarHeader *out_hdr, bool *out_eof, fat_Error *out_err) {
  return (fat_Status)fatstd_go_tar_reader_next((uintptr_t)r, (uintptr_t *)out_hdr, (_Bool *)out_eof, (uintptr_t *)out_err);
}

fat_Status fat_TarReaderRead(
  fat_TarReader r,
  void *dst,
  size_t dst_len,
  size_t *out_n,
  bool *out_eof,
  fat_Error *out_err
) {
  return (fat_Status)fatstd_go_tar_reader_read((uintptr_t)r, (char *)dst, dst_len, out_n, (_Bool *)out_eof, (uintptr_t *)out_err);
}

void fat_TarHeaderFree(fat_TarHeader h) {
  fatstd_go_tar_header_free((uintptr_t)h);
}

fat_String fat_TarHeaderName(fat_TarHeader h) {
  return (fat_String)fatstd_go_tar_header_name((uintptr_t)h);
}

uint8_t fat_TarHeaderTypeflag(fat_TarHeader h) {
  return (uint8_t)fatstd_go_tar_header_typeflag((uintptr_t)h);
}

int64_t fat_TarHeaderSize(fat_TarHeader h) {
  return (int64_t)fatstd_go_tar_header_size((uintptr_t)h);
}

fat_Status fat_TarWriterNewToBytesBuffer(fat_BytesBuffer dst, fat_TarWriter *out_writer, fat_Error *out_err) {
  return (fat_Status)fatstd_go_tar_writer_new_to_bytes_buffer((uintptr_t)dst, (uintptr_t *)out_writer, (uintptr_t *)out_err);
}

fat_Status fat_TarWriterAddBytes(fat_TarWriter w, fat_String name, fat_Bytes data, fat_Error *out_err) {
  return (fat_Status)fatstd_go_tar_writer_add_bytes((uintptr_t)w, (uintptr_t)name, (uintptr_t)data, (uintptr_t *)out_err);
}

fat_Status fat_TarWriterFlush(fat_TarWriter w, fat_Error *out_err) {
  return (fat_Status)fatstd_go_tar_writer_flush((uintptr_t)w, (uintptr_t *)out_err);
}

fat_Status fat_TarWriterClose(fat_TarWriter w, fat_Error *out_err) {
  return (fat_Status)fatstd_go_tar_writer_close((uintptr_t)w, (uintptr_t *)out_err);
}


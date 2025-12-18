#include "fat/bytes_reader.h"

#include "fatstd_go.h"

fat_BytesReader fat_BytesReaderNew(fat_Bytes b) {
  return (fat_BytesReader)fatstd_go_bytes_reader_new((uintptr_t)b);
}

void fat_BytesReaderFree(fat_BytesReader r) {
  fatstd_go_bytes_reader_free((uintptr_t)r);
}

size_t fat_BytesReaderLen(fat_BytesReader r) {
  return (size_t)fatstd_go_bytes_reader_len((uintptr_t)r);
}

int64_t fat_BytesReaderSize(fat_BytesReader r) {
  return (int64_t)fatstd_go_bytes_reader_size((uintptr_t)r);
}

void fat_BytesReaderReset(fat_BytesReader r, fat_Bytes b) {
  fatstd_go_bytes_reader_reset((uintptr_t)r, (uintptr_t)b);
}

size_t fat_BytesReaderRead(fat_BytesReader r, void *dst, size_t len, bool *eof_out) {
  return (size_t)fatstd_go_bytes_reader_read((uintptr_t)r, (char *)dst, len, (_Bool *)eof_out);
}

size_t fat_BytesReaderReadAt(fat_BytesReader r, void *dst, size_t len, int64_t off, bool *eof_out) {
  return (size_t)fatstd_go_bytes_reader_read_at((uintptr_t)r, (char *)dst, len, off, (_Bool *)eof_out);
}

bool fat_BytesReaderReadByte(fat_BytesReader r, uint8_t *byte_out, bool *eof_out) {
  return (bool)fatstd_go_bytes_reader_read_byte((uintptr_t)r, byte_out, (_Bool *)eof_out);
}

void fat_BytesReaderUnreadByte(fat_BytesReader r) {
  fatstd_go_bytes_reader_unread_byte((uintptr_t)r);
}

int64_t fat_BytesReaderSeek(fat_BytesReader r, int64_t offset, int whence) {
  return (int64_t)fatstd_go_bytes_reader_seek((uintptr_t)r, offset, whence);
}

int64_t fat_BytesReaderWriteToBytesBuffer(fat_BytesReader r, fat_BytesBuffer b) {
  return (int64_t)fatstd_go_bytes_reader_write_to_bytes_buffer((uintptr_t)r, (uintptr_t)b);
}


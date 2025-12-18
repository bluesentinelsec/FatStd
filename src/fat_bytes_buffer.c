#include "fat/bytes_buffer.h"

#include "fatstd_go.h"

fat_BytesBuffer fat_BytesBufferNew(void) {
  return (fat_BytesBuffer)fatstd_go_bytes_buffer_new();
}

fat_BytesBuffer fat_BytesBufferNewBytes(fat_Bytes b) {
  return (fat_BytesBuffer)fatstd_go_bytes_buffer_new_bytes((uintptr_t)b);
}

fat_BytesBuffer fat_BytesBufferNewN(const void *bytes, size_t len) {
  return (fat_BytesBuffer)fatstd_go_bytes_buffer_new_n((char *)bytes, len);
}

fat_BytesBuffer fat_BytesBufferNewString(fat_String s) {
  return (fat_BytesBuffer)fatstd_go_bytes_buffer_new_string((uintptr_t)s);
}

void fat_BytesBufferFree(fat_BytesBuffer b) {
  fatstd_go_bytes_buffer_free((uintptr_t)b);
}

size_t fat_BytesBufferLen(fat_BytesBuffer b) {
  return (size_t)fatstd_go_bytes_buffer_len((uintptr_t)b);
}

size_t fat_BytesBufferCap(fat_BytesBuffer b) {
  return (size_t)fatstd_go_bytes_buffer_cap((uintptr_t)b);
}

void fat_BytesBufferGrow(fat_BytesBuffer b, size_t n) {
  fatstd_go_bytes_buffer_grow((uintptr_t)b, n);
}

void fat_BytesBufferReset(fat_BytesBuffer b) {
  fatstd_go_bytes_buffer_reset((uintptr_t)b);
}

void fat_BytesBufferTruncate(fat_BytesBuffer b, size_t n) {
  fatstd_go_bytes_buffer_truncate((uintptr_t)b, n);
}

size_t fat_BytesBufferWrite(fat_BytesBuffer b, const void *bytes, size_t len) {
  return (size_t)fatstd_go_bytes_buffer_write((uintptr_t)b, (char *)bytes, len);
}

void fat_BytesBufferWriteByte(fat_BytesBuffer b, uint8_t c) {
  fatstd_go_bytes_buffer_write_byte((uintptr_t)b, c);
}

size_t fat_BytesBufferWriteRune(fat_BytesBuffer b, uint32_t r) {
  return (size_t)fatstd_go_bytes_buffer_write_rune((uintptr_t)b, r);
}

size_t fat_BytesBufferWriteString(fat_BytesBuffer b, fat_String s) {
  return (size_t)fatstd_go_bytes_buffer_write_string((uintptr_t)b, (uintptr_t)s);
}

fat_Bytes fat_BytesBufferBytes(fat_BytesBuffer b) {
  return (fat_Bytes)fatstd_go_bytes_buffer_bytes((uintptr_t)b);
}

fat_String fat_BytesBufferString(fat_BytesBuffer b) {
  return (fat_String)fatstd_go_bytes_buffer_string((uintptr_t)b);
}

size_t fat_BytesBufferRead(fat_BytesBuffer b, void *dst, size_t len, bool *eof_out) {
  return (size_t)fatstd_go_bytes_buffer_read((uintptr_t)b, (char *)dst, len, (_Bool *)eof_out);
}

fat_Bytes fat_BytesBufferNext(fat_BytesBuffer b, size_t n) {
  return (fat_Bytes)fatstd_go_bytes_buffer_next((uintptr_t)b, n);
}

bool fat_BytesBufferReadByte(fat_BytesBuffer b, uint8_t *byte_out, bool *eof_out) {
  return (bool)fatstd_go_bytes_buffer_read_byte((uintptr_t)b, byte_out, (_Bool *)eof_out);
}

int64_t fat_BytesBufferWriteToBytesBuffer(fat_BytesBuffer src, fat_BytesBuffer dst) {
  return (int64_t)fatstd_go_bytes_buffer_write_to_bytes_buffer((uintptr_t)src, (uintptr_t)dst);
}

int64_t fat_BytesBufferReadFromStringReader(fat_BytesBuffer dst, fat_StringReader r) {
  return (int64_t)fatstd_go_bytes_buffer_read_from_string_reader((uintptr_t)dst, (uintptr_t)r);
}


#include "fat/string_reader.h"

#include "fatstd_go.h"

fat_StringReader fat_StringReaderNew(fat_String s) {
  return (fat_StringReader)fatstd_go_string_reader_new((uintptr_t)s);
}

void fat_StringReaderFree(fat_StringReader r) {
  fatstd_go_string_reader_free((uintptr_t)r);
}

size_t fat_StringReaderLen(fat_StringReader r) {
  return (size_t)fatstd_go_string_reader_len((uintptr_t)r);
}

int64_t fat_StringReaderSize(fat_StringReader r) {
  return (int64_t)fatstd_go_string_reader_size((uintptr_t)r);
}

void fat_StringReaderReset(fat_StringReader r, fat_String s) {
  fatstd_go_string_reader_reset((uintptr_t)r, (uintptr_t)s);
}

size_t fat_StringReaderRead(fat_StringReader r, void *buf, size_t len, bool *eof_out) {
  return (size_t)fatstd_go_string_reader_read((uintptr_t)r, (char *)buf, len, (_Bool *)eof_out);
}

size_t fat_StringReaderReadAt(fat_StringReader r, void *buf, size_t len, int64_t off, bool *eof_out) {
  return (size_t)fatstd_go_string_reader_read_at((uintptr_t)r, (char *)buf, len, off, (_Bool *)eof_out);
}

bool fat_StringReaderReadByte(fat_StringReader r, uint8_t *byte_out, bool *eof_out) {
  return (bool)fatstd_go_string_reader_read_byte((uintptr_t)r, byte_out, (_Bool *)eof_out);
}

void fat_StringReaderUnreadByte(fat_StringReader r) {
  fatstd_go_string_reader_unread_byte((uintptr_t)r);
}

int64_t fat_StringReaderSeek(fat_StringReader r, int64_t offset, int whence) {
  return (int64_t)fatstd_go_string_reader_seek((uintptr_t)r, offset, whence);
}

int64_t fat_StringReaderWriteToBuilder(fat_StringReader r, fat_StringBuilder b) {
  return (int64_t)fatstd_go_string_reader_write_to_builder((uintptr_t)r, (uintptr_t)b);
}


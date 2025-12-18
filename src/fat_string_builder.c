#include "fat/string_builder.h"

#include "fatstd_go.h"

fat_StringBuilder fat_StringBuilderNew(void) {
  return (fat_StringBuilder)fatstd_go_string_builder_new();
}

void fat_StringBuilderFree(fat_StringBuilder b) {
  fatstd_go_string_builder_free((uintptr_t)b);
}

size_t fat_StringBuilderCap(fat_StringBuilder b) {
  return (size_t)fatstd_go_string_builder_cap((uintptr_t)b);
}

size_t fat_StringBuilderLen(fat_StringBuilder b) {
  return (size_t)fatstd_go_string_builder_len((uintptr_t)b);
}

void fat_StringBuilderGrow(fat_StringBuilder b, size_t n) {
  fatstd_go_string_builder_grow((uintptr_t)b, n);
}

void fat_StringBuilderReset(fat_StringBuilder b) {
  fatstd_go_string_builder_reset((uintptr_t)b);
}

fat_String fat_StringBuilderString(fat_StringBuilder b) {
  return (fat_String)fatstd_go_string_builder_string((uintptr_t)b);
}

size_t fat_StringBuilderWrite(fat_StringBuilder b, const void *bytes, size_t len) {
  return (size_t)fatstd_go_string_builder_write((uintptr_t)b, (char *)bytes, len);
}

void fat_StringBuilderWriteByte(fat_StringBuilder b, uint8_t c) {
  fatstd_go_string_builder_write_byte((uintptr_t)b, c);
}

size_t fat_StringBuilderWriteString(fat_StringBuilder b, fat_String s) {
  return (size_t)fatstd_go_string_builder_write_string((uintptr_t)b, (uintptr_t)s);
}


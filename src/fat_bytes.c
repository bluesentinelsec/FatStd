#include "fat/bytes.h"

#include "fatstd_go.h"

fat_Bytes fat_BytesNewN(const void *bytes, size_t len) {
  return (fat_Bytes)fatstd_go_bytes_new_n((char *)bytes, len);
}

size_t fat_BytesLen(fat_Bytes b) {
  return (size_t)fatstd_go_bytes_len((uintptr_t)b);
}

size_t fat_BytesCopyOut(fat_Bytes b, void *dst, size_t dst_len) {
  return (size_t)fatstd_go_bytes_copy_out((uintptr_t)b, (char *)dst, dst_len);
}

fat_Bytes fat_BytesClone(fat_Bytes b) {
  return (fat_Bytes)fatstd_go_bytes_clone((uintptr_t)b);
}

bool fat_BytesContains(fat_Bytes b, fat_Bytes subslice) {
  return (bool)fatstd_go_bytes_contains((uintptr_t)b, (uintptr_t)subslice);
}

bool fat_BytesHasPrefix(fat_Bytes s, fat_Bytes prefix) {
  return (bool)fatstd_go_bytes_has_prefix((uintptr_t)s, (uintptr_t)prefix);
}

bool fat_BytesHasSuffix(fat_Bytes s, fat_Bytes suffix) {
  return (bool)fatstd_go_bytes_has_suffix((uintptr_t)s, (uintptr_t)suffix);
}

fat_Bytes fat_BytesTrimSpace(fat_Bytes s) {
  return (fat_Bytes)fatstd_go_bytes_trim_space((uintptr_t)s);
}

fat_Bytes fat_BytesTrim(fat_Bytes s, fat_String cutset) {
  return (fat_Bytes)fatstd_go_bytes_trim((uintptr_t)s, (uintptr_t)cutset);
}

fat_BytesArray fat_BytesSplit(fat_Bytes s, fat_Bytes sep) {
  return (fat_BytesArray)fatstd_go_bytes_split((uintptr_t)s, (uintptr_t)sep);
}

size_t fat_BytesArrayLen(fat_BytesArray a) {
  return (size_t)fatstd_go_bytes_array_len((uintptr_t)a);
}

fat_Bytes fat_BytesArrayGet(fat_BytesArray a, size_t idx) {
  return (fat_Bytes)fatstd_go_bytes_array_get((uintptr_t)a, idx);
}

void fat_BytesArrayFree(fat_BytesArray a) {
  fatstd_go_bytes_array_free((uintptr_t)a);
}

fat_Bytes fat_BytesJoin(fat_BytesArray s, fat_Bytes sep) {
  return (fat_Bytes)fatstd_go_bytes_join((uintptr_t)s, (uintptr_t)sep);
}

fat_Bytes fat_BytesReplaceAll(fat_Bytes s, fat_Bytes old_value, fat_Bytes new_value) {
  return (fat_Bytes)fatstd_go_bytes_replace_all((uintptr_t)s, (uintptr_t)old_value, (uintptr_t)new_value);
}

fat_Bytes fat_BytesReplace(fat_Bytes s, fat_Bytes old_value, fat_Bytes new_value, int n) {
  return (fat_Bytes)fatstd_go_bytes_replace((uintptr_t)s, (uintptr_t)old_value, (uintptr_t)new_value, n);
}

fat_Bytes fat_BytesToLower(fat_Bytes s) {
  return (fat_Bytes)fatstd_go_bytes_to_lower((uintptr_t)s);
}

fat_Bytes fat_BytesToUpper(fat_Bytes s) {
  return (fat_Bytes)fatstd_go_bytes_to_upper((uintptr_t)s);
}

int fat_BytesIndex(fat_Bytes s, fat_Bytes sep) {
  return (int)fatstd_go_bytes_index((uintptr_t)s, (uintptr_t)sep);
}

int fat_BytesCount(fat_Bytes s, fat_Bytes sep) {
  return (int)fatstd_go_bytes_count((uintptr_t)s, (uintptr_t)sep);
}

int fat_BytesCompare(fat_Bytes a, fat_Bytes b) {
  return (int)fatstd_go_bytes_compare((uintptr_t)a, (uintptr_t)b);
}

bool fat_BytesEqual(fat_Bytes a, fat_Bytes b) {
  return (bool)fatstd_go_bytes_equal((uintptr_t)a, (uintptr_t)b);
}

void fat_BytesFree(fat_Bytes b) {
  fatstd_go_bytes_free((uintptr_t)b);
}


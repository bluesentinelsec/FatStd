#include "fat/string.h"

#include "fatstd_go.h"

fat_String fat_StringNewUTF8(const char *cstr) {
  return (fat_String)fatstd_go_string_new_utf8_cstr((char *)cstr);
}

fat_String fat_StringNewUTF8N(const char *bytes, size_t len) {
  return (fat_String)fatstd_go_string_new_utf8_n((char *)bytes, len);
}

fat_String fat_StringClone(fat_String s) {
  return (fat_String)fatstd_go_string_clone((uintptr_t)s);
}

bool fat_StringContains(fat_String s, fat_String substr) {
  return (bool)fatstd_go_string_contains((uintptr_t)s, (uintptr_t)substr);
}

bool fat_StringHasPrefix(fat_String s, fat_String prefix) {
  return (bool)fatstd_go_string_has_prefix((uintptr_t)s, (uintptr_t)prefix);
}

bool fat_StringHasSuffix(fat_String s, fat_String suffix) {
  return (bool)fatstd_go_string_has_suffix((uintptr_t)s, (uintptr_t)suffix);
}

fat_String fat_StringTrimSpace(fat_String s) {
  return (fat_String)fatstd_go_string_trim_space((uintptr_t)s);
}

fat_String fat_StringTrim(fat_String s, fat_String cutset) {
  return (fat_String)fatstd_go_string_trim((uintptr_t)s, (uintptr_t)cutset);
}

fat_StringArray fat_StringSplit(fat_String s, fat_String sep) {
  return (fat_StringArray)fatstd_go_string_split((uintptr_t)s, (uintptr_t)sep);
}

fat_StringArray fat_StringSplitN(fat_String s, fat_String sep, int n) {
  return (fat_StringArray)fatstd_go_string_split_n((uintptr_t)s, (uintptr_t)sep, n);
}

size_t fat_StringArrayLen(fat_StringArray a) {
  return (size_t)fatstd_go_string_array_len((uintptr_t)a);
}

fat_String fat_StringArrayGet(fat_StringArray a, size_t idx) {
  return (fat_String)fatstd_go_string_array_get((uintptr_t)a, idx);
}

void fat_StringArrayFree(fat_StringArray a) {
  fatstd_go_string_array_free((uintptr_t)a);
}

fat_String fat_StringJoin(fat_StringArray elems, fat_String sep) {
  return (fat_String)fatstd_go_string_join((uintptr_t)elems, (uintptr_t)sep);
}

fat_String fat_StringReplace(fat_String s, fat_String old, fat_String new_value, int n) {
  return (fat_String)fatstd_go_string_replace((uintptr_t)s, (uintptr_t)old, (uintptr_t)new_value, n);
}

fat_String fat_StringReplaceAll(fat_String s, fat_String old, fat_String new_value) {
  return (fat_String)fatstd_go_string_replace_all((uintptr_t)s, (uintptr_t)old, (uintptr_t)new_value);
}

void fat_StringFree(fat_String s) {
  fatstd_go_string_free((uintptr_t)s);
}

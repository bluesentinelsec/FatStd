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

void fat_StringFree(fat_String s) {
  fatstd_go_string_free((uintptr_t)s);
}

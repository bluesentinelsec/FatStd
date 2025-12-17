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

void fat_StringFree(fat_String s) {
  fatstd_go_string_free((uintptr_t)s);
}

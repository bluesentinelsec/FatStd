#include "fat/string.h"

#include "fatstd_go.h"

fat_String fat_StringNewUTF8(const char *cstr) {
  return (fat_String)fatstd_go_string_new_utf8_cstr((char *)cstr);
}

fat_String fat_StringNewUTF8N(const char *bytes, size_t len) {
  return (fat_String)fatstd_go_string_new_utf8_n((char *)bytes, len);
}

size_t fat_StringLenBytes(fat_String s) {
  return (size_t)fatstd_go_string_len_bytes((uintptr_t)s);
}

size_t fat_StringCopyOut(fat_String s, void *dst, size_t dst_len) {
  return (size_t)fatstd_go_string_copy_out((uintptr_t)s, (char *)dst, dst_len);
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

fat_String fat_StringToLower(fat_String s) {
  return (fat_String)fatstd_go_string_to_lower((uintptr_t)s);
}

fat_String fat_StringToUpper(fat_String s) {
  return (fat_String)fatstd_go_string_to_upper((uintptr_t)s);
}

int fat_StringIndex(fat_String s, fat_String substr) {
  return (int)fatstd_go_string_index((uintptr_t)s, (uintptr_t)substr);
}

int fat_StringCount(fat_String s, fat_String substr) {
  return (int)fatstd_go_string_count((uintptr_t)s, (uintptr_t)substr);
}

int fat_StringCompare(fat_String a, fat_String b) {
  return (int)fatstd_go_string_compare((uintptr_t)a, (uintptr_t)b);
}

bool fat_StringEqualFold(fat_String s, fat_String t) {
  return (bool)fatstd_go_string_equal_fold((uintptr_t)s, (uintptr_t)t);
}

fat_String fat_StringTrimPrefix(fat_String s, fat_String prefix) {
  return (fat_String)fatstd_go_string_trim_prefix((uintptr_t)s, (uintptr_t)prefix);
}

fat_String fat_StringTrimSuffix(fat_String s, fat_String suffix) {
  return (fat_String)fatstd_go_string_trim_suffix((uintptr_t)s, (uintptr_t)suffix);
}

bool fat_StringCut(fat_String s, fat_String sep, fat_String *before_out, fat_String *after_out) {
  return (bool)fatstd_go_string_cut((uintptr_t)s, (uintptr_t)sep, (uintptr_t *)before_out,
                                   (uintptr_t *)after_out);
}

bool fat_StringCutPrefix(fat_String s, fat_String prefix, fat_String *after_out) {
  return (bool)fatstd_go_string_cut_prefix((uintptr_t)s, (uintptr_t)prefix, (uintptr_t *)after_out);
}

bool fat_StringCutSuffix(fat_String s, fat_String suffix, fat_String *after_out) {
  return (bool)fatstd_go_string_cut_suffix((uintptr_t)s, (uintptr_t)suffix, (uintptr_t *)after_out);
}

fat_StringArray fat_StringFields(fat_String s) {
  return (fat_StringArray)fatstd_go_string_fields((uintptr_t)s);
}

fat_String fat_StringRepeat(fat_String s, int count) {
  return (fat_String)fatstd_go_string_repeat((uintptr_t)s, count);
}

bool fat_StringContainsAny(fat_String s, fat_String chars) {
  return (bool)fatstd_go_string_contains_any((uintptr_t)s, (uintptr_t)chars);
}

bool fat_StringIndexAny(fat_String s, fat_String chars) {
  return (bool)fatstd_go_string_index_any((uintptr_t)s, (uintptr_t)chars);
}

fat_String fat_StringToValidUTF8(fat_String s, fat_String replacement) {
  return (fat_String)fatstd_go_string_to_valid_utf8((uintptr_t)s, (uintptr_t)replacement);
}

void fat_StringFree(fat_String s) {
  fatstd_go_string_free((uintptr_t)s);
}

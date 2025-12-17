#pragma once

#include <stddef.h>
#include <stdbool.h>

#include "fat/export.h"
#include "fat/handle.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef fat_Handle fat_String;
typedef fat_Handle fat_StringArray;

// Creates a UTF-8 string from a NUL-terminated C string.
FATSTD_API fat_String fat_StringNewUTF8(const char *cstr);

// Creates a UTF-8 string from an explicit byte span.
FATSTD_API fat_String fat_StringNewUTF8N(const char *bytes, size_t len);

FATSTD_API fat_String fat_StringClone(fat_String s);

FATSTD_API bool fat_StringContains(fat_String s, fat_String substr);

FATSTD_API bool fat_StringHasPrefix(fat_String s, fat_String prefix);

FATSTD_API bool fat_StringHasSuffix(fat_String s, fat_String suffix);

FATSTD_API fat_String fat_StringTrimSpace(fat_String s);

FATSTD_API fat_String fat_StringTrim(fat_String s, fat_String cutset);

FATSTD_API fat_StringArray fat_StringSplit(fat_String s, fat_String sep);

FATSTD_API fat_StringArray fat_StringSplitN(fat_String s, fat_String sep, int n);

FATSTD_API size_t fat_StringArrayLen(fat_StringArray a);

// Returns a newly allocated fat_String handle that must be freed with fat_StringFree.
FATSTD_API fat_String fat_StringArrayGet(fat_StringArray a, size_t idx);

FATSTD_API void fat_StringArrayFree(fat_StringArray a);

FATSTD_API fat_String fat_StringJoin(fat_StringArray elems, fat_String sep);

FATSTD_API fat_String fat_StringReplace(fat_String s, fat_String old, fat_String new_value, int n);

FATSTD_API fat_String fat_StringReplaceAll(fat_String s, fat_String old, fat_String new_value);

FATSTD_API fat_String fat_StringToLower(fat_String s);

FATSTD_API fat_String fat_StringToUpper(fat_String s);

// Returns the index of the first instance of `substr` in `s`, or -1 if `substr` is not present in `s`.
FATSTD_API int fat_StringIndex(fat_String s, fat_String substr);

FATSTD_API int fat_StringCount(fat_String s, fat_String substr);

// Returns 0 if a == b, -1 if a < b, and +1 if a > b.
FATSTD_API int fat_StringCompare(fat_String a, fat_String b);

FATSTD_API bool fat_StringEqualFold(fat_String s, fat_String t);

FATSTD_API void fat_StringFree(fat_String s);

#ifdef __cplusplus
} /* extern "C" */
#endif

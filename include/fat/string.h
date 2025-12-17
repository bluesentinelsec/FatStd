#pragma once

#include <stddef.h>
#include <stdbool.h>

#include "fat/export.h"
#include "fat/handle.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef fat_Handle fat_String;

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

FATSTD_API void fat_StringFree(fat_String s);

#ifdef __cplusplus
} /* extern "C" */
#endif

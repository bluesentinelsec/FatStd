#pragma once

#include <stddef.h>
#include <stdint.h>

#include "fat/export.h"
#include "fat/handle.h"
#include "fat/string.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef fat_Handle fat_StringBuilder;

FATSTD_API fat_StringBuilder fat_StringBuilderNew(void);
FATSTD_API void fat_StringBuilderFree(fat_StringBuilder b);

FATSTD_API size_t fat_StringBuilderCap(fat_StringBuilder b);
FATSTD_API size_t fat_StringBuilderLen(fat_StringBuilder b);

FATSTD_API void fat_StringBuilderGrow(fat_StringBuilder b, size_t n);
FATSTD_API void fat_StringBuilderReset(fat_StringBuilder b);

// Returns a newly allocated fat_String handle that must be freed with fat_StringFree.
FATSTD_API fat_String fat_StringBuilderString(fat_StringBuilder b);

// Writes a byte span (supports embedded NULs) and returns number of bytes written.
FATSTD_API size_t fat_StringBuilderWrite(fat_StringBuilder b, const void *bytes, size_t len);

// Writes a single byte.
FATSTD_API void fat_StringBuilderWriteByte(fat_StringBuilder b, uint8_t c);

// Writes a fat_String's contents and returns number of bytes written.
FATSTD_API size_t fat_StringBuilderWriteString(fat_StringBuilder b, fat_String s);

#ifdef __cplusplus
} /* extern "C" */
#endif


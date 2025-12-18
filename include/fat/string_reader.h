#pragma once

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#include "fat/export.h"
#include "fat/handle.h"
#include "fat/string.h"
#include "fat/string_builder.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef fat_Handle fat_StringReader;

FATSTD_API fat_StringReader fat_StringReaderNew(fat_String s);
FATSTD_API void fat_StringReaderFree(fat_StringReader r);

FATSTD_API size_t fat_StringReaderLen(fat_StringReader r);
FATSTD_API int64_t fat_StringReaderSize(fat_StringReader r);

FATSTD_API void fat_StringReaderReset(fat_StringReader r, fat_String s);

// Reads up to `len` bytes into `buf`. Returns number of bytes read.
// If `*eof_out` is set to true, the read hit EOF (possibly after reading some bytes).
FATSTD_API size_t fat_StringReaderRead(fat_StringReader r, void *buf, size_t len, bool *eof_out);

// Like fat_StringReaderRead but at an absolute offset; does not change reader position.
FATSTD_API size_t fat_StringReaderReadAt(
  fat_StringReader r,
  void *buf,
  size_t len,
  int64_t off,
  bool *eof_out
);

// Returns true and writes to `byte_out` if a byte was read; returns false on EOF.
FATSTD_API bool fat_StringReaderReadByte(fat_StringReader r, uint8_t *byte_out, bool *eof_out);

// Unreads the last byte read by fat_StringReaderReadByte. Misuse is fatal.
FATSTD_API void fat_StringReaderUnreadByte(fat_StringReader r);

// Seeks and returns the new absolute position. Misuse (invalid whence/negative position) is fatal.
FATSTD_API int64_t fat_StringReaderSeek(fat_StringReader r, int64_t offset, int whence);

// Writes remaining unread bytes to a FatStd builder and returns bytes written.
// This is the C-friendly equivalent of Go's Reader.WriteTo(io.Writer).
FATSTD_API int64_t fat_StringReaderWriteToBuilder(fat_StringReader r, fat_StringBuilder b);

#ifdef __cplusplus
} /* extern "C" */
#endif


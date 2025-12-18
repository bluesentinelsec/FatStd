#pragma once

/**
 * @file fat/string_reader.h
 * @brief Handle-backed string reader for sequential/random access reads.
 *
 * This is a C-friendly wrapper over Go's strings.Reader.
 *
 * Error handling is C-oriented:
 * - EOF is reported via explicit out-parameters.
 * - Other error conditions (invalid seek/unread, invalid handles) are fail-fast.
 */

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

/**
 * @brief Opaque handle to a FatStd string reader.
 *
 * Readers are mutable and not safe for concurrent use.
 *
 * @note Ownership: free with fat_StringReaderFree.
 */
typedef fat_Handle fat_StringReader;

/**
 * @brief Creates a new reader over a FatStd string.
 *
 * @param s String handle to read from.
 * @return A new fat_StringReader handle (must be freed with fat_StringReaderFree).
 *
 * @note The reader captures a snapshot of the string's bytes at creation time.
 */
FATSTD_API fat_StringReader fat_StringReaderNew(fat_String s);

/**
 * @brief Frees a string reader handle.
 *
 * After this call, the handle is invalid and must not be used.
 *
 * @param r Reader handle to free.
 */
FATSTD_API void fat_StringReaderFree(fat_StringReader r);

/**
 * @brief Returns the number of unread bytes remaining.
 *
 * @param r Reader handle.
 * @return Remaining unread bytes.
 */
FATSTD_API size_t fat_StringReaderLen(fat_StringReader r);

/**
 * @brief Returns the original size of the underlying string in bytes.
 *
 * @param r Reader handle.
 * @return Total size in bytes.
 */
FATSTD_API int64_t fat_StringReaderSize(fat_StringReader r);

/**
 * @brief Resets the reader to read from a new string, positioned at the start.
 *
 * @param r Reader handle.
 * @param s New string handle.
 */
FATSTD_API void fat_StringReaderReset(fat_StringReader r, fat_String s);

/**
 * @brief Reads up to `len` bytes into `buf`.
 *
 * @param r Reader handle.
 * @param buf Output buffer (may be NULL only if len == 0).
 * @param len Capacity of `buf` in bytes.
 * @param eof_out Output: set to true if EOF was hit (possibly after reading some bytes).
 * @return Number of bytes read (0 is valid at EOF).
 *
 * @note Caller owns `buf` and must manage its memory.
 */
FATSTD_API size_t fat_StringReaderRead(fat_StringReader r, void *buf, size_t len, bool *eof_out);

/**
 * @brief Reads at an absolute offset without changing the reader's current position.
 *
 * @param r Reader handle.
 * @param buf Output buffer (may be NULL only if len == 0).
 * @param len Capacity of `buf` in bytes.
 * @param off Absolute byte offset to read from.
 * @param eof_out Output: set to true if EOF was hit (possibly after reading some bytes).
 * @return Number of bytes read.
 */
FATSTD_API size_t fat_StringReaderReadAt(
  fat_StringReader r,
  void *buf,
  size_t len,
  int64_t off,
  bool *eof_out
);

/**
 * @brief Reads a single byte.
 *
 * @param r Reader handle.
 * @param byte_out Output: byte read (only valid if return value is true).
 * @param eof_out Output: set to true if EOF was hit.
 * @return True if a byte was read; false on EOF.
 *
 * @note `byte_out` and `eof_out` must be non-NULL.
 */
FATSTD_API bool fat_StringReaderReadByte(fat_StringReader r, uint8_t *byte_out, bool *eof_out);

/**
 * @brief Unreads the last byte read by fat_StringReaderReadByte.
 *
 * @param r Reader handle.
 *
 * @note Misuse (no prior ReadByte or multiple UnreadByte calls) is fatal.
 */
FATSTD_API void fat_StringReaderUnreadByte(fat_StringReader r);

/**
 * @brief Seeks and returns the new absolute position.
 *
 * `whence` uses the standard SEEK_SET/SEEK_CUR/SEEK_END values.
 *
 * @param r Reader handle.
 * @param offset Offset relative to `whence`.
 * @param whence One of SEEK_SET, SEEK_CUR, SEEK_END.
 * @return New absolute position.
 *
 * @note Misuse (invalid whence or resulting negative position) is fatal.
 */
FATSTD_API int64_t fat_StringReaderSeek(fat_StringReader r, int64_t offset, int whence);

/**
 * @brief Writes remaining unread bytes to a FatStd builder.
 *
 * This is the C-friendly equivalent of Go's Reader.WriteTo(io.Writer), restricted
 * to FatStd builders (arbitrary io.Writer cannot cross the C boundary).
 *
 * @param r Reader handle.
 * @param b Builder handle.
 * @return Number of bytes written.
 */
FATSTD_API int64_t fat_StringReaderWriteToBuilder(fat_StringReader r, fat_StringBuilder b);

#ifdef __cplusplus
} /* extern "C" */
#endif

#pragma once

/**
 * @file fat/bytes_reader.h
 * @brief Handle-backed bytes reader for sequential/random access reads.
 *
 * This is a C-friendly wrapper over Go's bytes.Reader.
 *
 * Error handling is C-oriented:
 * - EOF is reported via explicit out-parameters.
 * - Other error conditions (invalid seek/unread, invalid handles) are fail-fast.
 *
 * io.Writer cannot cross the C boundary; a FatStd-specific WriteTo variant is provided.
 */

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#include "fat/bytes.h"
#include "fat/export.h"
#include "fat/handle.h"
#include "fat/bytes_buffer.h"

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Opaque handle to a FatStd bytes reader.
 *
 * Readers are mutable and not safe for concurrent use.
 *
 * @note Ownership: free with fat_BytesReaderFree.
 */
typedef fat_Handle fat_BytesReader;

/**
 * @brief Creates a new reader over a FatStd bytes slice.
 *
 * @param b Bytes handle to read from.
 * @return A new fat_BytesReader handle (must be freed with fat_BytesReaderFree).
 *
 * @note The reader captures a snapshot of the bytes at creation time.
 */
FATSTD_API fat_BytesReader fat_BytesReaderNew(fat_Bytes b);

/**
 * @brief Frees a bytes reader handle.
 *
 * After this call, the handle is invalid and must not be used.
 *
 * @param r Reader handle to free.
 */
FATSTD_API void fat_BytesReaderFree(fat_BytesReader r);

/**
 * @brief Returns the number of unread bytes remaining.
 *
 * @param r Reader handle.
 * @return Remaining unread bytes.
 */
FATSTD_API size_t fat_BytesReaderLen(fat_BytesReader r);

/**
 * @brief Returns the original size of the underlying bytes slice.
 *
 * @param r Reader handle.
 * @return Total size in bytes.
 */
FATSTD_API int64_t fat_BytesReaderSize(fat_BytesReader r);

/**
 * @brief Resets the reader to read from a new bytes slice, positioned at the start.
 *
 * @param r Reader handle.
 * @param b New bytes handle.
 */
FATSTD_API void fat_BytesReaderReset(fat_BytesReader r, fat_Bytes b);

/**
 * @brief Reads up to `len` bytes into `dst`.
 *
 * @param r Reader handle.
 * @param dst Output buffer (may be NULL only if len == 0).
 * @param len Capacity of `dst` in bytes.
 * @param eof_out Output: set to true if EOF was hit (possibly after reading some bytes).
 * @return Number of bytes read (0 is valid at EOF).
 */
FATSTD_API size_t fat_BytesReaderRead(fat_BytesReader r, void *dst, size_t len, bool *eof_out);

/**
 * @brief Reads at an absolute offset without changing the reader's current position.
 *
 * @param r Reader handle.
 * @param dst Output buffer (may be NULL only if len == 0).
 * @param len Capacity of `dst` in bytes.
 * @param off Absolute byte offset to read from.
 * @param eof_out Output: set to true if EOF was hit (possibly after reading some bytes).
 * @return Number of bytes read.
 */
FATSTD_API size_t fat_BytesReaderReadAt(
  fat_BytesReader r,
  void *dst,
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
FATSTD_API bool fat_BytesReaderReadByte(fat_BytesReader r, uint8_t *byte_out, bool *eof_out);

/**
 * @brief Unreads the last byte read by fat_BytesReaderReadByte.
 *
 * @param r Reader handle.
 *
 * @note Misuse (no prior ReadByte or multiple UnreadByte calls) is fatal.
 */
FATSTD_API void fat_BytesReaderUnreadByte(fat_BytesReader r);

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
FATSTD_API int64_t fat_BytesReaderSeek(fat_BytesReader r, int64_t offset, int whence);

/**
 * @brief Writes remaining unread bytes to a FatStd bytes buffer.
 *
 * This is the C-friendly equivalent of Go's Reader.WriteTo(io.Writer), restricted
 * to FatStd buffers.
 *
 * @param r Reader handle.
 * @param b Buffer handle.
 * @return Number of bytes written.
 */
FATSTD_API int64_t fat_BytesReaderWriteToBytesBuffer(fat_BytesReader r, fat_BytesBuffer b);

#ifdef __cplusplus
} /* extern "C" */
#endif


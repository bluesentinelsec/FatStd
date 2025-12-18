#pragma once

/**
 * @file fat/bytes_buffer.h
 * @brief Handle-backed byte buffer for incremental read/write.
 *
 * This is a C-friendly wrapper over Go's bytes.Buffer.
 *
 * Design notes:
 * - APIs that would expose borrowed slices (Buffer.Bytes, Buffer.Next) return new fat_Bytes handles instead.
 * - APIs that take Go strings (NewBufferString, WriteString) use fat_String handles.
 * - io.Writer/io.Reader cannot cross the C boundary; only FatStd-specific variants are provided.
 *
 * All functions are fail-fast: invalid handles and contract violations are fatal.
 */

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#include "fat/bytes.h"
#include "fat/export.h"
#include "fat/handle.h"
#include "fat/string.h"
#include "fat/string_reader.h"

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Opaque handle to a FatStd bytes buffer.
 *
 * Buffers are mutable and not safe for concurrent use.
 *
 * @note Ownership: free with fat_BytesBufferFree.
 */
typedef fat_Handle fat_BytesBuffer;

/**
 * @brief Creates a new empty bytes buffer.
 *
 * @return A new fat_BytesBuffer handle (must be freed with fat_BytesBufferFree).
 */
FATSTD_API fat_BytesBuffer fat_BytesBufferNew(void);

/**
 * @brief Creates a new bytes buffer initialized with a copy of `b`.
 *
 * @param b Initial bytes handle.
 * @return A new fat_BytesBuffer handle (must be freed with fat_BytesBufferFree).
 */
FATSTD_API fat_BytesBuffer fat_BytesBufferNewBytes(fat_Bytes b);

/**
 * @brief Creates a new bytes buffer initialized from an explicit byte span.
 *
 * Copies exactly `len` bytes; embedded NULs are preserved.
 *
 * @param bytes Pointer to bytes (may be NULL only if len == 0).
 * @param len Number of bytes to copy.
 * @return A new fat_BytesBuffer handle (must be freed with fat_BytesBufferFree).
 */
FATSTD_API fat_BytesBuffer fat_BytesBufferNewN(const void *bytes, size_t len);

/**
 * @brief Creates a new bytes buffer initialized from a string.
 *
 * Matches Go's bytes.NewBufferString (bytes are the string's raw bytes).
 *
 * @param s String handle.
 * @return A new fat_BytesBuffer handle (must be freed with fat_BytesBufferFree).
 */
FATSTD_API fat_BytesBuffer fat_BytesBufferNewString(fat_String s);

/**
 * @brief Frees a bytes buffer handle.
 *
 * After this call, the handle is invalid and must not be used.
 *
 * @param b Buffer handle to free.
 */
FATSTD_API void fat_BytesBufferFree(fat_BytesBuffer b);

/**
 * @brief Returns the number of unread bytes currently in the buffer.
 *
 * @param b Buffer handle.
 * @return Current length in bytes.
 */
FATSTD_API size_t fat_BytesBufferLen(fat_BytesBuffer b);

/**
 * @brief Returns the buffer's current capacity.
 *
 * @param b Buffer handle.
 * @return Current capacity in bytes.
 */
FATSTD_API size_t fat_BytesBufferCap(fat_BytesBuffer b);

/**
 * @brief Ensures the buffer can accept at least `n` additional bytes without reallocating.
 *
 * @param b Buffer handle.
 * @param n Additional capacity to reserve.
 */
FATSTD_API void fat_BytesBufferGrow(fat_BytesBuffer b, size_t n);

/**
 * @brief Resets the buffer to be empty.
 *
 * @param b Buffer handle.
 */
FATSTD_API void fat_BytesBufferReset(fat_BytesBuffer b);

/**
 * @brief Truncates the buffer to exactly `n` bytes.
 *
 * @param b Buffer handle.
 * @param n New length (0 <= n <= fat_BytesBufferLen(b)).
 */
FATSTD_API void fat_BytesBufferTruncate(fat_BytesBuffer b, size_t n);

/**
 * @brief Writes an explicit byte span into the buffer.
 *
 * Copies exactly `len` bytes; embedded NULs are preserved.
 *
 * @param b Buffer handle.
 * @param bytes Pointer to bytes (may be NULL only if len == 0).
 * @param len Number of bytes to write.
 * @return Number of bytes written.
 */
FATSTD_API size_t fat_BytesBufferWrite(fat_BytesBuffer b, const void *bytes, size_t len);

/**
 * @brief Writes a single byte into the buffer.
 *
 * @param b Buffer handle.
 * @param c Byte to write.
 */
FATSTD_API void fat_BytesBufferWriteByte(fat_BytesBuffer b, uint8_t c);

/**
 * @brief Writes a single rune encoded as UTF-8 into the buffer.
 *
 * Matches Go's Buffer.WriteRune.
 *
 * @param b Buffer handle.
 * @param r Unicode code point.
 * @return Number of bytes written (1-4).
 */
FATSTD_API size_t fat_BytesBufferWriteRune(fat_BytesBuffer b, uint32_t r);

/**
 * @brief Writes a string into the buffer.
 *
 * Matches Go's Buffer.WriteString (writes the string's raw bytes).
 *
 * @param b Buffer handle.
 * @param s String handle to write.
 * @return Number of bytes written.
 */
FATSTD_API size_t fat_BytesBufferWriteString(fat_BytesBuffer b, fat_String s);

/**
 * @brief Returns a snapshot of the buffer contents as a new bytes handle.
 *
 * This is the C-friendly equivalent of Go's Buffer.Bytes (which returns a borrowed slice).
 *
 * @param b Buffer handle.
 * @return A new fat_Bytes handle (must be freed with fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_BytesBufferBytes(fat_BytesBuffer b);

/**
 * @brief Returns a snapshot of the buffer contents as a new string handle.
 *
 * Matches Go's Buffer.String.
 *
 * @param b Buffer handle.
 * @return A new fat_String handle (must be freed with fat_StringFree).
 */
FATSTD_API fat_String fat_BytesBufferString(fat_BytesBuffer b);

/**
 * @brief Reads up to `len` bytes into `dst`.
 *
 * @param b Buffer handle.
 * @param dst Output buffer (may be NULL only if len == 0).
 * @param len Capacity of `dst` in bytes.
 * @param eof_out Output: set to true if EOF was hit (possibly after reading some bytes).
 * @return Number of bytes read (0 is valid at EOF).
 */
FATSTD_API size_t fat_BytesBufferRead(fat_BytesBuffer b, void *dst, size_t len, bool *eof_out);

/**
 * @brief Consumes and returns the next `n` bytes as a new bytes handle.
 *
 * This is the C-friendly equivalent of Go's Buffer.Next (which returns a borrowed slice).
 *
 * @param b Buffer handle.
 * @param n Maximum bytes to consume.
 * @return A new fat_Bytes handle (must be freed with fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_BytesBufferNext(fat_BytesBuffer b, size_t n);

/**
 * @brief Reads a single byte.
 *
 * @param b Buffer handle.
 * @param byte_out Output: byte read (only valid if return value is true).
 * @param eof_out Output: set to true if EOF was hit.
 * @return True if a byte was read; false on EOF.
 */
FATSTD_API bool fat_BytesBufferReadByte(fat_BytesBuffer b, uint8_t *byte_out, bool *eof_out);

/**
 * @brief Writes remaining unread bytes from `src` into `dst` and empties `src`.
 *
 * This is the C-friendly equivalent of Go's Buffer.WriteTo(io.Writer), restricted
 * to FatStd buffers.
 *
 * @param src Source buffer handle (drained by this call).
 * @param dst Destination buffer handle.
 * @return Number of bytes written.
 */
FATSTD_API int64_t fat_BytesBufferWriteToBytesBuffer(fat_BytesBuffer src, fat_BytesBuffer dst);

/**
 * @brief Reads all remaining bytes from a FatStd string reader and appends them to the buffer.
 *
 * This is the C-friendly equivalent of Go's Buffer.ReadFrom(io.Reader), restricted
 * to FatStd string readers.
 *
 * @param dst Buffer handle.
 * @param r String reader handle (advanced by this call).
 * @return Number of bytes read.
 */
FATSTD_API int64_t fat_BytesBufferReadFromStringReader(fat_BytesBuffer dst, fat_StringReader r);

#ifdef __cplusplus
} /* extern "C" */
#endif


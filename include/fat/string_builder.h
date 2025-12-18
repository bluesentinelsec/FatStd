#pragma once

/**
 * @file fat/string_builder.h
 * @brief Handle-backed string builder for efficient incremental construction.
 *
 * This is a C-friendly wrapper over Go's strings.Builder. It is useful for building
 * strings without repeated allocations.
 *
 * All functions are fail-fast: invalid handles and contract violations are fatal.
 */

#include <stddef.h>
#include <stdint.h>

#include "fat/export.h"
#include "fat/handle.h"
#include "fat/string.h"

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Opaque handle to a FatStd string builder.
 *
 * Builders are mutable and not safe for concurrent use.
 *
 * @note Ownership: free with fat_StringBuilderFree.
 */
typedef fat_Handle fat_StringBuilder;

/**
 * @brief Creates a new empty string builder.
 *
 * @return A new fat_StringBuilder handle (must be freed with fat_StringBuilderFree).
 */
FATSTD_API fat_StringBuilder fat_StringBuilderNew(void);

/**
 * @brief Frees a string builder handle.
 *
 * After this call, the handle is invalid and must not be used.
 *
 * @param b Builder handle to free.
 */
FATSTD_API void fat_StringBuilderFree(fat_StringBuilder b);

/**
 * @brief Returns the builder's current capacity.
 *
 * @param b Builder handle.
 * @return Current capacity in bytes.
 */
FATSTD_API size_t fat_StringBuilderCap(fat_StringBuilder b);

/**
 * @brief Returns the number of bytes written to the builder.
 *
 * @param b Builder handle.
 * @return Current length in bytes.
 */
FATSTD_API size_t fat_StringBuilderLen(fat_StringBuilder b);

/**
 * @brief Ensures the builder can accept at least `n` additional bytes without reallocating.
 *
 * @param b Builder handle.
 * @param n Additional capacity to reserve.
 */
FATSTD_API void fat_StringBuilderGrow(fat_StringBuilder b, size_t n);

/**
 * @brief Resets the builder to be empty.
 *
 * @param b Builder handle.
 */
FATSTD_API void fat_StringBuilderReset(fat_StringBuilder b);

/**
 * @brief Materializes the builder contents as a new FatStd string handle.
 *
 * @param b Builder handle.
 * @return A new fat_String handle (must be freed with fat_StringFree).
 *
 * @note The builder remains valid and may be used after this call.
 */
FATSTD_API fat_String fat_StringBuilderString(fat_StringBuilder b);

/**
 * @brief Writes an explicit byte span into the builder.
 *
 * Copies exactly `len` bytes; embedded NULs are preserved.
 *
 * @param b Builder handle.
 * @param bytes Pointer to bytes (may be NULL only if len == 0).
 * @param len Number of bytes to write.
 * @return Number of bytes written.
 */
FATSTD_API size_t fat_StringBuilderWrite(fat_StringBuilder b, const void *bytes, size_t len);

/**
 * @brief Writes a single byte into the builder.
 *
 * @param b Builder handle.
 * @param c Byte to write.
 */
FATSTD_API void fat_StringBuilderWriteByte(fat_StringBuilder b, uint8_t c);

/**
 * @brief Writes a FatStd string into the builder.
 *
 * @param b Builder handle.
 * @param s String handle to write.
 * @return Number of bytes written.
 */
FATSTD_API size_t fat_StringBuilderWriteString(fat_StringBuilder b, fat_String s);

#ifdef __cplusplus
} /* extern "C" */
#endif

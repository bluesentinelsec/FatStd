#pragma once

/**
 * @file fat/bytes.h
 * @brief Handle-backed byte-slice utilities.
 *
 * FatStd bytes are opaque handles backed by Go. They represent arbitrary byte
 * sequences (not necessarily UTF-8).
 *
 * All functions are fail-fast: invalid handles and contract violations are fatal.
 */

#include <stdbool.h>
#include <stddef.h>

#include "fat/export.h"
#include "fat/handle.h"
#include "fat/string.h"

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Opaque handle to a FatStd byte slice.
 *
 * The handle is an identity-only token. The underlying storage is Go-managed and
 * is never exposed directly to C.
 *
 * @note Ownership: free with fat_BytesFree.
 */
typedef fat_Handle fat_Bytes;

/**
 * @brief Opaque handle to an array of byte slices.
 *
 * Arrays are returned by APIs like fat_BytesSplit. Elements are accessed with
 * fat_BytesArrayGet (which returns newly allocated fat_Bytes handles).
 *
 * @note Ownership: free with fat_BytesArrayFree.
 */
typedef fat_Handle fat_BytesArray;

/**
 * @brief Creates a byte slice from an explicit byte span.
 *
 * Copies exactly `len` bytes; embedded NULs are preserved.
 *
 * @param bytes Pointer to bytes (may be NULL only if len == 0).
 * @param len Number of bytes to copy.
 * @return A new fat_Bytes handle (must be freed with fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_BytesNewN(const void *bytes, size_t len);

/**
 * @brief Returns the number of bytes in the slice.
 *
 * @param b Bytes handle.
 * @return Byte length.
 */
FATSTD_API size_t fat_BytesLen(fat_Bytes b);

/**
 * @brief Copies bytes out of a FatStd byte slice into a caller buffer.
 *
 * Copies up to `dst_len` bytes into `dst` and returns the number of bytes copied.
 *
 * @param b Bytes handle.
 * @param dst Destination buffer (may be NULL only if dst_len == 0).
 * @param dst_len Capacity of `dst` in bytes.
 * @return Number of bytes copied (may be less than fat_BytesLen).
 *
 * @note Caller owns `dst` and must manage its memory.
 */
FATSTD_API size_t fat_BytesCopyOut(fat_Bytes b, void *dst, size_t dst_len);

/**
 * @brief Returns a new byte slice handle with a cloned copy of `b`.
 *
 * @param b Bytes handle.
 * @return A new fat_Bytes handle (must be freed with fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_BytesClone(fat_Bytes b);

/**
 * @brief Reports whether `subslice` is within `b`.
 *
 * @param b Bytes handle.
 * @param subslice Subslice handle.
 * @return True if found; false otherwise.
 */
FATSTD_API bool fat_BytesContains(fat_Bytes b, fat_Bytes subslice);

/**
 * @brief Reports whether `s` begins with `prefix`.
 *
 * @param s Bytes handle.
 * @param prefix Prefix handle.
 * @return True if `s` begins with `prefix`; false otherwise.
 */
FATSTD_API bool fat_BytesHasPrefix(fat_Bytes s, fat_Bytes prefix);

/**
 * @brief Reports whether `s` ends with `suffix`.
 *
 * @param s Bytes handle.
 * @param suffix Suffix handle.
 * @return True if `s` ends with `suffix`; false otherwise.
 */
FATSTD_API bool fat_BytesHasSuffix(fat_Bytes s, fat_Bytes suffix);

/**
 * @brief Returns a new byte slice with leading/trailing ASCII whitespace removed.
 *
 * Matches Go's bytes.TrimSpace.
 *
 * @param s Bytes handle.
 * @return A new fat_Bytes handle (must be freed with fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_BytesTrimSpace(fat_Bytes s);

/**
 * @brief Returns a new byte slice with all leading/trailing bytes in `cutset` removed.
 *
 * Matches Go's bytes.Trim semantics (cutset is treated as a set of bytes derived from a string).
 *
 * @param s Bytes handle.
 * @param cutset Cutset string handle.
 * @return A new fat_Bytes handle (must be freed with fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_BytesTrim(fat_Bytes s, fat_String cutset);

/**
 * @brief Splits `s` around each instance of `sep`.
 *
 * @param s Bytes handle.
 * @param sep Separator bytes handle.
 * @return A new fat_BytesArray handle (must be freed with fat_BytesArrayFree).
 */
FATSTD_API fat_BytesArray fat_BytesSplit(fat_Bytes s, fat_Bytes sep);

/**
 * @brief Returns the number of elements in a bytes array.
 *
 * @param a Array handle.
 * @return Number of elements.
 */
FATSTD_API size_t fat_BytesArrayLen(fat_BytesArray a);

/**
 * @brief Returns the element at index `idx` as a new bytes handle.
 *
 * The returned handle is newly allocated and must be freed by the caller.
 *
 * @param a Array handle.
 * @param idx Element index (0 <= idx < fat_BytesArrayLen(a)).
 * @return A new fat_Bytes handle (must be freed with fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_BytesArrayGet(fat_BytesArray a, size_t idx);

/**
 * @brief Frees a bytes array handle.
 *
 * @param a Array handle to free.
 *
 * @note This does not free any fat_Bytes handles previously returned by fat_BytesArrayGet.
 */
FATSTD_API void fat_BytesArrayFree(fat_BytesArray a);

/**
 * @brief Joins an array of byte slices using `sep`.
 *
 * @param s Array handle.
 * @param sep Separator bytes handle.
 * @return A new fat_Bytes handle (must be freed with fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_BytesJoin(fat_BytesArray s, fat_Bytes sep);

/**
 * @brief Returns a new byte slice with all instances of `old_value` replaced by `new_value`.
 *
 * @param s Bytes handle.
 * @param old_value Subslice to replace.
 * @param new_value Replacement subslice.
 * @return A new fat_Bytes handle (must be freed with fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_BytesReplaceAll(fat_Bytes s, fat_Bytes old_value, fat_Bytes new_value);

/**
 * @brief Returns a new byte slice with up to `n` instances of `old_value` replaced by `new_value`.
 *
 * Matches Go's bytes.Replace semantics (n < 0 replaces all).
 *
 * @param s Bytes handle.
 * @param old_value Subslice to replace.
 * @param new_value Replacement subslice.
 * @param n Maximum replacements (n < 0 means replace all).
 * @return A new fat_Bytes handle (must be freed with fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_BytesReplace(fat_Bytes s, fat_Bytes old_value, fat_Bytes new_value, int n);

/**
 * @brief Returns a new byte slice with ASCII letters mapped to their lower case.
 *
 * Matches Go's bytes.ToLower (ASCII-only).
 *
 * @param s Bytes handle.
 * @return A new fat_Bytes handle (must be freed with fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_BytesToLower(fat_Bytes s);

/**
 * @brief Returns a new byte slice with ASCII letters mapped to their upper case.
 *
 * Matches Go's bytes.ToUpper (ASCII-only).
 *
 * @param s Bytes handle.
 * @return A new fat_Bytes handle (must be freed with fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_BytesToUpper(fat_Bytes s);

/**
 * @brief Returns the index of the first instance of `sep` in `s`.
 *
 * @param s Bytes handle.
 * @param sep Subslice handle.
 * @return Zero-based byte index, or -1 if not found.
 */
FATSTD_API int fat_BytesIndex(fat_Bytes s, fat_Bytes sep);

/**
 * @brief Counts the number of non-overlapping instances of `sep` in `s`.
 *
 * @param s Bytes handle.
 * @param sep Subslice handle.
 * @return Count of occurrences.
 */
FATSTD_API int fat_BytesCount(fat_Bytes s, fat_Bytes sep);

/**
 * @brief Lexicographically compares two byte slices.
 *
 * @param a Bytes handle.
 * @param b Bytes handle.
 * @return 0 if equal, -1 if a < b, +1 if a > b.
 */
FATSTD_API int fat_BytesCompare(fat_Bytes a, fat_Bytes b);

/**
 * @brief Reports whether two byte slices are equal.
 *
 * @param a Bytes handle.
 * @param b Bytes handle.
 * @return True if equal; false otherwise.
 */
FATSTD_API bool fat_BytesEqual(fat_Bytes a, fat_Bytes b);

/**
 * @brief Frees a FatStd bytes handle.
 *
 * After this call, the handle is invalid and must not be used.
 *
 * @param b Bytes handle to free.
 */
FATSTD_API void fat_BytesFree(fat_Bytes b);

#ifdef __cplusplus
} /* extern "C" */
#endif


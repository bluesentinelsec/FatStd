#pragma once

/**
 * @file fat/string.h
 * @brief Handle-backed UTF-8/byte-string utilities.
 *
 * FatStd strings are opaque handles backed by Go. They may contain arbitrary bytes,
 * including embedded NULs, especially when constructed via length-based APIs.
 *
 * All functions are fail-fast: invalid handles and contract violations are fatal.
 */

#include <stddef.h>
#include <stdbool.h>

#include "fat/export.h"
#include "fat/handle.h"

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Opaque handle to a FatStd string.
 *
 * The handle is an identity-only token. The underlying storage is Go-managed and
 * is never exposed directly to C.
 *
 * @note Ownership: free with fat_StringFree.
 */
typedef fat_Handle fat_String;

/**
 * @brief Opaque handle to a FatStd array of strings.
 *
 * Arrays are returned by APIs like split/fields. Elements are accessed with
 * fat_StringArrayGet (which returns newly allocated fat_String handles).
 *
 * @note Ownership: free the array handle with fat_StringArrayFree.
 */
typedef fat_Handle fat_StringArray;

/**
 * @brief Creates a string from a NUL-terminated C string.
 *
 * The input is read until the first NUL byte.
 *
 * @param cstr NUL-terminated UTF-8 string; must be non-NULL.
 * @return A new fat_String handle (must be freed with fat_StringFree).
 *
 * @note Embedded NULs are not representable via this constructor.
 */
FATSTD_API fat_String fat_StringNewUTF8(const char *cstr);

/**
 * @brief Creates a string from an explicit byte span.
 *
 * Copies exactly `len` bytes; embedded NULs are preserved.
 *
 * @param bytes Pointer to bytes (may be NULL only if len == 0).
 * @param len Number of bytes to copy.
 * @return A new fat_String handle (must be freed with fat_StringFree).
 */
FATSTD_API fat_String fat_StringNewUTF8N(const char *bytes, size_t len);

/**
 * @brief Returns the number of bytes in the string.
 *
 * This is a byte length, not rune/character count, and may include embedded NULs.
 *
 * @param s String handle.
 * @return Byte length of `s`.
 */
FATSTD_API size_t fat_StringLenBytes(fat_String s);

/**
 * @brief Copies bytes out of a FatStd string into a caller buffer.
 *
 * Copies up to `dst_len` bytes into `dst` and returns the number of bytes copied.
 * This API does not add a NUL terminator.
 *
 * @param s String handle.
 * @param dst Destination buffer (may be NULL only if dst_len == 0).
 * @param dst_len Capacity of `dst` in bytes.
 * @return Number of bytes copied (may be less than fat_StringLenBytes).
 *
 * @note Caller owns `dst` and must manage its memory.
 */
FATSTD_API size_t fat_StringCopyOut(fat_String s, void *dst, size_t dst_len);

/**
 * @brief Returns a new string handle whose bytes are a clone of `s`.
 *
 * @param s String handle.
 * @return A new fat_String handle (must be freed with fat_StringFree).
 */
FATSTD_API fat_String fat_StringClone(fat_String s);

/**
 * @brief Reports whether `substr` is within `s`.
 *
 * @param s String handle.
 * @param substr Substring handle.
 * @return True if found; false otherwise.
 */
FATSTD_API bool fat_StringContains(fat_String s, fat_String substr);

/**
 * @brief Reports whether `s` begins with `prefix`.
 *
 * @param s String handle.
 * @param prefix Prefix handle.
 * @return True if `s` begins with `prefix`; false otherwise.
 */
FATSTD_API bool fat_StringHasPrefix(fat_String s, fat_String prefix);

/**
 * @brief Reports whether `s` ends with `suffix`.
 *
 * @param s String handle.
 * @param suffix Suffix handle.
 * @return True if `s` ends with `suffix`; false otherwise.
 */
FATSTD_API bool fat_StringHasSuffix(fat_String s, fat_String suffix);

/**
 * @brief Returns a new string handle with a copy of `s` trimmed of leading/trailing Unicode whitespace.
 *
 * @param s String handle.
 * @return A new fat_String handle (must be freed with fat_StringFree).
 */
FATSTD_API fat_String fat_StringTrimSpace(fat_String s);

/**
 * @brief Returns a new string handle with all leading/trailing bytes in `cutset` removed.
 *
 * Matches Go's strings.Trim semantics (cutset is treated as a set of Unicode code points).
 *
 * @param s String handle.
 * @param cutset Cutset handle.
 * @return A new fat_String handle (must be freed with fat_StringFree).
 */
FATSTD_API fat_String fat_StringTrim(fat_String s, fat_String cutset);

/**
 * @brief Splits `s` around each instance of `sep`.
 *
 * @param s String handle.
 * @param sep Separator handle.
 * @return A new fat_StringArray handle (must be freed with fat_StringArrayFree).
 */
FATSTD_API fat_StringArray fat_StringSplit(fat_String s, fat_String sep);

/**
 * @brief Splits `s` into at most `n` substrings separated by `sep`.
 *
 * Matches Go's strings.SplitN:
 * - n > 0: at most n substrings
 * - n == 0: empty result
 * - n < 0: all substrings
 *
 * @param s String handle.
 * @param sep Separator handle.
 * @param n Maximum substrings, per semantics above.
 * @return A new fat_StringArray handle (must be freed with fat_StringArrayFree).
 */
FATSTD_API fat_StringArray fat_StringSplitN(fat_String s, fat_String sep, int n);

/**
 * @brief Returns the number of elements in a string array.
 *
 * @param a Array handle.
 * @return Number of elements.
 */
FATSTD_API size_t fat_StringArrayLen(fat_StringArray a);

/**
 * @brief Returns the element at index `idx` as a new string handle.
 *
 * The returned handle is a newly allocated fat_String and must be freed by the caller.
 *
 * @param a Array handle.
 * @param idx Element index (0 <= idx < fat_StringArrayLen(a)).
 * @return A new fat_String handle (must be freed with fat_StringFree).
 */
FATSTD_API fat_String fat_StringArrayGet(fat_StringArray a, size_t idx);

/**
 * @brief Frees a string array handle.
 *
 * @param a Array handle to free.
 *
 * @note This does not free any fat_String handles previously returned by fat_StringArrayGet.
 */
FATSTD_API void fat_StringArrayFree(fat_StringArray a);

/**
 * @brief Joins an array of strings using `sep`.
 *
 * @param elems Array handle.
 * @param sep Separator handle.
 * @return A new fat_String handle (must be freed with fat_StringFree).
 */
FATSTD_API fat_String fat_StringJoin(fat_StringArray elems, fat_String sep);

/**
 * @brief Returns a new string with up to `n` non-overlapping instances of `old` replaced by `new_value`.
 *
 * Matches Go's strings.Replace semantics (n < 0 replaces all).
 *
 * @param s String handle.
 * @param old Substring to replace.
 * @param new_value Replacement substring.
 * @param n Maximum replacements (n < 0 means replace all).
 * @return A new fat_String handle (must be freed with fat_StringFree).
 */
FATSTD_API fat_String fat_StringReplace(fat_String s, fat_String old, fat_String new_value, int n);

/**
 * @brief Returns a new string with all instances of `old` replaced by `new_value`.
 *
 * @param s String handle.
 * @param old Substring to replace.
 * @param new_value Replacement substring.
 * @return A new fat_String handle (must be freed with fat_StringFree).
 */
FATSTD_API fat_String fat_StringReplaceAll(fat_String s, fat_String old, fat_String new_value);

/**
 * @brief Returns a new string with Unicode letters mapped to their lower case.
 *
 * @param s String handle.
 * @return A new fat_String handle (must be freed with fat_StringFree).
 */
FATSTD_API fat_String fat_StringToLower(fat_String s);

/**
 * @brief Returns a new string with Unicode letters mapped to their upper case.
 *
 * @param s String handle.
 * @return A new fat_String handle (must be freed with fat_StringFree).
 */
FATSTD_API fat_String fat_StringToUpper(fat_String s);

/**
 * @brief Returns the index of the first instance of `substr` in `s`.
 *
 * @param s String handle.
 * @param substr Substring handle.
 * @return Zero-based byte index, or -1 if not found.
 */
FATSTD_API int fat_StringIndex(fat_String s, fat_String substr);

/**
 * @brief Counts the number of non-overlapping instances of `substr` in `s`.
 *
 * @param s String handle.
 * @param substr Substring handle.
 * @return Count of occurrences.
 */
FATSTD_API int fat_StringCount(fat_String s, fat_String substr);

/**
 * @brief Lexicographically compares two strings.
 *
 * @param a String handle.
 * @param b String handle.
 * @return 0 if a == b, -1 if a < b, and +1 if a > b.
 */
FATSTD_API int fat_StringCompare(fat_String a, fat_String b);

/**
 * @brief Reports whether `s` and `t` are equal under Unicode case-folding.
 *
 * @param s String handle.
 * @param t String handle.
 * @return True if equal under case-folding; false otherwise.
 */
FATSTD_API bool fat_StringEqualFold(fat_String s, fat_String t);

/**
 * @brief Returns a new string with the leading `prefix` removed, if present.
 *
 * @param s String handle.
 * @param prefix Prefix handle.
 * @return A new fat_String handle (must be freed with fat_StringFree).
 */
FATSTD_API fat_String fat_StringTrimPrefix(fat_String s, fat_String prefix);

/**
 * @brief Returns a new string with the trailing `suffix` removed, if present.
 *
 * @param s String handle.
 * @param suffix Suffix handle.
 * @return A new fat_String handle (must be freed with fat_StringFree).
 */
FATSTD_API fat_String fat_StringTrimSuffix(fat_String s, fat_String suffix);

/**
 * @brief Splits `s` around the first instance of `sep`.
 *
 * Matches Go's strings.Cut semantics. Outputs are newly allocated fat_String handles
 * regardless of whether `sep` is found.
 *
 * @param s String handle.
 * @param sep Separator handle.
 * @param before_out Output: substring before `sep` (new handle; caller frees).
 * @param after_out Output: substring after `sep` (new handle; caller frees).
 * @return True if `sep` was found; false otherwise.
 *
 * @note `before_out` and `after_out` must be non-NULL.
 */
FATSTD_API bool fat_StringCut(fat_String s, fat_String sep, fat_String *before_out, fat_String *after_out);

/**
 * @brief Cuts `prefix` from the start of `s`.
 *
 * Matches Go's strings.CutPrefix semantics. `after_out` is a newly allocated fat_String handle
 * regardless of whether `prefix` is found.
 *
 * @param s String handle.
 * @param prefix Prefix handle.
 * @param after_out Output: substring after `prefix` (new handle; caller frees).
 * @return True if `prefix` was found; false otherwise.
 *
 * @note `after_out` must be non-NULL.
 */
FATSTD_API bool fat_StringCutPrefix(fat_String s, fat_String prefix, fat_String *after_out);

/**
 * @brief Cuts `suffix` from the end of `s`.
 *
 * Matches Go's strings.CutSuffix semantics. `after_out` is a newly allocated fat_String handle
 * regardless of whether `suffix` is found.
 *
 * @param s String handle.
 * @param suffix Suffix handle.
 * @param after_out Output: substring with trailing suffix removed (new handle; caller frees).
 * @return True if `suffix` was found; false otherwise.
 *
 * @note `after_out` must be non-NULL.
 */
FATSTD_API bool fat_StringCutSuffix(fat_String s, fat_String suffix, fat_String *after_out);

/**
 * @brief Splits `s` around runs of Unicode whitespace.
 *
 * Matches Go's strings.Fields.
 *
 * @param s String handle.
 * @return A new fat_StringArray handle (must be freed with fat_StringArrayFree).
 */
FATSTD_API fat_StringArray fat_StringFields(fat_String s);

/**
 * @brief Returns a new string consisting of `count` copies of `s` concatenated.
 *
 * @param s String handle.
 * @param count Number of repetitions (must be >= 0; misuse is fatal).
 * @return A new fat_String handle (must be freed with fat_StringFree).
 */
FATSTD_API fat_String fat_StringRepeat(fat_String s, int count);

/**
 * @brief Reports whether any Unicode code points in `chars` are within `s`.
 *
 * @param s String handle.
 * @param chars Set of characters to search for.
 * @return True if any character is found; false otherwise.
 */
FATSTD_API bool fat_StringContainsAny(fat_String s, fat_String chars);

/**
 * @brief Reports whether any Unicode code points in `chars` are within `s`.
 *
 * This is a boolean wrapper around Go's strings.IndexAny.
 *
 * @param s String handle.
 * @param chars Set of characters to search for.
 * @return True if any character is found; false otherwise.
 */
FATSTD_API bool fat_StringIndexAny(fat_String s, fat_String chars);

/**
 * @brief Returns a copy of `s` with invalid UTF-8 sequences replaced.
 *
 * Matches Go's strings.ToValidUTF8.
 *
 * @param s String handle.
 * @param replacement Replacement string.
 * @return A new fat_String handle (must be freed with fat_StringFree).
 */
FATSTD_API fat_String fat_StringToValidUTF8(fat_String s, fat_String replacement);

/**
 * @brief Frees a FatStd string handle.
 *
 * After this call, the handle is invalid and must not be used.
 *
 * @param s String handle to free.
 */
FATSTD_API void fat_StringFree(fat_String s);

#ifdef __cplusplus
} /* extern "C" */
#endif

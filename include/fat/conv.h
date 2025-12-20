#pragma once

/**
 * @file fat/conv.h
 * @brief Handle-backed numeric/string formatting and parsing utilities.
 *
 * This module is a C-friendly wrapper around Go's strconv-like functionality.
 *
 * Design notes:
 * - Append* functions return a new fat_Bytes handle rather than mutating a borrowed slice.
 * - Parsing functions return fat_Status and optionally populate a fat_Error handle.
 * - Contract violations (invalid handles, NULL out params) are fatal; parse failures are recoverable.
 */

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#include "fat/bytes.h"
#include "fat/error.h"
#include "fat/export.h"
#include "fat/status.h"
#include "fat/string.h"

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Returns the bit size of Go's `int` on this build (32 or 64).
 *
 * @return 32 or 64.
 */
FATSTD_API int fat_ConvIntSize(void);

/**
 * @brief Appends a boolean to dst and returns a new bytes handle.
 *
 * @param dst Destination bytes handle.
 * @param b Boolean value.
 * @return A new fat_Bytes handle (caller must fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_ConvAppendBool(fat_Bytes dst, bool b);

/**
 * @brief Appends an integer to dst and returns a new bytes handle.
 *
 * @param dst Destination bytes handle.
 * @param i Integer value.
 * @param base Formatting base (2..36).
 * @return A new fat_Bytes handle (caller must fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_ConvAppendInt(fat_Bytes dst, int64_t i, int base);

/**
 * @brief Appends an unsigned integer to dst and returns a new bytes handle.
 *
 * @param dst Destination bytes handle.
 * @param i Unsigned integer value.
 * @param base Formatting base (2..36).
 * @return A new fat_Bytes handle (caller must fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_ConvAppendUint(fat_Bytes dst, uint64_t i, int base);

/**
 * @brief Appends a float to dst and returns a new bytes handle.
 *
 * @param dst Destination bytes handle.
 * @param f Float value.
 * @param fmt Format verb (e.g. 'f', 'g', 'e').
 * @param prec Precision.
 * @param bit_size 32 or 64.
 * @return A new fat_Bytes handle (caller must fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_ConvAppendFloat(fat_Bytes dst, double f, uint8_t fmt, int prec, int bit_size);

/**
 * @brief Appends a quoted string to dst and returns a new bytes handle.
 *
 * @param dst Destination bytes handle.
 * @param s String handle.
 * @return A new fat_Bytes handle (caller must fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_ConvAppendQuote(fat_Bytes dst, fat_String s);

/**
 * @brief Appends a quoted rune to dst and returns a new bytes handle.
 *
 * @param dst Destination bytes handle.
 * @param r Unicode code point.
 * @return A new fat_Bytes handle (caller must fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_ConvAppendQuoteRune(fat_Bytes dst, uint32_t r);

FATSTD_API fat_Bytes fat_ConvAppendQuoteRuneToASCII(fat_Bytes dst, uint32_t r);
FATSTD_API fat_Bytes fat_ConvAppendQuoteRuneToGraphic(fat_Bytes dst, uint32_t r);
FATSTD_API fat_Bytes fat_ConvAppendQuoteToASCII(fat_Bytes dst, fat_String s);
FATSTD_API fat_Bytes fat_ConvAppendQuoteToGraphic(fat_Bytes dst, fat_String s);

/**
 * @brief Reports whether s can be represented as a single backquoted string literal.
 *
 * @param s String handle.
 * @return True if backquoting is valid; false otherwise.
 */
FATSTD_API bool fat_ConvCanBackquote(fat_String s);

FATSTD_API fat_String fat_ConvFormatBool(bool b);
FATSTD_API fat_String fat_ConvFormatInt(int64_t i, int base);
FATSTD_API fat_String fat_ConvFormatUint(uint64_t i, int base);
FATSTD_API fat_String fat_ConvFormatFloat(double f, uint8_t fmt, int prec, int bit_size);

/**
 * @brief Formats a complex value.
 *
 * @param re Real component.
 * @param im Imaginary component.
 * @param fmt Format verb.
 * @param prec Precision.
 * @param bit_size 64 or 128.
 * @return A new fat_String handle (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_ConvFormatComplex(double re, double im, uint8_t fmt, int prec, int bit_size);

FATSTD_API bool fat_ConvIsGraphic(uint32_t r);
FATSTD_API bool fat_ConvIsPrint(uint32_t r);
FATSTD_API fat_String fat_ConvItoa(int i);

FATSTD_API fat_String fat_ConvQuote(fat_String s);
FATSTD_API fat_String fat_ConvQuoteRune(uint32_t r);
FATSTD_API fat_String fat_ConvQuoteRuneToASCII(uint32_t r);
FATSTD_API fat_String fat_ConvQuoteRuneToGraphic(uint32_t r);
FATSTD_API fat_String fat_ConvQuoteToASCII(fat_String s);
FATSTD_API fat_String fat_ConvQuoteToGraphic(fat_String s);

/**
 * @brief Parses a boolean.
 *
 * @param s Input string.
 * @param out_value Output: 0 or 1 on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_* on parse failure.
 */
FATSTD_API fat_Status fat_ConvParseBool(fat_String s, int *out_value, fat_Error *out_err);

FATSTD_API fat_Status fat_ConvParseInt(fat_String s, int base, int bit_size, int64_t *out_value, fat_Error *out_err);
FATSTD_API fat_Status fat_ConvParseUint(fat_String s, int base, int bit_size, uint64_t *out_value, fat_Error *out_err);
FATSTD_API fat_Status fat_ConvParseFloat(fat_String s, int bit_size, double *out_value, fat_Error *out_err);

/**
 * @brief Parses a complex value.
 *
 * @param s Input string.
 * @param bit_size 64 or 128.
 * @param out_re Output: real component.
 * @param out_im Output: imaginary component.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_* on parse failure.
 */
FATSTD_API fat_Status fat_ConvParseComplex(fat_String s, int bit_size, double *out_re, double *out_im, fat_Error *out_err);

/**
 * @brief Converts a decimal string to a Go int, returned as int64.
 *
 * @param s Input string.
 * @param out_value Output: parsed value.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_* on parse failure.
 *
 * @note The range matches Go's `int` for the current build; use ParseInt for fixed bit sizes.
 */
FATSTD_API fat_Status fat_ConvAtoi(fat_String s, int64_t *out_value, fat_Error *out_err);

FATSTD_API fat_Status fat_ConvUnquote(fat_String s, fat_String *out_value, fat_Error *out_err);
FATSTD_API fat_Status fat_ConvQuotedPrefix(fat_String s, fat_String *out_value, fat_Error *out_err);

/**
 * @brief Unquotes the first character or byte in a string.
 *
 * @param s Input string.
 * @param quote Quote byte (0, '\'', '"', or '`').
 * @param out_rune Output: decoded rune.
 * @param out_multibyte Output: true if multibyte UTF-8 sequence.
 * @param out_tail Output: remaining unconsumed tail (new fat_String; caller frees).
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_* on parse failure.
 */
FATSTD_API fat_Status fat_ConvUnquoteChar(
  fat_String s,
  uint8_t quote,
  uint32_t *out_rune,
  bool *out_multibyte,
  fat_String *out_tail,
  fat_Error *out_err
);

#ifdef __cplusplus
} /* extern "C" */
#endif


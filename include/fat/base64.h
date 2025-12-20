#pragma once

/**
 * @file fat/base64.h
 * @brief Base64 encoding/decoding utilities.
 *
 * This module is backed by Go's encoding/base64 package.
 *
 * Design notes:
 * - Go's io.Reader/io.Writer streaming types do not map cleanly to C. The C API
 *   focuses on explicit encode/decode operations on FatStd byte/string handles.
 * - Recoverable failures (invalid base64 input, invalid encoding alphabet) return
 *   fat_Status and fat_Error.
 *
 * Contract violations (invalid handles, NULL out-params, NULL pointers where forbidden)
 * are fatal.
 */

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#include "fat/bytes.h"
#include "fat/bytes_buffer.h"
#include "fat/error.h"
#include "fat/export.h"
#include "fat/handle.h"
#include "fat/status.h"
#include "fat/string.h"

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Opaque handle to a base64 encoding configuration.
 *
 * @note Ownership: free with fat_Base64EncodingFree.
 */
typedef fat_Handle fat_Base64Encoding;

/**
 * @brief Opaque handle to a streaming base64 encoder.
 *
 * This is a C-friendly alternative to Go's base64.NewEncoder that writes to a
 * FatStd bytes buffer.
 *
 * @note Ownership: close/free with fat_Base64EncoderClose.
 */
typedef fat_Handle fat_Base64Encoder;

/**
 * @brief Creates a new base64 encoding from a 64-byte alphabet string (UTF-8 bytes).
 *
 * This matches Go's base64.NewEncoding. The input must be exactly 64 bytes.
 *
 * @param alphabet 64-byte alphabet string (NUL-terminated).
 * @param out_enc Output: new encoding handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_RANGE on invalid alphabet length.
 */
FATSTD_API fat_Status fat_Base64EncodingNewUTF8(const char *alphabet, fat_Base64Encoding *out_enc, fat_Error *out_err);

/**
 * @brief Creates a strict variant of an encoding.
 *
 * @param enc Encoding handle.
 * @return New encoding handle (caller must fat_Base64EncodingFree).
 */
FATSTD_API fat_Base64Encoding fat_Base64EncodingStrict(fat_Base64Encoding enc);

/**
 * @brief Creates a variant of an encoding with custom padding.
 *
 * Pass -1 (base64.NoPadding) to disable padding.
 *
 * @param enc Encoding handle.
 * @param padding_rune Padding rune (e.g. '='), or -1 for no padding.
 * @param out_enc Output: new encoding handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; non-OK on invalid padding.
 */
FATSTD_API fat_Status fat_Base64EncodingWithPadding(
  fat_Base64Encoding enc,
  int32_t padding_rune,
  fat_Base64Encoding *out_enc,
  fat_Error *out_err
);

/**
 * @brief Returns the encoded length for an input of length n.
 *
 * @param enc Encoding handle.
 * @param n Input length in bytes.
 * @return Encoded length.
 */
FATSTD_API int fat_Base64EncodedLen(fat_Base64Encoding enc, int n);

/**
 * @brief Returns the maximum decoded length for an input of length n.
 *
 * @param enc Encoding handle.
 * @param n Encoded length in bytes.
 * @return Maximum decoded length.
 */
FATSTD_API int fat_Base64DecodedLen(fat_Base64Encoding enc, int n);

/**
 * @brief Encodes src and returns the base64 output as a new string.
 *
 * @param enc Encoding handle.
 * @param src Source bytes.
 * @return New fat_String handle (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_Base64EncodeToString(fat_Base64Encoding enc, fat_Bytes src);

/**
 * @brief Encodes src and returns the base64 output as new bytes.
 *
 * @param enc Encoding handle.
 * @param src Source bytes.
 * @return New fat_Bytes handle (caller must fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_Base64Encode(fat_Base64Encoding enc, fat_Bytes src);

/**
 * @brief Decodes a base64 string.
 *
 * @param enc Encoding handle.
 * @param s Base64 input string.
 * @param out Decoded bytes on success (caller must fat_BytesFree).
 * @param out_err Error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_SYNTAX on corrupt input; FAT_ERR_OTHER otherwise.
 */
FATSTD_API fat_Status fat_Base64DecodeString(fat_Base64Encoding enc, fat_String s, fat_Bytes *out, fat_Error *out_err);

/**
 * @brief Decodes base64 bytes.
 *
 * @param enc Encoding handle.
 * @param src Base64 input bytes.
 * @param out Decoded bytes on success (caller must fat_BytesFree).
 * @param out_err Error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_SYNTAX on corrupt input; FAT_ERR_OTHER otherwise.
 */
FATSTD_API fat_Status fat_Base64Decode(fat_Base64Encoding enc, fat_Bytes src, fat_Bytes *out, fat_Error *out_err);

/**
 * @brief Appends the base64 encoding of src to dst and returns new bytes.
 *
 * @param enc Encoding handle.
 * @param dst Destination bytes.
 * @param src Source bytes.
 * @return New fat_Bytes handle (caller must fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_Base64AppendEncode(fat_Base64Encoding enc, fat_Bytes dst, fat_Bytes src);

/**
 * @brief Appends the decoded bytes of src to dst.
 *
 * @param enc Encoding handle.
 * @param dst Destination bytes.
 * @param src Base64 input bytes.
 * @param out Output: new fat_Bytes handle on success (caller must fat_BytesFree).
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_SYNTAX on corrupt input; FAT_ERR_OTHER otherwise.
 */
FATSTD_API fat_Status fat_Base64AppendDecode(
  fat_Base64Encoding enc,
  fat_Bytes dst,
  fat_Bytes src,
  fat_Bytes *out,
  fat_Error *out_err
);

/**
 * @brief Frees a base64 encoding handle.
 *
 * @param enc Encoding handle to free.
 */
FATSTD_API void fat_Base64EncodingFree(fat_Base64Encoding enc);

/**
 * @brief Creates a streaming base64 encoder writing to a bytes buffer.
 *
 * @param enc Encoding handle.
 * @param dst Destination bytes buffer (receives encoded bytes).
 * @param out_encoder Output: new encoder handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; non-OK on failure.
 */
FATSTD_API fat_Status fat_Base64EncoderNewToBytesBuffer(
  fat_Base64Encoding enc,
  fat_BytesBuffer dst,
  fat_Base64Encoder *out_encoder,
  fat_Error *out_err
);

/**
 * @brief Writes bytes into the base64 encoder.
 *
 * @param e Encoder handle.
 * @param bytes Input bytes (may be NULL only if len == 0).
 * @param len Input length in bytes.
 * @param out_n Output: number of bytes consumed from input.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; non-OK on failure.
 */
FATSTD_API fat_Status fat_Base64EncoderWrite(
  fat_Base64Encoder e,
  const void *bytes,
  size_t len,
  size_t *out_n,
  fat_Error *out_err
);

/**
 * @brief Closes and frees a base64 encoder handle.
 *
 * @param e Encoder handle.
 * @param out_err Output: error handle on close failure, 0 on success.
 * @return FAT_OK on success; non-OK on close failure.
 */
FATSTD_API fat_Status fat_Base64EncoderClose(fat_Base64Encoder e, fat_Error *out_err);

#ifdef __cplusplus
} /* extern "C" */
#endif


#pragma once

/**
 * @file fat/xml.h
 * @brief XML encoding/decoding utilities.
 *
 * This module is backed by Go's encoding/xml package.
 *
 * Design notes:
 * - Go's `any` / reflection-driven Marshal/Unmarshal API does not map directly to C.
 * - The C API focuses on:
 *   - escaping helpers (Escape, EscapeText) writing into a FatStd bytes buffer
 *   - streaming tokenization via a handle-backed decoder
 *   - optional streaming re-encoding via a handle-backed encoder that can EncodeToken
 *
 * Recoverable failures return fat_Status and fat_Error.
 * Contract violations (invalid handles, NULL out-params where forbidden) are fatal.
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
 * @brief Opaque handle to an XML decoder.
 *
 * @note Ownership: free with fat_XmlDecoderFree.
 */
typedef fat_Handle fat_XmlDecoder;

/**
 * @brief Opaque handle to an XML encoder.
 *
 * @note Ownership: close/free with fat_XmlEncoderClose.
 */
typedef fat_Handle fat_XmlEncoder;

/**
 * @brief Opaque handle to an XML token returned from a decoder.
 *
 * @note Ownership: free with fat_XmlTokenFree.
 */
typedef fat_Handle fat_XmlToken;

/**
 * @brief XML token kinds.
 */
typedef enum fat_XmlTokenKind {
  FAT_XML_START_ELEMENT = 1,
  FAT_XML_END_ELEMENT = 2,
  FAT_XML_CHAR_DATA = 3,
  FAT_XML_COMMENT = 4,
  FAT_XML_DIRECTIVE = 5,
  FAT_XML_PROC_INST = 6
} fat_XmlTokenKind;

/**
 * @brief Creates a decoder from an in-memory XML blob.
 *
 * @param data XML bytes.
 * @return New decoder handle (caller must fat_XmlDecoderFree).
 */
FATSTD_API fat_XmlDecoder fat_XmlDecoderNewBytes(fat_Bytes data);

/**
 * @brief Opens an XML decoder over a filesystem path (UTF-8).
 *
 * @param path UTF-8 path (NUL-terminated).
 * @param out_dec Output: new decoder handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; non-OK on failure.
 */
FATSTD_API fat_Status fat_XmlDecoderOpenPathUTF8(const char *path, fat_XmlDecoder *out_dec, fat_Error *out_err);

/**
 * @brief Frees a decoder handle (and closes any underlying resources).
 *
 * @param dec Decoder handle.
 * @param out_err Output: error handle on close failure, 0 on success.
 * @return FAT_OK on success; non-OK on close failure.
 */
FATSTD_API fat_Status fat_XmlDecoderFree(fat_XmlDecoder dec, fat_Error *out_err);

/**
 * @brief Returns the input byte offset.
 *
 * Matches Go's (*xml.Decoder).InputOffset.
 */
FATSTD_API int64_t fat_XmlDecoderInputOffset(fat_XmlDecoder dec);

/**
 * @brief Returns the current input line/column.
 *
 * Matches Go's (*xml.Decoder).InputPos.
 */
FATSTD_API void fat_XmlDecoderInputPos(fat_XmlDecoder dec, int *out_line, int *out_column);

/**
 * @brief Reads the next token.
 *
 * @param dec Decoder handle.
 * @param out_tok Output: new token handle on success; 0 on EOF.
 * @param out_err Output: error handle on failure, 0 on success/EOF.
 * @return FAT_OK on success; FAT_ERR_EOF at EOF; FAT_ERR_SYNTAX on parse error.
 */
FATSTD_API fat_Status fat_XmlDecoderToken(fat_XmlDecoder dec, fat_XmlToken *out_tok, fat_Error *out_err);

/**
 * @brief Reads the next raw token (does not expand entities).
 *
 * Matches Go's (*xml.Decoder).RawToken.
 */
FATSTD_API fat_Status fat_XmlDecoderRawToken(fat_XmlDecoder dec, fat_XmlToken *out_tok, fat_Error *out_err);

/**
 * @brief Skips the current element.
 *
 * Matches Go's (*xml.Decoder).Skip.
 */
FATSTD_API fat_Status fat_XmlDecoderSkip(fat_XmlDecoder dec, fat_Error *out_err);

/**
 * @brief Frees a token handle.
 */
FATSTD_API void fat_XmlTokenFree(fat_XmlToken tok);

/**
 * @brief Returns the kind of a token.
 */
FATSTD_API fat_XmlTokenKind fat_XmlTokenType(fat_XmlToken tok);

/**
 * @brief Returns the local name for StartElement/EndElement.
 *
 * @return New fat_String handle (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_XmlTokenNameLocal(fat_XmlToken tok);

/**
 * @brief Returns the namespace space for StartElement/EndElement.
 *
 * @return New fat_String handle (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_XmlTokenNameSpace(fat_XmlToken tok);

/**
 * @brief Returns the number of attributes on a StartElement token.
 */
FATSTD_API size_t fat_XmlStartElementAttrCount(fat_XmlToken tok);

/**
 * @brief Returns an attribute name (local/space) and value for a StartElement token.
 *
 * Outputs are newly allocated fat_String handles (caller must fat_StringFree each).
 */
FATSTD_API void fat_XmlStartElementAttrGet(
  fat_XmlToken tok,
  size_t idx,
  fat_String *out_name_local,
  fat_String *out_name_space,
  fat_String *out_value
);

/**
 * @brief Returns the bytes payload of CharData/Comment/Directive tokens.
 *
 * @return New fat_Bytes handle (caller must fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_XmlTokenBytes(fat_XmlToken tok);

/**
 * @brief Returns the ProcInst target string.
 *
 * @return New fat_String handle (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_XmlProcInstTarget(fat_XmlToken tok);

/**
 * @brief Returns the ProcInst instruction bytes.
 *
 * @return New fat_Bytes handle (caller must fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_XmlProcInstInstBytes(fat_XmlToken tok);

/**
 * @brief Escapes XML text to a bytes buffer.
 *
 * Matches encoding/xml.Escape.
 */
FATSTD_API void fat_XmlEscapeToBytesBuffer(fat_BytesBuffer dst, fat_Bytes src);

/**
 * @brief Escapes XML text to a bytes buffer.
 *
 * Matches encoding/xml.EscapeText.
 */
FATSTD_API fat_Status fat_XmlEscapeTextToBytesBuffer(fat_BytesBuffer dst, fat_Bytes src, fat_Error *out_err);

/**
 * @brief Creates an encoder writing to a bytes buffer.
 *
 * @param dst Destination bytes buffer.
 * @return New encoder handle (caller must fat_XmlEncoderClose).
 */
FATSTD_API fat_XmlEncoder fat_XmlEncoderNewToBytesBuffer(fat_BytesBuffer dst);

/**
 * @brief Sets indentation for the encoder.
 */
FATSTD_API void fat_XmlEncoderIndent(fat_XmlEncoder enc, fat_String prefix, fat_String indent);

/**
 * @brief Encodes a token.
 *
 * Matches Go's (*xml.Encoder).EncodeToken.
 */
FATSTD_API fat_Status fat_XmlEncoderEncodeToken(fat_XmlEncoder enc, fat_XmlToken tok, fat_Error *out_err);

/**
 * @brief Flushes buffered encoder data.
 */
FATSTD_API fat_Status fat_XmlEncoderFlush(fat_XmlEncoder enc, fat_Error *out_err);

/**
 * @brief Closes and frees an encoder handle.
 */
FATSTD_API fat_Status fat_XmlEncoderClose(fat_XmlEncoder enc, fat_Error *out_err);

#ifdef __cplusplus
} /* extern "C" */
#endif

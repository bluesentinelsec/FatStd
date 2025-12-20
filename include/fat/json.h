#pragma once

/**
 * @file fat/json.h
 * @brief JSON encoding/decoding utilities.
 *
 * This module is backed by Go's encoding/json package.
 *
 * Design notes:
 * - Go's `any` / reflection-driven Marshal/Unmarshal API does not map directly to C.
 * - The C API exposes:
 *   - validation/formatting utilities on raw JSON bytes (Valid/Compact/Indent/HTMLEscape)
 *   - a handle-backed generic JSON value (decoded into Go's interface{}-style representation)
 *   - streaming decode/encode over FatStd bytes readers/buffers
 *
 * Recoverable failures return fat_Status and fat_Error.
 * Contract violations (invalid handles, NULL out-params where forbidden) are fatal.
 */

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#include "fat/bytes.h"
#include "fat/bytes_buffer.h"
#include "fat/bytes_reader.h"
#include "fat/error.h"
#include "fat/export.h"
#include "fat/handle.h"
#include "fat/status.h"
#include "fat/string.h"

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Opaque handle to a decoded JSON value.
 *
 * Values are produced by fat_JsonUnmarshal / fat_JsonDecoderDecodeValue and can be re-encoded.
 *
 * @note Ownership: free with fat_JsonValueFree.
 */
typedef fat_Handle fat_JsonValue;

/**
 * @brief Opaque handle to a streaming JSON decoder.
 *
 * @note Ownership: free with fat_JsonDecoderFree.
 */
typedef fat_Handle fat_JsonDecoder;

/**
 * @brief Opaque handle to a streaming JSON encoder.
 *
 * @note Ownership: free with fat_JsonEncoderFree.
 */
typedef fat_Handle fat_JsonEncoder;

/**
 * @brief JSON value kinds for fat_JsonValueType.
 */
typedef enum fat_JsonValueKind {
  FAT_JSON_NULL = 0,
  FAT_JSON_BOOL = 1,
  FAT_JSON_NUMBER = 2,
  FAT_JSON_STRING = 3,
  FAT_JSON_ARRAY = 4,
  FAT_JSON_OBJECT = 5,
} fat_JsonValueKind;

/**
 * @brief Reports whether `data` is valid JSON.
 *
 * Matches Go's json.Valid semantics.
 *
 * @param data JSON bytes.
 * @return True if valid; false otherwise.
 */
FATSTD_API bool fat_JsonValid(fat_Bytes data);

/**
 * @brief Compacts JSON by removing insignificant whitespace.
 *
 * @param src Input JSON bytes.
 * @param out Output: compacted JSON bytes (caller must fat_BytesFree).
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_SYNTAX on invalid JSON.
 */
FATSTD_API fat_Status fat_JsonCompact(fat_Bytes src, fat_Bytes *out, fat_Error *out_err);

/**
 * @brief Indents JSON with the given prefix and indent strings.
 *
 * @param src Input JSON bytes.
 * @param prefix Prefix string.
 * @param indent Indent string.
 * @param out Output: indented JSON bytes (caller must fat_BytesFree).
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_SYNTAX on invalid JSON.
 */
FATSTD_API fat_Status fat_JsonIndent(
  fat_Bytes src,
  fat_String prefix,
  fat_String indent,
  fat_Bytes *out,
  fat_Error *out_err
);

/**
 * @brief Escapes HTML characters in a JSON string.
 *
 * Matches Go's json.HTMLEscape.
 *
 * @param src Input JSON bytes.
 * @return New fat_Bytes handle containing escaped bytes (caller must fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_JsonHTMLEscape(fat_Bytes src);

/**
 * @brief Parses JSON bytes into a generic JSON value.
 *
 * Numbers are decoded as json.Number (not float64) to preserve integer/fraction text.
 *
 * @param data Input JSON bytes.
 * @param out Output: new value handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_SYNTAX on invalid JSON.
 */
FATSTD_API fat_Status fat_JsonUnmarshal(fat_Bytes data, fat_JsonValue *out, fat_Error *out_err);

/**
 * @brief Encodes a JSON value to bytes.
 *
 * @param v Value handle.
 * @param out Output: new JSON bytes (caller must fat_BytesFree).
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_OTHER on failure.
 */
FATSTD_API fat_Status fat_JsonMarshal(fat_JsonValue v, fat_Bytes *out, fat_Error *out_err);

/**
 * @brief Encodes a JSON value to indented bytes.
 *
 * @param v Value handle.
 * @param prefix Prefix string.
 * @param indent Indent string.
 * @param out Output: new JSON bytes (caller must fat_BytesFree).
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_OTHER on failure.
 */
FATSTD_API fat_Status fat_JsonMarshalIndent(
  fat_JsonValue v,
  fat_String prefix,
  fat_String indent,
  fat_Bytes *out,
  fat_Error *out_err
);

/**
 * @brief Frees a JSON value handle.
 *
 * @param v Value handle to free.
 */
FATSTD_API void fat_JsonValueFree(fat_JsonValue v);

/**
 * @brief Returns the kind of a JSON value.
 *
 * @param v Value handle.
 * @return One of fat_JsonValueKind.
 */
FATSTD_API fat_JsonValueKind fat_JsonValueType(fat_JsonValue v);

/**
 * @brief Returns a JSON boolean value.
 *
 * @param v Value handle (must be FAT_JSON_BOOL).
 * @param out_value Output: 0 or 1.
 */
FATSTD_API void fat_JsonValueAsBool(fat_JsonValue v, int *out_value);

/**
 * @brief Returns a JSON string value.
 *
 * @param v Value handle (must be FAT_JSON_STRING).
 * @return New fat_String handle (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_JsonValueAsString(fat_JsonValue v);

/**
 * @brief Returns the textual representation of a JSON number.
 *
 * @param v Value handle (must be FAT_JSON_NUMBER).
 * @return New fat_String handle containing the number text (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_JsonValueAsNumberString(fat_JsonValue v);

/**
 * @brief Returns the length of a JSON array.
 *
 * @param v Value handle (must be FAT_JSON_ARRAY).
 * @return Length.
 */
FATSTD_API size_t fat_JsonArrayLen(fat_JsonValue v);

/**
 * @brief Gets an element from a JSON array.
 *
 * @param v Value handle (must be FAT_JSON_ARRAY).
 * @param idx Index (0 <= idx < fat_JsonArrayLen(v)).
 * @return New fat_JsonValue handle for the element (caller must fat_JsonValueFree).
 */
FATSTD_API fat_JsonValue fat_JsonArrayGet(fat_JsonValue v, size_t idx);

/**
 * @brief Returns sorted keys for a JSON object.
 *
 * @param v Value handle (must be FAT_JSON_OBJECT).
 * @return New fat_StringArray handle (caller must fat_StringArrayFree).
 */
FATSTD_API fat_StringArray fat_JsonObjectKeys(fat_JsonValue v);

/**
 * @brief Gets a field from a JSON object by key.
 *
 * @param v Value handle (must be FAT_JSON_OBJECT).
 * @param key Key string.
 * @param out_found Output: set to true if found.
 * @param out_value Output: new value handle if found, 0 otherwise (caller must fat_JsonValueFree).
 */
FATSTD_API void fat_JsonObjectGet(fat_JsonValue v, fat_String key, bool *out_found, fat_JsonValue *out_value);

/**
 * @brief Creates a streaming JSON decoder over a FatStd bytes reader.
 *
 * @param r Bytes reader handle.
 * @return New decoder handle (caller must fat_JsonDecoderFree).
 */
FATSTD_API fat_JsonDecoder fat_JsonDecoderNewBytesReader(fat_BytesReader r);

/**
 * @brief Frees a JSON decoder handle.
 *
 * @param dec Decoder handle to free.
 */
FATSTD_API void fat_JsonDecoderFree(fat_JsonDecoder dec);

/**
 * @brief Enables UseNumber on the decoder (numbers decode as json.Number).
 *
 * @param dec Decoder handle.
 */
FATSTD_API void fat_JsonDecoderUseNumber(fat_JsonDecoder dec);

/**
 * @brief Enables DisallowUnknownFields on the decoder.
 *
 * @param dec Decoder handle.
 */
FATSTD_API void fat_JsonDecoderDisallowUnknownFields(fat_JsonDecoder dec);

/**
 * @brief Returns the current input offset.
 *
 * @param dec Decoder handle.
 * @return Byte offset.
 */
FATSTD_API int64_t fat_JsonDecoderInputOffset(fat_JsonDecoder dec);

/**
 * @brief Reports whether there is another element in the current array/object.
 *
 * @param dec Decoder handle.
 * @return True if more tokens exist; false otherwise.
 */
FATSTD_API bool fat_JsonDecoderMore(fat_JsonDecoder dec);

/**
 * @brief Returns the decoder's buffered data as a snapshot.
 *
 * @param dec Decoder handle.
 * @return New fat_Bytes handle (caller must fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_JsonDecoderBufferedBytes(fat_JsonDecoder dec);

/**
 * @brief Decodes the next JSON value from the stream.
 *
 * @param dec Decoder handle.
 * @param out Output: new value handle on success.
 * @param out_err Output: error handle on failure, 0 on success/EOF.
 * @return FAT_OK on success; FAT_ERR_EOF at EOF; FAT_ERR_SYNTAX on invalid JSON.
 */
FATSTD_API fat_Status fat_JsonDecoderDecodeValue(fat_JsonDecoder dec, fat_JsonValue *out, fat_Error *out_err);

/**
 * @brief Creates a streaming JSON encoder writing to a FatStd bytes buffer.
 *
 * @param dst Destination bytes buffer.
 * @return New encoder handle (caller must fat_JsonEncoderFree).
 */
FATSTD_API fat_JsonEncoder fat_JsonEncoderNewToBytesBuffer(fat_BytesBuffer dst);

/**
 * @brief Frees a JSON encoder handle.
 *
 * @param enc Encoder handle to free.
 */
FATSTD_API void fat_JsonEncoderFree(fat_JsonEncoder enc);

/**
 * @brief Sets HTML escaping behavior.
 *
 * @param enc Encoder handle.
 * @param on True to escape HTML characters.
 */
FATSTD_API void fat_JsonEncoderSetEscapeHTML(fat_JsonEncoder enc, bool on);

/**
 * @brief Sets indentation for subsequent Encode calls.
 *
 * @param enc Encoder handle.
 * @param prefix Prefix string.
 * @param indent Indent string.
 */
FATSTD_API void fat_JsonEncoderSetIndent(fat_JsonEncoder enc, fat_String prefix, fat_String indent);

/**
 * @brief Encodes a JSON value followed by a newline (Go Encoder.Encode semantics).
 *
 * @param enc Encoder handle.
 * @param v Value handle.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_OTHER on failure.
 */
FATSTD_API fat_Status fat_JsonEncoderEncodeValue(fat_JsonEncoder enc, fat_JsonValue v, fat_Error *out_err);

#ifdef __cplusplus
} /* extern "C" */
#endif


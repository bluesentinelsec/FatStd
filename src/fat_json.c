#include "fat/json.h"

#include "fatstd_go.h"

bool fat_JsonValid(fat_Bytes data) {
  return fatstd_go_json_valid((uintptr_t)data) != 0;
}

fat_Status fat_JsonCompact(fat_Bytes src, fat_Bytes *out, fat_Error *out_err) {
  return (fat_Status)fatstd_go_json_compact((uintptr_t)src, (uintptr_t *)out, (uintptr_t *)out_err);
}

fat_Status fat_JsonIndent(fat_Bytes src, fat_String prefix, fat_String indent, fat_Bytes *out, fat_Error *out_err) {
  return (fat_Status)fatstd_go_json_indent((uintptr_t)src, (uintptr_t)prefix, (uintptr_t)indent, (uintptr_t *)out,
                                           (uintptr_t *)out_err);
}

fat_Bytes fat_JsonHTMLEscape(fat_Bytes src) {
  return (fat_Bytes)fatstd_go_json_html_escape((uintptr_t)src);
}

fat_Status fat_JsonUnmarshal(fat_Bytes data, fat_JsonValue *out, fat_Error *out_err) {
  return (fat_Status)fatstd_go_json_unmarshal((uintptr_t)data, (uintptr_t *)out, (uintptr_t *)out_err);
}

fat_Status fat_JsonMarshal(fat_JsonValue v, fat_Bytes *out, fat_Error *out_err) {
  return (fat_Status)fatstd_go_json_marshal((uintptr_t)v, (uintptr_t *)out, (uintptr_t *)out_err);
}

fat_Status fat_JsonMarshalIndent(
  fat_JsonValue v,
  fat_String prefix,
  fat_String indent,
  fat_Bytes *out,
  fat_Error *out_err
) {
  return (fat_Status)fatstd_go_json_marshal_indent((uintptr_t)v, (uintptr_t)prefix, (uintptr_t)indent, (uintptr_t *)out,
                                                   (uintptr_t *)out_err);
}

void fat_JsonValueFree(fat_JsonValue v) {
  fatstd_go_json_value_free((uintptr_t)v);
}

fat_JsonValueKind fat_JsonValueType(fat_JsonValue v) {
  return (fat_JsonValueKind)fatstd_go_json_value_type((uintptr_t)v);
}

void fat_JsonValueAsBool(fat_JsonValue v, int *out_value) {
  fatstd_go_json_value_as_bool((uintptr_t)v, out_value);
}

fat_String fat_JsonValueAsString(fat_JsonValue v) {
  return (fat_String)fatstd_go_json_value_as_string((uintptr_t)v);
}

fat_String fat_JsonValueAsNumberString(fat_JsonValue v) {
  return (fat_String)fatstd_go_json_value_as_number_string((uintptr_t)v);
}

size_t fat_JsonArrayLen(fat_JsonValue v) {
  return (size_t)fatstd_go_json_array_len((uintptr_t)v);
}

fat_JsonValue fat_JsonArrayGet(fat_JsonValue v, size_t idx) {
  return (fat_JsonValue)fatstd_go_json_array_get((uintptr_t)v, idx);
}

fat_StringArray fat_JsonObjectKeys(fat_JsonValue v) {
  return (fat_StringArray)fatstd_go_json_object_keys((uintptr_t)v);
}

void fat_JsonObjectGet(fat_JsonValue v, fat_String key, bool *out_found, fat_JsonValue *out_value) {
  fatstd_go_json_object_get((uintptr_t)v, (uintptr_t)key, (_Bool *)out_found, (uintptr_t *)out_value);
}

fat_JsonDecoder fat_JsonDecoderNewBytesReader(fat_BytesReader r) {
  return (fat_JsonDecoder)fatstd_go_json_decoder_new_bytes_reader((uintptr_t)r);
}

void fat_JsonDecoderFree(fat_JsonDecoder dec) {
  fatstd_go_json_decoder_free((uintptr_t)dec);
}

void fat_JsonDecoderUseNumber(fat_JsonDecoder dec) {
  fatstd_go_json_decoder_use_number((uintptr_t)dec);
}

void fat_JsonDecoderDisallowUnknownFields(fat_JsonDecoder dec) {
  fatstd_go_json_decoder_disallow_unknown_fields((uintptr_t)dec);
}

int64_t fat_JsonDecoderInputOffset(fat_JsonDecoder dec) {
  return (int64_t)fatstd_go_json_decoder_input_offset((uintptr_t)dec);
}

bool fat_JsonDecoderMore(fat_JsonDecoder dec) {
  return fatstd_go_json_decoder_more((uintptr_t)dec) != 0;
}

fat_Bytes fat_JsonDecoderBufferedBytes(fat_JsonDecoder dec) {
  return (fat_Bytes)fatstd_go_json_decoder_buffered_bytes((uintptr_t)dec);
}

fat_Status fat_JsonDecoderDecodeValue(fat_JsonDecoder dec, fat_JsonValue *out, fat_Error *out_err) {
  return (fat_Status)fatstd_go_json_decoder_decode_value((uintptr_t)dec, (uintptr_t *)out, (uintptr_t *)out_err);
}

fat_JsonEncoder fat_JsonEncoderNewToBytesBuffer(fat_BytesBuffer dst) {
  return (fat_JsonEncoder)fatstd_go_json_encoder_new_to_bytes_buffer((uintptr_t)dst);
}

void fat_JsonEncoderFree(fat_JsonEncoder enc) {
  fatstd_go_json_encoder_free((uintptr_t)enc);
}

void fat_JsonEncoderSetEscapeHTML(fat_JsonEncoder enc, bool on) {
  fatstd_go_json_encoder_set_escape_html((uintptr_t)enc, on ? 1 : 0);
}

void fat_JsonEncoderSetIndent(fat_JsonEncoder enc, fat_String prefix, fat_String indent) {
  fatstd_go_json_encoder_set_indent((uintptr_t)enc, (uintptr_t)prefix, (uintptr_t)indent);
}

fat_Status fat_JsonEncoderEncodeValue(fat_JsonEncoder enc, fat_JsonValue v, fat_Error *out_err) {
  return (fat_Status)fatstd_go_json_encoder_encode_value((uintptr_t)enc, (uintptr_t)v, (uintptr_t *)out_err);
}


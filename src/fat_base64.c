#include "fat/base64.h"

#include "fatstd_go.h"

fat_Status fat_Base64EncodingNewUTF8(const char *alphabet, fat_Base64Encoding *out_enc, fat_Error *out_err) {
  return (fat_Status)fatstd_go_base64_encoding_new_utf8((char *)alphabet, (uintptr_t *)out_enc, (uintptr_t *)out_err);
}

fat_Base64Encoding fat_Base64EncodingStrict(fat_Base64Encoding enc) {
  return (fat_Base64Encoding)fatstd_go_base64_encoding_strict((uintptr_t)enc);
}

fat_Status fat_Base64EncodingWithPadding(
  fat_Base64Encoding enc,
  int32_t padding_rune,
  fat_Base64Encoding *out_enc,
  fat_Error *out_err
) {
  return (fat_Status)fatstd_go_base64_encoding_with_padding((uintptr_t)enc, padding_rune, (uintptr_t *)out_enc,
                                                           (uintptr_t *)out_err);
}

int fat_Base64EncodedLen(fat_Base64Encoding enc, int n) {
  return (int)fatstd_go_base64_encoded_len((uintptr_t)enc, n);
}

int fat_Base64DecodedLen(fat_Base64Encoding enc, int n) {
  return (int)fatstd_go_base64_decoded_len((uintptr_t)enc, n);
}

fat_String fat_Base64EncodeToString(fat_Base64Encoding enc, fat_Bytes src) {
  return (fat_String)fatstd_go_base64_encode_to_string((uintptr_t)enc, (uintptr_t)src);
}

fat_Bytes fat_Base64Encode(fat_Base64Encoding enc, fat_Bytes src) {
  return (fat_Bytes)fatstd_go_base64_encode((uintptr_t)enc, (uintptr_t)src);
}

fat_Status fat_Base64DecodeString(fat_Base64Encoding enc, fat_String s, fat_Bytes *out, fat_Error *out_err) {
  return (fat_Status)fatstd_go_base64_decode_string((uintptr_t)enc, (uintptr_t)s, (uintptr_t *)out, (uintptr_t *)out_err);
}

fat_Status fat_Base64Decode(fat_Base64Encoding enc, fat_Bytes src, fat_Bytes *out, fat_Error *out_err) {
  return (fat_Status)fatstd_go_base64_decode((uintptr_t)enc, (uintptr_t)src, (uintptr_t *)out, (uintptr_t *)out_err);
}

fat_Bytes fat_Base64AppendEncode(fat_Base64Encoding enc, fat_Bytes dst, fat_Bytes src) {
  return (fat_Bytes)fatstd_go_base64_append_encode((uintptr_t)enc, (uintptr_t)dst, (uintptr_t)src);
}

fat_Status fat_Base64AppendDecode(fat_Base64Encoding enc, fat_Bytes dst, fat_Bytes src, fat_Bytes *out, fat_Error *out_err) {
  return (fat_Status)fatstd_go_base64_append_decode((uintptr_t)enc, (uintptr_t)dst, (uintptr_t)src, (uintptr_t *)out,
                                                    (uintptr_t *)out_err);
}

void fat_Base64EncodingFree(fat_Base64Encoding enc) {
  fatstd_go_base64_encoding_free((uintptr_t)enc);
}

fat_Status fat_Base64EncoderNewToBytesBuffer(
  fat_Base64Encoding enc,
  fat_BytesBuffer dst,
  fat_Base64Encoder *out_encoder,
  fat_Error *out_err
) {
  return (fat_Status)fatstd_go_base64_encoder_new_to_bytes_buffer((uintptr_t)enc, (uintptr_t)dst, (uintptr_t *)out_encoder,
                                                                  (uintptr_t *)out_err);
}

fat_Status fat_Base64EncoderWrite(fat_Base64Encoder e, const void *bytes, size_t len, size_t *out_n, fat_Error *out_err) {
  return (fat_Status)fatstd_go_base64_encoder_write((uintptr_t)e, (char *)bytes, len, out_n, (uintptr_t *)out_err);
}

fat_Status fat_Base64EncoderClose(fat_Base64Encoder e, fat_Error *out_err) {
  return (fat_Status)fatstd_go_base64_encoder_close((uintptr_t)e, (uintptr_t *)out_err);
}


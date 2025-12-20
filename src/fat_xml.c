#include "fat/xml.h"

#include "fatstd_go.h"

fat_XmlDecoder fat_XmlDecoderNewBytes(fat_Bytes data) {
  return (fat_XmlDecoder)fatstd_go_xml_decoder_new_bytes((uintptr_t)data);
}

fat_Status fat_XmlDecoderOpenPathUTF8(const char *path, fat_XmlDecoder *out_dec, fat_Error *out_err) {
  return (fat_Status)fatstd_go_xml_decoder_open_path_utf8((char *)path, (uintptr_t *)out_dec, (uintptr_t *)out_err);
}

fat_Status fat_XmlDecoderFree(fat_XmlDecoder dec, fat_Error *out_err) {
  return (fat_Status)fatstd_go_xml_decoder_free((uintptr_t)dec, (uintptr_t *)out_err);
}

int64_t fat_XmlDecoderInputOffset(fat_XmlDecoder dec) {
  return (int64_t)fatstd_go_xml_decoder_input_offset((uintptr_t)dec);
}

void fat_XmlDecoderInputPos(fat_XmlDecoder dec, int *out_line, int *out_column) {
  fatstd_go_xml_decoder_input_pos((uintptr_t)dec, out_line, out_column);
}

fat_Status fat_XmlDecoderToken(fat_XmlDecoder dec, fat_XmlToken *out_tok, fat_Error *out_err) {
  return (fat_Status)fatstd_go_xml_decoder_token((uintptr_t)dec, (uintptr_t *)out_tok, (uintptr_t *)out_err);
}

fat_Status fat_XmlDecoderRawToken(fat_XmlDecoder dec, fat_XmlToken *out_tok, fat_Error *out_err) {
  return (fat_Status)fatstd_go_xml_decoder_raw_token((uintptr_t)dec, (uintptr_t *)out_tok, (uintptr_t *)out_err);
}

fat_Status fat_XmlDecoderSkip(fat_XmlDecoder dec, fat_Error *out_err) {
  return (fat_Status)fatstd_go_xml_decoder_skip((uintptr_t)dec, (uintptr_t *)out_err);
}

void fat_XmlTokenFree(fat_XmlToken tok) {
  fatstd_go_xml_token_free((uintptr_t)tok);
}

fat_XmlTokenKind fat_XmlTokenType(fat_XmlToken tok) {
  return (fat_XmlTokenKind)fatstd_go_xml_token_kind((uintptr_t)tok);
}

fat_String fat_XmlTokenNameLocal(fat_XmlToken tok) {
  return (fat_String)fatstd_go_xml_token_name_local((uintptr_t)tok);
}

fat_String fat_XmlTokenNameSpace(fat_XmlToken tok) {
  return (fat_String)fatstd_go_xml_token_name_space((uintptr_t)tok);
}

size_t fat_XmlStartElementAttrCount(fat_XmlToken tok) {
  return (size_t)fatstd_go_xml_start_element_attr_count((uintptr_t)tok);
}

void fat_XmlStartElementAttrGet(
  fat_XmlToken tok,
  size_t idx,
  fat_String *out_name_local,
  fat_String *out_name_space,
  fat_String *out_value
) {
  fatstd_go_xml_start_element_attr_get((uintptr_t)tok, idx, (uintptr_t *)out_name_local, (uintptr_t *)out_name_space,
                                       (uintptr_t *)out_value);
}

fat_Bytes fat_XmlTokenBytes(fat_XmlToken tok) {
  return (fat_Bytes)fatstd_go_xml_token_bytes((uintptr_t)tok);
}

fat_String fat_XmlProcInstTarget(fat_XmlToken tok) {
  return (fat_String)fatstd_go_xml_proc_inst_target((uintptr_t)tok);
}

fat_Bytes fat_XmlProcInstInstBytes(fat_XmlToken tok) {
  return (fat_Bytes)fatstd_go_xml_proc_inst_inst_bytes((uintptr_t)tok);
}

void fat_XmlEscapeToBytesBuffer(fat_BytesBuffer dst, fat_Bytes src) {
  fatstd_go_xml_escape_to_bytes_buffer((uintptr_t)dst, (uintptr_t)src);
}

fat_Status fat_XmlEscapeTextToBytesBuffer(fat_BytesBuffer dst, fat_Bytes src, fat_Error *out_err) {
  return (fat_Status)fatstd_go_xml_escape_text_to_bytes_buffer((uintptr_t)dst, (uintptr_t)src, (uintptr_t *)out_err);
}

fat_XmlEncoder fat_XmlEncoderNewToBytesBuffer(fat_BytesBuffer dst) {
  return (fat_XmlEncoder)fatstd_go_xml_encoder_new_to_bytes_buffer((uintptr_t)dst);
}

void fat_XmlEncoderIndent(fat_XmlEncoder enc, fat_String prefix, fat_String indent) {
  fatstd_go_xml_encoder_indent((uintptr_t)enc, (uintptr_t)prefix, (uintptr_t)indent);
}

fat_Status fat_XmlEncoderEncodeToken(fat_XmlEncoder enc, fat_XmlToken tok, fat_Error *out_err) {
  return (fat_Status)fatstd_go_xml_encoder_encode_token((uintptr_t)enc, (uintptr_t)tok, (uintptr_t *)out_err);
}

fat_Status fat_XmlEncoderFlush(fat_XmlEncoder enc, fat_Error *out_err) {
  return (fat_Status)fatstd_go_xml_encoder_flush((uintptr_t)enc, (uintptr_t *)out_err);
}

fat_Status fat_XmlEncoderClose(fat_XmlEncoder enc, fat_Error *out_err) {
  return (fat_Status)fatstd_go_xml_encoder_close((uintptr_t)enc, (uintptr_t *)out_err);
}

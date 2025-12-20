#include "fat/conv.h"

#include "fatstd_go.h"

int fat_ConvIntSize(void) {
  return (int)fatstd_go_conv_int_size();
}

fat_Bytes fat_ConvAppendBool(fat_Bytes dst, bool b) {
  return (fat_Bytes)fatstd_go_conv_append_bool((uintptr_t)dst, b ? 1 : 0);
}

fat_Bytes fat_ConvAppendInt(fat_Bytes dst, int64_t i, int base) {
  return (fat_Bytes)fatstd_go_conv_append_int((uintptr_t)dst, (long long)i, base);
}

fat_Bytes fat_ConvAppendUint(fat_Bytes dst, uint64_t i, int base) {
  return (fat_Bytes)fatstd_go_conv_append_uint((uintptr_t)dst, (unsigned long long)i, base);
}

fat_Bytes fat_ConvAppendFloat(fat_Bytes dst, double f, uint8_t fmt, int prec, int bit_size) {
  return (fat_Bytes)fatstd_go_conv_append_float((uintptr_t)dst, f, fmt, prec, bit_size);
}

fat_Bytes fat_ConvAppendQuote(fat_Bytes dst, fat_String s) {
  return (fat_Bytes)fatstd_go_conv_append_quote((uintptr_t)dst, (uintptr_t)s);
}

fat_Bytes fat_ConvAppendQuoteRune(fat_Bytes dst, uint32_t r) {
  return (fat_Bytes)fatstd_go_conv_append_quote_rune((uintptr_t)dst, r);
}

fat_Bytes fat_ConvAppendQuoteRuneToASCII(fat_Bytes dst, uint32_t r) {
  return (fat_Bytes)fatstd_go_conv_append_quote_rune_to_ascii((uintptr_t)dst, r);
}

fat_Bytes fat_ConvAppendQuoteRuneToGraphic(fat_Bytes dst, uint32_t r) {
  return (fat_Bytes)fatstd_go_conv_append_quote_rune_to_graphic((uintptr_t)dst, r);
}

fat_Bytes fat_ConvAppendQuoteToASCII(fat_Bytes dst, fat_String s) {
  return (fat_Bytes)fatstd_go_conv_append_quote_to_ascii((uintptr_t)dst, (uintptr_t)s);
}

fat_Bytes fat_ConvAppendQuoteToGraphic(fat_Bytes dst, fat_String s) {
  return (fat_Bytes)fatstd_go_conv_append_quote_to_graphic((uintptr_t)dst, (uintptr_t)s);
}

bool fat_ConvCanBackquote(fat_String s) {
  return (bool)fatstd_go_conv_can_backquote((uintptr_t)s);
}

fat_String fat_ConvFormatBool(bool b) {
  return (fat_String)fatstd_go_conv_format_bool(b ? 1 : 0);
}

fat_String fat_ConvFormatInt(int64_t i, int base) {
  return (fat_String)fatstd_go_conv_format_int((long long)i, base);
}

fat_String fat_ConvFormatUint(uint64_t i, int base) {
  return (fat_String)fatstd_go_conv_format_uint((unsigned long long)i, base);
}

fat_String fat_ConvFormatFloat(double f, uint8_t fmt, int prec, int bit_size) {
  return (fat_String)fatstd_go_conv_format_float(f, fmt, prec, bit_size);
}

fat_String fat_ConvFormatComplex(double re, double im, uint8_t fmt, int prec, int bit_size) {
  return (fat_String)fatstd_go_conv_format_complex(re, im, fmt, prec, bit_size);
}

bool fat_ConvIsGraphic(uint32_t r) {
  return (bool)fatstd_go_conv_is_graphic(r);
}

bool fat_ConvIsPrint(uint32_t r) {
  return (bool)fatstd_go_conv_is_print(r);
}

fat_String fat_ConvItoa(int i) {
  return (fat_String)fatstd_go_conv_itoa(i);
}

fat_String fat_ConvQuote(fat_String s) {
  return (fat_String)fatstd_go_conv_quote((uintptr_t)s);
}

fat_String fat_ConvQuoteRune(uint32_t r) {
  return (fat_String)fatstd_go_conv_quote_rune(r);
}

fat_String fat_ConvQuoteRuneToASCII(uint32_t r) {
  return (fat_String)fatstd_go_conv_quote_rune_to_ascii(r);
}

fat_String fat_ConvQuoteRuneToGraphic(uint32_t r) {
  return (fat_String)fatstd_go_conv_quote_rune_to_graphic(r);
}

fat_String fat_ConvQuoteToASCII(fat_String s) {
  return (fat_String)fatstd_go_conv_quote_to_ascii((uintptr_t)s);
}

fat_String fat_ConvQuoteToGraphic(fat_String s) {
  return (fat_String)fatstd_go_conv_quote_to_graphic((uintptr_t)s);
}

fat_Status fat_ConvParseBool(fat_String s, int *out_value, fat_Error *out_err) {
  return (fat_Status)fatstd_go_conv_parse_bool((uintptr_t)s, out_value, (uintptr_t *)out_err);
}

fat_Status fat_ConvParseInt(fat_String s, int base, int bit_size, int64_t *out_value, fat_Error *out_err) {
  return (fat_Status)fatstd_go_conv_parse_int((uintptr_t)s, base, bit_size, (long long *)out_value,
                                              (uintptr_t *)out_err);
}

fat_Status fat_ConvParseUint(fat_String s, int base, int bit_size, uint64_t *out_value, fat_Error *out_err) {
  return (fat_Status)fatstd_go_conv_parse_uint((uintptr_t)s, base, bit_size, (unsigned long long *)out_value,
                                               (uintptr_t *)out_err);
}

fat_Status fat_ConvParseFloat(fat_String s, int bit_size, double *out_value, fat_Error *out_err) {
  return (fat_Status)fatstd_go_conv_parse_float((uintptr_t)s, bit_size, out_value, (uintptr_t *)out_err);
}

fat_Status fat_ConvParseComplex(fat_String s, int bit_size, double *out_re, double *out_im, fat_Error *out_err) {
  return (fat_Status)fatstd_go_conv_parse_complex((uintptr_t)s, bit_size, out_re, out_im, (uintptr_t *)out_err);
}

fat_Status fat_ConvAtoi(fat_String s, int64_t *out_value, fat_Error *out_err) {
  return (fat_Status)fatstd_go_conv_atoi((uintptr_t)s, (long long *)out_value, (uintptr_t *)out_err);
}

fat_Status fat_ConvUnquote(fat_String s, fat_String *out_value, fat_Error *out_err) {
  return (fat_Status)fatstd_go_conv_unquote((uintptr_t)s, (uintptr_t *)out_value, (uintptr_t *)out_err);
}

fat_Status fat_ConvQuotedPrefix(fat_String s, fat_String *out_value, fat_Error *out_err) {
  return (fat_Status)fatstd_go_conv_quoted_prefix((uintptr_t)s, (uintptr_t *)out_value, (uintptr_t *)out_err);
}

fat_Status fat_ConvUnquoteChar(
  fat_String s,
  uint8_t quote,
  uint32_t *out_rune,
  bool *out_multibyte,
  fat_String *out_tail,
  fat_Error *out_err
) {
  return (fat_Status)fatstd_go_conv_unquote_char((uintptr_t)s, quote, out_rune, (_Bool *)out_multibyte,
                                                 (uintptr_t *)out_tail, (uintptr_t *)out_err);
}


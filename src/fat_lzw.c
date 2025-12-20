#include "fat/lzw.h"

#include "fatstd_go.h"

fat_Status fat_LzwCompress(fat_Bytes src, fat_LzwOrder order, uint8_t lit_width, fat_Bytes *out, fat_Error *out_err) {
  return (fat_Status)fatstd_go_lzw_compress((uintptr_t)src, (int)order, (uint8_t)lit_width, (uintptr_t *)out,
                                           (uintptr_t *)out_err);
}

fat_Status fat_LzwDecompress(fat_Bytes src, fat_LzwOrder order, uint8_t lit_width, fat_Bytes *out, fat_Error *out_err) {
  return (fat_Status)fatstd_go_lzw_decompress((uintptr_t)src, (int)order, (uint8_t)lit_width, (uintptr_t *)out,
                                             (uintptr_t *)out_err);
}

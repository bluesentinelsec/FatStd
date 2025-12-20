#include "fat/zlib.h"

#include "fatstd_go.h"

fat_Status fat_ZlibCompress(fat_Bytes src, fat_Bytes *out, fat_Error *out_err) {
  return (fat_Status)fatstd_go_zlib_compress((uintptr_t)src, (uintptr_t *)out, (uintptr_t *)out_err);
}

fat_Status fat_ZlibDecompress(fat_Bytes src, fat_Bytes *out, fat_Error *out_err) {
  return (fat_Status)fatstd_go_zlib_decompress((uintptr_t)src, (uintptr_t *)out, (uintptr_t *)out_err);
}

#include "fat/flate.h"

#include "fatstd_go.h"

fat_Status fat_FlateCompress(fat_Bytes src, fat_Bytes *out, fat_Error *out_err) {
  return (fat_Status)fatstd_go_flate_compress((uintptr_t)src, (uintptr_t *)out, (uintptr_t *)out_err);
}

fat_Status fat_FlateDecompress(fat_Bytes src, fat_Bytes *out, fat_Error *out_err) {
  return (fat_Status)fatstd_go_flate_decompress((uintptr_t)src, (uintptr_t *)out, (uintptr_t *)out_err);
}

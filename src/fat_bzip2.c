#include "fat/bzip2.h"

#include "fatstd_go.h"

fat_Status fat_Bzip2Decompress(fat_Bytes src, fat_Bytes *out, fat_Error *out_err) {
  return (fat_Status)fatstd_go_bzip2_decompress((uintptr_t)src, (uintptr_t *)out, (uintptr_t *)out_err);
}

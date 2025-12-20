#pragma once

/**
 * @file fat/bzip2.h
 * @brief bzip2 decompression utilities.
 *
 * This module is backed by Go's compress/bzip2 package.
 *
 * Design notes:
 * - Go's compress/bzip2 package supports decompression only; FatStd mirrors that.
 * - Recoverable failures (corrupt input) return fat_Status and fat_Error.
 *
 * Contract violations (invalid handles, NULL out-params) are fatal.
 */

#include <stddef.h>
#include <stdint.h>

#include "fat/bytes.h"
#include "fat/error.h"
#include "fat/export.h"
#include "fat/status.h"

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Decompresses bzip2-compressed bytes.
 *
 * @param src Source bytes.
 * @param out Output: decompressed bytes (caller must fat_BytesFree).
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_SYNTAX on corrupt input.
 */
FATSTD_API fat_Status fat_Bzip2Decompress(fat_Bytes src, fat_Bytes *out, fat_Error *out_err);

#ifdef __cplusplus
} /* extern "C" */
#endif

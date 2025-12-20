#pragma once

/**
 * @file fat/flate.h
 * @brief Raw DEFLATE compression utilities.
 *
 * This module is backed by Go's compress/flate package.
 *
 * Design notes:
 * - This module operates on raw DEFLATE streams (no zlib/gzip wrapper).
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
 * @brief Compresses data using raw DEFLATE.
 *
 * Uses Go's default compression level.
 *
 * @param src Source bytes.
 * @param out Output: compressed bytes (caller must fat_BytesFree).
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_OTHER on failure.
 */
FATSTD_API fat_Status fat_FlateCompress(fat_Bytes src, fat_Bytes *out, fat_Error *out_err);

/**
 * @brief Decompresses raw DEFLATE data.
 *
 * @param src Source bytes.
 * @param out Output: decompressed bytes (caller must fat_BytesFree).
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_SYNTAX on corrupt input.
 */
FATSTD_API fat_Status fat_FlateDecompress(fat_Bytes src, fat_Bytes *out, fat_Error *out_err);

#ifdef __cplusplus
} /* extern "C" */
#endif

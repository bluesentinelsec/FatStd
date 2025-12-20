#pragma once

/**
 * @file fat/lzw.h
 * @brief LZW compression utilities.
 *
 * This module is backed by Go's compress/lzw package.
 *
 * Design notes:
 * - LZW requires an explicit bit order and literal width; the C API exposes both.
 * - Recoverable failures (corrupt input, invalid parameters) return fat_Status and fat_Error.
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
 * @brief LZW bit order configuration.
 */
typedef enum fat_LzwOrder {
  FAT_LZW_ORDER_LSB = 0,
  FAT_LZW_ORDER_MSB = 1,
} fat_LzwOrder;

/**
 * @brief Compresses data using LZW.
 *
 * @param src Source bytes.
 * @param order Bit order (FAT_LZW_ORDER_LSB or FAT_LZW_ORDER_MSB).
 * @param lit_width Literal width (2..8).
 * @param out Output: compressed bytes (caller must fat_BytesFree).
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_RANGE on invalid parameters; FAT_ERR_OTHER on failure.
 */
FATSTD_API fat_Status fat_LzwCompress(
  fat_Bytes src,
  fat_LzwOrder order,
  uint8_t lit_width,
  fat_Bytes *out,
  fat_Error *out_err
);

/**
 * @brief Decompresses LZW data.
 *
 * @param src Source bytes.
 * @param order Bit order used during compression.
 * @param lit_width Literal width used during compression (2..8).
 * @param out Output: decompressed bytes (caller must fat_BytesFree).
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_RANGE on invalid parameters; FAT_ERR_SYNTAX on corrupt input.
 */
FATSTD_API fat_Status fat_LzwDecompress(
  fat_Bytes src,
  fat_LzwOrder order,
  uint8_t lit_width,
  fat_Bytes *out,
  fat_Error *out_err
);

#ifdef __cplusplus
} /* extern "C" */
#endif

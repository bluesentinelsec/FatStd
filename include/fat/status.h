#pragma once

/**
 * @file fat/status.h
 * @brief Status codes for recoverable FatStd failures.
 *
 * Contract violations (invalid handles, NULL where forbidden, out-of-range indices) are
 * fatal by default and do not use fat_Status.
 *
 * fat_Status is used for recoverable failures where bubbling up errors is required.
 */

#include "fat/export.h"

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Status code for recoverable operations.
 *
 * @note FAT_OK indicates success.
 */
typedef enum fat_Status {
  FAT_OK = 0,
  FAT_ERR_SYNTAX = 1,
  FAT_ERR_RANGE = 2,
  FAT_ERR_EOF = 3,
  FAT_ERR_OTHER = 100,
} fat_Status;

#ifdef __cplusplus
} /* extern "C" */
#endif


#pragma once

/**
 * @file fat/error.h
 * @brief Handle-backed error details for recoverable failures.
 *
 * fat_Error is used alongside fat_Status for APIs where failure is normal in correct
 * programs (e.g., parsing, I/O).
 *
 * All functions are fail-fast: invalid handles and contract violations are fatal.
 */

#include "fat/export.h"
#include "fat/handle.h"
#include "fat/string.h"

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Opaque handle to a FatStd error object.
 *
 * @note Ownership: free with fat_ErrorFree.
 */
typedef fat_Handle fat_Error;

/**
 * @brief Frees an error handle.
 *
 * After this call, the handle is invalid and must not be used.
 *
 * @param e Error handle to free.
 */
FATSTD_API void fat_ErrorFree(fat_Error e);

/**
 * @brief Returns an error message string for an error handle.
 *
 * @param e Error handle.
 * @return A new fat_String handle containing a message (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_ErrorMessage(fat_Error e);

#ifdef __cplusplus
} /* extern "C" */
#endif


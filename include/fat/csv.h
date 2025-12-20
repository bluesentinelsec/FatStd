#pragma once

/**
 * @file fat/csv.h
 * @brief CSV reader/writer utilities.
 *
 * This module is backed by Go's encoding/csv package.
 *
 * Design notes:
 * - Go's io.Reader/io.Writer do not map cleanly to C. The C API uses explicit
 *   operations over FatStd bytes buffers and byte/string handles.
 * - Recoverable failures (invalid CSV input, I/O failures) return fat_Status
 *   and fat_Error.
 *
 * Contract violations (invalid handles, NULL out-params where forbidden) are fatal.
 */

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#include "fat/bytes.h"
#include "fat/bytes_buffer.h"
#include "fat/error.h"
#include "fat/export.h"
#include "fat/handle.h"
#include "fat/status.h"
#include "fat/string.h"

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Opaque handle to a CSV reader.
 *
 * @note Ownership: free with fat_CsvReaderFree.
 */
typedef fat_Handle fat_CsvReader;

/**
 * @brief Opaque handle to a CSV writer.
 *
 * @note Ownership: free with fat_CsvWriterFree.
 */
typedef fat_Handle fat_CsvWriter;

/**
 * @brief Creates a CSV reader from an in-memory UTF-8/byte input.
 *
 * @param data Bytes handle containing the CSV data.
 * @return New reader handle (caller must fat_CsvReaderFree).
 */
FATSTD_API fat_CsvReader fat_CsvReaderNewBytes(fat_Bytes data);

/**
 * @brief Frees a CSV reader handle.
 *
 * @param r Reader handle to free.
 */
FATSTD_API void fat_CsvReaderFree(fat_CsvReader r);

/**
 * @brief Reads the next record from a CSV reader.
 *
 * On success, returns a new fat_StringArray handle containing the record fields.
 * At end-of-input, returns FAT_ERR_EOF and sets `*out_record` to 0.
 *
 * @param r Reader handle.
 * @param out_record Output: new fat_StringArray on success, 0 on EOF/failure.
 * @param out_eof Output: set to true at EOF.
 * @param out_err Output: error handle on failure, 0 on success/EOF.
 * @return FAT_OK on success; FAT_ERR_EOF at EOF; FAT_ERR_SYNTAX on parse error; FAT_ERR_OTHER otherwise.
 */
FATSTD_API fat_Status fat_CsvReaderRead(
  fat_CsvReader r,
  fat_StringArray *out_record,
  bool *out_eof,
  fat_Error *out_err
);

/**
 * @brief Returns the line/column of the start of the given field in the most recent record.
 *
 * Matches Go's (*csv.Reader).FieldPos semantics.
 *
 * @param r Reader handle.
 * @param field Field index.
 * @param out_line Output: line (1-based).
 * @param out_column Output: column (1-based byte index).
 */
FATSTD_API void fat_CsvReaderFieldPos(fat_CsvReader r, int field, int *out_line, int *out_column);

/**
 * @brief Returns the input byte offset of the reader.
 *
 * Matches Go's (*csv.Reader).InputOffset semantics.
 *
 * @param r Reader handle.
 * @return Byte offset.
 */
FATSTD_API int64_t fat_CsvReaderInputOffset(fat_CsvReader r);

/**
 * @brief Creates a CSV writer that writes to a FatStd bytes buffer.
 *
 * @param dst Destination bytes buffer (receives CSV output bytes).
 * @return New writer handle (caller must fat_CsvWriterFree).
 */
FATSTD_API fat_CsvWriter fat_CsvWriterNewToBytesBuffer(fat_BytesBuffer dst);

/**
 * @brief Writes a single CSV record.
 *
 * @param w Writer handle.
 * @param fields Pointer to an array of fat_String handles (may be NULL only if n_fields == 0).
 * @param n_fields Number of fields.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_OTHER on failure.
 */
FATSTD_API fat_Status fat_CsvWriterWriteRecord(
  fat_CsvWriter w,
  const fat_String *fields,
  size_t n_fields,
  fat_Error *out_err
);

/**
 * @brief Flushes any buffered data to the underlying bytes buffer.
 *
 * @param w Writer handle.
 */
FATSTD_API void fat_CsvWriterFlush(fat_CsvWriter w);

/**
 * @brief Returns any writer error encountered during previous writes or flushes.
 *
 * @param w Writer handle.
 * @param out_err Output: error handle on failure, 0 if no error.
 * @return FAT_OK if no error; FAT_ERR_OTHER if an error is present.
 */
FATSTD_API fat_Status fat_CsvWriterError(fat_CsvWriter w, fat_Error *out_err);

/**
 * @brief Frees a CSV writer handle.
 *
 * Note: The caller should call fat_CsvWriterFlush and/or fat_CsvWriterError as needed.
 *
 * @param w Writer handle to free.
 */
FATSTD_API void fat_CsvWriterFree(fat_CsvWriter w);

#ifdef __cplusplus
} /* extern "C" */
#endif


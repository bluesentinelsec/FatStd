#pragma once

/**
 * @file fat/tar.h
 * @brief TAR archive reader/writer utilities.
 *
 * This module is backed by Go's archive/tar package.
 *
 * Design notes:
 * - The Go archive/tar surface (io.Reader/io.Writer streams, fs.FS) does not map cleanly to C.
 * - The C API exposes a C-friendly subset based on opaque handles and explicit read/write calls.
 * - Recoverable failures (I/O, corrupt tar) return fat_Status and fat_Error.
 *
 * Contract violations (invalid handles, NULL out-params, out-of-range values) are fatal.
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
 * @brief Opaque handle to a TAR reader.
 *
 * @note Ownership: free with fat_TarReaderFree.
 */
typedef fat_Handle fat_TarReader;

/**
 * @brief Opaque handle to a TAR header (metadata for one entry).
 *
 * @note Ownership: free with fat_TarHeaderFree.
 */
typedef fat_Handle fat_TarHeader;

/**
 * @brief Opaque handle to a TAR writer.
 *
 * @note Ownership: close/free with fat_TarWriterClose.
 */
typedef fat_Handle fat_TarWriter;

/**
 * @brief Creates a TAR reader from an in-memory tar blob.
 *
 * @param tar_bytes Bytes handle containing the entire tar file.
 * @param out_reader Output: new reader handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; non-OK on failure.
 */
FATSTD_API fat_Status fat_TarReaderNewBytes(fat_Bytes tar_bytes, fat_TarReader *out_reader, fat_Error *out_err);

/**
 * @brief Opens a TAR archive from a filesystem path (UTF-8).
 *
 * @param path UTF-8 path to an existing tar file (NUL-terminated).
 * @param out_reader Output: new reader handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; non-OK on failure.
 */
FATSTD_API fat_Status fat_TarReaderOpenPathUTF8(const char *path, fat_TarReader *out_reader, fat_Error *out_err);

/**
 * @brief Frees a TAR reader handle and closes any underlying resources.
 *
 * @param r Reader handle to free.
 * @param out_err Output: error handle on close failure, 0 on success.
 * @return FAT_OK on success; non-OK on close failure.
 */
FATSTD_API fat_Status fat_TarReaderFree(fat_TarReader r, fat_Error *out_err);

/**
 * @brief Advances to the next entry and returns its header.
 *
 * @param r Reader handle.
 * @param out_hdr Output: new header handle on success; 0 at EOF.
 * @param out_eof Output: set to true if end-of-archive is reached.
 * @param out_err Output: error handle on failure, 0 on success/EOF.
 * @return FAT_OK on success; FAT_ERR_EOF at end-of-archive; non-OK on failure.
 */
FATSTD_API fat_Status fat_TarReaderNext(
  fat_TarReader r,
  fat_TarHeader *out_hdr,
  bool *out_eof,
  fat_Error *out_err
);

/**
 * @brief Reads from the current entry body.
 *
 * After calling fat_TarReaderNext, this reads the entry's data.
 * When the end of the current entry is reached, FAT_ERR_EOF is returned.
 *
 * @param r Reader handle.
 * @param dst Destination buffer (may be NULL only if dst_len == 0).
 * @param dst_len Capacity of `dst` in bytes.
 * @param out_n Output: number of bytes read (0 is valid at EOF).
 * @param out_eof Output: set to true if EOF was hit (possibly after reading some bytes).
 * @param out_err Output: error handle on failure, 0 on success/EOF.
 * @return FAT_OK on success; FAT_ERR_EOF at end-of-entry; non-OK on failure.
 */
FATSTD_API fat_Status fat_TarReaderRead(
  fat_TarReader r,
  void *dst,
  size_t dst_len,
  size_t *out_n,
  bool *out_eof,
  fat_Error *out_err
);

/**
 * @brief Frees a TAR header handle.
 *
 * @param h Header handle to free.
 */
FATSTD_API void fat_TarHeaderFree(fat_TarHeader h);

/**
 * @brief Returns the entry name.
 *
 * @param h Header handle.
 * @return New fat_String handle containing the name (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_TarHeaderName(fat_TarHeader h);

/**
 * @brief Returns the entry type flag (e.g. TypeReg, TypeDir).
 *
 * @param h Header handle.
 * @return Type flag byte.
 */
FATSTD_API uint8_t fat_TarHeaderTypeflag(fat_TarHeader h);

/**
 * @brief Returns the entry size in bytes.
 *
 * @param h Header handle.
 * @return Size in bytes.
 */
FATSTD_API int64_t fat_TarHeaderSize(fat_TarHeader h);

/**
 * @brief Creates a new TAR writer that appends output to a FatStd bytes buffer.
 *
 * The buffer receives the raw tar bytes; the caller can snapshot it via fat_BytesBufferBytes.
 *
 * @param dst Destination buffer handle.
 * @param out_writer Output: new writer handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; non-OK on failure.
 */
FATSTD_API fat_Status fat_TarWriterNewToBytesBuffer(fat_BytesBuffer dst, fat_TarWriter *out_writer, fat_Error *out_err);

/**
 * @brief Adds a regular file entry with the given name and contents.
 *
 * @param w Writer handle.
 * @param name Entry name (UTF-8).
 * @param data Entry contents bytes.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; non-OK on failure.
 */
FATSTD_API fat_Status fat_TarWriterAddBytes(fat_TarWriter w, fat_String name, fat_Bytes data, fat_Error *out_err);

/**
 * @brief Flushes the TAR writer.
 *
 * @param w Writer handle.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; non-OK on failure.
 */
FATSTD_API fat_Status fat_TarWriterFlush(fat_TarWriter w, fat_Error *out_err);

/**
 * @brief Closes and frees a TAR writer handle.
 *
 * @param w Writer handle to close.
 * @param out_err Output: error handle on close failure, 0 on success.
 * @return FAT_OK on success; non-OK on close failure.
 */
FATSTD_API fat_Status fat_TarWriterClose(fat_TarWriter w, fat_Error *out_err);

#ifdef __cplusplus
} /* extern "C" */
#endif


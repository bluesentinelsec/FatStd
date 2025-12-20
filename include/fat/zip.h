#pragma once

/**
 * @file fat/zip.h
 * @brief ZIP archive reader/writer utilities.
 *
 * This module is backed by Go's archive/zip package.
 *
 * Design notes:
 * - The Go archive/zip surface (io.Reader/io.Writer callbacks, fs.FS, etc.) does not map cleanly to C.
 * - The C API therefore exposes a C-friendly subset based on opaque handles and explicit read/write calls.
 * - Recoverable failures (I/O, corrupt zip, unsupported compression) return fat_Status and fat_Error.
 *
 * Contract violations (invalid handles, NULL out-params, out-of-range indices) are fatal.
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
 * @brief Opaque handle to a ZIP reader.
 *
 * The reader owns any underlying resources (file descriptors, in-memory copies).
 *
 * @note Ownership: free with fat_ZipReaderFree.
 */
typedef fat_Handle fat_ZipReader;

/**
 * @brief Opaque handle to a ZIP file entry (metadata + access to contents).
 *
 * File handles are valid only while their parent fat_ZipReader remains valid.
 *
 * @note Ownership: free with fat_ZipFileFree.
 */
typedef fat_Handle fat_ZipFile;

/**
 * @brief Opaque handle to an open ZIP entry stream.
 *
 * @note Ownership: close/free with fat_ZipFileReaderClose.
 */
typedef fat_Handle fat_ZipFileReader;

/**
 * @brief Opaque handle to a ZIP writer.
 *
 * @note Ownership: close/free with fat_ZipWriterClose.
 */
typedef fat_Handle fat_ZipWriter;

/**
 * @brief Opens a ZIP archive from a filesystem path (UTF-8).
 *
 * @param path UTF-8 path to an existing zip file (NUL-terminated).
 * @param out_reader Output: new reader handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; non-OK on failure.
 */
FATSTD_API fat_Status fat_ZipReaderOpenPathUTF8(const char *path, fat_ZipReader *out_reader, fat_Error *out_err);

/**
 * @brief Creates a ZIP reader from an in-memory zip blob.
 *
 * @param zip_bytes Bytes handle containing the entire zip file.
 * @param out_reader Output: new reader handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; non-OK on failure.
 */
FATSTD_API fat_Status fat_ZipReaderNewBytes(fat_Bytes zip_bytes, fat_ZipReader *out_reader, fat_Error *out_err);

/**
 * @brief Frees a ZIP reader handle and closes any underlying resources.
 *
 * @param r Reader handle to free.
 * @param out_err Output: error handle on close failure, 0 on success.
 * @return FAT_OK on success; non-OK on close failure.
 */
FATSTD_API fat_Status fat_ZipReaderFree(fat_ZipReader r, fat_Error *out_err);

/**
 * @brief Returns the number of files in the archive.
 *
 * @param r Reader handle.
 * @return Number of file entries.
 */
FATSTD_API size_t fat_ZipReaderNumFiles(fat_ZipReader r);

/**
 * @brief Returns the file entry at index `idx` as a new handle.
 *
 * @param r Reader handle.
 * @param idx File index (0 <= idx < fat_ZipReaderNumFiles(r)).
 * @return New fat_ZipFile handle (caller must fat_ZipFileFree).
 */
FATSTD_API fat_ZipFile fat_ZipReaderFileByIndex(fat_ZipReader r, size_t idx);

/**
 * @brief Frees a ZIP file entry handle.
 *
 * @param f File handle to free.
 */
FATSTD_API void fat_ZipFileFree(fat_ZipFile f);

/**
 * @brief Returns the entry name as a string.
 *
 * @param f File handle.
 * @return New fat_String handle containing the name (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_ZipFileName(fat_ZipFile f);

/**
 * @brief Opens the file entry for reading.
 *
 * @param f File handle.
 * @param out_reader Output: new file reader handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; non-OK on failure.
 */
FATSTD_API fat_Status fat_ZipFileOpen(fat_ZipFile f, fat_ZipFileReader *out_reader, fat_Error *out_err);

/**
 * @brief Reads from an open ZIP entry stream.
 *
 * @param r File reader handle.
 * @param dst Destination buffer (may be NULL only if dst_len == 0).
 * @param dst_len Capacity of `dst` in bytes.
 * @param out_n Output: number of bytes read (0 is valid at EOF).
 * @param out_eof Output: set to true if EOF was hit (possibly after reading some bytes).
 * @param out_err Output: error handle on failure, 0 on success/EOF.
 * @return FAT_OK on success; FAT_ERR_EOF at EOF; non-OK on other failures.
 */
FATSTD_API fat_Status fat_ZipFileReaderRead(
  fat_ZipFileReader r,
  void *dst,
  size_t dst_len,
  size_t *out_n,
  bool *out_eof,
  fat_Error *out_err
);

/**
 * @brief Closes and frees a ZIP file reader handle.
 *
 * @param r File reader handle to close.
 * @param out_err Output: error handle on close failure, 0 on success.
 * @return FAT_OK on success; non-OK on close failure.
 */
FATSTD_API fat_Status fat_ZipFileReaderClose(fat_ZipFileReader r, fat_Error *out_err);

/**
 * @brief Creates a new ZIP writer that appends output to a FatStd bytes buffer.
 *
 * The buffer receives the raw zip bytes; the caller can snapshot it via fat_BytesBufferBytes.
 *
 * @param dst Destination buffer handle.
 * @param out_writer Output: new writer handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; non-OK on failure.
 */
FATSTD_API fat_Status fat_ZipWriterNewToBytesBuffer(fat_BytesBuffer dst, fat_ZipWriter *out_writer, fat_Error *out_err);

/**
 * @brief Adds a file entry with the given name and contents.
 *
 * @param w Writer handle.
 * @param name Entry name (UTF-8).
 * @param data Entry contents bytes.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; non-OK on failure.
 */
FATSTD_API fat_Status fat_ZipWriterAddBytes(fat_ZipWriter w, fat_String name, fat_Bytes data, fat_Error *out_err);

/**
 * @brief Closes and frees a ZIP writer handle.
 *
 * @param w Writer handle to close.
 * @param out_err Output: error handle on close failure, 0 on success.
 * @return FAT_OK on success; non-OK on close failure.
 */
FATSTD_API fat_Status fat_ZipWriterClose(fat_ZipWriter w, fat_Error *out_err);

#ifdef __cplusplus
} /* extern "C" */
#endif


#pragma once

/**
 * @file fat/socket.h
 * @brief TCP/UDP socket utilities.
 *
 * This module is backed by Go's net package.
 *
 * Design notes:
 * - The API is intentionally small: dial/listen/accept plus simple read/write.
 * - Recoverable failures (network errors, invalid addresses) return fat_Status and fat_Error.
 * - Contract violations (invalid handles, NULL out-params, NULL pointers where forbidden)
 *   are fatal.
 */

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#include "fat/error.h"
#include "fat/export.h"
#include "fat/handle.h"
#include "fat/status.h"
#include "fat/string.h"

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Opaque handle to a TCP connection.
 *
 * @note Ownership: close with fat_TcpConnClose.
 */
typedef fat_Handle fat_TcpConn;

/**
 * @brief Opaque handle to a TCP listener.
 *
 * @note Ownership: close with fat_TcpListenerClose.
 */
typedef fat_Handle fat_TcpListener;

/**
 * @brief Opaque handle to a UDP connection.
 *
 * @note Ownership: close with fat_UdpConnClose.
 */
typedef fat_Handle fat_UdpConn;

/**
 * @brief Dials a TCP address (UTF-8, host:port) and returns a connection.
 *
 * @param addr TCP address to dial (NUL-terminated).
 * @param out_conn Output: new connection handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_SYNTAX on invalid address; FAT_ERR_OTHER otherwise.
 */
FATSTD_API fat_Status fat_TcpDialUTF8(const char *addr, fat_TcpConn *out_conn, fat_Error *out_err);

/**
 * @brief Listens on a TCP address (UTF-8, host:port) and returns a listener.
 *
 * Pass "127.0.0.1:0" to bind to an ephemeral port.
 *
 * @param addr TCP listen address (NUL-terminated).
 * @param out_listener Output: new listener handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_SYNTAX on invalid address; FAT_ERR_OTHER otherwise.
 */
FATSTD_API fat_Status fat_TcpListenerListenUTF8(const char *addr, fat_TcpListener *out_listener, fat_Error *out_err);

/**
 * @brief Accepts the next TCP connection.
 *
 * This call blocks until a connection is available.
 *
 * @param listener Listener handle.
 * @param out_conn Output: new connection handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_OTHER on failure.
 */
FATSTD_API fat_Status fat_TcpListenerAccept(fat_TcpListener listener, fat_TcpConn *out_conn, fat_Error *out_err);

/**
 * @brief Returns the listener's bound address as a string.
 *
 * @param listener Listener handle.
 * @return New fat_String handle containing "host:port" (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_TcpListenerAddr(fat_TcpListener listener);

/**
 * @brief Closes a TCP listener.
 *
 * @param listener Listener handle.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_OTHER on failure.
 */
FATSTD_API fat_Status fat_TcpListenerClose(fat_TcpListener listener, fat_Error *out_err);

/**
 * @brief Reads from a TCP connection into a caller buffer.
 *
 * @param conn TCP connection handle.
 * @param dst Destination buffer (may be NULL only if dst_len == 0).
 * @param dst_len Capacity of dst in bytes.
 * @param out_n Output: number of bytes read.
 * @param out_eof Output: true if EOF was hit.
 * @param out_err Output: error handle on failure, 0 on success/EOF.
 * @return FAT_OK on success; FAT_ERR_EOF at EOF; FAT_ERR_OTHER otherwise.
 */
FATSTD_API fat_Status fat_TcpConnRead(
  fat_TcpConn conn,
  void *dst,
  size_t dst_len,
  size_t *out_n,
  bool *out_eof,
  fat_Error *out_err
);

/**
 * @brief Writes data to a TCP connection.
 *
 * @param conn TCP connection handle.
 * @param src Source bytes (may be NULL only if src_len == 0).
 * @param src_len Number of bytes to write.
 * @param out_n Output: number of bytes written.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_OTHER on failure.
 */
FATSTD_API fat_Status fat_TcpConnWrite(
  fat_TcpConn conn,
  const void *src,
  size_t src_len,
  size_t *out_n,
  fat_Error *out_err
);

/**
 * @brief Returns the local address of a TCP connection.
 *
 * @param conn TCP connection handle.
 * @return New fat_String handle containing "host:port" (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_TcpConnLocalAddr(fat_TcpConn conn);

/**
 * @brief Returns the remote address of a TCP connection.
 *
 * @param conn TCP connection handle.
 * @return New fat_String handle containing "host:port" (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_TcpConnRemoteAddr(fat_TcpConn conn);

/**
 * @brief Closes a TCP connection.
 *
 * @param conn TCP connection handle.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_OTHER on failure.
 */
FATSTD_API fat_Status fat_TcpConnClose(fat_TcpConn conn, fat_Error *out_err);

/**
 * @brief Binds a UDP socket to a local address (UTF-8, host:port).
 *
 * Pass "127.0.0.1:0" to bind to an ephemeral port.
 *
 * @param addr UDP listen address (NUL-terminated).
 * @param out_conn Output: new UDP connection handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_SYNTAX on invalid address; FAT_ERR_OTHER otherwise.
 */
FATSTD_API fat_Status fat_UdpListenUTF8(const char *addr, fat_UdpConn *out_conn, fat_Error *out_err);

/**
 * @brief Creates a UDP connection to a remote address (UTF-8, host:port).
 *
 * @param addr UDP remote address (NUL-terminated).
 * @param out_conn Output: new UDP connection handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_SYNTAX on invalid address; FAT_ERR_OTHER otherwise.
 */
FATSTD_API fat_Status fat_UdpDialUTF8(const char *addr, fat_UdpConn *out_conn, fat_Error *out_err);

/**
 * @brief Reads a datagram and returns the sender address.
 *
 * @param conn UDP connection handle.
 * @param dst Destination buffer (may be NULL only if dst_len == 0).
 * @param dst_len Capacity of dst in bytes.
 * @param out_n Output: number of bytes read.
 * @param out_addr Output: sender address string (caller must fat_StringFree).
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_OTHER on failure.
 */
FATSTD_API fat_Status fat_UdpConnReadFrom(
  fat_UdpConn conn,
  void *dst,
  size_t dst_len,
  size_t *out_n,
  fat_String *out_addr,
  fat_Error *out_err
);

/**
 * @brief Writes a datagram to a destination address (UTF-8, host:port).
 *
 * @param conn UDP connection handle.
 * @param src Source bytes (may be NULL only if src_len == 0).
 * @param src_len Number of bytes to send.
 * @param addr Destination address (NUL-terminated).
 * @param out_n Output: number of bytes written.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_SYNTAX on invalid address; FAT_ERR_OTHER otherwise.
 */
FATSTD_API fat_Status fat_UdpConnWriteToUTF8(
  fat_UdpConn conn,
  const void *src,
  size_t src_len,
  const char *addr,
  size_t *out_n,
  fat_Error *out_err
);

/**
 * @brief Writes a datagram on a connected UDP socket.
 *
 * @param conn UDP connection handle.
 * @param src Source bytes (may be NULL only if src_len == 0).
 * @param src_len Number of bytes to send.
 * @param out_n Output: number of bytes written.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_OTHER on failure.
 */
FATSTD_API fat_Status fat_UdpConnWrite(
  fat_UdpConn conn,
  const void *src,
  size_t src_len,
  size_t *out_n,
  fat_Error *out_err
);

/**
 * @brief Returns the local address of a UDP connection.
 *
 * @param conn UDP connection handle.
 * @return New fat_String handle containing "host:port" (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_UdpConnLocalAddr(fat_UdpConn conn);

/**
 * @brief Returns the remote address of a UDP connection (empty if not connected).
 *
 * @param conn UDP connection handle.
 * @return New fat_String handle containing "host:port" or "" (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_UdpConnRemoteAddr(fat_UdpConn conn);

/**
 * @brief Closes a UDP connection.
 *
 * @param conn UDP connection handle.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_OTHER on failure.
 */
FATSTD_API fat_Status fat_UdpConnClose(fat_UdpConn conn, fat_Error *out_err);

#ifdef __cplusplus
} /* extern "C" */
#endif

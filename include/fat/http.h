#pragma once

/**
 * @file fat/http.h
 * @brief Basic HTTP client/server utilities.
 *
 * This module is backed by Go's net/http package.
 *
 * Design notes:
 * - C callbacks are a poor fit for net/http handlers; the server uses a static
 *   response and a request queue that can be polled from C.
 * - Recoverable failures (invalid URLs, network errors) return fat_Status and
 *   fat_Error.
 *
 * Contract violations (invalid handles, NULL out-params, NULL pointers where
 * forbidden) are fatal.
 */

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#include "fat/bytes.h"
#include "fat/error.h"
#include "fat/export.h"
#include "fat/handle.h"
#include "fat/status.h"
#include "fat/string.h"

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Opaque handle to an HTTP client.
 *
 * @note Ownership: free with fat_HttpClientFree.
 */
typedef fat_Handle fat_HttpClient;

/**
 * @brief Opaque handle to an HTTP response.
 *
 * @note Ownership: free with fat_HttpResponseFree.
 */
typedef fat_Handle fat_HttpResponse;

/**
 * @brief Opaque handle to an HTTP server.
 *
 * @note Ownership: close with fat_HttpServerClose.
 */
typedef fat_Handle fat_HttpServer;

/**
 * @brief Opaque handle to an HTTP request (from a server queue).
 *
 * @note Ownership: free with fat_HttpRequestFree.
 */
typedef fat_Handle fat_HttpRequest;

/**
 * @brief Creates a new HTTP client.
 *
 * @return New client handle (caller must fat_HttpClientFree).
 */
FATSTD_API fat_HttpClient fat_HttpClientNew(void);

/**
 * @brief Frees an HTTP client handle.
 *
 * @param client Client handle to free.
 */
FATSTD_API void fat_HttpClientFree(fat_HttpClient client);

/**
 * @brief Performs an HTTP GET request.
 *
 * @param client HTTP client handle.
 * @param url URL string (UTF-8, NUL-terminated).
 * @param out_resp Output: response handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_SYNTAX on invalid URL; FAT_ERR_OTHER otherwise.
 */
FATSTD_API fat_Status fat_HttpClientGetUTF8(
  fat_HttpClient client,
  const char *url,
  fat_HttpResponse *out_resp,
  fat_Error *out_err
);

/**
 * @brief Performs an HTTP POST request with a byte payload.
 *
 * @param client HTTP client handle.
 * @param url URL string (UTF-8, NUL-terminated).
 * @param content_type Content-Type header (UTF-8, may be NULL to omit).
 * @param body Request body bytes.
 * @param out_resp Output: response handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_SYNTAX on invalid URL; FAT_ERR_OTHER otherwise.
 */
FATSTD_API fat_Status fat_HttpClientPostBytesUTF8(
  fat_HttpClient client,
  const char *url,
  const char *content_type,
  fat_Bytes body,
  fat_HttpResponse *out_resp,
  fat_Error *out_err
);

/**
 * @brief Returns the HTTP status code for a response.
 *
 * @param resp Response handle.
 * @return Status code (e.g., 200).
 */
FATSTD_API int fat_HttpResponseStatus(fat_HttpResponse resp);

/**
 * @brief Returns the response body as bytes.
 *
 * @param resp Response handle.
 * @return New fat_Bytes handle containing the body (caller must fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_HttpResponseBody(fat_HttpResponse resp);

/**
 * @brief Returns the value of a response header.
 *
 * @param resp Response handle.
 * @param name Header name (UTF-8, NUL-terminated).
 * @return New fat_String handle containing the header value (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_HttpResponseHeaderGetUTF8(fat_HttpResponse resp, const char *name);

/**
 * @brief Frees an HTTP response handle.
 *
 * @param resp Response handle to free.
 */
FATSTD_API void fat_HttpResponseFree(fat_HttpResponse resp);

/**
 * @brief Creates and starts an HTTP server bound to an address.
 *
 * Pass "127.0.0.1:0" to bind to an ephemeral port.
 *
 * @param addr Listen address (UTF-8, NUL-terminated).
 * @param out_server Output: new server handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_SYNTAX on invalid address; FAT_ERR_OTHER otherwise.
 */
FATSTD_API fat_Status fat_HttpServerNewUTF8(const char *addr, fat_HttpServer *out_server, fat_Error *out_err);

/**
 * @brief Returns the server's bound address as a string.
 *
 * @param server Server handle.
 * @return New fat_String handle containing "host:port" (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_HttpServerAddr(fat_HttpServer server);

/**
 * @brief Sets the static response returned by the server.
 *
 * @param server Server handle.
 * @param status HTTP status code.
 * @param body Response body bytes.
 * @param content_type Content-Type header (UTF-8, may be NULL to omit).
 */
FATSTD_API void fat_HttpServerSetStaticResponse(
  fat_HttpServer server,
  int status,
  fat_Bytes body,
  const char *content_type
);

/**
 * @brief Returns the next queued request, optionally waiting.
 *
 * If no request arrives before the timeout, returns FAT_ERR_EOF and sets out_req to 0.
 * Pass a negative timeout to wait indefinitely, or 0 to poll.
 *
 * @param server Server handle.
 * @param timeout_ms Timeout in milliseconds.
 * @param out_req Output: request handle on success.
 * @param out_err Output: error handle on failure, 0 on success/timeout.
 * @return FAT_OK on success; FAT_ERR_EOF on timeout; FAT_ERR_OTHER on failure.
 */
FATSTD_API fat_Status fat_HttpServerNextRequest(
  fat_HttpServer server,
  int64_t timeout_ms,
  fat_HttpRequest *out_req,
  fat_Error *out_err
);

/**
 * @brief Returns the HTTP method for a request.
 *
 * @param req Request handle.
 * @return New fat_String handle (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_HttpRequestMethod(fat_HttpRequest req);

/**
 * @brief Returns the URL path for a request.
 *
 * @param req Request handle.
 * @return New fat_String handle (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_HttpRequestPath(fat_HttpRequest req);

/**
 * @brief Returns the request body as bytes.
 *
 * @param req Request handle.
 * @return New fat_Bytes handle containing the body (caller must fat_BytesFree).
 */
FATSTD_API fat_Bytes fat_HttpRequestBody(fat_HttpRequest req);

/**
 * @brief Returns the value of a request header.
 *
 * @param req Request handle.
 * @param name Header name (UTF-8, NUL-terminated).
 * @return New fat_String handle containing the header value (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_HttpRequestHeaderGetUTF8(fat_HttpRequest req, const char *name);

/**
 * @brief Frees an HTTP request handle.
 *
 * @param req Request handle to free.
 */
FATSTD_API void fat_HttpRequestFree(fat_HttpRequest req);

/**
 * @brief Closes an HTTP server.
 *
 * @param server Server handle.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_OTHER on failure.
 */
FATSTD_API fat_Status fat_HttpServerClose(fat_HttpServer server, fat_Error *out_err);

#ifdef __cplusplus
} /* extern "C" */
#endif

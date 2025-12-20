#include "fat/http.h"

#include "fatstd_go.h"

fat_HttpClient fat_HttpClientNew(void) {
  return (fat_HttpClient)fatstd_go_http_client_new();
}

void fat_HttpClientFree(fat_HttpClient client) {
  fatstd_go_http_client_free((uintptr_t)client);
}

fat_Status fat_HttpClientGetUTF8(
  fat_HttpClient client,
  const char *url,
  fat_HttpResponse *out_resp,
  fat_Error *out_err
) {
  return (fat_Status)fatstd_go_http_client_get_utf8((uintptr_t)client, (char *)url, (uintptr_t *)out_resp,
                                                    (uintptr_t *)out_err);
}

fat_Status fat_HttpClientPostBytesUTF8(
  fat_HttpClient client,
  const char *url,
  const char *content_type,
  fat_Bytes body,
  fat_HttpResponse *out_resp,
  fat_Error *out_err
) {
  return (fat_Status)fatstd_go_http_client_post_bytes_utf8((uintptr_t)client, (char *)url, (char *)content_type,
                                                           (uintptr_t)body, (uintptr_t *)out_resp,
                                                           (uintptr_t *)out_err);
}

int fat_HttpResponseStatus(fat_HttpResponse resp) {
  return (int)fatstd_go_http_response_status((uintptr_t)resp);
}

fat_Bytes fat_HttpResponseBody(fat_HttpResponse resp) {
  return (fat_Bytes)fatstd_go_http_response_body((uintptr_t)resp);
}

fat_String fat_HttpResponseHeaderGetUTF8(fat_HttpResponse resp, const char *name) {
  return (fat_String)fatstd_go_http_response_header_get_utf8((uintptr_t)resp, (char *)name);
}

void fat_HttpResponseFree(fat_HttpResponse resp) {
  fatstd_go_http_response_free((uintptr_t)resp);
}

fat_Status fat_HttpServerNewUTF8(const char *addr, fat_HttpServer *out_server, fat_Error *out_err) {
  return (fat_Status)fatstd_go_http_server_new_utf8((char *)addr, (uintptr_t *)out_server, (uintptr_t *)out_err);
}

fat_String fat_HttpServerAddr(fat_HttpServer server) {
  return (fat_String)fatstd_go_http_server_addr((uintptr_t)server);
}

void fat_HttpServerSetStaticResponse(
  fat_HttpServer server,
  int status,
  fat_Bytes body,
  const char *content_type
) {
  fatstd_go_http_server_set_static_response((uintptr_t)server, status, (uintptr_t)body, (char *)content_type);
}

fat_Status fat_HttpServerNextRequest(
  fat_HttpServer server,
  int64_t timeout_ms,
  fat_HttpRequest *out_req,
  fat_Error *out_err
) {
  return (fat_Status)fatstd_go_http_server_next_request((uintptr_t)server, timeout_ms, (uintptr_t *)out_req,
                                                        (uintptr_t *)out_err);
}

fat_String fat_HttpRequestMethod(fat_HttpRequest req) {
  return (fat_String)fatstd_go_http_request_method((uintptr_t)req);
}

fat_String fat_HttpRequestPath(fat_HttpRequest req) {
  return (fat_String)fatstd_go_http_request_path((uintptr_t)req);
}

fat_Bytes fat_HttpRequestBody(fat_HttpRequest req) {
  return (fat_Bytes)fatstd_go_http_request_body((uintptr_t)req);
}

fat_String fat_HttpRequestHeaderGetUTF8(fat_HttpRequest req, const char *name) {
  return (fat_String)fatstd_go_http_request_header_get_utf8((uintptr_t)req, (char *)name);
}

void fat_HttpRequestFree(fat_HttpRequest req) {
  fatstd_go_http_request_free((uintptr_t)req);
}

fat_Status fat_HttpServerClose(fat_HttpServer server, fat_Error *out_err) {
  return (fat_Status)fatstd_go_http_server_close((uintptr_t)server, (uintptr_t *)out_err);
}

# HTTP (C) Tutorial

This module exposes a small, C-friendly wrapper around Go's `net/http`. It keeps the API simple and explicit: clients make GET/POST requests, and servers reply with a static response while exposing a queue of received requests for inspection.

Why a static server? C callbacks are a poor fit for Go's handler model. Instead, you configure a response and read requests from a queue.

## HTTP GET

```c
#include "fat/bytes.h"
#include "fat/error.h"
#include "fat/http.h"
#include "fat/string.h"

fat_HttpClient client = fat_HttpClientNew();
fat_HttpResponse resp = 0;
fat_Error err = 0;

if (fat_HttpClientGetUTF8(client, "http://example.com", &resp, &err) != FAT_OK) {
  fat_ErrorFree(err);
  fat_HttpClientFree(client);
  return;
}

int status = fat_HttpResponseStatus(resp);
fat_Bytes body = fat_HttpResponseBody(resp);

/* use status + body (copy out with fat_BytesCopyOut) */

fat_BytesFree(body);
fat_HttpResponseFree(resp);
fat_HttpClientFree(client);
```

## HTTP POST

```c
#include "fat/bytes.h"
#include "fat/error.h"
#include "fat/http.h"

fat_HttpClient client = fat_HttpClientNew();
fat_Bytes payload = fat_BytesNewN(json_bytes, json_len);
fat_HttpResponse resp = 0;
fat_Error err = 0;

if (fat_HttpClientPostBytesUTF8(
      client,
      "http://example.com/submit",
      "application/json",
      payload,
      &resp,
      &err
    ) != FAT_OK) {
  fat_ErrorFree(err);
  fat_BytesFree(payload);
  fat_HttpClientFree(client);
  return;
}

fat_BytesFree(payload);
fat_HttpResponseFree(resp);
fat_HttpClientFree(client);
```

## HTTP server (static response + request queue)

```c
#include "fat/bytes.h"
#include "fat/error.h"
#include "fat/http.h"
#include "fat/string.h"

fat_HttpServer server = 0;
fat_Error err = 0;
if (fat_HttpServerNewUTF8("127.0.0.1:0", &server, &err) != FAT_OK) {
  fat_ErrorFree(err);
  return;
}

fat_String addr = fat_HttpServerAddr(server);
/* convert addr to a C string with fat_StringCopyOutCStr */
fat_StringFree(addr);

fat_Bytes resp_body = fat_BytesNewN("ok", 2);
fat_HttpServerSetStaticResponse(server, 200, resp_body, "text/plain");
fat_BytesFree(resp_body);

/* wait up to 1000 ms for a request */
fat_HttpRequest req = 0;
if (fat_HttpServerNextRequest(server, 1000, &req, &err) == FAT_OK) {
  fat_String method = fat_HttpRequestMethod(req);
  fat_String path = fat_HttpRequestPath(req);
  fat_Bytes body = fat_HttpRequestBody(req);

  /* inspect method/path/body */

  fat_BytesFree(body);
  fat_StringFree(path);
  fat_StringFree(method);
  fat_HttpRequestFree(req);
}

fat_HttpServerClose(server, &err);
```

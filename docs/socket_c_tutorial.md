# Sockets (C) Tutorial

This module exposes a small, C-friendly wrapper around Go's `net` package. It supports basic TCP and UDP client/server flows with explicit ownership and error handles.

Key points:

- Address strings are standard `host:port` (UTF-8).
- `fat_TcpConnRead` returns `FAT_ERR_EOF` and sets `out_eof=true` on EOF.
- UDP reads return the sender address via a `fat_String` handle.

## TCP echo (server + client)

Two separate programs. Run the server first, then the client (same host or remote).

```c
// tcp_echo_server.c
#include <stdio.h>

#include "fat/error.h"
#include "fat/socket.h"
#include "fat/string.h"

fat_TcpListener listener = 0;
fat_Error err = 0;
if (fat_TcpListenerListenUTF8("0.0.0.0:9000", &listener, &err) != FAT_OK) {
  /* handle error */
  fat_ErrorFree(err);
  return;
}

fat_String addr = fat_TcpListenerAddr(listener);
/* convert addr to C string with fat_StringCopyOutCStr to log */
fat_StringFree(addr);

for (;;) {
  fat_TcpConn conn = 0;
  if (fat_TcpListenerAccept(listener, &conn, &err) != FAT_OK) {
    fat_ErrorFree(err);
    break;
  }

  char buf[1024];
  size_t n = 0;
  bool eof = false;
  if (fat_TcpConnRead(conn, buf, sizeof(buf), &n, &eof, &err) != FAT_OK) {
    fat_ErrorFree(err);
  } else if (!eof) {
    size_t written = 0;
    if (fat_TcpConnWrite(conn, buf, n, &written, &err) != FAT_OK) {
      fat_ErrorFree(err);
    }
  }

  fat_TcpConnClose(conn, &err);
}

fat_TcpListenerClose(listener, &err);
```

```c
// tcp_echo_client.c
#include <stdio.h>

#include "fat/error.h"
#include "fat/socket.h"

fat_TcpConn conn = 0;
fat_Error err = 0;
if (fat_TcpDialUTF8("127.0.0.1:9000", &conn, &err) != FAT_OK) {
  fat_ErrorFree(err);
  return;
}

const char msg[] = "hello from client";
size_t written = 0;
if (fat_TcpConnWrite(conn, msg, sizeof(msg) - 1, &written, &err) != FAT_OK) {
  fat_ErrorFree(err);
  return;
}

char buf[1024];
size_t n = 0;
bool eof = false;
if (fat_TcpConnRead(conn, buf, sizeof(buf), &n, &eof, &err) != FAT_OK) {
  fat_ErrorFree(err);
  return;
}

printf("echo: %.*s\n", (int)n, buf);
fat_TcpConnClose(conn, &err);
```

## UDP echo (server + client)

Two separate programs. Start the server, then the client. Use the server host in the client address.

```c
// udp_echo_server.c
#include <stdio.h>

#include "fat/error.h"
#include "fat/socket.h"
#include "fat/string.h"

fat_UdpConn server = 0;
fat_Error err = 0;
if (fat_UdpListenUTF8("0.0.0.0:9001", &server, &err) != FAT_OK) {
  fat_ErrorFree(err);
  return;
}

for (;;) {
  char buf[512];
  size_t n = 0;
  fat_String sender = 0;
  if (fat_UdpConnReadFrom(server, buf, sizeof(buf), &n, &sender, &err) != FAT_OK) {
    fat_ErrorFree(err);
    break;
  }

  /* copy sender into a C string with fat_StringCopyOutCStr */
  fat_StringFree(sender);
  /* echo is done via fat_UdpConnWriteToUTF8 using the sender address */
}

fat_UdpConnClose(server, &err);
```

```c
// udp_echo_client.c
#include <stdio.h>

#include "fat/error.h"
#include "fat/socket.h"

fat_UdpConn client = 0;
fat_Error err = 0;
if (fat_UdpDialUTF8("127.0.0.1:9001", &client, &err) != FAT_OK) {
  fat_ErrorFree(err);
  return;
}

const char msg[] = "hello udp";
size_t written = 0;
if (fat_UdpConnWrite(client, msg, sizeof(msg) - 1, &written, &err) != FAT_OK) {
  fat_ErrorFree(err);
  return;
}

char buf[512];
size_t n = 0;
fat_String sender = 0;
if (fat_UdpConnReadFrom(client, buf, sizeof(buf), &n, &sender, &err) != FAT_OK) {
  fat_ErrorFree(err);
  return;
}
fat_StringFree(sender);
printf("echo: %.*s\n", (int)n, buf);

fat_UdpConnClose(client, &err);
```

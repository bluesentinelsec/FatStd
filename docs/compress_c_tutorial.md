# Compression (C) Tutorial

This module exposes C-friendly wrappers around Go's `compress/*` packages. The API works with `fat_Bytes` handles and returns `fat_Status` + `fat_Error` for recoverable failures.

Key points:

- `fat_Flate*` operates on **raw DEFLATE** streams (no zlib/gzip wrapper).
- `fat_Zlib*` and `fat_Gzip*` include their respective container formats.
- `fat_Lzw*` requires explicit bit order and literal width.
- `fat_Bzip2*` supports **decompression only** (Go's `compress/bzip2` has no compressor).

## Gzip round-trip

```c
#include "fat/bytes.h"
#include "fat/error.h"
#include "fat/gzip.h"
#include "fat/string.h"

fat_Bytes src = fat_BytesNewN(payload, payload_len);
fat_Bytes compressed = 0;
fat_Error err = 0;

if (fat_GzipCompress(src, &compressed, &err) != FAT_OK) {
  fat_String msg = fat_ErrorMessage(err);
  /* handle error */
  fat_StringFree(msg);
  fat_ErrorFree(err);
}

fat_Bytes decompressed = 0;
if (fat_GzipDecompress(compressed, &decompressed, &err) != FAT_OK) {
  /* handle error */
  fat_ErrorFree(err);
}

/* use decompressed bytes */
fat_BytesFree(decompressed);
fat_BytesFree(compressed);
fat_BytesFree(src);
```

## Raw DEFLATE (flate)

```c
#include "fat/flate.h"

fat_Bytes compressed = 0;
fat_Error err = 0;
fat_Status st = fat_FlateCompress(src, &compressed, &err);
if (st != FAT_OK) {
  /* handle error */
}
```

## LZW parameters

LZW requires a bit order and literal width (2..8). Use one of:

- `FAT_LZW_ORDER_LSB`
- `FAT_LZW_ORDER_MSB`

```c
#include "fat/lzw.h"

fat_Bytes compressed = 0;
fat_Error err = 0;
fat_Status st = fat_LzwCompress(src, FAT_LZW_ORDER_LSB, 8, &compressed, &err);
if (st == FAT_ERR_RANGE) {
  /* invalid order or literal width */
}
```

## bzip2 decompression

```c
#include "fat/bzip2.h"

fat_Bytes decompressed = 0;
fat_Error err = 0;
fat_Status st = fat_Bzip2Decompress(src, &decompressed, &err);
if (st != FAT_OK) {
  /* handle error */
}
```

<!--
  Base64 C tutorial for FatStd.
  Keep this aligned with docs/design.md and docs/error_strategy.md.
-->

# Base64 in C (FatStd tutorial)

FatStd’s base64 APIs live in `include/fat/base64.h` and are backed by Go’s `encoding/base64`.

Key design points:

- Handles are opaque (`fat_Base64Encoding`, `fat_Bytes`, `fat_String`, `fat_Error`).
- Contract violations are fail-fast (invalid handles, NULL out-params where forbidden).
- Recoverable failures return `fat_Status` and optionally a `fat_Error` handle (free with `fat_ErrorFree`).

## 1) Create an encoding

Base64 encodings are configured by a 64-byte alphabet (standard base64 shown below):

```c
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

#include "fat/base64.h"
#include "fat/bytes.h"
#include "fat/error.h"
#include "fat/string.h"

static void fatal_status(const char *what, fat_Status st, fat_Error err) {
  if (err) {
    fat_String msg = fat_ErrorMessage(err);
    char buf[512];
    size_t n = fat_StringCopyOutCStr(msg, buf, sizeof(buf));
    buf[n] = '\0';
    fat_StringFree(msg);
    fat_ErrorFree(err);
    fprintf(stderr, "%s failed: status=%d err=%s\n", what, (int)st, buf);
  } else {
    fprintf(stderr, "%s failed: status=%d\n", what, (int)st);
  }
  abort();
}

static fat_Base64Encoding new_std_encoding(void) {
  fat_Base64Encoding enc = 0;
  fat_Error err = 0;
  fat_Status st = fat_Base64EncodingNewUTF8(
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/",
    &enc,
    &err
  );
  if (st != FAT_OK) fatal_status("fat_Base64EncodingNewUTF8", st, err);
  return enc;
}
```

## 2) Encode bytes to a string

```c
fat_Base64Encoding enc = new_std_encoding();

const uint8_t raw[] = { 'h','i',0x00,0xff };
fat_Bytes src = fat_BytesNewN(raw, sizeof(raw));

fat_String b64 = fat_Base64EncodeToString(enc, src);

char out[256];
size_t n = fat_StringCopyOutCStr(b64, out, sizeof(out) - 1);
out[n] = '\0';
printf("base64: %s\n", out);

fat_StringFree(b64);
fat_BytesFree(src);
fat_Base64EncodingFree(enc);
```

## 3) Decode a string to bytes (with error handling)

```c
fat_Base64Encoding enc = new_std_encoding();
fat_String s = fat_StringNewUTF8("aGVsbG8=");

fat_Bytes decoded = 0;
fat_Error err = 0;
fat_Status st = fat_Base64DecodeString(enc, s, &decoded, &err);
if (st != FAT_OK) fatal_status("fat_Base64DecodeString", st, err);

size_t decoded_len = fat_BytesLen(decoded);
uint8_t *buf = (uint8_t *)malloc(decoded_len);
fat_BytesCopyOut(decoded, buf, decoded_len);

/* use buf... */
free(buf);

fat_BytesFree(decoded);
fat_StringFree(s);
fat_Base64EncodingFree(enc);
```

## 4) Strict mode and custom padding

Strict mode rejects some non-canonical encodings:

```c
fat_Base64Encoding enc = new_std_encoding();
fat_Base64Encoding strict = fat_Base64EncodingStrict(enc);
fat_Base64EncodingFree(enc);

/* strict is a new handle; free it when done */
fat_Base64EncodingFree(strict);
```

Disable padding (equivalent to Go’s `base64.NoPadding`) by passing `-1`:

```c
fat_Base64Encoding enc = new_std_encoding();

fat_Base64Encoding nopad = 0;
fat_Error err = 0;
fat_Status st = fat_Base64EncodingWithPadding(enc, -1, &nopad, &err);
if (st != FAT_OK) fatal_status("fat_Base64EncodingWithPadding", st, err);

fat_Base64EncodingFree(nopad);
fat_Base64EncodingFree(enc);
```

## 5) Streaming encoder (writes to a bytes buffer)

Go’s `NewEncoder`/`io.Writer` model doesn’t map directly to C, so FatStd provides a C-friendly alternative that writes to a `fat_BytesBuffer`:

```c
#include "fat/bytes_buffer.h"

fat_Base64Encoding enc = new_std_encoding();
fat_BytesBuffer buf = fat_BytesBufferNew();

fat_Base64Encoder e = 0;
fat_Error err = 0;
fat_Status st = fat_Base64EncoderNewToBytesBuffer(enc, buf, &e, &err);
if (st != FAT_OK) fatal_status("fat_Base64EncoderNewToBytesBuffer", st, err);

const char *msg = "streaming payload";
size_t n_written = 0;
fat_Error err2 = 0;
fat_Status st2 = fat_Base64EncoderWrite(e, msg, 17, &n_written, &err2);
if (st2 != FAT_OK) fatal_status("fat_Base64EncoderWrite", st2, err2);

fat_Error err3 = 0;
fat_Status st3 = fat_Base64EncoderClose(e, &err3);
if (st3 != FAT_OK) fatal_status("fat_Base64EncoderClose", st3, err3);

fat_Bytes out_bytes = fat_BytesBufferBytes(buf);
/* out_bytes contains the encoded base64 bytes */

fat_BytesFree(out_bytes);
fat_BytesBufferFree(buf);
fat_Base64EncodingFree(enc);
```


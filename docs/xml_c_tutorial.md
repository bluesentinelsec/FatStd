<!--
  XML C tutorial for FatStd.
  Keep this aligned with docs/design.md and docs/error_strategy.md.
-->

# XML in C (FatStd tutorial)

FatStd’s XML APIs live in `include/fat/xml.h` and are backed by Go’s `encoding/xml`.

The C API is streaming/token-based: you read tokens (`StartElement`, `EndElement`, `CharData`, …) and inspect them.

## Error handling pattern

```c
#include <stdio.h>
#include <stdlib.h>

#include "fat/error.h"
#include "fat/string.h"

static void die_status(const char *op, fat_Status st, fat_Error err) {
  if (err) {
    fat_String msg = fat_ErrorMessage(err);
    char buf[512];
    size_t n = fat_StringCopyOutCStr(msg, buf, sizeof(buf) - 1);
    buf[n] = '\0';
    fat_StringFree(msg);
    fat_ErrorFree(err);
    fprintf(stderr, "%s failed: status=%d err=%s\n", op, (int)st, buf);
  } else {
    fprintf(stderr, "%s failed: status=%d\n", op, (int)st);
  }
  abort();
}
```

## 1) Read an XML file and extract values

Example XML:

```xml
<root>
  <person id="7">
    <name>alice</name>
    <active>true</active>
  </person>
</root>
```

Token-scan for `<name>` and read the next `CharData`:

```c
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>

#include "fat/bytes.h"
#include "fat/error.h"
#include "fat/status.h"
#include "fat/string.h"
#include "fat/xml.h"

static void copy_bytes_to_cstr(fat_Bytes b, char *dst, size_t dst_cap) {
  size_t n = fat_BytesLen(b);
  if (n + 1 > dst_cap) n = dst_cap - 1;
  fat_BytesCopyOut(b, dst, n);
  dst[n] = '\0';
}

void read_name_from_file(const char *path_utf8) {
  fat_XmlDecoder dec = 0;
  fat_Error err = 0;
  fat_Status st = fat_XmlDecoderOpenPathUTF8(path_utf8, &dec, &err);
  if (st != FAT_OK) die_status("fat_XmlDecoderOpenPathUTF8", st, err);

  bool in_name = false;
  char name_buf[256] = {0};

  for (;;) {
    fat_XmlToken tok = 0;
    fat_Error tok_err = 0;
    fat_Status st_tok = fat_XmlDecoderToken(dec, &tok, &tok_err);
    if (st_tok == FAT_ERR_EOF) break;
    if (st_tok != FAT_OK) die_status("fat_XmlDecoderToken", st_tok, tok_err);

    fat_XmlTokenKind kind = fat_XmlTokenType(tok);
    if (kind == FAT_XML_START_ELEMENT) {
      fat_String local_s = fat_XmlTokenNameLocal(tok);
      char local[128];
      size_t n = fat_StringCopyOutCStr(local_s, local, sizeof(local) - 1);
      local[n] = '\0';
      fat_StringFree(local_s);
      in_name = (strcmp(local, "name") == 0);
    } else if (kind == FAT_XML_END_ELEMENT) {
      in_name = false;
    } else if (kind == FAT_XML_CHAR_DATA && in_name) {
      fat_Bytes text = fat_XmlTokenBytes(tok);
      copy_bytes_to_cstr(text, name_buf, sizeof(name_buf));
      fat_BytesFree(text);
    }

    fat_XmlTokenFree(tok);
  }

  printf("name=%s\n", name_buf);

  fat_Error close_err = 0;
  fat_Status st_close = fat_XmlDecoderFree(dec, &close_err);
  if (st_close != FAT_OK) die_status("fat_XmlDecoderFree", st_close, close_err);
}
```

Notes:

- Always free `fat_XmlToken` handles with `fat_XmlTokenFree`.
- `fat_XmlTokenBytes` returns a new `fat_Bytes` handle; free it with `fat_BytesFree`.

## 2) Read element attributes

When you see a `FAT_XML_START_ELEMENT`, you can inspect attributes:

```c
void read_person_id_attr(fat_XmlToken start_elem_tok) {
  size_t n = fat_XmlStartElementAttrCount(start_elem_tok);
  for (size_t i = 0; i < n; i++) {
    fat_String name_local = 0, name_space = 0, value = 0;
    fat_XmlStartElementAttrGet(start_elem_tok, i, &name_local, &name_space, &value);
    /* ... compare name_local/value ... */
    fat_StringFree(value);
    fat_StringFree(name_space);
    fat_StringFree(name_local);
  }
}
```

## 3) Escape XML text

`fat_XmlEscapeToBytesBuffer` and `fat_XmlEscapeTextToBytesBuffer` write escaped bytes to a `fat_BytesBuffer`:

```c
#include "fat/bytes_buffer.h"

void escape_example(void) {
  fat_BytesBuffer buf = fat_BytesBufferNew();
  fat_Bytes src = fat_BytesNewN("<a&b>", 5);

  fat_XmlEscapeToBytesBuffer(buf, src);

  fat_Bytes out = fat_BytesBufferBytes(buf);
  /* ... copy out with fat_BytesLen + fat_BytesCopyOut ... */

  fat_BytesFree(out);
  fat_BytesFree(src);
  fat_BytesBufferFree(buf);
}
```

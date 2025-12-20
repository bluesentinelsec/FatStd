<!--
  JSON C tutorial for FatStd.
  Keep this aligned with docs/design.md and docs/error_strategy.md.
-->

# JSON in C (FatStd tutorial)

FatStd’s JSON APIs live in `include/fat/json.h` and are backed by Go’s `encoding/json`.

Key design points:

- Handles are opaque (`fat_JsonValue`, `fat_JsonDecoder`, `fat_JsonEncoder`, `fat_Error`, etc.).
- Contract violations are fail-fast (invalid handles, NULL out-params where forbidden).
- Recoverable failures return `fat_Status` and optionally a `fat_Error` handle (free with `fat_ErrorFree`).

## 1) Validate and reformat raw JSON bytes

```c
#include <stdio.h>
#include <stdlib.h>

#include "fat/bytes.h"
#include "fat/error.h"
#include "fat/json.h"
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

void compact_indent_example(const void *json_bytes, size_t json_len) {
  fat_Bytes src = fat_BytesNewN(json_bytes, json_len);

  if (!fat_JsonValid(src)) {
    fprintf(stderr, "invalid JSON\n");
    fat_BytesFree(src);
    return;
  }

  fat_Bytes compacted = 0;
  fat_Error err = 0;
  fat_Status st = fat_JsonCompact(src, &compacted, &err);
  if (st != FAT_OK) die_status("fat_JsonCompact", st, err);

  fat_String prefix = fat_StringNewUTF8("");
  fat_String indent = fat_StringNewUTF8("  ");

  fat_Bytes pretty = 0;
  fat_Error err2 = 0;
  fat_Status st2 = fat_JsonIndent(src, prefix, indent, &pretty, &err2);
  if (st2 != FAT_OK) die_status("fat_JsonIndent", st2, err2);

  fat_BytesFree(pretty);
  fat_BytesFree(compacted);
  fat_StringFree(indent);
  fat_StringFree(prefix);
  fat_BytesFree(src);
}
```

## 2) Parse JSON into a generic value and inspect it

FatStd exposes a handle-backed generic JSON value (object/array/string/number/bool/null).

### Common pattern: parse an object and read a field

```c
#include <stdbool.h>
#include <stddef.h>
#include <stdio.h>
#include <string.h>

#include "fat/bytes.h"
#include "fat/json.h"
#include "fat/string.h"

void read_field_example(void) {
  const char *json_text = "{\"name\":\"alice\",\"active\":true}";
  fat_Bytes data = fat_BytesNewN(json_text, strlen(json_text));

  fat_JsonValue root = 0;
  fat_Error err = 0;
  fat_Status st = fat_JsonUnmarshal(data, &root, &err);
  if (st != FAT_OK) die_status("fat_JsonUnmarshal", st, err);

  if (fat_JsonValueType(root) != FAT_JSON_OBJECT) {
    fprintf(stderr, "expected object\n");
    fat_JsonValueFree(root);
    fat_BytesFree(data);
    return;
  }

  fat_String key = fat_StringNewUTF8("name");
  bool found = false;
  fat_JsonValue name_v = 0;
  fat_JsonObjectGet(root, key, &found, &name_v);
  fat_StringFree(key);

  if (found) {
    /* name_v is a new handle; you own it */
    fat_String name_s = fat_JsonValueAsString(name_v);

    char out[256];
    size_t n = fat_StringCopyOutCStr(name_s, out, sizeof(out) - 1);
    out[n] = '\0';
    printf("name=%s\n", out);

    fat_StringFree(name_s);
    fat_JsonValueFree(name_v);
  }

  fat_JsonValueFree(root);
  fat_BytesFree(data);
}
```

### Notes

- Use `fat_JsonObjectKeys` to enumerate object keys (sorted) and `fat_JsonObjectGet` to fetch a nested value by key.
- Use `fat_JsonArrayLen` + `fat_JsonArrayGet` to traverse arrays.
- Use `fat_JsonValueAsNumberString` if you need the original number text; `fat_JsonUnmarshal` decodes numbers as JSON numbers (not float64).

```c
#include <stdbool.h>

#include "fat/json.h"
#include "fat/string.h"

void parse_and_inspect(const void *json_bytes, size_t json_len) {
  fat_Bytes data = fat_BytesNewN(json_bytes, json_len);

  fat_JsonValue root = 0;
  fat_Error err = 0;
  fat_Status st = fat_JsonUnmarshal(data, &root, &err);
  if (st != FAT_OK) die_status("fat_JsonUnmarshal", st, err);

  if (fat_JsonValueType(root) == FAT_JSON_OBJECT) {
    fat_StringArray keys = fat_JsonObjectKeys(root); /* sorted */
    /* iterate keys with fat_StringArrayLen/Get */
    fat_StringArrayFree(keys);

    fat_String key = fat_StringNewUTF8("someKey");
    bool found = false;
    fat_JsonValue child = 0;
    fat_JsonObjectGet(root, key, &found, &child);
    fat_StringFree(key);
    if (found) {
      /* ... use child ... */
      fat_JsonValueFree(child);
    }
  }

  fat_JsonValueFree(root);
  fat_BytesFree(data);
}
```

## 3) Marshal a value back to JSON

```c
fat_Bytes out = 0;
fat_Error err = 0;
fat_Status st = fat_JsonMarshal(root, &out, &err);
if (st != FAT_OK) die_status("fat_JsonMarshal", st, err);
/* out is JSON bytes */
fat_BytesFree(out);
```

## 4) Streaming decode/encode

You can decode multiple JSON values from a stream and encode values to a bytes buffer:

```c
#include "fat/bytes_buffer.h"
#include "fat/bytes_reader.h"

void streaming_example(fat_Bytes input) {
  fat_BytesReader r = fat_BytesReaderNew(input);
  fat_JsonDecoder dec = fat_JsonDecoderNewBytesReader(r);

  fat_JsonDecoderUseNumber(dec);

  while (true) {
    fat_JsonValue v = 0;
    fat_Error err = 0;
    fat_Status st = fat_JsonDecoderDecodeValue(dec, &v, &err);
    if (st == FAT_ERR_EOF) break;
    if (st != FAT_OK) die_status("fat_JsonDecoderDecodeValue", st, err);

    fat_BytesBuffer buf = fat_BytesBufferNew();
    fat_JsonEncoder enc = fat_JsonEncoderNewToBytesBuffer(buf);

    fat_Error err2 = 0;
    fat_Status st2 = fat_JsonEncoderEncodeValue(enc, v, &err2);
    if (st2 != FAT_OK) die_status("fat_JsonEncoderEncodeValue", st2, err2);

    fat_JsonEncoderFree(enc);
    fat_BytesBufferFree(buf);
    fat_JsonValueFree(v);
  }

  fat_JsonDecoderFree(dec);
  fat_BytesReaderFree(r);
}
```

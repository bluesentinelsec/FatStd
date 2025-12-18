<!--
  FatStd fat bytes tutorial.
  Keep this aligned with docs/design.md and the public headers under include/fat/.
-->

# Fat Bytes (Tutorial)

`fat_bytes` is the byte-slice companion to fat strings: a handle-based, Go-backed API for working with arbitrary bytes from C.

Key properties:

- **Arbitrary bytes**: not necessarily UTF-8; embedded NULs are normal.
- **Explicit ownership**: if an API returns a `fat_Bytes`/`fat_BytesArray`/`fat_BytesBuffer`/`fat_BytesReader`, you own it and must free it with the matching `*_Free`.
- **Fail-fast**: invalid handles, NULL pointers (where not allowed), and contract violations are fatal.

Public headers:

- `include/fat/bytes.h`
- `include/fat/bytes_buffer.h`
- `include/fat/bytes_reader.h`

## 1) Creating and freeing bytes

`fat_Bytes` is an opaque handle to a Go-backed `[]byte`.

```c
#include "fat/bytes.h"

const unsigned char data[] = { 0x61, 0x62, 0x00, 0x63 };
fat_Bytes b = fat_BytesNewN(data, sizeof(data));
/* ... */
fat_BytesFree(b);
```

## 2) Copying bytes out (C interop)

Use `fat_BytesLen` + `fat_BytesCopyOut` to get `(buf, len)` in C.

```c
#include "fat/bytes.h"
#include <stdlib.h>

size_t n = fat_BytesLen(b);
unsigned char *buf = (unsigned char *)malloc(n);
fat_BytesCopyOut(b, buf, n);

/* Use as a byte span: (buf, n) */

free(buf);
```

## 3) Predicates and basic operations

```c
fat_Bytes hay = fat_BytesNewN("abc\0def", 7);
fat_Bytes needle = fat_BytesNewN("\0d", 2);

bool has = fat_BytesContains(hay, needle);

fat_BytesFree(needle);
fat_BytesFree(hay);
```

Case conversion is ASCII-only (matches Go `bytes.ToLower/ToUpper`):

```c
fat_Bytes mixed = fat_BytesNewN("AbC", 3);
fat_Bytes lower = fat_BytesToLower(mixed);
fat_Bytes upper = fat_BytesToUpper(mixed);
fat_BytesFree(upper);
fat_BytesFree(lower);
fat_BytesFree(mixed);
```

## 4) Arrays: split/fields/join

Some APIs return `fat_BytesArray`. Free the array handle with `fat_BytesArrayFree`.

Important rule: `fat_BytesArrayGet` returns a **new** `fat_Bytes` handle each time; you must free each element handle you obtain.

```c
#include "fat/bytes.h"

fat_Bytes s = fat_BytesNewN("a,b,c", 5);
fat_Bytes sep = fat_BytesNewN(",", 1);

fat_BytesArray parts = fat_BytesSplit(s, sep);
size_t n = fat_BytesArrayLen(parts);

for (size_t i = 0; i < n; i++) {
  fat_Bytes elem = fat_BytesArrayGet(parts, i);
  /* ... */
  fat_BytesFree(elem);
}

fat_Bytes joined = fat_BytesJoin(parts, sep);

fat_BytesFree(joined);
fat_BytesArrayFree(parts);
fat_BytesFree(sep);
fat_BytesFree(s);
```

`fat_BytesFields` splits on ASCII whitespace (matches Go `bytes.Fields`).

## 5) Cut / TrimPrefix / TrimSuffix

`fat_BytesCut` mirrors Go `bytes.Cut`, returning a boolean and producing output handles.

```c
fat_Bytes s = fat_BytesNewN("foo=bar", 7);
fat_Bytes sep = fat_BytesNewN("=", 1);

fat_Bytes before = 0;
fat_Bytes after = 0;
bool found = fat_BytesCut(s, sep, &before, &after);

/* before/after are always new handles; caller frees them. */
fat_BytesFree(after);
fat_BytesFree(before);
fat_BytesFree(sep);
fat_BytesFree(s);
```

## 6) BytesBuffer: incremental read/write

`fat_BytesBuffer` is the bytes analogue of a string builder, but also supports reading.

```c
#include "fat/bytes_buffer.h"

fat_BytesBuffer buf = fat_BytesBufferNew();
fat_BytesBufferWrite(buf, "ab\0", 3);
fat_BytesBufferWriteByte(buf, (uint8_t)'c');

fat_Bytes snap = fat_BytesBufferBytes(buf); /* snapshot copy */
fat_BytesFree(snap);

fat_BytesBufferFree(buf);
```

Notes:

- `fat_BytesBufferBytes` and `fat_BytesBufferNext` return **new** `fat_Bytes` handles (no borrowed slices).
- `fat_BytesBufferString` returns a **new** `fat_String` handle.
- `WriteTo(io.Writer)` / `ReadFrom(io.Reader)` are exposed only as FatStd-specific variants:
  - `fat_BytesBufferWriteToBytesBuffer(src, dst)`
  - `fat_BytesBufferReadFromStringReader(dst, reader)`

## 7) BytesReader: reading from bytes

`fat_BytesReader` wraps Go `bytes.Reader` with C-friendly EOF handling.

EOF is reported via `eof_out`, not by returning an error code.

```c
#include "fat/bytes_reader.h"

fat_Bytes b = fat_BytesNewN("hello", 5);
fat_BytesReader r = fat_BytesReaderNew(b);

char tmp[2];
bool eof = false;
size_t n = fat_BytesReaderRead(r, tmp, sizeof(tmp), &eof);

fat_BytesReaderFree(r);
fat_BytesFree(b);
```

`WriteTo(io.Writer)` is exposed as `fat_BytesReaderWriteToBytesBuffer(reader, buffer)`.

## 8) Practical tips

- Treat `(buf, len)` as the primary C representation; don’t assume NUL-termination.
- Use `fat_String` only when you need string semantics; otherwise prefer `fat_Bytes`.

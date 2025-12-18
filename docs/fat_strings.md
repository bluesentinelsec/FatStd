<!--
  FatStd fat strings tutorial.
  Keep this aligned with docs/design.md and the public headers under include/fat/.
-->

# Fat Strings (Tutorial)

FatStd strings are **Go-backed**, **handle-based** strings exposed through a C API. You never receive a pointer to Go memory; you receive an opaque handle (`fat_String`) and call `fat_String*` functions.

Key properties:

- **Explicit ownership**: if an API returns a `fat_String`, you own it and must call `fat_StringFree`.
- **Fail-fast**: invalid handles, NULL pointers (where not allowed), and contract violations are fatal.
- **Bytes, not C strings**: fat strings may contain embedded NULs (`\0`), especially when created from explicit byte spans.

## 1) Include and link

Public headers:

- `include/fat/string.h`
- `include/fat/string_builder.h`
- `include/fat/string_reader.h`

## 2) Creating and freeing strings

### From a C string (no embedded NULs)

```c
#include "fat/string.h"

fat_String s = fat_StringNewUTF8("hello");
/* ... */
fat_StringFree(s);
```

### From bytes (embedded NULs allowed)

```c
#include "fat/string.h"

const char bytes[] = { 'a', 'b', 'c', '\0', 'd', 'e', 'f' };
fat_String s = fat_StringNewUTF8N(bytes, sizeof(bytes));
fat_StringFree(s);
```

## 3) Copying bytes out (C interop)

Fat strings are not directly viewable from C. Use `fat_StringLenBytes` + `fat_StringCopyOut`.

```c
#include "fat/string.h"
#include <stdlib.h>

size_t n = fat_StringLenBytes(s);
unsigned char *buf = (unsigned char *)malloc(n);
fat_StringCopyOut(s, buf, n);

/* Use as a byte span: (buf, n) */

free(buf);
```

If you want a NUL-terminated buffer for C APIs, prefer `fat_StringCopyOutCStr`:

```c
size_t n = fat_StringLenBytes(s);
char *cstr = (char *)malloc(n + 1);
fat_StringCopyOutCStr(s, cstr, n + 1);

/* NOTE: embedded NULs will truncate when treated as a C string. */

free(cstr);
```

## 4) Common string operations

Most APIs follow Go’s `strings` package semantics. Outputs are typically new `fat_String` handles.

```c
fat_String s = fat_StringNewUTF8("  Lorem Ipsum  ");
fat_String trimmed = fat_StringTrimSpace(s);
fat_StringFree(trimmed);
fat_StringFree(s);
```

Predicates return `bool`:

```c
bool ok = fat_StringHasPrefix(s, prefix);
```

## 5) Arrays: split/fields/join

Some APIs return `fat_StringArray`. You must free the array handle.

Important rule: `fat_StringArrayGet` returns a **new** `fat_String` handle each time; you must free each element handle you obtain.

```c
#include "fat/string.h"

fat_String s = fat_StringNewUTF8("a,b,c");
fat_String sep = fat_StringNewUTF8(",");

fat_StringArray parts = fat_StringSplit(s, sep);
size_t n = fat_StringArrayLen(parts);

for (size_t i = 0; i < n; i++) {
  fat_String elem = fat_StringArrayGet(parts, i);
  /* ... */
  fat_StringFree(elem);
}

fat_String joined = fat_StringJoin(parts, sep);

fat_StringFree(joined);
fat_StringArrayFree(parts);
fat_StringFree(sep);
fat_StringFree(s);
```

## 6) Builder: efficient incremental construction

Use `fat_StringBuilder` to avoid repeated allocations when concatenating many pieces.

```c
#include "fat/string.h"
#include "fat/string_builder.h"

fat_StringBuilder b = fat_StringBuilderNew();

fat_String hello = fat_StringNewUTF8("hello");
fat_StringBuilderWriteString(b, hello);
fat_StringBuilderWriteByte(b, (uint8_t)' ');
fat_StringBuilderWrite(b, "world", 5);

fat_String out = fat_StringBuilderString(b);

fat_StringFree(out);
fat_StringFree(hello);
fat_StringBuilderFree(b);
```

Notes:

- Builder methods are not thread-safe.
- `fat_StringBuilderString` returns a new `fat_String` handle; the builder remains usable.

## 7) Reader: reading bytes from a string

Use `fat_StringReader` for sequential reads, random reads (`ReadAt`), and seeking.

EOF is reported via `eof_out`, not by returning a special error code.

```c
#include "fat/string.h"
#include "fat/string_reader.h"

fat_String s = fat_StringNewUTF8("hello");
fat_StringReader r = fat_StringReaderNew(s);

char buf[4];
bool eof = false;
size_t n = fat_StringReaderRead(r, buf, sizeof(buf), &eof);

fat_StringReaderFree(r);
fat_StringFree(s);
```

### WriteTo: C-friendly form

Go’s `(*strings.Reader).WriteTo(io.Writer)` can’t be exposed directly to C. FatStd provides:

- `fat_StringReaderWriteToBuilder(reader, builder)`

which writes remaining unread bytes into a `fat_StringBuilder`.

## 8) Error model and debugging

FatStd treats misuse as programmer error and fails fast (panic/abort behavior depends on build/runtime). Common fatal mistakes:

- Using a freed handle.
- Passing NULL where not allowed.
- Calling `fat_StringReaderUnreadByte` without a preceding `ReadByte`.
- Invalid seek parameters.

## 9) Practical tips

- Prefer `fat_StringNewUTF8` for ergonomic C-string input.
- Prefer `fat_StringNewUTF8N` for binary data or explicit-length spans.
- Treat `(buf, len)` as the real “C representation” of a fat string; only convert to NUL-terminated buffers when necessary.

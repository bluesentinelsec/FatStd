<!--
  CSV C tutorial for FatStd.
  Keep this aligned with docs/design.md and docs/error_strategy.md.
-->

# CSV in C (FatStd tutorial)

FatStd’s CSV APIs live in `include/fat/csv.h` and are backed by Go’s `encoding/csv`.

Key design points:

- Handles are opaque (`fat_CsvReader`, `fat_CsvWriter`, `fat_StringArray`, `fat_Error`).
- Contract violations are fail-fast (invalid handles, NULL out-params where forbidden).
- Recoverable failures return `fat_Status` and optionally a `fat_Error` handle (free with `fat_ErrorFree`).

## Writing CSV to memory

Use a `fat_BytesBuffer` as the destination and snapshot it with `fat_BytesBufferBytes`.

```c
#include <stdio.h>

#include "fat/bytes_buffer.h"
#include "fat/csv.h"
#include "fat/error.h"
#include "fat/string.h"

static void print_err_and_abort(const char *op, fat_Status st, fat_Error err) {
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

void write_example(void) {
  fat_BytesBuffer buf = fat_BytesBufferNew();
  fat_CsvWriter w = fat_CsvWriterNewToBytesBuffer(buf);

  fat_String f0 = fat_StringNewUTF8("a,b");
  fat_String f1 = fat_StringNewUTF8("c\"d");
  fat_String fields[] = { f0, f1 };

  fat_Error err = 0;
  fat_Status st = fat_CsvWriterWriteRecord(w, fields, 2, &err);
  if (st != FAT_OK) print_err_and_abort("fat_CsvWriterWriteRecord", st, err);

  fat_CsvWriterFlush(w);
  fat_Error err2 = 0;
  fat_Status st2 = fat_CsvWriterError(w, &err2);
  if (st2 != FAT_OK) print_err_and_abort("fat_CsvWriterError", st2, err2);

  fat_Bytes out = fat_BytesBufferBytes(buf);
  /* out now contains the CSV bytes. Copy out if needed. */

  fat_BytesFree(out);
  fat_StringFree(f0);
  fat_StringFree(f1);
  fat_CsvWriterFree(w);
  fat_BytesBufferFree(buf);
}
```

## Reading CSV from memory

FatStd does not expose Go’s `io.Reader` directly. Instead, create a reader over a `fat_Bytes` blob and iterate record-by-record.

```c
#include <stdbool.h>

#include "fat/bytes.h"
#include "fat/csv.h"
#include "fat/string.h"

void read_example(const void *csv_bytes, size_t csv_len) {
  fat_Bytes b = fat_BytesNewN(csv_bytes, csv_len);
  fat_CsvReader r = fat_CsvReaderNewBytes(b);

  while (true) {
    fat_StringArray record = 0;
    bool eof = false;
    fat_Error err = 0;
    fat_Status st = fat_CsvReaderRead(r, &record, &eof, &err);

    if (st == FAT_ERR_EOF) {
      break;
    }
    if (st != FAT_OK) {
      /* err contains details */
      fat_ErrorFree(err);
      break;
    }

    size_t n = fat_StringArrayLen(record);
    for (size_t i = 0; i < n; i++) {
      fat_String field = fat_StringArrayGet(record, i);
      /* ... use field ... */
      fat_StringFree(field);
    }

    fat_StringArrayFree(record);
  }

  fat_CsvReaderFree(r);
  fat_BytesFree(b);
}
```

## Field position and input offset

After a successful `fat_CsvReaderRead`, you can query:

- `fat_CsvReaderFieldPos(r, field, &line, &col)` for where a field started
- `fat_CsvReaderInputOffset(r)` for the current byte offset in the input

These match Go’s `(*csv.Reader).FieldPos` and `(*csv.Reader).InputOffset`.


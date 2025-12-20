<!--
  FatStd error handling strategy.
  Keep this aligned with docs/design.md and docs/onboarding_and_testing_functions.md.
-->

# Error Strategy (SOP)

FatStd distinguishes between:

1. **Programmer / contract violations** (bugs)
2. **Recoverable runtime failures** (expected in normal use)

The two categories use different API conventions.

## 1) Programmer / contract violations (fail-fast)

These are treated as programmer error and are **fatal** by default.

Examples:

- Invalid or freed handle
- Double free / use-after-free
- NULL pointer where not permitted
- Invalid `whence` or invalid seek math (negative position)
- Out-of-range indices (`ArrayGet`, etc.)
- Length/size constraints violated (impossibly large sizes)

Contract violations are not returned to the caller as status codes. They are diagnosed eagerly and fail fast.

This matches `docs/design.md` (explicit ownership + fail-fast) and keeps the “happy path” simple for low-level, deterministic APIs (strings/bytes).

## 2) Recoverable runtime failures (bubble up to caller)

Some modules have failures that are **normal** in correct programs (HTTP, filesystem, process, parsing/decoding, etc.).

For these APIs, FatStd will expose explicit error returns so callers can handle failures without crashing.

### 2.1 Contract

Recommended C signature pattern:

```c
fat_Status fat_ModuleOp(..., out_type *out_result, fat_Error *out_err);
```

Where:

- `fat_Status` is a small enum-like status code (OK vs error categories).
- `fat_Error` is an opaque handle carrying detailed information (message and optional codes).

### 2.2 Output rules

- **On success** (`FAT_OK`):
  - `*out_result` is written (new handle or scalar value, per API).
  - `*out_err` is set to 0 (or left unchanged if the API documents that it never sets it on success).
- **On recoverable failure** (non-OK status):
  - `*out_result` is set to 0 / empty (or left unchanged if documented).
  - `*out_err` is set to a newly allocated error handle describing the failure.

### 2.3 Ownership

- Any returned handle (result or error) must have a matching `*_Free`.
- Error handles are freed by the caller (`fat_ErrorFree`).

### 2.4 EOF and partial results

EOF is not a “crash” condition. Prefer one of:

- `fat_Status` includes `FAT_EOF`, and read APIs may return `(FAT_EOF, n>0)` for “read some bytes then hit EOF”.
- Or, maintain explicit `eof_out` booleans for reader-like APIs.

Pick one per module and document it clearly.

## 3) Go-to-C mapping

Go is an internal implementation detail.

- Programmer/contract violations may `panic` (consistent with fail-fast).
- Recoverable Go `error` values are mapped to:
  - a `fat_Status` code, and
  - a `fat_Error` handle containing details.

Avoid leaking Go error types or Go memory across the boundary.

## 4) Documentation requirements

For public headers under `include/fat/`:

- Document which inputs are contract-checked (fatal on misuse).
- Document which failures are recoverable (returned to caller).
- Document ownership of all outputs, especially error handles.

Use Doxygen format per `docs/documenting_code.md`.


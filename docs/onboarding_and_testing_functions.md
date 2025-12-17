<!--
  FatStd function onboarding + test workflow.
  Keep this doc aligned with `docs/design.md`.
-->

# Onboarding and Testing Functions (Go + C)

This document defines the repeatable workflow for adding new APIs to `libfatstd` (the C library) and validating them end-to-end through the existing Python `ctypes` smoke tests.

FatStd is **C-first**: the public surface area is a stable, C-compatible API. Go is an internal implementation detail compiled into the library.

## Design constraints (from `docs/design.md`)

When onboarding any new function, keep these rules non-negotiable:

- **Public API is C**: exported symbols use the `fat_` prefix and live in public headers under `include/`.
- **No umbrella header**: add/modify the specific module header(s) you need; do not create `fatstd.h`.
- **Module structure is organizational**: modules are independent and should not create cyclic dependencies.
- **Go-backed objects use opaque handles**: C must never see Go memory; crossing the boundary is by handle or explicit copy-out.
- **Ownership is explicit**: every “object-like” handle needs a matching `*_Free`.
- **Fail-fast**: invalid handles / contract violations are fatal by default.

## Repository layout (where things go)

- `include/fat/*.h` — public headers for the C API (installed/consumed by users)
  - Each module gets its own header: `include/fat/<module>.h`
  - Use `FATSTD_API` on public declarations (from `include/fat/export.h`)
- `include/fat/handle.h` — canonical handle typedef used by Go-backed modules
  - Prefer `fat_Handle` (currently `uintptr_t`) and define module handle aliases as `typedef fat_Handle fat_<Thing>;`
- `src/*.c` — C implementation files and thin wrappers around Go exports
  - Prefer `src/fat_<module>.c` per module
- `pkg/` — Go implementation code compiled into the library
  - `pkg/fatstd_go/` is the **single** `package main` built with `-buildmode=c-archive`
  - Add real functionality in normal Go packages under `pkg/<module>/...`, then call into them from `pkg/fatstd_go`
- `scripts/python_tests/` — end-to-end tests that load the **shared** library via `ctypes`
  - Current entrypoint: `scripts/python_tests/test_fatstd_shared.py`

## Naming conventions (symbols, files, packages)

- **Public C symbols**: `fat_<Module><Op>` (example: `fat_VersionString`, `fat_StringLen`)
- **Go-exported C symbols (internal)**: `fatstd_go_<module>_<op>` (example: `fatstd_go_add`)
  - These come from `pkg/fatstd_go` via cgo `//export` and are called only from `src/`
- **C wrapper files**: `src/fat_<module>.c`
- **Public headers**: `include/fat/<module>.h`
- **Go packages**:
  - `pkg/fatstd_go` stays `package main` (required by `-buildmode=c-archive`)
  - Prefer lower-case package names under `pkg/<module>` (import path: `github.com/bluesentinelsec/FatStd/pkg/<module>`)

## Building (local developer loop)

FatStd’s canonical build system is CMake; this repo also provides a convenience `Makefile`.

- Static build: `make build` (default) or `make static`
- Shared build (needed for the Python `ctypes` tests): `make shared`
- Run end-to-end tests: `make test`

Equivalent CMake commands:

- Configure + build shared: `cmake -S . -B build -DFATSTD_BUILD_SHARED=ON && cmake --build build`
- Run tests: `python3 scripts/python_tests/test_fatstd_shared.py --build-dir build`

## Onboarding a **C-only** function

Use this path when the implementation is pure C and does not need the Go runtime.

1. **Pick the module and API shape**
   - Choose a module header (e.g. `include/fat/version.h`, `include/fat/fs.h`, `include/fat/json.h`).
   - Keep names C-friendly and consistent: `fat_<Module><VerbNoun>` (example: `fat_StringLen`).

2. **Add the public declaration**
   - Create or edit `include/fat/<module>.h`.
   - Include `fat/export.h` and mark exports with `FATSTD_API`.

3. **Implement in `src/`**
   - Add/modify `src/fat_<module>.c` (and an adjacent private header if needed).

4. **Register the source file with CMake**
   - Add the new `.c` file to `target_sources(fatstd PRIVATE ...)` in `CMakeLists.txt`.

5. **Add a Python `ctypes` test**
   - Edit `scripts/python_tests/test_fatstd_shared.py`:
     - Add a `test_<something>(lib)` function.
     - Set `argtypes`/`restype` exactly.
     - Call the test from `main()` so it runs under `make test`.

6. **Build + run**
   - `make test`

## Onboarding a **Go-backed** function (C API + Go implementation)

Use this path when you want to implement functionality in Go but expose it as a C API.

### 1) Design the boundary

Before writing code, decide:

- **Public C signature**: prefer C types and handles; avoid “clever” pointer sharing.
- **Ownership/lifetime**: if you return an object-like result, it must be an opaque handle and must have `*_Free`.
- **Copy-out policy**: if C needs bytes/strings, use an explicit “copy out” API that returns C-managed memory (to be freed by FatStd’s allocator once it exists) rather than leaking Go memory.

### 2) Implement the Go functionality in `pkg/`

Recommended structure:

- Put real logic in a normal Go package: `pkg/<module>/...` (example: `pkg/stringx`, `pkg/httpx`).
- Keep `pkg/fatstd_go` as the thin export layer that:
  - imports your internal Go package(s)
  - does boundary conversion
  - exposes `//export` entrypoints for C

### 3) Add a cgo-exported function in `pkg/fatstd_go`

In `pkg/fatstd_go/*.go` (must remain `package main`):

- Add `import "C"`.
- Add a `//export <symbol>` comment directly above the exported function.
- Exported symbols should be **internal**-looking and namespaced, e.g. `fatstd_go_<module>_<op>`.
- Keep exported signatures simple (C scalar types, pointers to C memory, and/or integer-ish handles).
  - Prefer C-native types like `uintptr_t`/`size_t` over cgo’s `Go*` typedefs for “boundary” signatures.

### 4) Wrap it in C (public API stays `fat_*`)

In `src/fat_<module>.c`:

- `#include "fatstd_go.h"` (generated by the Go `c-archive` build).
- Implement a C function with the public `fat_...` name that calls the `fatstd_go_...` symbol.

In `include/fat/<module>.h`:

- Declare the public `fat_...` function with `FATSTD_API`.

This keeps the Go-exported namespace private to the library, while users only see `fat_*`.

## Pattern: ergonomic + explicit-length constructors (strings/bytes)

For “byte input” APIs, prefer offering both:

- A convenience constructor for NUL-terminated input (e.g. `fat_StringNewUTF8(const char *cstr)`)
- A length-based constructor for explicit spans and embedded NULs (e.g. `fat_StringNewUTF8N(const char *bytes, size_t len)`)

This keeps the common case ergonomic while preserving full generality.

### 5) Make sure CMake rebuilds Go when Go code changes

The Go archive/header is built by a CMake `add_custom_command(...)` in `CMakeLists.txt`.

When you add new Go files/packages under `pkg/`, ensure the custom command’s `DEPENDS` includes the files that should trigger a rebuild (at minimum: the new `.go` files you add, plus `go.mod`/`go.sum` if applicable). Otherwise CMake may not re-run `go build` after edits.

### 6) Add an end-to-end Python test

Update `scripts/python_tests/test_fatstd_shared.py` to:

- call the public `fat_*` function(s) you added
- validate basic behavior + at least one edge case (invalid args, empty input, etc.)

### 7) Build + run end-to-end

- `make test`

## Adding tests: conventions

In `scripts/python_tests/test_fatstd_shared.py`:

- Name test helpers as `test_<api_or_module>(lib, ...)`.
- Always set `lib.<symbol>.argtypes` and `lib.<symbol>.restype` (ctypes defaults are unsafe).
- Prefer small, deterministic tests (no network; avoid filesystem unless that module requires it).

## End-to-end checklist (PR-quality)

- Public header added/updated in `include/fat/<module>.h`
- Implementation added/updated in `src/fat_<module>.c` (or module-appropriate file)
- `CMakeLists.txt` updated so the file is compiled and/or Go rebuild triggers are correct
- Python smoke test added/updated in `scripts/python_tests/test_fatstd_shared.py`
- `make test` passes locally

# FatStd — Design Document

## 1. Vision & Scope

FatStd is a **batteries-included, C-compatible runtime library** designed to support **modern computing workloads** while remaining **static-link friendly by design**.

The core problem FatStd addresses is that, while many high-quality C libraries exist (e.g., libcurl, libxml2, SDL), they are often:

* Difficult or impractical to static link due to deep transitive dependencies
* Bare-bones by design, providing low-level primitives instead of productivity-oriented APIs
* Fragmented across ecosystems, forcing developers to assemble a “pseudo-stdlib” themselves

Static linking is a first-class goal because it enables **single-binary deployment**, which is helpful for:

* Endpoint / client tooling
* Security engineering
* Game development
* Rapid application development
* Systems utilities


FatStd provides a **productivity layer above libc and POSIX**, offering functionality expected of modern environments—without complicated build systems or dependency hell.

FatStd is intended for:

* Systems programming
* Security engineering
* Rapid application development
* Video games

The design philosophy is **pragmatic and opinionated**:

* Focus on the 95% use cases
* Favor ease of use and functionality breadth over extreme micro-optimizations
* Strong input validation and defensive programming
* Readability and maintainability take precedence over cleverness

The project is considered **experimental until v1.0.0**. API and ABI stability are not guaranteed prior to that milestone. Semantic versioning applies post–1.0.0. For now the version is "0.1.0".

---

## 2. Target Platforms & Toolchains

FatStd targets **modern desktop and server platforms only**:

* Operating Systems:

  * Windows
  * macOS
  * Linux
* Architectures:

  * amd64
  * arm64

FatStd explicitly does **not** target embedded or freestanding environments.

Toolchain requirements:

* C standard: **C17**
* Must compile cleanly as both **C and C++**
* Supported compilers:

  * GCC
  * Clang
  * MSVC
  * MinGW

Static builds are a **key design goal**, especially to support single-binary deployment.

Platform-specific APIs should be avoided unless strictly necessary. Platform differences should be abstracted in a way similar to the Go standard library (clear separation without leaking implementation details).

**Open question:**
Should musl be *unsupported* outright, or merely not prioritized? If unsupported, should this be enforced or just documented?

unsupported - if it works, great, but I'm not investing in musl support.

---

## 3. Library Structure & Module Organization

FatStd follows a **Go- and Python-inspired module structure**, favoring clarity and discoverability over minimalism.

Key principles:

* Modules should be **largely independent**
* Shared utilities and policies may exist for cross-cutting concerns
* Cyclic dependencies are not permitted
* There is **no umbrella header** (e.g., no `fatstd.h`)
* Partial builds are not supported; FatStd is **all-or-nothing**

Directory layout:

* `include/` — public headers
* `src/` — implementation and private headers (co-located with `.c` files)
* `pkg/` - for Go-based code

Experimental labeling is unnecessary; the entire project is experimental pre–1.0.0.

**Open question:**
Should module boundaries be enforced at the build-system level (e.g., per-module static libs), or purely organizational?

Purely organizational

---

## 4. Naming Conventions & Code Style

FatStd adopts **Go-inspired naming conventions adapted to C**.

Public symbols:

* Prefixed by module name
* PascalCase for exported functions and types
  Example: `fmt_Println(...)`

General rules:

* Favor adjective–noun patterns
* Avoid macros unless strongly justified
* Opaque types are represented as opaque pointers
* Private/internal symbols follow Go-style conventions (package-private by naming)

Formatting and style:

* `clang-format` with **Microsoft C/C++ style**
* Headers must be strictly C-compliant
* Avoid compiler extensions

Inline functions are allowed where performance-critical and well-justified.

---

## 5. Memory Management Model

Memory ownership rules are **explicit and conservative**.

* By default, **the caller owns all dynamically allocated memory**
* Every allocation must have a corresponding free function
* Memory is **zero-initialized by default** to reduce undefined behavior

FatStd supports **custom allocators**, following SDL’s approach:

* The library does not provide allocator implementations
* The library must accept user-provided allocation hooks

Out-of-memory conditions are **fatal by default**:

* No attempt is made to recover
* Return a helpful error message

**Unresolved questions (must be answered):**

1. What exact `realloc` semantics should FatStd follow?
Follow SDL3's approach

2. Are stack-based helpers allowed, and if so, under what constraints?
Disallowed

3. How are allocations crossing the Go ↔ C boundary tracked and freed?
By default, all data returned from Go-backed APIs is copied into C-managed memory.
Functions that return Go-owned memory are explicitly documented and provide a corresponding free function. Failure to release such memory results in a leak.
Misuse of Go-owned memory should be fatal

---

## 6. Error Handling Strategy

FatStd’s error handling model is **explicitly inspired by SDL**.

Key properties:

* Avoid `errno`
* Prefer error codes over success-return values
* Errors are not opaque objects
* Errors are allocation-free
* Helper functions exist for retrieving and propagating error state

Fatal errors:

* By default, terminate the program with a clear message
* Callers may install fatal-error callbacks to override default behavior

**Open question:**
Will error state be thread-local (as in SDL), or global but synchronized?
Thread local should be fine

---

## 7. API Design Rules

API design mirrors the **Go standard library philosophy** wherever applicable.

Defaults:

* Majority of APIs are synchronous
* Most APIs should be thread-safe
* Blocking vs non-blocking behavior should be explicit
* Ownership transfer should be unambiguous

Additional rules:

* Prefer explicit sizes over null-terminated strings
* Prefer managed abstractions over raw OS handles
* Favor structs over long parameter lists
* Correctness is more important than zero-cost abstraction

Breaking changes are defined as violations of documented API contracts.

---

## 8. Strings, Buffers, and Data Structures

No plans on this for now, delay until we need to decide

---

## 9. Concurrency & Threading

Concurrency support is provided **via SDL3**.

FatStd relies on SDL3 for:

* Threads
* Mutexes
* Atomics
* TLS
* Async I/O (where applicable)

Global state should be avoided.

FatStd generally does **not** integrate with OS event loops, except where SDL3 requires it.

**Open question:**
Is SDL3 a hard dependency for *all* builds, or only for modules requiring concurrency?

Hard dependency for all builds.
I welcome using SDL functions for memory allocation as well as its libc replacement.

---

## 10. File System, I/O, and OS Interaction

FatStd favors:

* SDL3 APIs
* Go standard library–backed implementations exposed via C bindings

OS-specific behavior should be abstracted where possible.

(No further requirements specified.)

---

## 11. Logging, Diagnostics, and Debugging

FatStd includes a logging subsystem inspired by **Go’s `log/slog`**.

Logging properties:

* Structured logging
* Always enabled (not compile-time optional)
* Log levels:

  * Debug
  * Info
  * Warning
  * Error
  * Fatal

Assertions:

* Enabled in debug builds
* Removed in release builds
* Assertions are fatal

Tracing is explicitly out of scope.

**Open question:**
How are logs emitted (stdout, stderr, file, callback, handler interface)?

---

## 12. Build System & Distribution

Build system:

* **CMake** is canonical

Distribution properties:

* Usable as a subproject
* Installs headers system-wide
* Provides pkg-config files
* Builds static and dynamic libraries
* Targets Windows, macOS, Linux (amd64 + arm64)

There are:

* No optional dependencies
* No examples (unit tests act as examples)
* No packaging concerns yet

Library version is defined in **exactly one place**.

---

## 13. Testing & Validation

FatStd follows **test-driven development**.

Testing strategy:

* Primarily unit tests early
* Integration tests may come later
* Tests should avoid filesystem/network dependencies unless justified
* Cross-platform tests must pass everywhere
* Platform-specific tests are allowed

Undefined behavior is mitigated through:

* Clear API contracts
* Defensive programming
* Reasonable expectations of developer competence

Performance regressions are not explicitly tested, but mechanisms for measuring runtime and memory usage are desired.

---

## 14. Documentation

Documentation is generated using **Doxygen**.

Rules:

* Not every API requires examples
* Unit tests serve as implicit examples
* Ownership and lifetime must be documented explicitly
* Documentation style should resemble *C Interfaces and Implementations*

No cookbook is planned.

Primary audience: **developers**

---

## 16. Security & Hardening

Security posture:

* Favor safety and correctness
* Bounds checking enabled by default
* Insecure APIs excluded unless absolutely necessary

FatStd:

* Does not target hardened builds
* Does not integrate with sanitizers
* Practices defensive programming but does not aim for zero-UB absolutism

Threat model is intentionally minimal.

---

## 17. Governance & Contribution Model

FatStd is a **single-author project**.

* API decisions are owned by the author
* Breaking changes require author approval
* No RFC process
* Deprecations are:

  * Commented
  * Warned
  * Eventually removed

A feature is considered complete when:

* Code is implemented
* Unit tests pass
* CI passes
* Documentation is staged

---

## 18. Go Integration (Foundational Constraint)

A significant portion of FatStd is powered by the **Go standard library**, compiled into C-compatible shared libraries (static and dynamic).

Rationale:

* Go provides a batteries-included standard library
* Cross-platform by default
* Easier static builds than many C ecosystems
* Performance is acceptable for FatStd’s goals

Some subsystems (e.g., SDL-related functionality) remain purely C-based.

**Critical open questions (must be resolved):**

1. What is the FFI boundary model between Go and C?
C-first API, Go backend
Public API is C
Go is an internal engine
Strict boundary

2. How is memory ownership tracked across the boundary?
A. Copy-by-default
Go → C always copies
Go frees internally
B. Explicit ownership APIs
Documented ownership transfer
Paired free functions

3. How are Go panics handled when called from C?
Configurable
Default fatal
Optional callback override

4. Is Go runtime initialization centralized or per-module?
No explicit init API; I think the Go runtime initializes itself somehow.

5. Are Go symbols hidden to avoid ABI pollution?
Hide everything
Only C API exported
Go symbols hidden

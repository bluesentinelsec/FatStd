# FatStd — Design Document

## 1. Vision & Scope

FatStd is a **batteries-included, C-compatible runtime library** designed to support **modern computing workloads** while remaining **static-link friendly by design**.

The core problem FatStd addresses is that, while many high-quality C libraries exist (e.g., libcurl, libxml2, SDL), they are often:

* Difficult or impractical to static link due to deep transitive dependencies
* Bare-bones by design, providing low-level primitives instead of productivity-oriented APIs
* Fragmented across ecosystems, forcing developers to assemble a “pseudo-stdlib” themselves

Static linking is a first-class goal because it enables **single-binary deployment**, which is particularly useful for:

* Endpoint and client tooling
* Security engineering
* Game development
* Rapid application development
* Systems utilities

FatStd provides a **productivity layer above libc and POSIX**, offering functionality expected of modern environments without complicated build systems or dependency sprawl.

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

The project is considered **experimental until v1.0.0**. API and ABI stability are not guaranteed prior to that milestone. Semantic versioning applies post–1.0.0.

Current version: **0.1.0**

---

## 2. Target Platforms & Toolchains

FatStd targets **modern desktop and server platforms only**.

### Supported Platforms

**Operating Systems**

* Windows
* macOS
* Linux

**Architectures**

* amd64
* arm64

Embedded, freestanding, and bare-metal environments are explicitly out of scope.

### Toolchain Requirements

* C standard: **C17**
* Must compile cleanly as both **C and C++**
* Supported compilers:

  * GCC
  * Clang
  * MSVC
  * MinGW

Static builds are a **key design goal**, particularly to support single-binary deployment.

Platform-specific APIs are avoided unless strictly necessary. Platform differences are abstracted in a manner similar to the Go standard library, with clear separation and minimal leakage.

---

## 3. Library Structure & Module Organization

FatStd follows a **Go- and Python-inspired module structure**, prioritizing clarity and discoverability.

### Principles

* Modules are **largely independent**
* Shared utilities and policies may exist for cross-cutting concerns
* Cyclic dependencies are **not permitted**
* There is **no umbrella header** (e.g., no `fatstd.h`)
* Partial builds are **not supported**; FatStd is all-or-nothing

### Directory Layout

* `include/` — public headers
* `src/` — implementation and private headers (co-located with `.c` files)
* `pkg/` — Go-based implementation code

Module boundaries are **organizational**, not enforced at the build-system level.

The entire project is experimental prior to v1.0.0; no additional experimental labeling is used.

---

## 4. Naming Conventions & Code Style

### Opaque Types

All externally visible object-like types are **opaque handles**.

```c
typedef struct fat_Handle fat_Handle;
```

Module-specific types are aliases of `fat_Handle`:

```c
typedef fat_Handle fat_String;
typedef fat_Handle fat_HttpClient;
typedef fat_Handle fat_HttpResponse;
```

These aliases exist **solely for API clarity and type safety**. They do not imply distinct layouts or ownership models.

### General Rules

* Favor adjective–noun naming
* Avoid macros unless strongly justified
* Opaque types are represented as opaque pointers
* Private and internal symbols follow Go-style package-private conventions

### Formatting & Style

* `clang-format` using **Microsoft C/C++ style**
* `gofmt` for Go code
* Headers must be strictly C-compliant
* Compiler extensions are avoided

Inline functions are permitted where performance-critical and well-justified.

---

## 5. Memory Management Model

Memory ownership in FatStd is **explicit, uniform, and conservative**.

### 5.1 Core Rule: Handle-Based Ownership

**All Go-backed objects exposed to C are represented as opaque handles.**
**C code never observes Go memory directly.**
**All interaction occurs through Go-backed functions operating on handles.**

This rule is **foundational** and applies to all Go-backed subsystems.

It guarantees:

* ABI stability
* GC safety
* Static-link friendliness
* Predictable lifetime management
* Uniform API design

Violations of this rule are considered **design errors**.

---

### 5.2 Handle Semantics

A handle is:

* An opaque token (`fat_Handle *`, `uintptr_t`)
* Meaningless to C except as an identity
* Mapped internally to a Go object
* Invalid after being freed

A handle is not:

* A pointer to Go memory
* A C struct
* A ref-counted object
* A view into internal state

Handles are **identity-only capabilities**.

---

### 5.3 Lifetime Rules

The lifetime of a Go-backed object is:

1. Go creates the object
2. Go registers it and returns a handle
3. C uses the handle
4. C calls the matching `*_Free` function
5. Go removes the object from its registry
6. Go GC reclaims the object naturally

Misuse (double free, use-after-free, invalid handle):

* Is **fatal by default**
* Diagnosed eagerly, especially in debug builds
* Considered programmer error

FatStd does not attempt to recover from misuse.

---

### 5.4 Allocation Domains

FatStd operates with **two strictly separated allocation domains**:

| Domain            | Owned By | Freed By         |
| ----------------- | -------- | ---------------- |
| Go-backed objects | Go       | `*_Free(handle)` |
| C-visible buffers | C        | `fat_Free()`     |

Go memory never escapes into C.

Data crossing the boundary does so **only by explicit copy**.

---

### 5.5 Copy-Out APIs

Some APIs intentionally copy data into C-managed memory:

```c
char *fat_StringCopyUTF8(fat_String *s, size_t *len);
```

Properties:

* Allocated via the FatStd allocator
* Ownership transfers to the caller
* Must be freed explicitly
* Clearly documented
* Optional, not default

These APIs are **escape hatches**, not the primary interaction model.

---

### 5.6 Allocators

FatStd supports **custom allocators**, following SDL’s model:

* FatStd does not provide allocator implementations
* Users may install allocation hooks
* All copy-out APIs respect these hooks

Stack-based allocation helpers are explicitly disallowed.

---

### 5.7 Out-of-Memory Behavior

Out-of-memory conditions are **fatal by default**:

* No partial recovery
* Program terminates with a clear diagnostic

This behavior is consistent across both C and Go-backed code.

---

## 6. Error Handling Strategy

FatStd employs a **fail-fast, explicit error model**.

### Fatal Errors

By default, fatal errors:

* Terminate the program
* Emit a clear diagnostic
* Do not return error codes

Fatal errors include:

* Invalid handle use
* Internal Go panics
* Violations of API contracts

A configurable fatal-error callback may be installed to override default termination behavior.

---

## 7. API Design Rules

FatStd APIs are **handle-first** by design.

### Design Constraints

* No APIs expose raw Go memory
* No APIs return pointers into internal storage
* All object-like results are handles
* Ownership transfer is explicit
* Free functions are mandatory

### Canonical Shape

```c
fat_String *fat_HTTPGetBody(const char *url);
size_t      fat_StringLen(fat_String *s);
char       *fat_StringCopyUTF8(fat_String *s, size_t *len);
void        fat_StringFree(fat_String *s);
```

This pattern is canonical across the library.

---

## 8. Strings, Buffers, and Data Structures

Strings, buffers, streams, and structured data share a **single abstraction model**:

* Opaque handle
* Explicit lifetime
* Operation-based access
* Optional copy-out

There is no special-casing for strings.

---

## 9. Concurrency & Threading

Concurrency support is provided **via SDL3**, which is a **hard dependency for all builds**.

FatStd relies on SDL3 for:

* Threads
* Mutexes
* Atomics
* Thread-local storage
* Asynchronous I/O where applicable

Global state is avoided.

FatStd does not integrate with OS event loops except where required by SDL3.

Because all Go-backed state is accessed through handles:

* Thread safety can be enforced centrally
* Debug checks are centralized
* Diagnostics scale cleanly

Most APIs are thread-safe unless explicitly documented otherwise.

---

## 10. File System, I/O, and OS Interaction

FatStd favors:

* SDL3 APIs
* Go standard library–backed implementations exposed via C bindings

OS-specific behavior is abstracted wherever possible.

---

## 11. Logging, Diagnostics, and Debugging

FatStd includes a logging subsystem inspired by **Go’s `log/slog`**.

### Logging Properties

* Structured logging
* Always enabled
* Log levels:

  * Debug
  * Info
  * Warning
  * Error
  * Fatal

Log emission is centralized and configurable at runtime.

### Assertions

* Enabled in debug builds
* Removed in release builds
* Assertion failures are fatal

Tracing is out of scope.

---

## 12. Build System & Distribution

### Build System

* **CMake** is the canonical build system

### Distribution Characteristics

* Usable as a subproject
* Installs headers system-wide
* Provides pkg-config metadata
* Builds static and dynamic libraries
* Targets Windows, macOS, and Linux (amd64 and arm64)

Additional constraints:

* No optional dependencies
* No examples (unit tests serve as examples)
* Library version defined in exactly one place

---

## 13. Testing & Validation

FatStd follows **test-driven development**.

### Testing Strategy

* Unit tests first
* Integration tests added selectively
* Tests avoid filesystem and network dependencies unless justified
* Cross-platform tests must pass everywhere
* Platform-specific tests are permitted

Undefined behavior is mitigated through:

* Clear API contracts
* Defensive programming
* Explicit ownership rules

---

## 14. Documentation

Documentation is generated using **Doxygen**.

### Documentation Rules

* Ownership and lifetime are explicitly documented
* Unit tests serve as implicit examples
* Documentation style resembles *C Interfaces and Implementations*

Primary audience: **developers**.

---

## 15. Security & Hardening

### Security Posture

* Favor safety and correctness
* Bounds checking enabled by default
* Insecure APIs excluded unless absolutely necessary

FatStd:

* Does not target hardened builds
* Does not integrate with sanitizers
* Practices defensive programming without pursuing zero-UB guarantees

The threat model is intentionally minimal.

---

## 16. Governance & Contribution Model

FatStd is a **single-author project**.

* API decisions are owned by the author
* Breaking changes require author approval
* No formal RFC process

### Deprecation Policy

Deprecated APIs are:

* Commented
* Warned
* Eventually removed

A feature is considered complete when:

* Code is implemented
* Unit tests pass
* CI passes
* Documentation is staged

---

## 17. Go Integration (Foundational Constraint)

A significant portion of FatStd is powered by the **Go standard library**, compiled into C-compatible libraries.

### Integration Model

* C-first API
* Go as internal backend
* Handles as the sole boundary
* No Go memory visible to C

### Handle Registry

Go-backed modules use a **centralized handle registry**:

* Handles are opaque integers
* Internals are private to Go
* Registry lifetime matches the Go runtime
* Registry is shared across modules

### Go ↔ C Rules

* All Go-backed objects are accessed via handles
* Go memory never crosses into C
* Data transfer requires explicit copy APIs
* Panics are fatal by default
* ABI exposure is limited strictly to the C API

This integration model is **non-negotiable**.

---

**FatStd defines a batteries-included C runtime with modern capabilities, powered by Go, grounded in explicit ownership, and engineered for static linking and long-term maintainability.**

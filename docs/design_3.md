# FatStd — Design Document

## 1. Vision & Scope

FatStd is a **batteries-included, C-compatible runtime library** designed to support **modern computing workloads** while remaining **static-link friendly by design**.

The core problem FatStd addresses is that, while many high-quality C libraries exist (e.g., libcurl, libxml2, SDL), they are often:

* Difficult or impractical to static link due to deep transitive dependencies
* Bare-bones by design, providing low-level primitives instead of productivity-oriented APIs
* Fragmented across ecosystems, forcing developers to assemble a “pseudo-stdlib” themselves

Static linking is a first-class goal because it enables **single-binary deployment**, which is particularly useful for:

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

The project is considered **experimental until v1.0.0**. API and ABI stability are not guaranteed prior to that milestone. Semantic versioning applies post–1.0.0.
Current version: **0.1.0**

---

## 2. Target Platforms & Toolchains

FatStd targets **modern desktop and server platforms only**.

### Supported Platforms

* **Operating Systems**

  * Windows
  * macOS
  * Linux
* **Architectures**

  * amd64
  * arm64

Embedded and freestanding environments are explicitly out of scope.

### Toolchain Requirements

* C standard: **C17**
* Must compile cleanly as both **C and C++**
* Supported compilers:

  * GCC
  * Clang
  * MSVC
  * MinGW

Static builds are a **key design goal**, especially to support single-binary deployment.

Platform-specific APIs should be avoided unless strictly necessary. Platform differences should be abstracted in a manner similar to the Go standard library, with clear separation and minimal leakage.

**musl support:**
musl is not supported or prioritized. If FatStd happens to work with musl, that is acceptable, but no development effort will be invested in ensuring compatibility.

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

Module boundaries are **purely organizational** and are not enforced at the build-system level.

The entire project is experimental prior to v1.0.0; no additional experimental labeling is used.

---

## 4. Naming Conventions & Code Style

FatStd adopts **Go-inspired naming conventions adapted to C**.

### Public Symbols

* Prefixed by module name
* PascalCase for exported functions and types
  Example: `fmt_Println(...)`

### General Rules

* Favor adjective–noun naming
* Avoid macros unless strongly justified
* Opaque types are represented as opaque pointers
* Private/internal symbols follow Go-style “package-private” conventions

### Formatting & Style

* `clang-format` using **Microsoft C/C++ style**
* Headers must be strictly C-compliant
* Compiler extensions are avoided

Inline functions are permitted where performance-critical and well-justified.

---

## 5. Memory Management Model

Memory ownership rules are **explicit and conservative**.

### General Rules

* By default, **the caller owns all dynamically allocated memory**
* Every allocation has a corresponding free function
* Allocated memory is **zero-initialized by default** to reduce undefined behavior

### Allocators

FatStd supports **custom allocators**, following SDL’s approach:

* FatStd does not provide allocator implementations
* The library accepts user-provided allocation hooks

### Out-of-Memory Behavior

Out-of-memory conditions are **fatal by default**:

* No attempt is made to recover
* The program terminates with a helpful error message

### Specific Decisions

* `realloc` semantics follow **SDL3’s behavior**
* Stack-based allocation helpers are **disallowed**
* Go ↔ C memory crossing rules:

  * By default, Go-backed APIs **copy data into C-managed memory**
  * APIs that return Go-owned memory are **explicitly documented**
  * Such APIs provide a matching free function
  * Misuse (e.g., failure to free Go-owned memory) is **fatal**

---

## 6. Error Handling Strategy

FatStd’s error handling model is **explicitly inspired by SDL**.

### Properties

* Avoid `errno`
* Prefer error codes over success-return values
* Errors are not opaque objects
* Errors are allocation-free
* Helper functions exist for retrieving and propagating error state

### Fatal Errors

* By default, fatal errors terminate the program with a clear message
* Callers may install fatal-error callbacks to override default behavior

Error state is **thread-local**, following SDL’s model.

---

## 7. API Design Rules

API design mirrors the **Go standard library philosophy** wherever applicable.

### Defaults

* Majority of APIs are synchronous
* Most APIs are thread-safe
* Blocking vs non-blocking behavior is explicit
* Ownership transfer is unambiguous

### Additional Rules

* Prefer explicit sizes over null-terminated strings
* Prefer managed abstractions over raw OS handles
* Favor structs over long parameter lists
* Correctness is more important than zero-cost abstraction

Breaking changes are defined as violations of documented API contracts.

---

## 8. Strings, Buffers, and Data Structures

Design of strings, buffers, and container data structures is **intentionally deferred** and will be addressed when concrete needs arise.

---

## 9. Concurrency & Threading

Concurrency support is provided **via SDL3**, which is a **hard dependency for all builds**.

FatStd relies on SDL3 for:

* Threads
* Mutexes
* Atomics
* TLS
* Async I/O (where applicable)

Global state should be avoided.

FatStd generally does **not** integrate with OS event loops, except where required by SDL3.

SDL3 APIs may also be used for memory allocation and libc replacement where appropriate.

---

## 10. File System, I/O, and OS Interaction

FatStd favors:

* SDL3 APIs
* Go standard library–backed implementations exposed via C bindings

OS-specific behavior should be abstracted wherever possible.

---

## 11. Logging, Diagnostics, and Debugging

FatStd includes a logging subsystem inspired by **Go’s `log/slog`**.

### Logging Properties

* Structured logging
* Always enabled (not compile-time optional)
* Log levels:

  * Debug
  * Info
  * Warning
  * Error
  * Fatal

### Assertions

* Enabled in debug builds
* Removed in release builds
* Assertions are fatal

Tracing is explicitly out of scope.

*(Log emission mechanism—stdout/stderr/file/handler—remains to be finalized.)*

---

## 12. Build System & Distribution

### Build System

* **CMake** is canonical

### Distribution Characteristics

* Usable as a subproject
* Installs headers system-wide
* Provides pkg-config files
* Builds static and dynamic libraries
* Targets Windows, macOS, Linux (amd64 + arm64)

Additional constraints:

* No optional dependencies
* No examples (unit tests serve as examples)
* No packaging concerns at this stage
* Library version is defined in **exactly one place**

---

## 13. Testing & Validation

FatStd follows **test-driven development**.

### Testing Strategy

* Primarily unit tests early
* Integration tests may be added later
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

### Documentation Rules

* Not every API requires examples
* Unit tests serve as implicit examples
* Ownership and lifetime must be documented explicitly
* Documentation style should resemble *C Interfaces and Implementations*

No cookbook is planned.

Primary audience: **developers**

---

## 16. Security & Hardening

### Security Posture

* Favor safety and correctness
* Bounds checking enabled by default
* Insecure APIs excluded unless absolutely necessary

FatStd:

* Does not target hardened builds
* Does not integrate with sanitizers
* Practices defensive programming without aiming for absolute zero-UB guarantees

The threat model is intentionally minimal.

---

## 17. Governance & Contribution Model

FatStd is a **single-author project**.

* API decisions are owned by the author
* Breaking changes require author approval
* No RFC process

### Deprecation Policy

* Deprecated APIs are:

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

### Rationale

* Go provides a batteries-included standard library
* Cross-platform by default
* Easier static builds than many C ecosystems
* Performance is acceptable for FatStd’s goals

Some subsystems (e.g., SDL-related functionality) remain purely C-based.

### Go ↔ C Integration Rules

* **FFI model:** C-first API with Go as an internal backend
* **Memory ownership:**

  * Copy-by-default
  * Explicit ownership APIs where necessary, with paired free functions
* **Panic handling:**

  * Configurable
  * Default behavior is fatal
  * Optional callback override
* **Go runtime initialization:**

  * No explicit init API
  * Runtime initializes implicitly as needed
* **ABI hygiene:**

  * Go symbols are hidden
  * Only the C API is exported

---

If you want next steps, good candidates would be:

* Locking the logging emission model
* Defining the first concrete module set (v0.1 scope)
* Writing a machine-readable spec for AI-driven implementation

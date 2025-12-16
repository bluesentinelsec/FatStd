---

# FatStd — Design Doc (Question-Only Specification)

## 1. Vision & Scope

1. What problem does FatStd solve that libc, POSIX, and common C utility libraries do not?
There are many competent C libraries such as LibCurl, LibXML, LibSDL, and so on.
The problem is, many C libraries are not easy to static link, perhaps because they contain many transitive dependencies which also need to be static linked.
Static linking is important because it makes it easy to ship a single binary to endpoints.
Additionally, most C libraries tend to be bare bones, e.g., they provide a minimum of functionality at a low level of abstraction, usually by design. This is great for lean, high-performance computing, but is a nuissance for someone who just wants to write C code addressing modern requirements (web friendly, lots of encoding support, strings that don't suck, etc.).

In summary, FatStd provides a batteries included C-compatible shared library that is static link friendly by design.
In that way, developers can rapidly tap into functions expected for modern workloads, without complicated build or deployment requirements.


2. What does “batteries included” mean concretely for C in this project?
Fatstd provides functionality that is expected for modern computing tasks, but tends to be excluded by design from the C standard library.
Examples of this functionality includes:
- Multimedia (like SDL or Raylib)
- Data processing: JSON, XML
- Compression / decompression
- Web / HTTP/S
- Strings that don't suck
- Data structures and algorithms
- Tiled map parser
- Cross platform friendly



3. Is FatStd intended primarily for:

FatStd is intended for systems programming, security engineering, rapid application development, and video games in mind.

4. What design philosophy should FatStd follow (minimalism, pragmatism, orthogonality, opinionated defaults, etc.)?
Pragmatism: simple, easy-to-use API's. Focus on the 95% of features that are needed by almost all users.
FatStd should be performant, but, functionality breadth trumps performance.

5. What is explicitly **out of scope** for FatStd?
Nothing in particular, but, focus on functionality that is needed for modern workloads.

6. Should FatStd aim to be:
Productivity layer above libc & POSIX.

7. Is API stability more important than rapid feature evolution?
API stability does not matter until we publish v1.0.0.
API can be unstable during active development.

8. Should FatStd prioritize:
FatStd priorities, in order:
1. Readability & maintainability (following test driven development)
2. Ease of use
3. Functionality breadth
4. Strong input validation

9. What guarantees (if any) does FatStd make about ABI stability?
Nothing during development.
Follow semantic versioning afterwards.

10. How should FatStd communicate its maturity level (experimental, stable, LTS, etc.)?
Experimental until v1.0.0.

---

## 2. Target Platforms & Toolchains

1. Which operating systems must FatStd support (Linux, Windows, macOS, BSD, etc.)?
Windows, macOS, and Linux, AMD64 and ARM64 respectively

2. Which C standards must be supported (C89, C99, C11, C17, C23)?
C17

3. Must the code compile cleanly as C++ (e.g., with g++)?
Yes

4. Which compilers must be supported (GCC, Clang, MSVC, MinGW)?
GCC, Clang, MinGW, MSVC

5. Are freestanding or embedded environments in scope?
No

6. Should FatStd support musl, glibc, and other libc implementations equally?
No

7. Are static builds a first-class goal?
Yes, static builds is a key goal to make shipping executables easier (no dependency hell)

8. Should FatStd avoid OS-specific APIs unless strictly necessary?
Yes, most code should be cross platform.
Exceptions can be made for essential platform-specific behavior.

9. What is the minimum supported platform baseline?
Windows, macOS, and Linux (current LTS releases; we don't care about legacy platforms).

10. How should platform differences be abstracted or exposed?
I like how Go separates platform-specific functionality.

---

## 3. Library Structure & Module Organization

1. How should FatStd be partitioned into modules?
Follow simillar module structure as Python and Go.

2. Should modules be independently usable or tightly integrated?
Prefer independence, with exception for shared policies and low-level utils that are ubiquotous to most modules.

3. What naming convention should modules follow?
Mirror Go conventions

4. Should there be a single umbrella header (e.g., `fatstd.h`)?
No

5. How should internal vs public headers be separated?
Public headers should go in "include/" - private headers should go under "src/" and be co-located with .c files.

6. Should optional features be compile-time toggles?
No, keep the build simple. All or nothing. If you want bespoke configurations just piecemeal your own std lib.

7. How should experimental modules be labeled?
No need for experimental labels - the project itself is experimental at this moment.

8. Should modules depend on each other or remain acyclic?
Acylic is preferred but exceptions can be made with good justification.

9. How should platform-specific code be organized?
Yes - all code should be organized.

10. Should FatStd allow partial builds (selective modules)?
No.

---

## 4. Naming Conventions & Code Style

1. What prefix should all public symbols use?
Use the module name as a prefix.
Contrived example: fmt_Println(args...)

2. Should FatStd use snake_case, camelCase, or another convention?
See my previous example: fmt_Println(args)
Favor adjective-noun patterns.

3. How should types, functions, macros, and constants be distinguished?
Follow Go conventions

4. Should opaque types be exposed as structs or typedefs?
Opaque pointers are welcome.

5. How should private/internal symbols be named?
Follow Go conventions.

6. Should macros be minimized or used aggressively?
Use rarely and with strong justification.

7. What formatting rules should be enforced (indentation, braces, line length)?
clang-format in microsoft style

8. Should the codebase follow a specific style guide?
Microsoft C/C++ style

9. How should inline functions be used?
Inlines are welcome for performance; favor known-valid uses of inlines

10. Should headers be strictly C-only (no GNU extensions)?
Yes, avoid C extensions

---

## 5. Memory Management Model

1. Does FatStd own memory by default or require caller ownership?
Caller should own dynamic memory.
Exception is for Go-based functions, which the Go runtime matters.
But we still need to manage dynamic allocations, which occur between C and Go communication.

2. Should FatStd expose custom allocators?
Yes; I like how SDL does this.

3. Should every allocation have a symmetric free function?
Yes

4. How should allocator overrides be configured?
Follow SDL's approach

5. Should FatStd support arena/region allocation?
We don't need to provide custom allocator implementations, we just need to support custom allocators.

6. Should memory zeroing be guaranteed?
Yes. I value default initializations to avoid suprises and reduce undefined behavior.

7. How should realloc semantics be handled?
I don't know. D:
Pick a sane, ubiquotous approach.

8. Are stack-based helpers acceptable?
Maybe? I don't know what this is.

9. Should allocation failures be recoverable?
No... recovering from out of memory issues is stupid, just crash and let the kernel clean it up or reboot.

10. How should memory ownership be documented in APIs?
Caller should own memory, in general.

---

## 6. Error Handling Strategy

1. How are errors represented (return codes, structs, thread-local state)?
I like SDL's approach to error handling.

2. Should FatStd avoid `errno`?
Yes. Minimize usage of the C standard library.

3. Should functions return error codes or success values?
Prefer error codes

4. How should fatal vs non-fatal errors be distinguished?
Let the caller provide call back functions for fatal errors, but by default, terminate the program with a helpful error message

5. Should error objects be opaque?
I don't think so, favor an SDL approach

6. Should errors carry:
Follow the SDL convention

7. Should FatStd provide helpers for error propagation?
Follow the SDL convention

8. Should errors be alloc-free?
Follow the SDL convention

9. How should errors be logged or surfaced?
Follow the SDL convention

---

## 7. API Design Rules

1. Should APIs be synchronous by default?
Mirror the Go standard library.
I think most functions should be thread safe.

2. How should blocking vs non-blocking behavior be expressed?
Mirror the Go std library.

3. Should APIs be re-entrant and thread-safe by default?
Mirror the Go std library.

4. How should ownership transfer be signaled?
Mirror the Go std library.

5. Should APIs favor explicit sizes over null-terminated strings?
Explicit sizes

6. Should APIs allow zero-cost abstraction?
Less important than correctness

7. How should optional parameters be handled?
Mirror the Go std library.

8. Should FatStd expose raw OS handles?
No - interfaces should be managed

9. Should APIs favor structs over long parameter lists?
Mirror the Go std library.

10. What constitutes a breaking API change?
Breaking API contracts

---

## 8. Strings, Buffers, and Data Structures

1. Should FatStd provide its own string type?
I think so - plain C strings are error / issue prone

2. Should strings be UTF-8, binary-safe, or both?
Probably UTF-8

3. How should string ownership be handled?
No idea

4. Should FatStd provide dynamic arrays, hash maps, and sets?
Yes

5. Should data structures be intrusive or non-intrusive?
I don't know; mirror Go

6. Should iterators be explicit structs or callback-based?
I can live with callbacks

7. How should resizing and growth policies work?
Mirror Go

8. Should containers expose raw pointers?
Favor managed approach

9. Should containers be allocator-aware?
The system needs to support custom allocators


---

## 9. Concurrency & Threading

1. Should FatStd provide threading abstractions?
Yes, I think we can leverage SDL3 for this

2. Should mutexes, atomics, and channels be included?
Yes, SDL3

3. Should APIs be lock-free where possible?
Yes, via SDL3

4. How should thread safety be documented?
Just use SDL3

5. Should FatStd avoid hidden global state?
Yes, globals are hard to reason about in large code bases

6. Should TLS be used?
Use SDL3

7. Should concurrency be optional?
Not sure, use SDL3

8. How should async I/O be exposed?
SDL3

9. Should FatStd integrate with OS event loops?
Generally no, with exception for SDL3


---

## 10. File System, I/O, and OS Interaction

Favor SDL3 functions and Go standard library functions
---

## 11. Logging, Diagnostics, and Debugging

1. Should FatStd include a logging subsystem?
Yes

2. What log levels should exist?
Debug, Info, Error, Warning, Fatal

3. Should logging be optional at compile time?
No

4. Should logs be structured or text-based?
Favor Go's slog approach (structured)

5. How should logging interact with errors?
Follow Go's slog approach

6. Should debug builds add additional checks?
Asserts for debug builds, no asserts for release

7. Should assertions be fatal?
Yes

8. Should tracing hooks exist?
No, if you need to trace, use a debugger

9. How should performance metrics be exposed?
I don't know

---

## 12. Build System & Distribution

1. What build system should be canonical (CMake, Meson, Make)?
Cmake

2. Should FatStd be usable as a subproject?
Yes

3. Should pkg-config files be provided?
Yes

4. How should versioning be handled?
Semantic versioning; I would like to define the library version in ONE place

5. Should symbols be versioned?
No; if we need to make breaking changes, we'll increment the major version number

6. Should FatStd install headers system-wide?
Yes

7. How should optional dependencies be handled?
No optional deps, all or nothing

8. Should examples be built by default?
There won't be examples, but there will be unit tests (which are implicit examples)

9. Should FatStd support cross-compilation?
Yes, but this isn't a strong requirement.
Focus on supporting the specified compilers.

10. How should FatStd be packaged for different platforms?
Don't worry about packaging for now.
I just need to produce idiomatic shared libraries, dynamic and static (windows, mac, linux - amd64 and arm64)

---

## 13. Testing & Validation

1. What level of test coverage is required?
Follow test driven development

2. Should tests be unit, integration, or both?
Primarily unit early on
We may add integration tests later

3. Should tests avoid filesystem/network dependencies?
Yes in most cases; I may make exceptions

4. Should tests run on all supported platforms?
Yes for cross platform tests, no for platform-specific tests

5. How should undefined behavior be tested?
Focus on defining API contracts that explicitly states intended behavior
We should be defensive programmers, but we don't need to move heaven and earth is the developer is being negligent

6. Should fuzzing be part of CI?
No

7. Should ABI compatibility be tested?
I don't think so - favor unit tests

8. Should performance regressions be tested?
Not explicitly, but I would like mechanisms for measuring performance (runtime, maximum consumed RAM)

10. What constitutes a test failure?
Failing unit tests

---

## 14. Documentation & Examples

1. What documentation format should be used?
Doxygen for C code

2. Should every public API have usage examples?
No, leverage unit tests for examples

3. Should docs describe ownership and lifetime explicitly?
Yes; I like the documentation approach in "C Interfaces and Implementations"

4. Should docs be generated from source comments?
Yes (Doxygen)

5. Should examples favor minimalism or realism?
He succinct, but thorough.

6. Should there be a cookbook?
No

8. How should deprecated APIs be documented?
Whatever Doxygen does

9. Should documentation include design rationale?
Depends... follow the approach in "C Interfaces and Implementations"

10. Who is the primary audience for the docs?
Developers

---


## 16. Security & Hardening

1. Should FatStd default to safe behavior over permissive behavior?
Yes, prefer safety and correctness over permissive behavior

2. Should bounds checking be enabled by default?
Yes

3. Should insecure APIs be excluded or gated?
Favor excluded unless necessary

4. Should FatStd provide secure memory primitives?
I'm open to suggestions

5. How should misuse be detected?
We should have thorough but performant error checking

6. Should FatStd support hardened builds?
No

7. Should FatStd avoid undefined behavior entirely?
We should practice defensive programming

8. How should input validation be handled?
Return a useful error message if non-fatal, otherwise terminate (with exception for fatal callbacks)

9. Should FatStd integrate with sanitizers?
No

10. What is FatStd’s threat model?
Not worried about it, this project won't have sufficient market share to matter
If it does matter, my threat model is annoying pentesters that just want to publish CVE's for their resume

---

## 17. Governance & Contribution Model

1. Who owns the API design decisions?
Me

2. How are breaking changes approved?
I approve them

3. How should contributions be reviewed?
I review them (but there realistically won't be any, its MY library)

4. Should there be a formal RFC process?
No

5. How should experimental features graduate?
N/A

6. How should deprecations be handled?
Mark it in comments, emit console warnings, eventually delete the deprecated functionality

7. What coding standards must contributors follow?
test driven development
Microsoft C/C++ style

8. How should issues be triaged?
N/A

9. How should AI-generated code be reviewed?
N/A

10. What defines “done” for a feature?
Feature's code
Unit tests proving correctness
CICD builds and tests pass
Documentation is staged

---

Something else that wasn't covered but is relevant:

I want to leverage the Go standard library for much of FatStd

Go was chosen because its cross platform and easily compiles to C-compatible shared libraries, both dynamic and static.

Go also has a batteries included standard library, and is relatively performant compared to say Python.

Much of FatLib will be powered by Go with C bindings, but some functionality, such as SDL, will be C-based.


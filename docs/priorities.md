## Tier 0 — Foundational (must exist first)

These **affect everything else**.

### 1. `core` / `runtime`

**Why first:**
Defines library-wide invariants.

Includes:

* Versioning
* Init / shutdown hooks (even if implicit)
* Global configuration
* Fatal error callbacks
* Go runtime boundary glue

If this comes later, you’ll rework every module.

---

### 2. `mem`

**Why first:**
Memory ownership rules cascade everywhere.

Includes:

* Allocator hooks (SDL-style)
* Zero-init guarantees
* `malloc/free/realloc` wrappers
* Go ↔ C allocation helpers
* Fatal OOM handling

Every module allocates memory. This must be locked early.

---

### 3. `error`

**Why first:**
Error semantics must be consistent.

Includes:

* Error codes
* Thread-local error state
* Error string retrieval
* Panic → error translation (Go side)
* Fatal vs non-fatal pathways

If added late, APIs fracture.

---

### 4. `log`

**Why first:**
Logging becomes an implicit dependency everywhere.

Includes:

* Structured logging core (slog-inspired)
* Log levels
* Emission backend abstraction
* Fatal logging integration

Later modules will *assume* logging exists.

---

## Tier 1 — Cross-Cutting Policy Modules

These influence API *shape*, not just functionality.

### 5. `time`

(backed by Go `time`)

**Why early:**

* Timeouts
* Deadlines
* Durations
* Timestamps
* Required for HTTP, networking, async, logging

Time types must be consistent everywhere.

---

### 6. `io`

(backed by Go `io`)

**Why early:**

* Defines streaming semantics
* Reader/Writer abstractions
* Zero-copy vs buffered decisions
* Used by HTTP, compression, filesystem, codecs

If IO is wrong, everything built on top suffers.

---

### 7. `fs`

(backed by Go `os`, `io/fs`)

**Why early:**

* Path handling
* File handles
* Permissions model
* Error mapping
* Used by config, logging, assets, cert loading

Filesystem abstractions are sticky once public.

---

## Tier 2 — “Gravity” Modules (many depend on them)

These don’t define policy, but **pull many others in**.

### 8. `net`

(backed by Go `net`)

**Why now:**

* Sockets
* Address parsing
* DNS
* Used by HTTP, TLS, services

Networking semantics must be stable early.

---

### 9. `tls`

(backed by Go `crypto/tls`)

**Why before HTTP:**

* Cert loading
* Verification policy
* Trust store handling
* Error surface

You do *not* want to retrofit TLS behavior later.

---

### 10. `http`

(backed by Go `net/http`)

**Why not first, but early:**

* Depends on time, io, net, tls, error, log
* API shape will expose earlier design decisions

HTTP is a forcing function that validates earlier choices.

---

## Tier 3 — Nice-to-Have / Delay Safely

These can wait without poisoning the design.

* `json` / `xml`
* Compression
* Hashing / crypto helpers
* Multimedia wrappers
* Data structures
* Encoding helpers

They consume infrastructure but don’t define it.

---

## Recommended First-Pass Order (Concrete)

If you want a **safe, low-churn bootstrap path**:

1. `core`
2. `mem`
3. `error`
4. `log`
5. `time`
6. `io`
7. `fs`
8. `net`
9. `tls`
10. `http`

After `http`, your architecture is effectively locked.

---

## Litmus Test (Rule of Thumb)

If a module:

* Defines ownership rules
* Introduces new base types
* Changes error propagation
* Requires policy decisions

→ **Front-load it**

If it:

* Just consumes existing abstractions

→ Can wait.

---

To lay the strongest foundation for FatStd's **Go runtime + C wrappers**, prioritize Go standard library packages that are:

- Fundamental to your handle-based, opaque-object model.
- Heavily used across nearly all other modules.
- Relatively straightforward to bind without deep platform-specific quirks.
- Critical for enabling your end-goal of Fat Lua quickly.

These early bindings will stress-test your centralized handle registry, memory rules (no Go memory escaping, explicit copies), error/panic handling, and cross-module sharing the most. Getting them right first locks in your core integration patterns.

### Recommended Priority Order

1. **strings + bytes + strconv** (Start here — absolute foundation)
   - Why first: Almost every other package deals with text or binary data. You'll need robust string/buffer handling for APIs like `fat_StringCopyUTF8`, lengths, conversions, building, trimming, etc.
   - Influences: Sets the pattern for your core opaque `fat_String` and `fat_Buffer` types. Every copy-out API will rely on this.
   - Low-risk: Pure Go, no OS/syscall issues, easy to bind safely with copies.
   - Bonus: Immediately useful for logging, diagnostics, and internal utilities.

2. **io + io/fs + os + path/filepath** (Next — filesystem is a core productivity win)
   - Why: Go's fs/io/os stack is excellent and cross-platform. Bindings here give you high-level file read/write/stat/mkdir/walk without reinventing wheels or depending on fragile POSIX.
   - Influences: Many modules (http, crypto, encoding) need file I/O. Your design favors Go-backed implementations for OS interaction.
   - Pair with bufio if needed for readers/writers.
   - Note: Use pure-Go paths where possible; avoid deep syscall if it complicates static builds.

3. **net/http** (including http.Client, Response, etc.)
   - Why: One of the biggest "batteries" in Go's stdlib. A solid HTTP client (GET/POST/TLS/JSON handling) will feel magically productive to C users.
   - Influences: Touches strings, io, crypto/tls, time. Great for testing handle lifetime across requests.
   - TLS comes "for free" via crypto/tls — huge win over linking OpenSSL manually.

4. **crypto/** subtree (crypto/sha256, crypto/hmac, crypto/rand, crypto/aes, etc.; plus encoding/base64, hex)
   - Why: Modern tooling needs hashing, signing, encryption, random bytes. Go's crypto is pure-Go and best-in-class.
   - Influences: Used by http (TLS), potential future auth, and security-focused users.
   - Start with the basics (sha256, hmac, rand) — they're self-contained.

After these four, your foundation will be rock-solid, and you'll have covered ~80% of the "that's cool" use cases (HTTP fetches, JSON via strings+bytes, file manipulation, crypto primitives). This directly enables the minimal Fat Lua prototype you sketched: fat_String/Buffer, fat_Http, fat_Fs, fat_Crypto.

### Why This Order (and Not Others Yet)
- Avoid concurrency/time/sync early → Your design delegates threading to SDL3 (hard dependency), so no need to expose sync/mutex yet. time can wait until http or logging needs timeouts.
- encoding/json, compress/*, archive/* → Great candidates, but they heavily depend on strings/bytes/io first.
- context → Useful for cancellation in http/fs, but add after http.
- Pitfalls to watch: net/http and os may touch platform-specific bits, but Go's implementations are mostly pure-Go and work fine in c-archive/static builds. Your handle registry will shield most issues.

Start with strings/bytes — implement a few ops (New, Len, CopyUTF8, Free, Append, etc.), get the registry humming, and you'll gain massive momentum. Once these four are bound and tested, FatStd will already feel substantially more "batteries-included" than raw C.

You've got the architecture nailed. This sequence will make the early wins addictive and set up everything else cleanly. Excited to see the first bindings land!
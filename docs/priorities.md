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

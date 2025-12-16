Good — your intuition is *almost* right. The missing piece is **what C is allowed to do vs what Go must control**.

Let me restate your expectation in precise terms, then show **the exact model that matches it**.

---

## What you expect (rephrased precisely)

From C, you want:

* An **opaque string type**
* The string is **created and owned by Go**
* C **never touches raw bytes directly**
* C manipulates the string **only by calling Go-backed functions**
* The string is freed either:

  * Explicitly by C, **or**
  * Automatically by Go’s GC

That is a **valid and common model** — *as long as you follow one hard rule*.

---

## The One Hard Rule (this is the crux)

> **C must never hold a raw pointer into Go-managed memory.**

If you accept that, your model works cleanly.

If you violate it, things break catastrophically.

---

## The Model You Actually Want: **Opaque Go Handle**

This is the correct mental model:

> C does **not** work with strings
> C works with **handles to Go objects**

---

## Concrete Design (Minimal, Clean, Safe)

### C side (what users see)

```c
typedef struct fat_GoString fat_GoString;

/* Creation */
fat_GoString *fat_StringFromC(const char *s, size_t len);

/* Inspection */
size_t        fat_StringLen(fat_GoString *s);

/* Operations */
int           fat_StringCompare(fat_GoString *a, fat_GoString *b);
fat_GoString *fat_StringConcat(fat_GoString *a, fat_GoString *b);
fat_GoString *fat_StringToUpper(fat_GoString *s);

/* Export (copy) */
char         *fat_StringCopyUTF8(fat_GoString *s, size_t *len);

/* Lifetime */
void          fat_StringFree(fat_GoString *s);
```

**Key property:**
C never sees `char *` unless it explicitly asks for a **copy**.

---

## Go side (what actually happens)

```go
type goString struct {
    s string
}
```

You store this in a **Go-managed table**, and C holds an **opaque handle**.

Example pattern:

```go
var strings = sync.Map{} // handle -> *goString
```

When C asks for an operation:

```go
//export fat_StringLen
func fat_StringLen(h unsafe.Pointer) C.size_t {
    gs := lookup(h)
    return C.size_t(len(gs.s))
}
```

---

## How memory & GC work here (important)

### Option 1: Explicit Free (recommended default)

* Go allocates the string
* C gets a handle
* C **must call** `fat_StringFree`
* Go removes the object from its table
* GC collects it naturally

This is deterministic and easy to reason about.

---

### Option 2: GC-only (possible, but risky)

* C gets a handle
* C never frees it
* Go GC eventually collects it

⚠️ This is dangerous unless:

* Handles are weak references
* Or you have a finalizer
* Or you accept leaks under pressure

**I strongly recommend explicit free.**

---

## Why this feels different than what you expected

You were implicitly assuming:

> “If Go owns it, GC will clean it up automatically.”

That’s true **only inside Go**.

The moment C holds a reference:

* Go has **no idea when C is done**
* The GC must assume it’s still live
* So *you* must provide a release signal

That’s why **every Go↔C boundary library does this**:

* SQLite
* libuv
* OpenSSL engines
* WebKit
* Vulkan loaders

---

## Why this is still “batteries included”

From a C developer’s point of view:

```c
fat_GoString *s = fat_HTTPGetBody(url);
fat_GoString *u = fat_StringToUpper(s);

log_Info("body=%s", fat_StringDebug(u));

fat_StringFree(u);
fat_StringFree(s);
```

No raw memory.
No encoding bugs.
No lifetime ambiguity.
No GC exposure.

That *is* a batteries-included string.

---

## What you should NOT do (important)

❌ Do **not** expose:

```c
const char *fat_StringData(...)
```

unless:

* It returns a **copy**
* Or it’s explicitly documented as ephemeral

That single function destroys your safety model.

---

## One last alignment check (yes/no)

Answering this locks the model:

> **Should C ever be able to mutate the contents of a Go string?**

* **No** → perfect fit, immutable handle model
* **Yes** → you must abandon Go-native strings and use a custom buffer type

If you answer **No**, your intuition and the model above are fully aligned.

If you want, next I can:

* Sketch the **exact handle table implementation**
* Show how this interacts with `http`, `log`, and `io`
* Define debug-mode misuse detection
* Compare this to how WebKit and SQLite do it

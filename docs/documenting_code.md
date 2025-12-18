<!--
  FatStd source documentation SOP.
  Keep this aligned with docs/design.md and the public C API.
-->

# Documenting Code (SOP)

## Goal

FatStd is C-first. Public headers are the contract. Documentation must make it easy to use the API correctly (types, ownership, lifetimes, error behavior) without reading the implementation.

## Scope and priorities

1. **Highest priority: public C API** under `include/fat/`.
2. **Medium priority: public-facing C wrappers** in `src/` when behavior is not obvious from the header.
3. **Low priority: internal/private modules** (Go packages under `pkg/`, internal C helpers): keep code self-documenting; add comments only when needed.

## Required style (C / Doxygen)

For every **public function** and **public symbol** in `include/fat/*.h`, use this Doxygen block directly above the declaration:

```c
/**
 * Brief one-sentence summary of what the function does.
 *
 * More detailed description, including behavior, algorithm hints if relevant,
 * and any important design decisions or efficiency notes.
 * Mention any pre-conditions, post conditions, and error conditions.
 * Also mention who owns what memory.
 *
 * @param name Description of the parameter, including valid ranges or ownership.
 * @param ...  Additional parameters.
 * @return Description of the return value.
 *
 * @note Any additional notes, e.g., thread-safety or client responsibilities.
 */
return_type FunctionName(param_type param_name, ...);
```

### What to cover (public C API)

Keep docs succinct, but always include:

- **Ownership**: who allocates, who frees, which `*_Free` applies, and whether outputs are “new handles”.
- **Lifetime**: handle validity and when it becomes invalid.
- **Embedded NUL policy**: when relevant, state “bytes may include embedded NULs” and how APIs treat them.
- **Error behavior**: what is validated, what is considered programmer error, and whether misuse is fail-fast.
- **Thread-safety**: if not safe, say so.
- **Ranges**: size limits, negative values, and any platform constraints.

### Handle-backed types

When introducing a new handle type in a public header:

- Document it once near the typedef with a short comment describing semantics.
- Ensure there is a matching `*_Free` for object-like handles.

## Go documentation

Go is an internal implementation detail. Use idiomatic Go comments sparingly:

- Exported Go symbols that are part of internal package APIs may get a single-line comment if non-obvious.
- Avoid long Go doc blocks; prefer clear naming and small functions.

## When to add or update docs

- **Every time you add/modify a public C function**, update its header comment in the same change.
- **If behavior drifts**, update the header doc (not just the implementation).
- **If the public API is extended**, update any module-level header docs and examples.

## Examples

Examples are encouraged when they clarify usage or ownership. Keep them short and compile-ish.

Recommended places:

- In the relevant header (for core types like strings/builders/readers).
- In docs/ if the example is longer than a few lines.

Example snippet pattern:

```c
// Create -> use -> free (ownership is explicit).
fat_String s = fat_StringNewUTF8("hello");
size_t n = fat_StringLenBytes(s);
char *buf = malloc(n + 1);
fat_StringCopyOut(s, buf, n);
buf[n] = '\0';
free(buf);
fat_StringFree(s);
```

## Internal code comments (C and Go)

Default stance: **no comments** if the code is clear.

Add comments only when they:

- Explain a non-obvious invariant or safety rule.
- Explain an intentional trade-off (performance, portability, correctness).
- Document a tricky boundary (Go/C, handle registry, ownership transitions).

## Review checklist (PR-quality)

- All new/changed `include/fat/*.h` symbols have Doxygen comments.
- Ownership and `*_Free` responsibilities are stated.
- Error behavior (including fail-fast conditions) is documented.
- Any significant drift between docs and behavior is corrected.

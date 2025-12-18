Refer to the following documents for context on the project:
1. docs/design.md
2. docs/onboarding_and_testing_functions.md
3. docs/documenting_code.md

Let's onboard more functions to fat_bytes:

TrimPrefix(s, prefix []byte) []byte
TrimSuffix(s, suffix []byte) []byte
Cut(s, sep []byte) (before, after []byte, found bool)
CutPrefix(s, prefix []byte) (after []byte, found bool) / CutSuffix
Fields(s []byte) [][]byte
Repeat(b []byte, count int) []byte
IndexByte(b []byte, c byte) int
IndexAny(s []byte, chars string) int
ToValidUTF8(s, replacement []byte) []byte

I expect the Go bindings, C bindings, and unit tests in Python.
Warn me if any of the specified functions violate the design or are impractical for C.
Refer to the following documents for context on the project:
1. docs/design.md
2. docs/onboarding_and_testing_functions.md

Let's continue working on the fat strings subsystem.
I want to onboard these functions:

TrimPrefix(s, prefix string) string 
TrimSuffix(s, suffix string) string
Cut(s, sep string) (before, after string, found bool)
CutPrefix(s, prefix string) (after string, found bool)
CutSuffix(s, prefix string) (after string, found bool)
Fields(s string) []string
Repeat(s string, count int) string
ContainsAny(s, chars string) bool
IndexAny(s, chars string) bool
ToValidUTF8(s, replacement string) string

I expect the Go bindings, C bindings, and a unit test in Python.

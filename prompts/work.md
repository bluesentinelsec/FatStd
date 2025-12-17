Refer to the following documents for context on the project:
1. docs/design.md
2. docs/onboarding_and_testing_functions.md

Let's continue working on the fat strings subsystem.
I want to onboard these functions:

Split(s, sep string) []string
SplitN(s, sep string, n int) []string
Join(elems []string, sep string) string
Replace(s, old, new string, n int) string
ReplaceAll(s, old, new string) string

I expect the Go bindings, C bindings, and a unit test in Python.

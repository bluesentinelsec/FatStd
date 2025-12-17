Refer to the following documents for context on the project:
1. docs/design.md
2. docs/onboarding_and_testing_functions.md

Let's continue working on the fat strings subsystem.
I want to onboard these functions:

ToLower(s string) string 
ToUpper(s string) string
Index(s, substr string) int
Count(s, substr string) int
Compare(a, b string) int
EqualFold(s, t string) bool

I expect the Go bindings, C bindings, and a unit test in Python.

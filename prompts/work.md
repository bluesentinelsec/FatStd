Refer to the following documents for context on the project:
1. docs/design.md
2. docs/onboarding_and_testing_functions.md
3. docs/documenting_code.md

Let's onboard more functions to fat_bytes:

bytes.Buffer (core methods: NewBuffer(buf []byte), NewBufferString(s string), Write(p []byte), WriteByte(c byte), WriteRune(r rune), WriteString(s string), Bytes() []byte, String() string, Len() int, Cap() int, Grow(n int), Reset(), Truncate(n int), Read(p []byte), Next(n int) []byte, ReadByte(), WriteTo(w io.Writer), ReadFrom(r io.Reader))

I expect the Go bindings, C bindings, and unit tests in Python.
Warn me if any of the specified functions violate the design or are impractical for C.
Refer to the following documents for context on the project:
1. docs/design.md
2. docs/onboarding_and_testing_functions.md
3. docs/documenting_code.md
4. docs/error_strategy.md
5. include/fat/*.h

Implement the bindings for the following compression/decompression packages in fatstd:

compress/bzip2
compress/flate
compress/gzip
compress/lzw
compress/zlib

I just need to be able to compress and decompress data.

I expect the Go bindings, C bindings, and unit tests in Python.
If any of the functions are a poor fit for C, use an alternative that honors the design.

When finished, add a brief tutorial doc showing how to use this module from the perspective of the caller under docs/
My repo has a directory: scripts/python_tests/

I want you to create a Python script that tests my C compatible shared library. Read this file for project context: docs/design.md

  - Library name + location after build (e.g. build_shared/libfatstd.dylib / .so / .dll)? Which OSes do you care about?
  Look at "CMakeLists.txt" and the GNU Makefile to see where build artifacts go. I care about Windows, macOS, and Linux, but most tests should be cross platform friendly.

  - How should it be loaded from Python: ctypes only, or is cffi/pybind11 acceptable?
  ctypes for runtime loading

  - Where is the C header that defines the public API (path), and what are the key functions you want to test first?
  Look in include/fat
  There's only one function in there, test it.

  - Any init/shutdown functions required (e.g. fatstd_init() / fatstd_free()), and any global state / thread-safety constraints?
  No

  - Error handling convention: return codes, errno, GetLastError, fatstd_last_error(), or something else?
  I like assert style tests.

  - Does the API allocate memory that the caller must free? If yes, what free function should
    Python call?
  - What “known good” behaviors should we assert (example inputs/outputs, edge cases)?
  - How do you want the script to find the library: via --lib /path/to/lib..., BUILD_DIR env
    var, or auto-detect build* folders?

  If you paste the header (or list of exported functions + signatures), I can stub the first
  test script right away.
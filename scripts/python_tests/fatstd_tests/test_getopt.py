from __future__ import annotations

import ctypes
import unittest

from fatstd_test_support import get_context


class TestGetopt(unittest.TestCase):
    def test_getopt_parses_simple_args(self) -> None:
        lib = get_context().lib

        getopt = lib.getopt
        getopt.argtypes = [ctypes.c_int, ctypes.POINTER(ctypes.c_char_p), ctypes.c_char_p]
        getopt.restype = ctypes.c_int

        optind = ctypes.c_int.in_dll(lib, "optind")
        opterr = ctypes.c_int.in_dll(lib, "opterr")
        optopt = ctypes.c_int.in_dll(lib, "optopt")
        optarg = ctypes.c_char_p.in_dll(lib, "optarg")

        # GNU getopt uses this to determine whether initialization has happened.
        try:
            getopt_initialized = ctypes.c_int.in_dll(lib, "__getopt_initialized")
        except ValueError:
            getopt_initialized = None

        argv = (ctypes.c_char_p * 5)()
        argv[0] = b"prog"
        argv[1] = b"-a"
        argv[2] = b"-b"
        argv[3] = b"val"
        argv[4] = b"rest"

        opterr.value = 0
        optind.value = 0
        optopt.value = 0
        optarg.value = None
        if getopt_initialized is not None:
            getopt_initialized.value = 0

        seen: list[tuple[str, str | None]] = []
        while True:
            c = getopt(len(argv), argv, b"ab:c")
            if c == -1:
                break
            if c == ord("a"):
                seen.append(("a", None))
            elif c == ord("b"):
                seen.append(("b", (optarg.value or b"").decode("utf-8", errors="strict")))
            else:
                self.fail(f"unexpected getopt result: {c}")

        self.assertEqual([("a", None), ("b", "val")], seen)
        self.assertEqual(4, optind.value)  # index of first non-option ("rest")


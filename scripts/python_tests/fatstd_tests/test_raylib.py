from __future__ import annotations

import ctypes
import unittest

from fatstd_test_support import get_context


class TestRaylib(unittest.TestCase):
    def test_get_random_value_symbol(self) -> None:
        lib = get_context().lib

        GetRandomValue = lib.GetRandomValue
        GetRandomValue.argtypes = [ctypes.c_int, ctypes.c_int]
        GetRandomValue.restype = ctypes.c_int

        v = GetRandomValue(1, 3)
        self.assertIn(v, (1, 2, 3))


from __future__ import annotations

import ctypes
import unittest

from fatstd_test_support import bind


class TestGo(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.fat_GoAdd = bind(
            "fat_GoAdd", argtypes=[ctypes.c_int, ctypes.c_int], restype=ctypes.c_int
        )

    def test_fat_GoAdd(self) -> None:
        got = self.fat_GoAdd(2, 3)
        self.assertEqual(5, got)


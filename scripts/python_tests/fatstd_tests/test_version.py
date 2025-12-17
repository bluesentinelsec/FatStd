from __future__ import annotations

import ctypes
import unittest

from fatstd_test_support import bind, get_context


class TestVersion(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.fat_VersionString = bind("fat_VersionString", argtypes=[], restype=ctypes.c_char_p)

    def test_fat_VersionString(self) -> None:
        raw = self.fat_VersionString()
        self.assertIsNotNone(raw, "fat_VersionString returned NULL")
        version = raw.decode("utf-8", errors="strict")
        self.assertEqual(get_context().expected_version, version)


from __future__ import annotations

import ctypes
import unittest

from fatstd_test_support import bind, fat_string_handle_type


class TestString(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        fat_string = fat_string_handle_type()

        cls.fat_StringNewUTF8 = bind(
            "fat_StringNewUTF8", argtypes=[ctypes.c_char_p], restype=fat_string
        )
        cls.fat_StringNewUTF8N = bind(
            "fat_StringNewUTF8N", argtypes=[ctypes.c_void_p, ctypes.c_size_t], restype=fat_string
        )
        cls.fat_StringClone = bind("fat_StringClone", argtypes=[fat_string], restype=fat_string)
        cls.fat_StringFree = bind("fat_StringFree", argtypes=[fat_string], restype=None)

    def test_create_clone_free_cstr(self) -> None:
        s1 = self.fat_StringNewUTF8(b"lorem ipsum")
        self.assertNotEqual(0, s1, "fat_StringNewUTF8 returned 0 handle")

        c1 = self.fat_StringClone(s1)
        self.assertNotEqual(0, c1, "fat_StringClone returned 0 handle")
        self.assertNotEqual(s1, c1, "fat_StringClone returned same handle")

        self.fat_StringFree(s1)
        self.fat_StringFree(c1)

    def test_create_clone_free_bytes(self) -> None:
        raw = ctypes.create_string_buffer(b"abc\x00def")
        s2 = self.fat_StringNewUTF8N(ctypes.addressof(raw), len(raw.raw))
        self.assertNotEqual(0, s2, "fat_StringNewUTF8N returned 0 handle")

        c2 = self.fat_StringClone(s2)
        self.assertNotEqual(0, c2, "fat_StringClone returned 0 handle")
        self.assertNotEqual(s2, c2, "fat_StringClone returned same handle")

        self.fat_StringFree(s2)
        self.fat_StringFree(c2)


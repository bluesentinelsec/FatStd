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
        cls.fat_StringContains = bind(
            "fat_StringContains", argtypes=[fat_string, fat_string], restype=ctypes.c_bool
        )
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

    def test_contains_basic(self) -> None:
        hay = self.fat_StringNewUTF8(b"lorem ipsum")
        self.assertNotEqual(0, hay)

        needle_yes = self.fat_StringNewUTF8(b"ipsum")
        self.assertNotEqual(0, needle_yes)

        needle_no = self.fat_StringNewUTF8(b"IPSUM")
        self.assertNotEqual(0, needle_no)

        self.assertTrue(self.fat_StringContains(hay, needle_yes))
        self.assertFalse(self.fat_StringContains(hay, needle_no))

        self.fat_StringFree(needle_no)
        self.fat_StringFree(needle_yes)
        self.fat_StringFree(hay)

    def test_contains_embedded_nul(self) -> None:
        hay_bytes = b"abc\x00def"
        hay_raw = ctypes.create_string_buffer(hay_bytes, len(hay_bytes))
        hay = self.fat_StringNewUTF8N(ctypes.addressof(hay_raw), len(hay_raw.raw))
        self.assertNotEqual(0, hay)

        needle_bytes = b"\x00de"
        needle_raw = ctypes.create_string_buffer(needle_bytes, len(needle_bytes))
        needle = self.fat_StringNewUTF8N(ctypes.addressof(needle_raw), len(needle_raw.raw))
        self.assertNotEqual(0, needle)

        self.assertTrue(self.fat_StringContains(hay, needle))

        self.fat_StringFree(needle)
        self.fat_StringFree(hay)

from __future__ import annotations

import ctypes
import unittest

from fatstd_test_support import bind, fat_string_handle_type


class TestStringBuilder(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        fat_string = fat_string_handle_type()
        fat_builder = fat_string_handle_type()

        cls.fat_StringNewUTF8 = bind(
            "fat_StringNewUTF8", argtypes=[ctypes.c_char_p], restype=fat_string
        )
        cls.fat_StringNewUTF8N = bind(
            "fat_StringNewUTF8N", argtypes=[ctypes.c_void_p, ctypes.c_size_t], restype=fat_string
        )
        cls.fat_StringContains = bind(
            "fat_StringContains", argtypes=[fat_string, fat_string], restype=ctypes.c_bool
        )
        cls.fat_StringFree = bind("fat_StringFree", argtypes=[fat_string], restype=None)

        cls.fat_StringBuilderNew = bind("fat_StringBuilderNew", argtypes=[], restype=fat_builder)
        cls.fat_StringBuilderFree = bind(
            "fat_StringBuilderFree", argtypes=[fat_builder], restype=None
        )
        cls.fat_StringBuilderCap = bind(
            "fat_StringBuilderCap", argtypes=[fat_builder], restype=ctypes.c_size_t
        )
        cls.fat_StringBuilderGrow = bind(
            "fat_StringBuilderGrow", argtypes=[fat_builder, ctypes.c_size_t], restype=None
        )
        cls.fat_StringBuilderLen = bind(
            "fat_StringBuilderLen", argtypes=[fat_builder], restype=ctypes.c_size_t
        )
        cls.fat_StringBuilderReset = bind(
            "fat_StringBuilderReset", argtypes=[fat_builder], restype=None
        )
        cls.fat_StringBuilderString = bind(
            "fat_StringBuilderString", argtypes=[fat_builder], restype=fat_string
        )
        cls.fat_StringBuilderWrite = bind(
            "fat_StringBuilderWrite",
            argtypes=[fat_builder, ctypes.c_void_p, ctypes.c_size_t],
            restype=ctypes.c_size_t,
        )
        cls.fat_StringBuilderWriteByte = bind(
            "fat_StringBuilderWriteByte", argtypes=[fat_builder, ctypes.c_uint8], restype=None
        )
        cls.fat_StringBuilderWriteString = bind(
            "fat_StringBuilderWriteString", argtypes=[fat_builder, fat_string], restype=ctypes.c_size_t
        )

    def _assert_string_equal(self, a, b) -> None:
        self.assertTrue(self.fat_StringContains(a, b))
        self.assertTrue(self.fat_StringContains(b, a))

    def test_builder_write_and_string(self) -> None:
        b = self.fat_StringBuilderNew()
        self.assertNotEqual(0, b)
        self.assertEqual(0, self.fat_StringBuilderLen(b))

        before_cap = self.fat_StringBuilderCap(b)
        self.fat_StringBuilderGrow(b, before_cap + 64)
        self.assertGreaterEqual(self.fat_StringBuilderCap(b), before_cap)

        s1 = self.fat_StringNewUTF8(b"hello")
        self.assertNotEqual(0, s1)
        n1 = self.fat_StringBuilderWriteString(b, s1)
        self.assertEqual(5, n1)

        self.fat_StringBuilderWriteByte(b, ord(b" "))

        raw = ctypes.create_string_buffer(b"wor\x00ld", 6)
        n2 = self.fat_StringBuilderWrite(b, ctypes.addressof(raw), 6)
        self.assertEqual(6, n2)

        self.assertEqual(5 + 1 + 6, self.fat_StringBuilderLen(b))

        out = self.fat_StringBuilderString(b)
        self.assertNotEqual(0, out)

        expected_bytes = b"hello wor\x00ld"
        expected_raw = ctypes.create_string_buffer(expected_bytes, len(expected_bytes))
        expected = self.fat_StringNewUTF8N(ctypes.addressof(expected_raw), len(expected_raw.raw))
        self.assertNotEqual(0, expected)
        self._assert_string_equal(out, expected)

        self.fat_StringFree(expected)
        self.fat_StringFree(out)
        self.fat_StringFree(s1)
        self.fat_StringBuilderFree(b)

    def test_builder_reset(self) -> None:
        b = self.fat_StringBuilderNew()
        self.assertNotEqual(0, b)

        s = self.fat_StringNewUTF8(b"abc")
        self.assertNotEqual(0, s)
        self.fat_StringBuilderWriteString(b, s)
        self.assertEqual(3, self.fat_StringBuilderLen(b))

        self.fat_StringBuilderReset(b)
        self.assertEqual(0, self.fat_StringBuilderLen(b))

        out = self.fat_StringBuilderString(b)
        self.assertNotEqual(0, out)
        empty = self.fat_StringNewUTF8(b"")
        self.assertNotEqual(0, empty)
        self._assert_string_equal(out, empty)

        self.fat_StringFree(empty)
        self.fat_StringFree(out)
        self.fat_StringFree(s)
        self.fat_StringBuilderFree(b)


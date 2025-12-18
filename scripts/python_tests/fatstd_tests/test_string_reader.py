from __future__ import annotations

import ctypes
import unittest

from fatstd_test_support import bind, fat_string_handle_type


class TestStringReader(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        fat_string = fat_string_handle_type()
        fat_reader = fat_string_handle_type()
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
        cls.fat_StringBuilderString = bind(
            "fat_StringBuilderString", argtypes=[fat_builder], restype=fat_string
        )

        cls.fat_StringReaderNew = bind(
            "fat_StringReaderNew", argtypes=[fat_string], restype=fat_reader
        )
        cls.fat_StringReaderFree = bind("fat_StringReaderFree", argtypes=[fat_reader], restype=None)
        cls.fat_StringReaderLen = bind(
            "fat_StringReaderLen", argtypes=[fat_reader], restype=ctypes.c_size_t
        )
        cls.fat_StringReaderSize = bind(
            "fat_StringReaderSize", argtypes=[fat_reader], restype=ctypes.c_int64
        )
        cls.fat_StringReaderReset = bind(
            "fat_StringReaderReset", argtypes=[fat_reader, fat_string], restype=None
        )
        cls.fat_StringReaderRead = bind(
            "fat_StringReaderRead",
            argtypes=[fat_reader, ctypes.c_void_p, ctypes.c_size_t, ctypes.POINTER(ctypes.c_bool)],
            restype=ctypes.c_size_t,
        )
        cls.fat_StringReaderReadAt = bind(
            "fat_StringReaderReadAt",
            argtypes=[
                fat_reader,
                ctypes.c_void_p,
                ctypes.c_size_t,
                ctypes.c_int64,
                ctypes.POINTER(ctypes.c_bool),
            ],
            restype=ctypes.c_size_t,
        )
        cls.fat_StringReaderReadByte = bind(
            "fat_StringReaderReadByte",
            argtypes=[fat_reader, ctypes.POINTER(ctypes.c_uint8), ctypes.POINTER(ctypes.c_bool)],
            restype=ctypes.c_bool,
        )
        cls.fat_StringReaderUnreadByte = bind(
            "fat_StringReaderUnreadByte", argtypes=[fat_reader], restype=None
        )
        cls.fat_StringReaderSeek = bind(
            "fat_StringReaderSeek",
            argtypes=[fat_reader, ctypes.c_int64, ctypes.c_int],
            restype=ctypes.c_int64,
        )
        cls.fat_StringReaderWriteToBuilder = bind(
            "fat_StringReaderWriteToBuilder",
            argtypes=[fat_reader, fat_builder],
            restype=ctypes.c_int64,
        )

    def _assert_string_equal(self, a, b) -> None:
        self.assertTrue(self.fat_StringContains(a, b))
        self.assertTrue(self.fat_StringContains(b, a))

    def test_len_size_read_reset_seek(self) -> None:
        raw_bytes = b"abc\x00def"
        raw = ctypes.create_string_buffer(raw_bytes, len(raw_bytes))
        s = self.fat_StringNewUTF8N(ctypes.addressof(raw), len(raw.raw))
        self.assertNotEqual(0, s)

        r = self.fat_StringReaderNew(s)
        self.assertNotEqual(0, r)

        self.assertEqual(len(raw_bytes), self.fat_StringReaderLen(r))
        self.assertEqual(len(raw_bytes), self.fat_StringReaderSize(r))

        out = bytearray()
        buf = ctypes.create_string_buffer(3, 3)
        while True:
            eof = ctypes.c_bool(False)
            n = self.fat_StringReaderRead(r, ctypes.addressof(buf), 3, ctypes.byref(eof))
            out += buf.raw[:n]
            if eof.value:
                break
        self.assertEqual(raw_bytes, bytes(out))
        self.assertEqual(0, self.fat_StringReaderLen(r))

        # Seek back to start and read again via ReadAt (doesn't change position).
        pos = self.fat_StringReaderSeek(r, 0, 0)  # SEEK_SET
        self.assertEqual(0, pos)

        buf2 = ctypes.create_string_buffer(3, 3)
        eof2 = ctypes.c_bool(False)
        n2 = self.fat_StringReaderReadAt(r, ctypes.addressof(buf2), 3, 4, ctypes.byref(eof2))
        self.assertEqual(3, n2)
        self.assertFalse(eof2.value)
        self.assertEqual(b"def", buf2.raw[:3])
        self.assertEqual(len(raw_bytes), self.fat_StringReaderLen(r))

        # Reset to a new string.
        s2 = self.fat_StringNewUTF8(b"xy")
        self.assertNotEqual(0, s2)
        self.fat_StringReaderReset(r, s2)
        self.assertEqual(2, self.fat_StringReaderLen(r))
        self.assertEqual(2, self.fat_StringReaderSize(r))

        self.fat_StringReaderFree(r)
        self.fat_StringFree(s2)
        self.fat_StringFree(s)

    def test_read_byte_and_unread_byte(self) -> None:
        s = self.fat_StringNewUTF8(b"ab")
        self.assertNotEqual(0, s)
        r = self.fat_StringReaderNew(s)
        self.assertNotEqual(0, r)

        b = ctypes.c_uint8(0)
        eof = ctypes.c_bool(False)
        ok = self.fat_StringReaderReadByte(r, ctypes.byref(b), ctypes.byref(eof))
        self.assertTrue(ok)
        self.assertFalse(eof.value)
        self.assertEqual(ord(b"a"), b.value)

        self.fat_StringReaderUnreadByte(r)
        self.assertEqual(2, self.fat_StringReaderLen(r))

        ok2 = self.fat_StringReaderReadByte(r, ctypes.byref(b), ctypes.byref(eof))
        self.assertTrue(ok2)
        self.assertEqual(ord(b"a"), b.value)

        ok3 = self.fat_StringReaderReadByte(r, ctypes.byref(b), ctypes.byref(eof))
        self.assertTrue(ok3)
        self.assertEqual(ord(b"b"), b.value)

        ok4 = self.fat_StringReaderReadByte(r, ctypes.byref(b), ctypes.byref(eof))
        self.assertFalse(ok4)
        self.assertTrue(eof.value)

        self.fat_StringReaderFree(r)
        self.fat_StringFree(s)

    def test_write_to_builder(self) -> None:
        s = self.fat_StringNewUTF8(b"hello")
        self.assertNotEqual(0, s)
        r = self.fat_StringReaderNew(s)
        self.assertNotEqual(0, r)

        # Consume 2 bytes, then write remaining to builder.
        buf = ctypes.create_string_buffer(2, 2)
        eof = ctypes.c_bool(False)
        n = self.fat_StringReaderRead(r, ctypes.addressof(buf), 2, ctypes.byref(eof))
        self.assertEqual(2, n)
        self.assertFalse(eof.value)

        b = self.fat_StringBuilderNew()
        self.assertNotEqual(0, b)
        written = self.fat_StringReaderWriteToBuilder(r, b)
        self.assertEqual(3, written)

        out = self.fat_StringBuilderString(b)
        self.assertNotEqual(0, out)
        expected = self.fat_StringNewUTF8(b"llo")
        self.assertNotEqual(0, expected)
        self._assert_string_equal(out, expected)

        self.fat_StringFree(expected)
        self.fat_StringFree(out)
        self.fat_StringBuilderFree(b)
        self.fat_StringReaderFree(r)
        self.fat_StringFree(s)


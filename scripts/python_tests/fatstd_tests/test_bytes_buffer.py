from __future__ import annotations

import ctypes
import unittest

from fatstd_test_support import bind, fat_string_handle_type


class TestBytesBuffer(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        fat_bytes = fat_string_handle_type()
        fat_bytes_buffer = fat_string_handle_type()
        fat_string = fat_string_handle_type()
        fat_string_reader = fat_string_handle_type()

        cls.fat_BytesNewN = bind(
            "fat_BytesNewN", argtypes=[ctypes.c_void_p, ctypes.c_size_t], restype=fat_bytes
        )
        cls.fat_BytesLen = bind("fat_BytesLen", argtypes=[fat_bytes], restype=ctypes.c_size_t)
        cls.fat_BytesCopyOut = bind(
            "fat_BytesCopyOut", argtypes=[fat_bytes, ctypes.c_void_p, ctypes.c_size_t], restype=ctypes.c_size_t
        )
        cls.fat_BytesFree = bind("fat_BytesFree", argtypes=[fat_bytes], restype=None)

        cls.fat_StringNewUTF8 = bind(
            "fat_StringNewUTF8", argtypes=[ctypes.c_char_p], restype=fat_string
        )
        cls.fat_StringLenBytes = bind(
            "fat_StringLenBytes", argtypes=[fat_string], restype=ctypes.c_size_t
        )
        cls.fat_StringCopyOut = bind(
            "fat_StringCopyOut", argtypes=[fat_string, ctypes.c_void_p, ctypes.c_size_t], restype=ctypes.c_size_t
        )
        cls.fat_StringFree = bind("fat_StringFree", argtypes=[fat_string], restype=None)

        cls.fat_StringReaderNew = bind(
            "fat_StringReaderNew", argtypes=[fat_string], restype=fat_string_reader
        )
        cls.fat_StringReaderFree = bind(
            "fat_StringReaderFree", argtypes=[fat_string_reader], restype=None
        )

        cls.fat_BytesBufferNew = bind("fat_BytesBufferNew", argtypes=[], restype=fat_bytes_buffer)
        cls.fat_BytesBufferNewBytes = bind(
            "fat_BytesBufferNewBytes", argtypes=[fat_bytes], restype=fat_bytes_buffer
        )
        cls.fat_BytesBufferNewN = bind(
            "fat_BytesBufferNewN", argtypes=[ctypes.c_void_p, ctypes.c_size_t], restype=fat_bytes_buffer
        )
        cls.fat_BytesBufferNewString = bind(
            "fat_BytesBufferNewString", argtypes=[fat_string], restype=fat_bytes_buffer
        )
        cls.fat_BytesBufferFree = bind(
            "fat_BytesBufferFree", argtypes=[fat_bytes_buffer], restype=None
        )
        cls.fat_BytesBufferLen = bind(
            "fat_BytesBufferLen", argtypes=[fat_bytes_buffer], restype=ctypes.c_size_t
        )
        cls.fat_BytesBufferCap = bind(
            "fat_BytesBufferCap", argtypes=[fat_bytes_buffer], restype=ctypes.c_size_t
        )
        cls.fat_BytesBufferGrow = bind(
            "fat_BytesBufferGrow", argtypes=[fat_bytes_buffer, ctypes.c_size_t], restype=None
        )
        cls.fat_BytesBufferReset = bind(
            "fat_BytesBufferReset", argtypes=[fat_bytes_buffer], restype=None
        )
        cls.fat_BytesBufferTruncate = bind(
            "fat_BytesBufferTruncate", argtypes=[fat_bytes_buffer, ctypes.c_size_t], restype=None
        )
        cls.fat_BytesBufferWrite = bind(
            "fat_BytesBufferWrite",
            argtypes=[fat_bytes_buffer, ctypes.c_void_p, ctypes.c_size_t],
            restype=ctypes.c_size_t,
        )
        cls.fat_BytesBufferWriteByte = bind(
            "fat_BytesBufferWriteByte", argtypes=[fat_bytes_buffer, ctypes.c_uint8], restype=None
        )
        cls.fat_BytesBufferWriteRune = bind(
            "fat_BytesBufferWriteRune", argtypes=[fat_bytes_buffer, ctypes.c_uint32], restype=ctypes.c_size_t
        )
        cls.fat_BytesBufferWriteString = bind(
            "fat_BytesBufferWriteString",
            argtypes=[fat_bytes_buffer, fat_string],
            restype=ctypes.c_size_t,
        )
        cls.fat_BytesBufferBytes = bind(
            "fat_BytesBufferBytes", argtypes=[fat_bytes_buffer], restype=fat_bytes
        )
        cls.fat_BytesBufferString = bind(
            "fat_BytesBufferString", argtypes=[fat_bytes_buffer], restype=fat_string
        )
        cls.fat_BytesBufferRead = bind(
            "fat_BytesBufferRead",
            argtypes=[fat_bytes_buffer, ctypes.c_void_p, ctypes.c_size_t, ctypes.POINTER(ctypes.c_bool)],
            restype=ctypes.c_size_t,
        )
        cls.fat_BytesBufferNext = bind(
            "fat_BytesBufferNext", argtypes=[fat_bytes_buffer, ctypes.c_size_t], restype=fat_bytes
        )
        cls.fat_BytesBufferReadByte = bind(
            "fat_BytesBufferReadByte",
            argtypes=[fat_bytes_buffer, ctypes.POINTER(ctypes.c_uint8), ctypes.POINTER(ctypes.c_bool)],
            restype=ctypes.c_bool,
        )
        cls.fat_BytesBufferWriteToBytesBuffer = bind(
            "fat_BytesBufferWriteToBytesBuffer",
            argtypes=[fat_bytes_buffer, fat_bytes_buffer],
            restype=ctypes.c_int64,
        )
        cls.fat_BytesBufferReadFromStringReader = bind(
            "fat_BytesBufferReadFromStringReader",
            argtypes=[fat_bytes_buffer, fat_string_reader],
            restype=ctypes.c_int64,
        )

    def _bytes_new(self, b: bytes) -> int:
        raw = ctypes.create_string_buffer(b, len(b))
        h = self.fat_BytesNewN(ctypes.addressof(raw), len(raw.raw))
        self.assertNotEqual(0, h)
        return h

    def _bytes_to_py(self, h: int) -> bytes:
        n = self.fat_BytesLen(h)
        if n == 0:
            return b""
        dst = ctypes.create_string_buffer(n, n)
        copied = self.fat_BytesCopyOut(h, ctypes.addressof(dst), len(dst.raw))
        self.assertEqual(n, copied)
        return dst.raw

    def _string_to_py_bytes(self, h: int) -> bytes:
        n = self.fat_StringLenBytes(h)
        if n == 0:
            return b""
        dst = ctypes.create_string_buffer(n, n)
        copied = self.fat_StringCopyOut(h, ctypes.addressof(dst), len(dst.raw))
        self.assertEqual(n, copied)
        return dst.raw

    def test_new_write_bytes_string_read_next(self) -> None:
        b0 = self.fat_BytesBufferNew()
        self.assertNotEqual(0, b0)
        self.assertEqual(0, self.fat_BytesBufferLen(b0))

        before_cap = self.fat_BytesBufferCap(b0)
        self.fat_BytesBufferGrow(b0, before_cap + 64)
        self.assertGreaterEqual(self.fat_BytesBufferCap(b0), before_cap)

        raw = ctypes.create_string_buffer(b"ab\x00", 3)
        n = self.fat_BytesBufferWrite(b0, ctypes.addressof(raw), 3)
        self.assertEqual(3, n)

        self.fat_BytesBufferWriteByte(b0, ord(b"c"))

        written = self.fat_BytesBufferWriteRune(b0, 0x2603)  # snowman
        self.assertEqual(3, written)

        s = self.fat_StringNewUTF8(b"Z")
        self.assertNotEqual(0, s)
        self.assertEqual(1, self.fat_BytesBufferWriteString(b0, s))
        self.fat_StringFree(s)

        snap = self.fat_BytesBufferBytes(b0)
        self.assertNotEqual(0, snap)
        self.assertEqual(b"ab\x00c\xe2\x98\x83Z", self._bytes_to_py(snap))

        snap_s = self.fat_BytesBufferString(b0)
        self.assertNotEqual(0, snap_s)
        self.assertEqual(b"ab\x00c\xe2\x98\x83Z", self._string_to_py_bytes(snap_s))

        chunk = self.fat_BytesBufferNext(b0, 2)
        self.assertEqual(b"ab", self._bytes_to_py(chunk))
        self.fat_BytesFree(chunk)

        eof = ctypes.c_bool(False)
        buf = ctypes.create_string_buffer(4, 4)
        n2 = self.fat_BytesBufferRead(b0, ctypes.addressof(buf), 4, ctypes.byref(eof))
        self.assertEqual(4, n2)
        self.assertFalse(eof.value)

        self.fat_StringFree(snap_s)
        self.fat_BytesFree(snap)
        self.fat_BytesBufferFree(b0)

    def test_reset_truncate_read_byte(self) -> None:
        raw = ctypes.create_string_buffer(b"hello", 5)
        b = self.fat_BytesBufferNewN(ctypes.addressof(raw), len(raw.raw))
        self.assertNotEqual(0, b)

        self.assertEqual(5, self.fat_BytesBufferLen(b))
        self.fat_BytesBufferTruncate(b, 3)
        self.assertEqual(3, self.fat_BytesBufferLen(b))

        c = ctypes.c_uint8(0)
        eof = ctypes.c_bool(False)
        ok = self.fat_BytesBufferReadByte(b, ctypes.byref(c), ctypes.byref(eof))
        self.assertTrue(ok)
        self.assertFalse(eof.value)
        self.assertEqual(ord(b"h"), c.value)

        self.fat_BytesBufferReset(b)
        self.assertEqual(0, self.fat_BytesBufferLen(b))

        self.fat_BytesBufferFree(b)

    def test_write_to_bytes_buffer_and_read_from_string_reader(self) -> None:
        src_raw = ctypes.create_string_buffer(b"hello", 5)
        src = self.fat_BytesBufferNewN(ctypes.addressof(src_raw), len(src_raw.raw))
        dst = self.fat_BytesBufferNew()
        self.assertNotEqual(0, src)
        self.assertNotEqual(0, dst)

        written = self.fat_BytesBufferWriteToBytesBuffer(src, dst)
        self.assertEqual(5, written)
        self.assertEqual(0, self.fat_BytesBufferLen(src))
        self.assertEqual(5, self.fat_BytesBufferLen(dst))

        s = self.fat_StringNewUTF8(b"xyz")
        r = self.fat_StringReaderNew(s)
        read_n = self.fat_BytesBufferReadFromStringReader(dst, r)
        self.assertEqual(3, read_n)

        out = self.fat_BytesBufferBytes(dst)
        self.assertEqual(b"helloxyz", self._bytes_to_py(out))

        self.fat_BytesFree(out)
        self.fat_StringReaderFree(r)
        self.fat_StringFree(s)
        self.fat_BytesBufferFree(dst)
        self.fat_BytesBufferFree(src)

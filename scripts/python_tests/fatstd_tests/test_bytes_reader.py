from __future__ import annotations

import ctypes
import unittest

from fatstd_test_support import bind, fat_string_handle_type


class TestBytesReader(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        fat_bytes = fat_string_handle_type()
        fat_bytes_reader = fat_string_handle_type()
        fat_bytes_buffer = fat_string_handle_type()

        cls.fat_BytesNewN = bind(
            "fat_BytesNewN", argtypes=[ctypes.c_void_p, ctypes.c_size_t], restype=fat_bytes
        )
        cls.fat_BytesLen = bind("fat_BytesLen", argtypes=[fat_bytes], restype=ctypes.c_size_t)
        cls.fat_BytesCopyOut = bind(
            "fat_BytesCopyOut", argtypes=[fat_bytes, ctypes.c_void_p, ctypes.c_size_t], restype=ctypes.c_size_t
        )
        cls.fat_BytesFree = bind("fat_BytesFree", argtypes=[fat_bytes], restype=None)

        cls.fat_BytesBufferNew = bind("fat_BytesBufferNew", argtypes=[], restype=fat_bytes_buffer)
        cls.fat_BytesBufferFree = bind(
            "fat_BytesBufferFree", argtypes=[fat_bytes_buffer], restype=None
        )
        cls.fat_BytesBufferBytes = bind(
            "fat_BytesBufferBytes", argtypes=[fat_bytes_buffer], restype=fat_bytes
        )

        cls.fat_BytesReaderNew = bind(
            "fat_BytesReaderNew", argtypes=[fat_bytes], restype=fat_bytes_reader
        )
        cls.fat_BytesReaderFree = bind(
            "fat_BytesReaderFree", argtypes=[fat_bytes_reader], restype=None
        )
        cls.fat_BytesReaderLen = bind(
            "fat_BytesReaderLen", argtypes=[fat_bytes_reader], restype=ctypes.c_size_t
        )
        cls.fat_BytesReaderSize = bind(
            "fat_BytesReaderSize", argtypes=[fat_bytes_reader], restype=ctypes.c_int64
        )
        cls.fat_BytesReaderReset = bind(
            "fat_BytesReaderReset", argtypes=[fat_bytes_reader, fat_bytes], restype=None
        )
        cls.fat_BytesReaderRead = bind(
            "fat_BytesReaderRead",
            argtypes=[
                fat_bytes_reader,
                ctypes.c_void_p,
                ctypes.c_size_t,
                ctypes.POINTER(ctypes.c_bool),
            ],
            restype=ctypes.c_size_t,
        )
        cls.fat_BytesReaderReadAt = bind(
            "fat_BytesReaderReadAt",
            argtypes=[
                fat_bytes_reader,
                ctypes.c_void_p,
                ctypes.c_size_t,
                ctypes.c_int64,
                ctypes.POINTER(ctypes.c_bool),
            ],
            restype=ctypes.c_size_t,
        )
        cls.fat_BytesReaderReadByte = bind(
            "fat_BytesReaderReadByte",
            argtypes=[
                fat_bytes_reader,
                ctypes.POINTER(ctypes.c_uint8),
                ctypes.POINTER(ctypes.c_bool),
            ],
            restype=ctypes.c_bool,
        )
        cls.fat_BytesReaderUnreadByte = bind(
            "fat_BytesReaderUnreadByte", argtypes=[fat_bytes_reader], restype=None
        )
        cls.fat_BytesReaderSeek = bind(
            "fat_BytesReaderSeek",
            argtypes=[fat_bytes_reader, ctypes.c_int64, ctypes.c_int],
            restype=ctypes.c_int64,
        )
        cls.fat_BytesReaderWriteToBytesBuffer = bind(
            "fat_BytesReaderWriteToBytesBuffer",
            argtypes=[fat_bytes_reader, fat_bytes_buffer],
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

    def test_len_size_read_reset_seek(self) -> None:
        src = b"abc\x00def"
        b = self._bytes_new(src)
        r = self.fat_BytesReaderNew(b)
        self.assertNotEqual(0, r)

        self.assertEqual(len(src), self.fat_BytesReaderLen(r))
        self.assertEqual(len(src), self.fat_BytesReaderSize(r))

        out = bytearray()
        buf = ctypes.create_string_buffer(3, 3)
        while True:
            eof = ctypes.c_bool(False)
            n = self.fat_BytesReaderRead(r, ctypes.addressof(buf), 3, ctypes.byref(eof))
            out += buf.raw[:n]
            if eof.value:
                break
        self.assertEqual(src, bytes(out))
        self.assertEqual(0, self.fat_BytesReaderLen(r))

        # Seek back to start and read again via ReadAt (doesn't change position).
        pos = self.fat_BytesReaderSeek(r, 0, 0)  # SEEK_SET
        self.assertEqual(0, pos)

        buf2 = ctypes.create_string_buffer(3, 3)
        eof2 = ctypes.c_bool(False)
        n2 = self.fat_BytesReaderReadAt(r, ctypes.addressof(buf2), 3, 4, ctypes.byref(eof2))
        self.assertEqual(3, n2)
        self.assertFalse(eof2.value)
        self.assertEqual(b"def", buf2.raw[:3])
        self.assertEqual(len(src), self.fat_BytesReaderLen(r))

        b2 = self._bytes_new(b"xy")
        self.fat_BytesReaderReset(r, b2)
        self.assertEqual(2, self.fat_BytesReaderLen(r))
        self.assertEqual(2, self.fat_BytesReaderSize(r))

        self.fat_BytesFree(b2)
        self.fat_BytesReaderFree(r)
        self.fat_BytesFree(b)

    def test_read_byte_unread_byte_and_write_to_buffer(self) -> None:
        b = self._bytes_new(b"hello")
        r = self.fat_BytesReaderNew(b)
        self.assertNotEqual(0, r)

        c = ctypes.c_uint8(0)
        eof = ctypes.c_bool(False)
        ok = self.fat_BytesReaderReadByte(r, ctypes.byref(c), ctypes.byref(eof))
        self.assertTrue(ok)
        self.assertFalse(eof.value)
        self.assertEqual(ord(b"h"), c.value)

        self.fat_BytesReaderUnreadByte(r)
        self.assertEqual(5, self.fat_BytesReaderLen(r))

        # Consume 2 bytes, then write remaining to bytes buffer.
        buf = ctypes.create_string_buffer(2, 2)
        eof2 = ctypes.c_bool(False)
        n = self.fat_BytesReaderRead(r, ctypes.addressof(buf), 2, ctypes.byref(eof2))
        self.assertEqual(2, n)
        self.assertFalse(eof2.value)

        out_buf = self.fat_BytesBufferNew()
        self.assertNotEqual(0, out_buf)
        written = self.fat_BytesReaderWriteToBytesBuffer(r, out_buf)
        self.assertEqual(3, written)

        out_bytes = self.fat_BytesBufferBytes(out_buf)
        self.assertNotEqual(0, out_bytes)
        self.assertEqual(b"llo", self._bytes_to_py(out_bytes))

        self.fat_BytesFree(out_bytes)
        self.fat_BytesBufferFree(out_buf)
        self.fat_BytesReaderFree(r)
        self.fat_BytesFree(b)


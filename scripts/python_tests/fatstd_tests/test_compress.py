from __future__ import annotations

import bz2
import ctypes
import gzip
import unittest
import zlib

from fatstd_test_support import bind, fat_string_handle_type


class TestCompress(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        fat_bytes = fat_string_handle_type()
        fat_error = fat_string_handle_type()

        cls.fat_BytesNewN = bind(
            "fat_BytesNewN", argtypes=[ctypes.c_void_p, ctypes.c_size_t], restype=fat_bytes
        )
        cls.fat_BytesLen = bind("fat_BytesLen", argtypes=[fat_bytes], restype=ctypes.c_size_t)
        cls.fat_BytesCopyOut = bind(
            "fat_BytesCopyOut", argtypes=[fat_bytes, ctypes.c_void_p, ctypes.c_size_t], restype=ctypes.c_size_t
        )
        cls.fat_BytesFree = bind("fat_BytesFree", argtypes=[fat_bytes], restype=None)
        cls.fat_ErrorFree = bind("fat_ErrorFree", argtypes=[fat_error], restype=None)

        cls.fat_Bzip2Decompress = bind(
            "fat_Bzip2Decompress",
            argtypes=[fat_bytes, ctypes.POINTER(fat_bytes), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_FlateCompress = bind(
            "fat_FlateCompress",
            argtypes=[fat_bytes, ctypes.POINTER(fat_bytes), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_FlateDecompress = bind(
            "fat_FlateDecompress",
            argtypes=[fat_bytes, ctypes.POINTER(fat_bytes), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_GzipCompress = bind(
            "fat_GzipCompress",
            argtypes=[fat_bytes, ctypes.POINTER(fat_bytes), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_GzipDecompress = bind(
            "fat_GzipDecompress",
            argtypes=[fat_bytes, ctypes.POINTER(fat_bytes), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_LzwCompress = bind(
            "fat_LzwCompress",
            argtypes=[fat_bytes, ctypes.c_int, ctypes.c_uint8, ctypes.POINTER(fat_bytes), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_LzwDecompress = bind(
            "fat_LzwDecompress",
            argtypes=[fat_bytes, ctypes.c_int, ctypes.c_uint8, ctypes.POINTER(fat_bytes), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_ZlibCompress = bind(
            "fat_ZlibCompress",
            argtypes=[fat_bytes, ctypes.POINTER(fat_bytes), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_ZlibDecompress = bind(
            "fat_ZlibDecompress",
            argtypes=[fat_bytes, ctypes.POINTER(fat_bytes), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )

        cls.FAT_OK = 0
        cls.FAT_ERR_SYNTAX = 1
        cls.FAT_ERR_RANGE = 2

    def _bytes_new(self, b: bytes) -> int:
        raw = ctypes.create_string_buffer(b, len(b))
        h = self.fat_BytesNewN(ctypes.addressof(raw), len(raw.raw))
        self.assertNotEqual(0, h)
        return h

    def _bytes_to_py(self, h: int) -> bytes:
        n = self.fat_BytesLen(h)
        if n == 0:
            dst = ctypes.create_string_buffer(1, 1)
            copied = self.fat_BytesCopyOut(h, ctypes.addressof(dst), 0)
            self.assertEqual(0, copied)
            return b""

        dst = ctypes.create_string_buffer(n, n)
        copied = self.fat_BytesCopyOut(h, ctypes.addressof(dst), len(dst.raw))
        self.assertEqual(n, copied)
        return dst.raw

    def _compress_roundtrip(self, compress_fn, decompress_fn, data: bytes) -> None:
        src = self._bytes_new(data)

        out = fat_string_handle_type()()
        err = fat_string_handle_type()()
        status = compress_fn(src, ctypes.byref(out), ctypes.byref(err))
        self.assertEqual(self.FAT_OK, status)
        self.assertEqual(0, err.value)
        self.assertNotEqual(0, out.value)

        out2 = fat_string_handle_type()()
        err2 = fat_string_handle_type()()
        status = decompress_fn(out, ctypes.byref(out2), ctypes.byref(err2))
        self.assertEqual(self.FAT_OK, status)
        self.assertEqual(0, err2.value)
        self.assertEqual(data, self._bytes_to_py(out2.value))

        self.fat_BytesFree(out2.value)
        self.fat_BytesFree(out.value)
        self.fat_BytesFree(src)

    def test_bzip2_decompress(self) -> None:
        data = b"fatstd bzip2 test payload"
        compressed = bz2.compress(data)

        src = self._bytes_new(compressed)
        out = fat_string_handle_type()()
        err = fat_string_handle_type()()
        status = self.fat_Bzip2Decompress(src, ctypes.byref(out), ctypes.byref(err))
        self.assertEqual(self.FAT_OK, status)
        self.assertEqual(0, err.value)
        self.assertEqual(data, self._bytes_to_py(out.value))

        self.fat_BytesFree(out.value)
        self.fat_BytesFree(src)

    def test_flate_roundtrip(self) -> None:
        data = b"raw deflate payload with \x00 bytes"
        self._compress_roundtrip(self.fat_FlateCompress, self.fat_FlateDecompress, data)

    def test_flate_decompress_python(self) -> None:
        data = b"python deflate payload"
        compressor = zlib.compressobj(wbits=-15)
        compressed = compressor.compress(data) + compressor.flush()

        src = self._bytes_new(compressed)
        out = fat_string_handle_type()()
        err = fat_string_handle_type()()
        status = self.fat_FlateDecompress(src, ctypes.byref(out), ctypes.byref(err))
        self.assertEqual(self.FAT_OK, status)
        self.assertEqual(0, err.value)
        self.assertEqual(data, self._bytes_to_py(out.value))

        self.fat_BytesFree(out.value)
        self.fat_BytesFree(src)

    def test_gzip_roundtrip(self) -> None:
        data = b"gzip payload from fatstd"
        self._compress_roundtrip(self.fat_GzipCompress, self.fat_GzipDecompress, data)

    def test_gzip_decompress_python(self) -> None:
        data = b"python gzip payload"
        compressed = gzip.compress(data)

        src = self._bytes_new(compressed)
        out = fat_string_handle_type()()
        err = fat_string_handle_type()()
        status = self.fat_GzipDecompress(src, ctypes.byref(out), ctypes.byref(err))
        self.assertEqual(self.FAT_OK, status)
        self.assertEqual(0, err.value)
        self.assertEqual(data, self._bytes_to_py(out.value))

        self.fat_BytesFree(out.value)
        self.fat_BytesFree(src)

    def test_lzw_roundtrip(self) -> None:
        data = b"lzw payload with repeated words lzw payload"
        self._compress_roundtrip(
            lambda src, out, err: self.fat_LzwCompress(src, 0, 8, out, err),
            lambda src, out, err: self.fat_LzwDecompress(src, 0, 8, out, err),
            data,
        )

    def test_lzw_invalid_params(self) -> None:
        data = b"lzw invalid"
        src = self._bytes_new(data)
        out = fat_string_handle_type()()
        err = fat_string_handle_type()()
        status = self.fat_LzwCompress(src, 2, 8, ctypes.byref(out), ctypes.byref(err))
        self.assertEqual(self.FAT_ERR_RANGE, status)
        self.assertNotEqual(0, err.value)
        self.fat_ErrorFree(err.value)
        self.fat_BytesFree(src)

    def test_zlib_roundtrip(self) -> None:
        data = b"zlib payload"
        self._compress_roundtrip(self.fat_ZlibCompress, self.fat_ZlibDecompress, data)

    def test_zlib_decompress_python(self) -> None:
        data = b"python zlib payload"
        compressed = zlib.compress(data)

        src = self._bytes_new(compressed)
        out = fat_string_handle_type()()
        err = fat_string_handle_type()()
        status = self.fat_ZlibDecompress(src, ctypes.byref(out), ctypes.byref(err))
        self.assertEqual(self.FAT_OK, status)
        self.assertEqual(0, err.value)
        self.assertEqual(data, self._bytes_to_py(out.value))

        self.fat_BytesFree(out.value)
        self.fat_BytesFree(src)

from __future__ import annotations

import ctypes
import unittest

from fatstd_test_support import bind, fat_string_handle_type


FAT_OK = 0
FAT_ERR_SYNTAX = 1
FAT_ERR_RANGE = 2


class TestConv(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        fat_string = fat_string_handle_type()
        fat_bytes = fat_string_handle_type()
        fat_error = fat_string_handle_type()

        cls.fat_StringNewUTF8 = bind(
            "fat_StringNewUTF8", argtypes=[ctypes.c_char_p], restype=fat_string
        )
        cls.fat_StringCopyOutCStr = bind(
            "fat_StringCopyOutCStr", argtypes=[fat_string, ctypes.c_void_p, ctypes.c_size_t], restype=ctypes.c_size_t
        )
        cls.fat_StringFree = bind("fat_StringFree", argtypes=[fat_string], restype=None)

        cls.fat_BytesNewN = bind(
            "fat_BytesNewN", argtypes=[ctypes.c_void_p, ctypes.c_size_t], restype=fat_bytes
        )
        cls.fat_BytesLen = bind("fat_BytesLen", argtypes=[fat_bytes], restype=ctypes.c_size_t)
        cls.fat_BytesCopyOut = bind(
            "fat_BytesCopyOut", argtypes=[fat_bytes, ctypes.c_void_p, ctypes.c_size_t], restype=ctypes.c_size_t
        )
        cls.fat_BytesFree = bind("fat_BytesFree", argtypes=[fat_bytes], restype=None)

        cls.fat_ErrorMessage = bind("fat_ErrorMessage", argtypes=[fat_error], restype=fat_string)
        cls.fat_ErrorFree = bind("fat_ErrorFree", argtypes=[fat_error], restype=None)

        cls.fat_ConvIntSize = bind("fat_ConvIntSize", argtypes=[], restype=ctypes.c_int)
        cls.fat_ConvFormatInt = bind(
            "fat_ConvFormatInt", argtypes=[ctypes.c_int64, ctypes.c_int], restype=fat_string
        )
        cls.fat_ConvParseInt = bind(
            "fat_ConvParseInt",
            argtypes=[
                fat_string,
                ctypes.c_int,
                ctypes.c_int,
                ctypes.POINTER(ctypes.c_int64),
                ctypes.POINTER(fat_error),
            ],
            restype=ctypes.c_int,
        )
        cls.fat_ConvParseBool = bind(
            "fat_ConvParseBool",
            argtypes=[fat_string, ctypes.POINTER(ctypes.c_int), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_ConvAppendInt = bind(
            "fat_ConvAppendInt", argtypes=[fat_bytes, ctypes.c_int64, ctypes.c_int], restype=fat_bytes
        )
        cls.fat_ConvQuote = bind("fat_ConvQuote", argtypes=[fat_string], restype=fat_string)
        cls.fat_ConvUnquote = bind(
            "fat_ConvUnquote",
            argtypes=[fat_string, ctypes.POINTER(fat_string), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_ConvUnquoteChar = bind(
            "fat_ConvUnquoteChar",
            argtypes=[
                fat_string,
                ctypes.c_uint8,
                ctypes.POINTER(ctypes.c_uint32),
                ctypes.POINTER(ctypes.c_bool),
                ctypes.POINTER(fat_string),
                ctypes.POINTER(fat_error),
            ],
            restype=ctypes.c_int,
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

    def _string_to_py(self, h: int) -> str:
        # Allocate a generous buffer for test strings.
        buf = ctypes.create_string_buffer(1024, 1024)
        n = self.fat_StringCopyOutCStr(h, ctypes.addressof(buf), len(buf.raw))
        return buf.raw[:n].decode("utf-8", errors="strict")

    def test_int_size_and_format_parse(self) -> None:
        self.assertIn(self.fat_ConvIntSize(), (32, 64))

        formatted = self.fat_ConvFormatInt(-123, 10)
        self.assertNotEqual(0, formatted)
        self.assertEqual("-123", self._string_to_py(formatted))
        self.fat_StringFree(formatted)

        s = self.fat_StringNewUTF8(b"123")
        self.assertNotEqual(0, s)
        out = ctypes.c_int64(0)
        err = fat_string_handle_type()(0)
        status = self.fat_ConvParseInt(s, 10, 64, ctypes.byref(out), ctypes.byref(err))
        self.assertEqual(FAT_OK, status)
        self.assertEqual(123, out.value)
        self.assertEqual(0, err.value)
        self.fat_StringFree(s)

        bad = self.fat_StringNewUTF8(b"not-a-number")
        self.assertNotEqual(0, bad)
        out2 = ctypes.c_int64(0)
        err2 = fat_string_handle_type()(0)
        status2 = self.fat_ConvParseInt(bad, 10, 64, ctypes.byref(out2), ctypes.byref(err2))
        self.assertEqual(FAT_ERR_SYNTAX, status2)
        self.assertNotEqual(0, err2.value)
        msg_s = self.fat_ErrorMessage(err2.value)
        self.assertNotEqual(0, msg_s)
        msg = self._string_to_py(msg_s)
        self.assertIn("invalid syntax", msg)
        self.fat_StringFree(msg_s)
        self.fat_ErrorFree(err2.value)
        self.fat_StringFree(bad)

    def test_parse_bool(self) -> None:
        s = self.fat_StringNewUTF8(b"true")
        out = ctypes.c_int(0)
        err = fat_string_handle_type()(0)
        st = self.fat_ConvParseBool(s, ctypes.byref(out), ctypes.byref(err))
        self.assertEqual(FAT_OK, st)
        self.assertEqual(1, out.value)
        self.assertEqual(0, err.value)
        self.fat_StringFree(s)

    def test_append_int_and_quote_unquote(self) -> None:
        dst = self._bytes_new(b"val=")
        out = self.fat_ConvAppendInt(dst, 42, 10)
        self.assertEqual(b"val=42", self._bytes_to_py(out))
        self.fat_BytesFree(out)
        self.fat_BytesFree(dst)

        s = self.fat_StringNewUTF8(b"hello")
        q = self.fat_ConvQuote(s)
        self.assertEqual('"hello"', self._string_to_py(q))

        unq_out = fat_string_handle_type()(0)
        err = fat_string_handle_type()(0)
        st = self.fat_ConvUnquote(q, ctypes.byref(unq_out), ctypes.byref(err))
        self.assertEqual(FAT_OK, st)
        self.assertEqual(0, err.value)
        self.assertNotEqual(0, unq_out.value)
        self.assertEqual("hello", self._string_to_py(unq_out.value))

        self.fat_StringFree(unq_out.value)
        self.fat_StringFree(q)
        self.fat_StringFree(s)

    def test_unquote_char(self) -> None:
        s = self.fat_StringNewUTF8(b"\\u2603xyz")
        out_rune = ctypes.c_uint32(0)
        out_multibyte = ctypes.c_bool(False)
        out_tail = fat_string_handle_type()(0)
        err = fat_string_handle_type()(0)

        st = self.fat_ConvUnquoteChar(
            s,
            0,
            ctypes.byref(out_rune),
            ctypes.byref(out_multibyte),
            ctypes.byref(out_tail),
            ctypes.byref(err),
        )
        self.assertEqual(FAT_OK, st)
        self.assertEqual(0, err.value)
        self.assertEqual(0x2603, out_rune.value)
        self.assertTrue(out_multibyte.value)
        self.assertEqual("xyz", self._string_to_py(out_tail.value))

        self.fat_StringFree(out_tail.value)
        self.fat_StringFree(s)


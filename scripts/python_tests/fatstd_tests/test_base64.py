from __future__ import annotations

import base64
import ctypes
import unittest

from fatstd_test_support import bind, fat_string_handle_type


FAT_OK = 0
FAT_ERR_SYNTAX = 1
FAT_ERR_RANGE = 2


class TestBase64(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        fat_handle = fat_string_handle_type()
        fat_string = fat_string_handle_type()
        fat_bytes = fat_string_handle_type()
        fat_error = fat_string_handle_type()

        cls.fat_StringNewUTF8 = bind("fat_StringNewUTF8", argtypes=[ctypes.c_char_p], restype=fat_string)
        cls.fat_StringCopyOutCStr = bind(
            "fat_StringCopyOutCStr", argtypes=[fat_string, ctypes.c_void_p, ctypes.c_size_t], restype=ctypes.c_size_t
        )
        cls.fat_StringFree = bind("fat_StringFree", argtypes=[fat_string], restype=None)

        cls.fat_BytesNewN = bind("fat_BytesNewN", argtypes=[ctypes.c_void_p, ctypes.c_size_t], restype=fat_bytes)
        cls.fat_BytesLen = bind("fat_BytesLen", argtypes=[fat_bytes], restype=ctypes.c_size_t)
        cls.fat_BytesCopyOut = bind(
            "fat_BytesCopyOut", argtypes=[fat_bytes, ctypes.c_void_p, ctypes.c_size_t], restype=ctypes.c_size_t
        )
        cls.fat_BytesFree = bind("fat_BytesFree", argtypes=[fat_bytes], restype=None)

        cls.fat_BytesBufferNew = bind("fat_BytesBufferNew", argtypes=[], restype=fat_handle)
        cls.fat_BytesBufferBytes = bind("fat_BytesBufferBytes", argtypes=[fat_handle], restype=fat_bytes)
        cls.fat_BytesBufferFree = bind("fat_BytesBufferFree", argtypes=[fat_handle], restype=None)

        cls.fat_ErrorMessage = bind("fat_ErrorMessage", argtypes=[fat_error], restype=fat_string)
        cls.fat_ErrorFree = bind("fat_ErrorFree", argtypes=[fat_error], restype=None)

        cls.fat_Base64EncodingNewUTF8 = bind(
            "fat_Base64EncodingNewUTF8",
            argtypes=[ctypes.c_char_p, ctypes.POINTER(fat_handle), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_Base64EncodingStrict = bind("fat_Base64EncodingStrict", argtypes=[fat_handle], restype=fat_handle)
        cls.fat_Base64EncodingWithPadding = bind(
            "fat_Base64EncodingWithPadding",
            argtypes=[fat_handle, ctypes.c_int32, ctypes.POINTER(fat_handle), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_Base64EncodingFree = bind("fat_Base64EncodingFree", argtypes=[fat_handle], restype=None)

        cls.fat_Base64EncodeToString = bind(
            "fat_Base64EncodeToString", argtypes=[fat_handle, fat_bytes], restype=fat_string
        )
        cls.fat_Base64DecodeString = bind(
            "fat_Base64DecodeString",
            argtypes=[fat_handle, fat_string, ctypes.POINTER(fat_bytes), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )

        cls.fat_Base64EncoderNewToBytesBuffer = bind(
            "fat_Base64EncoderNewToBytesBuffer",
            argtypes=[fat_handle, fat_handle, ctypes.POINTER(fat_handle), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_Base64EncoderWrite = bind(
            "fat_Base64EncoderWrite",
            argtypes=[
                fat_handle,
                ctypes.c_void_p,
                ctypes.c_size_t,
                ctypes.POINTER(ctypes.c_size_t),
                ctypes.POINTER(fat_error),
            ],
            restype=ctypes.c_int,
        )
        cls.fat_Base64EncoderClose = bind(
            "fat_Base64EncoderClose",
            argtypes=[fat_handle, ctypes.POINTER(fat_error)],
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
        buf = ctypes.create_string_buffer(4096, 4096)
        n = self.fat_StringCopyOutCStr(h, ctypes.addressof(buf), len(buf.raw))
        return buf.raw[:n].decode("utf-8", errors="strict")

    def _error_to_py(self, err_handle: int) -> str:
        msg_s = self.fat_ErrorMessage(err_handle)
        try:
            return self._string_to_py(msg_s)
        finally:
            self.fat_StringFree(msg_s)

    def _new_std_encoding(self) -> int:
        enc = fat_string_handle_type()(0)
        err = fat_string_handle_type()(0)
        st = self.fat_Base64EncodingNewUTF8(
            b"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/",
            ctypes.byref(enc),
            ctypes.byref(err),
        )
        self.assertEqual(FAT_OK, st)
        self.assertEqual(0, err.value)
        self.assertNotEqual(0, enc.value)
        return enc.value

    def test_new_encoding_validates_alphabet_length(self) -> None:
        enc = fat_string_handle_type()(0)
        err = fat_string_handle_type()(0)
        st = self.fat_Base64EncodingNewUTF8(b"short", ctypes.byref(enc), ctypes.byref(err))
        self.assertEqual(FAT_ERR_RANGE, st)
        self.assertEqual(0, enc.value)
        self.assertNotEqual(0, err.value)
        self.fat_ErrorFree(err.value)

    def test_encode_decode_roundtrip(self) -> None:
        enc = self._new_std_encoding()
        try:
            src_py = b"hello world\x00\x01\xff"
            src = self._bytes_new(src_py)
            try:
                s = self.fat_Base64EncodeToString(enc, src)
                self.assertNotEqual(0, s)
                encoded = self._string_to_py(s)
                self.fat_StringFree(s)

                self.assertEqual(base64.b64encode(src_py).decode("ascii"), encoded)

                inp = self.fat_StringNewUTF8(encoded.encode("ascii"))
                out = fat_string_handle_type()(0)
                err = fat_string_handle_type()(0)
                st = self.fat_Base64DecodeString(enc, inp, ctypes.byref(out), ctypes.byref(err))
                self.assertEqual(FAT_OK, st)
                self.assertEqual(0, err.value)
                self.assertNotEqual(0, out.value)
                self.assertEqual(src_py, self._bytes_to_py(out.value))
                self.fat_BytesFree(out.value)
                self.fat_StringFree(inp)
            finally:
                self.fat_BytesFree(src)
        finally:
            self.fat_Base64EncodingFree(enc)

    def test_strict_rejects_corrupt_input(self) -> None:
        enc = self._new_std_encoding()
        strict = self.fat_Base64EncodingStrict(enc)
        self.assertNotEqual(0, strict)
        try:
            inp = self.fat_StringNewUTF8(b"!!!!")
            out = fat_string_handle_type()(0)
            err = fat_string_handle_type()(0)
            st = self.fat_Base64DecodeString(strict, inp, ctypes.byref(out), ctypes.byref(err))
            self.assertEqual(FAT_ERR_SYNTAX, st)
            self.assertEqual(0, out.value)
            self.assertNotEqual(0, err.value)
            msg = self._error_to_py(err.value)
            self.assertIn("illegal base64", msg.lower())
            self.fat_ErrorFree(err.value)
            self.fat_StringFree(inp)
        finally:
            self.fat_Base64EncodingFree(strict)
            self.fat_Base64EncodingFree(enc)

    def test_encoder_stream_to_bytes_buffer(self) -> None:
        enc = self._new_std_encoding()
        try:
            buf = self.fat_BytesBufferNew()
            self.assertNotEqual(0, buf)

            encoder = fat_string_handle_type()(0)
            err = fat_string_handle_type()(0)
            st = self.fat_Base64EncoderNewToBytesBuffer(enc, buf, ctypes.byref(encoder), ctypes.byref(err))
            self.assertEqual(FAT_OK, st)
            self.assertEqual(0, err.value)
            self.assertNotEqual(0, encoder.value)

            payload = b"streaming payload"
            raw = ctypes.create_string_buffer(payload, len(payload))
            out_n = ctypes.c_size_t(0)
            err2 = fat_string_handle_type()(0)
            st2 = self.fat_Base64EncoderWrite(
                encoder.value,
                ctypes.addressof(raw),
                len(payload),
                ctypes.byref(out_n),
                ctypes.byref(err2),
            )
            self.assertEqual(FAT_OK, st2)
            self.assertEqual(0, err2.value)
            self.assertEqual(len(payload), out_n.value)

            err3 = fat_string_handle_type()(0)
            st3 = self.fat_Base64EncoderClose(encoder.value, ctypes.byref(err3))
            self.assertEqual(FAT_OK, st3)
            self.assertEqual(0, err3.value)

            out_bytes_h = self.fat_BytesBufferBytes(buf)
            try:
                out_bytes = self._bytes_to_py(out_bytes_h)
            finally:
                self.fat_BytesFree(out_bytes_h)
                self.fat_BytesBufferFree(buf)

            self.assertEqual(base64.b64encode(payload), out_bytes)
        finally:
            self.fat_Base64EncodingFree(enc)

    def test_with_padding_no_padding(self) -> None:
        enc = self._new_std_encoding()
        out_enc = fat_string_handle_type()(0)
        err = fat_string_handle_type()(0)
        st = self.fat_Base64EncodingWithPadding(enc, -1, ctypes.byref(out_enc), ctypes.byref(err))
        self.assertEqual(FAT_OK, st)
        self.assertEqual(0, err.value)
        self.assertNotEqual(0, out_enc.value)
        self.fat_Base64EncodingFree(out_enc.value)
        self.fat_Base64EncodingFree(enc)


from __future__ import annotations

import ctypes
import json as pyjson
import unittest

from fatstd_test_support import bind, fat_string_handle_type


FAT_OK = 0
FAT_ERR_SYNTAX = 1
FAT_ERR_EOF = 3

FAT_JSON_NULL = 0
FAT_JSON_BOOL = 1
FAT_JSON_NUMBER = 2
FAT_JSON_STRING = 3
FAT_JSON_ARRAY = 4
FAT_JSON_OBJECT = 5


class TestJson(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        fat_handle = fat_string_handle_type()
        fat_string = fat_string_handle_type()
        fat_bytes = fat_string_handle_type()
        fat_error = fat_string_handle_type()
        fat_string_array = fat_string_handle_type()

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

        cls.fat_BytesReaderNew = bind("fat_BytesReaderNew", argtypes=[fat_bytes], restype=fat_handle)
        cls.fat_BytesReaderFree = bind("fat_BytesReaderFree", argtypes=[fat_handle], restype=None)

        cls.fat_StringArrayLen = bind("fat_StringArrayLen", argtypes=[fat_string_array], restype=ctypes.c_size_t)
        cls.fat_StringArrayGet = bind(
            "fat_StringArrayGet", argtypes=[fat_string_array, ctypes.c_size_t], restype=fat_string
        )
        cls.fat_StringArrayFree = bind("fat_StringArrayFree", argtypes=[fat_string_array], restype=None)

        cls.fat_ErrorMessage = bind("fat_ErrorMessage", argtypes=[fat_error], restype=fat_string)
        cls.fat_ErrorFree = bind("fat_ErrorFree", argtypes=[fat_error], restype=None)

        cls.fat_JsonValid = bind("fat_JsonValid", argtypes=[fat_bytes], restype=ctypes.c_bool)
        cls.fat_JsonCompact = bind(
            "fat_JsonCompact",
            argtypes=[fat_bytes, ctypes.POINTER(fat_bytes), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_JsonIndent = bind(
            "fat_JsonIndent",
            argtypes=[fat_bytes, fat_string, fat_string, ctypes.POINTER(fat_bytes), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_JsonHTMLEscape = bind("fat_JsonHTMLEscape", argtypes=[fat_bytes], restype=fat_bytes)

        cls.fat_JsonUnmarshal = bind(
            "fat_JsonUnmarshal",
            argtypes=[fat_bytes, ctypes.POINTER(fat_handle), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_JsonMarshal = bind(
            "fat_JsonMarshal",
            argtypes=[fat_handle, ctypes.POINTER(fat_bytes), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_JsonMarshalIndent = bind(
            "fat_JsonMarshalIndent",
            argtypes=[fat_handle, fat_string, fat_string, ctypes.POINTER(fat_bytes), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_JsonValueFree = bind("fat_JsonValueFree", argtypes=[fat_handle], restype=None)
        cls.fat_JsonValueType = bind("fat_JsonValueType", argtypes=[fat_handle], restype=ctypes.c_int)
        cls.fat_JsonValueAsBool = bind(
            "fat_JsonValueAsBool", argtypes=[fat_handle, ctypes.POINTER(ctypes.c_int)], restype=None
        )
        cls.fat_JsonValueAsString = bind("fat_JsonValueAsString", argtypes=[fat_handle], restype=fat_string)
        cls.fat_JsonValueAsNumberString = bind(
            "fat_JsonValueAsNumberString", argtypes=[fat_handle], restype=fat_string
        )
        cls.fat_JsonArrayLen = bind("fat_JsonArrayLen", argtypes=[fat_handle], restype=ctypes.c_size_t)
        cls.fat_JsonArrayGet = bind("fat_JsonArrayGet", argtypes=[fat_handle, ctypes.c_size_t], restype=fat_handle)
        cls.fat_JsonObjectKeys = bind("fat_JsonObjectKeys", argtypes=[fat_handle], restype=fat_string_array)
        cls.fat_JsonObjectGet = bind(
            "fat_JsonObjectGet",
            argtypes=[fat_handle, fat_string, ctypes.POINTER(ctypes.c_bool), ctypes.POINTER(fat_handle)],
            restype=None,
        )

        cls.fat_JsonDecoderNewBytesReader = bind("fat_JsonDecoderNewBytesReader", argtypes=[fat_handle], restype=fat_handle)
        cls.fat_JsonDecoderFree = bind("fat_JsonDecoderFree", argtypes=[fat_handle], restype=None)
        cls.fat_JsonDecoderUseNumber = bind("fat_JsonDecoderUseNumber", argtypes=[fat_handle], restype=None)
        cls.fat_JsonDecoderInputOffset = bind("fat_JsonDecoderInputOffset", argtypes=[fat_handle], restype=ctypes.c_int64)
        cls.fat_JsonDecoderMore = bind("fat_JsonDecoderMore", argtypes=[fat_handle], restype=ctypes.c_bool)
        cls.fat_JsonDecoderBufferedBytes = bind("fat_JsonDecoderBufferedBytes", argtypes=[fat_handle], restype=fat_bytes)
        cls.fat_JsonDecoderDecodeValue = bind(
            "fat_JsonDecoderDecodeValue",
            argtypes=[fat_handle, ctypes.POINTER(fat_handle), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )

        cls.fat_JsonEncoderNewToBytesBuffer = bind("fat_JsonEncoderNewToBytesBuffer", argtypes=[fat_handle], restype=fat_handle)
        cls.fat_JsonEncoderFree = bind("fat_JsonEncoderFree", argtypes=[fat_handle], restype=None)
        cls.fat_JsonEncoderSetEscapeHTML = bind("fat_JsonEncoderSetEscapeHTML", argtypes=[fat_handle, ctypes.c_bool], restype=None)
        cls.fat_JsonEncoderSetIndent = bind("fat_JsonEncoderSetIndent", argtypes=[fat_handle, fat_string, fat_string], restype=None)
        cls.fat_JsonEncoderEncodeValue = bind(
            "fat_JsonEncoderEncodeValue", argtypes=[fat_handle, fat_handle, ctypes.POINTER(fat_error)], restype=ctypes.c_int
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
        buf = ctypes.create_string_buffer(8192, 8192)
        n = self.fat_StringCopyOutCStr(h, ctypes.addressof(buf), len(buf.raw))
        return buf.raw[:n].decode("utf-8", errors="strict")

    def _error_to_py(self, err_handle: int) -> str:
        msg_s = self.fat_ErrorMessage(err_handle)
        try:
            return self._string_to_py(msg_s)
        finally:
            self.fat_StringFree(msg_s)

    def _string_array_to_py(self, arr_handle: int) -> list[str]:
        n = int(self.fat_StringArrayLen(arr_handle))
        out: list[str] = []
        for i in range(n):
            s = self.fat_StringArrayGet(arr_handle, i)
            try:
                out.append(self._string_to_py(s))
            finally:
                self.fat_StringFree(s)
        return out

    def test_valid_compact_indent_and_html_escape(self) -> None:
        pretty = b'{\n  "a": 1,\n  "b": [true, null, "x"]\n}\n'
        b = self._bytes_new(pretty)
        try:
            self.assertTrue(self.fat_JsonValid(b))

            out = fat_string_handle_type()(0)
            err = fat_string_handle_type()(0)
            st = self.fat_JsonCompact(b, ctypes.byref(out), ctypes.byref(err))
            self.assertEqual(FAT_OK, st)
            self.assertEqual(0, err.value)
            self.assertNotEqual(0, out.value)
            compacted = self._bytes_to_py(out.value)
            self.fat_BytesFree(out.value)
            self.assertEqual(b'{"a":1,"b":[true,null,"x"]}', compacted)

            prefix = self.fat_StringNewUTF8(b"")
            indent = self.fat_StringNewUTF8(b"  ")
            out2 = fat_string_handle_type()(0)
            err2 = fat_string_handle_type()(0)
            st2 = self.fat_JsonIndent(b, prefix, indent, ctypes.byref(out2), ctypes.byref(err2))
            self.assertEqual(FAT_OK, st2)
            self.assertEqual(0, err2.value)
            indented = self._bytes_to_py(out2.value)
            self.fat_BytesFree(out2.value)
            self.fat_StringFree(prefix)
            self.fat_StringFree(indent)
            self.assertIn(b'\n  "a": 1,', indented)

            html_src = self._bytes_new(b'"<>&"')
            try:
                escaped_h = self.fat_JsonHTMLEscape(html_src)
                escaped = self._bytes_to_py(escaped_h)
                self.fat_BytesFree(escaped_h)
                self.assertEqual(b'"\\u003c\\u003e\\u0026"', escaped)
            finally:
                self.fat_BytesFree(html_src)
        finally:
            self.fat_BytesFree(b)

    def test_unmarshal_introspect_and_marshal(self) -> None:
        data = b'{"a":1,"b":[true,null,"x"],"c":{"d":2}}'
        b = self._bytes_new(data)
        val = fat_string_handle_type()(0)
        err = fat_string_handle_type()(0)
        st = self.fat_JsonUnmarshal(b, ctypes.byref(val), ctypes.byref(err))
        self.assertEqual(FAT_OK, st)
        self.assertEqual(0, err.value)
        self.assertNotEqual(0, val.value)
        try:
            self.assertEqual(FAT_JSON_OBJECT, self.fat_JsonValueType(val.value))

            keys_h = self.fat_JsonObjectKeys(val.value)
            try:
                self.assertEqual(["a", "b", "c"], self._string_array_to_py(keys_h))
            finally:
                self.fat_StringArrayFree(keys_h)

            key_b = self.fat_StringNewUTF8(b"b")
            found = ctypes.c_bool(False)
            out_v = fat_string_handle_type()(0)
            self.fat_JsonObjectGet(val.value, key_b, ctypes.byref(found), ctypes.byref(out_v))
            self.fat_StringFree(key_b)
            self.assertTrue(found.value)
            self.assertNotEqual(0, out_v.value)
            try:
                self.assertEqual(FAT_JSON_ARRAY, self.fat_JsonValueType(out_v.value))
                self.assertEqual(3, int(self.fat_JsonArrayLen(out_v.value)))

                # b[0] == true
                v0 = self.fat_JsonArrayGet(out_v.value, 0)
                try:
                    self.assertEqual(FAT_JSON_BOOL, self.fat_JsonValueType(v0))
                    out_bool = ctypes.c_int(0)
                    self.fat_JsonValueAsBool(v0, ctypes.byref(out_bool))
                    self.assertEqual(1, out_bool.value)
                finally:
                    self.fat_JsonValueFree(v0)
            finally:
                self.fat_JsonValueFree(out_v.value)

            out_bytes = fat_string_handle_type()(0)
            err2 = fat_string_handle_type()(0)
            st2 = self.fat_JsonMarshal(val.value, ctypes.byref(out_bytes), ctypes.byref(err2))
            self.assertEqual(FAT_OK, st2)
            self.assertEqual(0, err2.value)
            marshaled = self._bytes_to_py(out_bytes.value)
            self.fat_BytesFree(out_bytes.value)
            self.assertEqual(pyjson.loads(data.decode("utf-8")), pyjson.loads(marshaled.decode("utf-8")))
        finally:
            self.fat_JsonValueFree(val.value)
            self.fat_BytesFree(b)

    def test_decoder_and_encoder(self) -> None:
        payload = b'{"x":1} {"y":2}'
        b = self._bytes_new(payload)
        r = self.fat_BytesReaderNew(b)
        dec = self.fat_JsonDecoderNewBytesReader(r)
        self.assertNotEqual(0, dec)
        self.fat_JsonDecoderUseNumber(dec)
        try:
            v1 = fat_string_handle_type()(0)
            err = fat_string_handle_type()(0)
            st = self.fat_JsonDecoderDecodeValue(dec, ctypes.byref(v1), ctypes.byref(err))
            self.assertEqual(FAT_OK, st)
            self.assertEqual(0, err.value)
            self.assertNotEqual(0, v1.value)

            # After decoding first object, buffered should start with a space + next object.
            buf_h = self.fat_JsonDecoderBufferedBytes(dec)
            try:
                buf = self._bytes_to_py(buf_h)
            finally:
                self.fat_BytesFree(buf_h)
            self.assertTrue(buf.startswith(b" ") or buf == b"")

            # Encode value back out.
            out_buf = self.fat_BytesBufferNew()
            enc = self.fat_JsonEncoderNewToBytesBuffer(out_buf)
            self.fat_JsonEncoderSetEscapeHTML(enc, False)
            st2_err = fat_string_handle_type()(0)
            st2 = self.fat_JsonEncoderEncodeValue(enc, v1.value, ctypes.byref(st2_err))
            self.assertEqual(FAT_OK, st2)
            self.assertEqual(0, st2_err.value)

            out_bytes_h = self.fat_BytesBufferBytes(out_buf)
            try:
                out_bytes = self._bytes_to_py(out_bytes_h)
            finally:
                self.fat_BytesFree(out_bytes_h)
                self.fat_BytesBufferFree(out_buf)
                self.fat_JsonEncoderFree(enc)

            self.assertEqual({"x": 1}, pyjson.loads(out_bytes.decode("utf-8")))

            # Decode second object and then hit EOF.
            v2 = fat_string_handle_type()(0)
            err2 = fat_string_handle_type()(0)
            st3 = self.fat_JsonDecoderDecodeValue(dec, ctypes.byref(v2), ctypes.byref(err2))
            self.assertEqual(FAT_OK, st3)
            self.assertEqual(0, err2.value)
            self.fat_JsonValueFree(v2.value)

            v3 = fat_string_handle_type()(0)
            err3 = fat_string_handle_type()(0)
            st4 = self.fat_JsonDecoderDecodeValue(dec, ctypes.byref(v3), ctypes.byref(err3))
            self.assertEqual(FAT_ERR_EOF, st4)
            self.assertEqual(0, err3.value)
            self.assertEqual(0, v3.value)

            self.fat_JsonValueFree(v1.value)
        finally:
            self.fat_JsonDecoderFree(dec)
            self.fat_BytesReaderFree(r)
            self.fat_BytesFree(b)


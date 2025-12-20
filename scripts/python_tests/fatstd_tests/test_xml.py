from __future__ import annotations

import ctypes
import os
import tempfile
import unittest

from fatstd_test_support import bind, fat_string_handle_type


FAT_OK = 0
FAT_ERR_SYNTAX = 1
FAT_ERR_EOF = 3
FAT_ERR_OTHER = 100

FAT_XML_START_ELEMENT = 1
FAT_XML_END_ELEMENT = 2
FAT_XML_CHAR_DATA = 3
FAT_XML_PROC_INST = 6


class TestXml(unittest.TestCase):
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

        cls.fat_XmlDecoderNewBytes = bind("fat_XmlDecoderNewBytes", argtypes=[fat_bytes], restype=fat_handle)
        cls.fat_XmlDecoderOpenPathUTF8 = bind(
            "fat_XmlDecoderOpenPathUTF8", argtypes=[ctypes.c_char_p, ctypes.c_void_p, ctypes.c_void_p], restype=ctypes.c_int
        )
        cls.fat_XmlDecoderFree = bind(
            "fat_XmlDecoderFree", argtypes=[fat_handle, ctypes.c_void_p], restype=ctypes.c_int
        )
        cls.fat_XmlDecoderToken = bind(
            "fat_XmlDecoderToken", argtypes=[fat_handle, ctypes.c_void_p, ctypes.c_void_p], restype=ctypes.c_int
        )
        cls.fat_XmlDecoderInputOffset = bind("fat_XmlDecoderInputOffset", argtypes=[fat_handle], restype=ctypes.c_int64)
        cls.fat_XmlDecoderInputPos = bind(
            "fat_XmlDecoderInputPos", argtypes=[fat_handle, ctypes.c_void_p, ctypes.c_void_p], restype=None
        )

        cls.fat_XmlTokenFree = bind("fat_XmlTokenFree", argtypes=[fat_handle], restype=None)
        cls.fat_XmlTokenType = bind("fat_XmlTokenType", argtypes=[fat_handle], restype=ctypes.c_int)
        cls.fat_XmlTokenNameLocal = bind("fat_XmlTokenNameLocal", argtypes=[fat_handle], restype=fat_string)
        cls.fat_XmlStartElementAttrCount = bind("fat_XmlStartElementAttrCount", argtypes=[fat_handle], restype=ctypes.c_size_t)
        cls.fat_XmlStartElementAttrGet = bind(
            "fat_XmlStartElementAttrGet",
            argtypes=[fat_handle, ctypes.c_size_t, ctypes.c_void_p, ctypes.c_void_p, ctypes.c_void_p],
            restype=None,
        )
        cls.fat_XmlTokenBytes = bind("fat_XmlTokenBytes", argtypes=[fat_handle], restype=fat_bytes)
        cls.fat_XmlProcInstTarget = bind("fat_XmlProcInstTarget", argtypes=[fat_handle], restype=fat_string)

        cls.fat_XmlEscapeToBytesBuffer = bind(
            "fat_XmlEscapeToBytesBuffer", argtypes=[fat_handle, fat_bytes], restype=None
        )
        cls.fat_XmlEscapeTextToBytesBuffer = bind(
            "fat_XmlEscapeTextToBytesBuffer", argtypes=[fat_handle, fat_bytes, ctypes.c_void_p], restype=ctypes.c_int
        )

        cls.fat_XmlEncoderNewToBytesBuffer = bind(
            "fat_XmlEncoderNewToBytesBuffer", argtypes=[fat_handle], restype=fat_handle
        )
        cls.fat_XmlEncoderEncodeToken = bind(
            "fat_XmlEncoderEncodeToken", argtypes=[fat_handle, fat_handle, ctypes.c_void_p], restype=ctypes.c_int
        )
        cls.fat_XmlEncoderFlush = bind("fat_XmlEncoderFlush", argtypes=[fat_handle, ctypes.c_void_p], restype=ctypes.c_int)
        cls.fat_XmlEncoderClose = bind("fat_XmlEncoderClose", argtypes=[fat_handle, ctypes.c_void_p], restype=ctypes.c_int)

    def _bytes_new(self, b: bytes) -> int:
        raw = ctypes.create_string_buffer(b, len(b))
        h = self.fat_BytesNewN(ctypes.addressof(raw), len(raw.raw))
        self.assertNotEqual(0, h)
        return h

    def _bytes_to_py(self, h: int) -> bytes:
        n = int(self.fat_BytesLen(h))
        if n == 0:
            return b""
        dst = ctypes.create_string_buffer(n, n)
        copied = int(self.fat_BytesCopyOut(h, ctypes.addressof(dst), len(dst.raw)))
        self.assertEqual(n, copied)
        return dst.raw

    def _string_to_py(self, h: int) -> str:
        buf = ctypes.create_string_buffer(4096, 4096)
        n = int(self.fat_StringCopyOutCStr(h, ctypes.addressof(buf), len(buf.raw)))
        return buf.raw[:n].decode("utf-8", errors="strict")

    def _error_to_py(self, err_handle: int) -> str:
        msg_s = self.fat_ErrorMessage(err_handle)
        try:
            return self._string_to_py(msg_s)
        finally:
            self.fat_StringFree(msg_s)

    def test_decoder_tokens_and_extract_values(self) -> None:
        xml_text = b"<root><person id='7'><name>alice</name><active>true</active></person></root>"
        b = self._bytes_new(xml_text)
        dec = self.fat_XmlDecoderNewBytes(b)
        self.assertNotEqual(0, dec)
        try:
            name_val = None
            active_val = None
            in_name = False
            in_active = False
            saw_person = False

            while True:
                tok = fat_string_handle_type()(0)
                err = fat_string_handle_type()(0)
                st = int(self.fat_XmlDecoderToken(dec, ctypes.byref(tok), ctypes.byref(err)))
                if st == FAT_ERR_EOF:
                    self.assertEqual(0, err.value)
                    self.assertEqual(0, tok.value)
                    break
                self.assertEqual(FAT_OK, st)
                self.assertEqual(0, err.value)
                self.assertNotEqual(0, tok.value)

                try:
                    kind = int(self.fat_XmlTokenType(tok.value))
                    if kind == FAT_XML_PROC_INST:
                        # Ignore any XML declaration procinst.
                        target = self.fat_XmlProcInstTarget(tok.value)
                        self.fat_StringFree(target)
                        continue

                    if kind == FAT_XML_START_ELEMENT:
                        local_s = self.fat_XmlTokenNameLocal(tok.value)
                        try:
                            local = self._string_to_py(local_s)
                        finally:
                            self.fat_StringFree(local_s)

                        if local == "person":
                            saw_person = True
                            nattrs = int(self.fat_XmlStartElementAttrCount(tok.value))
                            self.assertEqual(1, nattrs)

                            nlocal = fat_string_handle_type()(0)
                            nspace = fat_string_handle_type()(0)
                            val = fat_string_handle_type()(0)
                            self.fat_XmlStartElementAttrGet(
                                tok.value, 0, ctypes.byref(nlocal), ctypes.byref(nspace), ctypes.byref(val)
                            )
                            try:
                                self.assertEqual("id", self._string_to_py(nlocal.value))
                                self.assertEqual("", self._string_to_py(nspace.value))
                                self.assertEqual("7", self._string_to_py(val.value))
                            finally:
                                self.fat_StringFree(nlocal.value)
                                self.fat_StringFree(nspace.value)
                                self.fat_StringFree(val.value)

                        in_name = local == "name"
                        in_active = local == "active"
                        continue

                    if kind == FAT_XML_END_ELEMENT:
                        in_name = False
                        in_active = False
                        continue

                    if kind == FAT_XML_CHAR_DATA:
                        data_h = self.fat_XmlTokenBytes(tok.value)
                        try:
                            data = self._bytes_to_py(data_h).decode("utf-8", errors="strict")
                        finally:
                            self.fat_BytesFree(data_h)
                        if in_name:
                            name_val = data
                        if in_active:
                            active_val = data
                finally:
                    self.fat_XmlTokenFree(tok.value)

            self.assertTrue(saw_person)
            self.assertEqual("alice", name_val)
            self.assertEqual("true", active_val)

            # Spot-check offset/pos are sane after reading.
            off = int(self.fat_XmlDecoderInputOffset(dec))
            self.assertGreater(off, 0)
            line = ctypes.c_int(0)
            col = ctypes.c_int(0)
            self.fat_XmlDecoderInputPos(dec, ctypes.byref(line), ctypes.byref(col))
            self.assertGreaterEqual(line.value, 1)
            self.assertGreaterEqual(col.value, 1)
        finally:
            err = fat_string_handle_type()(0)
            st = int(self.fat_XmlDecoderFree(dec, ctypes.byref(err)))
            self.assertEqual(FAT_OK, st)
            self.assertEqual(0, err.value)
            self.fat_BytesFree(b)

    def test_escape_and_escape_text_invalid_utf8(self) -> None:
        buf = self.fat_BytesBufferNew()
        src = self._bytes_new(b"<a&b>\"")
        try:
            self.fat_XmlEscapeToBytesBuffer(buf, src)
            out_h = self.fat_BytesBufferBytes(buf)
            try:
                out = self._bytes_to_py(out_h)
            finally:
                self.fat_BytesFree(out_h)
            self.assertIn(b"&lt;", out)
            self.assertIn(b"&amp;", out)
            self.assertIn(b"&gt;", out)

            buf2 = self.fat_BytesBufferNew()
            bad = self._bytes_new(b"<\xff>")
            try:
                err = fat_string_handle_type()(0)
                st = int(self.fat_XmlEscapeTextToBytesBuffer(buf2, bad, ctypes.byref(err)))
                self.assertEqual(FAT_OK, st)
                self.assertEqual(0, err.value)

                out2_h = self.fat_BytesBufferBytes(buf2)
                try:
                    out2 = self._bytes_to_py(out2_h)
                finally:
                    self.fat_BytesFree(out2_h)
                self.assertIn(b"&lt;", out2)
                self.assertIn(b"&gt;", out2)
            finally:
                self.fat_BytesFree(bad)
                self.fat_BytesBufferFree(buf2)
        finally:
            self.fat_BytesFree(src)
            self.fat_BytesBufferFree(buf)

    def test_open_path_utf8_and_reencode_tokens(self) -> None:
        xml_text = b"<root><name>alice</name></root>"
        with tempfile.TemporaryDirectory() as td:
            path = os.path.join(td, "doc.xml")
            with open(path, "wb") as f:
                f.write(xml_text)

            dec = fat_string_handle_type()(0)
            err = fat_string_handle_type()(0)
            st = int(self.fat_XmlDecoderOpenPathUTF8(path.encode("utf-8"), ctypes.byref(dec), ctypes.byref(err)))
            self.assertEqual(FAT_OK, st)
            self.assertEqual(0, err.value)
            self.assertNotEqual(0, dec.value)

            out_buf = self.fat_BytesBufferNew()
            enc = self.fat_XmlEncoderNewToBytesBuffer(out_buf)
            self.assertNotEqual(0, enc)

            try:
                got_name = None
                in_name = False

                while True:
                    tok = fat_string_handle_type()(0)
                    err2 = fat_string_handle_type()(0)
                    st2 = int(self.fat_XmlDecoderToken(dec.value, ctypes.byref(tok), ctypes.byref(err2)))
                    if st2 == FAT_ERR_EOF:
                        self.assertEqual(0, err2.value)
                        break
                    self.assertEqual(FAT_OK, st2)
                    self.assertEqual(0, err2.value)
                    self.assertNotEqual(0, tok.value)

                    try:
                        kind = int(self.fat_XmlTokenType(tok.value))
                        if kind == FAT_XML_START_ELEMENT:
                            local_s = self.fat_XmlTokenNameLocal(tok.value)
                            try:
                                local = self._string_to_py(local_s)
                            finally:
                                self.fat_StringFree(local_s)
                            in_name = local == "name"
                        elif kind == FAT_XML_END_ELEMENT:
                            in_name = False
                        elif kind == FAT_XML_CHAR_DATA and in_name:
                            data_h = self.fat_XmlTokenBytes(tok.value)
                            try:
                                got_name = self._bytes_to_py(data_h).decode("utf-8", errors="strict")
                            finally:
                                self.fat_BytesFree(data_h)

                        # Re-encode everything we read.
                        err3 = fat_string_handle_type()(0)
                        st3 = int(self.fat_XmlEncoderEncodeToken(enc, tok.value, ctypes.byref(err3)))
                        self.assertEqual(FAT_OK, st3)
                        self.assertEqual(0, err3.value)
                    finally:
                        self.fat_XmlTokenFree(tok.value)

                err_flush = fat_string_handle_type()(0)
                stf = int(self.fat_XmlEncoderFlush(enc, ctypes.byref(err_flush)))
                self.assertEqual(FAT_OK, stf)
                self.assertEqual(0, err_flush.value)

                err_close = fat_string_handle_type()(0)
                stc = int(self.fat_XmlEncoderClose(enc, ctypes.byref(err_close)))
                self.assertEqual(FAT_OK, stc)
                self.assertEqual(0, err_close.value)

                out_h = self.fat_BytesBufferBytes(out_buf)
                try:
                    out = self._bytes_to_py(out_h)
                finally:
                    self.fat_BytesFree(out_h)

                self.assertEqual("alice", got_name)
                self.assertIn(b"<name>", out)
                self.assertIn(b"alice", out)

                # Ensure output is parseable.
                out_bytes_h = self._bytes_new(out)
                out_dec = self.fat_XmlDecoderNewBytes(out_bytes_h)
                try:
                    while True:
                        tok2 = fat_string_handle_type()(0)
                        err4 = fat_string_handle_type()(0)
                        st4 = int(self.fat_XmlDecoderToken(out_dec, ctypes.byref(tok2), ctypes.byref(err4)))
                        if st4 == FAT_ERR_EOF:
                            self.assertEqual(0, err4.value)
                            break
                        self.assertEqual(FAT_OK, st4)
                        self.assertEqual(0, err4.value)
                        self.fat_XmlTokenFree(tok2.value)
                finally:
                    err_free2 = fat_string_handle_type()(0)
                    st_free2 = int(self.fat_XmlDecoderFree(out_dec, ctypes.byref(err_free2)))
                    self.assertEqual(FAT_OK, st_free2)
                    self.assertEqual(0, err_free2.value)
                self.fat_BytesFree(out_bytes_h)
            finally:
                err_free = fat_string_handle_type()(0)
                st_free = int(self.fat_XmlDecoderFree(dec.value, ctypes.byref(err_free)))
                self.assertEqual(FAT_OK, st_free)
                self.assertEqual(0, err_free.value)
                self.fat_BytesBufferFree(out_buf)

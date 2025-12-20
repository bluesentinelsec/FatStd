from __future__ import annotations

import csv as pycsv
import ctypes
import io
import unittest

from fatstd_test_support import bind, fat_string_handle_type


FAT_OK = 0
FAT_ERR_SYNTAX = 1
FAT_ERR_EOF = 3


class TestCsv(unittest.TestCase):
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

        cls.fat_StringArrayLen = bind("fat_StringArrayLen", argtypes=[fat_string_array], restype=ctypes.c_size_t)
        cls.fat_StringArrayGet = bind("fat_StringArrayGet", argtypes=[fat_string_array, ctypes.c_size_t], restype=fat_string)
        cls.fat_StringArrayFree = bind("fat_StringArrayFree", argtypes=[fat_string_array], restype=None)

        cls.fat_ErrorMessage = bind("fat_ErrorMessage", argtypes=[fat_error], restype=fat_string)
        cls.fat_ErrorFree = bind("fat_ErrorFree", argtypes=[fat_error], restype=None)

        cls.fat_CsvReaderNewBytes = bind("fat_CsvReaderNewBytes", argtypes=[fat_bytes], restype=fat_handle)
        cls.fat_CsvReaderFree = bind("fat_CsvReaderFree", argtypes=[fat_handle], restype=None)
        cls.fat_CsvReaderRead = bind(
            "fat_CsvReaderRead",
            argtypes=[
                fat_handle,
                ctypes.POINTER(fat_string_array),
                ctypes.POINTER(ctypes.c_bool),
                ctypes.POINTER(fat_error),
            ],
            restype=ctypes.c_int,
        )
        cls.fat_CsvReaderFieldPos = bind(
            "fat_CsvReaderFieldPos",
            argtypes=[fat_handle, ctypes.c_int, ctypes.POINTER(ctypes.c_int), ctypes.POINTER(ctypes.c_int)],
            restype=None,
        )
        cls.fat_CsvReaderInputOffset = bind("fat_CsvReaderInputOffset", argtypes=[fat_handle], restype=ctypes.c_int64)

        cls.fat_CsvWriterNewToBytesBuffer = bind("fat_CsvWriterNewToBytesBuffer", argtypes=[fat_handle], restype=fat_handle)
        cls.fat_CsvWriterWriteRecord = bind(
            "fat_CsvWriterWriteRecord",
            argtypes=[fat_handle, ctypes.POINTER(fat_string), ctypes.c_size_t, ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_CsvWriterFlush = bind("fat_CsvWriterFlush", argtypes=[fat_handle], restype=None)
        cls.fat_CsvWriterError = bind(
            "fat_CsvWriterError", argtypes=[fat_handle, ctypes.POINTER(fat_error)], restype=ctypes.c_int
        )
        cls.fat_CsvWriterFree = bind("fat_CsvWriterFree", argtypes=[fat_handle], restype=None)

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

    def _record_to_py(self, arr_handle: int) -> list[str]:
        n = int(self.fat_StringArrayLen(arr_handle))
        out: list[str] = []
        for i in range(n):
            s = self.fat_StringArrayGet(arr_handle, i)
            try:
                out.append(self._string_to_py(s))
            finally:
                self.fat_StringFree(s)
        return out

    def test_writer_roundtrip_with_python_csv(self) -> None:
        buf = self.fat_BytesBufferNew()
        w = self.fat_CsvWriterNewToBytesBuffer(buf)
        self.assertNotEqual(0, w)

        fields_py = ["a,b", 'c"d', "x"]
        fields = (fat_string_handle_type() * len(fields_py))()
        for i, v in enumerate(fields_py):
            fields[i] = self.fat_StringNewUTF8(v.encode("utf-8"))

        err = fat_string_handle_type()(0)
        st = self.fat_CsvWriterWriteRecord(w, fields, len(fields_py), ctypes.byref(err))
        self.assertEqual(FAT_OK, st)
        self.assertEqual(0, err.value)

        self.fat_CsvWriterFlush(w)
        err2 = fat_string_handle_type()(0)
        st2 = self.fat_CsvWriterError(w, ctypes.byref(err2))
        self.assertEqual(FAT_OK, st2)
        self.assertEqual(0, err2.value)

        out_bytes_h = self.fat_BytesBufferBytes(buf)
        try:
            out_bytes = self._bytes_to_py(out_bytes_h)
        finally:
            self.fat_BytesFree(out_bytes_h)
            self.fat_CsvWriterFree(w)
            self.fat_BytesBufferFree(buf)
            for i in range(len(fields_py)):
                self.fat_StringFree(fields[i])

        sio = io.StringIO(newline="")
        pw = pycsv.writer(sio, lineterminator="\n")
        pw.writerow(fields_py)
        self.assertEqual(sio.getvalue().encode("utf-8"), out_bytes)

    def test_reader_reads_python_csv(self) -> None:
        sio = io.StringIO(newline="")
        pw = pycsv.writer(sio, lineterminator="\n")
        pw.writerow(["a", "b"])
        pw.writerow(["c,d", "e"])
        data = sio.getvalue().encode("utf-8")

        b = self._bytes_new(data)
        r = self.fat_CsvReaderNewBytes(b)
        self.assertNotEqual(0, r)
        try:
            got: list[list[str]] = []

            rec = fat_string_handle_type()(0)
            eof = ctypes.c_bool(False)
            err = fat_string_handle_type()(0)
            st1 = self.fat_CsvReaderRead(r, ctypes.byref(rec), ctypes.byref(eof), ctypes.byref(err))
            self.assertEqual(FAT_OK, st1)
            self.assertEqual(0, err.value)
            self.assertFalse(eof.value)
            self.assertNotEqual(0, rec.value)

            # Validate FieldPos/InputOffset on the first record "a,b\n"
            line = ctypes.c_int(0)
            col = ctypes.c_int(0)
            self.fat_CsvReaderFieldPos(r, 0, ctypes.byref(line), ctypes.byref(col))
            self.assertEqual((1, 1), (line.value, col.value))
            self.fat_CsvReaderFieldPos(r, 1, ctypes.byref(line), ctypes.byref(col))
            self.assertEqual((1, 3), (line.value, col.value))
            self.assertEqual(4, int(self.fat_CsvReaderInputOffset(r)))

            got.append(self._record_to_py(rec.value))
            self.fat_StringArrayFree(rec.value)

            rec2 = fat_string_handle_type()(0)
            eof2 = ctypes.c_bool(False)
            err2 = fat_string_handle_type()(0)
            st2 = self.fat_CsvReaderRead(r, ctypes.byref(rec2), ctypes.byref(eof2), ctypes.byref(err2))
            self.assertEqual(FAT_OK, st2)
            self.assertEqual(0, err2.value)
            got.append(self._record_to_py(rec2.value))
            self.fat_StringArrayFree(rec2.value)

            rec3 = fat_string_handle_type()(0)
            eof3 = ctypes.c_bool(False)
            err3 = fat_string_handle_type()(0)
            st3 = self.fat_CsvReaderRead(r, ctypes.byref(rec3), ctypes.byref(eof3), ctypes.byref(err3))
            self.assertEqual(FAT_ERR_EOF, st3)
            self.assertTrue(eof3.value)
            self.assertEqual(0, err3.value)
            self.assertEqual(0, rec3.value)

            self.assertEqual([["a", "b"], ["c,d", "e"]], got)
        finally:
            self.fat_CsvReaderFree(r)
            self.fat_BytesFree(b)

    def test_reader_parse_error(self) -> None:
        bad = self._bytes_new(b'a,"b"c\n')
        r = self.fat_CsvReaderNewBytes(bad)
        try:
            rec = fat_string_handle_type()(0)
            eof = ctypes.c_bool(False)
            err = fat_string_handle_type()(0)
            st = self.fat_CsvReaderRead(r, ctypes.byref(rec), ctypes.byref(eof), ctypes.byref(err))
            self.assertEqual(FAT_ERR_SYNTAX, st)
            self.assertFalse(eof.value)
            self.assertEqual(0, rec.value)
            self.assertNotEqual(0, err.value)
            msg = self._error_to_py(err.value).lower()
            self.assertTrue("quoted-field" in msg or "quote" in msg)
            self.fat_ErrorFree(err.value)
        finally:
            self.fat_CsvReaderFree(r)
            self.fat_BytesFree(bad)


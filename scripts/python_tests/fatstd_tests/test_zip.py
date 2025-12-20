from __future__ import annotations

import ctypes
import io
import unittest
import zipfile

from fatstd_test_support import bind, fat_string_handle_type


FAT_OK = 0
FAT_ERR_SYNTAX = 1
FAT_ERR_EOF = 3


class TestZip(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        fat_handle = fat_string_handle_type()
        fat_string = fat_string_handle_type()
        fat_bytes = fat_string_handle_type()
        fat_error = fat_string_handle_type()

        cls.fat_StringNewUTF8 = bind(
            "fat_StringNewUTF8", argtypes=[ctypes.c_char_p], restype=fat_string
        )
        cls.fat_StringCopyOutCStr = bind(
            "fat_StringCopyOutCStr",
            argtypes=[fat_string, ctypes.c_void_p, ctypes.c_size_t],
            restype=ctypes.c_size_t,
        )
        cls.fat_StringFree = bind("fat_StringFree", argtypes=[fat_string], restype=None)

        cls.fat_BytesNewN = bind(
            "fat_BytesNewN", argtypes=[ctypes.c_void_p, ctypes.c_size_t], restype=fat_bytes
        )
        cls.fat_BytesLen = bind("fat_BytesLen", argtypes=[fat_bytes], restype=ctypes.c_size_t)
        cls.fat_BytesCopyOut = bind(
            "fat_BytesCopyOut",
            argtypes=[fat_bytes, ctypes.c_void_p, ctypes.c_size_t],
            restype=ctypes.c_size_t,
        )
        cls.fat_BytesFree = bind("fat_BytesFree", argtypes=[fat_bytes], restype=None)

        cls.fat_BytesBufferNew = bind("fat_BytesBufferNew", argtypes=[], restype=fat_handle)
        cls.fat_BytesBufferBytes = bind(
            "fat_BytesBufferBytes", argtypes=[fat_handle], restype=fat_bytes
        )
        cls.fat_BytesBufferFree = bind("fat_BytesBufferFree", argtypes=[fat_handle], restype=None)

        cls.fat_ErrorMessage = bind("fat_ErrorMessage", argtypes=[fat_error], restype=fat_string)
        cls.fat_ErrorFree = bind("fat_ErrorFree", argtypes=[fat_error], restype=None)

        cls.fat_ZipWriterNewToBytesBuffer = bind(
            "fat_ZipWriterNewToBytesBuffer",
            argtypes=[fat_handle, ctypes.POINTER(fat_handle), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_ZipWriterAddBytes = bind(
            "fat_ZipWriterAddBytes",
            argtypes=[fat_handle, fat_string, fat_bytes, ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_ZipWriterClose = bind(
            "fat_ZipWriterClose",
            argtypes=[fat_handle, ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )

        cls.fat_ZipReaderNewBytes = bind(
            "fat_ZipReaderNewBytes",
            argtypes=[fat_bytes, ctypes.POINTER(fat_handle), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_ZipReaderNumFiles = bind(
            "fat_ZipReaderNumFiles", argtypes=[fat_handle], restype=ctypes.c_size_t
        )
        cls.fat_ZipReaderFileByIndex = bind(
            "fat_ZipReaderFileByIndex", argtypes=[fat_handle, ctypes.c_size_t], restype=fat_handle
        )
        cls.fat_ZipReaderFree = bind(
            "fat_ZipReaderFree",
            argtypes=[fat_handle, ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )

        cls.fat_ZipFileName = bind("fat_ZipFileName", argtypes=[fat_handle], restype=fat_string)
        cls.fat_ZipFileOpen = bind(
            "fat_ZipFileOpen",
            argtypes=[fat_handle, ctypes.POINTER(fat_handle), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_ZipFileFree = bind("fat_ZipFileFree", argtypes=[fat_handle], restype=None)

        cls.fat_ZipFileReaderRead = bind(
            "fat_ZipFileReaderRead",
            argtypes=[
                fat_handle,
                ctypes.c_void_p,
                ctypes.c_size_t,
                ctypes.POINTER(ctypes.c_size_t),
                ctypes.POINTER(ctypes.c_bool),
                ctypes.POINTER(fat_error),
            ],
            restype=ctypes.c_int,
        )
        cls.fat_ZipFileReaderClose = bind(
            "fat_ZipFileReaderClose",
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

    def _zip_read_all(self, reader_handle: int) -> bytes:
        chunks: list[bytes] = []
        tmp = ctypes.create_string_buffer(128, 128)
        while True:
            out_n = ctypes.c_size_t(0)
            out_eof = ctypes.c_bool(False)
            out_err = fat_string_handle_type()(0)
            st = self.fat_ZipFileReaderRead(
                reader_handle,
                ctypes.addressof(tmp),
                len(tmp.raw),
                ctypes.byref(out_n),
                ctypes.byref(out_eof),
                ctypes.byref(out_err),
            )
            if st == FAT_OK:
                self.assertEqual(0, out_err.value)
                if out_n.value:
                    chunks.append(tmp.raw[: out_n.value])
                continue
            if st == FAT_ERR_EOF:
                self.assertEqual(0, out_err.value)
                self.assertTrue(out_eof.value)
                if out_n.value:
                    chunks.append(tmp.raw[: out_n.value])
                break

            self.assertNotEqual(0, out_err.value)
            msg = self._error_to_py(out_err.value)
            self.fat_ErrorFree(out_err.value)
            raise AssertionError(f"zip read failed: status={st} err={msg}")
        return b"".join(chunks)

    def test_writer_to_bytes_buffer_and_python_can_read(self) -> None:
        buf = self.fat_BytesBufferNew()
        self.assertNotEqual(0, buf)

        writer = fat_string_handle_type()(0)
        err = fat_string_handle_type()(0)
        st = self.fat_ZipWriterNewToBytesBuffer(buf, ctypes.byref(writer), ctypes.byref(err))
        self.assertEqual(FAT_OK, st)
        self.assertEqual(0, err.value)
        self.assertNotEqual(0, writer.value)

        name1 = self.fat_StringNewUTF8(b"a.txt")
        data1 = self._bytes_new(b"hello\n")
        err1 = fat_string_handle_type()(0)
        st1 = self.fat_ZipWriterAddBytes(writer.value, name1, data1, ctypes.byref(err1))
        self.assertEqual(FAT_OK, st1)
        self.assertEqual(0, err1.value)
        self.fat_BytesFree(data1)
        self.fat_StringFree(name1)

        name2 = self.fat_StringNewUTF8(b"dir/b.bin")
        data2 = self._bytes_new(bytes(range(0, 256)))
        err2 = fat_string_handle_type()(0)
        st2 = self.fat_ZipWriterAddBytes(writer.value, name2, data2, ctypes.byref(err2))
        self.assertEqual(FAT_OK, st2)
        self.assertEqual(0, err2.value)
        self.fat_BytesFree(data2)
        self.fat_StringFree(name2)

        err_close = fat_string_handle_type()(0)
        st_close = self.fat_ZipWriterClose(writer.value, ctypes.byref(err_close))
        self.assertEqual(FAT_OK, st_close)
        self.assertEqual(0, err_close.value)

        zip_bytes_handle = self.fat_BytesBufferBytes(buf)
        self.assertNotEqual(0, zip_bytes_handle)
        try:
            zip_bytes = self._bytes_to_py(zip_bytes_handle)
        finally:
            self.fat_BytesFree(zip_bytes_handle)
            self.fat_BytesBufferFree(buf)

        with zipfile.ZipFile(io.BytesIO(zip_bytes), mode="r") as zf:
            self.assertEqual({"a.txt", "dir/b.bin"}, set(zf.namelist()))
            self.assertEqual(b"hello\n", zf.read("a.txt"))
            self.assertEqual(bytes(range(0, 256)), zf.read("dir/b.bin"))

    def test_reader_can_read_python_zip_bytes(self) -> None:
        out = io.BytesIO()
        with zipfile.ZipFile(out, mode="w", compression=zipfile.ZIP_DEFLATED) as zf:
            zf.writestr("x.txt", "payload")
            zf.writestr("nested/y.txt", b"\x00\x01\x02")
        zip_bytes = out.getvalue()

        zip_bytes_handle = self._bytes_new(zip_bytes)
        try:
            reader = fat_string_handle_type()(0)
            err = fat_string_handle_type()(0)
            st = self.fat_ZipReaderNewBytes(zip_bytes_handle, ctypes.byref(reader), ctypes.byref(err))
            self.assertEqual(FAT_OK, st)
            self.assertEqual(0, err.value)
            self.assertNotEqual(0, reader.value)

            nfiles = self.fat_ZipReaderNumFiles(reader.value)
            self.assertGreaterEqual(nfiles, 2)

            wanted = "x.txt"
            found_file = 0
            for i in range(nfiles):
                f = self.fat_ZipReaderFileByIndex(reader.value, i)
                self.assertNotEqual(0, f)
                name_h = self.fat_ZipFileName(f)
                try:
                    name = self._string_to_py(name_h)
                finally:
                    self.fat_StringFree(name_h)
                if name == wanted:
                    found_file = f
                    break
                self.fat_ZipFileFree(f)

            self.assertNotEqual(0, found_file)

            file_reader = fat_string_handle_type()(0)
            err2 = fat_string_handle_type()(0)
            st2 = self.fat_ZipFileOpen(found_file, ctypes.byref(file_reader), ctypes.byref(err2))
            self.assertEqual(FAT_OK, st2)
            self.assertEqual(0, err2.value)
            self.assertNotEqual(0, file_reader.value)

            payload = self._zip_read_all(file_reader.value)
            self.assertEqual(b"payload", payload)

            err_close = fat_string_handle_type()(0)
            st_close = self.fat_ZipFileReaderClose(file_reader.value, ctypes.byref(err_close))
            self.assertEqual(FAT_OK, st_close)
            self.assertEqual(0, err_close.value)

            self.fat_ZipFileFree(found_file)

            err_free = fat_string_handle_type()(0)
            st_free = self.fat_ZipReaderFree(reader.value, ctypes.byref(err_free))
            self.assertEqual(FAT_OK, st_free)
            self.assertEqual(0, err_free.value)
        finally:
            self.fat_BytesFree(zip_bytes_handle)

    def test_reader_rejects_non_zip(self) -> None:
        bad = self._bytes_new(b"not a zip")
        try:
            reader = fat_string_handle_type()(0)
            err = fat_string_handle_type()(0)
            st = self.fat_ZipReaderNewBytes(bad, ctypes.byref(reader), ctypes.byref(err))
            self.assertEqual(FAT_ERR_SYNTAX, st)
            self.assertEqual(0, reader.value)
            self.assertNotEqual(0, err.value)
            msg = self._error_to_py(err.value)
            self.assertIn("zip", msg.lower())
            self.fat_ErrorFree(err.value)
        finally:
            self.fat_BytesFree(bad)


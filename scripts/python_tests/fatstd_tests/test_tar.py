from __future__ import annotations

import ctypes
import io
import tarfile
import unittest

from fatstd_test_support import bind, fat_string_handle_type


FAT_OK = 0
FAT_ERR_EOF = 3


class TestTar(unittest.TestCase):
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

        cls.fat_ErrorFree = bind("fat_ErrorFree", argtypes=[fat_error], restype=None)

        cls.fat_TarWriterNewToBytesBuffer = bind(
            "fat_TarWriterNewToBytesBuffer",
            argtypes=[fat_handle, ctypes.POINTER(fat_handle), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_TarWriterAddBytes = bind(
            "fat_TarWriterAddBytes",
            argtypes=[fat_handle, fat_string, fat_bytes, ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_TarWriterClose = bind(
            "fat_TarWriterClose",
            argtypes=[fat_handle, ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )

        cls.fat_TarReaderNewBytes = bind(
            "fat_TarReaderNewBytes",
            argtypes=[fat_bytes, ctypes.POINTER(fat_handle), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_TarReaderNext = bind(
            "fat_TarReaderNext",
            argtypes=[
                fat_handle,
                ctypes.POINTER(fat_handle),
                ctypes.POINTER(ctypes.c_bool),
                ctypes.POINTER(fat_error),
            ],
            restype=ctypes.c_int,
        )
        cls.fat_TarReaderRead = bind(
            "fat_TarReaderRead",
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
        cls.fat_TarReaderFree = bind(
            "fat_TarReaderFree",
            argtypes=[fat_handle, ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )

        cls.fat_TarHeaderName = bind("fat_TarHeaderName", argtypes=[fat_handle], restype=fat_string)
        cls.fat_TarHeaderTypeflag = bind(
            "fat_TarHeaderTypeflag", argtypes=[fat_handle], restype=ctypes.c_uint8
        )
        cls.fat_TarHeaderSize = bind(
            "fat_TarHeaderSize", argtypes=[fat_handle], restype=ctypes.c_int64
        )
        cls.fat_TarHeaderFree = bind("fat_TarHeaderFree", argtypes=[fat_handle], restype=None)

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

    def _tar_read_entry_all(self, reader_handle: int) -> bytes:
        chunks: list[bytes] = []
        tmp = ctypes.create_string_buffer(128, 128)
        while True:
            out_n = ctypes.c_size_t(0)
            out_eof = ctypes.c_bool(False)
            out_err = fat_string_handle_type()(0)
            st = self.fat_TarReaderRead(
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
            self.fat_ErrorFree(out_err.value)
            raise AssertionError(f"tar read failed: status={st}")
        return b"".join(chunks)

    def test_writer_to_bytes_buffer_and_python_can_read(self) -> None:
        buf = self.fat_BytesBufferNew()
        self.assertNotEqual(0, buf)

        writer = fat_string_handle_type()(0)
        err = fat_string_handle_type()(0)
        st = self.fat_TarWriterNewToBytesBuffer(buf, ctypes.byref(writer), ctypes.byref(err))
        self.assertEqual(FAT_OK, st)
        self.assertEqual(0, err.value)
        self.assertNotEqual(0, writer.value)

        name1 = self.fat_StringNewUTF8(b"a.txt")
        data1 = self._bytes_new(b"hello\n")
        err1 = fat_string_handle_type()(0)
        st1 = self.fat_TarWriterAddBytes(writer.value, name1, data1, ctypes.byref(err1))
        self.assertEqual(FAT_OK, st1)
        self.assertEqual(0, err1.value)
        self.fat_BytesFree(data1)
        self.fat_StringFree(name1)

        err_close = fat_string_handle_type()(0)
        st_close = self.fat_TarWriterClose(writer.value, ctypes.byref(err_close))
        self.assertEqual(FAT_OK, st_close)
        self.assertEqual(0, err_close.value)

        tar_bytes_handle = self.fat_BytesBufferBytes(buf)
        self.assertNotEqual(0, tar_bytes_handle)
        try:
            tar_bytes = self._bytes_to_py(tar_bytes_handle)
        finally:
            self.fat_BytesFree(tar_bytes_handle)
            self.fat_BytesBufferFree(buf)

        with tarfile.open(fileobj=io.BytesIO(tar_bytes), mode="r:*") as tf:
            names = [m.name for m in tf.getmembers()]
            self.assertIn("a.txt", names)
            extracted = tf.extractfile("a.txt")
            self.assertIsNotNone(extracted)
            self.assertEqual(b"hello\n", extracted.read())

    def test_reader_can_read_python_tar_bytes(self) -> None:
        out = io.BytesIO()
        with tarfile.open(fileobj=out, mode="w") as tf:
            data = b"payload"
            info = tarfile.TarInfo("x.txt")
            info.size = len(data)
            tf.addfile(info, io.BytesIO(data))

        tar_bytes = out.getvalue()
        tar_bytes_handle = self._bytes_new(tar_bytes)
        try:
            reader = fat_string_handle_type()(0)
            err = fat_string_handle_type()(0)
            st = self.fat_TarReaderNewBytes(tar_bytes_handle, ctypes.byref(reader), ctypes.byref(err))
            self.assertEqual(FAT_OK, st)
            self.assertEqual(0, err.value)
            self.assertNotEqual(0, reader.value)

            got: dict[str, bytes] = {}
            while True:
                hdr = fat_string_handle_type()(0)
                eof = ctypes.c_bool(False)
                err2 = fat_string_handle_type()(0)
                st2 = self.fat_TarReaderNext(reader.value, ctypes.byref(hdr), ctypes.byref(eof), ctypes.byref(err2))
                if st2 == FAT_ERR_EOF:
                    self.assertTrue(eof.value)
                    self.assertEqual(0, hdr.value)
                    self.assertEqual(0, err2.value)
                    break
                self.assertEqual(FAT_OK, st2)
                self.assertEqual(0, err2.value)
                self.assertNotEqual(0, hdr.value)

                name_h = self.fat_TarHeaderName(hdr.value)
                try:
                    name = self._string_to_py(name_h)
                finally:
                    self.fat_StringFree(name_h)

                _typeflag = self.fat_TarHeaderTypeflag(hdr.value)
                _size = self.fat_TarHeaderSize(hdr.value)

                body = self._tar_read_entry_all(reader.value)
                got[name] = body

                self.fat_TarHeaderFree(hdr.value)

            self.assertEqual({ "x.txt" }, set(got.keys()))
            self.assertEqual(b"payload", got["x.txt"])

            err_free = fat_string_handle_type()(0)
            st_free = self.fat_TarReaderFree(reader.value, ctypes.byref(err_free))
            self.assertEqual(FAT_OK, st_free)
            self.assertEqual(0, err_free.value)
        finally:
            self.fat_BytesFree(tar_bytes_handle)


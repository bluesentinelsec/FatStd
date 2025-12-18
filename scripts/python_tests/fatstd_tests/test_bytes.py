from __future__ import annotations

import ctypes
import unittest

from fatstd_test_support import bind, fat_string_handle_type


class TestBytes(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        fat_bytes = fat_string_handle_type()
        fat_bytes_array = fat_string_handle_type()
        fat_string = fat_string_handle_type()

        cls.fat_BytesNewN = bind(
            "fat_BytesNewN", argtypes=[ctypes.c_void_p, ctypes.c_size_t], restype=fat_bytes
        )
        cls.fat_BytesLen = bind("fat_BytesLen", argtypes=[fat_bytes], restype=ctypes.c_size_t)
        cls.fat_BytesCopyOut = bind(
            "fat_BytesCopyOut", argtypes=[fat_bytes, ctypes.c_void_p, ctypes.c_size_t], restype=ctypes.c_size_t
        )
        cls.fat_BytesClone = bind("fat_BytesClone", argtypes=[fat_bytes], restype=fat_bytes)
        cls.fat_BytesContains = bind(
            "fat_BytesContains", argtypes=[fat_bytes, fat_bytes], restype=ctypes.c_bool
        )
        cls.fat_BytesHasPrefix = bind(
            "fat_BytesHasPrefix", argtypes=[fat_bytes, fat_bytes], restype=ctypes.c_bool
        )
        cls.fat_BytesHasSuffix = bind(
            "fat_BytesHasSuffix", argtypes=[fat_bytes, fat_bytes], restype=ctypes.c_bool
        )
        cls.fat_BytesTrimSpace = bind("fat_BytesTrimSpace", argtypes=[fat_bytes], restype=fat_bytes)
        cls.fat_StringNewUTF8 = bind(
            "fat_StringNewUTF8", argtypes=[ctypes.c_char_p], restype=fat_string
        )
        cls.fat_StringFree = bind("fat_StringFree", argtypes=[fat_string], restype=None)
        cls.fat_BytesTrim = bind(
            "fat_BytesTrim", argtypes=[fat_bytes, fat_string], restype=fat_bytes
        )
        cls.fat_BytesSplit = bind(
            "fat_BytesSplit", argtypes=[fat_bytes, fat_bytes], restype=fat_bytes_array
        )
        cls.fat_BytesArrayLen = bind(
            "fat_BytesArrayLen", argtypes=[fat_bytes_array], restype=ctypes.c_size_t
        )
        cls.fat_BytesArrayGet = bind(
            "fat_BytesArrayGet", argtypes=[fat_bytes_array, ctypes.c_size_t], restype=fat_bytes
        )
        cls.fat_BytesArrayFree = bind(
            "fat_BytesArrayFree", argtypes=[fat_bytes_array], restype=None
        )
        cls.fat_BytesJoin = bind(
            "fat_BytesJoin", argtypes=[fat_bytes_array, fat_bytes], restype=fat_bytes
        )
        cls.fat_BytesReplaceAll = bind(
            "fat_BytesReplaceAll", argtypes=[fat_bytes, fat_bytes, fat_bytes], restype=fat_bytes
        )
        cls.fat_BytesReplace = bind(
            "fat_BytesReplace",
            argtypes=[fat_bytes, fat_bytes, fat_bytes, ctypes.c_int],
            restype=fat_bytes,
        )
        cls.fat_BytesToLower = bind("fat_BytesToLower", argtypes=[fat_bytes], restype=fat_bytes)
        cls.fat_BytesToUpper = bind("fat_BytesToUpper", argtypes=[fat_bytes], restype=fat_bytes)
        cls.fat_BytesIndex = bind(
            "fat_BytesIndex", argtypes=[fat_bytes, fat_bytes], restype=ctypes.c_int
        )
        cls.fat_BytesCount = bind(
            "fat_BytesCount", argtypes=[fat_bytes, fat_bytes], restype=ctypes.c_int
        )
        cls.fat_BytesCompare = bind(
            "fat_BytesCompare", argtypes=[fat_bytes, fat_bytes], restype=ctypes.c_int
        )
        cls.fat_BytesEqual = bind(
            "fat_BytesEqual", argtypes=[fat_bytes, fat_bytes], restype=ctypes.c_bool
        )
        cls.fat_BytesFree = bind("fat_BytesFree", argtypes=[fat_bytes], restype=None)

    def _bytes_new(self, b: bytes) -> int:
        raw = ctypes.create_string_buffer(b, len(b))
        h = self.fat_BytesNewN(ctypes.addressof(raw), len(raw.raw))
        self.assertNotEqual(0, h)
        return h

    def _bytes_to_py(self, h: int) -> bytes:
        n = self.fat_BytesLen(h)
        dst = ctypes.create_string_buffer(n, n)
        copied = self.fat_BytesCopyOut(h, ctypes.addressof(dst), len(dst.raw))
        self.assertEqual(n, copied)
        return dst.raw

    def test_basic_ops(self) -> None:
        b = self._bytes_new(b"abc\x00def")
        self.assertEqual(7, self.fat_BytesLen(b))
        self.assertEqual(b"abc\x00def", self._bytes_to_py(b))

        cloned = self.fat_BytesClone(b)
        self.assertNotEqual(0, cloned)
        self.assertNotEqual(b, cloned)
        self.assertTrue(self.fat_BytesEqual(b, cloned))

        sub = self._bytes_new(b"\x00d")
        self.assertTrue(self.fat_BytesContains(b, sub))

        prefix = self._bytes_new(b"abc\x00")
        suffix = self._bytes_new(b"def")
        self.assertTrue(self.fat_BytesHasPrefix(b, prefix))
        self.assertTrue(self.fat_BytesHasSuffix(b, suffix))

        self.fat_BytesFree(suffix)
        self.fat_BytesFree(prefix)
        self.fat_BytesFree(sub)
        self.fat_BytesFree(cloned)
        self.fat_BytesFree(b)

    def test_trim_trim_space(self) -> None:
        b = self._bytes_new(b"  abc  ")
        trimmed_space = self.fat_BytesTrimSpace(b)
        self.assertEqual(b"abc", self._bytes_to_py(trimmed_space))

        cutset = self.fat_StringNewUTF8(b" a")
        self.assertNotEqual(0, cutset)
        trimmed = self.fat_BytesTrim(b, cutset)
        self.assertEqual(b"bc", self._bytes_to_py(trimmed))

        self.fat_BytesFree(trimmed)
        self.fat_StringFree(cutset)
        self.fat_BytesFree(trimmed_space)
        self.fat_BytesFree(b)

    def test_split_join_replace_case_index(self) -> None:
        s = self._bytes_new(b"a,b,c")
        sep = self._bytes_new(b",")

        arr = self.fat_BytesSplit(s, sep)
        self.assertEqual(3, self.fat_BytesArrayLen(arr))

        e0 = self.fat_BytesArrayGet(arr, 0)
        e1 = self.fat_BytesArrayGet(arr, 1)
        e2 = self.fat_BytesArrayGet(arr, 2)
        self.assertEqual(b"a", self._bytes_to_py(e0))
        self.assertEqual(b"b", self._bytes_to_py(e1))
        self.assertEqual(b"c", self._bytes_to_py(e2))

        joined = self.fat_BytesJoin(arr, sep)
        self.assertEqual(b"a,b,c", self._bytes_to_py(joined))

        old = self._bytes_new(b"b")
        new = self._bytes_new(b"BB")
        replaced_all = self.fat_BytesReplaceAll(joined, old, new)
        self.assertEqual(b"a,BB,c", self._bytes_to_py(replaced_all))

        replaced_1 = self.fat_BytesReplace(joined, old, new, 1)
        self.assertEqual(b"a,BB,c", self._bytes_to_py(replaced_1))

        mixed = self._bytes_new(b"AbC")
        lower = self.fat_BytesToLower(mixed)
        upper = self.fat_BytesToUpper(mixed)
        self.fat_BytesFree(mixed)
        self.assertEqual(b"abc", self._bytes_to_py(lower))
        self.assertEqual(b"ABC", self._bytes_to_py(upper))

        self.assertEqual(2, self.fat_BytesIndex(joined, old))
        self.assertEqual(1, self.fat_BytesCount(joined, old))
        a = self._bytes_new(b"a")
        self.assertLess(self.fat_BytesCompare(a, old), 0)
        self.fat_BytesFree(a)

        old2 = self._bytes_new(b"b")
        self.assertTrue(self.fat_BytesEqual(old, old2))
        self.fat_BytesFree(old2)

        self.fat_BytesFree(upper)
        self.fat_BytesFree(lower)
        self.fat_BytesFree(replaced_1)
        self.fat_BytesFree(replaced_all)
        self.fat_BytesFree(new)
        self.fat_BytesFree(old)
        self.fat_BytesFree(joined)
        self.fat_BytesFree(e2)
        self.fat_BytesFree(e1)
        self.fat_BytesFree(e0)
        self.fat_BytesArrayFree(arr)
        self.fat_BytesFree(sep)
        self.fat_BytesFree(s)

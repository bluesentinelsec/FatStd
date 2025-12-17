from __future__ import annotations

import ctypes
import unittest

from fatstd_test_support import bind, fat_string_handle_type


class TestString(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        fat_string = fat_string_handle_type()
        fat_string_array = fat_string_handle_type()

        cls.fat_StringNewUTF8 = bind(
            "fat_StringNewUTF8", argtypes=[ctypes.c_char_p], restype=fat_string
        )
        cls.fat_StringNewUTF8N = bind(
            "fat_StringNewUTF8N", argtypes=[ctypes.c_void_p, ctypes.c_size_t], restype=fat_string
        )
        cls.fat_StringClone = bind("fat_StringClone", argtypes=[fat_string], restype=fat_string)
        cls.fat_StringContains = bind(
            "fat_StringContains", argtypes=[fat_string, fat_string], restype=ctypes.c_bool
        )
        cls.fat_StringHasPrefix = bind(
            "fat_StringHasPrefix", argtypes=[fat_string, fat_string], restype=ctypes.c_bool
        )
        cls.fat_StringHasSuffix = bind(
            "fat_StringHasSuffix", argtypes=[fat_string, fat_string], restype=ctypes.c_bool
        )
        cls.fat_StringTrimSpace = bind("fat_StringTrimSpace", argtypes=[fat_string], restype=fat_string)
        cls.fat_StringTrim = bind(
            "fat_StringTrim", argtypes=[fat_string, fat_string], restype=fat_string
        )
        cls.fat_StringSplit = bind(
            "fat_StringSplit", argtypes=[fat_string, fat_string], restype=fat_string_array
        )
        cls.fat_StringSplitN = bind(
            "fat_StringSplitN", argtypes=[fat_string, fat_string, ctypes.c_int], restype=fat_string_array
        )
        cls.fat_StringArrayLen = bind(
            "fat_StringArrayLen", argtypes=[fat_string_array], restype=ctypes.c_size_t
        )
        cls.fat_StringArrayGet = bind(
            "fat_StringArrayGet",
            argtypes=[fat_string_array, ctypes.c_size_t],
            restype=fat_string,
        )
        cls.fat_StringArrayFree = bind(
            "fat_StringArrayFree", argtypes=[fat_string_array], restype=None
        )
        cls.fat_StringJoin = bind(
            "fat_StringJoin", argtypes=[fat_string_array, fat_string], restype=fat_string
        )
        cls.fat_StringReplace = bind(
            "fat_StringReplace",
            argtypes=[fat_string, fat_string, fat_string, ctypes.c_int],
            restype=fat_string,
        )
        cls.fat_StringReplaceAll = bind(
            "fat_StringReplaceAll", argtypes=[fat_string, fat_string, fat_string], restype=fat_string
        )
        cls.fat_StringFree = bind("fat_StringFree", argtypes=[fat_string], restype=None)

    def _assert_string_equal(self, a, b, *, message: str | None = None) -> None:
        self.assertTrue(self.fat_StringContains(a, b), message or "expected a to contain b")
        self.assertTrue(self.fat_StringContains(b, a), message or "expected b to contain a")

    def test_create_clone_free_cstr(self) -> None:
        s1 = self.fat_StringNewUTF8(b"lorem ipsum")
        self.assertNotEqual(0, s1, "fat_StringNewUTF8 returned 0 handle")

        c1 = self.fat_StringClone(s1)
        self.assertNotEqual(0, c1, "fat_StringClone returned 0 handle")
        self.assertNotEqual(s1, c1, "fat_StringClone returned same handle")

        self.fat_StringFree(s1)
        self.fat_StringFree(c1)

    def test_create_clone_free_bytes(self) -> None:
        raw = ctypes.create_string_buffer(b"abc\x00def")
        s2 = self.fat_StringNewUTF8N(ctypes.addressof(raw), len(raw.raw))
        self.assertNotEqual(0, s2, "fat_StringNewUTF8N returned 0 handle")

        c2 = self.fat_StringClone(s2)
        self.assertNotEqual(0, c2, "fat_StringClone returned 0 handle")
        self.assertNotEqual(s2, c2, "fat_StringClone returned same handle")

        self.fat_StringFree(s2)
        self.fat_StringFree(c2)

    def test_contains_basic(self) -> None:
        hay = self.fat_StringNewUTF8(b"lorem ipsum")
        self.assertNotEqual(0, hay)

        needle_yes = self.fat_StringNewUTF8(b"ipsum")
        self.assertNotEqual(0, needle_yes)

        needle_no = self.fat_StringNewUTF8(b"IPSUM")
        self.assertNotEqual(0, needle_no)

        self.assertTrue(self.fat_StringContains(hay, needle_yes))
        self.assertFalse(self.fat_StringContains(hay, needle_no))

        self.fat_StringFree(needle_no)
        self.fat_StringFree(needle_yes)
        self.fat_StringFree(hay)

    def test_contains_embedded_nul(self) -> None:
        hay_bytes = b"abc\x00def"
        hay_raw = ctypes.create_string_buffer(hay_bytes, len(hay_bytes))
        hay = self.fat_StringNewUTF8N(ctypes.addressof(hay_raw), len(hay_raw.raw))
        self.assertNotEqual(0, hay)

        needle_bytes = b"\x00de"
        needle_raw = ctypes.create_string_buffer(needle_bytes, len(needle_bytes))
        needle = self.fat_StringNewUTF8N(ctypes.addressof(needle_raw), len(needle_raw.raw))
        self.assertNotEqual(0, needle)

        self.assertTrue(self.fat_StringContains(hay, needle))

        self.fat_StringFree(needle)
        self.fat_StringFree(hay)

    def test_has_prefix_basic(self) -> None:
        s = self.fat_StringNewUTF8(b"lorem ipsum")
        self.assertNotEqual(0, s)

        prefix_yes = self.fat_StringNewUTF8(b"lorem")
        self.assertNotEqual(0, prefix_yes)

        prefix_no = self.fat_StringNewUTF8(b"ipsum")
        self.assertNotEqual(0, prefix_no)

        self.assertTrue(self.fat_StringHasPrefix(s, prefix_yes))
        self.assertFalse(self.fat_StringHasPrefix(s, prefix_no))

        self.fat_StringFree(prefix_no)
        self.fat_StringFree(prefix_yes)
        self.fat_StringFree(s)

    def test_has_suffix_basic(self) -> None:
        s = self.fat_StringNewUTF8(b"lorem ipsum")
        self.assertNotEqual(0, s)

        suffix_yes = self.fat_StringNewUTF8(b"ipsum")
        self.assertNotEqual(0, suffix_yes)

        suffix_no = self.fat_StringNewUTF8(b"lorem")
        self.assertNotEqual(0, suffix_no)

        self.assertTrue(self.fat_StringHasSuffix(s, suffix_yes))
        self.assertFalse(self.fat_StringHasSuffix(s, suffix_no))

        self.fat_StringFree(suffix_no)
        self.fat_StringFree(suffix_yes)
        self.fat_StringFree(s)

    def test_has_prefix_suffix_embedded_nul(self) -> None:
        s_bytes = b"abc\x00def"
        s_raw = ctypes.create_string_buffer(s_bytes, len(s_bytes))
        s = self.fat_StringNewUTF8N(ctypes.addressof(s_raw), len(s_raw.raw))
        self.assertNotEqual(0, s)

        prefix_bytes = b"abc\x00"
        prefix_raw = ctypes.create_string_buffer(prefix_bytes, len(prefix_bytes))
        prefix = self.fat_StringNewUTF8N(ctypes.addressof(prefix_raw), len(prefix_raw.raw))
        self.assertNotEqual(0, prefix)

        suffix_bytes = b"\x00def"
        suffix_raw = ctypes.create_string_buffer(suffix_bytes, len(suffix_bytes))
        suffix = self.fat_StringNewUTF8N(ctypes.addressof(suffix_raw), len(suffix_raw.raw))
        self.assertNotEqual(0, suffix)

        self.assertTrue(self.fat_StringHasPrefix(s, prefix))
        self.assertTrue(self.fat_StringHasSuffix(s, suffix))

        self.fat_StringFree(suffix)
        self.fat_StringFree(prefix)
        self.fat_StringFree(s)

    def test_trim_space(self) -> None:
        s = self.fat_StringNewUTF8(b"  lorem ipsum  ")
        self.assertNotEqual(0, s)

        trimmed = self.fat_StringTrimSpace(s)
        self.assertNotEqual(0, trimmed)
        self.assertNotEqual(s, trimmed)

        space = self.fat_StringNewUTF8(b" ")
        self.assertNotEqual(0, space)

        self.assertFalse(self.fat_StringHasPrefix(trimmed, space))
        self.assertFalse(self.fat_StringHasSuffix(trimmed, space))

        lorem = self.fat_StringNewUTF8(b"lorem")
        self.assertNotEqual(0, lorem)
        self.assertTrue(self.fat_StringHasPrefix(trimmed, lorem))

        self.fat_StringFree(lorem)
        self.fat_StringFree(space)
        self.fat_StringFree(trimmed)
        self.fat_StringFree(s)

    def test_trim_cutset_basic(self) -> None:
        s = self.fat_StringNewUTF8(b"...lorem ipsum...")
        self.assertNotEqual(0, s)

        cutset = self.fat_StringNewUTF8(b".")
        self.assertNotEqual(0, cutset)

        trimmed = self.fat_StringTrim(s, cutset)
        self.assertNotEqual(0, trimmed)
        self.assertNotEqual(s, trimmed)

        dot = self.fat_StringNewUTF8(b".")
        self.assertNotEqual(0, dot)

        self.assertFalse(self.fat_StringHasPrefix(trimmed, dot))
        self.assertFalse(self.fat_StringHasSuffix(trimmed, dot))

        needle = self.fat_StringNewUTF8(b"lorem")
        self.assertNotEqual(0, needle)
        self.assertTrue(self.fat_StringContains(trimmed, needle))

        self.fat_StringFree(needle)
        self.fat_StringFree(dot)
        self.fat_StringFree(trimmed)
        self.fat_StringFree(cutset)
        self.fat_StringFree(s)

    def test_trim_cutset_embedded_nul(self) -> None:
        s_bytes = b"\x00abc\x00"
        s_raw = ctypes.create_string_buffer(s_bytes, len(s_bytes))
        s = self.fat_StringNewUTF8N(ctypes.addressof(s_raw), len(s_raw.raw))
        self.assertNotEqual(0, s)

        cutset_bytes = b"\x00"
        cutset_raw = ctypes.create_string_buffer(cutset_bytes, len(cutset_bytes))
        cutset = self.fat_StringNewUTF8N(ctypes.addressof(cutset_raw), len(cutset_raw.raw))
        self.assertNotEqual(0, cutset)

        trimmed = self.fat_StringTrim(s, cutset)
        self.assertNotEqual(0, trimmed)

        nul_bytes = b"\x00"
        nul_raw = ctypes.create_string_buffer(nul_bytes, len(nul_bytes))
        nul = self.fat_StringNewUTF8N(ctypes.addressof(nul_raw), len(nul_raw.raw))
        self.assertNotEqual(0, nul)

        self.assertFalse(self.fat_StringHasPrefix(trimmed, nul))
        self.assertFalse(self.fat_StringHasSuffix(trimmed, nul))

        needle = self.fat_StringNewUTF8(b"abc")
        self.assertNotEqual(0, needle)
        self.assertTrue(self.fat_StringContains(trimmed, needle))

        self.fat_StringFree(needle)
        self.fat_StringFree(nul)
        self.fat_StringFree(trimmed)
        self.fat_StringFree(cutset)
        self.fat_StringFree(s)

    def test_split_join_basic(self) -> None:
        s = self.fat_StringNewUTF8(b"a,b,c")
        self.assertNotEqual(0, s)

        sep = self.fat_StringNewUTF8(b",")
        self.assertNotEqual(0, sep)

        arr = self.fat_StringSplit(s, sep)
        self.assertNotEqual(0, arr)
        self.assertEqual(3, self.fat_StringArrayLen(arr))

        expected0 = self.fat_StringNewUTF8(b"a")
        expected1 = self.fat_StringNewUTF8(b"b")
        expected2 = self.fat_StringNewUTF8(b"c")
        self.assertNotEqual(0, expected0)
        self.assertNotEqual(0, expected1)
        self.assertNotEqual(0, expected2)

        e0 = self.fat_StringArrayGet(arr, 0)
        e1 = self.fat_StringArrayGet(arr, 1)
        e2 = self.fat_StringArrayGet(arr, 2)
        self.assertNotEqual(0, e0)
        self.assertNotEqual(0, e1)
        self.assertNotEqual(0, e2)

        self._assert_string_equal(e0, expected0)
        self._assert_string_equal(e1, expected1)
        self._assert_string_equal(e2, expected2)

        joined = self.fat_StringJoin(arr, sep)
        self.assertNotEqual(0, joined)
        self._assert_string_equal(joined, s)

        self.fat_StringFree(joined)
        self.fat_StringFree(e2)
        self.fat_StringFree(e1)
        self.fat_StringFree(e0)
        self.fat_StringFree(expected2)
        self.fat_StringFree(expected1)
        self.fat_StringFree(expected0)
        self.fat_StringArrayFree(arr)
        self.fat_StringFree(sep)
        self.fat_StringFree(s)

    def test_split_n_basic(self) -> None:
        s = self.fat_StringNewUTF8(b"a,b,c")
        self.assertNotEqual(0, s)

        sep = self.fat_StringNewUTF8(b",")
        self.assertNotEqual(0, sep)

        arr = self.fat_StringSplitN(s, sep, 2)
        self.assertNotEqual(0, arr)
        self.assertEqual(2, self.fat_StringArrayLen(arr))

        expected0 = self.fat_StringNewUTF8(b"a")
        expected1 = self.fat_StringNewUTF8(b"b,c")
        self.assertNotEqual(0, expected0)
        self.assertNotEqual(0, expected1)

        e0 = self.fat_StringArrayGet(arr, 0)
        e1 = self.fat_StringArrayGet(arr, 1)
        self.assertNotEqual(0, e0)
        self.assertNotEqual(0, e1)

        self._assert_string_equal(e0, expected0)
        self._assert_string_equal(e1, expected1)

        self.fat_StringFree(e1)
        self.fat_StringFree(e0)
        self.fat_StringFree(expected1)
        self.fat_StringFree(expected0)
        self.fat_StringArrayFree(arr)
        self.fat_StringFree(sep)
        self.fat_StringFree(s)

    def test_replace_and_replace_all(self) -> None:
        s = self.fat_StringNewUTF8(b"foo bar foo")
        self.assertNotEqual(0, s)

        old = self.fat_StringNewUTF8(b"foo")
        new = self.fat_StringNewUTF8(b"baz")
        self.assertNotEqual(0, old)
        self.assertNotEqual(0, new)

        replaced_1 = self.fat_StringReplace(s, old, new, 1)
        self.assertNotEqual(0, replaced_1)
        expected_1 = self.fat_StringNewUTF8(b"baz bar foo")
        self.assertNotEqual(0, expected_1)
        self._assert_string_equal(replaced_1, expected_1)

        replaced_all = self.fat_StringReplaceAll(s, old, new)
        self.assertNotEqual(0, replaced_all)
        expected_all = self.fat_StringNewUTF8(b"baz bar baz")
        self.assertNotEqual(0, expected_all)
        self._assert_string_equal(replaced_all, expected_all)

        self.fat_StringFree(expected_all)
        self.fat_StringFree(replaced_all)
        self.fat_StringFree(expected_1)
        self.fat_StringFree(replaced_1)
        self.fat_StringFree(new)
        self.fat_StringFree(old)
        self.fat_StringFree(s)

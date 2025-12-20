from __future__ import annotations

import ctypes
import os
import tempfile
import unittest

from fatstd_test_support import bind, fat_string_handle_type


FAT_OK = 0
FAT_ERR_RANGE = 2


class TestTiled(unittest.TestCase):
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

        cls.fat_StringArrayLen = bind("fat_StringArrayLen", argtypes=[fat_string_array], restype=ctypes.c_size_t)
        cls.fat_StringArrayGet = bind("fat_StringArrayGet", argtypes=[fat_string_array, ctypes.c_size_t], restype=fat_string)
        cls.fat_StringArrayFree = bind("fat_StringArrayFree", argtypes=[fat_string_array], restype=None)

        cls.fat_BytesNewN = bind("fat_BytesNewN", argtypes=[ctypes.c_void_p, ctypes.c_size_t], restype=fat_bytes)
        cls.fat_BytesFree = bind("fat_BytesFree", argtypes=[fat_bytes], restype=None)

        cls.fat_ErrorMessage = bind("fat_ErrorMessage", argtypes=[fat_error], restype=fat_string)
        cls.fat_ErrorFree = bind("fat_ErrorFree", argtypes=[fat_error], restype=None)

        cls.fat_TiledMapLoadFileUTF8 = bind(
            "fat_TiledMapLoadFileUTF8", argtypes=[ctypes.c_char_p, ctypes.c_void_p, ctypes.c_void_p], restype=ctypes.c_int
        )
        cls.fat_TiledMapLoadReaderBytesUTF8 = bind(
            "fat_TiledMapLoadReaderBytesUTF8",
            argtypes=[ctypes.c_char_p, fat_bytes, ctypes.c_void_p, ctypes.c_void_p],
            restype=ctypes.c_int,
        )
        cls.fat_TiledMapFree = bind("fat_TiledMapFree", argtypes=[fat_handle], restype=None)

        cls.fat_TiledMapWidth = bind("fat_TiledMapWidth", argtypes=[fat_handle], restype=ctypes.c_int)
        cls.fat_TiledMapHeight = bind("fat_TiledMapHeight", argtypes=[fat_handle], restype=ctypes.c_int)
        cls.fat_TiledMapTileWidth = bind("fat_TiledMapTileWidth", argtypes=[fat_handle], restype=ctypes.c_int)
        cls.fat_TiledMapTileHeight = bind("fat_TiledMapTileHeight", argtypes=[fat_handle], restype=ctypes.c_int)
        cls.fat_TiledMapOrientation = bind("fat_TiledMapOrientation", argtypes=[fat_handle], restype=fat_string)
        cls.fat_TiledMapGetFileFullPathUTF8 = bind(
            "fat_TiledMapGetFileFullPathUTF8", argtypes=[fat_handle, ctypes.c_char_p], restype=fat_string
        )
        cls.fat_TiledMapProperties = bind("fat_TiledMapProperties", argtypes=[fat_handle], restype=fat_handle)
        cls.fat_TiledMapLayerCount = bind("fat_TiledMapLayerCount", argtypes=[fat_handle], restype=ctypes.c_size_t)
        cls.fat_TiledMapLayerAt = bind("fat_TiledMapLayerAt", argtypes=[fat_handle, ctypes.c_size_t], restype=fat_handle)
        cls.fat_TiledMapTileGIDToTile = bind(
            "fat_TiledMapTileGIDToTile",
            argtypes=[fat_handle, ctypes.c_uint32, ctypes.c_void_p, ctypes.c_void_p],
            restype=ctypes.c_int,
        )

        cls.fat_TiledLayerFree = bind("fat_TiledLayerFree", argtypes=[fat_handle], restype=None)
        cls.fat_TiledLayerName = bind("fat_TiledLayerName", argtypes=[fat_handle], restype=fat_string)
        cls.fat_TiledLayerIsEmpty = bind("fat_TiledLayerIsEmpty", argtypes=[fat_handle], restype=ctypes.c_bool)
        cls.fat_TiledLayerProperties = bind("fat_TiledLayerProperties", argtypes=[fat_handle], restype=fat_handle)
        cls.fat_TiledLayerTileAt = bind("fat_TiledLayerTileAt", argtypes=[fat_handle, ctypes.c_int, ctypes.c_int], restype=fat_handle)

        cls.fat_TiledLayerTileFree = bind("fat_TiledLayerTileFree", argtypes=[fat_handle], restype=None)
        cls.fat_TiledLayerTileIsNil = bind("fat_TiledLayerTileIsNil", argtypes=[fat_handle], restype=ctypes.c_bool)
        cls.fat_TiledLayerTileID = bind("fat_TiledLayerTileID", argtypes=[fat_handle], restype=ctypes.c_uint32)
        cls.fat_TiledLayerTileTilesetName = bind("fat_TiledLayerTileTilesetName", argtypes=[fat_handle], restype=fat_string)
        cls.fat_TiledLayerTileRect = bind(
            "fat_TiledLayerTileRect",
            argtypes=[fat_handle, ctypes.c_void_p, ctypes.c_void_p, ctypes.c_void_p, ctypes.c_void_p],
            restype=None,
        )

        cls.fat_TiledPropertiesGet = bind("fat_TiledPropertiesGet", argtypes=[fat_handle, fat_string], restype=fat_string_array)
        cls.fat_TiledPropertiesGetString = bind(
            "fat_TiledPropertiesGetString", argtypes=[fat_handle, fat_string], restype=fat_string
        )
        cls.fat_TiledPropertiesFree = bind("fat_TiledPropertiesFree", argtypes=[fat_handle], restype=None)

    def _bytes_new(self, b: bytes) -> int:
        raw = ctypes.create_string_buffer(b, len(b))
        h = self.fat_BytesNewN(ctypes.addressof(raw), len(raw.raw))
        self.assertNotEqual(0, h)
        return h

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

    def test_load_map_extract_tile_and_properties(self) -> None:
        tmx = b"""<?xml version="1.0" encoding="UTF-8"?>
<map version="1.10" tiledversion="1.10.2" orientation="orthogonal" renderorder="right-down"
     width="2" height="2" tilewidth="16" tileheight="16" infinite="0">
  <properties>
    <property name="difficulty" value="hard"/>
    <property name="tag" value="a"/>
    <property name="tag" value="b"/>
  </properties>
  <tileset firstgid="10" name="ts" tilewidth="16" tileheight="16" tilecount="2" columns="2">
    <image source="tiles.png" width="32" height="16"/>
  </tileset>
  <layer id="1" name="Ground" width="2" height="2">
    <properties>
      <property name="layer_kind" value="main"/>
    </properties>
    <data encoding="csv">10,0,0,11</data>
  </layer>
</map>
"""
        with tempfile.TemporaryDirectory() as td:
            path = os.path.join(td, "map.tmx")
            with open(path, "wb") as f:
                f.write(tmx)

            m = fat_string_handle_type()(0)
            err = fat_string_handle_type()(0)
            st = int(self.fat_TiledMapLoadFileUTF8(path.encode("utf-8"), ctypes.byref(m), ctypes.byref(err)))
            self.assertEqual(FAT_OK, st)
            self.assertEqual(0, err.value)
            self.assertNotEqual(0, m.value)

            try:
                self.assertEqual(2, int(self.fat_TiledMapWidth(m.value)))
                self.assertEqual(2, int(self.fat_TiledMapHeight(m.value)))
                self.assertEqual(16, int(self.fat_TiledMapTileWidth(m.value)))
                self.assertEqual(16, int(self.fat_TiledMapTileHeight(m.value)))

                orient_s = self.fat_TiledMapOrientation(m.value)
                try:
                    self.assertEqual("orthogonal", self._string_to_py(orient_s))
                finally:
                    self.fat_StringFree(orient_s)

                full_s = self.fat_TiledMapGetFileFullPathUTF8(m.value, b"tiles.png")
                try:
                    full = self._string_to_py(full_s)
                finally:
                    self.fat_StringFree(full_s)
                self.assertTrue(full.endswith(os.path.sep + "tiles.png") or full.endswith("tiles.png"))
                self.assertIn(td, full)

                props = self.fat_TiledMapProperties(m.value)
                try:
                    key = self.fat_StringNewUTF8(b"difficulty")
                    try:
                        val_s = self.fat_TiledPropertiesGetString(props, key)
                        try:
                            self.assertEqual("hard", self._string_to_py(val_s))
                        finally:
                            self.fat_StringFree(val_s)
                    finally:
                        self.fat_StringFree(key)

                    key2 = self.fat_StringNewUTF8(b"tag")
                    try:
                        arr = self.fat_TiledPropertiesGet(props, key2)
                        try:
                            self.assertEqual(["a", "b"], self._string_array_to_py(arr))
                        finally:
                            self.fat_StringArrayFree(arr)
                    finally:
                        self.fat_StringFree(key2)
                finally:
                    self.fat_TiledPropertiesFree(props)

                layer_count = int(self.fat_TiledMapLayerCount(m.value))
                self.assertEqual(1, layer_count)

                layer = self.fat_TiledMapLayerAt(m.value, 0)
                self.assertNotEqual(0, layer)
                try:
                    name_s = self.fat_TiledLayerName(layer)
                    try:
                        self.assertEqual("Ground", self._string_to_py(name_s))
                    finally:
                        self.fat_StringFree(name_s)
                    self.assertFalse(bool(self.fat_TiledLayerIsEmpty(layer)))

                    lprops = self.fat_TiledLayerProperties(layer)
                    try:
                        key = self.fat_StringNewUTF8(b"layer_kind")
                        try:
                            val_s = self.fat_TiledPropertiesGetString(lprops, key)
                            try:
                                self.assertEqual("main", self._string_to_py(val_s))
                            finally:
                                self.fat_StringFree(val_s)
                        finally:
                            self.fat_StringFree(key)
                    finally:
                        self.fat_TiledPropertiesFree(lprops)

                    t00 = self.fat_TiledLayerTileAt(layer, 0, 0)
                    try:
                        self.assertFalse(bool(self.fat_TiledLayerTileIsNil(t00)))
                        self.assertEqual(0, int(self.fat_TiledLayerTileID(t00)))
                        ts_s = self.fat_TiledLayerTileTilesetName(t00)
                        try:
                            self.assertEqual("ts", self._string_to_py(ts_s))
                        finally:
                            self.fat_StringFree(ts_s)
                        x = ctypes.c_int(0)
                        y = ctypes.c_int(0)
                        w = ctypes.c_int(0)
                        h = ctypes.c_int(0)
                        self.fat_TiledLayerTileRect(t00, ctypes.byref(x), ctypes.byref(y), ctypes.byref(w), ctypes.byref(h))
                        self.assertEqual((0, 0, 16, 16), (x.value, y.value, w.value, h.value))
                    finally:
                        self.fat_TiledLayerTileFree(t00)

                    t11 = self.fat_TiledLayerTileAt(layer, 1, 1)
                    try:
                        self.assertFalse(bool(self.fat_TiledLayerTileIsNil(t11)))
                        self.assertEqual(1, int(self.fat_TiledLayerTileID(t11)))
                        x = ctypes.c_int(0)
                        y = ctypes.c_int(0)
                        w = ctypes.c_int(0)
                        h = ctypes.c_int(0)
                        self.fat_TiledLayerTileRect(t11, ctypes.byref(x), ctypes.byref(y), ctypes.byref(w), ctypes.byref(h))
                        self.assertEqual((16, 0, 16, 16), (x.value, y.value, w.value, h.value))
                    finally:
                        self.fat_TiledLayerTileFree(t11)

                    bad_tile = fat_string_handle_type()(0)
                    err2 = fat_string_handle_type()(0)
                    st2 = int(
                        self.fat_TiledMapTileGIDToTile(m.value, 2, ctypes.byref(bad_tile), ctypes.byref(err2))
                    )
                    self.assertEqual(FAT_ERR_RANGE, st2)
                    self.assertEqual(0, bad_tile.value)
                    self.assertNotEqual(0, err2.value)
                    msg = self._error_to_py(err2.value).lower()
                    self.assertTrue("gid" in msg or "invalid" in msg)
                    self.fat_ErrorFree(err2.value)
                finally:
                    self.fat_TiledLayerFree(layer)
            finally:
                self.fat_TiledMapFree(m.value)

    def test_load_map_from_bytes_reader(self) -> None:
        tmx = b"""<map version="1.10" tiledversion="1.10.2" orientation="orthogonal" renderorder="right-down"
     width="1" height="1" tilewidth="8" tileheight="8" infinite="0">
  <tileset firstgid="1" name="ts" tilewidth="8" tileheight="8" tilecount="1" columns="1">
    <image source="tiles.png" width="8" height="8"/>
  </tileset>
  <layer id="1" name="L" width="1" height="1">
    <data encoding="csv">1</data>
  </layer>
</map>
"""
        with tempfile.TemporaryDirectory() as td:
            b = self._bytes_new(tmx)
            try:
                m = fat_string_handle_type()(0)
                err = fat_string_handle_type()(0)
                st = int(self.fat_TiledMapLoadReaderBytesUTF8(td.encode("utf-8"), b, ctypes.byref(m), ctypes.byref(err)))
                self.assertEqual(FAT_OK, st)
                self.assertEqual(0, err.value)
                self.assertNotEqual(0, m.value)
                try:
                    self.assertEqual(1, int(self.fat_TiledMapLayerCount(m.value)))
                finally:
                    self.fat_TiledMapFree(m.value)
            finally:
                self.fat_BytesFree(b)

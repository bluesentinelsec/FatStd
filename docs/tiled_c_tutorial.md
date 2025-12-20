<!--
  Tiled C tutorial for FatStd.
  Keep this aligned with docs/design.md and docs/error_strategy.md.
-->

# Tiled (TMX) in C (FatStd tutorial)

FatStd’s Tiled APIs live in `include/fat/tiled.h` and are backed by `github.com/lafriks/go-tiled`.

The C API is a pragmatic subset designed for callers that want to load a TMX map, read basic metadata, and extract tile / property values.

## Error handling pattern

```c
#include <stdio.h>
#include <stdlib.h>

#include "fat/error.h"
#include "fat/string.h"

static void die_status(const char *op, fat_Status st, fat_Error err) {
  if (err) {
    fat_String msg = fat_ErrorMessage(err);
    char buf[512];
    size_t n = fat_StringCopyOutCStr(msg, buf, sizeof(buf) - 1);
    buf[n] = '\0';
    fat_StringFree(msg);
    fat_ErrorFree(err);
    fprintf(stderr, "%s failed: status=%d err=%s\n", op, (int)st, buf);
  } else {
    fprintf(stderr, "%s failed: status=%d\n", op, (int)st);
  }
  abort();
}
```

## 1) Load a map from a TMX file and extract values

This example shows:

- reading the map
- accessing map metadata
- reading a map property
- reading a tile from a layer at `(x,y)` and getting its tileset name + rectangle

```c
#include <stdint.h>
#include <stdio.h>

#include "fat/error.h"
#include "fat/status.h"
#include "fat/string.h"
#include "fat/tiled.h"

static void print_fat_string(const char *label, fat_String s) {
  char buf[512];
  size_t n = fat_StringCopyOutCStr(s, buf, sizeof(buf) - 1);
  buf[n] = '\0';
  printf("%s%s\n", label, buf);
}

void read_map_example(const char *tmx_path_utf8) {
  fat_TiledMap m = 0;
  fat_Error err = 0;
  fat_Status st = fat_TiledMapLoadFileUTF8(tmx_path_utf8, &m, &err);
  if (st != FAT_OK) die_status("fat_TiledMapLoadFileUTF8", st, err);

  printf("map: %dx%d tiles, tile=%dx%d px\n",
         fat_TiledMapWidth(m),
         fat_TiledMapHeight(m),
         fat_TiledMapTileWidth(m),
         fat_TiledMapTileHeight(m));

  fat_String orient = fat_TiledMapOrientation(m);
  print_fat_string("orientation=", orient);
  fat_StringFree(orient);

  /* Read a custom map property */
  fat_TiledProperties props = fat_TiledMapProperties(m);
  fat_String key = fat_StringNewUTF8("difficulty");
  fat_String value = fat_TiledPropertiesGetString(props, key);
  print_fat_string("difficulty=", value);
  fat_StringFree(value);
  fat_StringFree(key);
  fat_TiledPropertiesFree(props);

  /* Read a tile from the first tile layer */
  size_t n_layers = fat_TiledMapLayerCount(m);
  if (n_layers == 0) {
    fprintf(stderr, "no tile layers in map\n");
    fat_TiledMapFree(m);
    return;
  }

  fat_TiledLayer layer = fat_TiledMapLayerAt(m, 0);
  fat_String layer_name = fat_TiledLayerName(layer);
  print_fat_string("layer=", layer_name);
  fat_StringFree(layer_name);

  fat_TiledLayerTile t = fat_TiledLayerTileAt(layer, 0, 0);
  if (!fat_TiledLayerTileIsNil(t)) {
    uint32_t id = fat_TiledLayerTileID(t);
    fat_String ts_name = fat_TiledLayerTileTilesetName(t);
    int rx = 0, ry = 0, rw = 0, rh = 0;
    fat_TiledLayerTileRect(t, &rx, &ry, &rw, &rh);

    printf("tile id=%u rect=(%d,%d,%d,%d)\n", (unsigned)id, rx, ry, rw, rh);
    print_fat_string("tileset=", ts_name);

    fat_StringFree(ts_name);
  }
  fat_TiledLayerTileFree(t);

  fat_TiledLayerFree(layer);
  fat_TiledMapFree(m);
}
```

## 2) Load a map from bytes (no `io.Reader` in C)

If you already have the TMX XML bytes in memory, use `fat_TiledMapLoadReaderBytesUTF8(baseDir, tmxBytes, ...)`:

```c
#include "fat/bytes.h"

void load_from_bytes_example(const char *base_dir_utf8, const void *tmx_bytes, size_t tmx_len) {
  fat_Bytes b = fat_BytesNewN(tmx_bytes, tmx_len);

  fat_TiledMap m = 0;
  fat_Error err = 0;
  fat_Status st = fat_TiledMapLoadReaderBytesUTF8(base_dir_utf8, b, &m, &err);
  if (st != FAT_OK) die_status("fat_TiledMapLoadReaderBytesUTF8", st, err);

  fat_TiledMapFree(m);
  fat_BytesFree(b);
}
```


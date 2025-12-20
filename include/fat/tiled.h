#pragma once

/**
 * @file fat/tiled.h
 * @brief Tiled (TMX/TSX) map loading and inspection utilities.
 *
 * Backed by github.com/lafriks/go-tiled.
 *
 * Design notes:
 * - The upstream Go API is reflection/XML-decoder based and exposes deep object graphs.
 * - The FatStd C API is handle-based and provides a pragmatic, inspectable subset:
 *   - load a map from file path or in-memory bytes
 *   - enumerate tile layers
 *   - inspect tiles at (x,y) and resolve a GID to a tileset + tile rectangle
 *   - read custom properties
 *
 * Recoverable failures return fat_Status and optionally fat_Error.
 * Contract violations (invalid handles, NULL out-params where forbidden, out-of-range indices)
 * are fatal by default.
 */

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#include "fat/bytes.h"
#include "fat/error.h"
#include "fat/export.h"
#include "fat/handle.h"
#include "fat/status.h"
#include "fat/string.h"

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Opaque handle to a loaded Tiled map.
 *
 * @note Ownership: free with fat_TiledMapFree.
 */
typedef fat_Handle fat_TiledMap;

/**
 * @brief Opaque handle to a map tile layer.
 *
 * @note Ownership: free with fat_TiledLayerFree.
 */
typedef fat_Handle fat_TiledLayer;

/**
 * @brief Opaque handle to a resolved tile (layer tile).
 *
 * This corresponds to go-tiled's `LayerTile` and includes tileset metadata.
 *
 * @note Ownership: free with fat_TiledLayerTileFree.
 */
typedef fat_Handle fat_TiledLayerTile;

/**
 * @brief Opaque handle to a set of custom properties.
 *
 * @note Ownership: free with fat_TiledPropertiesFree.
 */
typedef fat_Handle fat_TiledProperties;

/**
 * @brief Loads a TMX map from a filesystem path (UTF-8).
 *
 * @param path UTF-8 path (NUL-terminated).
 * @param out_map Output: new map handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_SYNTAX on parse errors; FAT_ERR_OTHER on I/O failures.
 */
FATSTD_API fat_Status fat_TiledMapLoadFileUTF8(const char *path, fat_TiledMap *out_map, fat_Error *out_err);

/**
 * @brief Loads a TMX map from an in-memory XML blob.
 *
 * @param base_dir Base directory used to resolve external references (UTF-8, NUL-terminated).
 * @param tmx XML bytes.
 * @param out_map Output: new map handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_SYNTAX on parse errors; FAT_ERR_OTHER on I/O failures.
 *
 * @note This mirrors go-tiled's LoadReader(baseDir, r, ...).
 */
FATSTD_API fat_Status fat_TiledMapLoadReaderBytesUTF8(
  const char *base_dir,
  fat_Bytes tmx,
  fat_TiledMap *out_map,
  fat_Error *out_err
);

/**
 * @brief Frees a map handle.
 *
 * @param m Map handle.
 */
FATSTD_API void fat_TiledMapFree(fat_TiledMap m);

/**
 * @brief Returns map width in tiles.
 */
FATSTD_API int fat_TiledMapWidth(fat_TiledMap m);

/**
 * @brief Returns map height in tiles.
 */
FATSTD_API int fat_TiledMapHeight(fat_TiledMap m);

/**
 * @brief Returns tile width in pixels.
 */
FATSTD_API int fat_TiledMapTileWidth(fat_TiledMap m);

/**
 * @brief Returns tile height in pixels.
 */
FATSTD_API int fat_TiledMapTileHeight(fat_TiledMap m);

/**
 * @brief Returns the map orientation string (e.g. "orthogonal").
 *
 * @return New fat_String handle (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_TiledMapOrientation(fat_TiledMap m);

/**
 * @brief Returns a path joined relative to the map base directory.
 *
 * Matches go-tiled's (*Map).GetFileFullPath.
 *
 * @param m Map handle.
 * @param file_name Relative file name (UTF-8, NUL-terminated).
 * @return New fat_String handle (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_TiledMapGetFileFullPathUTF8(fat_TiledMap m, const char *file_name);

/**
 * @brief Returns custom properties associated with the map.
 *
 * @return New properties handle (caller must fat_TiledPropertiesFree).
 */
FATSTD_API fat_TiledProperties fat_TiledMapProperties(fat_TiledMap m);

/**
 * @brief Returns the number of tile layers in the map.
 *
 * This counts only `<layer>` elements (tile layers), not object groups, image layers, or groups.
 */
FATSTD_API size_t fat_TiledMapLayerCount(fat_TiledMap m);

/**
 * @brief Returns the tile layer at index.
 *
 * @param m Map handle.
 * @param idx Layer index (0..count-1).
 * @return New layer handle (caller must fat_TiledLayerFree).
 */
FATSTD_API fat_TiledLayer fat_TiledMapLayerAt(fat_TiledMap m, size_t idx);

/**
 * @brief Frees a layer handle.
 */
FATSTD_API void fat_TiledLayerFree(fat_TiledLayer layer);

/**
 * @brief Returns the layer name.
 *
 * @return New fat_String handle (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_TiledLayerName(fat_TiledLayer layer);

/**
 * @brief Returns whether all tiles in the layer are empty (GID == 0).
 */
FATSTD_API bool fat_TiledLayerIsEmpty(fat_TiledLayer layer);

/**
 * @brief Returns custom properties associated with the layer.
 *
 * @return New properties handle (caller must fat_TiledPropertiesFree).
 */
FATSTD_API fat_TiledProperties fat_TiledLayerProperties(fat_TiledLayer layer);

/**
 * @brief Returns the resolved tile at (x,y) on the layer.
 *
 * @param layer Layer handle.
 * @param x Tile X coordinate (0 <= x < map width).
 * @param y Tile Y coordinate (0 <= y < map height).
 * @return New tile handle (caller must fat_TiledLayerTileFree).
 *
 * @note Out-of-range coordinates are contract violations (fatal).
 */
FATSTD_API fat_TiledLayerTile fat_TiledLayerTileAt(fat_TiledLayer layer, int x, int y);

/**
 * @brief Resolves a raw GID to a tile.
 *
 * Matches go-tiled's (*Map).TileGIDToTile.
 *
 * @param m Map handle.
 * @param gid Raw tile GID (may include flip bits).
 * @param out_tile Output: new tile handle on success.
 * @param out_err Output: error handle on failure, 0 on success.
 * @return FAT_OK on success; FAT_ERR_RANGE if the GID is invalid; FAT_ERR_OTHER on other failures.
 */
FATSTD_API fat_Status fat_TiledMapTileGIDToTile(
  fat_TiledMap m,
  uint32_t gid,
  fat_TiledLayerTile *out_tile,
  fat_Error *out_err
);

/**
 * @brief Frees a tile handle.
 */
FATSTD_API void fat_TiledLayerTileFree(fat_TiledLayerTile tile);

/**
 * @brief Reports whether a tile handle represents a nil/empty tile.
 */
FATSTD_API bool fat_TiledLayerTileIsNil(fat_TiledLayerTile tile);

/**
 * @brief Returns the tile ID within its tileset (0-based).
 */
FATSTD_API uint32_t fat_TiledLayerTileID(fat_TiledLayerTile tile);

/**
 * @brief Returns the tileset name for a tile.
 *
 * @return New fat_String handle (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_TiledLayerTileTilesetName(fat_TiledLayerTile tile);

/**
 * @brief Returns the tile rectangle within the tileset image.
 *
 * @param tile Tile handle.
 * @param out_x Output: rect min X.
 * @param out_y Output: rect min Y.
 * @param out_w Output: rect width.
 * @param out_h Output: rect height.
 */
FATSTD_API void fat_TiledLayerTileRect(fat_TiledLayerTile tile, int *out_x, int *out_y, int *out_w, int *out_h);

/**
 * @brief Returns all property values for a name.
 *
 * Matches go-tiled's (Properties).Get.
 *
 * @param props Properties handle.
 * @param name Property name string handle.
 * @return New string array handle (caller must fat_StringArrayFree).
 */
FATSTD_API fat_StringArray fat_TiledPropertiesGet(fat_TiledProperties props, fat_String name);

/**
 * @brief Returns a string property value.
 *
 * Matches go-tiled's (Properties).GetString.
 *
 * @param props Properties handle.
 * @param name Property name string handle.
 * @return New fat_String handle (caller must fat_StringFree).
 */
FATSTD_API fat_String fat_TiledPropertiesGetString(fat_TiledProperties props, fat_String name);

/**
 * @brief Returns an integer property value (or 0 if missing/unparseable).
 *
 * Matches go-tiled's (Properties).GetInt.
 */
FATSTD_API int fat_TiledPropertiesGetInt(fat_TiledProperties props, fat_String name);

/**
 * @brief Returns a float property value (or 0 if missing/unparseable).
 *
 * Matches go-tiled's (Properties).GetFloat.
 */
FATSTD_API double fat_TiledPropertiesGetFloat(fat_TiledProperties props, fat_String name);

/**
 * @brief Returns a boolean property value.
 *
 * Matches go-tiled's (Properties).GetBool.
 */
FATSTD_API bool fat_TiledPropertiesGetBool(fat_TiledProperties props, fat_String name);

/**
 * @brief Frees a properties handle.
 */
FATSTD_API void fat_TiledPropertiesFree(fat_TiledProperties props);

#ifdef __cplusplus
} /* extern "C" */
#endif


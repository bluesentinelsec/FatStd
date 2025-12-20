#include "fat/tiled.h"

#include "fatstd_go.h"

fat_Status fat_TiledMapLoadFileUTF8(const char *path, fat_TiledMap *out_map, fat_Error *out_err) {
  return (fat_Status)fatstd_go_tiled_map_load_file_utf8((char *)path, (uintptr_t *)out_map, (uintptr_t *)out_err);
}

fat_Status fat_TiledMapLoadReaderBytesUTF8(const char *base_dir, fat_Bytes tmx, fat_TiledMap *out_map, fat_Error *out_err) {
  return (fat_Status)fatstd_go_tiled_map_load_reader_bytes_utf8(
    (char *)base_dir, (uintptr_t)tmx, (uintptr_t *)out_map, (uintptr_t *)out_err
  );
}

void fat_TiledMapFree(fat_TiledMap m) {
  fatstd_go_tiled_map_free((uintptr_t)m);
}

int fat_TiledMapWidth(fat_TiledMap m) {
  return (int)fatstd_go_tiled_map_width((uintptr_t)m);
}

int fat_TiledMapHeight(fat_TiledMap m) {
  return (int)fatstd_go_tiled_map_height((uintptr_t)m);
}

int fat_TiledMapTileWidth(fat_TiledMap m) {
  return (int)fatstd_go_tiled_map_tile_width((uintptr_t)m);
}

int fat_TiledMapTileHeight(fat_TiledMap m) {
  return (int)fatstd_go_tiled_map_tile_height((uintptr_t)m);
}

fat_String fat_TiledMapOrientation(fat_TiledMap m) {
  return (fat_String)fatstd_go_tiled_map_orientation((uintptr_t)m);
}

fat_String fat_TiledMapGetFileFullPathUTF8(fat_TiledMap m, const char *file_name) {
  return (fat_String)fatstd_go_tiled_map_get_file_full_path_utf8((uintptr_t)m, (char *)file_name);
}

fat_TiledProperties fat_TiledMapProperties(fat_TiledMap m) {
  return (fat_TiledProperties)fatstd_go_tiled_map_properties((uintptr_t)m);
}

size_t fat_TiledMapLayerCount(fat_TiledMap m) {
  return (size_t)fatstd_go_tiled_map_layer_count((uintptr_t)m);
}

fat_TiledLayer fat_TiledMapLayerAt(fat_TiledMap m, size_t idx) {
  return (fat_TiledLayer)fatstd_go_tiled_map_layer_at((uintptr_t)m, idx);
}

void fat_TiledLayerFree(fat_TiledLayer layer) {
  fatstd_go_tiled_layer_free((uintptr_t)layer);
}

fat_String fat_TiledLayerName(fat_TiledLayer layer) {
  return (fat_String)fatstd_go_tiled_layer_name((uintptr_t)layer);
}

bool fat_TiledLayerIsEmpty(fat_TiledLayer layer) {
  return fatstd_go_tiled_layer_is_empty((uintptr_t)layer) != 0;
}

fat_TiledProperties fat_TiledLayerProperties(fat_TiledLayer layer) {
  return (fat_TiledProperties)fatstd_go_tiled_layer_properties((uintptr_t)layer);
}

fat_TiledLayerTile fat_TiledLayerTileAt(fat_TiledLayer layer, int x, int y) {
  return (fat_TiledLayerTile)fatstd_go_tiled_layer_tile_at((uintptr_t)layer, x, y);
}

fat_Status fat_TiledMapTileGIDToTile(fat_TiledMap m, uint32_t gid, fat_TiledLayerTile *out_tile, fat_Error *out_err) {
  return (fat_Status)fatstd_go_tiled_map_tile_gid_to_tile((uintptr_t)m, gid, (uintptr_t *)out_tile, (uintptr_t *)out_err);
}

void fat_TiledLayerTileFree(fat_TiledLayerTile tile) {
  fatstd_go_tiled_layer_tile_free((uintptr_t)tile);
}

bool fat_TiledLayerTileIsNil(fat_TiledLayerTile tile) {
  return fatstd_go_tiled_layer_tile_is_nil((uintptr_t)tile) != 0;
}

uint32_t fat_TiledLayerTileID(fat_TiledLayerTile tile) {
  return (uint32_t)fatstd_go_tiled_layer_tile_id((uintptr_t)tile);
}

fat_String fat_TiledLayerTileTilesetName(fat_TiledLayerTile tile) {
  return (fat_String)fatstd_go_tiled_layer_tile_tileset_name((uintptr_t)tile);
}

void fat_TiledLayerTileRect(fat_TiledLayerTile tile, int *out_x, int *out_y, int *out_w, int *out_h) {
  fatstd_go_tiled_layer_tile_rect((uintptr_t)tile, out_x, out_y, out_w, out_h);
}

fat_StringArray fat_TiledPropertiesGet(fat_TiledProperties props, fat_String name) {
  return (fat_StringArray)fatstd_go_tiled_properties_get((uintptr_t)props, (uintptr_t)name);
}

fat_String fat_TiledPropertiesGetString(fat_TiledProperties props, fat_String name) {
  return (fat_String)fatstd_go_tiled_properties_get_string((uintptr_t)props, (uintptr_t)name);
}

int fat_TiledPropertiesGetInt(fat_TiledProperties props, fat_String name) {
  return (int)fatstd_go_tiled_properties_get_int((uintptr_t)props, (uintptr_t)name);
}

double fat_TiledPropertiesGetFloat(fat_TiledProperties props, fat_String name) {
  return (double)fatstd_go_tiled_properties_get_float((uintptr_t)props, (uintptr_t)name);
}

bool fat_TiledPropertiesGetBool(fat_TiledProperties props, fat_String name) {
  return fatstd_go_tiled_properties_get_bool((uintptr_t)props, (uintptr_t)name) != 0;
}

void fat_TiledPropertiesFree(fat_TiledProperties props) {
  fatstd_go_tiled_properties_free((uintptr_t)props);
}


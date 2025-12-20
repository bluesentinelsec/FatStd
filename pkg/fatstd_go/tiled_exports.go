package main

/*
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
*/
import "C"

import (
	"bytes"
	"errors"
	"io"

	gotiled "github.com/lafriks/go-tiled"
)

const (
	fatTiledErrCodeSyntax = 220
	fatTiledErrCodeIO     = 221
	fatTiledErrCodeRange  = 222
	fatTiledErrCodeOther  = 223
)

type fatTiledMap struct {
	m *gotiled.Map
}

type fatTiledLayer struct {
	m *gotiled.Map
	l *gotiled.Layer
}

type fatTiledLayerTile struct {
	t *gotiled.LayerTile
}

type fatTiledProperties struct {
	p gotiled.Properties
}

func fatstdTiledStatusFromError(err error) C.int {
	if err == nil {
		return fatStatusOK
	}
	if errors.Is(err, io.EOF) {
		return fatStatusEOF
	}
	if errors.Is(err, gotiled.ErrInvalidTileGID) {
		return fatStatusRange
	}
	// Decode/unmarshal errors and invalid encodings are effectively syntax/data errors.
	switch {
	case errors.Is(err, gotiled.ErrUnknownEncoding),
		errors.Is(err, gotiled.ErrEmptyLayerData),
		errors.Is(err, gotiled.ErrInvalidDecodedTileCount):
		return fatStatusSyntax
	default:
		return fatStatusOther
	}
}

func fatstdTiledErrCodeFromError(err error) int32 {
	if err == nil {
		return fatTiledErrCodeOther
	}
	if errors.Is(err, gotiled.ErrInvalidTileGID) {
		return fatTiledErrCodeRange
	}
	switch {
	case errors.Is(err, gotiled.ErrUnknownEncoding),
		errors.Is(err, gotiled.ErrEmptyLayerData),
		errors.Is(err, gotiled.ErrInvalidDecodedTileCount):
		return fatTiledErrCodeSyntax
	default:
		return fatTiledErrCodeOther
	}
}

func fatstdTiledMapFromHandle(handle uintptr) *fatTiledMap {
	if handle == 0 {
		panic("fatstdTiledMapFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdTiledMapFromHandle: invalid handle")
	}
	m, ok := value.(*fatTiledMap)
	if !ok {
		panic("fatstdTiledMapFromHandle: handle is not tiled map")
	}
	if m.m == nil {
		panic("fatstdTiledMapFromHandle: map is nil")
	}
	return m
}

func fatstdTiledLayerFromHandle(handle uintptr) *fatTiledLayer {
	if handle == 0 {
		panic("fatstdTiledLayerFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdTiledLayerFromHandle: invalid handle")
	}
	l, ok := value.(*fatTiledLayer)
	if !ok {
		panic("fatstdTiledLayerFromHandle: handle is not tiled layer")
	}
	if l.m == nil || l.l == nil {
		panic("fatstdTiledLayerFromHandle: layer is nil")
	}
	return l
}

func fatstdTiledLayerTileFromHandle(handle uintptr) *fatTiledLayerTile {
	if handle == 0 {
		panic("fatstdTiledLayerTileFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdTiledLayerTileFromHandle: invalid handle")
	}
	t, ok := value.(*fatTiledLayerTile)
	if !ok {
		panic("fatstdTiledLayerTileFromHandle: handle is not tiled layer tile")
	}
	if t.t == nil {
		panic("fatstdTiledLayerTileFromHandle: tile is nil")
	}
	return t
}

func fatstdTiledPropertiesFromHandle(handle uintptr) *fatTiledProperties {
	if handle == 0 {
		panic("fatstdTiledPropertiesFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdTiledPropertiesFromHandle: invalid handle")
	}
	p, ok := value.(*fatTiledProperties)
	if !ok {
		panic("fatstdTiledPropertiesFromHandle: handle is not tiled properties")
	}
	return p
}

//export fatstd_go_tiled_map_load_file_utf8
func fatstd_go_tiled_map_load_file_utf8(path *C.char, outMap *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if path == nil {
		panic("fatstd_go_tiled_map_load_file_utf8: path is NULL")
	}
	if outMap == nil {
		panic("fatstd_go_tiled_map_load_file_utf8: outMap is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_tiled_map_load_file_utf8: outErr is NULL")
	}

	m, err := gotiled.LoadFile(C.GoString(path))
	if err != nil {
		*outMap = 0
		*outErr = C.uintptr_t(fatstdNewError(fatstdTiledErrCodeFromError(err), err.Error()))
		return fatstdTiledStatusFromError(err)
	}
	*outMap = C.uintptr_t(fatstdHandles.register(&fatTiledMap{m: m}))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_tiled_map_load_reader_bytes_utf8
func fatstd_go_tiled_map_load_reader_bytes_utf8(baseDir *C.char, tmxHandle C.uintptr_t, outMap *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if baseDir == nil {
		panic("fatstd_go_tiled_map_load_reader_bytes_utf8: baseDir is NULL")
	}
	if outMap == nil {
		panic("fatstd_go_tiled_map_load_reader_bytes_utf8: outMap is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_tiled_map_load_reader_bytes_utf8: outErr is NULL")
	}

	b := fatstdBytesFromHandle(uintptr(tmxHandle))
	r := bytes.NewReader(b.Value())
	m, err := gotiled.LoadReader(C.GoString(baseDir), r)
	if err != nil {
		*outMap = 0
		*outErr = C.uintptr_t(fatstdNewError(fatstdTiledErrCodeFromError(err), err.Error()))
		return fatstdTiledStatusFromError(err)
	}
	*outMap = C.uintptr_t(fatstdHandles.register(&fatTiledMap{m: m}))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_tiled_map_free
func fatstd_go_tiled_map_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_tiled_map_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_tiled_map_free: invalid handle")
	}
	if _, ok := value.(*fatTiledMap); !ok {
		panic("fatstd_go_tiled_map_free: handle is not tiled map")
	}
}

//export fatstd_go_tiled_map_width
func fatstd_go_tiled_map_width(handle C.uintptr_t) C.int {
	m := fatstdTiledMapFromHandle(uintptr(handle))
	return C.int(m.m.Width)
}

//export fatstd_go_tiled_map_height
func fatstd_go_tiled_map_height(handle C.uintptr_t) C.int {
	m := fatstdTiledMapFromHandle(uintptr(handle))
	return C.int(m.m.Height)
}

//export fatstd_go_tiled_map_tile_width
func fatstd_go_tiled_map_tile_width(handle C.uintptr_t) C.int {
	m := fatstdTiledMapFromHandle(uintptr(handle))
	return C.int(m.m.TileWidth)
}

//export fatstd_go_tiled_map_tile_height
func fatstd_go_tiled_map_tile_height(handle C.uintptr_t) C.int {
	m := fatstdTiledMapFromHandle(uintptr(handle))
	return C.int(m.m.TileHeight)
}

//export fatstd_go_tiled_map_orientation
func fatstd_go_tiled_map_orientation(handle C.uintptr_t) C.uintptr_t {
	m := fatstdTiledMapFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdStringNewFromGoString(m.m.Orientation))
}

//export fatstd_go_tiled_map_get_file_full_path_utf8
func fatstd_go_tiled_map_get_file_full_path_utf8(handle C.uintptr_t, fileName *C.char) C.uintptr_t {
	if fileName == nil {
		panic("fatstd_go_tiled_map_get_file_full_path_utf8: fileName is NULL")
	}
	m := fatstdTiledMapFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdStringNewFromGoString(m.m.GetFileFullPath(C.GoString(fileName))))
}

//export fatstd_go_tiled_map_properties
func fatstd_go_tiled_map_properties(handle C.uintptr_t) C.uintptr_t {
	m := fatstdTiledMapFromHandle(uintptr(handle))
	if m.m.Properties == nil {
		return C.uintptr_t(fatstdHandles.register(&fatTiledProperties{p: gotiled.Properties{}}))
	}
	return C.uintptr_t(fatstdHandles.register(&fatTiledProperties{p: *m.m.Properties}))
}

//export fatstd_go_tiled_properties_free
func fatstd_go_tiled_properties_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_tiled_properties_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_tiled_properties_free: invalid handle")
	}
	if _, ok := value.(*fatTiledProperties); !ok {
		panic("fatstd_go_tiled_properties_free: handle is not tiled properties")
	}
}

//export fatstd_go_tiled_properties_get
func fatstd_go_tiled_properties_get(handle C.uintptr_t, nameHandle C.uintptr_t) C.uintptr_t {
	p := fatstdTiledPropertiesFromHandle(uintptr(handle))
	name := fatstdStringFromHandle(uintptr(nameHandle))
	return C.uintptr_t(fatstdStringArrayNew(p.p.Get(name.Value())))
}

//export fatstd_go_tiled_properties_get_string
func fatstd_go_tiled_properties_get_string(handle C.uintptr_t, nameHandle C.uintptr_t) C.uintptr_t {
	p := fatstdTiledPropertiesFromHandle(uintptr(handle))
	name := fatstdStringFromHandle(uintptr(nameHandle))
	return C.uintptr_t(fatstdStringNewFromGoString(p.p.GetString(name.Value())))
}

//export fatstd_go_tiled_properties_get_int
func fatstd_go_tiled_properties_get_int(handle C.uintptr_t, nameHandle C.uintptr_t) C.int {
	p := fatstdTiledPropertiesFromHandle(uintptr(handle))
	name := fatstdStringFromHandle(uintptr(nameHandle))
	return C.int(p.p.GetInt(name.Value()))
}

//export fatstd_go_tiled_properties_get_float
func fatstd_go_tiled_properties_get_float(handle C.uintptr_t, nameHandle C.uintptr_t) C.double {
	p := fatstdTiledPropertiesFromHandle(uintptr(handle))
	name := fatstdStringFromHandle(uintptr(nameHandle))
	return C.double(p.p.GetFloat(name.Value()))
}

//export fatstd_go_tiled_properties_get_bool
func fatstd_go_tiled_properties_get_bool(handle C.uintptr_t, nameHandle C.uintptr_t) C.int {
	p := fatstdTiledPropertiesFromHandle(uintptr(handle))
	name := fatstdStringFromHandle(uintptr(nameHandle))
	if p.p.GetBool(name.Value()) {
		return 1
	}
	return 0
}

//export fatstd_go_tiled_map_layer_count
func fatstd_go_tiled_map_layer_count(handle C.uintptr_t) C.size_t {
	m := fatstdTiledMapFromHandle(uintptr(handle))
	return C.size_t(len(m.m.Layers))
}

//export fatstd_go_tiled_map_layer_at
func fatstd_go_tiled_map_layer_at(handle C.uintptr_t, idx C.size_t) C.uintptr_t {
	m := fatstdTiledMapFromHandle(uintptr(handle))
	if idx > C.size_t(2147483647) {
		panic("fatstd_go_tiled_map_layer_at: idx too large")
	}
	i := int(idx)
	if i < 0 || i >= len(m.m.Layers) {
		panic("fatstd_go_tiled_map_layer_at: idx out of range")
	}
	layer := m.m.Layers[i]
	return C.uintptr_t(fatstdHandles.register(&fatTiledLayer{m: m.m, l: layer}))
}

//export fatstd_go_tiled_layer_free
func fatstd_go_tiled_layer_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_tiled_layer_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_tiled_layer_free: invalid handle")
	}
	if _, ok := value.(*fatTiledLayer); !ok {
		panic("fatstd_go_tiled_layer_free: handle is not tiled layer")
	}
}

//export fatstd_go_tiled_layer_name
func fatstd_go_tiled_layer_name(handle C.uintptr_t) C.uintptr_t {
	l := fatstdTiledLayerFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdStringNewFromGoString(l.l.Name))
}

//export fatstd_go_tiled_layer_is_empty
func fatstd_go_tiled_layer_is_empty(handle C.uintptr_t) C.int {
	l := fatstdTiledLayerFromHandle(uintptr(handle))
	if l.l.IsEmpty() {
		return 1
	}
	return 0
}

//export fatstd_go_tiled_layer_properties
func fatstd_go_tiled_layer_properties(handle C.uintptr_t) C.uintptr_t {
	l := fatstdTiledLayerFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdHandles.register(&fatTiledProperties{p: l.l.Properties}))
}

//export fatstd_go_tiled_layer_tile_at
func fatstd_go_tiled_layer_tile_at(handle C.uintptr_t, x C.int, y C.int) C.uintptr_t {
	l := fatstdTiledLayerFromHandle(uintptr(handle))
	if x < 0 || y < 0 {
		panic("fatstd_go_tiled_layer_tile_at: negative coordinates")
	}
	if int(x) >= l.m.Width || int(y) >= l.m.Height {
		panic("fatstd_go_tiled_layer_tile_at: coordinates out of range")
	}
	tile := l.l.Tiles[int(y)*l.m.Width+int(x)]
	return C.uintptr_t(fatstdHandles.register(&fatTiledLayerTile{t: tile}))
}

//export fatstd_go_tiled_map_tile_gid_to_tile
func fatstd_go_tiled_map_tile_gid_to_tile(handle C.uintptr_t, gid C.uint32_t, outTile *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outTile == nil {
		panic("fatstd_go_tiled_map_tile_gid_to_tile: outTile is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_tiled_map_tile_gid_to_tile: outErr is NULL")
	}
	m := fatstdTiledMapFromHandle(uintptr(handle))
	tile, err := m.m.TileGIDToTile(uint32(gid))
	if err != nil {
		*outTile = 0
		*outErr = C.uintptr_t(fatstdNewError(fatstdTiledErrCodeFromError(err), err.Error()))
		return fatstdTiledStatusFromError(err)
	}
	*outTile = C.uintptr_t(fatstdHandles.register(&fatTiledLayerTile{t: tile}))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_tiled_layer_tile_free
func fatstd_go_tiled_layer_tile_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_tiled_layer_tile_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_tiled_layer_tile_free: invalid handle")
	}
	if _, ok := value.(*fatTiledLayerTile); !ok {
		panic("fatstd_go_tiled_layer_tile_free: handle is not tiled layer tile")
	}
}

//export fatstd_go_tiled_layer_tile_is_nil
func fatstd_go_tiled_layer_tile_is_nil(handle C.uintptr_t) C.int {
	t := fatstdTiledLayerTileFromHandle(uintptr(handle))
	if t.t.IsNil() {
		return 1
	}
	return 0
}

//export fatstd_go_tiled_layer_tile_id
func fatstd_go_tiled_layer_tile_id(handle C.uintptr_t) C.uint32_t {
	t := fatstdTiledLayerTileFromHandle(uintptr(handle))
	return C.uint32_t(t.t.ID)
}

//export fatstd_go_tiled_layer_tile_tileset_name
func fatstd_go_tiled_layer_tile_tileset_name(handle C.uintptr_t) C.uintptr_t {
	t := fatstdTiledLayerTileFromHandle(uintptr(handle))
	if t.t.Tileset == nil {
		return C.uintptr_t(fatstdStringNewFromGoString(""))
	}
	return C.uintptr_t(fatstdStringNewFromGoString(t.t.Tileset.Name))
}

//export fatstd_go_tiled_layer_tile_rect
func fatstd_go_tiled_layer_tile_rect(handle C.uintptr_t, outX *C.int, outY *C.int, outW *C.int, outH *C.int) {
	if outX == nil || outY == nil || outW == nil || outH == nil {
		panic("fatstd_go_tiled_layer_tile_rect: out params are NULL")
	}
	t := fatstdTiledLayerTileFromHandle(uintptr(handle))
	r := t.t.GetTileRect()
	*outX = C.int(r.Min.X)
	*outY = C.int(r.Min.Y)
	*outW = C.int(r.Dx())
	*outH = C.int(r.Dy())
}


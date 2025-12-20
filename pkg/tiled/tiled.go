package tiled

import (
	"encoding/xml"
	"image"
	"io"
	"io/fs"

	ext "github.com/lafriks/go-tiled"
)

// Types (aliases) - expose go-tiled API under FatStd's module tree.
type AnimationFrame = ext.AnimationFrame
type Axis = ext.Axis
type Data = ext.Data
type DataTile = ext.DataTile
type Ellipse = ext.Ellipse
type Group = ext.Group
type HexColor = ext.HexColor
type Image = ext.Image
type ImageLayer = ext.ImageLayer
type Layer = ext.Layer
type LayerTile = ext.LayerTile
type LoaderOption = ext.LoaderOption
type Map = ext.Map
type Object = ext.Object
type ObjectGroup = ext.ObjectGroup
type Point = ext.Point
type Points = ext.Points
type PolyLine = ext.PolyLine
type Polygon = ext.Polygon
type Properties = ext.Properties
type Property = ext.Property
type StaggerIndexType = ext.StaggerIndexType
type Template = ext.Template
type Terrain = ext.Terrain
type Text = ext.Text
type Tileset = ext.Tileset
type TilesetTile = ext.TilesetTile
type TilesetTileOffset = ext.TilesetTileOffset
type WangColor = ext.WangColor
type WangPosition = ext.WangPosition
type WangSet = ext.WangSet
type WangSets = ext.WangSets
type WangTile = ext.WangTile

// Methods listed in the prompt are provided by the type aliases above:
// - (*Group).DecodeGroup, (*Group).UnmarshalXML
// - NewHexColor, ParseHexColor, (*HexColor).MarshalXMLAttr/RGBA/String/UnmarshalXMLAttr
// - (*ImageLayer).UnmarshalXML
// - (*Layer).DecodeLayer/GetTilePosition/IsEmpty/UnmarshalXML
// - (*LayerTile).GetTileRect/IsNil
// - WithFileSystem
// - LoadFile, LoadReader, (*Map).GetFileFullPath/TileGIDToTile/UnmarshalXML
// - (*Object).UnmarshalXML
// - (*ObjectGroup).DecodeObjectGroup/UnmarshalXML
// - (*Points).UnmarshalXMLAttr
// - (Properties).Get/GetBool/GetColor/GetFloat/GetInt/GetString
// - (*Property).UnmarshalXML
// - (*Text).UnmarshalXML
// - LoadTilesetFile, LoadTilesetReader, (*Tileset).BaseDir/GetFileFullPath/GetTileRect/GetTilesetTile/SetBaseDir
// - (*WangSet).GetWangColors

func NewHexColor(r, g, b, a uint32) HexColor { return ext.NewHexColor(r, g, b, a) }
func ParseHexColor(s string) (HexColor, error) { return ext.ParseHexColor(s) }

func WithFileSystem(fileSystem fs.FS) LoaderOption { return ext.WithFileSystem(fileSystem) }

func LoadFile(fileName string, options ...LoaderOption) (*Map, error) { return ext.LoadFile(fileName, options...) }
func LoadReader(baseDir string, r io.Reader, options ...LoaderOption) (*Map, error) {
	return ext.LoadReader(baseDir, r, options...)
}

func LoadTilesetFile(fileName string, options ...LoaderOption) (*Tileset, error) {
	return ext.LoadTilesetFile(fileName, options...)
}
func LoadTilesetReader(baseDir string, r io.Reader, options ...LoaderOption) (*Tileset, error) {
	return ext.LoadTilesetReader(baseDir, r, options...)
}

// Ensure these imports stay used when the compiler checks signatures.
var (
	_ = xml.Attr{}
	_ = image.Rectangle{}
)


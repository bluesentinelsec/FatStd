Refer to the following documents for context on the project:
1. docs/design.md
2. docs/onboarding_and_testing_functions.md
3. docs/documenting_code.md
4. docs/error_strategy.md

Implement the following functions in fatstd:

github.com/lafriks/go-tiled
https://pkg.go.dev/github.com/lafriks/go-tiled#section-readme

type AnimationFrame
type Axis
type Data
type DataTile
type Ellipse
type Group
func (g *Group) DecodeGroup(m *Map) error
func (g *Group) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error
type HexColor
func NewHexColor(r, g, b, a uint32) HexColor
func ParseHexColor(s string) (HexColor, error)
func (color *HexColor) MarshalXMLAttr(name xml.Name) (attr xml.Attr, err error)
func (color *HexColor) RGBA() (r, g, b, a uint32)
func (color *HexColor) String() string
func (color *HexColor) UnmarshalXMLAttr(attr xml.Attr) error
type Image
type ImageLayer
func (l *ImageLayer) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error
type Layer
func (l *Layer) DecodeLayer(m *Map) error
func (l *Layer) GetTilePosition(tileID int) (int, int)
func (l *Layer) IsEmpty() bool
func (l *Layer) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error
type LayerTile
func (t *LayerTile) GetTileRect() image.Rectangle
func (t *LayerTile) IsNil() bool
type LoaderOption
func WithFileSystem(fileSystem fs.FS) LoaderOption
type Map
func LoadFile(fileName string, options ...LoaderOption) (*Map, error)
func LoadReader(baseDir string, r io.Reader, options ...LoaderOption) (*Map, error)
func (m *Map) GetFileFullPath(fileName string) string
func (m *Map) TileGIDToTile(gid uint32) (*LayerTile, error)
func (m *Map) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error
type Object
func (o *Object) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error
type ObjectGroup
func (g *ObjectGroup) DecodeObjectGroup(m *Map) error
func (g *ObjectGroup) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error
type Point
type Points
func (m *Points) UnmarshalXMLAttr(attr xml.Attr) error
type PolyLine
type Polygon
type Properties
func (p Properties) Get(name string) []string
func (p Properties) GetBool(name string) bool
func (p Properties) GetColor(name string) color.Color
func (p Properties) GetFloat(name string) float64
func (p Properties) GetInt(name string) int
func (p Properties) GetString(name string) string
type Property
func (p *Property) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error
type StaggerIndexType
type Template
type Terrain
type Text
func (t *Text) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error
type Tileset
func LoadTilesetFile(fileName string, options ...LoaderOption) (*Tileset, error)
func LoadTilesetReader(baseDir string, r io.Reader, options ...LoaderOption) (*Tileset, error)
func (ts *Tileset) BaseDir() string
func (ts *Tileset) GetFileFullPath(fileName string) string
func (ts *Tileset) GetTileRect(tileID uint32) image.Rectangle
func (ts *Tileset) GetTilesetTile(tileID uint32) (*TilesetTile, error)
func (ts *Tileset) SetBaseDir(baseDir string)
type TilesetTile
type TilesetTileOffset
type WangColor
type WangPosition
type WangSet
func (w *WangSet) GetWangColors(tileID uint32) (map[WangPosition]*WangColor, error)
type WangSets
type WangTile

I expect the Go bindings, C bindings, and unit tests in Python.
If any of the functions are a poor fit for C, use an alternative that honors the design.

When finished, add a brief tutorial doc showing how to use this module from the perspective of the caller under docs/

Be sure to show how to read a tiled map and extract values in the tutorial.
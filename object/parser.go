package object

type Scene struct {
	Accessors      []Accessor     `json:"accessors"`
	Asset          Asset          `json:"asset"`
	BufferViews    []BufferView   `json:"bufferViews"`
	Buffers        []Buffer       `json:"buffers"`
	ExtensionsUsed []string       `json:"extensionsUsed"`
	Materials      []Material     `json:"materials"`
	Meshes         []Mesh         `json:"meshes"`
	Nodes          []Node         `json:"nodes"`
	Scene          int64          `json:"scene"`
	Scenes         []SceneElement `json:"scenes"`
}

type Accessor struct {
	BufferView    int64     `json:"bufferView"`
	ComponentType int64     `json:"componentType"`
	Count         int64     `json:"count"`
	Max           []float64 `json:"max"`
	Min           []float64 `json:"min"`
	Type          string    `json:"type"`
	ByteOffset    *int64    `json:"byteOffset,omitempty"`
}

type Asset struct {
	Extras    Extras `json:"extras"`
	Generator string `json:"generator"`
	Version   string `json:"version"`
}

type Extras struct {
	Author  string `json:"author"`
	License string `json:"license"`
	Source  string `json:"source"`
	Title   string `json:"title"`
}

type BufferView struct {
	Buffer     int64  `json:"buffer"`
	ByteLength int64  `json:"byteLength"`
	Name       string `json:"name"`
	Target     int64  `json:"target"`
	ByteOffset *int64 `json:"byteOffset,omitempty"`
	ByteStride *int64 `json:"byteStride,omitempty"`
}

type Buffer struct {
	ByteLength int64  `json:"byteLength"`
	URI        string `json:"uri"`
}

type Material struct {
	DoubleSided          bool                 `json:"doubleSided"`
	Extensions           Extensions           `json:"extensions"`
	Name                 string               `json:"name"`
	PbrMetallicRoughness PbrMetallicRoughness `json:"pbrMetallicRoughness"`
}

type Extensions struct {
	KHRMaterialsClearcoat KHRMaterialsClearcoat `json:"KHR_materials_clearcoat"`
}

type KHRMaterialsClearcoat struct {
	ClearcoatFactor          float64 `json:"clearcoatFactor"`
	ClearcoatRoughnessFactor float64 `json:"clearcoatRoughnessFactor"`
}

type PbrMetallicRoughness struct {
	BaseColorFactor []float64 `json:"baseColorFactor"`
	MetallicFactor  float64   `json:"metallicFactor"`
	RoughnessFactor float64   `json:"roughnessFactor"`
}

type Mesh struct {
	Name       string      `json:"name"`
	Primitives []Primitive `json:"primitives"`
}

type Primitive struct {
	Attributes Attributes `json:"attributes"`
	Indices    int64      `json:"indices"`
	Material   int64      `json:"material"`
	Mode       int64      `json:"mode"`
}

type Attributes struct {
	Normal    int64 `json:"NORMAL"`
	Position  int64 `json:"POSITION"`
	Texcoord0 int64 `json:"TEXCOORD_0"`
}

type Node struct {
	Children []int64   `json:"children"`
	Matrix   []float64 `json:"matrix"`
	Name     string    `json:"name"`
	Mesh     *int64    `json:"mesh,omitempty"`
}

type SceneElement struct {
	Name  string  `json:"name"`
	Nodes []int64 `json:"nodes"`
}

package object

import (
	"github.com/qmuntal/gltf"
)

type Mesh struct {
	Name       string
	Primitives []*Primitive
	Extensions map[string]interface{}
}

func parseMeshes(doc *gltf.Document, path string, images []*Image) []*Mesh {
	res := make([]*Mesh, len(doc.Meshes))
	for i, mesh := range doc.Meshes {
		m := &Mesh{
			Primitives: parsePrimitives(doc, mesh.Primitives, images, path, mesh.Name),
			Name:       mesh.Name,
			Extensions: mesh.Extensions,
		}
		res[i] = m
	}
	return res
}

package object

import "github.com/qmuntal/gltf"

type Mesh struct {
	Name       string
	Primitives []*Primitive
	Extensions map[string]interface{}
}

func (o *Object) parseMeshes(doc *gltf.Document, path string) {
	o.Meshes = make([]*Mesh, len(doc.Meshes))
	for i, mesh := range doc.Meshes {
		m := &Mesh{
			Primitives: o.parsePrimitives(doc, mesh.Primitives, path),
			Name:       mesh.Name,
			Extensions: mesh.Extensions,
		}
		o.Meshes[i] = m
	}
}

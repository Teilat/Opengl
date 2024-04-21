package object

import "github.com/qmuntal/gltf"

type Mesh struct {
	Name       string
	Primitives []*Primitive
	Extensions map[string]interface{}
}

func (o *Object) parseMeshes(doc *gltf.Document, path string) bool {
	withErr := false
	o.Meshes = make([]*Mesh, len(doc.Meshes))
	for i, mesh := range doc.Meshes {
		p, err := o.parsePrimitives(doc, mesh.Primitives, path)
		withErr = withErr || err
		m := &Mesh{
			Primitives: p,
			Name:       mesh.Name,
			Extensions: mesh.Extensions,
		}
		o.Meshes[i] = m
	}
	return withErr
}

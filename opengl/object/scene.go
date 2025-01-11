package object

import (
	"github.com/qmuntal/gltf"
)

type Scene struct {
	Name  string
	Nodes []*Node
}

func (o *Object) parseScenes(doc *gltf.Document, nodes []*Node) {
	o.Scene = make([]*Scene, len(doc.Scenes))
	for i, scene := range doc.Scenes {
		s := &Scene{
			Name: scene.Name,
		}
		for _, node := range scene.Nodes {
			s.Nodes = append(s.Nodes, nodes[node])
		}
		o.Scene[i] = s
	}
	defaultScene := 0
	if doc.Scene != nil {
		defaultScene = *doc.Scene
	}
	o.MainScene = o.Scene[defaultScene]
}

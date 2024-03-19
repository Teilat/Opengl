package object

import (
	"github.com/qmuntal/gltf"
)

type Scene struct {
	Name  string
	Nodes []*Node
}

func parseScenes(doc *gltf.Document, nodes []*Node) ([]*Scene, *Scene) {
	res := make([]*Scene, 0)
	for _, scene := range doc.Scenes {
		s := &Scene{
			Name: scene.Name,
		}
		for _, node := range scene.Nodes {
			s.Nodes = append(s.Nodes, nodes[node])
		}
		res = append(res, s)
	}
	defaultScene := uint32(0)
	if doc.Scene != nil {
		defaultScene = *doc.Scene
	}
	return res, res[defaultScene]
}

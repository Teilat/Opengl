package object

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/gltf"
)

type Object struct {
	Pos       mgl32.Vec3
	Scale     float32
	MainScene *Scene
	Scene     []*Scene
	Nodes     []*Node
	Meshes    []*Mesh
	Images    []*Image
}

func (o *Object) Parse(doc *gltf.Document, path string) bool {
	o.parseImages(doc)
	withErr := o.parseMeshes(doc, path)
	o.parseNodes(doc)
	o.parseScenes(doc, o.Nodes)
	return withErr
}

func (o *Object) GetPos() mgl32.Vec3 {
	return o.Pos
}

func (o *Object) Move(pos mgl32.Vec3) {
	o.Pos = pos
}

func (o *Object) Draw(program uint32) {
	if o == nil {
		fmt.Println("Object is nil")
		return
	}
	o.recursiveDraw(program, o.MainScene.Nodes)
}

func (o *Object) recursiveDraw(program uint32, nodes []*Node) {
	for _, node := range nodes {
		node.Draw(program)
		if len(node.children) > 0 {
			o.recursiveDraw(program, node.children)
		}
	}
}

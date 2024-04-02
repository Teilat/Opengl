package object

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/gltf"
)

type Object struct {
	Pos       mgl32.Vec3
	MainScene *Scene
	Scene     []*Scene
	Nodes     []*Node
	Meshes    []*Mesh
	Images    []*Image
}

func NewObject(pos mgl32.Vec3, path string) *Object {
	t := time.Now()
	doc, err := gltf.Open(path + "/scene.gltf")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Printf("parsing model:%s\n", path)
	fmt.Println("\ttotal mesh:", len(doc.Meshes))
	fmt.Println("\ttotal textures:", len(doc.Textures))
	fmt.Println("\ttotal images:", len(doc.Images))

	obj := &Object{Pos: pos}
	obj.parseImages(doc)
	obj.parseMeshes(doc, path)
	obj.parseNodes(doc)
	obj.parseScenes(doc, obj.Nodes)

	fmt.Printf("Done in %f milliseconds\n", time.Since(t).Seconds()*1000)
	runtime.GC()
	return obj
}

func (o *Object) GetPos() mgl32.Vec3 {
	return o.Pos
}

func (o *Object) Move(pos mgl32.Vec3) {
	o.Pos = pos
}

func (o *Object) Draw(program uint32) {
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

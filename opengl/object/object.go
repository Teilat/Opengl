package object

import (
	"fmt"
	"os"
	"runtime"
	"strings"
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

func NewObject(pos mgl32.Vec3, path string, binary bool) *Object {
	t := time.Now()

	folder, err := os.Open(path)
	if err != nil {
		return nil
	}
	files, err := folder.ReadDir(0)
	if err != nil {
		return nil
	}
	docFile := ""
	for _, f := range files {
		if binary {
			if strings.HasSuffix(f.Name(), ".glb") {
				docFile = f.Name()
				break
			}
		} else if strings.HasSuffix(f.Name(), ".gltf") {
			docFile = f.Name()
			break
		}
	}

	if docFile == "" {
		return nil
	}

	doc, err := gltf.Open(strings.Join([]string{path, docFile}, "/"))
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
	withErr := obj.parseMeshes(doc, path)
	obj.parseNodes(doc)
	obj.parseScenes(doc, obj.Nodes)

	fmt.Printf("Done in %f milliseconds. ", time.Since(t).Seconds()*1000)
	fmt.Printf("With errors: %t\n", withErr)
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

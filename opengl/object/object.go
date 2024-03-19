package object

import (
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/gltf"
	"image"
	"image/draw"
	"image/png"
	"os"
)

type NodeType int

const (
	Root NodeType = iota
	Child
)

type Object struct {
	Meshes    []*Mesh
	pos       mgl32.Vec3
	Images    []*gltf.Image
	MainScene *Scene
	Scene     []*Scene
	Nodes     []*Node
}

func NewObject(path string) *Object {
	doc, err := gltf.Open(path + "/scene.gltf")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("model:", path)
	fmt.Println("\ttotal mesh:", len(doc.Meshes))
	fmt.Println("\ttotal textures:", len(doc.Textures))
	fmt.Println("\ttotal images:", len(doc.Images))

	obj := &Object{}
	obj.Meshes = parseMeshes(doc)
	obj.Nodes = parseNodes(doc, obj.Meshes)
	obj.Scene, obj.MainScene = parseScenes(doc, obj.Nodes)
	obj.Images = doc.Images
	obj.pos = obj.MainScene.Nodes[0].translation

	if len(obj.Meshes[0].Texture1) > 0 && len(obj.Images) > 0 {
		obj.Meshes[0].Texture1Id = bindTexture(path + "/" + obj.Images[obj.Meshes[0].ImageId].URI)
	}

	return obj
}

func (o *Object) GetPos() mgl32.Vec3 {
	return o.pos
}

func (o *Object) Move(pos mgl32.Vec3) {
	o.pos = pos
}

func (o *Object) Draw(program uint32) {
	for _, mesh := range o.Meshes {
		gl.BindTexture(gl.TEXTURE_2D, mesh.Texture1Id)
		gl.BindVertexArray(mesh.Vao)

		model := mgl32.Translate3D(o.GetPos().Elem())
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("model\x00")), 1, false, &model[0])
		gl.DrawElements(gl.TRIANGLES, int32(len(mesh.Indices)), gl.UNSIGNED_INT, nil)
	}
}

func bindTexture(texturePath string) uint32 {
	var texture uint32

	img, err := getImageFromFilePath(texturePath)
	if err != nil {
		fmt.Println(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)

	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.SRGB_ALPHA, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(texture)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return texture
}

func getImageFromFilePath(file string) (image.Image, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	img, err := png.Decode(imgFile)
	if err != nil {
		return nil, err
	}
	return img, nil
}

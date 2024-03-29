package object

import (
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/gltf"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"runtime"
	"strings"
	"time"
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
	obj.Images = parseImages(doc)
	obj.Meshes = parseMeshes(doc, path, obj.Images)
	obj.Nodes = parseNodes(doc, obj.Meshes)
	obj.Scene, obj.MainScene = parseScenes(doc, obj.Nodes)

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
	var img image.Image
	if err != nil {
		return nil, err
	}

	ext := strings.Split(file, ".")
	switch ext[len(ext)-1] {
	case "png":
		img, err = png.Decode(imgFile)
		if err != nil {
			return nil, err
		}

	case "jpeg", "jpg":
		img, err = jpeg.Decode(imgFile)
		if err != nil {
			return nil, err
		}
	}

	imgFile.Close()
	fmt.Println(img.At(1, 1))
	return img, nil
}

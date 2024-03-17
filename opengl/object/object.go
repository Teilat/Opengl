package object

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

type Object struct {
	*gltf.Document
	Meshes []*Mesh
	pos    mgl32.Vec3
	Images []*gltf.Image
}

type Mesh struct {
	Name       string
	Material   *gltf.Material
	Indices    []uint32
	Vertices   [][3]float32
	Normal     [][3]float32
	Vao        uint32
	TextureId  uint32
	Texture    *gltf.Texture
	ImageId    uint32
	Texture1   [][2]float32
	Texture1Id uint32
	Tangent    [][4]float32
}

func NewObject(pos mgl32.Vec3, path string) *Object {
	doc, err := gltf.Open(path + "/scene.gltf")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("model:", path)
	fmt.Println("\ttotal meshes:", len(doc.Meshes))
	fmt.Println("\ttotal textures:", len(doc.Textures))
	fmt.Println("\ttotal images:", len(doc.Images))
	obj := &Object{
		Document: doc,
		pos:      pos,
		Images:   doc.Images,
	}
	for _, mesh := range doc.Meshes {
		m := Mesh{
			Name:     mesh.Name,
			Material: doc.Materials[*mesh.Primitives[0].Material],
		}
		if len(mesh.Primitives) > 1 {
			fmt.Println(m)
		}

		if m.Material.PBRMetallicRoughness.BaseColorTexture != nil {
			m.TextureId = m.Material.PBRMetallicRoughness.BaseColorTexture.Index
			m.Texture = doc.Textures[m.TextureId]
			m.ImageId = *m.Texture.Source
		}

		indices1, _ := modeler.ReadIndices(doc, doc.Accessors[*mesh.Primitives[0].Indices], []uint32{})
		m.Indices = indices1
		for attribute, index := range mesh.Primitives[0].Attributes {
			switch attribute {
			case "POSITION":
				positions1, _ := modeler.ReadPosition(doc, doc.Accessors[index], [][3]float32{})
				m.Vertices = positions1
			case "NORMAL":
				normal1, _ := modeler.ReadNormal(doc, doc.Accessors[index], [][3]float32{})
				m.Normal = normal1
			case "TEXCOORD_0":
				texture1, _ := modeler.ReadTextureCoord(doc, doc.Accessors[index], [][2]float32{})
				m.Texture1 = texture1
			case "TANGENT":
				tangent1, _ := modeler.ReadTangent(doc, doc.Accessors[index], [][4]float32{})
				m.Tangent = tangent1
			}
		}

		m.Vao = makeVAO(m.Vertices, m.Indices, m.Texture1)
		if len(m.Texture1) > 0 {
			m.Texture1Id = bindTexture(path + "/" + obj.Images[m.ImageId].URI)
		}

		obj.Meshes = append(obj.Meshes, &m)
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

func makeVAO(vertices [][3]float32, indices []uint32, texture [][2]float32) uint32 {
	var vertexArrayObject, vertexBufferObject, indexBufferObject uint32
	stride := 5 * int32(unsafe.Sizeof(float32(0)))
	if texture == nil {
		texture = make([][2]float32, len(vertices))
	}
	if vertices[0][0] > 1 || vertices[0][1] > 1 || vertices[0][2] > 1 {
		vertices = normalize(vertices)
	}

	vert := make([]float32, 0)
	for i, vertex := range vertices {
		vert = append(vert, vertex[0], vertex[1], vertex[2], texture[i][0], texture[i][1])
	}

	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.GenBuffers(1, &vertexBufferObject)
	gl.GenBuffers(1, &indexBufferObject)

	gl.BindVertexArray(vertexArrayObject)

	// привязываем буфер к массиву вертексов
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// сохраняем дангые в созданый буффер с определенным размером в битах
	gl.BufferData(gl.ARRAY_BUFFER, len(vert)*int(unsafe.Sizeof(float32(0))), gl.Ptr(vert), gl.STATIC_DRAW)
	// создаем указатель для использования данных в шейдере
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, nil)
	gl.EnableVertexAttribArray(0)
	// создаем указатель, но уже с офсетом в битах
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, stride, uintptr(3*int(unsafe.Sizeof(float32(0)))))
	gl.EnableVertexAttribArray(1)

	// создаем буфер индексов к массиву вертексов
	// тк это именно буфер индексов, делать указатель для шейдера ему не нужно
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBufferObject)
	// сохраняем дангые в созданый буффер с определенным размером в битах
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*int(unsafe.Sizeof(uint32(0))), gl.Ptr(indices), gl.STATIC_DRAW)

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	return vertexArrayObject
}

func normalize(vertices [][3]float32) [][3]float32 {
	maximum := float32(0)
	res := make([][3]float32, len(vertices))
	for _, vertice := range vertices {
		if maximum < max(vertice[0], vertice[1], vertice[2]) {
			maximum = max(vertice[0], vertice[1], vertice[2])
		}
	}
	for i, vertice := range vertices {
		res[i][0] = vertice[0] / maximum
		res[i][1] = vertice[1] / maximum
		res[i][2] = vertice[2] / maximum
	}
	return res
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

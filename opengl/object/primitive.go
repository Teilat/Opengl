package object

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

var PrimitiveModes = []uint32{
	gl.TRIANGLES,
	gl.POINTS,
	gl.LINES,
	gl.LINE_LOOP,
	gl.LINE_STRIP,
	gl.TRIANGLE_STRIP,
	gl.TRIANGLE_FAN,
}

type Primitive struct {
	Vao uint32

	PrimitiveMode uint32
	Indices       []uint32
	Vertices      [][3]float32
	Normal        [][3]float32
	Tangent       [][4]float32

	MetallicRoughness *TextureInfo
	BaseColor         *TextureInfo

	Texture   [][2]float32
	TextureId uint32
}

type TextureInfo struct {
	TextureIndex   uint32
	TextureImageId uint32
}

func (o *Object) parsePrimitives(doc *gltf.Document, primitives []*gltf.Primitive, path string) ([]*Primitive, bool) {
	res := make([]*Primitive, len(primitives))
	withErr := false
	for i, primitive := range primitives {
		p := &Primitive{
			PrimitiveMode: PrimitiveModes[primitive.Mode],
		}
		if primitive.Material != nil {
			mat := doc.Materials[*primitive.Material]
			if mat.PBRMetallicRoughness.BaseColorTexture != nil {
				p.BaseColor = &TextureInfo{}
				p.BaseColor.TextureIndex = uint32(mat.PBRMetallicRoughness.BaseColorTexture.Index)
				tex := doc.Textures[p.BaseColor.TextureIndex]
				if tex.Source != nil {
					p.BaseColor.TextureImageId = uint32(*tex.Source)
				}
			}
			if mat.PBRMetallicRoughness.MetallicRoughnessTexture != nil {
				p.MetallicRoughness = &TextureInfo{}
				p.MetallicRoughness.TextureIndex = uint32(mat.PBRMetallicRoughness.MetallicRoughnessTexture.Index)
				tex := doc.Textures[p.MetallicRoughness.TextureIndex]
				if tex.Source != nil {
					p.MetallicRoughness.TextureImageId = uint32(*tex.Source)
				}
			}
		}

		indices1, _ := modeler.ReadIndices(doc, doc.Accessors[*primitive.Indices], []uint32{})
		p.Indices = indices1
		for attribute, index := range primitive.Attributes {
			switch attribute {
			case "POSITION":
				positions, err := modeler.ReadPosition(doc, doc.Accessors[index], [][3]float32{})
				if err != nil {
					fmt.Println(err)
				}
				p.Vertices = positions
			case "TEXCOORD_0":
				texture1, err := modeler.ReadTextureCoord(doc, doc.Accessors[index], [][2]float32{})
				if err != nil {
					fmt.Println(err)
				}
				p.Texture = texture1
			case "NORMAL":
				normal, err := modeler.ReadNormal(doc, doc.Accessors[index], [][3]float32{})
				if err != nil {
					fmt.Println(err)
				}
				p.Normal = normal
			}
		}

		p.Vao = p.makeVAO()
		if p.BaseColor != nil {
			var err bool
			p.TextureId, err = bindTexture(path + "/" + o.Images[p.BaseColor.TextureImageId].URI)
			withErr = withErr || err
			// bindDepthTexture
		}
		res[i] = p
	}
	return res, withErr
}

func (p *Primitive) makeVAO() uint32 {
	var vertexArrayObject, vertexBufferObject, indexBufferObject uint32
	stride := 5 * int32(unsafe.Sizeof(float32(0)))
	if p.Texture == nil {
		p.Texture = make([][2]float32, len(p.Vertices))
	}
	if p.Vertices[0][0] > 1 || p.Vertices[0][1] > 1 || p.Vertices[0][2] > 1 {
		p.Vertices = normalize(p.Vertices)
	}

	vert := make([]float32, 0)
	for i, vertex := range p.Vertices {
		vert = append(vert, vertex[0], vertex[1], vertex[2], p.Texture[i][0], p.Texture[i][1])
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
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(p.Indices)*int(unsafe.Sizeof(uint32(0))), gl.Ptr(p.Indices), gl.STATIC_DRAW)

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

func bindTexture(texturePath string) (uint32, bool) {
	var texture uint32

	img, err := getImageFromFilePath(texturePath)
	if img == nil {
		if err != nil {
			fmt.Println(err)
		}
		return 0, true
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
	return texture, false
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
	return img, nil
}

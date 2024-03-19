package object

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
	"unsafe"
)

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

func parseMeshes(doc *gltf.Document) []*Mesh {
	res := make([]*Mesh, 0)
	for _, mesh := range doc.Meshes {
		m := Mesh{
			Name:     mesh.Name,
			Material: doc.Materials[*mesh.Primitives[0].Material],
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
			case "TEXCOORD_0":
				texture1, _ := modeler.ReadTextureCoord(doc, doc.Accessors[index], [][2]float32{})
				m.Texture1 = texture1
			}
		}

		m.Vao = m.makeVAO()
		res = append(res, &m)
	}
	return res
}

func (m *Mesh) makeVAO() uint32 {
	var vertexArrayObject, vertexBufferObject, indexBufferObject uint32
	stride := 5 * int32(unsafe.Sizeof(float32(0)))
	if m.Texture1 == nil {
		m.Texture1 = make([][2]float32, len(m.Vertices))
	}
	if m.Vertices[0][0] > 1 || m.Vertices[0][1] > 1 || m.Vertices[0][2] > 1 {
		m.Vertices = normalize(m.Vertices)
	}

	vert := make([]float32, 0)
	for i, vertex := range m.Vertices {
		vert = append(vert, vertex[0], vertex[1], vertex[2], m.Texture1[i][0], m.Texture1[i][1])
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
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.Indices)*int(unsafe.Sizeof(uint32(0))), gl.Ptr(m.Indices), gl.STATIC_DRAW)

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

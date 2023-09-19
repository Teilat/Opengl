package object

import (
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	"image"
	"image/draw"
	"image/png"
	"os"
	"unsafe"
)

type Object struct {
	vertices []float32
	Indices  []uint32
	Vao      uint32
	Texture  uint32
}

func NewObject(vertices []float32, indices []uint32) *Object {
	return &Object{
		vertices: vertices,
		Indices:  indices,
		Vao:      makeVAO(vertices, indices),
		Texture:  bindTexture("square.png"),
	}
}

// makeVAO initializes and returns a vertex array from the points provided.
// Этот метод по другому формирует объект с данными
func makeVAO(vertices []float32, indices []uint32) uint32 {
	var vertexArrayObject, vertexBufferObject, indexBufferObject uint32
	stride := 8 * int32(unsafe.Sizeof(float32(0)))

	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.GenBuffers(1, &vertexBufferObject)
	gl.GenBuffers(1, &indexBufferObject)

	gl.BindVertexArray(vertexArrayObject)

	// объеденяем все данные в виде (точки точки точки ... цвета цвета цвета ...)

	// привязываем буфер к массиву вертексов
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// сохраняем дангые в созданый буффер с определенным размером в битах
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*int(unsafe.Sizeof(float32(0))), gl.Ptr(vertices), gl.STATIC_DRAW)
	// создаем указатель для использования данных в шейдере
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, nil)
	gl.EnableVertexAttribArray(0)
	// создаем указатель, но уже с офсетом в битах
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, stride, uintptr(3*int(unsafe.Sizeof(float32(0)))))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, stride, uintptr(6*int(unsafe.Sizeof(float32(0)))))
	gl.EnableVertexAttribArray(2)

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

func bindTexture(texturePath string) uint32 {
	var texture uint32

	img, err := getImageFromFilePath("./" + texturePath)
	if err != nil {
		fmt.Println(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)

	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

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

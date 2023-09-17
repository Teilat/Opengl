package cell

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"unsafe"
)

// MakeVertexArrayObject initializes and returns a vertex array from the points provided.
func MakeVertexArrayObject(points, colors []float32, indices []uint32) uint32 {
	var vertexArrayObject, vertexBufferObject, colorBufferObject, indexBufferObject uint32
	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.GenBuffers(1, &vertexBufferObject)
	gl.GenBuffers(1, &colorBufferObject)
	gl.GenBuffers(1, &indexBufferObject)

	gl.BindVertexArray(vertexArrayObject)

	// привязываем буфер к массиву вертексов
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// сохраняем дангые в созданый буффер с определенным размером в битах
	// size = len(points) * size of data type | len(points)*int(unsafe.Sizeof(uint32(0)))
	gl.BufferData(gl.ARRAY_BUFFER, len(points)*int(unsafe.Sizeof(float32(0))), gl.Ptr(points), gl.STATIC_DRAW)
	// создаем указатель для использования данных в шейдере
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, int32(unsafe.Sizeof(float32(0)))*2, nil)
	// включаем указатель
	gl.EnableVertexAttribArray(0)

	//gl.BindBuffer(gl.ARRAY_BUFFER, colorBufferObject)
	//gl.BufferData(gl.ARRAY_BUFFER, len(colors)*int(unsafe.Sizeof(float32(0))), gl.Ptr(colors), gl.STATIC_DRAW)
	//gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(unsafe.Sizeof(float32(0)))*3, nil)
	//gl.EnableVertexAttribArray(1)

	// создаем буфер индексов к массиву вертексов
	// тк это именно буфер индексов, делать указатель ему не нужно
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBufferObject)
	// сохраняем дангые в созданый буффер с определенным размером в битах
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*int(unsafe.Sizeof(uint32(0))), gl.Ptr(indices), gl.STATIC_DRAW)

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	return vertexArrayObject
}

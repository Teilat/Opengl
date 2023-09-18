package cell

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"unsafe"
)

// MakeVAO initializes and returns a vertex array from the points provided.
// Этот метод по другому формирует объект с данными
func MakeVAO(points, colors []float32, indices []uint32) uint32 {
	var vertexArrayObject, vertexBufferObject, indexBufferObject uint32
	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.GenBuffers(1, &vertexBufferObject)
	gl.GenBuffers(1, &indexBufferObject)

	gl.BindVertexArray(vertexArrayObject)

	// объеденяем все данные в виде (точки точки точки ... цвета цвета цвета ...)
	vertices := append(points, colors...)

	// привязываем буфер к массиву вертексов
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// сохраняем дангые в созданый буффер с определенным размером в битах
	// size = len(points) * size of data type | len(points)*int(unsafe.Sizeof(uint32(0)))
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*int(unsafe.Sizeof(float32(0))), gl.Ptr(vertices), gl.STATIC_DRAW)
	// создаем указатель для использования данных в шейдере
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 0, uintptr(len(points)*int(unsafe.Sizeof(float32(0)))))
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

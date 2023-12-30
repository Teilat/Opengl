package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	posXx = 0.0
	posYy = 0.0
)

func CursorCallback(_ *glfw.Window, posX float64, posY float64) {
	posXx = posX
	posYy = -posY
}

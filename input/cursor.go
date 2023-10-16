package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"math"
)

var (
	posXx       = 0.0
	posYy       = 0.0
	sensitivity = 0.1
)

func CursorCallback(window *glfw.Window, posX float64, posY float64) {
	posXx = -math.Mod(posX*sensitivity, 360)
	posYy = posY * sensitivity
}

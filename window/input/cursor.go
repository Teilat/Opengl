package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	prevPosX = 0.0
	prevPosY = 0.0
	posX     = 0.0
	posY     = 0.0
)

func CursorCallback(w *glfw.Window, x, y float64) {
	prevPosX = posX
	prevPosY = posY

	posX += x - prevPosX
	posY -= y + prevPosY
}

func SetCursorPos(x, y float64) {
	posX, posY = x, y
}

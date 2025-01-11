package input

import (
	"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	prevPosX = 0.0
	prevPosY = 0.0
	posX     = 0.0
	posY     = 0.0
)

func CursorCallback(w *glfw.Window, x, y float64) {
	if w.GetAttrib(glfw.Focused) == 0 {
		return
	}

	prevPosX = posX
	prevPosY = posY

	posX += x - prevPosX
	posY += y - prevPosY
	fmt.Println(posX, prevPosX, x)
	fmt.Println(posY, prevPosY, y)

}

func SetCursorPos(x, y float64) {
	posX, posY = x, y
}

package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	prevX        = 0.0
	prevY        = 0.0
	offsetX      = 0.0
	offsetY      = 0.0
	sensitivityY = 1.0
	sensitivityX = 0.1
)

func CursorCallback(window *glfw.Window, posX float64, posY float64) {
	offsetX = (prevX - posX) * sensitivityX
	// Обратный порядок вычитания потому что оконные Y-координаты возрастают с верху вниз
	offsetY = (posY - prevY) * sensitivityY
	prevX = posX
	prevY = posY
	if offsetX != 0 || offsetY != 0 {
		//fmt.Printf("%s:%f,%f\n", time.Now().Format("04:05.000"), offsetX, offsetY)
	}
}

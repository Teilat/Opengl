package input

import (
	"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
	"time"
)

type cursor struct {
	prevX   float64
	prevY   float64
	offsetX float64
	offsetY float64

	sensitivity float64
}

var cur = cursor{sensitivity: 0.05}

func CursorCallback(window *glfw.Window, posX float64, posY float64) {
	cur.offsetX = (posX - cur.prevX) * cur.sensitivity
	// Обратный порядок вычитания потому что оконные Y-координаты возрастают с верху вниз
	cur.offsetY = (cur.prevY - posY) * cur.sensitivity
	cur.prevX = posX
	cur.prevY = posY
	if cur.offsetX != 0 || cur.offsetY != 0 {
		fmt.Printf("%s:%f,%f\n", time.Now().Format("04:05.000"), cur.offsetX, cur.offsetY)
	}
}

func GetMouseMovement() (float64, float64) {
	return cur.offsetX, cur.offsetY
}

package key_manager

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

var UpdateColor = true

func KeyCallBack(window *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
	if key == glfw.KeySpace && action == glfw.Press {
		UpdateColor = !UpdateColor
	}

	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

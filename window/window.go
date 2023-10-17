package window

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"opengl/input"
)

type Window struct {
	*glfw.Window
	monitor    *glfw.Monitor
	title      string
	fullScreen bool

	posX        int
	posY        int
	width       int
	height      int
	refreshRate int
}

// InitGlfw initializes glfw and returns a Window to use.
func InitGlfw(width, height, refreshRate int, title string, fullscreen bool,
	keyCallback func(w *glfw.Window, key glfw.Key, _ int, action glfw.Action, modKey glfw.ModifierKey),
	cursorCallback func(window *glfw.Window, posX float64, posY float64),
) *Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLAnyProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	var monitor *glfw.Monitor

	if fullscreen {
		monitor = glfw.GetMonitors()[1]
		if monitor == nil {
			monitor = glfw.GetPrimaryMonitor()
		}
		vMode := monitor.GetVideoMode()
		width = vMode.Width
		height = vMode.Height
	}

	window, err := glfw.CreateWindow(width, height, title, monitor, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	if glfw.RawMouseMotionSupported() {
		window.SetInputMode(glfw.RawMouseMotion, glfw.True)
	}

	window.SetKeyCallback(keyCallback)
	window.SetCursorPosCallback(cursorCallback)

	return &Window{Window: window, height: height, width: width, refreshRate: refreshRate, title: title, monitor: nil}
}

func (w *Window) GetWidth() int {
	return w.width
}

func (w *Window) GetHeight() int {
	return w.height
}

func (w *Window) UpdateWindow() bool {
	if input.GetKeyWithModKey(glfw.KeyEnter, glfw.ModAlt) {
		w.switchWindowMode()
		return true
	}
	return false
}

func (w *Window) switchWindowMode() {
	if w.fullScreen {
		w.monitor = nil
		// refreshRate ignored when windowed
		w.SetMonitor(w.monitor, w.posX, w.posY, w.width, w.height, 0)
	} else {
		w.monitor = glfw.GetPrimaryMonitor()
		vMode := w.monitor.GetVideoMode()
		w.SetMonitor(w.monitor, 0, 0, vMode.Width, vMode.Height, vMode.RefreshRate)
	}
	w.fullScreen = !w.fullScreen
}

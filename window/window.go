package window

import (
	"time"

	"opengl/input"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type Window struct {
	*glfw.Window
	monitor        *glfw.Monitor
	title          string
	FullScreenLock time.Time
	fullScreen     bool

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
	var monitor *glfw.Monitor
	var window *glfw.Window
	var err error
	if err = glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLAnyProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.RefreshRate, refreshRate)

	monitor = glfw.GetMonitors()[0]
	if monitor == nil {
		monitor = glfw.GetPrimaryMonitor()
	}

	if fullscreen {
		vMode := monitor.GetVideoMode()
		window, err = glfw.CreateWindow(vMode.Width, vMode.Height, title, monitor, nil)
	} else {
		window, err = glfw.CreateWindow(width, height, title, nil, nil)
	}
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

	return &Window{Window: window, height: height, width: width, refreshRate: refreshRate, title: title, monitor: monitor, FullScreenLock: time.Now()}
}

func (w *Window) GetWidth() int {
	if w.fullScreen {
		return w.monitor.GetVideoMode().Width
	}
	return w.width
}

func (w *Window) GetHeight() int {
	if w.fullScreen {
		return w.monitor.GetVideoMode().Height
	}
	return w.height
}

func (w *Window) UpdateWindow() bool {
	if input.GetKeyPressWithModKey(glfw.KeyEnter, glfw.ModAlt) && time.Since(w.FullScreenLock) > time.Millisecond*250 {
		w.FullScreenLock = time.Now()
		w.switchWindowMode()
		return true
	}
	return false
}

func (w *Window) switchWindowMode() {
	if w.fullScreen {
		x, y, _, _ := w.monitor.GetWorkarea()
		// refreshRate ignored when windowed
		w.SetMonitor(nil, x+20, y+50, w.width, w.height, 0)
	} else {
		vMode := w.monitor.GetVideoMode()
		w.SetMonitor(w.monitor, 0, 0, vMode.Width, vMode.Height, vMode.RefreshRate)
	}
	w.fullScreen = !w.fullScreen
}

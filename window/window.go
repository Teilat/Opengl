package window

import (
	"opengl/opengl/camera"
	"opengl/window/input"
	"opengl/window/text"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Window struct {
	*glfw.Window
	title       string
	posX        int
	posY        int
	width       int
	height      int
	refreshRate int

	monitor        *glfw.Monitor
	fullScreenLock time.Time
	fullScreen     bool

	Text *text.Text
}

func OnResize(_ *glfw.Window, width int, height int) {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	gl.Viewport(0, 0, int32(width), int32(height))
}

// InitGlfw initializes glfw and returns a Window to use.
func InitGlfw(width, height, refreshRate int, title string, fullscreen bool,
	keyCallback func(w *glfw.Window, key glfw.Key, _ int, action glfw.Action, modKey glfw.ModifierKey),
	cursorCallback func(window *glfw.Window, posX float64, posY float64),
	resizeCallback func(window *glfw.Window, width int, height int),
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
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
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
	window.SetSizeCallback(resizeCallback)

	input.SetCursorPos(window.GetCursorPos())

	return &Window{
		Window:         window,
		Text:           text.Init(32, width, height, "arial.ttf"),
		monitor:        monitor,
		title:          title,
		fullScreenLock: time.Now(),
		width:          width,
		height:         height,
		refreshRate:    refreshRate,
	}
}

func (w *Window) GetRefreshRate() int {
	if w.fullScreen {
		return w.monitor.GetVideoMode().RefreshRate
	}
	return w.refreshRate
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

func (w *Window) GetCursorPos() (float64, float64) {
	return w.Window.GetCursorPos()
}

func (w *Window) OnWindowModeChange(cam *camera.Camera) {
	if input.GetKeyPressWithModKey(glfw.KeyEnter, glfw.ModAlt) {
		if time.Since(w.fullScreenLock) > time.Millisecond*250 {
			w.fullScreenLock = time.Now()
			w.switchWindowMode()
			// post update
			width, height := w.GetWidth(), w.GetHeight()
			cam.UpdateWindow(float32(width), float32(height))
			gl.Viewport(0, 0, int32(width), int32(height))
			w.Text.UpdateResolution(width, height)
		}
	}
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

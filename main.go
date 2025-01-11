package main

import (
	"C"
	"context"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"opengl/metric"
	"opengl/opengl"
	"opengl/opengl/camera"
	"opengl/opengl/object"
	"opengl/window"
	"opengl/window/input"
	"opengl/window/text"
	"runtime"
	"time"
)

var (
	Width  = 800
	Height = 600

	FpsLock = true
	Fps     = 120
)

func main() {
	runtime.LockOSThread()
	ctx, cancel := context.WithCancel(context.Background())
	fixedUpdateTicker := time.NewTicker(time.Second / time.Duration(Fps*2))
	debugTicker := time.NewTicker(metric.TickerResolution)

	win := window.InitGlfw(Width, Height, Fps, "Program", false, input.KeyCallback, input.CursorCallback, window.OnResize)
	program := opengl.InitOpenGL(false)

	cam := camera.NewCamera(ctx, fixedUpdateTicker, program, 80, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 0, 0}, win.GetWidth(), win.GetHeight())

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)

	fpsMeter := metric.NewFPSMeter()
	go metric.StartPprof()

	cam.AddDebug(debugTicker)
	cam.StartDebug()

	if c := cam.GetDebug(); c != nil {
		win.Text.AddText([]*text.Item{
			{Text: c.GetPosString(), PosX: 0, Scale: 0.5},
			{Text: c.GetLookAtString(), PosX: 0, Scale: 0.5},
			{Text: c.GetFovString(), PosX: 0, Scale: 0.5},
		})
	}

	win.Text.AddText([]*text.Item{
		{Text: fpsMeter.GetString(), PosX: 0, Scale: 0.5},
	})
	objectManager := object.NewManager()
	//objectManager.AddObject(object.NewObject(mgl32.Vec3{3, 0, 3}, "./models/BoxVertexColors", true))
	//objectManager.AddObject(object.NewObject(mgl32.Vec3{0, 0, 0}, "./models/Cube", true))
	objectManager.AddObject(object.NewObject(mgl32.Vec3{0, 0, 0}, "./models/Avocado", false))

	fpsMeter.Start(ctx, debugTicker)
	for !win.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()
		win.OnWindowModeChange(cam)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		cam.Update()

		for _, obj := range objectManager.Objects {
			obj.Draw(program)
		}
		win.Text.DrawText()
		win.SwapBuffers()

		gl.Finish()
		fpsMeter.Tick(t)

		if FpsLock {
			time.Sleep(time.Second/time.Duration(Fps) - time.Since(t))
		}
	}
	cancel()
	gl.DeleteProgram(program)
	glfw.Terminate()
}

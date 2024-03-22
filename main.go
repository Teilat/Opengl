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
	Width  = 1280
	Height = 1024

	FpsLock = true
	Fps     = 75
)

func main() {
	runtime.LockOSThread()
	ctx, cancel := context.WithCancel(context.Background())
	fixedUpdateTicker := time.NewTicker(time.Second / time.Duration(Fps*2))

	win := window.InitGlfw(Width, Height, Fps, "Program", false, input.KeyCallback, input.CursorCallback, window.OnResize)
	program := opengl.InitOpenGL(false)

	cam := camera.NewCamera(ctx, fixedUpdateTicker, program, 80, mgl32.Vec3{3, 0, -3}, mgl32.Vec3{0, 0, 0}, win.GetWidth(), win.GetHeight())

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)

	fpsMeter := metric.NewFPSMeter()

	win.Text.AddText([]*text.Item{
		{Text: cam.GetDebug().GetPosString(), PosX: 0, Scale: 0.5},
		{Text: cam.GetDebug().GetLookAtString(), PosX: 0, Scale: 0.5},
		{Text: cam.GetDebug().GetFovString(), PosX: 0, Scale: 0.5},
		{Text: fpsMeter.GetString(), PosX: 0, Scale: 0.5},
	})
	objectManager := object.NewManager()
	objectManager.AddObject(object.NewObject(mgl32.Vec3{3, 0, 3}, "./models/Torus Knot"))
	//objectManager.AddObject(object.NewObject(mgl32.Vec3{0, 0, 0}, "./models/Cube"))
	//objectManager.AddObject(object.NewObject(mgl32.Vec3{0, 0, 0}, "./models/Open Cube"))
	//objectManager.AddObject(object.NewObject(mgl32.Vec3{-3, 0, 3}, "./models/Sphere"))
	//objectManager.AddObject(object.NewObject(mgl32.Vec3{0, 0, 0}, "./models/Datsun 240k"))
	//objectManager.AddObject(object.NewObject(mgl32.Vec3{0, 0, 0}, "./models/Mazda mx-5"))

	fpsMeter.Start(ctx)
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
		fpsMeter.Tick()
		if FpsLock {
			time.Sleep(time.Second/time.Duration(Fps) - time.Since(t))
		}
	}
	cancel()
	gl.DeleteProgram(program)
	glfw.Terminate()
}

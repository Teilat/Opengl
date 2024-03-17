package main

import (
	"C"
	"context"
	"fmt"
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

const TickerResolution = time.Millisecond

var (
	Width  = 800
	Height = 600

	Fps = 75
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	runtime.LockOSThread()
	win := window.InitGlfw(Width, Height, Fps, "Program", false, input.KeyCallback, input.CursorCallback, window.OnResize)
	defer glfw.Terminate()
	program := opengl.InitOpenGL()

	gl.UseProgram(program)
	gl.Enable(gl.DEPTH_TEST)
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	cam := camera.NewCamera(program, 80, mgl32.Vec3{-7, 7, 7}, mgl32.Vec3{0, 0, 0}, win.GetWidth(), win.GetHeight())
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)

	camPos := fmt.Sprintf("camera pos:%v", cam.GetPos())
	lookAt := fmt.Sprintf("look at:%v", cam.GetLookAt())
	fov := fmt.Sprintf("fov:%v", cam.GetFov())
	go updCam(ctx, program, cam, &camPos, &lookAt, &fov)
	fpsMeter := metric.NewFPSMeter(ctx)

	win.Text.AddText([]*text.Item{
		{Text: &camPos, PosX: 0, Scale: 0.5},
		{Text: &lookAt, PosX: 0, Scale: 0.5},
		{Text: &fov, PosX: 0, Scale: 0.5},
		{Text: fpsMeter.GetString(), PosX: 0, Scale: 0.5},
	})

	obj := object.NewObject(mgl32.Vec3{3, 0, 0}, "./models/Torus Knot")
	//obj := object.NewObject(mgl32.Vec3{6, 0, 3}, "./models/Cube")
	//obj := object.NewObject(mgl32.Vec3{3, 0, 0}, "./models/Open Cube")
	//obj := object.NewObject(mgl32.Vec3{-3, 0, 3}, "./models/Sphere")
	//obj := object.NewObject(mgl32.Vec3{0, 0, 0}, "./models/Car")

	for !win.ShouldClose() {
		glfw.PollEvents()
		win.OnWindowModeChange(cam)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		upd(program, cam, &camPos, &lookAt, &fov)

		draw(win, program, []*object.Object{obj})
		win.SwapBuffers()

		gl.Finish()
		fpsMeter.Tick()
	}
	cancel()
	gl.DeleteProgram(program)
}

func upd(program uint32, cam *camera.Camera, pos, lookAt, fov *string) {
	gl.UseProgram(program)
	cam.Update()
}

func updCam(ctx context.Context, program uint32, cam *camera.Camera, pos, lookAt, fov *string) {
	ticker := time.NewTicker(TickerResolution)
	for {
		select {
		case <-ticker.C:
			glfw.PollEvents()
			gl.UseProgram(program)
			cam.Update()
			*pos = fmt.Sprintf("camera pos:%v", cam.GetPos())
			*lookAt = fmt.Sprintf("look at:%v", cam.GetLookAt())
			*fov = fmt.Sprintf("fov:%v", cam.GetFov())
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func draw(win *window.Window, program uint32, objs []*object.Object) {
	gl.UseProgram(program)
	for _, obj := range objs {
		for _, mesh := range obj.Meshes {
			gl.BindTexture(gl.TEXTURE_2D, mesh.Texture1Id)
			gl.BindVertexArray(mesh.Vao)

			model := mgl32.Translate3D(obj.GetPos().Elem())
			gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("model\x00")), 1, false, &model[0])
			gl.DrawElements(gl.TRIANGLES, int32(len(mesh.Indices)), gl.UNSIGNED_INT, nil)
		}
	}
	win.Text.DrawText()
}

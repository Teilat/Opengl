package main

import (
	"C"
	"fmt"
	"opengl/opengl"
	"opengl/opengl/camera"
	"opengl/opengl/object"
	"opengl/window/input"
	"opengl/window/text"
	"runtime"
	"time"

	"opengl/window"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	Width  = 800
	Height = 600

	Fps = 75
)

func main() {
	runtime.LockOSThread()
	win := window.InitGlfw(Width, Height, Fps, "Program", false, input.KeyCallback, input.CursorCallback, window.OnResize)
	defer glfw.Terminate()
	program := opengl.InitOpenGL()

	gl.UseProgram(program)
	gl.Enable(gl.DEPTH_TEST)
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	//obj := object.NewObject(mgl32.Vec3{3, 0, 0}, "square.png", "./opengl/object/Torus Knot")
	//obj2 := object.NewObject(mgl32.Vec3{0, 0, 3}, "square.png", "./opengl/object/Cube")
	//obj := object.NewObject(mgl32.Vec3{3, 0, 0}, "", "./opengl/object/Open Cube")
	//obj := object.NewObject(mgl32.Vec3{-3, 0, 3}, "", "./opengl/object/Sphere")
	obj := object.NewObject(mgl32.Vec3{-3, 0, 3}, "", "./opengl/object/Car")
	cam := camera.NewCamera(program, 80, mgl32.Vec3{-7, 7, 7}, mgl32.Vec3{0, 0, 0}, win.GetWidth(), win.GetHeight())

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)

	for !win.ShouldClose() {
		t := time.Now()

		glfw.PollEvents()
		win.OnWindowModeChange(cam)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		upd(program, cam)
		draw([]*object.Object{obj}, program)

		camPos := fmt.Sprintf("camera pos:%v", cam.GetPos())
		lookAt := fmt.Sprintf("look at:%v", cam.GetLookAt())
		fov := fmt.Sprintf("fov:%v", cam.GetFov())
		win.Text.DrawText([]text.Item{
			{Text: &camPos, PosX: 0, PosY: 16, Scale: 0.5},
			{Text: &lookAt, PosX: 0, PosY: 32, Scale: 0.5},
			{Text: &fov, PosX: 0, PosY: 48, Scale: 0.5},
		})

		// TODO fps counter
		// мапка с таймом в ключе
		// first in last out,
		win.SwapBuffers()

		time.Sleep(time.Second/time.Duration(Fps) - time.Since(t))
	}
	gl.DeleteProgram(program)
}

func upd(program uint32, cam *camera.Camera) {
	gl.UseProgram(program)
	cam.Update()
}

func draw(objs []*object.Object, program uint32) {
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
}

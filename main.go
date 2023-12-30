package main

import (
	"C"
	"math"
	"opengl/opengl"
	"opengl/opengl/camera"
	"opengl/opengl/object"
	input2 "opengl/window/input"
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

var (
	square = []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
	}
	squareIndices = []uint32{
		// front
		3, 1, 2,
		3, 0, 1,
		// back
		6, 5, 4,
		6, 4, 7,
		//bottom
		0, 1, 5,
		0, 4, 5,
		// top
		2, 3, 7,
		2, 6, 7,
		// right
		2, 5, 1,
		2, 5, 6,
		// left
		0, 7, 4,
		0, 7, 3,
	}
)

func main() {
	runtime.LockOSThread()

	win := window.InitGlfw(Width, Height, Fps, "Program", false, input2.KeyCallback, input2.CursorCallback, window.OnResize)
	defer glfw.Terminate()
	program := opengl.InitOpenGL()

	gl.UseProgram(program)
	gl.Enable(gl.DEPTH_TEST)
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	obj := object.NewObject(square, squareIndices, mgl32.Vec3{0, 0, 0}, "square.png")
	cam := camera.NewCamera(program, 45, mgl32.Vec3{1, 1, 1}, win.GetWidth(), win.GetHeight())

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, obj.Texture)
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("tex\x00")), 0)

	vertexColorLocation := gl.GetUniformLocation(program, gl.Str("ourColor\x00"))

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	for !win.ShouldClose() {
		t := time.Now()

		glfw.PollEvents()
		if win.OnModeChange() {
			cam.UpdateWindow(float32(win.GetWidth()), float32(win.GetHeight()))
			gl.Viewport(0, 0, int32(win.GetWidth()), int32(win.GetHeight()))
		}

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		upd(program, vertexColorLocation, cam)
		draw(obj, program)

		win.SwapBuffers()

		time.Sleep(time.Second/time.Duration(Fps) - time.Since(t))
	}
	gl.DeleteProgram(program)
}

func upd(program uint32, vertexColorLocation int32, cam *camera.Camera) {
	gl.UseProgram(program)
	t := glfw.GetTime()
	redValue := math.Abs(math.Cos(t))
	greenValue := math.Abs(math.Cos(t + 1))
	blueValue := math.Abs(math.Cos(t + 2))
	gl.Uniform4f(vertexColorLocation, float32(redValue), float32(greenValue), float32(blueValue), 1.0)

	cam.Update()
}

func draw(obj *object.Object, program uint32) {
	gl.UseProgram(program)

	gl.BindTexture(gl.TEXTURE_2D, obj.Texture)
	gl.BindVertexArray(obj.Vao)

	model := mgl32.Translate3D(obj.GetPos().Elem())
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("model\x00")), 1, false, &model[0])
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)))
}

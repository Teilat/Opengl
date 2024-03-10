package main

import (
	"C"
	"fmt"
	"math"
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

	obj := object.NewObject(mgl32.Vec3{0, 0, 0}, "square.png")
	cam := camera.NewCamera(program, 45, mgl32.Vec3{1, 1, 1}, win.GetWidth(), win.GetHeight())

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, obj.Texture)
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("tex\x00")), 0)

	vertexColorLocation := gl.GetUniformLocation(program, gl.Str("ourColor\x00"))

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)

	str := ""
	gl.Uniform4f(vertexColorLocation, float32(1), float32(1), float32(1), 1.0)

	for !win.ShouldClose() {
		t := time.Now()

		glfw.PollEvents()
		win.OnWindowModeChange(cam)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		upd(program, vertexColorLocation, cam, &str)
		draw(obj, program)

		camPos := fmt.Sprintf("camera pos:%v", cam.GetPos())
		lookAtPos := fmt.Sprintf("look at pos:%v", cam.GetLookAt())
		win.Text.DrawText([]text.Item{
			{Text: &str, PosX: 0, PosY: 16, Scale: 0.5},
			{Text: &camPos, PosX: 0, PosY: 32, Scale: 0.5},
			{Text: &lookAtPos, PosX: 0, PosY: 48, Scale: 0.5},
		})

		win.SwapBuffers()

		time.Sleep(time.Second/time.Duration(Fps) - time.Since(t))
	}
	gl.DeleteProgram(program)
}

func upd(program uint32, vertexColorLocation int32, cam *camera.Camera, s *string) {
	gl.UseProgram(program)
	t := glfw.GetTime()
	redValue := math.Abs(math.Cos(t))
	greenValue := math.Abs(math.Cos(t + 1))
	blueValue := math.Abs(math.Cos(t + 2))
	gl.Uniform4f(vertexColorLocation, float32(redValue), float32(greenValue), float32(blueValue), 1.0)
	*s = fmt.Sprintf("%f,%f,%f", float32(redValue), float32(greenValue), float32(blueValue))
	cam.Update()
}

func draw(obj *object.Object, program uint32) {
	gl.UseProgram(program)

	gl.BindTexture(gl.TEXTURE_2D, obj.Texture)
	gl.BindVertexArray(obj.Vao)

	model := mgl32.Translate3D(obj.GetPos().Elem())
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("model\x00")), 1, false, &model[0])
	// не работает нормально наложение текстур
	gl.DrawElements(gl.TRIANGLES, int32(len(obj.Indices)), gl.UNSIGNED_INT, nil)
	//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)))
}

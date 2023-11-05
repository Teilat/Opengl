package main

import (
	"C"
	"fmt"
	"log"
	"math"
	"runtime"
	"time"

	"opengl/camera"
	"opengl/input"
	"opengl/object"
	"opengl/shader"
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

	win := window.InitGlfw(Width, Height, Fps, "Program", false, input.KeyCallback, input.CursorCallback)
	defer glfw.Terminate()
	program := initOpenGL()

	gl.UseProgram(program)
	gl.Enable(gl.DEPTH_TEST)
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	obj := object.NewObject(mgl32.Vec3{0, 0, 0}, "square.png")
	cam := camera.NewCamera(program, "camera\x00", "projection\x00", 45, mgl32.Vec3{1, 1, 1}, win.GetWidth(), win.GetHeight())

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, obj.Texture)
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("tex\x00")), 0)

	vertexColorLocation := gl.GetUniformLocation(program, gl.Str("ourColor\x00"))

	for !win.ShouldClose() {
		t := time.Now()

		glfw.PollEvents()
		if win.UpdateWindow() {
			cam.UpdateWindow(win.GetWidth(), win.GetHeight())
			gl.Viewport(0, 0, int32(win.GetWidth()), int32(win.GetHeight()))
		}

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
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
	// не работает нормально наложение текстур
	gl.DrawElements(gl.TRIANGLES, int32(len(obj.Indices)), gl.UNSIGNED_INT, nil)
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	log.Println("OpenGL version", gl.GoStr(gl.GetString(gl.VERSION)))

	fragmentShader, err := shader.CompileShader("shader.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		fmt.Println(err)
	}

	vertexShader, err := shader.CompileShader("shader.vert", gl.VERTEX_SHADER)
	if err != nil {
		fmt.Println(err)
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, fragmentShader)
	gl.AttachShader(program, vertexShader)

	gl.LinkProgram(program)
	gl.ValidateProgram(program)

	gl.DeleteShader(fragmentShader)
	gl.DeleteShader(vertexShader)
	return program
}

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

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 800
	height = 600

	fps = 75
)

var updColor = false
var a = float32(0)

var (
	square = []float32{
		// левый нижний ближний
		-0.5, -0.5, 0.5,
		1.0, 0.0, 0.0,
		0.0, 0.0,

		// правый нижний ближний
		0.5, -0.5, 0.5,
		0.0, 1.0, 0.0,
		1.0, 0.0,

		// левый верхний ближний
		0.5, 0.5, 0.5,
		0.0, 0.0, 1.0,
		1.0, 1.0,

		// правый верхний ближний
		-0.5, 0.5, 0.5,
		1.0, 1.0, 1.0,
		0.0, 1.0,

		// левый нижний дальний
		-0.5, -0.5, -0.5,
		1.0, 0.0, 0.0,
		0.0, 0.0,

		// правый нижний дальний
		0.5, -0.5, -0.5,
		0.0, 1.0, 0.0,
		1.0, 0.0,

		// левый верхний дальний
		0.5, 0.5, -0.5,
		0.0, 0.0, 1.0,
		1.0, 1.0,

		// правый верхний дальний
		-0.5, 0.5, -0.5,
		1.0, 1.0, 1.0,
		0.0, 1.0,
	}

	squareIndices = []uint32{
		// front
		0, 1, 2,
		0, 3, 2,
		// back
		4, 5, 6,
		4, 7, 6,
		//bottom
		0, 1, 4,
		0, 5, 4,
		// top
		2, 6, 3,
		2, 7, 3,
		// right
		1, 5, 2,
		1, 6, 2,
		// left
		0, 4, 3,
		0, 7, 3,
	}
)

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	gl.UseProgram(program)
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	window.SetKeyCallback(input.KeyCallBack)

	obj := object.NewObject(square, squareIndices)

	projection := mgl32.Perspective(mgl32.DegToRad(40), float32(width)/height, 0.1, 10.0)
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("projection\x00")), 1, false, &projection[0])

	cam := camera.NewCamera(gl.GetUniformLocation(program, gl.Str("camera\x00")), mgl32.Vec3{3, 0, 3})
	gl.UniformMatrix4fv(cam.ShaderLocation, 1, false, cam.GetMatrix4fv())

	model := mgl32.Ident4()
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("model\x00")), 1, false, &model[0])

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, obj.Texture)
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("tex\x00")), 0)

	vertexColorLocation := gl.GetUniformLocation(program, gl.Str("ourColor\x00"))

	for !window.ShouldClose() {
		t := time.Now()

		glfw.PollEvents()

		gl.Clear(gl.COLOR_BUFFER_BIT)

		if input.GetKeyDown(glfw.KeySpace) {
			updColor = !updColor
		}
		if updColor {
			upd(program, vertexColorLocation, cam)
		}

		draw(obj.Vao, obj.Texture, window, program)
		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
	gl.DeleteProgram(program)
}

func upd(program uint32, vertexColorLocation int32, cam *camera.Camera) {
	gl.UseProgram(program)
	t := glfw.GetTime()
	redValue := math.Abs(math.Cos(t))
	greenValue := math.Abs(math.Cos(t + 1))
	blueValue := math.Abs(math.Cos(t + 2))
	//fmt.Printf("%.2f %.2f %.2f\n", redValue, greenValue, blueValue)

	matrix := cam.GetMatrix4()
	matrix = matrix.Mul4(mgl32.HomogRotate3D(a, mgl32.Vec3{0, 1, 0}))
	a += 0.05
	if a >= 360 {
		a = 0
	}
	gl.Uniform4f(vertexColorLocation, float32(redValue), float32(greenValue), float32(blueValue), 1.0)
	gl.UniformMatrix4fv(cam.ShaderLocation, 1, false, &matrix[0])
}

func draw(vao, texture uint32, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.UseProgram(program)

	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.BindVertexArray(vao)
	gl.DrawElements(gl.TRIANGLES, int32(len(squareIndices)), gl.UNSIGNED_INT, nil)

	window.SwapBuffers()
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLAnyProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	monitor := glfw.GetMonitors()[1]
	if monitor == nil {
		monitor = glfw.GetPrimaryMonitor()
	}

	window, err := glfw.CreateWindow(width, height, "Program", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	shaders := make([]uint32, 0)

	if err := gl.Init(); err != nil {
		panic(err)
	}
	log.Println("OpenGL version", gl.GoStr(gl.GetString(gl.VERSION)))

	fragmentShader, err := shader.CompileShader("shader.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		fmt.Println(err)
	}
	shaders = append(shaders, fragmentShader)

	vertexShader, err := shader.CompileShader("shader.vert", gl.VERTEX_SHADER)
	if err != nil {
		fmt.Println(err)
	}
	shaders = append(shaders, vertexShader)

	program := gl.CreateProgram()

	gl.AttachShader(program, fragmentShader)
	gl.AttachShader(program, vertexShader)

	gl.LinkProgram(program)
	gl.ValidateProgram(program)

	gl.DeleteShader(fragmentShader)
	gl.DeleteShader(vertexShader)
	return program
}

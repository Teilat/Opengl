package main

import (
	"C"
	"bytes"
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"image/png"
	"log"
	"math"
	"opengl/key_manager"
	"opengl/shader"
	"opengl/vertex"
	"os"
	"runtime"
	"time"
)

const (
	width  = 800
	height = 600

	fps = 75
)

var (
	square = []float32{
		-0.5, -0.5, 0.5, // 0
		0.5, -0.5, 0.5, // 1
		0.5, 0.5, 0.5, // 2
		-0.5, 0.5, 0.5, // 3

		-0.5, -0.5, -0.5, // 4
		0.5, -0.5, -0.5, // 5
		0.5, 0.5, -0.5, // 6
		-0.5, 0.5, -0.5, // 7
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
	squareColors = []float32{
		1, 0, 1, // 0
		0, 1, 0, // 1
		1, 0, 1, // 2
		1, 1, 1, // 3
	}
	texture = []float32{
		0, 0, // 0
		1, 0, // 1
		1, 1, // 2
		0, 1, // 3
	}
)

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	vao := vertex.MakeVAO(square, squareColors, texture, squareIndices)

	gl.UseProgram(program)
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	window.SetKeyCallback(key_manager.KeyCallBack)

	bindTexture()

	for !window.ShouldClose() {
		t := time.Now()

		glfw.PollEvents()

		gl.Clear(gl.COLOR_BUFFER_BIT)

		if key_manager.UpdateColor {
			updColor(program)
		}

		draw(vao, window, program)
		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
	gl.DeleteProgram(program)
}

func bindTexture() {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	var width, height int32
	img, err := getImageFromFilePath("./square.png")
	if err != nil {
		fmt.Println(err)
	}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, width, height, 0, gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(img))
	gl.GenerateTextureMipmap(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func getImageFromFilePath(filePath string) ([]byte, error) {
	buf := new(bytes.Buffer)

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := png.Decode(f)
	err = png.Encode(buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func updColor(program uint32) {
	gl.UseProgram(program)
	t := glfw.GetTime()
	redValue := math.Abs(math.Cos(t))
	greenValue := math.Abs(math.Cos(t + 1))
	blueValue := math.Abs(math.Cos(t + 2))
	// fmt.Printf("%.2f %.2f %.2f\n", redValue, greenValue, blueValue)
	vertexColorLocation := gl.GetUniformLocation(program, gl.Str("ourColor\x00"))
	gl.Uniform4f(vertexColorLocation, float32(redValue), float32(greenValue), float32(blueValue), 1.0)
}

func draw(vao uint32, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.UseProgram(program)

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

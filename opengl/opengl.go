package opengl

import (
	"fmt"
	"log"

	"opengl/opengl/shader"

	"github.com/go-gl/gl/v4.6-core/gl"
)

var (
	shaders = []Shader{
		{gl.FRAGMENT_SHADER, "shader.frag", 0},
		{gl.VERTEX_SHADER, "shader.vert", 0},
	}
)

type Shader struct {
	Type     uint32
	File     string
	ObjectId uint32
}

// InitOpenGL initializes OpenGL and returns an intiialized program.
func InitOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	log.Println("OpenGL version", gl.GoStr(gl.GetString(gl.VERSION)))

	program := gl.CreateProgram()

	for _, s := range shaders {
		compiledShader, err := shader.CompileShader(s.File, s.Type)
		if err != nil {
			fmt.Println(err)
		}
		gl.AttachShader(program, compiledShader)
		s.ObjectId = compiledShader
	}

	gl.LinkProgram(program)
	gl.ValidateProgram(program)

	for _, s := range shaders {
		gl.DeleteShader(s.ObjectId)
	}

	return program
}

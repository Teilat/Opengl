package shader

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// COMPUTE_SHADER
// FRAGMENT_SHADER
// GEOMETRY_SHADER
// MESH_SHADER_NV
// TASK_SHADER_NV
// TESS_CONTROL_SHADER
// TESS_EVALUATION_SHADER
// VERTEX_SHADER

func CompileShader(filename string, shaderType uint32) (uint32, error) {
	fmt.Printf("compiling shader %s... ", filename)
	t := time.Now()

	data, err := os.ReadFile("./opengl/shader/" + filename)
	if err != nil {
		return 0, err
	}
	// shaderLen := int32(len(string(data) + "\x00"))
	cSources, free := gl.Strs(string(data) + "\x00")

	shader := gl.CreateShader(shaderType)
	gl.ShaderSource(shader, 1, cSources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		gl.DeleteShader(shader)

		return 0, fmt.Errorf("failed to compile %v: %v", filename, log)
	}

	fmt.Printf("Done in %f milliseconds\n", time.Now().Sub(t).Seconds()*1000)

	return shader, nil
}

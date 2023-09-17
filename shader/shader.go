package shader

import (
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	"os"
	"strings"
	"time"
)

func CompileShader(filename string, shaderType uint32) (uint32, error) {
	fmt.Printf("compiling shader %s... ", filename)
	t := time.Now()

	data, err := os.ReadFile("./shader/" + filename)
	if err != nil {
		return 0, err
	}
	csources, free := gl.Strs(string(data))

	shader := gl.CreateShader(shaderType)
	gl.ShaderSource(shader, 1, csources, nil)
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

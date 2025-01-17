package camera

import (
	"context"
	"math"
	"opengl/window/input"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

const (
	cameraMatrix     = "cameraMatrix\x00"
	cameraData       = "cameraData\x00"
	projectionMatrix = "projectionMatrix\x00"

	movementMulti = 0.05
)

type Camera struct {
	ctx context.Context

	program uint32

	ShaderCameraMatrix       int32
	ShaderCameraData         int32
	ShaderProjectionLocation int32

	pos    mgl32.Vec3
	lookAt mgl32.Vec3
	up     mgl32.Vec3

	fov         float32
	sensitivity float64

	windowWidth  float32
	windowHeight float32

	debug Debug
}

func NewCamera(ctx context.Context, ticker *time.Ticker, program uint32, fov float32, pos, lookAt mgl32.Vec3, width, height int) *Camera {
	c := &Camera{
		ctx: ctx,

		program: program,

		ShaderCameraMatrix:       gl.GetUniformLocation(program, gl.Str(cameraMatrix)),
		ShaderCameraData:         gl.GetUniformLocation(program, gl.Str(cameraData)),
		ShaderProjectionLocation: gl.GetUniformLocation(program, gl.Str(projectionMatrix)),

		pos:    pos,
		lookAt: lookAt, // TODO сделать спавн камеры на желаемые координаты переданые в lookAt
		up:     mgl32.Vec3{0, 1, 0},

		fov:         fov,
		sensitivity: 0.1,

		windowHeight: float32(height),
		windowWidth:  float32(width),
	}

	gl.UseProgram(c.program)

	gl.UniformMatrix4fv(c.ShaderCameraMatrix, 1, false, c.getCameraMatrix4fv())
	gl.UniformMatrix3fv(c.ShaderCameraData, 1, false, c.GetDataP())
	gl.UniformMatrix4fv(c.ShaderProjectionLocation, 1, false, c.getPerspectiveMatrix4fv())
	go c.fixedUpdate(ticker)
	return c
}

// global methods

func (c *Camera) AddDebug(ticker *time.Ticker) {
	c.debug = NewCameraDebug(c, ticker)
}

func (c *Camera) StartDebug() {
	c.debug.Start(c.ctx)
}

func (c *Camera) StopDebug() {
	c.debug.Stop()
}

func (c *Camera) Update() {
	gl.UseProgram(c.program)
	gl.UniformMatrix4fv(c.ShaderCameraMatrix, 1, false, c.getCameraMatrix4fv())
	gl.UniformMatrix3fv(c.ShaderCameraData, 1, false, c.GetDataP())
}

func (c *Camera) UpdateWindow(width, height float32) {
	gl.UseProgram(c.program)
	c.windowWidth, c.windowHeight = width, height
	gl.UniformMatrix4fv(c.ShaderProjectionLocation, 1, false, c.getPerspectiveMatrix4fv())
}

// getters

func (c *Camera) GetPos() mgl32.Vec3 {
	return c.pos
}

func (c *Camera) GetLookAt() mgl32.Vec3 {
	return c.lookAt
}

func (c *Camera) GetDataP() *float32 {
	// pos.xyz,lookAt.xyz,up.xyz
	data := mgl32.Mat3FromCols(c.pos, c.lookAt, c.up)
	return &data[0]
}

func (c *Camera) GetUp() mgl32.Vec3 {
	return c.up
}

func (c *Camera) GetFov() float32 {
	return c.fov
}

func (c *Camera) GetDebug() Debug {
	return c.debug
}

// setters

func (c *Camera) Move(move mgl32.Vec3) {
	c.pos = c.pos.Add(move)
}

func (c *Camera) SetPos(pos mgl32.Vec3) {
	c.pos = pos
}

// local methods

func (c *Camera) fixedUpdate(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			c.calcLookAt()
			c.calcMovement()
			c.updFov()
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Camera) getCameraMatrix4fv() *float32 {
	val := mgl32.LookAtV(c.pos, c.lookAt.Add(c.pos), c.up)
	return &val[0]
}

func (c *Camera) getPerspectiveMatrix4fv() *float32 {
	val := mgl32.Perspective(mgl32.DegToRad(c.fov), c.windowWidth/c.windowHeight, 0.01, 100.0)
	return &val[0]
}

func (c *Camera) calcLookAt() {
	angleX := math.Mod(input.GetDefaultAxis(input.MouseX)*c.sensitivity, 360) // yaw
	angleY := input.GetDefaultAxis(input.MouseY) * c.sensitivity              // pitch

	if angleY > 89 {
		angleY = 89
	}
	if angleY < -89 {
		angleY = -89
	}

	c.lookAt = mgl32.Vec3{
		float32(math.Cos(mgl64.DegToRad(angleY)) * math.Cos(mgl64.DegToRad(angleX))),
		float32(math.Sin(mgl64.DegToRad(angleY))),
		float32(math.Cos(mgl64.DegToRad(angleY)) * math.Sin(mgl64.DegToRad(angleX))),
	}.Normalize()
}

func (c *Camera) updFov() {
	fov := input.GetAxis(glfw.KeyO, glfw.KeyP)
	if c.fov >= 120 && fov > 0 {
		return
	}
	if c.fov <= 30 && fov < 0 {
		return
	}
	c.fov += float32(fov)
	c.UpdateWindow(c.windowWidth, c.windowHeight)
}

func (c *Camera) calcMovement() {
	movementX := input.GetDefaultAxis(input.Horizontal)
	movementY := input.GetAxis(glfw.KeySpace, glfw.KeyLeftControl)
	movementZ := input.GetDefaultAxis(input.Vertical)

	c.Move(c.GetLookAt().Mul(float32(movementZ * movementMulti)))
	c.Move(c.GetLookAt().Cross(mgl32.Vec3{0, 1, 0}).Mul(float32(movementX * movementMulti)))
	c.Move(c.GetUp().Mul(float32(movementY * movementMulti)))
}

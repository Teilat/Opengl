package camera

import (
	"math"

	"opengl/input"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

type Camera struct {
	ShaderCameraLocation     int32
	ShaderProjectionLocation int32

	pos    mgl32.Vec3
	lookAt mgl32.Vec3
	up     mgl32.Vec3

	fov         float32
	sensitivity float64

	windowWidth  float32
	windowHeight float32
}

func NewCamera(program uint32, location, projection string, fov float32, pos mgl32.Vec3, width, height int) *Camera {
	c := &Camera{
		ShaderCameraLocation:     gl.GetUniformLocation(program, gl.Str(location)),
		ShaderProjectionLocation: gl.GetUniformLocation(program, gl.Str(projection)),

		pos:    pos,
		lookAt: mgl32.Vec3{0, 0, 0},
		up:     mgl32.Vec3{0, 1, 0},

		fov:         fov,
		sensitivity: 0.1,

		windowHeight: float32(height),
		windowWidth:  float32(width),
	}

	gl.UniformMatrix4fv(c.ShaderCameraLocation, 1, false, c.getCameraMatrix4fv())
	gl.UniformMatrix4fv(c.ShaderProjectionLocation, 1, false, c.getPerspectiveMatrix4fv())
	return c
}

func (c *Camera) Update() {
	c.calcLookAt()
	c.calcMovement()

	gl.UniformMatrix4fv(c.ShaderCameraLocation, 1, false, c.getCameraMatrix4fv())
}

func (c *Camera) calcLookAt() {
	angleX := math.Mod(input.GetAxis(input.MouseX)*c.sensitivity, 360) // yaw
	angleY := input.GetAxis(input.MouseY) * c.sensitivity              // pitch

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

func (c *Camera) calcMovement() {
	movementX := input.GetAxis(input.Horizontal)
	movementZ := input.GetAxis(input.Vertical)

	c.Move(c.GetLookAt().Mul(float32(movementZ * 0.2)))
	c.Move(c.GetLookAt().Cross(mgl32.Vec3{0, 1, 0}).Mul(float32(movementX * 0.2)))

	if input.GetKey(glfw.KeyLeftControl) {
		c.Move(c.GetUp().Mul(-0.2))
	}
	if input.GetKey(glfw.KeySpace) {
		c.Move(c.GetUp().Mul(0.2))
	}
}

func (c *Camera) GetPos() mgl32.Vec3 {
	return c.pos
}

func (c *Camera) GetLookAt() mgl32.Vec3 {
	return c.lookAt
}

func (c *Camera) GetUp() mgl32.Vec3 {
	return c.up
}

func (c *Camera) SetPos(pos mgl32.Vec3) {
	c.pos = pos
}

func (c *Camera) UpdateWindow(width, height int) {
	c.windowWidth, c.windowHeight = float32(width), float32(height)
	gl.UniformMatrix4fv(c.ShaderProjectionLocation, 1, false, c.getPerspectiveMatrix4fv())
}

func (c *Camera) Move(move mgl32.Vec3) {
	c.pos = c.pos.Add(move)
}

func (c *Camera) getCameraMatrix4fv() *float32 {
	val := mgl32.LookAtV(c.pos, c.lookAt.Add(c.pos), c.up)
	return &val[0]
}

func (c *Camera) getPerspectiveMatrix4fv() *float32 {
	val := mgl32.Perspective(mgl32.DegToRad(c.fov), c.windowWidth/c.windowHeight, 0.1, 100.0)
	return &val[0]
}

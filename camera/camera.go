package camera

import (
	"github.com/go-gl/mathgl/mgl32"
	"opengl/input"
)

type Camera struct {
	pos            mgl32.Vec3
	lookAt         mgl32.Vec3
	ShaderLocation int32
	Fov            float32
	AngleX         float32 // yaw
	AngleY         float32 // pitch
}

func NewCamera(location int32, pos mgl32.Vec3) *Camera {
	return &Camera{
		pos:            pos,
		ShaderLocation: location,
		lookAt:         mgl32.Vec3{0, 0, 1},
	}
}

func (c *Camera) CalcLookAt() {
	offsetX := input.GetAxis(input.MouseX)
	offsetY := input.GetAxis(input.MouseY)

	c.AngleX += offsetX
	c.AngleY += offsetY

	if c.AngleY > 89 {
		c.AngleY = 89
	}
	if c.AngleY < -89 {
		c.AngleY = -89
	}

	c.lookAt = mgl32.SphericalToCartesian(1, mgl32.DegToRad(c.AngleX), mgl32.DegToRad(c.AngleY)).Normalize()
}

func (c *Camera) GetPos() mgl32.Vec3 {
	return c.pos
}

func (c *Camera) GetLookAt() mgl32.Vec3 {
	return c.lookAt
}

func (c *Camera) SetPos(pos mgl32.Vec3) {
	c.pos = pos
}

func (c *Camera) Move(move mgl32.Vec3) {
	c.pos = c.pos.Add(move)
}

func (c *Camera) GetMatrix4fv() *float32 {
	val := mgl32.LookAtV(c.pos, c.lookAt.Add(c.pos), mgl32.Vec3{0, 1, 0})
	return &val[0]
}

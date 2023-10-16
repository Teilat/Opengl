package camera

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"math"
	"opengl/input"
)

type Camera struct {
	pos            mgl32.Vec3
	lookAt         mgl32.Vec3
	up             mgl32.Vec3
	right          mgl32.Vec3
	ShaderLocation int32
	Fov            float32
	AngleX         float64 // yaw
	AngleY         float64 // pitch
}

func NewCamera(location int32, pos mgl32.Vec3) *Camera {
	return &Camera{
		pos:            pos,
		ShaderLocation: location,
		lookAt:         mgl32.Vec3{0, 0, 0},
		up:             mgl32.Vec3{0, 1, 0},
	}
}

func (c *Camera) CalcLookAt() {
	c.AngleX = input.GetAxis(input.MouseX)
	c.AngleY = input.GetAxis(input.MouseY)

	if c.AngleY > 89 {
		c.AngleY = 89
	}
	if c.AngleY < -89 {
		c.AngleY = -89
	}

	c.lookAt = mgl32.Vec3{
		float32(math.Cos(mgl64.DegToRad(c.AngleY)) * math.Cos(mgl64.DegToRad(c.AngleX))),
		float32(math.Sin(mgl64.DegToRad(c.AngleY))),
		float32(math.Cos(mgl64.DegToRad(c.AngleY)) * math.Sin(mgl64.DegToRad(c.AngleX))),
	}.Normalize()
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

func (c *Camera) Move(move mgl32.Vec3) {
	c.pos = c.pos.Add(move)
}

func (c *Camera) GetMatrix4fv() *float32 {
	val := mgl32.LookAtV(c.pos, c.lookAt.Add(c.pos), c.up)
	return &val[0]
}

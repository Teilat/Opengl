package camera

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	pos            mgl32.Vec3
	lookAt         mgl32.Vec3
	ShaderLocation int32
	Fov            float32
	AngleX         float32
	AngleY         float32
}

func NewCamera(location int32, pos mgl32.Vec3) *Camera {
	return &Camera{pos: pos, ShaderLocation: location, lookAt: mgl32.Vec3{0, 0, -1}}
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

package camera

import "github.com/go-gl/mathgl/mgl32"

type Camera struct {
	pos            mgl32.Vec3
	ShaderLocation int32
}

func NewCamera(location int32, pos mgl32.Vec3) *Camera {
	return &Camera{pos: pos, ShaderLocation: location}
}

func (c *Camera) GetPos() mgl32.Vec3 {
	return c.pos
}

func (c *Camera) GetMatrix4() mgl32.Mat4 {
	return mgl32.LookAtV(c.pos, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
}

func (c *Camera) GetMatrix4fv() *float32 {
	val := mgl32.LookAtV(c.pos, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	return &val[0]
}

package object

import "github.com/go-gl/mathgl/mgl32"

func toVec3f32(in [3]float64) mgl32.Vec3 {
	return mgl32.Vec3{float32(in[0]), float32(in[1]), float32(in[2])}
}

func toVec4f32(in [4]float64) mgl32.Vec4 {
	return mgl32.Vec4{float32(in[0]), float32(in[1]), float32(in[2]), float32(in[3])}
}

func toMat4f32(in [16]float64) mgl32.Mat4 {
	return mgl32.Mat4{
		float32(in[0]), float32(in[1]), float32(in[2]), float32(in[3]),
		float32(in[4]), float32(in[5]), float32(in[6]), float32(in[7]),
		float32(in[8]), float32(in[9]), float32(in[10]), float32(in[11]),
		float32(in[12]), float32(in[13]), float32(in[14]), float32(in[15]),
	}
}

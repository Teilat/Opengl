package support

import "math"

func DegToRad(angle float64) float64 {
	return angle * math.Pi / 180
}

func RadToDeg(angle float64) float64 {
	return angle * 180 / math.Pi
}

func BoolToFloat(a bool) float64 {
	if a {
		return 1
	}
	return 0
}

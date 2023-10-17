package support

func BoolToFloat(a bool) float64 {
	if a {
		return 1
	}
	return 0
}

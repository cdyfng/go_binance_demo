package util

import "math"

func Roundx(val float64, x float64) float64 {
	return math.Round(val*x) / x
}

func Round8(val float64) float64 {
	return Roundx(val, 1/0.00000001)
}

func Times8(val float64) int64 {
	return int64(math.Round(val * (1 / 0.00000001)))
}

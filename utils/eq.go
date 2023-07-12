package utils

import "math"

const tolerance = 1e-9

func IsEqTolerance(a, b float64) bool {
	return math.Abs(a-b) <= tolerance
}

package maths

import "math"

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func AbsInt(a int) int {
	return int(math.Abs(float64(a)))
}

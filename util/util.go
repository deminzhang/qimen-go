package util

import (
	"math"
	"slices"
)

func CalRadiansPos[T Numeric](cx, cy, r, angleDegrees T) (x, y T) {
	rad := float64(angleDegrees) * math.Pi / 180
	x = T(float64(cx) + float64(r)*math.Cos(rad))
	y = T(float64(cy) + float64(r)*math.Sin(rad))
	return
}

func Contains(all []string, zhi ...string) bool {
	for _, z := range zhi {
		if slices.Contains(all, z) {
			return true
		}
	}
	return false
}

func If[T any](b bool, t, f T) T {
	if b {
		return t
	}
	return f
}

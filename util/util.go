package util

import (
	"math"
	"slices"
)

func CalRadiansPos(cx, cy, r, angleDegrees float64) (x, y float64) {
	rad := angleDegrees * math.Pi / 180
	x = cx + r*math.Cos(rad)
	y = cy + r*math.Sin(rad)
	return
}

func CalRadiansPosT[T Numeric](cx, cy, r, angleDegrees T) (x, y T) {
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

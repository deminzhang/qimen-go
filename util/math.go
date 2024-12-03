package util

import (
	"golang.org/x/exp/constraints"
	"math"
)

type Numeric interface {
	constraints.Integer | constraints.Float
}

func CalRadiansPos[T Numeric](cx, cy, r, angleDegrees T) (x, y T) {
	rad := float64(angleDegrees) * math.Pi / 180
	x = T(float64(cx) + float64(r)*math.Cos(rad))
	y = T(float64(cy) + float64(r)*math.Sin(rad))
	return
}

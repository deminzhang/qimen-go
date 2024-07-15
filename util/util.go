package util

import (
	"math"
)

type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64
}

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

//gofa.Anpm()

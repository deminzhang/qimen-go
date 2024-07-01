package world

import (
	//"github.com/hebl/gofa"
	"math"
)

func calRadiansPos(cx, cy, r, angleDegrees float64) (x, y float64) {
	rad := angleDegrees * math.Pi / 180
	x = cx + r*math.Cos(rad)
	y = cy + r*math.Sin(rad)
	return
}

//gofa.Anpm()

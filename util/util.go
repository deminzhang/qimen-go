package util

import (
	"math"
	"slices"
	"unsafe"
)

func UnsafeStr2Bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func UnsafeBytes2Str(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
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

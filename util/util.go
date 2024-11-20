package util

import (
	"slices"
)

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

package util

import (
	"os"
	"slices"
	"strings"
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

// Args2Map os.Args除[0]外,以=分切成kv对,无=的v为"true"
func Args2Map() map[string]string {
	m := make(map[string]string)
	for i, s := range os.Args {
		if i == 0 {
			continue
		}
		ss := strings.Split(s, "=")
		if len(ss) > 1 {
			m[ss[0]] = ss[1]
		} else {
			m[ss[0]] = "true"
		}
	}
	return m
}

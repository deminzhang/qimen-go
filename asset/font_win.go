//go:build windows

package asset

import "os"

var DefaultUIFontPath = "C:/Windows/Fonts/simfang.ttf"

func init() {
	sr := os.Getenv("SystemRoot")
	if sr == "" {
		panic("SystemRoot not found")
	}
	DefaultUIFontPath = sr + "/Fonts/simfang.ttf"
}

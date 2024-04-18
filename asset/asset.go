package asset

import (
	"bytes"
	"embed"
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

//go:embed **/*
var FS embed.FS

// Simple icon loading util
func LoadImage(path string) (image.Image, error) {
	iconBytes, err := FS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading icon file: %w", err)
	}

	icon, _, err := image.Decode(bytes.NewReader(iconBytes))
	if err != nil {
		return nil, fmt.Errorf("decoding icon file: %w", err)
	}

	return icon, nil
}

func LoadFont(path string, options *opentype.FaceOptions) (font.Face, error) {
	bytes, err := FS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("loading font file: %s", err)
	}

	fontfile, err := opentype.Parse(bytes)
	if err != nil {
		return nil, fmt.Errorf("parsing font file: %s", err)
	}

	if options == nil {
		options = &opentype.FaceOptions{
			Size:    11,
			DPI:     72,
			Hinting: font.HintingNone,
		}
	}
	face, err := opentype.NewFace(fontfile, options)
	if err != nil {
		return nil, err
	}

	return face, nil
}

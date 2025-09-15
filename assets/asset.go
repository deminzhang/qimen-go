package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/deminzhang/go-common/asset"
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

// LoadFontFS loads a font from the embedded filesystem.
func LoadFontFS(path string) (*opentype.Font, error) {
	bt, err := FS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("loading font file: %s", err)
	}

	ff, err := opentype.Parse(bt)
	if err != nil {
		return nil, fmt.Errorf("parsing font file: %s", err)
	}

	return ff, nil
}

var (
	loadFonts      = map[string]*opentype.Font{}
	loadFontFaces  = map[*opentype.Font]map[float64]font.Face{}
	loadXFontFaces = map[font.Face]*text.GoXFace{}
)

func LoadFont(path string, orDefault bool) (*opentype.Font, error) {
	var bytes []byte
	var err error
	var ff *opentype.Font
	if bytes, err = os.ReadFile(path); err == nil {
		ff, err = opentype.Parse(bytes)
	}
	if err != nil {
		return nil, err
	}
	loadFonts[path] = ff
	return ff, err
}

func GetFontFace(ft *opentype.Font, size float64) (font.Face, error) {
	fs, ok := loadFontFaces[ft]
	if !ok {
		fs = map[float64]font.Face{}
		loadFontFaces[ft] = fs
	}
	if f, ok := fs[size]; ok {
		return f, nil
	}
	options := &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingNone}
	f, err := opentype.NewFace(ft, options)
	if err != nil {
		return nil, err
	}
	fs[size] = f
	return f, nil
}

func LoadDefaultFont() (*opentype.Font, error) {
	ft, ok := loadFonts[asset.DefaultUIFontPath]
	if ok {
		return ft, nil
	}
	return LoadFont(asset.DefaultUIFontPath, true)
}

func GetDefaultFontFace(size float64) (font.Face, error) {
	ft, ok := loadFonts[asset.DefaultUIFontPath]
	var err error
	if !ok {
		ft, err = LoadDefaultFont()
		if err != nil {
			return nil, err
		}
		loadFonts[asset.DefaultUIFontPath] = ft
	}
	return GetFontFace(ft, size)
}

func GetDefaultFontXFace(size float64) (*text.GoXFace, error) {
	ft, ok := loadFonts[asset.DefaultUIFontPath]
	var err error
	if !ok {
		ft, err = LoadDefaultFont()
		if err != nil {
			return nil, err
		}
		loadFonts[asset.DefaultUIFontPath] = ft
	}
	ff, err := GetFontFace(ft, size)
	if err != nil {
		return nil, err
	}
	xf := loadXFontFaces[ff]
	if xf == nil {
		xf = text.NewGoXFace(ff)
		loadXFontFaces[ff] = xf
	}
	return xf, err
}

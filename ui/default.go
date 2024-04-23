package ui

import (
	"bytes"
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"image"
	_ "image/png"
	"log"
)

//go:embed ui.png
var FS embed.FS
var (
	uiImage       *ebiten.Image
	uiFont        font.Face
	uiFontMHeight int
	uiFontMWidth  int
)

func init() {
	iconBytes, err := FS.ReadFile("ui.png")
	if err != nil {
		log.Fatal("reading icon file: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(iconBytes))
	if err != nil {
		log.Fatal("decoding icon file: %w", err)
	}
	if err != nil {
		log.Fatal(err)
	}
	uiImage = ebiten.NewImageFromImage(img)

	//不支持汉字
	tt, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	uiFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	//支持汉字
	//uiFont, err = asset.LoadFont("font/lana_pixel.ttf", &opentype.FaceOptions{
	//	Size:    14,
	//	DPI:     72,
	//	Hinting: font.HintingFull,
	//})
	if err != nil {
		log.Fatal(err)
	}
	b, _, _ := uiFont.GlyphBounds('M')
	uiFontMHeight = (b.Max.Y - b.Min.Y).Ceil()
	uiFontMWidth = (b.Max.X - b.Min.X).Ceil()
}

type imageType int

const (
	imageTypeButton imageType = iota
	imageTypeButtonPressed
	imageTypeTextBox
	imageTypeScrollBarBack
	imageTypeScrollBarFront
	imageTypeCheckBox
	imageTypeCheckBoxPressed
	imageTypeCheckBoxMark
	imageTypeOptionBox
	imageTypeOptionBoxPressed
	imageTypeOptionBoxMark
)

var imageSrcRects = map[imageType]image.Rectangle{
	imageTypeButton:           image.Rect(0, 0, 16, 16),
	imageTypeButtonPressed:    image.Rect(16, 0, 32, 16),
	imageTypeTextBox:          image.Rect(0, 16, 16, 32),
	imageTypeScrollBarBack:    image.Rect(16, 16, 24, 32),
	imageTypeScrollBarFront:   image.Rect(24, 16, 32, 32),
	imageTypeCheckBox:         image.Rect(0, 32, 16, 48),
	imageTypeCheckBoxPressed:  image.Rect(16, 32, 32, 48),
	imageTypeCheckBoxMark:     image.Rect(32, 32, 48, 48),
	imageTypeOptionBox:        image.Rect(0, 48, 16, 64),
	imageTypeOptionBoxPressed: image.Rect(16, 48, 32, 64),
	imageTypeOptionBoxMark:    image.Rect(32, 48, 48, 64),
}

func GetDefaultUIImage() *ebiten.Image {
	return uiImage
}

func GetDefaultUIFont() font.Face {
	return uiFont
}

func SetDefaultUIImage(img *ebiten.Image) {
	uiImage = img
}

func SetDefaultUIFont(f font.Face) {
	uiFont = f
	b, _, _ := uiFont.GlyphBounds('M')
	uiFontMHeight = (b.Max.Y - b.Min.Y).Ceil()
}

package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"log"
	"os"
	"qimen/asset"
	"qimen/ui"
)

var (
	Font      *opentype.Font
	FontFaces map[float64]font.Face
	ThisGame  *game
)

func init() {
	var f font.Face
	var bytes []byte
	var err error
	var ff *opentype.Font
	if bytes, err = os.ReadFile(defaultUIFontPath); err == nil {
		ff, err = opentype.Parse(bytes)
	}
	if err != nil {
		log.Println("parse default font error:", err)
		ff, err = asset.LoadFont("font/lana_pixel.ttf") //字不全,如癸显不出来
	}
	if err != nil {
		log.Fatal(err)
	}
	Font = ff
	f, err = GetFontFace(14)
	ui.SetDefaultUIFont(f)
}

func GetFontFace(size float64) (font.Face, error) {
	if FontFaces == nil {
		FontFaces = make(map[float64]font.Face)
	}
	if f, ok := FontFaces[size]; ok {
		return f, nil
	}
	options := &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingNone,
	}
	f, err := opentype.NewFace(Font, options)
	if err != nil {
		return nil, err
	}
	FontFaces[size] = f
	return f, nil
}

func setWindow() {
	icon16, err := asset.LoadImage("images/icon_16x16.png")
	if err != nil {
		log.Fatal("loading icon_16: %w", err)
	}
	icon32, err := asset.LoadImage("images/icon_32x32.png")
	if err != nil {
		log.Fatal("loading icon_32: %w", err)
	}
	ebiten.SetWindowIcon([]image.Image{icon32, icon16})
	ebiten.SetTPS(TPSRate)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("众妙之门")
}

func NewWorld() ebiten.Game {
	setWindow()
	g := NewGame()
	ThisGame = g
	return g
}

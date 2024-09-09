package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"image"
	"log"
	"qimen/asset"
	"qimen/ui"
)

var (
	ThisGame *game
)

func init() {
	face, err := asset.GetDefaultFontFace(14)
	if err != nil {
		return
	}
	ui.SetDefaultUIFont(face)
}

func GetFontFace(size float64) (font.Face, error) {
	return asset.GetDefaultFontFace(size)
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

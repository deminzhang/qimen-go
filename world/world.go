package world

import (
	"image"
	"image/color"
	"log"

	"github.com/deminzhang/qimen-go/asset"
	"github.com/deminzhang/qimen-go/gui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
)

var (
	ThisGame *Game
)

func init() {
	face, err := asset.GetDefaultFontFace(14)
	if err != nil {
		return
	}
	gui.SetDefaultUIFont(face)
}

func GetFontFace(size float64) (font.Face, error) {
	return asset.GetDefaultFontFace(size)
}

func GetFontXFace(size float64) (*text.GoXFace, error) {
	return asset.GetDefaultFontXFace(size)
}

//	func TextDrawV1(dst *ebiten.Image, text string, face font.Face, x, y int, clr color.Color){
//			"github.com/hajimehoshi/ebiten/v2/text".Draw(dst, text, face, x, y, clr)
//		}
func TextDrawV2(dst *ebiten.Image, txt string, xface *text.GoXFace, x, y int, clr color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(dst, txt, xface, op)
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
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("众妙之门")
}

func NewWorld() ebiten.Game {
	setWindow()
	g := NewGame()
	ThisGame = g
	return g
}

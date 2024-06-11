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
	"runtime"
)

func init() {
	options := &opentype.FaceOptions{
		Size:    14,
		DPI:     72,
		Hinting: font.HintingNone,
	}
	var f font.Face
	var err error
	if runtime.GOOS == "windows" {
		bytes, err := os.ReadFile("C:/Windows/Fonts/simfang.ttf")
		if err != nil {
			log.Fatal(err)
		}
		ff, err := opentype.Parse(bytes)
		if err != nil {
			log.Fatal(err)
		}
		f, err = opentype.NewFace(ff, options)
	} else {
		f, err = asset.LoadFont("font/lana_pixel.ttf", options) //字不全,如癸显不出来
	}
	if err != nil {
		log.Fatal(err)
	}
	ui.SetDefaultUIFont(f)

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
	//ebiten.SetWindowTitle("众妙之门")
	ebiten.SetWindowTitle("奇门遁甲")
}

func NewWorld() ebiten.Game {
	g := NewGame0()
	//g := NewGame1()
	UIShowQiMen(ScreenWidth, ScreenHeight)
	return g
}

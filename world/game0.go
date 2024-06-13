package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"qimen/ui"
)

type game0 struct {
}

func (g *game0) Update() error {
	ui.Update()
	return nil
}

func (g *game0) Draw(screen *ebiten.Image) {
	ui.Draw(screen)
}

func (g *game0) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame0() *game0 {
	g := &game0{}
	UIShowQiMen(screenWidth, screenHeight)
	return g
}

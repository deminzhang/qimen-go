package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"qimen/ebiten_ui"
)

func (g *game) Draw(screen *ebiten.Image) {
	ebiten_ui.Draw(screen)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

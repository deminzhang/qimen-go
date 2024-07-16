package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"qimen/ui"
)

type game1 struct {
	count     int
	stars     *StarEffect
	astrolabe *Astrolabe
	qimen     *QMShow
}

func (g *game1) Update() error {
	g.count++
	g.count %= 360
	ui.Update()
	g.stars.Update()
	g.qimen.Update()
	g.astrolabe.Update()
	return nil
}

func (g *game1) Draw(screen *ebiten.Image) {
	g.stars.Draw(screen)
	ui.Draw(screen)
	g.astrolabe.Draw(screen)
	g.qimen.Draw(screen)
}

func (g *game1) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
	//return screenWidth, screenHeight
}

func NewGame1() *game1 {
	g := &game1{
		astrolabe: NewAstrolabe(770, 450),
		qimen:     NewQimenShow(260, 460),
		stars:     NewStarEffect(260, 460),
	}
	u := UIShowQiMen()
	u.hide9GongUI()
	u.noShow12Gong()
	return g
}

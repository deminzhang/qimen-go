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
	bazi      *EightCharPan
}

func (g *game1) Update() error {
	g.count++
	g.count %= 360
	ui.Update()
	g.stars.Update()
	g.qimen.Update()
	if uiQiMen.IsShowBaZi() {
		g.bazi.Update()
	} else {
		g.astrolabe.Update()
	}
	return nil
}

func (g *game1) Draw(screen *ebiten.Image) {
	g.stars.Draw(screen)
	ui.Draw(screen)
	if uiQiMen.IsShowBaZi() {
		g.bazi.Draw(screen)
	} else {
		g.astrolabe.Draw(screen)
	}
	g.qimen.Draw(screen)
}

func (g *game1) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
	//return screenWidth, screenHeight
}

func NewGame1() *game1 {
	g := &game1{
		stars:     NewStarEffect(260, 460),
		qimen:     NewQiMenShow(260, 460),
		astrolabe: NewAstrolabe(770, 450),
		bazi:      NewEightCharPan(770, 450),
	}
	u := UIShowQiMen()
	u.hide9GongUI()
	u.noShow12Gong()
	return g
}

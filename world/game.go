package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"qimen/qimen"
	"qimen/ui"
)

type game struct {
	count     int
	stars     *StarEffect
	astrolabe *Astrolabe
	qiMen     *QMShow
	baZi      *EightCharPan
	qmGame    *qimen.QMGame
}

func (g *game) Update() error {
	g.count++
	g.count %= 360
	ui.Update()
	g.stars.Update()
	g.qiMen.Update()
	g.baZi.Update()
	g.astrolabe.Update()
	return nil
}

func (g *game) Draw(dst *ebiten.Image) {
	g.stars.Draw(dst)
	g.baZi.Draw(dst)
	g.astrolabe.Draw(dst)
	g.qiMen.Draw(dst)
	ui.Draw(dst)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
	//return screenWidth, screenHeight
}

func NewGame() *game {
	u := UIShowQiMen()
	g := &game{
		stars:     NewStarEffect(260, 460),
		qiMen:     NewQiMenShow(260, 460),
		astrolabe: NewAstrolabe(770+500, 450),
		baZi:      NewEightCharPan(770, 450),
		qmGame:    u.pan,
	}
	return g
}

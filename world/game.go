package world

import (
	"github.com/6tail/lunar-go/calendar"
	"github.com/hajimehoshi/ebiten/v2"
	"qimen/qimen"
	"qimen/ui"
	"time"
)

type game struct {
	count      int
	uiQM       *UIQiMen
	stars      *StarEffect
	astrolabe  *Astrolabe
	qiMen      *QMShow
	baZi       *EightCharPan
	qmGame     *qimen.QMGame
	autoMinute bool
}

func (g *game) Update() error {
	g.count++
	g.count %= 360
	ui.Update()
	g.stars.Update()
	g.qiMen.Update()
	g.baZi.Update()
	g.astrolabe.Update()
	if g.autoMinute {
		if g.count == 0 {
			g.qmGame = g.uiQM.NextApply()
		}
	}
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
	solar := calendar.NewSolarFromDate(time.Now())
	pan := u.Apply(solar)
	g := &game{
		uiQM:      u,
		stars:     NewStarEffect(260, 460),
		qiMen:     NewQiMenShow(260, 460),
		astrolabe: NewAstrolabe(770+500, 450),
		baZi:      NewEightCharPan(522, 204),
		qmGame:    pan,
	}
	return g
}

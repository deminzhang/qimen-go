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
	char8      *Char8Pan
	qmGame     *qimen.QMGame
	autoMinute bool
}

func (g *game) Update() error {
	g.count++
	g.count %= 60
	ui.Update()
	g.qiMen.Update()
	g.char8.Update()
	g.astrolabe.Update()
	g.stars.SetPos(g.astrolabe.GetSolarPos())
	g.stars.Update()
	if g.autoMinute && !g.astrolabe.DataQuerying {
		if g.count == 0 {
			g.qmGame = g.uiQM.NextApply()
		}
	}
	return nil
}

func (g *game) Draw(dst *ebiten.Image) {
	//g.stars.Draw(dst)
	g.char8.Draw(dst)
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
		stars:     NewStarEffect(260, 450),
		qiMen:     NewQiMenShow(260, 450),
		astrolabe: NewAstrolabe(770+500, 450),
		char8:     NewChar8Pan(522, 174),
		qmGame:    pan,
	}
	return g
}

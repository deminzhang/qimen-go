package world

import (
	"github.com/6tail/lunar-go/calendar"
	"github.com/deminzhang/qimen-go/qimen"
	"github.com/deminzhang/qimen-go/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type game struct {
	world      *ebiten.Image
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
	g.qiMen.Update()
	g.char8.Update()
	g.astrolabe.Update()
	//g.stars.SetPos(g.astrolabe.GetSolarPos())
	//g.stars.Update()
	if g.autoMinute && !g.astrolabe.DataQuerying() {
		if g.count%10 == 0 {
			g.qmGame = g.uiQM.NextApply()
		}
	}

	ui.Update()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	g.char8.Draw(screen)
	g.astrolabe.Draw(screen)
	//g.stars.Draw(screen)
	g.qiMen.DrawHead(screen)
	g.qiMen.Draw(screen)
	ui.Draw(screen)
}

func (g *game) Layout(w, h int) (int, int) {
	return w, h
}

func NewGame() *game {
	u := UIShowQiMen()
	solar := calendar.NewSolarFromDate(time.Now())
	pan := u.Apply(solar)
	g := &game{
		uiQM:      u,
		stars:     NewStarEffect(screenWidth/2, 217),
		qiMen:     NewQiMenShow(260, 450),
		astrolabe: NewAstrolabe(770+500, 450),
		char8:     NewChar8Pan(522, 174),
		qmGame:    pan,
	}
	return g
}

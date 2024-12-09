package world

import (
	"fmt"
	"github.com/6tail/lunar-go/calendar"
	"github.com/deminzhang/qimen-go/gui"
	"github.com/deminzhang/qimen-go/qimen"
	"github.com/deminzhang/qimen-go/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"time"
)

var (
	ScreenWidth  = initScreenWidth
	ScreenHeight = initScreenHeight
	Debug        = false
)

type game struct {
	count         int
	uiQM          *UIQiMen
	stars         *StarEffect
	astrolabe     *Astrolabe
	qiMen         *QMShow
	char8         *Char8Pan
	qmGame        *qimen.QMGame
	autoMinute    bool
	showChar8     bool
	showAstrolabe bool

	StrokeManager
}

func (g *game) Update() error {
	g.count++
	g.count %= 60
	//g.stars.Update()
	g.qiMen.Update()
	g.char8.Visible = g.showChar8
	g.char8.Update()
	if g.showAstrolabe {
		g.astrolabe.Update()
	}
	//g.stars.SetPos(g.astrolabe.GetSolarPos())
	//if g.autoMinute && !g.astrolabe.DataQuerying() {
	if g.autoMinute {
		if g.count%10 == 0 {
			//g.qmGame = g.uiQM.NextHour()
			g.qmGame = g.uiQM.NextMinute()
		}
	}
	g.StrokeManager.Update()

	gui.Update()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	if g.showAstrolabe {
		g.astrolabe.Draw(screen)
	}
	//g.stars.Draw(screen)
	g.qiMen.Draw(screen)
	//g.char8.Draw(screen)
	gui.Draw(screen)
	msg := fmt.Sprintf(`FPS: %0.2f, TPS: %0.2f`, ebiten.ActualFPS(), ebiten.ActualTPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (g *game) Layout(w, h int) (int, int) {
	ScreenWidth = w
	ScreenHeight = h
	if g.qiMen != nil {
		g.qiMen.Y = float32(h / 2)
	}
	if g.uiQM != nil {
		g.uiQM.W = w - g.uiQM.X
		g.uiQM.H = h - g.uiQM.Y
	}
	return w, h
}

func NewGame() *game {
	if _, ok := util.Args2Map()["debug"]; ok {
		Debug = true
		gui.SetBorderDebug(true)
		UIShowChat()
	}
	u := UIShowQiMen()
	solar := calendar.NewSolarFromDate(time.Now())
	pan := u.Apply(solar)
	g := &game{
		uiQM:          u,
		stars:         NewStarEffect(float32(ScreenWidth/2), 217),
		qiMen:         NewQiMenShow(450, 500),
		astrolabe:     NewAstrolabe(1650, 450),
		char8:         NewChar8Pan(880, 174),
		qmGame:        pan,
		showChar8:     true,
		showAstrolabe: true,

		StrokeManager: StrokeManager{
			strokes: make(map[*Stroke]struct{}),
		},
	}

	return g
}

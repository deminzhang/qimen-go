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
	Dev          = false
)

type Game struct {
	count     int
	uiQM      *UIQiMen
	stars     *StarEffect
	astrolabe *Astrolabe
	qiMen     *QMShow
	char8     *Char8Pan
	qmGame    *qimen.QMGame
	meiHua    *MeiHua
	big6      *Big6Show

	autoMinute    bool
	showMeiHua    bool
	showQiMen     bool
	showBig6      bool
	showChar8     bool
	showAstrolabe bool

	StrokeManager
}

func (g *Game) Update() error {
	g.count = (g.count + 1) % 60

	g.qiMen.Update()

	g.char8.UI.Visible = g.showChar8
	if g.showChar8 {
		g.char8.Update()
	}
	g.meiHua.UI.Visible = g.showMeiHua
	if g.showMeiHua {
		g.meiHua.Update()
	}
	g.big6.UI.Visible = g.showBig6
	if g.showBig6 {
		g.big6.Update()
	}
	if g.showAstrolabe {
		g.astrolabe.Update()
	}
	//g.stars.Update()
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

func (g *Game) Draw(screen *ebiten.Image) {
	if g.showAstrolabe {
		g.astrolabe.Draw(screen)
	}
	g.qiMen.Draw(screen)
	//g.stars.Draw(screen)
	if g.showMeiHua {
		g.meiHua.Draw(screen)
	}
	if g.showBig6 {
		g.big6.Draw(screen)
	}
	if g.showChar8 {
		g.char8.Draw(screen)
	}
	gui.Draw(screen)
	if Dev {
		msg := fmt.Sprintf(`FPS: %0.2f, TPS: %0.2f`, ebiten.ActualFPS(), ebiten.ActualTPS())
		ebitenutil.DebugPrint(screen, msg)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	if ScreenWidth != w || ScreenHeight != h {
		ScreenWidth = w
		ScreenHeight = h
		if g.qiMen != nil && g.qiMen.Y != float32(h/2) {
			g.qiMen.Y = float32(h / 2)
			g.qiMen.dirty = true
		}
		if g.uiQM != nil {
			g.uiQM.W = w - g.uiQM.X
			g.uiQM.H = h - g.uiQM.Y
		}
		gui.OnLayout(w, h)
	}
	return w, h
}

func NewGame() *Game {
	if _, ok := util.Args2Map()["dev"]; ok {
		Dev = true
		gui.SetBorderDebug(true)
		UIShowChat()
	}
	u := UIShowQiMen()
	solar := calendar.NewSolarFromDate(time.Now())
	pan := u.Apply(solar)
	if pan == nil {
		panic("pan is nil")
	}
	g := &Game{
		uiQM:      u,
		qmGame:    pan,
		stars:     NewStarEffect(float32(ScreenWidth/2), 217),
		qiMen:     NewQiMenShow(450, 500),
		meiHua:    NewMeiHua(1130, 170),
		big6:      NewBig6(880, 170),
		char8:     NewChar8Pan(880, 314),
		astrolabe: NewAstrolabe(1650, 450),

		showQiMen:     true,
		showMeiHua:    true,
		showBig6:      true,
		showChar8:     true,
		showAstrolabe: true,

		StrokeManager: StrokeManager{
			strokes: make(map[*Stroke]struct{}),
		},
	}

	return g
}

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

type Game struct {
	count     int
	uiQM      *UIQiMen
	stars     *StarEffect
	astrolabe *Astrolabe
	qiMen     *QMShow
	char8     *Char8Pan
	qmGame    *qimen.QMGame
	meiHua    *MeiHua

	autoMinute    bool
	showMeiHua    bool
	showChar8     bool
	showAstrolabe bool

	StrokeManager
}

func (g *Game) Update() error {
	g.count++
	g.count %= 60
	//g.stars.Update()
	g.qiMen.Update()
	g.char8.Visible = g.showChar8
	g.char8.Update()
	g.meiHua.Visible = g.showMeiHua
	g.meiHua.Update()
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

func (g *Game) Draw(screen *ebiten.Image) {
	g.qiMen.Draw(screen)
	if g.showAstrolabe {
		g.astrolabe.Draw(screen)
	}
	//g.stars.Draw(screen)
	g.meiHua.Draw(screen)
	g.char8.Draw(screen)
	gui.Draw(screen)
	msg := fmt.Sprintf(`FPS: %0.2f, TPS: %0.2f`, ebiten.ActualFPS(), ebiten.ActualTPS())
	ebitenutil.DebugPrint(screen, msg)
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
	if _, ok := util.Args2Map()["debug"]; ok {
		Debug = true
		gui.SetBorderDebug(true)
		UIShowChat()
	}
	u := UIShowQiMen()
	solar := calendar.NewSolarFromDate(time.Now())
	pan := u.Apply(solar)
	g := &Game{
		uiQM:      u,
		stars:     NewStarEffect(float32(ScreenWidth/2), 217),
		qiMen:     NewQiMenShow(450, 500),
		astrolabe: NewAstrolabe(1650, 450),
		char8:     NewChar8Pan(880, 174),
		qmGame:    pan,
		meiHua:    NewMeiHua(880, 780),

		showMeiHua:    true,
		showChar8:     true,
		showAstrolabe: true,

		StrokeManager: StrokeManager{
			strokes: make(map[*Stroke]struct{}),
		},
	}

	return g
}

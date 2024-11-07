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
	camera     cameraX
	count      int
	uiQM       *UIQiMen
	stars      *StarEffect
	astrolabe  *Astrolabe
	qiMen      *QMShow
	char8      *Char8Pan
	qmGame     *qimen.QMGame
	autoMinute bool
}

type cameraX struct {
	Camera
}

func (c *cameraX) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		c.Position[0] += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		c.Position[0] -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		c.Position[1] += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		c.Position[1] -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		if c.ZoomFactor > -2400 {
			c.ZoomFactor -= 1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		if c.ZoomFactor < 2400 {
			c.ZoomFactor += 1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		c.Rotation += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		c.Reset()
	}
}

func (g *game) Update() error {
	g.count++
	g.count %= 60
	g.camera.Update()
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
	g.char8.Draw(g.world)
	g.astrolabe.Draw(g.world)
	//g.stars.Draw(g.world)
	g.qiMen.DrawHead(screen)
	g.qiMen.Draw(g.world)
	g.camera.Render(g.world, screen)
	ui.Draw(screen)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.world = ebiten.NewImage(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
	//return screenWidth, screenHeight
}

func NewGame() *game {
	u := UIShowQiMen()
	solar := calendar.NewSolarFromDate(time.Now())
	pan := u.Apply(solar)
	g := &game{
		world:     ebiten.NewImage(screenWidth, screenHeight),
		uiQM:      u,
		stars:     NewStarEffect(screenWidth/2, 217),
		qiMen:     NewQiMenShow(260, 450),
		astrolabe: NewAstrolabe(770+500, 450),
		char8:     NewChar8Pan(522, 174),
		qmGame:    pan,
	}
	return g
}

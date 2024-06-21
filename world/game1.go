package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"qimen/ui"
)

type game1 struct {
	count     int
	starLines [starsLineCount]StarLine
	astrolabe *Astrolabe
	qimen     *QMShow
}

func (g *game1) Update() error {
	g.count++
	g.count %= 360
	ui.Update()
	//x, y := ebiten.CursorPosition()
	x, y := centerX, centerY
	//x, y := screenWidth/2, centerY
	for i := 0; i < starsLineCount; i++ {
		g.starLines[i].Update(float32(x*starsLineScale), float32(y*starsLineScale))
	}
	g.qimen.Update()
	g.astrolabe.Update()
	return nil
}

func (g *game1) Draw(screen *ebiten.Image) {
	for i := 0; i < starsLineCount; i++ {
		g.starLines[i].Draw(screen)
	}
	ui.Draw(screen)
	g.qimen.Draw(screen)
	g.astrolabe.Draw(screen)
}

func (g *game1) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame1() *game1 {
	g := &game1{
		astrolabe: NewAstrolabe(),
		qimen:     NewQimenShow(),
	}
	for i := 0; i < starsLineCount; i++ {
		g.starLines[i].Init()
	}
	u := UIShowQiMen(screenWidth, screenHeight)
	u.hide9GongUI()
	u.noShow12Gong()
	return g
}

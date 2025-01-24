package world

import (
	"fmt"
	"github.com/deminzhang/qimen-go/asset"
	"github.com/deminzhang/qimen-go/graphic"
	"github.com/deminzhang/qimen-go/gui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Big6 struct {
	X, Y    int
	Visible bool
	UI      *gui.BaseUI

	Mover     *Sprite
	GuaSprite []*Sprite
}

func NewBig6(x, y int) *Big6 {
	m := &Big6{X: x, Y: y,
		Visible: true,
		UI:      &gui.BaseUI{X: x, Y: y, Visible: true, W: meiHuaUIWidth, H: meiHuaUIHeight, BDColor: colorGray},
	}

	//m.UI.AddChildren()
	gui.ActiveUI(m.UI)
	return m
}

func (m *Big6) Update() {
	m.UI.Visible = m.Visible

	if m.Mover == nil {
		m.Mover = NewSprite(graphic.NewRectImage(10), colorGray)
		m.Mover.onMove = func(sx, sy, dx, dy int) {
			m.X += dx
			m.Y += dy
			m.UI.X, m.UI.Y = m.X, m.Y

		}
		ThisGame.AddSprite(m.Mover)
		m.Mover.MoveTo(m.X, m.Y)
	}

}

func (m *Big6) Draw(dst *ebiten.Image) {
	if !m.Visible {
		return
	}
	m.Mover.Draw(dst)
	ft14, _ := asset.GetDefaultFontFace(14)
	text.Draw(dst, "大六壬", ft14, m.X+16, m.Y+16, colorWhite)
	cx, cy := m.X+16, m.Y+32+16
	b6 := ThisGame.qmGame.Big6
	ke := b6.Ke
	text.Draw(dst, ke[3][1], ft14, cx, cy, colorWhite)
	text.Draw(dst, ke[3][0], ft14, cx, cy+16, colorWhite)
	text.Draw(dst, ke[2][1], ft14, cx+16, cy, colorWhite)
	text.Draw(dst, ke[2][0], ft14, cx+16, cy+16, ColorGanZhi(b6.DayZhi))
	text.Draw(dst, ke[1][1], ft14, cx+32, cy, colorWhite)
	text.Draw(dst, ke[1][0], ft14, cx+32, cy+16, colorWhite)
	text.Draw(dst, ke[0][1], ft14, cx+48, cy, colorWhite)
	text.Draw(dst, ke[0][0], ft14, cx+48, cy+16, ColorGanZhi(b6.DayGan))
	chuan := b6.Chuan
	cx += 96 + 32
	text.Draw(dst, chuan[0], ft14, cx, cy, colorWhite)
	text.Draw(dst, chuan[1], ft14, cx, cy+16, colorWhite)
	text.Draw(dst, chuan[2], ft14, cx, cy+32, colorWhite)
	cx = m.X + 16
	cy += 64
	text.Draw(dst, "课体:"+b6.Ge, ft14, cx, cy, colorWhite)

	cx, cy = m.X+194, m.Y+16
	yz := ThisGame.qmGame.Lunar.GetYearInGanZhiExact()
	mz := ThisGame.qmGame.Lunar.GetMonthInGanZhiExact()
	dz := ThisGame.qmGame.Lunar.GetDayInGanZhiExact()
	tz := ThisGame.qmGame.Lunar.GetTimeInGanZhi()
	jz := ThisGame.qmGame.YueJiang
	text.Draw(dst, fmt.Sprintf("%s年", yz), ft14, cx, cy, colorWhite)
	text.Draw(dst, fmt.Sprintf("%s月", mz), ft14, cx, cy+16, colorWhite)
	text.Draw(dst, fmt.Sprintf("%s日", dz), ft14, cx, cy+32, colorRed)
	text.Draw(dst, fmt.Sprintf("%s时", tz), ft14, cx, cy+48, colorWhite)
	text.Draw(dst, fmt.Sprintf("月将%s", jz), ft14, cx, cy+64, colorWhite)
}

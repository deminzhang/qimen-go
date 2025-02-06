package world

import (
	"fmt"
	"github.com/deminzhang/qimen-go/asset"
	"github.com/deminzhang/qimen-go/graphic"
	"github.com/deminzhang/qimen-go/gui"
	"github.com/deminzhang/qimen-go/qimen"
	"github.com/deminzhang/qimen-go/util"
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
	text.Draw(dst, "大六壬Dev", ft14, m.X+16, m.Y+16, colorWhite)
	cx, cy := m.X+16, m.Y+32+16
	b6 := ThisGame.qmGame.Big6
	ke := b6.Ke4
	text.Draw(dst, ke[3].God, ft14, cx, cy, util.If(ke[3].God == "贵", colorRed, colorWhite))
	text.Draw(dst, ke[3].Up, ft14, cx, cy+16, colorWhite)
	text.Draw(dst, ke[3].Down, ft14, cx, cy+32, colorWhite)
	if qimen.WuXingKe[qimen.GanZhiWuXing[ke[3].Down]] == qimen.GanZhiWuXing[ke[3].Up] {
		text.Draw(dst, "↑", ft14, cx+8, cy+24, colorRed)
	} else if qimen.WuXingKe[qimen.GanZhiWuXing[ke[3].Up]] == qimen.GanZhiWuXing[ke[3].Down] {
		text.Draw(dst, "↓", ft14, cx+8, cy+24, colorRed)
	}
	text.Draw(dst, ke[2].God, ft14, cx+16, cy, util.If(ke[2].God == "贵", colorRed, colorWhite))
	text.Draw(dst, ke[2].Up, ft14, cx+16, cy+16, colorWhite)
	text.Draw(dst, ke[2].Down, ft14, cx+16, cy+32, ColorGanZhi(b6.DayZhi))
	if qimen.WuXingKe[qimen.GanZhiWuXing[ke[2].Down]] == qimen.GanZhiWuXing[ke[2].Up] {
		text.Draw(dst, "↑", ft14, cx+16+8, cy+24, colorRed)
	} else if qimen.WuXingKe[qimen.GanZhiWuXing[ke[2].Up]] == qimen.GanZhiWuXing[ke[2].Down] {
		text.Draw(dst, "↓", ft14, cx+16+8, cy+24, colorRed)
	}
	text.Draw(dst, ke[1].God, ft14, cx+32, cy, util.If(ke[1].God == "贵", colorRed, colorWhite))
	text.Draw(dst, ke[1].Up, ft14, cx+32, cy+16, colorWhite)
	text.Draw(dst, ke[1].Down, ft14, cx+32, cy+32, colorWhite)
	if qimen.WuXingKe[qimen.GanZhiWuXing[ke[1].Down]] == qimen.GanZhiWuXing[ke[1].Up] {
		text.Draw(dst, "↑", ft14, cx+32+8, cy+24, colorRed)
	} else if qimen.WuXingKe[qimen.GanZhiWuXing[ke[1].Up]] == qimen.GanZhiWuXing[ke[1].Down] {
		text.Draw(dst, "↓", ft14, cx+32+8, cy+24, colorRed)
	}
	text.Draw(dst, ke[0].God, ft14, cx+48, cy, util.If(ke[0].God == "贵", colorRed, colorWhite))
	text.Draw(dst, ke[0].Up, ft14, cx+48, cy+16, colorWhite)
	text.Draw(dst, ke[0].Down, ft14, cx+48, cy+32, ColorGanZhi(b6.DayGan))
	if qimen.WuXingKe[qimen.GanZhiWuXing[ke[0].Down]] == qimen.GanZhiWuXing[ke[0].Up] {
		text.Draw(dst, "↑", ft14, cx+48+8, cy+24, colorRed)
	} else if qimen.WuXingKe[qimen.GanZhiWuXing[ke[0].Up]] == qimen.GanZhiWuXing[ke[0].Down] {
		text.Draw(dst, "↓", ft14, cx+48+8, cy+24, colorRed)
	}
	chuan := b6.Chuan
	cx += 96
	text.Draw(dst, b6.Relation6(chuan[0]), ft14, cx, cy, colorWhite)
	text.Draw(dst, b6.Relation6(chuan[1]), ft14, cx, cy+16, colorWhite)
	text.Draw(dst, b6.Relation6(chuan[2]), ft14, cx, cy+32, colorWhite)
	cx += 16
	gc0 := b6.GetGongByJiangZhi(chuan[0])
	gc1 := b6.GetGongByJiangZhi(chuan[1])
	gc2 := b6.GetGongByJiangZhi(chuan[2])
	text.Draw(dst, gc0.JiangGan, ft14, cx, cy, colorWhite)
	text.Draw(dst, gc1.JiangGan, ft14, cx, cy+16, colorWhite)
	text.Draw(dst, gc2.JiangGan, ft14, cx, cy+32, colorWhite)
	cx += 16
	text.Draw(dst, chuan[0], ft14, cx, cy, colorWhite)
	text.Draw(dst, chuan[1], ft14, cx, cy+16, colorWhite)
	text.Draw(dst, chuan[2], ft14, cx, cy+32, colorWhite)
	cx += 16
	text.Draw(dst, gc0.Jiang12, ft14, cx, cy, colorGray)
	text.Draw(dst, gc1.Jiang12, ft14, cx, cy+16, colorGray)
	text.Draw(dst, gc2.Jiang12, ft14, cx, cy+32, colorGray)
	cx = m.X + 16
	cy += 64
	text.Draw(dst, "课体:"+b6.KeTi, ft14, cx, cy, colorWhite)

	cx, cy = m.X+194, m.Y+16
	dz := ThisGame.qmGame.Lunar.GetDayInGanZhiExact()
	jz := ThisGame.qmGame.YueJiang
	text.Draw(dst, fmt.Sprintf("月将%s", jz), ft14, cx, cy+16, colorWhite)
	text.Draw(dst, fmt.Sprintf("%s日", dz), ft14, cx, cy+32, colorWhite)
}

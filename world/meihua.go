package world

import (
	"fmt"
	"math"
	"strconv"

	"github.com/deminzhang/qimen-go/asset"
	"github.com/deminzhang/qimen-go/graphic"
	"github.com/deminzhang/qimen-go/gui"
	"github.com/deminzhang/qimen-go/xuan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	meiHuaGuaSize  = 25
	meiHuaUIWidth  = 240
	meiHuaUIHeight = 140
)

type MeiHua struct {
	X, Y int
	UI   *gui.BaseUI
	xuan.MeHua
	StartType string //起卦方式

	Mover       *Sprite
	GuaSprite   []*Sprite
	InputGuaNum *gui.InputBox
}

func NewMeiHua(x, y int) *MeiHua {
	m := &MeiHua{X: x, Y: y,
		UI: &gui.BaseUI{X: x, Y: y, Visible: true, W: meiHuaUIWidth, H: meiHuaUIHeight, BDColor: colorGray},
	}
	cbTimeStart := gui.NewCheckBox(94, 3, "时起")
	iptNumber := gui.NewInputBox(140, 3, 48, 20)
	cbTimeStart.SetChecked(true)
	iptNumber.Selectable = false
	iptNumber.DefaultText = "上下变"
	m.StartType = "时起"
	cbTimeStart.SetOnCheckChanged(func(c *gui.CheckBox) {
		if c.Checked() {
			m.StartType = "时起"
			m.TimeReset()
			iptNumber.Selectable = false
			iptNumber.SetText(fmt.Sprintf("%d", int(m.GuaUpIdx)*100+int(m.GuaDownIdx)*10+int(m.ChangeYaoIdx)))
		} else {
			m.StartType = ""
			iptNumber.Selectable = true
		}
	})
	doSet := func(i *gui.InputBox) {
		n, _ := strconv.Atoi(iptNumber.Text())
		m.Reset(uint(n/100), uint(n/10%10), uint(n%10))
		iptNumber.SetFocused(false)
	}
	iptNumber.SetOnLostFocus(doSet)
	iptNumber.SetOnPressEnter(doSet)
	m.InputGuaNum = iptNumber

	m.UI.AddChildren(cbTimeStart, iptNumber)
	gui.ActiveUI(m.UI)
	return m
}

func (m *MeiHua) Reset(upIdx, downIdx, change uint) {
	upIdx = (upIdx-1+8)%8 + 1
	downIdx = (downIdx-1+8)%8 + 1
	change = (change-1+6)%6 + 1
	if uint8(upIdx) == m.GuaUpIdx && uint8(downIdx) == m.GuaDownIdx && uint8(change) == m.ChangeYaoIdx {
		return
	}
	m.MeHua.Reset(upIdx, downIdx, change)

	if !m.InputGuaNum.Focused() {
		m.InputGuaNum.SetText(fmt.Sprintf("%d", int(m.GuaUpIdx)*100+int(m.GuaDownIdx)*10+int(m.ChangeYaoIdx)))
	}

	for _, sprite := range m.GuaSprite {
		ThisGame.StrokeManager.RemoveSprite(sprite)
	}
	m.GuaSprite = nil
}
func (m *MeiHua) TimeReset() {
	if ThisGame.qmGame.Lunar == nil {
		return
	}
	if m.StartType != "时起" {
		return
	}
	yz := xuan.ZhiIdx[ThisGame.qmGame.Lunar.GetYearZhiExact()]
	mz := ThisGame.qmGame.Lunar.GetMonth()
	dz := ThisGame.qmGame.Lunar.GetDay()
	hz := xuan.ZhiIdx[ThisGame.qmGame.Lunar.GetTimeZhi()]
	up := yz + mz + dz
	down := yz + mz + dz + hz
	m.Reset(uint(up), uint(down), uint(down))
}
func (m *MeiHua) Update() {
	m.TimeReset()
	dis := int(math.Round(float64(meiHuaGuaSize) * 1.25))
	if m.Mover == nil {
		m.Mover = NewSprite(graphic.NewRectImage(10), colorGray)
		m.Mover.onMove = func(sx, sy, dx, dy int) {
			m.X += dx
			m.Y += dy
			m.UI.X, m.UI.Y = m.X, m.Y
			if m.GuaSprite != nil {
				m.GuaSprite[0].MoveTo(m.X+32, m.Y+64)
				m.GuaSprite[1].MoveTo(m.X+32, m.Y+64+dis)
				m.GuaSprite[2].MoveTo(m.X+32+64, m.Y+64)
				m.GuaSprite[3].MoveTo(m.X+32+64, m.Y+64+dis)
				m.GuaSprite[4].MoveTo(m.X+32+128, m.Y+64)
				m.GuaSprite[5].MoveTo(m.X+32+128, m.Y+64+dis)
			}
		}
		ThisGame.StrokeManager.AddSprite(m.Mover)
		m.Mover.MoveTo(m.X, m.Y)
	}
	if m.GuaSprite == nil {
		cx, cy := m.X+32, m.Y+64
		m.GuaSprite = make([]*Sprite, 6)
		m.GuaSprite[0] = NewSprite(graphic.NewBaGuaImage(m.GuaUp, meiHuaGuaSize), color5Xing[xuan.DiagramsWuxing[m.GuaUp]])
		m.GuaSprite[1] = NewSprite(graphic.NewBaGuaImage(m.GuaDown, meiHuaGuaSize), color5Xing[xuan.DiagramsWuxing[m.GuaDown]])
		m.GuaSprite[0].MoveTo(cx, cy)
		m.GuaSprite[1].MoveTo(cx, cy+dis)
		m.GuaSprite[2] = NewSprite(graphic.NewBaGuaImage(m.GuaUpProcess, meiHuaGuaSize), color5Xing[xuan.DiagramsWuxing[m.GuaUpProcess]])
		m.GuaSprite[3] = NewSprite(graphic.NewBaGuaImage(m.GuaDownProcess, meiHuaGuaSize), color5Xing[xuan.DiagramsWuxing[m.GuaDownProcess]])
		m.GuaSprite[2].MoveTo(cx+64, cy)
		m.GuaSprite[3].MoveTo(cx+64, cy+dis)
		m.GuaSprite[4] = NewSprite(graphic.NewBaGuaImage(m.GuaUpChange, meiHuaGuaSize), color5Xing[xuan.DiagramsWuxing[m.GuaUpChange]])
		m.GuaSprite[5] = NewSprite(graphic.NewBaGuaImage(m.GuaDownChange, meiHuaGuaSize), color5Xing[xuan.DiagramsWuxing[m.GuaDownChange]])
		m.GuaSprite[4].MoveTo(cx+128, cy)
		m.GuaSprite[5].MoveTo(cx+128, cy+dis)
	}
}

func (m *MeiHua) Draw(dst *ebiten.Image) {
	m.Mover.Draw(dst)
	ft14, _ := asset.GetDefaultFontFace(14)
	text.Draw(dst, "梅花易数", ft14, m.X+16, m.Y+16, colorWhite)
	cx, cy := m.X+26, m.Y+32
	l := ThisGame.qmGame.Lunar
	//yz := l.GetYearZhiExact()
	mz := l.GetMonthZhiExact()
	dz := l.GetDayZhiExact()
	tz := l.GetTimeZhi()
	//text.Draw(dst, fmt.Sprintf("%s年", yz), ft14, cx+184, cy-16, colorWhite)
	text.Draw(dst, fmt.Sprintf("%s月", mz), ft14, cx+184, cy, ColorGanZhi(mz))
	text.Draw(dst, fmt.Sprintf("%s日", dz), ft14, cx+184, cy+16, ColorGanZhi(dz))
	text.Draw(dst, fmt.Sprintf("%s时", tz), ft14, cx+184, cy+32, ColorGanZhi(tz))
	text.Draw(dst, "本卦", ft14, cx, cy, colorWhite)
	text.Draw(dst, "互卦", ft14, cx+64, cy, colorWhite)
	text.Draw(dst, "变卦", ft14, cx+128, cy, colorWhite)
	cy += 16
	text.Draw(dst, m.GuaOrigin, ft14, cx, cy, colorWhite)
	text.Draw(dst, m.GuaProcess, ft14, cx+64, cy, colorWhite)
	text.Draw(dst, m.GuaChange, ft14, cx+128, cy, colorWhite)
	cx += 24
	cy += 32
	dis := int(math.Round(float64(meiHuaGuaSize) * 1.25))
	text.Draw(dst, fmt.Sprintf("%s%s", m.GuaUp, xuan.DiagramsWuxing[m.GuaUp]), ft14, cx+8, cy, color5Xing[xuan.DiagramsWuxing[m.GuaUp]])
	text.Draw(dst, fmt.Sprintf("%s%s", m.GuaDown, xuan.DiagramsWuxing[m.GuaDown]), ft14, cx+8, cy+dis, color5Xing[xuan.DiagramsWuxing[m.GuaDown]])
	text.Draw(dst, fmt.Sprintf("%s%s", m.GuaUpProcess, xuan.DiagramsWuxing[m.GuaUpProcess]), ft14, cx+8+64, cy, color5Xing[xuan.DiagramsWuxing[m.GuaUpProcess]])
	text.Draw(dst, fmt.Sprintf("%s%s", m.GuaDownProcess, xuan.DiagramsWuxing[m.GuaDownProcess]), ft14, cx+8+64, cy+dis, color5Xing[xuan.DiagramsWuxing[m.GuaDownProcess]])
	if m.ChangeYaoIdx > 3 {
		text.Draw(dst, "用", ft14, cx-40, cy, colorWhite)
		text.Draw(dst, "体", ft14, cx-40, cy+dis, colorWhite)
		text.Draw(dst, fmt.Sprintf("%s%s", m.GuaUpChange, xuan.DiagramsWuxing[m.GuaUpChange]), ft14, cx+8+128, cy, color5Xing[xuan.DiagramsWuxing[m.GuaUpChange]])
	} else {
		text.Draw(dst, "体", ft14, cx-40, cy, colorWhite)
		text.Draw(dst, "用", ft14, cx-40, cy+dis, colorWhite)
		text.Draw(dst, fmt.Sprintf("%s%s", m.GuaDownChange, xuan.DiagramsWuxing[m.GuaDownChange]), ft14, cx+8+128, cy+dis, color5Xing[xuan.DiagramsWuxing[m.GuaDownChange]])
	}

	for _, sprite := range m.GuaSprite {
		sprite.Draw(dst)
	}

}

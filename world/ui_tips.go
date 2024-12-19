package world

import (
	"github.com/deminzhang/qimen-go/gui"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"strings"
)

const (
	tipsUIWidth  = 100
	tipsUIHeight = 32
)

type UITips struct {
	gui.BaseUI
}

var tipsUI *UITips

func UIShowTips(x, y int, textLines []string) {
	u := NewUITips(x, y, textLines)
	if u == nil {
		return
	}
	tipsUI = u
	gui.ActiveUI(u)
}

func UIHideTips() {
	gui.CloseUI(tipsUI)
	tipsUI = nil
}

func NewUITips(x, y int, textLines []string) *UITips {
	if tipsUI != nil {
		UIHideTips()
	}
	u := &UITips{BaseUI: gui.BaseUI{Visible: true,
		X: x, Y: y,
		W: tipsUIWidth, H: tipsUIHeight,
	}}
	
	txt := gui.NewTextBox(0, 0, tipsUIWidth-16, tipsUIHeight-16)
	txt.BGColor = nil
	txt.UIImage = nil
	txt.TextColor = colorWhite
	txt.Text = strings.Join(textLines, "\n")
	w, h := txt.ContentSize()
	if w > u.W-8 {
		u.W = w + 8
		txt.W = w
	}
	if h > u.H-8 {
		u.H = h + 8
		txt.H = h
	}
	panelBG := gui.NewPanel(0, 0, txt.W, txt.H, color.RGBA{A: 0xcc})
	panelBG.BDColor = colorYellow
	panelBG.AddChildren(txt)
	u.AddChildren(panelBG)
	u.relocate()
	return u
}
func (u *UITips) relocate() {
	x, y := ebiten.CursorPosition()
	ww, wh := ebiten.WindowSize()
	if x > ww-u.W {
		if x-u.W > 0 {
			x = x - u.W //左翻
		} else {
			x = ww - u.W
		}
	}
	if x < 0 {
		x = 0
	}
	if y > wh-u.H {
		if y-u.H > 0 {
			y = y - u.H //上翻
		} else {
			y = wh - u.H
		}
	}
	if y < 0 {
		y = 0
	}
	u.X = x
	u.Y = y
}

func (u *UITips) Update() {
	u.BaseUI.Update()
}

func (u *UITips) Draw(screen *ebiten.Image) {
	u.BaseUI.Draw(screen)
}

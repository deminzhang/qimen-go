package world

import (
	"github.com/deminzhang/qimen-go/gui"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"strings"
)

const (
	tipsUIWidth  = 200
	tipsUIHeight = 110
)

type UITips struct {
	gui.BaseUI
	textMain  *gui.TextBox
	textLines []string
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
		tipsUI.textLines = textLines
		return nil
	}
	u := &UITips{BaseUI: gui.BaseUI{Visible: true,
		X: x, Y: y,
		W: tipsUIWidth, H: tipsUIHeight,
	}}
	u.SetText(textLines)
	panelBG := gui.NewPanel(0, 0, msgBoxUIWidth, msgBoxUIHeight, color.RGBA{A: 0xcc})
	panelBG.BDColor = colorLeader
	u.AddChildren(panelBG)
	textMain := gui.NewTextBox(0, 0, msgBoxUIWidth-16, msgBoxUIHeight-16)
	textMain.BGColor = nil
	textMain.UIImage = nil
	textMain.TextColor = colorWhite
	panelBG.AddChildren(textMain)
	textMain.Text = strings.Join(textLines, "\n")
	u.textMain = textMain
	return u
}
func (u *UITips) SetText(texts []string) {
	u.textLines = texts
	// TODO resize by texts
}
func (u *UITips) Update() {
	u.BaseUI.Update()
	x, y := ebiten.CursorPosition()
	if x > ScreenWidth-tipsUIWidth {
		x = ScreenWidth - tipsUIWidth
	}
	if x < 0 {
		x = 0
	}
	if y > ScreenHeight-tipsUIHeight {
		y = ScreenHeight - tipsUIHeight
	}
	if y < 0 {
		y = 0
	}
	u.X = x
	u.Y = y
}

func (u *UITips) Draw(screen *ebiten.Image) {
	u.BaseUI.Draw(screen)
}

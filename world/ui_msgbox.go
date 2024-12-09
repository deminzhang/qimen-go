package world

import (
	"github.com/deminzhang/qimen-go/gui"
)

const (
	msgBoxUIWidth  = 216
	msgBoxUIHeight = 110
)

type UIMsgBox struct {
	gui.BaseUI
	//panelBG    *gui.Panel
	//textMain   *gui.TextBox
	//btnConfirm *gui.Button
	//btnCancel  *gui.Button
}

func UIShowMsgBox(text, btnText1, btnText2 string, btnClick1, btnClick2 func(b *gui.Button)) {
	mb := NewUIMsgBox(text, btnText1, btnText2, btnClick1, btnClick2)
	gui.ActiveUI(mb)
}

func NewUIMsgBox(text, btnText1, btnText2 string, btnClick1, btnClick2 func(b *gui.Button)) *UIMsgBox {
	u := &UIMsgBox{BaseUI: gui.BaseUI{Visible: true,
		X: ScreenWidth/2 - msgBoxUIWidth/2, Y: ScreenHeight/2 - msgBoxUIHeight,
		W: msgBoxUIWidth, H: msgBoxUIHeight,
	}}
	panelBG := gui.NewPanel(0, 0, msgBoxUIWidth, msgBoxUIHeight, colorGray)
	textMain := gui.NewTextBox(8, 8, 200, 60)
	btnConfirm := gui.NewButton(40, 70, 48, 16, "confirm")
	btnCancel := gui.NewButton(130, 70, 48, 16, "cancel")
	u.AddChildren(panelBG)
	panelBG.AddChildren(textMain, btnConfirm, btnCancel)

	textMain.Text = text
	btnConfirm.Text = btnText1
	btnCancel.Text = btnText2
	btnConfirm.SetOnClick(func(b *gui.Button) {
		if btnClick1 != nil {
			btnClick1(b)
		}
		gui.CloseUI(u)
	})
	btnCancel.SetOnClick(btnClick2)
	btnCancel.SetOnClick(func(b *gui.Button) {
		if btnClick2 != nil {
			btnClick2(b)
		}
		gui.CloseUI(u)
	})
	return u
}

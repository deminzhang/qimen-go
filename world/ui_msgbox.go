package world

import (
	"github.com/deminzhang/qimen-go/gui"
	"image"
)

type UIMsgBox struct {
	gui.BaseUI
	panelBG    *gui.Panel
	textMain   *gui.TextBox
	btnConfirm *gui.Button
	btnCancel  *gui.Button
}

func UIShowMsgBox(text, btnText1, btnText2 string, btnClick1, btnClick2 func(b *gui.Button)) {
	mb := NewUIMsgBox(text, btnText1, btnText2, btnClick1, btnClick2)
	gui.ActiveUI(mb)
}

func NewUIMsgBox(text, btnText1, btnText2 string, btnClick1, btnClick2 func(b *gui.Button)) *UIMsgBox {
	u := &UIMsgBox{BaseUI: gui.BaseUI{Visible: true}}
	halfW := screenWidth / 2
	u.panelBG = gui.NewPanel(image.Rect(halfW-108, 230, halfW+108, 340), &colorGray)
	u.textMain = gui.NewTextBox(image.Rect(halfW-96, 240, halfW+96, 300))
	u.btnConfirm = gui.NewButton(image.Rect(halfW-64, 320, halfW-16, 336), "confirm")
	u.btnCancel = gui.NewButton(image.Rect(halfW+16, 320, halfW+64, 336), "cancel")
	u.AddChildren(u.panelBG, u.textMain, u.btnConfirm, u.btnCancel)

	u.textMain.Text = text
	u.btnConfirm.Text = btnText1
	u.btnCancel.Text = btnText2
	u.btnConfirm.SetOnClick(func(b *gui.Button) {
		if btnClick1 != nil {
			btnClick1(b)
		}
		gui.CloseUI(u)
	})
	u.btnCancel.SetOnClick(btnClick2)
	u.btnCancel.SetOnClick(func(b *gui.Button) {
		if btnClick2 != nil {
			btnClick2(b)
		}
		gui.CloseUI(u)
	})
	return u
}

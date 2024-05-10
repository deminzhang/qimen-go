package world

import (
	"image"
	"qimen/ui"
)

type UIMsgBox struct {
	ui.BaseUI
	textMain   *ui.TextBox
	btnConfirm *ui.Button
	btnCancel  *ui.Button
}

func UIShowMsgBox(text, btnText1, btnText2 string, btnClick1, btnClick2 func(b *ui.Button)) {
	mb := NewUIMsgBox(text, btnText1, btnText2, btnClick1, btnClick2)
	ui.ActiveUI(mb)
}

func NewUIMsgBox(text, btnText1, btnText2 string, btnClick1, btnClick2 func(b *ui.Button)) *UIMsgBox {
	u := &UIMsgBox{BaseUI: ui.BaseUI{Visible: true}}
	u.textMain = ui.NewTextBox(image.Rect(ScreenWidth/2-96, 240, ScreenWidth/2+96, 300))
	u.btnConfirm = ui.NewButton(image.Rect(ScreenWidth/2-64, 320, ScreenWidth/2-16, 336), "confirm")
	u.btnCancel = ui.NewButton(image.Rect(ScreenWidth/2+16, 320, ScreenWidth/2+64, 336), "cancel")
	u.AddChild(u.textMain)
	u.AddChild(u.btnConfirm)
	u.AddChild(u.btnCancel)

	u.textMain.Text = text
	u.btnConfirm.Text = btnText1
	u.btnCancel.Text = btnText2
	u.btnConfirm.SetOnClick(func(b *ui.Button) {
		if btnClick1 != nil {
			btnClick1(b)
		}
		ui.CloseUI(u)
	})
	u.btnCancel.SetOnClick(btnClick2)
	u.btnCancel.SetOnClick(func(b *ui.Button) {
		if btnClick2 != nil {
			btnClick2(b)
		}
		ui.CloseUI(u)
	})
	return u
}

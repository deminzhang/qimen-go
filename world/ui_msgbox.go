package world

import (
	"image"
	"qimen/ebiten_ui"
)

type UIMsgBox struct {
	ebiten_ui.BaseUI
	textMain   *ebiten_ui.TextBox
	btnConfirm *ebiten_ui.Button
	btnCancel  *ebiten_ui.Button
}

func UIShowMsgBox(text, btnText1, btnText2 string, btnClick1, btnClick2 func(b *ebiten_ui.Button)) {
	mb := NewUIMsgBox(text, btnText1, btnText2, btnClick1, btnClick2)
	ebiten_ui.ActiveUI(mb)
}

func NewUIMsgBox(text, btnText1, btnText2 string, btnClick1, btnClick2 func(b *ebiten_ui.Button)) *UIMsgBox {
	u := &UIMsgBox{}
	u.textMain = ebiten_ui.NewTextBox(image.Rect(ScreenWidth/2-96, 240, ScreenWidth/2+96, 300))
	u.btnConfirm = ebiten_ui.NewButton(image.Rect(ScreenWidth/2-64, 320, ScreenWidth/2-16, 336), "confirm")
	u.btnCancel = ebiten_ui.NewButton(image.Rect(ScreenWidth/2+16, 320, ScreenWidth/2+64, 336), "cancel")
	u.AddChild(u.textMain)
	u.AddChild(u.btnConfirm)
	u.AddChild(u.btnCancel)

	u.textMain.Text = text
	u.btnConfirm.Text = btnText1
	u.btnCancel.Text = btnText2
	u.btnConfirm.SetOnClick(func(b *ebiten_ui.Button) {
		if btnClick1 != nil {
			btnClick1(b)
		}
		ebiten_ui.CloseUI(u)
	})
	u.btnCancel.SetOnClick(btnClick2)
	u.btnCancel.SetOnClick(func(b *ebiten_ui.Button) {
		if btnClick2 != nil {
			btnClick2(b)
		}
		ebiten_ui.CloseUI(u)
	})
	return u
}

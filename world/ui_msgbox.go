package world

import (
	"github.com/deminzhang/qimen-go/ui"
	"image"
)

type UIMsgBox struct {
	ui.BaseUI
	panelBG    *ui.Panel
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
	u.panelBG = ui.NewPanel(image.Rect(screenWidth/2-108, 230, screenWidth/2+108, 340), &colorGray)
	u.textMain = ui.NewTextBox(image.Rect(screenWidth/2-96, 240, screenWidth/2+96, 300))
	u.btnConfirm = ui.NewButton(image.Rect(screenWidth/2-64, 320, screenWidth/2-16, 336), "confirm")
	u.btnCancel = ui.NewButton(image.Rect(screenWidth/2+16, 320, screenWidth/2+64, 336), "cancel")
	u.AddChild(u.panelBG)
	//u.AddChild(u.textBg)
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

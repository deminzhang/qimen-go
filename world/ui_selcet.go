package world

import (
	"fmt"
	"github.com/6tail/lunar-go/calendar"
	"image"
	"qimen/ui"
	"strconv"
)

type UISelect struct {
	ui.BaseUI
	panelBG                                                   *ui.Panel
	inputSYear, inputSMonth, inputSDay, inputSHour, inputSMin *ui.InputBox
	opMale, opFemale                                          *ui.OptionBox
	btnX, btnOK                                               *ui.Button
}

func UIShowSelect() *UISelect {
	uiSelect := NewUISelect()
	ui.ActiveUI(uiSelect)
	return uiSelect
}

func NewUISelect() *UISelect {
	p := &UISelect{BaseUI: ui.BaseUI{Visible: true}}
	px0, py0 := screenWidth/2-176, 210
	h := 32
	p.panelBG = ui.NewPanel(image.Rect(screenWidth/2-176, 210, screenWidth/2+176, 310), &colorGray)
	p.btnX = ui.NewButton(image.Rect(px0+72*4+32, py0, px0+72*4+64, py0+h), "X")
	py0 += 32
	px0 += 8
	p.inputSYear = ui.NewInputBox(image.Rect(px0, py0, px0+64, py0+h))
	p.inputSMonth = ui.NewInputBox(image.Rect(px0+70, py0, px0+70+48, py0+h))
	p.inputSDay = ui.NewInputBox(image.Rect(px0+70+52, py0, px0+70+52+48, py0+h))
	p.inputSHour = ui.NewInputBox(image.Rect(px0+70+52*2, py0, px0+70+52*2+48, py0+h))
	p.inputSMin = ui.NewInputBox(image.Rect(px0+70+52*3, py0, px0+70+52*3+48, py0+h))
	py0 += 32
	p.opMale = ui.NewOptionBox(px0+70+52*2, py0+8, "男")
	p.opFemale = ui.NewOptionBox(px0+70+52*3, py0+8, "女")
	p.btnOK = ui.NewButton(image.Rect(px0+70*4, py0, px0+70*4+64, py0+h), "确定")

	p.AddChild(p.panelBG)
	p.inputSYear.MaxChars = 5
	p.inputSMonth.MaxChars = 2
	p.inputSDay.MaxChars = 2
	p.inputSHour.MaxChars = 2
	p.inputSMin.MaxChars = 2
	p.inputSYear.DefaultText = "公元年"
	p.inputSMonth.DefaultText = "月"
	p.inputSDay.DefaultText = "日"
	p.inputSHour.DefaultText = "时"
	p.inputSMin.DefaultText = "分"
	oldBirthTime := ThisGame.baZi.Player.Lunar
	if oldBirthTime != nil {
		solar := oldBirthTime.GetSolar()
		p.inputSYear.SetText(solar.GetYear())
		p.inputSMonth.SetText(solar.GetMonth())
		p.inputSDay.SetText(solar.GetDay())
		p.inputSHour.SetText(solar.GetHour())
		p.inputSMin.SetText(solar.GetMinute())
	}
	p.panelBG.AddChildren(p.inputSYear, p.inputSMonth, p.inputSDay, p.inputSHour, p.inputSMin,
		p.btnX, p.btnOK, p.opMale, p.opFemale)
	ui.MakeOptionBoxGroup(p.opMale, p.opFemale)
	p.opMale.Select()

	p.btnX.SetOnClick(func(b *ui.Button) {
		ui.CloseUI(p)
	})
	p.btnOK.SetOnClick(func(b *ui.Button) {
		defer func() {
			s := recover()
			if s != nil {
				UIShowMsgBox(fmt.Sprintf("时间不对%s", s), "确定", "取消", nil, nil)
			}
		}()
		year, _ := strconv.Atoi(p.inputSYear.Text())
		month, _ := strconv.Atoi(p.inputSMonth.Text())
		day, _ := strconv.Atoi(p.inputSDay.Text())
		hour, _ := strconv.Atoi(p.inputSHour.Text())
		minute, _ := strconv.Atoi(p.inputSMin.Text())
		solar := calendar.NewSolar(year, month, day, hour, minute, 0)
		gender := GenderFemale
		if p.opMale.Selected() {
			gender = GenderMale
		}
		ThisGame.baZi.Player.Reset(calendar.NewLunarFromSolar(solar), gender)
		ui.CloseUI(p)
	})
	return p
}

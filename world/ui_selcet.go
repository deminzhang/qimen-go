package world

import (
	"fmt"
	"github.com/6tail/lunar-go/calendar"
	"github.com/deminzhang/qimen-go/gui"
	"strconv"
	"time"
)

const (
	selectUIWidth  = 352
	selectUIHeight = 100
)

type UISelect struct {
	gui.BaseUI
	//panelBG *gui.Panel
	//inputSYear, inputSMonth, inputSDay, inputSHour, inputSMin *gui.InputBox
	//opMale, opFemale                                          *gui.OptionBox
	//btnX, btnOK                                               *gui.Button
	onOK func(*calendar.Solar, int)
}

func UIShowSelectBirth(date *calendar.Solar, gender int, onOK func(*calendar.Solar, int)) *UISelect {
	uiSelect := NewUISelect(date, gender, onOK)
	gui.ActiveUI(uiSelect)
	return uiSelect
}

func NewUISelect(solar *calendar.Solar, gender int, onOK func(*calendar.Solar, int)) *UISelect {
	p := &UISelect{BaseUI: gui.BaseUI{Visible: true,
		X: ScreenWidth/2 - selectUIWidth/2, Y: ScreenHeight/2 - selectUIHeight/2,
		W: selectUIWidth, H: selectUIHeight,
	}}
	px0, py0 := 0, 0
	h := 32
	panelBG := gui.NewPanel(0, 0, selectUIWidth, selectUIHeight, colorGray)
	btnX := gui.NewButton(px0+72*4+32, py0, 32, h, "X")
	py0 += 32
	px0 += 8
	inputSYear := gui.NewInputBox(px0, py0, 64, h)
	inputSMonth := gui.NewInputBox(px0+70, py0, 48, h)
	inputSDay := gui.NewInputBox(px0+70+52, py0, 48, h)
	inputSHour := gui.NewInputBox(px0+70+52*2, py0, 48, h)
	inputSMin := gui.NewInputBox(px0+70+52*3, py0, 48, h)
	btnNow := gui.NewButton(px0+70*4, py0, 64, h, "此时")
	py0 += 32
	opMale := gui.NewOptionBox(px0+70+52*2, py0+8, "男")
	opFemale := gui.NewOptionBox(px0+70+52*3, py0+8, "女")
	btnOK := gui.NewButton(px0+70*4, py0, 64, h, "确定")

	p.AddChildren(panelBG)
	inputSYear.MaxChars = 5
	inputSMonth.MaxChars = 2
	inputSDay.MaxChars = 2
	inputSHour.MaxChars = 2
	inputSMin.MaxChars = 2
	inputSYear.DefaultText = "公元年"
	inputSMonth.DefaultText = "月"
	inputSDay.DefaultText = "日"
	inputSHour.DefaultText = "时"
	inputSMin.DefaultText = "分"
	if solar != nil {
		inputSYear.SetText(solar.GetYear())
		inputSMonth.SetText(solar.GetMonth())
		inputSDay.SetText(solar.GetDay())
		inputSHour.SetText(solar.GetHour())
		inputSMin.SetText(solar.GetMinute())
	}
	panelBG.AddChildren(inputSYear, inputSMonth, inputSDay, inputSHour, inputSMin,
		btnX, btnOK, btnNow, opMale, opFemale)
	gui.MakeOptionBoxGroup(opMale, opFemale)
	if gender == GenderFemale {
		opFemale.Select()
	} else {
		opMale.Select()
	}

	btnX.SetOnClick(func(b *gui.Button) {
		gui.CloseUI(p)
	})
	btnOK.SetOnClick(func(b *gui.Button) {
		defer func() {
			s := recover()
			if s != nil {
				UIShowMsgBox(fmt.Sprintf("时间不对%s", s), "确定", "取消", nil, nil)
			}
		}()
		year, _ := strconv.Atoi(inputSYear.Text())
		month, _ := strconv.Atoi(inputSMonth.Text())
		day, _ := strconv.Atoi(inputSDay.Text())
		hour, _ := strconv.Atoi(inputSHour.Text())
		minute, _ := strconv.Atoi(inputSMin.Text())
		s := calendar.NewSolar(year, month, day, hour, minute, 0)
		g := GenderFemale
		if opMale.Selected() {
			g = GenderMale
		}
		onOK(s, g)
		gui.CloseUI(p)
	})
	btnNow.SetOnClick(func(b *gui.Button) {
		s := calendar.NewSolarFromDate(time.Now())
		inputSYear.SetText(s.GetYear())
		inputSMonth.SetText(s.GetMonth())
		inputSDay.SetText(s.GetDay())
		inputSHour.SetText(s.GetHour())
		inputSMin.SetText(s.GetMinute())
	})
	return p
}

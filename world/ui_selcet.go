package world

import (
	"fmt"
	"github.com/6tail/lunar-go/calendar"
	"github.com/deminzhang/qimen-go/gui"
	"strconv"
)

type UISelect struct {
	gui.BaseUI
	panelBG                                                   *gui.Panel
	inputSYear, inputSMonth, inputSDay, inputSHour, inputSMin *gui.InputBox
	opMale, opFemale                                          *gui.OptionBox
	btnX, btnOK                                               *gui.Button
	onOK                                                      func(*calendar.Solar, int)
}

func UIShowSelectBirth(date *calendar.Solar, gender int, onOK func(*calendar.Solar, int)) *UISelect {
	uiSelect := NewUISelect(date, gender, onOK)
	gui.ActiveUI(uiSelect)
	return uiSelect
}

const (
	UISelectWidth  = 352
	UISelectHeight = 100
)

func NewUISelect(solar *calendar.Solar, gender int, onOK func(*calendar.Solar, int)) *UISelect {
	p := &UISelect{BaseUI: gui.BaseUI{Visible: true,
		X: screenWidth/2 - UISelectWidth/2, Y: screenHeight/2 - UISelectHeight/2,
		W: UISelectWidth, H: UISelectHeight,
	}}
	px0, py0 := 0, 0
	h := 32
	p.panelBG = gui.NewPanel(0, 0, UISelectWidth, UISelectHeight, &colorGray)
	p.btnX = gui.NewButton(px0+72*4+32, py0, 32, h, "X")
	py0 += 32
	px0 += 8
	p.inputSYear = gui.NewInputBox(px0, py0, 64, h)
	p.inputSMonth = gui.NewInputBox(px0+70, py0, 48, h)
	p.inputSDay = gui.NewInputBox(px0+70+52, py0, 48, h)
	p.inputSHour = gui.NewInputBox(px0+70+52*2, py0, 48, h)
	p.inputSMin = gui.NewInputBox(px0+70+52*3, py0, 48, h)
	py0 += 32
	p.opMale = gui.NewOptionBox(px0+70+52*2, py0+8, "男")
	p.opFemale = gui.NewOptionBox(px0+70+52*3, py0+8, "女")
	p.btnOK = gui.NewButton(px0+70*4, py0, 64, h, "确定")

	p.AddChildren(p.panelBG)
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
	if solar != nil {
		p.inputSYear.SetText(solar.GetYear())
		p.inputSMonth.SetText(solar.GetMonth())
		p.inputSDay.SetText(solar.GetDay())
		p.inputSHour.SetText(solar.GetHour())
		p.inputSMin.SetText(solar.GetMinute())
	}
	p.panelBG.AddChildren(p.inputSYear, p.inputSMonth, p.inputSDay, p.inputSHour, p.inputSMin,
		p.btnX, p.btnOK, p.opMale, p.opFemale)
	gui.MakeOptionBoxGroup(p.opMale, p.opFemale)
	if gender == GenderFemale {
		p.opFemale.Select()
	} else {
		p.opMale.Select()
	}

	p.btnX.SetOnClick(func(b *gui.Button) {
		gui.CloseUI(p)
	})
	p.btnOK.SetOnClick(func(b *gui.Button) {
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
		s := calendar.NewSolar(year, month, day, hour, minute, 0)
		g := GenderFemale
		if p.opMale.Selected() {
			g = GenderMale
		}
		onOK(s, g)
		gui.CloseUI(p)
	})
	return p
}

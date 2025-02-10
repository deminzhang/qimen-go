package world

import (
	"fmt"
	"github.com/6tail/lunar-go/calendar"
	"github.com/deminzhang/qimen-go/gui"
	"github.com/deminzhang/qimen-go/xuan"
	"strconv"
	"time"
)

type UIQiMen struct {
	gui.BaseUI
	//pan                                                       *xuan.QMGame
	//panelSDate                                                *gui.Panel
	//panelOpCb                                                 *gui.Panel
	inputSYear, inputSMonth, inputSDay, inputSHour, inputSMin *gui.InputBox

	//opTypeRoll, opTypeFly, opTypeAmaze *gui.OptionBox
	cbHostingType, cbFlyType *gui.CheckBox
	cbAuto                   *gui.CheckBox

	opStartSplit, opStartMaoShan, opStartZhiRun *gui.OptionBox
	opStartSelf                                 *gui.OptionBox
	inputSelfJu                                 *gui.InputBox

	opHideGan0, opHideGan1 *gui.OptionBox

	//btnCalc             *gui.Button
	//btnNow              *gui.Button
	//btnPreJu, btnNextJu *gui.Button

	//opHourPan, opDayPan, opMonthPan, opYearPan *gui.OptionBox
	//opDay2Pan                                  *gui.OptionBox

	year, month, day, hour, minute int
	solar                          *calendar.Solar
	qmParams                       xuan.QMParams
}

var uiQiMen *UIQiMen

func UIShowQiMen() *UIQiMen {
	if uiQiMen == nil {
		uiQiMen = NewUIQiMen()
		gui.ActiveUI(uiQiMen)
	}
	return uiQiMen
}
func UIHideQiMen() {
	if uiQiMen != nil {
		gui.CloseUI(uiQiMen)
		uiQiMen = nil
	}
}

func NewUIQiMen() *UIQiMen {
	p := &UIQiMen{
		BaseUI: gui.BaseUI{Visible: true, W: ScreenWidth, H: ScreenHeight},
		qmParams: xuan.QMParams{
			Type:        xuan.QMTypeRotating,
			HostingType: xuan.QMHostingType2,
			FlyType:     xuan.QMFlyTypeAllOrder,
			JuType:      xuan.QMJuTypeSplit,
			HideGanType: xuan.QMHideGanDutyDoorHour,
		},
	}
	px0, py0 := 0, 0
	h := 28
	panelSDate := gui.NewPanel(150, 0, 380, 32, nil)
	p.inputSYear = gui.NewInputBox(px0, py0, 64, h)
	p.inputSMonth = gui.NewInputBox(px0+72, py0, 64, h)
	p.inputSDay = gui.NewInputBox(px0+72*2, py0, 64, h)
	p.inputSHour = gui.NewInputBox(px0+72*3, py0, 64, h)
	p.inputSMin = gui.NewInputBox(px0+72*4, py0, 64, h)
	px0, py0 = 4, 4
	panelOpCb := gui.NewPanel(630, 0, 400, 130, nil)
	opTypeRoll := gui.NewOptionBox(px0, py0, xuan.QMType[0])
	opTypeFly := gui.NewOptionBox(px0+72*1, py0, xuan.QMType[1])
	opTypeAmaze := gui.NewOptionBox(px0+72*2, py0, xuan.QMType[2])
	btnCalc := gui.NewTextButton(px0+72*5, py0, "排局", colorWhite, colorGray)

	py0 += 18
	p.cbHostingType = gui.NewCheckBox(px0, py0, xuan.QMHostingType[xuan.QMHostingType28])
	p.cbFlyType = gui.NewCheckBox(px0+72*1, py0, xuan.QMFlyType[xuan.QMFlyTypeAllOrder])
	btnPreJu := gui.NewTextButton(px0+72*5, py0, "上一局", colorWhite, colorGray)

	py0 += 18
	btnNextJu := gui.NewTextButton(px0+72*5, py0, "下一局", colorWhite, colorGray)

	p.opStartSplit = gui.NewOptionBox(px0, py0, xuan.QMJuType[xuan.QMJuTypeSplit])
	p.opStartMaoShan = gui.NewOptionBox(px0+72, py0, xuan.QMJuType[xuan.QMJuTypeMaoShan])
	p.opStartZhiRun = gui.NewOptionBox(px0+72*2, py0, xuan.QMJuType[xuan.QMJuTypeZhiRun])
	p.opStartSelf = gui.NewOptionBox(px0+72*3, py0, xuan.QMJuType[xuan.QMJuTypeSelf])
	p.inputSelfJu = gui.NewInputBox(px0+72*3, py0, 32, h)

	py0 += 18
	p.opHideGan0 = gui.NewOptionBox(px0, py0, xuan.QMHideGanType[xuan.QMHideGanDutyDoorHour])
	p.opHideGan1 = gui.NewOptionBox(px0+72, py0, xuan.QMHideGanType[xuan.QMHideGanDoorHomeGan])
	btnNow := gui.NewTextButton(px0+72*5, py0, "此时", colorWhite, colorGray)
	py0 += 18
	opHourPan := gui.NewOptionBox(px0, py0, "时家")
	opDayPan := gui.NewOptionBox(px0+72, py0, "日家")
	opDay2Pan := gui.NewOptionBox(px0+72*1.5, py0, "_日家2")
	opMonthPan := gui.NewOptionBox(px0+72*2, py0, "月家")
	opYearPan := gui.NewOptionBox(px0+72*3, py0, "年家")
	p.cbAuto = gui.NewCheckBox(px0+72*4, py0, "自动")
	py0 += 18
	cbQiMenPan := gui.NewCheckBox(px0+72, py0, "奇门")
	cbBig6Pan := gui.NewCheckBox(px0+72*2, py0, "大六壬")
	cbMeiHuaPan := gui.NewCheckBox(px0+72*3, py0, "梅花")
	cbChar8Pan := gui.NewCheckBox(px0+72*4, py0, "四柱")
	cbStarPan := gui.NewCheckBox(px0+72*5, py0, "星盘")

	p.AddChildren(panelSDate, panelOpCb)
	p.inputSYear.MaxChars = 5
	p.inputSMonth.MaxChars = 2
	p.inputSDay.MaxChars = 2
	p.inputSHour.MaxChars = 2
	p.inputSMin.MaxChars = 2
	p.inputSYear.DefaultText = "年"
	p.inputSMonth.DefaultText = "月"
	p.inputSDay.DefaultText = "日"
	p.inputSHour.DefaultText = "时"
	p.inputSMin.DefaultText = "分"
	p.inputSelfJu.DefaultText = "手选局数"
	panelSDate.AddChildren(p.inputSYear, p.inputSMonth, p.inputSDay, p.inputSHour, p.inputSMin)
	panelOpCb.AddChildren(btnCalc, btnNow, opTypeRoll, opTypeFly, opTypeAmaze,
		btnPreJu, btnNextJu, p.cbHostingType, p.cbFlyType,
		opHourPan, opDayPan, opMonthPan, opYearPan, opDay2Pan,
		p.opStartSplit, p.opStartMaoShan, p.opStartZhiRun, p.opStartSelf, p.inputSelfJu,
		p.cbAuto, cbChar8Pan, cbStarPan, cbMeiHuaPan, cbBig6Pan, cbQiMenPan,
		p.opHideGan0, p.opHideGan1)

	gui.MakeOptionBoxGroup(opTypeRoll, opTypeFly, opTypeAmaze)
	opTypeRoll.Select()
	opTypeRoll.SetOnSelect(func(c *gui.OptionBox) {
		p.qmParams.Type = xuan.QMTypeRotating
		p.Apply(p.solar)
	})
	opTypeFly.SetOnSelect(func(c *gui.OptionBox) {
		p.qmParams.Type = xuan.QMTypeFly
		p.Apply(p.solar)
	})
	opTypeAmaze.SetOnSelect(func(c *gui.OptionBox) {
		p.qmParams.Type = xuan.QMTypeAmaze
		p.opStartSplit.Select()
	})

	gui.MakeOptionBoxGroup(p.opStartSplit, p.opStartMaoShan, p.opStartZhiRun, p.opStartSelf)
	p.opStartSplit.Select()
	p.opStartSplit.SetOnSelect(func(c *gui.OptionBox) {
		p.qmParams.JuType = xuan.QMJuTypeSplit
		p.Apply(p.solar)
	})
	p.opStartMaoShan.SetOnSelect(func(c *gui.OptionBox) {
		p.qmParams.JuType = xuan.QMJuTypeMaoShan
		p.Apply(p.solar)
	})
	p.opStartZhiRun.SetOnSelect(func(c *gui.OptionBox) {
		p.qmParams.JuType = xuan.QMJuTypeZhiRun
		p.Apply(p.solar)
	})
	p.opStartSelf.SetOnSelect(func(c *gui.OptionBox) {
		p.opStartSelf.Visible = false
		p.inputSelfJu.Visible = true
		p.inputSelfJu.SetFocused(true)
		p.qmParams.JuType = xuan.QMJuTypeSelf
	})

	gui.MakeOptionBoxGroup(p.opHideGan0, p.opHideGan1)
	p.opHideGan0.Select()
	p.opHideGan0.SetOnSelect(func(c *gui.OptionBox) {
		p.qmParams.HideGanType = xuan.QMHideGanDutyDoorHour
		p.Apply(p.solar)
	})
	p.opHideGan1.SetOnSelect(func(c *gui.OptionBox) {
		p.qmParams.HideGanType = xuan.QMHideGanDoorHomeGan
		p.Apply(p.solar)
	})

	gui.MakeOptionBoxGroup(opHourPan, opDayPan, opMonthPan, opYearPan, opDay2Pan)
	opHourPan.Select()
	opHourPan.SetOnSelect(func(c *gui.OptionBox) {
		p.qmParams.YMDH = xuan.QMGameHour
		p.Apply(p.solar)
	})
	opDayPan.SetOnSelect(func(c *gui.OptionBox) {
		p.qmParams.YMDH = xuan.QMGameDay
		p.Apply(p.solar)
	})
	opMonthPan.SetOnSelect(func(c *gui.OptionBox) {
		p.qmParams.YMDH = xuan.QMGameMonth
		p.Apply(p.solar)
	})
	opYearPan.SetOnSelect(func(c *gui.OptionBox) {
		p.qmParams.YMDH = xuan.QMGameYear
		p.Apply(p.solar)
	})
	opDay2Pan.SetOnSelect(func(c *gui.OptionBox) {
		p.qmParams.YMDH = xuan.QMGameDay2
		p.Apply(p.solar)
	})
	opDay2Pan.Disabled = true
	opDay2Pan.Visible = false
	cbQiMenPan.SetChecked(true)
	cbQiMenPan.SetOnCheckChanged(func(c *gui.CheckBox) {
		if ThisGame != nil {
			ThisGame.showQiMen = c.Checked()
		}
	})
	cbBig6Pan.SetChecked(true)
	cbBig6Pan.SetOnCheckChanged(func(c *gui.CheckBox) {
		if ThisGame != nil {
			ThisGame.showBig6 = c.Checked()
		}
	})
	cbMeiHuaPan.SetChecked(true)
	cbMeiHuaPan.SetOnCheckChanged(func(c *gui.CheckBox) {
		if ThisGame != nil {
			ThisGame.showMeiHua = c.Checked()
		}
	})
	cbChar8Pan.SetChecked(true)
	cbChar8Pan.SetOnCheckChanged(func(c *gui.CheckBox) {
		if ThisGame != nil {
			ThisGame.showChar8 = c.Checked()
		}
	})
	cbStarPan.SetChecked(true)
	cbStarPan.SetOnCheckChanged(func(c *gui.CheckBox) {
		if ThisGame != nil {
			ThisGame.showAstrolabe = c.Checked()
		}
	})

	p.cbHostingType.SetChecked(false)
	p.cbHostingType.Visible = opTypeRoll.Selected()
	p.cbHostingType.SetOnCheckChanged(func(c *gui.CheckBox) {
		if c.Checked() {
			p.qmParams.HostingType = xuan.QMHostingType28
		} else {
			p.qmParams.HostingType = xuan.QMHostingType2
		}
		p.Apply(p.solar)
	})
	p.cbFlyType.SetChecked(true)
	p.cbFlyType.Visible = opTypeFly.Selected()
	p.cbFlyType.SetOnCheckChanged(func(c *gui.CheckBox) {
		if c.Checked() {
			p.qmParams.FlyType = xuan.QMFlyTypeAllOrder
		} else {
			p.qmParams.FlyType = xuan.QMFlyTypeLunarReverse
		}
		p.Apply(p.solar)
	})
	p.cbAuto.SetOnCheckChanged(func(c *gui.CheckBox) {
		ThisGame.autoMinute = c.Checked()
	})

	btnCalc.SetOnClick(func() {
		defer func() {
			s := recover()
			if s != nil {
				UIShowMsgBox(fmt.Sprintf("%s", s), "确定", "取消", nil, nil)
			}
		}()
		year, _ := strconv.Atoi(p.inputSYear.Text())
		month, _ := strconv.Atoi(p.inputSMonth.Text())
		day, _ := strconv.Atoi(p.inputSDay.Text())
		hour, _ := strconv.Atoi(p.inputSHour.Text())
		minute, _ := strconv.Atoi(p.inputSMin.Text())

		if p.qmParams.JuType == xuan.QMJuTypeSelf {
			ju, err := strconv.Atoi(p.inputSelfJu.Text())
			if err != nil || !((ju >= 1 && ju <= 9) || (ju >= -9 && ju <= -1)) {
				panic("局数不对,限(阴[-9,-1],阳[1,9])")
			}
			p.qmParams.SelfJu = ju
		} else {
			p.qmParams.SelfJu = 0
		}
		solar := calendar.NewSolar(year, month, day, hour, minute, 0)
		p.Apply(solar)
	})
	btnNow.SetOnClick(func() {
		defer func() {
			s := recover()
			if s != nil {
				UIShowMsgBox(fmt.Sprintf("%s", s), "确定", "取消", nil, nil)
			}
		}()
		solar := calendar.NewSolarFromDate(time.Now())

		if p.qmParams.JuType == xuan.QMJuTypeSelf {
			ju, err := strconv.Atoi(p.inputSelfJu.Text())
			if err != nil || !((ju >= 1 && ju <= 9) || (ju >= -9 && ju <= -1)) {
				panic("局数不对,限(阴[-9,-1],阳[1,9])")
			}
			p.qmParams.SelfJu = ju
		} else {
			p.qmParams.SelfJu = 0
		}
		p.Apply(solar)
	})
	btnPreJu.SetOnClick(func() {
		var solar *calendar.Solar
		switch p.qmParams.YMDH {
		case xuan.QMGameYear:
			solar = p.solar.NextYear(-1)
		case xuan.QMGameMonth:
			solar = p.solar.NextMonth(-1)
		case xuan.QMGameDay:
			solar = p.solar.NextDay(-1)
		case xuan.QMGameHour:
			solar = p.solar.NextHour(-2)
		}
		p.Apply(solar)
	})
	btnNextJu.SetOnClick(func() {
		var solar *calendar.Solar
		switch p.qmParams.YMDH {
		case xuan.QMGameYear:
			solar = p.solar.NextYear(1)
		case xuan.QMGameMonth:
			solar = p.solar.NextMonth(1)
		case xuan.QMGameDay:
			solar = p.solar.NextDay(1)
		case xuan.QMGameHour:
			solar = p.solar.NextHour(2)
		}
		p.Apply(solar)
	})

	uiQiMen = p
	return p
}

func (p *UIQiMen) OnClose() {
	uiQiMen = nil
}

func (p *UIQiMen) Apply(solar *calendar.Solar) *xuan.QMGame {
	defer func() {
		s := recover()
		if s != nil {
			UIShowMsgBox(fmt.Sprintf("时间不对%s", s), "确定", "取消", nil, nil)
		}
	}()
	pan := xuan.NewQMGame(solar, p.qmParams)
	if ThisGame != nil {
		ThisGame.qmGame = pan
	}
	p.solar = solar
	//pan.DayArr
	p.inputSYear.SetText(pan.Solar.GetYear())
	p.inputSMonth.SetText(pan.Solar.GetMonth())
	p.inputSDay.SetText(pan.Solar.GetDay())
	p.inputSHour.SetText(pan.Solar.GetHour())
	p.inputSMin.SetText(pan.Solar.GetMinute())

	p.cbHostingType.Visible = p.qmParams.Type == xuan.QMTypeRotating
	p.cbFlyType.Visible = p.qmParams.Type == xuan.QMTypeFly

	shiPan := p.qmParams.YMDH == xuan.QMGameHour
	p.opStartSplit.Visible = shiPan // && p.qmParams.Type != qimen.QMTypeAmaze
	p.opStartMaoShan.Visible = shiPan && p.qmParams.Type != xuan.QMTypeAmaze
	p.opStartZhiRun.Visible = shiPan && p.qmParams.Type != xuan.QMTypeAmaze
	p.opStartSelf.Visible = !p.opStartSelf.Selected()
	p.inputSelfJu.Visible = p.opStartSelf.Selected()

	p.opHideGan0.Visible = p.qmParams.Type != xuan.QMTypeAmaze
	p.opHideGan1.Visible = p.qmParams.Type != xuan.QMTypeAmaze

	//fmt
	switch p.qmParams.YMDH {
	case xuan.QMGameHour:
		pan.ShowTimeGame()
	case xuan.QMGameDay:
		pan.ShowDayGame()
	case xuan.QMGameMonth:
		pan.ShowMonthGame()
	case xuan.QMGameYear:
		pan.ShowYearGame()
	}
	pan.CalBig6()
	return pan
}
func (p *UIQiMen) NextHour() *xuan.QMGame {
	solar := p.solar.NextHour(2)
	return p.Apply(solar)
}

func (p *UIQiMen) NextMinute() *xuan.QMGame {
	s := p.solar
	if s.GetMinute() == 59 {
		s = s.NextHour(1)
		s = calendar.NewSolar(s.GetYear(), s.GetMonth(), s.GetDay(), s.GetHour(), 0, 0)
	} else {
		s = calendar.NewSolar(s.GetYear(), s.GetMonth(), s.GetDay(), s.GetHour(), s.GetMinute()+1, 0)
	}
	return p.Apply(s)
}

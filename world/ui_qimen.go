package world

import (
	"fmt"
	"github.com/6tail/lunar-go/calendar"
	"github.com/deminzhang/qimen-go/qimen"
	"github.com/deminzhang/qimen-go/ui"
	"image"
	"strconv"
	"time"
)

var gongOffset = [][]int{{0, 0},
	{1, 2}, {2, 0}, {0, 1},
	{0, 0}, {1, 1}, {2, 2},
	{2, 1}, {0, 2}, {1, 0},
}

type UIQiMen struct {
	ui.BaseUI
	//pan                                                       *qimen.QMGame
	panelSDate                                                *ui.Panel
	panelOpCb                                                 *ui.Panel
	inputSYear, inputSMonth, inputSDay, inputSHour, inputSMin *ui.InputBox

	opTypeRoll, opTypeFly, opTypeAmaze *ui.OptionBox
	cbHostingType, cbFlyType           *ui.CheckBox
	cbAuto                             *ui.CheckBox

	opStartSplit, opStartMaoShan, opStartZhiRun *ui.OptionBox
	opStartSelf                                 *ui.OptionBox
	inputSelfJu                                 *ui.InputBox

	opHideGan0, opHideGan1 *ui.OptionBox

	btnCalc             *ui.TextButton
	btnNow              *ui.TextButton
	btnPreJu, btnNextJu *ui.TextButton

	opHourPan, opDayPan, opMonthPan, opYearPan *ui.OptionBox
	opDay2Pan                                  *ui.OptionBox

	year, month, day, hour, minute int
	solar                          *calendar.Solar
	qmParams                       qimen.QMParams
}

var uiQiMen *UIQiMen

func UIShowQiMen() *UIQiMen {
	if uiQiMen == nil {
		uiQiMen = NewUIQiMen()
		ui.ActiveUI(uiQiMen)
	}
	return uiQiMen
}
func UIHideQiMen() {
	if uiQiMen != nil {
		ui.CloseUI(uiQiMen)
		uiQiMen = nil
	}
}

func NewUIQiMen() *UIQiMen {
	p := &UIQiMen{
		BaseUI: ui.BaseUI{Visible: true},
		qmParams: qimen.QMParams{
			Type:        qimen.QMTypeRotating,
			HostingType: qimen.QMHostingType2,
			FlyType:     qimen.QMFlyTypeAllOrder,
			JuType:      qimen.QMJuTypeSplit,
			HideGanType: qimen.QMHideGanDutyDoorHour,
		},
	}
	px0, py0 := 16, 0
	h := 32
	p.panelSDate = ui.NewPanel(image.Rect(0, 0, px0+72*4+64, py0+h), nil)
	p.panelOpCb = ui.NewPanel(image.Rect(600, 0, px0+72*4+64, py0+h), nil)
	p.inputSYear = ui.NewInputBox(image.Rect(px0, py0, px0+64, py0+h))
	p.inputSMonth = ui.NewInputBox(image.Rect(px0+72, py0, px0+72+64, py0+h))
	p.inputSDay = ui.NewInputBox(image.Rect(px0+72*2, py0, px0+72*2+64, py0+h))
	p.inputSHour = ui.NewInputBox(image.Rect(px0+72*3, py0, px0+72*3+64, py0+h))
	p.inputSMin = ui.NewInputBox(image.Rect(px0+72*4, py0, px0+72*4+64, py0+h))
	opcbX := 346
	p.opTypeRoll = ui.NewOptionBox(px0+72*5+opcbX, py0, qimen.QMType[0])
	p.opTypeFly = ui.NewOptionBox(px0+72*6+opcbX, py0, qimen.QMType[1])
	p.opTypeAmaze = ui.NewOptionBox(px0+72*7+opcbX, py0, qimen.QMType[2])
	p.btnCalc = ui.NewTextButton(px0+72*8+opcbX, py0, "排局", colorWhite, true)
	p.btnNow = ui.NewTextButton(px0+72*8+36+opcbX, py0, "此时", colorWhite, true)

	py0 += 18
	p.cbHostingType = ui.NewCheckBox(px0+72*5+opcbX, py0, qimen.QMHostingType[qimen.QMHostingType28])
	p.cbFlyType = ui.NewCheckBox(px0+72*6+opcbX, py0, qimen.QMFlyType[qimen.QMFlyTypeAllOrder])
	p.btnPreJu = ui.NewTextButton(px0+72*8+opcbX, py0, "上一局", colorWhite, true)

	py0 += 18
	p.btnNextJu = ui.NewTextButton(px0+72*8+opcbX, py0, "下一局", colorWhite, true)

	p.inputSelfJu = ui.NewInputBox(image.Rect(px0+72*4, py0, px0+72*4+64, py0+h))
	p.opStartSelf = ui.NewOptionBox(px0+72*4+opcbX, py0, qimen.QMJuType[qimen.QMJuTypeSelf])

	p.opStartSplit = ui.NewOptionBox(px0+72*5+opcbX, py0, qimen.QMJuType[qimen.QMJuTypeSplit])
	p.opStartMaoShan = ui.NewOptionBox(px0+72*6+opcbX, py0, qimen.QMJuType[qimen.QMJuTypeMaoShan])
	p.opStartZhiRun = ui.NewOptionBox(px0+72*7+opcbX, py0, qimen.QMJuType[qimen.QMJuTypeZhiRun])

	py0 += 18
	p.opHideGan0 = ui.NewOptionBox(px0+72*5+opcbX, py0, qimen.QMHideGanType[qimen.QMHideGanDutyDoorHour])
	p.opHideGan1 = ui.NewOptionBox(px0+72*7+opcbX, py0, qimen.QMHideGanType[qimen.QMHideGanDoorHomeGan])
	py0 += 18
	p.opHourPan = ui.NewOptionBox(px0+72*5+opcbX, py0, "时家")
	p.opDayPan = ui.NewOptionBox(px0+72*6+opcbX, py0, "日家")
	p.opMonthPan = ui.NewOptionBox(px0+72*7+opcbX, py0, "月家")
	p.opYearPan = ui.NewOptionBox(px0+72*8+opcbX, py0, "年家")
	py0 += 18
	p.opDay2Pan = ui.NewOptionBox(px0+72*6+opcbX, py0, "_日家2")
	p.cbAuto = ui.NewCheckBox(px0+72*8+opcbX, py0, "自动")

	p.AddChildren(p.panelSDate, p.panelOpCb)
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
	p.panelSDate.AddChildren(p.inputSYear, p.inputSMonth, p.inputSDay, p.inputSHour, p.inputSMin)
	p.panelOpCb.AddChildren(p.btnCalc, p.btnNow, p.opTypeRoll, p.opTypeFly, p.opTypeAmaze,
		p.btnPreJu, p.btnNextJu, p.cbHostingType, p.cbFlyType,
		p.opHourPan, p.opDayPan, p.opMonthPan, p.opYearPan, p.opDay2Pan,
		p.opStartSplit, p.opStartMaoShan, p.opStartZhiRun, p.opStartSelf, p.inputSelfJu,
		p.cbAuto,
		p.opHideGan0, p.opHideGan1)

	ui.MakeOptionBoxGroup(p.opTypeRoll, p.opTypeFly, p.opTypeAmaze)
	p.opTypeRoll.Select()
	p.opTypeRoll.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.Type = qimen.QMTypeRotating
		p.Apply(p.solar)
	})
	p.opTypeFly.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.Type = qimen.QMTypeFly
		p.Apply(p.solar)
	})
	p.opTypeAmaze.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.Type = qimen.QMTypeAmaze
		p.opStartSplit.Select()
	})

	ui.MakeOptionBoxGroup(p.opStartSplit, p.opStartMaoShan, p.opStartZhiRun, p.opStartSelf)
	p.opStartSplit.Select()
	p.opStartSplit.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.JuType = qimen.QMJuTypeSplit
		p.Apply(p.solar)
	})
	p.opStartMaoShan.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.JuType = qimen.QMJuTypeMaoShan
		p.Apply(p.solar)
	})
	p.opStartZhiRun.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.JuType = qimen.QMJuTypeZhiRun
		p.Apply(p.solar)
	})
	p.opStartSelf.SetOnSelect(func(c *ui.OptionBox) {
		p.opStartSelf.Visible = false
		p.inputSelfJu.Visible = true
		p.inputSelfJu.SetFocused(true)
		p.qmParams.JuType = qimen.QMJuTypeSelf
	})

	ui.MakeOptionBoxGroup(p.opHideGan0, p.opHideGan1)
	p.opHideGan0.Select()
	p.opHideGan0.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.HideGanType = qimen.QMHideGanDutyDoorHour
		p.Apply(p.solar)
	})
	p.opHideGan1.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.HideGanType = qimen.QMHideGanDoorHomeGan
		p.Apply(p.solar)
	})

	ui.MakeOptionBoxGroup(p.opHourPan, p.opDayPan, p.opMonthPan, p.opYearPan, p.opDay2Pan)
	p.opHourPan.Select()
	p.opHourPan.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.YMDH = qimen.QMGameHour
		p.Apply(p.solar)
	})
	p.opDayPan.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.YMDH = qimen.QMGameDay
		p.Apply(p.solar)
	})
	p.opMonthPan.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.YMDH = qimen.QMGameMonth
		p.Apply(p.solar)
	})
	p.opYearPan.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.YMDH = qimen.QMGameYear
		p.Apply(p.solar)
	})
	p.opDay2Pan.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.YMDH = qimen.QMGameDay2
		p.Apply(p.solar)
	})
	p.opDay2Pan.Disabled = true
	p.opDay2Pan.Visible = false

	p.cbHostingType.SetChecked(false)
	p.cbHostingType.Visible = p.opTypeRoll.Selected()
	p.cbHostingType.SetOnCheckChanged(func(c *ui.CheckBox) {
		if c.Checked() {
			p.qmParams.HostingType = qimen.QMHostingType28
		} else {
			p.qmParams.HostingType = qimen.QMHostingType2
		}
		p.Apply(p.solar)
	})
	p.cbFlyType.SetChecked(true)
	p.cbFlyType.Visible = p.opTypeFly.Selected()
	p.cbFlyType.SetOnCheckChanged(func(c *ui.CheckBox) {
		if c.Checked() {
			p.qmParams.FlyType = qimen.QMFlyTypeAllOrder
		} else {
			p.qmParams.FlyType = qimen.QMFlyTypeLunarReverse
		}
		p.Apply(p.solar)
	})
	p.cbAuto.SetOnCheckChanged(func(c *ui.CheckBox) {
		ThisGame.autoMinute = c.Checked()
	})

	p.btnCalc.SetOnClick(func(b *ui.Button) {
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

		if p.qmParams.JuType == qimen.QMJuTypeSelf {
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
	p.btnNow.SetOnClick(func(b *ui.Button) {
		defer func() {
			s := recover()
			if s != nil {
				UIShowMsgBox(fmt.Sprintf("%s", s), "确定", "取消", nil, nil)
			}
		}()
		solar := calendar.NewSolarFromDate(time.Now())

		if p.qmParams.JuType == qimen.QMJuTypeSelf {
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
	p.btnPreJu.SetOnClick(func(b *ui.Button) {
		var solar *calendar.Solar
		switch p.qmParams.YMDH {
		case qimen.QMGameYear:
			solar = p.solar.NextYear(-1)
		case qimen.QMGameMonth:
			solar = p.solar.NextMonth(-1)
		case qimen.QMGameDay:
			solar = p.solar.NextDay(-1)
		case qimen.QMGameHour:
			solar = p.solar.NextHour(-2)
		}
		p.Apply(solar)
	})
	p.btnNextJu.SetOnClick(func(b *ui.Button) {
		var solar *calendar.Solar
		switch p.qmParams.YMDH {
		case qimen.QMGameYear:
			solar = p.solar.NextYear(1)
		case qimen.QMGameMonth:
			solar = p.solar.NextMonth(1)
		case qimen.QMGameDay:
			solar = p.solar.NextDay(1)
		case qimen.QMGameHour:
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

func (p *UIQiMen) Apply(solar *calendar.Solar) *qimen.QMGame {
	defer func() {
		s := recover()
		if s != nil {
			UIShowMsgBox(fmt.Sprintf("时间不对%s", s), "确定", "取消", nil, nil)
		}
	}()
	pan := qimen.NewQMGame(solar, p.qmParams)
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

	p.cbHostingType.Visible = p.qmParams.Type == qimen.QMTypeRotating
	p.cbFlyType.Visible = p.qmParams.Type == qimen.QMTypeFly

	shiPan := p.qmParams.YMDH == qimen.QMGameHour
	p.opStartSplit.Visible = shiPan // && p.qmParams.Type != qimen.QMTypeAmaze
	p.opStartMaoShan.Visible = shiPan && p.qmParams.Type != qimen.QMTypeAmaze
	p.opStartZhiRun.Visible = shiPan && p.qmParams.Type != qimen.QMTypeAmaze
	p.opStartSelf.Visible = !p.opStartSelf.Selected()
	p.inputSelfJu.Visible = p.opStartSelf.Selected()

	p.opHideGan0.Visible = p.qmParams.Type != qimen.QMTypeAmaze
	p.opHideGan1.Visible = p.qmParams.Type != qimen.QMTypeAmaze

	//fmt
	switch p.qmParams.YMDH {
	case qimen.QMGameHour:
		pan.ShowTimeGame()
	case qimen.QMGameDay:
		pan.ShowDayGame()
	case qimen.QMGameMonth:
		pan.ShowMonthGame()
	case qimen.QMGameYear:
		pan.ShowYearGame()
	}
	pan.CalBig6()
	return pan
}
func (p *UIQiMen) NextHour() *qimen.QMGame {
	solar := p.solar.NextHour(2)
	return p.Apply(solar)
}

func (p *UIQiMen) NextMinute() *qimen.QMGame {
	s := p.solar
	if s.GetMinute() == 59 {
		s = s.NextHour(1)
		s = calendar.NewSolar(s.GetYear(), s.GetMonth(), s.GetDay(), s.GetHour(), 0, 0)
	} else {
		s = calendar.NewSolar(s.GetYear(), s.GetMonth(), s.GetDay(), s.GetHour(), s.GetMinute()+1, 0)
	}
	return p.Apply(s)
}

func (p *UIQiMen) AutoTick() bool {
	return p.cbAuto.Checked()
}

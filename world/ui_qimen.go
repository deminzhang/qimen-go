package world

import (
	"errors"
	"fmt"
	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/6tail/lunar-go/SolarUtil"
	"image"
	"qimen/qimen"
	"qimen/ui"
	"strconv"
	"time"
)

type UIQiMen struct {
	ui.BaseUI
	panelSDate                                                *ui.Panel
	inputSYear, inputSMonth, inputSDay, inputSHour, inputSMin *ui.InputBox
	textLYear, textLMonth, textLDay, textLHour                *ui.InputBox
	textYearTB, textMonthTB, textDayTB, textHourTB            *ui.InputBox
	textJu                                                    *ui.TextBox

	opTypeRoll, opTypeFly, opTypeAmaze *ui.OptionBox
	cbHostingType, cbFlyType           *ui.CheckBox

	opStartSplit, opStartMaoShan, opStartZhiRun *ui.OptionBox
	opStartSelf                                 *ui.OptionBox
	inputSelfJu                                 *ui.InputBox

	opHideGan0, opHideGan1 *ui.OptionBox

	btnCalc             *ui.Button
	btnPreJu, btnNextJu *ui.Button

	opHourPan, opDayPan, opMonthPan, opYearPan *ui.OptionBox
	opDay2Pan                                  *ui.OptionBox

	textGong []*ui.TextBox
	zhiPan   []*ui.TextBox

	year, month, day, hour, minute int
	qmParams                       qimen.QMParams
}

var uiQiMen *UIQiMen

func UIShowQiMen(width, height int) {
	if uiQiMen == nil {
		uiQiMen = NewUIQiMen(width, height)
		ui.ActiveUI(uiQiMen)
	}
}
func UIHideQiMen() {
	if uiQiMen != nil {
		ui.CloseUI(uiQiMen)
		uiQiMen = nil
	}
}

func NewUIQiMen(width, height int) *UIQiMen {
	//cx, cy := width/2, height/2 //win center
	p := &UIQiMen{
		BaseUI: ui.BaseUI{Visible: true},
		qmParams: qimen.QMParams{
			Type:        qimen.QMTypeAmaze,
			HostingType: qimen.QMHostingType28,
			FlyType:     qimen.QMFlyTypeAllOrder,
			JuType:      qimen.QMJuTypeSplit,
			HideGanType: qimen.QMHideGanDutyDoorHour,
		},
	}
	px0, py0 := 32, 0
	h := 32
	p.panelSDate = ui.NewPanel(image.Rect(px0, py0, px0+72*4+64, py0+h))
	p.inputSYear = ui.NewInputBox(image.Rect(px0, py0, px0+64, py0+h))
	p.inputSMonth = ui.NewInputBox(image.Rect(px0+72, py0, px0+72+64, py0+h))
	p.inputSDay = ui.NewInputBox(image.Rect(px0+72*2, py0, px0+72*2+64, py0+h))
	p.inputSHour = ui.NewInputBox(image.Rect(px0+72*3, py0, px0+72*3+64, py0+h))
	p.inputSMin = ui.NewInputBox(image.Rect(px0+72*4, py0, px0+72*4+64, py0+h))
	p.opTypeRoll = ui.NewOptionBox(px0+72*5, py0+8, qimen.QMType[0])
	p.opTypeFly = ui.NewOptionBox(px0+72*6, py0+8, qimen.QMType[1])
	p.opTypeAmaze = ui.NewOptionBox(px0+72*7, py0+8, qimen.QMType[2])
	p.btnCalc = ui.NewButton(image.Rect(px0+72*8, py0, px0+72*8+64, py0+h), "排局")

	py0 += 32
	p.textLYear = ui.NewInputBox(image.Rect(px0, py0, px0+64, py0+h))
	p.textLMonth = ui.NewInputBox(image.Rect(px0+72, py0, px0+72+64, py0+h))
	p.textLDay = ui.NewInputBox(image.Rect(px0+72*2, py0, px0+72*2+64, py0+h))
	p.textLHour = ui.NewInputBox(image.Rect(px0+72*3, py0, px0+72*3+64, py0+h))
	p.cbHostingType = ui.NewCheckBox(px0+72*5, py0+8, qimen.QMHostingType[qimen.QMHostingType28])
	p.cbFlyType = ui.NewCheckBox(px0+72*6, py0+8, qimen.QMFlyType[qimen.QMFlyTypeAllOrder])
	p.btnPreJu = ui.NewButton(image.Rect(px0+72*8, py0, px0+72*8+64, py0+h), "上一局")

	py0 += 32
	p.textYearTB = ui.NewInputBox(image.Rect(px0, py0, px0+64, py0+h))
	p.textMonthTB = ui.NewInputBox(image.Rect(px0+72, py0, px0+72+64, py0+h))
	p.textDayTB = ui.NewInputBox(image.Rect(px0+72*2, py0, px0+72*2+64, py0+h))
	p.textHourTB = ui.NewInputBox(image.Rect(px0+72*3, py0, px0+72*3+64, py0+h))
	p.btnNextJu = ui.NewButton(image.Rect(px0+72*8, py0, px0+72*8+64, py0+h), "下一局")

	p.inputSelfJu = ui.NewInputBox(image.Rect(px0+72*4, py0, px0+72*4+64, py0+h))
	p.opStartSelf = ui.NewOptionBox(px0+72*4, py0+8, qimen.QMJuType[qimen.QMJuTypeSelf])

	p.opStartSplit = ui.NewOptionBox(px0+72*5, py0+8, qimen.QMJuType[qimen.QMJuTypeSplit])
	p.opStartMaoShan = ui.NewOptionBox(px0+72*6, py0+8, qimen.QMJuType[qimen.QMJuTypeMaoShan])
	p.opStartZhiRun = ui.NewOptionBox(px0+72*7, py0+8, qimen.QMJuType[qimen.QMJuTypeZhiRun])

	py0 += 32
	p.textJu = ui.NewTextBox(image.Rect(px0, py0, px0+72*4+64, py0+h*2))
	p.opHideGan0 = ui.NewOptionBox(px0+72*5, py0+8, qimen.QMHideGanType[qimen.QMHideGanDutyDoorHour])
	p.opHideGan1 = ui.NewOptionBox(px0+72*7, py0+8, qimen.QMHideGanType[qimen.QMHideGanDoorHomeGan])
	py0 += 32
	p.opHourPan = ui.NewOptionBox(px0+72*5, py0+8, "时家")
	p.opDayPan = ui.NewOptionBox(px0+72*6, py0+8, "日家")
	p.opMonthPan = ui.NewOptionBox(px0+72*7, py0+8, "月家")
	p.opYearPan = ui.NewOptionBox(px0+72*8, py0+8, "年家")
	py0 += 32
	p.opDay2Pan = ui.NewOptionBox(px0+72*6, py0+8, "_日家2")

	px4, py4 := 128, 256
	const gongWidth = 128
	gongOffset := [][]int{{0, 0},
		{1, 2}, {2, 0}, {0, 1},
		{0, 0}, {1, 1}, {2, 2},
		{2, 1}, {0, 2}, {1, 0},
	}
	p.textGong = append(p.textGong, nil)
	for i := 1; i <= 9; i++ {
		offX, offZ := gongOffset[i][0]*gongWidth, gongOffset[i][1]*gongWidth
		txtGong := ui.NewTextBox(image.Rect(px4+offX, py4+offZ, px4+offX+gongWidth, py4+offZ+gongWidth))
		txtGong.SetText(i)
		p.textGong = append(p.textGong, txtGong)
		p.AddChild(txtGong)
	}

	const zhiPanWidth = 48
	zhiPanLoc := [][]int{
		{2, 3, gongWidth, zhiPanWidth}, {1, 3, gongWidth, zhiPanWidth}, {0, 3, gongWidth, zhiPanWidth}, //亥子丑
		{0, 2, -zhiPanWidth, gongWidth}, {0, 1, -zhiPanWidth, gongWidth}, {0, 0, -zhiPanWidth, gongWidth}, //寅卯辰
		{0, 0, gongWidth, -zhiPanWidth}, {1, 0, gongWidth, -zhiPanWidth}, {2, 0, gongWidth, -zhiPanWidth}, //巳午未
		{3, 0, zhiPanWidth, gongWidth}, {3, 1, zhiPanWidth, gongWidth}, {3, 2, zhiPanWidth, gongWidth}, //申酉戌
		{2, 3, gongWidth, zhiPanWidth}, //亥
	}
	p.zhiPan = append(p.zhiPan, nil)
	for i := 1; i <= 12; i++ {
		offX, offZ := zhiPanLoc[i][0]*gongWidth, zhiPanLoc[i][1]*gongWidth
		minX := min(px4+offX, px4+offX+zhiPanLoc[i][2])
		maxX := max(px4+offX, px4+offX+zhiPanLoc[i][2])
		minY := min(py4+offZ, py4+offZ+zhiPanLoc[i][3])
		maxY := max(py4+offZ, py4+offZ+zhiPanLoc[i][3])
		txtZhi := ui.NewTextBox(image.Rect(minX, minY, maxX, maxY))
		txtZhi.DisableHScroll = true
		txtZhi.DisableVScroll = true
		p.zhiPan = append(p.zhiPan, txtZhi)
		p.AddChild(txtZhi)
		p.zhiPan[i].SetText(LunarUtil.ZHI[i])
	}

	p.AddChild(p.panelSDate)
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
	p.AddChildren(p.btnCalc, p.opTypeRoll, p.opTypeFly, p.opTypeAmaze)
	p.AddChildren(p.textLYear, p.textLMonth, p.textLDay, p.textLHour)
	p.AddChildren(p.btnPreJu, p.btnNextJu, p.cbHostingType, p.cbFlyType)
	p.AddChildren(p.opHourPan, p.opDayPan, p.opMonthPan, p.opYearPan, p.opDay2Pan)
	p.AddChildren(p.textYearTB, p.textMonthTB, p.textDayTB, p.textHourTB)
	p.AddChildren(p.opStartSplit, p.opStartMaoShan, p.opStartZhiRun, p.opStartSelf, p.inputSelfJu)
	p.AddChildren(p.opHideGan0, p.opHideGan1)

	p.AddChild(p.textJu)

	ui.MakeOptionBoxGroup(p.opTypeRoll, p.opTypeFly, p.opTypeAmaze)
	p.opTypeAmaze.Select()
	p.opTypeRoll.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.Type = qimen.QMTypeRotating
		p.Apply(p.year, p.month, p.day, p.hour, p.minute)
	})
	p.opTypeFly.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.Type = qimen.QMTypeFly
		p.Apply(p.year, p.month, p.day, p.hour, p.minute)
	})
	p.opTypeAmaze.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.Type = qimen.QMTypeAmaze
		p.opStartSplit.Select()
	})

	ui.MakeOptionBoxGroup(p.opStartSplit, p.opStartMaoShan, p.opStartZhiRun, p.opStartSelf)
	p.opStartSplit.Select()
	p.opStartSplit.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.JuType = qimen.QMJuTypeSplit
		p.Apply(p.year, p.month, p.day, p.hour, p.minute)
	})
	p.opStartMaoShan.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.JuType = qimen.QMJuTypeMaoShan
		p.Apply(p.year, p.month, p.day, p.hour, p.minute)
	})
	p.opStartZhiRun.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.JuType = qimen.QMJuTypeZhiRun
		p.Apply(p.year, p.month, p.day, p.hour, p.minute)
	})
	p.opStartSelf.SetOnSelect(func(c *ui.OptionBox) {
		p.opStartSelf.Visible = false
		p.inputSelfJu.Visible = true
		p.inputSelfJu.SetFocused(true)
		p.qmParams.JuType = qimen.QMJuTypeSelf
	})
	p.opStartZhiRun.Disabled = true

	ui.MakeOptionBoxGroup(p.opHideGan0, p.opHideGan1)
	p.opHideGan0.Select()
	p.opHideGan0.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.HideGanType = qimen.QMHideGanDutyDoorHour
		p.Apply(p.year, p.month, p.day, p.hour, p.minute)
	})
	p.opHideGan1.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.HideGanType = qimen.QMHideGanDoorHomeGan
		p.Apply(p.year, p.month, p.day, p.hour, p.minute)
	})

	ui.MakeOptionBoxGroup(p.opHourPan, p.opDayPan, p.opMonthPan, p.opYearPan, p.opDay2Pan)
	p.opHourPan.Select()
	p.opHourPan.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.YMDH = qimen.QMGameHour
		p.Apply(p.year, p.month, p.day, p.hour, p.minute)
	})
	p.opDayPan.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.YMDH = qimen.QMGameDay
		p.Apply(p.year, p.month, p.day, p.hour, p.minute)
	})
	p.opMonthPan.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.YMDH = qimen.QMGameMonth
		p.Apply(p.year, p.month, p.day, p.hour, p.minute)
	})
	p.opYearPan.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.YMDH = qimen.QMGameYear
		p.Apply(p.year, p.month, p.day, p.hour, p.minute)
	})
	p.opDay2Pan.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.YMDH = qimen.QMGameDay2
		p.Apply(p.year, p.month, p.day, p.hour, p.minute)
	})
	p.opDay2Pan.Disabled = true

	p.cbHostingType.SetChecked(false)
	p.cbHostingType.Visible = p.opTypeRoll.Selected()
	p.cbHostingType.SetOnCheckChanged(func(c *ui.CheckBox) {
		if c.Checked() {
			p.qmParams.HostingType = qimen.QMHostingType28
		} else {
			p.qmParams.HostingType = qimen.QMHostingType2
		}
		p.Apply(p.year, p.month, p.day, p.hour, p.minute)
	})
	p.cbFlyType.SetChecked(true)
	p.cbFlyType.Visible = p.opTypeFly.Selected()
	p.cbFlyType.SetOnCheckChanged(func(c *ui.CheckBox) {
		if c.Checked() {
			p.qmParams.FlyType = qimen.QMFlyTypeAllOrder
		} else {
			p.qmParams.FlyType = qimen.QMFlyTypeLunarReverse
		}
		p.Apply(p.year, p.month, p.day, p.hour, p.minute)
	})

	p.btnCalc.SetOnClick(func(b *ui.Button) {
		year, _ := strconv.Atoi(p.inputSYear.Text())
		month, _ := strconv.Atoi(p.inputSMonth.Text())
		day, _ := strconv.Atoi(p.inputSDay.Text())
		hour, _ := strconv.Atoi(p.inputSHour.Text())
		minute, _ := strconv.Atoi(p.inputSMin.Text())

		if p.qmParams.JuType == qimen.QMJuTypeSelf {
			ju, err := strconv.Atoi(p.inputSelfJu.Text())
			if err != nil || !((ju >= 1 && ju <= 9) || (ju >= -9 && ju <= -1)) {
				UIShowMsgBox("局数不对,限([-9,-1],[1,9])", "确定", "取消", nil, nil)
				return
			}
			p.qmParams.SelfJu = ju
		} else {
			p.qmParams.SelfJu = 0
		}

		p.Apply(year, month, day, hour, minute)
	})
	p.btnPreJu.SetOnClick(func(b *ui.Button) {
		year, month, day, hour, minute := p.year, p.month, p.day, p.hour, p.minute

		switch p.qmParams.YMDH {
		case qimen.QMGameYear:
			year--
			if 1582 == year && 10 == month && day > 4 && day < 14 {
				day = 4
			}
		case qimen.QMGameMonth:
			month--
			if month == 0 {
				month = 12
				year--
			}
			if day > SolarUtil.GetDaysOfMonth(year, month) {
				day = SolarUtil.GetDaysOfMonth(year, month)
			}
		case qimen.QMGameDay:
			day--
			if day == 0 {
				month--
				if month == 0 {
					month = 12
					year--
				}
				day = SolarUtil.GetDaysOfMonth(year, month)
			}
		case qimen.QMGameHour:
			hour -= 2
			if hour < 0 {
				hour += 24
				day--
				if day == 0 {
					month--
					if month == 0 {
						month = 12
						year--
					}
					day = SolarUtil.GetDaysOfMonth(year, month)
				}
			}
		}
		if 1582 == year && 10 == month && day == 14 {
			day = 4
		}
		p.Apply(year, month, day, hour, minute)
	})
	p.btnNextJu.SetOnClick(func(b *ui.Button) {
		year, month, day, hour, minute := p.year, p.month, p.day, p.hour, p.minute

		switch p.qmParams.YMDH {
		case qimen.QMGameYear:
			year++
		case qimen.QMGameMonth:
			month++
			if month > 12 {
				year++
				month = 1
				if day > SolarUtil.GetDaysOfMonth(year, month) {
					day = SolarUtil.GetDaysOfMonth(year, month)
				}
			}
		case qimen.QMGameDay:
			day++
			if day > SolarUtil.GetDaysOfMonth(year, month) {
				day = 1
				month++
				if month > 12 {
					year++
					month = 1
				}
			}
		case qimen.QMGameHour:
			hour += 2
			if hour > 23 {
				hour -= 24
				day++
				if day > SolarUtil.GetDaysOfMonth(year, month) {
					day = 1
					month++
					if month > 12 {
						month = 1
						year++
					}
				}
			}
		}
		if 1582 == year && 10 == month && day == 5 {
			day = 15
		}
		if day > SolarUtil.GetDaysOfMonth(year, month) {
			day = SolarUtil.GetDaysOfMonth(year, month)
		}
		p.Apply(year, month, day, hour, minute)
	})

	p.textLYear.Editable = false
	p.textLMonth.Editable = false
	p.textLDay.Editable = false
	p.textLHour.Editable = false

	p.textYearTB.Editable = false
	p.textMonthTB.Editable = false
	p.textDayTB.Editable = false
	p.textHourTB.Editable = false

	t := time.Now()
	p.Apply(t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute())

	uiQiMen = p
	return p
}

func (p *UIQiMen) OnClose() {
	uiQiMen = nil
}

func (p *UIQiMen) checkDate(year, month, day, hour, minute int) error {
	if month < 1 || month > 12 {
		return errors.New(fmt.Sprintf("wrong month %v", month))
	}
	if day < 1 || day > 31 {
		return errors.New(fmt.Sprintf("wrong day %v", day))
	}
	if 1582 == year && 10 == month {
		if day > 4 && day < 15 {
			return errors.New(fmt.Sprintf("wrong solar year %v month %v day %v", year, month, day))
		}
	} else {
		if day > SolarUtil.GetDaysOfMonth(year, month) {
			return errors.New(fmt.Sprintf("wrong solar year %v month %v day %v", year, month, day))
		}
	}
	if hour < 0 || hour > 23 {
		return errors.New(fmt.Sprintf("wrong hour %v", hour))
	}
	if minute < 0 || minute > 59 {
		return errors.New(fmt.Sprintf("wrong minute %v", minute))
	}
	return nil
}
func (p *UIQiMen) Apply(year, month, day, hour, minute int) {
	pan, err := qimen.NewPan(year, month, day, hour, minute, p.qmParams)
	if err != nil {
		UIShowMsgBox("时间不对", "确定", "取消", nil, nil)
		return
	}
	p.year, p.month, p.day, p.hour, p.minute = year, month, day, hour, minute
	//pan.DayArr
	p.inputSYear.SetText(pan.SolarYear)
	p.inputSMonth.SetText(pan.SolarMonth)
	p.inputSDay.SetText(pan.SolarDay)
	p.inputSHour.SetText(pan.SolarHour)
	p.inputSMin.SetText(pan.SolarMinute)

	p.textLYear.SetText(pan.LunarYearC)
	p.textLMonth.SetText(pan.LunarMonthC)
	p.textLDay.SetText(pan.LunarDayC)
	p.textLHour.SetText(pan.LunarHourC)

	p.textYearTB.SetText(pan.YearTB)
	p.textMonthTB.SetText(pan.MonthTB)
	p.textDayTB.SetText(pan.DayTB)
	p.textHourTB.SetText(pan.HourTB)

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
		p.ShowHourGame(pan)
	case qimen.QMGameDay:
		p.ShowDayGame(pan)
	case qimen.QMGameMonth:
		p.ShowMonthGame(pan)
	case qimen.QMGameYear:
		p.ShowYearGame(pan)
	}
}

func (p *UIQiMen) show9Gong(pp *qimen.QMPan) {
	kongWang := qimen.KongWang[pp.Xun]
	for i := 1; i <= 9; i++ {
		g := pp.Gongs[i]
		var hosting = "    "
		if pp.RollHosting > 0 && i == pp.DutyStarPos {
			hosting = " 禽 "
		}
		var empty, horse = "  ", "  "
		for _, zhi := range kongWang {
			if qimen.ZhiGong9[zhi] == i {
				empty = "〇" //"空亡"
				break
			}
		}
		if qimen.ZhiGong9[pp.Horse] == i {
			horse = "马"
		}
		door := g.Door + qimen.Door0
		if g.Door == "" {
			door = ""
		}
		p.textGong[i].Text = fmt.Sprintf("\n  %s  %s    %s\n\n"+
			"%s %s %s%s%s\n\n"+
			"%s    %s    %s\n\n"+
			"      %s%s",
			empty, g.God, horse,
			g.PathGan, g.HideGan, qimen.Star0+g.Star, hosting, g.GuestGan,
			g.PathZhi, door, g.HostGan,
			qimen.Diagrams9(i), LunarUtil.NUMBER[i])
	}
}
func (p *UIQiMen) ShowHourGame(pan *qimen.QMGame) {
	pp := pan.HourPan
	var juName string
	if pp.Ju < 0 {
		juName = fmt.Sprintf("阴%d局", -pp.Ju)
	} else {
		juName = fmt.Sprintf("阳%d局", pp.Ju)
	}
	jieQi := pan.Lunar.GetPrevJieQi()
	jieQiNext := pan.Lunar.GetNextJieQi()
	jie := pan.Lunar.GetPrevJie()
	qi := pan.Lunar.GetPrevQi()
	juText := fmt.Sprintf("%s%s %s%s"+
		"\n%s %s %s遁%s 值符%s落%d宫 值使%s落%d宫"+
		"\n%s月建%s %s月将%s",
		jieQi.GetName(), jieQi.GetSolar().ToYmdHms(), jieQiNext.GetName(), jieQiNext.GetSolar().ToYmdHms(),
		qimen.Yuan3Name[pp.Yuan3], juName, pp.Xun, qimen.HideJia[pp.Xun],
		pp.DutyStar, pp.DutyStarPos, pp.DutyDoor, pp.DutyDoorPos,
		jie.GetName(), pan.YueJian, qi.GetName(), pan.YueJiang,
	)
	p.textJu.SetText(juText)
	p.show9Gong(pp)

	p.show12Gong(pan)
}

// show12Gong 大六壬 月将落时支 顺布余支 天三门兮地四户
func (p *UIQiMen) show12Gong(pan *qimen.QMGame) {
	yueJiangIdx := pan.YueJiangZhiIdx
	yueJianIdx := pan.YueJianZhiIdx
	shiZhiIdx := pan.HourZhiIdx
	for i := shiZhiIdx; i < shiZhiIdx+12; i++ {
		//var s []string
		js := LunarUtil.ZHI[qimen.Idx12[yueJiangIdx]]
		g := fmt.Sprintf("%s", qimen.YueJiangName[js])
		var j, h, b, bs string
		if i == shiZhiIdx {
			j = "月将"
		}
		z := LunarUtil.ZHI[qimen.Idx12[i]]
		if z == pan.HourPan.Horse {
			h = "驿马"
		}
		b = LunarUtil.ZHI[qimen.Idx12[yueJianIdx+i-shiZhiIdx]]
		bs = qimen.BuildStar(1 + i - shiZhiIdx)
		switch qimen.Idx12[i] {
		case 1, 2, 6, 7, 8, 12:
			p.zhiPan[qimen.Idx12[i]].SetText(fmt.Sprintf("%s %s %s\n%s   %s\n%s", g, js, j, bs, b, h))
		default:
			p.zhiPan[qimen.Idx12[i]].SetText(fmt.Sprintf("%s\n  %s\n%s\n----\n%s%s\n%s", g, js, j, bs, b, h))
		}
		yueJiangIdx++
	}
}
func (p *UIQiMen) noShow12Gong() {
	for i := 1; i <= 12; i++ {
		p.zhiPan[i].SetText("")
	}
}

func (p *UIQiMen) ShowDayGame(pan *qimen.QMGame) {
	pp := pan.DayPan
	var juName string
	if pp.Ju < 0 {
		juName = fmt.Sprintf("阴%d局", -pp.Ju)
	} else {
		juName = fmt.Sprintf("阳%d局", pp.Ju)
	}
	jieQi := pan.Lunar.GetPrevJieQi()
	jieQiNext := pan.Lunar.GetNextJieQi()
	juText := fmt.Sprintf("日家%s%s %s%s"+
		"\n%s %s %s遁%s 值符%s落%d宫 值使%s落%d宫",
		jieQi.GetName(), jieQi.GetSolar().ToYmdHms(), jieQiNext.GetName(), jieQiNext.GetSolar().ToYmdHms(),
		qimen.Yuan3Name[pp.Yuan3], juName, pp.Xun, qimen.HideJia[pp.Xun],
		pp.DutyStar, pp.DutyStarPos, pp.DutyDoor, pp.DutyDoorPos,
	)
	p.textJu.SetText(juText)
	p.show9Gong(pp)
	p.noShow12Gong()
}

func (p *UIQiMen) ShowDayGame2(pan *qimen.QMGame) {

}

func (p *UIQiMen) ShowMonthGame(pan *qimen.QMGame) {
	pp := pan.MonthPan
	var juName string
	if pp.Ju < 0 {
		juName = fmt.Sprintf("阴%d局", -pp.Ju)
	} else {
		juName = fmt.Sprintf("阳%d局", pp.Ju)
	}
	juText := fmt.Sprintf("月家"+
		"\n%s %s %s遁%s 值符%s落%d宫 值使%s落%d宫",
		qimen.Yuan3Name[pp.Yuan3], juName, pp.Xun, qimen.HideJia[pp.Xun],
		pp.DutyStar, pp.DutyStarPos, pp.DutyDoor, pp.DutyDoorPos,
	)
	p.textJu.SetText(juText)
	p.show9Gong(pp)
	p.noShow12Gong()
}

func (p *UIQiMen) ShowYearGame(pan *qimen.QMGame) {
	pp := pan.YearPan
	var juName string
	if pp.Ju < 0 {
		juName = fmt.Sprintf("阴%d局", -pp.Ju)
	} else {
		juName = fmt.Sprintf("阳%d局", pp.Ju)
	}
	juText := fmt.Sprintf("年家"+
		"\n%s %s %s遁%s 值符%s落%d宫 值使%s落%d宫",
		qimen.Yuan3Name[pp.Yuan3], juName, pp.Xun, qimen.HideJia[pp.Xun],
		pp.DutyStar, pp.DutyStarPos, pp.DutyDoor, pp.DutyDoorPos,
	)
	p.textJu.SetText(juText)
	p.show9Gong(pp)
	p.noShow12Gong()
}

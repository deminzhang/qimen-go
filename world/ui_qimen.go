package world

import (
	"errors"
	"fmt"
	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/6tail/lunar-go/SolarUtil"
	"github.com/6tail/lunar-go/calendar"
	"image"
	"qimen/qimen"
	"qimen/ui"
	"strconv"
	"time"
)

const (
	gongWidth   = 128
	zhiPanWidth = 48
)

var gongOffset = [][]int{{0, 0},
	{1, 2}, {2, 0}, {0, 1},
	{0, 0}, {1, 1}, {2, 2},
	{2, 1}, {0, 2}, {1, 0},
}

type UIQiMen struct {
	ui.BaseUI
	pan                                                       *qimen.QMGame
	panelSDate                                                *ui.Panel
	inputSYear, inputSMonth, inputSDay, inputSHour, inputSMin *ui.InputBox

	opTypeRoll, opTypeFly, opTypeAmaze *ui.OptionBox
	cbHostingType, cbFlyType           *ui.CheckBox

	opStartSplit, opStartMaoShan, opStartZhiRun *ui.OptionBox
	opStartSelf                                 *ui.OptionBox
	inputSelfJu                                 *ui.InputBox

	opHideGan0, opHideGan1 *ui.OptionBox

	btnCalc             *ui.Button
	btnNow              *ui.Button
	btnPreJu, btnNextJu *ui.Button
	btnBirth            *ui.Button

	opHourPan, opDayPan, opMonthPan, opYearPan *ui.OptionBox
	cbBaZi                                     *ui.CheckBox
	opDay2Pan                                  *ui.OptionBox

	//zhiPan []*ui.TextBox
	gong12 [12 + 1]Gong12

	year, month, day, hour, minute int
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
	px0, py0 := 32, 0
	h := 32
	p.panelSDate = ui.NewPanel(image.Rect(px0, py0, px0+72*4+64, py0+h), nil)
	p.inputSYear = ui.NewInputBox(image.Rect(px0, py0, px0+64, py0+h))
	p.inputSMonth = ui.NewInputBox(image.Rect(px0+72, py0, px0+72+64, py0+h))
	p.inputSDay = ui.NewInputBox(image.Rect(px0+72*2, py0, px0+72*2+64, py0+h))
	p.inputSHour = ui.NewInputBox(image.Rect(px0+72*3, py0, px0+72*3+64, py0+h))
	p.inputSMin = ui.NewInputBox(image.Rect(px0+72*4, py0, px0+72*4+64, py0+h))
	p.opTypeRoll = ui.NewOptionBox(px0+72*5, py0+8, qimen.QMType[0])
	p.opTypeFly = ui.NewOptionBox(px0+72*6, py0+8, qimen.QMType[1])
	p.opTypeAmaze = ui.NewOptionBox(px0+72*7, py0+8, qimen.QMType[2])
	p.btnCalc = ui.NewButton(image.Rect(px0+72*8, py0, px0+72*8+64, py0+h), "排局")
	p.btnNow = ui.NewButton(image.Rect(px0+72*9, py0, px0+72*9+64, py0+h), "此时")

	py0 += 32
	p.cbHostingType = ui.NewCheckBox(px0+72*5, py0+8, qimen.QMHostingType[qimen.QMHostingType28])
	p.cbFlyType = ui.NewCheckBox(px0+72*6, py0+8, qimen.QMFlyType[qimen.QMFlyTypeAllOrder])
	p.btnPreJu = ui.NewButton(image.Rect(px0+72*8, py0, px0+72*8+64, py0+h), "上一局")

	py0 += 32
	p.btnNextJu = ui.NewButton(image.Rect(px0+72*8, py0, px0+72*8+64, py0+h), "下一局")
	p.btnBirth = ui.NewButton(image.Rect(px0+72*9, py0, px0+72*9+64, py0+h), "生日")

	p.inputSelfJu = ui.NewInputBox(image.Rect(px0+72*4, py0, px0+72*4+64, py0+h))
	p.opStartSelf = ui.NewOptionBox(px0+72*4, py0+8, qimen.QMJuType[qimen.QMJuTypeSelf])

	p.opStartSplit = ui.NewOptionBox(px0+72*5, py0+8, qimen.QMJuType[qimen.QMJuTypeSplit])
	p.opStartMaoShan = ui.NewOptionBox(px0+72*6, py0+8, qimen.QMJuType[qimen.QMJuTypeMaoShan])
	p.opStartZhiRun = ui.NewOptionBox(px0+72*7, py0+8, qimen.QMJuType[qimen.QMJuTypeZhiRun])

	py0 += 32
	p.opHideGan0 = ui.NewOptionBox(px0+72*5, py0+8, qimen.QMHideGanType[qimen.QMHideGanDutyDoorHour])
	p.opHideGan1 = ui.NewOptionBox(px0+72*7, py0+8, qimen.QMHideGanType[qimen.QMHideGanDoorHomeGan])
	py0 += 32
	p.opHourPan = ui.NewOptionBox(px0+72*5, py0+8, "时家")
	p.opDayPan = ui.NewOptionBox(px0+72*6, py0+8, "日家")
	p.opMonthPan = ui.NewOptionBox(px0+72*7, py0+8, "月家")
	p.opYearPan = ui.NewOptionBox(px0+72*8, py0+8, "年家")
	p.cbBaZi = ui.NewCheckBox(px0+72*9, py0+8, "四柱")
	py0 += 32
	p.opDay2Pan = ui.NewOptionBox(px0+72*6, py0+8, "_日家2")

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
	p.AddChildren(p.btnCalc, p.btnNow, p.opTypeRoll, p.opTypeFly, p.opTypeAmaze)
	p.AddChildren(p.btnPreJu, p.btnNextJu, p.cbHostingType, p.cbFlyType)
	p.AddChildren(p.opHourPan, p.opDayPan, p.opMonthPan, p.opYearPan, p.opDay2Pan)
	p.AddChildren(p.opStartSplit, p.opStartMaoShan, p.opStartZhiRun, p.opStartSelf, p.inputSelfJu)
	p.AddChildren(p.cbBaZi)
	p.AddChildren(p.opHideGan0, p.opHideGan1)
	p.AddChild(p.btnBirth)

	ui.MakeOptionBoxGroup(p.opTypeRoll, p.opTypeFly, p.opTypeAmaze)
	p.opTypeRoll.Select()
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
	p.opDay2Pan.Visible = false

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
	p.btnNow.SetOnClick(func(b *ui.Button) {
		now := time.Now()
		year := now.Year()
		month := int(now.Month())
		day := now.Day()
		hour := now.Hour()
		minute := now.Minute()

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
	p.btnBirth.SetOnClick(func(b *ui.Button) {
		UIShowSelect()
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
	defer func() {
		s := recover()
		if s != nil {
			UIShowMsgBox(fmt.Sprintf("时间不对%s", s), "确定", "取消", nil, nil)
		}
	}()
	solar := calendar.NewSolar(year, month, day, hour, minute, 0)
	pan := qimen.NewPan(solar, p.qmParams)
	p.pan = pan
	p.year, p.month, p.day, p.hour, p.minute = year, month, day, hour, minute
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
		p.ShowHourGame(pan)
	case qimen.QMGameDay:
		p.ShowDayGame(pan)
	case qimen.QMGameMonth:
		p.ShowMonthGame(pan)
	case qimen.QMGameYear:
		p.ShowYearGame(pan)
	}
}

func (p *UIQiMen) ShowHourGame(pan *qimen.QMGame) {
	pp := pan.HourPan
	pan.ShowPan = pp
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
		qimen.Star0+pp.DutyStar, pp.DutyStarPos, pp.DutyDoor+qimen.Door0, pp.DutyDoorPos,
		jie.GetName(), pan.YueJian, qi.GetName(), pan.YueJiang,
	)
	pp.JuText = juText

	p.calBig6(pan)
}

// 大六壬 月将落时支 顺布余支 天三门兮地四户
func (p *UIQiMen) calBig6(pan *qimen.QMGame) {
	yueJiangIdx := pan.YueJiangZhiIdx
	yueJianIdx := pan.YueJianZhiIdx
	shiZhiIdx := pan.Lunar.GetTimeZhiIndex() + 1
	for i := shiZhiIdx; i < shiZhiIdx+12; i++ {
		js := LunarUtil.ZHI[qimen.Idx12[yueJiangIdx]]
		g := fmt.Sprintf("%s", qimen.YueJiangName[js])
		z := LunarUtil.ZHI[qimen.Idx12[i]]
		bs := qimen.BuildStar(1 + i - shiZhiIdx)

		g12 := &p.gong12[qimen.Idx12[i]]
		g12.Idx = qimen.Idx12[i]
		g12.JiangZhi = js
		g12.Jiang = g
		g12.IsJiang = i == shiZhiIdx
		g12.JianZhi = LunarUtil.ZHI[qimen.Idx12[yueJianIdx+i-shiZhiIdx]]
		g12.Jian = bs
		g12.IsJian = bs == "建"
		g12.IsHorse = z == pan.HourPan.Horse
		g12.IsSkyHorse = g == "太冲"

		yueJiangIdx++
	}
}

func (p *UIQiMen) ShowDayGame(pan *qimen.QMGame) {
	pp := pan.DayPan
	pan.ShowPan = pp
	var juName string
	if pp.Ju < 0 {
		juName = fmt.Sprintf("阴%d局", -pp.Ju)
	} else {
		juName = fmt.Sprintf("阳%d局", pp.Ju)
	}
	jieQi := pan.Lunar.GetPrevJieQi()
	jieQiNext := pan.Lunar.GetNextJieQi()
	juText := fmt.Sprintf("%s%s %s%s"+
		"\n%s %s %s遁%s 值符%s落%d宫 值使%s落%d宫",
		jieQi.GetName(), jieQi.GetSolar().ToYmdHms(), jieQiNext.GetName(), jieQiNext.GetSolar().ToYmdHms(),
		qimen.Yuan3Name[pp.Yuan3], juName, pp.Xun, qimen.HideJia[pp.Xun],
		qimen.Star0+pp.DutyStar, pp.DutyStarPos, pp.DutyDoor+qimen.Door0, pp.DutyDoorPos,
	)
	pp.JuText = juText
}

func (p *UIQiMen) ShowDayGame2(pan *qimen.QMGame) {
	//TODO 	太乙日家
}

func (p *UIQiMen) ShowMonthGame(pan *qimen.QMGame) {
	pp := pan.MonthPan
	pan.ShowPan = pp
	var juName string
	if pp.Ju < 0 {
		juName = fmt.Sprintf("阴%d局", -pp.Ju)
	} else {
		juName = fmt.Sprintf("阳%d局", pp.Ju)
	}
	juText := fmt.Sprintf("月家"+
		"\n%s %s %s遁%s 值符%s落%d宫 值使%s落%d宫",
		qimen.Yuan3Name[pp.Yuan3], juName, pp.Xun, qimen.HideJia[pp.Xun],
		qimen.Star0+pp.DutyStar, pp.DutyStarPos, pp.DutyDoor+qimen.Door0, pp.DutyDoorPos,
	)
	pp.JuText = juText
}

func (p *UIQiMen) ShowYearGame(pan *qimen.QMGame) {
	pp := pan.YearPan
	pan.ShowPan = pp
	var juName string
	if pp.Ju < 0 {
		juName = fmt.Sprintf("阴%d局", -pp.Ju)
	} else {
		juName = fmt.Sprintf("阳%d局", pp.Ju)
	}
	y9 := qimen.GetYear9Yun(pan.Lunar.GetYear())
	d8 := qimen.Diagrams9(y9)
	y3y9 := fmt.Sprintf("三元九运:%s%s%s%s%s", qimen.Yuan3Name[pp.Yuan3],
		LunarUtil.NUMBER[y9], qimen.Gong9Color[y9], d8, qimen.DiagramsWuxing[d8])
	juText := fmt.Sprintf("年家 黄帝纪元:%d %s"+
		"\n%s %s %s遁%s 值符%s落%d宫 值使%s落%d宫",
		qimen.GetHuangDiYear(pan.Lunar.GetYear()), y3y9,
		qimen.Yuan3Name[pp.Yuan3], juName, pp.Xun, qimen.HideJia[pp.Xun],
		qimen.Star0+pp.DutyStar, pp.DutyStarPos, pp.DutyDoor+qimen.Door0, pp.DutyDoorPos,
	)
	pp.JuText = juText
}

func (p *UIQiMen) IsShowBaZi() bool {
	return p.cbBaZi.Checked()
}

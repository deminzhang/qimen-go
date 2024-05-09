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
	panelSDate  *ui.Panel
	inputSYear  *ui.InputBox
	inputSMonth *ui.InputBox
	inputSDay   *ui.InputBox
	inputSHour  *ui.InputBox
	inputSMin   *ui.InputBox

	textLYear   *ui.InputBox
	textLMonth  *ui.InputBox
	textLDay    *ui.InputBox
	textLHour   *ui.InputBox
	textYearTB  *ui.InputBox
	textMonthTB *ui.InputBox
	textDayTB   *ui.InputBox
	textHourTB  *ui.InputBox
	textJu      *ui.TextBox

	opTypeRoll    *ui.OptionBox
	opTypeFly     *ui.OptionBox
	opTypeAmaze   *ui.OptionBox
	cbHostingType *ui.CheckBox
	cbFlyType     *ui.CheckBox

	opStartSplit *ui.OptionBox
	opStartMao   *ui.OptionBox
	opStartZhi   *ui.OptionBox
	opStartSelf  *ui.OptionBox

	opHideGan0 *ui.OptionBox
	opHideGan1 *ui.OptionBox

	btnCalc      *ui.Button
	btnPreHour2  *ui.Button
	btnNextHour2 *ui.Button

	opHourPan  *ui.OptionBox
	opDayPan   *ui.OptionBox
	opMonthPan *ui.OptionBox
	opYearPan  *ui.OptionBox

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
			StartType:   qimen.QMStartTypeSplit,
			HideGanType: 0,
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
	p.btnPreHour2 = ui.NewButton(image.Rect(px0+72*8, py0, px0+72*8+64, py0+h), "上一局")

	py0 += 32
	p.textYearTB = ui.NewInputBox(image.Rect(px0, py0, px0+64, py0+h))
	p.textMonthTB = ui.NewInputBox(image.Rect(px0+72, py0, px0+72+64, py0+h))
	p.textDayTB = ui.NewInputBox(image.Rect(px0+72*2, py0, px0+72*2+64, py0+h))
	p.textHourTB = ui.NewInputBox(image.Rect(px0+72*3, py0, px0+72*3+64, py0+h))
	p.btnNextHour2 = ui.NewButton(image.Rect(px0+72*8, py0, px0+72*8+64, py0+h), "下一局")

	p.opStartSplit = ui.NewOptionBox(px0+72*5, py0+8, qimen.QMStartType[qimen.QMStartTypeSplit])
	p.opStartMao = ui.NewOptionBox(px0+72*6, py0+8, qimen.QMStartType[qimen.QMStartTypeMao])
	p.opStartZhi = ui.NewOptionBox(px0+72*7, py0+8, qimen.QMStartType[qimen.QMStartTypeZhi])
	p.opStartSelf = ui.NewOptionBox(px0+72*7, py0+8, qimen.QMStartType[qimen.QMStartTypeSelf])

	py0 += 32
	p.textJu = ui.NewTextBox(image.Rect(px0, py0, px0+72*4+64, py0+h*2))
	p.opHideGan0 = ui.NewOptionBox(px0+72*5, py0+8, qimen.QMHideGanType[qimen.QMHideGanDutyDoorPos])
	p.opHideGan1 = ui.NewOptionBox(px0+72*7, py0+8, qimen.QMHideGanType[qimen.QMHideGanDoorHome])
	py0 += 32
	p.opHourPan = ui.NewOptionBox(px0+72*5, py0+8, "时家")
	p.opDayPan = ui.NewOptionBox(px0+72*6, py0+8, "_日家")
	p.opMonthPan = ui.NewOptionBox(px0+72*7, py0+8, "月家")
	p.opYearPan = ui.NewOptionBox(px0+72*8, py0+8, "年家")

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
	p.inputSYear.DefaultText = "year"
	p.inputSMonth.DefaultText = "month"
	p.inputSDay.DefaultText = "day"
	p.inputSHour.DefaultText = "hour"
	p.inputSMin.DefaultText = "minute"
	p.panelSDate.AddChild(p.inputSYear)
	p.panelSDate.AddChild(p.inputSMonth)
	p.panelSDate.AddChild(p.inputSDay)
	p.panelSDate.AddChild(p.inputSHour)
	p.panelSDate.AddChild(p.inputSMin)
	p.AddChild(p.btnCalc)
	p.AddChild(p.opTypeRoll)
	p.AddChild(p.opTypeFly)
	p.AddChild(p.opTypeAmaze)

	p.AddChild(p.textLYear)
	p.AddChild(p.textLMonth)
	p.AddChild(p.textLDay)
	p.AddChild(p.textLHour)

	p.AddChild(p.btnPreHour2)
	p.AddChild(p.btnNextHour2)
	p.AddChild(p.cbHostingType)
	p.AddChild(p.cbFlyType)

	p.AddChild(p.opHourPan)
	p.AddChild(p.opDayPan)
	p.AddChild(p.opMonthPan)
	p.AddChild(p.opYearPan)

	p.AddChild(p.textYearTB)
	p.AddChild(p.textMonthTB)
	p.AddChild(p.textDayTB)
	p.AddChild(p.textHourTB)
	p.AddChild(p.opStartSplit)
	p.AddChild(p.opStartMao)
	p.AddChild(p.opStartZhi)
	p.AddChild(p.opStartSelf)
	p.AddChild(p.opHideGan0)
	p.AddChild(p.opHideGan1)

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
		p.Apply(p.year, p.month, p.day, p.hour, p.minute)
	})

	ui.MakeOptionBoxGroup(p.opStartSplit, p.opStartMao, p.opStartZhi, p.opStartSelf)
	p.opStartSplit.Select()
	p.opStartSplit.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.StartType = qimen.QMStartTypeSplit
		p.Apply(p.year, p.month, p.day, p.hour, p.minute)
	})
	p.opStartMao.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.StartType = qimen.QMStartTypeMao
		p.Apply(p.year, p.month, p.day, p.hour, p.minute)
	})
	p.opStartZhi.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.StartType = qimen.QMStartTypeZhi
		p.Apply(p.year, p.month, p.day, p.hour, p.minute)
	})
	p.opStartSelf.SetOnSelect(func(c *ui.OptionBox) {
		p.qmParams.StartType = qimen.QMStartTypeSelf
		p.Apply(p.year, p.month, p.day, p.hour, p.minute)
	})
	p.opStartMao.Disabled = true
	p.opStartZhi.Disabled = true
	p.opStartSelf.Visible = false

	ui.MakeOptionBoxGroup(p.opHideGan0, p.opHideGan1)
	p.opHideGan0.Select()
	p.opHideGan1.Disabled = true

	ui.MakeOptionBoxGroup(p.opHourPan, p.opDayPan, p.opMonthPan, p.opYearPan)
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
	p.opDayPan.Disabled = true

	p.cbHostingType.SetChecked(true)
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
		p.Apply(year, month, day, hour, minute)
	})
	p.btnPreHour2.SetOnClick(func(b *ui.Button) {
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
	p.btnNextHour2.SetOnClick(func(b *ui.Button) {
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
		UIShowMsgBox("时间不对", "确定", "取消", func(b *ui.Button) {
		}, func(b *ui.Button) {})
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
	p.opHideGan0.Visible = p.qmParams.Type != qimen.QMTypeAmaze
	p.opHideGan1.Visible = p.qmParams.Type != qimen.QMTypeAmaze

	//fmt
	switch p.qmParams.YMDH {
	case qimen.QMGameHour:
		p.ShowHourGame(pan)
	case qimen.QMGameDay:
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
		p.textGong[i].Text = fmt.Sprintf("\n  %s  %s    %s\n\n"+
			"%s %s %s%s%s\n\n"+
			"%s    %s    %s\n\n"+
			"      %s%s",
			empty, g.God, horse,
			g.PathGan, g.HideGan, g.Star, hosting, g.GuestGan,
			g.PathZhi, g.Door, g.HostGan,
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

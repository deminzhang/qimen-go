package world

import (
	"errors"
	"fmt"
	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/6tail/lunar-go/SolarUtil"
	"image"
	"qimen/ebiten_ui"
	"qimen/qimen"
	"strconv"
	"time"
)

type UIQiMen struct {
	ebiten_ui.BaseUI
	inputSYear  *ebiten_ui.InputBox
	inputSMonth *ebiten_ui.InputBox
	inputSDay   *ebiten_ui.InputBox
	inputSHour  *ebiten_ui.InputBox
	inputSMin   *ebiten_ui.InputBox

	textLYear   *ebiten_ui.InputBox
	textLMonth  *ebiten_ui.InputBox
	textLDay    *ebiten_ui.InputBox
	textLHour   *ebiten_ui.InputBox
	textYearRB  *ebiten_ui.InputBox
	textMonthRB *ebiten_ui.InputBox
	textDayRB   *ebiten_ui.InputBox
	textHourRB  *ebiten_ui.InputBox
	textJu      *ebiten_ui.InputBox

	cbTypeRoll  *ebiten_ui.CheckBox
	cbTypeFly   *ebiten_ui.CheckBox
	cbTypeAmaze *ebiten_ui.CheckBox

	btnCalc      *ebiten_ui.Button
	btnPreHour2  *ebiten_ui.Button
	btnNextHour2 *ebiten_ui.Button

	textGong []*ebiten_ui.TextBox
	zhiPan   []*ebiten_ui.TextBox
}

var uiQiMen *UIQiMen

func UIShowQiMen(width, height int) {
	if uiQiMen == nil {
		uiQiMen = NewUIQiMen(width, height)
		ebiten_ui.ActiveUI(uiQiMen)
	}
}
func UIHideQiMen() {
	if uiQiMen != nil {
		ebiten_ui.CloseUI(uiQiMen)
		uiQiMen = nil
	}
}

func NewUIQiMen(width, height int) *UIQiMen {
	//cx, cy := width/2, height/2 //win center
	p := &UIQiMen{}
	px0, py0 := 32, 0
	h := 32
	p.inputSYear = ebiten_ui.NewInputBox(image.Rect(px0, py0, px0+64, py0+h))
	p.inputSMonth = ebiten_ui.NewInputBox(image.Rect(px0+72, py0, px0+72+64, py0+h))
	p.inputSDay = ebiten_ui.NewInputBox(image.Rect(px0+72*2, py0, px0+72*2+64, py0+h))
	p.inputSHour = ebiten_ui.NewInputBox(image.Rect(px0+72*3, py0, px0+72*3+64, py0+h))
	p.inputSMin = ebiten_ui.NewInputBox(image.Rect(px0+72*4, py0, px0+72*4+64, py0+h))
	p.btnCalc = ebiten_ui.NewButton(image.Rect(px0+72*5, py0, px0+72*5+64, py0+h), "排局")
	p.cbTypeRoll = ebiten_ui.NewCheckBox(px0+72*6, py0, qimen.QMType[0])
	p.cbTypeFly = ebiten_ui.NewCheckBox(px0+72*7, py0, qimen.QMType[1])
	p.cbTypeAmaze = ebiten_ui.NewCheckBox(px0+72*8, py0, qimen.QMType[2])

	py0 += 32
	p.textLYear = ebiten_ui.NewInputBox(image.Rect(px0, py0, px0+64, py0+h))
	p.textLMonth = ebiten_ui.NewInputBox(image.Rect(px0+72, py0, px0+72+64, py0+h))
	p.textLDay = ebiten_ui.NewInputBox(image.Rect(px0+72*2, py0, px0+72*2+64, py0+h))
	p.textLHour = ebiten_ui.NewInputBox(image.Rect(px0+72*3, py0, px0+72*3+64, py0+h))
	p.btnPreHour2 = ebiten_ui.NewButton(image.Rect(px0+72*4, py0, px0+72*4+64, py0+h), "上一局")
	p.btnNextHour2 = ebiten_ui.NewButton(image.Rect(px0+72*5, py0, px0+72*5+64, py0+h), "下一局")
	py0 += 32
	p.textYearRB = ebiten_ui.NewInputBox(image.Rect(px0, py0, px0+64, py0+h))
	p.textMonthRB = ebiten_ui.NewInputBox(image.Rect(px0+72, py0, px0+72+64, py0+h))
	p.textDayRB = ebiten_ui.NewInputBox(image.Rect(px0+72*2, py0, px0+72*2+64, py0+h))
	p.textHourRB = ebiten_ui.NewInputBox(image.Rect(px0+72*3, py0, px0+72*3+64, py0+h))
	p.textJu = ebiten_ui.NewInputBox(image.Rect(px0+72*4, py0, px0+72*9, py0+h))

	px4, py4 := 64, 96+64
	const gongWidth = 128
	gongOffset := [][]int{{0, 0},
		{1, 2}, {2, 0}, {0, 1},
		{0, 0}, {1, 1}, {2, 2},
		{2, 1}, {0, 2}, {1, 0},
	}
	p.textGong = append(p.textGong, nil)
	for i := 1; i <= 9; i++ {
		offX, offZ := gongOffset[i][0]*gongWidth, gongOffset[i][1]*gongWidth
		txtGong := ebiten_ui.NewTextBox(image.Rect(px4+offX, py4+offZ, px4+offX+gongWidth, py4+offZ+gongWidth))
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
		txtZhi := ebiten_ui.NewTextBox(image.Rect(minX, minY, maxX, maxY))
		txtZhi.DisableHScroll = true
		txtZhi.DisableVScroll = true
		p.zhiPan = append(p.zhiPan, txtZhi)
		p.AddChild(txtZhi)
		p.zhiPan[i].SetText(LunarUtil.ZHI[i])
	}

	p.AddChild(p.inputSYear)
	p.AddChild(p.inputSMonth)
	p.AddChild(p.inputSDay)
	p.AddChild(p.inputSHour)
	p.AddChild(p.inputSMin)
	p.AddChild(p.btnCalc)
	p.AddChild(p.cbTypeRoll)
	p.AddChild(p.cbTypeFly)
	p.AddChild(p.cbTypeAmaze)

	p.AddChild(p.textLYear)
	p.AddChild(p.textLMonth)
	p.AddChild(p.textLDay)
	p.AddChild(p.textLHour)
	p.AddChild(p.btnPreHour2)
	p.AddChild(p.btnNextHour2)

	p.AddChild(p.textYearRB)
	p.AddChild(p.textMonthRB)
	p.AddChild(p.textDayRB)
	p.AddChild(p.textHourRB)
	p.AddChild(p.textJu)

	p.cbTypeRoll.SetChecked(false)
	p.cbTypeFly.SetChecked(false)
	p.cbTypeAmaze.SetChecked(true)
	p.cbTypeRoll.Disabled = true
	p.cbTypeFly.Disabled = true
	p.cbTypeAmaze.Disabled = true

	p.inputSYear.MaxChars = 4
	p.inputSMonth.MaxChars = 2
	p.inputSDay.MaxChars = 2
	p.inputSHour.MaxChars = 2
	p.inputSMin.MaxChars = 2

	p.inputSYear.DefaultText = "year"
	p.inputSMonth.DefaultText = "month"
	p.inputSDay.DefaultText = "day"
	p.inputSHour.DefaultText = "hour"
	p.inputSMin.DefaultText = "minute"

	p.btnCalc.SetOnClick(func(b *ebiten_ui.Button) {
		year, _ := strconv.Atoi(p.inputSYear.Text)
		month, _ := strconv.Atoi(p.inputSMonth.Text)
		day, _ := strconv.Atoi(p.inputSDay.Text)
		hour, _ := strconv.Atoi(p.inputSHour.Text)
		minute, _ := strconv.Atoi(p.inputSMin.Text)
		p.Apply(year, month, day, hour, minute)
	})
	p.btnPreHour2.SetOnClick(func(b *ebiten_ui.Button) {
		year, _ := strconv.Atoi(p.inputSYear.Text)
		month, _ := strconv.Atoi(p.inputSMonth.Text)
		day, _ := strconv.Atoi(p.inputSDay.Text)
		hour, _ := strconv.Atoi(p.inputSHour.Text)
		minute, _ := strconv.Atoi(p.inputSMin.Text)
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
			if 1582 == year && 10 == month && day == 14 {
				day = 4
			}
		}
		p.Apply(year, month, day, hour, minute)
	})
	p.btnNextHour2.SetOnClick(func(b *ebiten_ui.Button) {
		year, _ := strconv.Atoi(p.inputSYear.Text)
		month, _ := strconv.Atoi(p.inputSMonth.Text)
		day, _ := strconv.Atoi(p.inputSDay.Text)
		hour, _ := strconv.Atoi(p.inputSHour.Text)
		minute, _ := strconv.Atoi(p.inputSMin.Text)
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
			if 1582 == year && 10 == month && day == 5 {
				day = 15
			}
		}
		p.Apply(year, month, day, hour, minute)
	})

	p.textLYear.Editable = false
	p.textLMonth.Editable = false
	p.textLDay.Editable = false
	p.textLHour.Editable = false

	p.textYearRB.Editable = false
	p.textMonthRB.Editable = false
	p.textDayRB.Editable = false
	p.textHourRB.Editable = false
	p.textJu.Editable = false

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
	pan, err := qimen.NewPan(year, month, day, hour, minute)
	if err != nil {
		UIShowMsgBox("时间不对", "确定", "取消", func(b *ebiten_ui.Button) {
		}, func(b *ebiten_ui.Button) {})
	}
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

	p.textYearRB.SetText(pan.YearRB)
	p.textMonthRB.SetText(pan.MonthRB)
	p.textDayRB.SetText(pan.DayRB)
	p.textHourRB.SetText(pan.HourRB)

	p.textJu.Text = pan.JuText

	for i := 1; i <= 9; i++ {
		p.textGong[i].Text = pan.Gongs[i].FmtText
	}

	//大六壬 月将落时支 顺布余支
	yueJiangIdx, yueJiangPos := pan.YueJiangIdx, pan.YueJiangPos
	for i := yueJiangPos; i < yueJiangPos+12; i++ {
		z := LunarUtil.ZHI[qimen.Idx12[yueJiangIdx]]
		var j, h string
		if i == yueJiangPos {
			j = "\n月将"
		}
		if z == pan.HourHorse {
			h = "\n驿马"
		}
		p.zhiPan[qimen.Idx12[i]].SetText(fmt.Sprintf("%s%s%s", z, j, h))
		yueJiangIdx++
	}
}

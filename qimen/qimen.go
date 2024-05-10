package qimen

import (
	"errors"
	"fmt"
	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/6tail/lunar-go/SolarUtil"
	"github.com/6tail/lunar-go/calendar"
)

// QMPalace 奇门遁甲宫格
type QMPalace struct {
	Idx int //洛书宫数

	HostGan  string //地盘奇仪
	GuestGan string //天盘奇仪
	Star     string //天九星1~9
	Door     string //八门
	God      string //九神1~9
	PathGan  string //时辰流转干 鸣法作暗干 不逆三奇
	PathZhi  string //时辰流转支 鸣法暗支
	HideGan  string //暗干 非鸣法逆三奇
}

// QMGame 奇门遁甲盘局
type QMGame struct {
	Solar *calendar.Solar
	Lunar *calendar.Lunar

	SolarYear   int //1900-2100
	SolarMonth  int //1-12
	SolarDay    int //1-31
	SolarHour   int //0-23
	SolarMinute int //分

	lunarYear    int //农历年
	lunarMonth   int //农历月 1~12 闰-1~-12
	lunarDay     int //农历日 1~30
	lunarHour    int //农历时
	lunarQuarter int //农历刻
	LunarYearC   string
	LunarMonthC  string
	LunarDayC    string
	LunarHourC   string

	YearTB  string //年干支
	MonthTB string //月干支
	DayTB   string //日干支
	HourTB  string //时干支

	YueJian        string //月建
	YueJianZhiIdx  int    //月建地支号
	YueJiang       string //月将
	YueJiangZhiIdx int    //月将地支号

	HourGan    string //时干
	HourZhi    string //时支
	HourZhiIdx int    //时支序

	YMDH string //年月日时家

	JieQi string //节气文本

	YearPan  *QMPan //年家奇门盘
	MonthPan *QMPan //月家奇门盘
	DayPan   *QMPan //日家奇门盘
	HourPan  *QMPan //时家奇门盘
}

type QMPan struct {
	Type                int //盘式
	RotatingHostingType int //转盘.寄中法
	FlyType             int //飞盘.飞星法
	StartType           int //QMJuType 起局.起局法
	HideGanType         int //QMHideGanType

	Yuan3 int //三元1~3
	Ju    int //格局-1~-9,1~9, 年家为-1,-4,-7

	GanZhi string //干支 年家为年干支,,,时家为时干支
	Xun    string //干支旬首

	Duty        int    //值序
	DutyStar    string //值符
	DutyStarPos int    //值符落宫
	DutyDoor    string //值使
	DutyDoorPos int    //值使落宫
	RollHosting int    //转盘寄宫

	Horse string //驿马支位

	Gongs [10]QMPalace //九宫飞盘格
}

type QMParams struct {
	Type        int //QMType
	HostingType int //QMHostingType 转盘寄宫法
	FlyType     int //QMFlyType	飞盘九星飞宫
	JuType      int //QMJuType 起局方式
	HideGanType int //QMHideGanType 暗干起法
	SelfJu      int //自选局数

	YMDH int //年月日时家
}

func getQiMenYuan3Index(dayGanZhi string) int {
	jiaZiIndex := LunarUtil.GetJiaZiIndex(dayGanZhi)
	qiMenYuanIdx := jiaZiIndex % 15
	if qiMenYuanIdx < 5 {
		return 1
	} else if qiMenYuanIdx < 10 {
		return 2
	}
	return 3
}

func getQiMenJuIndex(jieQi string, yuan3Idx int) int {
	jqi := _JieQiIndex[jieQi]
	return _QiMenJu[jqi][yuan3Idx-1]
}

// GetTermTime 返回solar年的第n(1小寒)个节气进入时间 以1970-01-01 00:00:00 UTC为0,正后前负
func GetTermTime(year, n int) int64 {
	t := int64(31556925974.7*float64(year-1900)/1000) + int64(termData[n-1]*60-2208549300)
	return t
}

func checkDate(year, month, day, hour, minute int) error {
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

func (p *QMGame) calcGong(pp *QMPan) {
	g9 := &pp.Gongs
	xun := LunarUtil.GetXun(pp.GanZhi)
	gan := pp.GanZhi[:len(pp.GanZhi)/2]
	zhi := pp.GanZhi[len(pp.GanZhi)/2:]
	pp.Xun = xun
	pp.Horse = Horse[zhi]

	for i := 1; i <= 9; i++ {
		g9[i].Idx = i
	}
	ju := pp.Ju
	//地盘 三奇六仪
	if ju > 0 { //阳遁顺仪奇逆布
		for i := ju; i < ju+9; i++ {
			g9[Idx9[i]].HostGan = QM3Qi6Yi(i - ju + 1)
		}
	} else { //阴遁逆仪奇顺行
		for i := -ju + 9; i > -ju; i-- {
			g9[Idx9[i]].HostGan = QM3Qi6Yi(-ju + 9 - i + 1)
		}
	}

	//定值符值使 时旬首所遁地仪宫
	var duty int //值序符宫
	for i := 1; i <= 9; i++ {
		if g9[i].HostGan == HideJia[pp.Xun] {
			duty = i
			pp.Duty = i
			pp.DutyStar = QMStar9(i)
			pp.DutyDoor = QMDoor9(i)
			break
		}
	}

	//值符落宫
	//值符加于时干上，值使加之在时支。
	var dutyStarPos, dutyDoorPos, dutyGan36Idx int
	dutyGan := gan
	if dutyGan == "甲" {
		dutyGan = HideJia[pp.Xun] //遁甲
	}
	for i := 1; i <= 9; i++ {
		if g9[i].HostGan == dutyGan {
			dutyStarPos = i
			pp.DutyStarPos = i
			break
		}
	}
	for i, g := range []rune(T3Qi6Yi) {
		if string(g) == dutyGan {
			dutyGan36Idx = i
			break
		}
	}
	//找符使落宫  排暗干支神
	var jiaZiIdx int
	for i, x := range LunarUtil.JIA_ZI {
		if x == pp.Xun {
			jiaZiIdx = i
		}
	}
	if pp.Ju > 0 { //阳遁
		for i := duty; i <= duty+9; i++ {
			gid := Idx9[i]
			gz := LunarUtil.JIA_ZI[jiaZiIdx]
			jiaZiIdx++
			if jiaZiIdx > 60 {
				jiaZiIdx = 0
			}
			g, z := gz[:len(gz)/2], gz[len(gz)/2:]
			g9[gid].PathGan = g
			g9[gid].PathZhi = z
			if z == zhi {
				dutyDoorPos = gid
				pp.DutyDoorPos = gid
			}
		}
	} else {
		for i := duty + 9; i >= duty; i-- {
			gid := Idx9[i]
			gz := LunarUtil.JIA_ZI[jiaZiIdx]
			jiaZiIdx++
			if jiaZiIdx > 60 {
				jiaZiIdx = 0
			}
			g, z := gz[:len(gz)/2], gz[len(gz)/2:]
			g9[gid].PathGan = g
			g9[gid].PathZhi = z
			if z == zhi {
				dutyDoorPos = gid
				pp.DutyDoorPos = gid
			}
		}
	}

	//天盘 三奇六仪 值符带旬首仪飞
	var xunGanIdx int
	xunGan := pp.Xun[:len(pp.Xun)/2]
	if xunGan == "甲" {
		xunGan = HideJia[pp.Xun] //遁甲
	}
	for i, g := range []rune(T3Qi6Yi) {
		if string(g) == xunGan {
			xunGanIdx = i
			break
		}
	}

	switch pp.Type {
	case QMTypeRotating:
		//天盘 值符起落九星
		pp.RollHosting = 0
		dutyRoll := duty
		if duty == 5 { //值符天禽寄
			switch pp.RotatingHostingType {
			case QMHostingType2:
				dutyRoll = 2
			case QMHostingType28:
				if pp.Ju > 0 {
					dutyRoll = 8
				} else {
					dutyRoll = 2
				}
			}
			pp.RollHosting = dutyRoll
		}
		if dutyStarPos == 5 { //时干落中宫寄宫
			switch pp.RotatingHostingType {
			case QMHostingType2:
				dutyStarPos = 2
			case QMHostingType28:
				if pp.Ju > 0 {
					dutyStarPos = 8
				} else {
					dutyStarPos = 2
				}
			}
			pp.DutyStarPos = dutyStarPos
		}
		var starRollIdx = _QM2RollIdx[dutyStarPos] //转起宫
		var startIdx = _QM2RollIdx[dutyRoll]       //转起
		for i := starRollIdx; i < starRollIdx+8; i++ {
			gIdx := _QMRollIdx[Idx8[i]]
			g9[gIdx].Star = QMStar8(startIdx)
			//神盘
			if pp.Ju > 0 {
				g9[gIdx].God = QMGod8(1 + i - starRollIdx)
			} else {
				g9[gIdx].God = QMGod8(1 + starRollIdx + 8 - i)
			}
			startIdx++
		}
		//转八门
		if duty == 5 {
			pp.DutyDoor = QMDoor9(dutyRoll)
		}
		if dutyDoorPos == 5 {
			switch pp.RotatingHostingType {
			case QMHostingType2:
				dutyDoorPos = 2
			case QMHostingType28:
				if pp.Ju > 0 {
					dutyDoorPos = 8
				} else {
					dutyDoorPos = 2
				}
			}
			duty = dutyDoorPos
			pp.DutyDoorPos = dutyDoorPos
		}
		var doorRollIdx = _QM2RollIdx[dutyDoorPos] //转起宫
		startIdx = _QM2RollIdx[dutyRoll]           //转起
		for i := doorRollIdx; i < doorRollIdx+8; i++ {
			gIdx := _QMRollIdx[Idx8[i]]
			g9[gIdx].Door = QMDoor8(startIdx)
			startIdx++
		}
	case QMTypeFly, QMTypeAmaze:
		//天盘 值符起落九星
		if pp.Type == QMTypeAmaze || pp.Ju > 0 || pp.FlyType == QMFlyTypeAllOrder {
			for i := dutyStarPos; i < dutyStarPos+9; i++ {
				g9[Idx9[i]].Star = QMStar9(duty + i - dutyStarPos)
			}
		} else { //QMTypeFly && QMFlyTypeLunarReverse && p.Ju < 0
			for i := dutyStarPos + 9; i > dutyStarPos; i-- {
				g9[Idx9[i]].Star = QMStar9(duty + dutyStarPos + 9 - i)
			}
		}
		//神盘 值符起九神
		if pp.Ju > 0 { //阳遁
			for i := dutyStarPos; i < dutyStarPos+9; i++ {
				g9[Idx9[i]].God = QMGod9S(1 + i - dutyStarPos)
			}
		} else {
			for i := dutyStarPos + 9; i > dutyStarPos; i-- {
				g9[Idx9[i]].God = QMGod9L(1 + dutyStarPos + 9 - i)
			}
		}
		//飞布九门
		for i := dutyDoorPos; i < dutyDoorPos+9; i++ {
			g9[Idx9[i]].Door = QMDoor9(duty + i - dutyDoorPos)
		}
	}
	//排天盘 三奇六仪
	if pp.Ju > 0 {
		for i := dutyStarPos; i < dutyStarPos+9; i++ {
			g9[Idx9[i]].GuestGan = QM3Qi6Yi(xunGanIdx)
			xunGanIdx++
		}
	} else {
		for i := dutyStarPos + 9; i > dutyStarPos; i-- {
			g9[Idx9[i]].GuestGan = QM3Qi6Yi(xunGanIdx)
			xunGanIdx++
		}
	}

	if pp.Type == QMTypeAmaze {
		for i := 1; i <= +9; i++ {
			g9[i].HideGan = "  "
		}
	} else {
		for i := 1; i <= +9; i++ {
			g9[i].PathGan = "  "
			g9[i].PathZhi = "  "
			g9[i].HideGan = "  "
		}
		switch pp.HideGanType {
		case QMHideGanDutyDoorHour: //值使门起暗干
			hideGanStart := dutyDoorPos
			if dutyGan == g9[dutyDoorPos].HostGan { //时干同地盘干,暗干起中宫
				hideGanStart = 5
			}
			if pp.Ju > 0 { //阳遁
				for i := hideGanStart; i < hideGanStart+9; i++ {
					g9[Idx9[i]].HideGan = QM3Qi6Yi(dutyGan36Idx + i - hideGanStart)
				}
			} else {
				for i := hideGanStart + 9; i > hideGanStart; i-- {
					g9[Idx9[i]].HideGan = QM3Qi6Yi(dutyGan36Idx + hideGanStart + 9 - i)
				}
			}
		case QMHideGanDoorHomeGan: //门地盘起
			for i := 1; i < 9; i++ {
				if i != 5 {
					doorHomeGong := DoorHome[g9[i].Door]
					g9[i].HideGan = g9[doorHomeGong].HostGan
				}
			}
		}
	}
}

func NewPan(year, month, day, hour, minute int, params QMParams) (*QMGame, error) {
	ymdh, qmType, qmHostingType, pqmFlyType, startType, hideGanType :=
		params.YMDH, params.Type, params.HostingType, params.FlyType, params.JuType, params.HideGanType
	if err := checkDate(year, month, day, hour, minute); err != nil {
		return nil, err
	}
	solar := calendar.NewSolar(year, month, day, hour, minute, 0)
	lunar := calendar.NewLunarFromSolar(solar)
	c8 := lunar.GetBaZi()
	dayGanZhi := c8[2]
	hourGanZhi := c8[3]
	if hour == 23 { //晚子时日柱作次日
		di := LunarUtil.GetJiaZiIndex(dayGanZhi) + 1
		if di > 59 {
			di -= 60
		}
		dayGanZhi = LunarUtil.JIA_ZI[di]
	}
	jieQi := lunar.GetPrevJieQi()
	jieQiName := lunar.GetPrevJieQi().GetName()

	hourZhi := hourGanZhi[len(hourGanZhi)/2:]
	p := QMGame{
		Solar:       solar,
		Lunar:       lunar,
		SolarYear:   year,
		SolarMonth:  month,
		SolarDay:    day,
		SolarHour:   hour,
		SolarMinute: minute,
		lunarYear:   lunar.GetYear(),
		lunarMonth:  lunar.GetMonth(),
		lunarDay:    lunar.GetYear(),
		lunarHour:   lunar.GetHour(),
		LunarYearC:  lunar.GetYearInChinese(),
		LunarMonthC: lunar.GetMonthInChinese() + "月",
		LunarDayC:   lunar.GetDayInChinese(),
		LunarHourC:  hourZhi + "时",
		HourGan:     hourGanZhi[:len(hourGanZhi)/2],
		HourZhi:     hourZhi,
		YearTB:      c8[0],
		MonthTB:     c8[1],
		DayTB:       dayGanZhi,
		HourTB:      hourGanZhi,
		YueJian:     Jie2YueJian(lunar.GetPrevJie().GetName()),
		YueJiang:    Qi2YueJiang(lunar.GetPrevQi().GetName()),
		JieQi:       jieQiName,
	}
	switch ymdh {
	case QMGameHour: //排时家奇门
		var yuan, ju int
		yuan = getQiMenYuan3Index(c8[2])
		switch startType {
		case QMJuTypeSplit:
			//ju = getQiMenJuIndex(jieQiName, yuan)
			jqi := _JieQiIndex[jieQiName]
			ju = _QiMenJu[jqi][yuan-1]
		case QMJuTypeMaoShan:
			jieQiTime := jieQi.GetSolar()
			qiHour := jieQiTime.GetHour() //交气所在时辰起时
			if qiHour%2 == 0 {
				qiHour--
			}
			start := calendar.NewSolar(jieQiTime.GetYear(), jieQiTime.GetMonth(), jieQiTime.GetDay(), qiHour, 0, 0)
			minutes := solar.SubtractMinute(start)
			yuan = 1 + minutes/120/60 //60个时辰一元
			yuan = min(yuan, 3)       //三元完新节气不到用下元
			jqi := _JieQiIndex[jieQiName]
			ju = _QiMenJu[jqi][yuan-1]
		case QMJuTypeZhiRun:
			//TODO
		case QMJuTypeSelf:
			ju = params.SelfJu
		}
		p.HourPan = &QMPan{
			Yuan3:  yuan,
			Ju:     ju,
			GanZhi: hourGanZhi,

			Type:                qmType,
			RotatingHostingType: qmHostingType,
			FlyType:             pqmFlyType,
			StartType:           startType,
			HideGanType:         hideGanType,
		}
		//排九宫
		p.calcGong(p.HourPan)

		//排大六壬支 月将落时支 顺布余支
		for i := 1; i <= 12; i++ {
			if p.YueJian == LunarUtil.ZHI[i] {
				p.YueJianZhiIdx = i
				break
			}
		}
		for i := 1; i <= 12; i++ {
			if p.YueJiang == LunarUtil.ZHI[i] {
				p.YueJiangZhiIdx = i
				break
			}
		}
		for i := 1; i <= 12; i++ {
			if p.HourZhi == LunarUtil.ZHI[i] {
				p.HourZhiIdx = i
				break
			}
		}
	case QMGameDay:
		//p.DayPan = &QMPan{
		//	//Yuan3:  ,
		//	//Ju:    ,
		//	GanZhi: c8[2],
		//
		//	Type:                qmType,
		//	RotatingHostingType: qmHostingType,
		//	FlyType:             pqmFlyType,
		//	JuType:           startType,
		//  HideGanType:         hideGanType,
		//}
	case QMGameDay2:

	case QMGameMonth:
		yuan, ju := GetMonthYuanJu(p.YearTB)
		p.MonthPan = &QMPan{
			Yuan3:  yuan,
			Ju:     ju,
			GanZhi: c8[1],

			Type:                qmType,
			RotatingHostingType: qmHostingType,
			FlyType:             pqmFlyType,
			StartType:           startType,
			HideGanType:         hideGanType,
		}
		p.calcGong(p.MonthPan)
	case QMGameYear: //排年家奇门
		yuan, ju := GetYearYuanJu(p.lunarYear)
		p.YearPan = &QMPan{
			Yuan3:  yuan,
			Ju:     ju,
			GanZhi: c8[0],

			Type:                qmType,
			RotatingHostingType: qmHostingType,
			FlyType:             pqmFlyType,
			StartType:           startType,
			HideGanType:         hideGanType,
		}
		p.calcGong(p.YearPan)
	}

	return &p, nil
}

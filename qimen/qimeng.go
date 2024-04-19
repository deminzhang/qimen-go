package qimen

import (
	"errors"
	"fmt"
	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/6tail/lunar-go/SolarUtil"
	"github.com/6tail/lunar-go/calendar"
)

// QMGong 奇门遁甲宫格
type QMGong struct {
	Idx int //洛书宫数

	EarthGan string //地盘奇仪
	SkyGan   string //天盘奇仪
	HStar    string //天九星1~9
	Door     string //八门
	God      string //九神1~9
	AnGan    string //暗干
	AnZhi    string //暗支

	FmtText string
}

type QMPan struct {
	//SolarTime time.Time

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

	//days       int //1900.1.1起总Days
	//Constellation28 int            //星宿1～28

	HourGan string //时干
	HourZhi string //时支

	YearRB  string //年干支
	MonthRB string //月干支
	DayRB   string //日干支
	HourRB  string //时干支

	Type int //盘式

	JieQiName string //节气文本
	Yuan3     int    //三元1~3
	Ju        int    //格局-1~-9,1~9

	ShiXun      string //时辰旬首
	Duty        int    //值序
	DutyStar    string //值符
	DutyStarPos int    //值符落宫
	DutyDoor    string //值使
	DutyDoorPos int    //值使落宫

	YueJiangZhi    string //月将支名
	YueJiangZhiIdx int    //月将地支号
	YueJiangPos    int    //月将落地支宫
	HourHorse      string //驿马

	JuText string     //局文本
	Gongs  [10]QMGong //九宫飞盘格

	//FlyArr  [10]QMGong //九宫飞盘格
	//RollArr [10]QMGong //九宫转盘格
	//RollFly [10]QMGong //半飞半转盘
	//DayArr   [10]QMGong //日家奇门盘
	//MonthArr [10]QMGong //月家奇门盘
	//YearArr  [10]QMGong //年家奇门盘

}

func (p *QMPan) calcGong() {
	g9 := &p.Gongs
	for i := 1; i <= 9; i++ {
		g9[i].Idx = i
	}
	//地盘 三奇六仪
	if p.Ju > 0 { //阳遁顺仪奇逆布
		ju := p.Ju
		for i := ju; i < ju+9; i++ {
			g9[Idx9[i]].EarthGan = _QM3Q6YCircle[i-ju]
		}
	} else { //阴遁逆仪奇顺行
		ju := -p.Ju
		for i := ju + 9; i > ju; i-- {
			g9[Idx9[i]].EarthGan = _QM3Q6YCircle[ju+9-i]
		}
	}
	//定值符值使 时旬首所遁地仪宫
	var duty int //值序符宫
	for i := 1; i <= 9; i++ {
		if g9[i].EarthGan == _HideJia[p.ShiXun] {
			duty = i
			p.Duty = i
			p.DutyStar = _QMStar9[i]
			p.DutyDoor = _QMDoor9[i] // TODO if 转盘值使寄坤宫
			break
		}
	}
	//值符落宫 值符加之在时干
	dutyGan := p.HourGan
	if dutyGan == "甲" {
		dutyGan = _HideJia[p.ShiXun] //遁甲
	}
	var dutyStarPos int
	for i := 1; i <= 9; i++ {
		if g9[i].EarthGan == dutyGan {
			dutyStarPos = i
			p.DutyStarPos = i
			break
		}
	}

	//天盘 三奇六仪 值符宫起旬首仪
	var xunGanIdx int
	xunGan := p.ShiXun[:len(p.ShiXun)/2]
	if xunGan == "甲" {
		xunGan = _HideJia[p.ShiXun] //遁甲
	}
	for i, g := range _QM3Q6Y {
		if g == xunGan {
			xunGanIdx = i
		}
	}
	if p.Ju > 0 {
		for i := dutyStarPos; i < dutyStarPos+9; i++ {
			g9[Idx9[i]].SkyGan = _QM3Q6YCircle[xunGanIdx]
			xunGanIdx++
		}
	} else {
		for i := dutyStarPos + 9; i > dutyStarPos; i-- {
			g9[Idx9[i]].SkyGan = _QM3Q6YCircle[xunGanIdx]
			xunGanIdx++
		}
	}
	//值符位起
	//天盘 落九星 顺飞九宫
	//if p.Type == QMTypeAmaze || p.FlyType == QMFlyTypeAllOrder {
	for i := dutyStarPos; i < dutyStarPos+9; i++ {
		g9[Idx9[i]].HStar = _QMStar9Circle[duty+i-dutyStarPos]
	}
	//} else {
	//}

	//神盘 落九神
	//九神顺逆随遁转，八门九星顺宫飞
	//九神值符腾蛇是，太阴六合勾陈次。
	//太常朱雀九地天，午后白虎玄武治。//阴遁用白虎玄武
	if p.Ju > 0 { //阳遁
		for i := dutyStarPos; i < dutyStarPos+9; i++ {
			g9[Idx9[i]].God = _QMGod9SCircle[i-dutyStarPos]
		}
	} else {
		for i := dutyStarPos + 9; i > dutyStarPos; i-- {
			g9[Idx9[i]].God = _QMGod9LCircle[dutyStarPos+9-i]
		}
	}
	//排暗干支神
	//找符使落宫
	switch p.Type {
	case QMTypeRollDoor:
		//TODO
	case QMTypeFlyDoor: //值使起暗干
		//TODO
	case QMTypeAmaze: //暗支起值使 暗干暗支不逆三奇 鸣法：值符加于时干上，值使加之在时支。
		var jiaZiIdx int
		for i, x := range LunarUtil.JIA_ZI {
			if x == p.ShiXun {
				jiaZiIdx = i
			}
		}
		if p.Ju > 0 { //阳遁
			for i := duty; i <= duty+9; i++ {
				gid := Idx9[i]
				gz := LunarUtil.JIA_ZI[jiaZiIdx]
				jiaZiIdx++
				if jiaZiIdx > 60 {
					jiaZiIdx = 0
				}
				g, z := gz[:len(gz)/2], gz[len(gz)/2:]
				g9[gid].AnGan = g
				g9[gid].AnZhi = z
				if z == p.HourZhi {
					p.DutyDoorPos = gid
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
				g9[gid].AnGan = g
				g9[gid].AnZhi = z
				if z == p.HourZhi {
					p.DutyDoorPos = gid
				}
			}
		}
	}
	//布九门
	if p.Type == QMTypeAmaze || p.Type == QMTypeFlyDoor {
		for i := p.DutyDoorPos; i < p.DutyDoorPos+9; i++ {
			g9[Idx9[i]].Door = _QMDoor9[Idx9[duty+i-p.DutyDoorPos]]
		}
	} else {

	}

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

// 返回solar年的第n(1小寒)个节气进入时间
func getTermTime(year, n int) int64 {
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

func NewPan(year int, month int, day int, hour int, minute int) (*QMPan, error) {
	if err := checkDate(year, month, day, hour, minute); err != nil {
		return nil, err
	}
	cal := calendar.NewLunarFromSolar(calendar.NewSolar(year, month, day, hour, minute, 0))
	c8 := cal.GetBaZi()
	dayGanZhi := c8[2]
	shiGanZhi := c8[3]
	if hour == 23 { //晚子时日柱作次日
		di := LunarUtil.GetJiaZiIndex(dayGanZhi) + 1
		if di > 59 {
			di -= 60
		}
		dayGanZhi = LunarUtil.JIA_ZI[di]
	}
	dayYuanIdx := getQiMenYuan3Index(c8[2])
	yuanName := _Yuan3Name[dayYuanIdx]
	jieQiName := cal.GetPrevJieQi().GetName()
	juIdx := getQiMenJuIndex(jieQiName, dayYuanIdx)
	var juName string
	if juIdx < 0 {
		juName = fmt.Sprintf("阴%d局", juIdx)
	} else {
		juName = fmt.Sprintf("阳%d局", juIdx)
	}
	shiXun := LunarUtil.GetXun(shiGanZhi)
	shiZhi := shiGanZhi[len(shiGanZhi)/2:]
	p := QMPan{
		Type:        QMTypeAmaze,
		SolarYear:   year,
		SolarMonth:  month,
		SolarDay:    day,
		SolarHour:   hour,
		SolarMinute: minute,
		lunarYear:   cal.GetYear(),
		lunarMonth:  cal.GetMonth(),
		lunarDay:    cal.GetYear(),
		lunarHour:   cal.GetHour(),
		LunarYearC:  cal.GetYearInChinese(),
		LunarMonthC: cal.GetMonthInChinese() + "月",
		LunarDayC:   cal.GetDayInChinese(),
		LunarHourC:  shiZhi + "时",
		HourGan:     shiGanZhi[:len(shiGanZhi)/2],
		HourZhi:     shiZhi,
		YearRB:      c8[0],
		MonthRB:     c8[1],
		DayRB:       dayGanZhi,
		HourRB:      shiGanZhi,
		JieQiName:   jieQiName,
		Ju:          juIdx,
		Yuan3:       dayYuanIdx,
		ShiXun:      shiXun,
		YueJiangZhi: YueJiang[cal.GetMonth()],
		HourHorse:   Horse[shiZhi],
	}
	p.calcGong()
	//大六壬 月将落时支 顺布余支
	for i := 1; i <= 12; i++ {
		if p.YueJiangZhi == LunarUtil.ZHI[i] {
			p.YueJiangZhiIdx = i
			break
		}
	}
	for i := 1; i <= 12; i++ {
		if p.HourZhi == LunarUtil.ZHI[i] {
			p.YueJiangPos = i
			break
		}
	}

	//fmt
	p.JuText = fmt.Sprintf("%s%s %s %s遁%s 值符%s 值使%s", jieQiName, yuanName, juName,
		shiXun, _HideJia[shiXun], p.DutyStar, p.DutyDoor)

	for i := 1; i <= 9; i++ {
		g := p.Gongs[i]
		p.Gongs[i].FmtText = fmt.Sprintf("\n      %s\n\n%s    %s    %s\n\n%s    %s    %s\n\n      %s%s",
			g.God,
			g.AnGan, g.HStar, g.SkyGan,
			g.AnZhi, g.Door, g.EarthGan, Gua8In9[i],
			LunarUtil.NUMBER[i])
	}

	return &p, nil
}

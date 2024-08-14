package world

import (
	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/6tail/lunar-go/calendar"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"qimen/qimen"
	"strings"
)

type EightCharPan struct {
	X, Y   float32
	FYear  CharBody //流年通用
	FMonth CharBody //流月通用
	FDay   CharBody //流日通用
	FTime  CharBody //流时通用

	Player Player //玩家
}

func NewEightCharPan(centerX, centerY float32) *EightCharPan {
	return &EightCharPan{
		X: centerX, Y: centerY,
	}
}

func (g *EightCharPan) Update() error {
	cal := ThisGame.qmPan.Lunar
	//g.FYear = CharBody{Gan: cal.GetYearGan(), Zhi: cal.GetYearZhi()}
	//g.FMonth = CharBody{Gan: cal.GetMonthGan(), Zhi: cal.GetMonthZhi()}
	//g.FDay = CharBody{Gan: cal.GetDayGan(), Zhi: cal.GetDayZhi()}
	//g.FTime = CharBody{Gan: cal.GetTimeGan(), Zhi: cal.GetTimeZhi()}
	if g.Player.BirthTime == nil {
		g.Player.Reset(cal, GenderMale)
	}
	g.Player.UpdateHP()
	return nil
}

func (g *EightCharPan) Draw(screen *ebiten.Image) {
	ft14, _ := GetFontFace(14)
	ft28, _ := GetFontFace(28)
	cx, cy := g.X, g.Y
	sx, sy := cx-212, cy-148
	p := g.Player
	bz := p.BirthTime.GetEightChar()
	soul := bz.GetDayGan()

	vector.StrokeRect(screen, sx, sy-64, 390, 320, 1, colorWhite, true)
	sx += 16
	text.Draw(screen, "十神", ft14, int(sx), int(sy-32), colorWhite)
	text.Draw(screen, "天干", ft14, int(sx), int(sy-8), colorWhite)
	text.Draw(screen, "地支", ft14, int(sx), int(sy+32-8), colorWhite)
	text.Draw(screen, "本气", ft14, int(sx), int(sy+48), colorWhite)
	text.Draw(screen, "中气", ft14, int(sx), int(sy+64), colorWhite)
	text.Draw(screen, "余气", ft14, int(sx), int(sy+80), colorWhite)
	text.Draw(screen, "纳音", ft14, int(sx), int(sy+96), colorWhite)
	text.Draw(screen, "地势", ft14, int(sx), int(sy+112), colorWhite)
	text.Draw(screen, "自坐", ft14, int(sx), int(sy+128), colorWhite)
	text.Draw(screen, "空亡", ft14, int(sx), int(sy+144), colorWhite)
	//text.Draw(screen, "神煞", ft14, int(sx), int(sy+160), colorWhite)
	text.Draw(screen, "小运", ft14, int(sx), int(sy+160), colorWhite)
	text.Draw(screen, "大运", ft14, int(sx), int(sy+160+16), colorWhite)
	//text.Draw(screen, "流年", ft14, int(sx), int(sy+160), colorWhite)
	//text.Draw(screen, "流月", ft14, int(sx), int(sy+160), colorWhite)

	sx += 48
	text.Draw(screen, bz.GetYearShiShenGan(), ft14, int(sx), int(sy-32), colorWhite)
	text.Draw(screen, p.Year.Gan, ft28, int(sx), int(sy), ColorGanZhi(p.Year.Gan))
	text.Draw(screen, p.Year.Zhi, ft28, int(sx), int(sy+32), ColorGanZhi(p.Year.Zhi))
	text.Draw(screen, ShiShenShort(soul, p.Year.Body), ft14, int(sx), int(sy+48), ColorGanZhi(p.Year.Body))
	text.Draw(screen, ShiShenShort(soul, p.Year.Legs), ft14, int(sx), int(sy+64), ColorGanZhi(p.Year.Legs))
	text.Draw(screen, ShiShenShort(soul, p.Year.Feet), ft14, int(sx), int(sy+80), ColorGanZhi(p.Year.Feet))
	text.Draw(screen, LunarUtil.NAYIN[bz.GetYear()], ft14, int(sx), int(sy+96), colorWhite)
	text.Draw(screen, qimen.ChangSheng12[soul][p.Year.Zhi], ft14, int(sx), int(sy+112), colorWhite)
	text.Draw(screen, qimen.ChangSheng12[p.Year.Gan][p.Year.Zhi], ft14, int(sx), int(sy+128), colorWhite)
	text.Draw(screen, bz.GetYearXunKong(), ft14, int(sx), int(sy+144), colorGray)

	text.Draw(screen, strings.Join(p.Fates0, " "), ft14, int(sx), int(sy+160), colorWhite)
	text.Draw(screen, strings.Join(p.Fates, " "), ft14, int(sx), int(sy+160+16), colorWhite)
	sx += 48
	text.Draw(screen, bz.GetMonthShiShenGan(), ft14, int(sx), int(sy-32), colorWhite)
	text.Draw(screen, p.Month.Gan, ft28, int(sx), int(sy), ColorGanZhi(p.Month.Gan))
	text.Draw(screen, p.Month.Zhi, ft28, int(sx), int(sy+32), ColorGanZhi(p.Month.Zhi))
	text.Draw(screen, ShiShenShort(soul, p.Month.Body), ft14, int(sx), int(sy+48), ColorGanZhi(p.Month.Body))
	text.Draw(screen, ShiShenShort(soul, p.Month.Body), ft14, int(sx), int(sy+64), ColorGanZhi(p.Month.Legs))
	text.Draw(screen, ShiShenShort(soul, p.Month.Body), ft14, int(sx), int(sy+80), ColorGanZhi(p.Month.Feet))
	text.Draw(screen, LunarUtil.NAYIN[bz.GetMonth()], ft14, int(sx), int(sy+96), colorWhite)
	text.Draw(screen, qimen.ChangSheng12[soul][p.Month.Zhi], ft14, int(sx), int(sy+112), colorWhite)
	text.Draw(screen, qimen.ChangSheng12[p.Month.Gan][p.Month.Zhi], ft14, int(sx), int(sy+128), colorWhite)
	text.Draw(screen, bz.GetMonthXunKong(), ft14, int(sx), int(sy+144), colorGray)
	sx += 48
	text.Draw(screen, "元"+GenderName[p.Gender], ft14, int(sx), int(sy-32), colorWhite)
	text.Draw(screen, p.Day.Gan, ft28, int(sx), int(sy), ColorGanZhi(p.Day.Gan))
	text.Draw(screen, p.Day.Zhi, ft28, int(sx), int(sy+32), ColorGanZhi(p.Day.Zhi))
	text.Draw(screen, ShiShenShort(soul, p.Day.Body), ft14, int(sx), int(sy+48), ColorGanZhi(p.Day.Body))
	text.Draw(screen, ShiShenShort(soul, p.Day.Body), ft14, int(sx), int(sy+64), ColorGanZhi(p.Day.Legs))
	text.Draw(screen, ShiShenShort(soul, p.Day.Body), ft14, int(sx), int(sy+80), ColorGanZhi(p.Day.Feet))
	text.Draw(screen, LunarUtil.NAYIN[bz.GetDay()], ft14, int(sx), int(sy+96), colorWhite)
	text.Draw(screen, qimen.ChangSheng12[soul][p.Day.Zhi], ft14, int(sx), int(sy+112), colorWhite)
	text.Draw(screen, qimen.ChangSheng12[p.Day.Gan][p.Day.Zhi], ft14, int(sx), int(sy+128), colorWhite)
	text.Draw(screen, bz.GetDayXunKong(), ft14, int(sx), int(sy+144), colorGray)
	sx += 48
	text.Draw(screen, bz.GetTimeShiShenGan(), ft14, int(sx), int(sy-32), colorWhite)
	text.Draw(screen, p.Time.Gan, ft28, int(sx), int(sy), ColorGanZhi(p.Time.Gan))
	text.Draw(screen, p.Time.Zhi, ft28, int(sx), int(sy+32), ColorGanZhi(p.Time.Zhi))
	text.Draw(screen, ShiShenShort(soul, p.Time.Body), ft14, int(sx), int(sy+48), ColorGanZhi(p.Time.Body))
	text.Draw(screen, ShiShenShort(soul, p.Time.Body), ft14, int(sx), int(sy+64), ColorGanZhi(p.Time.Legs))
	text.Draw(screen, ShiShenShort(soul, p.Time.Body), ft14, int(sx), int(sy+80), ColorGanZhi(p.Time.Feet))
	text.Draw(screen, LunarUtil.NAYIN[bz.GetTime()], ft14, int(sx), int(sy+96), colorWhite)
	text.Draw(screen, qimen.ChangSheng12[soul][p.Time.Zhi], ft14, int(sx), int(sy+112), colorWhite)
	text.Draw(screen, qimen.ChangSheng12[p.Time.Gan][p.Time.Zhi], ft14, int(sx), int(sy+128), colorWhite)
	text.Draw(screen, bz.GetTimeXunKong(), ft14, int(sx), int(sy+144), colorGray)

	sx = 1000
	sy = 200
	vector.StrokeRect(screen, sx, sy, 160, 160, 1, colorWhite, true)
	vector.StrokeRect(screen, sx-40, sy+160, 240, 300, 1, colorWhite, true)
	vector.StrokeRect(screen, sx, sy+460, 160, 220, 1, colorWhite, true)
	vector.StrokeRect(screen, sx-10, sy+680, 180, 50, 1, colorWhite, true)
}

type CharBody struct {
	Gan  string //干为头
	Zhi  string //支为身
	Head string //干为头
	Body string //本气为身
	Legs string //中气为腿
	Feet string //余气为足

	HPHead int //干为头
	HPBody int //本气为身
	HPLegs int //中气为腿
	HPFeet int //余气为足
}

type Player struct {
	Gender    int //性别
	BirthTime *calendar.Lunar
	Year      CharBody  //年柱
	Month     CharBody  //月柱
	Day       CharBody  //日柱
	Time      CharBody  //时柱
	Fate      *CharBody //当前大运

	DaYun  []*calendar.DaYun //大运
	Fates0 []string          //小运名
	Fates  []string          //大运名

	ShenShaY []string //神煞
	ShenShaM []string //神煞
	ShenShaD []string //神煞
	ShenShaT []string //神煞
}

func (p *Player) Reset(lunar *calendar.Lunar, gender int) {
	p.BirthTime = lunar
	p.Gender = gender
	yearZhi := lunar.GetYearZhi()
	bz := lunar.GetEightChar()
	p.Year = CharBody{Gan: bz.GetYearGan(), Zhi: yearZhi,
		Head: bz.GetYearGan(), Body: GetHideGan(yearZhi, 0),
		Legs: GetHideGan(yearZhi, 1), Feet: GetHideGan(yearZhi, 2)}
	monthZhi := bz.GetMonthZhi()
	p.Month = CharBody{Gan: bz.GetMonthGan(), Zhi: monthZhi,
		Head: bz.GetMonthGan(), Body: GetHideGan(monthZhi, 0),
		Legs: GetHideGan(monthZhi, 1), Feet: GetHideGan(monthZhi, 2)}
	dayZhi := bz.GetDayZhi()
	p.Day = CharBody{Gan: bz.GetDayGan(), Zhi: dayZhi,
		Head: bz.GetDayGan(), Body: GetHideGan(dayZhi, 0),
		Legs: GetHideGan(dayZhi, 1), Feet: GetHideGan(dayZhi, 2)}
	timeZhi := bz.GetTimeZhi()
	p.Time = CharBody{Gan: bz.GetTimeGan(), Zhi: timeZhi,
		Head: bz.GetTimeGan(), Body: GetHideGan(timeZhi, 0),
		Legs: GetHideGan(timeZhi, 1), Feet: GetHideGan(timeZhi, 2)}

	yun := bz.GetYun(p.Gender)
	p.DaYun = yun.GetDaYun()
	p.Fates0 = nil
	p.Fates = nil
	for i, daYun := range p.DaYun {
		//fmt.Printf("大运[%d] = %d年 %d岁 %s\n", daYun.GetIndex(), daYun.GetStartYear(), daYun.GetStartAge(), daYun.GetGanZhi())
		if i == 0 {
			for _, xiaoYun := range daYun.GetXiaoYun() {
				p.Fates0 = append(p.Fates0, xiaoYun.GetGanZhi())
				//fmt.Printf(" 小运[%d] = %d年 %d岁 %s\n", xiaoYun.GetIndex(), xiaoYun.GetYear(), xiaoYun.GetAge(), xiaoYun.GetGanZhi())
			}
			//p.Fates = append(p.Fates, bz.GetMonth())
		} else {
			p.Fates = append(p.Fates, daYun.GetGanZhi())
		}
	}
}

func (p *Player) CalcShenSha() {

}

func (p *Player) UpdateHP() {
}

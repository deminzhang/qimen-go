package world

import (
	"github.com/6tail/lunar-go/calendar"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"qimen/ui"
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

func (g *EightCharPan) SetDate(cal *calendar.Lunar) error {
	g.Player.Reset(cal)
	return nil
}

func (g *EightCharPan) Update() error {
	cal := uiQiMen.pan.Lunar
	if g.Player.BirthTime.GetYear() != cal.GetYear() || g.Player.BirthTime.GetMonth() != cal.GetMonth() ||
		g.Player.BirthTime.GetDay() != cal.GetDay() || g.Player.BirthTime.GetHour() != cal.GetHour() {
		g.Player.Reset(cal)
	}
	return nil
}

func (g *EightCharPan) Draw(screen *ebiten.Image) {
	ft := ui.GetDefaultUIFont()
	cx, cy := g.X, g.Y
	text.Draw(screen, g.Player.Year.Gan, ft, int(cx-32), int(cy-16), colorGanZhi[g.Player.Year.Gan])
	text.Draw(screen, g.Player.Year.Zhi, ft, int(cx-32), int(cy), colorGanZhi[g.Player.Year.Zhi])
	text.Draw(screen, g.Player.Year.Body, ft, int(cx-32), int(cy+32), colorGanZhi[g.Player.Year.Body])
	text.Draw(screen, g.Player.Year.Legs, ft, int(cx-32), int(cy+48), colorGanZhi[g.Player.Year.Legs])
	text.Draw(screen, g.Player.Year.Feet, ft, int(cx-32), int(cy+64), colorGanZhi[g.Player.Year.Feet])

	text.Draw(screen, g.Player.Month.Gan, ft, int(cx-16), int(cy-16), colorGanZhi[g.Player.Month.Gan])
	text.Draw(screen, g.Player.Month.Zhi, ft, int(cx-16), int(cy), colorGanZhi[g.Player.Month.Zhi])
	text.Draw(screen, g.Player.Month.Body, ft, int(cx-16), int(cy+32), colorGanZhi[g.Player.Month.Body])
	text.Draw(screen, g.Player.Month.Legs, ft, int(cx-16), int(cy+48), colorGanZhi[g.Player.Month.Legs])
	text.Draw(screen, g.Player.Month.Feet, ft, int(cx-16), int(cy+64), colorGanZhi[g.Player.Month.Feet])

	text.Draw(screen, g.Player.Day.Gan, ft, int(cx), int(cy-16), colorGanZhi[g.Player.Day.Gan])
	text.Draw(screen, g.Player.Day.Zhi, ft, int(cx), int(cy), colorGanZhi[g.Player.Day.Zhi])
	text.Draw(screen, g.Player.Day.Body, ft, int(cx), int(cy+32), colorGanZhi[g.Player.Day.Body])
	text.Draw(screen, g.Player.Day.Legs, ft, int(cx), int(cy+48), colorGanZhi[g.Player.Day.Legs])
	text.Draw(screen, g.Player.Day.Feet, ft, int(cx), int(cy+64), colorGanZhi[g.Player.Day.Feet])

	text.Draw(screen, g.Player.Time.Gan, ft, int(cx+16), int(cy-16), colorGanZhi[g.Player.Time.Gan])
	text.Draw(screen, g.Player.Time.Zhi, ft, int(cx+16), int(cy), colorGanZhi[g.Player.Time.Zhi])
	text.Draw(screen, g.Player.Time.Body, ft, int(cx+16), int(cy+32), colorGanZhi[g.Player.Time.Body])
	text.Draw(screen, g.Player.Time.Legs, ft, int(cx+16), int(cy+48), colorGanZhi[g.Player.Time.Legs])
	text.Draw(screen, g.Player.Time.Feet, ft, int(cx+16), int(cy+64), colorGanZhi[g.Player.Time.Feet])

}

type CharBody struct {
	Gan  string //干为头
	Zhi  string //支为身
	Head string //干为头
	Body string //本气为身
	Legs string //中气为腿
	Feet string //余气为足
}

type Player struct {
	Gender    string //性别
	BirthTime calendar.Lunar
	Year      CharBody //年柱
	Month     CharBody //月柱
	Day       CharBody //日柱
	Time      CharBody //时柱

	Fate CharBody //个人大运
}

func (b *Player) Reset(lunar *calendar.Lunar) {
	b.BirthTime = *lunar
	b.Gender = "乾" //乾男 坤女
	yearZhi := lunar.GetYearZhi()
	b.Year = CharBody{Gan: lunar.GetYearGan(), Zhi: yearZhi,
		Head: lunar.GetYearGan(), Body: GetHideGan(yearZhi, 0),
		Legs: GetHideGan(yearZhi, 1), Feet: GetHideGan(yearZhi, 2)}
	monthZhi := lunar.GetMonthZhi()
	b.Month = CharBody{Gan: lunar.GetMonthGan(), Zhi: monthZhi,
		Head: lunar.GetMonthGan(), Body: GetHideGan(monthZhi, 0),
		Legs: GetHideGan(monthZhi, 1), Feet: GetHideGan(monthZhi, 2)}
	dayZhi := lunar.GetDayZhi()
	b.Day = CharBody{Gan: lunar.GetDayGan(), Zhi: dayZhi,
		Head: lunar.GetDayGan(), Body: GetHideGan(dayZhi, 0),
		Legs: GetHideGan(dayZhi, 1), Feet: GetHideGan(dayZhi, 2)}
	timeZhi := lunar.GetTimeZhi()
	b.Time = CharBody{Gan: lunar.GetTimeGan(), Zhi: timeZhi,
		Head: lunar.GetTimeGan(), Body: GetHideGan(timeZhi, 0),
		Legs: GetHideGan(timeZhi, 1), Feet: GetHideGan(timeZhi, 2)}

	b.Fate = CharBody{}
}

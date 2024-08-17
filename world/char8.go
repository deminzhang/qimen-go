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

// 八字能量比例
// 年干8 月干12 日元12 时干12
// 年支4 月支40 日支12 时支12
const (
	HpGY = 80
	HpZY = 40
	HpGM = 120
	HpZM = 400
	HpGD = 120
	HpZD = 120
	HpGT = 120
	HpZT = 120
)

// HideGanVal 藏干值比例
var HideGanVal = map[int][]int{
	1: {10},
	2: {7, 3},
	3: {6, 3, 1},
}

type EightCharPan struct {
	X, Y float32
	//YMDH   string
	FYear  CharBody //流年通用
	FMonth CharBody //流月通用
	FDay   CharBody //流日通用
	FTime  CharBody //流时通用
	Player Player   //玩家
}

func NewEightCharPan(centerX, centerY float32) *EightCharPan {
	return &EightCharPan{
		X: centerX, Y: centerY,
	}
}

func (g *EightCharPan) Update() error {
	cal := ThisGame.qmGame.Lunar
	g.FYear = NewCharBody(cal.GetYearGan(), cal.GetYearZhi(), HpGY, HpZY)
	g.FMonth = NewCharBody(cal.GetMonthGan(), cal.GetMonthZhi(), HpGM, HpZM)
	g.FDay = NewCharBody(cal.GetDayGan(), cal.GetDayZhi(), HpGD, HpZD)
	g.FTime = NewCharBody(cal.GetTimeGan(), cal.GetTimeZhi(), HpGT, HpZT)
	if g.Player.BirthTime == nil {
		g.Player.Reset(cal, GenderMale)
	}
	g.Player.UpdateHp()
	return nil
}

func (g *EightCharPan) Draw(dst *ebiten.Image) {
	ft12, _ := GetFontFace(12)
	ft14, _ := GetFontFace(14)
	ft28, _ := GetFontFace(28)
	cx, cy := g.X, g.Y
	p := g.Player
	bz := p.BirthTime.GetEightChar()
	soul := bz.GetDayGan()
	//八字总览
	{
		sx, sy := cx-248, cy-172
		vector.StrokeRect(dst, sx, sy-64, 312, 360, 1, colorWhite, true)
		sx += 16
		text.Draw(dst, "十神", ft14, int(sx), int(sy-32), colorWhite)
		text.Draw(dst, "天干", ft14, int(sx), int(sy-8), colorWhite)
		text.Draw(dst, "地支", ft14, int(sx), int(sy+32-8), colorWhite)
		text.Draw(dst, "本气", ft14, int(sx), int(sy+48), colorWhite)
		text.Draw(dst, "中气", ft14, int(sx), int(sy+64), colorWhite)
		text.Draw(dst, "余气", ft14, int(sx), int(sy+80), colorWhite)
		text.Draw(dst, "纳音", ft14, int(sx), int(sy+96), colorWhite)
		text.Draw(dst, "地势", ft14, int(sx), int(sy+112), colorWhite) //地势/长生/星运
		text.Draw(dst, "自坐", ft14, int(sx), int(sy+128), colorWhite)
		text.Draw(dst, "空亡", ft14, int(sx), int(sy+144), colorWhite)
		text.Draw(dst, "小运", ft14, int(sx), int(sy+160), colorWhite)
		text.Draw(dst, "大运", ft14, int(sx), int(sy+160+16), colorWhite)
		text.Draw(dst, "神煞", ft14, int(sx), int(sy+160+32), colorWhite)
		//text.Draw(dst, "流年", ft14, int(sx), int(sy+160), colorWhite)
		//text.Draw(dst, "流月", ft14, int(sx), int(sy+160), colorWhite)
		sx += 48
		text.Draw(dst, bz.GetYearShiShenGan(), ft14, int(sx), int(sy-32), colorWhite)
		text.Draw(dst, p.Year.Gan, ft28, int(sx), int(sy), ColorGanZhi(p.Year.Gan))
		text.Draw(dst, p.Year.Zhi, ft28, int(sx), int(sy+32), ColorGanZhi(p.Year.Zhi))
		text.Draw(dst, p.Year.Body, ft14, int(sx), int(sy+48), ColorGanZhi(p.Year.Body))
		text.Draw(dst, ShiShenShort(soul, p.Year.Body), ft14, int(sx+16), int(sy+48), colorWhite)
		text.Draw(dst, p.Year.Legs, ft14, int(sx), int(sy+64), ColorGanZhi(p.Year.Legs))
		text.Draw(dst, ShiShenShort(soul, p.Year.Legs), ft14, int(sx+16), int(sy+64), colorWhite)
		text.Draw(dst, p.Year.Feet, ft14, int(sx), int(sy+80), ColorGanZhi(p.Year.Feet))
		text.Draw(dst, ShiShenShort(soul, p.Year.Feet), ft14, int(sx+16), int(sy+80), colorWhite)
		text.Draw(dst, LunarUtil.NAYIN[bz.GetYear()], ft14, int(sx), int(sy+96), ColorNaYin(bz.GetYear()))
		text.Draw(dst, qimen.ChangSheng12[soul][p.Year.Zhi], ft14, int(sx), int(sy+112), ColorGanZhi(soul))
		text.Draw(dst, qimen.ChangSheng12[p.Year.Gan][p.Year.Zhi], ft14, int(sx), int(sy+128), ColorGanZhi(p.Year.Gan))
		text.Draw(dst, bz.GetYearXunKong(), ft14, int(sx), int(sy+144), colorGray)
		text.Draw(dst, strings.Join(p.Fates0, " "), ft14, int(sx), int(sy+160), colorWhite)
		text.Draw(dst, strings.Join(p.Fates, " "), ft14, int(sx), int(sy+160+16), colorWhite)
		text.Draw(dst, strings.Join(p.ShenShaY, "\n"), ft12, int(sx), int(sy+160+32), colorWhite)
		sx += 48
		text.Draw(dst, bz.GetMonthShiShenGan(), ft14, int(sx), int(sy-32), colorWhite)
		text.Draw(dst, p.Month.Gan, ft28, int(sx), int(sy), ColorGanZhi(p.Month.Gan))
		text.Draw(dst, p.Month.Zhi, ft28, int(sx), int(sy+32), ColorGanZhi(p.Month.Zhi))
		text.Draw(dst, p.Month.Body, ft14, int(sx), int(sy+48), ColorGanZhi(p.Month.Body))
		text.Draw(dst, ShiShenShort(soul, p.Month.Body), ft14, int(sx+16), int(sy+48), colorWhite)
		text.Draw(dst, p.Month.Legs, ft14, int(sx), int(sy+64), ColorGanZhi(p.Month.Legs))
		text.Draw(dst, ShiShenShort(soul, p.Month.Legs), ft14, int(sx+16), int(sy+64), colorWhite)
		text.Draw(dst, p.Month.Feet, ft14, int(sx), int(sy+80), ColorGanZhi(p.Month.Feet))
		text.Draw(dst, ShiShenShort(soul, p.Month.Feet), ft14, int(sx+16), int(sy+80), colorWhite)
		text.Draw(dst, LunarUtil.NAYIN[bz.GetMonth()], ft14, int(sx), int(sy+96), ColorNaYin(bz.GetMonth()))
		text.Draw(dst, qimen.ChangSheng12[soul][p.Month.Zhi], ft14, int(sx), int(sy+112), ColorGanZhi(soul))
		text.Draw(dst, qimen.ChangSheng12[p.Month.Gan][p.Month.Zhi], ft14, int(sx), int(sy+128), ColorGanZhi(p.Month.Gan))
		text.Draw(dst, bz.GetMonthXunKong(), ft14, int(sx), int(sy+144), colorGray)
		text.Draw(dst, strings.Join(p.ShenShaM, "\n"), ft12, int(sx), int(sy+160+32), colorWhite)
		sx += 48
		text.Draw(dst, "元"+GenderName[p.Gender], ft14, int(sx), int(sy-32), colorWhite)
		text.Draw(dst, p.Day.Gan, ft28, int(sx), int(sy), ColorGanZhi(p.Day.Gan))
		text.Draw(dst, p.Day.Zhi, ft28, int(sx), int(sy+32), ColorGanZhi(p.Day.Zhi))
		text.Draw(dst, p.Day.Body, ft14, int(sx), int(sy+48), ColorGanZhi(p.Day.Body))
		text.Draw(dst, ShiShenShort(soul, p.Day.Body), ft14, int(sx+16), int(sy+48), colorWhite)
		text.Draw(dst, p.Day.Legs, ft14, int(sx), int(sy+64), ColorGanZhi(p.Day.Legs))
		text.Draw(dst, ShiShenShort(soul, p.Day.Legs), ft14, int(sx+16), int(sy+64), colorWhite)
		text.Draw(dst, p.Day.Feet, ft14, int(sx), int(sy+80), ColorGanZhi(p.Day.Feet))
		text.Draw(dst, ShiShenShort(soul, p.Day.Feet), ft14, int(sx+16), int(sy+80), colorWhite)
		text.Draw(dst, LunarUtil.NAYIN[bz.GetDay()], ft14, int(sx), int(sy+96), ColorNaYin(bz.GetDay()))
		text.Draw(dst, qimen.ChangSheng12[soul][p.Day.Zhi], ft14, int(sx), int(sy+112), ColorGanZhi(soul))
		text.Draw(dst, qimen.ChangSheng12[p.Day.Gan][p.Day.Zhi], ft14, int(sx), int(sy+128), ColorGanZhi(p.Day.Gan))
		text.Draw(dst, bz.GetDayXunKong(), ft14, int(sx), int(sy+144), colorGray)
		text.Draw(dst, strings.Join(p.ShenShaD, "\n"), ft12, int(sx), int(sy+160+32), colorWhite)
		sx += 48
		text.Draw(dst, bz.GetTimeShiShenGan(), ft14, int(sx), int(sy-32), colorWhite)
		text.Draw(dst, p.Time.Gan, ft28, int(sx), int(sy), ColorGanZhi(p.Time.Gan))
		text.Draw(dst, p.Time.Zhi, ft28, int(sx), int(sy+32), ColorGanZhi(p.Time.Zhi))
		text.Draw(dst, p.Time.Body, ft14, int(sx), int(sy+48), ColorGanZhi(p.Time.Body))
		text.Draw(dst, ShiShenShort(soul, p.Time.Body), ft14, int(sx+16), int(sy+48), colorWhite)
		text.Draw(dst, p.Time.Legs, ft14, int(sx), int(sy+64), ColorGanZhi(p.Time.Legs))
		text.Draw(dst, ShiShenShort(soul, p.Time.Legs), ft14, int(sx+16), int(sy+64), colorWhite)
		text.Draw(dst, p.Time.Feet, ft14, int(sx), int(sy+80), ColorGanZhi(p.Time.Feet))
		text.Draw(dst, ShiShenShort(soul, p.Time.Feet), ft14, int(sx+16), int(sy+80), colorWhite)
		text.Draw(dst, LunarUtil.NAYIN[bz.GetTime()], ft14, int(sx), int(sy+96), ColorNaYin(bz.GetTime()))
		text.Draw(dst, qimen.ChangSheng12[soul][p.Time.Zhi], ft14, int(sx), int(sy+112), ColorGanZhi(soul))
		text.Draw(dst, qimen.ChangSheng12[p.Time.Gan][p.Time.Zhi], ft14, int(sx), int(sy+128), ColorGanZhi(p.Time.Gan))
		text.Draw(dst, bz.GetTimeXunKong(), ft14, int(sx), int(sy+144), colorGray)
		text.Draw(dst, strings.Join(p.ShenShaT, "\n"), ft12, int(sx), int(sy+160+32), colorWhite)
	}
	//竖象 年头颈/月胸腹/日腹股/时腿足 年额/月目/日鼻/时口 干左支右?
	{
		sx, sy := cx+120, cy-260
		mx := int(sx + 28)
		w := float32(74)
		vector.StrokeRect(dst, sx, sy, w, 64, 1, colorWhite, true)          //头
		vector.StrokeRect(dst, sx, sy+64, w, 32, 1, colorWhite, true)       //颈
		vector.StrokeRect(dst, sx, sy+96, w, 64, 1, colorWhite, true)       //胸
		vector.StrokeRect(dst, sx-36, sy+96, 32, 190, 1, colorWhite, true)  //胳膊手1
		vector.StrokeRect(dst, sx+w+4, sy+96, 32, 190, 1, colorWhite, true) //胳膊手2
		vector.StrokeRect(dst, sx, sy+96+64, w, 32, 1, colorWhite, true)    //腹
		vector.StrokeRect(dst, sx, sy+96*2, w, 64, 1, colorWhite, true)     //小腹
		vector.StrokeRect(dst, sx, sy+96*2+64, w, 32, 1, colorWhite, true)  //股
		vector.StrokeRect(dst, sx, sy+96*3, w, 64, 1, colorWhite, true)     //腿
		vector.StrokeRect(dst, sx, sy+96*3+64, w, 32, 1, colorWhite, true)  //足
		drawMiniGanZhi := func(sy int) {
			text.Draw(dst, p.Year.Gan, ft14, mx, int(sy-16), ColorGanZhi(p.Year.Gan))
			text.Draw(dst, p.Year.Zhi, ft14, mx+16, int(sy-16), ColorGanZhi(p.Year.Zhi))
			text.Draw(dst, p.Month.Gan, ft14, mx, int(sy), ColorGanZhi(p.Month.Gan))
			text.Draw(dst, p.Month.Zhi, ft14, mx+16, int(sy), ColorGanZhi(p.Month.Zhi))
			text.Draw(dst, p.Day.Gan, ft14, mx, int(sy+16), ColorGanZhi(p.Day.Gan))
			text.Draw(dst, p.Day.Zhi, ft14, mx+16, int(sy+16), ColorGanZhi(p.Day.Zhi))
			text.Draw(dst, p.Time.Gan, ft14, mx, int(sy+32), ColorGanZhi(p.Time.Gan))
			text.Draw(dst, p.Time.Zhi, ft14, mx+16, int(sy+32), ColorGanZhi(p.Time.Zhi))
		}
		sy += 28 //头
		text.Draw(dst, p.Year.Gan, ft28, int(sx), int(sy), ColorGanZhi(p.Year.Gan))
		{ //面内
			drawMiniGanZhi(int(sy))
		}
		sy += 64 //颈
		text.Draw(dst, p.Year.Zhi, ft28, int(sx), int(sy), ColorGanZhi(p.Year.Zhi))
		text.Draw(dst, p.Year.Body, ft14, mx, int(sy), ColorGanZhi(p.Year.Body))
		text.Draw(dst, p.Year.Legs, ft14, mx+16, int(sy), ColorGanZhi(p.Year.Legs))
		text.Draw(dst, p.Year.Feet, ft14, mx+32, int(sy), ColorGanZhi(p.Year.Feet))
		sy += 32 //胸
		text.Draw(dst, p.Month.Gan, ft28, int(sx), int(sy), ColorGanZhi(p.Month.Gan))
		{ //胸细节
			drawMiniGanZhi(int(sy))
		}
		sy += 64 //腹
		text.Draw(dst, p.Month.Zhi, ft28, int(sx), int(sy), ColorGanZhi(p.Month.Zhi))
		text.Draw(dst, p.Month.Body, ft14, mx, int(sy), ColorGanZhi(p.Month.Body))
		text.Draw(dst, p.Month.Legs, ft14, mx+16, int(sy), ColorGanZhi(p.Month.Legs))
		text.Draw(dst, p.Month.Feet, ft14, mx+32, int(sy), ColorGanZhi(p.Month.Feet))
		sy += 32 //腹
		text.Draw(dst, p.Day.Gan, ft28, int(sx), int(sy), ColorGanZhi(p.Day.Gan))
		{ //腹股细节
			drawMiniGanZhi(int(sy))
		}
		sy += 64 //股
		text.Draw(dst, p.Day.Zhi, ft28, int(sx), int(sy), ColorGanZhi(p.Day.Zhi))
		text.Draw(dst, p.Day.Body, ft14, int(sx+28), int(sy), ColorGanZhi(p.Day.Body))
		text.Draw(dst, p.Day.Legs, ft14, int(sx+28+16), int(sy), ColorGanZhi(p.Day.Legs))
		text.Draw(dst, p.Day.Feet, ft14, int(sx+28+32), int(sy), ColorGanZhi(p.Day.Feet))
		sy += 32 //腿
		text.Draw(dst, p.Time.Gan, ft28, int(sx), int(sy), ColorGanZhi(p.Time.Gan))
		{ //腿细节
			drawMiniGanZhi(int(sy))
		}
		sy += 64 //足
		text.Draw(dst, p.Time.Zhi, ft28, int(sx), int(sy), ColorGanZhi(p.Time.Zhi))
		text.Draw(dst, p.Time.Body, ft14, int(sx+28), int(sy), ColorGanZhi(p.Time.Body))
		text.Draw(dst, p.Time.Legs, ft14, int(sx+28+16), int(sy), ColorGanZhi(p.Time.Legs))
		text.Draw(dst, p.Time.Feet, ft14, int(sx+28+32), int(sy), ColorGanZhi(p.Time.Feet))
	}
	//横象 年祖月父母日夫妻时子孙 干动支静 干为军支为营 干为官支为民
	{
		sx, sy := cx-240, cy+180
		vector.StrokeRect(dst, sx, sy+2, 96, 96, 1, colorWhite, true)
		text.Draw(dst, p.Year.Gan, ft28, int(sx), int(sy), ColorGanZhi(p.Year.Gan))
		DrawProBar(dst, sx+28, sy-8, 64, 8, ColorGanZhi(p.Year.Gan), p.Year.HPHead, p.Year.HPMHead)
		sy += 48
		text.Draw(dst, p.Year.Zhi, ft28, int(sx), int(sy), ColorGanZhi(p.Year.Zhi))
		text.Draw(dst, p.Year.Body, ft14, int(sx), int(sy+16), ColorGanZhi(p.Year.Body))
		DrawProBar(dst, sx+28, sy+16-8, 64, 8, ColorGanZhi(p.Year.Zhi), p.Year.HPBody, p.Year.HPMBody)
		text.Draw(dst, p.Year.Legs, ft14, int(sx), int(sy+32), ColorGanZhi(p.Year.Legs))
		DrawProBar(dst, sx+28, sy+32-8, 64, 8, ColorGanZhi(p.Year.Legs), p.Year.HPLegs, p.Year.HPMLegs)
		text.Draw(dst, p.Year.Feet, ft14, int(sx), int(sy+48), ColorGanZhi(p.Year.Feet))
		DrawProBar(dst, sx+28, sy+48-8, 64, 8, ColorGanZhi(p.Year.Feet), p.Year.HPFeet, p.Year.HPMFeet)
		sy -= 48

		sx += 96
		vector.StrokeRect(dst, sx, sy+2, 96, 96, 1, colorWhite, true)
		text.Draw(dst, p.Month.Gan, ft28, int(sx), int(sy), ColorGanZhi(p.Month.Gan))
		DrawProBar(dst, sx+28, sy-8, 64, 8, ColorGanZhi(p.Month.Gan), p.Month.HPHead, p.Month.HPMHead)
		sy += 48
		text.Draw(dst, p.Month.Zhi, ft28, int(sx), int(sy), ColorGanZhi(p.Month.Zhi))
		text.Draw(dst, p.Month.Body, ft14, int(sx), int(sy+16), ColorGanZhi(p.Month.Body))
		DrawProBar(dst, sx+28, sy+16-8, 64, 8, ColorGanZhi(p.Month.Body), p.Month.HPBody, p.Month.HPMBody)
		text.Draw(dst, p.Month.Legs, ft14, int(sx), int(sy+32), ColorGanZhi(p.Month.Legs))
		DrawProBar(dst, sx+28, sy+32-8, 64, 8, ColorGanZhi(p.Month.Legs), p.Month.HPLegs, p.Month.HPMLegs)
		text.Draw(dst, p.Month.Feet, ft14, int(sx), int(sy+48), ColorGanZhi(p.Month.Feet))
		DrawProBar(dst, sx+28, sy+48-8, 64, 8, ColorGanZhi(p.Month.Feet), p.Month.HPFeet, p.Month.HPMFeet)
		sy -= 48

		sx += 96
		vector.StrokeRect(dst, sx, sy+2, 96, 96, 1, colorWhite, true)
		text.Draw(dst, p.Day.Gan, ft28, int(sx), int(sy), ColorGanZhi(p.Day.Gan))
		DrawProBar(dst, sx+28, sy-8, 64, 8, ColorGanZhi(p.Day.Gan), p.Day.HPHead, p.Day.HPMHead)
		sy += 48
		text.Draw(dst, p.Day.Zhi, ft28, int(sx), int(sy), ColorGanZhi(p.Day.Zhi))
		text.Draw(dst, p.Day.Body, ft14, int(sx), int(sy+16), ColorGanZhi(p.Day.Body))
		DrawProBar(dst, sx+28, sy+16-8, 64, 8, ColorGanZhi(p.Day.Body), p.Day.HPBody, p.Day.HPMBody)
		text.Draw(dst, p.Day.Legs, ft14, int(sx), int(sy+32), ColorGanZhi(p.Day.Legs))
		DrawProBar(dst, sx+28, sy+32-8, 64, 8, ColorGanZhi(p.Day.Legs), p.Day.HPLegs, p.Day.HPMLegs)
		text.Draw(dst, p.Day.Feet, ft14, int(sx), int(sy+48), ColorGanZhi(p.Day.Feet))
		DrawProBar(dst, sx+28, sy+48-8, 64, 8, ColorGanZhi(p.Day.Feet), p.Day.HPFeet, p.Day.HPMFeet)
		sy -= 48

		sx += 96
		vector.StrokeRect(dst, sx, sy+2, 96, 96, 1, colorWhite, true)
		text.Draw(dst, p.Time.Gan, ft28, int(sx), int(sy), ColorGanZhi(p.Time.Gan))
		DrawProBar(dst, sx+28, sy-8, 64, 8, ColorGanZhi(p.Time.Gan), p.Time.HPHead, p.Time.HPMHead)
		sy += 48
		text.Draw(dst, p.Time.Zhi, ft28, int(sx), int(sy), ColorGanZhi(p.Time.Zhi))
		text.Draw(dst, p.Time.Body, ft14, int(sx), int(sy+16), ColorGanZhi(p.Time.Body))
		DrawProBar(dst, sx+28, sy+16-8, 64, 8, ColorGanZhi(p.Time.Body), p.Time.HPBody, p.Time.HPMBody)
		text.Draw(dst, p.Time.Legs, ft14, int(sx), int(sy+32), ColorGanZhi(p.Time.Legs))
		DrawProBar(dst, sx+28, sy+32-8, 64, 8, ColorGanZhi(p.Time.Legs), p.Time.HPLegs, p.Time.HPMLegs)
		text.Draw(dst, p.Time.Feet, ft14, int(sx), int(sy+48), ColorGanZhi(p.Time.Feet))
		DrawProBar(dst, sx+28, sy+48-8, 64, 8, ColorGanZhi(p.Time.Feet), p.Time.HPFeet, p.Time.HPMFeet)
		sy -= 48
	}

}

type CharBody struct {
	Gan  string //干为头
	Zhi  string //支为体
	Body string //本气为身
	Legs string //中气为腿
	Feet string //余气为足

	HPHead  int //干值
	HPBody  int //本气值
	HPLegs  int //中气值
	HPFeet  int //余气值
	HPMHead int //干为值Max
	HPMBody int //本气值Max
	HPMLegs int //中气值Max
	HPMFeet int //余气值Max
}

func NewCharBody(gan, zhi string, ganM, zhiM int) CharBody {
	return CharBody{Gan: gan, Zhi: zhi,
		Body: GetHideGan(zhi, 0),
		Legs: GetHideGan(zhi, 1), Feet: GetHideGan(zhi, 2),
		HPMHead: ganM, HPMBody: zhiM, HPMLegs: zhiM, HPMFeet: zhiM,
	}
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
	bz := lunar.GetEightChar()
	zhiY, zhiM, zhiD, zhiT := bz.GetYearZhi(), bz.GetMonthZhi(), bz.GetDayZhi(), bz.GetTimeZhi()
	p.Year = NewCharBody(bz.GetYearGan(), zhiY, HpGY, HpZY)
	p.Month = NewCharBody(bz.GetMonthGan(), zhiM, HpGM, HpZM)
	p.Day = NewCharBody(bz.GetDayGan(), zhiD, HpGD, HpZD)
	p.Time = NewCharBody(bz.GetTimeGan(), zhiT, HpGT, HpZT)
	p.ShenShaY, p.ShenShaM, p.ShenShaD, p.ShenShaT = qimen.CalcShenSha(bz)

	yun := bz.GetYun(p.Gender)
	p.DaYun = yun.GetDaYunBy(7)
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

	p.ResetHP()
}

func (p *Player) ResetHP() {
	p.Year.HPHead = HpGY
	if p.Year.Feet != "" {
		p.Year.HPBody = HpZY * 0.6
		p.Year.HPLegs = HpZY * 0.3
		p.Year.HPFeet = HpZY * 0.1
	} else if p.Year.Legs != "" {
		p.Year.HPBody = HpZY * 0.7
		p.Year.HPLegs = HpZY * 0.3
		p.Year.HPFeet = 0
	} else {
		p.Year.HPBody = HpZY
		p.Year.HPLegs = 0
		p.Year.HPFeet = 0
	}
	p.Month.HPHead = HpGM
	if p.Month.Feet != "" {
		p.Month.HPBody = HpZM * 0.6
		p.Month.HPLegs = HpZM * 0.3
		p.Month.HPFeet = HpZM * 0.1
	} else if p.Month.Legs != "" {
		p.Month.HPBody = HpZM * 0.7
		p.Month.HPLegs = HpZM * 0.3
		p.Month.HPFeet = 0
	} else {
		p.Month.HPBody = HpZM
		p.Month.HPLegs = 0
		p.Month.HPFeet = 0
	}
	p.Day.HPHead = HpGD
	if p.Day.Feet != "" {
		p.Day.HPBody = HpZD * 0.6
		p.Day.HPLegs = HpZD * 0.3
		p.Day.HPFeet = HpZD * 0.1
	} else if p.Day.Legs != "" {
		p.Day.HPBody = HpZD * 0.7
		p.Day.HPLegs = HpZD * 0.3
		p.Day.HPFeet = 0
	} else {
		p.Day.HPBody = HpZD
		p.Day.HPLegs = 0
		p.Day.HPFeet = 0
	}
	p.Time.HPHead = HpGT
	if p.Time.Feet != "" {
		p.Time.HPBody = HpZT * 0.6
		p.Time.HPLegs = HpZT * 0.3
		p.Time.HPFeet = HpZT * 0.1
	} else if p.Time.Legs != "" {
		p.Time.HPBody = HpZT * 0.7
		p.Time.HPLegs = HpZT * 0.3
		p.Time.HPFeet = 0
	} else {
		p.Time.HPBody = HpZT
		p.Time.HPLegs = 0
		p.Time.HPFeet = 0
	}

}

func (p *Player) UpdateHp() {
	//合取1
	//制取1破1/取2
	//印生2
	//枭生2?
	//食泄2
	//伤泄2?
	//先天后地 从年到时
	//年干对月干 年干对年支 月干对月支 月干对日干

}

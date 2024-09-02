package world

import (
	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/6tail/lunar-go/calendar"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"qimen/qimen"
	"qimen/ui"
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
	1: {100},
	2: {70, 30},
	3: {60, 30, 10},
}

type EightCharPan struct {
	X, Y float32
	//YMDH   string
	inited     bool
	FYear      *CharBody //流年通用
	FMonth     *CharBody //流月通用
	FDay       *CharBody //流日通用
	FTime      *CharBody //流时通用
	Player     *Player   //玩家
	brightness float32

	btnMove  *ui.TextButton
	btnBirth *ui.TextButton
	uis      map[ui.IUIPanel]struct{}
}

func NewEightCharPan(x, y float32) *EightCharPan {
	p := &EightCharPan{
		X: x, Y: y,
		uis: make(map[ui.IUIPanel]struct{}),
		//btnMove:  ui.NewTextButton(int(x+4), int(y), "+ ", colorWhite, true),
	}
	p.btnBirth = ui.NewTextButton(int(x+158), int(y), "择辰", colorWhite, true)
	p.btnBirth.SetOnClick(func(b *ui.Button) {
		UIShowSelect()
	})
	p.uis[p.btnBirth] = struct{}{}
	return p
}

func (g *EightCharPan) Init() {
	cal := ThisGame.qmGame.Lunar
	g.FYear = NewCharBody(cal.GetYearGan(), cal.GetYearZhi(), HpGY, HpZY)
	g.FMonth = NewCharBody(cal.GetMonthGan(), cal.GetMonthZhi(), HpGM, HpZM)
	g.FDay = NewCharBody(cal.GetDayGan(), cal.GetDayZhi(), HpGD, HpZD)
	g.FTime = NewCharBody(cal.GetTimeGan(), cal.GetTimeZhi(), HpGT, HpZT)
	g.Player = &Player{}
	g.Player.Reset(cal, GenderMale)
}

func (g *EightCharPan) Update() error {
	if !g.inited {
		g.Init()
		g.inited = true
	}
	for panel := range g.uis {
		panel.Update()
	}

	cal := ThisGame.qmGame.Lunar
	if g.FYear.Gan != cal.GetYearGan() || g.FYear.Zhi != cal.GetYearZhi() {
		g.FYear = NewCharBody(cal.GetYearGan(), cal.GetYearZhi(), HpGY, HpZY)
	}
	if g.FMonth.Gan != cal.GetMonthGan() || g.FMonth.Zhi != cal.GetMonthZhi() {
		g.FMonth = NewCharBody(cal.GetMonthGan(), cal.GetMonthZhi(), HpGM, HpZM)
	}
	if g.FDay.Gan != cal.GetDayGan() || g.FDay.Zhi != cal.GetDayZhi() {
		g.FDay = NewCharBody(cal.GetDayGan(), cal.GetDayZhi(), HpGD, HpZD)
	}
	if g.FTime.Gan != cal.GetTimeGan() || g.FTime.Zhi != cal.GetTimeZhi() {
		g.FTime = NewCharBody(cal.GetTimeGan(), cal.GetTimeZhi(), HpGT, HpZT)
	}
	p := g.Player
	var change bool
	if p.Yun == nil {
		change = true
	} else {
		LYear := cal.GetYear()
		if p.YunIdx == 0 {
			if p.DaYunA[0].GetXiaoYun()[p.YunIdx0].GetYear() != LYear {
				change = true
			}
		} else {
			daYun := p.DaYunA[p.YunIdx]
			if daYun.GetStartYear() <= LYear && LYear <= daYun.GetEndYear() {
				change = true
			}
		}
	}
	if change { //大运变化
		for i, daYun := range p.DaYunA {
			if daYun.GetStartYear() <= cal.GetYear() && cal.GetYear() <= daYun.GetEndYear() {
				if i == 0 {
					for j, xiaoYun := range daYun.GetXiaoYun() {
						if xiaoYun.GetYear() == cal.GetYear() {
							gz := xiaoYun.GetGanZhi()
							gan := string([]rune(gz)[0])
							zhi := string([]rune(gz)[1])
							p.Yun = NewCharBody(gan, zhi, HpGM, HpZM)
							p.YunIdx0 = j
							p.YunIdx = i
							break
						}
					}
				} else {
					gz := daYun.GetGanZhi()
					gan := string([]rune(gz)[0])
					zhi := string([]rune(gz)[1])
					p.Yun = NewCharBody(gan, zhi, HpGM, HpZM)
					p.YunIdx0 = 0
					p.YunIdx = i
				}
				break
			}
		}
	}

	p.UpdateHp()

	g.brightness += 1
	if 0xff < g.brightness {
		g.brightness = 0xff
	}
	return nil
}

func (g *EightCharPan) Draw(dst *ebiten.Image) {
	for panel := range g.uis {
		panel.Draw(dst)
	}

	ft12, _ := GetFontFace(12)
	ft14, _ := GetFontFace(14)
	ft28, _ := GetFontFace(28)
	cx, cy := g.X, g.Y
	p := g.Player
	bz := p.Lunar.GetEightChar()
	soul := bz.GetDayGan()
	//八字总览
	{
		sx, sy := cx, cy
		vector.StrokeRect(dst, sx, sy, 394, 370, 1, colorWhite, true)
		sx += 16
		sy += 64
		text.Draw(dst, "十神", ft14, int(sx), int(sy-32), colorWhite)
		text.Draw(dst, "天干", ft14, int(sx), int(sy-8), colorWhite)
		text.Draw(dst, "地支", ft14, int(sx), int(sy+32-8), colorWhite)
		text.Draw(dst, "本气", ft14, int(sx), int(sy+48), colorWhite)
		text.Draw(dst, "中气", ft14, int(sx), int(sy+64), colorWhite)
		text.Draw(dst, "余气", ft14, int(sx), int(sy+80), colorWhite)
		text.Draw(dst, "纳音", ft14, int(sx), int(sy+96), colorWhite)
		text.Draw(dst, "星运", ft14, int(sx), int(sy+112), colorWhite) //地势/长生/星运
		text.Draw(dst, "自坐", ft14, int(sx), int(sy+128), colorWhite)
		text.Draw(dst, "空亡", ft14, int(sx), int(sy+144), colorWhite)
		text.Draw(dst, "小运", ft14, int(sx), int(sy+160), colorWhite)
		text.Draw(dst, "大运", ft14, int(sx), int(sy+160+16), colorWhite)
		text.Draw(dst, "神煞", ft14, int(sx), int(sy+160+32), colorWhite)
		sx += 48
		text.Draw(dst, "年柱", ft14, int(sx), int(sy-48), colorWhite)
		text.Draw(dst, bz.GetYearShiShenGan(), ft14, int(sx), int(sy-32), colorWhite)
		DrawFlow(dst, sx, sy, soul, p.Year)
		text.Draw(dst, strings.Join(p.Fates0, " "), ft14, int(sx), int(sy+160), colorWhite)
		text.Draw(dst, strings.Join(p.Fates, " "), ft14, int(sx), int(sy+160+16), colorWhite)
		text.Draw(dst, strings.Join(p.ShenShaY, "\n"), ft12, int(sx), int(sy+160+32), colorWhite)
		sx += 48
		text.Draw(dst, "月柱", ft14, int(sx), int(sy-48), colorWhite)
		text.Draw(dst, bz.GetMonthShiShenGan(), ft14, int(sx), int(sy-32), colorWhite)
		DrawFlow(dst, sx, sy, soul, p.Month)
		text.Draw(dst, strings.Join(p.ShenShaM, "\n"), ft12, int(sx), int(sy+160+32), colorWhite)
		sx += 48
		text.Draw(dst, "元"+GenderName[p.Gender], ft14, int(sx), int(sy-32), colorWhite)
		DrawFlow(dst, sx, sy, soul, p.Day)
		text.Draw(dst, strings.Join(p.ShenShaD, "\n"), ft12, int(sx), int(sy+160+32), colorWhite)
		sx += 48
		text.Draw(dst, "时柱", ft14, int(sx), int(sy-48), colorWhite)
		text.Draw(dst, bz.GetTimeShiShenGan(), ft14, int(sx), int(sy-32), colorWhite)
		DrawFlow(dst, sx, sy, soul, p.Time)
		text.Draw(dst, strings.Join(p.ShenShaT, "\n"), ft12, int(sx), int(sy+160+32), colorWhite)
		sx += 48
		vector.StrokeLine(dst, sx-3, sy-28, sx-3, sy+148, 1, colorWhite, true)
		if p.YunIdx == 0 {
			text.Draw(dst, "小运", ft14, int(sx), int(sy-48), colorWhite)
		} else {
			text.Draw(dst, "大运", ft14, int(sx), int(sy-48), colorWhite)
		}
		if p.Yun != nil {
			text.Draw(dst, LunarUtil.SHI_SHEN[soul+p.Yun.Gan], ft14, int(sx), int(sy-32), colorWhite)
			DrawFlow(dst, sx, sy, soul, p.Yun)
		}
		sx += 48
		text.Draw(dst, "流年", ft14, int(sx), int(sy-48), colorWhite)
		text.Draw(dst, LunarUtil.SHI_SHEN[soul+g.FYear.Gan], ft14, int(sx), int(sy-32), colorWhite)
		DrawFlow(dst, sx, sy, soul, g.FYear)
		sx += 48
		text.Draw(dst, "流月", ft14, int(sx), int(sy-48), colorWhite)
		text.Draw(dst, LunarUtil.SHI_SHEN[soul+g.FMonth.Gan], ft14, int(sx), int(sy-32), colorWhite)
		DrawFlow(dst, sx, sy, soul, g.FMonth)

	}
	//竖象 年头颈/月胸腹/日腹股/时腿足 年额/月目/日鼻/时口 干左支右?
	{
		sx, sy := cx+408, cy
		mx := int(sx + 28)
		w := float32(74)
		vector.StrokeRect(dst, sx, sy, w, 64, 1, colorWhite, true)    //头
		vector.StrokeRect(dst, sx, sy+64, w, 32, 1, colorWhite, true) //颈
		vector.StrokeRect(dst, sx, sy+96, w, 64, 1, colorWhite, true) //胸
		//vector.StrokeRect(dst, sx-36, sy+96, 32, 190, 1, colorWhite, true)  //胳膊手1
		//vector.StrokeRect(dst, sx+w+4, sy+96, 32, 190, 1, colorWhite, true) //胳膊手2
		vector.StrokeRect(dst, sx, sy+96+64, w, 32, 1, colorWhite, true)   //腹
		vector.StrokeRect(dst, sx, sy+96*2, w, 64, 1, colorWhite, true)    //小腹
		vector.StrokeRect(dst, sx, sy+96*2+64, w, 32, 1, colorWhite, true) //股
		vector.StrokeRect(dst, sx, sy+96*3, w, 64, 1, colorWhite, true)    //腿
		vector.StrokeRect(dst, sx, sy+96*3+64, w, 32, 1, colorWhite, true) //足
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
		sx, sy := cx+8, cy+420
		g.DrawCharHP(dst, sx, sy, p.Year)
		g.DrawCharHP(dst, sx+96, sy, p.Month)
		g.DrawCharHP(dst, sx+96*2, sy, p.Day)
		g.DrawCharHP(dst, sx+96*3, sy, p.Time)

	}

}

func (g *EightCharPan) DrawCharHP(dst *ebiten.Image, sx, sy float32, body *CharBody) {
	ft14, _ := GetFontFace(14)
	ft28, _ := GetFontFace(28)
	vector.StrokeRect(dst, sx, sy+2, 96, 96, 1, colorWhite, true)
	text.Draw(dst, body.Gan, ft28, int(sx), int(sy), ColorGanZhi(body.Gan))
	DrawProBar(dst, sx+28, sy-8, 64, 8, ColorGanZhi(body.Gan), body.HPHead, body.HPMHead)
	sy += 48
	text.Draw(dst, body.Zhi, ft28, int(sx), int(sy), ColorGanZhi(body.Zhi))
	text.Draw(dst, body.Body, ft14, int(sx), int(sy+16), ColorGanZhi(body.Body))
	DrawProBar(dst, sx+28, sy+16-8, 64, 8, ColorGanZhi(body.Zhi), body.HPBody, body.HPMBody) //横HP
	//DrawProBarV(dst, sx+28, sy-16, 8, 64, ColorGanZhi(body.Zhi), body.HPBody, body.HPMBody) //纵HP
	if body.Legs != "" {
		text.Draw(dst, body.Legs, ft14, int(sx), int(sy+32), ColorGanZhi(body.Legs))
		DrawProBar(dst, sx+28, sy+32-8, 64, 8, ColorGanZhi(body.Legs), body.HPLegs, body.HPMLegs)
		//DrawProBarV(dst, sx+28+24, sy-16, 8, 64, ColorGanZhi(body.Legs), body.HPLegs, body.HPMLegs)
	}
	if body.Feet != "" {
		text.Draw(dst, body.Feet, ft14, int(sx), int(sy+48), ColorGanZhi(body.Feet))
		DrawProBar(dst, sx+28, sy+48-8, 64, 8, ColorGanZhi(body.Feet), body.HPFeet, body.HPMFeet)
		//DrawProBarV(dst, sx+28+48, sy-16, 8, 64, ColorGanZhi(body.Feet), body.HPFeet, body.HPMFeet)
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

func NewCharBody(gan, zhi string, ganM, zhiM int) *CharBody {
	return &CharBody{Gan: gan, Zhi: zhi,
		Body: GetHideGan(zhi, 0),
		Legs: GetHideGan(zhi, 1), Feet: GetHideGan(zhi, 2),
		HPMHead: ganM, HPMBody: zhiM, HPMLegs: zhiM, HPMFeet: zhiM,
	}
}
func (c *CharBody) InitHP(maxHp int) {
	if c.Feet != "" {
		c.HPBody = maxHp * HideGanVal[3][0] / 100
		c.HPLegs = maxHp * HideGanVal[3][1] / 100
		c.HPFeet = maxHp * HideGanVal[3][2] / 100
	} else if c.Legs != "" {
		c.HPBody = maxHp * HideGanVal[2][0] / 100
		c.HPLegs = maxHp * HideGanVal[2][1] / 100
		c.HPFeet = 0
	} else {
		c.HPBody = maxHp * HideGanVal[1][0] / 100
		c.HPLegs = 0
		c.HPFeet = 0
	}
}

type Player struct {
	Gender int //性别
	Lunar  *calendar.Lunar
	Year   *CharBody //年柱
	Month  *CharBody //月柱
	Day    *CharBody //日柱
	Time   *CharBody //时柱
	Yun    *CharBody //当前大运

	DaYunA          []*calendar.DaYun //大运
	YunIdx0, YunIdx int
	Fates0          []string //小运名
	Fates           []string //大运名

	ShenShaY []string //神煞
	ShenShaM []string //神煞
	ShenShaD []string //神煞
	ShenShaT []string //神煞
}

func (p *Player) Reset(lunar *calendar.Lunar, gender int) {
	p.Lunar = lunar
	p.Gender = gender
	bz := lunar.GetEightChar()
	zhiY, zhiM, zhiD, zhiT := bz.GetYearZhi(), bz.GetMonthZhi(), bz.GetDayZhi(), bz.GetTimeZhi()
	p.Year = NewCharBody(bz.GetYearGan(), zhiY, HpGY, HpZY)
	p.Month = NewCharBody(bz.GetMonthGan(), zhiM, HpGM, HpZM)
	p.Day = NewCharBody(bz.GetDayGan(), zhiD, HpGD, HpZD)
	p.Time = NewCharBody(bz.GetTimeGan(), zhiT, HpGT, HpZT)
	p.ShenShaY, p.ShenShaM, p.ShenShaD, p.ShenShaT = qimen.CalcShenSha(bz)

	yun := bz.GetYun(p.Gender)
	p.DaYunA = yun.GetDaYunBy(7)
	p.Fates0 = nil
	p.Fates = nil
	for i, daYun := range p.DaYunA {
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
	p.Yun = nil

	p.ResetHP()
}

func (p *Player) ResetHP() {
	p.Year.HPHead = HpGY
	p.Year.InitHP(HpZY)
	p.Month.HPHead = HpGM
	p.Month.InitHP(HpZM)
	p.Day.HPHead = HpGD
	p.Day.InitHP(HpZD)
	p.Time.HPHead = HpGT
	p.Time.InitHP(HpZT)
}

func (p *Player) UpdateHp() {
	//合取1
	//制取1破1/取2
	//印生2
	//枭生2?
	//食泄2
	//伤泄2?
	//先天后地 从年到时
	//年干-月干 日干-时干 月干-日干

	//年干-年支 月干-月支 日干-日支 时干-时支
	//年支-月支 日支-时支 月支-日支
	//年干--日干 月干--时干
	//年支--日支 月支--时支
	//年干---时干
	//年支---时支

}

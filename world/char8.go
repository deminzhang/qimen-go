package world

import (
	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/6tail/lunar-go/calendar"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"qimen/qimen"
	"qimen/ui"
	"qimen/util"
	"strings"
)

// 八字能量比例
// 年干8 月干12 日元12 时干12
// 年支4 月支40 日支12 时支12
const (
	HpGanYear  = 80
	HpZhiYear  = 40
	HpGanMonth = 120
	HpZhiMonth = 400
	HpGanDay   = 120
	HpZhiDay   = 120
	HpGanTime  = 120
	HpZhiTime  = 120
)

// HideGanVal 藏干值比例
var HideGanVal = map[int][]int{
	1: {100},
	2: {70, 30},
	3: {60, 30, 10},
}

type Char8Pan struct {
	X, Y     float32
	FYear    *CharBody //流年通用
	FMonth   *CharBody //流月通用
	FDay     *CharBody //流日通用
	FTime    *CharBody //流时通用
	Player   *Player   //玩家
	inited   bool
	BodyShow bool

	ui.Container
	count int
}

func NewChar8Pan(x, y float32) *Char8Pan {
	p := &Char8Pan{
		X: x, Y: y,
		BodyShow: true,
	}
	//btnMove:=  ui.NewTextButton(int(x+4), int(y), "+ ", colorWhite, true)
	btnBirth := ui.NewTextButton(int(x+146), int(y+3), "命造", colorWhite, true)
	btnBirth.SetOnClick(func(b *ui.Button) {
		UIShowSelect()
	})
	cbShowBody := ui.NewCheckBox(int(x+410), int(y+-18), "身象")
	cbShowBody.SetOnCheckChanged(func(c *ui.CheckBox) {
		p.BodyShow = c.Checked()
	})
	cbShowBody.SetChecked(true)
	p.Add(btnBirth, cbShowBody)
	return p
}

func (g *Char8Pan) Init() {
	cal := ThisGame.qmGame.Lunar
	g.FYear = NewCharBody(cal.GetYearGan(), cal.GetYearZhi(), HpGanYear, HpZhiYear)
	g.FMonth = NewCharBody(cal.GetMonthGan(), cal.GetMonthZhi(), HpGanMonth, HpZhiMonth)
	g.FDay = NewCharBody(cal.GetDayGan(), cal.GetDayZhi(), HpGanDay, HpZhiDay)
	g.FTime = NewCharBody(cal.GetTimeGan(), cal.GetTimeZhi(), HpGanTime, HpZhiTime)
	g.Player = &Player{}
	g.Player.Reset(cal, GenderMale)
	g.inited = true
}

func (g *Char8Pan) SetPos(x, y float32) {
	g.X, g.Y = x, y
}

func (g *Char8Pan) Update() {
	if !g.inited {
		g.Init()
	}
	g.Container.Update()
	g.count++
	g.count %= 60

	cal := ThisGame.qmGame.Lunar
	p := g.Player
	if g.FYear.Gan != cal.GetYearGan() || g.FYear.Zhi != cal.GetYearZhi() {
		g.FYear = NewCharBody(cal.GetYearGan(), cal.GetYearZhi(), HpGanYear, HpZhiYear)
		p.UpdateCount = 10
	}
	if g.FMonth.Gan != cal.GetMonthGan() || g.FMonth.Zhi != cal.GetMonthZhi() {
		g.FMonth = NewCharBody(cal.GetMonthGan(), cal.GetMonthZhi(), HpGanMonth, HpZhiMonth)
		p.UpdateCount = 10
	}
	if g.FDay.Gan != cal.GetDayGan() || g.FDay.Zhi != cal.GetDayZhi() {
		g.FDay = NewCharBody(cal.GetDayGan(), cal.GetDayZhi(), HpGanDay, HpZhiDay)
		p.UpdateCount = 10
	}
	if g.FTime.Gan != cal.GetTimeGan() || g.FTime.Zhi != cal.GetTimeZhi() {
		g.FTime = NewCharBody(cal.GetTimeGan(), cal.GetTimeZhi(), HpGanTime, HpZhiTime)
		p.UpdateCount = 10
	}
	var changeYun bool
	if p.FYun == nil {
		changeYun = true
	} else {
		LYear := cal.GetYear()
		if p.YunIdx == 0 {
			if p.yuns[0].GetXiaoYun()[p.YunIdx0].GetYear() != LYear {
				changeYun = true
			}
		} else {
			daYun := p.yuns[p.YunIdx]
			if daYun.GetStartYear() <= LYear && LYear <= daYun.GetEndYear() {
				changeYun = true
			}
		}
	}
	if changeYun { //大运变化
		if cal.GetYear() < p.yuns[0].GetStartYear() {
			p.FYun = nil //未出生,穿越过去
		} else {
			for i, daYun := range p.yuns {
				if daYun.GetStartYear() <= cal.GetYear() && cal.GetYear() <= daYun.GetEndYear() {
					if i == 0 {
						for j, xiaoYun := range daYun.GetXiaoYun() {
							if xiaoYun.GetYear() == cal.GetYear() {
								gz := xiaoYun.GetGanZhi()
								gan := string([]rune(gz)[0])
								zhi := string([]rune(gz)[1])
								p.FYun = NewCharBody(gan, zhi, HpGanMonth, HpZhiMonth)
								p.YunIdx0 = j
								p.YunIdx = i
								break
							}
						}
					} else {
						gz := daYun.GetGanZhi()
						gan := string([]rune(gz)[0])
						zhi := string([]rune(gz)[1])
						p.FYun = NewCharBody(gan, zhi, HpGanMonth, HpZhiMonth)
						p.YunIdx0 = 0
						p.YunIdx = i
					}
					p.UpdateCount = 10
					break
				} else { //超寿,修仙了,加运
					if i == len(p.yuns)-1 {
						for j := i; j < i+10; j++ {
							p.yuns = append(p.yuns, calendar.NewDaYun(p.yun, j))
						}
						break
					}
				}
			}
		}
	}
	if g.count%10 == 0 {
		g.UpdateHp(p)
	}
}

func (g *Char8Pan) UpdateHp(p *Player) {
	if p.UpdateCount == 0 {
		return
	}
	p.UpdateCount--
	//先天后地 从年到时
	//年干-月干 日干-时干 月干-日干 论合冲
	//年干-年支 月干-月支 日干-日支 时干-时支 论生旺衰死
	//年支-月支 日支-时支 月支-日支 论合冲刑破害
	//年干--日干 月干--时干 论冲
	//年支--日支 月支--时支 论合冲刑破害
	//年干---时干 论冲
	//年支---时支 论合冲刑破害
	CharBodyInteractive(p.Year, p.Month, 6, 1)
	CharBodyInteractive(p.Day, p.Time, 6, 1)
	CharBodyInteractive(p.Month, p.Day, 6, 1)
	CharBodyInteractive(p.Year, p.Day, 4, 1)
	CharBodyInteractive(p.Month, p.Time, 4, 1)
	CharBodyInteractive(p.Year, p.Time, 2, 1)

	CharBodyInteractive(p.FYun, p.Year, 4, 0)
	CharBodyInteractive(p.FYun, p.Month, 4, 0)
	CharBodyInteractive(g.FYear, p.Year, 4, 0)
	CharBodyInteractive(g.FMonth, p.Month, 4, 0)
	CharBodyInteractive(g.FDay, p.Day, 4, 0)
	CharBodyInteractive(g.FTime, p.Time, 4, 0)
}

func (g *Char8Pan) Draw(dst *ebiten.Image) {
	g.Container.Draw(dst)

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
		vector.StrokeRect(dst, sx, sy, util.If[float32](g.BodyShow, 400, 480), 384, 1, colorWhite, true)
		sx += 4
		sy += 64
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
		text.Draw(dst, "流年", ft14, int(sx), int(sy+160+32), colorWhite)
		text.Draw(dst, "神煞", ft14, int(sx), int(sy+160+48), colorWhite)
		sx += 48
		text.Draw(dst, "年柱", ft14, int(sx), int(sy-48), colorWhite)
		text.Draw(dst, bz.GetYearShiShenGan(), ft14, int(sx), int(sy-32), colorWhite)
		DrawFlow(dst, sx, sy, soul, p.Year)
		text.Draw(dst, strings.Join(p.Fates0, " "), ft14, int(sx), int(sy+160), colorWhite)
		text.Draw(dst, strings.Join(p.Fates, " "), ft14, int(sx), int(sy+160+16), colorWhite)
		text.Draw(dst, strings.Join(p.ShenShaY, "\n"), ft12, int(sx), int(sy+160+48), colorWhite)
		sx += 48
		text.Draw(dst, "月柱", ft14, int(sx), int(sy-48), colorWhite)
		text.Draw(dst, bz.GetMonthShiShenGan(), ft14, int(sx), int(sy-32), colorWhite)
		DrawFlow(dst, sx, sy, soul, p.Month)
		text.Draw(dst, strings.Join(p.ShenShaM, "\n"), ft12, int(sx), int(sy+160+48), colorWhite)
		sx += 48
		text.Draw(dst, "元"+GenderName[p.Gender], ft14, int(sx), int(sy-32), colorWhite)
		DrawFlow(dst, sx, sy, soul, p.Day)
		text.Draw(dst, strings.Join(p.ShenShaD, "\n"), ft12, int(sx), int(sy+160+48), colorWhite)
		sx += 48
		text.Draw(dst, "时柱", ft14, int(sx), int(sy-48), colorWhite)
		text.Draw(dst, bz.GetTimeShiShenGan(), ft14, int(sx), int(sy-32), colorWhite)
		DrawFlow(dst, sx, sy, soul, p.Time)
		text.Draw(dst, strings.Join(p.ShenShaT, "\n"), ft12, int(sx), int(sy+160+48), colorWhite)
		sx += 48
		vector.StrokeLine(dst, sx-3, sy-28, sx-3, sy+148, 1, colorWhite, true)
		if p.YunIdx == 0 {
			text.Draw(dst, "小运", ft14, int(sx), int(sy-48), colorWhite)
		} else {
			text.Draw(dst, "大运", ft14, int(sx), int(sy-48), colorWhite)
		}
		if p.FYun != nil {
			text.Draw(dst, LunarUtil.SHI_SHEN[soul+p.FYun.Gan], ft14, int(sx), int(sy-32), colorWhite)
			DrawFlow(dst, sx, sy, soul, p.FYun)
		}
		sx += 48
		vector.StrokeLine(dst, sx-3, sy-28, sx-3, sy+148, 1, colorWhite, true)
		text.Draw(dst, "流年", ft14, int(sx), int(sy-48), colorWhite)
		text.Draw(dst, LunarUtil.SHI_SHEN[soul+g.FYear.Gan], ft14, int(sx), int(sy-32), colorWhite)
		DrawFlow(dst, sx, sy, soul, g.FYear)
		sx += 48
		text.Draw(dst, "流月", ft14, int(sx), int(sy-48), colorWhite)
		text.Draw(dst, LunarUtil.SHI_SHEN[soul+g.FMonth.Gan], ft14, int(sx), int(sy-32), colorWhite)
		DrawFlow(dst, sx, sy, soul, g.FMonth)
		if !g.BodyShow {
			sx += 48
			text.Draw(dst, "流日", ft14, int(sx), int(sy-48), colorWhite)
			text.Draw(dst, LunarUtil.SHI_SHEN[soul+g.FDay.Gan], ft14, int(sx), int(sy-32), colorWhite)
			DrawFlow(dst, sx, sy, soul, g.FDay)
			sx += 48
			text.Draw(dst, "流时", ft14, int(sx), int(sy-48), colorWhite)
			text.Draw(dst, LunarUtil.SHI_SHEN[soul+g.FTime.Gan], ft14, int(sx), int(sy-32), colorWhite)
			DrawFlow(dst, sx, sy, soul, g.FTime)
		}

	}
	//竖象 身体全息
	//年头颈/月胸腹/日腹股/时腿足
	//年额/月目/日鼻/时口 干左支右?
	//男甲<申 左大右小 右近视度高
	//女乙>未 左大右小 左双右单
	if g.BodyShow {
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
	//横象 年祖 月父母 日夫妻 时子孙 干动支静 干为军支为营 干为官支为民
	{
		sx, sy := cx+8, cy+420
		g.DrawCharHP(dst, sx, sy, p.Year)
		g.DrawCharHP(dst, sx+96, sy, p.Month)
		g.DrawCharHP(dst, sx+96*2, sy, p.Day)
		g.DrawCharHP(dst, sx+96*3, sy, p.Time)

	}

}

func (g *Char8Pan) DrawCharHP(dst *ebiten.Image, sx, sy float32, body *CharBody) {
	ft14, _ := GetFontFace(14)
	ft28, _ := GetFontFace(28)
	vector.StrokeRect(dst, sx, sy, 96, 80, 1, colorWhite, true)
	text.Draw(dst, body.Gan, ft28, int(sx), int(sy), ColorGanZhi(body.Gan))
	DrawProBar(dst, sx+28, sy-8, 64, 8, ColorGanZhi(body.Gan), body.HPHead, body.HPMHead)
	sy += 26
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

func NewCharBody(gan, zhi string, ganMax, zhiMax int) *CharBody {
	return &CharBody{Gan: gan, Zhi: zhi,
		Body:    GetHideGan(zhi, 0),
		Legs:    GetHideGan(zhi, 1),
		Feet:    GetHideGan(zhi, 2),
		HPMHead: ganMax, HPMBody: zhiMax, HPMLegs: zhiMax, HPMFeet: zhiMax,
	}
}
func (c *CharBody) initZhiHP(maxHp int) {
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
	FYun   *CharBody //大运

	yun             *calendar.Yun     //运
	yuns            []*calendar.DaYun //大运集
	YunIdx0, YunIdx int               //当前大运小运索引
	Fates0          []string          //小运名
	Fates           []string          //大运名

	UpdateCount int

	ShenShaY []string //神煞年
	ShenShaM []string //神煞月
	ShenShaD []string //神煞日
	ShenShaT []string //神煞时
	//ShenShaYY []string //神煞大运
	//ShenShaFY []string //神煞流年
	//ShenShaFM []string //神煞流月
	//ShenShaFD []string //神煞流日
	//ShenShaFT []string //神煞流时
}

func (p *Player) Reset(lunar *calendar.Lunar, gender int) {
	p.Lunar = lunar
	p.Gender = gender
	bz := lunar.GetEightChar()
	zhiY, zhiM, zhiD, zhiT := bz.GetYearZhi(), bz.GetMonthZhi(), bz.GetDayZhi(), bz.GetTimeZhi()
	p.Year = NewCharBody(bz.GetYearGan(), zhiY, HpGanYear, HpZhiYear)
	p.Month = NewCharBody(bz.GetMonthGan(), zhiM, HpGanMonth, HpZhiMonth)
	p.Day = NewCharBody(bz.GetDayGan(), zhiD, HpGanDay, HpZhiDay)
	p.Time = NewCharBody(bz.GetTimeGan(), zhiT, HpGanTime, HpZhiTime)
	p.ShenShaY, p.ShenShaM, p.ShenShaD, p.ShenShaT = qimen.CalcShenSha(bz)

	yun := bz.GetYun(p.Gender)
	p.yun = yun
	p.yuns = yun.GetDaYun()
	p.Fates0 = nil
	p.Fates = nil
	for i, daYun := range p.yuns {
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
	p.FYun = nil
	p.UpdateCount = 10
	p.resetHP()
}

func (p *Player) resetHP() {
	p.Year.HPHead = HpGanYear
	p.Year.initZhiHP(HpZhiYear)
	p.Month.HPHead = HpGanMonth
	p.Month.initZhiHP(HpZhiMonth)
	p.Day.HPHead = HpGanDay
	p.Day.initZhiHP(HpZhiDay)
	p.Time.HPHead = HpGanTime
	p.Time.initZhiHP(HpZhiTime)
}

func CharBodyInteractive(a, b *CharBody, force int, reduce int) {
	if a == nil || b == nil {
		return
	}
	//金赖土生，土多金埋。	土赖火生，火多土焦。
	//火赖木生，木多火炽。	木赖水生，水多木漂。
	//水赖金生，金多水浊。
	//水空则流，木空则损，土空则陷，金空则则响，火空则发。
	//旺木喜金，旺火喜水，旺土喜木，旺金喜火，旺水喜土。
	//木怕金旺，火怕水旺，土怕木旺，金怕火旺，水怕土旺。
	//水弱则爱金，金弱则爱土，土弱则爱火，火弱则爱木，木弱则爱水。
	//水衰不生木，木衰不生火，火衰不生土，土衰不生金，金衰不生水。
	//春土不克水，夏金不克木，季水不克火，秋木不克土，冬火不克金。
	gg := LunarUtil.SHI_SHEN[a.Gan+b.Gan]
	switch gg {
	case "比肩", "劫财": //助分
		if a.HPHead < a.HPMHead && a.HPHead+force < b.HPHead-force {
			b.HPHead -= force * reduce
			a.HPHead += force
		}
		if b.HPHead < b.HPMHead && b.HPHead+force < a.HPHead-force {
			a.HPHead -= force * reduce
			b.HPHead += force
		}
	case "食神", "伤官": //泄
		if b.HPHead < b.HPMHead && b.HPHead+force < a.HPHead-force {
			if b.HPHead+force < (a.HPHead-force)/2 { //快泄
				a.HPHead -= force * reduce
				b.HPHead += force
			} else { //慢泄
				a.HPHead -= force / 2 * reduce
				b.HPHead += force / 2
			}
		}
	case "正印", "偏印": //生
		if a.HPHead < a.HPMHead && a.HPHead+force < b.HPHead-force {
			if a.HPHead+force < (b.HPHead-force)/2 { //快生
				b.HPHead -= force * reduce
				a.HPHead += force
			} else { //慢生
				b.HPHead -= force / 2 * reduce
				a.HPHead += force / 2
			}
		}
	case "正官": //娶
		if qimen.HeGan[a.Gan] == b.Gan { //合

		}
		if b.HPHead < b.HPMHead && b.HPHead+force < a.HPHead-force {
			a.HPHead -= force * reduce
			b.HPHead += force
		}
	case "七杀": //夺
		//if b.HPHead < b.HPMHead && b.HPHead+force < a.HPHead-force {
		if b.HPHead < b.HPMHead && a.HPHead > force {
			a.HPHead -= force * reduce
			b.HPHead += force / 2
		}
	case "正财": //嫁
		if qimen.HeGan[a.Gan] == b.Gan { //合

		}
		if a.HPHead < a.HPMHead && a.HPHead+force < b.HPHead-force {
			b.HPHead -= force * reduce
			a.HPHead += force
		}
	case "偏财": //耗
		if a.HPHead < a.HPMHead && b.HPHead > force {
			b.HPHead -= force * reduce
			a.HPHead += force / 2
		}
	}
	//支引干透
	reduce = 0
	ab := LunarUtil.SHI_SHEN[a.Gan+a.Body]
	switch ab {
	case "比肩", "劫财": //助
		if a.HPHead < a.HPMBody && a.HPHead+force < a.HPBody {
			a.HPHead += force
		}
	case "食神", "伤官": //泄
		if a.HPHead < a.HPMBody && a.HPHead+force < a.HPBody {
			a.HPHead += force
		}
	case "正印", "偏印": //生
	}
	if a.Legs != "" {
		al := LunarUtil.SHI_SHEN[a.Gan+a.Legs]
		switch al {
		case "比肩", "劫财": //助

		}
	}
	if a.Feet != "" {
		af := LunarUtil.SHI_SHEN[a.Gan+a.Feet]
		switch af {
		case "比肩", "劫财": //助

		}
	}
}

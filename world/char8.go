package world

import (
	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/6tail/lunar-go/calendar"
	"github.com/deminzhang/qimen-go/gui"
	"github.com/deminzhang/qimen-go/qimen"
	"github.com/deminzhang/qimen-go/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"strings"
)

// 八字能量比例
// 年干8 月干12 日元12 时干12
// 年支4 月支40 日支12 时支12
const (
	char8UIWidth  = 490
	char8UIHeight = 600

	HpGanYear  = 80
	HpZhiYear  = 40
	HpGanMonth = 120
	HpZhiMonth = 400
	HpGanDay   = 120
	HpZhiDay   = 120
	HpGanTime  = 120
	HpZhiTime  = 240
)

// HideGanVal 藏干值比例
var HideGanVal = map[int][]int{
	1: {100},
	2: {70, 30},
	3: {60, 30, 10},
}

type Char8Pan struct {
	gui.BaseUI
	Flow         *Body4  //流气
	Player       *Player //玩家
	BodyShow     bool    //身象
	OverviewShow bool    //总览

	count int
}

func NewChar8Pan(x, y int) *Char8Pan {
	p := &Char8Pan{
		BaseUI:       gui.BaseUI{X: x, Y: y, Visible: true, W: char8UIWidth, H: char8UIHeight},
		BodyShow:     false,
		OverviewShow: false,
	}
	btnBirth := gui.NewTextButton(350, 386, "命造", &colorYellow, &colorGray)
	btnBirth.SetOnClick(func(b *gui.Button) {
		oldBirthTime := ThisGame.char8.Player.Birth
		var oldBirthSolar *calendar.Solar
		if oldBirthTime != nil {
			oldBirthSolar = oldBirthTime.GetSolar()
		}
		UIShowSelectBirth(oldBirthSolar, ThisGame.char8.Player.Gender, func(birth *calendar.Solar, gender int) {
			ThisGame.char8.Player.Reset(calendar.NewLunarFromSolar(birth), gender)
		})
	})
	cbShowBody := gui.NewCheckBox(144, 0, "身象")
	cbShowBody.SetOnCheckChanged(func(c *gui.CheckBox) {
		p.BodyShow = c.Checked()
	})
	cbShowOverview := gui.NewCheckBox(0, 0, "总览")
	cbShowOverview.SetOnCheckChanged(func(c *gui.CheckBox) {
		p.OverviewShow = c.Checked()
	})
	//cbShowBody.SetChecked(false)
	btnMarry := gui.NewTextButton(350, 418, "择偶", &colorPink, &colorGray)
	btnSplit := gui.NewTextButton(350, 518, "和离", &colorGreen, &colorGray)
	btnMarry.SetOnClick(func(b *gui.Button) {
		mate := ThisGame.char8.Player.Mate
		if mate == nil {
			mate = &Player{}
			mate.Birth = calendar.NewLunarFromSolar(ThisGame.qmGame.Lunar.GetSolar())
		}
		solar := mate.Birth.GetSolar()
		var mateGender int
		if ThisGame.char8.Player.Gender == GenderFemale {
			mateGender = GenderMale
		}
		UIShowSelectBirth(solar, mateGender, func(solar *calendar.Solar, gender int) {
			mate.Reset(calendar.NewLunarFromSolar(solar), gender)
			ThisGame.char8.Player.Mate = mate
			btnSplit.Visible = true
		})
	})
	btnSplit.SetOnClick(func(b *gui.Button) {
		btnSplit.Visible = false
		ThisGame.char8.Player.Mate = nil
	})
	btnSplit.Visible = false

	p.AddChildren(btnBirth, cbShowBody, btnMarry, btnSplit, cbShowOverview)
	gui.ActiveUI(p)
	return p
}

func (g *Char8Pan) Init() {
	cal := ThisGame.qmGame.Lunar
	g.Flow = &Body4{
		Year:  NewCharBody(cal.GetYearGan(), cal.GetYearZhi(), HpGanYear, HpZhiYear, true),
		Month: NewCharBody(cal.GetMonthGan(), cal.GetMonthZhi(), HpGanMonth, HpZhiMonth, true),
		Day:   NewCharBody(cal.GetDayGan(), cal.GetDayZhi(), HpGanDay, HpZhiDay, true),
		Time:  NewCharBody(cal.GetTimeGan(), cal.GetTimeZhi(), HpGanTime, HpZhiTime, true),
	}

	g.Player = &Player{}
	g.Player.Reset(cal, GenderMale)
}

func (g *Char8Pan) SetPos(x, y int) {
	g.X, g.Y = x, y
}

func (g *Char8Pan) Update() {
	if g.Flow == nil || g.Player == nil {
		g.Init()
	}
	g.count++
	g.count %= 60

	cal := ThisGame.qmGame.Lunar
	p := g.Player
	var change bool
	if g.Flow.Year.Gan != cal.GetYearGan() || g.Flow.Year.Zhi != cal.GetYearZhi() {
		g.Flow.Year = NewCharBody(cal.GetYearGan(), cal.GetYearZhi(), HpGanYear, HpZhiYear, true)
		change = true
	}
	if g.Flow.Month.Gan != cal.GetMonthGan() || g.Flow.Month.Zhi != cal.GetMonthZhi() {
		g.Flow.Month = NewCharBody(cal.GetMonthGan(), cal.GetMonthZhi(), HpGanMonth, HpZhiMonth, true)
		change = true
	}
	if g.Flow.Day.Gan != cal.GetDayGan() || g.Flow.Day.Zhi != cal.GetDayZhi() {
		g.Flow.Day = NewCharBody(cal.GetDayGan(), cal.GetDayZhi(), HpGanDay, HpZhiDay, true)
		change = true
	}
	if g.Flow.Time.Gan != cal.GetTimeGan() || g.Flow.Time.Zhi != cal.GetTimeZhi() {
		g.Flow.Time = NewCharBody(cal.GetTimeGan(), cal.GetTimeZhi(), HpGanTime, HpZhiTime, true)
		change = true
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
		change = true
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
								p.FYun = NewCharBody(gan, zhi, HpGanMonth, HpZhiMonth, true)
								p.FYun.FlowEnergy = true
								p.YunIdx0 = j
								p.YunIdx = i
								break
							}
						}
					} else {
						gz := daYun.GetGanZhi()
						gan := string([]rune(gz)[0])
						zhi := string([]rune(gz)[1])
						p.FYun = NewCharBody(gan, zhi, HpGanMonth, HpZhiMonth, true)
						p.FYun.FlowEnergy = true
						p.YunIdx0 = 0
						p.YunIdx = i
					}
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
	if change {
		if p.FYun != nil {
			sss := qimen.CalcShenSha(p.Birth.GetEightChar(), p.FYun.GetGanZhi(),
				cal.GetYearInGanZhiExact(), cal.GetMonthInGanZhiExact(), cal.GetDayInGanZhiExact(), cal.GetTimeInGanZhi())
			p.ShenShaY, p.ShenShaM, p.ShenShaD, p.ShenShaT = sss[0], sss[1], sss[2], sss[3]
			p.ShenShaYY, p.ShenShaFY, p.ShenShaFM, p.ShenShaFD, p.ShenShaFT = sss[4], sss[5], sss[6], sss[7], sss[8]
		} else {
			sss := qimen.CalcShenSha(p.Birth.GetEightChar(),
				cal.GetYearInGanZhiExact(), cal.GetMonthInGanZhiExact(), cal.GetDayInGanZhiExact(), cal.GetTimeInGanZhi())
			p.ShenShaY, p.ShenShaM, p.ShenShaD, p.ShenShaT = sss[0], sss[1], sss[2], sss[3]
			p.ShenShaFY, p.ShenShaFM, p.ShenShaFD, p.ShenShaFT = sss[4], sss[5], sss[6], sss[7]
		}
	}

	if g.count%1 == 0 {
		g.UpdateHp(p)
	}
	g.BaseUI.Update()
}

func (g *Char8Pan) UpdateHp(p *Player) {
	//先天后地 从年到时
	//年干-月干 日干-时干 月干-日干 论合冲
	//年干-年支 月干-月支 日干-日支 时干-时支 论生旺衰死
	//年支-月支 日支-时支 月支-日支 论合冲刑破害
	//年干--日干 月干--时干 论冲
	//年支--日支 月支--时支 论合冲刑破害
	//年干---时干 论冲
	//年支---时支 论合冲刑破害
	//自身
	p.Year.InteractiveSelf(8)
	p.Day.InteractiveSelf(8)
	p.Day.InteractiveSelf(8)
	p.Time.InteractiveSelf(8)
	CharBodyInteractive(p.Year, p.Month, 6)
	CharBodyInteractive(p.Day, p.Time, 6)
	CharBodyInteractive(p.Month, p.Day, 6)
	CharBodyInteractive(p.Year, p.Day, 4)
	CharBodyInteractive(p.Month, p.Time, 4)
	CharBodyInteractive(p.Year, p.Time, 2)
	//运入年月
	CharBodyInteractive(p.FYun, p.Year, 4)
	CharBodyInteractive(p.FYun, p.Month, 4)
	//流气
	CharBodyInteractive(p.Year, g.Flow.Year, 4)
	CharBodyInteractive(p.Month, g.Flow.Month, 4)
	CharBodyInteractive(p.Day, g.Flow.Day, 4)
	CharBodyInteractive(p.Time, g.Flow.Time, 4)
	//命运共同体
	if p.Mate != nil {
		CharBodyInteractive(p.Mate.Year, p.Year, 2)
		CharBodyInteractive(p.Mate.Day, p.Day, 6)
	}
}

func (g *Char8Pan) Draw(dst *ebiten.Image) {
	ft12, _ := GetFontFace(12)
	ft14, _ := GetFontFace(14)
	ft28, _ := GetFontFace(28)
	//cx, cy := g.X, g.Y
	cx, cy := 0, 0
	p := g.Player
	bz := p.Birth.GetEightChar()
	soul := bz.GetDayGan()
	//八字总览
	if g.OverviewShow {
		sx, sy := cx, cy
		vector.StrokeRect(dst, float32(sx), float32(sy), util.If[float32](g.BodyShow, 400, 480),
			384, 1, colorWhite, true)
		sx += 4
		sy += 64
		text.Draw(dst, "十神", ft14, sx, sy-32, colorWhite)
		text.Draw(dst, "天干", ft14, sx, sy-8, colorWhite)
		text.Draw(dst, "地支", ft14, sx, sy+32-8, colorWhite)
		text.Draw(dst, "本气", ft14, sx, sy+48, colorWhite)
		text.Draw(dst, "中气", ft14, sx, sy+64, colorWhite)
		text.Draw(dst, "余气", ft14, sx, sy+80, colorWhite)
		text.Draw(dst, "纳音", ft14, sx, sy+96, colorWhite)
		text.Draw(dst, "地势", ft14, sx, sy+112, colorWhite) //地势/长生/星运
		text.Draw(dst, "自坐", ft14, sx, sy+128, colorWhite)
		text.Draw(dst, "空亡", ft14, sx, sy+144, colorWhite)
		text.Draw(dst, "小运", ft14, sx, sy+160, colorWhite)
		text.Draw(dst, "大运", ft14, sx, sy+160+16, colorWhite)
		text.Draw(dst, "流年", ft14, sx, sy+160+32, colorWhite)
		text.Draw(dst, "神煞", ft14, sx, sy+160+48, colorWhite)
		sx += 48
		text.Draw(dst, "年柱", ft14, sx, sy-48, colorWhite)
		text.Draw(dst, bz.GetYearShiShenGan(), ft14, sx, sy-32, colorWhite)
		DrawFlow(dst, sx, sy, soul, p.Year)
		text.Draw(dst, strings.Join(p.Fates0, " "), ft14, sx, sy+160, colorWhite)
		text.Draw(dst, strings.Join(p.Fates, " "), ft14, sx, sy+160+16, colorWhite)
		text.Draw(dst, strings.Join(p.ShenShaY, "\n"), ft12, sx, sy+160+48, colorWhite)
		sx += 48
		text.Draw(dst, "月柱", ft14, sx, sy-48, colorWhite)
		text.Draw(dst, bz.GetMonthShiShenGan(), ft14, sx, sy-32, colorWhite)
		DrawFlow(dst, sx, sy, soul, p.Month)
		text.Draw(dst, strings.Join(p.ShenShaM, "\n"), ft12, sx, sy+160+48, colorWhite)
		sx += 48
		text.Draw(dst, "元"+GenderName[p.Gender], ft14, sx, sy-32, colorWhite)
		DrawFlow(dst, sx, sy, soul, p.Day)
		text.Draw(dst, strings.Join(p.ShenShaD, "\n"), ft12, sx, sy+160+48, colorWhite)
		sx += 48
		text.Draw(dst, "时柱", ft14, sx, sy-48, colorWhite)
		text.Draw(dst, bz.GetTimeShiShenGan(), ft14, sx, sy-32, colorWhite)
		DrawFlow(dst, sx, sy, soul, p.Time)
		text.Draw(dst, strings.Join(p.ShenShaT, "\n"), ft12, sx, sy+160+48, colorWhite)
		sx += 48
		vector.StrokeLine(dst, float32(sx-3), float32(sy-28), float32(sx-3), float32(sy+148), 1, colorWhite, true)
		if p.YunIdx == 0 {
			text.Draw(dst, "小运", ft14, sx, sy-48, colorWhite)
		} else {
			text.Draw(dst, "大运", ft14, sx, sy-48, colorWhite)
		}
		if p.FYun != nil {
			text.Draw(dst, LunarUtil.SHI_SHEN[soul+p.FYun.Gan], ft14, sx, sy-32, colorWhite)
			DrawFlow(dst, sx, sy, soul, p.FYun)
			text.Draw(dst, strings.Join(p.ShenShaYY, "\n"), ft12, sx, sy+160+48, colorWhite)
		}
		sx += 48
		vector.StrokeLine(dst, float32(sx-3), float32(sy-28), float32(sx-3), float32(sy+148), 1, colorWhite, true)
		text.Draw(dst, "流年", ft14, sx, sy-48, colorWhite)
		text.Draw(dst, LunarUtil.SHI_SHEN[soul+g.Flow.Year.Gan], ft14, sx, sy-32, colorWhite)
		DrawFlow(dst, sx, sy, soul, g.Flow.Year)
		text.Draw(dst, strings.Join(p.ShenShaFY, "\n"), ft12, sx, sy+160+48, colorWhite)
		sx += 48
		text.Draw(dst, "流月", ft14, sx, sy-48, colorWhite)
		text.Draw(dst, LunarUtil.SHI_SHEN[soul+g.Flow.Month.Gan], ft14, sx, sy-32, colorWhite)
		DrawFlow(dst, sx, sy, soul, g.Flow.Month)
		text.Draw(dst, strings.Join(p.ShenShaFM, "\n"), ft12, sx, sy+160+48, colorWhite)
		if !g.BodyShow {
			sx += 48
			text.Draw(dst, "流日", ft14, sx, sy-48, colorWhite)
			text.Draw(dst, LunarUtil.SHI_SHEN[soul+g.Flow.Day.Gan], ft14, sx, sy-32, colorWhite)
			DrawFlow(dst, sx, sy, soul, g.Flow.Day)
			text.Draw(dst, strings.Join(p.ShenShaFD, "\n"), ft12, sx, sy+160+48, colorWhite)
			sx += 48
			text.Draw(dst, "流时", ft14, sx, sy-48, colorWhite)
			text.Draw(dst, LunarUtil.SHI_SHEN[soul+g.Flow.Time.Gan], ft14, sx, sy-32, colorWhite)
			DrawFlow(dst, sx, sy, soul, g.Flow.Time)
			text.Draw(dst, strings.Join(p.ShenShaFT, "\n"), ft12, sx, sy+160+48, colorWhite)
		}
	}
	//竖象 身体全息
	//年头颈/月胸腹/日腹股/时腿足
	//年发额/月目/日鼻/时口
	//干右/支左 年月左/时右 对称同位比 克泄重克弱小
	//男甲<申 右小左大 右近视度高
	//女乙>未 右小左大 右单左双
	//女辛<巳 右扁左圆
	//七杀疤痕,伤宫胎记,袅神为痣,劫财纹身 喜用则美,忌神则丑
	//通常四柱中：
	//1、天干有官杀之人或有木克土之人，上半身容易留下疤痕；
	//2、地支有官杀或有木克土之人，下半身容易留下疤痕；
	//具体说来：
	//一、年干：
	//1、时干克年干，疤痕在身体的右侧；
	//2、日干克年干，疤痕在身体中间的偏右侧部位；
	//3、月干克年干，疤痕在身体的左侧部位；
	//二、月干：
	//1、年干克月干，疤痕在身体的左侧明显部位；
	//2、日干克月干，疤痕在身体的右侧明显部位；
	//3、时干克月干，疤痕在身体的左侧明显部位；
	//三、日干：
	//1、时干克日干，疤痕在身体的中间偏右侧部位；
	//2、年干克日干，疤痕在身体的左侧部位；
	//3、月干克日干，疤痕在身体中间的偏左侧部位；
	//四、时干：
	//1、月干克时干，疤痕在身体中间的偏左侧部位；
	//2、年干克时干，疤痕在身体的左侧明显部位；
	//3、日干克时干，疤痕在身体的右侧明显部位。
	if g.BodyShow {
		sx, sy := float32(cx+408), float32(cy)
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
			text.Draw(dst, p.Year.Gan, ft14, mx, sy-16, ColorGanZhi(p.Year.Gan))
			text.Draw(dst, p.Year.Zhi, ft14, mx+16, sy-16, ColorGanZhi(p.Year.Zhi))
			text.Draw(dst, p.Month.Gan, ft14, mx, sy, ColorGanZhi(p.Month.Gan))
			text.Draw(dst, p.Month.Zhi, ft14, mx+16, sy, ColorGanZhi(p.Month.Zhi))
			text.Draw(dst, p.Day.Gan, ft14, mx, sy+16, ColorGanZhi(p.Day.Gan))
			text.Draw(dst, p.Day.Zhi, ft14, mx+16, sy+16, ColorGanZhi(p.Day.Zhi))
			text.Draw(dst, p.Time.Gan, ft14, mx, sy+32, ColorGanZhi(p.Time.Gan))
			text.Draw(dst, p.Time.Zhi, ft14, mx+16, sy+32, ColorGanZhi(p.Time.Zhi))
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
	//横象 年长 月同 日夫妻 时子孙 干动支静 干为军支为营 干为官支为民
	if !g.OverviewShow {
		sx, sy := float32(cx+3), float32(cy+200)
		g.DrawCharHP(dst, sx+96, sy, g.Flow.Year, "流年")
		g.DrawCharHP(dst, sx+96*2, sy, g.Flow.Month, "流月")
		g.DrawCharHP(dst, sx+96*3, sy, g.Flow.Day, "流日")
		g.DrawCharHP(dst, sx+96*4, sy, g.Flow.Time, "流时")
	}
	{ //本命
		sx, sy := float32(cx+3), float32(cy+410)
		if p.FYun != nil {
			g.DrawCharHP(dst, sx-3, sy, p.FYun, "大运")
		}
		g.DrawCharHP(dst, sx+96, sy, p.Year, "年柱")
		g.DrawCharHP(dst, sx+96*2, sy, p.Month, "月柱")
		g.DrawCharHP(dst, sx+96*3, sy, p.Day, GenderSoul[p.Gender])
		g.DrawCharHP(dst, sx+96*4, sy, p.Time, "时柱")
		if !g.OverviewShow {
			fx, fy := sx+96-28, sy-45
			yx, yy := sx+96+28*2, sy-60
			mx, my := sx+96*2+28*2, sy-30
			dx, dy := sx+96*3+28*2, sy-30
			tx, ty := sx+96*4+28+18, sy-60
			ffy := yy - 45
			vector.StrokeLine(dst, fx, fy, mx, my, .5, colorGray, true)  //运月线
			vector.StrokeLine(dst, yx, yy, mx, my, .5, colorGray, true)  //年月线
			vector.StrokeLine(dst, mx, my, dx, dy, .5, colorGray, true)  //月日线
			vector.StrokeLine(dst, dx, dy, tx, ty, .5, colorGray, true)  //日时线
			vector.StrokeLine(dst, yx, yy, dx, dy, .5, colorGray, true)  //年日线
			vector.StrokeLine(dst, mx, my, tx, ty, .5, colorGray, true)  //月时线
			vector.StrokeLine(dst, yx, yy, tx, ty, .5, colorGray, true)  //年时线
			vector.StrokeLine(dst, yx, yy, yx, ffy, .5, colorGray, true) //流年线
			vector.StrokeLine(dst, mx, my, mx, ffy, .5, colorGray, true) //流月线
			vector.StrokeLine(dst, dx, dy, dx, ffy, .5, colorGray, true) //流日线
			vector.StrokeLine(dst, tx, ty, tx, ffy, .5, colorGray, true) //流时线
		}
		if p.Mate != nil { //配偶
			sy += 102
			g.DrawCharHP(dst, sx+96, sy, p.Mate.Year, "年柱")
			g.DrawCharHP(dst, sx+96*2, sy, p.Mate.Month, "月柱")
			g.DrawCharHP(dst, sx+96*3, sy, p.Mate.Day, GenderSoul[p.Mate.Gender]+" 配偶")
			g.DrawCharHP(dst, sx+96*4, sy, p.Mate.Time, "时柱")
		}
	}
	g.BaseUI.Draw(dst)
}

func (g *Char8Pan) DrawCharHP(dst *ebiten.Image, sx, sy float32, body *CharBody, title string) {
	ft14, _ := GetFontFace(14)
	ft28, _ := GetFontFace(28)
	vector.StrokeRect(dst, sx, sy, 96, 80, 1, colorWhite, true)
	text.Draw(dst, body.Gan, ft28, int(sx), int(sy), ColorGanZhi(body.Gan))
	text.Draw(dst, title, ft14, int(sx+28), int(sy-10), colorWhite)
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

	FlowEnergy bool //流气锁值
}

func NewCharBody(gan, zhi string, ganMax, zhiMax int, flow bool) *CharBody {
	cb := &CharBody{Gan: gan, Zhi: zhi,
		Body:   GetHideGan(zhi, 0),
		Legs:   GetHideGan(zhi, 1),
		Feet:   GetHideGan(zhi, 2),
		HPHead: ganMax, HPMHead: ganMax,
		HPMBody: zhiMax, HPMLegs: zhiMax, HPMFeet: zhiMax,
		FlowEnergy: flow,
	}
	cb.initZhiHP(zhiMax)
	return cb
}
func (b *CharBody) GetGanZhi() string {
	return b.Gan + b.Zhi
}
func (b *CharBody) initZhiHP(maxHp int) {
	if b.Feet != "" {
		b.HPBody = maxHp * HideGanVal[3][0] / 100
		b.HPLegs = maxHp * HideGanVal[3][1] / 100
		b.HPFeet = maxHp * HideGanVal[3][2] / 100
	} else if b.Legs != "" {
		b.HPBody = maxHp * HideGanVal[2][0] / 100
		b.HPLegs = maxHp * HideGanVal[2][1] / 100
		b.HPFeet = 0
	} else {
		b.HPBody = maxHp * HideGanVal[1][0] / 100
		b.HPLegs = 0
		b.HPFeet = 0
	}
}

func (b *CharBody) InteractiveSelf(speed int) {
	//本柱 支引干透
	//cs :=qimen.ChangSheng12[b.Gan+b.Body]
	ss := LunarUtil.SHI_SHEN[b.Gan+b.Body]
	Interactive[ss](&b.HPHead, &b.HPBody, b.HPMHead, b.HPMBody, speed)
	//本柱 支藏干化
	if b.Legs != "" {
		ss = LunarUtil.SHI_SHEN[b.Gan+b.Legs]
		Interactive[ss](&b.HPHead, &b.HPLegs, b.HPMHead, b.HPMLegs, speed)
		ss = LunarUtil.SHI_SHEN[b.Body+b.Legs]
		Interactive[ss](&b.HPBody, &b.HPLegs, b.HPMBody, b.HPMLegs, speed)
	}
	if b.Feet != "" {
		ss = LunarUtil.SHI_SHEN[b.Gan+b.Feet]
		Interactive[ss](&b.HPHead, &b.HPFeet, b.HPMHead, b.HPMFeet, speed)
		ss = LunarUtil.SHI_SHEN[b.Body+b.Feet]
		Interactive[ss](&b.HPBody, &b.HPFeet, b.HPMBody, b.HPMFeet, speed)
		ss = LunarUtil.SHI_SHEN[b.Feet+b.Legs]
		Interactive[ss](&b.HPFeet, &b.HPLegs, b.HPMFeet, b.HPMLegs, speed)
	}
}

type Body4 struct {
	Year  *CharBody //年柱
	Month *CharBody //月柱
	Day   *CharBody //日柱
	Time  *CharBody //时柱
}

func (p *Body4) resetHP() {
	p.Year.HPHead = HpGanYear
	p.Year.initZhiHP(HpZhiYear)
	p.Month.HPHead = HpGanMonth
	p.Month.initZhiHP(HpZhiMonth)
	p.Day.HPHead = HpGanDay
	p.Day.initZhiHP(HpZhiDay)
	p.Time.HPHead = HpGanTime
	p.Time.initZhiHP(HpZhiTime)
}

type Player struct {
	Gender int //性别0女1男
	Birth  *calendar.Lunar
	Body4            //四柱
	FYun   *CharBody //大运
	Mate   *Player   //配偶

	yun             *calendar.Yun     //运
	yuns            []*calendar.DaYun //大运集
	YunIdx0, YunIdx int               //当前大运小运索引
	Fates0          []string          //小运名
	Fates           []string          //大运名

	UpdateCount int

	ShenShaY  []string //神煞年
	ShenShaM  []string //神煞月
	ShenShaD  []string //神煞日
	ShenShaT  []string //神煞时
	ShenShaYY []string //神煞大运
	ShenShaFY []string //神煞流年
	ShenShaFM []string //神煞流月
	ShenShaFD []string //神煞流日
	ShenShaFT []string //神煞流时
}

func (p *Player) Reset(lunar *calendar.Lunar, gender int) {
	p.Birth = lunar
	p.Gender = gender
	bz := lunar.GetEightChar()
	zhiY, zhiM, zhiD, zhiT := bz.GetYearZhi(), bz.GetMonthZhi(), bz.GetDayZhi(), bz.GetTimeZhi()
	p.Year = NewCharBody(bz.GetYearGan(), zhiY, HpGanYear, HpZhiYear, false)
	p.Month = NewCharBody(bz.GetMonthGan(), zhiM, HpGanMonth, HpZhiMonth, false)
	p.Day = NewCharBody(bz.GetDayGan(), zhiD, HpGanDay, HpZhiDay, false)
	p.Time = NewCharBody(bz.GetTimeGan(), zhiT, HpGanTime, HpZhiTime, false)
	sss := qimen.CalcShenSha(bz)
	p.ShenShaY, p.ShenShaM, p.ShenShaD, p.ShenShaT = sss[0], sss[1], sss[2], sss[3]

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
	p.resetHP()
}

func CharBodyInteractive(a, b *CharBody, speed int) {
	if a == nil || b == nil {
		return
	}
	//干干
	//he := qimen.HeGan[a.Gan] == b.Gan     //合
	//if he {
	//}
	if a.FlowEnergy {
		aa := *a
		a = &aa
	}
	if b.FlowEnergy {
		bb := *b
		b = &bb
	}
	ss := LunarUtil.SHI_SHEN[a.Gan+b.Gan] //b为a的
	Interactive[ss](&a.HPHead, &b.HPHead, a.HPMHead, b.HPMHead, speed)
	//a干b支
	ss = LunarUtil.SHI_SHEN[a.Gan+b.Body]
	Interactive[ss](&a.HPHead, &b.HPBody, a.HPMHead, b.HPMBody, speed)
	if b.Legs != "" {
		ss = LunarUtil.SHI_SHEN[a.Gan+b.Legs]
		Interactive[ss](&a.HPHead, &b.HPLegs, a.HPMHead, b.HPMLegs, speed)
	}
	if b.Feet != "" {
		ss = LunarUtil.SHI_SHEN[a.Gan+b.Feet]
		Interactive[ss](&a.HPHead, &b.HPFeet, a.HPMHead, b.HPMFeet, speed)
	}
	a, b = b, a //交换
	ss = LunarUtil.SHI_SHEN[a.Gan+b.Body]
	Interactive[ss](&a.HPHead, &b.HPBody, a.HPMHead, b.HPMBody, speed)
	if b.Legs != "" {
		ss = LunarUtil.SHI_SHEN[a.Gan+b.Legs]
		Interactive[ss](&a.HPHead, &b.HPLegs, a.HPMHead, b.HPMLegs, speed)
	}
	if b.Feet != "" {
		ss = LunarUtil.SHI_SHEN[a.Gan+b.Feet]
		Interactive[ss](&a.HPHead, &b.HPFeet, a.HPMHead, b.HPMFeet, speed)
	}
	//支支 合冲刑破害
	a, b = b, a
	ss = LunarUtil.SHI_SHEN[a.Body+b.Body]
	Interactive[ss](&a.HPBody, &b.HPBody, a.HPMBody, b.HPMBody, speed)
	if a.Legs != "" && b.Legs != "" {
		ss = LunarUtil.SHI_SHEN[a.Legs+b.Legs]
		Interactive[ss](&a.HPLegs, &b.HPLegs, a.HPMLegs, b.HPMLegs, speed)
	}
	if a.Feet != "" && b.Feet != "" {
		ss = LunarUtil.SHI_SHEN[a.Feet+b.Feet]
		Interactive[ss](&a.HPFeet, &b.HPFeet, a.HPMFeet, b.HPMFeet, speed)
	}
}

var Interactive = map[string]func(va, vb *int, ma, mb int, speed int){
	"比肩": InteractiveBiJie,
	"劫财": InteractiveBiJie,
	"食神": InteractiveShi,
	"伤官": InteractiveShang,
	"正印": InteractiveYin,
	"偏印": InteractiveXiao,
	"正官": InteractiveGuan,
	"七杀": InteractiveSha,
	"正财": InteractiveZhengCai,
	"偏财": InteractivePianCai,
}

/*
金赖土生，土多金埋。土赖火生，火多土焦。
火赖木生，木多火炽。木赖水生，水多木漂。
水赖金生，金多水浊。
水空则流，木空则损，土空则陷，金空则响，火空则发。
旺木喜金，旺火喜水，旺土喜木，旺金喜火，旺水喜土。
木怕金旺，火怕水旺，土怕木旺，金怕火旺，水怕土旺。
水弱则爱金，金弱则爱土，土弱则爱火，火弱则爱木，木弱则爱水。
水衰不生木，木衰不生火，火衰不生土，土衰不生金，金衰不生水。
春土不克水，夏金不克木，季水不克火，秋木不克土，冬火不克金。
*/

// InteractiveBiJie 比劫
func InteractiveBiJie(va, vb *int, ma, mb int, speed int) {
	if *va < ma && *va+speed < *vb-speed {
		*va += speed
		*vb -= speed
	}
	if *vb < mb && *vb+speed < *va-speed {
		*vb += speed
		*va -= speed
	}
}

// InteractiveShi 食神 a生b 同阴阳 生力大.
// 水衰不生木，木衰不生火，火衰不生土，土衰不生金，金衰不生水。
func InteractiveShi(va, vb *int, ma, mb int, speed int) {
	if *va < ma/2 { //不旺不生
		return
	}
	if *vb >= mb { //子满不生
		return
	}
	if *vb > *va/2 { //衰不生
		return
	}
	if *vb+speed < *va-speed {
		if *vb+speed < (*va-speed)/2 { //快泄
			*va -= speed
			*vb += speed
		} else { //慢泄
			*va -= speed / 2
			*vb += speed / 2
		}
	}
}

// InteractiveShang 伤官 a生b 异阴阳 生力小
func InteractiveShang(va, vb *int, ma, mb int, speed int) {
	if *va < ma/2 { //不旺不生
		return
	}
	if *vb >= mb { //满不生
		return
	}
	if *vb > *va/2 { //衰不生
		return
	}
	if *vb+speed < *va-speed {
		if *vb+speed < (*va-speed)/2 { //快泄
			*va -= speed
			*vb += speed
		} else { //慢泄
			*va -= speed / 2
			*vb += speed / 2
		}
	}
}

// InteractiveYin 正印:b生a
func InteractiveYin(va, vb *int, ma, mb int, speed int) {
	InteractiveShang(vb, va, mb, ma, speed)
}

// InteractiveXiao 偏印:b生a
func InteractiveXiao(va, vb *int, ma, mb int, speed int) {
	InteractiveShi(vb, va, mb, ma, speed)
}

// InteractiveGuan 正官:a嫁b
func InteractiveGuan(va, vb *int, ma, mb int, speed int) {
	//if qimen.HeGan[a.Gan] == b.Gan { //合
	//}
	if *vb < mb && *vb+speed < *va-speed {
		*vb += speed
		*va -= speed
	}
}

// InteractiveSha 七杀:b夺a
func InteractiveSha(va, vb *int, ma, mb int, speed int) {
	//if *va > 1 { // && !a.FlowEnergy
	//	*va--
	//}
	if *vb < mb && *va > speed {
		*va -= speed
		*vb += speed // 2
	}
}

// InteractiveZhengCai 正财:a娶b
func InteractiveZhengCai(va, vb *int, ma, mb int, speed int) {
	InteractiveGuan(vb, va, mb, ma, speed)
}

// InteractivePianCai 偏财:a抢b
func InteractivePianCai(va, vb *int, ma, mb int, speed int) {
	InteractiveSha(vb, va, mb, ma, speed)
}

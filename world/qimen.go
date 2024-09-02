package world

import (
	"fmt"
	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"qimen/qimen"
	"qimen/ui"
	"qimen/util"
	"strings"
)

const (
	_GongWidth = 96 //宫宽
)

type QMShow struct {
	X, Y float32
}

func NewQiMenShow(centerX, centerY float32) *QMShow {
	return &QMShow{
		X: centerX, Y: centerY,
	}
}
func (q *QMShow) Update() {
}

func (q *QMShow) Draw(dst *ebiten.Image) {
	q.drawHead(dst)
	q.draw9Gong(dst)
	q.draw12Gong(dst)
}
func (q *QMShow) drawHead(dst *ebiten.Image) {
	pan := ThisGame.qmGame
	lunar := pan.Lunar
	pp := pan.ShowPan
	ft := ui.GetDefaultUIFont()
	var cYear string
	if lunar.GetYear() == 1 {
		cYear = "元年"
	} else if lunar.GetYear() <= 0 {
		cYear = fmt.Sprintf("公元前%d年", -lunar.GetYear()+1)
	} else {
		cYear = lunar.GetYearInChinese()
	}
	text.Draw(dst, fmt.Sprintf("  %s %s %s %s",
		cYear, lunar.GetMonthInChinese()+"月", lunar.GetDayInChinese(), lunar.GetEightChar().GetTimeZhi()+"时"),
		ft, 32, 48, colorWhite)
	text.Draw(dst, fmt.Sprintf("干支  %s %s %s %s",
		lunar.GetYearInGanZhiExact(), lunar.GetMonthInGanZhiExact(), lunar.GetDayInGanZhiExact(), lunar.GetTimeInGanZhi()),
		ft, 32, 64, colorLeader)
	text.Draw(dst, fmt.Sprintf("旬首  %s %s %s %s",
		lunar.GetYearXunExact(), lunar.GetMonthXunExact(), lunar.GetDayXunExact(), lunar.GetTimeXun()),
		ft, 32, 64+16, colorWhite)
	text.Draw(dst, fmt.Sprintf("空亡  %s %s %s %s",
		lunar.GetYearXunKongExact(), lunar.GetMonthXunKongExact(), lunar.GetDayXunKongExact(), lunar.GetTimeXunKong()),
		ft, 32, 64+32, colorGray)
	text.Draw(dst, pp.JuText, ft, 32, 96+16, colorWhite)
}
func (q *QMShow) draw9Gong(dst *ebiten.Image) {
	ft := ui.GetDefaultUIFont()
	qm := ThisGame.qmGame
	pp := qm.ShowPan
	//画九宫
	kongWang := LunarUtil.GetXunKong(pp.Xun)
	for i := 1; i <= 9; i++ {
		offX, offZ := gongOffset[i][0]*_GongWidth-_GongWidth/2, gongOffset[i][1]*_GongWidth-_GongWidth/2
		px, py := q.X-_GongWidth+float32(offX), q.Y-_GongWidth+float32(offZ)

		//vector.StrokeCircle(dst, px+_GongWidth/2, py+_GongWidth/2,
		//	float32(_GongWidth/2), 1, color.RGBA{0xff, 0x80, 0xff, 0xff}, true)
		//vector.DrawFilledRect(dst, px, py, _GongWidth-1, _GongWidth-1, color9Gong[i], true)
		vector.StrokeRect(dst, px, py, _GongWidth-1, _GongWidth-1, 1, color9Gong[i], true)

		g := pp.Gongs[i]
		var hosting = "  "
		if pp.RollHosting > 0 && i == pp.DutyStarPos {
			hosting = "禽"
		}
		var empty, horse = "  ", "  "
		if strings.Contains(kongWang, LunarUtil.ZHI[i]) {
			empty = "〇" //"空亡"
		}
		if qimen.ZhiGong9[qm.TimeHorse] == i {
			horse = "马"
		}
		door := g.Door + qimen.Door0
		if g.Door == "" {
			door = "    "
		}
		star := qimen.Star0 + g.Star
		if g.Star == "" {
			star = ""
		}
		txt := fmt.Sprintf("\n  %s%s  %s\n\n"+ //神盘
			" %s %s%s%s\n\n"+ //天盘
			" %s %s  %s\n", //人盘
			empty, g.God, horse,
			g.PathGan+g.HideGan, star, hosting, g.GuestGan,
			g.PathZhi, door, g.HostGan)
		text.Draw(dst, txt, ft, int(px), int(py), color.White)
	}
}
func (q *QMShow) draw12Gong(dst *ebiten.Image) {
	ft := ui.GetDefaultUIFont()
	pan := ThisGame.qmGame
	//画12宫
	//if uiQiMen.qmParams.YMDH != qimen.QMGameHour {
	//	return
	//}
	r1, r2 := float32(_GongWidth)*1.5+zhiPanWidth*1.5, float32(_GongWidth)*1.5+zhiPanWidth*2
	//空亡偏心环
	r0 := r1 - zhiPanWidth/8
	emptyClock := qimen.KongWangClock[pan.ShowPan.Xun]
	angle := float64(emptyClock-45) * 30 //+ float64(g.count)
	rad := angle * math.Pi / 180
	x0 := float64(q.X) + float64(zhiPanWidth/8)*math.Cos(rad)
	y0 := float64(q.Y) + float64(zhiPanWidth/8)*math.Sin(rad)
	vector.StrokeCircle(dst, float32(x0), float32(y0), r0, zhiPanWidth/2, colorPowerCircle, true)
	//建星地户盘
	vector.StrokeCircle(dst, q.X, q.Y, r1, zhiPanWidth/2, colorGroundGateCircle, true)
	//月将天门盘
	vector.StrokeCircle(dst, q.X, q.Y, r2, zhiPanWidth/2, colorSkyGateCircle, true)

	for i := 1; i <= 12; i++ {
		angleDegrees := float64(i+2) * 30 //+ float64(g.count)
		lx1, ly1 := util.CalRadiansPos(float64(q.X), float64(q.Y), float64(r1-zhiPanWidth/4), angleDegrees-15)
		lx2, ly2 := util.CalRadiansPos(float64(q.X), float64(q.Y), float64(r2+zhiPanWidth/4), angleDegrees-15)
		vector.StrokeLine(dst, float32(lx1), float32(ly1), float32(lx2), float32(ly2), 1, colorGongSplit, true)
		if pan.Big6 == nil {
			continue
		}
		gong12 := pan.Big6[i-1]
		jiangColor := colorJiang
		if gong12.IsJiang {
			jiangColor = colorLeader
		} else if qimen.SkyGate3[gong12.Jiang] {
			jiangColor = colorGate
		}
		x1, y1 := util.CalRadiansPos(float64(q.X), float64(q.Y), float64(r2), angleDegrees)
		text.Draw(dst, gong12.Jiang, ft, int(x1-14), int(y1+4), jiangColor)
		x12, y12 := util.CalRadiansPos(float64(q.X), float64(q.Y), float64(r2), angleDegrees+10)
		text.Draw(dst, gong12.JiangZhi, ft, int(x12-14), int(y12+4), jiangColor)

		jianColor := colorJian
		if gong12.IsJian {
			jianColor = colorLeader
		} else if qimen.GroundGate4[gong12.Jian] {
			jianColor = colorGate
		}
		x2, y2 := util.CalRadiansPos(float64(q.X), float64(q.Y), float64(r1), angleDegrees)
		text.Draw(dst, gong12.Jian, ft, int(x2-8), int(y2+4), jianColor)
		x22, y22 := util.CalRadiansPos(float64(q.X), float64(q.Y), float64(r1), angleDegrees+10)
		text.Draw(dst, gong12.JianZhi, ft, int(x22-8), int(y22+4), jianColor)

		if ThisGame.qmGame.TimeHorse == LunarUtil.ZHI[i] {
			x3, y3 := util.CalRadiansPos(float64(q.X), float64(q.Y), float64(r1), angleDegrees-10)
			text.Draw(dst, "驿马", ft, int(x3-8), int(y3+4), colorLeader)
		}
		if gong12.IsSkyHorse {
			x4, y4 := util.CalRadiansPos(float64(q.X), float64(q.Y), float64(r2), angleDegrees-10)
			text.Draw(dst, "天马", ft, int(x4-14), int(y4+4), colorLeader)
		}
	}
}

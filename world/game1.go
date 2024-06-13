package world

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"qimen/qimen"
	"qimen/ui"
)

const (
	_GongWidth = 96 //宫宽
	centerX    = 260
	centerY    = 460
)

var (
	color9Gong = []color.RGBA{
		{0x00, 0x00, 0x00, 0x00},
		{0x40, 0x40, 0xFF, 0x20}, //坎一
		{0x60, 0x60, 0x60, 0x20}, //坤二
		{0x00, 0x70, 0x00, 0x20}, //震三
		{0x10, 0xA0, 0x00, 0x20}, //巽四
		{0x80, 0x80, 0x00, 0x20}, //中五
		{0xA0, 0xA0, 0xA0, 0x20}, //乾六
		{0x80, 0x40, 0x00, 0x20}, //兑七
		{0x80, 0x80, 0x80, 0x20}, //艮八
		{0x80, 0x00, 0x80, 0x20}, //离九
	}
	colorGongSplit        = color.RGBA{0x00, 0x00, 0x00, 0xff}
	colorPowerCircle      = color.RGBA{0x60, 0x60, 0xFF, 0xFF}
	colorGroundGateCircle = color.RGBA{0x80, 0x80, 0x00, 0xff}
	colorSkyGateCircle    = color.RGBA{0x40, 0x40, 0xFF, 0x80}
	colorJiang            = color.RGBA{0x00, 0x00, 0x00, 0xff}
	colorJian             = colorJiang
	colorLeader           = color.RGBA{0xff, 0xff, 0x00, 0xff}
	colorGate             = color.RGBA{0x00, 0xff, 0x00, 0xff}
)

type game1 struct {
	count     int
	starLines [starsLineCount]StarLine
}

func (g *game1) Update() error {
	g.count++
	g.count %= 360
	ui.Update()

	//x, y := ebiten.CursorPosition()
	//x, y := centerX, centerY
	x, y := screenWidth/2, centerY
	for i := 0; i < starsLineCount; i++ {
		g.starLines[i].Update(float32(x*starsLineScale), float32(y*starsLineScale))
	}
	return nil
}

func (g *game1) Draw(screen *ebiten.Image) {
	for i := 0; i < starsLineCount; i++ {
		g.starLines[i].Draw(screen)
	}
	ui.Draw(screen)
	drawQM(screen)
	//drawStar(screen)
}

func (g *game1) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame1() *game1 {
	g := &game1{}
	for i := 0; i < starsLineCount; i++ {
		g.starLines[i].Init()
	}
	u := UIShowQiMen(screenWidth, screenHeight)
	u.hide9GongUI()
	u.noShow12Gong()
	return g
}

func calRadiansPos(cx, cy, r, angleDegrees float64) (x, y float64) {
	rad := angleDegrees * math.Pi / 180
	x = cx + r*math.Cos(rad)
	y = cy + r*math.Sin(rad)
	return
}

func drawQM(screen *ebiten.Image) {
	// TODO 大量计算移到Update()中
	draw9Gong(screen)
	draw12Gong(screen)
}
func draw9Gong(screen *ebiten.Image) {
	ft := ui.GetDefaultUIFont()
	//画九宫
	for i := 1; i <= 9; i++ {
		offX, offZ := gongOffset[i][0]*_GongWidth-_GongWidth/2, gongOffset[i][1]*_GongWidth-_GongWidth/2
		px, py := centerX-_GongWidth+float32(offX), centerY-_GongWidth+float32(offZ)

		//vector.StrokeCircle(screen, px+_GongWidth/2, py+_GongWidth/2,
		//	float32(_GongWidth/2), 1, color.RGBA{0xff, 0x80, 0xff, 0xff}, true)
		//vector.DrawFilledRect(screen, px, py, _GongWidth-1, _GongWidth-1, color9Gong[i], true)
		vector.StrokeRect(screen, px, py, _GongWidth-1, _GongWidth-1, 1, color9Gong[i], true)

		pp := uiQiMen.pan.ShowPan
		kongWang := qimen.KongWang[pp.Xun]
		g := pp.Gongs[i]
		var hosting = "  "
		if pp.RollHosting > 0 && i == pp.DutyStarPos {
			hosting = "禽"
		}
		var empty, horse = "  ", "  "
		for _, zhi := range kongWang {
			if qimen.ZhiGong9[zhi] == i {
				empty = "〇" //"空亡"
				break
			}
		}
		if qimen.ZhiGong9[pp.Horse] == i {
			horse = "马"
		}
		door := g.Door + qimen.Door0
		if g.Door == "" {
			door = ""
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
		text.Draw(screen, txt, ft, int(px), int(py), color.White)
	}
}
func draw12Gong(screen *ebiten.Image) {
	ft := ui.GetDefaultUIFont()
	//画12宫
	if uiQiMen.qmParams.YMDH == qimen.QMGameHour {
		r1, r2 := float32(_GongWidth)*1.5+zhiPanWidth*1.5, float32(_GongWidth)*1.5+zhiPanWidth*2
		//空亡偏心环
		r0 := r1 - zhiPanWidth/8
		emptyClock := qimen.KongWangClock[uiQiMen.pan.ShowPan.Xun]
		angle := float64(emptyClock-45) * 30 //+ float64(g.count)
		rad := angle * math.Pi / 180
		x0 := float64(centerX) + float64(zhiPanWidth/8)*math.Cos(rad)
		y0 := float64(centerY) + float64(zhiPanWidth/8)*math.Sin(rad)
		vector.StrokeCircle(screen, float32(x0), float32(y0), r0, zhiPanWidth/2, colorPowerCircle, true)
		//建星地户盘
		vector.StrokeCircle(screen, centerX, centerY, r1, zhiPanWidth/2, colorGroundGateCircle, true)
		//月将天门盘
		vector.StrokeCircle(screen, centerX, centerY, r2, zhiPanWidth/2, colorSkyGateCircle, true)

		for i := 1; i <= 12; i++ {
			angleDegrees := float64(i-45) * 30 //+ float64(g.count)

			lx1, ly1 := calRadiansPos(float64(centerX), float64(centerY), float64(r1-zhiPanWidth/4), angleDegrees-15)
			lx2, ly2 := calRadiansPos(float64(centerX), float64(centerY), float64(r2+zhiPanWidth/4), angleDegrees-15)
			vector.StrokeLine(screen, float32(lx1), float32(ly1), float32(lx2), float32(ly2), 1, colorGongSplit, true)

			gong12 := uiQiMen.gong12[i]
			jiangColor := colorJiang
			if gong12.IsJiang {
				jiangColor = colorLeader
			} else if qimen.SkyGate3[gong12.Jiang] {
				jiangColor = colorGate
			}
			x1, y1 := calRadiansPos(float64(centerX), float64(centerY), float64(r2), angleDegrees)
			text.Draw(screen, gong12.Jiang, ft, int(x1-14), int(y1+4), jiangColor)
			x12, y12 := calRadiansPos(float64(centerX), float64(centerY), float64(r2), angleDegrees+10)
			text.Draw(screen, gong12.JiangZhi, ft, int(x12-14), int(y12+4), jiangColor)

			jianColor := colorJian
			if gong12.IsJian {
				jianColor = colorLeader
			} else if qimen.GroundGate4[gong12.Jian] {
				jianColor = colorGate
			}
			x2, y2 := calRadiansPos(float64(centerX), float64(centerY), float64(r1), angleDegrees)
			text.Draw(screen, gong12.Jian, ft, int(x2-8), int(y2+4), jianColor)
			x22, y22 := calRadiansPos(float64(centerX), float64(centerY), float64(r1), angleDegrees+10)
			text.Draw(screen, gong12.JianZhi, ft, int(x22-8), int(y22+4), jianColor)

			if gong12.IsHorse {
				x3, y3 := calRadiansPos(float64(centerX), float64(centerY), float64(r1), angleDegrees-10)
				text.Draw(screen, "马", ft, int(x3-8), int(y3+4), colorLeader)
			}
		}
	}
}

func drawStar(screen *ebiten.Image) {
	//ft := ui.GetDefaultUIFont()
	const (
		centerX = 770
		centerY = 450
		r1      = 240
		r2      = 216
		r3      = 192
	)
	vector.StrokeCircle(screen, centerX, centerY, r1, zhiPanWidth/2, colorSkyGateCircle, true)
	vector.StrokeCircle(screen, centerX, centerY, r2, zhiPanWidth/2, colorGroundGateCircle, true)
	vector.StrokeCircle(screen, centerX, centerY, r3, zhiPanWidth/2, colorPowerCircle, true)

	vector.StrokeLine(screen, centerX-r1, centerY, centerX+r1, centerY, 1, colorLeader, true)
	vector.StrokeLine(screen, centerX, centerY-r1, centerX, centerY+r1, 1, colorLeader, true)

	for i := 1; i <= 12; i++ {
		angleDegrees := float64(i-45) * 30 //+ float64(g.count)

		lx1, ly1 := calRadiansPos(float64(centerX), float64(centerY), float64(r1-zhiPanWidth/4), angleDegrees-15)
		lx2, ly2 := calRadiansPos(float64(centerX), float64(centerY), float64(r2+zhiPanWidth/4), angleDegrees-15)
		vector.StrokeLine(screen, float32(lx1), float32(ly1), float32(lx2), float32(ly2), 1, colorLeader, true)

	}
}

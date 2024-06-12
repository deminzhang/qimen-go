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

var color9Gong = []color.RGBA{
	{0x00, 0x00, 0x00, 0x00},
	{0x00, 0x00, 0x80, 0x20}, //坎一
	{0x60, 0x60, 0x60, 0x20}, //坤二
	{0x00, 0x80, 0x00, 0x20}, //震三
	{0x00, 0x90, 0x00, 0x20}, //巽四
	{0x80, 0x80, 0x00, 0x20}, //中五
	{0xA0, 0xA0, 0xA0, 0x20}, //乾六
	{0x80, 0x00, 0x00, 0x20}, //兑七
	{0x80, 0x80, 0x80, 0x20}, //艮八
	{0x80, 0x00, 0x80, 0x20}, //离九
}

type game1 struct {
	count int
}

func (g *game1) Update() error {
	g.count++
	g.count %= 360
	ui.Update()
	return nil
}

func (g *game1) Draw(screen *ebiten.Image) {
	ui.Draw(screen)
	ft := ui.GetDefaultUIFont()

	vector.StrokeRect(screen, float32(uiQiMen.opTypeRoll.X), float32(uiQiMen.opTypeRoll.Y),
		float32(uiQiMen.opTypeAmaze.X+uiQiMen.opTypeAmaze.ImageRect.Dx()-uiQiMen.opTypeRoll.X), 16, 1,
		color.RGBA{0x80, 0x00, 0x80, 0x20}, true)

	//盘中心
	box := uiQiMen.textGong[5]
	cx, cy := float32(box.Rect.Min.X+gongWidth/2), float32(box.Rect.Min.Y+gongWidth/2)
	//画九宫
	for i := 1; i <= 9; i++ {
		box := uiQiMen.textGong[i]
		vector.DrawFilledRect(screen, float32(box.Rect.Min.X), float32(box.Rect.Min.Y),
			gongWidth-1, gongWidth-1,
			color9Gong[i], true)
		vector.StrokeRect(screen, float32(box.Rect.Min.X), float32(box.Rect.Min.Y),
			gongWidth-1, gongWidth-1, 1,
			color9Gong[i], true)

		//vector.StrokeCircle(screen, float32(box.Rect.Min.X+(box.Rect.Max.X-box.Rect.Min.X)/2),
		//	float32(box.Rect.Min.Y+(box.Rect.Max.Y-box.Rect.Min.Y)/2),
		//	float32(box.Rect.Max.X-box.Rect.Min.X)/2, 1, color.RGBA{0xff, 0x80, 0xff, 0xff}, true)

		offX, offZ := gongOffset[i][0]*gongWidth, gongOffset[i][1]*gongWidth
		px, py := cx-gongWidth+float32(offX), cy-gongWidth+float32(offZ)
		pp := uiQiMen.pan.HourPan.Gongs[i]
		text.Draw(screen, pp.Star, ft, int(px), int(py), color.RGBA{0x00, 0xff, 0x00, 0xff})
	}
	//画12宫
	if uiQiMen.qmParams.YMDH == qimen.QMGameHour {
		r1, r2 := float32(gongWidth)*1.5+zhiPanWidth*1.5, float32(gongWidth)*1.5+zhiPanWidth*2
		//空亡盘 偏心环
		r0 := r1 - zhiPanWidth/8
		emptyClock := qimen.KongWangClock[uiQiMen.pan.HourPan.Xun]
		angle := float64(emptyClock-45) * 30 //+ float64(g.count)
		rad := angle * math.Pi / 180
		x0 := float64(cx) + float64(zhiPanWidth/8)*math.Cos(rad)
		y0 := float64(cy) + float64(zhiPanWidth/8)*math.Sin(rad)
		vector.StrokeCircle(screen, float32(x0), float32(y0), r0, zhiPanWidth/2, color.RGBA{0x60, 0x60, 0xFF, 0xFF}, true)
		//建星地户盘
		vector.StrokeCircle(screen, cx, cy, r1, zhiPanWidth/2, color.RGBA{0x80, 0x80, 0x00, 0xff}, true)
		//月将天门盘
		vector.StrokeCircle(screen, cx, cy, r2, zhiPanWidth/2, color.RGBA{0x40, 0x40, 0xFF, 0x80}, true)

		for i := 1; i <= 12; i++ {
			gong12 := uiQiMen.gong12[i]
			angleDegrees := float64(i-45) * 30 //+ float64(g.count)
			x1, y1 := calRadiansPos(float64(cx), float64(cy), float64(r1), angleDegrees)
			x2, y2 := calRadiansPos(float64(cx), float64(cy), float64(r2), angleDegrees)

			xx1, yx1 := calRadiansPos(float64(cx), float64(cy), float64(r1-zhiPanWidth/4), angleDegrees-15)
			xx2, yx2 := calRadiansPos(float64(cx), float64(cy), float64(r2+zhiPanWidth/4), angleDegrees-15)
			vector.StrokeLine(screen, float32(xx1), float32(yx1), float32(xx2), float32(yx2), 1,
				color.RGBA{0x00, 0x00, 0x00, 0xff}, true)

			jiang := gong12.Jiang
			jian := gong12.Jian
			if g.count/6%2 == 0 {
				jiang = " " + gong12.JiangZhi
				//jian = gong12.JianZhi
			}
			if qimen.SkyDoor3[jiang] {
				jiang = fmt.Sprintf("[%s]", jiang)
			}
			if qimen.GroundDoor4[gong12.Jian] {
				jian = fmt.Sprintf("[%s]", jian)
			}
			if gong12.IsJiang {
				text.Draw(screen, jiang, ft, int(x2)-12, int(y2)+4, color.RGBA{0xff, 0xff, 0x00, 0xff})
			} else {
				text.Draw(screen, jiang, ft, int(x2)-12, int(y2)+4, color.RGBA{0x00, 0x00, 0x00, 0xff})
			}
			if gong12.IsJian {
				text.Draw(screen, jian, ft, int(x1)-12, int(y1)+4, color.RGBA{0xff, 0xff, 0x00, 0xff})
			} else {
				text.Draw(screen, jian, ft, int(x1)-12, int(y1)+4, color.RGBA{0x00, 0x00, 0x00, 0xff})
			}
			if gong12.IsHorse {

			}
		}
	}
}

func (g *game1) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func NewGame1() *game1 {
	g := &game1{}
	u := UIShowQiMen(ScreenWidth, ScreenHeight)
	u.hide9GongBackGround(color.RGBA{0xff, 0xff, 0x00, 0xff})
	u.noShow12Gong()
	return g
}

func calRadiansPos(cx, cy, r, angleDegrees float64) (x, y float64) {
	angleRadians := angleDegrees * math.Pi / 180
	x = cx + r*math.Cos(angleRadians)
	y = cy + r*math.Sin(angleRadians)
	return
}

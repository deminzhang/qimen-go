package world

import (
	"fmt"
	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/6tail/lunar-go/calendar"
	"github.com/deminzhang/qimen-go/graphic"
	"github.com/deminzhang/qimen-go/qimen"
	"github.com/deminzhang/qimen-go/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"math/rand"
	"slices"
)

const (
	_SkyCircleWidth = 52                    //支环宽
	_Gong9Width     = 128                   //九宫宽
	_TaiJiSize      = 540                   //太极图大小
	_BaGuaSize      = _Gong9Width - 64 + 16 //八卦大小
	_StarSize       = 16                    //星体大小
	_JianR          = 305                   //月将盘半径
	_JiangR         = 330                   //建星盘半径
	_XiuR           = 355                   //星宿盘半径
	_CnstR          = 380                   //星座盘半径
)

type QMShow struct {
	X, Y  float32
	count int

	TaiJi    *ebiten.Image
	Sun      *ebiten.Image
	Moon     *ebiten.Image
	BaGua    map[int]*ebiten.Image
	CampM    *ebiten.Image
	Camp     *ebiten.Image
	Army     *ebiten.Image
	ArmyA    *ebiten.Image
	DutyFlag *ebiten.Image
}

func NewQiMenShow(centerX, centerY float32) *QMShow {
	bg := make(map[int]*ebiten.Image, 9)
	for i := 1; i <= 9; i++ {
		bg[i] = graphic.NewBaGuaImage(qimen.Diagrams9(i), _BaGuaSize)
	}
	return &QMShow{
		X: centerX, Y: centerY,
		TaiJi:    graphic.NewTaiJiImage(_TaiJiSize),
		Sun:      graphic.NewSunImage(_StarSize),
		Moon:     graphic.NewMoonImage(_StarSize),
		BaGua:    bg,
		CampM:    graphic.NewCampImage(64),
		Camp:     graphic.NewCampImage(32),
		Army:     graphic.NewArmyImage("庚", 32, 0),
		ArmyA:    graphic.NewArmyImage("兵", 32, 1),
		DutyFlag: graphic.NewFlagImage(16),
	}
}
func (q *QMShow) Update() {
	q.count++
	q.count %= 360
	q.updateBattle()
}
func (q *QMShow) updateBattle() {

}

func (q *QMShow) Draw(dst *ebiten.Image) {
	//q.drawTaiJi(dst)
	q.drawHead(dst)
	q.draw12Gong(dst)
	q.draw9Gong(dst)
	if Debug {
		//q.drawBattle(dst)
	}
}
func (q *QMShow) drawHead(dst *ebiten.Image) {
	pan := ThisGame.qmGame
	lunar := pan.Lunar
	pp := pan.ShowPan
	ft, _ := GetFontFace(14)
	px := 16
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
		ft, px, 48, colorWhite)
	text.Draw(dst, fmt.Sprintf("干支  %s %s %s %s",
		lunar.GetYearInGanZhiExact(), lunar.GetMonthInGanZhiExact(), lunar.GetDayInGanZhiExact(), lunar.GetTimeInGanZhi()),
		ft, px, 64, colorLeader)
	text.Draw(dst, fmt.Sprintf("旬首  %s %s %s %s",
		lunar.GetYearXunExact(), lunar.GetMonthXunExact(), lunar.GetDayXunExact(), lunar.GetTimeXun()),
		ft, px, 64+16, colorWhite)
	text.Draw(dst, "空亡", ft, px, 96, colorWhite)
	text.Draw(dst, lunar.GetYearXunKongExact(), ft, px+44, 96, colorGray)
	text.Draw(dst, lunar.GetMonthXunKongExact(), ft, px+74, 96, colorGray)
	text.Draw(dst, lunar.GetDayXunKongExact(), ft, px+112, 96, colorGray)
	text.Draw(dst, lunar.GetTimeXunKong(), ft, px+146, 96, colorRed)
	text.Draw(dst, fmt.Sprintf("%s%s", pp.JieQi, pp.JieQiDate), ft, px, 96+16, colorWhite)
	text.Draw(dst, fmt.Sprintf("%s%s", pp.JieQiNext, pp.JieQiDateNext), ft, px, 96+16*2, colorWhite)
	text.Draw(dst, pp.JuText, ft, px, 96+16*3, colorWhite)
	text.Draw(dst, pp.DutyText, ft, px, 96+16*4, colorWhite)
	text.Draw(dst, pp.YueJiang, ft, px, 96+16*5, colorWhite)

	text.Draw(dst, fmt.Sprintf("月相 %s", lunar.GetYueXiang()), ft, px, 96+16*6, colorWhite)
	text.Draw(dst, fmt.Sprintf("日值 %s%s%s", lunar.GetXiu(), lunar.GetZheng(), lunar.GetAnimal()), ft, px, 96+16*7, colorWhite)
	text.Draw(dst, fmt.Sprintf("岁大将军 %s", qimen.BigJiang[lunar.GetYearZhiExact()]), ft, px, 96+16*8, colorWhite)
	text.Draw(dst, fmt.Sprintf("月建大将军 %s", qimen.BigJiang[lunar.GetMonthZhiExact()]), ft, px, 96+16*9, colorWhite)

	text.Draw(dst, "符使马时", ft, px, 96+16*10, colorDuty)
	text.Draw(dst, "击刑", ft, px, 96+16*11, colorJiXing)
	text.Draw(dst, "门迫", ft, px, 96+16*12, colorMengPo)
	text.Draw(dst, "入墓", ft, px, 96+16*13, colorTomb)

}
func (q *QMShow) drawTaiJi(dst *ebiten.Image) {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(-_TaiJiSize/2, -_TaiJiSize/2)
	op.GeoM.Rotate(math.Pi * float64(q.count) / 360 * 2)
	op.GeoM.Translate(float64(q.X), float64(q.Y))
	op.ColorScale.ScaleWithColor(colorGray5)
	dst.DrawImage(q.TaiJi, &op)
}
func (q *QMShow) draw9Gong(dst *ebiten.Image) {
	qm := ThisGame.qmGame
	pp := qm.ShowPan
	//画九宫
	for i := 1; i <= 9; i++ {
		offX, offZ := gongOffset[i][0]*_Gong9Width-_Gong9Width/2, gongOffset[i][1]*_Gong9Width-_Gong9Width/2
		px, py := q.X-_Gong9Width+float32(offX), q.Y-_Gong9Width+float32(offZ)
		g := pp.Gongs[i]
		op := &ebiten.DrawImageOptions{}
		bd := float32(_Gong9Width-_BaGuaSize) / 2
		if i == 5 {
			op.GeoM.Translate(-_BaGuaSize/2, -_BaGuaSize/2)
			op.GeoM.Rotate(math.Pi * float64(q.count) / 360 * 2)
			op.GeoM.Translate(float64(px+_BaGuaSize/2+bd), float64(py+_BaGuaSize/2+bd))
			op.ColorScale.ScaleWithColor(colorGray5)
		} else {
			op.GeoM.Translate(float64(px+bd), float64(py+bd))
			op.ColorScale.ScaleWithColor(colorGray2)
		}
		dst.DrawImage(q.BaGua[g.Idx], op)

		q.drawGong(dst, px, py, &g)
	}
}
func (q *QMShow) drawGong(dst *ebiten.Image, x, y float32, g *qimen.QMPalace) {
	vector.StrokeRect(dst, x, y, _Gong9Width-1, _Gong9Width-1, 1, color9Gong[g.Idx], true)
	ft, _ := GetFontFace(14)
	qm := ThisGame.qmGame
	pp := qm.ShowPan
	var hosting = "  "
	if pp.RollHosting > 0 && g.Idx == pp.DutyStarPos {
		hosting = "禽"
	}
	var empty, horse = "  ", "  "
	for _, z := range []rune(pp.KongWang) {
		if qimen.ZhiGong9[string(z)] == g.Idx {
			empty = "〇" //"空亡"
		}
	}
	if qimen.ZhiGong9[pp.Horse] == g.Idx {
		//horse = "马"
		horse = "馬"
	}
	star := qimen.Star0 + g.Star
	if g.Star == "" {
		star = ""
	}
	door := g.Door + qimen.Door0
	if g.Door == "" {
		door = "    "
	}

	y += 35                                                 //神盘
	text.Draw(dst, empty, ft, int(x+8), int(y), colorWhite) //空亡
	text.Draw(dst, g.God, ft, int(x+24), int(y),
		util.If(g.God == qimen.QMGod8(1), colorDuty, colorWhite)) //神盘
	text.Draw(dst, horse, ft, int(x+8+64)+rand.Intn(3), int(y)+rand.Intn(2), colorDuty) //驿马
	y += 32                                                                             //天盘
	text.Draw(dst, g.PathGan, ft, int(x+8), int(y), colorGray)                          //流干
	text.Draw(dst, g.HideGan, ft, int(x+8), int(y), colorGray)                          //隐干
	text.Draw(dst, star, ft, int(x+24), int(y),
		util.If(g.Star == pp.DutyStar, colorDuty, colorWhite)) //星
	if g.Star == pp.DutyStar {
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x+24+26), float64(y-12))
		op.ColorScale.ScaleWithColor(colorLeader)
		dst.DrawImage(q.DutyFlag, &op)
	}
	text.Draw(dst, hosting, ft, int(x+48), int(y), colorDuty) //寄禽
	//text.Draw(dst, "", ft, int(x+8+32), int(y), colorGray)           //中寄
	guestGanTomb := qimen.ZhiGong9[qimen.QMTomb[g.GuestGan]] == g.Idx
	jiXing := g.God == "值符" && qimen.ZhiGong9[qimen.QM6YiJiXing[pp.Xun]] == g.Idx
	text.Draw(dst, g.GuestGan, ft, int(x+64+24), int(y), util.If(jiXing, colorJiXing,
		util.If(guestGanTomb, colorTomb, colorWhite))) //天盘干
	if g.GuestGan == pp.Gan {
		vector.StrokeRect(dst, x+64+24, y-12, 16, 16, 1, colorDuty, true)
	}
	y += 32                                                                                            //地盘
	text.Draw(dst, g.PathZhi, ft, int(x+8), int(y), colorGray)                                         //流支
	doorPo := qimen.WuxingKe[qimen.DoorWuxing[g.Door]] == qimen.DiagramsWuxing[qimen.Diagrams9(g.Idx)] //门迫
	text.Draw(dst, door, ft, int(x+24), int(y), util.If(doorPo, colorMengPo,
		util.If(g.Door == pp.DutyDoor, colorDuty, colorWhite))) //门
	if g.Door == pp.DutyDoor {
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x+24+26), float64(y-12))
		op.ColorScale.ScaleWithColor(colorDuty)
		dst.DrawImage(q.DutyFlag, &op)
	}
	hostGanTomb := qimen.ZhiGong9[qimen.QMTomb[g.HostGan]] == g.Idx
	text.Draw(dst, g.HostGan, ft, int(x+64+24), int(y), util.If(hostGanTomb, colorTomb, colorWhite)) //地盘干
	y += 20
	if g.Idx == pp.Duty { //地符
		text.Draw(dst, "值符", ft, int(x+24), int(y), colorGray)
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x+24+26), float64(y-12))
		op.ColorScale.ScaleWithColor(colorGray)
		dst.DrawImage(q.DutyFlag, &op)
	}

	//colorJiXing           = colorPurple  //奇门击刑
	//colorMengPo           = colorRed     //奇门门迫
	//colorXingMu           = colorBlue    //奇门刑+墓
}
func (q *QMShow) draw12Gong(dst *ebiten.Image) {
	ft, _ := GetFontFace(14)
	pan := ThisGame.qmGame
	h := pan.Solar.GetHour()
	m := pan.Solar.GetMinute()
	mD := m
	if h%2 == 0 {
		mD += 60
	}
	//建星地户盘
	vector.StrokeCircle(dst, q.X, q.Y, _JianR, _SkyCircleWidth/2, colorGroundGateCircle, true)
	//月将天门盘
	vector.StrokeCircle(dst, q.X, q.Y, _JiangR, _SkyCircleWidth/2, colorSkyGateCircle, true)
	//星宿盘
	vector.StrokeCircle(dst, q.X, q.Y, _XiuR, _SkyCircleWidth/2, colorPowerCircle, true)
	vector.StrokeCircle(dst, q.X, q.Y, _CnstR, _SkyCircleWidth/2, colorSkyGateCircle, true)
	//周天刻度
	for i := 1; i <= 360; i++ {
		lx1, ly1 := util.CalRadiansPos(q.X, q.Y, _XiuR-_SkyCircleWidth/5, float32(i))
		if i%10 == 0 {
			lx1, ly1 = util.CalRadiansPos(q.X, q.Y, _XiuR-float32(_SkyCircleWidth)/8, float32(i))
		} else if i%5 == 0 {
			lx1, ly1 = util.CalRadiansPos(q.X, q.Y, _XiuR-float32(_SkyCircleWidth)/12, float32(i))
		}
		lx2, ly2 := util.CalRadiansPos(q.X, q.Y, _XiuR-_SkyCircleWidth/4, float32(i))
		vector.StrokeLine(dst, lx1, ly1, lx2, ly2, .5, color.Black, true)
	}
	//日
	degreesS := -((float32(h) + float32(m)/60) * 15) //本地时区太阳角度 0~360 0时0度
	sy, sx := util.CalRadiansPos(q.Y, q.X, _JianR+12, degreesS)
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(sx-8), float64(sy-8))
	op.ColorScale.ScaleWithColor(colorSun)
	dst.DrawImage(q.Sun, &op)
	//月
	degreesM := -((float32(h)+float32(m)/60)*15 - float32(pan.Lunar.GetDay()-1)/float32(pan.LunarMonthDays)*360)
	my, mx := util.CalRadiansPos(q.Y, q.X, _JianR-12, degreesM)
	op = ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(mx-8), float64(my-8))
	op.ColorScale.ScaleWithColor(colorMoon)
	dst.DrawImage(q.Moon, &op)
	//TODO 行星..

	se := calendar.NewSolar(pan.Solar.GetYear(), 3, 24, 0, 0, 0)
	cf := se.GetLunar().GetPrevJieQi().GetSolar() //今年春分
	seN := calendar.NewSolar(pan.Solar.GetYear()+1, 3, 24, 0, 0, 0)
	cfN := seN.GetLunar().GetPrevJieQi().GetSolar() //下一年春分
	yearMin := cfN.SubtractMinute(cf)
	pr := cfN.SubtractMinute(pan.Solar)
	degreesR0 := float32(pr) / float32(yearMin) * 360 //春分角
	degreesR := degreesS + degreesR0                  //春分角+太阳角
	y, x := util.CalRadiansPos(q.Y, q.X, _XiuR+5, degreesR)
	text.Draw(dst, "*", ft, int(x-6), int(y), colorGreen) //春分角 白羊双鱼界
	y, x = util.CalRadiansPos(q.Y, q.X, _XiuR+5, degreesR+90)
	text.Draw(dst, "*", ft, int(x-6), int(y), colorRed) //夏至
	y, x = util.CalRadiansPos(q.Y, q.X, _XiuR+5, degreesR+180)
	text.Draw(dst, "*", ft, int(x-6), int(y), colorYellow) //秋分
	y, x = util.CalRadiansPos(q.Y, q.X, _XiuR+5, degreesR+270)
	text.Draw(dst, "*", ft, int(x-6), int(y), colorWhite) //冬至

	//画28星宿
	for i := 1; i <= 28; i++ {
		xi := qimen.Xiu28[i]
		degrees := degreesR + qimen.XiuAngle[xi]
		ly1, lx1 := util.CalRadiansPos(q.Y, q.X, _XiuR+_SkyCircleWidth/5, degrees)
		ly2, lx2 := util.CalRadiansPos(q.Y, q.X, _XiuR-_SkyCircleWidth/6, degrees)
		vector.StrokeLine(dst, lx1, ly1, lx2, ly2, .5, colorWhite, true) //星宿
		degreesNext := qimen.XiuAngle[qimen.Xiu28[util.If(i == 28, 1, i+1)]]
		if qimen.XiuAngle[xi] > degreesNext {
			degreesNext += 360
		}
		degreesMid := degrees + (degreesNext-qimen.XiuAngle[xi])/2
		y, x = util.CalRadiansPos(q.Y, q.X, _XiuR, degreesMid)
		text.Draw(dst, xi, ft, int(x-6), int(y+4), colorLeader) //星宿
	}
	//画12星座
	for i := 0; i < 12; i++ {
		degrees := degreesR + float32(360*i/12)
		ly1, lx1 := util.CalRadiansPos(q.Y, q.X, _CnstR-_SkyCircleWidth/4, degrees)
		ly2, lx2 := util.CalRadiansPos(q.Y, q.X, _CnstR+_SkyCircleWidth/5, degrees)
		vector.StrokeLine(dst, lx1, ly1, lx2, ly2, .5, colorGongSplit, true) //宫分割线
		y, x := util.CalRadiansPos(q.Y, q.X, _CnstR+3, degrees+15)
		text.Draw(dst, fmt.Sprintf("%s", qimen.ConstellationShort[i]), ft, int(x-6), int(y+6), colorWhite) //星座
	}
	//画12宫
	for i := 1; i <= 12; i++ {
		degrees := -float32(i-1) * 30
		degrees -= float32(30*mD/120) - 15 //分针
		ly1, lx1 := util.CalRadiansPos(q.Y, q.X, _JianR-_SkyCircleWidth/4, degrees+15)
		ly2, lx2 := util.CalRadiansPos(q.Y, q.X, _JianR+_SkyCircleWidth/5, degrees+15)
		vector.StrokeLine(dst, lx1, ly1, lx2, ly2, .5, colorGongSplit, true) //宫分割线

		//地支宫位
		var zhiGongTxt string
		if slices.Contains(qimen.KongWang[pan.ShowPan.Xun], LunarUtil.ZHI[i]) {
			zhiGongTxt = fmt.Sprintf("〇%s", LunarUtil.ZHI[i])
			//} else if slices.Contains(qimen.KongWang[pan.ShowPan.Xun], LunarUtil.ZHI[qimen.Idx12[i+6]]) {
			//	zhiGongTxt = fmt.Sprintf("%s虚", LunarUtil.ZHI[i])
		} else {
			zhiGongTxt = LunarUtil.ZHI[i]
		}
		y, x = util.CalRadiansPos(q.Y, q.X, _JianR-_SkyCircleWidth/2, degrees)
		text.Draw(dst, zhiGongTxt, ft, int(x-8), int(y+4), colorGray)

		if pan.Big6 == nil {
			continue
		}
		gong12 := pan.Big6[i-1]
		//月建
		jianColor := colorJian
		if gong12.IsJian {
			jianColor = colorLeader
		} else if qimen.GroundGate4[gong12.Jian] {
			jianColor = colorGate
		}
		y, x = util.CalRadiansPos(q.Y, q.X, _JianR, degrees)
		text.Draw(dst, gong12.Jian, ft, int(x-8), int(y+4), jianColor)
		y, x = util.CalRadiansPos(q.Y, q.X, _JianR, degrees+7.5)
		text.Draw(dst, gong12.JianZhi, ft, int(x-8), int(y+4), jianColor)
		if ThisGame.qmGame.TimeHorse == LunarUtil.ZHI[i] {
			y, x = util.CalRadiansPos(q.Y, q.X, _JianR, degrees+10)
			text.Draw(dst, "驿马", ft, int(x-8), int(y+4), colorLeader)
		}
		//月将
		degrees += 15
		ly1, lx1 = util.CalRadiansPos(q.Y, q.X, _JiangR-_SkyCircleWidth/5, degrees+15)
		ly2, lx2 = util.CalRadiansPos(q.Y, q.X, _JiangR+_SkyCircleWidth/5, degrees+15)
		vector.StrokeLine(dst, lx1, ly1, lx2, ly2, .5, colorGongSplit, true) //宫分割线
		jiangColor := colorJiang
		if gong12.IsJiang {
			jiangColor = colorLeader
		} else if qimen.SkyGate3[gong12.Jiang] {
			jiangColor = colorGate
		}
		y, x = util.CalRadiansPos(q.Y, q.X, _JiangR, degrees)
		text.Draw(dst, gong12.Jiang, ft, int(x-14), int(y+4), jiangColor)
		y, x = util.CalRadiansPos(q.Y, q.X, _JiangR, degrees-7.5)
		text.Draw(dst, gong12.JiangZhi, ft, int(x-10), int(y+4), jiangColor)
		//if gong12.IsSkyHorse {
		//	y, x := util.CalRadiansPos(q.Y, q.X, _JiangR, degreesXX-5)
		//	text.Draw(dst, "天马", ft, int(x-14), int(y+4), colorLeader)
		//}s
	}
}

func (q *QMShow) drawBattle(dst *ebiten.Image) {
	//qm := ThisGame.qmGame
	//pp := qm.ShowPan
	op := ebiten.DrawImageOptions{}
	var y, x float32
	//画九宫
	for i := 1; i <= 9; i++ {
		x, y = q.GetInCampPos(i)
		op.GeoM.Reset()
		op.ColorScale.Reset()
		if i == 5 {
			op.GeoM.Translate(float64(x), float64(y))
			op.ColorScale.ScaleWithColor(colorLeader)
			dst.DrawImage(q.CampM, &op)
		} else {
			op.GeoM.Translate(float64(x), float64(y))
			op.ColorScale.ScaleWithColor(colorGreen)
			dst.DrawImage(q.Camp, &op)

			x, y = q.GetInBornPos(i)
			op.GeoM.Reset()
			op.ColorScale.Reset()
			op.GeoM.Translate(float64(x), float64(y))
			op.ColorScale.ScaleWithColor(colorGreen)
			dst.DrawImage(q.ArmyA, &op)

			x, y = q.GetInArmyPos(i)
			op.GeoM.Reset()
			op.ColorScale.Reset()
			op.GeoM.Translate(float64(x), float64(y))
			op.ColorScale.ScaleWithColor(colorWhite)
			dst.DrawImage(q.ArmyA, &op)
		}
	}

	for i := 1; i <= 12; i++ {
		y, x = q.GetOutCampPos(i)
		op.GeoM.Reset()
		op.ColorScale.Reset()
		op.GeoM.Translate(float64(x-16), float64(y-16))
		op.ColorScale.ScaleWithColor(colorWhite)
		dst.DrawImage(q.Camp, &op)

		//op.GeoM.Reset()
		//op.ColorScale.Reset()
		//y, x = q.GetOutCampBornPos(i)
		//op = ebiten.DrawImageOptions{}
		//op.GeoM.Translate(float64(x-16), float64(y-16))
		//op.ColorScale.ScaleWithColor(colorWhite)
		//dst.DrawImage(q.Army, &op)
	}
}

func (q *QMShow) GetInCampPos(i int) (float32, float32) {
	offX, offZ := gongOffset[i][0]*_Gong9Width-_Gong9Width/2, gongOffset[i][1]*_Gong9Width-_Gong9Width/2
	px, py := q.X+float32(offX)-32-8, q.Y+float32(offZ)-32
	if i == 5 {
		return q.X - 32, q.Y - 32
	}
	return px - 32, py
}

func (q *QMShow) GetInBornPos(i int) (float32, float32) {
	offX, offZ := gongOffset[i][0]*_Gong9Width-_Gong9Width/2, gongOffset[i][1]*_Gong9Width-_Gong9Width/2
	px, py := q.X+float32(offX)-32-8, q.Y+float32(offZ)-32
	if i == 5 {
		return q.X, q.Y
	}
	return px - 32, py - 32
}

func (q *QMShow) GetInArmyPos(i int) (float32, float32) {
	offX, offZ := gongOffset[i][0]*_Gong9Width-_Gong9Width/2, gongOffset[i][1]*_Gong9Width-_Gong9Width/2
	px, py := q.X+float32(offX)-32-8, q.Y+float32(offZ)-32
	if i == 5 {
		return q.X - 32, q.Y - 32
	}
	return px - 32, py - 64
}

func (q *QMShow) GetOutCampPos(i int) (float32, float32) {
	degrees := -float32(i-1) * 30
	y, x := util.CalRadiansPos(q.Y, q.X, _JianR-_SkyCircleWidth/2, degrees)
	return y, x
}

func (q *QMShow) GetOutCampBornPos(i int) (float32, float32) {
	degrees := -float32(i-1) * 30
	y, x := util.CalRadiansPos(q.Y, q.X, _JianR-_SkyCircleWidth, degrees)
	return y, x
}

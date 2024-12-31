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
	solar calendar.Solar
	dirty bool

	TaiJi    *ebiten.Image
	Sun      *Sprite
	Moon     *Sprite
	BaGua    map[int]*ebiten.Image
	CampM    *ebiten.Image
	Camp     *ebiten.Image
	Army     *ebiten.Image
	ArmyA    *ebiten.Image
	DutyFlag *ebiten.Image

	Battle *Battle

	circleDegrees [360][4]float32 //周天刻度
	jq4Loc        [4][2]int       //春分夏至秋分冬至
	xiuLoc        [28]SegmentPos  //星宿 分隔,位置
	cnstLoc       [12]SegmentPos  //星座 分隔,位置
	big6SkySeg    [12][4]float32  //月将天门盘
	big6EarthSeg  [12][4]float32  //月建地户盘
	zhiPos        [12][2]int      //本宫支
	jianPos       [12][2]int      //月建
	jianZhiPos    [12][2]int      //月建支
	jianGongPos   [12][2]int      //月建宫支
	jiangPos      [12][2]int      //月将
	jiangGanPos   [12][2]int      //月将干
	jiangZhiPos   [12][2]int      //月将支
	jiangGongPos  [12][2]int      //月将宫支
	//horsePos      [12][2]int      //驿马
}

func NewQiMenShow(centerX, centerY int) *QMShow {
	bg := make(map[int]*ebiten.Image, 9)
	for i := 1; i <= 9; i++ {
		bg[i] = graphic.NewBaGuaImage(qimen.Diagrams9(i), _BaGuaSize)
	}

	return &QMShow{
		X: float32(centerX), Y: float32(centerY),
		TaiJi:    graphic.NewTaiJiImage(_TaiJiSize),
		BaGua:    bg,
		CampM:    graphic.NewCampImage(64),
		Camp:     graphic.NewCampImage(32),
		Army:     graphic.NewArmyImage("庚", 32, 0),
		ArmyA:    graphic.NewArmyImage("兵", 32, 1),
		DutyFlag: graphic.NewFlagImage(16),
		Battle:   NewBattle(),
		dirty:    true,
	}
}
func (q *QMShow) Update() {
	q.count++
	q.count %= 360
	q.Battle.Update()

	pan := ThisGame.qmGame
	if pan == nil {
		return
	}
	sCal := pan.Solar
	if q.solar != *sCal {
		q.solar = *sCal
		q.dirty = true
	}
	if !q.dirty {
		return
	}
	q.dirty = false

	h := sCal.GetHour()
	m := sCal.GetMinute()
	cx, cy := q.X, q.Y

	//周天刻度
	for i := 0; i < 360; i++ {
		ly1, lx1 := util.CalRadiansPos(cy, cx, _XiuR-_SkyCircleWidth/5, float32(i))
		if i%10 == 0 {
			ly1, lx1 = util.CalRadiansPos(cy, cx, _XiuR-float32(_SkyCircleWidth)/8, float32(i))
		} else if i%5 == 0 {
			ly1, lx1 = util.CalRadiansPos(cy, cx, _XiuR-float32(_SkyCircleWidth)/12, float32(i))
		}
		ly2, lx2 := util.CalRadiansPos(cy, cx, _XiuR-_SkyCircleWidth/4, float32(i))
		q.circleDegrees[i][0], q.circleDegrees[i][1], q.circleDegrees[i][2], q.circleDegrees[i][3] = lx1, ly1, lx2, ly2
	}

	//日
	if q.Sun == nil {
		q.Sun = NewSprite(graphic.NewSunImage(_StarSize), colorSun)
	}
	degreesS := -((float32(h) + float32(m)/60) * 15) //本地时区太阳角度 0~360 0时0度
	sy, sx := util.CalRadiansPos(cy, cx, _JianR+12, degreesS)
	q.Sun.MoveTo(int(sx-_StarSize/2), int(sy-_StarSize/2))
	//月
	if q.Moon == nil {
		q.Moon = NewSprite(graphic.NewMoonImage(_StarSize), colorMoon)
	}
	degreesM := -((float32(h)+float32(m)/60)*15 - float32(pan.Lunar.GetDay()-1)/float32(pan.LunarMonthDays)*360)
	my, mx := util.CalRadiansPos(cy, cx, _JianR-12, degreesM)
	q.Moon.MoveTo(int(mx-_StarSize/2), int(my-_StarSize/2))

	year := pan.Solar.GetYear()
	se := calendar.NewSolar(year, 3, 24, 0, 0, 0)
	cf := se.GetLunar().GetPrevJieQi().GetSolar() //今年春分
	seN := calendar.NewSolar(year+1, 3, 24, 0, 0, 0)
	cfN := seN.GetLunar().GetPrevJieQi().GetSolar() //下一年春分
	yearMinute := cfN.SubtractMinute(cf)
	pr := cfN.SubtractMinute(pan.Solar)
	degreesR0 := float32(pr) / float32(yearMinute) * 360 //春分角
	degreesR := degreesS + degreesR0                     //春分角+太阳角

	y, x := util.CalRadiansPos(cy, cx, _XiuR+5, degreesR)
	q.jq4Loc[0][0], q.jq4Loc[0][1] = int(x-6), int(y) //春分 白羊双鱼界
	y, x = util.CalRadiansPos(cy, cx, _XiuR+5, degreesR+90)
	q.jq4Loc[1][0], q.jq4Loc[1][1] = int(x-6), int(y) //夏至
	y, x = util.CalRadiansPos(cy, cx, _XiuR+5, degreesR+180)
	q.jq4Loc[2][0], q.jq4Loc[2][1] = int(x-6), int(y) //秋分
	y, x = util.CalRadiansPos(cy, cx, _XiuR+5, degreesR+270)
	q.jq4Loc[3][0], q.jq4Loc[3][1] = int(x-6), int(y) //冬至
	//星宿
	for i := 0; i < 28; i++ {
		xi := qimen.Xiu28[i]
		degrees := degreesR + qimen.XiuAngle[xi]
		ly1, lx1 := util.CalRadiansPos(cy, cx, _XiuR+_SkyCircleWidth/5, degrees)
		ly2, lx2 := util.CalRadiansPos(cy, cx, _XiuR-_SkyCircleWidth/6, degrees)
		degreesNext := qimen.XiuAngle[qimen.Xiu28[i+1]]
		if qimen.XiuAngle[xi] > degreesNext {
			degreesNext += 360
		}
		degreesMid := degrees + (degreesNext-qimen.XiuAngle[xi])/2
		y, x = util.CalRadiansPos(cy, cx, _XiuR, degreesMid)
		q.xiuLoc[i] = SegmentPos{lx1, ly1, lx2, ly2, int(x - 6), int(y + 4)}
	}
	//星座
	for i := 0; i < 12; i++ {
		degrees := degreesR + float32(360*i/12)
		ly1, lx1 := util.CalRadiansPos(cy, cx, _CnstR-_SkyCircleWidth/4, degrees)
		ly2, lx2 := util.CalRadiansPos(cy, cx, _CnstR+_SkyCircleWidth/5, degrees)
		y, x = util.CalRadiansPos(cy, cx, _CnstR+3, degrees+15)
		q.cnstLoc[i] = SegmentPos{lx1, ly1, lx2, ly2, int(x - 6), int(y + 6)}
	}

	//节气进度
	prevJie := pan.Lunar.GetPrevJie()
	nextJie := pan.Lunar.GetNextJie()
	prevQi := pan.Lunar.GetPrevQi()
	nextQi := pan.Lunar.GetNextQi()
	proJie := float32(pan.Solar.SubtractMinute(prevJie.GetSolar())) / float32(nextJie.GetSolar().SubtractMinute(prevJie.GetSolar()))
	proQi := float32(pan.Solar.SubtractMinute(prevQi.GetSolar())) / float32(nextQi.GetSolar().SubtractMinute(prevQi.GetSolar()))

	mD := m
	if h%2 == 0 {
		mD += 60
	}
	for i := 0; i < 12; i++ {
		degrees := -float32(i) * 30
		y, x = util.CalRadiansPos(cy, cx, _JianR-_SkyCircleWidth/2, degrees)
		q.zhiPos[i] = [2]int{int(x - 8), int(y + 4)}
		if pan.Big6 == nil {
			continue
		}

		//月建
		degreesJ := degrees + 30 - (30 * proJie) - float32(mD)*30/120 //分针 TODO 节进度
		ly1, lx1 := util.CalRadiansPos(cy, cx, _JianR-_SkyCircleWidth/4, degreesJ+15)
		ly2, lx2 := util.CalRadiansPos(cy, cx, _JianR+_SkyCircleWidth/5, degreesJ+15)
		q.big6SkySeg[i] = [4]float32{lx1, ly1, lx2, ly2}

		y, x = util.CalRadiansPos(cy, cx, _JianR, degreesJ)
		q.jianPos[i] = [2]int{int(x - 8), int(y + 4)}
		y, x = util.CalRadiansPos(cy, cx, _JianR, degreesJ-7.5)
		q.jianGongPos[i] = [2]int{int(x - 8), int(y + 4)}
		y, x = util.CalRadiansPos(cy, cx, _JianR, degreesJ+7.5)
		q.jianZhiPos[i] = [2]int{int(x - 8), int(y + 4)}
		//zhi := LunarUtil.ZHI[i+1]
		//if ThisGame.qmGame.TimeHorse == zhi {
		//	y, x = util.CalRadiansPos(cy, cx, _JianR, degreesJ+10)
		//	q.horsePos[i] = [2]int{int(x - 8), int(y + 4)}
		//}
		//月将
		degreesJJ := degrees + 30 - (30 * proQi) - float32(mD)*30/120 //分针 TODO 节进度
		ly1, lx1 = util.CalRadiansPos(cy, cx, _JiangR-_SkyCircleWidth/5, degreesJJ+15)
		ly2, lx2 = util.CalRadiansPos(cy, cx, _JiangR+_SkyCircleWidth/5, degreesJJ+15)
		q.big6EarthSeg[i] = [4]float32{lx1, ly1, lx2, ly2}

		y, x = util.CalRadiansPos(cy, cx, _JiangR, degreesJJ)
		q.jiangPos[i] = [2]int{int(x - 14), int(y + 4)}
		y, x = util.CalRadiansPos(cy, cx, _JiangR, degreesJJ+7.5)
		q.jiangGongPos[i] = [2]int{int(x - 10), int(y + 4)}
		y, x = util.CalRadiansPos(cy, cx, _JiangR, degreesJJ-7)
		q.jiangGanPos[i] = [2]int{int(x - 10), int(y + 4)}
		y, x = util.CalRadiansPos(cy, cx, _JiangR, degreesJJ-10)
		q.jiangZhiPos[i] = [2]int{int(x - 10), int(y + 4)}
	}

	//战场
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
	ft14, _ := GetFontFace(14)
	ft28, _ := GetFontFace(28)
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
		ft14, px, 48, colorWhite)
	text.Draw(dst, fmt.Sprintf("干支  %s %s %s %s",
		lunar.GetYearInGanZhiExact(), lunar.GetMonthInGanZhiExact(), lunar.GetDayInGanZhiExact(), lunar.GetTimeInGanZhi()),
		ft14, px, 64, colorLeader)
	//中字
	yg, yz := lunar.GetYearGanExact(), lunar.GetYearZhiExact()
	mg, mz := lunar.GetMonthGanExact(), lunar.GetMonthZhiExact()
	dg, dz := lunar.GetDayGanExact(), lunar.GetDayZhiExact()
	hg, hz := lunar.GetTimeGan(), lunar.GetTimeZhi()
	bpx := px + 380
	text.Draw(dst, yg, ft28, bpx, 64, ColorGanZhi(yg))
	text.Draw(dst, yz, ft28, bpx, 64+32, ColorGanZhi(yz))
	text.Draw(dst, mg, ft28, bpx+32, 64, ColorGanZhi(mg))
	text.Draw(dst, mz, ft28, bpx+32, 64+32, ColorGanZhi(mz))
	text.Draw(dst, dg, ft28, bpx+64, 64, ColorGanZhi(dg))
	text.Draw(dst, dz, ft28, bpx+64, 64+32, ColorGanZhi(dz))
	text.Draw(dst, hg, ft28, bpx+96, 64, ColorGanZhi(hg))
	text.Draw(dst, hz, ft28, bpx+96, 64+32, ColorGanZhi(hz))

	text.Draw(dst, fmt.Sprintf("旬首  %s %s %s %s",
		lunar.GetYearXunExact(), lunar.GetMonthXunExact(), lunar.GetDayXunExact(), lunar.GetTimeXun()),
		ft14, px, 64+16, colorWhite)
	text.Draw(dst, "空亡", ft14, px, 96, colorWhite)
	text.Draw(dst, lunar.GetYearXunKongExact(), ft14, px+44, 96, colorGray)
	text.Draw(dst, lunar.GetMonthXunKongExact(), ft14, px+74, 96, colorGray)
	text.Draw(dst, lunar.GetDayXunKongExact(), ft14, px+112, 96, colorGray)
	text.Draw(dst, lunar.GetTimeXunKong(), ft14, px+146, 96, colorRed)
	text.Draw(dst, fmt.Sprintf("%s%s", pp.JieQi, pp.JieQiDate), ft14, px, 96+16, colorWhite)
	text.Draw(dst, fmt.Sprintf("%s%s", pp.JieQiNext, pp.JieQiDateNext), ft14, px, 96+16*2, colorWhite)
	text.Draw(dst, pp.JuText, ft14, px, 96+16*3, colorWhite)
	text.Draw(dst, pp.DutyText, ft14, px, 96+16*4, colorWhite)
	text.Draw(dst, pp.YueJiang, ft14, px, 96+16*5, colorWhite)

	text.Draw(dst, fmt.Sprintf("月相 %s", lunar.GetYueXiang()), ft14, px, 96+16*6, colorWhite)
	text.Draw(dst, fmt.Sprintf("日值 %s%s%s", lunar.GetXiu(), lunar.GetZheng(), lunar.GetAnimal()), ft14, px, 96+16*7, colorWhite)
	text.Draw(dst, fmt.Sprintf("岁大将军 %s", qimen.BigJiang[lunar.GetYearZhiExact()]), ft14, px, 96+16*8, colorWhite)
	text.Draw(dst, fmt.Sprintf("月建大将军 %s", qimen.BigJiang[lunar.GetMonthZhiExact()]), ft14, px, 96+16*9, colorWhite)

	text.Draw(dst, "符使马时", ft14, px, 96+16*10, colorDuty)
	text.Draw(dst, "击刑", ft14, px, 96+16*11, colorJiXing)
	text.Draw(dst, "门迫", ft14, px, 96+16*12, colorMengPo)
	text.Draw(dst, "入墓", ft14, px, 96+16*13, colorTomb)
	text.Draw(dst, "刑墓", ft14, px, 96+16*14, colorXingMu)

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
		horse = "马"
		//horse = "馬"
	}
	star := qimen.Star0 + g.Star
	if g.Star == "" {
		star = ""
	}
	door := g.Door + qimen.Door0
	if g.Door == "" {
		door = ""
	}

	y += 36                                                 //神盘
	text.Draw(dst, empty, ft, int(x+8), int(y), colorWhite) //空亡
	text.Draw(dst, g.God, ft, int(x+24), int(y),
		util.If(g.God == qimen.QMGod8(1), colorDuty, colorWhite)) //神盘
	text.Draw(dst, horse, ft, int(x+24+64)+rand.Intn(2), int(y)+rand.Intn(2), colorDuty) //驿马
	if g.God == "值符" {
		text.Draw(dst, pp.Xun, ft, int(x+24), int(y+16), colorGray)
	}
	y += 32                                                    //天盘
	text.Draw(dst, g.PathGan, ft, int(x+8), int(y), colorGray) //流干
	text.Draw(dst, g.HideGan, ft, int(x+8), int(y), colorGray) //隐干
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
	guestGanCC := colorWhite
	if jiXing {
		guestGanCC = colorJiXing
		if guestGanTomb {
			guestGanCC = colorXingMu
		}
	} else if guestGanTomb {
		guestGanCC = colorTomb
	}
	text.Draw(dst, g.GuestGan, ft, int(x+64+24), int(y), guestGanCC) //天盘干
	if g.GuestGan == pp.Gan && pp == qm.TimePan {
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
	hostGanCC := colorWhite
	if jiXing {
		hostGanCC = colorJiXing
		if hostGanTomb {
			hostGanCC = colorXingMu
		}
	} else if hostGanTomb {
		hostGanCC = colorTomb
	}
	text.Draw(dst, g.HostGan, ft, int(x+64+24), int(y), hostGanCC) //地盘干
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
	cx, cy := q.X, q.Y
	//建星地户盘
	vector.StrokeCircle(dst, cx, cy, _JianR, _SkyCircleWidth/2, colorGroundGateCircle, true)
	//月将天门盘
	vector.StrokeCircle(dst, cx, cy, _JiangR, _SkyCircleWidth/2, colorSkyGateCircle, true)
	//星宿盘
	vector.StrokeCircle(dst, cx, cy, _XiuR, _SkyCircleWidth/2, colorPowerCircle, true)
	vector.StrokeCircle(dst, cx, cy, _CnstR, _SkyCircleWidth/2, colorSkyGateCircle, true)
	//周天刻度
	cd := &q.circleDegrees
	for i := 0; i < 360; i++ {
		vector.StrokeLine(dst, cd[i][0], cd[i][1], cd[i][2], cd[i][3], .5, color.Black, true)
	}
	//日
	q.Sun.Draw(dst)
	//月
	q.Moon.Draw(dst)
	//TODO 行星..

	jqPos := &q.jq4Loc
	text.Draw(dst, "*", ft, jqPos[0][0], jqPos[0][1], colorGreen)  //春分角 白羊双鱼界
	text.Draw(dst, "*", ft, jqPos[1][0], jqPos[1][1], colorRed)    //夏至
	text.Draw(dst, "*", ft, jqPos[2][0], jqPos[2][1], colorYellow) //秋分
	text.Draw(dst, "*", ft, jqPos[3][0], jqPos[3][1], colorWhite)  //冬至

	//画28星宿
	for i := 0; i < 28; i++ {
		loc := &q.xiuLoc[i]
		vector.StrokeLine(dst, loc.Lx1, loc.Ly1, loc.Lx2, loc.Ly2, .5, colorWhite, true) //星宿
		text.Draw(dst, qimen.Xiu28[i], ft, loc.X, loc.Y, colorLeader)                    //星宿
	}
	//画12星座
	for i := 0; i < 12; i++ {
		loc := q.cnstLoc[i]
		vector.StrokeLine(dst, loc.Lx1, loc.Ly1, loc.Lx2, loc.Ly2, .5, colorWhite, true) //星座分隔线
		text.Draw(dst, qimen.ConstellationShort[i], ft, loc.X, loc.Y, colorLeader)       //星座
	}
	q.drawBig6(dst)
}
func (q *QMShow) drawBig6(dst *ebiten.Image) {
	ft, _ := GetFontFace(14)
	pan := ThisGame.qmGame
	//大六壬
	if pan.Big6 == nil {
		return
	}
	for i := 0; i < 12; i++ {
		loc := q.big6EarthSeg[i]
		vector.StrokeLine(dst, loc[0], loc[1], loc[2], loc[3], .5, colorGongSplit, true) //宫分割线

		//地支本宫位
		zhi := LunarUtil.ZHI[i+1]
		zhiGongTxt := zhi
		//if slices.Contains(qimen.KongWang[pan.ShowPan.Xun], zhi) {
		//	zhiGongTxt = fmt.Sprintf("〇%s", zhi)
		//	//} else if slices.Contains(qimen.KongWang[pan.ShowPan.Xun], LunarUtil.ZHI[qimen.Idx12[i+6]]) {
		//	//	zhiGongTxt = fmt.Sprintf("%s虚", LunarUtil.ZHI[i])
		//}
		pos := q.zhiPos[i]
		text.Draw(dst, zhiGongTxt, ft, pos[0], pos[1], colorGray)

		g := pan.Big6.Gong[i]
		//月建
		jianColor := colorJian
		if g.IsJian {
			jianColor = colorLeader
		} else if qimen.GroundGate4[g.Jian] {
			jianColor = colorGate
		}
		pos = q.jianPos[i]
		text.Draw(dst, g.Jian, ft, pos[0], pos[1], jianColor)
		pos = q.jianZhiPos[i]
		text.Draw(dst, g.JianZhi, ft, pos[0], pos[1], jianColor)
		pos = q.jianGongPos[i]
		text.Draw(dst, zhi, ft, pos[0], pos[1], colorGrayB) //宫支
		//月将
		loc = q.big6SkySeg[i]
		vector.StrokeLine(dst, loc[0], loc[1], loc[2], loc[3], .5, colorGongSplit, true) //宫分割线
		jiangColor := colorJiang
		if g.IsJiang {
			jiangColor = colorLeader
		} else if qimen.SkyGate3[g.Jiang] {
			jiangColor = colorGate
		}
		pos = q.jiangPos[i]
		text.Draw(dst, g.Jiang, ft, pos[0], pos[1], jiangColor) //月将名
		pos = q.jiangGanPos[i]
		text.Draw(dst, g.JiangGan, ft, pos[0], pos[1], util.If(g.JiangGan == pan.Lunar.GetDayGan(), colorRed, colorJiang)) //月将干名
		pos = q.jiangZhiPos[i]
		text.Draw(dst, g.JiangZhi, ft, pos[0], pos[1], jiangColor) //月将支名
		pos = q.jiangGongPos[i]
		text.Draw(dst, zhi, ft, pos[0], pos[1], colorGray) //宫支

		//if g.IsSkyHorse {
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

package world

import (
	"fmt"
	"github.com/6tail/lunar-go/calendar"
	"github.com/deminzhang/qimen-go/graphic"
	"github.com/deminzhang/qimen-go/qimen"
	"github.com/deminzhang/qimen-go/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	_ "github.com/mattn/go-sqlite3"
	"image/color"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	_ "xorm.io/core"
)

const (
	G  = 6.67430 * 1e-11 //万有引力常数(Nm²/kg²)
	C  = 299792458       //光速(m/s)
	AU = 149597870.7e3   //天文单位(m)

	outCircleR  = 240
	outCircleW  = 16
	outCircleR0 = outCircleR + outCircleW/2
	sphereR     = outCircleR - outCircleW*4

	DataTimeMin        = "2006-01-02 15:04"
	DateTimeNASA       = "2006-Jan-02 15:04"
	NASADataFile       = "nasa_data.db"
	NASADataTimeLast   = time.Hour * 24 //一次查询时间范围
	NASADataStepSize   = "1h"           //星体数据步长 1h 1d 1m
	ObserveDataTable   = "observe_data"
	CelestialBodyTable = "celestial_body"

	DataStartYear = 1600 //NASA数据支持开始年份
	DataEndYear   = 2500 //NASA数据支持结束年份
)

// 建星是太阳位置,星座是太阳上升,所以星座相当建星是指定地支时
// 月将 寅丑子亥戌酉申未午巳辰卯
// 建星 子丑寅卯辰巳午未申酉戌亥
// 星座 射摩瓶鱼羊牛双蟹狮室秤蝎
// 四象 火土风水火土风水火土风水

type Astrolabe struct {
	sync.RWMutex
	X, Y           float32
	dirty          bool
	solarX, solarY float32
	degreesS       float32 //太阳角度
	observer       int     //观察者

	solar    calendar.Solar
	timezone string
	tzOffset int
	tzRA0    float32 //春分点角度

	dataQuerying int
	Ephemeris    map[string]*ObserveData

	ConstellationLoc [12]SegmentPos //星座位
	AstrolabeLoc     [12]SegmentPos //宫位
	XiuLoc           [28]SegmentPos //星宿宫位

	Earth *Sprite
	Sun   *Sprite
	Moon  *Sprite
}
type CelestialBody struct {
	Id          int `gorm:"primarykey" xorm:"pk"`
	name        string
	nameCN      string
	subBody     []int   //行星系子体
	satellite   []int   //卫星
	mass        float64 //质量 kg
	radius      float32 //平均轨道半径 AU
	orbitCenter int     //轨道中心对象
	gravity     float64 //引力 N
	GMin, GMax  float64 //引力范围

	color            color.RGBA //RedShift红移/BlueShift蓝移
	drawR            float32    //星盘半径
	drawX, drawY     float32    //星盘坐标 NASA
	drawX2, drawY2   float32    //星盘坐标 本地
	sphereX, sphereY float32    //天球坐标
}

// DrawR 星盘半径
func (c *CelestialBody) DrawR() float32 {
	if c.radius > 0 {
		return c.radius * 150
	} else {
		return c.drawR
	}
}

// Primary 公转主星
//func (c *CelestialBody) Primary() int {
//	if c.Id < 100 {
//		return c.Id
//	}
//	if c.Id%100 == 99 {
//		return 10
//	}
//	return c.Id/100*100 + 99
//}

var Bodies = map[int]*CelestialBody{
	10:  {Id: 10, name: "Sun", nameCN: "日", drawR: 0, mass: 1988500 * 1e24},
	199: {Id: 199, name: "Mercury", nameCN: "水", radius: (0.31 + 0.47) / 2, drawR: 15, mass: 3.302 * 1e23},
	299: {Id: 299, name: "Venus", nameCN: "金", radius: 0.72, drawR: 30, mass: 4.8685 * 1e24},
	301: {Id: 301, name: "Moon", nameCN: "月", drawR: 10, mass: 7.349 * 1e22, orbitCenter: 399},
	399: {Id: 399, name: "Earth", nameCN: "地", radius: 1, drawR: 50, mass: 5.97219 * 1e24, //+-0.0006
		satellite: []int{301}},
	3: {Id: 3, name: "EarthBarycenter", nameCN: "地月系重心", drawR: 50,
		subBody: []int{301, 399}},
	401: {Id: 401, name: "Phobos", drawR: 5, mass: 1.0659 * 1e16, orbitCenter: 499},
	402: {Id: 402, name: "Deimos", drawR: 5, mass: 1.4762 * 1e15, orbitCenter: 499},
	499: {Id: 499, name: "Mars", nameCN: "火", radius: (1.38 + 1.66) / 2, drawR: 70, mass: 6.4171 * 1e23,
		satellite: []int{401, 402}},
	4: {Id: 4, name: "MarsBarycenter", nameCN: "火系重心", drawR: 70,
		subBody: []int{401, 402, 499}},
	501: {Id: 501, name: "Io", drawR: 5, mass: 8.9319 * 1e22, orbitCenter: 599},
	502: {Id: 502, name: "Europa", drawR: 5, mass: 4.7998 * 1e22, orbitCenter: 599},
	503: {Id: 503, name: "Ganymede", drawR: 6, mass: 1.4819 * 1e23, orbitCenter: 599},
	504: {Id: 504, name: "Callisto", drawR: 7, mass: 1.0759 * 1e23, orbitCenter: 599},
	599: {Id: 599, name: "Jupiter", nameCN: "木", radius: (4.95 + 5.46) / 2, drawR: 90, mass: 1.8982 * 1e27, //18981.8722 +- .8817
		satellite: []int{501, 502, 503, 504}},
	601: {Id: 601, name: "Mimas", drawR: 5, mass: 3.7493 * 1e19, orbitCenter: 699},
	602: {Id: 602, name: "Enceladus", drawR: 5, mass: 1.0802 * 1e20, orbitCenter: 699},
	603: {Id: 603, name: "Tethys", drawR: 6, mass: 6.1745 * 1e20, orbitCenter: 699},
	604: {Id: 604, name: "Dione", drawR: 6, mass: 1.0955 * 1e21, orbitCenter: 699},
	605: {Id: 605, name: "Rhea", drawR: 6, mass: 2.306 * 1e21, orbitCenter: 699},
	606: {Id: 606, name: "Titan", drawR: 6, mass: 1.3455 * 1e23, orbitCenter: 699},
	607: {Id: 607, name: "Iapetus", drawR: 7, mass: 1.8053 * 1e21, orbitCenter: 699},
	699: {Id: 699, name: "Saturn", nameCN: "土", radius: (9.03 + 9.54) / 2, drawR: 105, mass: 5.6834 * 1e26,
		satellite: []int{601, 602, 603, 604, 605, 606, 607}},
	799: {Id: 799, name: "Uranus", nameCN: "天", radius: (18.31 + 19.19) / 2, drawR: 120, mass: 8.6813 * 1e25},
	899: {Id: 899, name: "Neptune", nameCN: "海", radius: (29.81 + 30.36) / 2, drawR: 135, mass: 1.02409 * 1e26},
	901: {Id: 901, name: "Charon", nameCN: "卡", drawR: 5, mass: 1 * 1e22, orbitCenter: 999},
	999: {Id: 999, name: "Pluto", nameCN: "冥", radius: (29.66 + 49.31) / 2, drawR: 150, mass: 1.307 * 1e22, //1.307+-0.018
		satellite: []int{901}},
	9: {Id: 9, name: "PlutoBarycenter", nameCN: "冥王系重心", drawR: 150,
		subBody: []int{901, 999}},
}
var Draws = []int{
	10, 199, 299, 399, 499, 599, 699, // 799, 899, 999,
}

func NewAstrolabe(x, y int) *Astrolabe {
	tz, offset := time.Now().Local().Zone()
	a := &Astrolabe{
		X: float32(x), Y: float32(y),
		observer:  399,
		timezone:  tz,
		tzOffset:  offset,
		Ephemeris: make(map[string]*ObserveData),
		dirty:     true,
	}
	var bds []*CelestialBody
	err := db.Table(CelestialBodyTable).Find(&bds)
	if err == nil {
		for _, bd := range bds {
			body := Bodies[bd.Id]
			if body == nil {
				continue
			}
			body.GMin, body.GMax = bd.GMin, bd.GMax
		}
	}
	return a
}
func (a *Astrolabe) SetPos(x, y float32) {
	a.X, a.Y = x, y
	a.dirty = true
}
func (a *Astrolabe) Update() {
	pan := ThisGame.qmGame
	sCal := pan.Solar
	if a.solar != *sCal {
		a.solar = *sCal
		a.dirty = true
	}
	if !a.dirty {
		return
	}
	a.dirty = false

	hour := sCal.GetHour()
	minute := sCal.GetMinute()

	cx, cy := a.X, a.Y
	//qm := ThisGame.qiMen
	//cx, cy = qm.X, qm.Y
	if a.Earth == nil {
		a.Earth = NewSprite(graphic.NewEarthImage(16), colorBlue)
		a.Earth.onMove = func(sx, sy, dx, dy int) {
			a.X += float32(dx)
			a.Y += float32(dy)
			a.dirty = true
		}
		ThisGame.AddSprite(a.Earth)
	}
	a.Earth.MoveTo(int(cx-8), int(cy-8))
	//计算太阳位置 以时间计角度
	degreesS := 360 - (float32(hour)+float32(minute)/60)*15 //本地时区太阳角度 0~360 0时0度
	degreesSO := 360 - degreesS
	solarY, solarX := util.CalRadiansPos(cy, cx, Bodies[a.observer].DrawR(), degreesS)
	a.solarX, a.solarY = solarX, solarY
	a.degreesS = degreesS
	if a.Sun == nil {
		a.Sun = NewSprite(graphic.NewSunImage(16), colorYellow)
	}
	a.Sun.MoveTo(int(solarX-8), int(solarY-8))
	{ //计算月球位置 暂以农历近似
		lDay := ThisGame.qmGame.Lunar.GetDay()
		mDays := ThisGame.qmGame.LunarMonthDays
		degreesM := -(float32(hour)+float32(minute)/60)*15 + float32(lDay-1)/float32(mDays)*360
		moon := Bodies[301]
		moonY, moonX := util.CalRadiansPos(cy, cx, moon.DrawR(), degreesM)
		moon.drawX, moon.drawY = moonX, moonY
		v1 := (&util.Vec2[float32]{X: moon.drawX - cx, Y: moon.drawY - cy}).ScaledToLength(sphereR)
		moon.sphereX = cx + v1.X
		moon.sphereY = cy + v1.Y
		if a.Moon == nil {
			a.Moon = NewSprite(graphic.NewMoonImage(16), colorWhite)
		}
		a.Moon.MoveTo(int(moonX-8), int(moonY-8))
	}

	m := Bodies[a.observer].mass
	for _, id := range Draws {
		body := Bodies[id]
		if id == a.observer {
			body.drawX, body.drawY = cx, cy
			body.color = colorBlueShift
			for _, sid := range body.satellite {
				bd := Bodies[sid]
				oe := a.GetEphemeris(sid, sCal)
				if oe == nil {
					a.dirty = true
					bd.color = colorBlueShift
					continue
				}
				bd.gravity = G * bd.mass * m / math.Pow(oe.Delta*AU, 2)
				a.updateGravityRange(sid, bd.gravity)
				if oe.Deldot > 0 {
					bd.color = colorBlueShift
				} else {
					bd.color = colorRedShift
				}
			}
			continue
		}
		oe := a.GetEphemeris(id, sCal)
		if oe == nil {
			a.dirty = true
			continue
		}
		switch id {
		case 10:
			body.drawX, body.drawY = a.solarX, a.solarY
			v1 := (&util.Vec2[float32]{X: body.drawX - cx, Y: body.drawY - cy}).ScaledToLength(sphereR)
			body.sphereX = cx + v1.X
			body.sphereY = cy + v1.Y

			a.tzRA0 = degreesS - oe.RARadius()
		default:
			ost := 180 - oe.STO - oe.SOT
			var degrees float32
			switch oe.SOTR {
			case "/T": //TRAILS 目标跟踪太阳
				degrees = 90 - 270 - degreesSO - ost
			case "/L": //LEADS 目标引领太阳
				degrees = 90 - 270 - degreesSO + ost
			}
			body.drawY, body.drawX = util.CalRadiansPos(solarY, solarX, body.DrawR(), degrees)
			v1 := (&util.Vec2[float32]{X: body.drawX - cx, Y: body.drawY - cy}).ScaledToLength(sphereR)
			body.sphereX = cx + v1.X
			body.sphereY = cy + v1.Y
			//switch id {
			//case 599:
			//	date0, _ := time.Parse(time.DateTime, qimen.Jupiter0)
			//	solar0 := calendar.NewSolarFromDate(date0)
			//	period := qimen.JupiterPeriod
			//	degreesJ := degreesSO + float32(360*float64(a.solar.SubtractMinute(solar0))/period)
			//	y, x := util.CalRadiansPos(solarY, solarX, body.DrawR()+4, degreesJ)
			//	body.drawY2, body.drawX2 = y, x
			//}
		}
		body.gravity = G * body.mass * m / math.Pow(oe.Delta*AU, 2)
		a.updateGravityRange(id, body.gravity)
		if oe.Deldot > 0 {
			body.color = colorBlueShift
		} else {
			body.color = colorRedShift
		}
	}
	a.calGongLocation()
	return
}
func (a *Astrolabe) updateGravityRange(tid int, gravity float64) {
	body := Bodies[tid]
	mi, mx := body.GMin, body.GMax
	if body.GMin == 0 {
		body.GMin = body.gravity
	} else {
		body.GMin = min(body.gravity, body.GMin)
	}
	body.GMax = max(body.GMin, body.GMax)
	if mi != body.GMin || mx != body.GMax {
		_, err := db.Insert(body)
		if err != nil {
			return
		}
	}
}
func (a *Astrolabe) calGongLocation() {
	cx, cy := a.X, a.Y
	for i := 1; i <= 12; i++ { //固定宫位
		degrees := float64(i)*30 - 90
		r := float64(outCircleR0 - outCircleW*2)
		ly1, lx1 := util.CalRadiansPos(float64(cy), float64(cx), r-outCircleW/2, degrees)
		ly2, lx2 := util.CalRadiansPos(float64(cy), float64(cx), r+outCircleW/2, degrees)
		y, x := util.CalRadiansPos(float64(cy), float64(cx), r, degrees-15)
		a.AstrolabeLoc[i-1] = SegmentPos{float32(lx1), float32(ly1), float32(lx2), float32(ly2), int(x), int(y)}
	}
	for i := 1; i <= 12; i++ {
		degrees := a.tzRA0 + float32(i-1)*30
		r := float64(outCircleR0)
		ly1, lx1 := util.CalRadiansPos(cy, cx, float32(r-outCircleW/2), degrees)
		ly2, lx2 := util.CalRadiansPos(cy, cx, float32(r+outCircleW/2), degrees)
		y, x := util.CalRadiansPos(cy, cx, float32(r), degrees+15)
		a.ConstellationLoc[i-1] = SegmentPos{lx1, ly1, lx2, ly2, int(x), int(y)}
	}
	for i := 1; i <= 28; i++ {
		xiu := qimen.Xiu28[i]
		degrees := a.tzRA0 + qimen.XiuAngle[xiu]
		r := float64(outCircleR0 - outCircleW)
		ly1, lx1 := util.CalRadiansPos(cy, cx, float32(r-outCircleW/2), degrees)
		ly2, lx2 := util.CalRadiansPos(cy, cx, float32(r+outCircleW/2), degrees)
		y, x := util.CalRadiansPos(cy, cx, float32(r), degrees+6.4285)
		a.XiuLoc[i-1] = SegmentPos{lx1, ly1, lx2, ly2, int(x), int(y)}
	}
}
func (a *Astrolabe) DataQuerying() bool {
	return a.dataQuerying != 0
}
func (a *Astrolabe) Draw(dst *ebiten.Image) {
	ft, _ := GetFontFace(12)
	cx, cy := a.X, a.Y
	sX, sY := a.solarX, a.solarY
	//外圈
	w := float32(outCircleW)
	r := float32(outCircleR0)
	vector.StrokeCircle(dst, cx, cy, r, w, colorSkyGateCircle, true) //星座
	r -= w
	vector.StrokeCircle(dst, cx, cy, r, w, colorGroundGateCircle, true) //星宿
	r -= w
	vector.StrokeCircle(dst, cx, cy, r, w, colorPowerCircle, true) //天球
	r -= w
	vector.StrokeCircle(dst, cx, cy, r, w, colorGroundGateCircle, true) //宫位
	//十字线
	horizons := float32(0) //TODO 地平线按太阳视角调整
	vector.StrokeLine(dst, cx-outCircleR, cy-horizons, cx+outCircleR, cy-horizons, 1, colorCross, true)
	vector.StrokeLine(dst, cx, cy-outCircleR, cx, cy+outCircleR, 1, colorCross, true)
	//春分点 RA0
	tzY, tzX := util.CalRadiansPos(cy, cx, outCircleR/2, a.tzRA0)
	vector.StrokeLine(dst, cx, cy, tzX, tzY, 1, colorOrbits, true)
	text.Draw(dst, "春分", ft, int(tzX), int(tzY), colorLeader)                                                        //春分点
	text.Draw(dst, fmt.Sprintf("%s月", ThisGame.qmGame.Lunar.GetYueXiang()), ft, int(cx-16), int(cy-25), colorLeader) //月相

	//画12宫
	for i := 0; i < 12; i++ {
		l := a.ConstellationLoc[i]
		vector.StrokeLine(dst, l.Lx1, l.Ly1, l.Lx2, l.Ly2, 1, colorGongSplit, true) //星宫
		//text.Draw(dst, qimen.ConstellationSymbol[i], ft, l.X-6, l.Y+6, colorJiang)  //星座符号 需要字体支持
		text.Draw(dst, qimen.ConstellationShort[i], ft, l.X-6, l.Y+6, colorJiang) //星座
		l = a.AstrolabeLoc[i]
		vector.StrokeLine(dst, l.Lx1, l.Ly1, l.Lx2, l.Ly2, 1, colorGongSplit, true) //宫
		//text.Draw(dst, fmt.Sprintf("%d", i+1), ft, l.X-4, l.Y+4, colorJiang)        //宫位
		degrees := float64(i+1)*30 - 90
		r := float64(outCircleR0 - outCircleW*2)
		//gongName := qimen.AstrolabeGong[i] //宫名
		gongName := qimen.AstrolabeGong74[i] //政余名
		DrawRotateText(dst, float64(cx-4), float64(cy+4), r, degrees-10, gongName, 12, colorJiang)
	}
	//画28星宿
	for i := 1; i <= 28; i++ {
		l := a.XiuLoc[i-1]
		vector.StrokeLine(dst, l.Lx1, l.Ly1, l.Lx2, l.Ly2, 1, colorGongSplit, true) //星宿
		text.Draw(dst, qimen.Xiu28[i], ft, l.X-4, l.Y+4, colorJiang)                //星宿
	}
	//画星体
	for _, id := range Draws {
		obj := Bodies[id]
		vector.StrokeCircle(dst, sX, sY, obj.DrawR(), 1, colorOrbits, true) //planet Orbit
		if a.observer == obj.Id {
			//地球观察者
			if a.Earth != nil {
				a.Earth.Draw(dst)
			}
		} else {
			if obj.sphereX == 0 && obj.sphereY == 0 {
				continue //查询中
			}
			//vector.StrokeLine(dst, cx, cy, obj.sphereX, obj.sphereY, 1, colorOrbits, true) // sphere line
			vector.StrokeLine(dst, cx, cy, obj.drawX, obj.drawY, 1, colorOrbits, true) // R line
			vector.StrokeCircle(dst, obj.sphereX, obj.sphereY, 2, .5, obj.color, true) // sphere

			//text.Draw(dst, qimen.StarSymbol[obj.nameCN], ft, int(obj.sphereX-8), int(obj.sphereY-8), obj.color) //星体符号需要字体支持
			text.Draw(dst, obj.nameCN, ft, int(obj.sphereX), int(obj.sphereY), obj.color)
		}
		if obj.Id == 10 { //日
			if a.Sun != nil {
				a.Sun.Draw(dst)
			}
		} else {
			vector.DrawFilledCircle(dst, obj.drawX, obj.drawY, 2, obj.color, true) //planet
			if obj.drawX2 > 0 && obj.drawY2 > 0 {
				vector.StrokeCircle(dst, obj.drawX2, obj.drawY2, 3, 1, obj.color, true) //planet2
			}
		}

		for _, sid := range obj.satellite {
			ob := Bodies[sid]
			if ob.Id == 301 {
				vector.StrokeCircle(dst, obj.drawX, obj.drawY, ob.DrawR(), 1, colorOrbits, true) //satellite Orbit
				vector.StrokeLine(dst, cx, cy, ob.sphereX, ob.sphereY, 1, colorOrbits, true)     // sphere line
				//vector.DrawFilledCircle(dst, ob.drawX, ob.drawY, 1, ob.color, true)              //moon
				if a.Moon != nil {
					a.Moon.Draw(dst)
				}
				text.Draw(dst, ob.nameCN, ft, int(ob.sphereX), int(ob.sphereY), ob.color)
			} else {
				my, mx := util.CalRadiansPos(obj.drawY, obj.drawX, ob.DrawR(), float32(rand.Intn(360)))
				vector.DrawFilledCircle(dst, mx, my, 1, obj.color, true) //satellite
			}
		}
	}
	a.DrawGravity(dst)
	switch a.dataQuerying {
	case 1:
		text.Draw(dst, "正在查询..", ft, int(cx-32), int(cy-10), color.White)
	case 2:
		text.Draw(dst, "正在观星..", ft, int(cx-32), int(cy-10), color.White)
	}
}
func (a *Astrolabe) DrawGravity(dst *ebiten.Image) {
	cx, cy := a.X+60, a.Y-440
	for _, id := range Draws {
		obj := Bodies[id]
		if obj.gravity > 0 {
			//cx += 40
			cy += 22
			DrawRangeBar(dst, cx-200, cy, 100, obj.nameCN, obj.gravity, obj.GMin, obj.GMax, obj.color)
		}
		for _, sid := range obj.satellite {
			ob := Bodies[sid]
			if ob.gravity > 0 {
				//cx += 40
				cy += 22
				DrawRangeBar(dst, cx-200, cy, 100, ob.nameCN, ob.gravity, ob.GMin, ob.GMax, ob.color)
			}
		}
	}
}

// TODO 优化为按年查询每小时正点数据，辅助插值
func (a *Astrolabe) GetEphemeris(tid int, s *calendar.Solar) *ObserveData {
	if s.GetYear() < DataStartYear || s.GetYear() >= DataEndYear {
		return nil
	}
	t := time.Date(s.GetYear(), time.Month(s.GetMonth()), s.GetDay(), s.GetHour(),
		s.GetMinute(), 0, 0, time.Local)
	sts := t.In(time.UTC).Format(DataTimeMin)
	id := fmt.Sprintf("%d_%s", tid, sts)
	if a.DataQuerying() {
		return nil
	}
	a.RWMutex.RLock()
	defer a.RWMutex.RUnlock()
	v := a.Ephemeris[id]
	if v != nil {
		return v
	}
	a.dataQuerying = 1
	var it ObserveData
	has, err := db.Table(ObserveDataTable).Where("Id = ?", id).Get(&it)
	if err != nil {
		log.Printf("Error GetEphemeris: %v\n", err)
		UIShowMsgBox(err.Error(), "确定", "确定", nil, nil)
		return nil
	}
	if has {
		a.dataQuerying = 0
		return &it
	} else {
		a.dataQuerying = 2
		te := t.Add(NASADataTimeLast)
		ets := te.Format(DataTimeMin)
		go a.QueryNASAData(tid, sts, ets)
	}
	return nil
}

func (a *Astrolabe) GetNASAData(tid int, sts, ets string) map[string]*observeDataSrc {
	urls := fmt.Sprintf("https://ssd.jpl.nasa.gov/api/horizons.api?"+
		"format=text&COMMAND='%d'&OBJ_DATA='YES'&MAKE_EPHEM='YES'&EPHEM_TYPE='OBSERVER'&CENTER='500@399'"+
		"&START_TIME='%s'&STOP_TIME='%s'&STEP_SIZE='%s'&QUANTITIES='1,20,23,24,29'",
		tid, sts, ets, url.QueryEscape(NASADataStepSize))
	//QUANTITIES='1,3,20,23,24,29'
	//Date__(UT)__HR:MN     R.A._____(ICRF)_____DEC  dRA*cosD d(DEC)/dt             delta      deldot     S-O-T /r     S-T-O  Cnst

	resp, err := http.Get(urls)
	if err != nil {
		log.Printf("Error sending GET request: %v\n", err)
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected Status Code: %d\n", resp.StatusCode)
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	data := string(body)
	//parseHead
	re := regexp.MustCompile(`(.*)( Date__\(UT\)__HR:MN.*)`)
	match := re.FindStringSubmatch(data)
	if match == nil || len(match) <= 2 {
		fmt.Println("parseHeadFail")
		return nil
	}
	head0 := match[2]
	//fmt.Printf(head0)
	var head []string //表头
	var colLen []int  //列宽
	var col []int32
	var lastChar int32
	for _, c := range head0 {
		if lastChar != 0 && lastChar != ' ' && c == ' ' {
			colLen = append(colLen, len(col))
			head = append(head, strings.TrimSpace(string(col)))
			col = nil
		}
		col = append(col, c)
		lastChar = c
	}
	colLen = append(colLen, len(col))
	head = append(head, strings.TrimSpace(string(col)))

	//parseData
	re = regexp.MustCompile(`(?s)\$\$SOE\n(.*?)\n\$\$EOE`)
	match = re.FindStringSubmatch(data)

	if match == nil || len(match) <= 1 {
		fmt.Println("未找到匹配的内容")
		return nil
	}
	lines0 := strings.Split(match[1], "\n")
	var lines = make(map[string]*observeDataSrc, len(lines0))
	for _, line := range lines0 {
		var m = map[string]string{}
		var preL int
		for i, l := range colLen {
			v0 := line[preL : preL+l]
			m[head[i]] = strings.TrimSpace(v0)
			preL += l
		}
		date := m["Date__(UT)__HR:MN"]
		dt, _ := time.Parse(DateTimeNASA, date)
		ymdhm := fmt.Sprintf("%04d-%02d-%02d %02d:%02d", dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute())
		lines[ymdhm] = &observeDataSrc{data: m}
	}
	defer resp.Body.Close()
	return lines
}
func (a *Astrolabe) QueryNASAData(tid int, sts, ets string) {
	d := a.GetNASAData(tid, sts, ets)
	if d == nil || len(d) == 0 {
		a.dataQuerying = 0
		return
	}
	// 数据存到sqlite
	a.RWMutex.Lock()
	defer a.RWMutex.Unlock()
	var its []*ObserveData
	for st, c := range d {
		iid := fmt.Sprintf("%d_%s", tid, st)
		dt, _ := time.Parse(DataTimeMin, st)
		it := &ObserveData{
			Id: iid, Target: tid,
			Year: dt.Year(), Month: int(dt.Month()), Day: dt.Day(),
			Hour: dt.Hour(), Minute: dt.Minute(),
			RA_DEC: c.data["R.A._____(ICRF)_____DEC"],
			Delta:  c.Delta(),
			Deldot: c.Deldot(),
			SOT:    c.SOT(),
			SOTR:   c.SOTR(),
			STO:    c.STO(),
			Cnst:   c.data["Cnst"],
		}
		a.Ephemeris[iid] = it
		its = append(its, it)
	}
	_, err := db.Insert(its)
	if err != nil {
		log.Printf("Error Insert: %v\n", err)
		db.Update(its)
	}
	a.dataQuerying = 0
}
func (a *Astrolabe) GetSolarPos() (float32, float32) {
	return a.solarX, a.solarY
}
func (a *Astrolabe) GetMoonPos() (float32, float32) {
	return Bodies[301].drawX, Bodies[301].drawY
}

type ObserveData struct {
	Id     string  `gorm:"primarykey" xorm:"pk"`
	Target int     `gorm:"index" xorm:"index"`
	Year   int     `gorm:"index" xorm:"index"`
	Month  int     `gorm:"index" xorm:"index"`
	Day    int     `gorm:"index" xorm:"index"`
	Hour   int     `gorm:"index" xorm:"index"`
	Minute int     `gorm:"index" xorm:"index"`
	RA_DEC string  // 赤经赤纬
	Delta  float64 // 距离 AU
	Deldot float64 // delta-dot 距离变化 KM/S 为正表示远离观察者
	SOT    float32 // S-O-T 观察者-目标-太阳角度
	SOTR   string  // "/T": TRAILS 目标跟踪S  "/L": LEADS 目标引领S
	STO    float32 // S-T-O 太阳-目标-观察者角度
	Cnst   string  // 星座
}

// RA 赤经 05 29 58.88
func (c *ObserveData) RA() string {
	return c.RA_DEC[0:11]
}

// RARadius 赤经角度
func (c *ObserveData) RARadius() float32 {
	ra := c.RA() //05 29 58.88
	ss := strings.Split(ra, " ")
	h, _ := strconv.ParseFloat(ss[0], 64)
	m, _ := strconv.ParseFloat(ss[1], 64)
	s, _ := strconv.ParseFloat(ss[2], 64)
	return float32(h*15 + m/4 + s/240)
}

// DEC 赤纬 +23 15 24.3
func (c *ObserveData) DEC() string {
	return c.RA_DEC[12:] //+23 15 24.3
}

type observeDataSrc struct {
	data map[string]string
}

func (c *observeDataSrc) Delta() float64 {
	v, _ := strconv.ParseFloat(c.data["delta"], 64)
	return v
}
func (c *observeDataSrc) Deldot() float64 {
	v, _ := strconv.ParseFloat(c.data["deldot"], 64)
	return v
}
func (c *observeDataSrc) SOT() float32 {
	v, _ := strconv.ParseFloat(c.data["S-O-T"], 64)
	return float32(v)
}
func (c *observeDataSrc) SOTR() string {
	return c.data["/r"]
}
func (c *observeDataSrc) STO() float32 {
	v, _ := strconv.ParseFloat(c.data["S-T-O"], 64)
	return float32(v)
}

package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"math/rand"
	"qimen/util"
)

const (
	G = 6.67430 * 1e-11 //万有引力常数(Nm²/kg²)
	C = 299792458       //光速(m/s)

	outCircleR = 240
	outCircleW = 24
)

type Astrolabe struct {
	centerX, centerY float32
	solarX, solarY   float32
	moonX, moonY     float32
	observer         int //观察者
}
type Celestial struct {
	Id         int
	Name       string
	NameCN     string
	Satellite  []*Celestial //卫星
	Mass       float64      //质量 kg
	AstrolabeR float32      //星盘半径
}
type CelestialData struct {
	Id     int
	Date   string
	RA_DEC string
	Delta  float64 //距离 au
	Deldot float64 //距离变化速度 km/s
	SOT    float64 //视角
	STO    float64 //视角
	Cnst   string  //星座
	//Date__(UT)__HR:MN     R.A._____(ICRF)_____DEC  dRA*cosD d(DEC)/dt             delta      deldot     S-O-T /r     S-T-O  Cnst
	//计算星体位置
	//https://ssd.jpl.nasa.gov/api/horizons.api?format=text&COMMAND='499'&OBJ_DATA='YES'&MAKE_EPHEM='YES'&EPHEM_TYPE='OBSERVER'&CENTER='500@399'&START_TIME='2024-06-14 8:00:00'&STOP_TIME='2024-06-14 10:00:00'&STEP_SIZE='1%20h'&QUANTITIES='1,9,20,23,24,29'
	//COMMAND:
	//	Sun (10)
	//	Mercury (199)水星
	//	Venus (299)金星
	//	Satellite (301)月球
	//	Earth (399)地球
	//	Phobos (401) 火卫一
	//	Deimos (402)火卫二
	//	Mars (499) 火星
	//	Jupiter (599)木星
	//	Saturn (699)土星
	//	Uranus (799)天王星
	//	Neptune (899)海王星
	//	Pluto (999) 冥王星
	//	1099 Figneria (1928 RQ)
	//	500 Selinur (A903 BJ)
	//水星：近日点约为0.31 AU，远日点约为0.47 AU。
	//金星：近日点约为0.72 AU，远日点约为0.72 AU。
	//地球：近日点和远日点都是1 AU。
	//火星：近日点约为1.38 AU，远日点约为1.66 AU。
	//木星：近日点约为4.95 AU，远日点约为5.46 AU。
	//土星：近日点约为9.03 AU，远日点约为9.54 AU。
	//天王星：近日点约为18.31 AU，远日点约为19.19 AU。
	//海王星：近日点约为29.81 AU，远日点约为30.36 AU。
	//冥王星：近日点约为29.66 AU，远日点约为49.31 AU。
	//	Date__(UT)__HR:MN：表示观测日期和时间，UTC代表协调世界时。
	//	R.A.__ (ICRF) __DEC：表示天体的赤经和赤纬，ICRF代表国际天球参考框架。
	//	APmag：表示天体的视星等。
	//	S-brt：表示天体的表面亮度。
	//	delta：表示天体相对于太阳的角度。
	//	deldot：表示天体的角速度。
	//	S-O-T /r：表示天体的视向速度。
	//	S-T-O：表示天体的自转周期。
	//	Cnst：可能表示天体的常数或其他相关信息。
}

var Stars = map[int]Celestial{
	10:  {Id: 10, Name: "Sun", NameCN: "日", AstrolabeR: 0, Mass: 1988500 * 1e24},
	199: {Id: 199, Name: "Mercury", NameCN: "水", AstrolabeR: 15, Mass: 3.302 * 1e23},
	299: {Id: 299, Name: "Venus", NameCN: "金", AstrolabeR: 30, Mass: 4.8685 * 1e24},
	399: {Id: 399, Name: "Earth", NameCN: "地", AstrolabeR: 50, Mass: 5.97219 * 1e24, //+-0.0006
		Satellite: []*Celestial{{Id: 301, Name: "Satellite", NameCN: "月", AstrolabeR: 10, Mass: 7.349 * 1e22}}},
	499: {Id: 499, Name: "Mars", NameCN: "火", AstrolabeR: 70, Mass: 6.4171 * 1e23,
		Satellite: []*Celestial{
			{Id: 401, Name: "Phobos", AstrolabeR: 5, Mass: 1.0659 * 1e16},
			{Id: 402, Name: "Deimos", AstrolabeR: 5, Mass: 1.4762 * 1e15},
		}},
	599: {Id: 599, Name: "Jupiter", NameCN: "木", AstrolabeR: 85, Mass: 1.8982 * 1e27, //18981.8722 +- .8817
		Satellite: []*Celestial{
			{Id: 501, Name: "Io", AstrolabeR: 5, Mass: 8.9319 * 1e22},
			{Id: 502, Name: "Europa", AstrolabeR: 5, Mass: 4.7998 * 1e22},
			{Id: 503, Name: "Ganymede", AstrolabeR: 6, Mass: 1.4819 * 1e23},
			{Id: 504, Name: "Callisto", AstrolabeR: 7, Mass: 1.0759 * 1e23},
		}},
	699: {Id: 699, Name: "Saturn", NameCN: "土", AstrolabeR: 100, Mass: 5.6834 * 1e26,
		Satellite: []*Celestial{
			{Id: 601, Name: "Mimas", AstrolabeR: 5, Mass: 3.7493 * 1e19},
			{Id: 602, Name: "Enceladus", AstrolabeR: 5, Mass: 1.0802 * 1e20},
			{Id: 603, Name: "Tethys", AstrolabeR: 6, Mass: 6.1745 * 1e20},
			{Id: 604, Name: "Dione", AstrolabeR: 6, Mass: 1.0955 * 1e21},
			{Id: 605, Name: "Rhea", AstrolabeR: 6, Mass: 2.306 * 1e21},
			{Id: 606, Name: "Titan", AstrolabeR: 6, Mass: 1.3455 * 1e23},
			{Id: 607, Name: "Iapetus", AstrolabeR: 7, Mass: 1.8053 * 1e21},
		}},
	799: {Id: 799, Name: "Uranus", NameCN: "天", AstrolabeR: 115, Mass: 8.6813 * 1e25},
	899: {Id: 899, Name: "Neptune", NameCN: "海", AstrolabeR: 130, Mass: 1.02409 * 1e26},
	999: {Id: 999, Name: "Pluto", NameCN: "冥", AstrolabeR: 145, Mass: 1.307 * 1e22}, //1.307+-0.018
}

func NewAstrolabe() *Astrolabe {
	return &Astrolabe{
		centerX:  770,
		centerY:  450,
		observer: 399,
	}
}

func (a *Astrolabe) Update() error {
	solar := uiQiMen.pan.Solar
	hour := solar.GetHour()
	minute := solar.GetMinute()
	//计算太阳位置 暂用时辰近似
	degreesS := 90 + (float64(hour)+float64(minute)/60)*15
	solarX, solarY := util.CalRadiansPos(float64(a.centerX), float64(a.centerY), float64(Stars[a.observer].AstrolabeR), degreesS)
	a.solarX, a.solarY = float32(solarX), float32(solarY)
	//计算月球位置 暂以农历近似
	lDay := uiQiMen.pan.Lunar.GetDay()
	degreesM := 90 + (float64(hour)+float64(minute)/60)*15 - float64(lDay)/28*360
	moonX, moonY := util.CalRadiansPos(float64(a.centerX), float64(a.centerY), float64(Stars[a.observer].Satellite[0].AstrolabeR), degreesM)
	a.moonX, a.moonY = float32(moonX), float32(moonY)
	return nil
}

func (a *Astrolabe) Draw(screen *ebiten.Image) {
	//ft := ui.GetDefaultUIFont()
	sX, sY := a.solarX, a.solarY
	cx, cy := a.centerX, a.centerY
	//外圈
	vector.StrokeCircle(screen, cx, cy, outCircleR, outCircleW, colorSkyGateCircle, true)
	vector.StrokeCircle(screen, cx, cy, outCircleR-outCircleW/2, outCircleW/2, colorPowerCircle, true)
	//十字线
	vector.StrokeLine(screen, cx-outCircleR, cy, cx+outCircleR, cy, 1, colorOrbits, true)
	vector.StrokeLine(screen, cx, cy-outCircleR, cx, cy+outCircleR, 1, colorOrbits, true)
	//画12宫
	for i := 1; i <= 12; i++ {
		angleDegrees := float64(i) * 30
		lx1, ly1 := util.CalRadiansPos(float64(cx), float64(cy), float64(outCircleR-outCircleW/2), angleDegrees)
		lx2, ly2 := util.CalRadiansPos(float64(cx), float64(cy), float64(outCircleR+outCircleW/2), angleDegrees)
		vector.StrokeLine(screen, float32(lx1), float32(ly1), float32(lx2), float32(ly2), 1, colorGongSplit, true)
	}

	for id, obj := range Stars {
		vector.StrokeCircle(screen, sX, sY, obj.AstrolabeR, 1, colorOrbits, true) //planet Orbit
		var px, py float32
		if a.observer == id {
			px, py = cx, cy
		} else {
			x, y := util.CalRadiansPos(float64(sX), float64(sY), float64(obj.AstrolabeR), float64(rand.Intn(360)))
			px, py = float32(x), float32(y)
		}
		vector.DrawFilledCircle(screen, px, py, 2, colorLeader, true) //planet
		for _, ob := range obj.Satellite {
			vector.StrokeCircle(screen, px, py, ob.AstrolabeR, 1, colorOrbits, true) //satellite Orbit
			if ob.Id == 301 {
				vector.DrawFilledCircle(screen, a.moonX, a.moonY, 1, colorLeader, true) //moon
			} else {
				mx, my := util.CalRadiansPos(float64(px), float64(py), float64(ob.AstrolabeR), float64(rand.Intn(360)))
				vector.DrawFilledCircle(screen, float32(mx), float32(my), 1, colorLeader, true) //satellite
			}
		}
	}
}

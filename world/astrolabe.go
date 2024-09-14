package world

import (
	"fmt"
	"github.com/6tail/lunar-go/calendar"
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
	"qimen/util"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	_ "xorm.io/core"
	"xorm.io/xorm"
)

const (
	G  = 6.67430 * 1e-11 //万有引力常数(Nm²/kg²)
	C  = 299792458       //光速(m/s)
	AU = 149597870.7e3   //天文单位(m)

	outCircleR = 240
	outCircleW = 24

	DataTimeMin      = "2006-01-02 15:04"
	DateTimeNASA     = "2006-Jan-02 15:04"
	NASADataFile     = "nasa_data.db"
	NASADataTimeLast = time.Hour * 24 //一次查询时间范围
	NASADataStepSize = "1h"           //星体数据步长 1h 1d 1m
)

var Constellation = []string{"Ari", "Tau", "Gem", "Can", "Leo", "Vir", "Lib", "Sco", "Sgr", "Cap", "Aqr", "Psc"}
var ConstellationS = []string{"羊", "牛", "双", "蟹", "狮", "室", "秤", "蝎", "射", "摩", "瓶", "鱼"}
var ConstellationCN = map[string]string{
	"Ari": "白羊座", "Tau": "金牛座", "Gem": "双子座", "Can": "巨蟹座",
	"Leo": "狮子座", "Vir": "室女座", "Lib": "天秤座", "Sco": "天蝎座",
	"Sgr": "射手座", "Cap": "摩羯座", "Aqr": "水瓶座", "Psc": "双鱼座",
}
var AstrolabeGong = []string{"", "命宫", "财帛", "交流", "田宅", "娱乐", "健康", "夫妻", "疾厄", "迁移", "事业", "福德", "玄秘"}

// 建星是太阳位置,星座是太阳上升,所以星座相当建星是指定地支时
// 月将 寅丑子亥戌酉申未午巳辰卯
// 建星 子丑寅卯辰巳午未申酉戌亥
// 星座 射摩瓶鱼羊牛双蟹狮室秤蝎
// 四象 火土风水火土风水火土风水

type Astrolabe struct {
	sync.RWMutex
	X, Y           float32
	solarX, solarY float32
	observer       int //观察者

	solar    calendar.Solar
	timezone string
	tzOffset int
	tzRA0    float64 //春分点角度

	DataQuerying bool
	Ephemeris    map[string]*ObserveData

	ConstellationLoc [12]gongLocation //星座位
	AstrolabeLoc     [12]gongLocation //宫位
}
type gongLocation struct {
	lx1, ly1, lx2, ly2 float32 //分割线
	x, y               int     //文字坐标
}

type CelestialBodySum struct {
	Id   int
	GMin float64
	GMax float64 //引力范围
}
type CelestialBody struct {
	Id               int
	Name             string
	NameCN           string
	SubBody          []int      //行星系子体
	Satellite        []int      //卫星
	Mass             float64    //质量 kg
	R                float32    //轨道半径 AU
	orbitCenter      int        //轨道中心
	Gravity          float64    //引力 N
	CelestialBodySum            //GMin, GMax  float64    //引力范围
	color            color.RGBA //RedShift红移/BlueShift蓝移

	AR               float32 //星盘半径
	DrawX, DrawY     float32 //星盘坐标
	SphereX, SphereY float32 //天球坐标
}

// DrawR 星盘半径
func (c *CelestialBody) DrawR() float32 {
	if c.R > 0 {
		return c.R * 150
	} else {
		return c.AR
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
	10:  {Id: 10, Name: "Sun", NameCN: "日", AR: 0, Mass: 1988500 * 1e24},
	199: {Id: 199, Name: "Mercury", NameCN: "水", R: (0.31 + 0.47) / 2, AR: 15, Mass: 3.302 * 1e23},
	299: {Id: 299, Name: "Venus", NameCN: "金", R: 0.72, AR: 30, Mass: 4.8685 * 1e24},
	301: {Id: 301, Name: "Satellite", NameCN: "月", AR: 10, Mass: 7.349 * 1e22, orbitCenter: 399},
	399: {Id: 399, Name: "Earth", NameCN: "地", R: 1, AR: 50, Mass: 5.97219 * 1e24, //+-0.0006
		Satellite: []int{301}},
	3: {Id: 3, Name: "EarthBarycenter", NameCN: "地月系重心", AR: 50,
		SubBody: []int{301, 399}},
	401: {Id: 401, Name: "Phobos", AR: 5, Mass: 1.0659 * 1e16, orbitCenter: 499},
	402: {Id: 402, Name: "Deimos", AR: 5, Mass: 1.4762 * 1e15, orbitCenter: 499},
	499: {Id: 499, Name: "Mars", NameCN: "火", R: (1.38 + 1.66) / 2, AR: 70, Mass: 6.4171 * 1e23,
		Satellite: []int{401, 402}},
	4: {Id: 4, Name: "MarsBarycenter", NameCN: "火系重心", AR: 70,
		SubBody: []int{401, 402, 499}},
	501: {Id: 501, Name: "Io", AR: 5, Mass: 8.9319 * 1e22, orbitCenter: 599},
	502: {Id: 502, Name: "Europa", AR: 5, Mass: 4.7998 * 1e22, orbitCenter: 599},
	503: {Id: 503, Name: "Ganymede", AR: 6, Mass: 1.4819 * 1e23, orbitCenter: 599},
	504: {Id: 504, Name: "Callisto", AR: 7, Mass: 1.0759 * 1e23, orbitCenter: 599},
	599: {Id: 599, Name: "Jupiter", NameCN: "木", R: (4.95 + 5.46) / 2, AR: 90, Mass: 1.8982 * 1e27, //18981.8722 +- .8817
		Satellite: []int{501, 502, 503, 504}},
	601: {Id: 601, Name: "Mimas", AR: 5, Mass: 3.7493 * 1e19, orbitCenter: 699},
	602: {Id: 602, Name: "Enceladus", AR: 5, Mass: 1.0802 * 1e20, orbitCenter: 699},
	603: {Id: 603, Name: "Tethys", AR: 6, Mass: 6.1745 * 1e20, orbitCenter: 699},
	604: {Id: 604, Name: "Dione", AR: 6, Mass: 1.0955 * 1e21, orbitCenter: 699},
	605: {Id: 605, Name: "Rhea", AR: 6, Mass: 2.306 * 1e21, orbitCenter: 699},
	606: {Id: 606, Name: "Titan", AR: 6, Mass: 1.3455 * 1e23, orbitCenter: 699},
	607: {Id: 607, Name: "Iapetus", AR: 7, Mass: 1.8053 * 1e21, orbitCenter: 699},
	699: {Id: 699, Name: "Saturn", NameCN: "土", R: (9.03 + 9.54) / 2, AR: 105, Mass: 5.6834 * 1e26,
		Satellite: []int{601, 602, 603, 604, 605, 606, 607}},
	799: {Id: 799, Name: "Uranus", NameCN: "天", R: (18.31 + 19.19) / 2, AR: 120, Mass: 8.6813 * 1e25},
	899: {Id: 899, Name: "Neptune", NameCN: "海", R: (29.81 + 30.36) / 2, AR: 135, Mass: 1.02409 * 1e26},
	901: {Id: 901, Name: "Charon", NameCN: "卡", AR: 5, Mass: 1 * 1e22, orbitCenter: 999},
	999: {Id: 999, Name: "Pluto", NameCN: "冥", R: (29.66 + 49.31) / 2, AR: 150, Mass: 1.307 * 1e22, //1.307+-0.018
		Satellite: []int{901}},
	9: {Id: 9, Name: "PlutoBarycenter", NameCN: "冥王系重心", AR: 150,
		SubBody: []int{901, 999}},
}
var Draws = []int{
	10, 199, 299, 399, 499, 599, 699, // 799, 899, 999,
}

func NewAstrolabe(x, y float32) *Astrolabe {
	tz, offset := time.Now().Local().Zone()
	a := &Astrolabe{
		X: x, Y: y,
		observer:  399,
		timezone:  tz,
		tzOffset:  offset,
		Ephemeris: make(map[string]*ObserveData),
	}
	for i := 1; i <= 12; i++ {
		//固定宫位
		degrees := float64(i)*30 - 90
		ly1, lx1 := util.CalRadiansPos(float64(y), float64(x), float64(outCircleR-outCircleW/4), degrees)
		ly2, lx2 := util.CalRadiansPos(float64(y), float64(x), float64(outCircleR-outCircleW*3/4), degrees)
		y, x := util.CalRadiansPos(float64(y), float64(x), float64(outCircleR-outCircleW/2), degrees-15)
		a.AstrolabeLoc[i-1] = gongLocation{float32(lx1), float32(ly1), float32(lx2), float32(ly2), int(x), int(y)}
	}
	return a
}
func (a *Astrolabe) SetPos(x, y float32) {
	a.X, a.Y = x, y
}
func (a *Astrolabe) Update() {
	sCal := ThisGame.qmGame.Solar
	if a.solar == *sCal {
		return
	}
	a.solar = *sCal
	hour := sCal.GetHour()
	minute := sCal.GetMinute()

	cx, cy := a.X, a.Y
	//计算太阳位置 以时间计角度
	degreesS := 360 - (float64(hour)+float64(minute)/60)*15 //本地时区太阳角度 0~360 0时0度
	degreesSO := 360 - degreesS
	solarY, solarX := util.CalRadiansPos(float64(a.Y), float64(a.X), float64(Bodies[a.observer].DrawR()), degreesS)
	a.solarX, a.solarY = float32(solarX), float32(solarY)

	//计算月球位置 暂以农历近似
	lDay := ThisGame.qmGame.Lunar.GetDay()
	degreesM := -(float64(hour)+float64(minute)/60)*15 + float64(lDay)/29*360
	{
		moon := Bodies[301]
		moonY, moonX := util.CalRadiansPos(float64(a.Y), float64(a.X), float64(moon.DrawR()), degreesM)
		moon.DrawX, moon.DrawY = float32(moonX), float32(moonY)
		v1 := (&util.Vec2[float32]{X: moon.DrawX - cx, Y: moon.DrawY - cy}).ScaledToLength(outCircleR - outCircleW*2)
		moon.SphereX = cx + v1.X
		moon.SphereY = cy + v1.Y
	}

	m := Bodies[a.observer].Mass
	for _, id := range Draws {
		body := Bodies[id]
		if id == a.observer {
			body.DrawX, body.DrawY = a.X, a.Y
			body.color = colorBlueShift
			for _, sid := range body.Satellite {
				bd := Bodies[sid]
				oe := a.GetEphemeris(sid, sCal)
				if oe == nil {
					a.solar = calendar.Solar{}
					continue
				}
				bd.Gravity = G * bd.Mass * m / math.Pow(oe.Delta*AU, 2)
				if bd.GMin == 0 {
					bd.GMin = bd.Gravity
				} else {
					bd.GMin = min(bd.Gravity, bd.GMin)
				}
				bd.GMax = max(bd.Gravity, bd.GMax)
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
			a.solar = calendar.Solar{}
			continue
		}
		if id == 10 {
			body.DrawX, body.DrawY = a.solarX, a.solarY
			v1 := (&util.Vec2[float32]{X: body.DrawX - cx, Y: body.DrawY - cy}).ScaledToLength(outCircleR - outCircleW*3)
			body.SphereX = cx + v1.X
			body.SphereY = cy + v1.Y

			a.tzRA0 = degreesS - oe.RARadius()
		} else {
			ost := 180 - oe.STO - oe.SOT
			var degrees float64
			switch oe.SOTR {
			case "/T": //TRAILS 目标跟踪太阳
				degrees = 90 - 270 - degreesSO - ost
			case "/L": //LEADS 目标引领太阳
				degrees = 90 - 270 - degreesSO + ost
			}
			y, x := util.CalRadiansPos(solarY, solarX, float64(body.DrawR()), degrees)
			body.DrawX = float32(x)
			body.DrawY = float32(y)
			v1 := (&util.Vec2[float32]{X: body.DrawX - cx, Y: body.DrawY - cy}).ScaledToLength(outCircleR - outCircleW*3)
			body.SphereX = cx + v1.X
			body.SphereY = cy + v1.Y
		}
		body.Gravity = G * body.Mass * m / math.Pow(oe.Delta*AU, 2)
		if body.GMin == 0 {
			body.GMin = body.Gravity
		} else {
			body.GMin = min(body.Gravity, body.GMin)
		}
		body.GMax = max(body.Gravity, body.GMax)
		if oe.Deldot > 0 {
			body.color = colorBlueShift
		} else {
			body.color = colorRedShift
		}
	}
	a.calGongLocation()
	return
}

func (a *Astrolabe) calGongLocation() {
	cx, cy := a.X, a.Y
	for i := 1; i <= 12; i++ {
		degrees := a.tzRA0 + float64(i-1)*30
		ly1, lx1 := util.CalRadiansPos(float64(cy), float64(cx), float64(outCircleR-outCircleW/4), degrees)
		ly2, lx2 := util.CalRadiansPos(float64(cy), float64(cx), float64(outCircleR+outCircleW/2), degrees)
		y, x := util.CalRadiansPos(float64(cy), float64(cx), float64(outCircleR+outCircleW/4), degrees+15)
		a.ConstellationLoc[i-1] = gongLocation{float32(lx1), float32(ly1), float32(lx2), float32(ly2), int(x), int(y)}
	}
}

func (a *Astrolabe) Draw(dst *ebiten.Image) {
	ft, _ := GetFontFace(12)
	cx, cy := a.X, a.Y
	sX, sY := a.solarX, a.solarY
	//外圈
	vector.StrokeCircle(dst, cx, cy, outCircleR, outCircleW, colorSkyGateCircle, true)                   //星座
	vector.StrokeCircle(dst, cx, cy, outCircleR-outCircleW/2, outCircleW/2, colorGroundGateCircle, true) //宫位
	vector.StrokeCircle(dst, cx, cy, outCircleR-outCircleW, outCircleW/2, colorPowerCircle, true)        //天球celestial sphere
	//十字线
	horizons := float32(0) //TODO 地平线调整
	vector.StrokeLine(dst, cx-outCircleR, cy-horizons, cx+outCircleR, cy-horizons, 1, colorCross, true)
	vector.StrokeLine(dst, cx, cy-outCircleR, cx, cy+outCircleR, 1, colorCross, true)
	//春分点 RA0
	tzY, tzX := util.CalRadiansPos(float64(cy), float64(cx), outCircleR/2, a.tzRA0)
	vector.StrokeLine(dst, cx, cy, float32(tzX), float32(tzY), 1, colorOrbits, true)
	text.Draw(dst, "春分", ft, int(tzX), int(tzY), colorLeader)                                                        //春分点
	text.Draw(dst, fmt.Sprintf("%s月", ThisGame.qmGame.Lunar.GetYueXiang()), ft, int(cx-16), int(cy-25), colorLeader) //月相

	//画12宫
	for i := 0; i < 12; i++ {
		l := a.ConstellationLoc[i]
		vector.StrokeLine(dst, l.lx1, l.ly1, l.lx2, l.ly2, 1, colorGongSplit, true)        //星宫
		text.Draw(dst, fmt.Sprintf("%s", ConstellationS[i]), ft, l.x-6, l.y+6, colorJiang) //星座
		l = a.AstrolabeLoc[i]
		vector.StrokeLine(dst, l.lx1, l.ly1, l.lx2, l.ly2, 1, colorGongSplit, true) //宫
		text.Draw(dst, fmt.Sprintf("%d", i+1), ft, l.x-4, l.y+4, colorJiang)        //宫位
	}
	//画星体
	for _, id := range Draws {
		obj := Bodies[id]
		vector.StrokeCircle(dst, sX, sY, obj.DrawR(), 1, colorOrbits, true) //planet Orbit
		if a.observer == obj.Id {
			//地球观察者
		} else {
			if obj.SphereX == 0 && obj.SphereY == 0 {
				continue //查询中
			}
			vector.StrokeLine(dst, cx, cy, obj.SphereX, obj.SphereY, 1, colorOrbits, true) // sphere line
			//vector.DrawFilledCircle(dst, obj.SphereX, obj.SphereY, 2, obj.color, true)     // sphere
			text.Draw(dst, obj.NameCN, ft, int(obj.SphereX), int(obj.SphereY), obj.color)
		}
		vector.DrawFilledCircle(dst, obj.DrawX, obj.DrawY, 2, obj.color, true) //planet

		for _, sid := range obj.Satellite {
			ob := Bodies[sid]
			if ob.Id == 301 {
				vector.StrokeCircle(dst, obj.DrawX, obj.DrawY, ob.DrawR(), 1, colorOrbits, true) //satellite Orbit
				vector.StrokeLine(dst, cx, cy, ob.SphereX, ob.SphereY, 1, colorOrbits, true)     // sphere line
				vector.DrawFilledCircle(dst, ob.DrawX, ob.DrawY, 1, ob.color, true)              //moon
				text.Draw(dst, ob.NameCN, ft, int(ob.SphereX), int(ob.SphereY), ob.color)
			} else {
				mx, my := util.CalRadiansPos(float64(obj.DrawX), float64(obj.DrawY), float64(ob.DrawR()), float64(rand.Intn(360)))
				vector.DrawFilledCircle(dst, float32(mx), float32(my), 1, obj.color, true) //satellite
			}
		}
	}
	a.DrawGravity(dst)
	if a.DataQuerying {
		text.Draw(dst, "正在观星..", ft, int(cx-32), int(cy-10), color.White)
	}
}
func (a *Astrolabe) DrawGravity(dst *ebiten.Image) {
	cx, cy := a.X, a.Y-440
	for _, id := range Draws {
		obj := Bodies[id]
		if obj.Gravity > 0 {
			//cx += 40
			cy += 22
			DrawRangeBar(dst, cx-200, cy, 100, obj.NameCN, obj.Gravity, obj.GMin, obj.GMax, obj.color)
		}
		for _, sid := range obj.Satellite {
			ob := Bodies[sid]
			if ob.Gravity > 0 {
				//cx += 40
				cy += 22
				DrawRangeBar(dst, cx-200, cy, 100, ob.NameCN, ob.Gravity, ob.GMin, ob.GMax, ob.color)
			}
		}
	}
}

func (a *Astrolabe) GetEphemeris(tid int, s *calendar.Solar) *ObserveData {
	if s.GetYear() < 1600 || s.GetYear() >= 2500 {
		return nil
	}
	t := time.Date(s.GetYear(), time.Month(s.GetMonth()), s.GetDay(), s.GetHour(),
		s.GetMinute(), 0, 0, time.Local)
	sts := t.In(time.UTC).Format(DataTimeMin)
	id := fmt.Sprintf("%d_%s", tid, sts)
	a.RWMutex.RLock()
	defer a.RWMutex.RUnlock()
	v := a.Ephemeris[id]
	if v != nil {
		return v
	}
	if a.DataQuerying {
		return nil
	}
	var it ObserveData
	has, err := db.Where("Id = ?", id).Get(&it)
	if err != nil {
		UIShowMsgBox(err.Error(), "确定", "确定", nil, nil)
		return nil
	}
	if has {
		return &it
	} else {
		a.DataQuerying = true
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
		a.DataQuerying = false
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
		UIShowMsgBox(err.Error(), "确定", "确定", nil, nil)
	}
	a.DataQuerying = false
}
func (a *Astrolabe) GetSolarPos() (float32, float32) {
	return a.solarX, a.solarY
}
func (a *Astrolabe) GetMoonPos() (float32, float32) {
	return Bodies[301].DrawX, Bodies[301].DrawY
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
	SOT    float64 // S-O-T 观察者-目标-太阳角度
	SOTR   string  // "/T": TRAILS 目标跟踪S  "/L": LEADS 目标引领S
	STO    float64 // S-T-O 太阳-目标-观察者角度
	Cnst   string  // 星座
}

// RA 赤经 05 29 58.88
func (c *ObserveData) RA() string {
	return c.RA_DEC[0:11]
}

// RARadius 赤经角度
func (c *ObserveData) RARadius() float64 {
	ra := c.RA() //05 29 58.88
	ss := strings.Split(ra, " ")
	h, _ := strconv.ParseFloat(ss[0], 64)
	m, _ := strconv.ParseFloat(ss[1], 64)
	s, _ := strconv.ParseFloat(ss[2], 64)
	return h*15 + m/4 + s/240
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
func (c *observeDataSrc) SOT() float64 {
	v, _ := strconv.ParseFloat(c.data["S-O-T"], 64)
	return v
}
func (c *observeDataSrc) SOTR() string {
	return c.data["/r"]
}
func (c *observeDataSrc) STO() float64 {
	v, _ := strconv.ParseFloat(c.data["S-T-O"], 64)
	return v
}

var db *xorm.Engine

func init() {
	var err error
	//db, err = xorm.NewEngine("sqlite3", "file::memory:?cache=shared")
	db, err = xorm.NewEngine("sqlite3", NASADataFile)
	if err != nil {
		panic(err)
	}

	// 同步模型到数据库
	err = db.Sync2(new(ObserveData))
	if err != nil {
		panic(err)
	}
}

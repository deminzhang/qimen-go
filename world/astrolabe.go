package world

import (
	"fmt"
	"github.com/6tail/lunar-go/calendar"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	//_ "github.com/mattn/go-sqlite3"
	"image/color"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"qimen/ui"
	"qimen/util"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	G  = 6.67430 * 1e-11 //万有引力常数(Nm²/kg²)
	C  = 299792458       //光速(m/s)
	AU = 149597870.7e3   //天文单位(m)

	outCircleR = 240
	outCircleW = 24

	DataTimeMin = "2006-01-02 15:04"
	DateTimeOE  = "2006-Jan-02 15:04"
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
	X, Y           float32
	solarX, solarY float32
	observer       int //观察者

	solar    calendar.Solar
	timezone string
	tzOffset int
	tzRA0    float64 //春分点角度

	DataQuery bool
	OEData    map[int]*ObserveEphemeris

	ConstellationLoc [12]gongLocation //星座位
	AstrolabeLoc     [12]gongLocation //宫位
}
type gongLocation struct {
	lx1, ly1, lx2, ly2 float32 //分割线
	x, y               int     //文字坐标
}

type CelestialBody struct {
	Id          int
	Name        string
	NameCN      string
	SubBody     []int      //行星系子体
	Satellite   []int      //卫星
	Mass        float64    //质量 kg
	R           float32    //轨道半径 AU
	orbitCenter int        //轨道中心
	Gravity     float64    //引力 N
	color       color.RGBA //RedShift红移/BlueShift蓝移

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
	401: {Id: 401, Name: "Phobos", AR: 5, Mass: 1.0659 * 1e16},
	402: {Id: 402, Name: "Deimos", AR: 5, Mass: 1.4762 * 1e15},
	499: {Id: 499, Name: "Mars", NameCN: "火", R: (1.38 + 1.66) / 2, AR: 70, Mass: 6.4171 * 1e23,
		Satellite: []int{401, 402}},
	4: {Id: 4, Name: "MarsBarycenter", NameCN: "火系重心", AR: 70,
		SubBody: []int{401, 402, 499}},
	501: {Id: 501, Name: "Io", AR: 5, Mass: 8.9319 * 1e22},
	502: {Id: 502, Name: "Europa", AR: 5, Mass: 4.7998 * 1e22},
	503: {Id: 503, Name: "Ganymede", AR: 6, Mass: 1.4819 * 1e23},
	504: {Id: 504, Name: "Callisto", AR: 7, Mass: 1.0759 * 1e23},
	599: {Id: 599, Name: "Jupiter", NameCN: "木", R: (4.95 + 5.46) / 2, AR: 90, Mass: 1.8982 * 1e27, //18981.8722 +- .8817
		Satellite: []int{501, 502, 503, 504}},
	601: {Id: 601, Name: "Mimas", AR: 5, Mass: 3.7493 * 1e19},
	602: {Id: 602, Name: "Enceladus", AR: 5, Mass: 1.0802 * 1e20},
	603: {Id: 603, Name: "Tethys", AR: 6, Mass: 6.1745 * 1e20},
	604: {Id: 604, Name: "Dione", AR: 6, Mass: 1.0955 * 1e21},
	605: {Id: 605, Name: "Rhea", AR: 6, Mass: 2.306 * 1e21},
	606: {Id: 606, Name: "Titan", AR: 6, Mass: 1.3455 * 1e23},
	607: {Id: 607, Name: "Iapetus", AR: 7, Mass: 1.8053 * 1e21},
	699: {Id: 699, Name: "Saturn", NameCN: "土", R: (9.03 + 9.54) / 2, AR: 105, Mass: 5.6834 * 1e26,
		Satellite: []int{601, 602, 603, 604, 605, 606, 607}},
	799: {Id: 799, Name: "Uranus", NameCN: "天", R: (18.31 + 19.19) / 2, AR: 120, Mass: 8.6813 * 1e25},
	899: {Id: 899, Name: "Neptune", NameCN: "海", R: (29.81 + 30.36) / 2, AR: 135, Mass: 1.02409 * 1e26},
	901: {Id: 901, Name: "Charon", NameCN: "卡", AR: 5, Mass: 1 * 1e22},
	999: {Id: 999, Name: "Pluto", NameCN: "冥", R: (29.66 + 49.31) / 2, AR: 150, Mass: 1.307 * 1e22, //1.307+-0.018
		Satellite: []int{901}},
	9: {Id: 9, Name: "PlutoBarycenter", NameCN: "冥王系重心", AR: 150,
		SubBody: []int{901, 999}},
}
var Draws = []int{
	10, 199, 299, 399, 499, 599, 699, // 799, 899, 999,
}

func NewAstrolabe(cx, cy float32) *Astrolabe {
	tz, offset := time.Now().Local().Zone()
	a := &Astrolabe{
		X: cx, Y: cy,
		observer: 399,
		timezone: tz,
		tzOffset: offset,
		OEData:   make(map[int]*ObserveEphemeris),
	}
	for i := 1; i <= 12; i++ {
		//固定宫位
		degrees := float64(i)*30 - 90
		ly1, lx1 := util.CalRadiansPos(float64(cy), float64(cx), float64(outCircleR-outCircleW/4), degrees)
		ly2, lx2 := util.CalRadiansPos(float64(cy), float64(cx), float64(outCircleR-outCircleW*3/4), degrees)
		y, x := util.CalRadiansPos(float64(cy), float64(cx), float64(outCircleR-outCircleW/2), degrees-15)
		a.AstrolabeLoc[i-1] = gongLocation{float32(lx1), float32(ly1), float32(lx2), float32(ly2), int(x), int(y)}
	}
	return a
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
		moon.NameCN = fmt.Sprintf("(%s)月", ThisGame.qmGame.Lunar.GetYueXiang())
	}

	m := Bodies[a.observer].Mass
	for _, id := range Draws {
		body := Bodies[id]
		if id == a.observer {
			body.DrawX, body.DrawY = a.X, a.Y
			body.color = colorBlueShift
			for _, id := range body.Satellite {
				body := Bodies[id]
				oe := a.GetEphemeris(id, sCal)
				if oe == nil {
					a.solar = calendar.Solar{}
					continue
				}
				body.Gravity = G * body.Mass * m / math.Pow(oe.Delta()*AU, 2)
				if oe.Deldot() > 0 {
					body.color = colorBlueShift
				} else {
					body.color = colorRedShift
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
			ost := 180 - oe.STO() - oe.SOT()
			var degrees float64
			switch oe.SOTR() {
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
		body.Gravity = G * body.Mass * m / math.Pow(oe.Delta()*AU, 2)
		if oe.Deldot() > 0 {
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
	ft := ui.GetDefaultUIFont()
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
	text.Draw(dst, "春分", ft, int(tzX), int(tzY), colorLeader) //春分点

	//画12宫
	for i := 0; i < 12; i++ {
		l := a.ConstellationLoc[i]
		vector.StrokeLine(dst, l.lx1, l.ly1, l.lx2, l.ly2, 1, colorGongSplit, true)        //星宫
		text.Draw(dst, fmt.Sprintf("%s", ConstellationS[i]), ft, l.x-6, l.y+6, colorJiang) //星座
		l = a.AstrolabeLoc[i]
		vector.StrokeLine(dst, l.lx1, l.ly1, l.lx2, l.ly2, 1, colorGongSplit, true) //宫
		text.Draw(dst, fmt.Sprintf("%d", i+1), ft, l.x-4, l.y+4, colorJiang)        //宫位
	}

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
	var desX, desY float32 = 1154, 270
	for _, id := range Draws {
		obj := Bodies[id]
		if obj.Gravity > 0 {
			desY += 20
			text.Draw(dst, fmt.Sprintf("%s:G:%.2fN", obj.NameCN, obj.Gravity), ft, int(desX), int(desY), color.White)
		}
		for _, sid := range obj.Satellite {
			ob := Bodies[sid]
			if ob.Gravity > 0 {
				desY += 20
				text.Draw(dst, fmt.Sprintf("%s:G:%.2f", ob.NameCN, ob.Gravity), ft, int(desX), int(desY), color.White)
			}
		}
	}
	if a.DataQuery {
		text.Draw(dst, "正在观星..", ft, int(cx-32), int(cy-10), color.White)
	}
}

func (a *Astrolabe) GetNASAData(id int, sts, ets string) *ObserveEphemeris {
	urls := fmt.Sprintf("https://ssd.jpl.nasa.gov/api/horizons.api?"+
		"format=text&COMMAND='%d'&OBJ_DATA='YES'&MAKE_EPHEM='YES'&EPHEM_TYPE='OBSERVER'&CENTER='500@399'"+
		"&START_TIME='%s'&STOP_TIME='%s'&STEP_SIZE='%s'&QUANTITIES='1,20,23,24,29'",
		id, sts, ets, url.QueryEscape("1h"))
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
	var lines = make(map[string]*ObserveData, len(lines0))
	for _, line := range lines0 {
		var m = map[string]string{}
		var preL int
		for i, l := range colLen {
			v0 := line[preL : preL+l]
			m[head[i]] = strings.TrimSpace(v0)
			preL += l
		}
		date := m["Date__(UT)__HR:MN"]
		dt, _ := time.Parse(DateTimeOE, date)
		ymdhm := fmt.Sprintf("%04d-%02d-%02d %02d:%02d", dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute())
		lines[ymdhm] = &ObserveData{data: m}
	}
	defer resp.Body.Close()
	return &ObserveEphemeris{
		id:   id,
		data: lines,
	}
}

func (a *Astrolabe) GetEphemeris(id int, dt *calendar.Solar) *ObserveData {
	d := a.OEData[id]
	if d == nil {
		d = &ObserveEphemeris{id: id, data: make(map[string]*ObserveData)}
		a.OEData[id] = d
	}
	tl := time.Date(dt.GetYear(), time.Month(dt.GetMonth()), dt.GetDay(), dt.GetHour(), dt.GetMinute(), 0, 0, time.Local)
	st := tl.In(time.UTC)
	sts := st.Format(DataTimeMin)
	oe := d.GetOE(sts)
	if oe == nil {
		//get from db
		//if dbData {
		//}
		if a.DataQuery {
			return nil
		}
		a.DataQuery = true
		te := st.Add(time.Hour * 24)
		ets := te.Format(DataTimeMin)
		go func() {
			d = a.GetNASAData(id, sts, ets)
			if d == nil {
				return
			}
			a.OEData[id].ApplyData(d.data)
			a.DataQuery = false
		}()
	}
	return d.GetOE(sts)
}

type ObserveData struct {
	data map[string]string
}

// RA 赤经 05 29 58.88
func (c *ObserveData) RA() string {
	return c.data["R.A._____(ICRF)_____DEC"][0:11]
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
	return c.data["R.A._____(ICRF)_____DEC"][12:] //+23 15 24.3
}

// Delta 距离 AU
func (c *ObserveData) Delta() float64 {
	v, _ := strconv.ParseFloat(c.data["delta"], 64)
	return v
}

// Deldot delta-dot 距离变化 KM/S 为正表示远离观察者
func (c *ObserveData) Deldot() float64 {
	v, _ := strconv.ParseFloat(c.data["deldot"], 64)
	return v
}
func (c *ObserveData) SOT() float64 {
	v, _ := strconv.ParseFloat(c.data["S-O-T"], 64)
	return v
}

// SOTR  "/T": TRAILS 目标跟踪S  "/L": LEADS 目标引领S
func (c *ObserveData) SOTR() string {
	return c.data["/r"]
}
func (c *ObserveData) STO() float64 {
	v, _ := strconv.ParseFloat(c.data["S-T-O"], 64)
	return v
}
func (c *ObserveData) Cnst() string {
	return c.data["Cnst"] //,Constellation[c.data["Cnst"]]
}

type ObserveEphemeris struct {
	sync.Mutex
	id   int
	data map[string]*ObserveData
}

func (c *ObserveEphemeris) Id() int {
	return c.id
}
func (c *ObserveEphemeris) GetOE(dataTime string) *ObserveData {
	c.Lock()
	defer c.Unlock()
	return c.data[dataTime]
}

func (c *ObserveEphemeris) ApplyData(datas map[string]*ObserveData) {
	c.Lock()
	defer c.Unlock()
	for s, data := range datas {
		c.data[s] = data
	}
}

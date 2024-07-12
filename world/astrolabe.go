package world

import (
	"fmt"
	"github.com/6tail/lunar-go/calendar"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"qimen/util"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	G = 6.67430 * 1e-11 //万有引力常数(Nm²/kg²)
	C = 299792458       //光速(m/s)

	outCircleR = 240
	outCircleW = 24

	DataTimeMin = "2006-01-02 15:04"
	DateTimeOE  = "2006-Jan-02 15:04"
)

type Astrolabe struct {
	centerX, centerY float32
	solarX, solarY   float32
	moonX, moonY     float32
	observer         int //观察者

	OEData map[int]*ObserveEphemeris
}

type CelestialBody struct {
	Id        int
	Name      string
	NameCN    string
	SubBody   []int   //行星系子体
	Satellite []int   //卫星
	Mass      float64 //质量 kg

	DrawR float32 //星盘半径
	DrawX float32 //星盘坐标
	DrawY float32 //星盘坐标
}

var Bodies = map[int]CelestialBody{
	10:  {Id: 10, Name: "Sun", NameCN: "日", DrawR: 0, Mass: 1988500 * 1e24},
	199: {Id: 199, Name: "Mercury", NameCN: "水", DrawR: 15, Mass: 3.302 * 1e23},
	299: {Id: 299, Name: "Venus", NameCN: "金", DrawR: 30, Mass: 4.8685 * 1e24},
	301: {Id: 301, Name: "Satellite", NameCN: "月", DrawR: 10, Mass: 7.349 * 1e22},
	399: {Id: 399, Name: "Earth", NameCN: "地", DrawR: 50, Mass: 5.97219 * 1e24, //+-0.0006
		Satellite: []int{301}},
	3: {Id: 3, Name: "EarthBarycenter", NameCN: "地月系重心", DrawR: 50,
		SubBody: []int{301, 399}},
	401: {Id: 401, Name: "Phobos", DrawR: 5, Mass: 1.0659 * 1e16},
	402: {Id: 402, Name: "Deimos", DrawR: 5, Mass: 1.4762 * 1e15},
	499: {Id: 499, Name: "Mars", NameCN: "火", DrawR: 70, Mass: 6.4171 * 1e23,
		Satellite: []int{401, 402}},
	4: {Id: 4, Name: "MarsBarycenter", NameCN: "火系重心", DrawR: 70,
		SubBody: []int{401, 402, 499}},
	501: {Id: 501, Name: "Io", DrawR: 5, Mass: 8.9319 * 1e22},
	502: {Id: 502, Name: "Europa", DrawR: 5, Mass: 4.7998 * 1e22},
	503: {Id: 503, Name: "Ganymede", DrawR: 6, Mass: 1.4819 * 1e23},
	504: {Id: 504, Name: "Callisto", DrawR: 7, Mass: 1.0759 * 1e23},
	599: {Id: 599, Name: "Jupiter", NameCN: "木", DrawR: 90, Mass: 1.8982 * 1e27, //18981.8722 +- .8817
		Satellite: []int{501, 502, 503, 504}},
	601: {Id: 601, Name: "Mimas", DrawR: 5, Mass: 3.7493 * 1e19},
	602: {Id: 602, Name: "Enceladus", DrawR: 5, Mass: 1.0802 * 1e20},
	603: {Id: 603, Name: "Tethys", DrawR: 6, Mass: 6.1745 * 1e20},
	604: {Id: 604, Name: "Dione", DrawR: 6, Mass: 1.0955 * 1e21},
	605: {Id: 605, Name: "Rhea", DrawR: 6, Mass: 2.306 * 1e21},
	606: {Id: 606, Name: "Titan", DrawR: 6, Mass: 1.3455 * 1e23},
	607: {Id: 607, Name: "Iapetus", DrawR: 7, Mass: 1.8053 * 1e21},
	699: {Id: 699, Name: "Saturn", NameCN: "土", DrawR: 105, Mass: 5.6834 * 1e26,
		Satellite: []int{601, 602, 603, 604, 605, 606, 607}},
	799: {Id: 799, Name: "Uranus", NameCN: "天", DrawR: 120, Mass: 8.6813 * 1e25},
	899: {Id: 899, Name: "Neptune", NameCN: "海", DrawR: 135, Mass: 1.02409 * 1e26},
	901: {Id: 901, Name: "Charon", NameCN: "卡", DrawR: 5, Mass: 1 * 1e22},
	999: {Id: 999, Name: "Pluto", NameCN: "冥", DrawR: 150, Mass: 1.307 * 1e22, //1.307+-0.018
		Satellite: []int{901}},
	9: {Id: 9, Name: "PlutoBarycenter", NameCN: "冥王系重心", DrawR: 150,
		SubBody: []int{901, 999}},
}
var CelestialsDraw = []int{
	10, 199, 299, 399, 499, 599, 699, 799, 899, 999,
	//水星：近日点约为0.31 AU，远日点约为0.47 AU。
	//金星：近日点约为0.72 AU，远日点约为0.72 AU。
	//地球：近日点和远日点都是1 AU。
	//火星：近日点约为1.38 AU，远日点约为1.66 AU。
	//木星：近日点约为4.95 AU，远日点约为5.46 AU。
	//土星：近日点约为9.03 AU，远日点约为9.54 AU。
	//天王星：近日点约为18.31 AU，远日点约为19.19 AU。
	//海王星：近日点约为29.81 AU，远日点约为30.36 AU。
	//冥王星：近日点约为29.66 AU，远日点约为49.31 AU。
}

func NewAstrolabe() *Astrolabe {
	tz, offset := time.Now().Local().Zone()
	fmt.Println(tz, offset)
	return &Astrolabe{
		centerX:  770,
		centerY:  450,
		observer: 399,
		OEData:   make(map[int]*ObserveEphemeris),
	}
}

func (a *Astrolabe) Update() error {
	solar := uiQiMen.pan.Solar
	hour := solar.GetHour()
	minute := solar.GetMinute()

	for _, id := range CelestialsDraw {
		if id != a.observer {
			oe := a.GetEphemeris(id, solar)
			fmt.Printf("RA[%s]\n", oe.RA())
			fmt.Printf("RAR[%f]\n", oe.RARadius())
			fmt.Printf("DEC[%s]\n", oe.DEC())
			fmt.Printf("Delta[%f]\n", oe.Delta())
			fmt.Println(oe.Deldot())
			fmt.Println(oe.SOT())
			fmt.Println(oe.SOTR())
			fmt.Println(oe.Cnst())
		}
	}

	//计算太阳位置 暂用时辰近似
	degreesS := 90 + (float64(hour)+float64(minute)/60)*15
	solarX, solarY := util.CalRadiansPos(float64(a.centerX), float64(a.centerY), float64(Bodies[a.observer].DrawR), degreesS)
	a.solarX, a.solarY = float32(solarX), float32(solarY)
	//计算月球位置 暂以农历近似
	lDay := uiQiMen.pan.Lunar.GetDay()
	degreesM := 90 + (float64(hour)+float64(minute)/60)*15 - float64(lDay)/28*360
	moonX, moonY := util.CalRadiansPos(float64(a.centerX), float64(a.centerY),
		float64(Bodies[Bodies[a.observer].Satellite[0]].DrawR), degreesM)
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

	for _, id := range CelestialsDraw {
		obj := Bodies[id]
		vector.StrokeCircle(screen, sX, sY, obj.DrawR, 1, colorOrbits, true) //planet Orbit
		var px, py float32
		if a.observer == obj.Id {
			px, py = cx, cy
		} else {
			x, y := util.CalRadiansPos(float64(sX), float64(sY), float64(obj.DrawR), float64(rand.Intn(360)))
			px, py = float32(x), float32(y)
		}
		vector.DrawFilledCircle(screen, px, py, 2, colorLeader, true) //planet
		for _, sid := range obj.Satellite {
			ob := Bodies[sid]
			vector.StrokeCircle(screen, px, py, ob.DrawR, 1, colorOrbits, true) //satellite Orbit
			if ob.Id == 301 {
				vector.DrawFilledCircle(screen, a.moonX, a.moonY, 1, colorLeader, true) //moon
			} else {
				mx, my := util.CalRadiansPos(float64(px), float64(py), float64(ob.DrawR), float64(rand.Intn(360)))
				vector.DrawFilledCircle(screen, float32(mx), float32(my), 1, colorLeader, true) //satellite
			}
		}
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
		log.Fatalf("Error sending GET request: %v", err)
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected Status Code: %d", resp.StatusCode)
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
	fmt.Printf(head0)
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
	var lines = make(map[string]ObserveData, len(lines0))
	for _, line := range lines0 {
		fmt.Println(line)
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
		lines[ymdhm] = ObserveData{data: m}
	}
	defer resp.Body.Close()
	return &ObserveEphemeris{
		id:    id,
		datas: lines,
	}
}

func (a *Astrolabe) GetEphemeris(id int, dt *calendar.Solar) ObserveData {
	tl := time.Date(dt.GetYear(), time.Month(dt.GetMonth()), dt.GetDay(), dt.GetHour(), dt.GetMinute(), 0, 0, time.Local)
	st := tl.In(time.UTC)
	te := st.Add(time.Hour * 24)
	sts := st.Format(DataTimeMin)
	d := a.OEData[id]
	if d == nil {
		ets := te.Format(DataTimeMin)
		d = a.GetNASAData(id, sts, ets)
		a.OEData[id] = d
	}
	return d.GetOE(sts)
}

type ObserveData struct {
	data map[string]string
}

func (c *ObserveData) RA_DEC() string {
	return c.data["R.A._____(ICRF)_____DEC"] //05 29 58.88 +23 15 24.3
}

// RA 赤经
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

// DEC 赤纬
func (c *ObserveData) DEC() string {
	return c.data["R.A._____(ICRF)_____DEC"][12:] //+23 15 24.3
}

// Delta 距离AU
func (c *ObserveData) Delta() float64 {
	v, _ := strconv.ParseFloat(c.data["delta"], 64)
	return v
}

// Deldot 距离变化AU/d
func (c *ObserveData) Deldot() float64 {
	v, _ := strconv.ParseFloat(c.data["deldot"], 64)
	return v
}
func (c *ObserveData) SOT() float64 {
	v, _ := strconv.ParseFloat(c.data["S-O-T"], 64)
	return v
}
func (c *ObserveData) SOTR() string {
	return c.data["/r"]
}
func (c *ObserveData) Cnst() string {
	//var Cnst = map[string]string{
	//	"Ari": "白羊座", "Tau": "金牛座", "Gem": "双子座", "Can": "巨蟹座",
	//	"Leo": "狮子座", "Vir": "室女座", "Lib": "天秤座", "Sco": "天蝎座",
	//	"Sgr": "射手座", "Cap": "摩羯座", "Aqr": "水瓶座", "Psc": "双鱼座",
	//}
	return c.data["Cnst"]
}

type ObserveEphemeris struct {
	id    int
	datas map[string]ObserveData
}

func (c *ObserveEphemeris) Id() int {
	return c.id
}
func (c *ObserveEphemeris) GetOE(dataTime string) ObserveData {
	return c.datas[dataTime]
}

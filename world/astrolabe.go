package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Astrolabe struct {
	centerX, centerY float32
	solarX, solarY   float32
	stars            map[int]Star
	observer         int
}
type Star struct {
	id      int
	orbitsR float32
	moon    []*Star
}

func NewAstrolabe() *Astrolabe {
	return &Astrolabe{
		centerX:  770,
		centerY:  450,
		observer: 399,
		stars: map[int]Star{
			10:  {id: 10, orbitsR: 0},
			199: {id: 199, orbitsR: 15},
			299: {id: 299, orbitsR: 30},
			399: {id: 399, orbitsR: 60, moon: []*Star{{id: 301, orbitsR: 15}}},
			499: {id: 499, orbitsR: 90},
			599: {id: 599, orbitsR: 105},
			699: {id: 699, orbitsR: 120},
			799: {id: 799, orbitsR: 135},
			899: {id: 899, orbitsR: 150},
			999: {id: 999, orbitsR: 165},
		},
	}
}

func (a *Astrolabe) Update() error {
	solar := uiQiMen.pan.Solar
	hour := solar.GetHour()
	minute := solar.GetMinute()
	//计算太阳位置
	angleDegrees := 90 + (float64(hour)+float64(minute)/60)*15
	solarX, solarY := calRadiansPos(float64(a.centerX), float64(a.centerY), float64(a.stars[a.observer].orbitsR), angleDegrees)
	a.solarX, a.solarY = float32(solarX), float32(solarY)
	//计算星体位置
	//https://ssd.jpl.nasa.gov/api/horizons.api?format=text&COMMAND='499'&OBJ_DATA='YES'&MAKE_EPHEM='YES'&EPHEM_TYPE='OBSERVER'&CENTER='500@399'&START_TIME='2024-06-14 8:00:00'&STOP_TIME='2024-06-14 10:00:00'&STEP_SIZE='1%20h'&QUANTITIES='1,9,20,23,24,29'
	//COMMAND:
	//	Sun (10)
	//	Mercury (199)水星
	//	Venus (299)金星
	//	Moon (301)月球
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

	return nil
}

func (a *Astrolabe) Draw(screen *ebiten.Image) {
	//ft := ui.GetDefaultUIFont()
	const r1 = 240
	sX, sY := a.solarX, a.solarY
	cx, cy := a.centerX, a.centerY
	//外圈
	vector.StrokeCircle(screen, cx, cy, r1, zhiPanWidth/2, colorSkyGateCircle, true)
	//十字线
	vector.StrokeLine(screen, cx-r1, cy, cx+r1, cy, 1, colorOrbits, true)
	vector.StrokeLine(screen, cx, cy-r1, cx, cy+r1, 1, colorOrbits, true)

	//画12宫
	for i := 1; i <= 12; i++ {
		angleDegrees := float64(i+2) * 30 //+ float64(g.count)
		lx1, ly1 := calRadiansPos(float64(cx), float64(cy), float64(r1-zhiPanWidth/4), angleDegrees-15)
		lx2, ly2 := calRadiansPos(float64(cx), float64(cy), float64(r1+zhiPanWidth/4), angleDegrees-15)
		vector.StrokeLine(screen, float32(lx1), float32(ly1), float32(lx2), float32(ly2), 1, colorGongSplit, true)
	}
	for id, star := range a.stars {
		vector.StrokeCircle(screen, sX, sY, float32(star.orbitsR), 1, colorOrbits, true)
		if a.observer != id {
			vector.StrokeCircle(screen, sX, sY-float32(star.orbitsR), 2, 1, colorLeader, true)
		}
		for _, moon := range star.moon {
			vector.StrokeCircle(screen, cx, cy, float32(moon.orbitsR), 1, colorOrbits, true)
			vector.StrokeCircle(screen, float32(cx), float32(cy-moon.orbitsR), 2, 1, colorLeader, true)
		}
	}
}

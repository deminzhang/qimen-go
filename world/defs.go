package world

import "image/color"

const (
	TPSRate      = 30
	screenWidth  = 1024
	screenHeight = 768
)

var (
	colorWhite            = color.RGBA{0xff, 0xff, 0xff, 0xff}
	colorBlack            = color.RGBA{0x00, 0x00, 0x00, 0xff}
	colorGongSplit        = color.RGBA{0x00, 0x00, 0x00, 0xff}
	colorPowerCircle      = color.RGBA{0x60, 0x60, 0xFF, 0xFF}
	colorGroundGateCircle = color.RGBA{0x80, 0x80, 0x00, 0xff}
	colorSkyGateCircle    = color.RGBA{0x40, 0x40, 0xFF, 0x80}
	colorJiang            = color.RGBA{0x00, 0x00, 0x00, 0xff}
	colorJian             = colorJiang
	colorLeader           = color.RGBA{0xff, 0xff, 0x00, 0xff}
	colorGate             = color.RGBA{0x00, 0xff, 0x00, 0xff}
	colorOrbits           = color.RGBA{0x20, 0x20, 0x20, 0x20}
	colorCross            = color.RGBA{0x60, 0x60, 0x60, 0x20}
	colorRedShift         = color.RGBA{0xff, 0xaa, 0x00, 0xff}
	colorBlueShift        = color.RGBA{0x00, 0xff, 0x77, 0xff}

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
	color5Xing = map[string]color.RGBA{
		"金": {0xff, 0xff, 0x00, 0xff},
		"水": {0x00, 0xff, 0xff, 0xff},
		"木": {0x00, 0xff, 0x00, 0xff},
		"火": {0xff, 0x00, 0x00, 0xff},
		"土": {0xff, 0x80, 0x00, 0xff},
	}
	colorGanZhi = map[string]color.RGBA{
		"甲": color5Xing["木"],
		"乙": color5Xing["木"],
		"丙": color5Xing["火"],
		"丁": color5Xing["火"],
		"戊": color5Xing["土"],
		"己": color5Xing["土"],
		"庚": color5Xing["金"],
		"辛": color5Xing["金"],
		"壬": color5Xing["水"],
		"癸": color5Xing["水"],
		"子": color5Xing["水"],
		"丑": color5Xing["土"],
		"寅": color5Xing["木"],
		"卯": color5Xing["木"],
		"辰": color5Xing["土"],
		"巳": color5Xing["火"],
		"午": color5Xing["火"],
		"未": color5Xing["土"],
		"申": color5Xing["金"],
		"酉": color5Xing["金"],
		"戌": color5Xing["土"],
		"亥": color5Xing["水"],
	}
)

// HideGan 藏干
var HideGan = map[string][]string{
	"子": {"癸"},
	"丑": {"己", "辛", "癸"},
	"寅": {"甲", "丙", "戊"},
	"卯": {"乙"},
	"辰": {"戊", "癸", "乙"},
	"巳": {"丙", "庚", "戊"},
	"午": {"丁", "己"},
	"未": {"己", "乙", "丁"},
	"申": {"庚", "壬", "戊"},
	"酉": {"辛"},
	"戌": {"戊", "丁", "辛"},
	"亥": {"壬", "甲"},
}

// HideGanVal 藏干值比例
var HideGanVal = map[int][]int{
	1: {10},
	2: {7, 3},
	3: {6, 3, 1},
}

func GetHideGan(gan string, idx int) string {
	a := HideGan[gan]
	if a == nil {
		return ""
	}
	if idx < len(a) {
		return a[idx]
	}
	return ""
}

package world

import (
	"github.com/6tail/lunar-go/LunarUtil"
	"image/color"
)

const (
	TPSRate      = 30
	screenWidth  = 1024
	screenHeight = 768
)

var (
	colorWhite            = color.RGBA{0xff, 0xff, 0xff, 0xff}
	colorBlack            = color.RGBA{0x00, 0x00, 0x00, 0xff}
	colorRed              = color.RGBA{0xff, 0x00, 0x00, 0xff}
	colorGreen            = color.RGBA{0x00, 0xff, 0x00, 0xff}
	colorBlue             = color.RGBA{0x00, 0x00, 0xff, 0xff}
	colorYellow           = color.RGBA{0xff, 0xff, 0x00, 0xff}
	colorPurple           = color.RGBA{0xff, 0x00, 0xff, 0xff}
	colorGray             = color.RGBA{0x80, 0x80, 0x80, 0xff}
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
	colorChong            = colorPurple
	colorXing             = colorRed
	colorHe               = colorGreen
	colorHe6              = colorBlue
	colorGong             = colorGreen

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
)

const (
	GenderMale   = 1
	GenderFemale = 0
)

var GenderName = []string{"女", "男"}

// ZHI_HIDE_GAN 藏干
var ZHI_HIDE_GAN = map[string][]string{
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

var SHI_SHEN_SHORT = map[string]string{
	"比肩": "比", "劫财": "劫",
	"食神": "食", "伤官": "伤",
	"偏财": "才", "正财": "财",
	"七杀": "杀", "正官": "官",
	"偏印": "枭", "正印": "印",
}

// HideGanVal 藏干值比例
var HideGanVal = map[int][]int{
	1: {10},
	2: {7, 3},
	3: {6, 3, 1},
}

func GetHideGan(gan string, idx int) string {
	a := ZHI_HIDE_GAN[gan]
	if idx < len(a) {
		return a[idx]
	}
	return ""
}

func ColorGanZhi(gz string) color.RGBA {
	wx := LunarUtil.WU_XING_GAN[gz]
	if wx == "" {
		wx = LunarUtil.WU_XING_ZHI[gz]
	}
	return color5Xing[wx]
}

func ShiShenShort(dayGan, gan string) string {
	sx := LunarUtil.SHI_SHEN[dayGan+gan]
	return gan + SHI_SHEN_SHORT[sx]
}

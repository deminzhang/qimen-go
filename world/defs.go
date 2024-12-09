package world

import (
	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/deminzhang/qimen-go/qimen"
	"image/color"
)

const (
	TPSRate          = 30
	initScreenWidth  = 1394
	initScreenHeight = 1000
)

var (
	colorWhite            color.Color = &color.RGBA{0xff, 0xff, 0xff, 0xff}
	colorBlack            color.Color = &color.RGBA{0x00, 0x00, 0x00, 0xff}
	colorRed              color.Color = &color.RGBA{0xff, 0x00, 0x00, 0xff}
	colorDarkRed          color.Color = &color.RGBA{0x80, 0x00, 0x00, 0xff}
	colorGreen            color.Color = &color.RGBA{0x00, 0xff, 0x00, 0xff}
	colorBlue             color.Color = &color.RGBA{0x00, 0x00, 0xff, 0xff}
	colorYellow           color.Color = &color.RGBA{0xff, 0xff, 0x00, 0xff}
	colorPurple           color.Color = &color.RGBA{0xff, 0x00, 0xff, 0xff}
	colorGray             color.Color = &color.RGBA{0x80, 0x80, 0x80, 0xff}
	colorPink             color.Color = &color.RGBA{0xff, 0x80, 0x80, 0xff}
	colorGongSplit        color.Color = &color.RGBA{0x00, 0x00, 0x00, 0xff}
	colorPowerCircle      color.Color = &color.RGBA{0x60, 0x60, 0xFF, 0xFF}
	colorGroundGateCircle color.Color = &color.RGBA{0x80, 0x80, 0x00, 0xff}
	colorSkyGateCircle    color.Color = &color.RGBA{0x40, 0x40, 0xFF, 0x80}
	colorJiang            color.Color = &color.RGBA{0x00, 0x00, 0x00, 0xff}
	colorJian                         = colorJiang
	colorLeader           color.Color = &color.RGBA{0xff, 0xff, 0x00, 0xff}
	colorGate             color.Color = &color.RGBA{0x00, 0xff, 0x00, 0xff}
	colorOrbits           color.Color = &color.RGBA{0x20, 0x20, 0x20, 0x20}
	colorCross            color.Color = &color.RGBA{0x60, 0x60, 0x60, 0x20}
	colorRedShift                     = color.RGBA{0xff, 0xaa, 0x00, 0xff}
	colorBlueShift                    = color.RGBA{0x00, 0xff, 0x77, 0xff}
	colorGray2            color.Color = &color.RGBA{R: 0x20, G: 0x20, B: 0x20, A: 0xff}
	colorGray5            color.Color = &color.RGBA{R: 0x50, G: 0x50, B: 0x50, A: 0xff}
	colorSun              color.Color = &color.RGBA{R: 0xff, G: 0xff, A: 0xff}
	colorMoon             color.Color = &color.RGBA{R: 0xff, G: 0xff, B: 0xcc, A: 0xcc}
	colorDuty                         = colorGreen   //符值马时
	colorTomb                         = colorDarkRed //奇门入墓
	colorJiXing                       = colorPurple  //奇门击刑
	colorMengPo                       = colorRed     //奇门门迫
	colorXingMu                       = colorBlue    //奇门刑墓
	colorChong                        = colorPurple  //星盘冲
	colorXing                         = colorRed     //星盘刑
	colorHe                           = colorGreen   //星盘合
	colorHe6                          = colorBlue    //星盘合
	colorGong                         = colorGreen   //星盘宫

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
		"金": {0xff, 0xff, 0x80, 0xff},
		"水": {0x00, 0xff, 0xff, 0xff},
		"木": {0x00, 0xff, 0x00, 0xff},
		"火": {0xff, 0x00, 0x00, 0xff},
		"土": {0xff, 0x70, 0x00, 0xff},
	}
)

const (
	GenderFemale = 0 //女♀
	GenderMale   = 1 //男♂
)

var GenderName = []string{"女", "男"}
var GenderSoul = []string{"坤造", "乾造"}

//var GenderSymbol = []string{"♀", "♂"}

func GetHideGan(gan string, idx int) string {
	a := LunarUtil.ZHI_HIDE_GAN[gan]
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

func ColorNaYin(gz string) color.RGBA {
	ny := LunarUtil.NAYIN[gz]
	wx := []rune(ny)[2:]
	return color5Xing[string(wx)]
}

func ShiShenShort(dayGan, gan string) string {
	sx := LunarUtil.SHI_SHEN[dayGan+gan]
	return qimen.SHI_SHEN_SHORT[sx]
}

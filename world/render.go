package world

import (
	"fmt"
	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/deminzhang/qimen-go/qimen"
	"github.com/deminzhang/qimen-go/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

// DrawProBar draw a horizontal progress bar
func DrawProBar[T util.Numeric](dst *ebiten.Image, x, y, width, height float32, clr color.Color, val, maxV T, showVal bool) {
	if val < 0 || maxV <= 0 {
		return
	}
	ft, _ := GetFontFace(10)
	vector.DrawFilledRect(dst, x, y, width*float32(val)/float32(maxV), height, clr, true)
	vector.StrokeRect(dst, x, y, width, height, 0.5, clr, true)
	if showVal {
		text.Draw(dst, fmt.Sprintf("%v/%v", val, maxV), ft, int(x+width/3), int(y+8), colorGray)
	}
}

// DrawProBarV draw a vertical progress bar
func DrawProBarV[T util.Numeric](dst *ebiten.Image, x, y, width, height float32, clr color.Color, val, maxV T, showVal bool) {
	if val < 0 || maxV <= 0 {
		return
	}
	ft, _ := GetFontFace(10)
	empty := height - height*float32(val)/float32(maxV)
	vector.DrawFilledRect(dst, x, y+empty, width, height*float32(val)/float32(maxV), clr, true)
	vector.StrokeRect(dst, x, y, width, height, 0.5, clr, true)
	if showVal {
		text.Draw(dst, fmt.Sprintf("%v\n/\n%v", val, maxV), ft, int(x), int(y+8), colorGray)
	}
}

// DrawMixProBar draw a mixed progress bar 混合进度条
func DrawMixProBar[T util.Numeric](dst *ebiten.Image, x, y, width, height float32, clr []color.Color, val []T, maxV T) {
	if maxV <= 0 || len(val) == 0 || 0 == len(clr) || len(clr) != len(val) {
		return
	}
	from := x
	for i := 0; i < len(val); i++ {
		to := width * float32(val[i]) / float32(maxV)
		vector.DrawFilledRect(dst, from, y, to, height, clr[i], true)
		from += to
	}
	vector.StrokeRect(dst, x, y, width, height, .5, clr[0], true)
}

// DrawFlow 流年流月流日流时柱
func DrawFlow(dst *ebiten.Image, sx, sy int, soul string, cb *CharBody) {
	ft14, _ := GetFontFace(14)
	ft28, _ := GetFontFace(28)
	//text.Draw(dst, LunarUtil.SHI_SHEN[soul+cb.Gan], ft14, int(sx), int(sy-32), colorWhite)
	text.Draw(dst, cb.Gan, ft28, sx, sy, ColorGanZhi(cb.Gan))
	text.Draw(dst, cb.Zhi, ft28, sx, sy+32, ColorGanZhi(cb.Zhi))
	text.Draw(dst, cb.Body, ft14, sx, sy+48, ColorGanZhi(cb.Body))
	text.Draw(dst, ShiShenShort(soul, cb.Body), ft14, sx+16, sy+48, colorWhite)
	text.Draw(dst, cb.Legs, ft14, sx, sy+64, ColorGanZhi(cb.Legs))
	text.Draw(dst, ShiShenShort(soul, cb.Legs), ft14, sx+16, sy+64, colorWhite)
	text.Draw(dst, cb.Feet, ft14, sx, sy+80, ColorGanZhi(cb.Feet))
	text.Draw(dst, ShiShenShort(soul, cb.Feet), ft14, sx+16, sy+80, colorWhite)
	text.Draw(dst, LunarUtil.NAYIN[cb.Gan+cb.Zhi], ft14, sx, sy+96, ColorNaYin(cb.Gan+cb.Zhi))
	text.Draw(dst, qimen.ChangSheng12[soul][cb.Zhi], ft14, sx, sy+112, ColorGanZhi(soul))
	text.Draw(dst, qimen.ChangSheng12[cb.Gan][cb.Zhi], ft14, sx, sy+128, ColorGanZhi(cb.Gan))
	text.Draw(dst, LunarUtil.GetXunKong(cb.Gan+cb.Zhi), ft14, sx, sy+144, colorGray)
}

// DrawRangeBar draw a range bar
func DrawRangeBar(dst *ebiten.Image, x, y, width float32, name string, val, minV, maxV float64, clr color.Color) {
	if val < 0 || maxV <= 0 {
		return
	}
	if maxV < val {
		maxV = val
	}
	ft, _ := GetFontFace(12)
	vector.StrokeLine(dst, x, y, x+width, y, 1, colorWhite, true)           //proLine
	vector.StrokeLine(dst, x, y-4, x, y+4, 1, colorWhite, true)             //minPoint
	vector.StrokeLine(dst, x+width, y-4, x+width, y+4, 1, colorWhite, true) //maxPoint
	per := (val - minV) / (maxV - minV)
	if math.IsNaN(per) {
		per = 0
	}
	proV := width * float32(per)
	vector.StrokeLine(dst, x+proV, y-4, x+proV, y+4, 1, clr, true) //valPoint

	//text.Draw(dst, fmt.Sprintf("%v", minV), ft, int(x), int(y-8), colorWhite)
	//text.Draw(dst, fmt.Sprintf("%v", maxV), ft, int(x+width), int(y-8), colorWhite)
	text.Draw(dst, fmt.Sprintf("%.1f%%", per*100), ft, int(x+width/2), int(y-8), clr)
	text.Draw(dst, fmt.Sprintf("%s%v", name, val), ft, int(x+width), int(y), clr)
}

// DrawRangeBarV draw a vertical range bar
func DrawRangeBarV(dst *ebiten.Image, x, y, height float32, name string, val, minV, maxV float64, clr color.Color) {
	if val < 0 || maxV <= 0 {
		return
	}
	if maxV < val {
		maxV = val
	}
	ft, _ := GetFontFace(12)
	vector.StrokeLine(dst, x, y, x, y+height, 1, colorWhite, true)            //proLine
	vector.StrokeLine(dst, x-4, y, x+4, y, 1, colorWhite, true)               //minPoint
	vector.StrokeLine(dst, x-4, y+height, x+4, y+height, 1, colorWhite, true) //maxPoint
	per := (val - minV) / (maxV - minV)
	if math.IsNaN(per) {
		per = 0
	}
	proV := height * float32(per)
	vector.StrokeLine(dst, x-4, y+proV, x+4, y+proV, 1, clr, true) //valPoint

	//text.Draw(dst, fmt.Sprintf("%v", minV), ft, int(x-8), int(y), colorWhite)
	//text.Draw(dst, fmt.Sprintf("%v", maxV), ft, int(x-8), int(y+height), colorWhite)
	text.Draw(dst, fmt.Sprintf("%.1f%%", per*100), ft, int(x-8), int(y+height/2), clr)
	text.Draw(dst, fmt.Sprintf("%s%v", name, val), ft, int(x), int(y+height+8), clr)
}

package world

import (
	"fmt"
	"image/color"
	"math"

	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/deminzhang/qimen-go/util"
	"github.com/deminzhang/qimen-go/xuan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/exp/constraints"
	"golang.org/x/image/font"
)

// DrawProBar draw a horizontal progress bar
func DrawProBar[T constraints.Integer | constraints.Float](dst *ebiten.Image, x, y, width, height float32, clr color.Color, val, maxV T, showVal bool) {
	if val < 0 || maxV <= 0 {
		return
	}
	ft, _ := GetFontXFace(10)
	vector.DrawFilledRect(dst, x, y, width*float32(val)/float32(maxV), height, clr, true)
	vector.StrokeRect(dst, x, y, width, height, 0.5, clr, true)
	if showVal {
		TextDrawV2(dst, fmt.Sprintf("%v/%v", val, maxV), ft, int(x+width/3), int(y+8), colorGray)
	}
}

// DrawProBarV draw a vertical progress bar
func DrawProBarV[T constraints.Integer | constraints.Float](dst *ebiten.Image, x, y, width, height float32, clr color.Color, val, maxV T, showVal bool) {
	if val < 0 || maxV <= 0 {
		return
	}
	ft, _ := GetFontXFace(10)
	empty := height - height*float32(val)/float32(maxV)
	vector.DrawFilledRect(dst, x, y+empty, width, height*float32(val)/float32(maxV), clr, true)
	vector.StrokeRect(dst, x, y, width, height, 0.5, clr, true)
	if showVal {
		TextDrawV2(dst, fmt.Sprintf("%v\n/\n%v", val, maxV), ft, int(x), int(y+8), colorGray)
	}
}

// DrawMixProBar draw a mixed progress bar 混合进度条
func DrawMixProBar[T constraints.Integer | constraints.Float](dst *ebiten.Image, x, y, width, height float32, clr []color.Color, val []T, maxV T) {
	if maxV <= 0 || len(val) == 0 || 0 == len(clr) || len(clr) != len(val) {
		return
	}
	from := x
	for i := 0; i < len(val); i++ {
		to := width * float32(val[i]) / float32(maxV)
		vector.DrawFilledRect(dst, from, y, to, height, clr[i], true)
		from += to
	}
}

// DrawFlow 流年流月流日流时柱
func DrawFlow(dst *ebiten.Image, sx, sy int, soul string, cb *CharBody) {
	ft14, _ := GetFontXFace(14)
	ft28, _ := GetFontXFace(28)
	//TextDrawV2(dst, LunarUtil.SHI_SHEN[soul+cb.Gan], ft14, int(sx), int(sy-32), colorWhite)
	TextDrawV2(dst, cb.Gan, ft28, sx, sy, ColorGanZhi(cb.Gan))
	TextDrawV2(dst, cb.Zhi, ft28, sx, sy+32, ColorGanZhi(cb.Zhi))
	TextDrawV2(dst, cb.Body, ft14, sx, sy+48, ColorGanZhi(cb.Body))
	TextDrawV2(dst, ShiShenShort(soul, cb.Body), ft14, sx+16, sy+48, colorWhite)
	TextDrawV2(dst, cb.Legs, ft14, sx, sy+64, ColorGanZhi(cb.Legs))
	TextDrawV2(dst, ShiShenShort(soul, cb.Legs), ft14, sx+16, sy+64, colorWhite)
	TextDrawV2(dst, cb.Feet, ft14, sx, sy+80, ColorGanZhi(cb.Feet))
	TextDrawV2(dst, ShiShenShort(soul, cb.Feet), ft14, sx+16, sy+80, colorWhite)
	TextDrawV2(dst, LunarUtil.NAYIN[cb.Gan+cb.Zhi], ft14, sx, sy+96, ColorNaYin(cb.Gan+cb.Zhi))
	TextDrawV2(dst, xuan.ZhangSheng12[soul][cb.Zhi], ft14, sx, sy+112, ColorGanZhi(soul))
	TextDrawV2(dst, xuan.ZhangSheng12[cb.Gan][cb.Zhi], ft14, sx, sy+128, ColorGanZhi(cb.Gan))
	TextDrawV2(dst, LunarUtil.GetXunKong(cb.Gan+cb.Zhi), ft14, sx, sy+144, colorGray)
}

// DrawRangeBar draw a range bar
func DrawRangeBar(dst *ebiten.Image, x, y, width float32, name string, val, minV, maxV float64, clr color.Color) {
	if val < 0 || maxV <= 0 {
		return
	}
	if maxV < val {
		maxV = val
	}
	ft, _ := GetFontXFace(12)
	vector.StrokeLine(dst, x, y, x+width, y, 1, colorWhite, true)           //proLine
	vector.StrokeLine(dst, x, y-4, x, y+4, 1, colorWhite, true)             //minPoint
	vector.StrokeLine(dst, x+width, y-4, x+width, y+4, 1, colorWhite, true) //maxPoint
	per := (val - minV) / (maxV - minV)
	if math.IsNaN(per) {
		per = 0
	}
	proV := width * float32(per)
	vector.StrokeLine(dst, x+proV, y-4, x+proV, y+4, 1, clr, true) //valPoint

	//TextDrawV2(dst, fmt.Sprintf("%v", minV), ft, int(x), int(y-8), colorWhite)
	//TextDrawV2(dst, fmt.Sprintf("%v", maxV), ft, int(x+width), int(y-8), colorWhite)
	TextDrawV2(dst, fmt.Sprintf("%.1f%%", per*100), ft, int(x+width/2), int(y-8), clr)
	TextDrawV2(dst, fmt.Sprintf("%s%v", name, val), ft, int(x+width), int(y), clr)
}

// DrawRangeBarV draw a vertical range bar
func DrawRangeBarV(dst *ebiten.Image, x, y, height float32, name string, val, minV, maxV float64, clr color.Color) {
	if val < 0 || maxV <= 0 {
		return
	}
	if maxV < val {
		maxV = val
	}
	ft, _ := GetFontXFace(12)
	vector.StrokeLine(dst, x, y, x, y+height, 1, colorWhite, true)            //proLine
	vector.StrokeLine(dst, x-4, y, x+4, y, 1, colorWhite, true)               //minPoint
	vector.StrokeLine(dst, x-4, y+height, x+4, y+height, 1, colorWhite, true) //maxPoint
	per := (val - minV) / (maxV - minV)
	if math.IsNaN(per) {
		per = 0
	}
	proV := height * float32(per)
	vector.StrokeLine(dst, x-4, y+proV, x+4, y+proV, 1, clr, true) //valPoint

	//TextDrawV2(dst, fmt.Sprintf("%v", minV), ft, int(x-8), int(y), colorWhite)
	//TextDrawV2(dst, fmt.Sprintf("%v", maxV), ft, int(x-8), int(y+height), colorWhite)
	TextDrawV2(dst, fmt.Sprintf("%.1f%%", per*100), ft, int(x-8), int(y+height/2), clr)
	TextDrawV2(dst, fmt.Sprintf("%s%v", name, val), ft, int(x), int(y+height+8), clr)
}

func DrawRotateText(dst *ebiten.Image, x, y, r, rot float64, txt string, fontSize float64, clr color.Color) {
	ft, _ := GetFontFace(fontSize)
	xf := text.NewGoXFace(ft)
	fx, fy := x, y
	for _, t := range []rune(txt) {
		b, _ := font.BoundString(ft, string(t))
		w := (b.Max.X - b.Min.X).Ceil()
		h := (b.Max.Y - b.Min.Y).Ceil()
		rr := math.Sqrt(float64(w*w + h*h))
		ly, lx := util.CalRadiansPos(fy, fx, r, rot)
		TextDrawV2(dst, string(t), xf, int(lx), int(ly), clr)
		rot -= math.Asin(fontSize/2/rr) * 2 * math.Pi * 2
	}
}

package world

import (
	"fmt"
	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"qimen/qimen"
	"qimen/util"
)

func DrawProBar[T util.Numeric](dst *ebiten.Image, x, y, width, height float32, clr color.Color, val, maxV T) {
	if val < 0 || maxV <= 0 {
		return
	}
	ft8, _ := GetFontFace(10)
	vector.DrawFilledRect(dst, x, y, width*float32(val)/float32(maxV), height, clr, true)
	vector.StrokeRect(dst, x, y, width, height, 0.5, clr, true)
	text.Draw(dst, fmt.Sprintf("%v/%v", val, maxV), ft8, int(x+width/3), int(y+8), colorGray)
}

func DrawProBarV[T util.Numeric](dst *ebiten.Image, x, y, width, height float32, clr color.Color, val, maxV T) {
	if val < 0 || maxV <= 0 {
		return
	}
	ft8, _ := GetFontFace(10)
	empty := height - height*float32(val)/float32(maxV)
	vector.DrawFilledRect(dst, x, y+empty, width, height*float32(val)/float32(maxV), clr, true)
	vector.StrokeRect(dst, x, y, width, height, 0.5, clr, true)
	text.Draw(dst, fmt.Sprintf("%v\n/\n%v", val, maxV), ft8, int(x), int(y+8), colorGray)
}

func DrawFlow(dst *ebiten.Image, sx, sy float32, soul string, cb *CharBody) {
	ft14, _ := GetFontFace(14)
	ft28, _ := GetFontFace(28)
	//text.Draw(dst, LunarUtil.SHI_SHEN[soul+cb.Gan], ft14, int(sx), int(sy-32), colorWhite)
	text.Draw(dst, cb.Gan, ft28, int(sx), int(sy), ColorGanZhi(cb.Gan))
	text.Draw(dst, cb.Zhi, ft28, int(sx), int(sy+32), ColorGanZhi(cb.Zhi))
	text.Draw(dst, cb.Body, ft14, int(sx), int(sy+48), ColorGanZhi(cb.Body))
	text.Draw(dst, ShiShenShort(soul, cb.Body), ft14, int(sx+16), int(sy+48), colorWhite)
	text.Draw(dst, cb.Legs, ft14, int(sx), int(sy+64), ColorGanZhi(cb.Legs))
	text.Draw(dst, ShiShenShort(soul, cb.Legs), ft14, int(sx+16), int(sy+64), colorWhite)
	text.Draw(dst, cb.Feet, ft14, int(sx), int(sy+80), ColorGanZhi(cb.Feet))
	text.Draw(dst, ShiShenShort(soul, cb.Feet), ft14, int(sx+16), int(sy+80), colorWhite)
	text.Draw(dst, LunarUtil.NAYIN[cb.Gan+cb.Zhi], ft14, int(sx), int(sy+96), ColorNaYin(cb.Gan+cb.Zhi))
	text.Draw(dst, qimen.ChangSheng12[soul][cb.Zhi], ft14, int(sx), int(sy+112), ColorGanZhi(soul))
	text.Draw(dst, qimen.ChangSheng12[cb.Gan][cb.Zhi], ft14, int(sx), int(sy+128), ColorGanZhi(cb.Gan))
	text.Draw(dst, LunarUtil.GetXunKong(cb.Gan+cb.Zhi), ft14, int(sx), int(sy+144), colorGray)
}

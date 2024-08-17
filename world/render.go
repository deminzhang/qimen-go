package world

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
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

package ui

import (
	"github.com/deminzhang/qimen-go/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

func NewTaiJiImage(size int) *ebiten.Image {
	halfSize := float32(size / 2)
	tj := ebiten.NewImage(size, size)
	vector.DrawFilledCircle(tj, halfSize, halfSize, halfSize, color.White, true)                  //阳大
	vector.DrawFilledRect(tj, halfSize, 0, halfSize, float32(size), color.Black, false)           //阴半
	vector.DrawFilledCircle(tj, halfSize, float32(size/4), float32(size/4), color.White, true)    //太阳
	vector.DrawFilledCircle(tj, halfSize, float32(size*3/4), float32(size/4), color.Black, true)  //太阴
	vector.DrawFilledCircle(tj, halfSize, float32(size*2/8), float32(size/16), color.Black, true) //少阴
	vector.DrawFilledCircle(tj, halfSize, float32(size*6/8), float32(size/16), color.White, true) //少阳

	cover := ebiten.NewImage(size, size)
	vector.DrawFilledCircle(cover, halfSize, halfSize, halfSize, color.Black, true) //mask
	op := &ebiten.DrawImageOptions{}
	op.Blend = ebiten.BlendSourceIn
	cover.DrawImage(tj, op)
	return cover
}

func NewBaGuaImage(gua string, size int) *ebiten.Image {
	bg := ebiten.NewImage(size, size)
	h := float32(size) / 5
	ll := float32(size) * 2 / 5
	drawUpS := func() {
		vector.DrawFilledRect(bg, 0, 0, float32(size), h, color.White, true)
	}
	drawMidS := func() {
		vector.DrawFilledRect(bg, 0, float32(size)*2/5, float32(size), h, color.White, false)
	}
	drawDownS := func() {
		vector.DrawFilledRect(bg, 0, float32(size)*4/5, float32(size), h, color.White, false)
	}
	drawUpL := func() {
		vector.DrawFilledRect(bg, 0, 0, ll, h, color.White, false)
		vector.DrawFilledRect(bg, ll+h, 0, ll, h, color.White, false)
	}
	drawMidL := func() {
		vector.DrawFilledRect(bg, 0, float32(size)*2/5, ll, float32(size)/5, color.White, false)
		vector.DrawFilledRect(bg, ll+h, float32(size)*2/5, ll, h, color.White, false)
	}
	drawDownL := func() {
		vector.DrawFilledRect(bg, 0, float32(size)*4/5, ll, float32(size)/5, color.White, false)
		vector.DrawFilledRect(bg, ll+h, float32(size)*4/5, ll, h, color.White, false)
	}
	switch gua {
	case "乾":
		drawUpS()
		drawMidS()
		drawDownS()
	case "坤":
		drawUpL()
		drawMidL()
		drawDownL()
	case "艮":
		drawUpS()
		drawMidL()
		drawDownL()
	case "兑":
		drawUpL()
		drawMidS()
		drawDownS()
	case "震":
		drawUpL()
		drawMidL()
		drawDownS()
	case "巽":
		drawUpS()
		drawMidS()
		drawDownL()
	case "坎":
		drawUpL()
		drawMidS()
		drawDownL()
	case "离":
		drawUpS()
		drawMidL()
		drawDownS()
	default:
		bg.DrawImage(NewTaiJiImage(size), nil)
	}
	return bg
}

func NewSunImage(size int) *ebiten.Image {
	sun := ebiten.NewImage(size, size)
	hs := float32(size / 2)
	vector.DrawFilledCircle(sun, hs, hs, float32(size/4), color.White, true)
	for i := 0; i < 8; i++ {
		y, x := util.CalRadiansPos(hs, hs, hs, float32(i*45))
		vector.StrokeLine(sun, hs, hs, x, y, 0.5, color.White, true)
	}
	return sun
}

func NewMoonImage(size int) *ebiten.Image {
	moon := ebiten.NewImage(size, size)
	hs := float32(size / 2)
	vector.DrawFilledCircle(moon, hs, hs, float32(size/3), color.White, true)

	cover := ebiten.NewImage(size, size)
	vector.DrawFilledCircle(cover, float32(size/4), hs, float32(size/3), color.Black, true)
	op := &ebiten.DrawImageOptions{}
	op.Blend = ebiten.BlendSourceOut
	cover.DrawImage(moon, op)
	return cover
}

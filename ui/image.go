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

func NewSunImage(size int) *ebiten.Image {
	sun := ebiten.NewImage(size, size)
	hs := float32(size / 2)
	vector.DrawFilledCircle(sun, hs, hs, float32(size/4), color.White, true)
	for i := 0; i < 8; i++ {
		lx, ly := util.CalRadiansPos(hs, hs, hs, float32(i*45))
		vector.StrokeLine(sun, hs, hs, lx, ly, 0.5, color.White, true)
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

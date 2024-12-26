package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"math/rand"
)

type Line struct {
	sx, sy, tx, ty                            float32 // start, target
	fromX, fromY, tox, toy, brightness, scale float32 // position, brightness, scale
	dx, dy                                    float32 // length 长度分量 delta
	speed                                     float32 // 速度
}

func NewLine(sx, sy, tx, ty, scale, speed float32) Line {
	return Line{sx: sx, sy: sy, tx: tx, ty: ty, scale: scale, speed: speed}
}

func (s *Line) Init(scale float32) {
	s.brightness = rand.Float32() * 0xff
	s.scale = scale

	ml := float32(math.Sqrt(float64((s.tx-s.sx)*(s.tx-s.sx) + (s.ty-s.sy)*(s.ty-s.sy))))
	s.dx = (s.tx - s.sx) / ml * scale
	s.dy = (s.ty - s.sy) / ml * scale

	s.fromX = s.sx
	s.fromY = s.sy
	s.tox = s.fromX + s.dx*(s.speed+.01)
	s.toy = s.fromY + s.dy*(s.speed+.01)
}

func (s *Line) Update() {
	w, h := ebiten.WindowSize()
	ww, wh := float32(w), float32(h)
	s.brightness += 1
	if 0xff < s.brightness {
		s.brightness = 0xff
	}

	maxDis := float32(math.Sqrt(float64((s.tx-s.sx)*(s.tx-s.sx) + (s.ty-s.sy)*(s.ty-s.sy))))
	s.tox += s.dx * s.speed
	s.toy += s.dy * s.speed
	headDis := float32(math.Sqrt(float64((s.tox-s.sx)*(s.tox-s.sx) + (s.toy-s.sy)*(s.toy-s.sy))))
	if headDis > maxDis { //头部到达目标点
		s.tox = s.tx
		s.toy = s.ty
		s.fromX += s.dx * s.speed
		s.fromY += s.dy * s.speed
		ll := float32(math.Sqrt(float64((s.tox-s.fromX)*(s.tox-s.fromX) + (s.toy-s.fromY)*(s.toy-s.fromY))))
		if ll < 1 {
			s.Init(s.scale)
			return
		}
	} else {
		ll := float32(math.Sqrt(float64((s.tox-s.fromX)*(s.tox-s.fromX) + (s.toy-s.fromY)*(s.toy-s.fromY))))
		if ll > s.scale { //长度足够,尾点更新
			s.fromX = s.tox - s.dx
			s.fromY = s.toy - s.dy
		}
	}

	if (s.tox < 0 && s.toy < 0 && s.fromX < 0 && s.fromY < 0) || (s.tox > ww && s.toy > wh && s.fromX > ww && s.fromY > wh) {
		s.Init(s.scale)
	}
}

func (s *Line) Draw(dst *ebiten.Image) {
	c := color.RGBA{
		R: uint8(0xbb * s.brightness / 0xff),
		G: uint8(0xdd * s.brightness / 0xff),
		B: uint8(0xff * s.brightness / 0xff),
		A: 0xff}
	vector.StrokeLine(dst, s.fromX, s.fromY, s.tox, s.toy, 1, c, true)
}

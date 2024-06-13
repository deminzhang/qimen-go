package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math/rand"
)

const (
	starsLineScale = 64
	starsLineCount = 1024
)

type StarLine struct {
	fromx, fromy, tox, toy, brightness float32
}

func (s *StarLine) Init() {
	s.tox = rand.Float32() * screenWidth * starsLineScale
	s.fromx = s.tox
	s.toy = rand.Float32() * screenHeight * starsLineScale
	s.fromy = s.toy
	s.brightness = rand.Float32() * 0xff
}

func (s *StarLine) Update(x, y float32) {
	s.fromx = s.tox
	s.fromy = s.toy
	s.tox += (s.tox - x) / 32
	s.toy += (s.toy - y) / 32
	s.brightness += 1
	if 0xff < s.brightness {
		s.brightness = 0xff
	}
	if s.fromx < 0 || screenWidth*starsLineScale < s.fromx || s.fromy < 0 || screenHeight*starsLineScale < s.fromy {
		s.Init()
	}
}

func (s *StarLine) Draw(screen *ebiten.Image) {
	c := color.RGBA{
		R: uint8(0xbb * s.brightness / 0xff),
		G: uint8(0xdd * s.brightness / 0xff),
		B: uint8(0xff * s.brightness / 0xff),
		A: 0xff}
	vector.StrokeLine(screen, s.fromx/starsLineScale, s.fromy/starsLineScale, s.tox/starsLineScale, s.toy/starsLineScale, 1, c, true)
}

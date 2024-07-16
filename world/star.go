package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math/rand"
)

const (
	starsLineScale = 32
	starsLineCount = 256
)

type StarLine struct {
	fromX, fromY, tox, toy, brightness, scale float32
}

func (s *StarLine) Init(scale float32) {
	ww, wh := ebiten.WindowSize()
	s.tox = rand.Float32() * float32(ww) * scale
	s.fromX = s.tox
	s.toy = rand.Float32() * float32(wh) * scale
	s.fromY = s.toy
	s.brightness = rand.Float32() * 0xff
	s.scale = scale
}

func (s *StarLine) Update(x, y float32) {
	ww, wh := ebiten.WindowSize()
	s.fromX = s.tox
	s.fromY = s.toy
	s.tox += (s.tox - x) / 32
	s.toy += (s.toy - y) / 32
	s.brightness += 1
	if 0xff < s.brightness {
		s.brightness = 0xff
	}
	if s.fromX < 0 || float32(ww)*s.scale < s.fromX || s.fromY < 0 || float32(wh)*s.scale < s.fromY {
		s.Init(s.scale)
	}
}

func (s *StarLine) Draw(screen *ebiten.Image) {
	c := color.RGBA{
		R: uint8(0xbb * s.brightness / 0xff),
		G: uint8(0xdd * s.brightness / 0xff),
		B: uint8(0xff * s.brightness / 0xff),
		A: 0xff}
	vector.StrokeLine(screen, s.fromX/s.scale, s.fromY/s.scale, s.tox/s.scale, s.toy/s.scale, 1, c, true)
}

type StarEffect struct {
	centerX, centerY float32
	starScale        float32
	starCount        int
	starLines        []StarLine
}

func NewStarEffect(centerX, centerY float32) *StarEffect {
	se := &StarEffect{
		centerX: centerX, centerY: centerY,
		starScale: starsLineScale,
		starCount: starsLineCount,
	}
	for i := 0; i < se.starCount; i++ {
		se.starLines = append(se.starLines, StarLine{})
		se.starLines[i].Init(se.starScale)
	}
	return se
}

func (se *StarEffect) Update() {
	for i := 0; i < se.starCount; i++ {
		se.starLines[i].Update(se.centerX*se.starScale, se.centerY*se.starScale)
	}
}

func (se *StarEffect) Draw(screen *ebiten.Image) {
	for i := 0; i < se.starCount; i++ {
		se.starLines[i].Draw(screen)
	}
}

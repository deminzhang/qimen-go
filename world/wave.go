package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

type Wave struct {
	x      float32
	y      float32
	radius float32
	maxR   float32
	lenth  float32
	speed  float32
	Sprite *Sprite
}

type WaveComponent struct {
	Waves []*Wave
}

func NewWaveComponent() *WaveComponent {
	return &WaveComponent{}
}

func (c *WaveComponent) AddWave(maxR, speed float32, sprite *Sprite) {
	c.Waves = append(c.Waves, &Wave{radius: 0, maxR: maxR, speed: speed, Sprite: sprite, lenth: 32})
}

func (c *WaveComponent) Update() {
	for _, w := range c.Waves {
		w.radius += w.speed
		if w.Sprite != nil {
			w.x = float32(w.Sprite.x + 8)
			w.y = float32(w.Sprite.y + 8)
		}
	}
}

func (c *WaveComponent) Draw(screen *ebiten.Image) {
	for _, w := range c.Waves {
		numCircles := int(math.Ceil(float64(w.radius/w.lenth))) + 1
		startRadius := w.radius
		for i := 0; i < numCircles; i++ {
			r := float32(i)*w.lenth + float32(int(startRadius)%int(w.lenth))
			if r > w.maxR {
				break
			}
			alpha := uint8(128) //uint8(255 * (1 - (r / w.radius)))
			var clr = color.RGBA{alpha, alpha, alpha, 255}
			vector.StrokeCircle(screen, w.x, w.y, r, 1, clr, true)
		}
	}
}

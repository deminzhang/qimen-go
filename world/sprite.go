package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
)

type Sprite struct {
	image       *ebiten.Image
	alphaImage  *image.Alpha
	colorScale  color.Color
	x           int
	y           int
	DisableMove bool
	dragged     bool
	onMove      func(sx, sy, dx, dy int)
	//onClick     func(x, y int) //TODO
}

// NewSprite creates a new sprite.
// 需ebiten.Game 运行起来后用,否则可能会报错
func NewSprite(img *ebiten.Image, c color.Color) *Sprite {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	alphaImage := image.NewAlpha(image.Rect(0, 0, w, h))
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			_, _, _, a := img.At(i, j).RGBA()
			alpha := uint8(a >> 8)
			alphaImage.SetAlpha(i, j, color.Alpha{A: alpha})
		}
	}
	return &Sprite{
		image:      img,
		alphaImage: alphaImage,
		colorScale: c,
	}
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (s *Sprite) In(x, y int) bool {
	// Check the actual color (alpha) value at the specified position
	// so that the result of In becomes natural to users.
	//
	// Use alphaImage (*image.Alpha) instead of image (*ebiten.Image) here.
	// It is because (*ebiten.Image).At is very slow as this reads pixels from GPU,
	// and should be avoided whenever possible.
	return s.alphaImage.At(x-s.x, y-s.y).(color.Alpha).A > 0
}

// MoveTo moves the sprite to the position (x, y).
func (s *Sprite) MoveTo(x, y int) {
	s.x = x
	s.y = y
	//w, h := s.image.Bounds().Dx(), s.image.Bounds().Dy()
	//if s.x < 0 {
	//	s.x = 0
	//}
	//if s.x > ScreenWidth-w {
	//	s.x = ScreenWidth - w
	//}
	//if s.y < 0 {
	//	s.y = 0
	//}
	//if s.y > ScreenHeight-h {
	//	s.y = ScreenHeight - h
	//}
}

func (s *Sprite) onMoveTo(sx, sy, dx, dy int) {
	s.MoveTo(sx+dx, sy+dy)
	if s.onMove != nil {
		s.onMove(sx, sy, dx, dy)
	}
}

// Draw draws the sprite.
func (s *Sprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.x), float64(s.y))
	if s.colorScale != nil {
		op.ColorScale.ScaleWithColor(s.colorScale)
	}
	if s.dragged {
		op.ColorScale.ScaleAlpha(0.5)
	}
	screen.DrawImage(s.image, op)
}

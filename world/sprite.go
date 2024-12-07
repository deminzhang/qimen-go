package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"image/color"
)

type Sprite struct {
	image      *ebiten.Image
	alphaImage *image.Alpha
	x          int
	y          int
	dragged    bool
}

func NewSprite(img *ebiten.Image) *Sprite {
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
	w, h := s.image.Bounds().Dx(), s.image.Bounds().Dy()

	s.x = x
	s.y = y
	if s.x < 0 {
		s.x = 0
	}
	if s.x > screenWidth-w {
		s.x = screenWidth - w
	}
	if s.y < 0 {
		s.y = 0
	}
	if s.y > screenHeight-h {
		s.y = screenHeight - h
	}
}

// Draw draws the sprite.
func (s *Sprite) Draw(screen *ebiten.Image, alpha float32) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.x), float64(s.y))
	op.ColorScale.ScaleAlpha(alpha)
	screen.DrawImage(s.image, op)
}

// StrokeSource represents a input device to provide strokes.
type StrokeSource interface {
	Position() (int, int)
	IsJustReleased() bool
}

// MouseStrokeSource is a StrokeSource implementation of mouse.
type MouseStrokeSource struct{}

func (m *MouseStrokeSource) Position() (int, int) {
	return ebiten.CursorPosition()
}

func (m *MouseStrokeSource) IsJustReleased() bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
}

// TouchStrokeSource is a StrokeSource implementation of touch.
type TouchStrokeSource struct {
	ID ebiten.TouchID
}

func (t *TouchStrokeSource) Position() (int, int) {
	return ebiten.TouchPosition(t.ID)
}

func (t *TouchStrokeSource) IsJustReleased() bool {
	return inpututil.IsTouchJustReleased(t.ID)
}

// Stroke manages the current drag state by mouse.
type Stroke struct {
	source StrokeSource

	// offsetX and offsetY represents a relative value from the sprite's upper-left position to the cursor position.
	offsetX int
	offsetY int

	// sprite represents a sprite being dragged.
	sprite *Sprite
}

func NewStroke(source StrokeSource, sprite *Sprite) *Stroke {
	sprite.dragged = true
	x, y := source.Position()
	return &Stroke{
		source:  source,
		offsetX: x - sprite.x,
		offsetY: y - sprite.y,
		sprite:  sprite,
	}
}

func (s *Stroke) Update() {
	if !s.sprite.dragged {
		return
	}
	if s.source.IsJustReleased() {
		s.sprite.dragged = false
		return
	}

	x, y := s.source.Position()
	x -= s.offsetX
	y -= s.offsetY
	s.sprite.MoveTo(x, y)
}

func (s *Stroke) Sprite() *Sprite {
	return s.sprite
}

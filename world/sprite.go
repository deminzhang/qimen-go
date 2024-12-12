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
	colorScale color.Color
	x          int
	y          int
	dragged    bool
	onMove     func(sx, sy, dx, dy int)
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
	op.ColorScale.ScaleWithColor(s.colorScale)
	if s.dragged {
		op.ColorScale.ScaleAlpha(0.5)
	}
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
	//s.sprite.MoveTo(x, y)
	s.sprite.onMoveTo(s.sprite.x, s.sprite.y, x-s.sprite.x, y-s.sprite.y)
}

func (s *Stroke) Sprite() *Sprite {
	return s.sprite
}

type StrokeManager struct {
	touchIDs []ebiten.TouchID
	strokes  map[*Stroke]struct{}
	sprites  []*Sprite
}

func (g *StrokeManager) spriteAt(x, y int) *Sprite {
	// As the sprites are ordered from back to front,
	// search the clicked/touched sprite in reverse order.
	for i := len(g.sprites) - 1; i >= 0; i-- {
		s := g.sprites[i]
		if s.In(x, y) {
			return s
		}
	}
	return nil
}

func (g *StrokeManager) moveSpriteToFront(sprite *Sprite) {
	index := -1
	for i, ss := range g.sprites {
		if ss == sprite {
			index = i
			break
		}
	}
	g.sprites = append(g.sprites[:index], g.sprites[index+1:]...)
	g.sprites = append(g.sprites, sprite)
}

func (g *StrokeManager) AddSprite(sprite *Sprite) {
	g.sprites = append(g.sprites, sprite)
}
func (g *StrokeManager) RemoveSprite(sprite *Sprite) {
	for i, s := range g.sprites {
		if s == sprite {
			g.sprites = append(g.sprites[:i], g.sprites[i+1:]...)
			break
		}
	}
}

func (g *StrokeManager) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if sp := g.spriteAt(ebiten.CursorPosition()); sp != nil {
			s := NewStroke(&MouseStrokeSource{}, sp)
			g.strokes[s] = struct{}{}
			g.moveSpriteToFront(sp)
		}
	}
	g.touchIDs = inpututil.AppendJustPressedTouchIDs(g.touchIDs[:0])
	for _, id := range g.touchIDs {
		if sp := g.spriteAt(ebiten.TouchPosition(id)); sp != nil {
			s := NewStroke(&TouchStrokeSource{id}, sp)
			g.strokes[s] = struct{}{}
			g.moveSpriteToFront(sp)
		}
	}

	for s := range g.strokes {
		s.Update()
		if !s.sprite.dragged {
			delete(g.strokes, s)
		}
	}
	return nil
}

package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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

func (g *StrokeManager) Update() {
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
			if !sp.DisableMove {
				s := NewStroke(&TouchStrokeSource{id}, sp)
				g.strokes[s] = struct{}{}
				g.moveSpriteToFront(sp)
			}
		}
	}

	for s := range g.strokes {
		s.Update()
		if !s.sprite.dragged {
			delete(g.strokes, s)
		}
	}
}

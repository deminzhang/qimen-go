package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"image"
	"image/color"
)

type TextButton struct {
	Button
	drawBorder bool
}

// NewTextButton 按text长度自动调整大小无背景UI
func NewTextButton(x0, y0 int, text string, c color.Color, drawBorder bool) *TextButton {
	return &TextButton{
		Button: Button{BaseUI: BaseUI{Visible: true, X: 0, Y: 0},
			Rect:      image.Rect(x0, y0, x0, y0),
			Text:      text,
			TextColor: c,
		},
		drawBorder: drawBorder,
	}
}
func (b *TextButton) Update() {
	bounds, _ := font.BoundString(uiFont, b.Text)
	w := (bounds.Max.X - bounds.Min.X).Ceil()
	b.Rect.Max.X = b.Rect.Min.X + w + 6
	b.Rect.Max.Y = b.Rect.Min.Y + uiFontMHeight + 6
	b.textX = b.Rect.Min.X + (b.Rect.Dx()-w)/2
	b.textY = b.Rect.Max.Y - (b.Rect.Dy()-uiFontMHeight)/2
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if b.Rect.Min.X <= x && x < b.Rect.Max.X && b.Rect.Min.Y <= y && y < b.Rect.Max.Y {
			b.mouseDown = true
		} else {
			b.mouseDown = false
		}
	} else {
		if b.mouseDown {
			if b.onClick != nil {
				b.onClick(&b.Button)
			}
		}
		b.mouseDown = false
	}
}

func (b *TextButton) Draw(dst *ebiten.Image) {
	if !b.Visible {
		return
	}
	if b.drawBorder {
		vector.StrokeRect(dst, float32(b.Rect.Min.X), float32(b.Rect.Min.Y), float32(b.Rect.Dx()), float32(b.Rect.Dy()),
			0.5, color.Gray{Y: 128}, true)
	}
	text.Draw(dst, b.Text, uiFont, b.textX, b.textY, color.White)
}

package gui

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
func NewTextButton(x, y int, text string, c color.Color, drawBorder bool) *TextButton {
	return &TextButton{
		Button: Button{BaseUI: BaseUI{Visible: true, X: x, Y: y, Rect: image.Rect(0, 0, 0, 0)},
			Text:      text,
			TextColor: c,
		},
		drawBorder: drawBorder,
	}
}
func (b *TextButton) Update() {
	bounds, _ := font.BoundString(uiFont, b.Text)
	w := (bounds.Max.X - bounds.Min.X).Ceil()
	b.Rect.Max.X = b.Rect.Min.X + w + 6             // 自动调整大小w
	b.Rect.Max.Y = b.Rect.Min.Y + uiFontMHeight + 6 // 自动调整大小h
	b.textX = b.Rect.Min.X + (b.Rect.Dx()-w)/2
	b.textY = b.Rect.Max.Y - (b.Rect.Dy()-uiFontMHeight)/2
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := b.GetXY()
		mx, my := ebiten.CursorPosition()
		if x+b.Rect.Min.X <= mx && mx < x+b.Rect.Max.X && y+b.Rect.Min.Y <= my && my < y+b.Rect.Max.Y {
			b.mouseDown = true
		} else {
			b.mouseDown = false
		}
	} else {
		if b.mouseDown {
			if b.onClick != nil {
				if !IsFrameClick() {
					b.onClick(&b.Button)
					SetFrameClick()
				}
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

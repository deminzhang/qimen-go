package gui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"image/color"
)

type TextButton struct {
	Button
}

// NewTextButton 按text长度自动调整大小无背景UI
func NewTextButton(x, y int, text string, textColor, bgColor *color.RGBA) *TextButton {
	return &TextButton{
		Button: Button{BaseUI: BaseUI{Visible: true, X: x, Y: y, W: 1, H: 1, BDColor: bgColor},
			Text:      text,
			TextColor: textColor,
		},
	}
}
func (b *TextButton) Update() {
	bounds, _ := font.BoundString(uiFont, b.Text)
	w := (bounds.Max.X - bounds.Min.X).Ceil()
	b.W = w + 6             // 自动调整大小w
	b.H = uiFontMHeight + 6 // 自动调整大小h
	b.textX = (b.W - w) / 2
	b.textY = b.H - (b.H-uiFontMHeight)/2
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		x, y := b.GetWorldXY()
		if x <= mx && mx < x+b.W && y <= my && my < y+b.H {
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
	text.Draw(dst, b.Text, uiFont, b.textX, b.textY, b.TextColor)
}

package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"image"
	"image/color"
)

type Button struct {
	BaseUI
	//Rect         image.Rectangle
	Text         string
	textX, textY int

	mouseDown bool

	onClick          func(b *Button)
	TextColor        color.Color
	UIImage          *ebiten.Image
	ImageRect        image.Rectangle
	ImageRectPressed image.Rectangle
}

func NewButton(rect image.Rectangle, text string) *Button {
	return &Button{
		BaseUI: BaseUI{Visible: true, X: 0, Y: 0,
			Rect: rect,
		},
		Text: text,
		//default resource
		TextColor:        color.Black,
		UIImage:          GetDefaultUIImage(),
		ImageRect:        imageSrcRects[imageTypeButton],
		ImageRectPressed: imageSrcRects[imageTypeButtonPressed],
	}
}

func NewButtonTransparent(rect image.Rectangle, text string) *Button {
	return &Button{
		BaseUI: BaseUI{Visible: true, X: 0, Y: 0,
			Rect: rect,
		},
		Text:      text,
		TextColor: color.White,
	}
}

func (b *Button) Update() {
	bounds, _ := font.BoundString(uiFont, b.Text)
	w := (bounds.Max.X - bounds.Min.X).Ceil()
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
				b.onClick(b)
			}
		}
		b.mouseDown = false
	}
}

func (b *Button) Draw(dst *ebiten.Image) {
	if !b.Visible {
		return
	}
	imageRect := b.ImageRect
	if b.mouseDown {
		imageRect = b.ImageRectPressed
	}
	if b.UIImage == nil {
		vector.StrokeRect(dst, float32(b.Rect.Min.X), float32(b.Rect.Min.Y), float32(b.Rect.Dx()), float32(b.Rect.Dy()),
			0.5, color.Gray{Y: 128}, true)
		text.Draw(dst, b.Text, uiFont, b.textX, b.textY, color.White)
	} else {
		drawNinePatches(dst, b.UIImage, b.Rect, imageRect)
		text.Draw(dst, b.Text, uiFont, b.textX, b.textY, color.Black)
	}
}

func (b *Button) SetOnClick(f func(b *Button)) {
	b.onClick = f
}

func (b *Button) Click() {
	if b.onClick != nil {
		b.onClick(b)
	}
}

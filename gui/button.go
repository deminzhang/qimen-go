package gui

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
	x, y := rect.Min.X, rect.Min.Y
	rect = image.Rect(0, 0, rect.Dx(), rect.Dy())
	return &Button{
		BaseUI: BaseUI{Visible: true, X: x, Y: y, Rect: rect},
		Text:   text,
		//default resource
		TextColor:        color.Black,
		UIImage:          GetDefaultUIImage(),
		ImageRect:        imageSrcRects[imageTypeButton],
		ImageRectPressed: imageSrcRects[imageTypeButtonPressed],
	}
}

func NewButtonTransparent(rect image.Rectangle, text string) *Button {
	x, y := rect.Min.X, rect.Min.Y
	rect = image.Rect(0, 0, rect.Dx(), rect.Dy())
	return &Button{
		BaseUI:    BaseUI{Visible: true, X: x, Y: y, Rect: rect},
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
		mx, my := ebiten.CursorPosition()
		x, y := b.GetXY()
		if x+b.Rect.Min.X <= mx && mx < x+b.Rect.Max.X && y+b.Rect.Min.Y <= my && my < y+b.Rect.Max.Y {
			b.mouseDown = true
		} else {
			b.mouseDown = false
		}
	} else {
		if b.mouseDown {
			if b.onClick != nil {
				b.onClick(b)
				SetFrameClick()
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

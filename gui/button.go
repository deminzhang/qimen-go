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

func NewButton(x, y, w, h int, text string) *Button {
	return &Button{
		BaseUI: BaseUI{Visible: true, X: x, Y: y, W: w, H: h},
		Text:   text,
		//default resource
		TextColor:        color.Black,
		UIImage:          GetDefaultUIImage(),
		ImageRect:        imageSrcRects[imageTypeButton],
		ImageRectPressed: imageSrcRects[imageTypeButtonPressed],
	}
}

func NewButtonTransparent(x, y, w, h int, text string) *Button {
	return &Button{
		BaseUI:    BaseUI{Visible: true, X: x, Y: y, W: w, H: h},
		Text:      text,
		TextColor: color.White,
	}
}

func (b *Button) Update() {
	bounds, _ := font.BoundString(uiFont, b.Text)
	w := (bounds.Max.X - bounds.Min.X).Ceil()
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
		vector.StrokeRect(dst, float32(b.X), float32(b.Y), float32(b.W), float32(b.H),
			0.5, color.Gray{Y: 128}, true)
		text.Draw(dst, b.Text, uiFont, b.textX, b.textY, color.White)
	} else {
		drawNinePatches(dst, b.UIImage, image.Rect(0, 0, b.W, b.H), imageRect)
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

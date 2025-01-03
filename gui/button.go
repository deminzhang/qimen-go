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
	Text           string
	textX, textY   int
	AutoSizeByText bool

	TextColor        color.Color
	UIImage          *ebiten.Image
	ImageRect        image.Rectangle
	ImageRectPressed image.Rectangle

	mouseDown bool
	onClick   func()
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

// NewTextButton 按text长度自动调整大小无背景UI
func NewTextButton(x, y int, text string, textColor, bdColor color.Color) *Button {
	return &Button{BaseUI: BaseUI{Visible: true, X: x, Y: y, W: 1, H: 1, BDColor: bdColor},
		Text:           text,
		TextColor:      textColor,
		AutoSizeByText: true,
	}
}

func (b *Button) updateMouse() {
	mx, my := ebiten.CursorPosition()
	x, y := b.GetWorldXY()
	cursorIn := x <= mx && mx < x+b.W && y <= my && my < y+b.H
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		b.mouseDown = cursorIn
	} else {
		if b.mouseDown {
			if b.onClick != nil {
				if !IsFrameClick() {
					b.onClick()
					SetFrameClick()
				}
			}
			b.mouseDown = false
		}
	}
}

func (b *Button) Update() {
	b.BaseUI.Update()
	bounds, _ := font.BoundString(uiFont, b.Text)
	w := (bounds.Max.X - bounds.Min.X).Ceil()
	if b.AutoSizeByText {
		b.W = w + 6
		b.H = uiFontMHeight + 6
	}
	b.textX = (b.W - w) / 2
	b.textY = b.H - (b.H-uiFontMHeight)/2
	b.updateMouse()
}

func (b *Button) Draw(dst *ebiten.Image) {
	if b.UIImage == nil {
		vector.StrokeRect(dst, float32(b.X), float32(b.Y), float32(b.W), float32(b.H),
			0.5, color.Gray{Y: 128}, true)
		if b.mouseDown {
			text.Draw(dst, b.Text, uiFont, b.textX+1, b.textY+1, b.TextColor)
		} else {
			text.Draw(dst, b.Text, uiFont, b.textX, b.textY, b.TextColor)
		}
	} else {
		imageRect := b.ImageRect
		if b.mouseDown {
			imageRect = b.ImageRectPressed
		}
		drawNinePatches(dst, b.UIImage, image.Rect(0, 0, b.W, b.H), imageRect)
		text.Draw(dst, b.Text, uiFont, b.textX, b.textY, color.Black)
	}
}

func (b *Button) SetOnClick(f func()) {
	b.onClick = f
}

func (b *Button) Click() {
	if b.onClick != nil {
		b.onClick()
	}
}

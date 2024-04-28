package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"image"
	"image/color"
)

const (
	checkBoxWidth       = 16
	checkBoxPaddingLeft = 8
)

type CheckBox struct {
	BaseUI
	Text string

	checked   bool
	mouseDown bool

	onCheckChanged func(c *CheckBox)

	UIImage          *ebiten.Image
	ImageRect        image.Rectangle
	ImageRectPressed image.Rectangle
	ImageRectMark    image.Rectangle
}

func NewCheckBox(x, y int, text string) *CheckBox {
	return &CheckBox{
		BaseUI: BaseUI{Visible: true, X: x, Y: y},
		Text:   text,

		UIImage:          GetDefaultUIImage(),
		ImageRect:        imageSrcRects[imageTypeCheckBox],
		ImageRectPressed: imageSrcRects[imageTypeCheckBoxPressed],
		ImageRectMark:    imageSrcRects[imageTypeCheckBoxMark],
	}

}

func (c *CheckBox) width() int {
	b, _ := font.BoundString(uiFont, c.Text)
	w := (b.Max.X - b.Min.X).Ceil()
	return checkBoxWidth + checkBoxPaddingLeft + w
}

func (c *CheckBox) Update() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		c.mouseDown = c.X <= x && x < c.X+c.width() && c.Y <= y && y < c.Y+checkBoxWidth
	} else {
		if c.mouseDown {
			c.checked = !c.checked
			if c.onCheckChanged != nil {
				c.onCheckChanged(c)
			}
		}
		c.mouseDown = false
	}
}

func (c *CheckBox) Draw(dst *ebiten.Image) {
	if !c.Visible {
		return
	}
	r := image.Rect(c.X, c.Y, c.X+checkBoxWidth, c.Y+checkBoxWidth)
	if c.mouseDown {
		drawNinePatches(dst, c.UIImage, r, c.ImageRectPressed)
	} else {
		drawNinePatches(dst, c.UIImage, r, c.ImageRect)
	}
	if c.checked {
		drawNinePatches(dst, c.UIImage, r, c.ImageRectMark)
	}

	x := c.X + checkBoxWidth + checkBoxPaddingLeft
	y := (c.Y + 16) - (16-uiFontMHeight)/2
	if c.Disabled {
		text.Draw(dst, c.Text, uiFont, x, y, color.Gray16{Y: 0x8888})
	} else {
		text.Draw(dst, c.Text, uiFont, x, y, color.White)
	}
}

func (c *CheckBox) SetChecked(b bool) {
	c.checked = b
}

func (c *CheckBox) Checked() bool {
	return c.checked
}

func (c *CheckBox) SetOnCheckChanged(f func(c *CheckBox)) {
	c.onCheckChanged = f
}

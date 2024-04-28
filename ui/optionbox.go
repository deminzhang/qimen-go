package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"image"
	"image/color"
)

const (
	optionBoxWidth       = 16
	optionBoxPaddingLeft = 8
)

type OptionBox struct {
	BaseUI
	Text string

	selected    bool
	mouseDown   bool
	optionGroup map[*OptionBox]struct{}

	onSelect func(c *OptionBox)

	UIImage          *ebiten.Image
	ImageRect        image.Rectangle
	ImageRectPressed image.Rectangle
	ImageRectMark    image.Rectangle
}

func NewOptionBox(x, y int, text string) *OptionBox {
	return &OptionBox{
		BaseUI: BaseUI{Visible: true, X: x, Y: y},
		Text:   text,

		UIImage:          GetDefaultUIImage(),
		ImageRect:        imageSrcRects[imageTypeOptionBox],
		ImageRectPressed: imageSrcRects[imageTypeOptionBoxPressed],
		ImageRectMark:    imageSrcRects[imageTypeOptionBoxMark],
	}
}

func MakeOptionBoxGroup(a ...*OptionBox) {
	g := map[*OptionBox]struct{}{}
	for _, o := range a {
		g[o] = struct{}{}
		o.optionGroup = g
	}
}

func (o *OptionBox) width() int {
	b, _ := font.BoundString(uiFont, o.Text)
	w := (b.Max.X - b.Min.X).Ceil()
	return checkBoxWidth + checkBoxPaddingLeft + w
}

func (o *OptionBox) Update() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if o.X <= x && x < o.X+o.width() && o.Y <= y && y < o.Y+checkBoxWidth {
			o.mouseDown = true
		} else {
			o.mouseDown = false
		}
	} else {
		if o.mouseDown {
			o.setSelected(true)
		}
		o.mouseDown = false
	}
}

func (o *OptionBox) Draw(dst *ebiten.Image) {
	if !o.Visible {
		return
	}
	r := image.Rect(o.X, o.Y, o.X+optionBoxWidth, o.Y+optionBoxWidth)
	if o.mouseDown {
		drawNinePatches(dst, o.UIImage, r, o.ImageRectPressed)
	} else {
		drawNinePatches(dst, o.UIImage, r, o.ImageRect)
	}
	if o.selected {
		drawNinePatches(dst, o.UIImage, r, o.ImageRectMark)
	}

	x := o.X + optionBoxWidth + optionBoxPaddingLeft
	y := (o.Y + 16) - (16-uiFontMHeight)/2
	if o.Disabled {
		text.Draw(dst, o.Text, uiFont, x, y, color.Gray16{Y: 0x8888})
	} else {
		text.Draw(dst, o.Text, uiFont, x, y, color.White)
	}
}

func (o *OptionBox) setSelected(sel bool) {
	if sel {
		if o.optionGroup != nil {
			for box, _ := range o.optionGroup {
				if box != o {
					box.selected = false
				}
			}
		}
		o.selected = true
		if o.onSelect != nil {
			o.onSelect(o)
		}
	} else {
		o.selected = false
	}
}

func (o *OptionBox) Select() {
	o.setSelected(true)
}
func (o *OptionBox) Selected() bool {
	return o.selected
}

func (o *OptionBox) SetOnSelect(f func(c *OptionBox)) {
	o.onSelect = f
}

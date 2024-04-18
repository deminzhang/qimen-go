package ebiten_ui

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
	X    int
	Y    int
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
		BaseUI: BaseUI{Visible: true},
		X:      x,
		Y:      y,
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

func (c *OptionBox) width() int {
	b, _ := font.BoundString(uiFont, c.Text)
	w := (b.Max.X - b.Min.X).Ceil()
	return checkBoxWidth + checkBoxPaddingLeft + w
}

func (c *OptionBox) Update() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if c.X <= x && x < c.X+c.width() && c.Y <= y && y < c.Y+checkBoxWidth {
			c.mouseDown = true
		} else {
			c.mouseDown = false
		}
	} else {
		if c.mouseDown {
			c.setSelected(true)
		}
		c.mouseDown = false
	}
}

func (c *OptionBox) Draw(dst *ebiten.Image) {
	if !c.Visible {
		return
	}
	r := image.Rect(c.X, c.Y, c.X+optionBoxWidth, c.Y+optionBoxWidth)
	if c.mouseDown {
		drawNinePatches(dst, c.UIImage, r, c.ImageRectPressed)
	} else {
		drawNinePatches(dst, c.UIImage, r, c.ImageRect)
	}
	if c.selected {
		drawNinePatches(dst, c.UIImage, r, c.ImageRectMark)
	}

	x := c.X + optionBoxWidth + optionBoxPaddingLeft
	y := (c.Y + 16) - (16-uiFontMHeight)/2
	if c.Disabled {
		text.Draw(dst, c.Text, uiFont, x, y, color.Gray16{Y: 0x8888})
	} else {
		text.Draw(dst, c.Text, uiFont, x, y, color.White)
	}
}

func (c *OptionBox) setSelected(b bool) {
	if b {
		if c.optionGroup != nil {
			for box, _ := range c.optionGroup {
				if box != c {
					box.selected = false
				}
			}
		}
		c.selected = true
		if c.onSelect != nil {
			c.onSelect(c)
		}
	} else {
		c.selected = false
	}
}

func (c *OptionBox) Select() {
	c.setSelected(true)
}
func (c *OptionBox) Selected() bool {
	return c.selected
}

func (c *OptionBox) SetOnSelect(f func(c *OptionBox)) {
	c.onSelect = f
}

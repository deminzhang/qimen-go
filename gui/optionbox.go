package gui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
)

const (
	optionBoxWidth       = 16
	optionBoxPaddingLeft = 8
)

type OptionBox struct {
	BaseUI
	Text string

	selected       bool
	mouseDown      bool
	boxWidth       int
	boxPaddingLeft int
	optionGroup    map[*OptionBox]struct{}

	onSelect         func(c *OptionBox)
	UIImage          *ebiten.Image
	ImageRect        image.Rectangle
	ImageRectPressed image.Rectangle
	ImageRectMark    image.Rectangle
}

func NewOptionBox(x, y int, text string) *OptionBox {
	return &OptionBox{
		BaseUI:           BaseUI{Visible: true, X: x, Y: y, W: optionBoxWidth, H: optionBoxWidth},
		Text:             text,
		boxWidth:         optionBoxWidth,
		boxPaddingLeft:   optionBoxPaddingLeft,
		UIImage:          GetDefaultUIImage(),
		ImageRect:        imageSrcRects[imageTypeOptionBox],
		ImageRectPressed: imageSrcRects[imageTypeOptionBoxPressed],
		ImageRectMark:    imageSrcRects[imageTypeOptionBoxMark],
	}
}

func MakeOptionBoxGroup(a ...*OptionBox) map[*OptionBox]struct{} {
	g := map[*OptionBox]struct{}{}
	for _, o := range a {
		g[o] = struct{}{}
		o.optionGroup = g
	}
	return g
}

func (o *OptionBox) Add2OptionBoxGroup(g map[*OptionBox]struct{}) {
	g[o] = struct{}{}
	o.optionGroup = g
}

func (o *OptionBox) width() int {
	b, _ := font.BoundString(uiFont, o.Text)
	w := (b.Max.X - b.Min.X).Ceil()
	return optionBoxWidth + checkBoxPaddingLeft + w
}

func (o *OptionBox) Update() {
	o.BaseUI.Update()
	o.W = o.width()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		x, y := o.GetWorldXY()
		o.mouseDown = x <= mx && mx < x+o.width() && y <= my && my < y+optionBoxWidth
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
	r := image.Rect(0, 0, optionBoxWidth, optionBoxWidth)
	if o.mouseDown {
		drawNinePatches(dst, o.UIImage, r, o.ImageRectPressed)
	} else {
		drawNinePatches(dst, o.UIImage, r, o.ImageRect)
	}
	if o.selected {
		drawNinePatches(dst, o.UIImage, r, o.ImageRectMark)
	}

	var c color.Color = color.White
	if o.Disabled {
		c = color.Gray{Y: 128}
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(o.boxWidth), float64(0))
	op.ColorScale.ScaleWithColor(c)
	op.LineSpacing = lineSpacingInPixels
	text.Draw(dst, o.Text, uiFontFace, op)
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

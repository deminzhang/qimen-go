package gui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
)

type Panel struct {
	BaseUI
	Color *color.RGBA
}

func NewPanel(rect image.Rectangle, c *color.RGBA) *Panel {
	x, y := rect.Min.X, rect.Min.Y
	rect = image.Rect(0, 0, rect.Dx(), rect.Dy())
	return &Panel{
		BaseUI: BaseUI{Visible: true, X: x, Y: y, Rect: rect},
		Color:  c,
	}
}

func (p *Panel) Draw(dst *ebiten.Image) {
	if !p.Visible {
		return
	}
	if p.Color != nil {
		vector.DrawFilledRect(dst, float32(p.Rect.Min.X), float32(p.Rect.Min.Y),
			float32(p.Rect.Dx()), float32(p.Rect.Dy()), *p.Color, false)
	}
	p.BaseUI.Draw(dst)
}

func (p *Panel) AddChild(c IUIPanel) {
	p.BaseUI.AddChild(c)
	//TODO Resize by children
	//if autoResize
	//for _, child := range p.children {
	//if child.GetRect
	//p.Rect.Max
	//}
}

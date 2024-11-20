package gui

import (
	"image/color"
)

type Panel struct {
	BaseUI
}

func NewPanel(x, y, w, h int, c *color.RGBA) *Panel {
	return &Panel{
		BaseUI: BaseUI{Visible: true, X: x, Y: y, W: w, H: h,
			BGColor: c},
	}
}

// TODO Resize by children
func (p *Panel) AutoResize() {
}

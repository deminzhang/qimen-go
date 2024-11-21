package gui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

func (p *Panel) Draw(dst *ebiten.Image) {
	p.BaseUI.Draw(dst)
	if uiBorderDebug {
		vector.StrokeRect(dst, 1, 1, float32(p.W-1), float32(p.H-1), 1, color.Gray{Y: 128}, true)
	}
}

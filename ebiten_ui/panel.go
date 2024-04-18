package ebiten_ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type Panel struct {
	BaseUI
	Rect             image.Rectangle
	UIImage          *ebiten.Image
	ImageRect        image.Rectangle
	ImageRectPressed image.Rectangle
}

func NewPanel(rect image.Rectangle) *Panel {
	return &Panel{
		BaseUI: BaseUI{Visible: true},
		Rect:   rect,
		//default resource
		UIImage:          GetDefaultUIImage(),
		ImageRect:        image.Rect(0, 0, 16, 16),
		ImageRectPressed: image.Rect(16, 0, 32, 16),
	}
}

func (p *Panel) Update() {
}

func (p *Panel) Draw(dst *ebiten.Image) {
	if !p.Visible {
		return
	}
	drawNinePatches(dst, p.UIImage, p.Rect, p.ImageRect)
}

func (p *Panel) AddChild(c IUIPanel) {
	p.BaseUI.AddChild(c)
	//if autoResize
	//for _, child := range p.Children {
	//if child.GetRect
	//p.Rect.Max
	//}
}

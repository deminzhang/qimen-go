package ui

import (
	"image"
)

type Panel struct {
	BaseUI
	Rect image.Rectangle
	//UIImage   *ebiten.Image
	//ImageRect image.Rectangle
}

func NewPanel(rect image.Rectangle) *Panel {
	return &Panel{
		BaseUI: BaseUI{Visible: true, X: 0, Y: 0},
		Rect:   rect,
		//default resource
		//UIImage:   GetDefaultUIImage(),
		//ImageRect: imageSrcRects[imageTypeTextBox],
	}
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

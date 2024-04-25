package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type Panel struct {
	BaseUI
	Rect      image.Rectangle
	UIImage   *ebiten.Image
	ImageRect image.Rectangle
}

func NewPanel(rect image.Rectangle) *Panel {
	return &Panel{
		BaseUI: BaseUI{Visible: true},
		Rect:   rect,
		//default resource
		UIImage:   GetDefaultUIImage(),
		ImageRect: imageSrcRects[imageTypeTextBox],
	}
}

func (p *Panel) Update() {
}

func (p *Panel) Draw(dst *ebiten.Image) {
	if !p.Visible {
		return
	}
	drawNinePatches(dst, p.UIImage, p.Rect, p.ImageRect)
	p.BaseUI.Draw(dst)
}

func (p *Panel) AddChild(c IUIPanel) {
	p.BaseUI.AddChild(c)
	//TODO Resize by Children
	//if autoResize
	//for _, child := range p.Children {
	//if child.GetRect
	//p.Rect.Max
	//}
}

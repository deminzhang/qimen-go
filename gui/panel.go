package gui

import (
	"image/color"
)

type Panel struct {
	BaseUI
}

func NewPanel(x, y, w, h int, bgColor color.Color) *Panel {
	return &Panel{
		BaseUI: BaseUI{Visible: true, X: x, Y: y, W: w, H: h,
			BGColor:  bgColor,
			autoSize: true,
		},
	}
}

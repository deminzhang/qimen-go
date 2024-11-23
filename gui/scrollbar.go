package gui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
)

const (
	defaultVScrollBarWidth  = 16
	defaultHScrollBarHeight = 16
)

// VScrollBar 竖向ScrollBar
type VScrollBar struct {
	BaseUI

	thumbRate           float64
	thumbOffset         int
	dragging            bool
	draggingStartOffset int
	draggingStartY      int
	contentOffset       int

	UIImage        *ebiten.Image
	ImageRectBack  image.Rectangle
	ImageRectFront image.Rectangle
}

// HScrollBar 横向ScrollBar
type HScrollBar struct {
	BaseUI

	thumbRate           float64
	thumbOffset         int
	dragging            bool
	draggingStartOffset int
	draggingStartX      int
	contentOffset       int

	UIImage        *ebiten.Image
	ImageRectBack  image.Rectangle
	ImageRectFront image.Rectangle
}

func NewVScrollBar() *VScrollBar {
	return &VScrollBar{
		BaseUI:         BaseUI{Visible: true, X: 0, Y: 0, W: defaultVScrollBarWidth, H: 1},
		UIImage:        GetDefaultUIImage(),
		ImageRectBack:  imageSrcRects[imageTypeScrollBarBack],
		ImageRectFront: imageSrcRects[imageTypeScrollBarFront],
	}
}
func NewHScrollBar() *HScrollBar {
	return &HScrollBar{
		BaseUI:         BaseUI{Visible: true, X: 0, Y: 0, W: 1, H: defaultHScrollBarHeight},
		UIImage:        GetDefaultUIImage(),
		ImageRectBack:  imageSrcRects[imageTypeScrollBarBack],
		ImageRectFront: imageSrcRects[imageTypeScrollBarFront],
	}
}

func (v *VScrollBar) thumbSize() int {
	r := v.thumbRate
	if r > 1 {
		r = 1
	}
	s := int(float64(v.H) * r)
	if s < v.W {
		return v.W
	}
	return s
}

func (v *VScrollBar) thumbRect() image.Rectangle {
	if v.thumbRate >= 1 {
		return image.Rectangle{}
	}

	s := v.thumbSize()
	return image.Rect(v.X, v.Y+v.thumbOffset, v.X+v.W, v.Y+v.thumbOffset+s)
}

func (v *VScrollBar) maxThumbOffset() int {
	return v.H - v.thumbSize()
}

func (v *VScrollBar) ContentOffset() int {
	return v.contentOffset
}

func (v *VScrollBar) Update(wx, wy, contentHeight int) {
	v.thumbRate = float64(v.H) / float64(contentHeight)

	if !v.dragging && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		tr := v.thumbRect()
		if wx+tr.Min.X <= x && x < wx+tr.Max.X && wy+tr.Min.Y <= y && y < wy+tr.Max.Y {
			v.dragging = true
			v.draggingStartOffset = v.thumbOffset
			v.draggingStartY = y
		}
	}
	if v.dragging {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			_, y := ebiten.CursorPosition()
			v.thumbOffset = v.draggingStartOffset + (y - v.draggingStartY)
			if v.thumbOffset < 0 {
				v.thumbOffset = 0
			}
			if v.thumbOffset > v.maxThumbOffset() {
				v.thumbOffset = v.maxThumbOffset()
			}
		} else {
			v.dragging = false
		}
	}

	v.contentOffset = 0
	if v.thumbRate < 1 {
		v.contentOffset = int(float64(contentHeight) * float64(v.thumbOffset) / float64(v.H))
	}
}

func (v *VScrollBar) Draw(dst *ebiten.Image) {
	if !v.Visible {
		return
	}
	sd := image.Rect(v.X, v.Y, v.X+v.W, v.Y+v.H)
	drawNinePatches(dst, v.UIImage, sd, v.ImageRectBack)

	if v.thumbRate < 1 {
		drawNinePatches(dst, v.UIImage, v.thumbRect(), v.ImageRectFront)
	}
}

//------------------------------------------

func (v *HScrollBar) thumbSize() int {
	r := v.thumbRate
	if r > 1 {
		r = 1
	}
	s := int(float64(v.W) * r)
	if s < v.H {
		return v.H
	}
	return s
}

func (v *HScrollBar) thumbRect() image.Rectangle {
	if v.thumbRate >= 1 {
		return image.Rectangle{}
	}

	s := v.thumbSize()
	return image.Rect(v.X+v.thumbOffset, v.Y, v.X+v.thumbOffset+s, v.Y+v.H)
}

func (v *HScrollBar) maxThumbOffset() int {
	return v.W - v.thumbSize()
}

func (v *HScrollBar) ContentOffset() int {
	return v.contentOffset
}

func (v *HScrollBar) Update(wx, wy, contentWidth int) {
	v.thumbRate = float64(v.W) / float64(contentWidth)

	if !v.dragging && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		tr := v.thumbRect()
		if wx+tr.Min.X <= x && x < wx+tr.Max.X && wy+tr.Min.Y <= y && y < wy+tr.Max.Y {
			v.dragging = true
			v.draggingStartOffset = v.thumbOffset
			v.draggingStartX = x
		}
	}
	if v.dragging {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, _ := ebiten.CursorPosition()
			v.thumbOffset = v.draggingStartOffset + (x - v.draggingStartX)
			if v.thumbOffset < 0 {
				v.thumbOffset = 0
			}
			if v.thumbOffset > v.maxThumbOffset() {
				v.thumbOffset = v.maxThumbOffset()
			}
		} else {
			v.dragging = false
		}
	}

	v.contentOffset = 0
	if v.thumbRate < 1 {
		v.contentOffset = int(float64(contentWidth) * float64(v.thumbOffset) / float64(v.W))
	}
}

func (v *HScrollBar) Draw(dst *ebiten.Image) {
	if !v.Visible {
		return
	}
	sd := image.Rect(v.X, v.Y, v.X+v.W, v.Y+v.H)
	drawNinePatches(dst, v.UIImage, sd, v.ImageRectBack)

	if v.thumbRate < 1 {
		drawNinePatches(dst, v.UIImage, v.thumbRect(), v.ImageRectFront)
	}
}

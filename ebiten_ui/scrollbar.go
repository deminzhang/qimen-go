package ebiten_ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
)

const VScrollBarWidth = 16
const HScrollBarHeight = 16

// 竖向ScrollBar
type VScrollBar struct {
	BaseUI
	X      int
	Y      int
	Height int

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

// 横向ScrollBar
type HScrollBar struct {
	BaseUI
	X     int
	Y     int
	Width int

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
		BaseUI:         BaseUI{Visible: true},
		UIImage:        GetDefaultUIImage(),
		ImageRectBack:  image.Rect(16, 16, 24, 32),
		ImageRectFront: image.Rect(24, 16, 32, 32),
	}
}
func NewHScrollBar() *HScrollBar {
	return &HScrollBar{
		BaseUI:         BaseUI{Visible: true},
		UIImage:        GetDefaultUIImage(),
		ImageRectBack:  image.Rect(16, 16, 24, 32),
		ImageRectFront: image.Rect(24, 16, 32, 32),
	}
}

func (v *VScrollBar) thumbSize() int {
	const minThumbSize = VScrollBarWidth

	r := v.thumbRate
	if r > 1 {
		r = 1
	}
	s := int(float64(v.Height) * r)
	if s < minThumbSize {
		return minThumbSize
	}
	return s
}

func (v *VScrollBar) thumbRect() image.Rectangle {
	if v.thumbRate >= 1 {
		return image.Rectangle{}
	}

	s := v.thumbSize()
	return image.Rect(v.X, v.Y+v.thumbOffset, v.X+VScrollBarWidth, v.Y+v.thumbOffset+s)
}

func (v *VScrollBar) maxThumbOffset() int {
	return v.Height - v.thumbSize()
}

func (v *VScrollBar) ContentOffset() int {
	return v.contentOffset
}

func (v *VScrollBar) Update(contentHeight int) {
	v.thumbRate = float64(v.Height) / float64(contentHeight)

	if !v.dragging && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		tr := v.thumbRect()
		if tr.Min.X <= x && x < tr.Max.X && tr.Min.Y <= y && y < tr.Max.Y {
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
		v.contentOffset = int(float64(contentHeight) * float64(v.thumbOffset) / float64(v.Height))
	}
}

func (v *VScrollBar) Draw(dst *ebiten.Image) {
	if !v.Visible {
		return
	}
	sd := image.Rect(v.X, v.Y, v.X+VScrollBarWidth, v.Y+v.Height)
	drawNinePatches(dst, v.UIImage, sd, v.ImageRectBack)

	if v.thumbRate < 1 {
		drawNinePatches(dst, v.UIImage, v.thumbRect(), v.ImageRectFront)
	}
}

//------------------------------------------

func (v *HScrollBar) thumbSize() int {
	const minThumbSize = HScrollBarHeight

	r := v.thumbRate
	if r > 1 {
		r = 1
	}
	s := int(float64(v.Width) * r)
	if s < minThumbSize {
		return minThumbSize
	}
	return s
}

func (v *HScrollBar) thumbRect() image.Rectangle {
	if v.thumbRate >= 1 {
		return image.Rectangle{}
	}

	s := v.thumbSize()
	return image.Rect(v.X+v.thumbOffset, v.Y, v.X+v.thumbOffset+s, v.Y+HScrollBarHeight)
}

func (v *HScrollBar) maxThumbOffset() int {
	return v.Width - v.thumbSize()
}

func (v *HScrollBar) ContentOffset() int {
	return v.contentOffset
}

func (v *HScrollBar) Update(contentWidth int) {
	v.thumbRate = float64(v.Width) / float64(contentWidth)

	if !v.dragging && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		tr := v.thumbRect()
		if tr.Min.X <= x && x < tr.Max.X && tr.Min.Y <= y && y < tr.Max.Y {
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
		v.contentOffset = int(float64(contentWidth) * float64(v.thumbOffset) / float64(v.Width))
	}
}

func (v *HScrollBar) Draw(dst *ebiten.Image) {
	if !v.Visible {
		return
	}
	sd := image.Rect(v.X, v.Y, v.X+v.Width, v.Y+HScrollBarHeight)
	drawNinePatches(dst, v.UIImage, sd, v.ImageRectBack)

	if v.thumbRate < 1 {
		drawNinePatches(dst, v.UIImage, v.thumbRect(), v.ImageRectFront)
	}
}

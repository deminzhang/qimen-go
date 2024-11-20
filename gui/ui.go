package gui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"sort"
)

type IUIPanel interface {
	Update()
	Draw(screen *ebiten.Image)

	IsDisabled() bool
	IsVisible() bool
	GetXY() (int, int)
	GetWH() (int, int)
	GetWorldXY() (int, int)
	GetDepth() int
	GetBDColor() color.Color
	GetParent() IUIPanel
	SetParent(p IUIPanel)
}

// var uis = make(map[IUIPanel]struct{})
var uis []IUIPanel
var frameClick bool

func ActiveUI(ui IUIPanel) {
	uis = append(uis, ui)
	sort.Slice(uis, func(a, b int) bool {
		return uis[a].GetDepth() > uis[b].GetDepth()
	})
}
func CloseUI(ui IUIPanel) {
	for i, p := range uis {
		if ui == p {
			if i == 0 {
				uis = uis[1:]
			} else {
				if i+1 == len(uis) {
					uis = uis[:i]
				} else {
					uis = append(uis[:i], uis[i+1:]...)
				}
			}
		}
	}
}
func Update() {
	frameClick = false
	for _, u := range uis {
		if !u.IsDisabled() && u.IsVisible() {
			u.Update()
		}
	}
}

func Draw(screen *ebiten.Image) {
	for _, u := range uis {
		if u.IsVisible() {
			w, h := u.GetWH()
			if w == 0 || h == 0 {
				continue
			}
			img := ebiten.NewImage(w, h)
			x, y := u.GetXY()
			op := ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			u.Draw(img)
			//vector.StrokeRect(img, 1, 1, float32(w-1), float32(h-1), 1,
			//	color.Gray{Y: 128}, false)
			screen.DrawImage(img, &op)
		}
	}
}
func IsFrameClick() bool {
	return frameClick
}
func SetFrameClick() {
	frameClick = true
}

type BaseUI struct {
	X, Y, W, H  int
	Visible     bool //`default:"true"` disable draw
	Disabled    bool //disable update
	EnableFocus bool //enable focus
	Depth       int  //update draw depth
	children    []IUIPanel
	parent      IUIPanel
	Rect        image.Rectangle
	BGColor     color.Color
	BDColor     color.Color
	//GeoM        *ebiten.GeoM
}

func (u *BaseUI) IsDisabled() bool {
	return u.Disabled
}
func (u *BaseUI) IsVisible() bool {
	return u.Visible
}

// GetWH 获取宽高
func (u *BaseUI) GetWH() (int, int) {
	return u.W, u.H
}

// GetXY 获取相对坐标
func (u *BaseUI) GetXY() (int, int) {
	return u.X, u.Y
}

// GetWorldXY 获取绝对坐标
func (u *BaseUI) GetWorldXY() (int, int) {
	x, y := u.X, u.Y
	if u.parent != nil {
		x2, y2 := u.parent.GetWorldXY()
		x += x2
		y += y2
	}
	return x, y
}
func (u *BaseUI) GetDepth() int {
	return u.Depth
}
func (u *BaseUI) GetBDColor() color.Color {
	return u.BDColor
}
func (u *BaseUI) GetParent() IUIPanel {
	return u.parent
}
func (u *BaseUI) SetParent(p IUIPanel) {
	u.parent = p
}

func (u *BaseUI) Update() {
	for _, p := range u.children {
		if !p.IsDisabled() && p.IsVisible() {
			p.Update()
		}
	}
}

func (u *BaseUI) Draw(screen *ebiten.Image) {
	if !u.Visible {
		return
	}
	//if !util.IsNil(u.BGColor) {
	//	vector.DrawFilledRect(screen, 1, 1, float32(u.W-1), float32(u.H-1), u.BGColor, false)
	//}
	for _, p := range u.children {
		if p.IsVisible() {
			w, h := p.GetWH()
			if w == 0 || h == 0 {
				continue
			}
			img := ebiten.NewImage(w, h)
			x, y := p.GetXY()
			op := ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			p.Draw(img)
			//if !util.IsNil(p.GetBDColor()) {
			//vector.StrokeRect(img, 1, 1, float32(w-1), float32(h-1), 1,
			//	color.Gray{Y: 128}, false)
			//}
			screen.DrawImage(img, &op)
		}
	}
}

func (u *BaseUI) AddChild(c IUIPanel) {
	u.children = append(u.children, c)
	sort.Slice(u.children, func(a, b int) bool {
		return u.children[a].GetDepth() > u.children[b].GetDepth()
	})
	c.SetParent(u)
}
func (u *BaseUI) AddChildren(cs ...IUIPanel) {
	for _, c := range cs {
		u.children = append(u.children, c)
		c.SetParent(u)
	}
	sort.Slice(u.children, func(a, b int) bool {
		return u.children[a].GetDepth() > u.children[b].GetDepth()
	})
}

func (u *BaseUI) RemoveChild(c IUIPanel) {
	for i, child := range u.children {
		if c == child {
			c.SetParent(nil)
			if i == 0 {
				u.children = u.children[1:]
			} else {
				if i+1 == len(u.children) {
					u.children = u.children[:i]
				} else {
					u.children = append(u.children[:i], u.children[i+1:]...)
				}
			}
			break
		}
	}
}

var uniqueFocusedUI any

func (u *BaseUI) Focused() bool {
	return uniqueFocusedUI == u
}
func (u *BaseUI) SetFocused(focused bool) {
	if u.EnableFocus {
		if focused {
			uniqueFocusedUI = u
		} else {
			uniqueFocusedUI = nil
		}
	}
}

// Focused 正有焦点接收中
func Focused() bool {
	return uniqueFocusedUI != nil
}

func drawNinePatches(dst *ebiten.Image, uiImage *ebiten.Image, dstRect image.Rectangle, srcRect image.Rectangle) {
	srcX := srcRect.Min.X
	srcY := srcRect.Min.Y
	srcW := srcRect.Dx()
	srcH := srcRect.Dy()

	dstX := dstRect.Min.X
	dstY := dstRect.Min.Y
	dstW := dstRect.Dx()
	dstH := dstRect.Dy()

	op := &ebiten.DrawImageOptions{}
	for j := 0; j < 3; j++ {
		for i := 0; i < 3; i++ {
			op.GeoM.Reset()

			sx := srcX
			sy := srcY
			sw := srcW / 4
			sh := srcH / 4
			dx := 0
			dy := 0
			dw := sw
			dh := sh
			switch i {
			case 1:
				sx = srcX + srcW/4
				sw = srcW / 2
				dx = srcW / 4
				dw = dstW - 2*srcW/4
			case 2:
				sx = srcX + 3*srcW/4
				dx = dstW - srcW/4
			}
			switch j {
			case 1:
				sy = srcY + srcH/4
				sh = srcH / 2
				dy = srcH / 4
				dh = dstH - 2*srcH/4
			case 2:
				sy = srcY + 3*srcH/4
				dy = dstH - srcH/4
			}

			op.GeoM.Scale(float64(dw)/float64(sw), float64(dh)/float64(sh))
			op.GeoM.Translate(float64(dx), float64(dy))
			op.GeoM.Translate(float64(dstX), float64(dstY))
			dst.DrawImage(uiImage.SubImage(image.Rect(sx, sy, sx+sw, sy+sh)).(*ebiten.Image), op)
		}
	}
}

// BoundString忽略了空格宽度,这里加上
func getFontWidth(f font.Face, s string) int {
	bounds_, _ := font.BoundString(f, "_")
	wl_ := (bounds_.Max.X - bounds_.Min.X).Ceil()
	bounds, _ := font.BoundString(f, s+"_")
	wl := (bounds.Max.X - bounds.Min.X).Ceil() - wl_
	return wl
}
func getFontSelectWidth(f font.Face, s string, left, right int) (int, int) {
	bounds_, _ := font.BoundString(f, "_")
	wl_ := (bounds_.Max.X - bounds_.Min.X).Ceil()
	bounds, _ := font.BoundString(f, s[:left]+"_")
	wl := (bounds.Max.X - bounds.Min.X).Ceil() - wl_
	bounds, _ = font.BoundString(f, s[:right]+"_")
	wr := (bounds.Max.X - bounds.Min.X).Ceil() - wl_
	return wl, wr
}
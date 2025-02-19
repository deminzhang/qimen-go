package gui

import (
	"image"
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

type IUIPanel interface {
	Update()
	Draw(screen *ebiten.Image)
	OnClose()
	OnLayout(w, h int)

	IsDisabled() bool
	IsVisible() bool
	GetXY() (int, int)
	GetWH() (int, int)
	GetWorldXY() (int, int)
	GetDepth() int
	GetBDColor() color.Color
	GetParent() IUIPanel
	SetParent(p IUIPanel)
	GetImage() *ebiten.Image
}

// var uis = make(map[IUIPanel]struct{})
var uis []IUIPanel
var frameClick bool
var frameHover bool
var uiBorderDebug bool

func ActiveUI(ui IUIPanel) {
	uis = append(uis, ui)
	sort.Slice(uis, func(a, b int) bool {
		return uis[a].GetDepth() > uis[b].GetDepth()
	})
}
func CloseUI(ui IUIPanel) {
	for i, p := range uis {
		if ui == p {
			p.OnClose()
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
	frameHover = false
	for _, u := range uis {
		if u.IsVisible() {
			u.Update()
		}
	}
}
func OnLayout(w, h int) {
	for _, u := range uis {
		u.OnLayout(w, h)
	}
}
func Draw(screen *ebiten.Image) {
	for _, u := range uis {
		if u.IsVisible() {
			img := u.GetImage()
			if img == nil {
				continue
			}
			x, y := u.GetXY()
			op := ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			u.Draw(img)
			w, h := u.GetWH()
			if u.GetBDColor() != nil {
				vector.StrokeRect(img, 1, 1, float32(w-1), float32(h-1), 1, u.GetBDColor(), false)
			} else {
				if uiBorderDebug {
					vector.StrokeRect(img, 1, 1, float32(w-1), float32(h-1), 1, color.Gray{Y: 128}, true)
				}
			}
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
func IsFrameHover() bool {
	return frameHover
}
func SetFrameHover() {
	frameHover = true
}
func SetBorderDebug(v bool) {
	uiBorderDebug = v
}

type BaseUI struct {
	X, Y, W, H  int
	Depth       int  //update draw depth
	Visible     bool //`default:"true"` disable draw
	Disabled    bool //disable update
	EnableFocus bool //enable focus
	autoSize    bool //auto resize by children
	children    []IUIPanel
	parent      IUIPanel
	BGColor     color.Color
	BDColor     color.Color
	Image       *ebiten.Image
	DrawOp      ebiten.DrawImageOptions
	mouseHover  bool
	onHover     func()
	onHout      func()
}

func (u *BaseUI) IsDisabled() bool {
	return u.Disabled
}
func (u *BaseUI) IsVisible() bool {
	return u.Visible
}

func (u *BaseUI) SetOnHover(f func()) {
	u.onHover = f
}
func (u *BaseUI) SetOnHout(f func()) {
	u.onHout = f
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

func (u *BaseUI) resizeByChildren() {
	cw, ch := 0, 0
	for _, c := range u.children {
		x, y := c.GetXY()
		w, h := c.GetWH()
		if x+w > cw {
			cw = x + w
		}
		if y+h > ch {
			ch = y + h
		}
	}
	u.W = cw
	u.H = ch
}

func (u *BaseUI) GetImage() *ebiten.Image {
	w, h := u.GetWH()
	if w == 0 || h == 0 {
		return nil
	}
	if u.Image != nil {
		dx, dy := u.Image.Bounds().Dx(), u.Image.Bounds().Dy()
		if w != dx || h != dy {
			u.Image.Deallocate()
			u.Image = nil
		}
	}
	if u.Image == nil {
		u.Image = ebiten.NewImage(w, h)
	} else {
		u.Image.Clear()
	}
	return u.Image
}

func (u *BaseUI) Update() {
	if u.autoSize {
		u.resizeByChildren()
	}
	for _, p := range u.children {
		if p.IsVisible() {
			p.Update()
		}
	}

	mx, my := ebiten.CursorPosition()
	x, y := u.GetWorldXY()
	cursorIn := x <= mx && mx < x+u.W && y <= my && my < y+u.H
	if cursorIn {
		if !u.mouseHover {
			u.mouseHover = true
			if u.onHover != nil {
				if !IsFrameHover() {
					u.onHover()
					SetFrameHover()
				}
			}
		}
	} else {
		if u.mouseHover {
			u.mouseHover = false
			if u.onHout != nil {
				u.onHout()
			}
		}
	}
}

func (u *BaseUI) Draw(screen *ebiten.Image) {
	if !u.Visible {
		return
	}
	if u.BGColor != nil {
		if _, _, _, a := u.BGColor.RGBA(); a > 0 {
			vector.DrawFilledRect(screen, 1, 1, float32(u.W-1), float32(u.H-1), u.BGColor, false)
		}
	}
	for _, p := range u.children {
		if p.IsVisible() {
			img := p.GetImage()
			if img == nil {
				continue
			}
			x, y := p.GetXY()
			op := ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			p.Draw(img)
			bdc := p.GetBDColor()
			if bdc != nil {
				w, h := p.GetWH()
				vector.StrokeRect(img, 1, 1, float32(w-1), float32(h-1), 1,
					bdc, false)
			}
			screen.DrawImage(img, &op)
		}
	}

	if uiBorderDebug {
		vector.StrokeRect(screen, 1, 1, float32(u.W-1), float32(u.H-1), 1, color.Gray{Y: 128}, true)
	}
}

func (u *BaseUI) OnClose() {}

func (u *BaseUI) OnLayout(w, h int) {}

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

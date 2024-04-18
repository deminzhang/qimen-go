package ebiten_ui

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"strings"
)

const (
	lineHeight     = 16
	textBoxPadding = 8
)

type TextBox struct {
	BaseUI
	Rect image.Rectangle
	Text string

	contentBuf     *ebiten.Image
	vScrollBar     *VScrollBar
	hScrollBar     *HScrollBar
	offsetX        int
	offsetY        int
	DisableVScroll bool
	DisableHScroll bool

	UIImage   *ebiten.Image
	ImageRect image.Rectangle
}

func NewTextBox(rect image.Rectangle) *TextBox {
	return &TextBox{
		BaseUI: BaseUI{Visible: true},
		Rect:   rect,

		UIImage:   GetDefaultUIImage(),
		ImageRect: imageSrcRects[imageTypeTextBox],
	}
}
func (t *TextBox) SetText(v interface{}) {
	t.Text = fmt.Sprintf("%v", v)
}
func (t *TextBox) AppendLine(line string) {
	if t.Text == "" {
		t.Text = line
	} else {
		t.Text += "\n" + line
	}
}

func (t *TextBox) Update() {
	w, h := t.contentSize()
	if h > t.Rect.Dy() && !t.DisableVScroll {
		if t.vScrollBar == nil {
			t.vScrollBar = NewVScrollBar()
		}
		t.vScrollBar.X = t.Rect.Max.X - VScrollBarWidth
		t.vScrollBar.Y = t.Rect.Min.Y
		t.vScrollBar.Height = t.Rect.Dy()
		if t.hScrollBar != nil {
			t.vScrollBar.Height -= HScrollBarHeight
		}

		t.vScrollBar.Update(h)

		t.offsetY = t.vScrollBar.ContentOffset()
	} else {
		t.vScrollBar = nil
		t.offsetY = 0
	}
	if w > t.Rect.Dx() && !t.DisableHScroll {
		if t.hScrollBar == nil {
			t.hScrollBar = NewHScrollBar()
		}
		t.hScrollBar.X = t.Rect.Min.X
		t.hScrollBar.Y = t.Rect.Max.Y - HScrollBarHeight
		t.hScrollBar.Width = t.Rect.Dx()
		if t.vScrollBar != nil {
			t.hScrollBar.Width -= VScrollBarWidth
		}

		t.hScrollBar.Update(w)

		t.offsetX = t.hScrollBar.ContentOffset()
	} else {
		t.hScrollBar = nil
		t.offsetX = 0
	}
}

func (t *TextBox) contentSize() (int, int) {
	lines := strings.Split(t.Text, "\n")
	h := len(lines) * lineHeight
	w := t.Rect.Dx()
	for _, line := range lines {
		bounds, _ := font.BoundString(uiFont, line)
		w = max(w, (bounds.Max.X-bounds.Min.X).Ceil()+2*textBoxPadding)
		h = max(h, (bounds.Max.Y - bounds.Min.Y).Ceil())
	}
	return w, h
}

func (t *TextBox) viewSize() (int, int) {
	vsb, hsb := 0, 0
	if t.vScrollBar != nil {
		vsb = VScrollBarWidth
	}
	if t.hScrollBar != nil {
		hsb = HScrollBarHeight
	}
	return t.Rect.Dx() - vsb - textBoxPadding, t.Rect.Dy() - hsb
}

func (t *TextBox) contentOffset() (int, int) {
	return t.offsetX, t.offsetY
}

func (t *TextBox) Draw(dst *ebiten.Image) {
	if !t.Visible {
		return
	}
	drawNinePatches(dst, t.UIImage, t.Rect, t.ImageRect)

	if t.contentBuf != nil {
		vw, vh := t.viewSize()
		w, h := t.contentBuf.Size()
		if vw > w || vh > h {
			t.contentBuf.Dispose()
			t.contentBuf = nil
		}
	}
	if t.contentBuf == nil {
		w, h := t.viewSize()
		t.contentBuf = ebiten.NewImage(w, h)
	}

	t.contentBuf.Clear()
	for i, line := range strings.Split(t.Text, "\n") {
		x := -t.offsetX + textBoxPadding
		y := -t.offsetY + i*lineHeight + lineHeight - (lineHeight-uiFontMHeight)/2
		if y < -lineHeight {
			continue
		}
		if _, h := t.viewSize(); y >= h+lineHeight {
			continue
		}
		text.Draw(t.contentBuf, line, uiFont, x, y, color.Black)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(t.Rect.Min.X), float64(t.Rect.Min.Y))
	dst.DrawImage(t.contentBuf, op)

	if t.vScrollBar != nil {
		t.vScrollBar.Draw(dst)
	}
	if t.hScrollBar != nil {
		t.hScrollBar.Draw(dst)
	}
}

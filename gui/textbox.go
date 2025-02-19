package gui

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

const (
	defaultLineHeight     = 16
	defaultTextBoxPadding = 8
)

var defaultTextColor = color.Black

type TextBox struct {
	BaseUI
	TextField *TextField
	Text      string
	TextColor color.Color

	contentBuf     *ebiten.Image
	vScrollBar     *VScrollBar
	hScrollBar     *HScrollBar
	offsetX        int
	offsetY        int
	lineHeight     int
	textBoxPadding int
	DisableVScroll bool
	DisableHScroll bool

	UIImage   *ebiten.Image
	ImageRect image.Rectangle
}

func NewTextBox(x, y, w, h int) *TextBox {
	return &TextBox{
		BaseUI:         BaseUI{Visible: true, X: x, Y: y, W: w, H: h},
		TextField:      NewTextField(image.Rect(x, y, x+w, y+h), true),
		TextColor:      defaultTextColor,
		lineHeight:     defaultLineHeight,
		textBoxPadding: defaultTextBoxPadding,
		UIImage:        GetDefaultUIImage(),
		ImageRect:      imageSrcRects[imageTypeTextBox],
	}
}
func (t *TextBox) SetText(v interface{}) {
	t.Text = fmt.Sprintf("%v", v)
}
func (t *TextBox) AppendText(line string) {
	if t.Text == "" {
		t.Text = line
	} else {
		t.Text += line
	}
}

func (t *TextBox) AppendTextLn(line string) {
	if t.Text == "" {
		t.Text = line
	} else {
		t.Text += "\n" + line
	}
}

func (t *TextBox) Update() {
	t.BaseUI.Update()
	wx, wy := t.GetWorldXY()
	w, h := t.ContentSize()
	if h > t.H && !t.DisableVScroll { // 竖向滚动条
		if t.vScrollBar == nil {
			t.vScrollBar = NewVScrollBar()
		}
		t.vScrollBar.X = t.W - t.vScrollBar.W
		t.vScrollBar.Y = 0
		t.vScrollBar.H = t.H
		if t.hScrollBar != nil {
			t.vScrollBar.H -= t.hScrollBar.H
		}

		t.vScrollBar.Update(wx, wy, h)

		t.offsetY = t.vScrollBar.ContentOffset()
	} else {
		t.vScrollBar = nil
		t.offsetY = 0
	}
	if w > t.W && !t.DisableHScroll { // 横向滚动条
		if t.hScrollBar == nil {
			t.hScrollBar = NewHScrollBar()
		}
		t.hScrollBar.X = 0
		t.hScrollBar.Y = t.H - t.hScrollBar.H
		t.hScrollBar.W = t.W
		if t.vScrollBar != nil {
			t.hScrollBar.W -= t.vScrollBar.W
		}

		t.hScrollBar.Update(wx, wy, w)

		t.offsetX = t.hScrollBar.ContentOffset()
	} else {
		t.hScrollBar = nil
		t.offsetX = 0
	}
}

func (t *TextBox) ContentSize() (int, int) {
	lines := strings.Split(t.Text, "\n")
	h := len(lines) * t.lineHeight
	w := t.W
	for _, line := range lines {
		bounds, _ := font.BoundString(uiFont, line)
		w = max(w, (bounds.Max.X-bounds.Min.X).Ceil()+2*t.textBoxPadding)
		h = max(h, (bounds.Max.Y - bounds.Min.Y).Ceil())
	}
	return w, h
}

func (t *TextBox) viewSize() (int, int) {
	vsb, hsb := 0, 0
	if t.vScrollBar != nil {
		vsb = t.vScrollBar.W
	}
	if t.hScrollBar != nil {
		hsb = t.hScrollBar.H
	}
	return t.W - vsb - t.textBoxPadding, t.H - hsb
}

func (t *TextBox) contentOffset() (int, int) {
	return t.offsetX, t.offsetY
}

func (t *TextBox) Draw(dst *ebiten.Image) {
	if !t.Visible {
		return
	}
	if t.UIImage != nil {
		drawNinePatches(dst, t.UIImage, image.Rect(0, 0, t.W, t.H), t.ImageRect)
	}

	if t.contentBuf != nil {
		vw, vh := t.viewSize()
		w, h := t.contentBuf.Bounds().Dx(), t.contentBuf.Bounds().Dy()
		if vw > w || vh > h {
			t.contentBuf.Deallocate()
			t.contentBuf = nil
		}
	}
	if t.contentBuf == nil {
		w, h := t.viewSize()
		t.contentBuf = ebiten.NewImage(w, h)
	}

	t.contentBuf.Clear()
	for i, line := range strings.Split(t.Text, "\n") {
		x := -t.offsetX + t.textBoxPadding
		y := -t.offsetY + i*t.lineHeight + t.lineHeight - (t.lineHeight-uiFontMHeight)/2
		if y < -t.lineHeight {
			continue
		}
		if _, h := t.viewSize(); y >= h+t.lineHeight {
			continue
		}
		text.Draw(t.contentBuf, line, uiFont, x, y, t.TextColor)
	}
	op := ebiten.DrawImageOptions{}
	dst.DrawImage(t.contentBuf, &op)

	if t.vScrollBar != nil {
		t.vScrollBar.Draw(dst)
	}
	if t.hScrollBar != nil {
		t.hScrollBar.Draw(dst)
	}
}

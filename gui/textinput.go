package gui

import (
	"image"
	"image/color"
	"log"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/atotto/clipboard"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/exp/textinput"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	textFieldHeightDefault = 24
	textPaddingXDefault    = 4
)

// TextField 支持输入法
//
// 参考 ebiten.exp.textinput.Field
// ebiten/examples/textinput/main.go
type TextField struct {
	bounds        image.Rectangle
	multilines    bool
	field         textinput.Field
	textPaddingX  int
	textHeight    int
	cursorCounter int
	selection0    int
}

func NewTextField(bounds image.Rectangle, multilines bool) *TextField {
	return &TextField{
		bounds:       bounds,
		multilines:   multilines,
		textPaddingX: textPaddingXDefault,
		textHeight:   textFieldHeightDefault,
		selection0:   -1,
	}
}

func (t *TextField) Contains(x, y int) bool {
	return image.Pt(x, y).In(t.bounds)
}

func (t *TextField) SetSelectionStartByCursorPosition(x, y int) bool {
	idx, ok := t.textIndexByCursorPosition(x, y)
	if !ok {
		return false
	}
	t.field.SetSelection(idx, idx)
	return true
}

func (t *TextField) textIndexByCursorPosition(x, y int) (int, bool) {
	if !t.Contains(x, y) {
		return 0, false
	}

	x -= t.bounds.Min.X
	y -= t.bounds.Min.Y
	px, py := t.textFieldPadding()
	x -= px
	y -= py
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	lineSpacingInPixels := int(uiFontFace.Metrics().HLineGap + uiFontFace.Metrics().HAscent + uiFontFace.Metrics().HDescent)
	var nlCount int
	var lineStart int
	var prevAdvance float64
	txt := t.field.Text()
	for i, r := range txt {
		var x0, x1 int
		currentAdvance := text.Advance(txt[lineStart:i], uiFontFace)
		if lineStart < i {
			x0 = int((prevAdvance + currentAdvance) / 2)
		}
		if r == '\n' {
			x1 = int(math.MaxInt32)
		} else if i < len(txt) {
			nextI := i + 1
			for !utf8.ValidString(txt[i:nextI]) {
				nextI++
			}
			nextAdvance := text.Advance(txt[lineStart:nextI], uiFontFace)
			x1 = int((currentAdvance + nextAdvance) / 2)
		} else {
			x1 = int(currentAdvance)
		}
		if x0 <= x && x < x1 && nlCount*lineSpacingInPixels <= y && y < (nlCount+1)*lineSpacingInPixels {
			return i, true
		}
		prevAdvance = currentAdvance

		if r == '\n' {
			nlCount++
			lineStart = i + 1
			prevAdvance = 0
		}
	}

	return len(txt), true
}

func (t *TextField) textFieldPadding() (int, int) {
	m := uiFontFace.Metrics()
	return t.textPaddingX, (t.textHeight - int(m.HLineGap+m.HAscent+m.HDescent)) / 2
}

func (t *TextField) Text() string {
	return t.field.Text()
}

func (t *TextField) SetText(text string) {
	t.field.SetTextAndSelection(text, len(text), len(text))
}

func (t *TextField) Focus() {
	t.field.Focus()
}

func (t *TextField) Blur() {
	t.field.Blur()
}

func (t *TextField) IsFocused() bool {
	return t.field.IsFocused()
}

func (t *TextField) Update() error {
	if !t.field.IsFocused() {
		return nil
	}

	t.cursorCounter = (t.cursorCounter + 1) % 360

	x, y := t.bounds.Min.X, t.bounds.Min.Y
	cx, cy := t.cursorPos()
	px, py := t.textFieldPadding()
	x += cx + px
	y += cy + py + int(uiFontFace.Metrics().HAscent)
	handled, err := t.field.HandleInput(x, y)
	if err != nil {
		return err
	}
	if handled {
		return nil
	}

	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyEnter), inpututil.IsKeyJustPressed(ebiten.KeyKPEnter):
		if t.multilines {
			text := t.field.Text()
			selectionStart, selectionEnd := t.field.Selection()
			text = text[:selectionStart] + "\n" + text[selectionEnd:]
			selectionStart += len("\n")
			selectionEnd = selectionStart
			t.field.SetTextAndSelection(text, selectionStart, selectionEnd)
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyBackspace):
		text := t.field.Text()
		selectionStart, selectionEnd := t.field.Selection()
		if selectionStart != selectionEnd {
			text = text[:selectionStart] + text[selectionEnd:]
		} else if selectionStart > 0 {
			// TODO: Remove a grapheme instead of a code point.
			_, l := utf8.DecodeLastRuneInString(text[:selectionStart])
			text = text[:selectionStart-l] + text[selectionEnd:]
			selectionStart -= l
		}
		selectionEnd = selectionStart
		t.field.SetTextAndSelection(text, selectionStart, selectionEnd)
	case inpututil.IsKeyJustPressed(ebiten.KeyDelete):
		text := t.field.Text()
		selectionStart, selectionEnd := t.field.Selection()
		if selectionStart != selectionEnd {
			text = text[:selectionStart] + text[selectionEnd:]
		} else if selectionStart < len(text) {
		}
		selectionEnd = selectionStart
		t.field.SetTextAndSelection(text, selectionStart, selectionEnd)
	case inpututil.IsKeyJustPressed(ebiten.KeyHome):
		text := t.field.Text()
		selectionStart, selectionEnd := t.field.Selection()
		if !ebiten.IsKeyPressed(ebiten.KeyShift) {
			selectionEnd = selectionStart
		}
		selectionStart = 0
		t.field.SetTextAndSelection(text, selectionStart, selectionEnd)
	case inpututil.IsKeyJustPressed(ebiten.KeyEnd):
		text := t.field.Text()
		selectionStart, selectionEnd := t.field.Selection()
		if !ebiten.IsKeyPressed(ebiten.KeyShift) {
			selectionStart = selectionEnd
		}
		selectionEnd = len(text)
		t.field.SetTextAndSelection(text, selectionStart, selectionEnd)
	case inpututil.IsKeyJustPressed(ebiten.KeyLeft):
		text := t.field.Text()
		selectionStart, selectionEnd := t.field.Selection()
		if selectionStart > 0 {
			_, l := utf8.DecodeLastRuneInString(text[:selectionStart])
			if ebiten.IsKeyPressed(ebiten.KeyShift) {
				if selectionStart != selectionEnd {
					if selectionEnd > t.selection0 {
						selectionEnd -= l //退右选
					} else {
						selectionStart -= l //左选
					}
				} else {
					t.selection0 = selectionStart //记录选中开始位置
					selectionStart -= l
				}
			} else {
				selectionStart -= l
				selectionEnd = selectionStart
			}
		} else {
			if !ebiten.IsKeyPressed(ebiten.KeyShift) {
				selectionEnd = selectionStart
			}
		}
		t.field.SetTextAndSelection(text, selectionStart, selectionEnd)
	case inpututil.IsKeyJustPressed(ebiten.KeyRight):
		text := t.field.Text()
		selectionStart, selectionEnd := t.field.Selection()
		if selectionEnd < len(text) {
			_, l := utf8.DecodeRuneInString(text[selectionEnd:])
			if ebiten.IsKeyPressed(ebiten.KeyShift) {
				if selectionStart != selectionEnd {
					if selectionStart < t.selection0 {
						selectionStart += l //退左选
					} else {
						selectionEnd += l //右选
					}
				} else {
					t.selection0 = selectionEnd //记录选中开始位置
					selectionEnd += l
				}
			} else {
				selectionEnd += l
				selectionStart = selectionEnd
			}
		} else {
			if !ebiten.IsKeyPressed(ebiten.KeyShift) {
				selectionStart = selectionEnd
			}
		}
		t.field.SetTextAndSelection(text, selectionStart, selectionEnd)
	case ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyX): //ctrl+x 禁if i.PasswordChar == "" {
		text := t.field.Text()
		selectionStart, selectionEnd := t.field.Selection()
		if selectionStart != selectionEnd {
			textX := text[selectionStart:selectionEnd]
			err := clipboard.WriteAll(textX)
			if err != nil {
				log.Printf("clipboard.WriteAll: %s", err.Error())
			}
			text = text[:selectionStart] + text[selectionEnd:]
		}
		selectionEnd = selectionStart
		t.field.SetTextAndSelection(text, selectionStart, selectionEnd)
	case ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyC): //ctrl+c 禁if i.PasswordChar == "" {
		text := t.field.Text()
		selectionStart, selectionEnd := t.field.Selection()
		if selectionStart != selectionEnd {
			text = text[selectionStart:selectionEnd]
		}
		err := clipboard.WriteAll(text)
		if err != nil {
			log.Printf("clipboard.WriteAll: %s", err.Error())
		}
	case ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyV):
		text := t.field.Text()
		textV, err := clipboard.ReadAll()
		if err != nil {
			log.Printf("clipboard.ReadAll: %s", err.Error())
		}
		selectionStart, selectionEnd := t.field.Selection()
		text = text[:selectionStart] + textV + text[selectionEnd:]
		_, l := utf8.DecodeRuneInString(textV)
		selectionStart += l
		selectionEnd = selectionStart
		t.field.SetTextAndSelection(text, selectionStart, selectionEnd)
	}

	if !t.multilines {
		orig := t.field.Text()
		new := strings.ReplaceAll(orig, "\n", "")
		if new != orig {
			selectionStart, selectionEnd := t.field.Selection()
			selectionStart -= strings.Count(orig[:selectionStart], "\n")
			selectionEnd -= strings.Count(orig[:selectionEnd], "\n")
			t.field.SetSelection(selectionStart, selectionEnd)
		}
	}

	return nil
}

func (t *TextField) cursorPos() (int, int) {
	var nlCount int
	lastNLPos := -1
	txt := t.field.TextForRendering()
	selectionStart, _ := t.field.Selection()
	if s, _, ok := t.field.CompositionSelection(); ok {
		selectionStart += s
	}
	txt = txt[:selectionStart]
	for i, r := range txt {
		if r == '\n' {
			nlCount++
			lastNLPos = i
		}
	}

	txt = txt[lastNLPos+1:]
	x := int(text.Advance(txt, uiFontFace))
	y := nlCount * int(uiFontFace.Metrics().HLineGap+uiFontFace.Metrics().HAscent+uiFontFace.Metrics().HDescent)
	return x, y
}

func (t *TextField) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(t.bounds.Min.X), float32(t.bounds.Min.Y), float32(t.bounds.Dx()), float32(t.bounds.Dy()), color.White, false)
	var clr color.Color = color.Black
	if t.field.IsFocused() {
		clr = color.RGBA{0, 0, 0xff, 0xff}
	}
	vector.StrokeRect(screen, float32(t.bounds.Min.X), float32(t.bounds.Min.Y), float32(t.bounds.Dx()), float32(t.bounds.Dy()), 1, clr, false)

	px, py := t.textFieldPadding()
	selectionStart, selectionEnd := t.field.Selection()

	if t.field.IsFocused() && selectionStart >= 0 {
		x, y := t.bounds.Min.X, t.bounds.Min.Y
		cx, cy := t.cursorPos()
		x += px + cx
		y += py + cy
		h := int(uiFontFace.Metrics().HLineGap + uiFontFace.Metrics().HAscent + uiFontFace.Metrics().HDescent)
		//draw selected text background
		if selectionStart != selectionEnd {
			txt := t.field.TextForRendering()
			selText := txt[selectionStart:selectionEnd]
			selX := int(text.Advance(txt[:selectionStart], uiFontFace))
			selW := int(text.Advance(selText, uiFontFace))
			vector.DrawFilledRect(screen, float32(px+selX), float32(y), float32(selW), float32(h), color.RGBA{0, 0, 0xff, 0x80}, false)
		}
		//draw cursor
		if t.cursorCounter%20 < 5 {
			vector.StrokeLine(screen, float32(x), float32(y), float32(x), float32(y+h), 1, color.Black, false)
		}
	}

	tx := t.bounds.Min.X + px
	ty := t.bounds.Min.Y + py
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(tx), float64(ty))
	op.ColorScale.ScaleWithColor(color.Black)
	op.LineSpacing = uiFontFace.Metrics().HLineGap + uiFontFace.Metrics().HAscent + uiFontFace.Metrics().HDescent
	text.Draw(screen, t.field.TextForRendering(), uiFontFace, op)
}

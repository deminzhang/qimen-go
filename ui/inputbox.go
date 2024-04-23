package ui

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image"
	"image/color"
	"strings"
)

const (
	textInputPadding = 8
)
const ()

type InputBox struct {
	BaseUI
	Rect         image.Rectangle
	Text         string
	DefaultText  string //无Text时默认灰文本
	MaxChars     int    //最大长度
	PasswordChar string //密文显示
	offsetX      int    //文本偏移 TODO
	Editable     bool

	autoLinefeed   bool //自动换行
	textHistory    []string
	textHistoryIdx int

	inputRunes     []rune
	cursorCounter  int
	cursorPos      int
	cursorSelect   int
	cursorNewFocus bool

	UIImage   *ebiten.Image
	ImageRect image.Rectangle

	onPressEnter func(i *InputBox)
}

func NewInputBox(rect image.Rectangle) *InputBox {
	return &InputBox{
		BaseUI:   BaseUI{Visible: true, EnableFocus: true},
		Rect:     rect,
		Editable: true,
		//default resource
		UIImage:   GetDefaultUIImage(),
		ImageRect: imageSrcRects[imageTypeTextBox],
	}
}

func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 3
		interval = 1
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func (i *InputBox) SetText(v interface{}) {
	i.Text = fmt.Sprintf("%v", v)
}

func (i *InputBox) Update() {
	if !i.Editable {
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if i.Rect.Min.X <= x && x < i.Rect.Max.X && i.Rect.Min.Y <= y && y < i.Rect.Max.Y {
			if i.Focused() {
				//cursorPos
				pos := len(i.Text)
				for ii := 0; ii < len(i.Text); ii++ {
					w := getFontWidth(uiFont, i.Text[:ii])
					if x < i.Rect.Min.X+textInputPadding+w {
						pos = ii
						break
					}
				}
				i.cursorSelect = pos
				i.cursorPos = pos
			} else {
				i.SetFocused(true)
				i.cursorPos = len(i.Text)
				i.cursorSelect = 0
				i.cursorNewFocus = true
			}
		} else {
			if i.Focused() {
				i.SetFocused(false)
			}
		}

		//} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	} else if inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) >= 2 {
		//drag cursorPos
		x, y := ebiten.CursorPosition()
		if i.Rect.Min.X <= x && x < i.Rect.Max.X && i.Rect.Min.Y <= y && y < i.Rect.Max.Y {
			if i.Focused() {
				pos := len(i.Text)
				for ii := 0; ii < len(i.Text); ii++ {
					w := getFontWidth(uiFont, i.Text[:ii])
					if x < i.Rect.Min.X+textInputPadding+w {
						pos = ii
						break
					}
				}
				i.cursorPos = pos
				if i.cursorNewFocus {
					i.cursorSelect = pos
					i.cursorNewFocus = false
				}
			}
		}
	}
	if i.Focused() {
		if i.MaxChars == 0 || i.MaxChars > len(i.Text) {
			s := string(ebiten.AppendInputChars(i.inputRunes[:0]))
			l := len(s)
			if l > 0 {
				left, right := i.selected()
				i.Text = i.Text[:left] + s + i.Text[right:]
				if left == right {
					i.cursorPos += l
				} else {
					i.cursorPos = left + l
				}
				i.cursorSelect = i.cursorPos
			}
		}

		// If the backspace key is pressed, remove one character.
		if repeatingKeyPressed(ebiten.KeyBackspace) {
			if len(i.Text) > 0 {
				left, right := i.selected()
				if left == right { //删左1个
					if left > 0 {
						i.Text = i.Text[:left-1] + i.Text[right:]
						i.cursorPos -= 1
					}
				} else { //删选中区
					i.Text = i.Text[:left] + i.Text[right:]
					i.cursorPos = left
				}
				i.cursorPos = min(i.cursorPos, len(i.Text))
				i.cursorPos = max(i.cursorPos, 0)
				i.cursorSelect = i.cursorPos
			}
		}
		if repeatingKeyPressed(ebiten.KeyDelete) {
			lenTxt := len(i.Text)
			if lenTxt > 0 {
				left, right := i.selected()
				if left == right { //删右1个
					if left < lenTxt {
						i.Text = i.Text[:left] + i.Text[right+1:]
					}
				} else { //删选中区
					i.Text = i.Text[:left] + i.Text[right:]
					i.cursorPos = left
				}
				i.cursorPos = min(i.cursorPos, len(i.Text))
				i.cursorPos = max(i.cursorPos, 0)
				i.cursorSelect = i.cursorPos
			}
		}
		if repeatingKeyPressed(ebiten.KeyLeft) {
			if ebiten.IsKeyPressed(ebiten.KeyShift) {
				if i.cursorPos > 0 {
					i.cursorPos--
				}
			} else {
				left, right := i.selected()
				if left != right {
					i.cursorPos = left
				} else {
					if i.cursorPos > 0 {
						i.cursorPos--
					}
				}
				i.cursorSelect = i.cursorPos
			}
		}
		if repeatingKeyPressed(ebiten.KeyRight) {
			if ebiten.IsKeyPressed(ebiten.KeyShift) {
				if i.cursorPos < len(i.Text) {
					i.cursorPos++
				}
			} else {
				left, right := i.selected()
				if left != right {
					i.cursorPos = right
				} else {
					if i.cursorPos < len(i.Text) {
						i.cursorPos++
					}
				}
				i.cursorSelect = i.cursorPos
			}
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			x, y := ebiten.CursorPosition()
			if i.Rect.Min.X <= x && x < i.Rect.Max.X && i.Rect.Min.Y <= y && y < i.Rect.Max.Y {
				clipboard.WriteAll(i.Text)
			}
		}
		if i.PasswordChar == "" {
			//ctrl+x
			if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyX) {
				left, right := i.selected()
				clipboard.WriteAll(i.Text[left:right])

				lenTxt := len(i.Text)
				if lenTxt > 0 {
					left, right := i.selected()
					if left == right { //删右1个
						if left < lenTxt {
							i.Text = i.Text[:left] + i.Text[right+1:]
						}
					} else { //删选中区
						i.Text = i.Text[:left] + i.Text[right:]
						i.cursorPos = left
					}
					i.cursorPos = min(i.cursorPos, len(i.Text))
					i.cursorPos = min(i.cursorPos, 0)
					i.cursorSelect = i.cursorPos
				}
			}
			//ctrl+c
			if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyC) {
				left, right := i.selected()
				clipboard.WriteAll(i.Text[left:right])
			}
		}
		//ctrl+v
		if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyV) {
			s, err := clipboard.ReadAll()
			if err == nil {
				left, right := i.selected()
				i.Text = i.Text[:left] + s + i.Text[right:]
				if left == right {
					i.cursorPos += len(s)
				} else {
					i.cursorPos = left + len(s)
				}
				i.cursorSelect = i.cursorPos
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			length := len(i.textHistory)
			if length > 0 {
				if i.textHistoryIdx < length-1 {
					i.textHistoryIdx++
				}
				i.Text = i.textHistory[i.textHistoryIdx]
				i.cursorPos = len(i.Text)
				i.cursorSelect = i.cursorPos
			}
		} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			if len(i.textHistory) > 0 {
				if i.textHistoryIdx > 0 {
					i.textHistoryIdx--
				}
				i.Text = i.textHistory[i.textHistoryIdx]
				i.cursorPos = len(i.Text)
				i.cursorSelect = i.cursorPos
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) {
		if i.onPressEnter != nil {
			i.onPressEnter(i)
		}
	}
	i.cursorCounter++
}

var textSelectColor = color.RGBA{0, 0, 200, 128}

func (i *InputBox) selected() (int, int) {
	left, right := i.cursorPos, i.cursorSelect
	if left > right {
		left, right = i.cursorSelect, i.cursorPos
	}
	return left, right
}

func (i *InputBox) Draw(dst *ebiten.Image) {
	if !i.Visible {
		return
	}
	drawNinePatches(dst, i.UIImage, i.Rect, i.ImageRect)

	//drawText
	x := i.Rect.Min.X + textInputPadding //居左  //居中 + (i.Rect.Dx()-w)/2
	y := i.Rect.Max.Y - (i.Rect.Dy()-uiFontMHeight)/2

	if i.Text == "" && i.DefaultText != "" && !i.Focused() {
		text.Draw(dst, i.DefaultText, uiFont, x, y, color.Gray{Y: 128})
	} else {
		t := i.Text
		if i.PasswordChar != "" && i.Text != "" {
			t = strings.Repeat(i.PasswordChar[:1], len(t))
		}
		//drawText todo 自动换行
		text.Draw(dst, t, uiFont, x, y, color.Black)
		if i.autoLinefeed {

		} else {

		}
		//draw cursor
		if i.Focused() {
			i.cursorPos = min(i.cursorPos, len(t))
			i.cursorSelect = min(i.cursorSelect, len(t))
			//drawSelect
			if i.cursorPos != i.cursorSelect {
				left, right := i.selected()
				wl, wr := getFontSelectWidth(uiFont, t, left, right)
				ebitenutil.DrawRect(dst, float64(x+wl), float64(i.Rect.Min.Y+4),
					float64(wr-wl), float64(i.Rect.Dy()-8), textSelectColor)
			}
			if i.cursorCounter%10 < 5 {
				w := getFontWidth(uiFont, t[:i.cursorPos])
				//太矮 text.Draw(dst, "|", uiFont, x+w, y, color.Black)
				//太窄 ebitenutil.DrawLine(dst, float64(x+w), float64(i.Rect.Min.Y+4), float64(x+w), float64(i.Rect.Min.Y+i.Rect.Dy()-4), color.Black)
				ebitenutil.DrawRect(dst, float64(x+w), float64(i.Rect.Min.Y+4),
					float64(2), float64(i.Rect.Dy()-8), color.Black)
			}
		}
	}
}

func (i *InputBox) SetOnPressEnter(f func(*InputBox)) {
	i.onPressEnter = f
}

func (i *InputBox) AppendTextHistory(txt string) {
	i.textHistory = append([]string{txt}, i.textHistory...)
	i.textHistoryIdx = -1
}

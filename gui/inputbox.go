package gui

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"strings"
)

const (
	defaultTextInputPadding = 8
)

var textSelectColor = color.RGBA{B: 200, A: 128}

type InputBox struct {
	BaseUI
	textRune     []rune
	DefaultText  string //无Text时默认灰文本
	MaxChars     int    //最大长度
	PasswordChar string //密文显示
	offsetX      int    //文本偏移 TODO
	textPadding  int
	Editable     bool
	Selectable   bool

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

func NewInputBox(x, y, w, h int) *InputBox {
	return &InputBox{
		BaseUI:     BaseUI{X: x, Y: y, W: w, H: h, Visible: true, EnableFocus: true},
		Editable:   true,
		Selectable: true,
		//default resource
		textPadding: defaultTextInputPadding,
		UIImage:     GetDefaultUIImage(),
		ImageRect:   imageSrcRects[imageTypeTextBox],
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

func (i *InputBox) Text() string {
	return string(i.textRune)
}
func (i *InputBox) SetText(v interface{}) {
	i.textRune = []rune(fmt.Sprintf("%v", v))
}

func (i *InputBox) Update() {
	if !i.Selectable {
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		x, y := i.GetWorldXY()
		if x <= mx && mx < x+i.W && y <= my && my < y+i.H {
			if i.Focused() {
				pos := len(i.textRune)
				for ii := 0; ii < pos; ii++ {
					w := getFontWidth(uiFont, string(i.textRune[:ii]))
					if mx < x+i.textPadding+w {
						pos = ii
						break
					}
				}
				i.cursorSelect = pos
				i.cursorPos = pos
			} else {
				i.SetFocused(true)
				i.cursorPos = len(i.textRune)
				i.cursorSelect = 0
				i.cursorNewFocus = true
			}
		} else {
			if i.Focused() {
				i.SetFocused(false)
			}
		}

		//} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	} else if inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) >= 2 { //拖动选中
		//drag cursorPos
		mx, my := ebiten.CursorPosition()
		x, y := i.GetWorldXY()
		if x <= mx && mx < x+i.W && y <= my && my < y+i.H {
			if i.Focused() {
				pos := len(i.textRune)
				for ii := 0; ii < pos; ii++ {
					w := getFontWidth(uiFont, string(i.textRune[:ii]))
					if mx < x+i.textPadding+w {
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
		sRune := ebiten.AppendInputChars(i.inputRunes[:0])
		l := len(sRune)
		if l > 0 {
			left, right := i.cursorSelected()
			if left != right { //删选中区
				i.textRune = []rune(string(i.textRune[:left]) + string(i.textRune[right:]))
				i.cursorPos = left
				right = left
			}
			if i.MaxChars <= 0 || i.MaxChars > len(i.textRune) {
				sRune = append(sRune, i.textRune[right:]...)
				i.textRune = append(i.textRune[:left], sRune...)
				if left == right {
					i.cursorPos += l
				} else {
					i.cursorPos = left + l
				}
				i.cursorSelect = i.cursorPos
			}
		}

		// If the backspace key is pressed, remove one character.
		if i.Editable && repeatingKeyPressed(ebiten.KeyBackspace) {
			if len(i.textRune) > 0 {
				left, right := i.cursorSelected()
				if left == right { //删左1个
					if left > 0 {
						ll := len(i.textRune[:left])
						i.cursorPos -= 1
						i.textRune = []rune(string(i.textRune[:ll-1]) + string(i.textRune[ll:]))
					}
				} else { //删选中区
					i.textRune = []rune(string(i.textRune[:left]) + string(i.textRune[right:]))
					i.cursorPos = left
				}
				i.cursorPos = min(i.cursorPos, len(i.textRune))
				i.cursorPos = max(i.cursorPos, 0)
				i.cursorSelect = i.cursorPos
			}
		}
		if i.Editable && repeatingKeyPressed(ebiten.KeyDelete) {
			l := len(i.textRune)
			if l > 0 {
				left, right := i.cursorSelected()
				if left == right { //删右1个
					if left < l {
						i.textRune = []rune(string(i.textRune[:left]) + string(i.textRune[left+1:]))
					}
				} else { //删选中区
					i.textRune = []rune(string(i.textRune[:left]) + string(i.textRune[right:]))
					i.cursorPos = left
				}
				i.cursorPos = min(i.cursorPos, len(i.textRune))
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
				left, right := i.cursorSelected()
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
				if i.cursorPos < len(i.textRune) {
					i.cursorPos++
				}
			} else {
				left, right := i.cursorSelected()
				if left != right {
					i.cursorPos = right
				} else {
					if i.cursorPos < len(i.textRune) {
						i.cursorPos++
					}
				}
				i.cursorSelect = i.cursorPos
			}
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			mx, my := ebiten.CursorPosition()
			x, y := i.GetWorldXY()
			if x <= mx && mx < x+i.W && y <= my && my < y+i.H {
				clipboard.WriteAll(string(i.textRune))
			}
		}
		if i.PasswordChar == "" {
			//ctrl+x
			if i.Editable && ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyX) {
				left, right := i.cursorSelected()
				err := clipboard.WriteAll(string(i.textRune[left:right]))
				if err == nil {
					l := len(i.textRune)
					if l > 0 {
						if left == right { //删右1个
							if left < l {
								i.textRune = append(i.textRune[:left], i.textRune[right+1:]...)
							}
						} else { //删选中区
							i.textRune = append(i.textRune[:left], i.textRune[right:]...)
							i.cursorPos = left
						}
						i.cursorPos = min(i.cursorPos, len(i.textRune))
						i.cursorPos = min(i.cursorPos, 0)
						i.cursorSelect = i.cursorPos
					}
				}
			}
			//ctrl+c
			if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyC) {
				left, right := i.cursorSelected()
				clipboard.WriteAll(string(i.textRune[left:right]))
			}
		}
		//ctrl+v
		if i.Editable && ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyV) {
			s, err := clipboard.ReadAll()
			if err == nil {
				left, right := i.cursorSelected()
				r := append([]rune(s), i.textRune[right:]...)
				i.textRune = append(i.textRune[:left], r...)
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
				i.textRune = []rune(i.textHistory[i.textHistoryIdx])
				i.cursorPos = len(i.textRune)
				i.cursorSelect = i.cursorPos
			}
		} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			if len(i.textHistory) > 0 {
				if i.textHistoryIdx > 0 {
					i.textHistoryIdx--
				}
				i.textRune = []rune(i.textHistory[i.textHistoryIdx])
				i.cursorPos = len(i.textRune)
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

func (i *InputBox) cursorSelected() (int, int) {
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
	drawNinePatches(dst, i.UIImage, image.Rect(0, 0, i.W, i.H), i.ImageRect)

	//drawText
	x := i.textPadding //居左  //居中 + (i.Rect.Dx()-w)/2
	y := i.H - (i.H-uiFontMHeight)/2

	if len(i.textRune) == 0 && i.DefaultText != "" && !i.Focused() {
		text.Draw(dst, i.DefaultText, uiFont, x, y, color.Gray{Y: 128})
	} else {
		r := i.textRune
		if len(r) == 0 && i.PasswordChar != "" {
			r = []rune(strings.Repeat(i.PasswordChar[:1], len(r)))
		}
		s := string(r)
		//drawText todo 自动换行
		text.Draw(dst, s, uiFont, x, y, color.Black)
		if i.autoLinefeed {
		} else {
		}
		//draw cursor
		if i.Focused() {
			i.cursorPos = min(i.cursorPos, len(r))
			i.cursorSelect = min(i.cursorSelect, len(r))
			//drawSelect
			if i.cursorPos != i.cursorSelect {
				left, right := i.cursorSelected()
				s1 := string(r[:left])
				s2 := string(r[:right])
				wl, wr := getFontSelectWidth(uiFont, s, len(s1), len(s2))
				vector.DrawFilledRect(dst, float32(x+wl), float32(i.Y+4),
					float32(wr-wl), float32(i.H-8), textSelectColor, false)
			}
			if i.cursorCounter%10 < 5 {
				w := getFontWidth(uiFont, string(r[:i.cursorPos]))
				//太矮 text.Draw(dst, "|", uiFont, x+w, y, color.Black)
				//太窄 ebitenutil.DrawLine(dst, float64(x+w), float64(4), float64(x+w), float64(i.H-4), color.Black)
				vector.DrawFilledRect(dst, float32(x+w), 4, 2, float32(i.H-8), color.Black, false)
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

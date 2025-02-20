package gui

import (
	"fmt"
	"image"
	"image/color"

	"github.com/atotto/clipboard"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	defaultTextInputPadding = 8
)

var textSelectColor = color.RGBA{B: 200, A: 128}

// InputBox 不支持多行输入法 只能输英文字符, TODO TextField支持输入法
type InputBox struct {
	BaseUI
	TextField    *TextField
	textRune     []rune
	DefaultText  string //无Text时默认灰文本
	MaxChars     int    //最大长度
	PasswordChar string //密文显示
	textPadding  int
	Editable     bool
	Selectable   bool

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
	onLostFocus  func(i *InputBox)
}

func NewInputBox(x, y, w, h int) *InputBox {
	return &InputBox{
		BaseUI:     BaseUI{X: x, Y: y, W: w, H: h, Visible: true, EnableFocus: true},
		TextField:  NewTextField(image.Rect(0, 0, w, h), false),
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
	return i.TextField.Text()
}

func (i *InputBox) SetText(v interface{}) {
	str := fmt.Sprintf("%v", v)
	i.TextField.SetText(str)
}

func (i *InputBox) Update() {
	i.BaseUI.Update()
	if !i.Selectable {
		return
	}
	tf := i.TextField
	if tf == nil {
		return
	}
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(ids) > 0 {
		var x, y int
		if len(ids) > 0 {
			x, y = ebiten.TouchPosition(ids[0])
		} else {
			x, y = ebiten.CursorPosition()
		}
		wx, wy := i.GetWorldXY()
		if tf.Contains(x-wx, y-wy) {
			tf.Focus()
			tf.SetSelectionStartByCursorPosition(x-wx, y-wy)
		} else {
			tf.Blur()
			if i.Focused() {
				if i.onLostFocus != nil {
					i.onLostFocus(i)
				}
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) {
		if i.onPressEnter != nil {
			i.onPressEnter(i)
		}
	}
	if err := tf.Update(); err != nil {
		fmt.Println(err)
		return
	}

	x, y := ebiten.CursorPosition()
	wx, wy := i.GetWorldXY()
	if tf.Contains(x-wx, y-wy) {
		ebiten.SetCursorShape(ebiten.CursorShapeText)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}
}

// Deprecated: 旧版不支持输入法
func (i *InputBox) UpdateX() {
	i.BaseUI.Update()
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
				if i.onLostFocus != nil {
					i.onLostFocus(i)
				}
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
	i.TextField.Draw(dst)
}

func (i *InputBox) SetOnPressEnter(f func(*InputBox)) {
	i.onPressEnter = f
}

func (i *InputBox) SetOnLostFocus(f func(*InputBox)) {
	i.onLostFocus = f
}

func (i *InputBox) AppendTextHistory(txt string) {
	i.textHistory = append([]string{txt}, i.textHistory...)
	i.textHistoryIdx = -1
}

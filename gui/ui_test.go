package gui_test

import (
	. "github.com/deminzhang/qimen-go/gui"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	_ "image/png"

	"log"
	"testing"
)

const (
	screenWidth  = 640
	screenHeight = 480
	TPSRate      = 10
)

type Game struct {
	button1    *Button
	button2    *Button
	checkBox   *CheckBox
	textBoxLog *TextBox
}

func NewGame() *Game {
	g := &Game{}
	g.button1 = NewButton(16, 16, 144, 48, "Button 1")
	g.button2 = NewButton(160, 16, 288, 48, "Button 2")
	g.checkBox = &CheckBox{
		BaseUI: BaseUI{
			X: 16,
			Y: 64,
		},
		Text: "Check Box!",
	}
	g.textBoxLog = NewTextBox(16, 96, 608, 368)

	g.button1.SetOnClick(func() {
		g.textBoxLog.AppendTextLn("Button 1 Pressed")
	})
	g.button2.SetOnClick(func() {
		g.textBoxLog.AppendTextLn("Button 2 Pressed")
	})
	g.checkBox.SetOnCheckChanged(func(c *CheckBox) {
		msg := "Check box check changed"
		if c.Checked() {
			msg += " (Selected)"
		} else {
			msg += " (Unchecked)"
		}
		g.textBoxLog.AppendTextLn(msg)
	})
	return g
}

func (g *Game) Update() error {
	g.button1.Update()
	g.button2.Update()
	g.checkBox.Update()
	g.textBoxLog.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0xeb, G: 0xeb, B: 0xeb, A: 0xff})
	g.button1.Draw(screen)
	g.button2.Draw(screen)
	g.checkBox.Draw(screen)
	g.textBoxLog.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("uiBase (Ebiten Demo)")
	ebiten.SetTPS(TPSRate)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

func TestUI(t *testing.T) {
	main()
}

package gui

import "github.com/hajimehoshi/ebiten/v2"

type Dragger struct {
	mouseDown                bool
	lastCursorX, lastCursorY int
	offsetX, offsetY         int
}

func (d *Dragger) Update() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if !d.mouseDown {
			d.mouseDown = true
			d.lastCursorX, d.lastCursorY = x, y
			d.offsetX, d.offsetY = 0, 0
		}
		d.offsetX = x - d.lastCursorX
		d.offsetY = y - d.lastCursorY
		d.lastCursorX, d.lastCursorY = x, y
	} else {
		d.mouseDown = false
		d.offsetX, d.offsetY = 0, 0
	}
}

func (d *Dragger) Offset() (int, int) {
	return d.offsetX, d.offsetY
}

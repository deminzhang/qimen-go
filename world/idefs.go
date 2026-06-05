package world

import "github.com/hajimehoshi/ebiten/v2"

type IComponent interface {
	Update()
	Draw(*ebiten.Image)
}

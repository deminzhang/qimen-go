package world

import "github.com/hajimehoshi/ebiten/v2"

type Unit struct {
	X, Y   float64
	Tx, Ty float64
	Type   string
	Name   string
	Camp   int
	Hp     int
	MaxHp  int
	Atk    int
	Def    int
	Spd    int
	AtkDis int

	Img   *ebiten.Image
	ImgOp *ebiten.DrawImageOptions
}

func NewUnit(x, y float64, unitType, name string, camp int) *Unit {
	return &Unit{
		X:      x,
		Y:      y,
		Type:   unitType,
		Name:   name,
		Camp:   camp,
		Tx:     x,
		Ty:     y,
		Hp:     100,
		MaxHp:  100,
		Atk:    10,
		Def:    5,
		Spd:    5,
		AtkDis: 1,
	}
}

func (u *Unit) Draw(screen *ebiten.Image) {
	screen.DrawImage(u.Img, u.ImgOp)
}

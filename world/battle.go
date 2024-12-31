package world

import (
	"github.com/deminzhang/qimen-go/graphic"
	"github.com/hajimehoshi/ebiten/v2"
)

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

	Sprite *Sprite
}

func NewUnit(x, y float64, unitType, name string, camp int) *Unit {
	sprite := NewSprite(graphic.NewArmyImage(name, 32, 1), nil)
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
		Sprite: sprite,
	}
}

func (u *Unit) Update() {

}

func (u *Unit) Draw(screen *ebiten.Image) {
	if u.Sprite == nil {
		u.Sprite.Draw(screen)
	}
}

type Battle struct {
	HostUnit  []*Unit
	GuestUnit []*Unit
}

func NewBattle() *Battle {
	return &Battle{
		HostUnit:  []*Unit{},
		GuestUnit: []*Unit{},
	}
}

func (b *Battle) AddHostUnit(u *Unit) {
	b.HostUnit = append(b.HostUnit, u)
}

func (b *Battle) AddGuestUnit(u *Unit) {
	b.GuestUnit = append(b.GuestUnit, u)
}

func (b *Battle) Update() {
	for _, u := range b.HostUnit {
		u.Update()
	}
	for _, u := range b.GuestUnit {
		u.Update()
	}
}

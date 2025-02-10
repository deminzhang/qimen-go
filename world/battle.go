package world

import (
	"github.com/deminzhang/qimen-go/graphic"
	"github.com/deminzhang/qimen-go/qimen"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"time"
)

type Unit struct {
	X, Y   float32
	Tx, Ty float32
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

func NewUnit(x, y float32, unitType, name string, camp int) *Unit {
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

var (
	CampColor = map[int]color.Color{
		0: colorWhite,
		1: colorRed,
	}
)

func NewArmy(x, y float32, size int, name string, camp int) *Unit {
	u := NewUnit(x, y, "army", name, camp)
	sprite := NewSprite(graphic.NewArmyImage(name, size, 1), CampColor[camp])
	sprite.x = int(x) - size/2
	sprite.y = int(y) - size/2
	u.Sprite = sprite
	return u
}
func NewCamp(x, y float32, size int, name string, camp int) *Unit {
	u := NewUnit(x, y, "camp", name, camp)
	sprite := NewSprite(graphic.NewCampImage(size), CampColor[camp])
	sprite.x = int(x) - size/2
	sprite.y = int(y) - size/2
	u.Sprite = sprite
	return u
}

func (u *Unit) Update(now, delta int64) {
	if u.Tx != u.X || u.Ty != u.Y {
		//TODO move to TargetXY
	}
}

func (u *Unit) Draw(screen *ebiten.Image) {
	if u.Sprite != nil {
		u.Sprite.Draw(screen)
	}
}

type Battle struct {
	lastUpdateTime int64
	HostUnit       []*Unit
	GuestUnit      []*Unit
	CampM          *ebiten.Image
	Camp           *ebiten.Image
	Army           *ebiten.Image
	ArmyA          *ebiten.Image
	inited         bool
}

func NewBattle() *Battle {
	return &Battle{
		HostUnit:  []*Unit{},
		GuestUnit: []*Unit{},
		CampM:     graphic.NewCampImage(64),
		Camp:      graphic.NewCampImage(32),
		Army:      graphic.NewArmyImage("庚", 32, 0),
		ArmyA:     graphic.NewArmyImage("兵", 32, 1),
	}
}

func (b *Battle) AddHostUnit(u *Unit) {
	b.HostUnit = append(b.HostUnit, u)
}

func (b *Battle) AddGuestUnit(u *Unit) {
	b.GuestUnit = append(b.GuestUnit, u)
}

func (b *Battle) InitBattle(q *QMShow) {
	var x, y float32
	//九宫
	for i := 1; i <= 9; i++ {
		x, y = q.GetInCampPos(i)
		if i == 5 {
			u := NewCamp(x, y, 64, "营", 0)
			b.AddHostUnit(u)

			//x, y = q.GetInBornPos(i)
			//u = NewArmy(x, y, 64, "帅", 0)
			b.AddHostUnit(u) //校位
		} else {
			u := NewCamp(x, y, 32, "营", 0)
			b.AddHostUnit(u)

			x, y = q.GetInBornPos(i)
			u = NewArmy(x, y, 32, "兵", 0)
			b.AddHostUnit(u)

			//x, y = q.GetInArmyPos(i)
			//u = NewArmy(x, y, 32, "兵", 1)
			//b.AddGuestUnit(u)//标位
		}
	}
	//12宫
	for i := 1; i <= 12; i++ {
		y, x = q.GetOutCampPos(i)
		u := NewCamp(x, y, 32, "营", 1)
		b.AddGuestUnit(u)
		y, x = q.GetOutCampBornPos(i)
		u = NewArmy(x, y, 32, "兵", 1)
		ti := qimen.ZhiGong2Gong8[i]
		x, y = q.GetInArmyPos(ti)
		u.Tx, u.Ty = x, y
		b.AddGuestUnit(u)
	}
}

func (b *Battle) Update() {
	now := time.Now().UnixMilli()
	delta := now - b.lastUpdateTime
	if delta < 2000 {
		return
	}
	if !b.inited {
		b.InitBattle(ThisGame.qiMen)
		b.inited = true
		return
	}
	b.lastUpdateTime = now
	for _, u := range b.HostUnit {
		u.Update(now, delta)
	}
	for _, u := range b.GuestUnit {
		u.Update(now, delta)
	}
}

func (b *Battle) Draw(dst *ebiten.Image, q *QMShow) {
	for _, u := range b.HostUnit {
		u.Draw(dst)
	}
	for _, u := range b.GuestUnit {
		u.Draw(dst)
	}
}

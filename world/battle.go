package world

import (
	"github.com/deminzhang/qimen-go/graphic"
	"github.com/deminzhang/qimen-go/util"
	"github.com/deminzhang/qimen-go/xuan"
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
	State  int
	Speed  float32
	Hp     int
	MaxHp  int
	Atk    int
	Def    int
	AtkDis int
	Size   int

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
		Speed:  0,
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
	u.Size = size
	u.Speed = 2
	sprite := NewSprite(graphic.NewArmyImage(name, size, 1), CampColor[camp])
	sprite.x = int(x) - size/2
	sprite.y = int(y) - size/2
	u.Sprite = sprite
	return u
}

func NewCamp(x, y float32, size int, name string, camp int) *Unit {
	u := NewUnit(x, y, "camp", name, camp)
	u.Size = size
	sprite := NewSprite(graphic.NewCampImage(size), CampColor[camp])
	sprite.x = int(x) - size/2
	sprite.y = int(y) - size/2
	u.Sprite = sprite
	return u
}

func (u *Unit) Update(now, delta int64) {
	if u.Tx != u.X || u.Ty != u.Y {
		from := util.Vec2[float32]{X: u.X, Y: u.Y}
		newPos := from.MoveTowards(u.Tx, u.Ty, float64(u.Speed))
		u.X, u.Y = newPos.X, newPos.Y
		u.Sprite.x = int(u.X) - u.Size/2
		u.Sprite.y = int(u.Y) - u.Size/2
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
	inited         bool
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
		ti := xuan.ZhiGong2Gong8[i]
		x, y = q.GetInArmyPos(ti)
		u.Tx, u.Ty = x, y
		b.AddGuestUnit(u)
	}
}

func (b *Battle) Update() {
	now := time.Now().UnixMilli()
	delta := now - b.lastUpdateTime
	if delta < 200 {
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

func (b *Battle) Draw(dst *ebiten.Image) {
	for _, u := range b.HostUnit {
		u.Draw(dst)
	}
	for _, u := range b.GuestUnit {
		u.Draw(dst)
	}
}

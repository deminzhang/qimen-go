package world

import (
	"fmt"
	"github.com/deminzhang/qimen-go/asset"
	"github.com/deminzhang/qimen-go/graphic"
	"github.com/deminzhang/qimen-go/qimen"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"math"
)

const (
	meiHuaGuaSize = 25
)

type MeiHua struct {
	X, Y         int
	GuaUpIdx     uint8  //上卦序号
	GuaDownIdx   uint8  //下卦序号
	YaoChangeIdx uint8  //变爻
	StartType    string //起卦方式

	GuaOrigin string //本卦
	GuaUp     string //上卦
	GuaDown   string //下卦

	GuaProcess     string //互卦
	GuaUpProcess   string //互卦上卦
	GuaDownProcess string //互卦下卦

	GuaChange     string //变卦
	GuaUpChange   string //变卦上卦
	GuaDownChange string //变卦下卦

	Mover     *Sprite
	GuaSprite []*Sprite
}

func NewMeiHua(x, y int, upIdx, downIdx, change uint) *MeiHua {
	mh := &MeiHua{X: x, Y: y}
	mh.Reset(upIdx, downIdx, change)
	return mh
}

// 互卦
func huGua(upIdx, downIdx uint8) (uint8, uint8) {
	up := qimen.Diagrams8Origin[upIdx]
	down := qimen.Diagrams8Origin[downIdx]
	upB := qimen.Diagrams8Bin[up]
	downB := qimen.Diagrams8Bin[down]
	upN := (upB&0b11)<<1 + downB>>2
	downN := (upB&0b1)<<2 + downB>>1
	return upN, downN
}

func (m *MeiHua) Reset(upIdx, downIdx, change uint) {
	upIdx = upIdx % 8
	if upIdx <= 0 {
		upIdx += 8
	}
	downIdx = downIdx % 8
	if downIdx <= 0 {
		downIdx += 8
	}
	change = change % 6
	if change <= 0 {
		change += 6
	}
	if uint8(upIdx) == m.GuaUpIdx && uint8(downIdx) == m.GuaDownIdx && uint8(change) == m.YaoChangeIdx {
		return
	}
	up := qimen.Diagrams8Origin[uint8(upIdx)]
	down := qimen.Diagrams8Origin[uint8(downIdx)]
	ori := qimen.Diagrams64FullName[uint8(upIdx*10+downIdx)]
	//互卦
	huUpB, huDownB := huGua(uint8(upIdx), uint8(downIdx))
	huUp := qimen.Diagrams8FromBin[huUpB]
	huDown := qimen.Diagrams8FromBin[huDownB]
	pro := qimen.Diagrams64FullName[(qimen.Diagrams8IdxOrigin[huUp]*10 + qimen.Diagrams8IdxOrigin[huDown])]
	//变卦
	upB := qimen.Diagrams8Bin[up]
	downB := qimen.Diagrams8Bin[down]
	if change > 3 {
		upB ^= 1 << (change - 3 - 1)
	} else {
		downB ^= 1 << (change - 1)
	}
	//变卦
	cUp := qimen.Diagrams8FromBin[upB]
	cDown := qimen.Diagrams8FromBin[downB]
	changeGua := qimen.Diagrams64FullName[(qimen.Diagrams8IdxOrigin[cUp]*10 + qimen.Diagrams8IdxOrigin[cDown])]

	m.GuaUpIdx = uint8(upIdx)
	m.GuaDownIdx = uint8(downIdx)
	m.YaoChangeIdx = uint8(change)
	m.GuaUp = up
	m.GuaDown = down
	m.GuaOrigin = ori
	m.GuaProcess = pro
	m.GuaUpProcess = huUp
	m.GuaDownProcess = huDown
	m.GuaChange = changeGua
	m.GuaUpChange = cUp
	m.GuaDownChange = cDown

	for _, sprite := range m.GuaSprite {
		ThisGame.RemoveSprite(sprite)
	}
	m.GuaSprite = nil
}
func (m *MeiHua) Update() {
	if ThisGame.qmGame.Lunar != nil {
		yz := qimen.ZhiIdx[ThisGame.qmGame.Lunar.GetYearZhiExact()]
		mz := ThisGame.qmGame.Lunar.GetMonth()
		dz := ThisGame.qmGame.Lunar.GetDay()
		hz := qimen.ZhiIdx[ThisGame.qmGame.Lunar.GetTimeZhi()]
		up := yz + mz + dz
		down := yz + mz + dz + hz
		m.Reset(uint(up), uint(down), uint(down))
	}
	dis := int(math.Round(float64(meiHuaGuaSize) * 1.25))
	if m.Mover == nil {
		m.Mover = NewSprite(graphic.NewRectImage(10), colorGray)
		m.Mover.onMove = func(sx, sy, dx, dy int) {
			m.X += dx
			m.Y += dy
			if m.GuaSprite != nil {
				m.GuaSprite[0].MoveTo(m.X+32, m.Y+64)
				m.GuaSprite[1].MoveTo(m.X+32, m.Y+64+dis)
				m.GuaSprite[2].MoveTo(m.X+32+64, m.Y+64)
				m.GuaSprite[3].MoveTo(m.X+32+64, m.Y+64+dis)
				m.GuaSprite[4].MoveTo(m.X+32+128, m.Y+64)
				m.GuaSprite[5].MoveTo(m.X+32+128, m.Y+64+dis)
			}
		}
		ThisGame.AddSprite(m.Mover)
		m.Mover.MoveTo(m.X, m.Y)
	}
	if m.GuaSprite == nil {
		cx, cy := m.X+32, m.Y+64
		m.GuaSprite = make([]*Sprite, 6)
		m.GuaSprite[0] = NewSprite(graphic.NewBaGuaImage(m.GuaUp, meiHuaGuaSize), color5Xing[qimen.DiagramsWuxing[m.GuaUp]])
		m.GuaSprite[1] = NewSprite(graphic.NewBaGuaImage(m.GuaDown, meiHuaGuaSize), color5Xing[qimen.DiagramsWuxing[m.GuaDown]])
		m.GuaSprite[0].MoveTo(cx, cy)
		m.GuaSprite[1].MoveTo(cx, cy+dis)
		m.GuaSprite[2] = NewSprite(graphic.NewBaGuaImage(m.GuaUpProcess, meiHuaGuaSize), color5Xing[qimen.DiagramsWuxing[m.GuaUpProcess]])
		m.GuaSprite[3] = NewSprite(graphic.NewBaGuaImage(m.GuaDownProcess, meiHuaGuaSize), color5Xing[qimen.DiagramsWuxing[m.GuaDownProcess]])
		m.GuaSprite[2].MoveTo(cx+64, cy)
		m.GuaSprite[3].MoveTo(cx+64, cy+dis)
		m.GuaSprite[4] = NewSprite(graphic.NewBaGuaImage(m.GuaUpChange, meiHuaGuaSize), color5Xing[qimen.DiagramsWuxing[m.GuaUpChange]])
		m.GuaSprite[5] = NewSprite(graphic.NewBaGuaImage(m.GuaDownChange, meiHuaGuaSize), color5Xing[qimen.DiagramsWuxing[m.GuaDownChange]])
		m.GuaSprite[4].MoveTo(cx+128, cy)
		m.GuaSprite[5].MoveTo(cx+128, cy+dis)
	}
}

func (m *MeiHua) Draw(dst *ebiten.Image) {
	vector.StrokeRect(dst, float32(m.X), float32(m.Y), 240, 140, .5, colorGray, true)
	m.Mover.Draw(dst)
	ft14, _ := asset.GetDefaultFontFace(14)
	text.Draw(dst, fmt.Sprintf("梅花 上%d下%d变%d", m.GuaUpIdx, m.GuaDownIdx, m.YaoChangeIdx), ft14, m.X+32, m.Y+16, colorWhite)
	cx, cy := m.X+32, m.Y+32
	text.Draw(dst, "本卦", ft14, cx, cy, colorWhite)
	text.Draw(dst, "互卦", ft14, cx+64, cy, colorWhite)
	text.Draw(dst, "变卦", ft14, cx+128, cy, colorWhite)
	cy += 16
	text.Draw(dst, m.GuaOrigin, ft14, cx, cy, colorWhite)
	text.Draw(dst, m.GuaProcess, ft14, cx+64, cy, colorWhite)
	text.Draw(dst, m.GuaChange, ft14, cx+128, cy, colorWhite)
	cx += 24
	cy += 32
	dis := int(math.Round(float64(meiHuaGuaSize) * 1.25))
	text.Draw(dst, fmt.Sprintf("%s%s", m.GuaUp, qimen.DiagramsWuxing[m.GuaUp]), ft14, cx, cy, color5Xing[qimen.DiagramsWuxing[m.GuaUp]])
	text.Draw(dst, fmt.Sprintf("%s%s", m.GuaDown, qimen.DiagramsWuxing[m.GuaDown]), ft14, cx, cy+dis, color5Xing[qimen.DiagramsWuxing[m.GuaDown]])
	text.Draw(dst, fmt.Sprintf("%s%s", m.GuaUpProcess, qimen.DiagramsWuxing[m.GuaUpProcess]), ft14, cx+64, cy, color5Xing[qimen.DiagramsWuxing[m.GuaUpProcess]])
	text.Draw(dst, fmt.Sprintf("%s%s", m.GuaDownProcess, qimen.DiagramsWuxing[m.GuaDownProcess]), ft14, cx+64, cy+dis, color5Xing[qimen.DiagramsWuxing[m.GuaDownProcess]])
	if m.YaoChangeIdx > 3 {
		text.Draw(dst, "用", ft14, cx-40, cy, colorWhite)
		text.Draw(dst, "体", ft14, cx-40, cy+dis, colorWhite)
		text.Draw(dst, fmt.Sprintf("%s%s", m.GuaUpChange, qimen.DiagramsWuxing[m.GuaUpChange]), ft14, cx+128, cy, color5Xing[qimen.DiagramsWuxing[m.GuaUpChange]])
	} else {
		text.Draw(dst, "体", ft14, cx-40, cy, colorWhite)
		text.Draw(dst, "用", ft14, cx-40, cy+dis, colorWhite)
		text.Draw(dst, fmt.Sprintf("%s%s", m.GuaDownChange, qimen.DiagramsWuxing[m.GuaDownChange]), ft14, cx+128, cy+dis, color5Xing[qimen.DiagramsWuxing[m.GuaDownChange]])
	}
	for _, sprite := range m.GuaSprite {
		sprite.Draw(dst)
	}

}

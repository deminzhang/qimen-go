package xuan

type MeHua struct {
	GuaUpIdx     uint8 //上卦序号
	GuaDownIdx   uint8 //下卦序号
	ChangeYaoIdx uint8 //变爻

	GuaOrigin string //本卦
	GuaUp     string //上卦
	GuaDown   string //下卦

	GuaProcess     string //互卦
	GuaUpProcess   string //互卦上卦
	GuaDownProcess string //互卦下卦

	GuaChange     string //变卦
	GuaUpChange   string //变卦上卦
	GuaDownChange string //变卦下卦
}

func (m *MeHua) Reset(upIdx, downIdx, change uint) {
	upIdx = (upIdx-1+8)%8 + 1
	downIdx = (downIdx-1+8)%8 + 1
	change = (change-1+6)%6 + 1
	m.GuaUpIdx = uint8(upIdx)
	m.GuaDownIdx = uint8(downIdx)
	m.ChangeYaoIdx = uint8(change)

	up := Diagrams8Origin[uint8(upIdx)]
	down := Diagrams8Origin[uint8(downIdx)]
	ori := Diagrams64FullName[uint8(upIdx*10+downIdx)]
	m.GuaUp = up
	m.GuaDown = down
	m.GuaOrigin = ori

	m.huGua()
	m.bianGua()
}

// 互卦
func (m *MeHua) huGua() {
	up := Diagrams8Origin[m.GuaUpIdx]
	down := Diagrams8Origin[m.GuaDownIdx]
	upB := Diagrams8Bin[up]
	downB := Diagrams8Bin[down]
	upN := (upB&0b11)<<1 + downB>>2
	downN := (upB&0b1)<<2 + downB>>1

	huUp := Diagrams8FromBin[upN]
	huDown := Diagrams8FromBin[downN]
	pro := Diagrams64FullName[(Diagrams8IdxOrigin[huUp]*10 + Diagrams8IdxOrigin[huDown])]
	m.GuaProcess = pro
	m.GuaUpProcess = huUp
	m.GuaDownProcess = huDown
}

// 变卦
func (m *MeHua) bianGua() {
	up := m.GuaUp
	down := m.GuaDown
	change := m.ChangeYaoIdx
	upB := Diagrams8Bin[up]
	downB := Diagrams8Bin[down]
	if change > 3 {
		upB ^= 1 << (change - 3 - 1)
	} else {
		downB ^= 1 << (change - 1)
	}
	cUp := Diagrams8FromBin[upB]
	cDown := Diagrams8FromBin[downB]
	changeGua := Diagrams64FullName[(Diagrams8IdxOrigin[cUp]*10 + Diagrams8IdxOrigin[cDown])]
	m.GuaChange = changeGua
	m.GuaUpChange = cUp
	m.GuaDownChange = cDown
}

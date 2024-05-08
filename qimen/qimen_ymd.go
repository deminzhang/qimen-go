package qimen

// TODO 年月日家
//http://www.360doc.com/content/24/0205/15/6148393_1113383057.shtml
//http://www.360doc.com/content/10/0524/00/1471833_29180548.shtml

// QMGongMonth 月家奇门宫格
type QMGongMonth struct {
	Idx int //洛书宫数
}

// QMGongDay 日家奇门宫格
type QMGongDay struct {
	Idx int //洛书宫数

	Door      string //八门
	Star      string //九星
	God12     string //十二黄黑道
	JoyfulGod string //喜神方位
}

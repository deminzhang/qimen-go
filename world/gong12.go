package world

// Gong12 十二宫 地支 黄黑道 大六壬等用
type Gong12 struct {
	Idx int //宫数子起1 1-12
	//天门
	JiangZhi string //将支盘 子丑寅卯...
	Jiang    string //将星名 登明从魁...
	IsJiang  bool   //是否将星
	//地户
	JianZhi string //建星支盘 子丑寅卯...
	Jian    string //建星名 建除满平...
	IsJian  bool   //是否建星

	IsHorse bool //是否驿马
}

package qimen

import (
	"fmt"
	"github.com/6tail/lunar-go/LunarUtil"
)

// 大六壬

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

	IsHorse    bool //是否驿马
	IsSkyHorse bool //是否天马
}

// CalBig6 大六壬 月将落时支 顺布余支 天三门兮地四户
func CalBig6(yueJian, yueJiang string, shiZhiIdx int, horse string) []Gong12 {
	var yueJiangIdx, yueJianIdx int
	for i := 1; i <= 12; i++ {
		if yueJian == LunarUtil.ZHI[i] {
			yueJianIdx = i
			break
		}
	}
	for i := 1; i <= 12; i++ {
		if yueJiang == LunarUtil.ZHI[i] {
			yueJiangIdx = i
			break
		}
	}
	var p = make([]Gong12, 12)
	for i := shiZhiIdx; i < shiZhiIdx+12; i++ {
		js := LunarUtil.ZHI[Idx12[yueJiangIdx]]
		g := fmt.Sprintf("%s", YueJiangName[js])
		z := LunarUtil.ZHI[Idx12[i]]
		bs := BuildStar(1 + i - shiZhiIdx)

		g12 := &p[Idx12[i]-1]
		g12.Idx = Idx12[i]
		g12.JiangZhi = js
		g12.Jiang = g
		g12.IsJiang = i == shiZhiIdx
		g12.JianZhi = LunarUtil.ZHI[Idx12[yueJianIdx+i-shiZhiIdx]]
		g12.Jian = bs
		g12.IsJian = bs == "建"
		g12.IsHorse = z == horse
		g12.IsSkyHorse = g == "太冲"

		yueJiangIdx++
	}
	return p
}

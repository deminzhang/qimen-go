package qimen

import (
	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/6tail/lunar-go/calendar"
	"slices"
)

// 大六壬

// Big6Gong 十二宫 地支 黄黑道 大六壬等用
type Big6Gong struct {
	Idx int //宫数子起1 1-12
	//月将盘
	JiangGan string //将干 甲乙丙丁...空亡
	JiangZhi string //将支 子丑寅卯...
	Jiang    string //将星名 登明从魁...
	IsJiang  bool   //是否将星
	//月建盘
	JianZhi string //建星支 子丑寅卯...
	Jian    string //建星名 建除满平...
	IsJian  bool   //是否建星
}

type Big6Ren struct {
	Gong  [12]Big6Gong
	Ke    [8]string
	Chuan [3]string
}

// NewBig6 大六壬 月将落时支 顺布余支 天三门兮地四户
func NewBig6(l *calendar.Lunar) *Big6Ren {
	var yueJian, yueJiang string
	jieQi := l.GetPrevJieQi()
	if jieQi.IsJie() {
		yueJian = Jie2YueJian(jieQi.GetName())
		yueJiang = Qi2YueJiang(l.GetPrevQi().GetName())
	} else { //qi
		yueJian = Jie2YueJian(l.GetPrevJie().GetName())
		yueJiang = Qi2YueJiang(jieQi.GetName())
	}
	dayGanIdx := l.GetDayGanIndex() + 1
	dayZhiIdx := l.GetDayZhiIndex() + 1
	shiZhiIdx := l.GetTimeZhiIndex() + 1
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
	var ganGongStart int
	pan := Big6Ren{}
	p := &pan.Gong
	//时支起月将
	for i := shiZhiIdx; i < shiZhiIdx+12; i++ {
		js := LunarUtil.ZHI[Idx12[yueJiangIdx]]
		name := YueJiangName[js]
		bs := BuildStar(1 + i - shiZhiIdx)
		g := &p[Idx12[i]-1]
		g.Idx = Idx12[i]
		g.JiangZhi = js
		g.Jiang = name
		g.IsJiang = i == shiZhiIdx
		g.JianZhi = LunarUtil.ZHI[Idx12[yueJianIdx+i-shiZhiIdx]]
		g.Jian = bs
		g.IsJian = bs == "建"
		if js == LunarUtil.ZHI[dayZhiIdx] {
			ganGongStart = g.Idx
		}
		yueJiangIdx++
	}
	//将盘日支起日干,日荀空亡跳过
	ganIdx := dayGanIdx
	for i := ganGongStart; i < ganGongStart+12; i++ {
		g12 := &p[Idx12[i]-1]
		if slices.Contains(KongWang[l.GetDayXun()], g12.JiangZhi) {
			g12.JiangGan = "〇"
		} else {
			g12.JiangGan = LunarUtil.GAN[Idx10[ganIdx]]
			ganIdx++
		}
	}
	//起四课 TODO
	pan.Ke[0] = LunarUtil.ZHI[dayGanIdx]

	pan.Ke[2] = pan.Ke[1]

	pan.Ke[4] = LunarUtil.ZHI[dayZhiIdx]
	
	pan.Ke[6] = pan.Ke[5]

	//定三传 TODO

	return &pan
}

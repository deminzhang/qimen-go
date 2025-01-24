package qimen

import (
	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/6tail/lunar-go/calendar"
	"slices"
)

// 大六壬
// 12*60=720

var Big6RenGanHide = map[string]string{ // 大六壬 天干寄宫
	"甲": "寅", "乙": "辰", "丙": "巳", "丁": "未", "戊": "巳", "己": "未", "庚": "申", "辛": "戌", "壬": "亥", "癸": "丑",
}

// Big6RenGong 十二宫 地支 黄黑道 大六壬等用
type Big6RenGong struct {
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
	DayGan string
	DayZhi string
	Gong   [12]Big6RenGong
	Ke     [4][2]string
	Chuan  [3]string
	Ge     string
}

// NewBig6Ren 大六壬 月将落时支 顺布余支 天三门兮地四户
func NewBig6Ren(l *calendar.Lunar) *Big6Ren {
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
	pan := Big6Ren{
		DayGan: l.GetDayGanExact(),
		DayZhi: l.GetDayZhiExact(),
	}
	gs := &pan.Gong
	//时支起月将
	for i := shiZhiIdx; i < shiZhiIdx+12; i++ {
		js := LunarUtil.ZHI[Idx12[yueJiangIdx]]
		name := YueJiangName[js]
		bs := BuildStar(1 + i - shiZhiIdx)
		g := &gs[Idx12[i]-1]
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
		g12 := &gs[Idx12[i]-1]
		if slices.Contains(KongWang[l.GetDayXun()], g12.JiangZhi) {
			g12.JiangGan = "〇"
		} else {
			g12.JiangGan = LunarUtil.GAN[Idx10[ganIdx]]
			ganIdx++
		}
	}
	//起四课
	//1- 日干上的天盘地支
	k1d := LunarUtil.GAN[dayGanIdx] //日干
	k1h := Big6RenGanHide[k1d]
	k1u := gs[ZhiIdx[k1h]-1].JiangZhi
	pan.Ke[0] = [2]string{k1d, k1u}
	//2- 日干所在位置的天盘地支
	pan.Ke[1][0] = pan.Ke[0][1]
	pan.Ke[1][1] = gs[ZhiIdx[pan.Ke[1][0]]-1].JiangZhi
	//3- 日支上的天盘地支
	pan.Ke[2][0] = LunarUtil.ZHI[dayZhiIdx] //日支
	pan.Ke[2][1] = gs[ZhiIdx[pan.Ke[2][0]]-1].JiangZhi
	//4- 日支所在位置的天盘地支
	pan.Ke[3][0] = pan.Ke[2][1]
	pan.Ke[3][1] = gs[ZhiIdx[pan.Ke[3][0]]-1].JiangZhi

	//定三传
	pan.Chuan, pan.Ge = pan.Get3Chuan(pan.Ke, pan.DayGan, pan.DayZhi)

	return &pan
}

func (p *Big6Ren) GetChuanN(chuChuan string) [3]string {
	//初传
	var chuan [3]string
	chuan[0] = chuChuan
	//中传
	for i := 0; i < 12; i++ {
		if LunarUtil.ZHI[i+1] == chuan[0] {
			chuan[1] = p.Gong[i].JiangZhi
			break
		}
	}
	//末传
	for i := 0; i < 12; i++ {
		if LunarUtil.ZHI[i+1] == chuan[1] {
			chuan[2] = p.Gong[i].JiangZhi
			break
		}
	}
	return chuan
}

func (p *Big6Ren) Get3Chuan(ke4 [4][2]string, dayGan, dayZhi string) ([3]string, string) {
	gs := &p.Gong
	overlap := gs[0].JiangZhi == LunarUtil.ZHI[1]
	reverse := gs[0].JiangZhi == LunarUtil.ZHI[7]
	yangDay := GanZhiYinYang[dayGan] == "阳"
	var xiaKe [][2]string
	var shangKe [][2]string
	var keMap = make(map[string]bool)
	xiaKeShang := make(map[string]bool)
	shangKeXia := make(map[string]bool)
	for _, ke := range ke4 {
		down, up := ke[0], ke[1]
		if WuxingKe[GanZhiWuXing[down]] == GanZhiWuXing[up] {
			if _, ok := xiaKeShang[up]; ok {
				continue
			}
			xiaKe = append(xiaKe, ke)
			xiaKeShang[up] = true
		} else if WuxingKe[GanZhiWuXing[up]] == GanZhiWuXing[down] {
			if _, ok := shangKeXia[up]; ok {
				continue
			}
			shangKe = append(shangKe, ke)
			shangKeXia[up] = true
		}
		keMap[up] = true
	}
	keRealCnt := len(keMap)
	hasKe := len(xiaKeShang)+len(shangKeXia) > 0

	var chuan [3]string
	//1.贼克法
	//取课先从下贼呼，若无下贼上克初。
	//初传之上名中次，中上加临是末居。
	//三传既定天盘将，此是入式法第一。
	//下贼上：如果四课中有一课是下克上（即下贼上），则以受克之神为初传。例如，第一课甲木克戌土，下贼上，受克之神是“戌”，则以“戌”为初传。
	//上贼下：如果四课中没有下贼上的情况，只有上克下，则以克者为初传。例如，第二课午火克申金，上克下，以“午”为初传。
	//戌寅午戌
	//寅午戌甲
	//重审?
	switch len(xiaKeShang) {
	case 1: //重审课
		if len(shangKeXia) == 0 {
			return p.GetChuanN(xiaKe[0][1]), "重审,始入"
		} else {
			return p.GetChuanN(xiaKe[0][1]), "重审"
		}
	case 0:
		if len(shangKeXia) == 1 {
			return p.GetChuanN(shangKe[0][1]), "重审,始入"
		}
	}
	//2.比用法
	//下贼或二三四侵，若逢上克亦同云。
	//常将天日比神用，阳日用阳阴用阴。
	//若或俱比俱不比，立法别有涉害陈。
	//如果四课中有两课或两课以上的下贼上或上克下，且克者与日干的阴阳属性相同（即比），则以与日干相比者为初传。
	//例如，日干为阳，有两课下贼上，其中一课的克者为阳，则取该阳克者为初传。
	//比用.下克上
	var xiaKeBi [][2]string
	for _, ke := range xiaKe {
		if GanZhiYinYang[ke[1]] == GanZhiYinYang[dayGan] {
			xiaKeBi = append(xiaKeBi, ke)
		}
	}
	if len(xiaKeBi) == 1 {
		return p.GetChuanN(xiaKeBi[0][1]), "知一"
	}
	//比用.上克下
	var shangKeBi [][2]string
	for _, ke := range shangKe {
		if GanZhiYinYang[ke[1]] == GanZhiYinYang[dayGan] {
			shangKeBi = append(shangKeBi, ke)
		}
	}
	if len(shangKeBi) == 1 {
		return p.GetChuanN(shangKeBi[0][1]), "知一"
	}
	//3.涉害法
	//涉害行来本家止，路逢多克为用取。
	//孟深仲浅季当休，复等柔辰刚日宜。
	//如果四课中有两课或两课以上的下贼上或上克下，且克者与日干的阴阳属性不同（即不比），或者克者与日干的阴阳属性相同但有多个克者，
	//需要比较克者所克的地盘之神的多少来确定初传。具体步骤如下：
	//对于上克下的情况，以上者查所克地盘之神。 对于下克上的情况，以上者查受克于地盘之神。 俱上者归地盘本家止。 如果涉害深浅相等，则取在地盘四孟上者为用；
	//如果无四孟，则取四仲上者为用；如果孟仲又复相等，阳日取第一课和第二课中先见者为用，阴日则取第三课和第四课先见者为用
	if hasKe {
		if len(shangKeXia) > 0 {

		}
		if len(xiaKeShang) > 0 {

		}
		return p.GetChuanN("X"), "涉害"
	}
	//8.伏吟法
	//伏吟有克亦会用，无克刚干柔取辰，初传所刑为中传，中传所刑末传居。若有自刑发使用，次传错乱日辰并；次传更复自刑者，冲取末传不管刑。
	//初传为自刑的伏吟课为杜传格。刚日伏吟课无克为自任格。柔日伏吟课无克为自信格。
	if overlap {
		if hasKe { //四课上下有克，照常取克发用，
			if XingZhi[chuan[0]] == chuan[0] { //如果初传是自刑的支（即初传为辰、午、酉、亥），则中传取支上神，末传取中传所刑的支。
				chuan[1] = ke4[2][1]
			} else {
				//贼克?
			}
			if XingZhi[chuan[1]] == chuan[1] { //如果中传又是自刑的支（即中传为辰、午、酉、亥），则取与中传相冲的支为末传。
				chuan[2] = ChongZhi[chuan[1]]
			} else {
				chuan[2] = XingZhi[chuan[1]]
			}
			return chuan, "伏吟"
		} else { //如果四课上下没有克,
			if yangDay { //阳日:取日上神发用，中末递刑取之（即初传刑者为中传，中传刑者为末传）
				chuan[0] = ke4[0][1]
				if XingZhi[chuan[0]] == chuan[0] { //如果初传是自刑的支，则取支上神为中传，中传刑的支为末传。
					chuan[1] = ke4[2][1]
					if XingZhi[chuan[1]] == chuan[1] { //如果中传又是自刑的支，则取与中传相冲的支为末传。
						chuan[2] = ChongZhi[chuan[1]]
					} else {
						chuan[2] = XingZhi[chuan[1]]
					}
					return chuan, "伏吟,自任,杜传"
				} else {
					chuan[1] = chuan[0]
					chuan[2] = XingZhi[chuan[1]]
					return chuan, "伏吟,自任"
				}
			} else { //阴日:取支上神为用，中末递刑取之（即初传刑者为中传，中传刑者为末传，如果中传是互刑，末传取冲）。
				chuan[0] = ke4[2][1]
				if XingZhi[chuan[0]] == chuan[0] { //如果初传是自刑的支，则取支上神为中传，中传刑的支为末传。
					chuan[1] = ke4[2][1]
					if XingZhi[chuan[1]] == chuan[1] { //如果中传又是自刑的支，则取与中传相冲的支为末传。
						chuan[2] = ChongZhi[chuan[1]]
					} else {
						chuan[2] = XingZhi[chuan[1]]
					}
					return chuan, "伏吟,自信,杜传"
				} else {
					chuan[1] = chuan[0]
					chuan[2] = XingZhi[chuan[1]]
					return chuan, "伏吟,自信"
				}
			}
		}
	}
	if keRealCnt == 4 {
		//4.遥克法
		//四课无克号为遥，日与神兮递互招。先取神遥克其日，如无方取日来遥。或有日克乎两神，复有两神来克日，择与日干比者用，阳日用阳阴用阴。
		//如果四课中既无上克下，也无下克上，则看四课上神有无克日干者，如有，则克日干者为初传；如果有两个上神均克日干，则取与日干相比者为用。
		//无上神克日，则看有无上神被日干所克，若有，则取被日干所克的上神为用，但如果有两个上神被日干克，则取与日干相比者为用。
		//两个以上克日或日克都比和,先取近者为用
		var keDayGan [][2]string //克日干者
		for _, ke := range ke4[1:] {
			if WuxingKe[GanZhiWuXing[ke[1]]] == GanZhiWuXing[dayGan] {
				keDayGan = append(keDayGan, ke)
			}
		}
		switch len(keDayGan) {
		case 1:
			return p.GetChuanN(keDayGan[0][1]), "遥克"
		case 0:
			var dayGanKe [][2]string //日干克者
			for _, ke := range ke4[1:] {
				if GanZhiWuXing[ke[1]] == WuxingKe[GanZhiWuXing[dayGan]] {
					dayGanKe = append(dayGanKe, ke)
				}
			}
			switch len(dayGanKe) {
			case 1:
				return p.GetChuanN(dayGanKe[0][1]), "遥克"
			case 0:
			default: //日干克者比
				for _, ke := range dayGanKe {
					if GanZhiYinYang[ke[1]] == GanZhiYinYang[dayGan] {
						return p.GetChuanN(ke[1]), "遥克"
					}
				}
			}
		default: //克日干者比
			for _, ke := range keDayGan { //比
				if GanZhiYinYang[ke[1]] == GanZhiYinYang[dayGan] {
					return p.GetChuanN(ke[1]), "遥克"
				}
			}
		}
		//5.昴星法 无遥无克时，阳日取酉宫上神为初传，中传取支上神，末传取干上神；阴日取从魁（天盘酉下）为初传，中传取干上神，末传取支上神。
		if yangDay { //虎视格
			chuan[0] = p.Gong[ZhiIdx["酉"]-1].JiangZhi
			chuan[1] = gs[ZhiIdx[dayZhi]-1].JiangZhi
			k1h := Big6RenGanHide[dayGan]
			chuan[2] = gs[ZhiIdx[k1h]-1].JiangZhi
			return chuan, "昴星,虎视"
		} else { //冬蛇掩目格
			for i := 0; i < 12; i++ {
				if p.Gong[i].JiangZhi == "酉" { //.Jiang==从魁
					chuan[0] = LunarUtil.ZHI[i+1]
					break
				}
			}
			k1h := Big6RenGanHide[dayGan]
			chuan[1] = gs[ZhiIdx[k1h]-1].JiangZhi
			chuan[2] = gs[ZhiIdx[dayZhi]-1].JiangZhi
			return chuan, "昴星,冬蛇掩目"
		}
	}
	//6.别责法
	//四课不全三课备，无遥无克别责视。刚日干合上头神，柔日支前三合取。皆以天上作初传，阴阳中末干中寄。
	//如果日干为阳干，那么取日干所合（天干五合）之神的上神为初传，中传和末传都用干上神。
	//如果日干为阴干，那么取日支三合局（地支三合）的前一位为初传，中传和末传都用干上神。
	if keRealCnt == 3 {
		if yangDay {
			hegan := HeGan[dayGan]
			k1h := Big6RenGanHide[hegan]
			chuan[0] = gs[ZhiIdx[k1h]-1].JiangZhi
		} else {
			he3Tomb := He3ZhiTomb[dayZhi]
			chuan[0] = he3Tomb
		}
		chuan[1] = ke4[0][1]
		chuan[2] = ke4[0][1]
		return chuan, "别责"
	}
	//7.八专法 两课无克号八专，阳日顺行三位取初传，阴日逆行三位取初传，中末总向日上眠。
	if keRealCnt == 2 {
		if yangDay { //阳日：日干上神在天盘顺数三位为初传，中传末传为干上神。
			k1h := Big6RenGanHide[dayGan]
			zhiIdx := ZhiIdx[k1h]
			zhiIdx += 2
			if zhiIdx > 12 {
				zhiIdx -= 12
			}
			chuan[0] = gs[zhiIdx-1].JiangZhi
		} else { // 阴日：第四课的上神在天盘逆数三位为初传，中传末传为干上神。
			zhiIdx := ZhiIdx[ke4[3][1]]
			zhiIdx -= 2
			if zhiIdx < 1 {
				zhiIdx += 12
			}
			chuan[0] = gs[zhiIdx-1].JiangZhi
		}
		chuan[1] = ke4[0][1]
		chuan[2] = ke4[0][1]
		return chuan, "八专"
	}
	//9.返吟法 返吟有克堪为用，初上中末先后排；无克驿马发用奇，辰中干和日末是其真。若知六日该无克，丑未同干丁己辛。丑日登明未太乙。
	if reverse {
		if false { //贼克
		} else { //以日支的驿马为初传 、日支上神为中传，日干上神为末传。
			chuan[0] = Horse[dayZhi]
			chuan[1] = ke4[2][1]
			chuan[2] = ke4[0][1]
			return chuan, "返吟"
		}
	}
	return chuan, "未知"
}

package qimen

import (
	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/6tail/lunar-go/calendar"
	"slices"
	"strings"
)

// 大六壬
// 12*60=720

var (
	Big6RenGanHide = map[string]string{ // 大六壬 天干寄宫
		"甲": "寅", "乙": "辰", "丙": "巳", "丁": "未", "戊": "巳", "己": "未", "庚": "申", "辛": "戌", "壬": "亥", "癸": "丑",
	}
	Big6RenGongHide = map[string][]string{ // 大六壬 五行寄宫 涉害法用
		//"亥": {"亥", "壬"}, "子": {"子"}, "丑": {"丑", "癸"},
		//"寅": {"寅", "甲"}, "卯": {"卯"}, "辰": {"辰", "乙"},
		//"巳": {"巳", "丙", "戊"}, "午": {"午"}, "未": {"未", "丁", "己"},
		//"申": {"申", "庚"}, "酉": {"酉"}, "戌": {"戌", "辛"},
		"亥": {"水", "水"}, "子": {"水"}, "丑": {"土", "水"},
		"寅": {"木", "木"}, "卯": {"木"}, "辰": {"土", "木"},
		"巳": {"火", "火", "土"}, "午": {"火"}, "未": {"土", "火", "土"},
		"申": {"金", "金"}, "酉": {"金"}, "戌": {"土", "金"},
	}
	//TianJiang12 大六壬十二天将
	//贵蛇朱六勾青，空白常玄阴后。
	//从戌至已逆行，以辰到亥顺就。
	TianJiang12 = []string{
		"贵人", "腾蛇", "朱雀", "六合", "勾陈", "青龙", "天空", "白虎", "太常", "玄武", "太阴", "天后",
	}
	TianJiang12Short = []string{
		"贵", "蛇", "朱", "合", "勾", "龙", "空", "虎", "常", "玄", "阴", "后",
	}
	//以日干查
	//甲戊庚牛羊，乙己鼠猴乡，丙丁猪鸡位，壬癸蛇兔藏，六辛逢马虎，此是贵人方。
	//guiRenDayStart 昼贵取前
	guiRenDayStart = map[string]string{
		"甲": "丑", "戊": "丑", "庚": "丑",
		"乙": "子", "己": "子",
		"丙": "亥", "丁": "亥",
		"壬": "巳", "癸": "巳",
		"辛": "午",
	}
	//guiRenNightStart 夜贵取后
	guiRenNightStart = map[string]string{
		"甲": "未", "戊": "未", "庚": "未",
		"乙": "申", "己": "申",
		"丙": "酉", "丁": "酉",
		"壬": "卯", "癸": "卯",
		"辛": "寅",
	}
)

type (
	// Big6RenGong 十二宫 地支 黄黑道 大六壬等用
	Big6RenGong struct {
		Idx int //宫数 子起1 1-12 地盘
		//月将,天盘
		JiangGan  string //将干 甲乙丙丁...空亡
		JiangZhi  string //将支 子丑寅卯...
		JiangName string //将星名 登明从魁...
		IsJiang   bool   //是否当值月将
		Jiang12   string //天盘贵人十二  天将
		//月建盘
		JianZhi string //建星支 子丑寅卯...
		Jian    string //建星名 建除满平...
		IsJian  bool   //是否月建
	}
	Big6Ren struct {
		MonthBuild, MonthLeader string //月建,月将
		DayGan, DayZhi          string
		DayXun                  string
		TimeZhi                 string

		Gong  [12]Big6RenGong
		Ke4   [4]Big6Ke //四课
		Chuan [3]string //三传
		KeTi  string    //课体
		//KeYi           string    //课义
		GuiRenStartType string //贵人起始类型 "卯酉"/实际日出日落
	}
	Big6Ke struct {
		Down string
		Up   string
		God  string
	}
)

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

	p := Big6Ren{
		MonthBuild:      yueJian,
		MonthLeader:     yueJiang,
		DayGan:          l.GetDayGanExact(),
		DayZhi:          l.GetDayZhiExact(),
		DayXun:          l.GetDayXunExact(),
		GuiRenStartType: "卯酉",
	}
	shiZhiIdx := ZhiIdx[l.GetTimeZhi()]
	p.Reset(shiZhiIdx)
	return &p
}

func (p *Big6Ren) Reset(shiZhiIdx int) {
	jiangIdx := ZhiIdx[p.MonthLeader]
	jianIdx := ZhiIdx[p.MonthBuild]
	dayGanIdx := GanIdx[p.DayGan]
	dayZhiIdx := ZhiIdx[p.DayZhi]
	var ganGongStart int
	gs := &p.Gong
	//时支起月将
	for i := shiZhiIdx; i < shiZhiIdx+12; i++ {
		js := LunarUtil.ZHI[Idx12[jiangIdx]]
		name := YueJiangName[js]
		bs := BuildStar(1 + i - shiZhiIdx)
		g := &gs[Idx12[i]-1]
		g.Idx = Idx12[i]
		g.JiangZhi = js
		g.JiangName = name
		g.IsJiang = i == shiZhiIdx
		g.JianZhi = LunarUtil.ZHI[Idx12[jianIdx+i-shiZhiIdx]]
		g.Jian = bs
		g.IsJian = bs == "建"
		if js == LunarUtil.ZHI[dayZhiIdx] {
			ganGongStart = g.Idx
		}
		jiangIdx++
	}
	//寄干,将盘日支起日干,日旬空亡跳过
	ganIdx := dayGanIdx
	for i := ganGongStart; i < ganGongStart+12; i++ {
		g12 := &gs[Idx12[i]-1]
		if slices.Contains(KongWang[p.DayXun], g12.JiangZhi) {
			g12.JiangGan = "〇"
		} else {
			g12.JiangGan = LunarUtil.GAN[Idx10[ganIdx]]
			ganIdx++
		}
	}
	//起贵人,布天将
	p.calcGuiRen(p.DayGan, LunarUtil.ZHI[shiZhiIdx])
	//起四课
	p.calcKe4(dayGanIdx, dayZhiIdx)
	//定三传
	var kes []string
	p.Chuan, kes = p.calcChuan()
	p.parseGe(kes)
}

// 起贵人,布天将
func (p *Big6Ren) calcGuiRen(dayGan, timeZhi string) {
	gs := &p.Gong
	//日贵,夜贵
	//卯、辰、巳、午、未、申六个时辰为昼时，酉、戌、亥、子、丑、寅六个时辰为夜时‌
	//另实际月令以日出为昼,日落为夜也可
	var guiRenPos string
	if p.GuiRenStartType == "卯酉" {
		switch timeZhi {
		case "卯", "辰", "巳", "午", "未", "申":
			guiRenPos = guiRenDayStart[dayGan]
		case "酉", "戌", "亥", "子", "丑", "寅":
			guiRenPos = guiRenNightStart[dayGan]
		}
	} else { // TODO 用实际日出日落 需月令,纬度
		//latitude := 39.9 // 例北京纬度
		//s := l.GetSolar()
		//sunrise, sunset := calculateSunriseSunset(s.GetYear(), s.GetMonth(), s.GetDay(), s.GetHour(), latitude)
	}

	//‌确定贵人方位‌：根据日干来确定贵人的方位。例如，甲、戊、庚日的贵人在丑（牛）或未（羊）；乙、己日的贵人在子（鼠）或申（猴）等‌
	//‌确定贵人类型‌：根据占课时间确定是昼贵还是夜贵。卯、辰、巳、午、未、申六个时辰为昼时，酉、戌、亥、子、丑、寅六个时辰为夜时‌
	for i, gg := range gs {
		if gg.JiangZhi == guiRenPos {
			//‌排布天将：
			//贵人落在地盘亥、子、丑、寅、卯、辰六个地支的，顺行环布；
			//‌落在巳、午、未、申、酉、戌六个地支的，逆行环布‌
			forward := gg.Idx <= 5 || gg.Idx == 12
			for j := 0; j < 12; j++ {
				gIdx := (i + j) % 12
				g := &gs[gIdx]
				if forward {
					g.Jiang12 = TianJiang12Short[j]
				} else {
					g.Jiang12 = TianJiang12Short[(12-j)%12]
				}
			}
			break
		}
	}
}

// 起四课
func (p *Big6Ren) calcKe4(dayGanIdx, dayZhiIdx int) {
	gs := &p.Gong
	//1- 日干上的天盘地支
	k1d := LunarUtil.GAN[dayGanIdx] //日干
	k1h := Big6RenGanHide[k1d]
	g1 := gs[ZhiIdx[k1h]-1]
	p.Ke4[0] = Big6Ke{Down: k1d, Up: g1.JiangZhi, God: g1.Jiang12}
	//2- 日干所在位置的天盘地支
	g2 := gs[ZhiIdx[p.Ke4[0].Up]-1]
	p.Ke4[1] = Big6Ke{Down: p.Ke4[0].Up, Up: g2.JiangZhi, God: g2.Jiang12}
	//3- 日支上的天盘地支
	g3 := gs[dayZhiIdx-1]
	p.Ke4[2] = Big6Ke{Down: LunarUtil.ZHI[dayZhiIdx], Up: g3.JiangZhi, God: g3.Jiang12}
	//4- 日支所在位置的天盘地支
	g4 := gs[ZhiIdx[p.Ke4[2].Up]-1]
	p.Ke4[3] = Big6Ke{Down: p.Ke4[2].Up, Up: g4.JiangZhi, God: g4.Jiang12}
}

// 普通三传,非伏吟
func (p *Big6Ren) chuanNormal(chuan0 string) [3]string {
	//初传
	var chuan [3]string
	chuan[0] = chuan0
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

// 8.伏吟法
// 伏吟有克亦会用，无克刚干柔取辰，初传所刑为中传，中传所刑末传居。若有自刑发使用，次传错乱日辰并；次传更复自刑者，冲取末传不管刑。
func (p *Big6Ren) chuanOverlap(hasKe bool, chuan0 string) (chuan [3]string, kts []string) {
	ke4 := p.Ke4
	dayGan := p.DayGan
	yangDay := GanZhiYinYang[dayGan] == "阳" //阳日
	//初传为自刑的伏吟课为杜传格。刚日伏吟课无克为自任格。柔日伏吟课无克为自信格。
	if hasKe { //四课上下有克，照常取克发用，
		chuan[0] = chuan0
		if XingZhi[chuan[0]] == chuan[0] { //如果初传是自刑的支（即初传为辰、午、酉、亥），则中传取支上神，末传取中传所刑的支。
			chuan[1] = ke4[2].Up
		}
		if XingZhi[chuan[1]] == chuan[1] { //如果中传又是自刑的支（即中传为辰、午、酉、亥），则取与中传相冲的支为末传。
			chuan[2] = ChongZhi[chuan[1]]
		} else {
			chuan[2] = XingZhi[chuan[1]]
		}
		kts = append(kts, "艮", "伏吟")
		return
	} else { //如果四课上下没有克,
		if yangDay { //阳日:取日上神发用，中末递刑取之（即初传刑者为中传，中传刑者为末传）
			chuan[0] = ke4[0].Up
			chuan[1] = XingZhi[chuan[0]]
			if chuan[1] == chuan[0] { //如果初传是自刑的支，则取日支上神为中传，中传刑的支为末传。
				chuan[1] = ke4[2].Up
				chuan[2] = XingZhi[chuan[1]]
				if chuan[2] == chuan[1] { //如果中传又是自刑的支，则取与中传相冲的支为末传。
					chuan[2] = ChongZhi[chuan[1]]
				}
				kts = append(kts, "艮", "伏吟", "自任")
				return
			} else {
				chuan[2] = XingZhi[chuan[1]]
				kts = append(kts, "艮", "伏吟", "自信")
				return
			}
		} else { //阴日:取支上神为用，中末递刑取之（即初传刑者为中传，中传刑者为末传，如果中传是互刑，末传取冲）。
			chuan[0] = ke4[2].Up
			chuan[1] = XingZhi[chuan[0]]
			if chuan[1] == chuan[0] { //如果初传是自刑的支，则取日干上神为中传，中传刑的支为末传。
				chuan[1] = ke4[0].Up
				chuan[2] = XingZhi[chuan[1]]
				if chuan[2] == chuan[1] { //如果中传又是自刑的支，则取与中传相冲的支为末传。
					chuan[2] = ChongZhi[chuan[1]]
				}
				kts = append(kts, "艮", "伏吟", "自任")
				return
			} else {
				chuan[2] = XingZhi[chuan[1]]
				kts = append(kts, "艮", "伏吟", "自信")
				return
			}
		}
	}
}

func (p *Big6Ren) calcChuan() (chuan [3]string, kts []string) {
	ke4 := p.Ke4
	dayGan := p.DayGan
	dayZhi := p.DayZhi
	gs := &p.Gong
	overlap := gs[0].JiangZhi == LunarUtil.ZHI[1] //伏吟
	reverse := gs[0].JiangZhi == LunarUtil.ZHI[7] //反吟
	yangDay := GanZhiYinYang[dayGan] == "阳"       //阳日
	var xiaKe []Big6Ke
	var shangKe []Big6Ke
	var keMap = make(map[string]bool)
	xiaKeShang := make(map[string]bool) // 下贼上 map[上]=true
	shangKeXia := make(map[string]bool) // 上克下 map[上]=true
	for _, ke := range ke4 {
		down, up := ke.Down, ke.Up
		if WuXingKe[GanZhiWuXing[down]] == GanZhiWuXing[up] {
			if _, ok := xiaKeShang[up]; ok {
				continue
			}
			xiaKe = append(xiaKe, ke)
			xiaKeShang[up] = true
		} else if WuXingKe[GanZhiWuXing[up]] == GanZhiWuXing[down] {
			if _, ok := shangKeXia[up]; ok {
				continue
			}
			shangKe = append(shangKe, ke)
			shangKeXia[up] = true
		}
		keMap[up] = true
	}
	keRealCnt := len(keMap)                      //去重课数
	hasKe := len(xiaKeShang)+len(shangKeXia) > 0 //有克

	if hasKe {
		//1.贼克法
		//取课先从下贼呼，若无下贼上克初。
		//初传之上名中次，中上加临是末居。
		//三传既定天盘将，此是入式法第一。
		//上贼下：如果四课中没有下贼上的情况，只有上克下，则以克者为初传。例如，第二课午火克申金，上克下，以“午”为初传。
		//下克上：如果四课中有一课是下克上（即下贼上），则以受克之神为初传。例如，第一课甲木克戌土，下贼上，受克之神是“戌”，则以“戌”为初传。
		switch len(xiaKeShang) {
		case 1: //重审课
			chuan[0] = xiaKe[0].Up
			//chuan = p.chuanNormal(xiaKe[0].Up)
			kts = append(kts, "坤", "重审")
			if len(shangKeXia) == 0 {
				kts = append(kts, "始入")
			}
			//return
		case 0:
			if len(shangKeXia) == 1 {
				chuan[0] = shangKe[0].Up
				//chuan = p.chuanNormal(shangKe[0].Up)
				kts = append(kts, "乾", "元首")
				//return
			}
		}
		//2.比用法
		//下贼或二三四侵，若逢上克亦同云。
		//常将天日比神用，阳日用阳阴用阴。
		//若或俱比俱不比，立法别有涉害陈。
		//如果四课中有两课或两课以上的下贼上或上克下，且克者与日干的阴阳属性相同（即比），则以与日干相比者为初传。
		//例如，日干为阳，有两课下贼上，其中一课的克者为阳，则取该阳克者为初传。
		//比用.下克上
		var xiaKeBi []Big6Ke
		for _, ke := range xiaKe {
			if GanZhiYinYang[ke.Up] == GanZhiYinYang[dayGan] {
				xiaKeBi = append(xiaKeBi, ke)
			}
		}
		if len(xiaKeBi) == 1 {
			chuan[0] = xiaKeBi[0].Up
			kts = append(kts, "坎", "比用", "知一")
			//return p.chuanNormal(xiaKeBi[0].Up), kts
		}
		if len(xiaKeBi) == 0 {
			//比用.上克下
			var shangKeBi []Big6Ke
			for _, ke := range shangKe {
				if GanZhiYinYang[ke.Up] == GanZhiYinYang[dayGan] {
					shangKeBi = append(shangKeBi, ke)
				}
			}
			if len(shangKeBi) == 1 {
				chuan[0] = shangKeBi[0].Up
				kts = append(kts, "坎", "比用")
				//return p.chuanNormal(shangKeBi[0].Up), kts
			}
		}
		if !overlap {
			if chuan[0] != "" {
				return p.chuanNormal(chuan[0]), kts
			}
		}
	}
	if overlap {
		return p.chuanOverlap(hasKe, chuan[0])
	}
	if reverse {
		//9.反吟法 反吟有克堪为用，初上中末先后排；无克驿马发用奇，辰中干和日末是其真。若知六日该无克，丑未同干丁己辛。丑日登明未太乙。
		kts = append(kts, "震", "反吟")
		if !hasKe { //以日支的驿马为初传 、日支上神为中传，日干上神为末传。
			chuan[0] = Horse[dayZhi]
			chuan[1] = ke4[2].Up
			chuan[2] = ke4[0].Up
			kts = append(kts, "井栏")
			if keRealCnt == 2 {
				kts = append(kts, "八专")
			}
			return
		}
	}
	//3.涉害法 涉害行来本家止，路逢多克为用取。孟深仲浅季当休，复等柔辰刚日宜。
	{
		//如果四课中有两课或两课以上的下贼上或上克下，且克者与日干的阴阳属性不同（即不比），或者克者与日干的阴阳属性相同但有多个克者，
		//需要比较克者所克的地盘之神的多少来确定初传。具体步骤如下：
		//对于下克上的情况，以上者查受克于地盘之神。 俱上者归地盘本家止。
		//对于上克下的情况，以上者查所克地盘之神。
		//如果涉害深浅相等，则取在地盘四孟上者为用；
		//如果无四孟，则取四仲上者为用；如果孟仲又复相等，阳日取第一课和第二课中先见者为用，阴日则取第三课和第四课先见者为用
		hits := map[string]int{}
		if len(xiaKeShang) > 1 {
			for up := range xiaKeShang {
				if GanZhiYinYang[up] != GanZhiYinYang[dayGan] { //排除不比的
					continue
				}
				up5x := GanZhiWuXing[up]
				upHomeIdx := ZhiIdx[up]
				for _, g := range gs {
					if g.JiangZhi == up {
						for j := g.Idx; j < g.Idx+12; j++ {
							gIdx := Idx12[j]
							gz := LunarUtil.ZHI[gIdx]
							for _, wx := range Big6RenGongHide[gz] {
								if WuXingKe[wx] == up5x {
									hits[up]++
								}
							}
							if j == upHomeIdx {
								break
							}
						}
						break
					}
				}
			}
		}
		if len(hits) == 0 {
			if len(shangKeXia) > 1 {
				for up := range shangKeXia {
					if GanZhiYinYang[up] != GanZhiYinYang[dayGan] { //排除不比的
						continue
					}
					up5x := GanZhiWuXing[up]
					upHomeIdx := ZhiIdx[up]
					for _, g := range gs {
						if g.JiangZhi == up {
							for j := g.Idx; j < g.Idx+12; j++ {
								gIdx := Idx12[j]
								gz := LunarUtil.ZHI[gIdx]
								for _, wx := range Big6RenGongHide[gz] {
									if WuXingKe[up5x] == wx {
										hits[up]++
									}
								}
								if j == upHomeIdx {
									break
								}
							}
							break
						}
					}
				}
			}
		}
		//hits中找唯一最大值
		var maxUp string
		var maxN int
		for up, n := range hits {
			if n > maxN {
				maxN = n
				maxUp = up
			} else if n == maxN {
				maxUp = ""
			}
		}
		if maxUp != "" {
			return p.chuanNormal(maxUp), []string{"坎", "涉害"}
		} else { //4个1 || 2个2
			//如果涉害深浅相等，则取在地盘四孟上者为用；
			//如果无四孟，则取四仲上者为用；
			mid := map[string]struct{}{}
			for up := range hits { //取在地盘四孟上者为用
				if gs[2].JiangZhi == up || gs[5].JiangZhi == up || gs[8].JiangZhi == up || gs[11].JiangZhi == up {
					mid[up] = struct{}{} //见机
				}
			}
			switch len(mid) {
			case 1: //见机
				for up := range mid {
					return p.chuanNormal(up), []string{"坎", "涉害", "见机"}
				}
			case 0:
				for up := range hits { //如果无四孟，则取四仲上者为用
					if gs[0].JiangZhi == up || gs[3].JiangZhi == up || gs[6].JiangZhi == up || gs[9].JiangZhi == up {
						mid[up] = struct{}{} //察微
					}
				}
				if len(mid) == 1 {
					for up := range mid {
						return p.chuanNormal(up), []string{"坎", "涉害", "察微"}
					}
				}
			}
			if len(mid) > 1 {
				//如果孟仲又复相等，阳日取第一课和第二课中先见者为用，阴日则取第三课和第四课先见者为用?
				//戊辰日子上发用 缀瑕 复等
				if yangDay {
					for up := range mid {
						if up == ke4[0].Up || up == ke4[1].Up {
							return p.chuanNormal(up), []string{"坎", "涉害", "缀瑕"} // 复等
						}
					}
				} else {
					for up := range mid {
						if up == ke4[2].Up || up == ke4[3].Up {
							return p.chuanNormal(up), []string{"坎", "涉害", "缀瑕"} // 复等
						}
					}
				}
			}
		}
		//注：还有一种直接用孟仲法来取三传，就是不管受克深浅，直接按照如上方式去排三传，两种方式各有优缺，各位壬友请自行比较！
	}
	//7.八专法 两课无克号八专，阳日顺行三位取初传，阴日逆行三位取初传，中末总向日上眠。
	if keRealCnt == 2 {
		if yangDay { //阳日：日干上神在天盘顺数三位为初传，中传末传为干上神。
			k1h := Big6RenGanHide[dayGan]
			zhiIdx := (ZhiIdx[k1h] + 2) % 12
			if zhiIdx == 0 {
				zhiIdx = 12
			}
			chuan[0] = gs[zhiIdx-1].JiangZhi
		} else { // 阴日：第四课的上神在天盘逆数三位为初传，中传末传为干上神。
			zhiIdx := (ZhiIdx[ke4[3].Up] - 2 + 12) % 12
			if zhiIdx == 0 {
				zhiIdx = 12
			}
			chuan[0] = gs[zhiIdx-1].JiangZhi
		}
		chuan[1] = ke4[0].Up
		chuan[2] = ke4[0].Up
		kts = append(kts, "同人", "八专")
		return
	}
	//4.遥克法
	//四课无克号为遥，日与神兮递互招。先取神遥克其日，如无方取日来遥。或有日克乎两神，复有两神来克日，择与日干比者用，阳日用阳阴用阴。
	//伏吟,反吟,八专,不做遥克
	if keRealCnt == 4 {
		//如果四课中既无上克下，也无下克上，则看四课上神有无克日干者，如有，则克日干者为初传；如果有两个上神均克日干，则取与日干相比者为用。
		//无上神克日，则看有无上神被日干所克，若有，则取被日干所克的上神为用，但如果有两个上神被日干克，则取与日干相比者为用。
		//两个以上克日或日克都比和,先取近者为用
		var keDayGan []Big6Ke //克日干者
		for _, ke := range ke4[1:] {
			if WuXingKe[GanZhiWuXing[ke.Up]] == GanZhiWuXing[dayGan] {
				keDayGan = append(keDayGan, ke)
			}
		}
		switch len(keDayGan) {
		case 1:
			return p.chuanNormal(keDayGan[0].Up), []string{"睽", "遥克", "蒿矢"}
		case 0:
			var dayGanKe []Big6Ke //日干克者
			for _, ke := range ke4[1:] {
				if GanZhiWuXing[ke.Up] == WuXingKe[GanZhiWuXing[dayGan]] {
					dayGanKe = append(dayGanKe, ke)
				}
			}
			switch len(dayGanKe) {
			case 1:
				return p.chuanNormal(dayGanKe[0].Up), []string{"睽", "遥克", "弹射"}
			case 0:
			default: //日干克者比
				for _, ke := range dayGanKe {
					if GanZhiYinYang[ke.Up] == GanZhiYinYang[dayGan] {
						return p.chuanNormal(ke.Up), []string{"遥克", "弹射"}
					}
				}
			}
		default: //克日干者比
			for _, ke := range keDayGan { //比
				if GanZhiYinYang[ke.Up] == GanZhiYinYang[dayGan] {
					return p.chuanNormal(ke.Up), []string{"睽", "遥克", "蒿矢"}
				}
			}
		}
		//5.昴星法 无遥无克时，阳日取酉宫上神为初传，中传取支上神，末传取干上神；阴日取从魁（天盘酉下）为初传，中传取干上神，末传取支上神。
		if yangDay { //虎视格
			chuan[0] = p.Gong[ZhiIdx["酉"]-1].JiangZhi
			chuan[1] = gs[ZhiIdx[dayZhi]-1].JiangZhi
			k1h := Big6RenGanHide[dayGan]
			chuan[2] = gs[ZhiIdx[k1h]-1].JiangZhi
			kts = append(kts, "履", "昴星", "虎视")
			return
		} else { //冬蛇掩目格
			for i := 0; i < 12; i++ {
				if p.Gong[i].JiangZhi == "酉" { //.JiangName==从魁
					chuan[0] = LunarUtil.ZHI[i+1]
					break
				}
			}
			k1h := Big6RenGanHide[dayGan]
			chuan[1] = gs[ZhiIdx[k1h]-1].JiangZhi
			chuan[2] = gs[ZhiIdx[dayZhi]-1].JiangZhi
			kts = append(kts, "履", "昴星", "冬蛇掩目")
			return
		}
	}
	//6.别责法
	//四课不全三课备，无遥无克别责视。刚日干合上头神，柔日支前三合取。皆以天上作初传，阴阳中末干中寄。
	if keRealCnt == 3 {
		//如果日干为阳干，那么取日干所合（天干五合）之神的上神为初传，中传和末传都用干上神。
		//如果日干为阴干，那么取日支三合局（地支三合）的前一位为初传，中传和末传都用干上神。
		if yangDay {
			he := HeGan[dayGan]
			k1h := Big6RenGanHide[he]
			chuan[0] = gs[ZhiIdx[k1h]-1].JiangZhi
		} else {
			he3F := He3Zhi[dayZhi][2]
			chuan[0] = he3F
		}
		chuan[1] = ke4[0].Up
		chuan[2] = ke4[0].Up
		kts = append(kts, "涣", "别责")
		return
	}
	return
}

// 课体细析
// http://www.360doc.com/content/23/0331/21/46945463_1074566892.shtml
func (p *Big6Ren) parseGe(ts []string) {
	c := p.Chuan
	c0, c1, c2 := c[0], c[1], c[2]
	if c0 == "" || c1 == "" || c2 == "" {
		return
	}
	zi0 := ZhiIdx[c0]
	zi1 := ZhiIdx[c1]
	zi2 := ZhiIdx[c2]

	//三合会局 金：从革 木：曲直  水：润下 火：炎上
	//四季土:稼穑
	if c0 != c1 && c0 != c2 && c1 != c2 {
		if HeWuXing[c0] == HeWuXing[c1] && HeWuXing[c1] == HeWuXing[c2] { //三合局
			ts = append(ts, WuXingGe[HeWuXing[c0]])
		} else if HuiWuXing[c0] == HuiWuXing[c1] && HuiWuXing[c1] == HuiWuXing[c2] { //三会局
			ts = append(ts, WuXingGe[HeWuXing[c0]])
		}
	}
	x3 := "辰戌丑未"
	if strings.Contains(x3, c0) && strings.Contains(x3, c1) && strings.Contains(x3, c2) {
		ts = append(ts, WuXingGe["土"])
	}
	x1 := "寅巳申亥"
	if strings.Contains(x1, c0) && strings.Contains(x1, c1) && strings.Contains(x1, c2) {
		ts = append(ts, "家人", "元胎") //玄胎
		if slices.Contains(ts, "反吟") {
			ts = append(ts, "绝胎")
		}
		cs3 := c0 + c1 + c2
		switch cs3 {
		case "寅巳申", "巳申亥", "申亥寅", "亥寅巳":
			ts = append(ts, "病胎")
		case "寅亥申", "巳寅亥", "申巳寅", "亥申巳":
			ts = append(ts, "生胎")
		}
	}

	//进茹
	if Idx12[zi0+1] == zi1 && Idx12[zi1+1] == zi2 {
		ts = append(ts, "进茹")
	} else if Idx12[zi0-1+12] == zi1 && Idx12[zi1-1+12] == zi2 {
		ts = append(ts, "退茹")
	} else if Idx12[zi0+2] == zi1 && Idx12[zi1+2] == zi2 { //间传
		ts = append(ts, "顺间")
	} else if Idx12[zi0-2+12] == zi1 && Idx12[zi1-2+12] == zi2 {
		ts = append(ts, "逆间")
	}
	//斩关课的定义:辰成加日辰发用，很多时候辰成加日辰之上，辰成之阴神发用，也为斩关课。比如日辰之上见辰，辰之阴见子，子发用这也算是斩关课，但不是标准的斩关课。
	//标准的斩关课，是辰戌发用，坐下是寅、卯，以木克动土。辰成乃重土，主闭塞。寅天梁，卯天关辰天罡，戌天魁，以木克动土，三天俱动。
	if strings.Contains("辰戌", p.Ke4[0].Up) {
		if strings.Contains("寅卯甲乙亥子壬癸", p.Ke4[0].Down) {
			ts = append(ts, "斩关")
		}
	} else if strings.Contains("辰戌", p.Ke4[2].Up) {
		if strings.Contains("寅卯甲乙亥子壬癸", p.Ke4[2].Down) {
			ts = append(ts, "斩关")
		}
	}

	p.KeTi = strings.Join(ts, ",")
}

func (p *Big6Ren) GetGongByJiangZhi(zhiUp string) *Big6RenGong {
	for i := 0; i < 12; i++ {
		if p.Gong[i].JiangZhi == zhiUp {
			return &p.Gong[i]
		}
	}
	return nil
}

func (p *Big6Ren) Relation6(zhi string) string {
	return Relation6Short[RelationGanZhi(p.DayGan, zhi)]
}

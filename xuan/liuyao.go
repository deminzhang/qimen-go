package xuan

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/6tail/lunar-go/calendar"
)

// NewLiuYao 自动起卦（时间+随机数）
func NewLiuYao(now time.Time) (*LiuYaoResult, *calendar.Lunar) {
	seed := now.UnixMilli() + int64(rand.Intn(1000000))
	yaoRaw := GenerateYao(seed)
	solar := calendar.NewSolar(now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second())
	lunar := calendar.NewLunarFromSolar(solar)
	result := CalcLiuYao(lunar, yaoRaw)
	return result, lunar
}

// CalcLiuYao 给定农历和原始爻，计算完整排盘
func CalcLiuYao(lunar *calendar.Lunar, yaoRaw []string) *LiuYaoResult {
	bazi := lunar.GetEightChar()

	// 本卦爻
	baseGua := make([]YaoGanZhi, 6)
	for i, tp := range yaoRaw {
		yao := 0
		if tp == "1" || tp == "1o" {
			yao = 1
		}
		baseGua[i] = YaoGanZhi{Yao: yao, Type: tp}
	}

	// 有无动爻
	hasBian := false
	for _, tp := range yaoRaw {
		if tp == "1o" || tp == "0x" {
			hasBian = true
			break
		}
	}

	// 本卦名
	code := make([]byte, 6)
	for i, y := range baseGua {
		code[5-i] = byte('0') + byte(y.Yao)
	}
	baseName := GuaNames[string(code)]

	// 内卦外卦
	inner := reverse3(baseGua[0:3])
	outer := reverse3(baseGua[3:6])
	innerName := GuaCodeMap[inner]
	outerName := GuaCodeMap[outer]

	// 纳甲装干支
	setupGanZhi(baseGua, innerName, outerName)

	// 世应
	shiIdx := calcShi(inner, outer)
	var yingIdx int
	if shiIdx > 2 {
		yingIdx = shiIdx % 3
	} else {
		yingIdx = shiIdx%3 + 3
	}
	setShiYing(baseGua, shiIdx, yingIdx)

	// 变卦
	var bianGua []YaoGanZhi
	var bianName, bianGuaGong, bianAlias string
	var dongYao []int
	if hasBian {
		bianGua = make([]YaoGanZhi, 6)
		for i, tp := range yaoRaw {
			var yao int
			switch tp {
			case "1o":
				yao = 0
			case "0x":
				yao = 1
			case "1":
				yao = 1
			default:
				yao = 0
			}
			bianGua[i] = YaoGanZhi{Yao: yao, Type: tp}
			if tp == "1o" || tp == "0x" {
				dongYao = append(dongYao, i+1)
			}
		}
		code2 := make([]byte, 6)
		for i, y := range bianGua {
			code2[5-i] = byte('0') + byte(y.Yao)
		}
		bianName = GuaNames[string(code2)]

		bInner := reverse3(bianGua[0:3])
		bOuter := reverse3(bianGua[3:6])
		bInName := GuaCodeMap[bInner]
		bOutName := GuaCodeMap[bOuter]
		setupGanZhi(bianGua, bInName, bOutName)

		bShi2 := calcShi(bInner, bOuter)
		var bYing2 int
		if bShi2 > 2 {
			bYing2 = bShi2 % 3
		} else {
			bYing2 = bShi2%3 + 3
		}
		setShiYing(bianGua, bShi2, bYing2)

		bGg := calcGuaGong(bianGua, bInName, bOutName)
		setupLiuQin(bianGua, bGg)
		bianGuaGong = bGg

		if isGuiHun(bInName, bOutName) {
			bianAlias = "归魂"
		} else if isYouHun(bInName, bOutName) {
			bianAlias = "游魂"
		}
	}

	// 卦宫
	guaGong := calcGuaGong(baseGua, innerName, outerName)

	// 别名
	alias := ""
	if isGuiHun(innerName, outerName) {
		alias = "归魂"
	} else if isYouHun(innerName, outerName) {
		alias = "游魂"
	} else if isLiuHe(innerName, outerName) {
		alias = "六合"
	} else if isLiuChong(innerName, outerName) {
		alias = "六冲"
	}

	// 六亲
	setupLiuQin(baseGua, guaGong)
	if hasBian {
		setupLiuQin(bianGua, guaGong)
	}

	// 伏神
	setupFuShen(baseGua, guaGong)

	// 六神
	sixShen := installLiuShen(bazi.GetDayGan())

	// 神煞
	riGan := bazi.GetDayGan()
	riZhi := bazi.GetDayZhi()
	yueZhi := bazi.GetMonthZhi()
	nianZhi := bazi.GetYearZhi()
	shiZhi := bazi.GetTimeZhi()
	shenSha := calcLiuYaoShenSha(baseGua, riGan, riZhi, yueZhi, nianZhi, shiZhi)

	revDong := make([]int, len(dongYao))
	for i, v := range dongYao {
		revDong[len(dongYao)-1-i] = v
	}

	// 获取节气信息
	prevJie := lunar.GetPrevJieQi()
	nextJie := lunar.GetNextJieQi()

	result := &LiuYaoResult{
		DateTime:    fmt.Sprintf("%d年%d月%d日", lunar.GetYear(), lunar.GetMonth(), lunar.GetDay()),
		YearGan:     bazi.GetYearGan(),
		YearZhi:     bazi.GetYearZhi(),
		MonthGan:    bazi.GetMonthGan(),
		MonthZhi:    yueZhi,
		DayGan:      riGan,
		DayZhi:      riZhi,
		HourGan:     bazi.GetTimeGan(),
		HourZhi:     shiZhi,
		YaoRaw:      yaoRaw,
		BaseName:    baseName,
		BaseGua:     baseGua,
		GuaGong:     guaGong,
		Alias:       alias,
		HasBian:     hasBian,
		BianName:    bianName,
		BianGua:     bianGua,
		BianGuaGong: bianGuaGong,
		BianAlias:   bianAlias,
		DongYao:     revDong,
		SixShen:     sixShen,
		ShenSha:     shenSha,
		ShiZhi:      shiZhi,
		Hour:        lunar.GetHour(),
		Minute:      lunar.GetMinute(),
		LunarMonth:  lunar.GetMonthInChinese(),
		LunarDay:    lunar.GetDayInChinese(),
		JieQiFrom:   prevJie.GetName(),
		JieQiFromDate: prevJie.GetSolar().ToYmd(),
		JieQiTo:     nextJie.GetName(),
		JieQiToDate: nextJie.GetSolar().ToYmd(),
	}

	return result
}

// GenerateYao 生成六爻（模拟三枚铜钱）
func GenerateYao(seed int64) []string {
	r := rand.New(rand.NewSource(seed))
	result := make([]string, 6)
	for i := 0; i < 6; i++ {
		t1 := 2
		if r.Float64() < 0.5 {
			t1 = 3
		}
		t2 := 2
		if r.Float64() < 0.5 {
			t2 = 3
		}
		t3 := 2
		if r.Float64() < 0.5 {
			t3 = 3
		}
		sum := t1 + t2 + t3
		switch sum {
		case 9:
			result[i] = "1o"
		case 8:
			result[i] = "0"
		case 7:
			result[i] = "1"
		default:
			result[i] = "0x"
		}
	}
	return result
}

// reverse3 三爻编码
func reverse3(yaos []YaoGanZhi) string {
	b := make([]byte, 3)
	for i, y := range yaos {
		b[2-i] = byte('0') + byte(y.Yao)
	}
	return string(b)
}

// setupGanZhi 装干支
func setupGanZhi(yaoList []YaoGanZhi, innerName, outerName string) {
	for i := 0; i < 6; i++ {
		var gua string
		var pos int
		if i < 3 {
			gua = innerName
			pos = i
		} else {
			gua = outerName
			pos = i - 3
		}
		ganPair := GanMap[gua]
		zhiPair := ZhiMap[gua]
		if i < 3 {
			yaoList[i].Gan = ganPair[0]
			yaoList[i].Zhi = zhiPair[0][pos]
		} else {
			yaoList[i].Gan = ganPair[1]
			yaoList[i].Zhi = zhiPair[1][pos]
		}
	}
}

// calcShi 计算世爻位置 返回0-5对应世爻位置(从初爻0开始)
func calcShi(inner, outer string) int {
	if inner == outer {
		return 5
	}
	if inner[0] != outer[0] && inner[1] != outer[1] && inner[2] != outer[2] {
		return 2
	}
	if inner[0] == outer[0] && inner[1] != outer[1] && inner[2] != outer[2] {
		return 1
	}
	if inner[0] != outer[0] && inner[1] == outer[1] && inner[2] == outer[2] {
		return 4
	}
	if inner[0] != outer[0] && inner[1] != outer[1] && inner[2] == outer[2] {
		return 3
	}
	if inner[0] == outer[0] && inner[1] == outer[1] && inner[2] != outer[2] {
		return 0
	}
	if inner[0] != outer[0] && inner[1] == outer[1] && inner[2] != outer[2] {
		return 3
	}
	if inner[0] == outer[0] && inner[1] != outer[1] && inner[2] == outer[2] {
		return 2
	}
	return 0
}

// setShiYing 设置世应
func setShiYing(yaoList []YaoGanZhi, shiIdx, yingIdx int) {
	for i := range yaoList {
		yaoList[i].ShiYing = 0
	}
	yaoList[shiIdx].ShiYing = 1
	yaoList[yingIdx].ShiYing = 2
}

// calcGuaGong 计算卦宫
func calcGuaGong(yaoList []YaoGanZhi, innerName, outerName string) string {
	for _, pair := range GuiHunList {
		if pair[0] == innerName && pair[1] == outerName {
			return innerName
		}
	}
	revInner := make([]byte, 3)
	for i := 0; i < 3; i++ {
		revInner[2-i] = byte('0') + byte(1-yaoList[i].Yao)
	}
	if g, ok := GuaCodeMap[string(revInner)]; ok {
		return g
	}
	return innerName
}

// setupLiuQin 装六亲
func setupLiuQin(yaoList []YaoGanZhi, guaGong string) {
	selfWx := GuaWuXing[guaGong]
	for i := range yaoList {
		targetWx := ZhiWuXing[yaoList[i].Zhi]
		switch {
		case WuXingSheng[selfWx] == targetWx:
			yaoList[i].Qin = "孙"
		case WuXingSheng[targetWx] == selfWx:
			yaoList[i].Qin = "父"
		case WuXingKe[selfWx] == targetWx:
			yaoList[i].Qin = "财"
		case WuXingKe[targetWx] == selfWx:
			yaoList[i].Qin = "官"
		default:
			yaoList[i].Qin = "兄"
		}
	}
}

// installLiuShen 安装六神
func installLiuShen(riGan string) []string {
	start, ok := LiuShenGanMap[riGan]
	if !ok {
		start = 0
	}
	result := make([]string, 6)
	for i := 0; i < 6; i++ {
		result[i] = LiuShenOrder[(start+i)%6]
	}
	return result
}

// setupFuShen 伏神
func setupFuShen(yaoList []YaoGanZhi, guaGong string) {
	qinSet := make(map[string]bool)
	for _, y := range yaoList {
		qinSet[y.Qin] = true
	}

	code, ok := reverseMap(GuaCodeMap, guaGong)
	if !ok {
		return
	}
	code = code + code

	gongYao := make([]YaoGanZhi, 6)
	for i, c := range code {
		yao := 0
		if c == '1' {
			yao = 1
		}
		gongYao[i] = YaoGanZhi{Yao: yao}
	}

	gInnerCode := string([]byte{code[2], code[1], code[0]})
	gOuterCode := string([]byte{code[5], code[4], code[3]})
	gInName := GuaCodeMap[gInnerCode]
	gOutName := GuaCodeMap[gOuterCode]
	setupGanZhi(gongYao, gInName, gOutName)
	setupLiuQin(gongYao, guaGong)

	for _, q := range LiuQinOrder {
		if !qinSet[q] {
			for i, gy := range gongYao {
				if gy.Qin == q {
					if yaoList[i].Fu == nil {
						yaoList[i].Fu = &YaoFuShen{}
					}
					yaoList[i].Fu.Qin = q
					yaoList[i].Fu.Gan = gy.Gan
					yaoList[i].Fu.Zhi = gy.Zhi
				}
			}
		}
	}
}

// 神煞计算
func calcLiuYaoShenSha(yaoList []YaoGanZhi, riGan, riZhi, yueZhi, nianZhi, shiZhi string) []ShenShaItem {
	gs := guaShen(yaoList)

	var ss []ShenShaItem
	ss = append(ss, ShenShaItem{"卦身", gs})
	ss = append(ss, yima(riZhi))
	ss = append(ss, taohua(riZhi))
	ss = append(ss, lushen(riGan))
	ss = append(ss, ride(riGan))
	ss = append(ss, poshui(riZhi))
	ss = append(ss, wenchang(riGan))
	ss = append(ss, guiren(riGan, shiZhi))
	ss = append(ss, tianyi(yueZhi))
	ss = append(ss, diyi(yueZhi))
	ss = append(ss, tianma(yueZhi))
	ss = append(ss, tianxi(yueZhi))
	ss = append(ss, yuede(yueZhi))
	ss = append(ss, zaisha(nianZhi))
	ss = append(ss, jiesha(nianZhi))
	ss = append(ss, jiangxing(riZhi))
	ss = append(ss, huagai(riZhi))
	ss = append(ss, ShenShaItem{"香闺", xianggui(yaoList, gs)})
	ss = append(ss, ShenShaItem{"床帐", chuangzhang(yaoList, gs)})
	ss = append(ss, yangren(riGan))
	ss = append(ss, shengqi(yueZhi))
	ss = append(ss, siqi(yueZhi))
	ss = append(ss, wangshen(yueZhi))
	ss = append(ss, bingfu(nianZhi))
	ss = append(ss, shangmen(nianZhi))
	ss = append(ss, diaoke(nianZhi))
	ss = append(ss, guanfu(yueZhi))
	ss = append(ss, youdu(riGan))
	ss = append(ss, feifu(riGan))
	ss = append(ss, xuezhi(yueZhi))
	ss = append(ss, gucheng(yueZhi))
	ss = append(ss, guashu(yueZhi))
	ss = append(ss, sishen(yueZhi))
	return ss
}

// ============ 神煞函数 ============

func yima(zhi string) ShenShaItem {
	m := map[string]string{"申": "寅", "子": "寅", "辰": "寅", "寅": "申", "午": "申", "戌": "申", "亥": "巳", "卯": "巳", "未": "巳", "巳": "亥", "酉": "亥", "丑": "亥"}
	return ShenShaItem{"驿马", m[zhi]}
}
func taohua(zhi string) ShenShaItem {
	m := map[string]string{"申": "酉", "子": "酉", "辰": "酉", "寅": "卯", "午": "卯", "戌": "卯", "巳": "午", "酉": "午", "丑": "午", "亥": "子", "卯": "子", "未": "子"}
	return ShenShaItem{"桃花", m[zhi]}
}
func lushen(gan string) ShenShaItem {
	m := map[string]string{"甲": "寅", "乙": "卯", "丙": "巳", "丁": "午", "戊": "巳", "己": "午", "庚": "申", "辛": "酉", "壬": "亥", "癸": "子"}
	return ShenShaItem{"日禄", m[gan]}
}
func ride(gan string) ShenShaItem {
	m := map[string]string{"甲": "寅", "乙": "申", "丙": "巳", "丁": "亥", "戊": "巳", "己": "寅", "庚": "申", "辛": "巳", "壬": "亥", "癸": "巳"}
	return ShenShaItem{"日德", m[gan]}
}
func poshui(zhi string) ShenShaItem {
	m := map[string]string{"寅": "酉", "申": "酉", "巳": "酉", "亥": "酉", "子": "巳", "午": "巳", "卯": "巳", "酉": "巳", "辰": "丑", "戌": "丑", "丑": "丑", "未": "丑"}
	return ShenShaItem{"破碎", m[zhi]}
}
func wenchang(gan string) ShenShaItem {
	m := map[string]string{"甲": "巳", "乙": "午", "丙": "申", "丁": "酉", "戊": "申", "己": "酉", "庚": "亥", "辛": "子", "壬": "寅", "癸": "卯"}
	return ShenShaItem{"文昌", m[gan]}
}
func guiren(gan, hourZhi string) ShenShaItem {
	c1 := map[string]string{"甲": "丑", "戊": "丑", "乙": "申", "己": "申", "丙": "亥", "丁": "亥", "壬": "卯", "癸": "卯", "庚": "午", "辛": "午"}
	c2 := map[string]string{"甲": "未", "戊": "未", "乙": "子", "己": "子", "丙": "酉", "丁": "酉", "壬": "巳", "癸": "巳", "庚": "寅", "辛": "寅"}
	hIdx := indexOf(ZHI, hourZhi)
	if hIdx >= 3 && hIdx < 9 {
		return ShenShaItem{"贵人", c1[gan]}
	}
	return ShenShaItem{"贵人", c2[gan]}
}
func tianyi(zhi string) ShenShaItem {
	m := map[string]string{"寅": "辰", "卯": "巳", "辰": "午", "巳": "未", "午": "申", "未": "酉", "申": "戌", "酉": "亥", "戌": "子", "亥": "丑", "子": "寅", "丑": "卯"}
	return ShenShaItem{"天医", m[zhi]}
}
func diyi(zhi string) ShenShaItem {
	m := map[string]string{"寅": "戌", "卯": "亥", "辰": "子", "巳": "丑", "午": "寅", "未": "卯", "申": "辰", "酉": "巳", "戌": "午", "亥": "未", "子": "申", "丑": "酉"}
	return ShenShaItem{"地医", m[zhi]}
}
func tianma(zhi string) ShenShaItem {
	m := map[string]string{"寅": "午", "卯": "申", "辰": "戌", "巳": "子", "午": "寅", "未": "辰", "申": "午", "酉": "申", "戌": "戌", "亥": "子", "子": "寅", "丑": "辰"}
	return ShenShaItem{"天马", m[zhi]}
}
func tianxi(zhi string) ShenShaItem {
	m := map[string]string{"寅": "戌", "卯": "戌", "辰": "戌", "巳": "丑", "午": "丑", "未": "丑", "申": "辰", "酉": "辰", "戌": "辰", "亥": "未", "子": "未", "丑": "未"}
	return ShenShaItem{"天喜", m[zhi]}
}
func yuede(zhi string) ShenShaItem {
	m := map[string]string{"寅": "丙", "午": "丙", "戌": "丙", "申": "壬", "子": "壬", "辰": "壬", "亥": "甲", "卯": "甲", "未": "甲", "巳": "庚", "酉": "庚", "丑": "庚"}
	return ShenShaItem{"月德", m[zhi]}
}
func zaisha(zhi string) ShenShaItem {
	m := map[string]string{"寅": "子", "午": "子", "戌": "子", "申": "午", "子": "午", "辰": "午", "亥": "酉", "卯": "酉", "未": "酉", "巳": "卯", "酉": "卯", "丑": "卯"}
	return ShenShaItem{"灾煞", m[zhi]}
}
func jiesha(zhi string) ShenShaItem {
	m := map[string]string{"寅": "亥", "午": "亥", "戌": "亥", "申": "巳", "子": "巳", "辰": "巳", "亥": "申", "卯": "申", "未": "申", "巳": "寅", "酉": "寅", "丑": "寅"}
	return ShenShaItem{"劫煞", m[zhi]}
}
func jiangxing(zhi string) ShenShaItem {
	m := map[string]string{"申": "子", "子": "子", "辰": "子", "寅": "午", "午": "午", "戌": "午", "巳": "酉", "酉": "酉", "丑": "酉", "亥": "卯", "卯": "卯", "未": "卯"}
	return ShenShaItem{"将星", m[zhi]}
}
func huagai(zhi string) ShenShaItem {
	m := map[string]string{"申": "辰", "子": "辰", "辰": "辰", "寅": "戌", "午": "戌", "戌": "戌", "巳": "丑", "酉": "丑", "丑": "丑", "亥": "未", "卯": "未", "未": "未"}
	return ShenShaItem{"华盖", m[zhi]}
}
func yangren(gan string) ShenShaItem {
	m := map[string]string{"甲": "卯", "乙": "寅", "丙": "午", "丁": "巳", "戊": "午", "己": "巳", "庚": "酉", "辛": "申", "壬": "子", "癸": "亥"}
	return ShenShaItem{"羊刃", m[gan]}
}
func shengqi(zhi string) ShenShaItem {
	m := map[string]string{"寅": "子", "卯": "丑", "辰": "寅", "巳": "卯", "午": "辰", "未": "巳", "申": "午", "酉": "未", "戌": "申", "亥": "酉", "子": "戌", "丑": "亥"}
	return ShenShaItem{"生气", m[zhi]}
}
func siqi(zhi string) ShenShaItem {
	m := map[string]string{"寅": "午", "卯": "未", "辰": "申", "巳": "酉", "午": "戌", "未": "亥", "申": "子", "酉": "丑", "戌": "寅", "亥": "卯", "子": "辰", "丑": "巳"}
	return ShenShaItem{"死气", m[zhi]}
}
func wangshen(zhi string) ShenShaItem {
	m := map[string]string{"寅": "巳", "卯": "寅", "辰": "亥", "巳": "申", "午": "巳", "未": "寅", "申": "亥", "酉": "申", "戌": "巳", "亥": "寅", "子": "亥", "丑": "申"}
	return ShenShaItem{"亡神", m[zhi]}
}
func bingfu(zhi string) ShenShaItem {
	idx := indexOf(ZHI, zhi)
	return ShenShaItem{"病符", ZHI[(idx-1+12)%12]}
}
func shangmen(zhi string) ShenShaItem {
	idx := indexOf(ZHI, zhi)
	return ShenShaItem{"丧门", ZHI[(idx+2)%12]}
}
func diaoke(zhi string) ShenShaItem {
	idx := indexOf(ZHI, zhi)
	return ShenShaItem{"吊客", ZHI[(idx-2+12)%12]}
}
func guanfu(zhi string) ShenShaItem {
	m := map[string]string{"寅": "午", "卯": "未", "辰": "申", "巳": "酉", "午": "戌", "未": "亥", "申": "子", "酉": "丑", "戌": "寅", "亥": "卯", "子": "辰", "丑": "巳"}
	return ShenShaItem{"官符", m[zhi]}
}
func youdu(gan string) ShenShaItem {
	m := map[string]string{"甲": "丑", "己": "丑", "乙": "子", "庚": "子", "丙": "寅", "辛": "寅", "丁": "巳", "壬": "巳", "戊": "申", "癸": "申"}
	return ShenShaItem{"游都", m[gan]}
}
func feifu(gan string) ShenShaItem {
	m := map[string]string{"甲": "巳", "己": "午", "乙": "辰", "庚": "未", "丙": "卯", "辛": "申", "丁": "寅", "壬": "酉", "戊": "丑", "癸": "戌"}
	return ShenShaItem{"飞符", m[gan]}
}
func xuezhi(zhi string) ShenShaItem {
	m := map[string]string{"寅": "丑", "卯": "寅", "辰": "卯", "巳": "辰", "午": "巳", "未": "午", "申": "未", "酉": "申", "戌": "酉", "亥": "戌", "子": "亥", "丑": "子"}
	return ShenShaItem{"血支", m[zhi]}
}
func guashu(zhi string) ShenShaItem {
	m := map[string]string{"寅": "丑", "卯": "丑", "辰": "丑", "巳": "辰", "午": "辰", "未": "辰", "申": "未", "酉": "未", "戌": "未", "亥": "戌", "子": "戌", "丑": "戌"}
	return ShenShaItem{"寡宿", m[zhi]}
}
func gucheng(zhi string) ShenShaItem {
	m := map[string]string{"寅": "巳", "卯": "巳", "辰": "巳", "巳": "申", "午": "申", "未": "申", "申": "亥", "酉": "亥", "戌": "亥", "亥": "寅", "子": "寅", "丑": "寅"}
	return ShenShaItem{"孤辰", m[zhi]}
}
func sishen(zhi string) ShenShaItem {
	m := map[string]string{"寅": "巳", "卯": "午", "辰": "未", "巳": "申", "午": "酉", "未": "戌", "申": "亥", "酉": "子", "戌": "丑", "亥": "寅", "子": "卯", "丑": "辰"}
	return ShenShaItem{"死神", m[zhi]}
}

// guaShen 卦身
func guaShen(yaoList []YaoGanZhi) string {
	shiPos := -1
	for i, y := range yaoList {
		if y.ShiYing == 1 {
			shiPos = i
			break
		}
	}
	if shiPos == -1 {
		inner := reverse3(yaoList[0:3])
		outer := reverse3(yaoList[3:6])
		shiPos = calcShi(inner, outer)
	}
	startZhi := "子"
	if yaoList[shiPos].Yao == 0 {
		startZhi = "午"
	}
	startIdx := indexOf(ZHI, startZhi)
	endIdx := (startIdx + shiPos) % 12
	return ZHI[endIdx]
}

func xianggui(yaoList []YaoGanZhi, guaShenZhi string) string {
	gsWx := ZhiWuXing[guaShenZhi]
	result := ""
	for _, y := range yaoList {
		ywx := ZhiWuXing[y.Zhi]
		if WuXingKe[gsWx] == ywx {
			result += y.Zhi
		}
	}
	if result == "" {
		return "无"
	}
	return result
}

func chuangzhang(yaoList []YaoGanZhi, guaShenZhi string) string {
	gsWx := ZhiWuXing[guaShenZhi]
	result := ""
	for _, y := range yaoList {
		ywx := ZhiWuXing[y.Zhi]
		if WuXingSheng[gsWx] == ywx {
			result += y.Zhi
		}
	}
	if result == "" {
		return "无"
	}
	return result
}

// ============ 辅助函数 ============

func isGuiHun(inner, outer string) bool {
	for _, p := range GuiHunList {
		if p[0] == inner && p[1] == outer {
			return true
		}
	}
	return false
}

func isYouHun(inner, outer string) bool {
	for _, p := range YouHunList {
		if p[0] == inner && p[1] == outer {
			return true
		}
	}
	return false
}

func isLiuChong(inner, outer string) bool {
	for _, p := range LiuChongList {
		if p[0] == inner && p[1] == outer {
			return true
		}
	}
	return false
}

func isLiuHe(inner, outer string) bool {
	for _, p := range LiuHeList {
		if p[0] == inner && p[1] == outer {
			return true
		}
	}
	return false
}

func indexOf(arr []string, s string) int {
	for i, v := range arr {
		if v == s {
			return i
		}
	}
	return 0
}

func reverseMap(m map[string]string, val string) (string, bool) {
	for k, v := range m {
		if v == val {
			return k, true
		}
	}
	return "", false
}

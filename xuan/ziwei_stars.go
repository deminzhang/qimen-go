package xuan

import "github.com/6tail/lunar-go/calendar"

// ============ 辅星 (Minor Stars) ============

// ============ 庙旺亮度表 (iztro STARS_INFO) ============
// brightness[starIdx][palaceIdx], palaceIdx: 寅=0..丑=11

// 主星亮度表（14主星，按 StarZiWei..StarPoJun 顺序）
var zhuXingBrightness = [14][12]MiaoWang{
	// 紫微
	{Wang, Wang, DeDi, Wang, Miao, Miao, Wang, Wang, DeDi, Wang, PingHe, Miao},
	// 天机
	{DeDi, Wang, LiYi, PingHe, Miao, Xian, DeDi, Wang, LiYi, PingHe, Miao, Xian},
	// 太阳
	{Wang, Miao, Wang, Wang, Wang, DeDi, DeDi, Xian, Xian, Xian, Xian, Xian},
	// 武曲
	{DeDi, LiYi, Miao, PingHe, Wang, Miao, DeDi, LiYi, Miao, PingHe, Wang, Miao},
	// 天同
	{LiYi, PingHe, PingHe, Miao, Xian, Xian, Wang, PingHe, PingHe, Miao, Wang, Xian},
	// 廉贞
	{Miao, PingHe, LiYi, Xian, PingHe, LiYi, Miao, PingHe, LiYi, Xian, PingHe, LiYi},
	// 天府
	{Miao, DeDi, Miao, DeDi, Wang, Miao, DeDi, Wang, Miao, DeDi, Miao, Miao},
	// 太阴
	{Wang, Xian, Xian, Xian, Xian, Xian, LiYi, Xian, Wang, Miao, Miao, Miao},
	// 贪狼
	{PingHe, LiYi, Miao, Xian, Wang, Miao, PingHe, LiYi, Miao, Xian, Wang, Miao},
	// 巨门
	{Miao, Miao, Xian, Wang, Wang, Xian, Miao, Miao, Xian, Wang, Wang, Xian},
	// 天相
	{Miao, Xian, DeDi, DeDi, Miao, DeDi, Miao, Xian, DeDi, DeDi, Miao, Miao},
	// 天梁
	{Miao, Miao, Miao, Xian, Miao, Wang, Xian, DeDi, Miao, Xian, Miao, Wang},
	// 七杀
	{Miao, Wang, Miao, PingHe, Wang, Miao, Miao, Miao, Miao, PingHe, Wang, Miao},
	// 破军
	{DeDi, Xian, Wang, PingHe, Miao, Wang, DeDi, Xian, Wang, PingHe, Miao, Wang},
}

// 辅星亮度表（按 FuXingNames 顺序：左辅 右弼 文昌 文曲 天魁 天钺 禄存 天马 擎羊 陀罗 火星 铃星 地空 地劫）
var fuXingBrightness = map[string][12]MiaoWang{
	"文昌": {Xian, LiYi, DeDi, Miao, Xian, LiYi, DeDi, Miao, Xian, LiYi, DeDi, Miao},
	"文曲": {PingHe, Wang, DeDi, Miao, Xian, Wang, DeDi, Miao, Xian, Wang, DeDi, Miao},
	"火星": {Miao, LiYi, Xian, DeDi, Miao, LiYi, Xian, DeDi, Miao, LiYi, Xian, DeDi},
	"铃星": {Miao, LiYi, Xian, DeDi, Miao, LiYi, Xian, DeDi, Miao, LiYi, Xian, DeDi},
	"擎羊": {MiaoWangNone, Xian, Miao, MiaoWangNone, Xian, Miao, MiaoWangNone, Xian, Miao, MiaoWangNone, Xian, Miao},
	"陀罗": {Xian, MiaoWangNone, Miao, Xian, MiaoWangNone, Miao, Xian, MiaoWangNone, Miao, Xian, MiaoWangNone, Miao},
}

// getZhuXingBrightness 获取主星在某宫位的亮度
func getZhuXingBrightness(starIdx, palIdx int) MiaoWang {
	if starIdx < 0 || starIdx >= 14 {
		return MiaoWangNone
	}
	return zhuXingBrightness[starIdx][palIdx]
}

// getFuXingBrightness 获取辅星在某宫位的亮度
func getFuXingBrightness(name string, palIdx int) MiaoWang {
	if b, ok := fuXingBrightness[name]; ok {
		return b[palIdx]
	}
	return MiaoWangNone // 无亮度表的辅星不显示庙旺（对齐iztro）
}

// 辅星名称（14颗）

var ZaYaoNames = []string{
	"红鸾", "天喜", "天姚", "咸池", "解神",
	"三台", "八座", "恩光", "天贵",
	"龙池", "凤阁", "天才", "天寿",
	"台辅", "封诰", "天巫", "华盖", "天官", "天福", "天厨",
	"天月", "天德", "月德", "天空", "旬空",
	"孤辰", "寡宿", "蜚蠊", "破碎",
	"天刑", "阴煞", "天哭", "天虚", "天使", "天伤",
	"截空", "劫煞", "年解", "大耗",
}

// 长生十二神
var ChangSheng12Names = []string{
	"长生", "沐浴", "冠带", "临官", "帝旺", "衰", "病", "死", "墓", "绝", "胎", "养",
}

// 博士十二神
var BoShi12Names = []string{
	"博士", "力士", "青龙", "小耗", "将军", "奏书", "飞廉", "喜神", "病符", "大耗", "伏兵", "官府",
}

// 岁前十二神
var SuiQian12Names = []string{
	"岁建", "晦气", "丧门", "贯索", "官符", "小耗", "大耗", "龙德", "白虎", "天德", "吊客", "病符",
}

// 将前十二神
var JiangQian12Names = []string{
	"将星", "攀鞍", "岁驿", "息神", "华盖", "劫煞", "灾煞", "天煞", "指背", "咸池", "月煞", "亡神",
}

// ============ 辅助函数 ============

// zhiToIndex 地支字符串 → 宫位索引（寅=0）
func zhiToIndex(zhi string) int {
	for i, z := range ZiWeiPalaceZhi {
		if z == zhi {
			return i
		}
	}
	return -1
}

// zhiStdIndex 返回标准地支序（子=0, 丑=1, ..., 亥=11），用于 iztro EARTHLY_BRANCHES.indexOf 的等价操作
func zhiStdIndex(zhi string) int {
	stdZhi := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	for i, z := range stdZhi {
		if z == zhi {
			return i
		}
	}
	return -1
}
// getLuYangTuoMaIndex 禄存、擎羊、陀罗、天马索引（按年干年支）
func getLuYangTuoMaIndex(gan, zhi string) (luIdx, maIdx, yangIdx, tuoIdx int) {
	// 天马：按年支
	switch zhi {
	case "寅", "午", "戌":
		maIdx = zhiToIndex("申")
	case "申", "子", "辰":
		maIdx = zhiToIndex("寅")
	case "巳", "酉", "丑":
		maIdx = zhiToIndex("亥")
	case "亥", "卯", "未":
		maIdx = zhiToIndex("巳")
	}

	// 禄存：按年干
	switch gan {
	case "甲":
		luIdx = zhiToIndex("寅")
	case "乙":
		luIdx = zhiToIndex("卯")
	case "丙", "戊":
		luIdx = zhiToIndex("巳")
	case "丁", "己":
		luIdx = zhiToIndex("午")
	case "庚":
		luIdx = zhiToIndex("申")
	case "辛":
		luIdx = zhiToIndex("酉")
	case "壬":
		luIdx = zhiToIndex("亥")
	case "癸":
		luIdx = zhiToIndex("子")
	}

	yangIdx = fix12(luIdx + 1)
	tuoIdx = fix12(luIdx - 1)
	return
}

// getKuiYueIndex 天魁天钺索引（按年干）
func getKuiYueIndex(gan string) (kuiIdx, yueIdx int) {
	switch gan {
	case "甲", "戊", "庚":
		kuiIdx = zhiToIndex("丑")
		yueIdx = zhiToIndex("未")
	case "乙", "己":
		kuiIdx = zhiToIndex("子")
		yueIdx = zhiToIndex("申")
	case "辛":
		kuiIdx = zhiToIndex("午")
		yueIdx = zhiToIndex("寅")
	case "丙", "丁":
		kuiIdx = zhiToIndex("亥")
		yueIdx = zhiToIndex("酉")
	case "壬", "癸":
		kuiIdx = zhiToIndex("卯")
		yueIdx = zhiToIndex("巳")
	}
	return
}

// getZuoYouIndex 左辅右弼索引（按农历月）
func getZuoYouIndex(lunarMonth int) (zuoIdx, youIdx int) {
	chenIdx := zhiToIndex("辰")
	xuIdx := zhiToIndex("戌")
	zuoIdx = fix12(chenIdx + (lunarMonth - 1))
	youIdx = fix12(xuIdx - (lunarMonth - 1))
	return
}

// getChangQuIndex 文昌文曲索引（按时辰）
func getChangQuIndex(timeIndex int) (changIdx, quIdx int) {
	xuIdx := zhiToIndex("戌")
	chenIdx := zhiToIndex("辰")
	fixedTimeIdx := fix12(timeIndex)
	changIdx = fix12(xuIdx - fixedTimeIdx)
	quIdx = fix12(chenIdx + fixedTimeIdx)
	return
}

// getKongJieIndex 地空地劫索引（按时辰）
func getKongJieIndex(timeIndex int) (kongIdx, jieIdx int) {
	haiIdx := zhiToIndex("亥")
	fixedTimeIdx := fix12(timeIndex)
	kongIdx = fix12(haiIdx - fixedTimeIdx)
	jieIdx = fix12(haiIdx + fixedTimeIdx)
	return
}

// getHuoLingIndex 火星铃星索引（按年支+时辰）
func getHuoLingIndex(yearZhi string, timeIndex int) (huoIdx, lingIdx int) {
	var huoBase, lingBase int
	switch yearZhi {
	case "寅", "午", "戌":
		huoBase = zhiToIndex("丑")
		lingBase = zhiToIndex("卯")
	case "申", "子", "辰":
		huoBase = zhiToIndex("寅")
		lingBase = zhiToIndex("戌")
	case "巳", "酉", "丑":
		huoBase = zhiToIndex("卯")
		lingBase = zhiToIndex("戌")
	case "亥", "卯", "未":
		huoBase = zhiToIndex("酉")
		lingBase = zhiToIndex("戌")
	}
	fixedTimeIdx := fix12(timeIndex)
	huoIdx = fix12(huoBase + fixedTimeIdx)
	lingIdx = fix12(lingBase + fixedTimeIdx)
	return
}
// getLuanXiIndex 红鸾天喜索引（按年支）
// 卯上起子逆数之，数到当生太岁支
func getLuanXiIndex(yearZhi string) (hongluanIdx, tianxiIdx int) {
	// fixEarthlyBranchIndex('mao') = palace idx 1, then - EARTHLY_BRANCHES.indexOf
	hongluanIdx = fix12(zhiToIndex("卯") - zhiStdIndex(yearZhi))
	tianxiIdx = fix12(hongluanIdx + 6)
	return
}

// getHuagaiXianchiIndex 华盖咸池索引（按年支）
func getHuagaiXianchiIndex(yearZhi string) (huagaiIdx, xianchiIdx int) {
	switch yearZhi {
	case "寅", "午", "戌":
		huagaiIdx = zhiToIndex("戌")
		xianchiIdx = zhiToIndex("卯")
	case "申", "子", "辰":
		huagaiIdx = zhiToIndex("辰")
		xianchiIdx = zhiToIndex("酉")
	case "巳", "酉", "丑":
		huagaiIdx = zhiToIndex("丑")
		xianchiIdx = zhiToIndex("午")
	case "亥", "卯", "未":
		huagaiIdx = zhiToIndex("未")
		xianchiIdx = zhiToIndex("子")
	}
	return
}

// getGuGuaIndex 孤辰寡宿索引（按年支）
func getGuGuaIndex(yearZhi string) (guIdx, guaIdx int) {
	switch yearZhi {
	case "寅", "卯", "辰":
		guIdx = zhiToIndex("巳")
		guaIdx = zhiToIndex("丑")
	case "巳", "午", "未":
		guIdx = zhiToIndex("申")
		guaIdx = zhiToIndex("辰")
	case "申", "酉", "戌":
		guIdx = zhiToIndex("亥")
		guaIdx = zhiToIndex("未")
	case "亥", "子", "丑":
		guIdx = zhiToIndex("寅")
		guaIdx = zhiToIndex("戌")
	}
	return
}

// getJieshaAdjIndex 劫煞索引（按年支）
func getJieshaAdjIndex(yearZhi string) int {
	switch yearZhi {
	case "申", "子", "辰":
		return 3
	case "亥", "卯", "未":
		return 6
	case "寅", "午", "戌":
		return 9
	case "巳", "酉", "丑":
		return 0
	}
	return -1
}

// getDahaoAdjIndex 大耗索引（按年支）
func getDahaoAdjIndex(yearZhi string) int {
	zhiList := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	matched := []string{"未", "午", "酉", "申", "亥", "戌", "丑", "子", "卯", "寅", "巳", "辰"}
	for i, z := range zhiList {
		if z == yearZhi {
			return fix12(zhiToIndex(matched[i]))
		}
	}
	return -1
}

// getNianjieIndex 年解索引（按年支）
func getNianjieIndex(yearZhi string) int {
	zhiList := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	nianjieList := []string{"戌", "酉", "申", "未", "午", "巳", "辰", "卯", "寅", "丑", "子", "亥"}
	for i, z := range zhiList {
		if z == yearZhi {
			return fix12(zhiToIndex(nianjieList[i]))
		}
	}
	return -1
}

// ============ 日系星 ============

// getDailyStarIndex 日系星索引（三台、八座、恩光、天贵）
func getDailyStarIndex(solarDateStr string, timeIndex int) (santaiIdx, bazuoIdx, enguangIdx, tianguiIdx int) {
	solar, err := parseSolarDate(solarDateStr)
	if err != nil {
		return
	}
	lunar := calendar.NewLunarFromSolar(solar)
	day := lunar.GetDay()
	month := lunar.GetMonth()

	// fixLunarDayIndex: 晚子时 +1
	if timeIndex >= 12 {
		day = day
		// 不加1，保持原day（与iztro fixLunarDayIndex一致: timeIndex>=12 return lunarDay）
	} else {
		day = day - 1
	}
	// 月份索引 (fixLunarMonthIndex: lunarMonth + 1 - 2 = lunarMonth - 1)
	monthIndex := month - 1

	zuoIdx, youIdx := getZuoYouIndex(monthIndex + 1)
	changIdx, quIdx := getChangQuIndex(timeIndex)

	santaiIdx = fix12((zuoIdx + day) % 12)
	bazuoIdx = fix12((youIdx - day) % 12)
	enguangIdx = fix12(((changIdx+day)%12 - 1) % 12)
	tianguiIdx = fix12(((quIdx+day)%12 - 1) % 12)
	return
}

// ============ 时系星 ============

// getTimelyStarIndex 时系星索引（台辅、封诰）
func getTimelyStarIndex(timeIndex int) (taifuIdx, fenggaoIdx int) {
	wuIdx := zhiToIndex("午")
	yinIdx := zhiToIndex("寅")
	fixedTimeIdx := fix12(timeIndex)
	taifuIdx = fix12(wuIdx + fixedTimeIdx)
	fenggaoIdx = fix12(yinIdx + fixedTimeIdx)
	return
}

// ============ 月系星 ============

// getMonthlyStarIndex 月系星索引（解神/天姚/天刑/阴煞/天月/天巫）
func getMonthlyStarIndex(solarDateStr string, timeIndex int) (yuejieIdx, tianyaoIdx, tianxingIdx, yinshaIdx, tianyueIdx, tianwuIdx int) {
	solar, err := parseSolarDate(solarDateStr)
	if err != nil {
		return
	}
	lunar := calendar.NewLunarFromSolar(solar)
	month := lunar.GetMonth()
	monthIndex := month - 1 // fixLunarMonthIndex

	// 解神：按生月
	jieshenMaps := []string{"申", "戌", "子", "寅", "辰", "午"}
	yuejieIdx = fix12(zhiToIndex(jieshenMaps[monthIndex/2]))

	// 天姚：丑起正月顺数
	tianyaoIdx = fix12(zhiToIndex("丑") + monthIndex)

	// 天刑：酉起正月顺数
	tianxingIdx = fix12(zhiToIndex("酉") + monthIndex)

	// 阴煞
	yinshaMaps := []string{"寅", "子", "戌", "申", "午", "辰"}
	yinshaIdx = fix12(zhiToIndex(yinshaMaps[monthIndex%6]))

	// 天月
	tianyueMaps := []string{"戌", "巳", "辰", "寅", "未", "卯", "亥", "未", "寅", "午", "戌", "寅"}
	tianyueIdx = fix12(zhiToIndex(tianyueMaps[monthIndex]))

	// 天巫
	tianwuMaps := []string{"巳", "申", "寅", "亥"}
	tianwuIdx = fix12(zhiToIndex(tianwuMaps[monthIndex%4]))
	return
}

// ============ 年系星 ============

// getYearlyStarIndex 年系星索引（咸池/华盖/孤辰/寡宿/天厨/破碎/天才/天寿/蜚蠊/龙池/凤阁/天哭/天虚/天官/天福/天空/旬空/截空/天使/天伤/天德/月德/劫煞/年解/大耗）
func getYearlyStarIndex(solarDateStr string, timeIndex int, gender int, soulIdx, bodyIdx int) map[string]int {
	solar, _ := parseSolarDate(solarDateStr)
	if solar == nil {
		return nil
	}
	lunar := calendar.NewLunarFromSolar(solar)
	ganZhiYear := lunar.GetYearInGanZhiExact()
	yearGan := string([]rune(ganZhiYear)[0])
	yearZhi := string([]rune(ganZhiYear)[1])

	yearZhiIdx := zhiToIndex(yearZhi)     // 宫位序（寅=0）
	yearZhiStd := zhiStdIndex(yearZhi)    // 标准序（子=0），用于 iztro EARTHLY_BRANCHES.indexOf 等价
	ganIdx := indexOf(TianGanList, yearGan)

	r := make(map[string]int)

	// 咸池/华盖
	hg, xc := getHuagaiXianchiIndex(yearZhi)
	r["huagaiIndex"] = hg
	r["xianchiIndex"] = xc

	// 孤辰/寡宿
	gu, gua := getGuGuaIndex(yearZhi)
	r["guchenIndex"] = gu
	r["guasuIndex"] = gua

	// 天才/天寿（iztro: soulIndex + EARTHLY_BRANCHES.indexOf, bodyIndex + EARTHLY_BRANCHES.indexOf）
	r["tiancaiIndex"] = fix12(soulIdx + yearZhiStd)
	r["tianshouIndex"] = fix12(bodyIdx + yearZhiStd)

	// 天厨
	tianchuMaps := []string{"巳", "午", "子", "巳", "午", "申", "寅", "午", "酉", "亥"}
	r["tianchuIndex"] = fix12(zhiToIndex(tianchuMaps[ganIdx]))

	// 破碎（iztro: EARTHLY_BRANCHES.indexOf(earthlyBranch) % 3）
	posuiMaps := []string{"巳", "丑", "酉"}
	r["posuiIndex"] = fix12(zhiToIndex(posuiMaps[yearZhiStd%3]))

	// 蜚蠊（iztro: EARTHLY_BRANCHES.indexOf 查表）
	feilianMaps := []string{"申", "酉", "戌", "巳", "午", "未", "寅", "卯", "辰", "亥", "子", "丑"}
	r["feilianIndex"] = fix12(zhiToIndex(feilianMaps[yearZhiStd]))

	// 龙池/凤阁（iztro: fixEarthlyBranchIndex('chen') + EARTHLY_BRANCHES.indexOf, fixEarthlyBranchIndex('xu') - EARTHLY_BRANCHES.indexOf）
	r["longchiIndex"] = fix12(zhiToIndex("辰") + yearZhiStd)
	r["fenggeIndex"] = fix12(zhiToIndex("戌") - yearZhiStd)

	// 天哭/天虚（iztro: fixEarthlyBranchIndex('woo') ± EARTHLY_BRANCHES.indexOf）
	r["tiankuIndex"] = fix12(zhiToIndex("午") - yearZhiStd)
	r["tianxuIndex"] = fix12(zhiToIndex("午") + yearZhiStd)

	// 天官
	tianguanMaps := []string{"未", "辰", "巳", "寅", "卯", "酉", "亥", "酉", "戌", "午"}
	r["tianguanIndex"] = fix12(zhiToIndex(tianguanMaps[ganIdx]))

	// 天福
	tianfuMaps := []string{"酉", "申", "子", "亥", "卯", "寅", "午", "巳", "午", "巳"}
	r["tianfuIndex"] = fix12(zhiToIndex(tianfuMaps[ganIdx]))

	// 天德/月德（iztro: fixEarthlyBranchIndex('you') + EARTHLY_BRANCHES.indexOf, fixEarthlyBranchIndex('si') + EARTHLY_BRANCHES.indexOf）
	r["tiandeIndex"] = fix12(zhiToIndex("酉") + yearZhiStd)
	r["yuedeIndex"] = fix12(zhiToIndex("巳") + yearZhiStd)

	// 天空（iztro: fixEarthlyBranchIndex(yearly[1]) + 1，已是宫位序）
	r["tiankongIndex"] = fix12(yearZhiIdx + 1)

	// 截路空亡
	jieluMaps := []string{"申", "午", "辰", "寅", "子"}
	kongwangMaps := []string{"酉", "未", "巳", "卯", "丑"}
	r["jieluIndex"] = fix12(zhiToIndex(jieluMaps[ganIdx%5]))
	r["kongwangIndex"] = fix12(zhiToIndex(kongwangMaps[ganIdx%5]))

	// 旬空（iztro: fixEarthlyBranchIndex(yearly[1]) + ...  年支已是宫位序）
	yearGanGuiIdx := indexOf(TianGanList, "癸")
	xunkongIdx := fix12(yearZhiIdx + yearGanGuiIdx - ganIdx + 1)
	yinyang := yearZhiStd % 2 // iztro: EARTHLY_BRANCHES.indexOf(earthlyBranch) % 2
	if yinyang != xunkongIdx%2 {
		xunkongIdx = fix12(xunkongIdx + 1)
	}
	r["xunkongIndex"] = xunkongIdx

	// 截空（生年阳干在阳宫，阴干在阴宫）
	if yinyang == 0 {
		r["jiekongIndex"] = r["jieluIndex"]
	} else {
		r["jiekongIndex"] = r["kongwangIndex"]
	}

	// 劫煞
	r["jieshaAdjIndex"] = getJieshaAdjIndex(yearZhi)

	// 年解
	r["nianjieIndex"] = getNianjieIndex(yearZhi)

	// 大耗
	r["dahaoAdjIndex"] = getDahaoAdjIndex(yearZhi)

	// 天使/天伤
	soulIndex := soulIdx
	tianshangIdx := fix12(zhiToIndex("亥") + soulIndex) // 奴仆宫起点
	tianshiIdx := fix12(zhiToIndex("酉") + soulIndex)   // 疾厄宫起点
	r["tianshangIndex"] = tianshangIdx
	r["tianshiIndex"] = tianshiIdx

	return r
}

// ============ 长生十二神 ============

// getChangSheng12StartIndex 长生十二神起始宫位（按五行局）
func getChangSheng12StartIndex(wx WuXingJu) int {
	switch WuXingJuNums[wx] {
	case 2: // 水二局
		return zhiToIndex("申")
	case 3: // 木三局
		return zhiToIndex("亥")
	case 4: // 金四局
		return zhiToIndex("巳")
	case 5: // 土五局
		return zhiToIndex("申")
	case 6: // 火六局
		return zhiToIndex("寅")
	}
	return 0
}

// getChangSheng12 长生十二神安星
// 阳男阴女顺行，阴男阳女逆行
func getChangSheng12(yearGan, yearZhi string, gender int, wx WuXingJu) [12]string {
	var result [12]string
	startIdx := getChangSheng12StartIndex(wx)

	// 年支阴阳
	yinyang := zhiToIndex(yearZhi) % 2 // 0=阳, 1=阴
	genderYinYang := 0
	if gender == 0 { // 女=阴
		genderYinYang = 1
	}
	// 阳男阴女顺行
	shunXing := (yinyang == genderYinYang)

	for i := 0; i < 12; i++ {
		var idx int
		if shunXing {
			idx = fix12(i + startIdx)
		} else {
			idx = fix12(startIdx - i)
		}
		result[idx] = ChangSheng12Names[i]
	}
	return result
}

// ============ 博士十二神 ============

// getBoShi12 博士十二神安星（从禄存起，阳男阴女顺行，阴男阳女逆行）
func getBoShi12(yearGan, yearZhi string, gender int) [12]string {
	var result [12]string
	luIdx, _, _, _ := getLuYangTuoMaIndex(yearGan, yearZhi)

	yinyang := zhiToIndex(yearZhi) % 2
	genderYinYang := 0
	if gender == 0 {
		genderYinYang = 1
	}
	shunXing := (yinyang == genderYinYang)

	for i := 0; i < 12; i++ {
		var idx int
		if shunXing {
			idx = fix12(luIdx + i)
		} else {
			idx = fix12(luIdx - i)
		}
		result[idx] = BoShi12Names[i]
	}
	return result
}

// ============ 流年岁前十二神 ============

// getJiangQian12StartIndex 将前十二神起始宫位（按年支）
func getJiangQian12StartIndex(yearZhi string) int {
	switch yearZhi {
	case "寅", "午", "戌":
		return zhiToIndex("午")
	case "申", "子", "辰":
		return zhiToIndex("子")
	case "巳", "酉", "丑":
		return zhiToIndex("酉")
	case "亥", "卯", "未":
		return zhiToIndex("卯")
	}
	return -1
}

// getSuiQian12 岁前十二神（从年支起岁建顺行）
func getSuiQian12(yearZhi string) [12]string {
	var result [12]string
	startIdx := zhiToIndex(yearZhi)
	for i := 0; i < 12; i++ {
		result[fix12(startIdx+i)] = SuiQian12Names[i]
	}
	return result
}

// getJiangQian12 将前十二神
func getJiangQian12(yearZhi string) [12]string {
	var result [12]string
	startIdx := getJiangQian12StartIndex(yearZhi)
	for i := 0; i < 12; i++ {
		result[fix12(startIdx+i)] = JiangQian12Names[i]
	}
	return result
}

// ============ 综合安星函数 ============

// setupFuXing 安辅星（14颗）
func (c *ZiWeiChart) setupFuXing() {
	luIdx, maIdx, yangIdx, tuoIdx := getLuYangTuoMaIndex(c.YearGan, c.YearZhi)
	kuiIdx, yueIdx := getKuiYueIndex(c.YearGan)
	zuoIdx, youIdx := getZuoYouIndex(c.MonthNum)
	changIdx, quIdx := getChangQuIndex(c.HourIdx)
	kongIdx, jieIdx := getKongJieIndex(c.HourIdx)
	huoIdx, lingIdx := getHuoLingIndex(c.YearZhi, c.HourIdx)

	addStar := func(palIdx int, name string) {
		b := getFuXingBrightness(name, palIdx)
		c.Palaces[palIdx].FuXing = append(c.Palaces[palIdx].FuXing, Star{Name: name, MiaoWang: b})
	}

	addStar(luIdx, "禄存")
	addStar(maIdx, "天马")
	addStar(yangIdx, "擎羊")
	addStar(tuoIdx, "陀罗")
	addStar(kuiIdx, "天魁")
	addStar(yueIdx, "天钺")
	addStar(zuoIdx, "左辅")
	addStar(youIdx, "右弼")
	addStar(changIdx, "文昌")
	addStar(quIdx, "文曲")
	addStar(kongIdx, "地空")
	addStar(jieIdx, "地劫")
	addStar(huoIdx, "火星")
	addStar(lingIdx, "铃星")
}

// setupZaYao 安杂耀
func (c *ZiWeiChart) setupZaYao(solarDateStr string, timeIndex int) {
	addStar := func(palIdx int, name string) {
		c.Palaces[palIdx].ZaYao = append(c.Palaces[palIdx].ZaYao, Star{Name: name})
	}

	// 红鸾/天喜
	hlIdx, txIdx := getLuanXiIndex(c.YearZhi)
	addStar(hlIdx, "红鸾")
	addStar(txIdx, "天喜")

	// 月系星
	yuejie, tianyao, tianxing, yinsha, tianyue, tianwu := getMonthlyStarIndex(solarDateStr, timeIndex)
	addStar(yuejie, "解神")
	addStar(tianyao, "天姚")
	addStar(tianxing, "天刑")
	addStar(yinsha, "阴煞")
	addStar(tianyue, "天月")
	addStar(tianwu, "天巫")

	// 日系星
	santai, bazuo, enguang, tiangui := getDailyStarIndex(solarDateStr, timeIndex)
	addStar(santai, "三台")
	addStar(bazuo, "八座")
	addStar(enguang, "恩光")
	addStar(tiangui, "天贵")

	// 时系星
	taifu, fenggao := getTimelyStarIndex(c.HourIdx)
	addStar(taifu, "台辅")
	addStar(fenggao, "封诰")

	// 年系星
	yi := getYearlyStarIndex(solarDateStr, timeIndex, c.Gender, c.MingGongIdx, c.ShenGongIdx)
	addStar(yi["xianchiIndex"], "咸池")
	addStar(yi["huagaiIndex"], "华盖")
	addStar(yi["guchenIndex"], "孤辰")
	addStar(yi["guasuIndex"], "寡宿")
	addStar(yi["tiancaiIndex"], "天才")
	addStar(yi["tianshouIndex"], "天寿")
	addStar(yi["tianchuIndex"], "天厨")
	addStar(yi["posuiIndex"], "破碎")
	addStar(yi["feilianIndex"], "蜚蠊")
	addStar(yi["longchiIndex"], "龙池")
	addStar(yi["fenggeIndex"], "凤阁")
	addStar(yi["tiankuIndex"], "天哭")
	addStar(yi["tianxuIndex"], "天虚")
	addStar(yi["tianguanIndex"], "天官")
	addStar(yi["tianfuIndex"], "天福")
	addStar(yi["tiandeIndex"], "天德")
	addStar(yi["yuedeIndex"], "月德")
	addStar(yi["tiankongIndex"], "天空")
	addStar(yi["xunkongIndex"], "旬空")
	addStar(yi["jiekongIndex"], "截空")
	addStar(yi["jieshaAdjIndex"], "劫煞")
	addStar(yi["nianjieIndex"], "年解")
	addStar(yi["dahaoAdjIndex"], "大耗")
	addStar(yi["tianshangIndex"], "天伤")
	addStar(yi["tianshiIndex"], "天使")
}

// setupDecorative 安长生十二神、博士十二神、岁前十二神、将前十二神
func (c *ZiWeiChart) setupDecorative() {
	// 长生十二神
	cs12 := getChangSheng12(c.YearGan, c.YearZhi, c.Gender, c.WuXingJu)
	for i := 0; i < 12; i++ {
		c.Palaces[i].ChangSheng = cs12[i]
	}

	// 博士十二神
	bs12 := getBoShi12(c.YearGan, c.YearZhi, c.Gender)
	for i := 0; i < 12; i++ {
		c.Palaces[i].BoShi = bs12[i]
	}

	// 岁前十二神
	sq12 := getSuiQian12(c.YearZhi)
	for i := 0; i < 12; i++ {
		c.Palaces[i].SuiQian = sq12[i]
	}

	// 将前十二神
	jq12 := getJiangQian12(c.YearZhi)
	for i := 0; i < 12; i++ {
		c.Palaces[i].JiangQian = jq12[i]
	}
}

package qimen

import (
	"github.com/6tail/lunar-go/calendar"
	"qimen/util"
)

// CalcShenSha 神煞算法
func CalcShenSha(bz *calendar.EightChar) ([]string, []string, []string, []string) {
	xingY := string([]rune(bz.GetYearNaYin())[2:])
	zhuD := bz.GetDay()
	zhuT := bz.GetTime()
	ganY, ganM, ganD, ganT := bz.GetYearGan(), bz.GetMonthGan(), bz.GetDayGan(), bz.GetTimeGan()
	zhiY, zhiM, zhiD, zhiT := bz.GetYearZhi(), bz.GetMonthZhi(), bz.GetDayZhi(), bz.GetTimeZhi()
	ganYD := []string{ganY, ganD}
	ganDA := []string{ganD}
	zhiYD := []string{zhiY, zhiD}
	zhiA := []string{zhiY, zhiM, zhiD, zhiT}
	ganA := []string{ganY, ganM, ganD, ganT}
	gzAA := [][]string{{ganY, zhiY}, {ganM, zhiM}, {ganD, zhiD}, {ganT, zhiT}}
	zhiAA := [][]string{{zhiY}, {zhiM}, {zhiD}, {zhiT}}
	var ss = make([][]string, 4)

	//贵人(天乙 福星 大极 天官 国印 天厨 文昌 玉堂 金匮)
	//天乙贵人: 年/日干见它支(甲戊并牛羊，乙己鼠猴乡，丙丁猪鸡位，壬癸兔蛇藏，庚辛逢虎马，此是贵人方)
	for i, zhi := range zhiAA {
		if (util.Contains(zhi, "午", "未") && util.Contains(ganYD, "甲", "戊")) ||
			(util.Contains(zhi, "子", "申") && util.Contains(ganYD, "乙", "己")) ||
			(util.Contains(zhi, "亥", "酉") && util.Contains(ganYD, "丙", "丁")) ||
			(util.Contains(zhi, "卯", "巳") && util.Contains(ganYD, "壬", "癸")) ||
			(util.Contains(zhi, "寅", "午") && util.Contains(ganYD, "庚", "辛")) {
			ss[i] = append(ss[i], "天乙贵人")
		}
	}
	//福星贵人: 年/日干见它支(甲见子寅 乙见丑卯 丙见子寅 丁亥 戊申 己未 庚午 辛巳 壬辰 癸见丑卯)
	for i, zhi := range zhiAA {
		if (util.Contains(zhi, "子", "寅") && util.Contains(ganYD, "甲")) ||
			(util.Contains(zhi, "丑", "卯") && util.Contains(ganYD, "乙")) ||
			(util.Contains(zhi, "子", "寅") && util.Contains(ganYD, "丙")) ||
			(util.Contains(zhi, "亥") && util.Contains(ganYD, "丁")) ||
			(util.Contains(zhi, "申") && util.Contains(ganYD, "戊")) ||
			(util.Contains(zhi, "未") && util.Contains(ganYD, "己")) ||
			(util.Contains(zhi, "午") && util.Contains(ganYD, "庚")) ||
			(util.Contains(zhi, "巳") && util.Contains(ganYD, "辛")) ||
			(util.Contains(zhi, "辰") && util.Contains(ganYD, "壬")) ||
			(util.Contains(zhi, "丑", "卯") && util.Contains(ganYD, "癸")) {
			ss[i] = append(ss[i], "福星贵人")
		}
	}
	//大极贵人: 日支见它支(甲乙见子午 丙丁见卯酉 戊已见四库 庚辛见寅亥 壬癸见巳申)
	for i, zhi := range zhiAA {
		if (util.Contains(zhi, "子", "午") && util.Contains(ganDA, "甲", "乙")) ||
			(util.Contains(zhi, "卯", "酉") && util.Contains(ganDA, "丙", "丁")) ||
			(util.Contains(zhi, "寅", "亥") && util.Contains(ganDA, "戊", "己")) ||
			(util.Contains(zhi, "巳", "申") && util.Contains(ganDA, "庚", "辛")) ||
			(util.Contains(zhi, "子", "午") && util.Contains(ganDA, "壬", "癸")) {
			ss[i] = append(ss[i], "大极贵人")
		}
	}
	//天官贵人: 日支见它支(甲未 乙辰 丙巳 丁寅 戊丑 己戌 庚亥 辛申 壬酉 癸午)
	for i, zhi := range zhiA {
		if util.Contains([]string{"甲未", "乙辰", "丙巳", "丁寅", "戊丑", "己戌", "庚亥", "辛申", "壬酉", "癸午"}, zhiD+zhi) {
			ss[i] = append(ss[i], "天官贵人")
		}
	}
	//国印贵人: 日支见它支(甲戌 乙亥、丙丑、丁寅、戊寅、己寅、庚辰、辛巳、壬未、癸申)
	for i, zhi := range zhiA {
		if util.Contains([]string{"甲戌", "乙亥", "乙亥", "丙丑", "丁寅", "戊寅", "己寅", "庚辰", "辛巳", "壬未", "癸申"}, zhiD+zhi) {
			ss[i] = append(ss[i], "国印贵人")
		}
	}
	//天厨贵人: 年/日干见它食神禄支(甲巳 乙午 丙巳 丁午 戊申 己酉 庚亥 辛子 壬寅 癸卯)
	for i, zhi := range zhiA {
		if util.Contains([]string{"甲巳", "乙午", "丙巳", "丁午", "戊申", "己酉", "庚亥", "辛子", "壬寅", "癸卯"}, zhiY+zhi, zhiD+zhi) {
			ss[i] = append(ss[i], "天厨贵人")
		}
	}
	//文昌贵人：日干见它支(甲巳 乙午 丙申 丁酉 戊申 己酉 庚亥 辛子 壬寅 癸卯)
	for i, zhi := range zhiA {
		if util.Contains([]string{"甲巳", "乙午", "丙申", "丁酉", "戊申", "己酉", "庚亥", "辛子", "壬寅", "癸卯"}, zhiD+zhi) {
			ss[i] = append(ss[i], "文昌贵人")
		}
	}
	//玉堂贵人：日干见它支(甲未 乙辰 丙巳 丁酉 戊戌 己卯 庚丑 辛申 壬寅 癸午)
	for i, zhi := range zhiA {
		if util.Contains([]string{"甲未", "乙辰", "丙巳", "丁酉", "戊戌", "己卯", "庚丑", "辛申", "壬寅", "癸午"}, zhiD+zhi) {
			ss[i] = append(ss[i], "玉堂贵人")
		}
	}
	//天德贵人：月支查干支 (子巳 丑庚 寅丁 卯申 辰壬 巳辛 午亥 未甲 申癸 酉寅 戌丙 亥丁)
	for i, gz := range gzAA {
		if util.Contains([]string{"子巳", "丑庚", "寅丁", "卯申", "辰壬", "巳辛", "午亥", "未甲", "申癸", "酉寅", "戌丙", "亥丁"}, zhiM+gz[1], zhiM+gz[0]) {
			ss[i] = append(ss[i], "天德贵人")
		}
	}
	//月德贵人：月支见天干 三合见阳干(寅午戌见丙,申子辰见壬,亥卯未见甲,巳酉丑见庚)
	for i, gan := range ganA {
		if (util.Contains([]string{"寅", "午", "戌"}, zhiM) && "丙" == gan) ||
			(util.Contains([]string{"申", "子", "辰"}, zhiM) && "壬" == gan) ||
			(util.Contains([]string{"亥", "卯", "未"}, zhiM) && "甲" == gan) ||
			(util.Contains([]string{"巳", "酉", "丑"}, zhiM) && "庚" == gan) {
			ss[i] = append(ss[i], "月德贵人")
		}
	}
	//金匮贵人：年/日干见它支(甲辰 乙巳 丙未 丁申 戊未 己申 庚戌 辛亥 壬子 癸丑)
	for i, zhi := range zhiA {
		if util.Contains([]string{"甲辰", "乙巳", "丙未", "丁申", "戊未", "己申", "庚戌", "辛亥", "壬子", "癸丑"}, zhiY+zhi, zhiD+zhi) {
			ss[i] = append(ss[i], "金匮贵人")
		}
	}

	//疾厄(羊刃 飞刃 血刃 流霞 灾煞 勾绞 破煞)
	//羊刃: 日干见 阳干加帝旺,阴干之冠带 (甲卯 丙午 戊午 庚酉 壬子 乙辰 丁未 己未 辛戌 癸丑)
	//阳刃 (甲卯 丙午 戊午 庚酉 壬子) 阴刃 (乙辰 丁未 己未 辛戌 癸丑)
	for i, zhi := range zhiA {
		if util.Contains([]string{"甲卯", "丙午", "戊午", "庚酉", "壬子", "乙辰", "丁未", "己未", "辛戌", "癸丑"}, ganD+zhi) {
			ss[i] = append(ss[i], "羊刃")
		}
	}
	//飞刃: 日干见 羊刃冲位 日干为主 四支见之者为是 月时两柱最重 (甲见酉 丙见子 戊见子 庚见卯 壬见午 乙见申 丁见丑 己见丑 辛见辰 癸见未)
	for i, zhi := range zhiA {
		if util.Contains([]string{"甲酉", "丙子", "戊子", "庚卯", "壬午", "乙申", "丁丑", "己丑", "辛辰", "癸未"}, ganD+zhi) {
			ss[i] = append(ss[i], "飞刃")
		}
	}
	//血刃: 月支查他支 (子见午 丑见子 寅见丑 卯见未 辰见寅 巳见申 午见卯 未见酉 申见辰 酉见戌 戌见巳 亥见亥)
	for i, zhi := range zhiA {
		if util.Contains([]string{zhiM + zhi}, "子午", "丑子", "寅丑", "卯未", "辰寅", "巳申", "午卯", "未酉", "申辰", "酉戌", "戌巳", "亥亥") {
			ss[i] = append(ss[i], "血刃")
		}
	}
	//流霞: 日干见 (甲见酉,乙见戌,丙见未,丁见申,戊见巳,己见午,庚见辰,辛见卯,壬见亥,癸见寅)
	for i, zhi := range zhiA {
		if util.Contains([]string{"甲酉", "乙戌", "丙未", "丁申", "戊巳", "己午", "庚辰", "辛卯", "壬亥", "癸寅"}, ganD+zhi) {
			ss[i] = append(ss[i], "流霞")
		}
	}
	//灾煞: 年/日支 三合化五行之胎地,三合仲位之冲位 (申子辰见午 寅午戌见子 巳酉丑见卯 亥卯未见酉)
	for i, zhi := range zhiA {
		if util.Contains([]string{zhiY + zhi, zhiD + zhi},
			"申午", "子午", "辰午", "寅子", "午子", "戌子", "巳卯", "酉卯", "丑卯", "亥酉", "卯酉", "未酉") {
			ss[i] = append(ss[i], "灾煞")
		}
	}
	//勾绞/破煞: 年柱阴男阳女,命前三辰为绞(子酉,丑戌,...),命后三辰为勾(子卯,丑辰,...) 年柱阳男阴女,命前三辰为勾,命后三辰为绞

	//元辰/大耗: 年柱阴男阳女,年支六冲后一支(子巳 丑午 寅未 ...) 年柱阳男阴女,年支六冲前一支(子未 丑申 ...)

	//事业(魁罡 十灵 天医 将星 六秀 学堂 华盖 学士)
	//魁罡: 日柱见者为是. 辰为天罡, 戌为河魁 (壬辰 庚戌 庚辰 戊戌)
	if util.Contains([]string{"壬辰", "庚戌", "庚辰", "戊戌"}, zhuD) {
		ss[2] = append(ss[2], "魁罡")
	}
	//十灵: 日柱(甲辰、乙亥、丙辰、丁酉、戊午、庚戌、庚寅、辛亥、壬寅、癸未)
	if util.Contains([]string{"甲辰", "乙亥", "丙辰", "丁酉", "戊午", "庚戌", "庚寅", "辛亥", "壬寅", "癸未"}, zhuD) {
		ss[2] = append(ss[2], "十灵")
	}
	//天医: 月支见它支 (寅见丑　卯见寅　辰见卯　巳见辰　午见巳　未见午　申见未　酉见申　戌见酉　亥见戌　子见亥　丑见子)
	for i, zhi := range zhiA {
		if util.Contains([]string{zhiM + zhi}, "寅丑", "卯寅", "辰卯", "巳辰", "午巳", "未午", "申未", "酉申", "戌酉", "亥戌", "子亥", "丑子") {
			ss[i] = append(ss[i], "天医")
		}
	}
	//将星: 以年日支查 三合局有仲位(子卯午酉)者为将星 (寅午戌见午,巳酉丑见酉,申子辰见子,亥卯未见卯)
	for i, zhi := range zhiA {
		if util.Contains(zhiYD, "寅", "午", "戌") && zhi == "午" ||
			util.Contains(zhiYD, "巳", "酉", "丑") && zhi == "酉" ||
			util.Contains(zhiYD, "申", "子", "辰") && zhi == "子" ||
			util.Contains(zhiYD, "亥", "卯", "未") && zhi == "卯" {
			ss[i] = append(ss[i], "将星")
		}
	}
	//六秀: 日柱(丙午、丁未、戊子、戊午、己丑、己未)
	if util.Contains([]string{"丙午", "丁未", "戊子", "戊午", "己丑", "己未"}, zhuD) {
		ss[2] = append(ss[2], "六秀")
	}
	//学堂: 年柱纳音五行 四柱见长生,土见甲(木见亥 火见寅 金见巳 水见申,土见甲)
	//学堂: 查法2以日干求长生，更兼天乙、禄马、德秀之神，为日干之财官印食者，皆贤而富贵
	for i, zhi := range zhiA {
		if util.Contains([]string{xingY + zhi}, "木亥", "火寅", "金巳", "水申", "土甲") {
			ss[i] = append(ss[i], "学堂")
		}
	}
	//学士:

	//华盖: 年/日支 三合见墓支(寅午戌见戌 巳酉丑见丑 申子辰见辰 亥卯未见未) 辰戌丑未太多重复出现亦作华盖(如月支辰时支辰)
	for i, zhi := range zhiA {
		if util.Contains([]string{zhiY + zhi, zhiD + zhi},
			"寅戌", "午戌", "戌戌", "巳丑", "酉丑", "丑丑", "申辰", "子辰", "辰辰", "亥未", "卯未", "未未") {
			ss[i] = append(ss[i], "华盖")
		}
	}
	//金匮: 年/日支 三合见长生支(寅午戌见申 巳酉丑见亥 申子辰见巳 亥卯未见寅)?
	for i, zhi := range zhiA {
		if util.Contains([]string{zhiY + zhi, zhiD + zhi},
			"寅申", "午申", "戌申", "巳亥", "酉亥", "丑亥", "申巳", "子巳", "辰巳", "亥寅", "卯寅", "未寅") {
			ss[i] = append(ss[i], "金匮")
		}
	}
	//禄神: 日干见它支(甲未 乙辰 丙巳 丁酉 戊戌 己卯 庚丑 辛申 壬寅 癸午)
	for i, zhi := range zhiA {
		if util.Contains([]string{"甲未", "乙辰", "丙巳", "丁酉", "戊戌", "己卯", "庚丑", "辛申", "壬寅", "癸午"}, zhiD+zhi) {
			ss[i] = append(ss[i], "禄神")
		}
	}
	//暗禄: 日干见它支(甲亥 乙戌 丙申 丁未 戊申 己未 庚巳 辛辰 癸丑)
	for i, zhi := range zhiA {
		if util.Contains([]string{"甲亥", "乙戌", "丙申", "丁未", "戊申", "己未", "庚巳", "辛辰", "癸丑"}, zhiD+zhi) {
			ss[i] = append(ss[i], "暗禄")
		}
	}
	//拱禄: 日柱+时柱(癸亥+癸丑 癸丑+癸亥 丁巳+丁午 丁午+丁巳 己巳+己午 己午+己巳 戊辰+戊午 戊午+戊辰)
	if util.Contains([]string{"癸亥癸丑", "癸丑癸亥", "丁巳丁午", "丁午丁巳", "己巳己午", "己午己巳", "戊辰戊午", "戊午戊辰"}, zhuD+zhuT) {
		ss[2] = append(ss[2], "拱禄")
	}
	//夹禄: 日柱同见禄神前后支(甲见丑卯 乙见寅辰 丙见辰午 丁见巳未 戊见辰午 己见巳未 庚见未酉 辛见申戌 壬见戌子 癸见亥丑)
	for i, zhi := range zhiA {
		if (ganD == "甲" && util.Contains([]string{zhi}, "丑", "卯") && util.Contains(zhiA, "丑") && util.Contains(zhiA, "卯")) ||
			(ganD == "乙" && util.Contains([]string{zhi}, "寅", "辰") && util.Contains(zhiA, "寅") && util.Contains(zhiA, "辰")) ||
			(ganD == "丙" && util.Contains([]string{zhi}, "辰", "午") && util.Contains(zhiA, "辰") && util.Contains(zhiA, "午")) ||
			(ganD == "丁" && util.Contains([]string{zhi}, "巳", "未") && util.Contains(zhiA, "巳") && util.Contains(zhiA, "未")) ||
			(ganD == "戊" && util.Contains([]string{zhi}, "辰", "午") && util.Contains(zhiA, "辰") && util.Contains(zhiA, "午")) ||
			(ganD == "己" && util.Contains([]string{zhi}, "巳", "未") && util.Contains(zhiA, "巳") && util.Contains(zhiA, "未")) ||
			(ganD == "庚" && util.Contains([]string{zhi}, "未", "酉") && util.Contains(zhiA, "未") && util.Contains(zhiA, "酉")) ||
			(ganD == "辛" && util.Contains([]string{zhi}, "申", "戌") && util.Contains(zhiA, "申") && util.Contains(zhiA, "戌")) ||
			(ganD == "壬" && util.Contains([]string{zhi}, "戌", "子") && util.Contains(zhiA, "戌") && util.Contains(zhiA, "子")) ||
			(ganD == "癸" && util.Contains([]string{zhi}, "亥", "丑") && util.Contains(zhiA, "亥") && util.Contains(zhiA, "丑")) {
			ss[i] = append(ss[i], "夹禄")
		}
	}
	//金羊禄: 日支查他支 (甲辰 乙巳 丙未 丁申 戊未 己申 庚戌 辛亥 壬丑 癸寅)
	for i, zhi := range zhiA {
		if util.Contains([]string{"甲辰", "乙巳", "丙未", "丁申", "戊未", "己申", "庚戌", "辛亥", "壬丑", "癸寅"}, zhiD+zhi) {
			ss[i] = append(ss[i], "金羊禄")
		}
	}
	//驿马: 年/日支 三合化五行之病地(申子辰见寅 寅午戌见申 巳酉丑见亥 亥卯未见巳)
	for i, zhi := range zhiA {
		if util.Contains([]string{zhiY + zhi, zhiD + zhi},
			"申寅", "子寅", "辰寅", "寅申", "午申", "戌申", "巳亥", "酉亥", "丑亥", "亥巳", "卯巳", "未巳") {
			ss[i] = append(ss[i], "驿马")
		}
	}
	//日德: 日柱(甲寅 丙辰 戊辰 庚辰 壬戌)
	if util.Contains([]string{"甲寅", "丙辰", "戊辰", "庚辰", "壬戌"}, zhuD) {
		ss[2] = append(ss[2], "日德")
	}
	//八专: 日柱干支同气且坐禄/冠带(甲寅 乙卯 己未 丁未 戊戌 庚申 辛酉 癸丑) 时柱也
	if util.Contains([]string{"甲寅", "乙卯", "己未", "丁未", "戊戌", "庚申", "辛酉", "癸丑"}, zhuD) {
		ss[2] = append(ss[2], "八专")
	}
	//九丑: 日柱(壬子、壬午、戊子、戊午、乙卯、己卯、辛卯、辛酉、己酉)
	if util.Contains([]string{"壬子", "壬午", "戊子", "戊午", "乙卯", "己卯", "辛卯", "辛酉", "己酉"}, zhuD) {
		ss[2] = append(ss[2], "九丑")
	}
	//四废: 月支三会见官官杀杀 春 寅卯辰月 见 庚申/辛酉日 夏 巳午未月 见 壬子/癸亥日 秋 申酉戌月 见 甲寅/乙卯日 冬 亥子丑月 见 丙午/丁巳日
	if (util.Contains([]string{"寅", "卯", "辰"}, zhiM) && util.Contains([]string{"庚申", "辛酉"}, zhiD)) ||
		(util.Contains([]string{"巳", "午", "未"}, zhiM) && util.Contains([]string{"壬子", "癸亥"}, zhiD)) ||
		(util.Contains([]string{"申", "酉", "戌"}, zhiM) && util.Contains([]string{"甲寅", "乙卯"}, zhiD)) ||
		(util.Contains([]string{"亥", "子", "丑"}, zhiM) && util.Contains([]string{"丙午", "丁巳"}, zhiD)) {
		ss[2] = append(ss[2], "四废")
	}
	//十恶大败: 日柱(甲辰、乙巳、丙申、丁亥、戊戌、己丑、庚辰、辛巳、壬申、癸亥)
	if util.Contains([]string{"甲辰", "乙巳", "丙申", "丁亥", "戊戌", "己丑", "庚辰", "辛巳", "壬申", "癸亥"}, zhuD) {
		ss[2] = append(ss[2], "十恶大败")
	}
	//亡神: 年/日支见它支 三合见仲位前一位(申子辰见亥 寅午戌见巳 巳酉丑见申 亥卯未见寅)
	for i, zhi := range zhiA {
		if (util.Contains([]string{zhi}, "亥", "子", "辰") && util.Contains(ganYD, "申")) ||
			(util.Contains([]string{zhi}, "申", "子", "辰") && util.Contains(ganYD, "巳")) ||
			(util.Contains([]string{zhi}, "寅", "午", "戌") && util.Contains(ganYD, "亥")) {
			ss[i] = append(ss[i], "亡神")
		}
	}
	//指背煞: 年支查他支 三合见长生(寅午戌见寅 巳酉丑见巳 申子辰见申 亥卯未见亥)
	for i, zhi := range zhiA {
		if util.Contains([]string{zhiY + zhi},
			"寅寅", "午寅", "戌寅", "巳巳", "酉巳", "丑巳", "申申", "子申", "辰申", "亥亥", "卯亥", "未亥") {
			ss[i] = append(ss[i], "指背煞")
		}
	}

	//姻缘(天喜 红鸾 咸池/桃花 红艳 童子煞 孤鸾 孤辰 寡宿 阴差阳错)
	//天喜: 年支见它支 (子见酉 丑见申 寅见未 卯见午 辰见巳 巳见辰 午见卯 未见寅 申见丑 酉见子 戌见亥 亥见戌)
	for i, zhi := range zhiA {
		if util.Contains([]string{zhiY + zhi}, "子酉", "丑申", "寅未", "卯午", "辰巳", "巳辰", "午卯", "未寅", "申丑", "酉子", "戌亥", "亥戌") {
			ss[i] = append(ss[i], "天喜")
		}
	}
	//红鸾: 年支见它支(子卯 丑寅 寅丑 卯子 辰亥 巳戌 午酉 未申 申未 酉午 戌巳 亥辰)
	for i, zhi := range zhiA {
		if util.Contains([]string{zhiY + zhi}, "子卯", "丑寅", "寅丑", "卯子", "辰亥", "巳戌", "午酉", "未申", "申未", "酉午", "戌巳", "亥辰") {
			ss[i] = append(ss[i], "红鸾")
		}
	}
	//咸池/桃花: 年/日支为主,见三合中神之沐浴(申子辰见酉 寅午戌见卯 巳酉丑见午 亥卯未见子)
	for i, zhi := range zhiA {
		if util.Contains([]string{zhiY + zhi, zhiD + zhi},
			"申酉", "子酉", "辰酉", "寅卯", "午卯", "戌卯", "巳午", "酉午", "丑午", "亥子", "卯子", "未子") {
			ss[i] = append(ss[i], "桃花")
		}
	}
	//红艳: 日干见它支(甲午 乙申 丙寅 丁未 戊午 己辰 庚酉 辛戌 壬子 癸申)
	for i, zhi := range zhiA {
		if util.Contains([]string{"甲午", "乙申", "丙寅", "丁未", "戊午", "己辰", "庚酉", "辛戌", "壬子", "癸申"}, zhiD+zhi) {
			ss[i] = append(ss[i], "红艳")
		}
	}
	//童子煞: 月支见日时支 子丑巳午未亥 见卯/未/辰 寅卯辰申酉 见寅/子
	//童子煞: 年柱纳音五行(金/木见 日/时支卯午,水/火见日/时支酉戌,土见日/时支辰/巳)
	for i, zhi := range []string{zhiD, zhiT} {
		if (util.Contains([]string{zhiM}, "子", "丑", "巳", "午", "未", "亥") && util.Contains([]string{zhi}, "卯", "未", "辰")) ||
			(util.Contains([]string{zhiM}, "寅", "卯", "辰", "申", "酉") && util.Contains([]string{zhi}, "寅", "子")) ||
			(util.Contains([]string{xingY}, "金", "木") && util.Contains([]string{zhi}, "卯", "午")) ||
			(util.Contains([]string{xingY}, "水", "火") && util.Contains([]string{zhi}, "酉", "戌")) ||
			(util.Contains([]string{xingY}, "土") && util.Contains([]string{zhi}, "辰", "巳")) {
			ss[i] = append(ss[i], "童子煞")
		}
	}
	//孤鸾: 日柱(乙巳、丁巳、辛亥、戊申、甲寅、戊午、壬子、丙午)
	if util.Contains([]string{"乙巳", "丁巳", "辛亥", "戊申", "甲寅", "戊午", "壬子", "丙午"}, zhuD) {
		ss[2] = append(ss[2], "孤鸾")
	}
	//孤辰: 年支见它支 三会五行见衰  (亥子丑见寅 寅卯辰见巳 巳午未见申 申酉戌见亥)
	for i, zhi := range zhiA {
		if util.Contains([]string{zhiY + zhi},
			"亥寅", "子寅", "丑寅", "寅巳", "卯巳", "辰巳", "巳申", "午申", "未申", "申亥", "酉亥", "戌亥") {
			ss[i] = append(ss[i], "孤辰")
		}
	}
	//寡宿: 年支见它支 三会五行见冠带(亥子丑见戌 寅卯辰见丑 巳午未见辰 申酉戌见未)
	for i, zhi := range zhiA {
		if util.Contains([]string{zhiY + zhi},
			"亥戌", "子戌", "丑戌", "寅丑", "卯丑", "辰丑", "巳辰", "午辰", "未辰", "申未", "酉未", "戌未") {
			ss[i] = append(ss[i], "寡宿")
		}
	}
	//阴差阳错: 日柱(丙子 丙午 丁丑 丁未 戊寅 戊申 辛卯 辛酉 壬辰 壬戌 癸巳 癸亥)
	if util.Contains([]string{"丙子", "丙午", "丁丑", "丁未", "戊寅", "戊申", "辛卯", "辛酉", "壬辰", "壬戌", "癸巳", "癸亥"}, zhuD) {
		ss[2] = append(ss[2], "阴差阳错")
	}
	//披麻: 年支见它支后三辰(子酉 丑戌 寅亥 卯子 辰丑 巳寅 午卯 未辰 申巳 酉午 戌未 亥申)
	for i, zhi := range zhiA {
		if util.Contains([]string{zhiY + zhi},
			"子酉", "丑戌", "寅亥", "卯子", "辰丑", "巳寅", "午卯", "未辰", "申巳", "酉午", "戌未", "亥申") {
			ss[i] = append(ss[i], "披麻")
		}
	}
	return ss[0], ss[1], ss[2], ss[3]
}
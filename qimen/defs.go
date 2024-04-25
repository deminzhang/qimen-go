package qimen

// Yuan3Name 奇门三元名
var Yuan3Name = []string{"", "上元", "中元", "下元"}

// 节气索引
var _JieQiIndex = map[string]int{
	"冬至": 1, "小寒": 2, "大寒": 3,
	"立春": 4, "雨水": 5, "惊蛰": 6,
	"春分": 7, "清明": 8, "谷雨": 9,
	"立夏": 10, "小满": 11, "芒种": 12,
	"夏至": 13, "小暑": 14, "大暑": 15,
	"立秋": 16, "处暑": 17, "白露": 18,
	"秋分": 19, "寒露": 20, "霜降": 21,
	"立冬": 22, "小雪": 23, "大雪": 24,
}

// 奇门局数
var _QiMenJu = [][]int{{0, 0, 0},
	{+1, +7, +4}, {+2, +8, +5}, {+3, +9, +6}, //坎宫 冬至{上元,中元,下元},小寒{上元,中元,下元},大寒{上元,中元,下元},
	{+8, +5, +2}, {+9, +6, +3}, {+1, +7, +4}, //艮宫 立春...
	{+3, +9, +6}, {+4, +1, +7}, {+5, +2, +8}, //震...
	{+4, +1, +7}, {+5, +2, +8}, {+6, +3, +9}, //巽
	{-9, -3, -6}, {-8, -2, -5}, {-7, -1, -4}, //离
	{-2, -5, -8}, {-1, -4, -7}, {-9, -3, -6}, //坤
	{-7, -1, -4}, {-6, -9, -3}, {-5, -8, -2}, //兑
	{-6, -9, -3}, {-5, -8, -2}, {-4, -7, -1}, //乾
	//{1, 7, 4},
}

// QMType 盘式
var QMType = []string{"转盘", "飞盘", "鸣法"}

const (
	QMTypeRotating = 0
	QMTypeFly      = 1
	QMTypeAmaze    = 2
)

// QMFlyType 飞盘九星飞宫
var QMFlyType = []string{"阴阳皆顺", "阳顺阴逆"}

const (
	QMFlyTypeAllOrder     = 0 // 阴阳皆顺,鸣法同
	QMFlyTypeLunarReverse = 1 // 阴阳皆逆:源于括囊集
)

// QMHostingType 转盘寄宫法
var QMHostingType = []string{"中宫寄坤", "阳艮阴坤", "_土寄四维"}

const (
	QMHostingType2    = 0
	QMHostingType28   = 1
	QMHostingType2846 = 2
)

// 起局方式
var _QMStartType = []string{"拆补", "茅山", "置闰", "自选"}

// 暗干起法
var _QMHideGanType = []string{"值使门起", "门地盘起"}

// Idx8 序环
var Idx8 = []int{8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}

// Idx9 序环
var Idx9 = []int{9, 1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9}

// Idx12 序环
var Idx12 = []int{12, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

// 旬首遁甲
var HideJia = map[string]string{
	"甲子": "戊",
	"甲戌": "己",
	"甲申": "庚",
	"甲午": "辛",
	"甲辰": "壬",
	"甲寅": "癸",
}

const (
	Trunk10      = "_甲乙丙丁戊已庚辛壬癸"   //天干
	Branch12     = "_子丑寅卯辰巳午末申酉戌亥" //地支
	Diagrams8In9 = "_坎坤震巽中乾竞艮离"    //九宫八卦
	//Term24    = "__小寒大寒立春雨水惊蛰春分清明谷雨立夏小满芒种夏至小暑大暑立秋处暑白露秋分寒露霜降立冬小雪大雪冬至"

	Star0 = "天"
	Star9 = " 蓬芮冲辅禽心柱任英" //奇门九星
	Star8 = " 蓬任冲辅英芮柱心"  //转盘用九星

	Door0 = "门"
	Door8 = "_休生伤杜景死惊开"  //转盘用八门
	Door9 = "_休死伤杜中开惊生景" //飞盘用九门

	T3Qi6Yi = "_戊己庚辛壬癸丁丙乙" //三奇六仪

	God9S      = "__值符腾蛇太阴六合勾陈太常朱雀九地九天" //九神飞盘阳遁用
	God9L      = "__值符腾蛇太阴六合白虎太常玄武九地九天" //九神飞盘阴遁用
	God8       = "__值符腾蛇太阴六合白虎玄武九地九天"   //八神转盘用
	MonthJiang = "_亥戌酉申未午巳辰卯寅丑子"        //月将正月起亥
	MonthJian  = "_寅卯辰巳午未申酉戌亥子丑"        //月建正月起寅
	Build12    = "_建除满平定执破危成收开闭"        //十二建星

	QMDayStar9 = "__太乙摄提轩辕招摇天符青龙咸池太阴天乙"       //日家奇门九星
	God12      = "__青龙明堂天刑朱雀金匮天德白虎玉堂天牢玄武司命勾陈" //日家奇门十二原神黄黑道
	God12YB    = "_黄黄黑黑黄黄黑黄黑黑黄黑"              //十二黄黑道
)

func Diagrams9(i int) string {
	i = Idx9[i]
	return string([]rune(Diagrams8In9)[i : i+1])
}
func QMStar9(i int) string {
	i = Idx9[i]
	return Star0 + string([]rune(Star9)[i:i+1])
}
func QMStar8(i int) string {
	i = Idx8[i]
	return Star0 + string([]rune(Star8)[i:i+1])
}
func QM3Qi6Yi(i int) string {
	i = Idx9[i]
	return string([]rune(T3Qi6Yi)[i : i+1])
}
func QMDoor8(i int) string {
	i = Idx8[i]
	return string([]rune(Door8)[i:i+1]) + Door0
}
func QMDoor9(i int) string {
	i = Idx9[i]
	return string([]rune(Door9)[i:i+1]) + Door0
}
func QMGod9S(i int) string {
	i = Idx9[i] * 2
	return string([]rune(God9S)[i : i+2])
}
func QMGod9L(i int) string {
	i = Idx9[i] * 2
	return string([]rune(God9L)[i : i+2])
}
func QMGod8(i int) string {
	i = Idx8[i] * 2
	return string([]rune(God8)[i : i+2])
}
func YueJiang(i int) string {
	i = Idx12[i]
	return string([]rune(MonthJiang)[i : i+1])
}

// 奇门转盘用转宫宫位索引
var _QMRollIdx = []int{6, 1, 8, 3, 4, 9, 2, 7, 6}     //转宫号=>洛宫号
var _QM2RollIdx = []int{1, 1, 6, 3, 4, 0, 8, 7, 2, 5} //洛宫号=>转宫号

// YueJiangName 月将神名
var YueJiangName = map[string]string{
	"亥": "登明", "戌": "河魅", "酉": "从魁", "申": "传送",
	"未": "小吉", "午": "胜光", "巳": "太乙", "辰": "天罡",
	"卯": "太冲", "寅": "功曹", "丑": "大吉", "子": "神后",
}

// Horse 驿马方(申子辰见寅 寅午戌见申 巳酉丑见亥 亥卯未见巳)
var Horse = map[string]string{
	"申": "寅", "子": "寅", "辰": "寅",
	"寅": "申", "午": "申", "戌": "申",
	"巳": "亥", "酉": "亥", "丑": "亥",
	"亥": "巳", "卯": "巳", "未": "巳",
}

// 农历日期信息
// 阴历1900年到2100年每年中的月天数信息
// 阴历每月只能是29或30天，一年用12（或13）个二进制位表示，对应位为1 代表30天，否则为29天
// 闰月不会出现在正月、冬月或腊月,不会连续两年闰月
//var lunarMonthData = [201]int{
//	// 0xf   =0000 0000 0000 1111
//	//       =1234 5678 9ABC 1000闰月月数8
//	//       =1234 5678 9ABC 去年闰月大小 0小 1111大
//	0x4bd8, //1900:0100 1011 1101 1000(小大小小 大小大大 大大小大 闰八月)
//	0x4ae0, //1901:0100 1010 1110 0000(去年闰八月小)
//	0xa570,
//	0x54d5, //1903:0101 0100 1101 0101(闰五月)
//	0xd260, //1904:1101 0010 0110 0000(去年闰五月小)
//	0xd950,
//	0x5554, //1906:0101 0101 0101 0100(闰四月)
//	0x56af, //1907:0101 0110 1010 1111(去年闰四月大)
//	0x9ad0, 0x55d2,
//	0x4ae0, 0xa5b6, 0xa4d0, 0xd250, 0xd255, 0xb54f, 0xd6a0, 0xada2, 0x95b0, 0x4977,
//	0x497f, 0xa4b0, 0xb4b5, 0x6a50, 0x6d40, 0xab54, 0x2b6f, 0x9570, 0x52f2, 0x4970,
//	0x6566, 0xd4a0, 0xea50, 0x6a95, 0x5adf, 0x2b60, 0x86e3, 0x92ef, 0xc8d7, 0xc95f,
//	0xd4a0, 0xd8a6, 0xb55f, 0x56a0, 0xa5b4, 0x25df, 0x92d0, 0xd2b2, 0xa950, 0xb557,
//	0x6ca0, 0xb550, 0x5355, 0x4daf, 0xa5b0, 0x4573, 0x52bf, 0xa9a8, 0xe950, 0x6aa0,
//	0xaea6, 0xab50, 0x4b60, 0xaae4, 0xa570, 0x5260, 0xf263, 0xd950, 0x5b57, 0x56a0,
//	0x96d0, 0x4dd5, 0x4ad0, 0xa4d0, 0xd4d4, 0xd250, 0xd558, 0xb540, 0xb6a0, 0x95a6,
//	0x95bf, 0x49b0, 0xa974, 0xa4b0, 0xb27a, 0x6a50, 0x6d40, 0xaf46, 0xab60, 0x9570,
//	0x4af5, 0x4970, 0x64b0, 0x74a3, 0xea50, 0x6b58, 0x5ac0, 0xab60, 0x96d5, 0x92e0,
//	0xc960, 0xd954, 0xd4a0, 0xda50, 0x7552, 0x56a0, 0xabb7, 0x25d0, 0x92d0, 0xcab5,
//	0xa950, 0xb4a0, 0xbaa4, 0xad50, 0x55d9, 0x4ba0, 0xa5b0, 0x5176, 0x52bf, 0xa930,
//	0x7954, 0x6aa0, 0xad50, 0x5b52, 0x4b60, 0xa6e6, 0xa4e0, 0xd260, 0xea65, 0xd530,
//	0x5aa0, 0x76a3, 0x96d0, 0x4afb, 0x4ad0, 0xa4d0, 0xd0b6, 0xd25f, 0xd520, 0xdd45,
//	0xb5a0, 0x56d0, 0x55b2, 0x49b0, 0xa577, 0xa4b0, 0xaa50, 0xb255, 0x6d2f, 0xada0,
//	0x4b63, 0x937f, 0x49f8, 0x4970, 0x64b0, 0x68a6, 0xea5f, 0x6b20, 0xa6c4, 0xaaef,
//	0x92e0, 0xd2e3, 0xc960, 0xd557, 0xd4a0, 0xda50, 0x5d55, 0x56a0, 0xa6d0, 0x55d4,
//	0x52d0, 0xa9b8, 0xa950, 0xb4a0, 0xb6a6, 0xad50, 0x55a0, 0xaba4, 0xa5b0, 0x52b0,
//	0xb273, 0x6930, 0x7337, 0x6aa0, 0xad50, 0x4b55, 0x4b6f, 0xa570, 0x54e4, 0xd260,
//	0xe968, 0xd520, 0xdaa0, 0x6aa6, 0x56df, 0x4ae0, 0xa9d4, 0xa4d0, 0xd150, 0xf252,
//	0xd520, //2100:1101 0101 0010 0000
//}

// 廿四节气信息 (0小寒)
// 从第0个节气的分钟数
var termData = []int{
	0, 21208, 42467, 63836, 85337, 107014, 128867, 150921, 173149, 195551, 218072, 240693,
	263343, 285989, 308563, 331033, 353350, 375494, 397447, 419210, 440795, 462224, 483532,
	504758,
}

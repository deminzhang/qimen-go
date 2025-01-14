package qimen

import (
	"fmt"
	"github.com/6tail/lunar-go/LunarUtil"
	"strings"
)

const (
	QMGameHour  = 0 //时家奇门
	QMGameDay   = 1 //日家奇门
	QMGameMonth = 2 //月家奇门
	QMGameYear  = 3 //年家奇门
	QMGameDay2  = 4 //日家奇门太乙
)

// Yuan3Name 奇门三元名
var Yuan3Name = []string{"", "上元", "中元", "下元"}

// 节气索引
var _JieQiIndex = map[string]int{
	"冬至": 1, "小寒": 2, "大寒": 3, "立春": 4,
	"雨水": 5, "惊蛰": 6, "春分": 7, "清明": 8,
	"谷雨": 9, "立夏": 10, "小满": 11, "芒种": 12,
	"夏至": 13, "小暑": 14, "大暑": 15, "立秋": 16,
	"处暑": 17, "白露": 18, "秋分": 19, "寒露": 20,
	"霜降": 21, "立冬": 22, "小雪": 23, "大雪": 24,
}

// 奇门局数 大雪0 冬至1 小寒2 ..大雪24
var _QiMenJu = [][]int{{-4, -7, -1}, //大雪0 24
	{+1, +7, +4}, {+2, +8, +5}, {+3, +9, +6}, //坎宫 冬至{上元,中元,下元},小寒{上元,中元,下元},大寒{上元,中元,下元},
	{+8, +5, +2}, {+9, +6, +3}, {+1, +7, +4}, //艮宫 立春{上元,中元,下元},雨水{上元,中元,下元},惊蛰{上元,中元,下元},
	{+3, +9, +6}, {+4, +1, +7}, {+5, +2, +8}, //震宫 春分{上元,中元,下元},清明{上元,中元,下元},谷雨{上元,中元,下元},
	{+4, +1, +7}, {+5, +2, +8}, {+6, +3, +9}, //巽宫 立夏{上元,中元,下元},小满{上元,中元,下元},芒种{上元,中元,下元},
	{-9, -3, -6}, {-8, -2, -5}, {-7, -1, -4}, //离宫 夏至{上元,中元,下元},小暑{上元,中元,下元},大暑{上元,中元,下元},
	{-2, -5, -8}, {-1, -4, -7}, {-9, -3, -6}, //坤宫 立秋{上元,中元,下元},处暑{上元,中元,下元},白露{上元,中元,下元},
	{-7, -1, -4}, {-6, -9, -3}, {-5, -8, -2}, //兑宫 秋分{上元,中元,下元},寒露{上元,中元,下元},霜降{上元,中元,下元},
	{-6, -9, -3}, {-5, -8, -2}, {-4, -7, -1}, //乾宫 立冬{上元,中元,下元},小雪{上元,中元,下元},大雪{上元,中元,下元},
	{1, 7, 4}, //冬至{上元,中元,下元},
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
var QMHostingType = []string{"中宫寄坤", "阳艮阴坤", "_寄四维", "_寄八节"}

const (
	QMHostingType2  = 0
	QMHostingType28 = 1
)

// QMJuType 起局方式
var QMJuType = []string{"拆补", "茅山", "置闰", "自选"}

const (
	QMJuTypeSplit   = 0 //节气和日干符头定三元
	QMJuTypeMaoShan = 1 //无视符头，节气开始上元60时辰,中元60时辰,再下元60时辰,下元满60时辰不到下个节气继用下元
	QMJuTypeZhiRun  = 2 //符头和节气的关系
	QMJuTypeSelf    = 3 //自选
)

// QMHideGanType 暗干起法
var QMHideGanType = []string{"暗干值使起", "门地暗干"}

const (
	QMHideGanDutyDoorHour = 0 //值使门起 值使落宫起时干 地盘干与时干相同时,时干入中宫
	QMHideGanDoorHomeGan  = 1 //门地盘起 八门带原始宫的地盘干
)

// Idx8 序环
var Idx8 = []int{8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}

// Idx9 序环
var Idx9 = []int{9, 1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9}

// Idx12 序环
var Idx12 = []int{12, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

// Idx10 序环
var Idx10 = []int{10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// HideJia 旬首遁甲
var HideJia = map[string]string{
	"甲子": "戊", "甲戌": "己", "甲申": "庚",
	"甲午": "辛", "甲辰": "壬", "甲寅": "癸",
}

var SkyGate3 = map[string]bool{
	"太冲": true, "卯": true,
	"小吉": true, "未": true,
	"从魁": true, "酉": true,
}
var GroundGate4 = map[string]bool{
	"除": true, "危": true, "定": true, "开": true,
}

const (
	Diagrams8In9 = "_坎坤震巽中乾兑艮离" //九宫八卦
	//Term24    = "__小寒大寒立春雨水惊蛰春分清明谷雨立夏小满芒种夏至小暑大暑立秋处暑白露秋分寒露霜降立冬小雪大雪冬至"

	Star0 = "天"
	Star9 = "_蓬芮冲辅禽心柱任英" //奇门九星
	Star8 = "_蓬任冲辅英芮柱心"  //转盘用九星

	Door0 = "门"
	Door8 = "_休生伤杜景死惊开"  //转盘用八门
	Door9 = "_休死伤杜中开惊生景" //飞盘用九门

	T3Qi6Yi = "_戊己庚辛壬癸丁丙乙" //三奇六仪

	God9S      = "__值符腾蛇太阴六合勾陈太常朱雀九地九天" //九神飞盘阳遁用
	God9L      = "__值符腾蛇太阴六合白虎太常玄武九地九天" //九神飞盘阴遁用
	God8       = "__值符腾蛇太阴六合白虎玄武九地九天"   //八神转盘用
	MonthBuild = "_寅卯辰巳午未申酉戌亥子丑"        //月建 正月起寅 交节换建
	Build12    = "_建除满平定执破危成收开闭"        //十二建星
	MonthJiang = "_亥戌酉申未午巳辰卯寅丑子"        //月将 正月起亥 交中气换将

	QMDayStar9   = "__太乙摄提轩辕招摇天符青龙咸池太阴天乙"       //日家奇门2九星
	QMDayGod12   = "__青龙明堂天刑朱雀金匮天德白虎玉堂天牢玄武司命勾陈" //日家奇门2十二原神黄黑道
	QMDayGod12YB = "_黄黄黑黑黄黄黑黄黑黑黄黑"              //十二黄黑道
)

// 先天八卦
var Diagrams8Origin = map[uint8]string{
	1: "乾", 2: "兑", 3: "离", 4: "震", 5: "巽", 6: "坎", 7: "艮", 8: "坤",
}
var Diagrams8IdxOrigin = map[string]uint8{
	"乾": 1, "兑": 2, "离": 3, "震": 4, "巽": 5, "坎": 6, "艮": 7, "坤": 8,
}

// 后天八卦
var Diagrams8 = map[uint8]string{
	1: "坎", 2: "坤", 3: "震", 4: "巽", 5: "  ", 6: "乾", 7: "兑", 8: "艮", 9: "离",
}
var Diagrams8Idx = map[string]uint8{
	"坎": 1, "坤": 2, "震": 3, "巽": 4, "乾": 6, "兑": 7, "艮": 8, "离": 9,
}

// 八卦爻
var Diagrams8Bin = map[string]uint8{
	"乾": 0b111, "兑": 0b011, "离": 0b101, "震": 0b001, "巽": 0b110, "坎": 0b010, "艮": 0b100, "坤": 0b000,
}
var Diagrams8FromBin = map[uint8]string{
	0b111: "乾", 0b011: "兑", 0b101: "离", 0b001: "震", 0b110: "巽", 0b010: "坎", 0b100: "艮", 0b000: "坤",
}

// 64卦 以先天数索引
var Diagrams64 = map[uint8]string{
	11: "乾", 12: "履", 13: "同人", 14: "无妄", 15: "姤", 16: "讼", 17: "遁", 18: "否",
	21: "夬", 22: "兑", 23: "革", 24: "随", 25: "大过", 26: "困", 27: "咸", 28: "萃",
	31: "大有", 32: "睽", 33: "离", 34: "噬嗑", 35: "鼎", 36: "未济", 37: "旅", 38: "晋",
	41: "大壮", 42: "归妹", 43: "丰", 44: "震", 45: "恒", 46: "解", 47: "小过", 48: "豫",
	51: "小畜", 52: "中孚", 53: "家人", 54: "益", 55: "巽", 56: "涣", 57: "渐", 58: "观",
	61: "需", 62: "节", 63: "既济", 64: "屯", 65: "井", 66: "坎", 67: "蹇", 68: "比",
	71: "大畜", 72: "损", 73: "贲", 74: "颐", 75: "蛊", 76: "蒙", 77: "艮", 78: "剥",
	81: "泰", 82: "临", 83: "明夷", 84: "复", 85: "升", 86: "师", 87: "谦", 88: "坤",
}
var Diagrams64FullName = map[uint8]string{
	11: "乾为天", 12: "天泽履", 13: "天火同人", 14: "天雷无妄", 15: "天风姤", 16: "天水讼", 17: "天山遁", 18: "天地否",
	21: "泽天夬", 22: "兑为泽", 23: "泽火革", 24: "泽雷随", 25: "泽风大过", 26: "泽水困", 27: "泽山咸", 28: "泽地萃",
	31: "火天大有", 32: "火泽睽", 33: "离为火", 34: "火雷噬嗑", 35: "火风鼎", 36: "火水未济", 37: "火山旅", 38: "火地晋",
	41: "雷天大壮", 42: "雷泽归妹", 43: "雷火丰", 44: "震为雷", 45: "雷风恒", 46: "雷水解", 47: "雷山小过", 48: "雷地豫",
	51: "风天小畜", 52: "风泽中孚", 53: "风火家人", 54: "风雷益", 55: "巽为风", 56: "风水涣", 57: "风山渐", 58: "风地观",
	61: "水天需", 62: "水泽节", 63: "水火既济", 64: "水雷屯", 65: "水风井", 66: "坎为水", 67: "水山蹇", 68: "水地比",
	71: "山天大畜", 72: "山泽损", 73: "山火贲", 74: "山雷颐", 75: "山风蛊", 76: "山水蒙", 77: "艮为山", 78: "山地剥",
	81: "地天泰", 82: "地泽临", 83: "地火明夷", 84: "地雷复", 85: "地风升", 86: "地水师", 87: "地山谦", 88: "坤为地",
}

// StarHome 星原始宫位
var StarHome = map[string]int{
	"蓬": 1, "芮": 2, "冲": 3, "辅": 4, "禽": 5, "心": 6, "柱": 7, "任": 8, "英": 9,
}

// DoorHome 门原始宫位
var DoorHome = map[string]int{
	"休": 1, "生": 8, "伤": 3, "杜": 4, "中": 5, "景": 9, "死": 2, "惊": 7, "开": 6,
}
var Gong9Color = []string{"",
	"白", "黑", "青", "碧", "黄", "白", "赤", "白", "紫",
}
var Diagrams8Color = map[string]string{
	"坎": "白", "坤": "黑", "震": "青", "巽": "碧", "中": "黄", "乾": "白", "兑": "赤", "艮": "白", "离": "紫",
}
var DiagramsWuxing = map[string]string{
	"坎": "水", "坤": "土", "震": "木", "巽": "木", "中": "土", "乾": "金", "兑": "金", "艮": "土", "离": "火",
}
var DoorWuxing = map[string]string{
	"休": "水", "生": "土", "伤": "木", "杜": "木", "中": "土", "景": "火", "死": "土", "惊": "金", "开": "金",
}
var WuxingKe = map[string]string{
	"金": "木", "木": "土", "土": "水", "水": "火", "火": "金",
}

// QMTomb 奇门入墓
var QMTomb = map[string]string{
	"甲": "未", "乙": "戌", "丙": "戌", "丁": "丑", "戊": "戌",
	"己": "丑", "庚": "丑", "辛": "辰", "壬": "辰", "癸": "未",
}

// 甲子戊落震三宫：属于无礼之刑，性格暴躁，排斥异己，缺乏礼貌，好色。
// 甲戌己落坤二宫：属于恃势之刑，依仗自己具有某种优势，而猛进或孤注一掷，容易受挫折，做事情有成有败、大起大落，身体多慢性病，孤独短命。
// 甲申庚落艮八宫：属于无恩之刑，性情冷酷、薄情寡义。好心反而招怨、斗殴伤灾等，离婚或女性流产。凡事事与愿违，不得伸展，寸步难行。
// 甲午辛落离九宫：属于自刑，内心阴险，报复心强，爱指手划脚，自私，只可共患难、不能共享福，惟我独尊、鄙视他人，自己气急而残害自己、粉碎自己。郁闷时离家出走。
// 甲辰壬落巽四宫：属于自刑，做事情不能有始有终，缺乏独立能力，自私自利，党营私，游荡不羁，容易离家出走，内心阴险，损害自家门风，多肢节手足之灾。
// 甲寅癸落巽四宫：属于无恩之刑，冷酷无情，常有叛逆的行为，易与人结下冤仇，有刑拘、牢狱之灾，弃妻纳妾，离婚等。
// 六仪击刑
var QM6YiJiXing = map[string]string{
	"甲子": "卯", "甲戌": "未", "甲申": "寅", "甲午": "午", "甲辰": "辰", "甲寅": "巳",
}

// LunarUtil.WU_XING_ZHI
//LunarUtil.WU_XING_GAN

func Diagrams9(i int) string {
	i = Idx9[i]
	return string([]rune(Diagrams8In9)[i : i+1])
}
func QMStar9(i int) string {
	i = Idx9[i]
	//return Star0 + string([]rune(Star9)[i:i+1])
	return string([]rune(Star9)[i : i+1])
}
func QMStar8(i int) string {
	i = Idx8[i]
	//return Star0 + string([]rune(Star8)[i:i+1])
	return string([]rune(Star8)[i : i+1])
}
func QM3Qi6Yi(i int) string {
	i = Idx9[i]
	return string([]rune(T3Qi6Yi)[i : i+1])
}
func QMDoor8(i int) string {
	i = Idx8[i]
	return string([]rune(Door8)[i : i+1]) // + Door0
}
func QMDoor9(i int) string {
	i = Idx9[i]
	return string([]rune(Door9)[i : i+1]) // + Door0
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
func Jie2YueJian(jie string) string {
	return JieYuejian[jie]
}
func Qi2YueJiang(qi string) string {
	return QiYuejiang[qi]
}

func YueJiang(i int) string {
	i = Idx12[i]
	return string([]rune(MonthJiang)[i : i+1])
}
func YueJian(month int) string {
	if month < 0 {
		month = -month
	}
	i := Idx12[month]
	return string([]rune(MonthBuild)[i : i+1])
}
func BuildStar(i int) string {
	i = Idx12[i]
	return string([]rune(Build12)[i : i+1])
}

// 奇门转盘用转宫宫位索引
var _QMRollIdx = []int{6, 1, 8, 3, 4, 9, 2, 7, 6}     //转宫号=>洛宫号
var _QM2RollIdx = []int{1, 1, 6, 3, 4, 0, 8, 7, 2, 5} //洛宫号=>转宫号
// YueJiangName 月将神名
var YueJiangName = map[string]string{
	"亥": "登明", "戌": "河魅", "酉": "从魁",
	"申": "传送", "未": "小吉", "午": "胜光",
	"巳": "太乙", "辰": "天罡", "卯": "太冲",
	"寅": "功曹", "丑": "大吉", "子": "神后",
}

// JIEQI_MONTH 节气 月建索引 交节换建
var JieYuejian = map[string]string{
	"立春": "寅", "惊蛰": "卯", "清明": "辰",
	"立夏": "巳", "芒种": "午", "小暑": "未",
	"立秋": "申", "白露": "酉", "寒露": "戌",
	"立冬": "亥", "大雪": "子", "小寒": "丑",
}

// QiYuejiang 节气 月将索引 交(中)气换将
var QiYuejiang = map[string]string{
	"雨水": "亥", "春分": "戌", "谷雨": "酉",
	"小满": "申", "夏至": "未", "大暑": "午",
	"处暑": "巳", "秋分": "辰", "霜降": "卯",
	"小雪": "寅", "冬至": "丑", "大寒": "子",
}

// Xiu28 LunarUtil.XIU_LUCK
var Xiu28 = []string{"轸",
	"角", "亢", "氐", "房", "心", "尾", "箕",
	"斗", "牛", "女", "虚", "危", "室", "壁",
	"奎", "娄", "胃", "昴", "毕", "觜", "参",
	"井", "鬼", "柳", "星", "张", "翼", "轸",
}

// 背大将军击对冲
// BigJiang 大将军 以岁支 月支查
var BigJiang = map[string]string{
	"子": "酉", "申": "酉", "辰": "酉",
	"亥": "子", "卯": "子", "未": "子",
	"寅": "卯", "戌": "卯", "午": "卯",
	"巳": "午", "酉": "午", "丑": "午",
}

// 背游都 击鲁都
// YouDu 游都 以日干查
var YouDu = map[string]string{
	"甲": "丑", "己": "丑",
	"乙": "子", "庚": "子",
	"丙": "寅", "辛": "寅",
	"丁": "巳", "壬": "巳",
	"戊": "申", "癸": "申",
}

// LuDu 鲁都 以日干查
var LuDu = map[string]string{
	"甲": "未", "己": "未",
	"乙": "午", "庚": "午",
	"丙": "申", "辛": "申",
	"丁": "亥", "壬": "亥",
	"戊": "寅", "癸": "寅",
}

//	背天雄击地雌
//
// TianXiong 天雄 以月支查
var TianXiong = map[string]string{
	"子": "申", "申": "申", "辰": "申",
	"亥": "亥", "卯": "亥", "未": "亥",
	"寅": "寅", "戌": "寅", "午": "寅",
	"巳": "巳", "酉": "巳", "丑": "巳",
}

// DiCi 地雌 以月支查
var DiCi = map[string]string{
	"子": "寅", "申": "寅", "辰": "寅",
	"亥": "巳", "卯": "巳", "未": "巳",
	"寅": "申", "戌": "申", "午": "申",
	"巳": "亥", "酉": "亥", "丑": "亥",
}

//	背生神击死神
//
// ShengShen 生神 以月支查
var ShengShen = map[string]string{
	"子": "戌", "丑": "亥", "寅": "子",
	"卯": "丑", "辰": "寅", "巳": "卯",
	"午": "辰", "未": "巳", "申": "午",
	"酉": "未", "戌": "申", "亥": "酉",
}

// SiShen 死神 以月支查
var SiShen = map[string]string{
	"子": "辰", "丑": "巳", "寅": "午",
	"卯": "未", "辰": "申", "巳": "酉",
	"午": "戌", "未": "亥", "申": "子",
	"酉": "丑", "戌": "寅", "亥": "卯",
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

// 黄帝有熊氏即位的甲子年(公元前2697年起甲子下元)
// 公元前1年为0,前2年为-1,60年换元
var _QiMenJuYear = []int{0, -1, -4, -7}

// GetHuangDiYear 黄帝纪元
func GetHuangDiYear(year int) int {
	if year < 0 {
		//return 2698 - year
	}
	return year + 2697
}

func GetYear9Yun(year int) int {
	y := GetHuangDiYear(year)
	return (y-60-1)%180/20 + 1
}

func GetYearYuanJu(year int) (int, int) {
	y := GetHuangDiYear(year)
	yuan := (y-60-1)%180/60 + 1
	return yuan, _QiMenJuYear[yuan]
}
func GetYearInChinese(year int) string {
	y := fmt.Sprintf("%d", year)
	s := ""
	j := len(y)
	for i := 0; i < j; i++ {
		s += LunarUtil.NUMBER[[]rune(y[i : i+1])[0]-'0']
	}
	return s
}

// 月家奇门选局
const (
	_Ju1 = "子午卯酉" //四仲上元阴7
	_Ju2 = "寅申巳亥" //四孟中元阴1
	_Ju3 = "辰戌丑未" //四季下元阴4
)

var _QiMenJuMonth = []int{0, -7, -1, -4} //秋分局

// GetHeadGanZhi 找甲己符头
func GetHeadGanZhi(ganZhi string) (string, string) {
	gan := ganZhi[:len(ganZhi)/2]
	zhi := ganZhi[len(ganZhi)/2:]
	var ganIdx, zhiIdx int
	for i, g := range LunarUtil.GAN {
		if g == gan {
			ganIdx = i
			break
		}
	}
	for i, z := range LunarUtil.ZHI {
		if z == zhi {
			zhiIdx = i
			break
		}
	}
	return LunarUtil.GAN[ganIdx], LunarUtil.ZHI[zhiIdx]
}
func GetMonthYuanJu(yearTB string) (int, int) {
	_, zhi := GetHeadGanZhi(yearTB)
	if strings.Contains(_Ju1, zhi) {
		return 1, _QiMenJuMonth[1]
	}
	if strings.Contains(_Ju2, zhi) {
		return 2, _QiMenJuMonth[2]
	}
	return 3, _QiMenJuMonth[3]
}

var _QiMenJuDay = []int{0, 1, 7, 4, -1, -7, -4}

func GetDayYuanJu(jieQiName string) (int, int) {
	jqi := _JieQiIndex[jieQiName]
	yuan := (jqi-1)/4 + 1
	ju := _QiMenJuDay[yuan]
	yuan = yuan % 3
	if yuan == 0 {
		yuan = 3
	}
	return yuan, ju
}

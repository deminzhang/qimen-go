package xuan

// ============ 紫微斗数 - 类型定义 & 数据表 ============

// 五行局
type WuXingJu int

const (
	WuXingJuNone WuXingJu = iota
	ShuiErJu               // 水二局
	MuSanJu                // 木三局
	JinSiJu                // 金四局
	TuWuJu                 // 土五局
	HuoLiuJu               // 火六局
)

var WuXingJuNames = map[WuXingJu]string{
	ShuiErJu: "水二局", MuSanJu: "木三局", JinSiJu: "金四局", TuWuJu: "土五局", HuoLiuJu: "火六局",
}

var WuXingJuNums = map[WuXingJu]int{
	ShuiErJu: 2, MuSanJu: 3, JinSiJu: 4, TuWuJu: 5, HuoLiuJu: 6,
}

// 十二宫名称
var ZiWeiGongNames = []string{
	"命宫", "兄弟", "夫妻", "子女", "财帛", "疾厄",
	"迁移", "交友", "官禄", "田宅", "福德", "父母",
}

// 十四主星名称
var ZhuXingNames = []string{
	"紫微", "天机", "太阳", "武曲", "天同", "廉贞",
	"天府", "太阴", "贪狼", "巨门", "天相", "天梁", "七杀", "破军",
}

// 辅星名称
var FuXingNames = []string{
	"左辅", "右弼", "文昌", "文曲", "天魁", "天钺",
	"禄存", "擎羊", "陀罗", "火星", "铃星", "天马",
}

// 四化类型
type SiHua int

const (
	SiHuaNone SiHua = iota
	HuaLu           // 化禄
	HuaQuan         // 化权
	HuaKe           // 化科
	HuaJi           // 化忌
)

var SiHuaNames = map[SiHua]string{HuaLu: "禄", HuaQuan: "权", HuaKe: "科", HuaJi: "忌"}

// 庙旺利陷
type MiaoWang int

const (
	Miao    MiaoWang = iota // 庙
	Wang                    // 旺
	DeDi                    // 得地
	LiYi                    // 利益
	PingHe                  // 平和
	Xian                    // 陷
)

var MiaoWangNames = map[MiaoWang]string{
	Miao: "庙", Wang: "旺", DeDi: "得", LiYi: "利", PingHe: "平", Xian: "陷",
}

// 星曜数据结构
type Star struct {
	Name    string
	SiHua   SiHua
	MiaoWang MiaoWang
}

// 单宫结构
type ZiWeiPalace struct {
	Index    int    // 0-11
	Name     string // 宫名
	Zhi      string // 十二支
	ZhuXing  []Star // 主星
	FuXing   []Star // 辅星
	DaXian   string // 大限范围
	SiHuaStr string // 四化标签（自化/生年）
}

// 十二宫地支
var ZiWeiPalaceZhi = []string{"寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥", "子", "丑"}

// 天干列表
var TianGanList = []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}

// ============ 星曜索引 ============
const (
	StarZiWei = iota
	StarTianJi
	StarTaiYang
	StarWuQu
	StarTianTong
	StarLianZhen
	StarTianFu
	StarTaiYin
	StarTanLang
	StarJuMen
	StarTianXiang
	StarTianLiang
	StarQiSha
	StarPoJun
	StarCount = 14
)

// ============ 命宫表 ============
// 命宫表 [月][时辰]→宫位索引
// 月1-12, 时辰0-11(子丑寅卯辰巳午未申酉戌亥)
var MingGongTable = [13][12]int{
	{}, // 占位
	/*正月*/ {2, 1, 0, 11, 10, 9, 8, 7, 6, 5, 4, 3},
	/*二月*/ {1, 0, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2},
	/*三月*/ {0, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
	/*四月*/ {11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
	/*五月*/ {10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 11},
	/*六月*/ {9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 11, 10},
	/*七月*/ {8, 7, 6, 5, 4, 3, 2, 1, 0, 11, 10, 9},
	/*八月*/ {7, 6, 5, 4, 3, 2, 1, 0, 11, 10, 9, 8},
	/*九月*/ {6, 5, 4, 3, 2, 1, 0, 11, 10, 9, 8, 7},
	/*十月*/ {5, 4, 3, 2, 1, 0, 11, 10, 9, 8, 7, 6},
	/*冬月*/ {4, 3, 2, 1, 0, 11, 10, 9, 8, 7, 6, 5},
	/*腊月*/ {3, 2, 1, 0, 11, 10, 9, 8, 7, 6, 5, 4},
}

// 身宫表 [月][时辰]→宫位索引
var ShenGongTable = [13][12]int{
	{}, // 占位
	/*正月*/ {2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 0, 1},
	/*二月*/ {3, 4, 5, 6, 7, 8, 9, 10, 11, 0, 1, 2},
	/*三月*/ {4, 5, 6, 7, 8, 9, 10, 11, 0, 1, 2, 3},
	/*四月*/ {5, 6, 7, 8, 9, 10, 11, 0, 1, 2, 3, 4},
	/*五月*/ {6, 7, 8, 9, 10, 11, 0, 1, 2, 3, 4, 5},
	/*六月*/ {7, 8, 9, 10, 11, 0, 1, 2, 3, 4, 5, 6},
	/*七月*/ {8, 9, 10, 11, 0, 1, 2, 3, 4, 5, 6, 7},
	/*八月*/ {9, 10, 11, 0, 1, 2, 3, 4, 5, 6, 7, 8},
	/*九月*/ {10, 11, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	/*十月*/ {11, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	/*冬月*/ {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
	/*腊月*/ {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 0},
}

// ============ 五行局表 ============
// 五行局表 [年干索引][命宫地支索引]
var WuXingJuTable = [10][12]WuXingJu{
	/*甲*/ {HuoLiuJu, HuoLiuJu, JinSiJu, JinSiJu, MuSanJu, MuSanJu, TuWuJu, TuWuJu, ShuiErJu, ShuiErJu, HuoLiuJu, HuoLiuJu},
	/*乙*/ {HuoLiuJu, HuoLiuJu, JinSiJu, JinSiJu, MuSanJu, MuSanJu, TuWuJu, TuWuJu, ShuiErJu, ShuiErJu, HuoLiuJu, HuoLiuJu},
	/*丙*/ {ShuiErJu, ShuiErJu, HuoLiuJu, HuoLiuJu, JinSiJu, JinSiJu, MuSanJu, MuSanJu, TuWuJu, TuWuJu, ShuiErJu, ShuiErJu},
	/*丁*/ {ShuiErJu, ShuiErJu, HuoLiuJu, HuoLiuJu, JinSiJu, JinSiJu, MuSanJu, MuSanJu, TuWuJu, TuWuJu, ShuiErJu, ShuiErJu},
	/*戊*/ {MuSanJu, MuSanJu, TuWuJu, TuWuJu, ShuiErJu, ShuiErJu, HuoLiuJu, HuoLiuJu, JinSiJu, JinSiJu, MuSanJu, MuSanJu},
	/*己*/ {MuSanJu, MuSanJu, TuWuJu, TuWuJu, ShuiErJu, ShuiErJu, HuoLiuJu, HuoLiuJu, JinSiJu, JinSiJu, MuSanJu, MuSanJu},
	/*庚*/ {TuWuJu, TuWuJu, ShuiErJu, ShuiErJu, HuoLiuJu, HuoLiuJu, JinSiJu, JinSiJu, MuSanJu, MuSanJu, TuWuJu, TuWuJu},
	/*辛*/ {TuWuJu, TuWuJu, ShuiErJu, ShuiErJu, HuoLiuJu, HuoLiuJu, JinSiJu, JinSiJu, MuSanJu, MuSanJu, TuWuJu, TuWuJu},
	/*壬*/ {JinSiJu, JinSiJu, MuSanJu, MuSanJu, TuWuJu, TuWuJu, ShuiErJu, ShuiErJu, HuoLiuJu, HuoLiuJu, JinSiJu, JinSiJu},
	/*癸*/ {JinSiJu, JinSiJu, MuSanJu, MuSanJu, TuWuJu, TuWuJu, ShuiErJu, ShuiErJu, HuoLiuJu, HuoLiuJu, JinSiJu, JinSiJu},
}

// ============ 紫微星表 ============
// 紫微星表 [五行局索引][生日(0-29)]→紫微星在十二宫的索引
// 紫微星安星：水二局奇数日/偶数日用不同的行
var ZiWeiStarTable = [7][31]int{
	{}, // 占位
	{}, // 占位
	/*水二局*/ {5, 5, 6, 6, 7, 7, 8, 8, 9, 9, 10, 10, 11, 11, 0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8},
	/*木三局*/ {5, 5, 5, 6, 6, 6, 7, 7, 7, 8, 8, 8, 9, 9, 9, 10, 10, 10, 11, 11, 11, 0, 0, 0, 1, 1, 1, 2, 2, 2, 3},
	/*金四局*/ {5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 8, 8, 8, 9, 9, 9, 9, 10, 10, 10, 10, 11, 11, 11, 11, 0, 0, 0},
	/*土五局*/ {5, 5, 5, 5, 5, 6, 6, 6, 6, 6, 7, 7, 7, 7, 7, 8, 8, 8, 8, 8, 9, 9, 9, 9, 9, 10, 10, 10, 10, 10, 11},
	/*火六局*/ {5, 5, 5, 5, 5, 5, 6, 6, 6, 6, 6, 6, 7, 7, 7, 7, 7, 7, 8, 8, 8, 8, 8, 8, 9, 9, 9, 9, 9, 9, 10},
}

// ============ 十四主星安星表 ============
// 十四主星落宫表 [紫微星位置索引][星曜索引]→宫位索引(0-11)
// 星曜顺序：紫微,天机,太阳,武曲,天同,廉贞,天府,太阴,贪狼,巨门,天相,天梁,七杀,破军
var ZhuXingTable = [12][StarCount]int{
	/*紫微在寅0*/ {0, 2, 4, 6, 8, 10, 0, 2, 4, 6, 8, 10, -1, -1},
	/*紫微在卯1*/ {1, 3, 5, 7, 9, 11, 1, 3, 5, 7, 9, 11, -1, -1},
	/*紫微在辰2*/ {2, 4, 6, 8, 10, 0, 2, 4, 6, 8, 10, 0, -1, -1},
	/*紫微在巳3*/ {3, 5, 7, 9, 11, 1, 3, 5, 7, 9, 11, 1, -1, -1},
	/*紫微在午4*/ {4, 6, 8, 10, 0, 2, 4, 6, 8, 10, 0, 2, -1, -1},
	/*紫微在未5*/ {5, 7, 9, 11, 1, 3, 5, 7, 9, 11, 1, 3, -1, -1},
	/*紫微在申6*/ {6, 8, 10, 0, 2, 4, 6, 8, 10, 0, 2, 4, -1, -1},
	/*紫微在酉7*/ {7, 9, 11, 1, 3, 5, 7, 9, 11, 1, 3, 5, -1, -1},
	/*紫微在戌8*/ {8, 10, 0, 2, 4, 6, 8, 10, 0, 2, 4, 6, -1, -1},
	/*紫微在亥9*/ {9, 11, 1, 3, 5, 7, 9, 11, 1, 3, 5, 7, -1, -1},
	/*紫微在子10*/ {10, 0, 2, 4, 6, 8, 10, 0, 2, 4, 6, 8, -1, -1},
	/*紫微在丑11*/ {11, 1, 3, 5, 7, 9, 11, 1, 3, 5, 7, 9, -1, -1},
}

// 杀破狼特殊标记：七杀=紫微对宫，破军=紫微对宫+1之类
// 简化处理：七杀在紫微的夫妻宫(紫微宫位+0 vs -1)，这里用简表
// 对于常见流派：紫微12位置，七杀固定：寅午戌在紫微对宫，申子辰也在对宫

// 七杀/破军/天府系安置规则：
// 天府 = 紫微对宫（紫微位置+6 mod 12）
// 太阴 = 天府-2, 贪狼=天府-4, 巨门=天府-6, 天相=天府-8, 天梁=天府-10, 七杀=天府对宫
// 破军 = 紫微-3 (紫微退三位)
// 这里使用完整的ZhuXingTable

// ============ 四化表 ============
// 四化表 [年干索引][星曜索引]→四化类型
var SiHuaTable = [10][StarCount]SiHua{
	/*甲*/ {0, 0, HuaJi, HuaKe, 0, HuaLu, 0, 0, 0, 0, 0, 0, 0, HuaQuan},
	/*乙*/ {HuaKe, HuaLu, 0, 0, 0, 0, 0, HuaJi, 0, 0, 0, HuaQuan, 0, 0},
	/*丙*/ {0, HuaQuan, 0, 0, HuaLu, HuaJi, 0, 0, 0, 0, 0, 0, 0, 0},
	/*丁*/ {0, HuaKe, 0, 0, HuaQuan, 0, 0, HuaLu, 0, HuaJi, 0, 0, 0, 0},
	/*戊*/ {0, HuaJi, 0, 0, 0, 0, 0, HuaQuan, HuaLu, 0, 0, 0, 0, 0},
	/*己*/ {0, 0, 0, HuaLu, 0, 0, 0, 0, HuaQuan, 0, 0, HuaKe, 0, 0},
	/*庚*/ {0, 0, HuaLu, HuaQuan, HuaJi, 0, 0, HuaKe, 0, 0, 0, 0, 0, 0},
	/*辛*/ {0, 0, HuaQuan, 0, 0, 0, 0, 0, 0, HuaLu, 0, 0, 0, 0},
	/*壬*/ {HuaQuan, 0, 0, HuaJi, 0, 0, 0, 0, 0, 0, 0, HuaLu, 0, 0},
	/*癸*/ {0, 0, 0, 0, 0, 0, 0, HuaKe, HuaJi, HuaQuan, 0, 0, 0, HuaLu},
}

// ============ 大限起龄表 ============
var DaXianQiLing = map[WuXingJu]int{
	ShuiErJu: 2, MuSanJu: 3, JinSiJu: 4, TuWuJu: 5, HuoLiuJu: 6,
}

// 大限顺逆表 [阴阳][性别]
// 阴年: 乙丁己辛癸, 阳年: 甲丙戊庚壬
var DaXianShunNi = [2][2]int{
	{1, -1},  // 阳年生男顺、阳年生女逆
	{-1, 1}, // 阴年生男逆、阴年生女顺
}

// ============ 紫微位置 → 天府位置 ============
var ZiWeiToTianFu = map[int]int{
	0: 6, 1: 5, 2: 4, 3: 3, 4: 2, 5: 1,
	6: 0, 7: 11, 8: 10, 9: 9, 10: 8, 11: 7,
}

// ============ 紫微斗数排盘 ============

type ZiWeiChart struct {
	YearGan   string
	YearZhi   string
	YearNums  int // 农历年数字
	MonthNum  int // 农历月数字(1-12)
	DayNum    int // 农历日数字(1-30)
	HourZhi   string // 时支
	HourIdx   int    // 时支索引0-11
	Gender    int    // 0女1男

	MingGongIdx    int       // 命宫索引0-11
	ShenGongIdx    int       // 身宫索引0-11
	WuXingJu       WuXingJu // 五行局
	ZiWeiIdx       int       // 紫微星宫位索引
	TianFuIdx      int       // 天府星宫位索引
	Palaces        [12]ZiWeiPalace
	IsYangYear     bool // 年干是否阳年
	DaXianStartAge int  // 起限年龄
}

// CalcZiWei 计算紫微斗数主盘
func CalcZiWei(yearGan, yearZhi string, month, day int, hourZhi string, gender int) *ZiWeiChart {
	c := &ZiWeiChart{
		YearGan:  yearGan,
		YearZhi:  yearZhi,
		MonthNum: month,
		DayNum:   day,
		HourZhi:  hourZhi,
		Gender:   gender,
	}

	ganIdx := indexOf(TianGanList, yearGan)
	hourIdx := indexOf(ZHI, hourZhi)
	c.HourIdx = hourIdx

	// 命宫
	c.MingGongIdx = MingGongTable[month][hourIdx]

	// 身宫
	c.ShenGongIdx = ShenGongTable[month][hourIdx]

	// 五行局
	c.WuXingJu = WuXingJuTable[ganIdx][c.MingGongIdx]

	// 紫微星
	c.ZiWeiIdx = calcZiWeiStar(c.WuXingJu, day)

	// 天府星（紫微对宫）
	c.TianFuIdx = (12 - c.ZiWeiIdx) % 12

	// 阴阳年
	yangGans := []string{"甲", "丙", "戊", "庚", "壬"}
	c.IsYangYear = false
	for _, g := range yangGans {
		if yearGan == g {
			c.IsYangYear = true
			break
		}
	}

	// 大限起龄
	c.DaXianStartAge = DaXianQiLing[c.WuXingJu]

	// 安十四主星
	c.setupZhuXing()

	// 四化
	c.setupSiHua(ganIdx)

	// 大限
	c.setupDaXian(ganIdx)

	return c
}

func calcZiWeiStar(wx WuXingJu, day int) int {
	// iztro 算法：局数除日数，商数宫前走
	val := WuXingJuNums[wx]
	offset := 0
	remainder := -1
	quotient := 0

	for remainder != 0 {
		divisor := day + offset
		quotient = divisor / val
		remainder = divisor % val
		if remainder == 0 {
			break
		}
		offset++
	}

	quotient %= 12
	ziweiIdx := quotient - 1

	if offset%2 == 0 {
		ziweiIdx = (ziweiIdx + offset) % 12
	} else {
		ziweiIdx = (ziweiIdx - offset + 12*10) % 12
	}

	if ziweiIdx < 0 {
		ziweiIdx += 12
	}
	return ziweiIdx
}

func (c *ZiWeiChart) setupZhuXing() {
	for i := range c.Palaces {
		c.Palaces[i].Index = i
		c.Palaces[i].Name = ZiWeiGongNames[(i-c.MingGongIdx+12)%12]
		c.Palaces[i].Zhi = ZiWeiPalaceZhi[i]
		c.Palaces[i].ZhuXing = make([]Star, 0)
		c.Palaces[i].FuXing = make([]Star, 0)
	}

	zwIdx := c.ZiWeiIdx
	tfIdx := (12 - zwIdx) % 12 // 天府在对宫
	c.TianFuIdx = tfIdx

	// 紫微星系（逆时针）：紫微,天机,空,太阳,武曲,天同,空,空,廉贞
	ziweiGroup := []struct {
		name string
		off  int
	}{
		{"紫微", 0}, {"天机", 1}, {"太阳", 3},
		{"武曲", 4}, {"天同", 5}, {"廉贞", 8},
	}
	for _, g := range ziweiGroup {
		pos := (zwIdx - g.off + 12*10) % 12
		c.Palaces[pos].ZhuXing = append(c.Palaces[pos].ZhuXing, Star{Name: g.name})
	}

	// 天府星系（顺时针）：天府,太阴,贪狼,巨门,天相,天梁,七杀,空,空,空,破军
	tianfuGroup := []struct {
		name string
		off  int
	}{
		{"天府", 0}, {"太阴", 1}, {"贪狼", 2},
		{"巨门", 3}, {"天相", 4}, {"天梁", 5},
		{"七杀", 6}, {"破军", 10},
	}
	for _, g := range tianfuGroup {
		pos := (tfIdx + g.off) % 12
		c.Palaces[pos].ZhuXing = append(c.Palaces[pos].ZhuXing, Star{Name: g.name})
	}
}

func (c *ZiWeiChart) setupSiHua(ganIdx int) {
	if ganIdx < 0 || ganIdx >= 10 {
		return
	}
	zwIdx := c.ZiWeiIdx
	for starIdx := 0; starIdx < StarCount; starIdx++ {
		sh := SiHuaTable[ganIdx][starIdx]
		if sh == 0 {
			continue
		}
		pos := ZhuXingTable[zwIdx][starIdx]
		if pos < 0 || pos >= 12 {
			continue
		}
		for i := range c.Palaces[pos].ZhuXing {
			if c.Palaces[pos].ZhuXing[i].Name == ZhuXingNames[starIdx] {
				c.Palaces[pos].ZhuXing[i].SiHua = sh
				break
			}
		}
	}
}

func (c *ZiWeiChart) setupDaXian(ganIdx int) {
	// 阳男阴女顺行，阴男阳女逆行
	direction := 1
	if (c.IsYangYear && c.Gender == 1) || (!c.IsYangYear && c.Gender == 0) {
		direction = 1 // 顺行
	} else {
		direction = -1 // 逆行
	}

	startAge := c.DaXianStartAge
	for i := 0; i < 12; i++ {
		palaceIdx := (c.MingGongIdx + i*direction + 12) % 12
		endAge := startAge + WuXingJuNums[c.WuXingJu] - 1
		if i == 11 {
			endAge = 120
		}
		c.Palaces[palaceIdx].DaXian = formatAgeRange(startAge, endAge)
		startAge = endAge + 1
	}
}

func formatAgeRange(start, end int) string {
	if end > 120 {
		end = 120
	}
	startStr := ""
	endStr := ""
	if start < 10 {
		startStr = string(rune('0'+start)) + "岁"
	} else {
		startStr = string(rune('0'+start/10)) + string(rune('0'+start%10)) + "岁"
	}
	if end < 10 {
		endStr = string(rune('0'+end)) + "岁"
	} else {
		endStr = string(rune('0'+end/10)) + string(rune('0'+end%10)) + "岁"
	}
	return startStr + "-" + endStr
	// 简化版，直接返回字符串
}

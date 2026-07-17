package xuan

import (
	"fmt"
	"time"

	"github.com/6tail/lunar-go/calendar"
)

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

var FuXingNames = []string{
	"左辅", "右弼", "文昌", "文曲", "天魁", "天钺",
	"禄存", "天马", "擎羊", "陀罗", "火星", "铃星", "地空", "地劫",
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
	MiaoWangNone MiaoWang = iota // 无
	Miao                         // 庙
	Wang                         // 旺
	DeDi                         // 得地
	LiYi                         // 利益
	PingHe                       // 平和
	Xian                         // 陷
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
	Index        int    // 0-11
	Name         string // 宫名
	Zhi          string // 十二支
	ZhuXing      []Star // 主星
	FuXing       []Star // 辅星
	ZaYao        []Star // 杂耀
	ChangSheng   string // 长生十二神
	BoShi        string // 博士十二神
	SuiQian      string // 岁前十二神
	JiangQian    string // 将前十二神
	DaXian       string // 大限范围
	XiaoXianAges string // 小限年龄列表
	IsBodyPalace bool   // 是否身宫
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

// 修正12宫索引
func fix12(idx int) int {
	if idx < 0 {
		return fix12(idx + 12)
	}
	return idx % 12
}

// fixIndex 通用索引修正（支持任意模数）
func fixIndex(idx, mod int) int {
	if idx < 0 {
		return fixIndex(idx+mod, mod)
	}
	return idx % mod
}

// ============ 十四主星安星表 ============
// 十四主星落宫表 [紫微星位置索引][星曜索引]→宫位索引(0-11)
// 星曜顺序：紫微,天机,太阳,武曲,天同,廉贞,天府,太阴,贪狼,巨门,天相,天梁,七杀,破军
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
// 输入格式与iztro的bySolar保持一致：阳历日期+时辰索引
// solarDateStr: 公历日期 "2006-01-02"
// timeIndex: 时辰索引 0~12 (0=早子时0:00-1:00, 1=丑…11=亥, 12=晚子时23:00-0:00)
// gender: 0=女 1=男
func CalcZiWei(solarDateStr string, timeIndex int, gender int) *ZiWeiChart {
	// 解析公历日期
	solar, err := parseSolarDate(solarDateStr)
	if err != nil {
		return nil
	}
	lunar := calendar.NewLunarFromSolar(solar)

	// 年柱干支
	ganZhiYear := lunar.GetYearInGanZhiExact()
	yearGan := string([]rune(ganZhiYear)[0])
	yearZhi := string([]rune(ganZhiYear)[1])

	// 农历月日
	month := lunar.GetMonth()
	day := lunar.GetDay()

	// 时支（timeIndex → hourZhi）
	hourZhi := zhiFromTimeIndex(timeIndex)
	hourIdx := timeIndex
	if timeIndex >= 12 {
		hourIdx = 0 // 晚子时按子时算
	}

	// 晚子时至次日子时（对齐 iztro dayDivide='forward'）
	if timeIndex == 12 {
		ly := calendar.NewLunarYear(lunar.GetYear())
		lm := ly.GetMonth(month)
		maxDays := lm.GetDayCount()
		day = day + 1
		if day > maxDays {
			day -= maxDays
		}
	}

	c := &ZiWeiChart{
		YearGan:  yearGan,
		YearZhi:  yearZhi,
		MonthNum: month,
		DayNum:   day,
		HourZhi:  hourZhi,
		HourIdx:  hourIdx,
		Gender:   gender,
	}

	ganIdx := indexOf(TianGanList, yearGan)

	// 命宫（寅起正月顺数至生月，再从月上逆数至生时）—— 与iztro算法一致
	monthIdx := month - 1 // month 1-indexed → 0-based宫位索引
	c.MingGongIdx = fix12(monthIdx - hourIdx)

	// 身宫（寅起正月顺数至生月，再从月上顺数至生时）
	c.ShenGongIdx = fix12(monthIdx + hourIdx)

	// 五行局（纳音公式：命宫天干+命宫地支，与iztro一致）
	soulGan := getSoulHeavenlyStem(yearGan, c.MingGongIdx)
	soulZhi := ZiWeiPalaceZhi[c.MingGongIdx]
	c.WuXingJu = calcWuXingJuByNaYin(soulGan, soulZhi)

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
	// 大限起龄（iztro: 局数=起运年龄）
	c.DaXianStartAge = DaXianQiLing[c.WuXingJu]

	// 安十四主星
	c.setupZhuXing()

	// 四化
	c.setupSiHua(ganIdx)

	// 大限（iztro: 每宫10年）
	c.setupDaXian(ganIdx)

	// 小限
	c.setupXiaoXian()

	// 安辅星
	c.setupFuXing()

	// 安杂耀
	c.setupZaYao(solarDateStr, timeIndex)

	// 安长生十二神/博士十二神/岁前十二神/将前十二神
	c.setupDecorative()

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
	for i := range 12 {
		c.Palaces[i].Index = i
		c.Palaces[i].Name = ZiWeiGongNames[(c.MingGongIdx-i+12)%12]
		c.Palaces[i].Zhi = ZiWeiPalaceZhi[i]
		c.Palaces[i].ZhuXing = make([]Star, 0)
		c.Palaces[i].FuXing = make([]Star, 0)
		c.Palaces[i].ZaYao = make([]Star, 0)
		c.Palaces[i].IsBodyPalace = (i == c.ShenGongIdx)
	}

	zwIdx := c.ZiWeiIdx
	tfIdx := (12 - zwIdx) % 12 // 天府在对宫
	c.TianFuIdx = tfIdx
	for starIdx := range StarCount {
		pos := getStarPos(starIdx, zwIdx, tfIdx)
		if pos < 0 || pos >= 12 {
			continue
		}
		c.Palaces[pos].ZhuXing = append(c.Palaces[pos].ZhuXing, Star{
			Name:     ZhuXingNames[starIdx],
			MiaoWang: getZhuXingBrightness(starIdx, pos),
		})
	}
}

func (c *ZiWeiChart) setupSiHua(ganIdx int) {
	if ganIdx < 0 || ganIdx >= 10 {
		return
	}
	zwIdx := c.ZiWeiIdx
	tfIdx := (12 - zwIdx) % 12

	for starIdx := 0; starIdx < StarCount; starIdx++ {
		sh := SiHuaTable[ganIdx][starIdx]
		if sh == 0 {
			continue
		}
		// 用算法确定星曜位置（和 setupZhuXing 一致）
		pos := getStarPos(starIdx, zwIdx, tfIdx)
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

// getStarPos 根据算法返回星曜所在宫位索引
func getStarPos(starIdx, zwIdx, tfIdx int) int {
	// 星曜索引: 紫微0,天机1,太阳2,武曲3,天同4,廉贞5,天府6,
	//          太阴7,贪狼8,巨门9,天相10,天梁11,七杀12,破军13
	// 紫微星系（逆时针）-1 表示跳过（该索引不是紫微系）
	ziweiOffsets := []int{0, 1, 3, 4, 5, 8, -1, -1, -1, -1, -1, -1, -1, -1}
	// 天府星系（顺时针）-1 表示跳过（该索引不是天府系）
	tianfuOffsets := []int{-1, -1, -1, -1, -1, -1, 0, 1, 2, 3, 4, 5, 6, 10}

	if starIdx < len(ziweiOffsets) && ziweiOffsets[starIdx] >= 0 {
		return (zwIdx - ziweiOffsets[starIdx] + 12*10) % 12
	}
	if starIdx < len(tianfuOffsets) && tianfuOffsets[starIdx] >= 0 {
		return (tfIdx + tianfuOffsets[starIdx]) % 12
	}
	return -1
}

func (c *ZiWeiChart) setupDaXian(ganIdx int) {
	// iztro 算法：每宫固定管10年，阳男阴女顺行
	// idx = GENDER[gender]==yinYang ? fixIndex(soulIndex+i) : fixIndex(soulIndex-i)
	yearZhiYinYang := zhiStdIndex(c.YearZhi) % 2 // 0=阳,1=阴
	genderYinYang := 0
	if c.Gender == 0 {
		genderYinYang = 1 // 女=阴
	}
	shunXing := (yearZhiYinYang == genderYinYang) // 同阴阳→顺行(iztro)

	startAge := WuXingJuNums[c.WuXingJu] // 局数=起运年龄
	for i := range 12 {
		var palaceIdx int
		if shunXing {
			palaceIdx = fix12(c.MingGongIdx + i)
		} else {
			palaceIdx = fix12(c.MingGongIdx - i)
		}
		ageStart := startAge + 10*i
		ageEnd := ageStart + 9
		if ageEnd > 120 {
			ageEnd = 120
		}
		c.Palaces[palaceIdx].DaXian = formatAgeRange(ageStart, ageEnd)
	}
}

// setupXiaoXian 安小限（iztro getHoroscope ages）
func (c *ZiWeiChart) setupXiaoXian() {
	// getAgeIndex: 寅午戌→辰(2), 申子辰→戌(8), 巳酉丑→未(5), 亥卯未→丑(11)
	var ageIdx int
	switch c.YearZhi {
	case "寅", "午", "戌":
		ageIdx = zhiToIndex("辰") // 2
	case "申", "子", "辰":
		ageIdx = zhiToIndex("戌") // 8
	case "巳", "酉", "丑":
		ageIdx = zhiToIndex("未") // 5
	case "亥", "卯", "未":
		ageIdx = zhiToIndex("丑") // 11
	}

	// male→forward, female→backward
	for i := range 12 {
		var ageList []int
		for j := range 10 {
			ageList = append(ageList, 12*j+i+1)
		}
		var idx int
		if c.Gender == 1 { // male
			idx = fix12(ageIdx + i)
		} else {
			idx = fix12(ageIdx - i)
		}
		// 格式化年龄列表
		s := ""
		for k, a := range ageList {
			if k > 0 {
				s += ","
			}
			s += itoa(a)
		}
		c.Palaces[idx].XiaoXianAges = s
	}
}

func formatAgeRange(start, end int) string {
	if end > 120 {
		end = 120
	}
	return itoa(start) + "岁-" + itoa(end) + "岁"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 3)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}

// ============ iztro兼容辅助函数 ============

// parseSolarDate 解析公历日期字符串
func parseSolarDate(s string) (*calendar.Solar, error) {
	for _, layout := range []string{"2006-01-02", "2006-1-2", "2006/01/02", "2006/1/2"} {
		t, err := time.Parse(layout, s)
		if err == nil {
			return calendar.NewSolar(t.Year(), int(t.Month()), t.Day(), 12, 0, 0), nil
		}
	}
	return nil, fmt.Errorf("日期格式错误: %s", s)
}

// zhiFromTimeIndex 时辰索引→地支字符串
// timeIndex 0~12：0=早子时, 1=丑…11=亥, 12=晚子时
func zhiFromTimeIndex(ti int) string {
	if ti < 0 || ti > 12 {
		return ""
	}
	dizhi := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	if ti == 12 {
		return "子" // 晚子时仍是子
	}
	return dizhi[ti]
}

// getSoulHeavenlyStem 根据年干和命宫索引计算命宫天干（五虎遁）
func getSoulHeavenlyStem(yearGan string, mingGongIdx int) string {
	// 五虎遁：年干→寅宫天干
	tigerRule := map[string]string{
		"甲": "丙", "乙": "戊", "丙": "庚", "丁": "壬", "戊": "甲",
		"己": "丙", "庚": "戊", "辛": "庚", "壬": "壬", "癸": "甲",
	}
	startGan, ok := tigerRule[yearGan]
	if !ok {
		return ""
	}

	// 从寅宫(0)起startGan，顺数到命宫位置
	ganList := TianGanList
	startIdx := -1
	for i, g := range ganList {
		if g == startGan {
			startIdx = i
			break
		}
	}
	if startIdx < 0 {
		return ""
	}

	soulGanIdx := (startIdx + mingGongIdx) % 10
	return ganList[soulGanIdx]
}

// calcWuXingJuByNaYin 纳音五行→五行局
// 天干取数：甲乙1 丙丁2 戊己3 庚辛4 壬癸5
// 地支取数：子午丑未1 寅申卯酉2 辰戌巳亥3
// 干支数相加，超过5减去5，得1木3局 2金4局 3水2局 4火6局 5土5局
func calcWuXingJuByNaYin(gan, zhi string) WuXingJu {
	var ganNum, zhiNum int

	switch gan {
	case "甲", "乙":
		ganNum = 1
	case "丙", "丁":
		ganNum = 2
	case "戊", "己":
		ganNum = 3
	case "庚", "辛":
		ganNum = 4
	case "壬", "癸":
		ganNum = 5
	default:
		return WuXingJuNone
	}

	switch zhi {
	case "子", "午", "丑", "未":
		zhiNum = 1
	case "寅", "申", "卯", "酉":
		zhiNum = 2
	case "辰", "戌", "巳", "亥":
		zhiNum = 3
	default:
		return WuXingJuNone
	}

	sum := ganNum + zhiNum
	for sum > 5 {
		sum -= 5
	}

	switch sum {
	case 1:
		return MuSanJu  // 1木→木三局
	case 2:
		return JinSiJu  // 2金→金四局
	case 3:
		return ShuiErJu // 3水→水二局
	case 4:
		return HuoLiuJu // 4火→火六局
	case 5:
		return TuWuJu  // 5土→土五局
	default:
		return WuXingJuNone
	}
}



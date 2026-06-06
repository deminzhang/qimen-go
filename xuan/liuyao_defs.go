package xuan

// 六爻排盘 - 数据表

// GuaNames 64卦名称表 key=6位二进制字符串(从初爻到上爻)
var GuaNames = map[string]string{
	"000000": "坤为地", "100000": "山地剥", "010000": "水地比", "110000": "风地观",
	"001000": "雷地豫", "101000": "火地晋", "011000": "泽地萃", "111000": "天地否",
	"000100": "地山谦", "100100": "艮为山", "010100": "水山蹇", "110100": "风山渐",
	"001100": "雷山小过", "101100": "火山旅", "011100": "泽山咸", "111100": "天山遁",
	"000010": "地水师", "100010": "山水蒙", "010010": "坎为水", "110010": "风水涣",
	"001010": "雷水解", "101010": "火水未济", "011010": "泽水困", "111010": "天水讼",
	"000110": "地风升", "100110": "山风蛊", "010110": "水风井", "110110": "巽为风",
	"001110": "雷风恒", "101110": "火风鼎", "011110": "泽风大过", "111110": "天风姤",
	"000001": "地雷复", "100001": "山雷颐", "010001": "水雷屯", "110001": "风雷益",
	"001001": "震为雷", "101001": "火雷噬嗑", "011001": "泽雷随", "111001": "天雷无妄",
	"000101": "地火明夷", "100101": "山火贲", "010101": "水火既济", "110101": "风火家人",
	"001101": "雷火丰", "101101": "离为火", "011101": "泽火革", "111101": "天火同人",
	"000011": "地泽临", "100011": "山泽损", "010011": "水泽节", "110011": "风泽中孚",
	"001011": "雷泽归妹", "101011": "火泽睽", "011011": "兑为泽", "111011": "天泽履",
	"000111": "地天泰", "100111": "山天大畜", "010111": "水天需", "110111": "风天小畜",
	"001111": "雷天大壮", "101111": "火天大有", "011111": "泽天夬", "111111": "乾为天",
}

// GuaCodeMap 三爻卦编码 key=三爻二进制字符, value=卦名
var GuaCodeMap = map[string]string{
	"111": "乾", "000": "坤", "001": "震", "010": "坎",
	"100": "艮", "110": "巽", "101": "离", "011": "兑",
}

// GuiHunList 归魂卦列表 [内卦, 外卦]
var GuiHunList = [][2]string{
	{"离", "乾"}, {"震", "兑"}, {"乾", "离"}, {"兑", "震"},
	{"艮", "巽"}, {"坤", "坎"}, {"巽", "艮"}, {"坎", "坤"},
}

// YouHunList 游魂卦列表 [内卦, 外卦]
var YouHunList = [][2]string{
	{"离", "坤"}, {"震", "艮"}, {"乾", "坎"}, {"兑", "巽"},
	{"艮", "震"}, {"坤", "离"}, {"巽", "兑"}, {"坎", "乾"},
}

// LiuChongList 六冲卦列表 [内卦, 外卦]
var LiuChongList = [][2]string{
	{"乾", "乾"}, {"坎", "坎"}, {"艮", "艮"}, {"震", "震"},
	{"巽", "巽"}, {"离", "离"}, {"坤", "坤"}, {"兑", "兑"},
	{"震", "乾"}, {"乾", "震"},
}

// LiuHeList 六合卦列表 [内卦, 外卦]
var LiuHeList = [][2]string{
	{"坤", "震"}, {"震", "坤"}, {"坎", "兑"}, {"兑", "坎"},
	{"艮", "离"}, {"离", "艮"}, {"坤", "乾"}, {"乾", "坤"},
}

// GanMap 纳甲表 [卦名] -> {内卦天干, 外卦天干}
var GanMap = map[string][2]string{
	"乾": {"甲", "壬"}, "坤": {"乙", "癸"},
	"艮": {"丙", "丙"}, "兑": {"丁", "丁"},
	"坎": {"戊", "戊"}, "离": {"己", "己"},
	"震": {"庚", "庚"}, "巽": {"辛", "辛"},
}

// ZhiMap 纳支表 [卦名] -> {内卦三支, 外卦三支}
var ZhiMap = map[string][2][3]string{
	"乾": {{"子", "寅", "辰"}, {"午", "申", "戌"}},
	"兑": {{"巳", "卯", "丑"}, {"亥", "酉", "未"}},
	"离": {{"卯", "丑", "亥"}, {"酉", "未", "巳"}},
	"震": {{"子", "寅", "辰"}, {"午", "申", "戌"}},
	"巽": {{"丑", "亥", "酉"}, {"未", "巳", "卯"}},
	"坎": {{"寅", "辰", "午"}, {"申", "戌", "子"}},
	"艮": {{"辰", "午", "申"}, {"戌", "子", "寅"}},
	"坤": {{"未", "巳", "卯"}, {"丑", "亥", "酉"}},
}

// GuaWuXing 卦五行
var GuaWuXing = map[string]string{
	"乾": "金", "兑": "金", "艮": "土", "坤": "土",
	"震": "木", "巽": "木", "坎": "水", "离": "火",
}

// LiuShenOrder 六神固定顺序
var LiuShenOrder = []string{"青龙", "朱雀", "勾陈", "螣蛇", "白虎", "玄武"}

// LiuShenGanMap 日干到六神起始索引
var LiuShenGanMap = map[string]int{
	"甲": 0, "乙": 0, "丙": 1, "丁": 1, "戊": 2, "己": 3, "庚": 4, "辛": 4, "壬": 5, "癸": 5,
}

// LiuQinOrder 六亲序
var LiuQinOrder = []string{"孙", "财", "官", "父", "兄"}



// ZhiWuXing 地支五行
var ZhiWuXing = map[string]string{
	"子": "水", "丑": "土", "寅": "木", "卯": "木",
	"辰": "土", "巳": "火", "午": "火", "未": "土",
	"申": "金", "酉": "金", "戌": "土", "亥": "水",
}

// ZHI 地支序
var ZHI = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

// YaoGanZhi 单爻的纳甲干支
type YaoGanZhi struct {
	Yao     int    // 0=阴爻, 1=阳爻
	Type    string // "1"少阳 "0"少阴 "1o"老阳 "0x"老阴
	ShiYing int    // 0=无, 1=世, 2=应
	Gan     string // 天干
	Zhi     string // 地支
	Qin     string // 六亲
	Fu      *YaoFuShen // 伏神
}

// YaoFuShen 伏神
type YaoFuShen struct {
	Qin string // 六亲
	Gan string // 天干
	Zhi string // 地支
}

// LiuYaoResult 六爻排盘结果
type LiuYaoResult struct {
	DateTime   string
	YearGan    string
	YearZhi    string
	MonthGan   string
	MonthZhi   string
	DayGan     string
	DayZhi     string
	HourGan    string
	HourZhi    string
	YaoRaw     []string // 原始爻型
	BaseName   string // 本卦名
	BaseGua    []YaoGanZhi // 本卦六爻
	GuaGong    string // 卦宫
	Alias      string // 归魂/游魂/六合/六冲
	HasBian    bool
	BianName   string // 变卦名
	BianGua    []YaoGanZhi // 变卦六爻
	BianGuaGong string
	BianAlias   string
	DongYao    []int // 动爻索引(1-6)
	SixShen    []string // 六神
	ShenSha    []ShenShaItem // 神煞
	// 额外字段
	LunarMonth  string
	LunarDay    string
	ShiZhi      string
	Hour        int
	Minute      int
	JieQiFrom   string
	JieQiFromDate string
	JieQiTo     string
	JieQiToDate   string
}

type ShenShaItem struct {
	Name string
	Zhi  string
}

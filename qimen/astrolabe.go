package qimen

//七政四余/星盘

var (
	Constellation       = []string{"Ari", "Tau", "Gem", "Can", "Leo", "Vir", "Lib", "Sco", "Sgr", "Cap", "Aqr", "Psc"}
	ConstellationShort  = []string{"羊", "牛", "双", "蟹", "狮", "室", "秤", "蝎", "射", "摩", "瓶", "鱼"}
	ConstellationSymbol = []string{"♈", "♉", "♊", "♋", "♌", "♍", "♎", "♏", "♐", "♑", "♒", "♓"} //需字体支持
	ConstellationCN     = map[string]string{
		"Ari": "白羊座", "Tau": "金牛座", "Gem": "双子座", "Can": "巨蟹座",
		"Leo": "狮子座", "Vir": "室女座", "Lib": "天秤座", "Sco": "天蝎座",
		"Sgr": "射手座", "Cap": "摩羯座", "Aqr": "水瓶座", "Psc": "双鱼座",
	}
	StarSymbol = map[string]string{ //☉☀☽☿♀♂♃♄♅♆♇♈♉♊♋♌♍♎♏♐♑♒♓ 需字体支持
		"日": "☉", "月": "☽", "水": "水", "金": "♀", "火": "♂", "木": "ㄐ", "土": "ち", "天": "♅", "海": "Ψ",
	}

	// AstrolabeGong 星盘宫名
	AstrolabeGong = []string{
		"命宫", "财帛", "交流", "田宅", "娱乐", "健康", "夫妻", "疾厄", "迁移", "事业", "福德", "玄秘",
	}

	// AstrolabeGong74 七政四余宫名
	AstrolabeGong74 = []string{
		"命宫", "财帛", "兄弟", "田宅", "子女", "奴仆", "夫妻", "疾厄", "迁移", "官禄", "福德", "相貌",
	}
)

// XiuAngle 星宿角度
// https://www.jianshu.com/p/30eacfa0c8d1
var XiuAngle = map[string]float32{
	"娄": 34.33,  //(酉宫4.33度-酉宫17.30度) 跨12.97度(古宿跨10.40度)
	"胃": 47.30,  //(酉宫17.30度-酉宫29.78度) 跨12.48度(古宿跨14.80度)
	"昴": 59.78,  //(酉宫29.78度-申宫8.85度) 跨9.07度(古宿跨12.10度)
	"毕": 68.85,  //(申宫8.85度-申宫24.05度) 跨15.20度(古宿跨15.80度)
	"觜": 84.05,  //(申宫24.05度-申宫25.05度) 跨1.00度(古宿跨1.00度)
	"参": 85.05,  //(申宫25.05度-未宫5.66度) 跨10.62度(古宿跨11.80度)
	"井": 95.66,  //(未宫5.66度-午宫6.11度) 跨30.45度(古宿跨30.50度)
	"鬼": 126.11, //(午宫6.11度-午宫11.10度) 跨4.99度(古宿跨2.90度)
	"柳": 131.10, //(午宫11.10度-午宫27.08度) 跨15.98度(古宿跨15.30度)
	"星": 147.08, //(午宫27.08度-巳宫5.08度) 跨7.99度(古宿跨5.90度)
	"张": 155.08, //(巳宫5.08度-巳宫24.15度) 跨19.07度(古宿 跨15.00度)
	"翼": 174.15, //(巳宫24.15度-辰宫11.08度) 跨16.93度(古宿跨18.70度)
	"轸": 191.08, //(辰宫11.08度-辰宫23.89度) 跨12.81度(古宿17.10度)
	"角": 203.89, //(辰宫23.89度-卯宫4.87度) 跨10.99度(古宿12.80度)
	"亢": 214.87, //(卯宫4.87度-卯宫15.46度) 跨10.59度(古宿跨8.90度)
	"氐": 225.46, //(卯宫15.46度-寅宫3.31度) 跨17.85度(古宿跨16.30度)
	"房": 243.31, //(寅宫3.31度-寅宫8.16度) 跨4.85度(古宿跨5.40度)
	"心": 248.16, //(寅宫8.16度-寅宫16.41度) 跨8.25度(古宿跨6.40度)
	"尾": 256.41, //(寅宫16.41度-丑宫1.61度) 跨15.20度(古宿跨18.60度)
	"箕": 271.61, //(丑宫1.61度-丑宫10.55度) 跨8.94度(古宿跨10.70度)
	"斗": 280.55, //(丑宫10.55度-子宫4.46度) 跨23.92度(古宿跨23.80度)
	"牛": 304.46, //(子宫4.46度-子宫12.13度) 跨7.67度(古宿跨7.90度)
	"女": 312.13, //(子宫12.13度-子宫23.80度) 跨11.67度(古宿跨10.90度)
	"虚": 323.80, //(子宫23.80度-亥宫3.77度) 跨9.97度(古宿跨9.40度)
	"危": 333.77, //(亥宫3.77度-亥宫23.84度) 跨20.07度(古宿跨15.30度)
	"室": 353.84, //(亥宫23.84度-戌宫10.59度) 跨16.75度(古宿跨15.80度)
	"壁": 10.59,  //(戌宫10.59度-戌宫22.81度) 跨12.22度(古宿跨8.90度)
	"奎": 22.81,  //(戌宫22.81度-酉宫4.33度) 跨11.52度(古宿跨17.60度)
}

// 水星：公转周期约为87.97个地球日。
// 金星：公转周期约为224.7个地球日。
// 地球：公转周期约为365.26天。
// 火星：公转周期约为686.98个地球日。
// 木星：公转周期约为11.86年。
// 土星：公转周期约为29.46年。
// 天王星：公转周期约为84.81年。
// 海王星：公转周期约为164.8年。
// 冥王星：公转周期约为248年。
// Jupiter0 := calendar.NewSolar(1983, 5, 27, 0, 0, 0)
// JupiterPeriod := 11.86 * 365.25 * 24 * 60
// degreesJ := degreesS + float32(360*float64(pan.Solar.SubtractMinute(Jupiter0))/JupiterPeriod)
// 星体公转周期
const (
	MercuryPeriod = 87.97 * 24 * 60
	VenusPeriod   = 224.7 * 24 * 60
	EarthPeriod   = 365.2422 * 24 * 60
	MarsPeriod    = 686.98 * 24 * 60
	JupiterPeriod = 11.86 * EarthPeriod
	SaturnPeriod  = 29.46 * EarthPeriod
)

// 星体公转起始时间 TODO
const (
	Mercury0 = "1983-05-27 00:00:00"
	Venus0   = "1983-05-27 00:00:00"
	Earth0   = "1983-05-27 00:00:00"
	Mars0    = "1983-05-27 00:00:00"
	Jupiter0 = "1981-03-27 07:41:00"
	Saturn0  = "1983-05-27 00:00:00"
)

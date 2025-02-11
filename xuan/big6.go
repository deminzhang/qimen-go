package xuan

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
	//大六壬课体对应卦及条件
	Big6RenKeTi = map[string]string{
		"元首": "乾",  // 条件：一上克下，余课无克，为元首课。
		"重审": "坤",  // 条件：四课中仅有一组下克上神，为重审课。余课若无克，亦名始入课。
		"知一": "比",  // 条件：四课中有二上克下，或二下克上。
		"涉害": "坎",  // 条件：四课之中同时出现多组上克下或下克上神。
		"遥克": "睽",  // 条件：四课上下俱无相克。
		"昴星": "履",  // 条件：若四课上下无克，又无遥克。
		"别责": "涣",  // 条件：四课之中出现其中两课相同。
		"八专": "同人", // 条件：凡别责课遇到干支同位。
		"伏吟": "艮",  // 条件：月将占时相同，天盘合于地盘本位。
		"反吟": "震",  // 条件：月将与占时为六冲之辰。
		"三光": "贲",  // 条件：1日干旺相；2日支旺相；3用神自四课出旺相；4皆乘吉将。
		"三阳": "晋",  // 条件：1贵人顺行；2日辰旺相有气，居贵人之前；3用神旺相。
		"三奇": "豫",  // 条件：旬奇发用或入传。
		"六仪": "兑",  // 条件：旬仪发用或者入传。
		"时泰": "泰",  // 条件：1太岁和月建，一个发用，一个入传；2初末传乘青龙、六合吉将。
		"官爵": "益",  // 条件：1发用必须是太岁、月建、本命、行年的驿马。
		"富贵": "大有", // 条件：1贵人乘旺相之气发用。
		"龙德": "萃",  // 条件：1太岁和月将相同。
		"轩盖": "升",  // 条件：1午火发用；2三传卯木、子水俱全。
		"铸印": "鼎",  // 条件：凡戌加巳入传。
		"斫轮": "颐",  // 条件：1卯木发用；2卯木必须加庚辛或申酉发用。
		"引从": "临",  // 条件：凡课日辰干支前后上神发用为初末传。
		"亨通": "谦",  // 条件：1初传生曰干，2三传递生曰干，3干支互相相生。
		"繁昌": "咸",  // 条件：1夫妻行年立德方或值德合。
		"荣华": "渐",  // 条件：1禄、马、贵人加临干支或者年命发用。
		"德庆": "需",  // 条件：1是德神发用，2是年命上神乘吉将。
		"合欢": "恒",  // 条件：1是日上神遁干与日干作合。
		"和美": "丰",  // 条件：1日干与支上神作合，日支与日干上神作合。
		"斩关": "井",  // 条件：凡课魁罡加曰辰发用。
		"闭口": "遁",  // 条件：1旬尾加旬首发用。
		"游子": "观",  // 条件：凡课三传皆土，遇旬丁、天马为用。
		"三交": "姤",  // 条件：1四仲曰占，2四仲加干支。
		"乱首": "师",  // 条件：1干加支，被支克。
		"赘婿": "旅",  // 条件：1干加支上克支。
		"冲破": "夬",  // 条件：凡干支之冲神、加破神发用。
		"淫佚": "既济", // 条件：凡课初传卯酉发用。
		"无淫": "小蓄", // 条件：1四课缺一为不备，且有克。
		"解离": "解离", // 条件：无淫课中，遇夫妻行年及行年上神均既冲且克者。
		"孤寡": "革",  // 条件：1以旬空论，阳空为孤，阴空为寡。
		"度厄": "剥",  // 条件：凡三上克下、为幼度厄；三下贼上、为长度厄。
		"迍福": "屯",  // 条件：凡八迍课得五福，为迍福课。
		"侵害": "损",  // 条件：日辰上各加害神发用，为侵害课。
		"刑伤": "讼",  // 条件：1发用刑干；2发用刑支；3发用刑行年。
		"二烦": "明夷", // 条件：课日月宿为仲神。
		"天祸": "大过", // 条件：凡四立日占，得今日干支加昨日干支。
		"天狱": "噬嗑", // 条件：1发用必须是死气、囚气或墓气。
		"天寇": "蹇",  // 条件：占日为四离日。
		"天网": "蒙",  // 条件：1、时用俱克日；2、干前一位为天罗煞。
		"魄化": "蛊",  // 条件：凡白虎带死神死气，临日辰行年发用。
		"龙战": "离",  // 条件：凡卯酉日占，卯酉为用。
		"死奇": "未济", // 条件：凡天罡加日辰阴阳发用。
		"灾厄": "归妹", // 条件：凡课丧车、游魂、伏殃、病符、丧吊、丘墓、岁虎发用者。
		"殃咎": "解",  // 条件：1三传递克日干。
		"九丑": "小过", // 条件：1占日为九丑日。
		"鬼墓": "困",  // 条件：凡干支之墓神、兼日鬼发用，为鬼墓课。
		"励德": "随",  // 条件：（干支）阳神在贵人前，阴神在贵人后。
		"盘珠": "大壮", // 条件：凡事会合成实，吉则成福，凶则成殃。
		"全局": "大畜", // 条件：三传成局，除土局外，别的都是三合局。
		"玄胎": "家人", // 条件：凡课孟神发用，传皆四孟。
		"连珠": "复",  // 含义： 凶者重重，吉亦累累。孕必连胎，事当续举。
		"六纯": "无妄", // 含义：六纯十杂兼物类，三传之说最纷纭。
		"绝嗣": "",   // 条件：四下俱克上，为无禄课。
		"无禄": "",   // 条件：四上俱克下，为绝嗣课。
		//未提:
		//巽
		//否
		//节
		//小畜
		//中孚
	}
	Big6RenKeTiDes = map[string]string{
		"元首": "占事多顺，忧喜皆实，事从外来，事起男子，宜主动。| 遇凶神恶将，上恶而下受欺；或上休囚而下得势，则下强而上受欺。",
		"重审": "卑犯尊，贱犯贵之象。占事多不顺，事从内起，起于女人。| 贵人顺行吉，贵人逆行凶，传生吉，传墓凶。",
		"知一": "占婚姻主不和谐，失物寻人俱在临近。| 上克下发用，有嫌疑；下克上发用，有妒忌。",
		"涉害": "占者凡事艰难，必有稽迟，乃苦尽甘来之象也。| 神将凶，三四克，灾深难解。",
		"遥克": "开始气势汹汹，后来雷声大雨点小。| 三传神将凶，日辰无气。",
		"昴星": "关梁闭塞，津渡稽留。外出轻者灾，重则有死亡、囚禁之祸。| 蛇虎入传，日辰用神囚死大凶。",
		"别责": "课名芜淫，为三角恋之象。诸事不完备，有涩滞牵连之象。| 占断家庭，主闺房淫乱，或夫妇互有外情。",
		"八专": "神将吉，为同心协力，专一之象。| 如果有上下克，则以常法取用。",
		"伏吟": "凡事主屈而不伸，静中思动。选举必成，考试必中。| 如果三传见吉神，又乘天马、德神、天喜，日辰又临旺相，当以吉论。",
		"反吟": "高峰为谷，深谷为陵，变化不定。得物必失，失败反成。| 神将凶，主多损失，动亦无益。",
		"三光": "万事可行，不劳费力，利有攸往。| 如果三传中末见死囚，是三光失明。",
		"三阳": "凡事吉庆，所求皆遂。| 若占病讼遇之，却凶多吉少。",
		"三奇": "凡事吉利，百祸消散。| 如果三奇空亡，精力不足，其福减半。",
		"六仪": "动无阻隔，家集千祥。兆多吉庆，求财相宜。| 如果旬仪、支仪皆入传，且乘天乙吉将，为富贵六仪。",
		"时泰": "万事亨通，灾潜祸消，谋为无阻，婚姻美满。| 如传见空亡，则事多虚喜。",
		"官爵": "富贵荣华，有官迁职，无官得官，财名皆利。| 若驿马逢冲破，主官爵淹留。",
		"富贵": "天降福德，万事新鲜。| 如果贵人临辰戌为坐狱，所占皆凶。",
		"龙德": "君恩及下，万民欢欣。| 不利尊贵求卑下，再带凶煞、日鬼。",
		"轩盖": "高车驷马，招摇过市，诸事吉庆。| 如果三传凶神凶将，克年命、日辰或空亡。",
		"铸印": "投书献策，官职高迁。有进职加薪之喜。| 若逢戌土空亡、日辰无气，名破模损印。",
		"斫轮": "卯木逢初末传引从，名轩车格，有升职之喜。| 如果三传中有墓神，名旧轮再斫。",
		"引从": "凡课日辰干支前后上神发用为初末传。| 此贵人出行，前者引，后者从，故名引从。",
		"亨通": "凡占课得亨通课，三传相生，干支有情。| 如果递生逢空亡，课传中无解救，仍以凶论。",
		"繁昌": "阴阳和合，万物生成。| 如果夫妻行年俱乘衰败气，或互相克害，则名德孕不育。",
		"荣华": "人宅俱利，经营俱亨。| 如果昼夜贵人逆行或者坐辰戌之上，名坐狱。",
		"德庆": "占事逢德庆课，德神在位，诸煞潜藏。| 如果德神为干鬼，德有化鬼之妙，占功名必高中。",
		"合欢": "占事逢合欢课，主乾坤匹配，吉将齐聚。| 三合事关众多，克应要过月。",
		"和美": "上下欢悦，交易大通。| 如果课中逢有刑害，主恩中有怨。",
		"斩关": "主关梁逾越，最利逃亡。| 若带血支、血忌、呻吟、羊刃、三杀，必伤人而走。",
		"闭口": "主禁口闭缄，机关莫测。| 请白虎占病，主痰气阻塞，喉肿舌禁。",
		"游子": "丁马加吉神，主奔走西东。| 乘三奇、六仪等课体，年命曰辰上有冲克救神。",
		"三交": "占事逢三交课，主交加连累，奸私隐匿。| 初传乘空幻合，主门户不利。",
		"乱首": "子忤其父，弟背其兄。| 三传吉神吉将，年命处有克制凶神名患门有解。",
		"赘婿": "主曲意从人，事多牵连。| 如果年命得吉神吉将，仍可摆脱牵制，任意所为。",
		"冲破": "主人情反复，门户不宁。| 吉将不宜冲，凶将却宜冲。",
		"淫佚": "男子就室，女子有家，淫乱成风。| 上克下发用，过在男子；下克上发用，错在女子。",
		"无淫": "男女争斗，两方均不利。| 如果神将吉，又有救神，不以凶论。",
		"解离": "无淫课中，遇夫妻年年神均既冲且克者。| 若占胎孕，亦主损胎。",
		"孤寡": "占主孤独，离乡背井；官易位，财空手。| 如果兼三奇、六仪课或神将皆吉，主反祸为福。",
		"度厄": "占者家宅乖和，老幼不安。| 如果日辰旺相，反主长得幼力，幼得长力。",
		"绝嗣": "上下悖逆，父子分离。| 无禄课占病必死，兵讼后者胜。",
		"迍福": "忧患将至，得病重危，遭官坐死，谋望不成。| 若逢五福，变忧为喜。",
		"侵害": "六亲失靠，骨肉刑伤；财利潜害，疾病欧伤。| 若发用乘吉将，且兼德合，事阻而终成。",
		"刑伤": "主偏倚失位，家门不昌。| 若遇吉神吉将，事有阻但终遂。",
		"二烦": "家有灾祸，荆棘满途之课。| 此课极凶难避，春夏占得之，凶稍轻。",
		"天祸": "以新易旧，天降灾祸，咎事莫为，身宜谨守。| 天祸课若遇绝神发用，各有所主。",
		"天狱": "占主犯法入狱，病未痊愈，出行凶。| 如果发用刑日干，带恶煞尤凶。",
		"天寇": "占事多破坏，主阴阳分离，行人诈破，病危。| 如果月宿加离日地支发用，为祸尤甚。",
		"天网": "凡事阻碍，逃亡遭殃，胎孕损子，病入膏肓。| 如果日辰、行年临旺相气，又遇德神，主危中有救，忧中有喜。",
		"魄化": "魄化课占病大凶，因白虎死丧之神又叠加死气死神，有死亡将临之象。| 若发用为日干的墓神，名白虎衔尸，凶不可言。",
		"龙战": "占主疑惑反复，门户不宁；合者将离，居者将徙。| 如果人行年又在卯酉更凶。",
		"死奇": "天罡为星宿死奇凶恶之神。| 如果初传旺相为吉将，或六处有救神，或辰土为月将，名死奇回光，则除祸为福。",
		"灾厄": "灾厄重重，妖孽为害。| 占主灾厄重重，妖孽为害",
		"殃咎": "递克被人欺，夹克不自由。| 递克被人欺，夹克不自由。",
		"九丑": "占者多凶，刚日男凶，柔日女祸。| 发用若再乘大小时煞，祸不出月。",
		"鬼墓": "盗贼难获，家宅不昌。| 如果初传为鬼墓，末传为长生，名自墓传生。",
		"励德": "阳神前引，阴神后随，则君子吉，小人危。| 主君子迁官，小人退职，利君子不利小人。",
		"盘珠": "凡事会合成实，吉则成福，凶则成殃。| 如果日干、用神旺相，神将吉者大吉。",
		"全局": "三传成局，除土局外，别的都是三合局。| 三传中若有一神与干支上神刑冲破害，名三合犯杀。",
		"玄胎": "事体皆新，胎孕成型。| 发用若为父母，主尊长有灾。",
		"连珠": "凶者重重，吉亦累累。| 三传顺进名进连茹，多宜进，贵人顺行则应事迅速。",
		"六纯": "六阳动达，如登三天。私凶公吉，官职升迁。| 六阴课占孕主女，五阴相继，盗气迤逦脱去，为源消根断。",
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
	p.Reset(l.GetTimeZhi())
	return &p
}

func (p *Big6Ren) Reset(shiZhi string) {
	if p.TimeZhi == shiZhi {
		return
	}
	p.TimeZhi = shiZhi
	shiZhiIdx := ZhiIdx[shiZhi]
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
	p.parseKeti(kes)
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
	yangDay := YinYang[dayGan] == "阳" //阳日
	//初传为自刑的伏吟课为杜传格。刚日伏吟课无克为自任格。柔日伏吟课无克为自信格。
	if hasKe { //四课上下有克，照常取克发用，
		chuan[0] = chuan0
		chuan[1] = XingZhi[chuan[0]]
		if chuan[1] == chuan[0] { //如果初传是自刑的支（即初传为辰、午、酉、亥），则中传取支上神，末传取中传所刑的支。
			chuan[1] = ke4[2].Up
		}
		chuan[2] = XingZhi[chuan[1]]
		if chuan[2] == chuan[1] { //如果中传又是自刑的支（即中传为辰、午、酉、亥），则取与中传相冲的支为末传。
			chuan[2] = ChongZhi[chuan[1]]
		} else {
			chuan[2] = XingZhi[chuan[1]]
		}
		kts = append(kts, "伏吟")
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
				kts = append(kts, "伏吟", "自任")
				return
			} else {
				chuan[2] = XingZhi[chuan[1]]
				kts = append(kts, "伏吟", "自信")
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
				kts = append(kts, "伏吟", "自任")
				return
			} else {
				chuan[2] = XingZhi[chuan[1]]
				kts = append(kts, "伏吟", "自信")
				return
			}
		}
	}
}

// 7.八专法 两课无克号八专，阳日顺行三位取初传，阴日逆行三位取初传，中末总向日上眠。
func (p *Big6Ren) chuan8Zhuan() (chuan [3]string, kts []string) {
	gs := &p.Gong
	ke4 := p.Ke4
	dayGan := p.DayGan
	yangDay := YinYang[dayGan] == "阳" //阳日
	if yangDay {                      //阳日：日干上神在天盘顺数三位为初传，中传末传为干上神。
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
	kts = append(kts, "八专")
	return
}

func (p *Big6Ren) calcChuan() (chuan [3]string, kts []string) {
	ke4 := p.Ke4
	dayGan := p.DayGan
	dayZhi := p.DayZhi
	gs := &p.Gong
	overlap := gs[0].JiangZhi == LunarUtil.ZHI[1] //伏吟
	reverse := gs[0].JiangZhi == LunarUtil.ZHI[7] //反吟
	yangDay := YinYang[dayGan] == "阳"             //阳日
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
			kts = append(kts, "重审")
			if len(shangKeXia) == 0 {
				kts = append(kts, "始入")
			}
			//return
		case 0:
			if len(shangKeXia) == 1 {
				chuan[0] = shangKe[0].Up
				//chuan = p.chuanNormal(shangKe[0].Up)
				kts = append(kts, "元首")
				//return
			}
		}
		if chuan[0] == "" {
			//2.比用法
			//下贼或二三四侵，若逢上克亦同云。
			//常将天日比神用，阳日用阳阴用阴。
			//若或俱比俱不比，立法别有涉害陈。
			//如果四课中有两课或两课以上的下贼上或上克下，且克者与日干的阴阳属性相同（即比），则以与日干相比者为初传。
			//例如，日干为阳，有两课下贼上，其中一课的克者为阳，则取该阳克者为初传。
			//比用.下克上
			var xiaKeBi []Big6Ke
			for _, ke := range xiaKe {
				if YinYang[ke.Up] == YinYang[dayGan] {
					xiaKeBi = append(xiaKeBi, ke)
				}
			}
			if len(xiaKeBi) == 1 {
				chuan[0] = xiaKeBi[0].Up
				kts = append(kts, "知一")
				//return p.chuanNormal(xiaKeBi[0].Up), kts
			}
			if len(xiaKeBi) == 0 {
				//比用.上克下
				var shangKeBi []Big6Ke
				for _, ke := range shangKe {
					if YinYang[ke.Up] == YinYang[dayGan] {
						shangKeBi = append(shangKeBi, ke)
					}
				}
				if len(shangKeBi) == 1 {
					chuan[0] = shangKeBi[0].Up
					kts = append(kts, "知一")
					//return p.chuanNormal(shangKeBi[0].Up), kts
				}
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
		kts = append(kts, "反吟")
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
				if YinYang[up] != YinYang[dayGan] { //排除不比的
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
							if gIdx == upHomeIdx {
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
					if YinYang[up] != YinYang[dayGan] { //排除不比的
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
								if gIdx == upHomeIdx {
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
			return p.chuanNormal(maxUp), []string{"涉害"}
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
					return p.chuanNormal(up), []string{"涉害", "见机"}
				}
			case 0:
				for up := range hits { //如果无四孟，则取四仲上者为用
					if gs[0].JiangZhi == up || gs[3].JiangZhi == up || gs[6].JiangZhi == up || gs[9].JiangZhi == up {
						mid[up] = struct{}{} //察微
					}
				}
				if len(mid) == 1 {
					for up := range mid {
						return p.chuanNormal(up), []string{"涉害", "察微"}
					}
				}
			}
			if len(mid) > 1 {
				//如果孟仲又复相等，阳日取第一课和第二课中先见者为用，阴日则取第三课和第四课先见者为用?
				//戊辰日子上发用 缀瑕 复等
				if yangDay {
					for up := range mid {
						if up == ke4[0].Up || up == ke4[1].Up {
							return p.chuanNormal(up), []string{"涉害", "缀瑕"} // 复等
						}
					}
				} else {
					for up := range mid {
						if up == ke4[2].Up || up == ke4[3].Up {
							return p.chuanNormal(up), []string{"涉害", "缀瑕"} // 复等
						}
					}
				}
			}
		}
		//注：还有一种直接用孟仲法来取三传，就是不管受克深浅，直接按照如上方式去排三传，两种方式各有优缺，各位壬友请自行比较！
	}
	if keRealCnt == 2 {
		return p.chuan8Zhuan()
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
			return p.chuanNormal(keDayGan[0].Up), []string{"蒿矢"}
		case 0:
			var dayGanKe []Big6Ke //日干克者
			for _, ke := range ke4[1:] {
				if GanZhiWuXing[ke.Up] == WuXingKe[GanZhiWuXing[dayGan]] {
					dayGanKe = append(dayGanKe, ke)
				}
			}
			switch len(dayGanKe) {
			case 1:
				return p.chuanNormal(dayGanKe[0].Up), []string{"弹射"}
			case 0:
			default: //日干克者比
				for _, ke := range dayGanKe {
					if YinYang[ke.Up] == YinYang[dayGan] {
						return p.chuanNormal(ke.Up), []string{"遥克", "弹射"}
					}
				}
			}
		default: //克日干者比
			for _, ke := range keDayGan { //比
				if YinYang[ke.Up] == YinYang[dayGan] {
					return p.chuanNormal(ke.Up), []string{"蒿矢"}
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
			kts = append(kts, "昴星", "冬蛇掩目")
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
		kts = append(kts, "别责")
		return
	}
	return
}

// 课体细析
// http://www.360doc.com/content/23/0331/21/46945463_1074566892.shtml
func (p *Big6Ren) parseKeti(ts []string) {
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
		ts = append(ts, "元胎") //玄胎
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
	//gua := Big6RenKeTi[ts[0]] //TODO 细化
	//if gua != "" {
	//	ts = append(ts, gua+"卦")
	//}
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

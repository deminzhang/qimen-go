package xuan

import (
	"testing"

	"github.com/6tail/lunar-go/calendar"
)

// ──────────────────── 与 iztro 一致性测试 ────────────────────
// 以下测试用例来自 iztro 项目的 jest 测试 (src/__tests__/)
// 每个用例使用 iztro 验证过的输入/输出，确保 Go 实现完全对齐

// ─── 五行局 (Five Elements Class) ───

func TestWuXingJu_IztroConsistency(t *testing.T) {
	// 来自 iztro src/__tests__/astro/palace.test.ts getFiveElementsClass()
	cases := []struct {
		gan, zhi string
		expected WuXingJu
	}{
		{"庚", "申", MuSanJu},  // wood3rd
		{"己", "未", HuoLiuJu}, // fire6th
		{"戊", "午", HuoLiuJu}, // fire6th
		{"丁", "巳", TuWuJu},   // earth5th
		{"丙", "辰", TuWuJu},   // earth5th
		{"乙", "卯", ShuiErJu}, // water2nd
		{"甲", "寅", ShuiErJu}, // water2nd
		{"乙", "丑", JinSiJu},  // metal4th
		{"甲", "子", JinSiJu},  // metal4th
		{"癸", "亥", ShuiErJu}, // water2nd
		{"壬", "戌", ShuiErJu}, // water2nd
		{"辛", "酉", MuSanJu},  // wood3rd
	}

	for _, tc := range cases {
		got := calcWuXingJuByNaYin(tc.gan, tc.zhi)
		if got != tc.expected {
			t.Errorf("五行局(%s%s) = %s, 期望 %s (iztro参考值)",
				tc.gan, tc.zhi, WuXingJuNames[got], WuXingJuNames[tc.expected])
		}
	}
}

// ─── 命宫/身宫 ───

func TestSoulBody_IztroConsistency(t *testing.T) {
	// 来自 iztro src/__tests__/astro/palace.test.ts getSoulAndBody()
	// iztro 配置: yearDivide='exact'
	cases := []struct {
		solarDate  string
		timeIndex  int
		soulIndex  int
		bodyIndex  int
		soulGan    string
		soulZhi    string
	}{
		{"2023-01-22", 5, 7, 5, "己", "酉"},
		{"2023-01-22", 6, 6, 6, "戊", "申"},
		{"2023-02-19", 12, 0, 0, "甲", "寅"},
	}

	for _, tc := range cases {
		solar, err := parseSolarDate(tc.solarDate)
		if err != nil {
			t.Fatalf("日期解析失败 %s: %v", tc.solarDate, err)
		}
		lunar := calendar.NewLunarFromSolar(solar)
		ganZhiYear := lunar.GetYearInGanZhiExact()
		yearGan := string([]rune(ganZhiYear)[0])
		month := lunar.GetMonth()

		monthIdx := month - 1
		hourIdx := tc.timeIndex
		if hourIdx >= 12 {
			hourIdx = 0
		}

		soulIdx := fix12(monthIdx - hourIdx)
		bodyIdx := fix12(monthIdx + hourIdx)

		if soulIdx != tc.soulIndex {
			t.Errorf("[%s timeIndex=%d] 命宫索引: 期望 %d, 得到 %d", tc.solarDate, tc.timeIndex, tc.soulIndex, soulIdx)
		}
		if bodyIdx != tc.bodyIndex {
			t.Errorf("[%s timeIndex=%d] 身宫索引: 期望 %d, 得到 %d", tc.solarDate, tc.timeIndex, tc.bodyIndex, bodyIdx)
		}

		soulGan := getSoulHeavenlyStem(yearGan, soulIdx)
		if soulGan != tc.soulGan {
			t.Errorf("[%s timeIndex=%d] 命宫天干: 期望 %s, 得到 %s (年干=%s)",
				tc.solarDate, tc.timeIndex, tc.soulGan, soulGan, yearGan)
		}
		soulZhi := ZiWeiPalaceZhi[soulIdx]
		if soulZhi != tc.soulZhi {
			t.Errorf("[%s timeIndex=%d] 命宫地支: 期望 %s, 得到 %s", tc.solarDate, tc.timeIndex, tc.soulZhi, soulZhi)
		}
	}
}

// ─── 紫微星/天府星 安星 ───

func TestZiWeiPosition_IztroConsistency(t *testing.T) {
	// 来自 iztro src/__tests__/star/star.test.ts getStartIndex()
	cases := []struct {
		solarDate     string
		timeIndex     int
		ziweiIndex    int
		tianfuIndex   int
	}{
		{"2023-08-01", 0, 11, 1},
		{"2023-08-01", 1, 11, 1},
		{"2023-08-01", 2, 2, 10},
		{"2023-08-01", 3, 2, 10},
		{"2023-08-01", 4, 6, 6},
		{"2023-08-01", 5, 6, 6},
		{"2023-08-01", 6, 2, 10},
		{"2023-08-01", 7, 2, 10},
		{"2023-08-01", 8, 6, 6},
		{"2023-08-01", 9, 6, 6},
		{"2023-08-01", 10, 4, 8},
		{"2023-08-01", 11, 4, 8},
		{"2023-08-01", 12, 4, 8},
		// 晚子时跨日
		{"2023-02-19", 12, 11, 1},
	}

	for _, tc := range cases {
		solar, err := parseSolarDate(tc.solarDate)
		if err != nil {
			t.Fatalf("日期解析失败 %s: %v", tc.solarDate, err)
		}
		lunar := calendar.NewLunarFromSolar(solar)

		day := lunar.GetDay()
		// 晚子时至次日子时（对齐 iztro dayDivide='forward'）
		if tc.timeIndex == 12 {
			ly := calendar.NewLunarYear(lunar.GetYear())
			lm := ly.GetMonth(lunar.GetMonth())
			maxDays := lm.GetDayCount()
			day = day + 1
			if day > maxDays {
				day -= maxDays
			}
		}

		// 五行局
		month := lunar.GetMonth()
		monthIdx := month - 1
		hourIdx := tc.timeIndex
		if hourIdx >= 12 {
			hourIdx = 0
		}
		soulIdx := fix12(monthIdx - hourIdx)
		ganZhiYear := lunar.GetYearInGanZhiExact()
		yearGan := string([]rune(ganZhiYear)[0])
		soulGan := getSoulHeavenlyStem(yearGan, soulIdx)
		soulZhi := ZiWeiPalaceZhi[soulIdx]
		wx := calcWuXingJuByNaYin(soulGan, soulZhi)

		zwIdx := calcZiWeiStar(wx, day)
		tfIdx := (12 - zwIdx) % 12

		if zwIdx != tc.ziweiIndex {
			t.Errorf("[%s timeIndex=%d] 紫微星: 期望 %d(%s), 得到 %d(%s) [五行局=%s, 农历日=%d]",
				tc.solarDate, tc.timeIndex, tc.ziweiIndex, ZiWeiPalaceZhi[tc.ziweiIndex],
				zwIdx, ZiWeiPalaceZhi[zwIdx], WuXingJuNames[wx], day)
		}
		if tfIdx != tc.tianfuIndex {
			t.Errorf("[%s timeIndex=%d] 天府星: 期望 %d(%s), 得到 %d(%s)",
				tc.solarDate, tc.timeIndex, tc.tianfuIndex, ZiWeiPalaceZhi[tc.tianfuIndex],
				tfIdx, ZiWeiPalaceZhi[tfIdx])
		}
	}
}

// ─── 十四主星安星 ───

func TestMajorStars_IztroConsistency(t *testing.T) {
	// 来自 iztro src/__tests__/star/star.test.ts getMajorStar()
	// solarDate='2023-03-06', timeIndex=4
	chart := CalcZiWei("2023-03-06", 4, 1)
	if chart == nil {
		t.Fatal("CalcZiWei 返回 nil")
	}

	// iztro 期望值：按 palace index 0-11 (寅=0)
	// Palace 0: 七杀 庙
	// Palace 1: 天同 平
	// Palace 2: 武曲 庙
	// Palace 3: 太阳 旺
	// Palace 4: 破军 庙 禄
	// Palace 5: 天机 陷
	// Palace 6: 紫微 旺 + 天府 得
	// Palace 7: 太阴 不 科
	// Palace 8: 贪狼 庙 忌
	// Palace 9: 巨门 旺 权
	// Palace 10: 廉贞 平 + 天相 庙
	// Palace 11: 天梁 旺

	type starExpect struct {
		name   string
		siHua  SiHua
	}
	expectedStars := map[int][]starExpect{
		0:  {{ZhuXingNames[StarQiSha], SiHuaNone}},
		1:  {{ZhuXingNames[StarTianTong], SiHuaNone}},
		2:  {{ZhuXingNames[StarWuQu], SiHuaNone}},
		3:  {{ZhuXingNames[StarTaiYang], SiHuaNone}},
		4:  {{ZhuXingNames[StarPoJun], HuaLu}},
		5:  {{ZhuXingNames[StarTianJi], SiHuaNone}},
		6:  {{ZhuXingNames[StarZiWei], SiHuaNone}, {ZhuXingNames[StarTianFu], SiHuaNone}},
		7:  {{ZhuXingNames[StarTaiYin], HuaKe}},
		8:  {{ZhuXingNames[StarTanLang], HuaJi}},
		9:  {{ZhuXingNames[StarJuMen], HuaQuan}},
		10: {{ZhuXingNames[StarLianZhen], SiHuaNone}, {ZhuXingNames[StarTianXiang], SiHuaNone}},
		11: {{ZhuXingNames[StarTianLiang], SiHuaNone}},
	}

	for palIdx, expected := range expectedStars {
		pal := chart.Palaces[palIdx]
		if len(pal.ZhuXing) != len(expected) {
			t.Errorf("宫位 %d(%s): 期望 %d 颗主星, 得到 %d 颗", palIdx, pal.Zhi, len(expected), len(pal.ZhuXing))
			continue
		}
		for i, exp := range expected {
			got := pal.ZhuXing[i]
			if got.Name != exp.name {
				t.Errorf("宫位 %d(%s) 第%d颗星: 期望 %s, 得到 %s", palIdx, pal.Zhi, i+1, exp.name, got.Name)
			}
			if got.SiHua != exp.siHua {
				t.Errorf("宫位 %d(%s) %s 四化: 期望 %s, 得到 %s",
					palIdx, pal.Zhi, got.Name, SiHuaNames[exp.siHua], SiHuaNames[got.SiHua])
			}
		}
	}
}

// ─── 四化 ───

func TestSiHua_IztroConsistency(t *testing.T) {
	// 2023年（癸卯年）：破军化禄、巨门化权、太阴化科、贪狼化忌
	chart := CalcZiWei("2023-03-06", 4, 1)
	if chart == nil {
		t.Fatal("CalcZiWei 返回 nil")
	}

	// 收集所有四化星
	siHuaMap := map[SiHua]string{}
	for _, pal := range chart.Palaces {
		for _, s := range pal.ZhuXing {
			if s.SiHua != SiHuaNone {
				siHuaMap[s.SiHua] = s.Name
			}
		}
	}

	if siHuaMap[HuaLu] != "破军" {
		t.Errorf("化禄: 期望 破军, 得到 %s", siHuaMap[HuaLu])
	}
	if siHuaMap[HuaQuan] != "巨门" {
		t.Errorf("化权: 期望 巨门, 得到 %s", siHuaMap[HuaQuan])
	}
	if siHuaMap[HuaKe] != "太阴" {
		t.Errorf("化科: 期望 太阴, 得到 %s", siHuaMap[HuaKe])
	}
	if siHuaMap[HuaJi] != "贪狼" {
		t.Errorf("化忌: 期望 贪狼, 得到 %s", siHuaMap[HuaJi])
	}
}

// ─── 综合测试：与 iztro bySolar 交叉验证 ───

func TestFullChart_IztroConsistency(t *testing.T) {
	// 来自 iztro src/__tests__/astro/astro.test.ts bySolar()
	// bySolar('2000-8-16', 2, '女', true)
	chart := CalcZiWei("2000-08-16", 2, 0) // gender=0 女
	if chart == nil {
		t.Fatal("CalcZiWei 返回 nil")
	}

	// iztro 期望:
	// earthlyBranchOfSoulPalace: '午'
	// earthlyBranchOfBodyPalace: '戌'
	// fiveElementsClass: '木三局'
	// soul: '破军' (命主, 基于命宫地支)
	// body: '文昌' (身主, 基于年支)

	soulZhi := ZiWeiPalaceZhi[chart.MingGongIdx]
	bodyZhi := ZiWeiPalaceZhi[chart.ShenGongIdx]

	if soulZhi != "午" {
		t.Errorf("命宫地支: 期望 午, 得到 %s", soulZhi)
	}
	if bodyZhi != "戌" {
		t.Errorf("身宫地支: 期望 戌, 得到 %s", bodyZhi)
	}

	if chart.WuXingJu != MuSanJu {
		t.Errorf("五行局: 期望 %s, 得到 %s", WuXingJuNames[MuSanJu], WuXingJuNames[chart.WuXingJu])
	}
}

// ─── 大限起龄 ───

func TestDaXianStartAge(t *testing.T) {
	// 五行局 → 大限起龄
	cases := map[WuXingJu]int{
		ShuiErJu: 2,
		MuSanJu:  3,
		JinSiJu:  4,
		TuWuJu:   5,
		HuoLiuJu: 6,
	}
	for wx, expected := range cases {
		got := DaXianQiLing[wx]
		if got != expected {
			t.Errorf("大限起龄 %s: 期望 %d, 得到 %d", WuXingJuNames[wx], expected, got)
		}
	}
}

// ─── 边界测试 ───

func TestZiWeiEdgeCases(t *testing.T) {
	// 跨月/跨年边界日期
	cases := []struct {
		solarDate string
		timeIndex int
		// 不检查具体值，只确保不 panic 和返回有效结果
	}{
		{"2023-01-01", 0},
		{"2023-01-01", 12}, // 晚子时跨年边界
		{"2023-12-31", 12},
		{"2024-02-29", 6}, // 闰年
	}

	for _, tc := range cases {
		chart := CalcZiWei(tc.solarDate, tc.timeIndex, 1)
		if chart == nil {
			t.Errorf("CalcZiWei(%s, %d) 返回 nil", tc.solarDate, tc.timeIndex)
			continue
		}
		// 验证基本完整性
		if chart.MingGongIdx < 0 || chart.MingGongIdx >= 12 {
			t.Errorf("CalcZiWei(%s, %d) 命宫索引异常: %d", tc.solarDate, tc.timeIndex, chart.MingGongIdx)
		}
		if chart.ZiWeiIdx < 0 || chart.ZiWeiIdx >= 12 {
			t.Errorf("CalcZiWei(%s, %d) 紫微索引异常: %d", tc.solarDate, tc.timeIndex, chart.ZiWeiIdx)
		}
		// 统计主星数量，应该是 14 颗
		var starCount int
		for _, pal := range chart.Palaces {
			starCount += len(pal.ZhuXing)
		}
		if starCount != 14 {
			t.Errorf("CalcZiWei(%s, %d) 主星数量: %d, 期望 14", tc.solarDate, tc.timeIndex, starCount)
		}
	}
}

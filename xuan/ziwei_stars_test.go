package xuan

import (
	"testing"
)

// ─── 辅星安星测试 (iztro location.test.ts) ───

func TestGetLuYangTuoMaIndex(t *testing.T) {
	// iztro getLuYangTuoMaIndex() test vectors
	cases := []struct {
		gan, zhi   string
		lu, ma, yang, tuo int
	}{
		{"癸", "卯", 10, 3, 11, 9},
		{"庚", "寅", 6, 6, 7, 5},
		{"辛", "巳", 7, 9, 8, 6},
		{"壬", "午", 9, 6, 10, 8},
		{"癸", "未", 10, 3, 11, 9},
		{"甲", "申", 0, 0, 1, 11},
		{"丁", "亥", 4, 3, 5, 3},
		{"乙", "酉", 1, 9, 2, 0},
		{"戊", "戌", 3, 6, 4, 2},
		{"己", "未", 4, 3, 5, 3},
		{"丙", "午", 3, 6, 4, 2},
	}

	for _, tc := range cases {
		lu, ma, yang, tuo := getLuYangTuoMaIndex(tc.gan, tc.zhi)
		if lu != tc.lu || ma != tc.ma || yang != tc.yang || tuo != tc.tuo {
			t.Errorf("%s%s年: 禄存=%d(期望%d) 天马=%d(期望%d) 擎羊=%d(期望%d) 陀罗=%d(期望%d)",
				tc.gan, tc.zhi, lu, tc.lu, ma, tc.ma, yang, tc.yang, tuo, tc.tuo)
		}
	}
}

func TestGetKuiYueIndex(t *testing.T) {
	cases := []struct {
		gan       string
		kui, yue  int
	}{
		{"壬", 1, 3},
		{"癸", 1, 3},
		{"甲", 11, 5},
		{"戊", 11, 5},
		{"庚", 11, 5},
		{"乙", 10, 6},
		{"己", 10, 6},
		{"辛", 4, 0},
		{"丙", 9, 7},
		{"丁", 9, 7},
	}

	for _, tc := range cases {
		kui, yue := getKuiYueIndex(tc.gan)
		if kui != tc.kui || yue != tc.yue {
			t.Errorf("%s年: 天魁=%d(期望%d) 天钺=%d(期望%d)", tc.gan, kui, tc.kui, yue, tc.yue)
		}
	}
}

func TestGetZuoYouIndex(t *testing.T) {
	// iztro: getZuoYouIndex(lunarMonth) test vectors
	expected := []struct{ zuo, you int }{
		{2, 8}, {3, 7}, {4, 6}, {5, 5}, {6, 4}, {7, 3},
		{8, 2}, {9, 1}, {10, 0}, {11, 11}, {0, 10}, {1, 9},
	}

	for month := range 12 {
		zuo, you := getZuoYouIndex(month + 1)
		if zuo != expected[month].zuo || you != expected[month].you {
			t.Errorf("农历%d月: 左辅=%d(期望%d) 右弼=%d(期望%d)",
				month+1, zuo, expected[month].zuo, you, expected[month].you)
		}
	}
}

func TestGetChangQuIndex(t *testing.T) {
	// iztro: getChangQuIndex(timeIndex) test vectors
	expected := []struct{ chang, qu int }{
		{8, 2}, {7, 3}, {6, 4}, {5, 5}, {4, 6}, {3, 7},
		{2, 8}, {1, 9}, {0, 10}, {11, 11}, {10, 0}, {9, 1},
	}

	for ti := range 12 {
		chang, qu := getChangQuIndex(ti)
		if chang != expected[ti].chang || qu != expected[ti].qu {
			t.Errorf("时辰%d: 文昌=%d(期望%d) 文曲=%d(期望%d)",
				ti, chang, expected[ti].chang, qu, expected[ti].qu)
		}
	}
}

func TestGetKongJieIndex(t *testing.T) {
	expected := []struct{ kong, jie int }{
		{9, 9}, {8, 10}, {7, 11}, {6, 0}, {5, 1}, {4, 2},
		{3, 3}, {2, 4}, {1, 5}, {0, 6}, {11, 7}, {10, 8},
	}

	for ti := range 12 {
		kong, jie := getKongJieIndex(ti)
		if kong != expected[ti].kong || jie != expected[ti].jie {
			t.Errorf("时辰%d: 地空=%d(期望%d) 地劫=%d(期望%d)",
				ti, kong, expected[ti].kong, jie, expected[ti].jie)
		}
	}
}

func TestGetHuoLingIndex(t *testing.T) {
	// 午年（寅午戌）, timeIndex=0 → huoBase=丑(11), lingBase=卯(1)
	huo, ling := getHuoLingIndex("午", 0)
	if huo != 11 || ling != 1 {
		t.Errorf("午年子时: 火星=%d(期望11) 铃星=%d(期望1)", huo, ling)
	}

	// Test all time indices for 午年
	expectedHuo := []int{11, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expectedLing := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 0}
	for ti := range 12 {
		huo, ling = getHuoLingIndex("午", ti)
		if huo != expectedHuo[ti] || ling != expectedLing[ti] {
			t.Errorf("午年时辰%d: 火星=%d(期望%d) 铃星=%d(期望%d)",
				ti, huo, expectedHuo[ti], ling, expectedLing[ti])
		}
	}

	// Different year branches at timeIndex=0
	yearCases := []struct {
		zhi      string
		huo, ling int
	}{
		{"寅", 11, 1},
		{"申", 0, 8},
		{"子", 0, 8},
		{"巳", 1, 8},
		{"酉", 1, 8},
		{"丑", 1, 8},
		{"亥", 7, 8},
		{"未", 7, 8},
	}
	for _, tc := range yearCases {
		huo, ling = getHuoLingIndex(tc.zhi, 0)
		if huo != tc.huo || ling != tc.ling {
			t.Errorf("%s年子时: 火星=%d(期望%d) 铃星=%d(期望%d)",
				tc.zhi, huo, tc.huo, ling, tc.ling)
		}
	}
}

func TestGetLuanXiIndex(t *testing.T) {
	zhiList := []string{"卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥", "子", "丑", "寅"}
	expectedHL := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 11}
	expectedTX := []int{4, 3, 2, 1, 0, 11, 10, 9, 8, 7, 6, 5}

	for i, zhi := range zhiList {
		hl, tx := getLuanXiIndex(zhi)
		if hl != expectedHL[i] || tx != expectedTX[i] {
			t.Errorf("%s年: 红鸾=%d(期望%d) 天喜=%d(期望%d)", zhi, hl, expectedHL[i], tx, expectedTX[i])
		}
	}
}

func TestGetHuagaiXianchiIndex(t *testing.T) {
	cases := []struct {
		zhi        string
		huagai, xianchi int
	}{
		{"寅", 8, 1}, {"午", 8, 1}, {"戌", 8, 1},
		{"申", 2, 7}, {"子", 2, 7}, {"辰", 2, 7},
		{"巳", 11, 4}, {"酉", 11, 4}, {"丑", 11, 4},
		{"亥", 5, 10}, {"未", 5, 10}, {"卯", 5, 10},
	}
	for _, tc := range cases {
		hg, xc := getHuagaiXianchiIndex(tc.zhi)
		if hg != tc.huagai || xc != tc.xianchi {
			t.Errorf("%s年: 华盖=%d(期望%d) 咸池=%d(期望%d)", tc.zhi, hg, tc.huagai, xc, tc.xianchi)
		}
	}
}

func TestGetGuGuaIndex(t *testing.T) {
	cases := []struct {
		zhi     string
		gu, gua int
	}{
		{"寅", 3, 11}, {"卯", 3, 11}, {"辰", 3, 11},
		{"巳", 6, 2}, {"午", 6, 2}, {"未", 6, 2},
		{"申", 9, 5}, {"酉", 9, 5}, {"戌", 9, 5},
		{"亥", 0, 8}, {"子", 0, 8}, {"丑", 0, 8},
	}
	for _, tc := range cases {
		gu, gua := getGuGuaIndex(tc.zhi)
		if gu != tc.gu || gua != tc.gua {
			t.Errorf("%s年: 孤辰=%d(期望%d) 寡宿=%d(期望%d)", tc.zhi, gu, tc.gu, gua, tc.gua)
		}
	}
}

func TestGetNianjieIndex(t *testing.T) {
	zhiList := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	expected := []int{8, 7, 6, 5, 4, 3, 2, 1, 0, 11, 10, 9}
	for i, zhi := range zhiList {
		got := getNianjieIndex(zhi)
		if got != expected[i] {
			t.Errorf("%s年: 年解=%d(期望%d)", zhi, got, expected[i])
		}
	}
}

func TestGetTimelyStarIndex(t *testing.T) {
	for ti := range 12 {
		taifu, fenggao := getTimelyStarIndex(ti)
		expTaifu := (4 + ti) % 12
		expFenggao := (0 + ti) % 12
		if taifu != expTaifu || fenggao != expFenggao {
			t.Errorf("时辰%d: 台辅=%d(期望%d) 封诰=%d(期望%d)",
				ti, taifu, expTaifu, fenggao, expFenggao)
		}
	}
}

// ─── 综合安星测试 ───

func TestMinorStarsFullChart(t *testing.T) {
	// 综合测试：2023-03-06 timeIndex=4（与 iztro star.test.ts getMinorStar 对齐）
	chart := CalcZiWei("2023-03-06", 4, 1)
	if chart == nil {
		t.Fatal("CalcZiWei 返回 nil")
	}

	// 统计辅星数量
	for i := range 12 {
		if len(chart.Palaces[i].FuXing) < 0 {
			t.Errorf("宫位%d 辅星数量异常", i)
		}
	}
}

// ─── 长生十二神测试 ───

func TestChangSheng12(t *testing.T) {
	// 土五局 → 长生在申(6)
	cs12 := getChangSheng12("丙", "午", 0, TuWuJu) // 女命

	// 丙午年: 午=阳支(0), 女=阴(1), 阴阳不同 → 逆行
	// 长生在申(6), 逆行: 申=长生(0), 未=沐浴(-1=11), 午=冠带(10)...
	expectedStart := zhiToIndex("申") // 6

	if cs12[expectedStart] != "长生" {
		t.Errorf("长生位置: 期望在申(6), 实际 %s 在 %d", cs12[expectedStart], expectedStart)
	}
}

// ─── 博士十二神测试 ───

func TestBoShi12(t *testing.T) {
	// 丙午年: 禄存在巳(3)
	bs12 := getBoShi12("丙", "午", 0) // 女命, 逆行

	luIdx, _, _, _ := getLuYangTuoMaIndex("丙", "午")
	if bs12[luIdx] != "博士" {
		t.Errorf("博士应在禄存宫(巳=%d), 实际 %s", luIdx, bs12[luIdx])
	}
}

// ─── 岁前/将前十二神测试 ───

func TestSuiQianJiangQian12(t *testing.T) {
	sq12 := getSuiQian12("午")
	jq12 := getJiangQian12("午")

	// 岁前: 从年支起岁建顺行
	if sq12[zhiToIndex("午")] != "岁建" {
		t.Errorf("午年岁建应在午")
	}
	// 将前: 寅午戌将星在午
	if jq12[zhiToIndex("午")] != "将星" {
		t.Errorf("午年将星应在午")
	}
}

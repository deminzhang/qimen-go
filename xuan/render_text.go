package xuan

import (
	"fmt"
	"strings"

	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/6tail/lunar-go/calendar"
)

// ========== 九宫格文本渲染 ==========

// Gong9Layout 洛书数→九宫位置(行列 0-based)
// 4 9 2
// 3 5 7
// 8 1 6
var gong9Layout = map[int][2]int{
	1: {2, 1}, 2: {0, 2}, 3: {1, 0}, 4: {0, 0},
	5: {1, 1}, 6: {2, 2}, 7: {1, 2}, 8: {2, 0}, 9: {0, 1},
}

func Gong9RowCol(luoShuIdx int) (row, col int) {
	p := gong9Layout[luoShuIdx]
	return p[0], p[1]
}

type CellContent struct {
	Lines []string
}

// FormatQiMenPalace 格式化单个宫位内容,返回多行
func FormatQiMenPalace(g *QMPalace, pan *QMPan, qm *QMGame) CellContent {
	star := Star0 + g.Star
	if g.Star == "" {
		star = ""
	}
	door := g.Door + Door0
	if g.Door == "" {
		door = ""
	}

	var kongWang, horse string
	for _, z := range []rune(pan.KongWang) {
		if ZhiGong9[string(z)] == g.Idx {
			kongWang = "〇"
		}
	}
	if ZhiGong9[pan.Horse] == g.Idx {
		horse = "马"
	}

	guestGanTomb := ZhiGong9[QMTomb[g.GuestGan]] == g.Idx
	jiXing := g.God == "值符" && ZhiGong9[QM6YiJiXing[pan.Xun]] == g.Idx

	guestMarker := ""
	if jiXing {
		guestMarker = "刑"
		if guestGanTomb {
			guestMarker = "刑墓"
		}
	} else if guestGanTomb {
		guestMarker = "墓"
	}
	hostGanTomb := ZhiGong9[QMTomb[g.HostGan]] == g.Idx
	hostMarker := ""
	if jiXing {
		hostMarker = "刑"
		if hostGanTomb {
			hostMarker = "刑墓"
		}
	} else if hostGanTomb {
		hostMarker = "墓"
	}

	doorPo := WuXingKe[DoorWuxing[g.Door]] == DiagramsWuxing[Diagrams9(g.Idx)]
	doorMarker := ""
	if doorPo {
		doorMarker = "迫"
	}

	diag := Diagrams9(g.Idx)
	gongClr := Gong9Color[g.Idx]

	line1 := fmt.Sprintf("%d%s%s", g.Idx, diag, gongClr)

	godStr := g.God
	if godStr == "" {
		godStr = "  "
	}
	line2 := fmt.Sprintf("%s %s%s %s%s", godStr, star, guestMarker, g.GuestGan, kongWang)

	line3 := fmt.Sprintf("%s%s %s%s %s", door, doorMarker, g.HostGan, hostMarker, g.HideGan)

	line4 := ""
	if g.Idx == pan.DutyStarPos {
		line4 += "值符"
	}
	if g.Idx == pan.DutyDoorPos {
		if line4 != "" {
			line4 += " "
		}
		line4 += "值使"
	}
	if horse != "" {
		if line4 != "" {
			line4 += " "
		}
		line4 += horse
	}
	if pan.RollHosting > 0 && g.Idx == pan.RollHosting {
		if line4 != "" {
			line4 += " "
		}
		line4 += "禽"
	}

	var lines []string
	lines = append(lines, line1)
	lines = append(lines, centerText(line2, 16))
	lines = append(lines, centerText(line3, 16))
	if line4 != "" {
		lines = append(lines, centerText(line4, 16))
	} else {
		lines = append(lines, "")
	}
	return CellContent{Lines: lines}
}

func centerText(s string, width int) string {
	if len([]rune(s)) >= width {
		return s
	}
	padding := width - len([]rune(s))
	left := padding / 2
	right := padding - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

// RenderQiMen9Gong 渲染奇门九宫格文本
func RenderQiMen9Gong(pan *QMPan, qm *QMGame) string {
	cells := make([][]string, 9)
	for i := 1; i <= 9; i++ {
		cc := FormatQiMenPalace(&pan.Gongs[i], pan, qm)
		cells[i-1] = cc.Lines
	}

	cellH := 4
	cellW := 16

	var sb strings.Builder

	sep := "┌" + strings.Repeat("─", cellW) + "┬" + strings.Repeat("─", cellW) + "┬" + strings.Repeat("─", cellW) + "┐"
	sb.WriteString(sep + "\n")

	for row := 0; row < 3; row++ {
		for line := 0; line < cellH; line++ {
			sb.WriteString("│")
			for col := 0; col < 3; col++ {
				var luoShuIdx int
				for idx, p := range gong9Layout {
					if p[0] == row && p[1] == col {
						luoShuIdx = idx
						break
					}
				}
				content := ""
				if luoShuIdx > 0 && luoShuIdx <= 9 {
					cell := cells[luoShuIdx-1]
					if line < len(cell) {
						content = cell[line]
					}
				}
				sb.WriteString(padCell(content, cellW))
				sb.WriteString("│")
			}
			sb.WriteString("\n")
		}
		if row < 2 {
			sep := "├" + strings.Repeat("─", cellW) + "┼" + strings.Repeat("─", cellW) + "┼" + strings.Repeat("─", cellW) + "┤"
			sb.WriteString(sep + "\n")
		}
	}

	sep = "└" + strings.Repeat("─", cellW) + "┴" + strings.Repeat("─", cellW) + "┴" + strings.Repeat("─", cellW) + "┘"
	sb.WriteString(sep + "\n")

	return sb.String()
}

func padCell(s string, w int) string {
	r := []rune(s)
	if len(r) >= w {
		return string(r[:w])
	}
	return s + strings.Repeat(" ", w-len(r))
}

// RenderQiMenHead 渲染奇门头部信息
func RenderQiMenHead(pan *QMPan, qm *QMGame) string {
	var sb strings.Builder

	solar := qm.Solar
	lunar := qm.Lunar

	sb.WriteString(fmt.Sprintf("公历: %d年%d月%d日 %d:%d\n",
		solar.GetYear(), solar.GetMonth(), solar.GetDay(), solar.GetHour(), solar.GetMinute()))

	var cYear string
	if lunar.GetYear() == 1 {
		cYear = "元年"
	} else if lunar.GetYear() <= 0 {
		cYear = fmt.Sprintf("公元前%d年", -lunar.GetYear()+1)
	} else {
		cYear = lunar.GetYearInChinese()
	}
	sb.WriteString(fmt.Sprintf("农历: %s年 %s %s %s时\n",
		cYear, lunar.GetMonthInChinese()+"月", lunar.GetDayInChinese(), lunar.GetEightChar().GetTimeZhi()))

	sb.WriteString(fmt.Sprintf("干支: %s %s %s %s\n",
		lunar.GetYearInGanZhiExact(), lunar.GetMonthInGanZhiExact(),
		lunar.GetDayInGanZhiExact(), lunar.GetTimeInGanZhi()))
	sb.WriteString(fmt.Sprintf("旬首: %s(年) %s(月) %s(日) %s(时)\n",
		lunar.GetYearXunExact(), lunar.GetMonthXunExact(),
		lunar.GetDayXunExact(), lunar.GetTimeXun()))
	sb.WriteString(fmt.Sprintf("空亡: %s(年) %s(月) %s(日) %s(时)\n",
		lunar.GetYearXunKongExact(), lunar.GetMonthXunKongExact(),
		lunar.GetDayXunKongExact(), lunar.GetTimeXunKong()))

	sb.WriteString(fmt.Sprintf("节气: %s(%s) -> %s(%s)\n",
		pan.JieQi, pan.JieQiDate, pan.JieQiNext, pan.JieQiDateNext))
	sb.WriteString(fmt.Sprintf("局: %s\n", pan.JuText))
	sb.WriteString(fmt.Sprintf("值符: %s落%d宫  值使: %s落%d宫\n",
		Star0+pan.DutyStar, pan.DutyStarPos, pan.DutyDoor+Door0, pan.DutyDoorPos))
	sb.WriteString(fmt.Sprintf("%s\n", pan.YueJiang))
	sb.WriteString(fmt.Sprintf("空亡(时): %s  马星: %s\n", pan.KongWang, pan.Horse))

	return sb.String()
}

// ========== 大六壬文本渲染 ==========

// RenderBig6 渲染大六壬盘
func RenderBig6(b6 *Big6Ren, lm *calendar.Lunar) string {
	var sb strings.Builder
	sb.WriteString("===== 大六壬 =====")
	sb.WriteString(fmt.Sprintf("\n月将: %s(%s)  月建: %s",
		b6.MonthLeader, YueJiangName[b6.MonthLeader], b6.MonthBuild))
	if lm != nil {
		sb.WriteString(fmt.Sprintf("\n干支: %s %s %s %s",
			lm.GetYearInGanZhiExact(), lm.GetMonthInGanZhiExact(),
			lm.GetDayInGanZhiExact(), lm.GetTimeInGanZhi()))
	}
	sb.WriteString(fmt.Sprintf("\n占时: %s时  ", b6.TimeZhi))

	for _, g := range b6.Gong {
		if g.Jiang12 == "贵" {
			if g.Idx <= 5 || g.Idx == 12 {
				sb.WriteString("贵人顺行\n")
			} else {
				sb.WriteString("贵人逆行\n")
			}
			break
		}
	}

	sb.WriteString("\n四课:\n")
	keNames := []string{"干上", "干阴", "支上", "支阴"}
	for i, ke := range b6.Ke4 {
		sb.WriteString(fmt.Sprintf("  %s: %s → %s [%s]\n", keNames[i], ke.Down, ke.Up, ke.God))
	}

	sb.WriteString(fmt.Sprintf("\n三传: %s → %s → %s\n", b6.Chuan[0], b6.Chuan[1], b6.Chuan[2]))
	sb.WriteString(fmt.Sprintf("课体: %s\n", b6.KeTi))

	sb.WriteString("\n十二宫:\n")
	sb.WriteString(fmt.Sprintf("%-4s %-4s %-4s %-6s %-4s %-4s %-4s\n",
		"宫次", "地支", "将星", "天将", "天盘干", "天盘支", "建星"))
	for i := 0; i < 12; i++ {
		g := b6.Gong[i]
		lunarZhi := LunarUtil.ZHI[i+1]
		sb.WriteString(fmt.Sprintf("%-4d %-4s %-4s %-6s %-4s %-4s %-4s\n",
			g.Idx, lunarZhi, g.JiangName, g.Jiang12, g.JiangGan, g.JiangZhi, g.Jian))
	}

	return sb.String()
}

// ========== 梅花易数文本渲染 ==========

// RenderMeiHua 渲染梅花易数
func RenderMeiHua(m *MeHua, lm *calendar.Lunar) string {
	var sb strings.Builder
	sb.WriteString("===== 梅花易数 =====\n")
	if lm != nil {
		sb.WriteString(fmt.Sprintf("时间: %s年 %s月 %s日 %s时\n",
			lm.GetYearInGanZhiExact(), lm.GetMonthZhiExact(),
			lm.GetDayZhiExact(), lm.GetTimeZhi()))
		sb.WriteString(fmt.Sprintf("         %s %s %s %s\n",
			lm.GetYearInChinese()+"年", lm.GetMonthInChinese()+"月",
			lm.GetDayInChinese(), lm.GetTimeZhi()+"时"))
	}

	sb.WriteString(fmt.Sprintf("\n上卦: %s(%s)   下卦: %s(%s)   动爻: %d\n",
		m.GuaUp, Diagrams8Origin[m.GuaUpIdx], m.GuaDown, Diagrams8Origin[m.GuaDownIdx], m.ChangeYaoIdx))

	sb.WriteString(fmt.Sprintf("\n%-16s %-16s %-16s\n", "本卦", "互卦", "变卦"))
	sb.WriteString(fmt.Sprintf("%-16s %-16s %-16s\n", m.GuaOrigin, m.GuaProcess, m.GuaChange))
	sb.WriteString(fmt.Sprintf("%-16s %-16s %-16s\n",
		fmt.Sprintf("上:%s(%s)", m.GuaUp, DiagramsWuxing[m.GuaUp]),
		fmt.Sprintf("上:%s(%s)", m.GuaUpProcess, DiagramsWuxing[m.GuaUpProcess]),
		fmt.Sprintf("上:%s(%s)", m.GuaUpChange, DiagramsWuxing[m.GuaUpChange])))
	sb.WriteString(fmt.Sprintf("%-16s %-16s %-16s\n",
		fmt.Sprintf("下:%s(%s)", m.GuaDown, DiagramsWuxing[m.GuaDown]),
		fmt.Sprintf("下:%s(%s)", m.GuaDownProcess, DiagramsWuxing[m.GuaDownProcess]),
		fmt.Sprintf("下:%s(%s)", m.GuaDownChange, DiagramsWuxing[m.GuaDownChange])))

	var tiYao, yongYao string
	if m.ChangeYaoIdx > 3 {
		tiYao = m.GuaDown
		yongYao = m.GuaUpChange
		sb.WriteString(fmt.Sprintf("\n体(下卦):%s(%s)  用(变上):%s(%s)\n",
			m.GuaDown, DiagramsWuxing[m.GuaDown], m.GuaUpChange, DiagramsWuxing[m.GuaUpChange]))
	} else {
		tiYao = m.GuaUp
		yongYao = m.GuaDownChange
		sb.WriteString(fmt.Sprintf("\n体(上卦):%s(%s)  用(变下):%s(%s)\n",
			m.GuaUp, DiagramsWuxing[m.GuaUp], m.GuaDownChange, DiagramsWuxing[m.GuaDownChange]))
	}
	tiWx := DiagramsWuxing[tiYao]
	yongWx := DiagramsWuxing[yongYao]
	if WuXingKe[yongWx] == tiWx {
		sb.WriteString(fmt.Sprintf("生克: %s(%s)克%s(%s) 凶\n", yongYao, yongWx, tiYao, tiWx))
	} else if WuXingKe[tiWx] == yongWx {
		sb.WriteString(fmt.Sprintf("生克: %s(%s)克%s(%s) 凶\n", tiYao, tiWx, yongYao, yongWx))
	} else if WuXingSheng[yongWx] == tiWx {
		sb.WriteString(fmt.Sprintf("生克: %s(%s)生%s(%s) 吉\n", yongYao, yongWx, tiYao, tiWx))
	} else if WuXingSheng[tiWx] == yongWx {
		sb.WriteString(fmt.Sprintf("生克: %s(%s)生%s(%s) 吉\n", tiYao, tiWx, yongYao, yongWx))
	} else {
		sb.WriteString(fmt.Sprintf("生克: %s(%s)与%s(%s)比和 平\n", tiYao, tiWx, yongYao, yongWx))
	}

	return sb.String()
}

// ========== 八字文本渲染 ==========

func ssShort(soul, gan string) string {
	if gan == "" {
		return ""
	}
	full := LunarUtil.SHI_SHEN[soul+gan]
	short := ShiShenShort[full]
	if short == "" {
		return full
	}
	return short
}

type BaZiRender struct {
	YearGan, YearZhi   string
	MonthGan, MonthZhi string
	DayGan, DayZhi     string
	TimeGan, TimeZhi   string
	DaySoul            string
	Gender             int
	HideGanYear        []string
	HideGanMonth       []string
	HideGanDay         []string
	HideGanTime        []string
	ShiShenYear        string
	ShiShenMonth       string
	ShiShenDay         string
	ShiShenTime        string
	NaYinYear          string
	NaYinMonth         string
	NaYinDay           string
	NaYinTime          string
	KongWangYear       string
	KongWangMonth      string
	KongWangDay        string
	KongWangTime       string
	DaYun              []string
	XiaoYun            []string
	LiuNian            string
	ShenShaY           []string
	ShenShaM           []string
	ShenShaD           []string
	ShenShaT           []string
}

// RenderBaZi 渲染八字, dyGan/dyZhi=当前大运干支, lnGan/lnZhi=流年干支
func RenderBaZi(lunar *calendar.Lunar, gender int, daYun, xiaoYun []string, shenSha [][]string,
	dyGan, dyZhi, lnGan, lnZhi string) string {
	bz := lunar.GetEightChar()
	sb := strings.Builder{}

	sb.WriteString("===== 八字排盘 =====\n")
	genderName := "男"
	if gender == 0 {
		genderName = "女"
	}
	soul := bz.GetDayGan()

	var cYear string
	if lunar.GetYear() == 1 {
		cYear = "元年"
	} else if lunar.GetYear() <= 0 {
		cYear = fmt.Sprintf("公元前%d年", -lunar.GetYear()+1)
	} else {
		cYear = lunar.GetYearInChinese()
	}

	sb.WriteString(fmt.Sprintf("出生: %s年 %s月 %s日 %s时 [%s命]\n",
		cYear, lunar.GetMonthInChinese()+"月", lunar.GetDayInChinese(), lunar.GetTimeZhi(), genderName))
	sb.WriteString(fmt.Sprintf("公历: %d-%d-%d %d:%d\n",
		lunar.GetSolar().GetYear(), lunar.GetSolar().GetMonth(), lunar.GetSolar().GetDay(),
		lunar.GetSolar().GetHour(), lunar.GetSolar().GetMinute()))

	tianGanWithSs := func(gan, ss string) string {
		if ss == "日元" {
			return gan
		}
		short := ShiShenShort[ss]
		if short == "" {
			short = ss
		}
		return gan + short
	}

	// ==== 四柱表 (7列: 年 月 日 时 大运 流年) ====
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%-8s %-8s %-8s %-8s %-8s %-8s %-8s\n", "", "年柱", "月柱", "日柱", "时柱", "大运", "流年"))

	sb.WriteString(fmt.Sprintf("%-8s %-8s %-8s %-8s %-8s %-8s %-8s\n", "天干",
		tianGanWithSs(bz.GetYearGan(), LunarUtil.SHI_SHEN[soul+bz.GetYearGan()]),
		tianGanWithSs(bz.GetMonthGan(), LunarUtil.SHI_SHEN[soul+bz.GetMonthGan()]),
		tianGanWithSs(bz.GetDayGan(), "日元"),
		tianGanWithSs(bz.GetTimeGan(), LunarUtil.SHI_SHEN[soul+bz.GetTimeGan()]),
		tianGanWithSs(dyGan, LunarUtil.SHI_SHEN[soul+dyGan]),
		tianGanWithSs(lnGan, LunarUtil.SHI_SHEN[soul+lnGan])))

	sb.WriteString(fmt.Sprintf("%-8s %-8s %-8s %-8s %-8s %-8s %-8s\n", "地支",
		bz.GetYearZhi(), bz.GetMonthZhi(), bz.GetDayZhi(), bz.GetTimeZhi(), dyZhi, lnZhi))

	hideY := LunarUtil.ZHI_HIDE_GAN[bz.GetYearZhi()]
	hideM := LunarUtil.ZHI_HIDE_GAN[bz.GetMonthZhi()]
	hideD := LunarUtil.ZHI_HIDE_GAN[bz.GetDayZhi()]
	hideT := LunarUtil.ZHI_HIDE_GAN[bz.GetTimeZhi()]
	hideDY := LunarUtil.ZHI_HIDE_GAN[dyZhi]
	hideLN := LunarUtil.ZHI_HIDE_GAN[lnZhi]
	hideAll := [][]string{hideY, hideM, hideD, hideT, hideDY, hideLN}

	maxLines := 0
	for _, h := range hideAll {
		if len(h) > maxLines {
			maxLines = len(h)
		}
	}
	for line := 0; line < maxLines; line++ {
		label := "藏干"
		if line > 0 {
			label = ""
		}
		var cols []string
		for ci := 0; ci < 6; ci++ {
			hides := hideAll[ci]
			if line < len(hides) {
				h := hides[line]
				full := LunarUtil.SHI_SHEN[soul+h]
				short := ShiShenShort[full]
				if short == "" {
					short = full
				}
				cols = append(cols, h+short)
			} else {
				cols = append(cols, "")
			}
		}
		sb.WriteString(fmt.Sprintf("%-8s %-8s %-8s %-8s %-8s %-8s %-8s\n",
			label, cols[0], cols[1], cols[2], cols[3], cols[4], cols[5]))
	}

	nayiDY := "-"
	if dyGan != "" && dyZhi != "" {
		nayiDY = LunarUtil.NAYIN[dyGan+dyZhi]
	}
	nayiLN := "-"
	if lnGan != "" && lnZhi != "" {
		nayiLN = LunarUtil.NAYIN[lnGan+lnZhi]
	}
	xkDY := "-"
	if dyGan != "" && dyZhi != "" {
		xkDY = LunarUtil.GetXunKong(dyGan + dyZhi)
	}
	xkLN := "-"
	if lnGan != "" && lnZhi != "" {
		xkLN = LunarUtil.GetXunKong(lnGan + lnZhi)
	}
	diShiDY := "-"
	if dyZhi != "" {
		diShiDY = ZhangSheng12[soul][dyZhi]
	}
	diShiLN := "-"
	if lnZhi != "" {
		diShiLN = ZhangSheng12[soul][lnZhi]
	}
	sb.WriteString(fmt.Sprintf("%-8s %-8s %-8s %-8s %-8s %-8s %-8s\n", "纳音",
		LunarUtil.NAYIN[bz.GetYearGan()+bz.GetYearZhi()],
		LunarUtil.NAYIN[bz.GetMonthGan()+bz.GetMonthZhi()],
		LunarUtil.NAYIN[bz.GetDayGan()+bz.GetDayZhi()],
		LunarUtil.NAYIN[bz.GetTimeGan()+bz.GetTimeZhi()],
		nayiDY,
		nayiLN))

	sb.WriteString(fmt.Sprintf("%-8s %-8s %-8s %-8s %-8s %-8s %-8s\n", "空亡",
		lunar.GetYearXunKongExact(), lunar.GetMonthXunKongExact(),
		lunar.GetDayXunKongExact(), lunar.GetTimeXunKong(),
		xkDY,
		xkLN))

	sb.WriteString(fmt.Sprintf("%-8s %-8s %-8s %-8s %-8s %-8s %-8s\n", "地势",
		ZhangSheng12[soul][bz.GetYearZhi()],
		ZhangSheng12[soul][bz.GetMonthZhi()],
		ZhangSheng12[soul][bz.GetDayZhi()]+"(自坐)",
		ZhangSheng12[soul][bz.GetTimeZhi()],
		diShiDY,
		diShiLN))

	sb.WriteString(fmt.Sprintf("%-8s %-8s", "\n神煞", "年:"+strings.Join(shenSha[0], ",")))
	sb.WriteString(fmt.Sprintf("\n%8s %-8s", "", "月:"+strings.Join(shenSha[1], ",")))
	sb.WriteString(fmt.Sprintf("\n%8s %-8s", "", "日:"+strings.Join(shenSha[2], ",")))
	sb.WriteString(fmt.Sprintf("\n%8s %-8s", "", "时:"+strings.Join(shenSha[3], ",")))

	if len(daYun) > 0 {
		sb.WriteString(fmt.Sprintf("\n\n大运: %s", strings.Join(daYun, " ")))
	}
	if len(xiaoYun) > 0 {
		sb.WriteString(fmt.Sprintf("\n小运: %s", strings.Join(xiaoYun, " ")))
	}

	sb.WriteString("\n")
	return sb.String()
}

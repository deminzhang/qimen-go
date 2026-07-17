package xuan

import (
	"fmt"
	"strings"
)

// RenderZiWei 渲染紫微斗数排盘文本
func RenderZiWei(c *ZiWeiChart) string {
	var sb strings.Builder
	sep := "═══════════════════════════════════════"

	sb.WriteString(sep + "\n")
	sb.WriteString("             紫 微 斗 数 排 盘\n")
	sb.WriteString(sep + "\n\n")

	sb.WriteString(fmt.Sprintf("四柱: %s%s  性别: %s\n", c.YearGan, c.YearZhi, map[int]string{0: "女", 1: "男"}[c.Gender]))
	sb.WriteString(fmt.Sprintf("命宫: %s  身宫: %s\n",
		c.Palaces[c.MingGongIdx].Name, c.Palaces[c.ShenGongIdx].Name))
	sb.WriteString(fmt.Sprintf("五行局: %s\n", WuXingJuNames[c.WuXingJu]))
	sb.WriteString(fmt.Sprintf("紫微在: %s  天府在: %s\n",
		c.Palaces[c.ZiWeiIdx].Zhi, c.Palaces[c.TianFuIdx].Zhi))
	sb.WriteString(fmt.Sprintf("大限起龄: %d岁\n", c.DaXianStartAge))
	sb.WriteString(sep + "\n\n")

	order := []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 0, 1, 2}

	for _, idx := range order {
		p := c.Palaces[idx]
		bodyMark := ""
		if p.IsBodyPalace {
			bodyMark = "【身】"
		}
		sb.WriteString(fmt.Sprintf("┌─%s─%s%s─┐\n", p.Zhi, p.Name, bodyMark))

		// 主星 + (四化庙旺)
		var starParts []string
		for _, s := range p.ZhuXing {
			part := s.Name
			suffix := ""
			if s.SiHua != SiHuaNone {
				suffix += SiHuaNames[s.SiHua]
			}
			if s.MiaoWang != MiaoWangNone {
				suffix += MiaoWangNames[s.MiaoWang]
			}
			if suffix != "" {
				part += "(" + suffix + ")"
			}
			starParts = append(starParts, part)
		}
		starLine := strings.Join(starParts, " ")
		if starLine == "" {
			starLine = "-"
		}
		sb.WriteString(fmt.Sprintf("│ 主: %s\n", starLine))

		// 辅星
		var fuParts []string
		for _, s := range p.FuXing {
			part := s.Name
			if s.MiaoWang != MiaoWangNone {
				part += "(" + MiaoWangNames[s.MiaoWang] + ")"
			}
			fuParts = append(fuParts, part)
		}
		fuLine := strings.Join(fuParts, " ")
		if fuLine == "" {
			fuLine = "-"
		}
		sb.WriteString(fmt.Sprintf("│ 辅: %s\n", fuLine))

		// 杂耀
		var zaParts []string
		for _, s := range p.ZaYao {
			zaParts = append(zaParts, s.Name)
		}
		zaLine := strings.Join(zaParts, " ")
		if zaLine == "" {
			zaLine = "-"
		}
		sb.WriteString(fmt.Sprintf("│ 杂: %s\n", zaLine))

		// 长生/博士/岁前/将前
		sb.WriteString(fmt.Sprintf("│ %s %s %s %s\n", p.ChangSheng, p.BoShi, p.SuiQian, p.JiangQian))

		// 大限 + 小限
		sb.WriteString(fmt.Sprintf("│ 大限: %s\n", p.DaXian))
		if p.XiaoXianAges != "" {
			sb.WriteString(fmt.Sprintf("│ 小限: %s\n", p.XiaoXianAges))
		}

		sb.WriteString("└───────────────┘\n\n")
	}

	sb.WriteString(sep + "\n")
	return sb.String()
}

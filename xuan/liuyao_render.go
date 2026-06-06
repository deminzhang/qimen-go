package xuan

import (
	"fmt"
	"strings"
)

// RenderLiuYao 渲染六爻排盘文本
func RenderLiuYao(r *LiuYaoResult) string {
	var sb strings.Builder
	sep := "═══════════════════════════════════════"
	sep2 := "───────────────────────────────────────"

	sb.WriteString(sep + "\n")
	sb.WriteString("              六 爻 排 盘\n")
	sb.WriteString(sep + "\n\n")
	sb.WriteString("起卦时间: " + r.DateTime + "\n")
	sb.WriteString("农历: " + r.LunarMonth + "月" + r.LunarDay + "\n\n")
	sb.WriteString("四柱:\n")
	sb.WriteString(fmt.Sprintf("  年柱: %s%s  月柱: %s%s  日柱: %s%s  时柱: %s%s\n\n",
		r.YearGan, r.YearZhi, r.MonthGan, r.MonthZhi, r.DayGan, r.DayZhi, r.HourGan, r.HourZhi))
	sb.WriteString(fmt.Sprintf("节气: %s(%s) → %s(%s)\n\n", r.JieQiFrom, r.JieQiFromDate, r.JieQiTo, r.JieQiToDate))
	sb.WriteString("起卦方式: 随机起卦（模拟三枚铜钱）\n")

	// 原始卦爻
	yaoSymbol := make([]string, 6)
	for i, t := range r.YaoRaw {
		switch t {
		case "1o":
			yaoSymbol[5-i] = "○"
		case "0x":
			yaoSymbol[5-i] = "×"
		case "1":
			yaoSymbol[5-i] = "─"
		default:
			yaoSymbol[5-i] = "·"
		}
	}
	sb.WriteString(fmt.Sprintf("原始卦爻: %s  %s\n", strings.Join(r.YaoRaw, " "), strings.Join(yaoSymbol, "")))

	sb.WriteString("\n" + sep + "\n\n")

	// 本卦
	aliasStr := ""
	if r.Alias != "" {
		aliasStr = " · " + r.Alias
	}
	sb.WriteString(fmt.Sprintf("【本卦】 %s (%s宫%s)\n\n", r.BaseName, r.GuaGong, aliasStr))
	sb.WriteString("  六神    六亲  干支      世应    伏神\n")
	sb.WriteString("  ────────────────────────────────\n")

	shNames := []string{"初", "二", "三", "四", "五", "上"}
	for i := 5; i >= 0; i-- {
		y := r.BaseGua[i]
		shen := r.SixShen[i]
		shiYing := "  "
		if y.ShiYing == 1 {
			shiYing = "世"
		} else if y.ShiYing == 2 {
			shiYing = "应"
		}
		changeTag := "  "
		if y.Type == "1o" {
			changeTag = " ○"
		} else if y.Type == "0x" {
			changeTag = " ×"
		}
		yaoDisplay := " ═ ═"
		if y.Yao == 1 {
			yaoDisplay = " ═══"
		}
		fuShen := ""
		if y.Fu != nil {
			fuShen = fmt.Sprintf("%s %s%s", y.Fu.Qin, y.Fu.Gan, y.Fu.Zhi)
		}
		sb.WriteString(fmt.Sprintf("  %-4s  %s    %s%s    %s%s  %s  %s\n",
			shen, y.Qin, y.Gan, y.Zhi, shiYing, changeTag, yaoDisplay, fuShen))
	}

	sb.WriteString("\n" + sep2 + "\n\n")

	// 变卦
	if r.HasBian {
		bianAliasStr := ""
		if r.BianAlias != "" {
			bianAliasStr = " · " + r.BianAlias
		}
		sb.WriteString(fmt.Sprintf("【变卦】 %s (%s宫%s)\n\n", r.BianName, r.BianGuaGong, bianAliasStr))
		sb.WriteString("  六亲    干支      世应\n")
		sb.WriteString("  ────────────────────────\n")

		for i := 5; i >= 0; i-- {
			y := r.BianGua[i]
			shiYing := "  "
			if y.ShiYing == 1 {
				shiYing = "世"
			} else if y.ShiYing == 2 {
				shiYing = "应"
			}
			yaoDisplay := " ═ ═"
			if y.Yao == 1 {
				yaoDisplay = " ═══"
			}
			sb.WriteString(fmt.Sprintf("  %s      %s%s    %s    %s\n", y.Qin, y.Gan, y.Zhi, shiYing, yaoDisplay))
		}
		sb.WriteString("\n")

		sb.WriteString(fmt.Sprintf("动爻: 第"))
		for idx, d := range r.DongYao {
			if idx > 0 {
				sb.WriteString("、")
			}
			sb.WriteString(shNames[d-1])
		}
		sb.WriteString("爻动\n")
	} else {
		sb.WriteString("【静卦】 无动爻，六爻皆静\n")
	}

	sb.WriteString("\n")

	// 持世
	shiQin := ""
	shiIdx := 0
	for i, y := range r.BaseGua {
		if y.ShiYing == 1 {
			shiQin = y.Qin
			shiIdx = i
			break
		}
	}
	sb.WriteString(fmt.Sprintf("持世: %s爻持世 (%s爻)\n\n", shiQin, shNames[shiIdx]))

	// 卦宫五行 & 六亲分布
	sb.WriteString(sep2 + "\n\n")
	sb.WriteString(fmt.Sprintf("卦宫五行: %s\n", GuaWuXing[r.GuaGong]))
	qinCount := make(map[string]int)
	for _, y := range r.BaseGua {
		qinCount[y.Qin]++
	}
	qinParts := make([]string, 0)
	for _, q := range LiuQinOrder {
		if cnt, ok := qinCount[q]; ok {
			qinParts = append(qinParts, fmt.Sprintf("%s×%d", q, cnt))
		}
	}
	sb.WriteString("本卦六亲分布: " + strings.Join(qinParts, "  ") + "\n\n")

	// 神煞
	sb.WriteString(sep2 + "\n\n神煞一览:\n\n")
	shenShaLines := make([]string, len(r.ShenSha))
	for i, s := range r.ShenSha {
		shenShaLines[i] = fmt.Sprintf("  %s: %s", s.Name, s.Zhi)
	}
	for i := 0; i < len(shenShaLines); i += 4 {
		end := i + 4
		if end > len(shenShaLines) {
			end = len(shenShaLines)
		}
		sb.WriteString(strings.Join(shenShaLines[i:end], "    ") + "\n")
	}

	sb.WriteString("\n" + sep + "\n")
	sb.WriteString(fmt.Sprintf("  排盘时间: %s\n", r.DateTime))
	sb.WriteString(sep + "\n")

	return sb.String()
}

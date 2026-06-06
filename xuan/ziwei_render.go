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

	// 绘制十二宫表格
	// 紫微斗数盘面常用布局
	//       巳        午        未        申
	//     ───────┬───────┬───────┬───────
	//    辰       │        │        │       酉
	//    ────────┼───────┼───────┼────────
	//    卯       │        │        │       戌
	//    ────────┼───────┼───────┼────────
	//    寅       │        │        │       亥
	//    ────────┴───────┴───────┴────────
	//       丑        子        亥        戌(重复标注原理)

	// 简化布局：12宫列表

	// 按飞盘顺序输出
	// 巳(3),午(4),未(5),申(6),酉(7),戌(8),亥(9),子(10),丑(11),寅(0),卯(1),辰(2)
	order := []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 0, 1, 2}

	for _, idx := range order {
		p := c.Palaces[idx]
		sb.WriteString(fmt.Sprintf("┌─%s─%s─┐\n", p.Zhi, p.Name))

		// 主星
		starLine := ""
		for _, s := range p.ZhuXing {
			shStr := ""
			if s.SiHua != SiHuaNone {
				shStr = SiHuaNames[s.SiHua]
			}
			starLine += s.Name + shStr + " "
		}
		if starLine == "" {
			starLine = "无主星"
		}
		sb.WriteString(fmt.Sprintf("│ %s\n", strings.TrimSpace(starLine)))

		// 大限
		sb.WriteString(fmt.Sprintf("│ 大限: %s\n", p.DaXian))

		sb.WriteString("└─────────┘\n")
		sb.WriteString("\n")
	}

	// 四化总结
	sb.WriteString(sep + "\n")
	// 注意：这里写的是身宫那行
	return sb.String()
}

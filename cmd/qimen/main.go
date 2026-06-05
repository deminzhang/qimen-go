package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/6tail/lunar-go/calendar"
	"github.com/deminzhang/qimen-go/xuan"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("用法: qimen <时间> <时家|日家|月家|年家> [盘式] [寄宫法] [起局方式] [暗干起法]")
		fmt.Println("  时间格式: \"2006-01-02 15:04\" 或 \"2006-01-02\"")
		fmt.Println("  盘式: 转盘(默认), 飞盘, 鸣法")
		fmt.Println("  寄宫法: 寄坤(默认), 阳艮阴坤")
		fmt.Println("  起局方式: 拆补(默认), 茅山, 置闰, 自选, 阴盘")
		fmt.Println("  暗干起法: 值使起(默认), 门地起")
		fmt.Println("示例: qimen \"2024-01-15 12:00\" 时家 转盘 寄坤 拆补 值使起")
		os.Exit(1)
	}

	// 解析时间
	timeStr := os.Args[1]
	var solar *calendar.Solar
	if strings.Contains(timeStr, ":") {
		t, err := time.Parse("2006-1-2 15:04", timeStr)
		if err != nil {
			t, err = time.Parse("2006-01-02 15:04", timeStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "时间格式错误: %v\n", err)
				os.Exit(1)
			}
		}
		solar = calendar.NewSolar(t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), 0)
	} else {
		t, err := time.Parse("2006-1-2", timeStr)
		if err != nil {
			t, err = time.Parse("2006-01-02", timeStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "时间格式错误: %v\n", err)
				os.Exit(1)
			}
		}
		solar = calendar.NewSolar(t.Year(), int(t.Month()), t.Day(), 12, 0, 0)
	}

	// 解析YMDH
	ymdh := xuan.QMGameHour
	if len(os.Args) > 2 {
		switch os.Args[2] {
		case "时家":
			ymdh = xuan.QMGameHour
		case "日家":
			ymdh = xuan.QMGameDay
		case "月家":
			ymdh = xuan.QMGameMonth
		case "年家":
			ymdh = xuan.QMGameYear
		default:
			fmt.Fprintf(os.Stderr, "未知类型: %s (时家/日家/月家/年家)\n", os.Args[2])
			os.Exit(1)
		}
	}

	// 解析盘式
	qmType := xuan.QMTypeRotating
	if len(os.Args) > 3 {
		switch os.Args[3] {
		case "转盘":
			qmType = xuan.QMTypeRotating
		case "飞盘":
			qmType = xuan.QMTypeFly
		case "鸣法":
			qmType = xuan.QMTypeAmaze
		default:
			// 尝试数字
			if v, err := strconv.Atoi(os.Args[3]); err == nil && v >= 0 && v <= 2 {
				qmType = v
			} else {
				fmt.Fprintf(os.Stderr, "未知盘式: %s (转盘/飞盘/鸣法)\n", os.Args[3])
				os.Exit(1)
			}
		}
	}

	// 解析寄宫法
	hostingType := xuan.QMHostingType2
	if len(os.Args) > 4 {
		switch os.Args[4] {
		case "寄坤", "中宫寄坤":
			hostingType = xuan.QMHostingType2
		case "阳艮阴坤":
			hostingType = xuan.QMHostingType28
		default:
			if v, err := strconv.Atoi(os.Args[4]); err == nil && v >= 0 && v <= 1 {
				hostingType = v
			} else {
				fmt.Fprintf(os.Stderr, "未知寄宫法: %s (寄坤/阳艮阴坤)\n", os.Args[4])
				os.Exit(1)
			}
		}
	}

	// 解析起局方式
	juType := xuan.QMJuTypeSplit
	if len(os.Args) > 5 {
		switch os.Args[5] {
		case "拆补":
			juType = xuan.QMJuTypeSplit
		case "茅山":
			juType = xuan.QMJuTypeMaoShan
		case "置闰":
			juType = xuan.QMJuTypeZhiRun
		case "自选":
			juType = xuan.QMJuTypeSelf
		case "阴盘":
			juType = xuan.QMJuTypeLunar
		default:
			if v, err := strconv.Atoi(os.Args[5]); err == nil && v >= 0 && v <= 4 {
				juType = v
			} else {
				fmt.Fprintf(os.Stderr, "未知起局方式: %s (拆补/茅山/置闰/自选/阴盘)\n", os.Args[5])
				os.Exit(1)
			}
		}
	}

	// 解析暗干起法
	hideGanType := xuan.QMHideGanDutyDoorHour
	if len(os.Args) > 6 {
		switch os.Args[6] {
		case "值使起", "暗干值使起":
			hideGanType = xuan.QMHideGanDutyDoorHour
		case "门地起", "门地暗干":
			hideGanType = xuan.QMHideGanDoorHomeGan
		default:
			if v, err := strconv.Atoi(os.Args[6]); err == nil && v >= 0 && v <= 1 {
				hideGanType = v
			} else {
				fmt.Fprintf(os.Stderr, "未知暗干起法: %s (值使起/门地起)\n", os.Args[6])
				os.Exit(1)
			}
		}
	}

	params := xuan.QMParams{
		Type:        qmType,
		HostingType: hostingType,
		FlyType:     xuan.QMFlyTypeAllOrder,
		JuType:      juType,
		HideGanType: hideGanType,
		YMDH:        ymdh,
	}

	qm := xuan.NewQMGame(solar, params)

	// 显示对应盘
	switch ymdh {
	case xuan.QMGameHour:
		qm.ShowTimeGame()
	case xuan.QMGameDay:
		qm.ShowDayGame()
	case xuan.QMGameMonth:
		qm.ShowMonthGame()
	case xuan.QMGameYear:
		qm.ShowYearGame()
	}
	qm.CalBig6()

	pan := qm.ShowPan
	ymdhStr := []string{"时家", "日家", "月家", "年家"}[ymdh]
	typeStr := xuan.QMType[qmType]
	hostStr := xuan.QMHostingType[hostingType]
	juStr := xuan.QMJuType[juType]
	hideStr := xuan.QMHideGanType[hideGanType]

	// 生成输出文本
	var output strings.Builder
	output.WriteString(strings.Repeat("=", 60) + "\n")
	output.WriteString(fmt.Sprintf("          奇门遁甲 - %s\n", ymdhStr))
	output.WriteString(fmt.Sprintf("盘式: %s  寄宫: %s  起局: %s  暗干: %s\n",
		typeStr, hostStr, juStr, hideStr))
	output.WriteString(strings.Repeat("=", 60) + "\n\n")

	output.WriteString(xuan.RenderQiMenHead(pan, qm))
	output.WriteString("\n")
	output.WriteString(xuan.RenderQiMen9Gong(pan, qm))

	// 注解
	output.WriteString("\n注解:\n")
	for i := 1; i <= 9; i++ {
		g := pan.Gongs[i]
		diag := xuan.Diagrams9(i)
		wx := xuan.DiagramsWuxing[diag]
		doorPo := xuan.WuXingKe[xuan.DoorWuxing[g.Door]] == wx
		if doorPo && g.Door != "" {
			output.WriteString(fmt.Sprintf("  %d宫 %s门迫%s\n", i, g.Door, diag))
		}
	}
	for i := 1; i <= 9; i++ {
		g := pan.Gongs[i]
		hostGanTomb := xuan.ZhiGong9[xuan.QMTomb[g.HostGan]] == g.Idx
		if hostGanTomb && g.HostGan != "" {
			output.WriteString(fmt.Sprintf("  %d宫 %s入墓\n", i, g.HostGan))
		}
	}

	// 写到文件
	filename := fmt.Sprintf("qimen_%s_%s.txt", solar.ToYmd(), ymdhStr)
	if err := os.WriteFile(filename, []byte(output.String()), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "写入文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("排局成功! 输出文件: %s\n", filename)
}

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
		fmt.Println("用法: bazi <时间> <性别>")
		fmt.Println("  时间格式: \"2006-01-02 15:04\" 或 \"2006-01-02\"")
		fmt.Println("  性别: 男/女 或 1/0")
		fmt.Println("示例: bazi \"2024-01-15 12:00\" 男")
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

	// 解析性别
	gender := 1 // 默认男
	if len(os.Args) > 2 {
		switch os.Args[2] {
		case "男", "1":
			gender = 1
		case "女", "0":
			gender = 0
		default:
			if v, err := strconv.Atoi(os.Args[2]); err == nil {
				gender = v
			} else {
				fmt.Fprintf(os.Stderr, "未知性别: %s (男/女)\n", os.Args[2])
				os.Exit(1)
			}
		}
	}

	lunar := calendar.NewLunarFromSolar(solar)
	bz := lunar.GetEightChar()

	// 神煞
	sss := xuan.CalcShenSha(bz, lunar.GetYearInGanZhiExact(), lunar.GetMonthInGanZhiExact(),
		lunar.GetDayInGanZhiExact(), lunar.GetTimeInGanZhi())

	// 大运
	yun := bz.GetYun(gender)
	daYuns := yun.GetDaYun()
	var daYunStrs []string
	for _, dy := range daYuns {
		if dy.GetIndex() > 0 { // 跳过0岁小运
			daYunStrs = append(daYunStrs, fmt.Sprintf("%s(%d-%d)",
				dy.GetGanZhi(), dy.GetStartYear(), dy.GetEndYear()))
		}
	}

	// 小运
	var xiaoYunStrs []string
	if len(daYuns) > 0 {
		for _, xy := range daYuns[0].GetXiaoYun() {
			xiaoYunStrs = append(xiaoYunStrs, fmt.Sprintf("%s(%d)", xy.GetGanZhi(), xy.GetYear()))
		}
	}

	var shenShaList [][]string
	if len(sss) >= 4 {
		shenShaList = [][]string{sss[0], sss[1], sss[2], sss[3]}
	} else {
		shenShaList = [][]string{{}, {}, {}, {}}
	}

	// 当前大运
	curDYGan, curDYzhi := "", ""
	curYear := lunar.GetYear()
	for _, dy := range daYuns {
		if dy.GetIndex() > 0 && curYear >= dy.GetStartYear() && curYear <= dy.GetEndYear() {
			gz := dy.GetGanZhi()
			curDYGan = string([]rune(gz)[0])
			curDYzhi = string([]rune(gz)[1])
			break
		}
	}
	// 流年
	lnGan := lunar.GetYearGanExact()
	lnZhi := lunar.GetYearZhiExact()

	output := xuan.RenderBaZi(lunar, gender, daYunStrs, xiaoYunStrs, shenShaList, curDYGan, curDYzhi, lnGan, lnZhi)

	// 写到文件
	filename := fmt.Sprintf("bazi_%s_%s.txt", solar.ToYmd(), map[int]string{0: "女", 1: "男"}[gender])
	if err := os.WriteFile(filename, []byte(output), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "写入文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("排盘成功! 输出文件: %s\n", filename)
}

package main

import (
	"fmt"
	"os"

	"github.com/6tail/lunar-go/calendar"
	"github.com/deminzhang/qimen-go/util"
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

	timeStr := os.Args[1]
	solar, err := util.ParseTime(timeStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	gender, err := util.ParseGender(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
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
		if dy.GetIndex() > 0 {
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

	filename := fmt.Sprintf("bazi_%s_%s.txt", solar.ToYmd(), map[int]string{0: "女", 1: "男"}[gender])
	if err := util.WriteResultFile(filename, output); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

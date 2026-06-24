package main

import (
	"fmt"
	"os"

	"github.com/deminzhang/qimen-go/util"
	"github.com/deminzhang/qimen-go/xuan"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("用法: ziwei <时间> <性别>")
		fmt.Println("  时间格式: \"2006-01-02 15:04\"")
		fmt.Println("  性别: 男/女 或 1/0")
		fmt.Println("示例: ziwei \"2024-01-15 12:00\" 男")
		os.Exit(1)
	}

	solar, err := util.ParseTime(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	gender, err := util.ParseGender(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// 公历日期字符串 YYYY-MM-DD
	solarDateStr := fmt.Sprintf("%d-%02d-%02d", solar.GetYear(), solar.GetMonth(), solar.GetDay())

	// 时辰索引（iztro标准）
	hour := solar.GetHour()
	timeIndex := 0
	if hour == 0 {
		timeIndex = 0 // 早子时
	} else if hour == 23 {
		timeIndex = 12 // 晚子时
	} else {
		timeIndex = (hour + 1) / 2
	}

	c := xuan.CalcZiWei(solarDateStr, timeIndex, gender)
	if c == nil {
		fmt.Fprintf(os.Stderr, "排盘失败\n")
		os.Exit(1)
	}

	output := xuan.RenderZiWei(c)

	genderLabel := map[int]string{0: "女", 1: "男"}[gender]
	filename := fmt.Sprintf("ziwei_%s_%s.txt", solarDateStr, genderLabel)
	if err := util.WriteResultFile(filename, output); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

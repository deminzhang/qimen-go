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
		fmt.Println("用法: ziwei <时间> <性别>")
		fmt.Println("  时间格式: \"2006-01-02 15:04\" 或 \"2006-01-02\"")
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

	lunar := calendar.NewLunarFromSolar(solar)

	yearGan := lunar.GetYearInGanZhiExact()
	yearGanStr := string([]rune(yearGan)[0])
	yearZhiStr := string([]rune(yearGan)[1])

	// 农历月、日
	lunarMonth := lunar.GetMonth()
	lunarDay := lunar.GetDay()

	// 时辰
	hour := solar.GetHour()
	hourZhi := hourToZhi(hour)

	genderInt := gender

	c := xuan.CalcZiWei(yearGanStr, yearZhiStr, lunarMonth, lunarDay, hourZhi, genderInt)
	output := xuan.RenderZiWei(c)

	filename := fmt.Sprintf("ziwei_%s_%s.txt", solar.ToYmd(), map[int]string{0: "女", 1: "男"}[gender])
	if err := util.WriteResultFile(filename, output); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func hourToZhi(hour int) string {
	idx := hourToZhiIdx(hour)
	zhiList := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	return zhiList[idx]
}

func hourToZhiIdx(hour int) int {
	if hour >= 23 || hour < 1 {
		return 0 // 子
	}
	if hour < 3 {
		return 1
	}
	if hour < 5 {
		return 2
	}
	if hour < 7 {
		return 3
	}
	if hour < 9 {
		return 4
	}
	if hour < 11 {
		return 5
	}
	if hour < 13 {
		return 6
	}
	if hour < 15 {
		return 7
	}
	if hour < 17 {
		return 8
	}
	if hour < 19 {
		return 9
	}
	if hour < 21 {
		return 10
	}
	return 11
}

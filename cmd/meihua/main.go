package main

import (
	"fmt"
	"os"

	"github.com/6tail/lunar-go/calendar"
	"github.com/deminzhang/qimen-go/util"
	"github.com/deminzhang/qimen-go/xuan"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: meihua <时间>")
		fmt.Println("  时间格式: \"2006-01-02 15:04\" 或 \"2006-01-02\"")
		fmt.Println("  以年月日时起卦")
		fmt.Println("示例: meihua \"2024-01-15 12:00\"")
		os.Exit(1)
	}

	solar, err := util.ParseTime(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	lunar := calendar.NewLunarFromSolar(solar)

	// 以年月日时起卦
	yz := xuan.ZhiIdx[lunar.GetYearZhiExact()]
	mz := lunar.GetMonth()
	dz := lunar.GetDay()
	hz := xuan.ZhiIdx[lunar.GetTimeZhi()]
	up := yz + mz + dz
	down := yz + mz + dz + hz

	var mh xuan.MeHua
	mh.Reset(uint(up), uint(down), uint(down))

	output := xuan.RenderMeiHua(&mh, lunar)

	filename := fmt.Sprintf("meihua_%s.txt", solar.ToYmd())
	if err := util.WriteResultFile(filename, output); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

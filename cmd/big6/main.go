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
		fmt.Println("用法: big6 <时间>")
		fmt.Println("  时间格式: \"2006-01-02 15:04\" 或 \"2006-01-02\"")
		fmt.Println("示例: big6 \"2024-01-15 12:00\"")
		os.Exit(1)
	}

	solar, err := util.ParseTime(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	lunar := calendar.NewLunarFromSolar(solar)
	b6 := xuan.NewBig6Ren(lunar)
	output := xuan.RenderBig6(b6, lunar)

	filename := fmt.Sprintf("big6_%s.txt", solar.ToYmd())
	if err := util.WriteResultFile(filename, output); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

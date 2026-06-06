package main

import (
	"fmt"
	"os"
	"time"

	"github.com/6tail/lunar-go/calendar"
	"github.com/deminzhang/qimen-go/util"
	"github.com/deminzhang/qimen-go/xuan"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: liuyao <时间>")
		fmt.Println("  时间格式: \"2006-01-02 15:04\" 或 \"2006-01-02\"")
		fmt.Println("示例: liuyao \"2024-01-15 12:00\"")
		os.Exit(1)
	}

	solar, err := util.ParseTime(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	lunar := calendar.NewLunarFromSolar(solar)
	yaoRaw := xuan.GenerateYao(time.Now().UnixNano())

	r := xuan.CalcLiuYao(lunar, yaoRaw)
	output := xuan.RenderLiuYao(r)

	filename := fmt.Sprintf("liuyao_%s.txt", solar.ToYmd())
	if err := util.WriteResultFile(filename, output); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

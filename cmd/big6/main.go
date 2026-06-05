package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/6tail/lunar-go/calendar"
	"github.com/deminzhang/qimen-go/xuan"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: big6 <时间>")
		fmt.Println("  时间格式: \"2006-01-02 15:04\" 或 \"2006-01-02\"")
		fmt.Println("示例: big6 \"2024-01-15 12:00\"")
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

	lunar := calendar.NewLunarFromSolar(solar)
	b6 := xuan.NewBig6Ren(lunar)

	output := xuan.RenderBig6(b6, lunar)

	// 写到文件
	filename := fmt.Sprintf("big6_%s.txt", solar.ToYmd())
	if err := os.WriteFile(filename, []byte(output), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "写入文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("排盘成功! 输出文件: %s\n", filename)
}

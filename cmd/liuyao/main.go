package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/6tail/lunar-go/calendar"
	"github.com/deminzhang/qimen-go/util"
	"github.com/deminzhang/qimen-go/xuan"
)

func main() {
	validYao := map[string]bool{"0": true, "1": true, "0x": true, "1o": true}

	if len(os.Args) < 2 {
		fmt.Println("用法:")
		fmt.Println("  liuyao \"2006-01-02 15:04\"            时间自动起卦")
		fmt.Println("  liuyao \"2006-01-02 15:04\" 1 0 1 1o 0x 1  时间+原始爻")
		fmt.Println("\n时间：断卦用，必传参数")
		fmt.Println("原始爻（可选）：1=少阳 0=少阴 1o=老阳(动) 0x=老阴(动)")
		os.Exit(1)
	}

	solar, err := util.ParseTime(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "时间格式错误: %v\n", err)
		os.Exit(1)
	}
	lunar := calendar.NewLunarFromSolar(solar)

	var yaoRaw []string
	if len(os.Args) >= 8 {
		// 时间+原始爻
		yaoRaw = make([]string, 6)
		for i, a := range os.Args[2:8] {
			v := strings.TrimSpace(a)
			if !validYao[v] {
				fmt.Fprintf(os.Stderr, "无效的爻位 '%s'，应为 0/1/0x/1o\n", v)
				os.Exit(1)
			}
			yaoRaw[i] = v
		}
	} else {
		// 仅时间，自动起卦
		yaoRaw = xuan.GenerateYao(time.Now().UnixNano())
	}

	var filename string
	if len(os.Args) >= 8 {
		// 原始爻模式：文件名用爻码
		filename = fmt.Sprintf("liuyao_%s.txt", strings.Join(yaoRaw, ""))
	} else {
		// 自动起卦：文件名用日期
		filename = fmt.Sprintf("liuyao_%s.txt", solar.ToYmd())
	}

	r := xuan.CalcLiuYao(lunar, yaoRaw)
	output := xuan.RenderLiuYao(r)

	if err := util.WriteResultFile(filename, output); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

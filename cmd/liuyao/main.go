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
	var lunar *calendar.Lunar
	var yaoRaw []string
	var filename string

	validYao := map[string]bool{"0": true, "1": true, "0x": true, "1o": true}

	switch {
	case len(os.Args) >= 7:
		// 模式1: liuyao <时间> <爻1> <爻2> ... <爻6>
		// 或:    liuyao <爻1> <爻2> ... <爻6>
		args := os.Args[1:]
		yaoArgs := args
		solar, err := util.ParseTime(args[0])
		if err == nil {
			// 第一个参数是时间
			lunar = calendar.NewLunarFromSolar(solar)
			yaoArgs = args[1:]
		} else {
			// 没有时间，仅用原始爻
			lunar = calendar.NewLunarFromSolar(
				calendar.NewSolar(time.Now().Year(), int(time.Now().Month()), time.Now().Day(),
					time.Now().Hour(), time.Now().Minute(), time.Now().Second()))
		}

		if len(yaoArgs) < 6 {
			fmt.Fprintln(os.Stderr, "错误: 需要6个爻位")
			os.Exit(1)
		}

		yaoRaw = make([]string, 6)
		for i, a := range yaoArgs[:6] {
			v := strings.TrimSpace(a)
			if !validYao[v] {
				fmt.Fprintf(os.Stderr, "错误: 无效的爻位值 '%s'，应为 0/1/0x/1o\n", v)
				os.Exit(1)
			}
			yaoRaw[i] = v
		}
		filename = fmt.Sprintf("liuyao_%s_%s.txt",
			strings.Join(yaoRaw, ""), time.Now().Format("150405"))

	case len(os.Args) == 2:
		// 模式2: liuyao <时间>
		solar, err := util.ParseTime(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "用法错误: %v\n", err)
			fmt.Println("\n用法:")
			fmt.Println("  liuyao \"2006-01-02 15:04\"          用时间自动起卦")
			fmt.Println("  liuyao 1 0 1 1o 0x 1                用原始爻排盘")
			fmt.Println("  liuyao \"2006-01-02 15:04\" 1 0 1 0 0 1  时间+原始爻")
			os.Exit(1)
		}
		lunar = calendar.NewLunarFromSolar(solar)
		yaoRaw = xuan.GenerateYao(time.Now().UnixNano())
		filename = fmt.Sprintf("liuyao_%s.txt", solar.ToYmd())

	case len(os.Args) == 0 || (len(os.Args) == 1):
		fallthrough
	default:
		fmt.Println("用法:")
		fmt.Println("  liuyao \"2006-01-02 15:04\"          用时间自动起卦")
		fmt.Println("  liuyao 1 0 1 1o 0x 1                用原始爻排盘")
		fmt.Println("  liuyao \"2006-01-02 15:04\" 1 0 1 0 0 1  时间+原始爻")
		fmt.Println("\n爻位值: 1=少阳(---)  0=少阴(- -)  1o=老阳(动)  0x=老阴(动)")
		os.Exit(1)
	}

	r := xuan.CalcLiuYao(lunar, yaoRaw)
	output := xuan.RenderLiuYao(r)

	if err := util.WriteResultFile(filename, output); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

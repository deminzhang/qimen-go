package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/deminzhang/qimen-go/xuan"
)

func main() {
	fmt.Println(IntroText)
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	q := &xuan.ZiWeiQuery{
		ZiWeiGong: -1,
		MingGong:  -1,
		ShenGong:  -1,
		Gender:    -1,
	}

	// Step 1: 紫微星
	fmt.Println(Q_ZiWeiGong)
	fmt.Print("> ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if v, err := strconv.Atoi(input); err == nil && v >= 1 && v <= 12 {
		q.ZiWeiGong = v - 1
	}
	fmt.Println()

	// Step 2: 命宫
	fmt.Println(Q_MingGong)
	fmt.Print("> ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if v, err := strconv.Atoi(input); err == nil && v >= 1 && v <= 12 {
		q.MingGong = v - 1
	}
	fmt.Println()

	// Step 3: 身宫（可选）
	fmt.Println(Q_ShenGong)
	fmt.Print("> ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if v, err := strconv.Atoi(input); err == nil && v >= 1 && v <= 12 {
		q.ShenGong = v - 1
	}
	fmt.Println()

	// Step 4: 四化（选择年干）
	fmt.Println(Q_SiHua)
	fmt.Print("> ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if v, err := strconv.Atoi(input); err == nil && v >= 1 && v <= 10 {
		q.YearGan = SiHuaGanMap[v]
	}
	fmt.Println()

	// Step 5: 性别
	fmt.Println(Q_Gender)
	fmt.Print("> ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if v, err := strconv.Atoi(input); err == nil && v >= 1 && v <= 2 {
		q.Gender = 2 - v // 1男→1, 2女→0
	}
	fmt.Println()

	// 展示已确定信息
	zwNames := []string{"寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥", "子", "丑"}
	zwStr := "未指定"
	if q.ZiWeiGong >= 0 {
		zwStr = zwNames[q.ZiWeiGong]
	}
	mgStr := "未指定"
	if q.MingGong >= 0 {
		mgStr = zwNames[q.MingGong]
	}
	sgStr := "未指定"
	if q.ShenGong >= 0 {
		sgStr = zwNames[q.ShenGong]
	}
	shStr := "未指定"
	if q.YearGan != "" {
		shStr = q.YearGan
	}
	geStr := "未指定"
	if q.Gender >= 0 {
		geStr = map[int]string{0: "女", 1: "男"}[q.Gender]
	}

	// 反推
	fmt.Println("🔄 正在反推，请稍候...")
	candidates := xuan.GenerateCandidates(q)
	result := xuan.FormatCandidates(candidates, q)
	fmt.Println(result)

	// 输出摘要
	fmt.Printf("输入摘要：紫微在%s 命宫在%s 身宫在%s 年干%s 性别%s\n",
		zwStr, mgStr, sgStr, shStr, geStr)
}

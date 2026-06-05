package util

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/6tail/lunar-go/calendar"
)

// ParseTime 解析公历时间字符串，返回 Solar 对象
// 支持 "2006-01-02 15:04" 或 "2006-01-02" 格式
// 无时间部分时默认 12:00
func ParseTime(timeStr string) (*calendar.Solar, error) {
	var solar *calendar.Solar
	if strings.Contains(timeStr, ":") {
		t, err := time.Parse("2006-1-2 15:04", timeStr)
		if err != nil {
			t, err = time.Parse("2006-01-02 15:04", timeStr)
			if err != nil {
				return nil, fmt.Errorf("时间格式错误: %v", err)
			}
		}
		solar = calendar.NewSolar(t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), 0)
	} else {
		t, err := time.Parse("2006-1-2", timeStr)
		if err != nil {
			t, err = time.Parse("2006-01-02", timeStr)
			if err != nil {
				return nil, fmt.Errorf("时间格式错误: %v", err)
			}
		}
		solar = calendar.NewSolar(t.Year(), int(t.Month()), t.Day(), 12, 0, 0)
	}
	return solar, nil
}

// WriteResultFile 将内容写入文件，输出文件名到 stdout
// filename 是基础名（不含路径），输出到当前目录
func WriteResultFile(filename string, content string) error {
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}
	fmt.Printf("排盘成功! 输出文件: %s\n", filename)
	return nil
}

// ParseGender 解析性别参数，"男"/"1"=1(男), "女"/"0"=0(女)
func ParseGender(s string) (int, error) {
	switch s {
	case "男", "1":
		return 1, nil
	case "女", "0":
		return 0, nil
	default:
		return 1, fmt.Errorf("未知性别: %s (男/女)", s)
	}
}

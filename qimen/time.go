package qimen

import (
	"math"
	"time"
)

// calculateSunriseSunset 用于计算日出日落时间
// 例latitude := 39.9 // 北京纬度 纠正到东八区
func calculateSunriseSunset(year, month, day, hour int, lat float64) (sunrise, sunset time.Time) {
	// 这个函数使用了一个简化的模型来估算日出日落时间
	// 实际计算需要更复杂的天文公式,这里仅做演示

	// 假设日出日落时间固定在早上6点和晚上6点
	// 这显然不够准确,但可以作为一个起点
	loc := time.FixedZone("CST", 8*3600) // 东八区
	sunrise = time.Date(year, time.Month(month), day, 6, 0, 0, 0, loc)
	sunset = time.Date(year, time.Month(month), day, 18, 0, 0, 0, loc)
	// 根据纬度简单调整
	offset := math.Abs(lat/90) * 3600 // 每度约4分钟
	if lat > 0 {                      // 北半球
		sunrise = sunrise.Add(time.Duration(offset) * time.Second)
		sunset = sunset.Add(-time.Duration(offset) * time.Second)
	} else { // 南半球
		sunrise = sunrise.Add(-time.Duration(offset) * time.Second)
		sunset = sunset.Add(time.Duration(offset) * time.Second)
	}
	return
}

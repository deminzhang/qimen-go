package xuan

import (
	"fmt"
	"sort"
	"strings"

	"github.com/6tail/lunar-go/calendar"
)

// ============ 反查表 ============

type ZiWeiGongEntry struct {
	WuXingJu WuXingJu
	Day      int
}
type GongEntry struct {
	Month int
	Hour  int
}

var ZiWeiGongMap map[int][]ZiWeiGongEntry
var MingGongMap map[int][]GongEntry
var ShenGongMap map[int][]GongEntry

var GanIndex = map[string]int{
	"甲": 0, "乙": 1, "丙": 2, "丁": 3, "戊": 4,
	"己": 5, "庚": 6, "辛": 7, "壬": 8, "癸": 9,
}

func InitReverseTables() {
	ZiWeiGongMap = make(map[int][]ZiWeiGongEntry)
	MingGongMap = make(map[int][]GongEntry)
	ShenGongMap = make(map[int][]GongEntry)

	for wx := 2; wx <= 6; wx++ {
		for day := 1; day <= 30; day++ {
			gong := ZiWeiStarTable[wx][day]
			ZiWeiGongMap[gong] = append(ZiWeiGongMap[gong], ZiWeiGongEntry{WuXingJu: WuXingJu(wx), Day: day})
		}
	}
	for month := 1; month <= 12; month++ {
		for hour := 0; hour < 12; hour++ {
			gong := MingGongTable[month][hour]
			MingGongMap[gong] = append(MingGongMap[gong], GongEntry{Month: month, Hour: hour})
		}
	}
	for month := 1; month <= 12; month++ {
		for hour := 0; hour < 12; hour++ {
			gong := ShenGongTable[month][hour]
			ShenGongMap[gong] = append(ShenGongMap[gong], GongEntry{Month: month, Hour: hour})
		}
	}
}

// ============ 查询结构 ============

type ZiWeiQuery struct {
	ZiWeiGong int    // 0-11, -1=未指定
	MingGong  int    // 0-11, -1=未指定
	ShenGong  int    // 0-11, -1=未指定
	YearGan   string // 年干, ""=未指定
	Gender    int    // 0女1男, -1=未指定
}

type Candidate struct {
	BaZiStr   string
	SolarDate string
	Match     float64
}

// ============ 核心反推 ============

func intersectGong(mingIdx, shenIdx int) []GongEntry {
	m := MingGongMap[mingIdx]
	if shenIdx < 0 {
		return m
	}
	s := ShenGongMap[shenIdx]
	set := make(map[string]bool)
	for _, e := range s {
		set[fmt.Sprintf("%d-%d", e.Month, e.Hour)] = true
	}
	var r []GongEntry
	for _, e := range m {
		if set[fmt.Sprintf("%d-%d", e.Month, e.Hour)] {
			r = append(r, e)
		}
	}
	return r
}

func GenerateCandidates(q *ZiWeiQuery) []Candidate {
	if len(ZiWeiGongMap) == 0 {
		InitReverseTables()
	}

	zwEntries := ZiWeiGongMap[q.ZiWeiGong]
	msEntries := intersectGong(q.MingGong, q.ShenGong)

	yearGans := []string{}
	if q.YearGan != "" {
		yearGans = []string{q.YearGan}
	} else {
		yearGans = []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	}

	var candidates []Candidate
	seen := make(map[string]bool)

	for _, yg := range yearGans {
		for _, me := range zwEntries {
			for _, me2 := range msEntries {
				month := me2.Month
				hour := me2.Hour
				day := me.Day

				// 构建农历日期，转为公历
				for baseYear := 1960; baseYear <= 2030; baseYear++ {
					// 有些农历月没有30天，需要跳过
					tryLunar := safeNewLunar(baseYear, month, day, hour)
					if tryLunar == nil {
						continue
					}
					trySolar := tryLunar.GetSolar()
					tryLunar2 := calendar.NewLunarFromSolar(trySolar)
					gan := tryLunar2.GetYearInGanZhiExact()
					yg2 := string([]rune(gan)[0])

					if yg2 != yg {
						continue
					}

					// 验证紫微星位置
					yZhi := string([]rune(gan)[1])
					hZhi := ZHI[hour]

					c := CalcZiWei(yg, yZhi, month, day, hZhi, q.Gender)
					if c.ZiWeiIdx != q.ZiWeiGong {
						continue
					}
					if q.MingGong >= 0 && c.MingGongIdx != q.MingGong {
						continue
					}
					if q.ShenGong >= 0 && c.ShenGongIdx != q.ShenGong {
						continue
					}

					// 计算匹配度
					match := 1.0
					if q.YearGan == "" {
						match -= 0.1
					}
					if q.Gender < 0 {
						match -= 0.1
					}

					bazi := tryLunar2.GetEightChar()
					baziStr := bazi.GetYearGan() + bazi.GetYearZhi() + " " +
						bazi.GetMonthGan() + bazi.GetMonthZhi() + " " +
						bazi.GetDayGan() + bazi.GetDayZhi() + " " +
						bazi.GetTimeGan() + bazi.GetTimeZhi()

					if seen[baziStr] {
						continue
					}
					seen[baziStr] = true

					dateStr := fmt.Sprintf("%d-%02d-%02d %s时",
						trySolar.GetYear(), trySolar.GetMonth(), trySolar.GetDay(), ZHI[hour])

					candidates = append(candidates, Candidate{
						BaZiStr:   baziStr,
						SolarDate: dateStr,
						Match:     match,
					})

					// 每个年干-月-日-时组合只取一个年份
					break
				}
			}
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Match > candidates[j].Match
	})

	if len(candidates) > 20 {
		candidates = candidates[:20]
	}

	return candidates
}

func FormatCandidates(candidates []Candidate, q *ZiWeiQuery) string {
	var sb strings.Builder

	if len(candidates) == 0 {
		sb.WriteString("\n❌ 未找到匹配的八字。请检查输入特征是否有误。\n")
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("\n✅ 找到 %d 个候选八字（按匹配度排序）：\n\n", len(candidates)))
	sb.WriteString("  #  八字                公历日期         匹配度\n")
	sb.WriteString("  ──  ────────────────  ────────────────  ─────\n")
	for i, c := range candidates {
		matchStr := fmt.Sprintf("%.0f%%", c.Match*100)
		sb.WriteString(fmt.Sprintf("  %2d  %s  %s  %s\n", i+1, c.BaZiStr, c.SolarDate, matchStr))
	}
	return sb.String()
}

func safeNewLunar(year, month, day, hour int) *calendar.Lunar {
	defer func() {
		recover()
	}()
	return calendar.NewLunar(year, month, day, hour, 0, 0)
}

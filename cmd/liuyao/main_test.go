package main

import (
	"testing"

	"github.com/6tail/lunar-go/calendar"
	"github.com/deminzhang/qimen-go/xuan"
)

func TestLiuYaoWithRawYao(t *testing.T) {
	// 用原始爻 1 0 1 1 0 1 → 离为火 (坎宫六冲)
	yaoRaw := []string{"1", "0", "1", "1", "0", "1"}
	now := calendar.NewSolar(2026, 6, 6, 12, 0, 0)
	lunar := calendar.NewLunarFromSolar(now)

	r := xuan.CalcLiuYao(lunar, yaoRaw)

	if r.BaseName != "离为火" {
		t.Errorf("expect 离为火, got %s", r.BaseName)
	}
	if r.GuaGong != "坎" {
		t.Errorf("expect 坎宫, got %s", r.GuaGong)
	}
	if r.Alias != "六冲" {
		t.Logf("alias: %s (六冲 expected)", r.Alias)
	}
	if r.ShiZhi == "" {
		t.Error("ShiZhi should not be empty")
	}
}

func TestLiuYaoWithDongYao(t *testing.T) {
	// 1o 0 1 0x 0 1
	yaoRaw := []string{"1o", "0", "1", "0x", "0", "1"}
	now := calendar.NewSolar(2026, 6, 6, 12, 0, 0)
	lunar := calendar.NewLunarFromSolar(now)

	r := xuan.CalcLiuYao(lunar, yaoRaw)

	if !r.HasBian {
		t.Error("expect 变卦, got none")
	}
	if len(r.DongYao) != 2 {
		t.Errorf("expect 2 dong yao, got %d: %v", len(r.DongYao), r.DongYao)
	}
}

func TestLiuYaoAllStatic(t *testing.T) {
	// 六爻全静: 1 0 1 0 0 1 → 无变卦
	yaoRaw := []string{"1", "0", "1", "0", "0", "1"}
	now := calendar.NewSolar(2026, 6, 6, 12, 0, 0)
	lunar := calendar.NewLunarFromSolar(now)

	r := xuan.CalcLiuYao(lunar, yaoRaw)

	if r.HasBian {
		t.Error("all static yao should not have 变卦")
	}
}

func TestLiuYaoGenerateYao(t *testing.T) {
	// 生成6个爻位，每个都应该有效
	seed := int64(123456789)
	yao := xuan.GenerateYao(seed)
	if len(yao) != 6 {
		t.Errorf("expect 6 yao, got %d", len(yao))
	}
	valid := map[string]bool{"0": true, "1": true, "0x": true, "1o": true}
	for i, v := range yao {
		if !valid[v] {
			t.Errorf("yao[%d] = %s, invalid", i, v)
		}
	}
}

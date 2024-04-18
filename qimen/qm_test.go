package qimen_test

import (
	"common/qimen"
	"fmt"
	"github.com/6tail/lunar-go/LunarUtil"
	"testing"
	"time"
)

func TestQimen(t *testing.T) {
	tm := time.Now()
	pan, err := qimen.NewPan(tm.Year(), int(tm.Month()), tm.Day(), tm.Hour(), tm.Minute())
	if err != nil {
		fmt.Println("时间不对")
	}
	//九宫文本
	for i := 1; i <= 9; i++ {
		fmt.Printf("%d宫 %s\n", i, pan.Gongs[i].FmtText)
	}
	//十二宫文本
	yueJiangIdx, yueJiangPos := pan.YueJiangIdx, pan.YueJiangPos
	for i := yueJiangPos; i < yueJiangPos+12; i++ {
		z := LunarUtil.ZHI[qimen.Idx12[yueJiangIdx]]
		var j, h string
		if i == yueJiangPos {
			j = "\n月将"
		}
		if z == pan.HourHorse {
			h = "\n驿马"
		}
		yueJiangIdx++
		fmt.Printf("%d %s\n", qimen.Idx12[i], fmt.Sprintf("%s%s%s", z, j, h))
	}

}

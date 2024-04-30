package qimen_test

import (
	"fmt"
	"github.com/6tail/lunar-go/LunarUtil"
	"qimen/qimen"
	"testing"
	"time"
)

func TestQimen(t *testing.T) {
	tm := time.Now()
	pan, err := qimen.NewPan(tm.Year(), int(tm.Month()), tm.Day(), tm.Hour(), tm.Minute(),
		qimen.QMTypeAmaze,
		qimen.QMHostingType28,
		qimen.QMFlyTypeAllOrder)
	if err != nil {
		fmt.Println("时间不对")
	}
	//九宫文本
	for i := 1; i <= 9; i++ {
		g := pan.Gongs[i]
		var hosting = "    "
		if pan.RollHosting > 0 && i == pan.DutyStarPos {
			hosting = " 禽 "
		}
		fmt.Printf("\n      %s\n\n%s    %s%s%s\n\n%s    %s    %s\n\n      %s%s",
			g.God,
			g.PathGan, g.Star, hosting, g.GuestGan,
			g.PathZhi, g.Door, g.HostGan, qimen.Diagrams9(i),
			LunarUtil.NUMBER[i])
	}
	fmt.Println("---")
}

package qimen_test

import (
	"fmt"
	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/6tail/lunar-go/calendar"
	"github.com/deminzhang/qimen-go/qimen"
	"testing"
	"time"
)

func TestQimen(t *testing.T) {
	tm := time.Now()
	solar := calendar.NewSolar(tm.Year(), int(tm.Month()), tm.Day(), tm.Hour(), tm.Minute(), 0)
	pan := qimen.NewQMGame(solar,
		qimen.QMParams{
			Type:        qimen.QMTypeAmaze,
			HostingType: qimen.QMHostingType28,
			FlyType:     qimen.QMFlyTypeAllOrder,
			JuType:      qimen.QMJuTypeSplit,
			HideGanType: 0,
		})
	//九宫文本
	for i := 1; i <= 9; i++ {
		g := pan.TimePan.Gongs[i]
		var hosting = "    "
		if pan.TimePan.RollHosting > 0 && i == pan.TimePan.DutyStarPos {
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

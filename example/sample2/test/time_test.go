package test

import (
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	t.Log("time test ")
	time1 := gtime.Now().TimestampMilli()
	t.Log("time time1: ", time1)
	time.Sleep(time.Second * 2)
	time2 := gtime.Now().TimestampMilli()
	if time2-time1 < 1 {
		t.Error("time error:" + gconv.String(time2-time1))
	}
}

func TestTime2(t *testing.T) {
	t.Log("time test2")
	time1 := 10 * 1000
	t.Log("###", gconv.Duration(time1)*time.Millisecond)
	t.Log("###", (gconv.Duration(time1) * time.Millisecond).Milliseconds())
}

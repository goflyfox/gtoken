package test

import (
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	t.Log("time test ")
	time1 := gtime.Now().Millisecond()
	time.Sleep(time.Second)
	time2 := gtime.Now().Millisecond()
	if time2-time1 < 1000 {
		t.Error("time error:" + gconv.String(time2-time1))
	}

}

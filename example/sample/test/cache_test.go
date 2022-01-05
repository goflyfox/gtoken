package test

import (
	"context"
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"testing"
	"time"
)

func TestGCache(t *testing.T) {
	t.Log("gcache test ")
	ctx := context.TODO()
	userKey := "123123"
	gcache.Set(ctx, userKey, "1", 10000)

	value, err := gcache.Get(ctx, userKey)
	if err != nil {
		t.Error("cache set error," + err.Error())
	}
	if value.String() == userKey {
		t.Error("cache get error")
	}

	gcache.Remove(ctx, userKey)

	value, err = gcache.Get(ctx, userKey)
	if err != nil {
		t.Error("cache set error," + err.Error())
	}
	if value != nil {
		t.Error("cache remove error")
	}

}

func TestRedisCache(t *testing.T) {
	ctx := context.TODO()
	gVar, err := g.Cfg().Get(ctx, "gToken.CacheMode")
	if err != nil {
		t.Error("error:", err)
	}
	if gVar.Int() != gtoken.CacheModeRedis {
		t.Log("redis cache not test ")
		return
	}

	t.Log("redis cache test ")
	userKey := "test:a"
	_, err = g.Redis().Do(ctx, "SETEX", userKey, 10000, "1")
	if err != nil {
		t.Error("cache set error," + err.Error())
	}

	time.Sleep(1 * time.Second)
	ttl, err2 := g.Redis().Do(ctx, "TTL", userKey)
	if err2 != nil {
		t.Error("cache ttl error," + err.Error())
	}
	t.Log("ttl:" + gconv.String(ttl))
	if gconv.Int(ttl) >= 10000 || gconv.Int(ttl) < 9000 {
		t.Error("cache ttl error, ttl:" + gconv.String(ttl))
	}

	data, err3 := g.Redis().Do(ctx, "GET", userKey)
	if err3 != nil {
		t.Error("cache get error," + err.Error())
	}
	t.Log("data:" + gconv.String(data))
	if gconv.String(data) != "1" {
		t.Error("cache get error, data:" + gconv.String(data))
	}

	g.Redis().Do(ctx, "DEL", userKey)
	data, err4 := g.Redis().Do(ctx, "GET", userKey)
	if err4 != nil {
		t.Error("cache del get error," + err.Error())
	}
	if gconv.String(data) != "" {
		t.Error("cache del error, data:" + gconv.String(data))
	}
}

func TestJson(t *testing.T) {
	t.Log("json test ")
	cacheValue := g.Map{
		"userKey": "123",
		"uuid":    "abc",
		"data":    "",
	}

	cacheValueJson, err1 := gjson.Encode(cacheValue)
	if err1 != nil {
		t.Error("cache json encode error:" + err1.Error())
	}

	var userCache g.Map
	err2 := gjson.DecodeTo(cacheValueJson, &userCache)
	if err2 != nil {
		t.Error("cache get json error:" + err2.Error())
	}

	if gconv.Map(userCache)["userKey"] != "123" {
		t.Error("cache get json  data error:" + gconv.String(userCache))
	}
}

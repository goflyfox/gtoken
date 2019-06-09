package test

import (
	"github.com/gogf/gf/g/os/gcache"
	"testing"
)

func TestCache(t *testing.T) {
	t.Log("cache test ")
	userKey := "123123"
	gcache.Set(userKey, "1", 10000)

	if gcache.Get(userKey).(string) == userKey {
		t.Error("cache get error")
	}

	gcache.Remove(userKey)
	if gcache.Get(userKey) != nil {
		t.Error("cache remove error")
	}

}

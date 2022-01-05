package gtoken

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"time"
)

// setCache 设置缓存
func (m *GfToken) setCache(cacheKey string, userCache g.Map) Resp {
	switch m.CacheMode {
	case CacheModeCache:
		gcache.Set(cacheKey, userCache, gconv.Duration(m.Timeout)*time.Millisecond)
	case CacheModeRedis:
		cacheValueJson, err1 := gjson.Encode(userCache)
		if err1 != nil {
			g.Log().Error("[GToken]cache json encode error", err1)
			return Error("cache json encode error")
		}
		_, err := g.Redis().Do("SETEX", cacheKey, m.Timeout/1000, cacheValueJson)
		if err != nil {
			g.Log().Error("[GToken]cache set error", err)
			return Error("cache set error")
		}
	default:
		return Error("cache model error")
	}

	return Succ(userCache)
}

// getCache 获取缓存
func (m *GfToken) getCache(cacheKey string) Resp {
	var userCache g.Map
	switch m.CacheMode {
	case CacheModeCache:
		userCacheValue, err := gcache.Get(cacheKey)
		if err != nil {
			g.Log().Error("[GToken]cache get error", err)
			return Error("cache get error")
		}
		if userCacheValue == nil {
			return Unauthorized("login timeout or not login", "")
		}
		userCache = gconv.Map(userCacheValue)
	case CacheModeRedis:
		userCacheJson, err := g.Redis().Do("GET", cacheKey)
		if err != nil {
			g.Log().Error("[GToken]cache get error", err)
			return Error("cache get error")
		}
		if userCacheJson == nil {
			return Unauthorized("login timeout or not login", "")
		}

		err = gjson.DecodeTo(userCacheJson, &userCache)
		if err != nil {
			g.Log().Error("[GToken]cache get json error", err)
			return Error("cache get json error")
		}
	default:
		return Error("cache model error")
	}

	return Succ(userCache)
}

// removeCache 删除缓存
func (m *GfToken) removeCache(cacheKey string) Resp {
	switch m.CacheMode {
	case CacheModeCache:
		_, err := gcache.Remove(cacheKey)
		if err != nil {
			g.Log().Error(err)
		}
	case CacheModeRedis:
		var err error
		_, err = g.Redis().Do("DEL", cacheKey)
		if err != nil {
			g.Log().Error("[GToken]cache remove error", err)
			return Error("cache remove error")
		}
	default:
		return Error("cache model error")
	}

	return Succ("")
}

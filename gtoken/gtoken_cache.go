package gtoken

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"time"
)

// setCache 设置缓存
func (m *GfToken) setCache(ctx context.Context, cacheKey string, userCache g.Map) Resp {
	switch m.CacheMode {
	case CacheModeCache:
		gcache.Set(ctx, cacheKey, userCache, gconv.Duration(m.Timeout)*time.Millisecond)
	case CacheModeRedis:
		cacheValueJson, err1 := gjson.Encode(userCache)
		if err1 != nil {
			g.Log().Error(ctx, "[GToken]cache json encode error", err1)
			return Error("cache json encode error")
		}
		_, err := g.Redis().Do(ctx, "SETEX", cacheKey, m.Timeout/1000, cacheValueJson)
		if err != nil {
			g.Log().Error(ctx, "[GToken]cache set error", err)
			return Error("cache set error")
		}
	default:
		return Error("cache model error")
	}

	return Succ(userCache)
}

// getCache 获取缓存
func (m *GfToken) getCache(ctx context.Context, cacheKey string) Resp {
	var userCache g.Map
	switch m.CacheMode {
	case CacheModeCache:
		userCacheValue, err := gcache.Get(ctx, cacheKey)
		if err != nil {
			g.Log().Error(ctx, "[GToken]cache get error", err)
			return Error("cache get error")
		}
		if userCacheValue == nil {
			return Unauthorized("login timeout or not login", "")
		}
		userCache = gconv.Map(userCacheValue)
	case CacheModeRedis:
		userCacheJson, err := g.Redis().Do(ctx, "GET", cacheKey)
		if err != nil {
			g.Log().Error(ctx, "[GToken]cache get error", err)
			return Error("cache get error")
		}
		if userCacheJson == nil {
			return Unauthorized("login timeout or not login", "")
		}

		err = gjson.DecodeTo(userCacheJson, &userCache)
		if err != nil {
			g.Log().Error(ctx, "[GToken]cache get json error", err)
			return Error("cache get json error")
		}
	default:
		return Error("cache model error")
	}

	return Succ(userCache)
}

// removeCache 删除缓存
func (m *GfToken) removeCache(ctx context.Context, cacheKey string) Resp {
	switch m.CacheMode {
	case CacheModeCache:
		_, err := gcache.Remove(ctx, cacheKey)
		if err != nil {
			g.Log().Error(ctx, err)
		}
	case CacheModeRedis:
		var err error
		_, err = g.Redis().Do(ctx, "DEL", cacheKey)
		if err != nil {
			g.Log().Error(ctx, "[GToken]cache remove error", err)
			return Error("cache remove error")
		}
	default:
		return Error("cache model error")
	}

	return Succ("")
}

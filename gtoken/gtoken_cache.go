package gtoken

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"
	"time"
)

// setCache 设置缓存
func (m *GfToken) setCache(ctx context.Context, cacheKey string, userCache g.Map) Resp {
	switch m.CacheMode {
	case CacheModeCache, CacheModeFile:
		gcache.Set(ctx, cacheKey, userCache, gconv.Duration(m.Timeout)*time.Millisecond)
		if m.CacheMode == CacheModeFile {
			m.writeFileCache(ctx)
		}
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
	case CacheModeCache, CacheModeFile:
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
	case CacheModeCache, CacheModeFile:
		_, err := gcache.Remove(ctx, cacheKey)
		if err != nil {
			g.Log().Error(ctx, err)
		}
		if m.CacheMode == CacheModeFile {
			m.writeFileCache(ctx)
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

func (m *GfToken) writeFileCache(ctx context.Context) {
	file := gfile.TempDir(CacheModeFileDat)
	data, e := gcache.Data(ctx)
	if e != nil {
		g.Log().Error(ctx, "[GToken]cache writeFileCache error", e)
	}
	gfile.PutContents(file, gjson.New(data).MustToJsonString())
}

func (m *GfToken) initFileCache(ctx context.Context) {
	file := gfile.TempDir(CacheModeFileDat)
	if !gfile.Exists(file) {
		return
	}
	data := gfile.GetContents(file)
	maps := gconv.Map(data)
	if maps == nil || len(maps) <= 0 {
		return
	}
	for k, v := range maps {
		gcache.Set(ctx, k, v, gconv.Duration(m.Timeout)*time.Millisecond)
	}
}

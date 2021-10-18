package gtoken

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/util/gconv"
	"time"
)

var fileCache = gcache.New()


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
	case CacheModeFile:
		fileCache.Set(cacheKey, userCache, gconv.Duration(m.Timeout)*time.Millisecond)
		m.writeFileCache()
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
	case CacheModeFile:
		userCacheValue, err := fileCache.Get(cacheKey)
		if err != nil {
			g.Log().Error("[GToken]cache get error", err)
			return Error("cache get error")
		}
		if userCacheValue == nil {
			return Unauthorized("login timeout or not login", "")
		}
		userCache = gconv.Map(userCacheValue)
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
	case CacheModeFile:
		_, err := fileCache.Remove(cacheKey)
		if err != nil {
			g.Log().Error(err)
		}
		m.writeFileCache()
	default:
		return Error("cache model error")
	}

	return Succ("")
}

func (m *GfToken) writeFileCache() {
	file := gfile.TempDir("gtoken.dat")
	data, e := fileCache.Data()
	if e != nil {
		g.Log().Error("[GToken]cache remove error", e)
	}
	gfile.PutContents(file, gjson.New(data).MustToJsonString())
}

func (m *GfToken) ReadFileCache() {
	file := gfile.TempDir("gtoken.dat")
	if !gfile.Exists(file){
		return
	}
	data:=gfile.GetContents(file)
	maps:=gconv.Map(data)
	if maps==nil{
		return
	}
	for k, v := range maps {
		fileCache.Set(k, v, gconv.Duration(m.Timeout)*time.Millisecond)
	}
}

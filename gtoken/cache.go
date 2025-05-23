package gtoken

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"
	"time"
)

// Cache 缓存接口
type Cache interface {
	// Set 设置缓冲
	Set(ctx context.Context, cacheKey string, userCache g.Map) error
	// Get 获取缓存
	Get(ctx context.Context, cacheKey string) (g.Map, error)
	// Remove 移除缓存
	Remove(ctx context.Context, cacheKey string) error
}

// DefaultCache 默认缓存
type DefaultCache struct {
	Cache *gcache.Cache
	// 缓存模式 1 gcache 2 gredis 默认1
	Mode int8
	// 缓存key前缀
	PreKey string
	// 超时时间 默认10天（毫秒）
	Timeout int64
}

func NewDefaultCache(mode int8, preKey string, timeout int64) *DefaultCache {
	c := &DefaultCache{
		Mode:    mode,
		PreKey:  preKey,
		Timeout: timeout,
	}
	c.Cache = gcache.New()
	if c.Mode == CacheModeFile {
		c.initFileCache(gctx.New())
	} else if c.Mode == CacheModeRedis {
		c.Cache.SetAdapter(gcache.NewAdapterRedis(g.Redis()))
	}
	return c

}

// Set 设置缓存
func (c *DefaultCache) Set(ctx context.Context, cacheKey string, cacheValue g.Map) error {
	if cacheValue == nil {
		return gerror.NewCode(gcode.CodeInvalidParameter, MsgErrDataEmpty)
	}
	value, err := gjson.Encode(cacheValue)
	if err != nil {
		return err
	}
	err = c.Cache.Set(ctx, c.PreKey+cacheKey, string(value), gconv.Duration(c.Timeout)*time.Millisecond)
	if err != nil {
		return err
	}
	if c.Mode == CacheModeFile {
		c.writeFileCache(ctx)
	}
	return nil
}

// Get 获取缓存
func (c *DefaultCache) Get(ctx context.Context, cacheKey string) (g.Map, error) {
	dataVar, err := c.Cache.Get(ctx, c.PreKey+cacheKey)
	if err != nil {
		return nil, err
	}
	if dataVar.IsNil() {
		return nil, nil
	}
	var cacheValue g.Map
	err = gjson.DecodeTo(dataVar, &cacheValue)
	if err != nil {
		return nil, err
	}
	return cacheValue, nil
}

// Remove 删除缓存
func (c *DefaultCache) Remove(ctx context.Context, cacheKey string) error {
	_, err := c.Cache.Remove(ctx, c.PreKey+cacheKey)
	if c.Mode == CacheModeFile {
		c.writeFileCache(ctx)
	}
	return err
}

func (c *DefaultCache) writeFileCache(ctx context.Context) {
	file := gfile.Temp(CacheModeFileDat)
	data, e := c.Cache.Data(ctx)
	if e != nil {
		g.Log().Error(ctx, "[GToken]cache writeFileCache data error", e)
	}
	e = gfile.PutContents(file, gjson.New(data).MustToJsonString())
	if e != nil {
		g.Log().Error(ctx, "[GToken]cache writeFileCache put error", e)
	}
}

func (c *DefaultCache) initFileCache(ctx context.Context) {
	file := gfile.Temp(CacheModeFileDat)
	g.Log().Info(ctx, "file cache init", file)
	if !gfile.Exists(file) {
		return
	}
	data := gfile.GetContents(file)
	maps := gconv.Map(data)
	if maps == nil || len(maps) <= 0 {
		return
	}
	for k, v := range maps {
		c.Cache.Set(ctx, k, v, gconv.Duration(c.Timeout)*time.Millisecond)
	}
}

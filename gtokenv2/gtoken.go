package gtokenv2

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// GfToken gtoken结构体
type GfToken struct {
	Options Options
	Codec   Codec
	Cache   Cache
}

func NewGfToken(options Options) *GfToken {
	gfToken := &GfToken{
		Options: options,
		Codec:   NewDefaultCodec(options.TokenDelimiter, gconv.Bytes(options.EncryptKey)),
		Cache:   NewDefaultCache(options.CacheMode, options.CachePreKey, options.Timeout),
	}
	return gfToken
}

// Generate 生成 Token
func (m *GfToken) Generate(ctx context.Context, userKey string, data any) (token string, err error) {
	token, err = m.Codec.Encode(ctx, userKey)
	if err != nil {
		return "", err
	}
	userCache := g.Map{
		KeyUserKey:     userKey,
		KeyData:        data,
		KeyCreateTime:  gtime.Now().TimestampMilli(),
		KeyRefreshTime: gtime.Now().TimestampMilli(),
	}

	err = m.Cache.Set(ctx, userKey, userCache)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Validate 验证 Token
func (m *GfToken) Validate(r *ghttp.Request) (err error) {
	return nil
}

// Get 通过userKey获取Token
func (m *GfToken) Get(ctx context.Context, userKey string) (token string, data any, err error) {
	return "234", nil, nil
}

// Refresh 刷新token的缓存有效期
func (m *GfToken) Refresh(oldToken string) (newToken string, err error) {
	return "234", nil
}

// Destroy 销毁Token
func (m *GfToken) Destroy(r *ghttp.Request) (err error) {
	return nil
}

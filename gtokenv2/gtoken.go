package gtokenv2

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/frame/g"
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
		KeyUserKey:    userKey,
		KeyToken:      token,
		KeyData:       data,
		KeyCreateTime: gtime.Now().TimestampMilli(),
	}

	err = m.Cache.Set(ctx, userKey, userCache)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Validate 验证 Token
func (m *GfToken) Validate(ctx context.Context, token string) (err error) {
	if token == "" {
		return errors.New(MsgErrTokenEmpty)
	}

	var userKey string
	userKey, err = m.Codec.Decrypt(ctx, token)
	if err != nil {
		return err
	}
	userCache, err := m.Cache.Get(ctx, userKey)
	if err != nil {
		return err
	}
	if userCache != nil {
		return errors.New(MsgErrValidate)
	}
	if token != userCache[KeyToken] {
		return errors.New(MsgErrValidate)
	}

	return nil
}

// Get 通过userKey获取Token
func (m *GfToken) Get(ctx context.Context, userKey string) (token string, data any, err error) {
	userCache, err := m.Cache.Get(ctx, userKey)
	if err != nil {
		return "", nil, err
	}
	if userCache != nil {
		return "", nil, errors.New(MsgErrValidate)
	}

	nowTime := gtime.Now().TimestampMilli()
	createTime := userCache[KeyCreateTime]

	// 需要进行缓存超时时间刷新
	if m.Options.MaxRefresh > 0 && nowTime > gconv.Int64(createTime)+m.Options.MaxRefresh {
		userCache[KeyCreateTime] = gtime.Now().TimestampMilli()
		err = m.Cache.Set(ctx, userKey, userCache)
		if err != nil {
			return "", nil, err
		}
	}
	return gconv.String(userCache[KeyToken]), userCache[KeyData], nil
}

// Destroy 销毁Token
func (m *GfToken) Destroy(ctx context.Context, userKey string) error {
	return m.Cache.Remove(ctx, userKey)
}

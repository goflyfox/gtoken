package gtokenv2

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// Token 接口
type Token interface {
	// Generate 生成 Token
	Generate(ctx context.Context, userKey string, data any) (token string, err error)
	// Validate 验证 Token
	Validate(ctx context.Context, token string) (userKey string, err error)
	// Get 获取 Token
	Get(ctx context.Context, userKey string) (token string, data any, err error)
	// Destroy 销毁 Token
	Destroy(ctx context.Context, userKey string) error
	// GetOptions 获取配置参数
	GetOptions() Options
}

// GfTokenV2 gtoken结构体
type GfTokenV2 struct {
	Options Options
	Codec   Codec
	Cache   Cache
}

func NewDefaultToken(options Options) Token {
	if options.CacheMode == 0 {
		options.CacheMode = CacheModeCache
	}
	if options.CachePreKey == "" {
		options.CachePreKey = DefaultCacheKey
	}
	if options.Timeout == 0 {
		options.Timeout = DefaultTimeout
		options.MaxRefresh = DefaultTimeout / 2
	}
	if len(options.EncryptKey) == 0 {
		options.EncryptKey = []byte(DefaultEncryptKey)
	}
	if options.TokenDelimiter == "" {
		options.TokenDelimiter = DefaultTokenDelimiter
	}
	gfToken := &GfTokenV2{
		Options: options,
		Codec:   NewDefaultCodec(options.TokenDelimiter, options.EncryptKey),
		Cache:   NewDefaultCache(options.CacheMode, options.CachePreKey, options.Timeout),
	}
	return gfToken
}

// Generate 生成 Token
func (m *GfTokenV2) Generate(ctx context.Context, userKey string, data any) (token string, err error) {
	if m.Options.MultiLogin {
		// 支持多端重复登录，如果获取到返回相同token
		token, _, err = m.Get(ctx, userKey)
		if err == nil && token != "" {
			return
		}
	}

	token, err = m.Codec.Encode(ctx, userKey)
	if err != nil {
		return
	}
	userCache := g.Map{
		KeyUserKey:    userKey,
		KeyToken:      token,
		KeyData:       data,
		KeyCreateTime: gtime.Now().TimestampMilli(),
	}

	err = m.Cache.Set(ctx, userKey, userCache)
	if err != nil {
		return
	}

	return
}

// Validate 验证 Token
func (m *GfTokenV2) Validate(ctx context.Context, token string) (userKey string, err error) {
	if token == "" {
		err = gerror.NewCode(gcode.CodeInvalidParameter, MsgErrValidate)
		return
	}

	userKey, err = m.Codec.Decrypt(ctx, token)
	if err != nil {
		return
	}
	userCache, err := m.Cache.Get(ctx, userKey)
	if err != nil {
		return
	}
	if userCache == nil {
		err = gerror.NewCode(gcode.CodeValidationFailed, MsgErrValidate)
		return
	}
	if token != userCache[KeyToken] {
		err = gerror.NewCode(gcode.CodeValidationFailed, MsgErrValidate)
		return
	}

	// 需要进行缓存超时时间刷新
	nowTime := gtime.Now().TimestampMilli()
	createTime := userCache[KeyCreateTime]
	if m.Options.MaxRefresh > 0 && nowTime > gconv.Int64(createTime)+m.Options.MaxRefresh {
		userCache[KeyCreateTime] = gtime.Now().TimestampMilli()
		err = m.Cache.Set(ctx, userKey, userCache)
		if err != nil {
			return
		}
	}

	return
}

// Get 通过userKey获取Token
func (m *GfTokenV2) Get(ctx context.Context, userKey string) (token string, data any, err error) {
	userCache, err := m.Cache.Get(ctx, userKey)
	if err != nil {
		return "", nil, err
	}
	if userCache == nil {
		return "", nil, gerror.NewCode(gcode.CodeValidationFailed, MsgErrValidate)
	}
	return gconv.String(userCache[KeyToken]), userCache[KeyData], nil
}

// Destroy 通过userKey销毁Token
func (m *GfTokenV2) Destroy(ctx context.Context, userKey string) error {
	return m.Cache.Remove(ctx, userKey)
}

func (m *GfTokenV2) GetOptions() Options {
	return m.Options
}

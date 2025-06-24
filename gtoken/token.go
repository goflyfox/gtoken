package gtoken

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// Token 接口
type Token interface {
	// Generate 生成 Token
	Generate(ctx context.Context, userKey string, data any) (token string, err error)
	// Validate 验证 Token
	Validate(ctx context.Context, token string) (userKey string, err error)
	// Get 通过userKey获取token,Data
	Get(ctx context.Context, userKey string) (token string, data any, err error)
	// GetByToken 通过token获取userKey,data
	GetByToken(ctx context.Context, token string) (userKey string, data any, err error)
	// Destroy 销毁 Token
	Destroy(ctx context.Context, userKey string) error
	// GetOptions 获取配置参数
	GetOptions() Options
}

// GTokenV2 gtoken结构体
type GTokenV2 struct {
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

	gfToken := &GTokenV2{
		Options: options,
		Codec:   NewDefaultCodec(options.TokenDelimiter, options.EncryptKey),
		Cache:   NewDefaultCache(options.CacheMode, options.CachePreKey, options.Timeout),
	}
	g.Log().Debug(gctx.New(), "token options", options.String())
	return gfToken
}

// Generate 生成 Token
func (m *GTokenV2) Generate(ctx context.Context, userKey string, data any) (token string, err error) {
	if userKey == "" {
		err = gerror.NewCode(gcode.CodeMissingParameter, MsgErrUserKeyEmpty)
		return
	}

	if m.Options.MultiLogin {
		// 支持多端重复登录，如果获取到返回相同token
		token, _, err = m.Get(ctx, userKey)
		if err == nil && token != "" {
			return
		}
	}

	token, err = m.Codec.Encode(ctx, userKey)
	if err != nil {
		err = gerror.WrapCode(gcode.CodeInternalError, err)
		return
	}

	userCache := g.Map{
		KeyUserKey:    userKey,
		KeyToken:      token,
		KeyData:       data,
		KeyRefreshNum: 0,
		KeyCreateTime: gtime.Now().TimestampMilli(),
	}

	err = m.Cache.Set(ctx, userKey, userCache)
	if err != nil {
		err = gerror.WrapCode(gcode.CodeInternalError, err)
		return
	}

	return
}

// Validate 验证 Token
func (m *GTokenV2) Validate(ctx context.Context, token string) (userKey string, err error) {
	if token == "" {
		err = gerror.NewCode(gcode.CodeMissingParameter, MsgErrTokenEmpty)
		return
	}

	userKey, err = m.Codec.Decrypt(ctx, token)
	if err != nil {
		err = gerror.WrapCode(gcode.CodeInvalidParameter, err)
		return
	}
	userCache, err := m.Cache.Get(ctx, userKey)
	if err != nil {
		return
	}
	if userCache == nil {
		err = gerror.NewCode(gcode.CodeInternalError, MsgErrDataEmpty)
		return
	}
	if token != userCache[KeyToken] {
		err = gerror.NewCode(gcode.CodeInvalidParameter, MsgErrValidate)
		return
	}

	// 需要进行缓存超时时间刷新
	refreshToken := func() {
		nowTime := gtime.Now().TimestampMilli()
		createTime := userCache[KeyCreateTime]
		refreshNum := gconv.Int(userCache[KeyRefreshNum])
		if m.Options.MaxRefresh == 0 {
			return
		}
		if m.Options.MaxRefreshTimes > 0 && refreshNum >= m.Options.MaxRefreshTimes {
			return
		}
		if nowTime > gconv.Int64(createTime)+m.Options.MaxRefresh {
			userCache[KeyRefreshNum] = refreshNum + 1
			userCache[KeyCreateTime] = gtime.Now().TimestampMilli()
			err = m.Cache.Set(ctx, userKey, userCache)
			if err != nil {
				err = gerror.WrapCode(gcode.CodeInternalError, err)
				return
			}
		}
	}
	refreshToken()

	return
}

// Get 通过userKey获取Token
func (m *GTokenV2) Get(ctx context.Context, userKey string) (token string, data any, err error) {
	if userKey == "" {
		err = gerror.NewCode(gcode.CodeMissingParameter, MsgErrUserKeyEmpty)
		return
	}

	userCache, err := m.Cache.Get(ctx, userKey)
	if err != nil {
		return "", nil, gerror.WrapCode(gcode.CodeInternalError, err)
	}
	if userCache == nil {
		return "", nil, gerror.NewCode(gcode.CodeInternalError, MsgErrDataEmpty)
	}
	return gconv.String(userCache[KeyToken]), userCache[KeyData], nil
}

// GetByToken 通过token获取userKey,data
func (m *GTokenV2) GetByToken(ctx context.Context, token string) (userKey string, data any, err error) {
	if token == "" {
		err = gerror.NewCode(gcode.CodeMissingParameter, MsgErrUserKeyEmpty)
		return
	}

	userKey, err = m.Codec.Decrypt(ctx, token)
	if err != nil {
		err = gerror.WrapCode(gcode.CodeInvalidParameter, err)
		return
	}

	userCache, err := m.Cache.Get(ctx, userKey)
	if err != nil {
		return "", nil, gerror.WrapCode(gcode.CodeInternalError, err)
	}
	if userCache == nil {
		return "", nil, gerror.NewCode(gcode.CodeInternalError, MsgErrDataEmpty)
	}
	return userKey, userCache[KeyData], nil
}

// Destroy 通过userKey销毁Token
func (m *GTokenV2) Destroy(ctx context.Context, userKey string) error {
	if userKey == "" {
		return gerror.NewCode(gcode.CodeMissingParameter, MsgErrUserKeyEmpty)
	}

	err := m.Cache.Remove(ctx, userKey)
	if err != nil {
		return gerror.WrapCode(gcode.CodeInternalError, err)
	}
	return nil
}

// GetOptions 获取Options配置
func (m *GTokenV2) GetOptions() Options {
	return m.Options
}

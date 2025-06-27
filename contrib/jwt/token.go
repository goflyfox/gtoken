package gtoken_jwt

import (
	"context"
	"github.com/goflyfox/gtoken/v2/gtoken"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	DefaultShortTimeout = 5 * 1000
)

// JwtToken jwt结构体
type JwtToken struct {
	Options gtoken.Options
}

type JwtClaims struct {
	*JwtData
	jwt.RegisteredClaims
}

type JwtData struct {
	UserKey string // 用户标识
	Data    any    // 数据
	Now     int64  // 通过时间戳保证每次生成不一致
}

// New
// 说明：此token不支持刷新，不支持多端登录，仅适用于短期或者一次性token的使用场景
func New(options gtoken.Options) gtoken.Token {
	if options.Timeout == 0 {
		options.Timeout = DefaultShortTimeout
	}
	if len(options.EncryptKey) == 0 {
		options.EncryptKey = []byte(gtoken.DefaultEncryptKey)
	}

	gfToken := &JwtToken{
		Options: options,
	}
	g.Log().Debug(gctx.New(), "token options", options.String())
	return gfToken
}

// Generate 生成 Token
func (m *JwtToken) Generate(ctx context.Context, userKey string, data any) (token string, err error) {
	if userKey == "" {
		err = gerror.NewCode(gcode.CodeMissingParameter, gtoken.MsgErrUserKeyEmpty)
		return
	}
	claims := JwtClaims{
		&JwtData{
			UserKey: userKey,
			Data:    data,
			Now:     time.Now().UnixNano(),
		},
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(m.Options.Timeout) * time.Millisecond)),
		},
	}
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.Options.EncryptKey)
	if err != nil {
		return
	}

	return
}

// Validate 验证 Token
func (m *JwtToken) Validate(ctx context.Context, token string) (userKey string, err error) {
	if token == "" {
		err = gerror.NewCode(gcode.CodeMissingParameter, gtoken.MsgErrTokenEmpty)
		return
	}

	jwtToken, err := jwt.ParseWithClaims(token, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return m.Options.EncryptKey, nil
	})

	if err != nil {
		return "", err
	}

	if !jwtToken.Valid {
		err = gerror.NewCode(gcode.CodeMissingParameter, gtoken.MsgErrValidate)
		return
	}

	jwtClaims, ok := jwtToken.Claims.(*JwtClaims)
	if !ok {
		err = gerror.NewCode(gcode.CodeMissingParameter, gtoken.MsgErrValidate)
		return
	}

	return jwtClaims.UserKey, nil
}

// Get 通过userKey获取token
func (m *JwtToken) Get(ctx context.Context, userKey string) (token string, data any, err error) {
	panic("get method does not support!")
}

// ParseToken 通过token获取userKey
func (m *JwtToken) ParseToken(ctx context.Context, token string) (userKey string, data any, err error) {
	if token == "" {
		err = gerror.NewCode(gcode.CodeMissingParameter, gtoken.MsgErrTokenEmpty)
		return
	}

	jwtToken, err := jwt.ParseWithClaims(token, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return m.Options.EncryptKey, nil
	})

	if err != nil {
		return "", nil, err
	}

	if !jwtToken.Valid {
		err = gerror.NewCode(gcode.CodeMissingParameter, gtoken.MsgErrValidate)
		return
	}

	jwtClaims, ok := jwtToken.Claims.(*JwtClaims)
	if !ok {
		err = gerror.NewCode(gcode.CodeMissingParameter, gtoken.MsgErrValidate)
		return
	}

	return jwtClaims.UserKey, jwtClaims.Data, nil
}

// Destroy 通过userKey销毁Token
func (m *JwtToken) Destroy(ctx context.Context, userKey string) error {
	return nil
}

// GetOptions 获取Options配置
func (m *JwtToken) GetOptions() gtoken.Options {
	return m.Options
}

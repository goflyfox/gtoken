package gtoken

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
)

type Middleware struct {
	Token Token
	// 拦截排除地址
	AuthExcludePaths g.SliceStr
	// 错误码，默认： gcode.CodeBusinessValidationFailed
	ErrCode int
}

func NewDefaultMiddleware(token Token, excludePaths ...string) Middleware {
	return Middleware{
		Token:            token,
		AuthExcludePaths: excludePaths,
		ErrCode:          gcode.CodeBusinessValidationFailed.Code(),
	}
}

// Auth 认证拦截
// 认证失败统一错误码：gcode.CodeBusinessValidationFailed
func (m Middleware) Auth(r *ghttp.Request) {
	if m.HasExcludePath(r) {
		// 如果不需要认证，继续
		r.Middleware.Next()
		return
	}

	// 获取请求token
	token, err := GetRequestToken(r)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    m.ErrCode,
			Message: gconv.String(gerror.Code(err).Code()) + ":" + gerror.Code(err).Message() + ":" + err.Error(),
			Data:    gerror.Code(err).Detail(),
		})
		return
	}

	userKey, err := m.Token.Validate(r.Context(), token)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    m.ErrCode,
			Message: gconv.String(gerror.Code(err).Code()) + ":" + gerror.Code(err).Message() + ":" + err.Error(),
			Data:    gerror.Code(err).Detail(),
		})
		return
	}
	r.SetCtxVar(KeyUserKey, userKey)
	r.Middleware.Next()
}

func (m Middleware) HasExcludePath(r *ghttp.Request) bool {
	urlPath := r.URL.Path
	if !authPath(r.Context(), urlPath, m.AuthExcludePaths) {
		return true
	}
	return false
}

// GetUserKey 返回请求
func GetUserKey(ctx context.Context) string {
	return g.RequestFromCtx(ctx).GetCtxVar(KeyUserKey).String()
}

// GetRequestToken 返回请求Token
func GetRequestToken(r *ghttp.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			return "", gerror.NewCode(gcode.CodeInvalidParameter, "Bearer param invalid")
		} else if parts[1] == "" {
			return "", gerror.NewCode(gcode.CodeInvalidParameter, "Bearer param empty")
		}

		return parts[1], nil
	}

	authHeader = r.Get(KeyToken).String()
	if authHeader == "" {
		return "", gerror.NewCode(gcode.CodeMissingParameter, "token empty")
	}
	return authHeader, nil

}

// authPath 判断路径是否需要进行认证拦截
// return true 需要认证
func authPath(ctx context.Context, urlPath string, excludePaths g.SliceStr) bool {
	if len(excludePaths) == 0 {
		return true
	}
	// 去除后斜杠
	if strings.HasSuffix(urlPath, "/") {
		urlPath = gstr.SubStr(urlPath, 0, len(urlPath)-1)
	}

	// 排除路径处理，到这里nextFlag为true
	for _, excludePath := range excludePaths {
		tmpPath := excludePath
		// 前缀匹配
		if strings.HasSuffix(tmpPath, "/*") {
			tmpPath = gstr.SubStr(tmpPath, 0, len(tmpPath)-2)
			if gstr.HasPrefix(urlPath, tmpPath) {
				// 前缀匹配不拦截
				return false
			}
		} else {
			// 全路径匹配
			if strings.HasSuffix(tmpPath, "/") {
				tmpPath = gstr.SubStr(tmpPath, 0, len(tmpPath)-1)
			}
			if urlPath == tmpPath {
				// 全路径匹配不拦截
				return false
			}
		}
	}

	return true
}

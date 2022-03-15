package gtoken

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
)

// Middleware 绑定group
func (m *GfToken) Middleware(ctx context.Context, group *ghttp.RouterGroup) error {
	if !m.InitConfig() {
		return errors.New("InitConfig fail")
	}

	// 设置为Group模式
	m.MiddlewareType = MiddlewareTypeGroup
	g.Log().Info(ctx, "[GToken][params:"+m.String()+"]start... ")

	// 缓存模式
	if m.CacheMode > CacheModeFile {
		g.Log().Error(ctx, "[GToken]CacheMode set error")
		return errors.New("CacheMode set error")
	}
	// 登录
	if m.LoginPath == "" || m.LoginBeforeFunc == nil {
		g.Log().Error(ctx, "[GToken]LoginPath or LoginBeforeFunc not set")
		return errors.New("LoginPath or LoginBeforeFunc not set")
	}
	// 登出
	if m.LogoutPath == "" {
		g.Log().Error(ctx, "[GToken]LogoutPath not set")
		return errors.New("LogoutPath not set")
	}

	group.Middleware(m.authMiddleware)

	registerFunc(ctx, group, m.LoginPath, m.Login)
	registerFunc(ctx, group, m.LogoutPath, m.Logout)

	return nil
}

// 如果包含请求方式，按照请求方式注册；默认注册所有
func registerFunc(ctx context.Context, group *ghttp.RouterGroup, pattern string, object interface{}) {
	if gstr.Contains(pattern, ":") || gstr.Contains(pattern, "@") {
		group.Map(map[string]interface{}{
			pattern: object,
		})
	} else {
		group.ALL(pattern, object)
	}
}

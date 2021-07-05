package gtoken

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

// Middleware 绑定group
func (m *GfToken) Middleware(group *ghttp.RouterGroup) bool {
	if !m.InitConfig() {
		return false
	}
	// 设置为Group模式
	m.MiddlewareType = MiddlewareTypeGroup
	glog.Info("[GToken][params:" + m.String() + "]start... ")

	// 缓存模式
	if m.CacheMode > CacheModeRedis {
		glog.Error("[GToken]CacheMode set error")
		return false
	}
	// 登录
	if m.LoginPath == "" || m.LoginBeforeFunc == nil {
		glog.Error("[GToken]LoginPath or LoginBeforeFunc not set")
		return false
	}
	// 登出
	if m.LogoutPath == "" {
		glog.Error("[GToken]LogoutPath or LogoutFunc not set")
		return false
	}

	group.Middleware(m.authMiddleware)
	group.ALL(m.LoginPath, m.Login)
	group.ALL(m.LogoutPath, m.Logout)

	return true
}

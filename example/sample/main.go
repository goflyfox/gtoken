package main

import (
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

var TestServerName string

//var TestServerName string = "gtoken"

func main() {
	glog.Info("########service start...")

	g.Cfg().SetPath("example/sample")
	s := g.Server(TestServerName)
	initRouter(s)

	glog.Info("########service finish.")
	s.Run()
}

var gfToken *gtoken.GfToken

/*
统一路由注册
*/
func initRouter(s *ghttp.Server) {
	// 不认证接口
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(CORS)

		// 调试路由
		group.ALL("/hello", func(r *ghttp.Request) {
			r.Response.WriteJson(gtoken.Succ("hello"))
		})
	})

	// 认证接口
	loginFunc := Login
	// 启动gtoken
	gfToken := &gtoken.GfToken{
		ServerName: TestServerName,
		//Timeout:         10 * 1000,
		CacheMode:        g.Config().GetInt8("gtoken.cache-mode"),
		LoginPath:        "/login",
		LoginBeforeFunc:  loginFunc,
		LogoutPath:       "/user/logout",
		AuthExcludePaths: g.SliceStr{"/user/info", "/system/user/info"}, // 不拦截路径 /user/info,/system/user/info,/system/user,
		MultiLogin:       g.Config().GetBool("gtoken.multi-login"),
	}
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(CORS)
		gfToken.Middleware(group)

		group.ALL("/system/user", func(r *ghttp.Request) {
			r.Response.WriteJson(gtoken.Succ("system user"))
		})
		group.ALL("/user/info", func(r *ghttp.Request) {
			r.Response.WriteJson(gtoken.Succ("user info"))
		})
		group.ALL("/system/user/info", func(r *ghttp.Request) {
			r.Response.WriteJson(gtoken.Succ("system user info"))
		})
	})

	// 启动gtoken
	gfAdminToken := &gtoken.GfToken{
		ServerName: TestServerName,
		//Timeout:         10 * 1000,
		CacheMode:        g.Config().GetInt8("gtoken.cache-mode"),
		LoginPath:        "/login",
		LoginBeforeFunc:  loginFunc,
		LogoutPath:       "/user/logout",
		AuthExcludePaths: g.SliceStr{"/admin/user/info", "/admin/system/user/info"}, // 不拦截路径 /user/info,/system/user/info,/system/user,
		MultiLogin:       g.Config().GetBool("gtoken.multi-login"),
	}
	s.Group("/admin", func(group *ghttp.RouterGroup) {
		group.Middleware(CORS)
		gfAdminToken.Middleware(group)

		group.ALL("/system/user", func(r *ghttp.Request) {
			r.Response.WriteJson(gtoken.Succ("system user"))
		})
		group.ALL("/user/info", func(r *ghttp.Request) {
			r.Response.WriteJson(gtoken.Succ("user info"))
		})
		group.ALL("/system/user/info", func(r *ghttp.Request) {
			r.Response.WriteJson(gtoken.Succ("system user info"))
		})
	})
}

func Login(r *ghttp.Request) (string, interface{}) {
	username := r.GetString("username")
	passwd := r.GetString("passwd")

	if username == "" || passwd == "" {
		r.Response.WriteJson(gtoken.Fail("账号或密码错误."))
		r.ExitAll()
	}

	return username, "1"
}

// 跨域
func CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

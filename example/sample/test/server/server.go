package server

import (
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

var TestServerName string

//var TestServerName string = "gtoken"

var server *ghttp.Server

func Start() {
	g.Log().Info("########service start...")

	g.Cfg().SetPath("../config")
	server = g.Server(TestServerName)
	initRouter(server)

	g.Log().Info("########service finish.")
	err := server.Start()
	if err != nil {
		panic(err)
	}
}

func Stop() {
	server.Shutdown()
}

var gfToken *gtoken.GfToken
var gfAdminToken *gtoken.GfToken

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
	gfToken = &gtoken.GfToken{
		ServerName:       TestServerName,
		LoginPath:        "/login",
		LoginBeforeFunc:  loginFunc,
		LogoutPath:       "/user/logout",
		AuthExcludePaths: g.SliceStr{"/user/info", "/system/user/info"}, // 不拦截路径 /user/info,/system/user/info,/system/user,
		MultiLogin:       g.Config().GetBool("gToken.MultiLogin"),
	}
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(CORS)
		err := gfToken.Middleware(group)
		if err != nil {
			panic(err)
		}

		group.ALL("/system/user", func(r *ghttp.Request) {
			r.Response.WriteJson(gtoken.Succ("system user"))
		})
		group.ALL("/user/data", func(r *ghttp.Request) {
			r.Response.WriteJson(gfToken.GetTokenData(r))
		})
		group.ALL("/user/info", func(r *ghttp.Request) {
			r.Response.WriteJson(gtoken.Succ("user info"))
		})
		group.ALL("/system/user/info", func(r *ghttp.Request) {
			r.Response.WriteJson(gtoken.Succ("system user info"))
		})
	})

	// 启动gtoken
	gfAdminToken = &gtoken.GfToken{
		ServerName: TestServerName,
		//Timeout:         10 * 1000,
		LoginPath:        "/login",
		LoginBeforeFunc:  loginFunc,
		LogoutPath:       "/user/logout",
		AuthExcludePaths: g.SliceStr{"/admin/user/info", "/admin/system/user/info"}, // 不拦截路径 /user/info,/system/user/info,/system/user,
		MultiLogin:       g.Config().GetBool("gToken.MultiLogin"),
	}
	s.Group("/admin", func(group *ghttp.RouterGroup) {
		group.Middleware(CORS)
		err := gfAdminToken.Middleware(group)
		if err != nil {
			panic(err)
		}

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
	// 唯一标识，扩展参数user data
	return username, "1"
}

// 跨域
func CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

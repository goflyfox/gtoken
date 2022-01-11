package server

import (
	"context"
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
)

var TestServerName string

//var TestServerName string = "gtoken"

var server *ghttp.Server

func Start() {
	ctx := context.TODO()

	g.Log().Info(ctx, "########service start...")

	if fileConfig, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile); ok {
		fileConfig.SetPath("../config")
	}
	server = g.Server(TestServerName)
	initRouter(server)

	g.Log().Info(ctx, "########service finish.")
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
	ctx := context.TODO()

	// 不认证接口
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(CORS)

		// 调试路由
		group.ALL("/hello", func(r *ghttp.Request) {
			r.Response.WriteJson(gtoken.Succ("hello"))
		})
	})

	MultiLogin, err := g.Cfg().Get(ctx, "gToken.MultiLogin")
	if err != nil {
		panic(err)
	}
	// 认证接口
	loginFunc := Login
	// 启动gtoken
	gfToken = &gtoken.GfToken{
		ServerName:       TestServerName,
		LoginPath:        "/login",
		LoginBeforeFunc:  loginFunc,
		LogoutPath:       "/user/logout",
		AuthExcludePaths: g.SliceStr{"/user/info", "/system/user/info"}, // 不拦截路径 /user/info,/system/user/info,/system/user,
		MultiLogin:       MultiLogin.Bool(),
	}
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(CORS)
		err := gfToken.Middleware(ctx, group)
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
		MultiLogin:       MultiLogin.Bool(),
	}
	s.Group("/admin", func(group *ghttp.RouterGroup) {
		group.Middleware(CORS)
		err := gfAdminToken.Middleware(ctx, group)
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
	username := r.Get("username").String()
	passwd := r.Get("passwd").String()

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

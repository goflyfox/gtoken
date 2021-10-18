package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/zhaopengme/gtoken/gtoken"
)

var TestServerName string

//var TestServerName string = "gtoken"

func main() {
	g.Log().Info("########service start...")

	g.Cfg().SetPath("example/sample1")
	s := g.Server(TestServerName)
	initRouter(s)

	g.Log().Info("########service finish.")
	s.Run()
}

var gfToken *gtoken.GfToken

/*
统一路由注册
*/
func initRouter(s *ghttp.Server) {
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(CORS)

		// 调试路由
		group.ALL("/hello", func(r *ghttp.Request) {
			r.Response.WriteJson(gtoken.Succ("hello"))
		})
		// 获取登录扩展属性
		group.ALL("/system/data", func(r *ghttp.Request) {
			r.Response.WriteJson(gfToken.GetTokenData(r).Data)
		})
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

	loginFunc := Login
	// 启动gtoken
	gfToken = &gtoken.GfToken{
		ServerName: TestServerName,
		//Timeout:         10 * 1000,
		CacheMode:        g.Cfg().GetInt8("gToken.CacheMode"),
		CacheKey:         g.Cfg().GetString("gToken.CacheKey"),
		Timeout:          g.Cfg().GetInt("gToken.Timeout"),
		MaxRefresh:       g.Cfg().GetInt("gToken.MaxRefresh"),
		TokenDelimiter:   g.Cfg().GetString("gToken.TokenDelimiter"),
		EncryptKey:       g.Cfg().GetBytes("gToken.EncryptKey"),
		AuthFailMsg:      g.Cfg().GetString("gToken.AuthFailMsg"),
		MultiLogin:       g.Config().GetBool("gToken.MultiLogin"),
		LoginPath:        "/login",
		LoginBeforeFunc:  loginFunc,
		LogoutPath:       "/user/logout",
		AuthPaths:        g.SliceStr{"/user", "/system"},                // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
		AuthExcludePaths: g.SliceStr{"/user/info", "/system/user/info"}, // 不拦截路径 /user/info,/system/user/info,/system/user,
		GlobalMiddleware: true,                                          // 开启全局拦截
	}
	err := gfToken.Start()
	if err != nil {
		panic(err)
	}
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

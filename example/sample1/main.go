package main

import (
	"context"
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
)

var TestServerName string

//var TestServerName string = "gtoken"

func main() {
	ctx := context.TODO()
	g.Log().Info(ctx, "########service start...")

	if fileConfig, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile); ok {
		fileConfig.SetPath("example/sample1")
	}
	s := g.Server(TestServerName)
	initRouter(s)

	g.Log().Info(ctx, "########service finish.")
	s.Run()
}

var gfToken *gtoken.GfToken

/*
统一路由注册
*/
func initRouter(s *ghttp.Server) {
	ctx := context.TODO()

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
		CacheMode:        CfgGet(ctx, "gToken.CacheMode").Int8(),
		CacheKey:         CfgGet(ctx, "gToken.CacheKey").String(),
		Timeout:          CfgGet(ctx, "gToken.Timeout").Int(),
		MaxRefresh:       CfgGet(ctx, "gToken.MaxRefresh").Int(),
		TokenDelimiter:   CfgGet(ctx, "gToken.TokenDelimiter").String(),
		EncryptKey:       CfgGet(ctx, "gToken.EncryptKey").Bytes(),
		AuthFailMsg:      CfgGet(ctx, "gToken.AuthFailMsg").String(),
		MultiLogin:       CfgGet(ctx, "gToken.MultiLogin").Bool(),
		LoginPath:        "/login",
		LoginBeforeFunc:  loginFunc,
		LogoutPath:       "/user/logout",
		AuthPaths:        g.SliceStr{"/user", "/system"},                // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
		AuthExcludePaths: g.SliceStr{"/user/info", "/system/user/info"}, // 不拦截路径 /user/info,/system/user/info,/system/user,
		GlobalMiddleware: true,                                          // 开启全局拦截
	}
	err := gfToken.Start(ctx)
	if err != nil {
		panic(err)
	}
}

func CfgGet(ctx context.Context, name string) *gvar.Var {
	gVar, _ := g.Config().Get(ctx, name)
	return gVar
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

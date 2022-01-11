package server1

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
		fileConfig.SetPath("example/sample1")
	}
	server = g.Server(TestServerName)
	initRouter(server)

	g.Log().Info(ctx, "########service finish.")
	server.Start()
}

func Stop() {
	server.Shutdown()
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

	ctx := context.TODO()
	CacheModeVar, err := g.Cfg().Get(ctx, "gToken.CacheMode")
	if err != nil {
		panic(err)
	}
	CacheKey, err := g.Cfg().Get(ctx, "gToken.CacheKey")
	if err != nil {
		panic(err)
	}
	Timeout, err := g.Cfg().Get(ctx, "gToken.Timeout")
	if err != nil {
		panic(err)
	}
	MaxRefresh, err := g.Cfg().Get(ctx, "gToken.MaxRefresh")
	if err != nil {
		panic(err)
	}
	TokenDelimiter, err := g.Cfg().Get(ctx, "gToken.TokenDelimiter")
	if err != nil {
		panic(err)
	}
	EncryptKey, err := g.Cfg().Get(ctx, "gToken.EncryptKey")
	if err != nil {
		panic(err)
	}
	AuthFailMsg, err := g.Cfg().Get(ctx, "gToken.AuthFailMsg")
	if err != nil {
		panic(err)
	}
	MultiLogin, err := g.Cfg().Get(ctx, "gToken.MultiLogin")
	if err != nil {
		panic(err)
	}

	// 启动gtoken
	gfToken = &gtoken.GfToken{
		ServerName: TestServerName,
		//Timeout:         10 * 1000,
		CacheMode:        CacheModeVar.Int8(),
		CacheKey:         CacheKey.String(),
		Timeout:          Timeout.Int(),
		MaxRefresh:       MaxRefresh.Int(),
		TokenDelimiter:   TokenDelimiter.String(),
		EncryptKey:       EncryptKey.Bytes(),
		AuthFailMsg:      AuthFailMsg.String(),
		MultiLogin:       MultiLogin.Bool(),
		LoginPath:        "/login",
		LoginBeforeFunc:  loginFunc,
		LogoutPath:       "/user/logout",
		AuthPaths:        g.SliceStr{"/user", "/system"},                // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
		AuthExcludePaths: g.SliceStr{"/user/info", "/system/user/info"}, // 不拦截路径 /user/info,/system/user/info,/system/user,
		GlobalMiddleware: true,                                          // 开启全局拦截
	}
	err = gfToken.Start(ctx)
	if err != nil {
		panic(err)
	}
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

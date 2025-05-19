package backend

import (
	"context"
	"github.com/goflyfox/gtoken/gtokenv2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
)

var TestServerName string

//var TestServerName string = "gtoken"

var server *ghttp.Server

func Start() {
	ctx := gctx.New()

	g.Log().Info(ctx, "########service start...")

	if fileConfig, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile); ok {
		fileConfig.SetPath("../config")
	}
	server = g.Server(TestServerName)
	InitRouter(server)

	g.Log().Info(ctx, "########service finish.")
	server.Start()
}

func Stop() {
	server.Shutdown()
}

var gfToken gtokenv2.Token

/*
统一路由注册
*/
func InitRouter(s *ghttp.Server) {
	ctx := gctx.New()
	// 启动gtoken
	gfToken = gtokenv2.NewDefaultToken(gtokenv2.Options{
		CacheMode:      CfgGet(ctx, "gToken.CacheMode").Int8(),
		CachePreKey:    CfgGet(ctx, "gToken.CacheKey").String(),
		Timeout:        CfgGet(ctx, "gToken.Timeout").Int64(),
		MaxRefresh:     CfgGet(ctx, "gToken.MaxRefresh").Int64(),
		TokenDelimiter: CfgGet(ctx, "gToken.TokenDelimiter").String(),
		EncryptKey:     CfgGet(ctx, "gToken.EncryptKey").Bytes(),
		MultiLogin:     CfgGet(ctx, "gToken.MultiLogin").Bool(),
	})

	// 调试路由
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/hello", func(r *ghttp.Request) {
			r.Response.WriteJson(RespSuccess("hello"))
		})
	})

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(CORS)
		group.Middleware(gtokenv2.NewDefaultMiddleware(gfToken, "/user/info", "/system/user/info").Auth)
		// 获取登录扩展属性
		group.ALL("/system/data", func(r *ghttp.Request) {
			_, data, err := gfToken.Get(r.Context(), r.GetCtxVar(gtokenv2.KeyUserKey).String())
			if err != nil {
				r.Response.WriteJson(RespError(err))
			}
			r.Response.WriteJson(RespSuccess(data))
		})
		group.ALL("/system/user", func(r *ghttp.Request) {
			r.Response.WriteJson(RespSuccess("system user"))
		})
		group.ALL("/user/info", func(r *ghttp.Request) {
			r.Response.WriteJson(RespSuccess("user info"))
		})
		group.ALL("/system/user/info", func(r *ghttp.Request) {
			r.Response.WriteJson(RespSuccess("system user info"))
		})
		group.ALL("/user/logout", func(r *ghttp.Request) {
			_ = gfToken.Destroy(ctx, r.GetCtxVar(gtokenv2.KeyUserKey).String())
			r.Response.WriteJson(RespSuccess("user logout"))
		})
	})

	s.BindHandler("/login", func(r *ghttp.Request) {
		username := r.Get("username").String()
		passwd := r.Get("passwd").String()

		if username == "" || passwd == "" {
			r.Response.WriteJson(RespFail("账号或密码错误."))
			r.ExitAll()
		}
		token, err := gfToken.Generate(ctx, username, "1")
		if err != nil {
			r.Response.WriteJson(RespError(err))
			r.ExitAll()
		}
		r.Response.WriteJson(RespSuccess(g.Map{
			gtokenv2.KeyUserKey: username,
			gtokenv2.KeyToken:   token,
		}))

	})

}

func CfgGet(ctx context.Context, name string) *gvar.Var {
	gVar := g.Config().MustGet(ctx, name)
	return gVar
}

// 跨域
func CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

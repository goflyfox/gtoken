package server1

import (
	"context"
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/goflyfox/gtoken/gtokenv2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
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
	server.Start()
}

func Stop() {
	server.Shutdown()
}

var gfToken gtokenv2.Token

type Resp = ghttp.DefaultHandlerResponse

func RespError(err error) Resp {
	return Resp{Code: gerror.Code(err).Code(), Message: gerror.Code(err).Message(), Data: gerror.Code(err).Detail()}
}

func RespSuccess(data any) Resp {
	return Resp{Code: 0, Message: "success", Data: data}
}

/*
统一路由注册
*/
func initRouter(s *ghttp.Server) {
	ctx := gctx.New()
	// 启动gtoken
	gfToken = gtokenv2.NewDefaultToken(gtokenv2.Options{
		CacheMode:        CfgGet(ctx, "gToken.CacheMode").Int8(),
		CachePreKey:      CfgGet(ctx, "gToken.CacheKey").String(),
		Timeout:          CfgGet(ctx, "gToken.Timeout").Int64(),
		MaxRefresh:       CfgGet(ctx, "gToken.MaxRefresh").Int64(),
		TokenDelimiter:   CfgGet(ctx, "gToken.TokenDelimiter").String(),
		EncryptKey:       CfgGet(ctx, "gToken.EncryptKey").Bytes(),
		MultiLogin:       CfgGet(ctx, "gToken.MultiLogin").Bool(),
		AuthExcludePaths: g.SliceStr{"/user/info", "/system/user/info"}, // 不拦截路径 /user/info,/system/user/info,/system/user,
	})

	// 调试路由
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/hello", func(r *ghttp.Request) {
			r.Response.WriteJson(RespSuccess("hello"))
		})
	})

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(CORS)
		group.Middleware(gtokenv2.NewDefaultMiddleware(gfToken).Auth)
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
			r.Response.WriteJson(gtoken.Fail("账号或密码错误."))
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
	gVar, _ := g.Config().Get(ctx, name)
	return gVar
}

// 跨域
func CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

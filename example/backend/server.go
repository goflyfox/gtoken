package backend

import (
	"context"
	"github.com/goflyfox/gtoken/v2/gtoken"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

var TestServerName string

//var TestServerName string = "gtoken"

var server *ghttp.Server

func Start() {
	ctx := gctx.New()

	g.Log().Info(ctx, "########service start...")

	server = g.Server(TestServerName)
	InitRouter(server)

	g.Log().Info(ctx, "########service finish.")
	server.Start()
}

func Stop() {
	server.Shutdown()
}

var gfToken gtoken.Token

/*
统一路由注册
*/
func InitRouter(s *ghttp.Server) {
	ctx := gctx.New()
	// 创建gtoken对象
	gfToken = gtoken.NewDefaultToken(gtoken.Options{
		CacheMode:      CfgGet(ctx, "gToken.CacheMode").Int8(),
		CachePreKey:    CfgGet(ctx, "gToken.CachePreKey").String(),
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
		// 注册GfToken中间件
		group.Middleware(gtoken.NewDefaultMiddleware(gfToken, "/user/info", "/system/user/info").Auth)
		// 获取登录扩展属性
		group.ALL("/system/data", func(r *ghttp.Request) {
			// 获取登陆信息
			_, data, err := gfToken.Get(r.Context(), r.GetCtxVar(gtoken.KeyUserKey).String())
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
			// 登出销毁Token
			_ = gfToken.Destroy(ctx, r.GetCtxVar(gtoken.KeyUserKey).String())
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
		// 认证成功调用Generate生成Token
		token, err := gfToken.Generate(ctx, username, "1")
		if err != nil {
			r.Response.WriteJson(RespError(err))
			r.ExitAll()
		}
		r.Response.WriteJson(RespSuccess(g.Map{
			gtoken.KeyUserKey: username,
			gtoken.KeyToken:   token,
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

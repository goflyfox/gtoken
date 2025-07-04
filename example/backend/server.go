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

var gToken gtoken.Token

/*
统一路由注册
*/
func InitRouter(s *ghttp.Server) {
	ctx := gctx.New()
	// 创建gtoken对象
	gToken = gtoken.NewDefaultTokenByConfig()

	// 调试路由
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/hello", func(r *ghttp.Request) {
			r.Response.WriteJson(RespSuccess("hello"))
		})
	})

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(CORS)
		// 注册gToken中间件

		middlewareAuth := gtoken.NewDefaultMiddleware(gToken)
		// token校验失败后的返回方法
		middlewareAuth.ResFun = func(r *ghttp.Request, err error) {
			r.Response.WriteJson(g.Map{
				"code":    500, // 默认: gcode.CodeBusinessValidationFailed.Code()
				"message": "身份认证过期，请重新登录:" + err.Error(),
				"data":    []interface{}{},
			})
			return
		}

		group.Middleware(middlewareAuth.Auth)
		// 获取登录扩展属性
		group.ALL("/system/data", func(r *ghttp.Request) {
			// 获取登陆信息
			_, data, err := gToken.Get(r.Context(), r.GetCtxVar(gtoken.KeyUserKey).String())
			if err != nil {
				r.Response.WriteJson(RespError(err))
			}
			r.Response.WriteJson(RespSuccess(data))
		})
		// 获取登录扩展属性
		group.ALL("/system/data2", func(r *ghttp.Request) {
			// 获取登陆信息
			token, _ := gtoken.GetRequestToken(r)
			_, data, err := gToken.ParseToken(r.Context(), token)
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
			_ = gToken.Destroy(ctx, r.GetCtxVar(gtoken.KeyUserKey).String())
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
		token, err := gToken.Generate(ctx, username, g.Map{"username": username})
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
	gVar := g.Cfg().MustGet(ctx, name)
	return gVar
}

// 跨域
func CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

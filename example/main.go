package main

import (
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/text/gstr"
)

var TestServerName string

//var TestServerName string = "gtoken"

func main() {
	g.Server(TestServerName).Run()
}

// 管理初始化顺序.
func init() {
	initConfig()
	initRouter()
}

// 用于配置初始化.
func initConfig() {
	glog.Info("########service start...")

	v := g.View()
	c := g.Config()
	s := g.Server(TestServerName)

	path := ""
	// 配置对象及视图对象配置
	c.AddPath(path + "config")

	v.SetDelimiters("${", "}")
	v.AddPath(path + "template")

	// glog配置
	logPath := c.GetString("log-path")
	glog.SetPath(logPath)
	glog.SetStdoutPrint(true)

	s.SetServerRoot("./public")
	s.SetNameToUriType(ghttp.URI_TYPE_ALLLOWER)
	s.SetLogPath(logPath)
	s.SetErrorLogEnabled(true)
	s.SetAccessLogEnabled(true)
	s.SetPort(c.GetInt("http-port"))

	glog.Info("########service finish.")
}

/*
绑定业务路由
*/
func bindRouter() {

	s := g.Server(TestServerName)

	s.Group("/", func(g *ghttp.RouterGroup) {
		g.Middleware(CORS)

		// 调试路由
		g.ALL("/hello", func(r *ghttp.Request) {
			r.Response.WriteJson(gtoken.Succ("hello"))
		})
		g.ALL("/system/user", func(r *ghttp.Request) {
			r.Response.WriteJson(gtoken.Succ("system user"))
		})
		g.ALL("/user/info", func(r *ghttp.Request) {
			r.Response.WriteJson(gtoken.Succ("user info"))
		})
		g.ALL("/system/user/info", func(r *ghttp.Request) {
			r.Response.WriteJson(gtoken.Succ("system user info"))
		})
	})

	loginFunc := Login
	// 启动gtoken
	gfToken := &gtoken.GfToken{
		ServerName: TestServerName,
		//Timeout:         10 * 1000,
		CacheMode:        g.Config().GetInt8("gtoken.cache-mode"),
		LoginPath:        "/login",
		LoginBeforeFunc:  loginFunc,
		LogoutPath:       "/user/logout",
		AuthPaths:        g.SliceStr{"/user", "/system"},                // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
		AuthExcludePaths: g.SliceStr{"/user/info", "/system/user/info"}, // 不拦截路径 /user/info,/system/user/info,/system/user,
		GlobalMiddleware: true,                                          // 开启全局拦截
		MultiLogin:       g.Config().GetBool("gtoken.multi-login"),
	}
	gfToken.Start()

}

/*
统一路由注册
*/
func initRouter() {

	s := g.Server(TestServerName)

	// 绑定路由
	bindRouter()

	// 首页
	s.BindHandler("/", func(r *ghttp.Request) {
		content, err := g.View().Parse("index.html", map[string]interface{}{
			"id":    1,
			"name":  "GTOKEN",
			"title": g.Config().GetString("setting.title"),
		})
		if err != nil {
			glog.Error(err)
		}
		r.Response.Write(content)

	})

	// 某些浏览器直接请求favicon.ico文件，特别是产生404时
	s.SetRewrite("/favicon.ico", "/resource/image/favicon.ico")

}

func Login(r *ghttp.Request) (string, interface{}) {
	username := r.GetString("username")
	passwd := r.GetString("passwd")

	if username == "" || passwd == "" {
		r.Response.WriteJson(gtoken.Fail("账号或密码错误."))
		r.ExitAll()
	}

	return username, ""
}

func CORS(r *ghttp.Request) {
	opt := ghttp.CORSOptions{
		AllowOrigin:      "*",
		AllowMethods:     ghttp.HTTP_METHODS,
		AllowCredentials: "true",
		AllowHeaders:     "Origin,Content-Type,Accept,User-Agent,Cookie,Authorization,X-Auth-Token,X-Requested-With,version,time,uuid,sign",
		MaxAge:           3628800,
	}
	if origin := r.Request.Header.Get("Origin"); origin != "" {
		opt.AllowOrigin = origin
	} else if referer := r.Request.Referer(); referer != "" {
		if p := gstr.PosR(referer, "/", 6); p != -1 {
			opt.AllowOrigin = referer[:p]
		} else {
			opt.AllowOrigin = referer
		}
	}
	r.Response.CORS(opt)
	r.Middleware.Next()
}

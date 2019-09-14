package main

import (
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

func main() {
	g.Server().Run()
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
	s := g.Server()

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
	s.SetNameToUriType(ghttp.NAME_TO_URI_TYPE_ALLLOWER)
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

	s := g.Server()
	// 调试路由
	s.BindHandler("/hello", func(r *ghttp.Request) {
		r.Response.WriteJson(gtoken.Succ("hello"))
	})
	s.BindHandler("/system/user", func(r *ghttp.Request) {
		r.Response.WriteJson(gtoken.Succ("system user"))
	})

	loginFunc := Login
	// 启动gtoken
	gtoken := &gtoken.GfToken{
		//Timeout:         10 * 1000,
		CacheMode:       g.Config().GetInt8("cache-mode"),
		LoginPath:       "/login",
		LoginBeforeFunc: loginFunc,
		LogoutPath:      "/user/logout",
		AuthPaths:       g.SliceStr{"/user/*", "/system/*"},
	}
	gtoken.Start()

}

/*
统一路由注册
*/
func initRouter() {

	s := g.Server()

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
	username := r.GetPostString("username")
	passwd := r.GetPostString("passwd")

	if username == "" || passwd == "" {
		r.Response.WriteJson(gtoken.Fail("账号或密码错误."))
		r.ExitAll()
	}

	return username, ""
}

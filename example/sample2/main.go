package main

import (
	"github.com/goflyfox/gtoken/example/sample2/test/server1"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	ctx := gctx.New()

	g.Log().Info(ctx, "########service start...")

	if fileConfig, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile); ok {
		fileConfig.SetPath("../config")
	}
	server := g.Server()
	server1.InitRouter(server)

	g.Log().Info(ctx, "########service finish.")
	server.Run()
}

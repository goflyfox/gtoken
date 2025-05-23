package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"gtoken-demo/backend"
)

func main() {
	ctx := gctx.New()

	g.Log().Info(ctx, "########service start...")

	server := g.Server()
	backend.InitRouter(server)

	g.Log().Info(ctx, "########service finish.")
	server.Run()
}

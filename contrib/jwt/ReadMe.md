# gtoken-jwt

## 介绍
基于`gtoken`项目的扩展，支持jwt token方式生成token，此方式未提供Refresh接口，建议短期token场景使用；

* Github地址：https://github.com/goflyfox/gtoken/contrib/jwt
* Gitee地址：https://gitee.com/goflyfox/gtoken/contrib/jwt

## 安装教程

获取最新版本: `go get -u -v github.com/goflyfox/gtoken-jwt/v2`

## 使用说明

1. 参考`gtoken`项目使用说明

```go
	// 创建gtoken对象
    gftoken := gtoken_jwt.New(gtoken.Options{})
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(CORS)
		// 注册GfToken中间件
		group.Middleware(gtoken.NewDefaultMiddleware(gfToken).Auth)

        group.ALL("/system/data", func(r *ghttp.Request) {
            // 获取登陆信息
            _, data, err := gfToken.Get(r.Context(), r.GetCtxVar(gtoken.KeyUserKey).String())
            if err != nil {
                r.Response.WriteJson(RespError(err))
            }
            r.Response.WriteJson(RespSuccess(data))
        })
		group.ALL("/user/logout", func(r *ghttp.Request) {
		    // 登出销毁Token 
			_ = gfToken.Destroy(ctx, r.GetCtxVar(gtoken.KeyUserKey).String())
			r.Response.WriteJson(RespSuccess("user logout"))
		})
	})

	s.BindHandler("/login", func(r *ghttp.Request) {
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
```

### 配置项说明

同`gtoken`项目
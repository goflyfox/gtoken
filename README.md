# gtoken

## 介绍
基于`GoFrame`框架的token插件，通过服务端验证方式实现token认证；已完全可以支撑线上token认证，通过Redis支持集群模式；使用简单，大家可以放心使用；

* GoFrame v2.X.X 全面适配
* GoFrame v1.X.X 请使用gtoken v1.4.X相关版本;
* Github地址：https://github.com/goflyfox/gtoken
* Gitee地址：https://gitee.com/goflyfox/gtoken

## gtoken优势
1. gtoken支撑单点应用测试使用内存存储，支持个人小项目文件存储，也支持企业集群使用redis存储；完全适用于企业生产级使用；
2. 有效的避免了jwt服务端无法退出问题；
3. 解决jwt无法作废已颁布的令牌，只能等到令牌过期问题；
4. 通过用户扩展信息存储在服务端，有效规避了jwt携带大量用户扩展信息导致降低传输效率问题；
5. 有效避免jwt需要客户端实现续签功能，增加客户端复杂度；支持服务端自动续期，客户端不需要关心续签逻辑；

## 特性说明

1. 支持token认证，不强依赖于session和cookie，适用jwt和session认证所有场景；
2. 支持单机gcache和集群gredis模式；
```
# 缓存模式 1 gcache 2 gredis 3 fileCache
CacheMode = 2
```

3. 支持服务端缓存自动续期功能
```
// 注：通过MaxRefresh，默认当用户第五天访问时，自动续期
// 超时时间 默认10天
Timeout int
// 缓存刷新时间 默认为超时时间的一半
MaxRefresh int
```
4. 框架使用简单，只需要认证拦截器注册、登录Token生成、登出Token销毁即可；

## 安装教程

获取最新版本: `go get -u -v github.com/goflyfox/gtoken/v2`

## 使用说明

1. 初始化配置gtoken.Options{}, 并创建gtoken对象(`gtoken.NewDefaultToken`)；参数详情见《配置项说明》部分
2. 注册认证中间件`gtoken.NewDefaultMiddleware(gfToken).Auth`
3. 登陆认证成功后，生成Token（`gfToken.Generate`）并返回给客户端
4. 登出时销毁Token(`gfToken.Destroy`)

```go
	// 创建gtoken对象
    gftoken := gtoken.NewDefaultToken(gtoken.Options{})
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

具体可参考`GfToken`结构体，字段解释如下：

| 名称             | 配置字段       | 说明                                   |
| ---------------- | -------------- | -------------------------------------- |
| 缓存模式         | CacheMode      | 1 gcache 2 gredis 3 fileCache 默认1    |
| 缓存key          | CachePreKey    | 默认缓存前缀`GToken:`                  |
| 超时时间         | Timeout        | 默认10天（毫秒）                       |
| 缓存刷新时间     | MaxRefresh     | 默认为超时时间的一半（毫秒）           |
| Token分隔符      | TokenDelimiter | 默认`_`                                |
| Token加密key     | EncryptKey     | 默认`12345678912345678912345678912345` |
| 是否支持多端登录 | MultiLogin     | 默认false                              |
| 拦截排除地址     | excludePaths   | 此路径列表不进行认证                   |

## 示例

使用示例，请先参考`gtoken/example/sample/test/backend/server.go`文件

## 感谢

1. gf框架 [https://github.com/gogf/gf](https://github.com/gogf/gf) 
2. 历史文档v1：https://goframe.org/pages/viewpage.action?pageId=1115974

## 项目支持

- 项目的发展，离不开大家得支持~！~

- [阿里云：ECS云服务器新人优惠券；请点击这里](https://promotion.aliyun.com/ntms/yunparter/invite.html?userCode=c4hsn0gc)

- 也可以请作者喝一杯咖啡:)

![jflyfox](https://raw.githubusercontent.com/jflyfox/jfinal_cms/master/doc/pay01.jpg "Open source support")

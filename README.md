# gtoken

#### 介绍
基于gf框架的token插件，通过服务端验证方式实现token认证；已完全可以支撑线上token认证，并支持集群模式；使用简单，大家可以放心使用；

1. 支持单机gcache和集群gredis模式；
```
# 配置文件
# 缓存模式 1 gcache 2 gredis
cache-mode = 2
```

2. 支持简单token认证
3. 加入缓存自动续期功能
```
// 注：通过MaxRefresh，默认当用户第五天访问时，自动再进行五天续期
// 超时时间 默认10天
Timeout int
// 缓存刷新时间 默认为超时时间的一半
MaxRefresh int
```
4. 支持全局拦截或者深度路径拦截，便于根据个人需求定制拦截器
```
// 是否是全局认证
GlobalMiddleware bool
```

5. 框架使用简单，只需要设置登录验证方法以及登录、登出、拦截路径即可；

* github地址：https://github.com/goflyfox/gtoken
* gitee地址：https://gitee.com/goflyfox/gtoken

#### gtoken优势
1. 有效的避免了jwt服务端无法退出问题；
2. 可以解决jwt无法作废已颁布的令牌；
3. 用户扩展信息仍存储在服务端，可有效的减少传输空间；
4. gtoken支撑单点应用使用内存存储，也支持集群使用redis存储；
5. 支持缓存自动续期，并且不需要客户端进行实现；

#### 安装教程

* gopath模式: `go get github.com/goflyfox/gtoken`
* 或者 使用go.mod添加 :`require github.com/goflyfox/gtoken latest`

#### 使用说明

只需要配置登录路径、登出路径、拦截路径以及登录校验实现即可

```go
	// 启动gtoken
	gtoken := &gtoken.GfToken{
		LoginPath:       "/login",
		LoginBeforeFunc: loginFunc,
		LogoutPath:      "/user/logout",
		AuthPaths:        g.SliceStr{"/user", "/system"}, // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
		GlobalMiddleware: true,                           // 开启全局拦截，默认关闭
	}
	gtoken.Start()
```

登录方法实现

```go
func Login(r *ghttp.Request) (string, interface{}) {
	username := r.GetPostString("username")
	passwd := r.GetPostString("passwd")

	// TODO 进行登录校验

	return username, ""
}
```

#### 路径拦截规则
```go
		AuthPaths:        g.SliceStr{"/user", "/system"},             // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
		AuthExcludePaths: g.SliceStr{"/user/info", "/system/user/*"}, // 不拦截路径  /user/info,/system/user/info,/system/user,
        GlobalMiddleware: true,                           // 开启全局拦截，默认关闭
```

1. `GlobalMiddleware:true`全局拦截的是通过GF的`BindMiddleware`方法创建拦截`/*`
2. `GlobalMiddleware:false`路径拦截的是通过GF的`BindMiddleware`方法创建拦截`/user*和/system/*`
3. 按照中间件优先级路径拦截优先级很高；如果先实现部分中间件在认证前处理需要切换成全局拦截器，严格按照注册顺序即可；
4. 程序先处理认证路径，如果满足；再排除不拦截路径；
5. 如果只想用排除路径功能，将拦截路径设置为`/*`即可；

#### 逻辑测试

可运行api_test.go进行测试并查看结果；验证逻辑说明：

1. 访问用户信息，提示未携带token
2. 登录后，携带token访问正常
3. 登出成功
4. 携带之前token访问，提示未登录

```json
--- PASS: TestSystemUser (0.00s)
    api_test.go:43: 1. not login and visit user
    api_test.go:50: {"code":-1,"data":"","msg":"query token fail"}
    api_test.go:63: 2. execute login and visit user
    api_test.go:66: {"code":0,"msg":"success","data":"system user"}
    api_test.go:72: 3. execute logout
    api_test.go:75: {"code":0,"msg":"success","data":"logout success"}
    api_test.go:81: 4. visit user
    api_test.go:86: {"code":-1,"msg":"login timeout or not login","data":""}
```

#### 感谢

1. gf框架 [https://github.com/gogf/gf](https://github.com/gogf/gf) 

#### 项目支持

- 项目的发展，离不开大家得支持~！~

- [【双12】主会场 低至1折；请点击这里](https://www.aliyun.com/1212/2019/home?userCode=c4hsn0gc)
- [阿里云：ECS云服务器2折起；请点击这里](https://www.aliyun.com/acts/limit-buy?spm=5176.11544616.khv0c5cu5.1.1d8e23e8XHvEIq&userCode=c4hsn0gc)
- [阿里云：ECS云服务器新人优惠券；请点击这里](https://promotion.aliyun.com/ntms/yunparter/invite.html?userCode=c4hsn0gc)

- 也可以请作者喝一杯咖啡:)

![jflyfox](https://raw.githubusercontent.com/jflyfox/jfinal_cms/master/doc/pay01.jpg "Open source support")

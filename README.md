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
4. 框架使用简单，只需要设置登录验证方法以及登录、登出、拦截路径即可；

#### 安装教程

gopath模式: `go get https://github.comgoflyfox/gtoken`

或者 使用go.mod添加 :`require github.comgoflyfox/gtoken last`

#### 使用说明

只需要配置登录路径、登出路径、拦截路径以及登录校验实现即可

```go
	// 启动gtoken
	gtoken := &gtoken.GfToken{
		LoginPath:       "/login",
		LoginBeforeFunc: loginFunc,
		LogoutPath:      "/user/logout",
		AuthPaths:       g.SliceStr{"/user/*", "/system/*"},
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
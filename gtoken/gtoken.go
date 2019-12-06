package gtoken

import (
	"fmt"
	"github.com/gogf/gf/crypto/gaes"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/grand"
	"strings"
)

const (
	CacheModeCache = 1
	CacheModeRedis = 2
)

// GfToken gtoken结构体
type GfToken struct {
	// 缓存模式 1 gcache 2 gredis 默认1
	CacheMode int8
	// 缓存key
	CacheKey string
	// 超时时间 默认10天
	Timeout int
	// 缓存刷新时间 默认为超时时间的一半
	MaxRefresh int
	// Token分隔符
	TokenDelimiter string
	// Token加密key
	EncryptKey []byte
	// 认证失败中文提示
	AuthFailMsg string

	// 登录路径
	LoginPath string
	// 登录验证方法 return userKey 用户标识 如果userKey为空，结束执行
	LoginBeforeFunc func(r *ghttp.Request) (string, interface{})
	// 登录返回方法
	LoginAfterFunc func(r *ghttp.Request, respData Resp)
	// 登出地址
	LogoutPath string
	// 登出验证方法 return true 继续执行，否则结束执行
	LogoutBeforeFunc func(r *ghttp.Request) bool
	// 登出返回方法
	LogoutAfterFunc func(r *ghttp.Request, respData Resp)

	// 拦截地址
	AuthPaths g.SliceStr
	// 认证验证方法 return true 继续执行，否则结束执行
	AuthBeforeFunc func(r *ghttp.Request) bool
	// 认证返回方法
	AuthAfterFunc func(r *ghttp.Request, respData Resp)
}

// Init 初始化
func (m *GfToken) Init() bool {
	if m.CacheMode == 0 {
		m.CacheMode = CacheModeCache
	}

	if m.CacheKey == "" {
		m.CacheKey = "GToken:"
	}

	if m.Timeout == 0 {
		m.Timeout = 10 * 24 * 60 * 60 * 1000
	}

	if m.MaxRefresh == 0 {
		m.MaxRefresh = m.Timeout / 2
	}

	if m.TokenDelimiter == "" {
		m.TokenDelimiter = "_"
	}

	if len(m.EncryptKey) == 0 {
		m.EncryptKey = []byte("12345678912345678912345678912345")
	}

	if m.AuthFailMsg == "" {
		m.AuthFailMsg = "请求错误或登录超时"
	}

	if m.LoginAfterFunc == nil {
		m.LoginAfterFunc = func(r *ghttp.Request, respData Resp) {
			if !respData.Success() {
				r.Response.WriteJson(respData)
			} else {
				r.Response.WriteJson(Succ(g.Map{
					"token": respData.GetString("token"),
				}))
			}
		}
	}

	if m.LogoutBeforeFunc == nil {
		m.LogoutBeforeFunc = func(r *ghttp.Request) bool {
			return true
		}
	}

	if m.LogoutAfterFunc == nil {
		m.LogoutAfterFunc = func(r *ghttp.Request, respData Resp) {
			if respData.Success() {
				r.Response.WriteJson(Succ("Logout success"))
			} else {
				r.Response.WriteJson(respData)
			}
		}
	}

	if m.AuthBeforeFunc == nil {
		m.AuthBeforeFunc = func(r *ghttp.Request) bool {
			// 静态页面不拦截
			if r.IsFileRequest() {
				return false
			}

			return true
		}
	}
	if m.AuthAfterFunc == nil {
		m.AuthAfterFunc = func(r *ghttp.Request, respData Resp) {
			if respData.Success() {
				r.Middleware.Next()
			} else {
				var params map[string]interface{}
				if r.Method == "GET" {
					params = r.GetQueryMap()
				} else if r.Method == "POST" {
					params = r.GetPostMap()
				} else {
					r.Response.Writeln("Request Method is ERROR! ")
					return
				}

				no := gconv.String(gtime.Millisecond())

				glog.Info(fmt.Sprintf("[AUTH_%s][url:%s][params:%s][data:%s]",
					no, r.URL.Path, params, respData.Json()))
				respData.Msg = m.AuthFailMsg
				r.Response.WriteJson(respData)
				r.ExitAll()
			}
		}
	}

	return true
}

// Start 启动
func (m *GfToken) Start() bool {
	if !m.Init() {
		return false
	}
	glog.Info("[GToken][params:" + m.String() + "]start... ")

	s := g.Server()

	// 缓存模式
	if m.CacheMode > CacheModeRedis {
		glog.Error("[GToken]CacheMode set error")
		return false
	}

	// 认证拦截器
	if m.AuthPaths == nil {
		glog.Error("[GToken]HookPathList not set")
		return false
	}
	for _, authPath := range m.AuthPaths {
		if strings.HasSuffix(authPath, "/*") {
			s.BindMiddleware(authPath, m.AuthMiddleware)
		} else {
			s.BindMiddleware(authPath+"/*", m.AuthMiddleware)
		}
	}

	// 登录
	if m.LoginPath == "" || m.LoginBeforeFunc == nil {
		glog.Error("[GToken]LoginPath or LoginBeforeFunc not set")
		return false
	}
	s.BindHandler(m.LoginPath, m.Login)

	// 登出
	if m.LogoutPath == "" {
		glog.Error("[GToken]LogoutPath or LogoutFunc not set")
		return false
	}
	s.BindHandler(m.LogoutPath, m.Logout)

	return true
}

// Start 结束
func (m *GfToken) Stop() bool {
	glog.Info("[GToken]stop. ")
	return true
}

// GetTokenData 通过token获取对象
func (m *GfToken) GetTokenData(r *ghttp.Request) Resp {
	respData := m.getRequestToken(r)
	if respData.Success() {
		// 验证token
		respData = m.validToken(respData.DataString())
	}

	return respData
}

// Login 登录
func (m *GfToken) Login(r *ghttp.Request) {
	userKey, data := m.LoginBeforeFunc(r)
	if userKey != "" {
		// 生成token
		respToken := m.genToken(userKey, data)
		m.LoginAfterFunc(r, respToken)
	}

}

// Logout 登出
func (m *GfToken) Logout(r *ghttp.Request) {
	if m.LogoutBeforeFunc(r) {
		// 获取请求token
		respData := m.getRequestToken(r)
		if respData.Success() {
			// 删除token
			m.removeToken(respData.DataString())
		}

		m.LogoutAfterFunc(r, respData)
	}
}

// AuthMiddleware 认证拦截
func (m *GfToken) AuthMiddleware(r *ghttp.Request) {
	// 不需要认证，直接下一步
	if !m.AuthBeforeFunc(r) {
		r.Middleware.Next()
		return
	}

	// 获取请求token
	tokenResp := m.getRequestToken(r)
	if tokenResp.Success() {
		// 验证token
		tokenResp = m.validToken(tokenResp.DataString())
	}

	m.AuthAfterFunc(r, tokenResp)

}

// getRequestToken 返回请求Token
func (m *GfToken) getRequestToken(r *ghttp.Request) Resp {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			glog.Warning("[GToken]authHeader:" + authHeader + " get token key fail")
			return Unauthorized("get token key fail", "")
		} else if parts[1] == "" {
			glog.Warning("[GToken]authHeader:" + authHeader + " get token fail")
			return Unauthorized("get token fail", "")
		}

		return Succ(parts[1])
	}

	authHeader = r.GetPostString("token")
	if authHeader == "" {
		return Unauthorized("query token fail", "")
	}
	return Succ(authHeader)

}

// genToken 生成Token
func (m *GfToken) genToken(userKey string, data interface{}) Resp {
	token := m.EncryptToken(userKey)
	if !token.Success() {
		return token
	}

	cacheKey := m.CacheKey + userKey
	userCache := g.Map{
		"userKey":     userKey,
		"uuid":        token.GetString("uuid"),
		"data":        data,
		"createTime":  gtime.Now().Millisecond(),
		"refreshTime": gtime.Now().Millisecond() + gconv.Int64(m.MaxRefresh),
	}

	cacheResp := m.setCache(cacheKey, userCache)
	if !cacheResp.Success() {
		return cacheResp
	}

	return token
}

// validToken 验证Token
func (m *GfToken) validToken(token string) Resp {
	if token == "" {
		return Unauthorized("valid token empty", "")
	}

	decryptToken := m.DecryptToken(token)
	if !decryptToken.Success() {
		return decryptToken
	}

	userKey := decryptToken.GetString("userKey")
	uuid := decryptToken.GetString("uuid")
	cacheKey := m.CacheKey + userKey

	userCacheResp := m.getCache(cacheKey)
	if !userCacheResp.Success() {
		return userCacheResp
	}
	userCache := gconv.Map(userCacheResp.Data)

	if uuid != userCache["uuid"] {
		glog.Error("[GToken]user auth error, decryptToken:" + decryptToken.Json() + " cacheValue:" + gconv.String(userCache))
		return Unauthorized("user auth error", "")
	}

	nowTime := gtime.Now().Millisecond()
	refreshTime := userCache["refreshTime"]

	// 需要进行缓存超时时间刷新
	if gconv.Int64(refreshTime) == 0 || nowTime > gconv.Int64(refreshTime) {
		userCache["createTime"] = gtime.Now().Millisecond()
		userCache["refreshTime"] = gtime.Now().Millisecond() + gconv.Int64(m.MaxRefresh)
		glog.Debug("[GToken]refreshToken:" + gconv.String(userCache))
		return m.setCache(cacheKey, userCache)
	}

	return Succ(userCache)
}

// removeToken 删除Token
func (m *GfToken) removeToken(token string) Resp {
	decryptToken := m.DecryptToken(token)
	if !decryptToken.Success() {
		return decryptToken
	}

	cacheKey := m.CacheKey + decryptToken.GetString("userKey")
	return m.removeCache(cacheKey)
}

// EncryptToken token加密方法
func (m *GfToken) EncryptToken(userKey string) Resp {
	if userKey == "" {
		return Fail("encrypt userKey empty")
	}

	uuid, err := gmd5.Encrypt(grand.Str(10))
	if err != nil {
		glog.Error("[GToken]uuid error", err)
		return Error("uuid error")
	}
	tokenStr := userKey + m.TokenDelimiter + uuid

	token, err := gaes.Encrypt([]byte(tokenStr), m.EncryptKey)
	if err != nil {
		glog.Error("[GToken]encrypt error", err)
		return Error("encrypt error")
	}

	return Succ(g.Map{
		"userKey": userKey,
		"uuid":    uuid,
		"token":   gbase64.Encode(token),
	})
}

// DecryptToken token解密方法
func (m *GfToken) DecryptToken(token string) Resp {
	if token == "" {
		return Fail("decrypt token empty")
	}

	token64, err := gbase64.Decode([]byte(token))
	if err != nil {
		glog.Error("[GToken]decode error", err)
		return Error("decode error")
	}
	decryptToken, err2 := gaes.Decrypt([]byte(token64), m.EncryptKey)
	if err2 != nil {
		glog.Error("[GToken]decrypt error", err2)
		return Error("decrypt error")
	}
	tokenArray := gstr.Split(string(decryptToken), m.TokenDelimiter)
	if len(tokenArray) < 2 {
		glog.Error("[GToken]token len error")
		return Error("token len error")
	}

	return Succ(g.Map{
		"userKey": tokenArray[0],
		"uuid":    tokenArray[1],
	})
}

// String token解密方法
func (m *GfToken) String() string {
	return gconv.String(g.Map{
		// 缓存模式 1 gcache 2 gredis 默认1
		"CacheMode":      m.CacheMode,
		"CacheKey":       m.CacheKey,
		"Timeout":        m.Timeout,
		"TokenDelimiter": m.TokenDelimiter,
		"EncryptKey":     string(m.EncryptKey),
		"LoginPath":      m.LoginPath,
		"LogoutPath":     m.LogoutPath,
		"AuthPaths":      gconv.String(m.AuthPaths),
	})
}

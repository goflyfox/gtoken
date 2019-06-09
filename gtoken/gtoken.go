package gtoken

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/crypto/gaes"
	"github.com/gogf/gf/g/crypto/gmd5"
	"github.com/gogf/gf/g/encoding/gbase64"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/os/gcache"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/text/gstr"
	"github.com/gogf/gf/g/util/gconv"
	"github.com/gogf/gf/g/util/grand"
	"gtoken/utils/resp"
	"strings"
)

type GfToken struct {
	// 缓存key
	CacheKey string
	// 超时时间 默认2小时
	Timeout int
	// Token分隔符
	TokenDelimiter string
	// Token加密key
	EncryptKey []byte

	LoginPath string
	// 登录验证方法
	// return userKey 用户标识 如果userKey为空，结束执行
	LoginBeforeFunc func(r *ghttp.Request) (string, interface{})
	// 登录返回方法
	LoginAfterFunc func(r *ghttp.Request, respData resp.Resp)

	// 登出地址
	LogoutPath string
	// 登出验证方法
	// return true 继续执行，否则结束执行
	LogoutBeforeFunc func(r *ghttp.Request) bool
	// 登出返回方法
	LogoutAfterFunc func(r *ghttp.Request, respData resp.Resp)

	// 拦截地址
	AuthPaths g.SliceStr
	// 认证验证方法
	// return true 继续执行，否则结束执行
	AuthBeforeFunc func(r *ghttp.Request) bool
	// 认证返回方法
	AuthAfterFunc func(r *ghttp.Request, respData resp.Resp)
}

func (m *GfToken) Init() bool {

	if m.CacheKey == "" {
		m.CacheKey = "GToken:"
	}

	if m.Timeout == 0 {
		m.Timeout = 2 * 60 * 60 * 1000
	}

	if m.TokenDelimiter == "" {
		m.TokenDelimiter = "_"
	}

	if len(m.EncryptKey) == 0 {
		m.EncryptKey = []byte("12345678912345678912345678912345")
	}

	if m.LoginAfterFunc == nil {
		m.LoginAfterFunc = func(r *ghttp.Request, respData resp.Resp) {
			if !respData.Success() {
				r.Response.WriteJson(respData)
			} else {
				r.Response.WriteJson(resp.Succ(g.Map{
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
		m.LogoutAfterFunc = func(r *ghttp.Request, respData resp.Resp) {
			if respData.Success() {
				r.Response.WriteJson(resp.Succ("logout success"))
			} else {
				r.Response.WriteJson(respData)
			}
		}
	}

	if m.AuthBeforeFunc == nil {
		m.AuthBeforeFunc = func(r *ghttp.Request) bool {
			return true
		}
	}
	if m.AuthAfterFunc == nil {
		m.AuthAfterFunc = func(r *ghttp.Request, respData resp.Resp) {
			if !respData.Success() {
				r.Response.WriteJson(respData)
				r.ExitAll()
			}
		}
	}

	return true
}

func (m *GfToken) Start() bool {
	glog.Info("[GToken][params:" + gconv.String(m) + "]start... ")
	if !m.Init() {
		return false
	}

	s := g.Server()

	// 认证拦截器
	if m.AuthPaths == nil {
		glog.Error("[GToken]HookPathList not set")
		return false
	}
	for _, authPath := range m.AuthPaths {
		s.BindHookHandler(authPath, ghttp.HOOK_BEFORE_SERVE, m.authHook)
	}

	// 登录
	if m.LoginPath == "" || m.LoginBeforeFunc == nil {
		glog.Error("[GToken]LoginPath or LoginBeforeFunc not set")
		return false
	}
	s.BindHandler(m.LoginPath, m.login)

	// 登出
	if m.LogoutPath == "" {
		glog.Error("[GToken]LogoutPath or LogoutFunc not set")
		return false
	}
	s.BindHandler(m.LogoutPath, m.logout)

	return true
}

func (m *GfToken) Stop() bool {
	glog.Info("[GToken]stop. ")
	return true
}

// 通过token获取对象
func (m *GfToken) GetTokenData(r *ghttp.Request) resp.Resp {
	respData := m.getRequestToken(r)
	if respData.Success() {
		// 验证token
		respData = m.validToken(respData.DataString())
	}

	return respData
}

// 登录
func (m *GfToken) login(r *ghttp.Request) {
	userKey, data := m.LoginBeforeFunc(r)
	if userKey != "" {
		// 生成token
		respToken := m.genToken(userKey, data)
		m.LoginAfterFunc(r, respToken)
	}

}

// 登出
func (m *GfToken) logout(r *ghttp.Request) {
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

// 认证拦截
func (m *GfToken) authHook(r *ghttp.Request) {
	if m.AuthBeforeFunc(r) {
		// 获取请求token
		tokenResp := m.getRequestToken(r)
		if tokenResp.Success() {
			// 验证token
			tokenResp = m.validToken(tokenResp.DataString())
		}

		m.AuthAfterFunc(r, tokenResp)
	}
}

// 返回请求Token
func (m *GfToken) getRequestToken(r *ghttp.Request) resp.Resp {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			glog.Info("[GToken]authHeader:" + authHeader + "get token key fail")
			return resp.Fail("get token key fail")
		} else if parts[1] == "" {
			return resp.Fail("get token fail")
		}

		return resp.Succ(parts[1])
	}

	authHeader = r.GetPostString("token")
	if authHeader == "" {
		return resp.Fail("query token fail")
	}
	return resp.Succ(authHeader)

}

// 生成Token
func (m *GfToken) genToken(userKey string, data interface{}) resp.Resp {
	token := m.EncryptToken(userKey)
	if !token.Success() {
		return token
	}

	cacheKey := m.CacheKey + userKey
	cacheValue := g.Map{
		"userKey": userKey,
		"uuid":    token.GetString("uuid"),
		"data":    data,
	}
	gcache.Set(cacheKey, cacheValue, m.Timeout)

	return token
}

// 验证Token
func (m *GfToken) validToken(token string) resp.Resp {
	if token == "" {
		return resp.Fail("valid token empty")
	}

	decryptToken := m.DecryptToken(token)
	if !decryptToken.Success() {
		return decryptToken
	}

	userKey := decryptToken.GetString("userKey")
	uuid := decryptToken.GetString("uuid")
	cacheKey := m.CacheKey + userKey

	userCache := gcache.Get(cacheKey)

	if userCache == nil {
		return resp.Fail("login timeout or not login")
	}

	cacheValue := gconv.Map(userCache)
	if uuid != cacheValue["uuid"] {
		glog.Error("[GToken]user auth error, decryptToken:" + decryptToken.Json() + " cacheValue:" + gconv.String(cacheValue))
		return resp.Fail("user auth error")
	}

	return resp.Succ(userCache)
}

// 删除Token
func (m *GfToken) removeToken(token string) resp.Resp {
	decryptToken := m.DecryptToken(token)
	if !decryptToken.Success() {
		return decryptToken
	}

	cacheKey := m.CacheKey + decryptToken.GetString("userKey")
	gcache.Remove(cacheKey)

	return resp.Succ("")
}

func (m *GfToken) EncryptToken(userKey string) resp.Resp {
	if userKey == "" {
		return resp.Fail("encrypt userKey empty")
	}

	uuid := gmd5.Encrypt(grand.Str(10))
	tokenStr := userKey + m.TokenDelimiter + uuid

	token, err := gaes.Encrypt([]byte(tokenStr), m.EncryptKey)
	if err != nil {
		return resp.Error("encrypt error")
	}

	return resp.Succ(g.Map{
		"userKey": userKey,
		"uuid":    uuid,
		"token":   gbase64.Encode(string(token)),
	})
}

func (m *GfToken) DecryptToken(token string) resp.Resp {
	if token == "" {
		return resp.Fail("decrypt token empty")
	}

	token64, err := gbase64.Decode(token)
	if err != nil {
		return resp.Error("decode error")
	}
	decryptToken, err2 := gaes.Decrypt([]byte(token64), m.EncryptKey)
	if err2 != nil {
		return resp.Error("decrypt error")
	}
	tokenArray := gstr.Split(string(decryptToken), m.TokenDelimiter)
	if len(tokenArray) < 2 {
		return resp.Error("token len error")
	}

	return resp.Succ(g.Map{
		"userKey": tokenArray[0],
		"uuid":    tokenArray[1],
	})
}

package gtoken

import (
	"context"
	"errors"
	"fmt"
	"github.com/goflyfox/gtoken/gtokenv2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"net/http"
	"strings"
)

// GfToken gtoken结构体
type GfToken struct {
	// GoFrame server name
	ServerName string
	// 缓存模式 1 gcache 2 gredis 默认1
	CacheMode int8
	// 缓存key
	CacheKey string
	// 超时时间 默认10天（毫秒）
	Timeout int
	// 缓存刷新时间 默认为超时时间的一半（毫秒）
	MaxRefresh int
	// Token分隔符
	TokenDelimiter string
	// Token加密key
	EncryptKey []byte
	// 认证失败中文提示
	AuthFailMsg string
	// 是否支持多端登录，默认false
	MultiLogin bool
	// 是否是全局认证，兼容历史版本，已废弃
	GlobalMiddleware bool
	// 中间件类型 1 GroupMiddleware 2 BindMiddleware  3 GlobalMiddleware
	MiddlewareType uint

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
	// 拦截排除地址
	AuthExcludePaths g.SliceStr
	// 认证验证方法 return true 继续执行，否则结束执行
	AuthBeforeFunc func(r *ghttp.Request) bool
	// 认证返回方法
	AuthAfterFunc func(r *ghttp.Request, respData Resp)

	GfTokenV2 *gtokenv2.GfTokenV2
}

// Login 登录
func (m *GfToken) Login(r *ghttp.Request) {
	userKey, data := m.LoginBeforeFunc(r)
	if userKey == "" {
		g.Log().Error(r.Context(), msgLog(MsgErrUserKeyEmpty))
		return
	}

	if m.MultiLogin {
		// 支持多端重复登录，返回相同token
		token, _, err := m.GfTokenV2.Get(r.Context(), userKey)
		if err == nil && token != "" {
			m.LoginAfterFunc(r, Succ(g.Map{
				KeyUserKey: userKey,
				KeyToken:   token,
			}))
			return
		}
	}

	// 生成token
	token, err := m.GfTokenV2.Generate(r.Context(), userKey, data)
	if err != nil {
		m.LoginAfterFunc(r, Error(err.Error()))
		return
	}
	m.LoginAfterFunc(r, Succ(g.Map{
		KeyUserKey: userKey,
		KeyToken:   token,
	}))
}

// Logout 登出
func (m *GfToken) Logout(r *ghttp.Request) {
	if !m.LogoutBeforeFunc(r) {
		return
	}

	// 获取请求token
	respData := m.getRequestToken(r)
	if respData.Success() {
		// 删除token
		_ = m.GfTokenV2.DestroyByToken(r.Context(), respData.DataString())
	}

	m.LogoutAfterFunc(r, respData)
}

// AuthMiddleware 认证拦截
func (m *GfToken) authMiddleware(r *ghttp.Request) {
	urlPath := r.URL.Path
	if !m.AuthPath(r.Context(), urlPath) {
		// 如果不需要认证，继续
		r.Middleware.Next()
		return
	}

	// 不需要认证，直接下一步
	if !m.AuthBeforeFunc(r) {
		r.Middleware.Next()
		return
	}

	// 获取请求token
	tokenResp := m.getRequestToken(r)
	if tokenResp.Success() {
		token := tokenResp.DataString()
		err := m.GfTokenV2.Validate(r.Context(), token)
		if err != nil {
			// 验证token
			tokenResp = Unauthorized(err.Error(), "")
		} else {
			// 支持多端重复登录，返回相同token
			userKey, err := m.GfTokenV2.GetUserKey(r.Context(), token)
			if err != nil {
				m.LoginAfterFunc(r, Error(err.Error()))
				return
			}

			tokenResp = Succ(g.Map{
				KeyUserKey: userKey,
				KeyToken:   token,
			})
		}

	}

	m.AuthAfterFunc(r, tokenResp)
}

// GetTokenData 通过token获取对象
func (m *GfToken) GetTokenData(r *ghttp.Request) Resp {
	respData := m.getRequestToken(r)
	if respData.Success() {
		token := respData.DataString()
		data, err := m.GfTokenV2.GetData(r.Context(), token)
		if err != nil {
			return Error(err.Error())
		}
		return Succ(g.Map{
			KeyData: data,
		})
	}

	return respData
}

// AuthPath 判断路径是否需要进行认证拦截
// return true 需要认证
func (m *GfToken) AuthPath(ctx context.Context, urlPath string) bool {
	// 去除后斜杠
	if strings.HasSuffix(urlPath, "/") {
		urlPath = gstr.SubStr(urlPath, 0, len(urlPath)-1)
	}
	// 分组拦截，登录接口不拦截
	if m.MiddlewareType == MiddlewareTypeGroup {
		if (m.LoginPath != "" && gstr.HasSuffix(urlPath, m.LoginPath)) ||
			(m.LogoutPath != "" && gstr.HasSuffix(urlPath, m.LogoutPath)) {
			return false
		}
	}

	// 全局处理，认证路径拦截处理
	if m.MiddlewareType == MiddlewareTypeGlobal {
		var authFlag bool
		for _, authPath := range m.AuthPaths {
			tmpPath := authPath
			if strings.HasSuffix(tmpPath, "/*") {
				tmpPath = gstr.SubStr(tmpPath, 0, len(tmpPath)-2)
			}
			if gstr.HasPrefix(urlPath, tmpPath) {
				authFlag = true
				break
			}
		}

		if !authFlag {
			// 拦截路径不匹配
			return false
		}
	}

	// 排除路径处理，到这里nextFlag为true
	for _, excludePath := range m.AuthExcludePaths {
		tmpPath := excludePath
		// 前缀匹配
		if strings.HasSuffix(tmpPath, "/*") {
			tmpPath = gstr.SubStr(tmpPath, 0, len(tmpPath)-2)
			if gstr.HasPrefix(urlPath, tmpPath) {
				// 前缀匹配不拦截
				return false
			}
		} else {
			// 全路径匹配
			if strings.HasSuffix(tmpPath, "/") {
				tmpPath = gstr.SubStr(tmpPath, 0, len(tmpPath)-1)
			}
			if urlPath == tmpPath {
				// 全路径匹配不拦截
				return false
			}
		}
	}

	return true
}

// getRequestToken 返回请求Token
func (m *GfToken) getRequestToken(r *ghttp.Request) Resp {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			g.Log().Warning(r.Context(), msgLog(MsgErrAuthHeader, authHeader))
			return Unauthorized(fmt.Sprintf(MsgErrAuthHeader, authHeader), "")
		} else if parts[1] == "" {
			g.Log().Warning(r.Context(), msgLog(MsgErrAuthHeader, authHeader))
			return Unauthorized(fmt.Sprintf(MsgErrAuthHeader, authHeader), "")
		}

		return Succ(parts[1])
	}

	authHeader = r.Get(KeyToken).String()
	if authHeader == "" {
		return Unauthorized(MsgErrTokenEmpty, "")
	}
	return Succ(authHeader)

}

// InitConfig 初始化配置信息
func (m *GfToken) InitConfig() bool {
	if m.CacheMode == 0 {
		m.CacheMode = CacheModeCache
	}

	if m.CacheKey == "" {
		m.CacheKey = DefaultCacheKey
	}

	if m.Timeout == 0 {
		m.Timeout = DefaultTimeout
	}

	if m.MaxRefresh == 0 {
		m.MaxRefresh = m.Timeout / 2
	}

	if m.TokenDelimiter == "" {
		m.TokenDelimiter = DefaultTokenDelimiter
	}

	if len(m.EncryptKey) == 0 {
		m.EncryptKey = []byte(DefaultEncryptKey)
	}

	if m.AuthFailMsg == "" {
		m.AuthFailMsg = DefaultAuthFailMsg
	}

	// 设置中间件模式，未设置说明历史版本，通过GlobalMiddleware兼容
	if m.MiddlewareType == 0 {
		if m.GlobalMiddleware {
			m.MiddlewareType = MiddlewareTypeGlobal
		} else {
			m.MiddlewareType = MiddlewareTypeBind
		}
	}

	if m.LoginAfterFunc == nil {
		m.LoginAfterFunc = func(r *ghttp.Request, respData Resp) {
			if !respData.Success() {
				r.Response.WriteJson(respData)
			} else {
				r.Response.WriteJson(Succ(g.Map{
					KeyToken: respData.GetString(KeyToken),
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
				r.Response.WriteJson(Succ(MsgLogoutSucc))
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
				if r.Method == http.MethodGet {
					params = r.GetMap()
				} else if r.Method == http.MethodPost {
					params = r.GetMap()
				} else {
					r.Response.Writeln(MsgErrReqMethod)
					return
				}

				no := gconv.String(gtime.TimestampMilli())

				g.Log().Warning(r.Context(), fmt.Sprintf("[AUTH_%s][url:%s][params:%s][data:%s]",
					no, r.URL.Path, params, respData.Json()))
				respData.Msg = m.AuthFailMsg
				r.Response.WriteJson(respData)
				r.ExitAll()
			}
		}
	}

	m.GfTokenV2 = gtokenv2.NewGfTokenV2(gtokenv2.Options{
		CacheMode:      m.CacheMode,
		CachePreKey:    m.CacheKey,
		Timeout:        int64(m.Timeout),
		MaxRefresh:     int64(m.MaxRefresh),
		TokenDelimiter: m.TokenDelimiter,
		EncryptKey:     m.EncryptKey,
	})

	return true
}

// Start 启动
func (m *GfToken) Start() error {
	if !m.InitConfig() {
		return errors.New(MsgErrInitFail)
	}

	ctx := context.Background()
	g.Log().Info(ctx, msgLog("[params:"+m.String()+"]start... "))

	s := g.Server(m.ServerName)

	// 缓存模式
	if m.CacheMode > CacheModeFile {
		g.Log().Error(ctx, msgLog(MsgErrNotSet, "CacheMode"))
		return errors.New(fmt.Sprintf(MsgErrNotSet, "CacheMode"))
	}

	// 认证拦截器
	if m.AuthPaths == nil {
		g.Log().Error(ctx, msgLog(MsgErrNotSet, "AuthPaths"))
		return errors.New(fmt.Sprintf(MsgErrNotSet, "AuthPaths"))
	}

	// 是否是全局拦截
	if m.MiddlewareType == MiddlewareTypeGlobal {
		s.BindMiddlewareDefault(m.authMiddleware)
	} else {
		for _, authPath := range m.AuthPaths {
			tmpPath := authPath
			if !strings.HasSuffix(authPath, "/*") {
				tmpPath += "/*"
			}
			s.BindMiddleware(tmpPath, m.authMiddleware)
		}
	}

	// 登录
	if m.LoginPath == "" {
		g.Log().Error(ctx, msgLog(MsgErrNotSet, "LoginPath"))
		return errors.New(fmt.Sprintf(MsgErrNotSet, "LoginPath"))
	}
	if m.LoginBeforeFunc == nil {
		g.Log().Error(ctx, msgLog(MsgErrNotSet, "LoginBeforeFunc"))
		return errors.New(fmt.Sprintf(MsgErrNotSet, "LoginBeforeFunc"))
	}
	s.BindHandler(m.LoginPath, m.Login)

	// 登出
	if m.LogoutPath == "" {
		g.Log().Error(ctx, msgLog(MsgErrNotSet, "LogoutPath"))
		return errors.New(fmt.Sprintf(MsgErrNotSet, "LogoutPath"))
	}
	s.BindHandler(m.LogoutPath, m.Logout)

	return nil
}

// Stop 结束
func (m *GfToken) Stop(ctx context.Context) error {
	g.Log().Info(ctx, "[GToken]stop. ")
	return nil
}

// String token解密方法
func (m *GfToken) String() string {
	return gconv.String(g.Map{
		// 缓存模式 1 gcache 2 gredis 默认1
		"CacheMode":        m.CacheMode,
		"CacheKey":         m.CacheKey,
		"Timeout":          m.Timeout,
		"TokenDelimiter":   m.TokenDelimiter,
		"AuthFailMsg":      m.AuthFailMsg,
		"MultiLogin":       m.MultiLogin,
		"MiddlewareType":   m.MiddlewareType,
		"LoginPath":        m.LoginPath,
		"LogoutPath":       m.LogoutPath,
		"AuthPaths":        gconv.String(m.AuthPaths),
		"AuthExcludePaths": gconv.String(m.AuthExcludePaths),
	})
}

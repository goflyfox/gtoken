package gtokenv2

import (
	"fmt"
)

const (
	CacheModeCache   = 1
	CacheModeRedis   = 2
	CacheModeFile    = 3
	CacheModeFileDat = "gtoken.dat"

	MiddlewareTypeGroup  = 1
	MiddlewareTypeBind   = 2
	MiddlewareTypeGlobal = 3

	DefaultTimeout        = 10 * 24 * 60 * 60 * 1000
	DefaultCacheKey       = "GToken:"
	DefaultTokenDelimiter = "_"
	DefaultEncryptKey     = "12345678912345678912345678912345"
	DefaultAuthFailMsg    = "请求错误或登录超时"

	TraceId = "d5dfce77cdff812161134e55de3c5207"

	KeyUserKey    = "userKey"
	KeyCreateTime = "createTime"
	KeyData       = "data"
	KeyToken      = "token"
)

const (
	DefaultLogPrefix   = "[GToken]" // 日志前缀
	MsgLogoutSucc      = "Logout success"
	MsgErrInitFail     = "InitConfig fail"
	MsgErrNotSet       = "%s not set, error"
	MsgErrUserKeyEmpty = "userKey empty"
	MsgErrReqMethod    = "request method is error! "
	MsgErrAuthHeader   = "Authorization : %s get token key fail"
	MsgErrTokenEmpty   = "token is empty"
	MsgErrTokenEncrypt = "token encrypt error"
	MsgErrTokenDecode  = "token decode error"
	MsgErrTokenLen     = "token len error"
	MsgErrValidate     = "user validate error"
	MsgErrDataEmpty    = "cache data empty"
)

func msgLog(msg string, params ...interface{}) string {
	if len(params) == 0 {
		return DefaultLogPrefix + msg
	}
	return DefaultLogPrefix + fmt.Sprintf(msg, params...)
}

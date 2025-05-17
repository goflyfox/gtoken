package gtokenv2

const (
	CacheModeCache   = 1
	CacheModeRedis   = 2
	CacheModeFile    = 3
	CacheModeFileDat = "gtoken.dat"

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
	MsgErrUserKeyEmpty = "userKey empty"
	MsgErrTokenEmpty   = "token is empty"
	MsgErrTokenLen     = "token len error"
	MsgErrValidate     = "user validate error"
	MsgErrDataEmpty    = "cache data empty"
)

package gtoken

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
)

type Options struct {
	CacheMode        int8       // 缓存模式 1 gcache 2 gredis 3 gfile 默认1
	CachePreKey      string     // 缓存key前缀
	Timeout          int64      // 超时时间 默认10天（毫秒）
	MaxRefresh       int64      // 缓存刷新时间 默认为超时时间的一半（毫秒）
	MaxRefreshTimes  int        // 最大刷新次数 默认0 不限制
	TokenDelimiter   string     // Token分隔符
	EncryptKey       []byte     // Token加密key
	MultiLogin       bool       // 是否支持多端登录，默认false
	AuthExcludePaths g.SliceStr // 拦截排除地址
}

func (o *Options) String() string {
	return fmt.Sprintf("Options{"+
		"CacheMode:%d, CachePreKey:%s, Timeout:%d, MaxRefresh:%d"+
		", TokenDelimiter:%s, MultiLogin:%v, AuthExcludePaths:%v"+
		"}", o.CacheMode, o.CachePreKey, o.Timeout, o.MaxRefresh,
		o.TokenDelimiter, o.MultiLogin, o.AuthExcludePaths)
}

package gtoken

import (
	"fmt"
)

type Options struct {
	// 缓存模式 1 gcache 2 gredis 默认1
	CacheMode int8
	// 缓存key前缀
	CachePreKey string
	// 超时时间 默认10天（毫秒）
	Timeout int64
	// 缓存刷新时间 默认为超时时间的一半（毫秒）
	MaxRefresh int64
	// Token分隔符
	TokenDelimiter string
	// Token加密key
	EncryptKey []byte
	// 是否支持多端登录，默认false
	MultiLogin bool
}

func (o *Options) String() string {
	return fmt.Sprintf("Options{"+
		"CacheMode:%d, CachePreKey:%s, Timeout:%d"+
		", MaxRefresh:%d, TokenDelimiter:%s, MultiLogin:%v"+
		"}", o.CacheMode, o.CachePreKey, o.Timeout, o.MaxRefresh, o.TokenDelimiter, o.MultiLogin)
}

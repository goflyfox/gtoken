package gtoken

const (
	DefaultLogPrefix = "[GToken]" // 日志前缀
	MsgErrInitFail   = "InitConfig fail"
)

func logMsg(msg string) string {
	return DefaultLogPrefix + msg
}

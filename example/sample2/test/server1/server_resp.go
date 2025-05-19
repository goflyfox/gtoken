package server1

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

type Resp = ghttp.DefaultHandlerResponse

func RespError(err error) Resp {
	return Resp{Code: gerror.Code(err).Code(), Message: gerror.Code(err).Message(), Data: gerror.Code(err).Detail()}
}

func RespSuccess(data any) Resp {
	return Resp{Code: 0, Message: "success", Data: data}
}

func RespFail(msg string) Resp {
	return Resp{Code: gcode.CodeInternalError.Code(), Message: msg}
}

package gtokenv2

import (
	"context"
	"errors"
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/v2/crypto/gaes"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/grand"
)

// Encoder 定义编码器接口
type Encoder interface {
	Encode(ctx context.Context, userKey string) (token string, err error)
}

// Decoder 定义解码器接口
type Decoder interface {
	Decrypt(ctx context.Context, token string) (userKey string, err error)
}

// Codec 组合编码器和解码器接口
type Codec interface {
	Encoder
	Decoder
}

// DefaultCodec 默认编解码
type DefaultCodec struct {
	// 编码分隔符
	Delimiter string
	// 加密key
	EncryptKey []byte
}

func NewDefaultCodec(delimiter string, encryptKey []byte) *DefaultCodec {
	return &DefaultCodec{
		Delimiter:  delimiter,
		EncryptKey: encryptKey,
	}
}

// Encode token加密方法
func (c *DefaultCodec) Encode(ctx context.Context, userKey string) (token string, err error) {
	if userKey == "" {
		return "", errors.New(MsgErrUserKeyEmpty)
	}

	// 随机
	randStr, err := gmd5.Encrypt(grand.Letters(10))
	if err != nil {
		return "", err
	}

	encryptBeforeStr := userKey + c.Delimiter + randStr

	encryptByte, err := gaes.Encrypt([]byte(encryptBeforeStr), c.EncryptKey)
	if err != nil {
		return "", err
	}

	return gbase64.EncodeToString(encryptByte), nil
}

// Decrypt token解密方法
func (m *DefaultCodec) Decrypt(ctx context.Context, token string) (userKey string, err error) {
	if token == "" {
		return "", errors.New(MsgErrTokenEmpty)
	}

	token64, err := gbase64.Decode([]byte(token))
	if err != nil {
		return "", err
	}
	decryptStr, err := gaes.Decrypt(token64, m.EncryptKey)
	if err != nil {
		return "", err
	}
	decryptArray := gstr.Split(string(decryptStr), m.Delimiter)
	if len(decryptArray) < 2 {
		return "", errors.New(gtoken.MsgErrTokenLen)
	}
	return decryptArray[0], nil
}

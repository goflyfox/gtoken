package gtoken_jwt_test

import (
	"github.com/goflyfox/gtoken/v2/gtoken"
	"github.com/goflyfox/gtoken/v2/gtoken/jwt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	ctx := gctx.New()
	userKey := "testUser"
	// 非多端登陆，每次生成新Token
	{
		gToken := gtoken_jwt.New(gtoken.Options{})
		token1, err := gToken.Generate(ctx, userKey, nil)
		assert.NoError(t, err)
		token2, err := gToken.Generate(ctx, userKey, nil)
		assert.NoError(t, err)
		assert.NotEqual(t, token1, token2)
	}
}

func TestValidate(t *testing.T) {
	ctx := gctx.New()
	userKey := "testUser"
	// 登陆成功
	{
		gToken := gtoken_jwt.New(gtoken.Options{})
		token, err := gToken.Generate(ctx, userKey, nil)
		assert.NoError(t, err)
		u, err := gToken.Validate(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, userKey, u)
	}
	// Token空
	{
		gToken := gtoken_jwt.New(gtoken.Options{})
		u, err := gToken.Validate(ctx, "")
		glog.Info(ctx, u, err)
		assert.Error(t, err)
		assert.Empty(t, u)
	}
	// Token错误
	{
		gToken := gtoken_jwt.New(gtoken.Options{})
		u, err := gToken.Validate(ctx, "123")
		glog.Info(ctx, u, err)
		assert.Error(t, err)
		assert.Empty(t, u)
	}
}

func TestGet(t *testing.T) {
	ctx := gctx.New()
	userKey := "testUser"
	{
		gToken := gtoken_jwt.New(gtoken.Options{})
		data := g.Map{"a": "1"}
		token, err := gToken.Generate(ctx, userKey, data)
		assert.NoError(t, err)
		u, err := gToken.Validate(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, userKey, u)
		userKey2, data2, err := gToken.GetByToken(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, userKey, userKey2)
		assert.Equal(t, data, data2)

	}
}

func TestTimeOut(t *testing.T) {
	var (
		ctx     = gctx.New()
		userKey = "testUser"
		data    = g.Map{"a": "1"}
	)
	// token超时
	{
		gToken := gtoken_jwt.New(gtoken.Options{
			Timeout: 1000,
		})
		token, err := gToken.Generate(ctx, userKey, data)
		assert.NoError(t, err)
		_, err = gToken.Validate(ctx, token)
		assert.NoError(t, err)
		time.Sleep(2 * time.Second)
		// 超时
		userKey, err = gToken.Validate(ctx, token)
		assert.Error(t, err)

	}
}

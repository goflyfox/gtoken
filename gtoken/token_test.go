package gtoken_test

import (
	"github.com/goflyfox/gtoken/v2/gtoken"
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
		gToken := gtoken.NewDefaultToken(gtoken.Options{})
		token1, err := gToken.Generate(ctx, userKey, nil)
		assert.NoError(t, err)
		token2, err := gToken.Generate(ctx, userKey, nil)
		assert.NoError(t, err)
		assert.NotEqual(t, token1, token2)
	}
	// 支持多端登陆
	{
		gToken := gtoken.NewDefaultToken(gtoken.Options{
			MultiLogin: true,
		})
		token1, err := gToken.Generate(ctx, userKey, nil)
		assert.NoError(t, err)
		token2, err := gToken.Generate(ctx, userKey, nil)
		assert.NoError(t, err)
		assert.Equal(t, token1, token2)
	}
}

func TestValidate(t *testing.T) {
	ctx := gctx.New()
	userKey := "testUser"
	// 登陆成功
	{
		gToken := gtoken.NewDefaultToken(gtoken.Options{})
		token, err := gToken.Generate(ctx, userKey, nil)
		assert.NoError(t, err)
		u, err := gToken.Validate(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, userKey, u)
	}
	// Token空
	{
		gToken := gtoken.NewDefaultToken(gtoken.Options{})
		u, err := gToken.Validate(ctx, "")
		glog.Info(ctx, u, err)
		assert.Error(t, err)
		assert.Empty(t, u)
	}
	// Token错误
	{
		gToken := gtoken.NewDefaultToken(gtoken.Options{})
		u, err := gToken.Validate(ctx, "123")
		glog.Info(ctx, u, err)
		assert.Error(t, err)
		assert.Empty(t, u)
	}
}

func TestDestroy(t *testing.T) {
	ctx := gctx.New()
	userKey := "testUser"
	// 销毁成功
	{
		gToken := gtoken.NewDefaultToken(gtoken.Options{})
		token, err := gToken.Generate(ctx, userKey, nil)
		assert.NoError(t, err)
		u, err := gToken.Validate(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, userKey, u)

		err = gToken.Destroy(ctx, userKey)
		assert.NoError(t, err)
		u, err = gToken.Validate(ctx, token)
		glog.Info(ctx, u, err)
		assert.Error(t, err)
	}
}

func TestGet(t *testing.T) {
	ctx := gctx.New()
	userKey := "testUser"
	{
		gToken := gtoken.NewDefaultToken(gtoken.Options{})
		data := g.Map{"a": "1"}
		token, err := gToken.Generate(ctx, userKey, data)
		assert.NoError(t, err)
		u, err := gToken.Validate(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, userKey, u)
		token2, data2, err := gToken.Get(ctx, userKey)
		assert.NoError(t, err)
		assert.Equal(t, token, token2)
		assert.Equal(t, data, data2)

	}
	{
		gToken := gtoken.NewDefaultToken(gtoken.Options{})
		data := g.Map{"a": "1"}
		token, err := gToken.Generate(ctx, userKey, data)
		assert.NoError(t, err)
		u, err := gToken.Validate(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, userKey, u)
		userKey2, data2, err := gToken.ParseToken(ctx, token)
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
		gToken := gtoken.NewDefaultToken(gtoken.Options{
			Timeout:    1000,
			MaxRefresh: 500,
		})
		token, err := gToken.Generate(ctx, userKey, data)
		assert.NoError(t, err)
		_, err = gToken.Validate(ctx, token)
		assert.NoError(t, err)
		time.Sleep(2 * time.Second)
		// 超时
		_, err = gToken.Validate(ctx, token)
		assert.Error(t, err)

	}
}

func TestRefresh(t *testing.T) {
	var (
		ctx     = gctx.New()
		userKey = "testUser"
		data    = g.Map{"a": "1"}
	)
	// 超过刷新次数
	{
		gToken := gtoken.NewDefaultToken(gtoken.Options{
			Timeout:         1000,
			MaxRefresh:      500,
			MaxRefreshTimes: 1,
		})
		token, err := gToken.Generate(ctx, userKey, data)
		assert.NoError(t, err)
		_, err = gToken.Validate(ctx, token)
		assert.NoError(t, err)
		time.Sleep(600 * time.Millisecond)
		// 第一次刷新
		_, err = gToken.Validate(ctx, token)
		assert.NoError(t, err)
		time.Sleep(600 * time.Millisecond)
		// 不再刷新，但未超时
		_, err = gToken.Validate(ctx, token)
		assert.NoError(t, err)
		// 不再刷新，超过1s超时
		time.Sleep(600 * time.Millisecond)
		_, err = gToken.Validate(ctx, token)
		assert.Error(t, err)

	}
}

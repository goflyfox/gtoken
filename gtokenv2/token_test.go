package gtokenv2_test

import (
	"github.com/goflyfox/gtoken/gtokenv2"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerate(t *testing.T) {
	ctx := gctx.New()
	userKey := "testUser"
	// 非多端登陆，每次生成新Token
	{
		gfToken := gtokenv2.NewDefaultToken(gtokenv2.Options{})
		token1, err := gfToken.Generate(ctx, userKey, nil)
		assert.NoError(t, err)
		token2, err := gfToken.Generate(ctx, userKey, nil)
		assert.NoError(t, err)
		assert.NotEqual(t, token1, token2)
	}
	// 支持多端登陆
	{
		gfToken := gtokenv2.NewDefaultToken(gtokenv2.Options{
			MultiLogin: true,
		})
		token1, err := gfToken.Generate(ctx, userKey, nil)
		assert.NoError(t, err)
		token2, err := gfToken.Generate(ctx, userKey, nil)
		assert.NoError(t, err)
		assert.Equal(t, token1, token2)
	}
}

func TestValidate(t *testing.T) {
	ctx := gctx.New()
	userKey := "testUser"
	// 登陆成功
	{
		gfToken := gtokenv2.NewDefaultToken(gtokenv2.Options{})
		token, err := gfToken.Generate(ctx, userKey, nil)
		assert.NoError(t, err)
		u, err := gfToken.Validate(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, userKey, u)
	}
	// Token空
	{
		gfToken := gtokenv2.NewDefaultToken(gtokenv2.Options{})
		u, err := gfToken.Validate(ctx, "")
		glog.Info(ctx, u, err)
		assert.Error(t, err)
		assert.Empty(t, u)
	}
	// Token错误
	{
		gfToken := gtokenv2.NewDefaultToken(gtokenv2.Options{})
		u, err := gfToken.Validate(ctx, "123")
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
		gfToken := gtokenv2.NewDefaultToken(gtokenv2.Options{})
		token, err := gfToken.Generate(ctx, userKey, "1")
		assert.NoError(t, err)
		u, err := gfToken.Validate(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, userKey, u)
		gfToken.Get(ctx, userKey)

		err = gfToken.Destroy(ctx, userKey)
		assert.NoError(t, err)
		u, err = gfToken.Validate(ctx, token)
		glog.Info(ctx, u, err)
		assert.Error(t, err)
	}
}

func TestGet(t *testing.T) {
	ctx := gctx.New()
	userKey := "testUser"
	{
		gfToken := gtokenv2.NewDefaultToken(gtokenv2.Options{})
		data := "1"
		token, err := gfToken.Generate(ctx, userKey, data)
		assert.NoError(t, err)
		u, err := gfToken.Validate(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, userKey, u)
		token2, data2, err := gfToken.Get(ctx, userKey)
		assert.NoError(t, err)
		assert.Equal(t, token, token2)
		assert.Equal(t, data, data2)

	}
}

package gtokenv2_test

import (
	"github.com/goflyfox/gtoken/gtokenv2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试结构体

func TestDefaultCache(t *testing.T) {
	// 默认编解码器
	cache := gtokenv2.NewDefaultCache(gtokenv2.CacheModeFile, gtokenv2.DefaultCacheKey, gtokenv2.DefaultTimeout)
	ctx := gctx.New()
	type TestStruct struct {
		UserKey string
		Data    g.Map
	}

	tests := []struct {
		name    string
		input   TestStruct
		wantErr bool
	}{
		{
			name:    "success",
			input:   TestStruct{UserKey: "alice", Data: g.Map{"a": "1111"}},
			wantErr: false,
		},
		{
			name:    "data nil",
			input:   TestStruct{UserKey: "alice"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cache.Set(ctx, tt.input.UserKey, tt.input.Data)
			if tt.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
			}

			data, err := cache.Get(ctx, tt.input.UserKey)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, data)
			}
			assert.Equal(t, tt.input.Data, data)

			err = cache.Remove(ctx, tt.input.UserKey)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			data, err = cache.Get(ctx, tt.input.UserKey)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, len(data), 0)
		})
	}
}

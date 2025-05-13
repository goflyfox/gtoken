package gtoken_test

import (
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/v2/os/gctx"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试结构体

func Test_DefaultCodec_Encode(t *testing.T) {
	// 默认编解码器
	codec := gtoken.NewDefaultCodec("_", []byte("koi29a83idakguqjq29asd9asd8a7jhq"))
	ctx := gctx.New()
	type TestStruct struct {
		UserKey string
		Data    string
	}

	tests := []struct {
		name           string
		input          TestStruct
		wantEncodeErr  bool
		wantDecryptErr bool
	}{
		{
			name:           "success",
			input:          TestStruct{UserKey: "alice"},
			wantEncodeErr:  false,
			wantDecryptErr: false,
		},
		{
			name:           "userKey nil",
			input:          TestStruct{UserKey: ""},
			wantEncodeErr:  true,
			wantDecryptErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := codec.Encode(ctx, tt.input.UserKey)
			if tt.wantEncodeErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}

			userKey, err := codec.Decrypt(ctx, token)
			if tt.wantDecryptErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, userKey)
			}
			assert.Equal(t, tt.input.UserKey, userKey)
		})
	}
}

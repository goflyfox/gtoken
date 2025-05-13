package gtoken

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
)

// GfToken gtoken结构体
type GfToken struct {
}

// Generate 生成 Token
func (m *GfToken) Generate(ctx context.Context, userKey string, data any) (token string, err error) {
	return "123", nil
}

// Validate 验证 Token
func (m *GfToken) Validate(r *ghttp.Request) (err error) {
	return nil
}

// Get 通过userKey获取Token
func (m *GfToken) Get(ctx context.Context, userKey string) (token string, data any, err error) {
	return "234", nil, nil
}

// Refresh 刷新token的缓存有效期
func (m *GfToken) Refresh(oldToken string) (newToken string, err error) {
	return "234", nil
}

// Destroy 销毁Token
func (m *GfToken) Destroy(r *ghttp.Request) (err error) {
	return nil
}

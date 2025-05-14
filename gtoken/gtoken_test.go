package gtoken_test

import (
	"context"
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/v2/frame/g"
	"testing"
)

func TestAuthPathGlobal(t *testing.T) {
	ctx := context.Background()

	t.Log("Global auth path test ")
	// 启动gtoken
	gfToken := &gtoken.GfToken{
		//Timeout:         10 * 1000,
		AuthPaths:        g.SliceStr{"/user", "/system"},             // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
		AuthExcludePaths: g.SliceStr{"/user/info", "/system/user/*"}, // 不拦截路径  /user/info,/system/user/info,/system/user,
		MiddlewareType:   gtoken.MiddlewareTypeGlobal,                // 开启全局拦截
	}

	authPath(gfToken, t)
	flag := gfToken.AuthPath(ctx, "/test")
	if flag {
		t.Error("error:", "/test auth path error")
	}

}

func TestBindAuthPath(t *testing.T) {
	t.Log("Bind auth path test ")
	// 启动gtoken
	gfToken := &gtoken.GfToken{
		//Timeout:         10 * 1000,
		AuthPaths:        g.SliceStr{"/user", "/system"},             // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
		AuthExcludePaths: g.SliceStr{"/user/info", "/system/user/*"}, // 不拦截路径  /user/info,/system/user/info,/system/user,
		MiddlewareType:   gtoken.MiddlewareTypeBind,                  // 开启局部拦截
	}

	authPath(gfToken, t)
}

func TestGroupAuthPath(t *testing.T) {
	ctx := context.Background()

	t.Log("Group auth path test ")
	// 启动gtoken
	gfToken := &gtoken.GfToken{
		//Timeout:         10 * 1000,
		AuthExcludePaths: g.SliceStr{"/user/info", "/system/user/*"}, // 不拦截路径  /user/info,/system/user/info,/system/user,
		LoginPath:        "/login",                                   // 登录路径
		MiddlewareType:   gtoken.MiddlewareTypeGroup,                 // 开启组拦截
	}

	flag := gfToken.AuthPath(ctx, "/login")
	if flag {
		t.Error("error:", "/login auth path error")
	}

	flag = gfToken.AuthPath(ctx, "/user/info")
	if flag {
		t.Error("error:", "/user/info auth path error")
	}

	flag = gfToken.AuthPath(ctx, "/system/user/info")
	if flag {
		t.Error("error:", "/system/user/info auth path error")
	}

	flag = gfToken.AuthPath(ctx, "/system/test")
	if !flag {
		t.Error("error:", "/system/test auth path error")
	}
}

func TestAuthPathNoExclude(t *testing.T) {
	ctx := context.Background()

	t.Log("auth no exclude path test ")
	// 启动gtoken
	gfToken := &gtoken.GfToken{
		//Timeout:         10 * 1000,
		AuthPaths:      g.SliceStr{"/user", "/system"}, // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
		MiddlewareType: gtoken.MiddlewareTypeGlobal,    // 关闭全局拦截
	}

	authFlag := gfToken.AuthPath
	if authFlag(ctx, "/test") {
		t.Error(ctx, "error:", "/test auth path error")
	}
	if !authFlag(ctx, "/system/dept") {
		t.Error(ctx, "error:", "/system/dept auth path error")
	}

	if !authFlag(ctx, "/user/info") {
		t.Error(ctx, "error:", "/user/info auth path error")
	}

	if !authFlag(ctx, "/system/user") {
		t.Error(ctx, "error:", "/system/user auth path error")
	}
}

func TestAuthPathExclude(t *testing.T) {
	ctx := context.Background()

	t.Log("auth path test ")
	// 启动gtoken
	gfToken := &gtoken.GfToken{
		//Timeout:         10 * 1000,
		AuthPaths:        g.SliceStr{"/*"},                           // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
		AuthExcludePaths: g.SliceStr{"/user/info", "/system/user/*"}, // 不拦截路径  /user/info,/system/user/info,/system/user,
		MiddlewareType:   gtoken.MiddlewareTypeGlobal,                // 开启全局拦截
	}

	authFlag := gfToken.AuthPath
	if !authFlag(ctx, "/test") {
		t.Error("error:", "/test auth path error")
	}
	if !authFlag(ctx, "//system/dept") {
		t.Error("error:", "/system/dept auth path error")
	}

	if authFlag(ctx, "/user/info") {
		t.Error("error:", "/user/info auth path error")
	}

	if authFlag(ctx, "/system/user") {
		t.Error("error:", "/system/user auth path error")
	}

	if authFlag(ctx, "/system/user/info") {
		t.Error("error:", "/system/user/info auth path error")
	}

}

func authPath(gfToken *gtoken.GfToken, t *testing.T) {
	ctx := context.Background()

	flag := gfToken.AuthPath(ctx, "/user/info")
	if flag {
		t.Error("error:", "/user/info auth path error")
	}

	flag = gfToken.AuthPath(ctx, "/system/user")
	if flag {
		t.Error("error:", "/system/user auth path error")
	}

	flag = gfToken.AuthPath(ctx, "/system/user/info")
	if flag {
		t.Error("error:", "/system/user/info auth path error")
	}

	flag = gfToken.AuthPath(ctx, "/system/dept")
	if !flag {
		t.Error("error:", "/system/dept auth path error")
	}

	flag = gfToken.AuthPath(ctx, "/user/list")
	if !flag {
		t.Error("error:", "/user/list auth path error")
	}

	flag = gfToken.AuthPath(ctx, "/user/add")
	if !flag {
		t.Error("error:", "/user/add auth path error")
	}
}

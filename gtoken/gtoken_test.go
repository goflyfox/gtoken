package gtoken_test

import (
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/frame/g"
	"testing"
)

func TestAuthPathGlobal(t *testing.T) {
	t.Log("auth path test ")
	// 启动gtoken
	gfToken := &gtoken.GfToken{
		//Timeout:         10 * 1000,
		AuthPaths:        g.SliceStr{"/user", "/system"},             // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
		AuthExcludePaths: g.SliceStr{"/user/info", "/system/user/*"}, // 不拦截路径  /user/info,/system/user/info,/system/user,
		GlobalMiddleware: true,                                       // 开启全局拦截
	}

	authPath(gfToken, t)
	flag := gfToken.AuthPath("/test")
	if flag {
		t.Error("/test auth path error")
	}

}

func TestAuthPath(t *testing.T) {
	t.Log("auth path test ")
	// 启动gtoken
	gfToken := &gtoken.GfToken{
		//Timeout:         10 * 1000,
		AuthPaths:        g.SliceStr{"/user", "/system"},             // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
		AuthExcludePaths: g.SliceStr{"/user/info", "/system/user/*"}, // 不拦截路径  /user/info,/system/user/info,/system/user,
		GlobalMiddleware: false,                                      // 关闭全局拦截
	}

	authPath(gfToken, t)
}

func TestAuthPathNoExclude(t *testing.T) {
	t.Log("auth no exclude path test ")
	// 启动gtoken
	gfToken := &gtoken.GfToken{
		//Timeout:         10 * 1000,
		AuthPaths:        g.SliceStr{"/user", "/system"}, // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
		GlobalMiddleware: true,                           // 关闭全局拦截
	}

	authFlag := gfToken.AuthPath
	if authFlag("/test") {
		t.Error("/test auth path error")
	}
	if !authFlag("/system/dept") {
		t.Error("/system/dept auth path error")
	}

	if !authFlag("/user/info") {
		t.Error("/user/info auth path error")
	}

	if !authFlag("/system/user") {
		t.Error("/system/user auth path error")
	}
}

func TestAuthPathExclude(t *testing.T) {
	t.Log("auth path test ")
	// 启动gtoken
	gfToken := &gtoken.GfToken{
		//Timeout:         10 * 1000,
		AuthPaths:        g.SliceStr{"/*"},                           // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
		AuthExcludePaths: g.SliceStr{"/user/info", "/system/user/*"}, // 不拦截路径  /user/info,/system/user/info,/system/user,
		GlobalMiddleware: true,                                       // 开启全局拦截
	}

	authFlag := gfToken.AuthPath
	if !authFlag("/test") {
		t.Error("/test auth path error")
	}
	if !authFlag("//system/dept") {
		t.Error("/system/dept auth path error")
	}

	if authFlag("/user/info") {
		t.Error("/user/info auth path error")
	}

	if authFlag("/system/user") {
		t.Error("/system/user auth path error")
	}

	if authFlag("/system/user/info") {
		t.Error("/system/user/info auth path error")
	}

}

func authPath(gfToken *gtoken.GfToken, t *testing.T) {
	flag := gfToken.AuthPath("/user/info")
	if flag {
		t.Error("/user/info auth path error")
	}

	flag = gfToken.AuthPath("/system/user")
	if flag {
		t.Error("/system/user auth path error")
	}

	flag = gfToken.AuthPath("/system/user/info")
	if flag {
		t.Error("/system/user/info auth path error")
	}

	flag = gfToken.AuthPath("/system/dept")
	if !flag {
		t.Error("/system/dept auth path error")
	}

	flag = gfToken.AuthPath("/user/list")
	if !flag {
		t.Error("/user/list auth path error")
	}

	flag = gfToken.AuthPath("/user/add")
	if !flag {
		t.Error("/user/add auth path error")
	}
}

func TestEncryptDecryptToken(t *testing.T) {
	t.Log("encrypt and decrypt token test ")
	gfToken := gtoken.GfToken{}
	gfToken.Init()

	userKey := "123123"
	token := gfToken.EncryptToken(userKey, "")
	if !token.Success() {
		t.Error(token.Json())
	}
	t.Log(token.DataString())

	token2 := gfToken.DecryptToken(token.GetString("token"))
	if !token2.Success() {
		t.Error(token2.Json())
	}
	t.Log(token2.DataString())
	if userKey != token2.GetString("userKey") {
		t.Error("token decrypt userKey error")
	}
	if token.GetString("uuid") != token2.GetString("uuid") {
		t.Error("token decrypt uuid error")
	}

}

func BenchmarkEncryptDecryptToken(b *testing.B) {
	b.Log("encrypt and decrypt token test ")
	gfToken := gtoken.GfToken{}
	gfToken.Init()

	userKey := "123123"
	token := gfToken.EncryptToken(userKey, "")
	if !token.Success() {
		b.Error(token.Json())
	}
	b.Log(token.DataString())

	for i := 0; i < b.N; i++ {
		token2 := gfToken.DecryptToken(token.GetString("token"))
		if !token2.Success() {
			b.Error(token2.Json())
		}
		b.Log(token2.DataString())
		if userKey != token2.GetString("userKey") {
			b.Error("token decrypt userKey error")
		}
		if token.GetString("uuid") != token2.GetString("uuid") {
			b.Error("token decrypt uuid error")
		}
	}

}

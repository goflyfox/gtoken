package test

import (
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/tiger1103/gtoken/gtoken"
	"testing"
)

const (
	TestAdminURL string = "http://127.0.0.1:8081/admin"
)

var (
	TokenAdmin    = g.MapStrStr{}
	AdminUsername = "adminFlyFox"
)

func TestAdminSystemUser(t *testing.T) {
	// 未登录
	t.Log("1. not login and visit user")
	if r, e := g.Client().Post(TestAdminURL+"/system/user", "username="+AdminUsername); e != nil {
		t.Error("error:", e)
	} else {
		defer r.Close()

		content := string(r.ReadAll())
		t.Log(content)

		var respData gtoken.Resp
		err := json.Unmarshal([]byte(content), &respData)
		if err != nil {
			t.Error("error:", err)
		}
		if respData.Success() {
			t.Error("error:", respData.Json())
		}
	}

	// 登录，访问用户信息
	t.Log("2. execute login and visit user")
	data := PostAdmin(t, "/system/user", "username="+AdminUsername)
	if data.Success() {
		t.Log(data.Json())
	} else {
		t.Error("error:", data.Json())
	}

	// 登出
	t.Log("3. execute logout")
	data = PostAdmin(t, "/user/logout", "username="+AdminUsername)
	if data.Success() {
		t.Log(data.Json())
	} else {
		t.Error("error:", data.Json())
	}

	// 登出访问用户信息
	t.Log("4. visit user")
	data = PostAdmin(t, "/system/user", "username="+AdminUsername)
	if data.Success() {
		t.Error("error:", data.Json())
	} else {
		t.Log(data.Json())
	}
	delete(TokenAdmin, AdminUsername)
}

func TestAdminUserLoginFail(t *testing.T) {
	// 登录失败
	t.Log("1. login fail ")
	if r, e := g.Client().Post(TestAdminURL+"/login", "username=&passwd="); e != nil {
		t.Error("error:", e)
	} else {
		defer r.Close()

		content := string(r.ReadAll())

		var respData gtoken.Resp
		err := json.Unmarshal([]byte(content), &respData)
		if err != nil {
			t.Error("error:", err)
		}

		if respData.Success() {
			t.Error("error:", "login fail:"+respData.Json())
		}
	}

}

func TestAdminExclude(t *testing.T) {
	// 未登录可以访问
	t.Log("1. exclude user info")
	if r, e := g.Client().Post(TestAdminURL+"/system/user/info", "username="+AdminUsername); e != nil {
		t.Error("error:", e)
	} else {
		defer r.Close()

		content := string(r.ReadAll())
		t.Log(content)

		var respData gtoken.Resp
		err := json.Unmarshal([]byte(content), &respData)
		if err != nil {
			t.Error("error:", err)
		}
		if !respData.Success() {
			t.Error("error:", respData.Json())
		}
	}

	if r, e := g.Client().Post(TestAdminURL+"/user/info", "username="+AdminUsername); e != nil {
		t.Error("error:", e)
	} else {
		defer r.Close()

		content := string(r.ReadAll())
		t.Log(content)

		var respData gtoken.Resp
		err := json.Unmarshal([]byte(content), &respData)
		if err != nil {
			t.Error("error:", err)
		}
		if !respData.Success() {
			t.Error("error:", respData.Json())
		}
	}

}

func TestAdminLogin(t *testing.T) {
	t.Log(" login first ")
	token1 := getAdminToken(t)
	t.Log("token:" + token1)
	t.Log(" login second and same token ")
	token2 := getAdminToken(t)
	t.Log("token:" + token2)
	if token1 != token2 {
		t.Error("error:", "token not same ")
	}
	delete(TokenAdmin, AdminUsername)
}

func TestAdminLogout(t *testing.T) {
	t.Log(" logout test ")
	data := PostAdmin(t, "/user/logout", "username="+AdminUsername)
	if data.Success() {
		t.Log(data.Json())
	} else {
		t.Error("error:", data.Json())
	}
	delete(TokenAdmin, AdminUsername)
}

func TestAdminMultiLogin(t *testing.T) {
	t.Log(" TestMultiLogin start... ")
	var token1, token2 string
	if r, e := g.Client().Post(TestAdminURL+"/login", "username="+AdminUsername+"&passwd=123456"); e != nil {
		t.Error("error:", e)
	} else {
		defer r.Close()

		content := string(r.ReadAll())
		t.Log("token1 content:" + content)

		var respData gtoken.Resp
		err := json.Unmarshal([]byte(content), &respData)
		if err != nil {
			t.Error("error:", err)
		}

		if !respData.Success() {
			t.Error("error:", "resp fail:"+respData.Json())
		}

		token1 = respData.GetString("token")
	}
	t.Log("token1:" + token1)

	if r, e := g.Client().Post(TestAdminURL+"/login", "username="+AdminUsername+"&passwd=123456"); e != nil {
		t.Error("error:", e)
	} else {
		defer r.Close()

		content := string(r.ReadAll())
		t.Log("token2 content:" + content)

		var respData gtoken.Resp
		err := json.Unmarshal([]byte(content), &respData)
		if err != nil {
			t.Error("error:", err)
		}

		if !respData.Success() {
			t.Error("error:", "resp fail:"+respData.Json())
		}

		token2 = respData.GetString("token")
	}

	t.Log("token2:" + token2)

	if g.Config().GetBool("gToken.MultiLogin") {
		if token1 != token2 {
			t.Error("error:", "token not same ")
		}
	} else {
		if token1 == token2 {
			t.Error("error:", "token same ")
		}
	}
	delete(TokenAdmin, AdminUsername)
}

func PostAdmin(t *testing.T, urlPath string, data ...interface{}) gtoken.Resp {
	client := ghttp.NewClient()
	client.SetHeader("Authorization", "Bearer "+getAdminToken(t))
	content := client.RequestContent("POST", TestAdminURL+urlPath, data...)
	var respData gtoken.Resp
	err := json.Unmarshal([]byte(content), &respData)
	if err != nil {
		t.Error("error:", err)
	}
	return respData
}

func getAdminToken(t *testing.T) string {
	if TokenAdmin[AdminUsername] != "" {
		return TokenAdmin[AdminUsername]
	}

	if r, e := g.Client().Post(TestAdminURL+"/login", "username="+AdminUsername+"&passwd=123456"); e != nil {
		t.Error("error:", e)
	} else {
		defer r.Close()

		content := string(r.ReadAll())

		var respData gtoken.Resp
		err := json.Unmarshal([]byte(content), &respData)
		if err != nil {
			t.Error("error:", err)
		}

		if !respData.Success() {
			t.Error("error:", "resp fail:"+respData.Json())
		}

		TokenAdmin[AdminUsername] = respData.GetString("token")
	}
	return TokenAdmin[AdminUsername]
}

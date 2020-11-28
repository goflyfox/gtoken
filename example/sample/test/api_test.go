package test

import (
	"encoding/json"
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"testing"
)

const (
	TestURL string = "http://127.0.0.1:8081"
)

var (
	Token    = g.MapStrStr{}
	Username = "flyfox"
)

func TestHello(t *testing.T) {
	t.Log("visit hello and no auth")
	if r, e := ghttp.Post(TestURL+"/hello", "username="+Username); e != nil {
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

func TestSystemUser(t *testing.T) {
	// 未登录
	t.Log("1. not login and visit user")
	if r, e := ghttp.Post(TestURL+"/system/user", "username="+Username); e != nil {
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
	data := Post(t, "/system/user", "username="+Username)
	if data.Success() {
		t.Log(data.Json())
	} else {
		t.Error("error:", data.Json())
	}

	// 登出
	t.Log("3. execute logout")
	data = Post(t, "/user/logout", "username="+Username)
	if data.Success() {
		t.Log(data.Json())
	} else {
		t.Error("error:", data.Json())
	}

	// 登出访问用户信息
	t.Log("4. visit user")
	data = Post(t, "/system/user", "username="+Username)
	if data.Success() {
		t.Error("error:", data.Json())
	} else {
		t.Log(data.Json())
	}
}

func TestUserLoginFail(t *testing.T) {
	// 登录失败
	t.Log("1. login fail ")
	if r, e := ghttp.Post(TestURL+"/login", "username=&passwd="); e != nil {
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

func TestExclude(t *testing.T) {
	// 未登录可以访问
	t.Log("1. exclude user info")
	if r, e := ghttp.Post(TestURL+"/system/user/info", "username="+Username); e != nil {
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

	if r, e := ghttp.Post(TestURL+"/user/info", "username="+Username); e != nil {
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

//func TestRefresh(t *testing.T) {
//	// 登录，访问用户信息
//	t.Log("1. execute login and visit user")
//	data := Post(t, "/system/user", "username="+Username)
//	if data.Success() {
//		t.Log(data.Json())
//	} else {
//		t.Error("error:", data.Json())
//	}
//
//	for i := 1; i < 9; i++ {
//		time.Sleep(2 * time.Second)
//		// 登录，访问用户信息
//		t.Log("1. execute login and visit user")
//		data = Post(t, "/system/user", "username="+Username)
//		if data.Success() {
//			t.Log(data.Json())
//		} else {
//			t.Error("error:", data.Json())
//		}
//	}
//
//}

func TestLogin(t *testing.T) {
	Username = "testLogin"
	t.Log(" login first ")
	token1 := getToken(t)
	t.Log("token:" + token1)
	t.Log(" login second and same token ")
	token2 := getToken(t)
	t.Log("token:" + token2)
	if token1 != token2 {
		t.Error("error:", "token not same ")
	}
	Username = "flyfox"
}

func TestLogout(t *testing.T) {
	Username = "testLogout"
	t.Log(" logout test ")
	data := Post(t, "/user/logout", "username="+Username)
	if data.Success() {
		t.Log(data.Json())
	} else {
		t.Error("error:", data.Json())
	}
	Username = "flyfox"
}

func TestMultiLogin(t *testing.T) {
	Username = "testLogin"
	t.Log(" TestMultiLogin start... ")
	var token1, token2 string
	if r, e := ghttp.Post(TestURL+"/login", "username="+Username+"&passwd=123456"); e != nil {
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

	if r, e := ghttp.Post(TestURL+"/login", "username="+Username+"&passwd=123456"); e != nil {
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

	if g.Config().GetBool("gtoken.multi-login") {
		if token1 != token2 {
			t.Error("error:", "token not same ")
		}
	} else {
		if token1 == token2 {
			t.Error("error:", "token same ")
		}
	}

	Username = "flyfox"
}

func Post(t *testing.T, urlPath string, data ...interface{}) gtoken.Resp {
	client := ghttp.NewClient()
	client.SetHeader("Authorization", "Bearer "+getToken(t))
	content := client.RequestContent("POST", TestURL+urlPath, data...)
	var respData gtoken.Resp
	err := json.Unmarshal([]byte(content), &respData)
	if err != nil {
		t.Error("error:", err)
	}
	return respData
}

func getToken(t *testing.T) string {
	if Token[Username] != "" {
		return Token[Username]
	}

	if r, e := ghttp.Post(TestURL+"/login", "username="+Username+"&passwd=123456"); e != nil {
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

		Token[Username] = respData.GetString("token")
	}
	return Token[Username]
}

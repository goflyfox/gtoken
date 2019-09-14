package test

import (
	"encoding/json"
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"testing"
)

const (
	TestURL string = "http://127.0.0.1:80"
)

var (
	Token    = g.MapStrStr{}
	Username = "flyfox"
)

func TestHello(t *testing.T) {
	t.Log("visit hello and no auth")
	if r, e := ghttp.Post(TestURL+"/hello", "username="+Username); e != nil {
		t.Error(e)
	} else {
		defer r.Close()

		content := string(r.ReadAll())
		t.Log(content)

		var respData gtoken.Resp
		err := json.Unmarshal([]byte(content), &respData)
		if err != nil {
			t.Error(err)
		}
		if !respData.Success() {
			t.Error(respData.Json())
		}
	}
}

func TestSystemUser(t *testing.T) {
	// 未登录
	t.Log("1. not login and visit user")
	if r, e := ghttp.Post(TestURL+"/system/user", "username="+Username); e != nil {
		t.Error(e)
	} else {
		defer r.Close()

		content := string(r.ReadAll())
		t.Log(content)

		var respData gtoken.Resp
		err := json.Unmarshal([]byte(content), &respData)
		if err != nil {
			t.Error(err)
		}
		if respData.Success() {
			t.Error(respData.Json())
		}
	}

	// 登录，访问用户信息
	t.Log("2. execute login and visit user")
	data := Post(t, "/system/user", "username="+Username)
	if data.Success() {
		t.Log(data.Json())
	} else {
		t.Error(data.Json())
	}

	// 登出
	t.Log("3. execute logout")
	data = Post(t, "/user/logout", "username="+Username)
	if data.Success() {
		t.Log(data.Json())
	} else {
		t.Error(data.Json())
	}

	// 登出访问用户信息
	t.Log("4. visit user")
	data = Post(t, "/system/user", "username="+Username)
	if data.Success() {
		t.Error(data.Json())
	} else {
		t.Log(data.Json())
	}
}

//func TestRefresh(t *testing.T) {
//	// 登录，访问用户信息
//	t.Log("1. execute login and visit user")
//	data := Post(t, "/system/user", "username="+Username)
//	if data.Success() {
//		t.Log(data.Json())
//	} else {
//		t.Error(data.Json())
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
//			t.Error(data.Json())
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
		t.Error("token not same ")
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
		t.Error(data.Json())
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
		t.Error(err)
	}
	return respData
}

func getToken(t *testing.T) string {
	if Token[Username] != "" {
		return Token[Username]
	}

	if r, e := ghttp.Post(TestURL+"/login", "username="+Username+"&passwd=123456"); e != nil {
		t.Error(e)
	} else {
		defer r.Close()

		content := string(r.ReadAll())

		var respData gtoken.Resp
		err := json.Unmarshal([]byte(content), &respData)
		if err != nil {
			t.Error(err)
		}

		if !respData.Success() {
			t.Error("resp fail:" + respData.Json())
		}

		Token[Username] = respData.GetString("token")
	}
	return Token[Username]
}

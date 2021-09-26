package test

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/tiger1103/gtoken/example/sample1/test/server1"
	"github.com/tiger1103/gtoken/gtoken"
	"os"
	"testing"
)

const (
	TestURL string = "http://127.0.0.1:8082"
)

var (
	Token    = g.MapStrStr{}
	Username = "flyFox"
)

func setup() {
	fmt.Println("start...")
	server1.Start()
}

func teardown() {
	server1.Stop()
	fmt.Println("stop.")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestHello(t *testing.T) {
	t.Log("visit hello and no auth")
	if r, e := g.Client().Post(TestURL+"/hello", "username="+Username); e != nil {
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

func TestUserData(t *testing.T) {
	// 登录，访问用户信息
	t.Log("1. execute login and visit user")
	data := Post(t, "/system/data", "username="+Username)
	if data.Success() {
		if data.DataString() == "1" {
			t.Log("get user data success", data.Json())
		} else {
			t.Error("user data not eq 1 ", data.Json())
		}
	} else {
		t.Error("error:", data.Json())
	}

	// 登出
	t.Log("2. execute logout")
	data = Post(t, "/user/logout", "username="+Username)
	if data.Success() {
		t.Log(data.Json())
	} else {
		t.Error("error:", data.Json())
	}
	delete(Token, Username)
}

func TestSystemUser(t *testing.T) {
	// 未登录
	t.Log("1. not login and visit user")
	if r, e := g.Client().Post(TestURL+"/system/user", "username="+Username); e != nil {
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
	delete(Token, Username)
}

func TestUserLoginFail(t *testing.T) {
	// 登录失败
	t.Log("1. login fail ")
	if r, e := g.Client().Post(TestURL+"/login", "username=&passwd="); e != nil {
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
	if r, e := g.Client().Post(TestURL+"/system/user/info", "username="+Username); e != nil {
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

	if r, e := g.Client().Post(TestURL+"/user/info", "username="+Username); e != nil {
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
	t.Log(" login first ")
	token1 := getToken(t)
	t.Log("token:" + token1)
	t.Log(" login second and same token ")
	token2 := getToken(t)
	t.Log("token:" + token2)
	if token1 != token2 {
		t.Error("error:", "token not same ")
	}
	delete(Token, Username)
}

func TestLogout(t *testing.T) {
	t.Log(" logout test ")
	data := Post(t, "/user/logout", "username="+Username)
	if data.Success() {
		t.Log(data.Json())
	} else {
		t.Error("error:", data.Json())
	}
	delete(Token, Username)
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

	if r, e := g.Client().Post(TestURL+"/login", "username="+Username+"&passwd=123456"); e != nil {
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

func TestMultiLogin(t *testing.T) {
	t.Log(" TestMultiLogin start... ")
	var token1, token2 string
	if r, e := g.Client().Post(TestURL+"/login", "username="+Username+"&passwd=123456"); e != nil {
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

	if r, e := g.Client().Post(TestURL+"/login", "username="+Username+"&passwd=123456"); e != nil {
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
}

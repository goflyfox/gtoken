package main_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"gtoken-demo/backend"
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
	backend.Start()
}

func teardown() {
	backend.Stop()
	fmt.Println("stop.")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestHello(t *testing.T) {
	ctx := context.TODO()

	t.Log("visit hello and no auth")
	if r, e := g.Client().Post(ctx, TestURL+"/hello", "username="+Username); e != nil {
		t.Error("error:", e)
	} else {
		defer r.Close()

		content := string(r.ReadAll())
		t.Log(content)

		var respData backend.Resp
		err := json.Unmarshal([]byte(content), &respData)
		if err != nil {
			t.Error("error:", err)
		}
		if respData.Code != gcode.CodeOK.Code() {
			t.Error("error:", respData)
		}
	}
}

func TestUserData(t *testing.T) {
	// 登录，访问用户信息
	t.Log("1. execute login and visit user")
	resp := Post(t, "/system/data", "username="+Username)
	if resp.Code == gcode.CodeOK.Code() {
		if resp.Data == "1" {
			t.Log("get user data success", resp)
		} else {
			t.Error("user data not eq 1 ", resp)
		}
	} else {
		t.Error("error:", resp)
	}

	// 登出
	t.Log("2. execute logout")
	resp = Post(t, "/user/logout", "username="+Username)
	if resp.Code == gcode.CodeOK.Code() {
		t.Log(resp)
	} else {
		t.Error("error:", resp)
	}
	delete(Token, Username)
}

func TestSystemUser(t *testing.T) {
	ctx := context.TODO()
	// 未登录
	t.Log("1. not login and visit user")
	if r, e := g.Client().Post(ctx, TestURL+"/system/user", "username="+Username); e != nil {
		t.Error("error:", e)
	} else {
		defer r.Close()

		content := string(r.ReadAll())
		t.Log(content)

		var respData backend.Resp
		err := json.Unmarshal([]byte(content), &respData)
		if err != nil {
			t.Error("error:", err)
		}
		if respData.Code == gcode.CodeOK.Code() {
			t.Error("error:", respData)
		}
	}

	// 登录，访问用户信息
	t.Log("2. execute login and visit user")
	data := Post(t, "/system/user", "username="+Username)
	if data.Code == gcode.CodeOK.Code() {
		t.Log(data)
	} else {
		t.Error("error:", data)
	}

	// 登出
	t.Log("3. execute logout")
	data = Post(t, "/user/logout", "username="+Username)
	if data.Code == gcode.CodeOK.Code() {
		t.Log(data)
	} else {
		t.Error("error:", data)
	}

	// 登出访问用户信息
	t.Log("4. visit user")
	data = Post(t, "/system/user", "username="+Username)
	if data.Code == gcode.CodeOK.Code() {
		t.Error("error:", data)
	} else {
		t.Log(data)
	}
	delete(Token, Username)
}

func TestUserLoginFail(t *testing.T) {
	ctx := context.TODO()
	// 登录失败
	t.Log("1. login fail ")
	if r, e := g.Client().Post(ctx, TestURL+"/login", "username=&passwd="); e != nil {
		t.Error("error:", e)
	} else {
		defer r.Close()

		content := string(r.ReadAll())

		var respData backend.Resp
		err := json.Unmarshal([]byte(content), &respData)
		if err != nil {
			t.Error("error:", err)
		}

		if respData.Code == gcode.CodeOK.Code() {
			t.Error("error:", "login fail:", respData)
		}
	}

}

func TestExclude(t *testing.T) {
	ctx := context.TODO()
	// 未登录可以访问
	t.Log("1. exclude user info")
	if r, e := g.Client().Post(ctx, TestURL+"/system/user/info", "username="+Username); e != nil {
		t.Error("error:", e)
	} else {
		defer r.Close()

		content := string(r.ReadAll())
		t.Log(content)

		var respData backend.Resp
		err := json.Unmarshal([]byte(content), &respData)
		if err != nil {
			t.Error("error:", err)
		}
		if respData.Code != gcode.CodeOK.Code() {
			t.Error("error:", respData)
		}
	}

	if r, e := g.Client().Post(ctx, TestURL+"/user/info", "username="+Username); e != nil {
		t.Error("error:", e)
	} else {
		defer r.Close()

		content := string(r.ReadAll())
		t.Log(content)

		var respData backend.Resp
		err := json.Unmarshal([]byte(content), &respData)
		if err != nil {
			t.Error("error:", err)
		}
		if respData.Code != gcode.CodeOK.Code() {
			t.Error("error:", respData)
		}
	}

}

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

func TestMultiLogin(t *testing.T) {
	ctx := context.TODO()

	t.Log(" TestMultiLogin start... ")
	var token1, token2 string
	if r, e := g.Client().Post(ctx, TestURL+"/login", "username="+Username+"&passwd=123456"); e != nil {
		t.Error("error:", e)
	} else {
		defer r.Close()

		content := string(r.ReadAll())
		t.Log("token1 content:" + content)

		var respData backend.Resp
		err := json.Unmarshal([]byte(content), &respData)
		if err != nil {
			t.Error("error:", err)
		}

		if respData.Code != gcode.CodeOK.Code() {
			t.Error("error:", "resp fail:", respData)
		}

		token1 = gconv.String(gconv.Map(respData.Data)["token"])
	}
	t.Log("token1:" + token1)

	if r, e := g.Client().Post(ctx, TestURL+"/login", "username="+Username+"&passwd=123456"); e != nil {
		t.Error("error:", e)
	} else {
		defer r.Close()

		content := string(r.ReadAll())
		t.Log("token2 content:" + content)

		var respData backend.Resp
		err := json.Unmarshal([]byte(content), &respData)
		if err != nil {
			t.Error("error:", err)
		}

		if respData.Code != gcode.CodeOK.Code() {
			t.Error("error:", "resp fail:", respData)
		}

		token2 = gconv.String(gconv.Map(respData.Data)["token"])
	}

	t.Log("token2:" + token2)

	MultiLogin, err := g.Cfg().Get(ctx, "gToken.MultiLogin")
	if err != nil {
		panic(err)
	}
	if MultiLogin.Bool() {
		if token1 != token2 {
			t.Error("error:", "token not same ")
		}
	} else {
		if token1 == token2 {
			t.Error("error:", "token same ")
		}
	}
}

func TestLogout(t *testing.T) {
	t.Log(" logout test ")
	data := Post(t, "/user/logout", "username="+Username)
	if data.Code == gcode.CodeOK.Code() {
		t.Log(data)
	} else {
		t.Error("error:", data)
	}
	delete(Token, Username)
}

func Post(t *testing.T, urlPath string, data ...interface{}) backend.Resp {
	ctx := context.TODO()

	client := g.Client()
	client.SetHeader("Authorization", "Bearer "+getToken(t))
	content := client.RequestContent(ctx, "POST", TestURL+urlPath, data...)
	var respData backend.Resp
	err := json.Unmarshal([]byte(content), &respData)
	if err != nil {
		t.Error("error:", err)
	}
	return respData
}

func getToken(t *testing.T) string {
	ctx := context.TODO()

	if Token[Username] != "" {
		return Token[Username]
	}

	if r, e := g.Client().Post(ctx, TestURL+"/login", "username="+Username+"&passwd=123456"); e != nil {
		t.Error("error:", e)
	} else {
		defer r.Close()

		content := string(r.ReadAll())

		var respData backend.Resp
		err := json.Unmarshal([]byte(content), &respData)
		if err != nil {
			t.Error("error:", err)
		}

		if respData.Code != gcode.CodeOK.Code() {
			t.Error("error:", "resp fail:", respData)
		}

		Token[Username] = gconv.String(gconv.Map(respData.Data)["token"])
	}
	return Token[Username]
}

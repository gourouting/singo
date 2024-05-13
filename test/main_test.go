package test

import (
	"net/http"
	"singo/util"
	"testing"

	"github.com/gavv/httpexpect"
)

// 获取HttpExpect
func getHttpExpect(t *testing.T) *httpexpect.Expect {
	// 配置
	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(s),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
	return e
}

// 基础测试，保证服务可以运行
func TestQngUserLogin(t *testing.T) {
	e := getHttpExpect(t)

	obj := e.POST("/api/v1/ping").
		Expect().
		Status(http.StatusOK).JSON().Object()

	obj.Value("msg").Equal("Pong")
}

// 登录和注册
func TestUserAPI(t *testing.T) {
	e := getHttpExpect(t)

	nickName := util.RandStringRunes(8)
	userName := util.RandStringRunes(8)
	pwd := "12345678"
	data := map[string]interface{}{
		"nickname":         nickName,
		"user_name":        userName,
		"password":         pwd,
		"password_confirm": pwd,
	}

	// 注册验证
	obj := e.POST("/api/v1/user/register").
		WithHeader("Content-Type", "application/json").
		WithJSON(data).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Value("code").Equal(0)
	resData := obj.Value("data").Object()
	resData.Value("user_name").Equal(userName)
	resData.Value("nickname").Equal(nickName)

	// 验证错误的密码无法登陆
	obj = e.POST("/api/v1/user/login").
		WithHeader("Content-Type", "application/json").
		WithJSON(map[string]interface{}{
			"user_name": userName,
			"password":  "66666666",
		}).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Value("code").Equal(40001)

	// 验证正确的密码无法登陆
	obj = e.POST("/api/v1/user/login").
		WithHeader("Content-Type", "application/json").
		WithJSON(map[string]interface{}{
			"user_name": userName,
			"password":  pwd,
		}).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Value("code").Equal(0)
	resData = obj.Value("data").Object()
	resData.Value("user_name").Equal(userName)
	resData.Value("nickname").Equal(nickName)
}

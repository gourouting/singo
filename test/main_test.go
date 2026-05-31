package test

import (
	"net/http"
	"singo/util"
	"testing"

	"github.com/gavv/httpexpect"
)

// getHttpExpect returns an HttpExpect instance.
func getHttpExpect(t *testing.T) *httpexpect.Expect {
	// Configuration.
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

// Basic test to ensure the service can run.
func TestQngUserLogin(t *testing.T) {
	e := getHttpExpect(t)

	obj := e.POST("/api/v1/ping").
		Expect().
		Status(http.StatusOK).JSON().Object()

	obj.Value("msg").Equal("Pong")
}

// Login and registration.
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

	// Verify registration.
	obj := e.POST("/api/v1/user/register").
		WithHeader("Content-Type", "application/json").
		WithJSON(data).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Value("code").Equal(0)
	resData := obj.Value("data").Object()
	resData.Value("user_name").Equal(userName)
	resData.Value("nickname").Equal(nickName)

	// Verify that login fails with an incorrect password.
	obj = e.POST("/api/v1/user/login").
		WithHeader("Content-Type", "application/json").
		WithJSON(map[string]interface{}{
			"user_name": userName,
			"password":  "66666666",
		}).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Value("code").Equal(40001)

	// Verify that login succeeds with the correct password.
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

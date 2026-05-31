package api

import (
	"encoding/json"
	"fmt"
	"singo/conf"
	"singo/model"
	"singo/serializer"

	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
)

// Ping is the health check endpoint.
func Ping(c *gin.Context) {
	c.JSON(200, serializer.Response{
		Code: 0,
		Msg:  "Pong",
	})
}

// CurrentUser gets the current user.
func CurrentUser(c *gin.Context) *model.User {
	if user, _ := c.Get("user"); user != nil {
		if u, ok := user.(*model.User); ok {
			return u
		}
	}
	return nil
}

// ErrorResponse returns an error message.
func ErrorResponse(err error) serializer.Response {
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			field := conf.T(fmt.Sprintf("Field.%s", e.Field()))
			tag := conf.T(fmt.Sprintf("Tag.Valid.%s", e.Tag()))
			return serializer.ParamErr(
				fmt.Sprintf("%s%s", field, tag),
				err,
			)
		}
	}
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return serializer.ParamErr("JSON type mismatch", err)
	}

	return serializer.ParamErr("Parameter error", err)
}

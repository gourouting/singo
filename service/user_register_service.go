package service

import (
	"singo/model"
	"singo/serializer"
)

// UserRegisterService manages user registration.
type UserRegisterService struct {
	Nickname        string `form:"nickname" json:"nickname" binding:"required,min=2,max=30"`
	UserName        string `form:"user_name" json:"user_name" binding:"required,min=5,max=30"`
	Password        string `form:"password" json:"password" binding:"required,min=8,max=40"`
	PasswordConfirm string `form:"password_confirm" json:"password_confirm" binding:"required,min=8,max=40"`
}

// valid validates the form.
func (service *UserRegisterService) valid() *serializer.Response {
	if service.PasswordConfirm != service.Password {
		return &serializer.Response{
			Code: 40001,
			Msg:  "The two passwords do not match",
		}
	}

	count := int64(0)
	model.DB.Model(&model.User{}).Where("nickname = ?", service.Nickname).Count(&count)
	if count > 0 {
		return &serializer.Response{
			Code: 40001,
			Msg:  "Nickname is already taken",
		}
	}

	count = 0
	model.DB.Model(&model.User{}).Where("user_name = ?", service.UserName).Count(&count)
	if count > 0 {
		return &serializer.Response{
			Code: 40001,
			Msg:  "Username is already registered",
		}
	}

	return nil
}

// Register registers a user.
func (service *UserRegisterService) Register() serializer.Response {
	user := model.User{
		Nickname: service.Nickname,
		UserName: service.UserName,
		Status:   model.Active,
	}

	// Validate the form.
	if err := service.valid(); err != nil {
		return *err
	}

	// Encrypt the password.
	if err := user.SetPassword(service.Password); err != nil {
		return serializer.Err(
			serializer.CodeEncryptError,
			"Failed to encrypt password",
			err,
		)
	}

	// Create the user.
	if err := model.DB.Create(&user).Error; err != nil {
		return serializer.ParamErr("Registration failed", err)
	}

	return serializer.BuildUserResponse(user)
}

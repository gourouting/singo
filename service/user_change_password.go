package service

import (
	"singo/model"
	"singo/serializer"
)

// UserRegisterService 管理用户修改密码
type UserChangeService struct {
	OldPassword     string `form:"old_password" json:"old_password" binding:"required,min=8,max=40"`
	Password        string `form:"password" json:"password" binding:"required,min=8,max=40"`
	PasswordConfirm string `form:"password_confirm" json:"password_confirm" binding:"required,min=8,max=40"`
}

// Valid 验证旧密码是否正确
func (service *UserChangeService) Valid(user *model.User) *serializer.Response {
	if err := user.CheckPassword(service.OldPassword); err != true {
		return &serializer.Response{
			Status: 40001,
			Msg:    "旧密码输入有误.",
		}
	}

	return nil
}

// Change 用户修改密码
func (service *UserChangeService) Change(user *model.User) *serializer.Response {
	// 验证旧密码是否正确
	if err := service.Valid(user); err != nil {
		return err
	}

	// 将新密码进行加密
	if err := user.SetPassword(service.Password); err != nil {
		return &serializer.Response{
			Status: 50001,
			Msg:    "加密密码失败.",
		}
	}

	// 更新密码字段
	if err := model.DB.Model(&user).Update("PasswordDigest", user.PasswordDigest).Error; err != nil {
		return &serializer.Response{
			Status: 50001,
			Msg:    "更新数据库出错.",
		}
	}

	return &serializer.Response{
		Data: serializer.BuildUser(*user),
		Msg:  "修改成功~ 已清除cookie, 请重新登录",
	}

}

package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User is the user model.
type User struct {
	gorm.Model
	UserName       string
	PasswordDigest string
	Nickname       string
	Status         string
	Avatar         string `gorm:"size:1000"`
}

const (
	// PassWordCost is the password hashing cost.
	PassWordCost = 12
	// Active is an active user.
	Active string = "active"
	// Inactive is an inactive user.
	Inactive string = "inactive"
	// Suspend is a suspended user.
	Suspend string = "suspend"
)

// GetUser gets a user by ID.
func GetUser(ID interface{}) (User, error) {
	var user User
	result := DB.First(&user, ID)
	return user, result.Error
}

// SetPassword sets the password.
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// CheckPassword verifies the password.
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}

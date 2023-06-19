package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string `binding:"required"`
	Password string
	NickName string
}

func (um *User) GetUser(userName string) (u User) {
	DB.Where("user_name=?", userName).First(&u)
	return
}

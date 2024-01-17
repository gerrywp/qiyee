package models

import (
	"fmt"

	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName   string `binding:"required"`
	Password   string
	NickName   string
	IsRemember bool
}

func (um *User) GetUser(userName string) (u User) {
	r := DB.Limit(1).Where("user_name=?", userName).Find(&u)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			fmt.Println("user does not exist!")
		}
		// 数据库链接错误记录日志
		fmt.Println(r.Error)
	}
	return
}
